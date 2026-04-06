package bridgr_worker

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	workerconfig "github.com/Kanishkmittal55/bridgr-api/internal/bridgr_worker/config"
	"github.com/Kanishkmittal55/bridgr-api/internal/cloud"
	ierrr "github.com/Kanishkmittal55/bridgr-api/internal/errors"

	jobsearchv1 "github.com/Kanishkmittal55/bridgr-api/internal/gen/radar/services/job_search/v1"
	"github.com/Kanishkmittal55/bridgr-api/internal/llm/openaijson"
	"github.com/Kanishkmittal55/bridgr-api/internal/logger"
	"github.com/Kanishkmittal55/bridgr-api/internal/radar"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository/sqlc"
	"github.com/Kanishkmittal55/bridgr-api/internal/uuid"
	guuid "github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const maxJobCandidateJdTextRunes = 16000

// Processor handles job discovery work from SQS (Radar FindJobs + candidate persistence).
type Processor struct {
	repo            *repository.Repo
	q               sqlc.Querier
	jobSearchClient *radar.JobSearchClient
	radarClient     *radar.Client // PDF extract + Discovery; optional when RADAR_ADDR unset
	s3              cloud.Interface
	s3Bucket        string // HASSLE_SKIP_S3_BUCKET — JD upload for skill-gap / feed when inline jd_text
	workerOpts      workerconfig.WorkerOpts
	openAIKey       string
	prefilterLLM    *openaijson.Client
}

// NewProcessor builds a processor. jobSearch may be nil when RADAR_ADDR is unset.
func NewProcessor(repo *repository.Repo, q sqlc.Querier, jobSearch *radar.JobSearchClient, radarClient *radar.Client, s3 cloud.Interface, s3Bucket string, opts workerconfig.WorkerOpts, openAIAPIKey string) *Processor {
	return &Processor{
		repo:            repo,
		q:               q,
		jobSearchClient: jobSearch,
		radarClient:     radarClient,
		s3:              s3,
		s3Bucket:        strings.TrimSpace(s3Bucket),
		workerOpts:      opts,
		openAIKey:       strings.TrimSpace(openAIAPIKey),
		prefilterLLM:    &openaijson.Client{HTTP: &http.Client{Timeout: 90 * time.Second}},
	}
}

// Process unmarshals the queue body, runs job discovery, always Acks the message (no SQS retries),
// and returns the processing error when there was one so callers can log or metrics.
func (p *Processor) Process(ctx context.Context, msg *cloud.QueueMessage) error {
	log := logger.Get(ctx)

	var payload cloud.QueuePayload
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		log.Errorw("bridgr-worker: invalid json", "error", err)
		_ = msg.Ack()
		return fmt.Errorf("unmarshal: %w", err)
	}

	err := p.processJobDiscovery(ctx, payload)
	if err != nil {
		_ = msg.Ack() // We ack when an error occurs, no retry as of now.
		return err
	}
	return msg.Ack() // We ack when a message is successfully processed
}

func (p *Processor) processJobDiscovery(ctx context.Context, payload cloud.QueuePayload) error {
	log := logger.Get(ctx)

	runUID, err := guuid.FromString(payload.RunUUID)
	if err != nil {
		return p.finishDiscoveryRunFailed(ctx, nil, pgtype.Timestamp{}, "wrong_uuid_bad_request", fmt.Errorf("bad discovery run uuid in enqueued job: %w: %w", err, ierrr.ErrBadRequest))
	}

	pgid := uuid.ToPgUuid(runUID)
	row, err := p.repo.GetJobSearchDiscoveryRunByUUID(ctx, p.q, pgid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return p.finishDiscoveryRunFailed(ctx, nil, pgtype.Timestamp{}, "discovery_run_not_found", fmt.Errorf("discovery run not found in db: %w", ierrr.ErrNotFound))
		}
		return p.finishDiscoveryRunFailed(ctx, nil, pgtype.Timestamp{}, "discovery_run_found_but_still_error", fmt.Errorf("load discovery run from db failed: %w", ierrr.ErrNotFound))
	}

	switch row.Status {
	case "completed", "failed", "cancelled":
		return p.finishDiscoveryRunFailed(ctx, nil, pgtype.Timestamp{}, "discovery_run_found_but_status_is_wrong", fmt.Errorf("bad data returned: %w", ierrr.ErrInternal))
	}

	log.Infow("bridgr-worker: valid job discovery run received", "run_uuid", payload.RunUUID, "status", row.Status, "run_id", row.ID)

	now := pgtype.Timestamp{Time: time.Now().UTC(), Valid: true}
	// pending is not redundant here.
	if row.Status == "pending" || row.Status == "queued" {
		// TODO : where NULL means “leave unchanged”.
		_, err := p.repo.PatchJobSearchDiscoveryRun(ctx, p.q, sqlc.PatchJobSearchDiscoveryRunParams{
			ID:                row.ID,
			Status:            "running",
			StartedAt:         now,
			CompletedAt:       row.CompletedAt,
			RawCandidateCount: row.RawCandidateCount,
			NewCandidateCount: row.NewCandidateCount,
			RadarMeta:         row.RadarMeta,
			ErrorCode:         row.ErrorCode,
			ErrorDetail:       row.ErrorDetail,
			SqsMessageID:      row.SqsMessageID,
		})
		if err != nil {
			return p.finishDiscoveryRunFailed(ctx, nil, pgtype.Timestamp{}, "error_updating_discovery_run_status_running", fmt.Errorf("failed to update discovery run status to running: %w", ierrr.ErrInternal))
		}

		// Changing the values in memory for more downstream operations, this is not redundant.
		row.Status = "running"
		row.StartedAt = now
	}

	// TODO: convert this to a health check to see if radar client is alive and healthy, if not then finishDiscoveryRun sooner.
	if p.jobSearchClient == nil {
		return p.finishDiscoveryRunFailed(ctx, nil, pgtype.Timestamp{}, "radar_client_failed_health_check", fmt.Errorf("radar client failed health check before invocation: %w", ierrr.ErrInternal))
	}

	prefs, rawParams, err := ParseDiscoveryRequestParams(row.RequestParams)
	if err != nil {
		return p.finishDiscoveryRunFailed(ctx, row, now, "wrong_discovery_params_bad_request", fmt.Errorf("bad discovery params in enqueued job: %w: %w", err, ierrr.ErrBadRequest))
	}
	if prefs.UserID != 0 && prefs.UserID != row.UserID {
		return p.finishDiscoveryRunFailed(ctx, row, now, "discovery_params_user_mismatch", fmt.Errorf("request_params user_id %d does not match run user_id %d: %w", prefs.UserID, row.UserID, ierrr.ErrBadRequest))
	}

	resumeText := ""
	canonical := strings.TrimSpace(prefs.CanonicalCvAnalysisUUID)
	if canonical == "" {
		if p.workerOpts.DiscoveryRequireCanonicalCV {
			return p.finishDiscoveryRunFailed(ctx, row, now, "missing_canonical_cv_analysis_uuid", fmt.Errorf("canonical_cv_analysis_uuid required when discovery_require_canonical_cv is true: %w", ierrr.ErrBadRequest))
		}
		log.Infow("bridgr-worker: discovery run without canonical_cv_analysis_uuid; FindJobs resume_text empty", "run_uuid", payload.RunUUID, "user_id", row.UserID)
	} else {
		txt, rerr := resolveResumeTextForDiscovery(ctx, p.repo, p.q, p.radarClient, p.s3, row.UserID, canonical)
		if rerr != nil {
			if p.workerOpts.DiscoveryRequireCanonicalCV {
				return p.finishDiscoveryRunFailed(ctx, row, now, "cv_resume_unavailable", fmt.Errorf("resolve CV/resume for discovery: %w", rerr))
			}
			log.Warnw("bridgr-worker: CV resume not loaded; continuing with empty resume_text", "run_uuid", payload.RunUUID, "user_id", row.UserID, "canonical_cv_analysis_uuid", canonical, "error", rerr)
		} else {
			resumeText = txt
		}
	}

	findJobReq, sourceBoard, err := BuildFindJobsRequest(prefs, rawParams, payload.RunUUID, resumeText)
	if err != nil {
		return p.finishDiscoveryRunFailed(ctx, row, now, "wrong_discovery_params_bad_request", fmt.Errorf("build FindJobs request: %w: %w", err, ierrr.ErrBadRequest))
	}

	if (p.workerOpts.DiscoveryLLMPrefilterEnabled || p.workerOpts.DiscoveryLLMScoreEnabled || p.workerOpts.DiscoveryLLMQueryRefineEnabled) && p.openAIKey == "" {
		return p.finishDiscoveryRunFailed(ctx, row, now, "llm_discovery_misconfigured", fmt.Errorf("discovery LLM features require BRIDGR_OPENAI_API_KEY or OPENAI_API_KEY: %w", ierrr.ErrInternal))
	}

	profileMax := findJobReq.GetMaxResults()
	if profileMax <= 0 {
		profileMax = 10
	}
	targetAccept := profileMax

	maxIter := p.workerOpts.DiscoveryMaxIterations
	if maxIter < 1 {
		maxIter = 1
	}

	var deadline time.Time
	if p.workerOpts.DiscoveryIterationTimeBudgetSec > 0 {
		deadline = time.Now().UTC().Add(time.Duration(p.workerOpts.DiscoveryIterationTimeBudgetSec) * time.Second)
	}

	baseReq := copyFindJobsRequest(findJobReq)
	currentQuery := strings.TrimSpace(baseReq.GetSearchQuery())
	currentRoles := append([]string(nil), baseReq.GetTargetRoles()...)
	if len(currentRoles) == 0 && currentQuery != "" {
		currentRoles = []string{currentQuery}
	}

	var accum discoveryRunAccum
	var iterationRecords []map[string]interface{}
	var totalRaw int32
	var lastSearchQuery string
	var lastTargetRoles []string
	jobSeq := 0

	log.Infow("bridgr-worker: discovery iterative FindJobs",
		"run_uuid", payload.RunUUID,
		"location", findJobReq.GetLocation(),
		"profile_max_results", profileMax,
		"target_accept", targetAccept,
		"max_iterations", maxIter,
		"iteration_time_budget_sec", p.workerOpts.DiscoveryIterationTimeBudgetSec,
		"findjobs_page_cap", p.workerOpts.DiscoveryFindJobsPageCap,
		"source_board", sourceBoard,
		"resume_text_runes", utf8.RuneCountInString(findJobReq.GetResumeText()),
		"canonical_cv_analysis_uuid", canonical,
		"discovery_require_canonical_cv", p.workerOpts.DiscoveryRequireCanonicalCV,
		"discovery_score_threshold", p.workerOpts.DiscoveryScoreThreshold,
		"discovery_llm_query_refine_enabled", p.workerOpts.DiscoveryLLMQueryRefineEnabled,
	)

	for iter := 1; iter <= int(maxIter) && accum.ingested < targetAccept; iter++ {
		if !deadline.IsZero() && time.Now().UTC().After(deadline) {
			iterationRecords = append(iterationRecords, map[string]interface{}{
				"iteration":  iter,
				"early_exit": "time_budget",
			})
			break
		}

		need := targetAccept - accum.ingested
		pageMax := discoveryPageMaxResults(need, profileMax, p.workerOpts.DiscoveryFindJobsPageCap)

		req := copyFindJobsRequest(baseReq)
		req.SearchQuery = currentQuery
		req.TargetRoles = append([]string(nil), currentRoles...)
		if len(req.TargetRoles) == 0 && req.SearchQuery != "" {
			req.TargetRoles = []string{req.SearchQuery}
		}
		req.MaxResults = pageMax
		lastSearchQuery = req.GetSearchQuery()
		lastTargetRoles = append([]string(nil), req.GetTargetRoles()...)

		log.Infow("bridgr-worker: calling Radar JobSearchService.FindJobs",
			"run_uuid", payload.RunUUID,
			"iteration", iter,
			"iteration_max", maxIter,
			"accepted_so_far", accum.ingested,
			"target_accept", targetAccept,
			"search_query", req.GetSearchQuery(),
			"target_roles", req.GetTargetRoles(),
			"max_results_this_round", pageMax,
			"location", req.GetLocation(),
			"source_board", sourceBoard,
		)

		resp, err := p.jobSearchClient.FindJobs(ctx, req)
		if err != nil {
			return p.finishDiscoveryRunFailed(ctx, row, now, "radar_find_jobs_rpc_call_error", fmt.Errorf("radar client find jobs rpc call failed: %w: %w", err, ierrr.ErrInternal))
		}

		jobs := resp.GetJobs()
		if jobs == nil {
			jobs = []*jobsearchv1.JobInfo{}
		}

		if len(jobs) == 0 {
			if iter == 1 {
				stubMeta, _ := json.Marshal(map[string]interface{}{
					"stub": true, "worker": "bridgr", "radar_skipped": false, "did_not_find_any_jobs": true,
					"iterations": []interface{}{map[string]interface{}{"iteration": 1, "jobs_returned": 0, "search_query": req.GetSearchQuery()}},
				})
				return p.finishDiscoveryRun(ctx, row, now, 0, 0, stubMeta)
			}
			iterationRecords = append(iterationRecords, map[string]interface{}{
				"iteration":     iter,
				"jobs_returned": 0,
				"search_query":  req.GetSearchQuery(),
				"target_roles":  req.GetTargetRoles(),
				"early_exit":    "no_jobs_returned",
			})
			break
		}

		totalRaw += int32(len(jobs))

		batchHashes := uniqueURLHashesFromJobInfos(jobs)
		excluded, xerr := p.repo.BatchExcludedURLHashes(ctx, p.q, row.UserID, batchHashes)
		if xerr != nil {
			return p.finishDiscoveryRunFailed(ctx, row, now, "exclusion_batch_lookup_error", fmt.Errorf("batch exclusion lookup: %w", xerr))
		}
		if len(batchHashes) > 0 {
			log.Infow("bridgr-worker: exclusion batch lookup", "run_uuid", payload.RunUUID, "user_id", row.UserID, "batch_url_hashes", len(batchHashes), "already_excluded", len(excluded))
		}

		bIn := accum.ingested
		bEx := accum.exclusionSkipped
		bPr := accum.prefilterRejected
		bLo := accum.lowScoreRejected

		batchStats, perr := p.processDiscoveryJobsIteration(ctx, row, payload.RunUUID, prefs, canonical, resumeText, sourceBoard, pgid, now, jobs, excluded, iter, &jobSeq, &accum)
		if perr != nil {
			return p.finishDiscoveryRunFailed(ctx, row, now, "discovery_iteration_error", perr)
		}

		rec := map[string]interface{}{
			"iteration":             iter,
			"jobs_returned":         len(jobs),
			"accepted_delta":        accum.ingested - bIn,
			"exclusion_skipped":     accum.exclusionSkipped - bEx,
			"prefilter_rejected":    accum.prefilterRejected - bPr,
			"low_score_rejected":    accum.lowScoreRejected - bLo,
			"search_query":          req.GetSearchQuery(),
			"target_roles":          req.GetTargetRoles(),
			"max_results_requested": pageMax,
			"batch_stats_for_refine": map[string]interface{}{
				"accepted":           batchStats.AcceptedDelta,
				"exclusion_skipped":  batchStats.ExclusionSkipped,
				"prefilter_rejected": batchStats.PrefilterRejected,
				"low_score_rejected": batchStats.LowScoreRejected,
			},
		}

		if accum.ingested >= targetAccept {
			iterationRecords = append(iterationRecords, rec)
			break
		}

		if iter < int(maxIter) && accum.ingested < targetAccept {
			nextQ, nextR, ruleNote, llmRat := refineForNextIteration(ctx, p, prefs, currentQuery, currentRoles, batchStats, &accum.llmCallsThisRun, &accum)
			rec["refine_rule_note"] = ruleNote
			if llmRat != "" {
				rec["llm_query_refine_snippet"] = llmRat
			}
			rec["next_search_query"] = nextQ
			rec["next_target_roles"] = nextR
			currentQuery = nextQ
			currentRoles = nextR
		}

		iterationRecords = append(iterationRecords, rec)
	}

	meta, _ := json.Marshal(map[string]interface{}{
		"radar":                              true,
		"find_jobs":                          true,
		"iterative_find_jobs":                true,
		"iterations":                         iterationRecords,
		"iteration_count":                    len(iterationRecords),
		"jobs_returned_total":                totalRaw,
		"target_accept":                      targetAccept,
		"accepted_count":                     accum.ingested,
		"excluded_this_run_total":            accum.exclusionSkipped + accum.prefilterRejected + accum.lowScoreRejected,
		"exclusion_skipped":                  accum.exclusionSkipped,
		"prefilter_rejected":                 accum.prefilterRejected,
		"low_score_rejected":                 accum.lowScoreRejected,
		"llm_calls_this_run":                 accum.llmCallsThisRun,
		"llm_prefilter_calls":                accum.llmPrefilterCalls,
		"llm_score_calls":                    accum.llmScoreCalls,
		"llm_query_refine_calls":             accum.llmQueryRefineCalls,
		"discovery_llm_prefilter_enabled":    p.workerOpts.DiscoveryLLMPrefilterEnabled,
		"discovery_llm_score_enabled":        p.workerOpts.DiscoveryLLMScoreEnabled,
		"discovery_llm_query_refine_enabled": p.workerOpts.DiscoveryLLMQueryRefineEnabled,
		"candidates_upserted":                accum.ingested,
		"feed_items_upserted":                accum.feedUpserted,
		"source_board":                       sourceBoard,
		"last_search_query":                  lastSearchQuery,
		"last_target_roles":                  lastTargetRoles,
		"location":                           findJobReq.GetLocation(),
		"resume_text_runes":                  utf8.RuneCountInString(findJobReq.GetResumeText()),
		"canonical_cv_analysis_uuid":         canonical,
		"discovery_pipeline": map[string]interface{}{
			"require_canonical_cv":      p.workerOpts.DiscoveryRequireCanonicalCV,
			"score_threshold":           p.workerOpts.DiscoveryScoreThreshold,
			"prefilter_min_score":       p.workerOpts.DiscoveryPrefilterMinScore,
			"max_iterations":            p.workerOpts.DiscoveryMaxIterations,
			"iteration_time_budget_sec": p.workerOpts.DiscoveryIterationTimeBudgetSec,
			"findjobs_page_cap":         p.workerOpts.DiscoveryFindJobsPageCap,
			"max_llm_calls_per_run":     p.workerOpts.DiscoveryMaxLLMCallsPerRun,
			"llm_prefilter_enabled":     p.workerOpts.DiscoveryLLMPrefilterEnabled,
			"llm_score_enabled":         p.workerOpts.DiscoveryLLMScoreEnabled,
			"llm_query_refine_enabled":  p.workerOpts.DiscoveryLLMQueryRefineEnabled,
			"openai_model":              p.workerOpts.DiscoveryOpenAIModel,
			"scoring_version_static":    discoveryScoringVersion,
		},
	})

	return p.finishDiscoveryRun(ctx, row, now, totalRaw, accum.ingested, meta)
}

type discoveryRunAccum struct {
	ingested            int32
	feedUpserted        int32
	exclusionSkipped    int32
	prefilterRejected   int32
	lowScoreRejected    int32
	llmCallsThisRun     int32
	llmPrefilterCalls   int32
	llmScoreCalls       int32
	llmQueryRefineCalls int32
}

func (p *Processor) processDiscoveryJobsIteration(
	ctx context.Context,
	row *sqlc.BridgrJobSearchDiscoveryRun,
	runUUID string,
	prefs DiscoveryRequestParams,
	canonical string,
	resumeText string,
	sourceBoard string,
	pgid pgtype.UUID,
	now pgtype.Timestamp,
	jobs []*jobsearchv1.JobInfo,
	excluded map[string]struct{},
	iteration int,
	jobSeq *int,
	accum *discoveryRunAccum,
) (discoveryBatchStats, error) {
	log := logger.Get(ctx)
	bIn := accum.ingested
	bEx := accum.exclusionSkipped
	bPr := accum.prefilterRejected
	bLo := accum.lowScoreRejected

	for _, j := range jobs {
		jobURL := strings.TrimSpace(j.GetJobUrl())
		if jobURL == "" {
			continue
		}

		urlHash := urlHashHex(jobURL)
		if _, isExcluded := excluded[urlHash]; isExcluded {
			accum.exclusionSkipped++
			log.Infow("bridgr-worker: skip job (URL on user discovery exclusion list)", "run_uuid", runUUID, "user_id", row.UserID, "url_hash", urlHash)
			continue
		}

		reqs := j.GetRequirements()
		if reqs == nil {
			reqs = []string{}
		}

		jdInline := jdTextFromJobInfo(j)

		if p.workerOpts.DiscoveryLLMPrefilterEnabled && p.openAIKey != "" {
			budget := p.workerOpts.DiscoveryMaxLLMCallsPerRun
			if budget > 0 && accum.llmCallsThisRun >= budget {
				log.Warnw("bridgr-worker: LLM prefilter budget exhausted; skipping prefilter for this job",
					"run_uuid", runUUID, "user_id", row.UserID, "url_hash", urlHash,
					"budget", budget)
			} else {
				accum.llmCallsThisRun++
				accum.llmPrefilterCalls++
				userPayload, perr := buildPrefilterUserPayload(prefs, j.GetTitle(), j.GetCompany(), jdInline)
				if perr != nil {
					return discoveryBatchStats{}, fmt.Errorf("prefilter payload: %w", perr)
				}
				res, rerr := runJobPrefilterLLM(ctx, p.prefilterLLM, p.openAIKey, p.workerOpts.DiscoveryOpenAIModel, p.workerOpts.DiscoveryOpenAIBaseURL, userPayload)
				if rerr != nil {
					return discoveryBatchStats{}, fmt.Errorf("prefilter llm: %w", rerr)
				}
				if !prefilterConfidenceAccept(res, p.workerOpts.DiscoveryPrefilterMinScore) {
					accum.prefilterRejected++
					if exErr := p.repo.AddJobDiscoveryExclusion(ctx, p.q, row.UserID, urlHash, "prefilter", pgid); exErr != nil {
						log.Errorw("bridgr-worker: add prefilter exclusion failed", "url_hash", urlHash, "error", exErr)
					}
					log.Infow("bridgr-worker: prefilter rejected job", "run_uuid", runUUID, "user_id", row.UserID, "url_hash", urlHash,
						"pass", res.Pass, "confidence", res.Confidence, "by_dimension", res.ByDimension)
					continue
				}
			}
		}

		scoreOut := computeDiscoveryScore(prefs, resumeText, jdInline, j.GetTitle(), j.GetCompany())
		if p.workerOpts.DiscoveryLLMScoreEnabled && p.openAIKey != "" {
			budget := p.workerOpts.DiscoveryMaxLLMCallsPerRun
			if budget > 0 && accum.llmCallsThisRun >= budget {
				log.Warnw("bridgr-worker: LLM score budget exhausted; using deterministic score only",
					"run_uuid", runUUID, "user_id", row.UserID, "url_hash", urlHash,
					"budget", budget)
			} else {
				accum.llmCallsThisRun++
				accum.llmScoreCalls++
				userPayload, sbuildErr := buildScoreLLMUserPayloadWithCV(prefs, j.GetTitle(), j.GetCompany(), jdInline, resumeText, scoreOut)
				if sbuildErr != nil {
					return discoveryBatchStats{}, fmt.Errorf("score llm payload: %w", sbuildErr)
				}
				res, rerr := runJobScoreLLM(ctx, p.prefilterLLM, p.openAIKey, p.workerOpts.DiscoveryOpenAIModel, p.workerOpts.DiscoveryOpenAIBaseURL, userPayload)
				if rerr != nil {
					log.Warnw("bridgr-worker: score LLM failed; using deterministic score only",
						"run_uuid", runUUID, "user_id", row.UserID, "url_hash", urlHash, "error", rerr)
				} else {
					mn := strings.TrimSpace(p.workerOpts.DiscoveryOpenAIModel)
					if mn == "" {
						mn = "openai"
					}
					scoreOut = mergeDeterministicWithLLM(scoreOut, res, mn)
				}
			}
		}

		if float64(scoreOut.CompositeScore)+1e-9 < p.workerOpts.DiscoveryScoreThreshold {
			accum.lowScoreRejected++
			if exErr := p.repo.AddJobDiscoveryExclusion(ctx, p.q, row.UserID, urlHash, "low_score", pgid); exErr != nil {
				log.Errorw("bridgr-worker: add low_score exclusion failed", "url_hash", urlHash, "error", exErr)
			}
			log.Infow("bridgr-worker: discovery score below threshold; excluded",
				"run_uuid", runUUID, "user_id", row.UserID, "url_hash", urlHash,
				"composite_score", scoreOut.CompositeScore, "threshold", p.workerOpts.DiscoveryScoreThreshold,
				"scoring_model", scoreOut.ScoringModel)
			continue
		}

		pl, _ := json.Marshal(map[string]interface{}{
			"title":              j.GetTitle(),
			"company":            j.GetCompany(),
			"requirements":       reqs,
			"discovery_run_uuid": runUUID,
			"composite_score":    scoreOut.CompositeScore,
			"gap_severity":       scoreOut.GapSeverity,
			"findjobs_iteration": iteration,
		})

		*jobSeq++
		sourceID := fmt.Sprintf("findjobs-%d-it%d-%d", row.ID, iteration, *jobSeq)

		cand, uerr := p.repo.UpsertJobCandidate(ctx, p.q, sqlc.UpsertJobCandidateParams{
			UserID:           row.UserID,
			DiscoveryRunUuid: pgid,
			SourceBoard:      sourceBoard,
			SourceJobID:      pgtype.Text{String: sourceID, Valid: true},
			JobUrl:           jobURL,
			UrlHash:          urlHash,
			ContentHash:      pgtype.Text{},
			Title:            pgtype.Text{String: j.GetTitle(), Valid: j.GetTitle() != ""},
			Company:          pgtype.Text{String: j.GetCompany(), Valid: j.GetCompany() != ""},
			Location:         pgtype.Text{},
			JdText:           pgtype.Text{String: jdInline, Valid: jdInline != ""},
			JdS3Uri:          pgtype.Text{},
			FetchedAt:        now,
			IngestionStatus:  "fetched",
			RadarPayload:     pl,
			ApplicationUrl:   pgtype.Text{String: jobURL, Valid: true},
		})
		if uerr != nil {
			log.Errorw("bridgr-worker: upsert job candidate failed", "job_url", jobURL, "error", uerr)
			continue
		}
		enrichUUID := pgtype.UUID{Valid: false}
		scoreRow, serr := p.repo.UpsertJobScore(ctx, p.q, sqlc.UpsertJobScoreParams{
			UserID:               row.UserID,
			JobCandidateUuid:     cand.Uuid,
			EnrichmentUuid:       enrichUUID,
			SkillMatchScore:      scoreOut.SkillMatchScore,
			ExperienceMatchScore: scoreOut.ExperienceMatchScore,
			LocationMatchScore:   scoreOut.LocationMatchScore,
			RecencyScore:         scoreOut.RecencyScore,
			BoardQualityScore:    scoreOut.BoardQualityScore,
			CompositeScore:       scoreOut.CompositeScore,
			MatchedSkills:        scoreOut.MatchedSkillsJSON,
			GapSkills:            scoreOut.GapSkillsJSON,
			GapSeverity:          pgtype.Text{String: scoreOut.GapSeverity, Valid: scoreOut.GapSeverity != ""},
			ScoringModel:         pgtype.Text{String: scoreOut.ScoringModel, Valid: scoreOut.ScoringModel != ""},
			ScoringVersion:       pgtype.Text{String: scoreOut.ScoringVersion, Valid: scoreOut.ScoringVersion != ""},
		})
		if serr != nil {
			log.Errorw("bridgr-worker: upsert job score failed", "job_url", jobURL, "error", serr)
			continue
		}
		if strings.TrimSpace(canonical) == "" {
			log.Warnw("bridgr-worker: skip feed + skill-gap row (no canonical_cv_analysis_uuid)", "run_uuid", runUUID, "user_id", row.UserID, "url_hash", urlHash)
			accum.ingested++
			continue
		}
		if err := p.runDiscoveryFeedPipeline(ctx, row.UserID, pgid, canonical, cand, scoreRow, scoreOut, jdInline, jobURL, prefs.Location, now); err != nil {
			log.Errorw("bridgr-worker: discovery feed pipeline failed", "run_uuid", runUUID, "user_id", row.UserID, "url_hash", urlHash, "error", err)
			continue
		}
		accum.feedUpserted++
		accum.ingested++
	}

	return discoveryBatchStats{
		AcceptedDelta:     accum.ingested - bIn,
		ExclusionSkipped:  accum.exclusionSkipped - bEx,
		PrefilterRejected: accum.prefilterRejected - bPr,
		LowScoreRejected:  accum.lowScoreRejected - bLo,
	}, nil
}

// ################################
// Helpers
// ################################

func (p *Processor) finishDiscoveryRunFailed(ctx context.Context, row *sqlc.BridgrJobSearchDiscoveryRun, now pgtype.Timestamp, code string, cause error) error {
	log := logger.Get(ctx)
	if row == nil {
		if cause == nil {
			return nil
		}
		return cause
	}
	detail := pgtype.Text{}
	if cause != nil {
		detail = pgtype.Text{String: cause.Error(), Valid: true}
	}
	if _, err := p.repo.PatchJobSearchDiscoveryRun(ctx, p.q, sqlc.PatchJobSearchDiscoveryRunParams{
		ID:                row.ID,
		Status:            "failed",
		StartedAt:         row.StartedAt,
		CompletedAt:       now,
		RawCandidateCount: row.RawCandidateCount,
		NewCandidateCount: row.NewCandidateCount,
		RadarMeta:         row.RadarMeta,
		ErrorCode:         pgtype.Text{String: code, Valid: true},
		ErrorDetail:       detail,
		SqsMessageID:      row.SqsMessageID,
	}); err != nil {
		log.Errorw("bridgr-worker: discovery failed state persist error", "error", err)
		return fmt.Errorf("persist discovery failure: %w", err)
	}
	if cause != nil {
		return cause
	}
	return nil
}

func (p *Processor) finishDiscoveryRun(ctx context.Context, row *sqlc.BridgrJobSearchDiscoveryRun, now pgtype.Timestamp, rawCount, newCount int32, radarMeta []byte) error {
	log := logger.Get(ctx)
	if _, err := p.repo.PatchJobSearchDiscoveryRun(ctx, p.q, sqlc.PatchJobSearchDiscoveryRunParams{
		ID:                row.ID,
		Status:            "completed",
		StartedAt:         row.StartedAt,
		CompletedAt:       now,
		RawCandidateCount: rawCount,
		NewCandidateCount: newCount,
		RadarMeta:         radarMeta,
		ErrorCode:         pgtype.Text{},
		ErrorDetail:       pgtype.Text{},
		SqsMessageID:      row.SqsMessageID,
	}); err != nil {
		log.Errorw("bridgr-worker: discovery finish failed", "run_id", row.ID, "error", err)
		if _, uerr := p.repo.PatchJobSearchDiscoveryRun(ctx, p.q, sqlc.PatchJobSearchDiscoveryRunParams{
			ID:                row.ID,
			Status:            "failed",
			StartedAt:         row.StartedAt,
			CompletedAt:       now,
			RawCandidateCount: 0,
			NewCandidateCount: 0,
			RadarMeta:         radarMeta,
			ErrorCode:         pgtype.Text{String: "worker_error", Valid: true},
			ErrorDetail:       pgtype.Text{String: err.Error(), Valid: true},
			SqsMessageID:      row.SqsMessageID,
		}); uerr != nil {
			return fmt.Errorf("%w: persist discovery failure: %w", err, uerr)
		}
		return nil
	}
	return nil
}

// jdTextFromJobInfo flattens Radar JobInfo into stored jd_text (Indeed list crawl yields title/company/requirements, not full HTML JD).
func jdTextFromJobInfo(j *jobsearchv1.JobInfo) string {
	if j == nil {
		return ""
	}
	var parts []string
	for _, r := range j.GetRequirements() {
		t := strings.TrimSpace(r)
		if t != "" {
			parts = append(parts, t)
		}
	}
	if len(parts) == 0 {
		return ""
	}
	s := strings.Join(parts, "\n")
	runes := []rune(s)
	if len(runes) > maxJobCandidateJdTextRunes {
		return string(runes[:maxJobCandidateJdTextRunes]) + "\n…"
	}
	return s
}

// urlHashHex is a stable dedupe key for (user_id, url_hash).
func urlHashHex(jobURL string) string {
	s := strings.TrimSpace(strings.ToLower(jobURL))
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func uniqueURLHashesFromJobInfos(jobs []*jobsearchv1.JobInfo) []string {
	if len(jobs) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(jobs))
	out := make([]string, 0, len(jobs))
	for _, j := range jobs {
		if j == nil {
			continue
		}
		u := strings.TrimSpace(j.GetJobUrl())
		if u == "" {
			continue
		}
		h := urlHashHex(u)
		if _, ok := seen[h]; ok {
			continue
		}
		seen[h] = struct{}{}
		out = append(out, h)
	}
	return out
}
