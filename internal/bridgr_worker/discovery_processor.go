package bridgr_worker

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Kanishkmittal55/bridgr-api/internal/cloud"
	ierrr "github.com/Kanishkmittal55/bridgr-api/internal/errors"

	jobsearchv1 "github.com/Kanishkmittal55/bridgr-api/internal/gen/radar/services/job_search/v1"
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
}

// NewProcessor builds a processor. jobSearch may be nil when RADAR_ADDR is unset.
func NewProcessor(repo *repository.Repo, q sqlc.Querier, jobSearch *radar.JobSearchClient) *Processor {
	return &Processor{repo: repo, q: q, jobSearchClient: jobSearch}
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
		// TODO : where NULL means “leave unchanged”).
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

	findJobReq, sourceBoard, err := findJobsRequestFromJobDiscoveryRequestParams(row.RequestParams)
	if err != nil {
		return p.finishDiscoveryRunFailed(ctx, row, now, "wrong_discovery_params_bad_request", fmt.Errorf("bad discovery params in enqueued job: %w: %w", err, ierrr.ErrBadRequest))
	}

	log.Infow("bridgr-worker: calling Radar JobSearchService.FindJobs", "run_uuid", payload.RunUUID, "target_role", findJobReq.GetSearchQuery(), "location", findJobReq.GetLocation(), "max_results", findJobReq.GetMaxResults(), "source_board", sourceBoard)

	resp, err := p.jobSearchClient.FindJobs(ctx, findJobReq)
	if err != nil {
		return p.finishDiscoveryRunFailed(ctx, row, now, "radar_find_jobs_rpc_call_error", fmt.Errorf("radar client find jobs rpc call failed: %w: %w", err, ierrr.ErrInternal))
	}

	jobs := resp.GetJobs()
	if jobs == nil {
		jobs = []*jobsearchv1.JobInfo{}
		stubMeta, _ := json.Marshal(map[string]interface{}{"stub": true, "worker": "bridgr", "radar_skipped": false, "did_not_find_any_jobs": true})
		return p.finishDiscoveryRun(ctx, row, now, 0, 0, stubMeta) // There was no error radar just did not find any jobs for that request job search profile.
	}

	rawCount := int32(len(jobs))
	var ingested int32
	for i, j := range jobs {

		jobURL := strings.TrimSpace(j.GetJobUrl())
		if jobURL == "" {
			continue
		}

		reqs := j.GetRequirements()
		if reqs == nil {
			reqs = []string{}
		}

		pl, _ := json.Marshal(map[string]interface{}{
			"title":              j.GetTitle(),
			"company":            j.GetCompany(),
			"requirements":       reqs,
			"discovery_run_uuid": payload.RunUUID,
		})

		jdInline := jdTextFromJobInfo(j)

		suffix := ""

		if i > 0 {
			suffix = fmt.Sprintf("-%d", i)
		}

		// TODO: Finally at this point when we have the job description from one job we need to enforce the request params we set
		// If the job follows the discovery params it is worthy for a job score
		// If the job score is above the threshold then we do a skill analyses and upsert ? Why upsert here we should only insert
		// Also a limit on when to stop , till the time we dont find 10 good candidates belong to the same runUUID we keep looking
		// We have a growing exclusion list such that if a URL was visited and rejected or accepted earlier then dont vist that again

		_, uerr := p.repo.UpsertJobCandidate(ctx, p.q, sqlc.UpsertJobCandidateParams{
			UserID:           row.UserID,
			DiscoveryRunUuid: pgid,
			SourceBoard:      sourceBoard,
			SourceJobID:      pgtype.Text{String: fmt.Sprintf("findjobs-%d%s", row.ID, suffix), Valid: true},
			JobUrl:           jobURL,
			UrlHash:          urlHashHex(jobURL),
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
		ingested++
	}

	// Then we finally add the metadata to the run uuid
	meta, _ := json.Marshal(map[string]interface{}{
		"radar":               true,
		"find_jobs":           true,
		"jobs_returned":       len(jobs),
		"candidates_upserted": ingested,
		"source_board":        sourceBoard,
		"target_role":         findJobReq.GetSearchQuery(),
		"location":            findJobReq.GetLocation(),
	})

	// And we successfully finish the run
	return p.finishDiscoveryRun(ctx, row, now, rawCount, ingested, meta)
}

// ################################
// Helpers
// ################################

// findJobsRequestFromDiscoveryParams maps bridgr.job_search_discovery_runs.request_params JSON
// (from BuildDiscoveryRequestParams) to Radar FindJobsRequest. Returns preferred source_board label.
func findJobsRequestFromJobDiscoveryRequestParams(requestParams []byte) (*jobsearchv1.FindJobsRequest, string, error) {
	if len(requestParams) == 0 {
		return nil, "radar", fmt.Errorf("empty request_params")
	}
	var m map[string]interface{}
	if err := json.Unmarshal(requestParams, &m); err != nil {
		return nil, "radar", err
	}

	searchQuery := strings.TrimSpace(stringFromMap(m, "target_role"))
	if searchQuery == "" {
		searchQuery = strings.TrimSpace(stringFromMap(m, "search_query"))
	}
	var targetRoles []string
	if searchQuery != "" {
		targetRoles = []string{searchQuery}
	} else {
		targetRoles = stringSliceFromMixed("target_roles", m)
		searchQuery = strings.TrimSpace(strings.Join(targetRoles, " "))
	}
	if searchQuery == "" {
		searchQuery = "software engineer"
		targetRoles = []string{searchQuery}
	}

	loc := strings.TrimSpace(stringFromMap(m, "location"))
	if loc == "" {
		loc = firstLocationString(m)
	}

	maxN := maxResultsFromParams(m)

	board := strings.TrimSpace(stringFromMap(m, "source_board"))
	if board == "" {
		board = firstBoard(m)
	}

	req := &jobsearchv1.FindJobsRequest{
		TargetRoles: targetRoles,
		SearchQuery: searchQuery,
		Location:    loc,
		MaxResults:  maxN,
	}
	return req, board, nil
}

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

func stringFromMap(m map[string]interface{}, key string) string {
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}

func stringSliceFromMixed(key string, m map[string]interface{}) []string {
	v, ok := m[key]
	if !ok || v == nil {
		return nil
	}
	raw, ok := v.([]interface{})
	if !ok {
		return nil
	}
	out := make([]string, 0, len(raw))
	for _, it := range raw {
		switch x := it.(type) {
		case string:
			if s := strings.TrimSpace(x); s != "" {
				out = append(out, s)
			}
		case float64:
			out = append(out, strconv.FormatInt(int64(x), 10))
		}
	}
	return out
}

func firstLocationString(m map[string]interface{}) string {
	v, ok := m["locations"]
	if !ok || v == nil {
		return ""
	}
	raw, ok := v.([]interface{})
	if !ok {
		return ""
	}
	for _, it := range raw {
		switch x := it.(type) {
		case string:
			if s := strings.TrimSpace(x); s != "" {
				return s
			}
		case map[string]interface{}:
			if s, _ := x["location"].(string); strings.TrimSpace(s) != "" {
				return strings.TrimSpace(s)
			}
		}
	}
	return ""
}

func maxResultsFromParams(m map[string]interface{}) int32 {
	v, ok := m["max_surfaced_jobs"]
	if !ok {
		return 10
	}
	switch x := v.(type) {
	case float64:
		if x > 0 && x <= 100 {
			return int32(x)
		}
	case int:
		if x > 0 && x <= 100 {
			return int32(x)
		}
	case int32:
		if x > 0 && x <= 100 {
			return x
		}
	}
	return 10
}

func firstBoard(m map[string]interface{}) string {
	v, ok := m["boards_enabled"]
	if !ok || v == nil {
		return "indeed"
	}
	raw, ok := v.([]interface{})
	if !ok {
		return "indeed"
	}
	for _, it := range raw {
		if s, ok := it.(string); ok && strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
	}
	return "indeed"
}
