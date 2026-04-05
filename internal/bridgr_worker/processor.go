package bridgr_worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

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

// Processor runs Bridgr skill-gap pipeline stages for one analysis.
// Stubs LLM/extraction/graph construction until wired to real services; still drives status + summary fields.
type Processor struct {
	repo      *repository.Repo
	q         sqlc.Querier
	jobSearch *radar.JobSearchClient // optional; nil => job discovery completes without Radar (stub)
}

// NewProcessor builds a processor. jobSearch may be nil when RADAR_ADDR is unset.
func NewProcessor(repo *repository.Repo, q sqlc.Querier, jobSearch *radar.JobSearchClient) *Processor {
	return &Processor{repo: repo, q: q, jobSearch: jobSearch}
}

// Process handles one SQS message. Ack on success or unrecoverable bad payload; Nack on transient DB errors.
func (p *Processor) Process(ctx context.Context, msg *Message) error {
	log := logger.Get(ctx)

	var payload QueuePayload
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		log.Errorw("bridgr-worker: invalid json", "error", err)
		_ = msg.Ack()
		return fmt.Errorf("unmarshal: %w", err)
	}
	if payload.Kind == KindJobDiscovery {
		if payload.RunUUID == "" || payload.UserID == 0 {
			log.Errorw("bridgr-worker: job_discovery missing run_uuid or user_id")
			_ = msg.Ack()
			return errors.New("job_discovery: missing fields")
		}
		return p.processJobDiscovery(ctx, msg, payload)
	}
	if payload.Kind == "" {
		payload.Kind = KindSkillGapAnalysis
	}
	if payload.AnalysisUUID == "" {
		log.Errorw("bridgr-worker: missing analysis_uuid")
		_ = msg.Ack()
		return errors.New("missing analysis_uuid")
	}

	uid, err := guuid.FromString(payload.AnalysisUUID)
	if err != nil {
		log.Errorw("bridgr-worker: bad uuid", "analysis_uuid", payload.AnalysisUUID, "error", err)
		_ = msg.Ack()
		return err
	}
	pgid := uuid.ToPgUuid(uid)

	row, err := p.repo.GetSkillGapAnalysisByUUID(ctx, p.q, pgid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warnw("bridgr-worker: analysis gone, acking", "analysis_uuid", payload.AnalysisUUID)
			return msg.Ack()
		}
		msg.Nack()
		return fmt.Errorf("load analysis: %w", err)
	}

	// Idempotency: finished terminal states
	switch row.Status {
	case "completed", "failed":
		log.Infow("bridgr-worker: skip terminal analysis", "status", row.Status, "id", row.ID)
		return msg.Ack()
	}

	if err := p.runStubPipeline(ctx, row); err != nil {
		log.Errorw("bridgr-worker: pipeline failed", "analysis_id", row.ID, "error", err)
		_, uerr := p.repo.UpdateSkillGapAnalysisError(ctx, p.q, sqlc.UpdateSkillGapAnalysisErrorParams{
			ID:          row.ID,
			ErrorCode:   pgtype.Text{String: "pipeline_error", Valid: true},
			ErrorDetail: pgtype.Text{String: err.Error(), Valid: true},
		})
		if uerr != nil {
			log.Errorw("bridgr-worker: failed to persist error", "error", uerr)
			msg.Nack()
			return fmt.Errorf("%w: update error row: %w", err, uerr)
		}
		return msg.Ack()
	}

	return msg.Ack()
}

func (p *Processor) processJobDiscovery(ctx context.Context, msg *Message, payload QueuePayload) error {
	log := logger.Get(ctx)
	runUID, err := guuid.FromString(payload.RunUUID)
	if err != nil {
		log.Errorw("bridgr-worker: bad run uuid", "run_uuid", payload.RunUUID, "error", err)
		_ = msg.Ack()
		return err
	}
	pgid := uuid.ToPgUuid(runUID)
	row, err := p.repo.GetJobSearchDiscoveryRunByUUID(ctx, p.q, pgid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warnw("bridgr-worker: discovery run gone, acking", "run_uuid", payload.RunUUID)
			return msg.Ack()
		}
		msg.Nack()
		return fmt.Errorf("load discovery run: %w", err)
	}
	if row.UserID != payload.UserID {
		log.Errorw("bridgr-worker: job_discovery user_id mismatch", "run_user", row.UserID, "payload_user", payload.UserID)
		_ = msg.Ack()
		return errors.New("user_id mismatch")
	}

	switch row.Status {
	case "completed", "failed", "cancelled":
		log.Infow("bridgr-worker: skip terminal discovery run", "status", row.Status, "id", row.ID)
		return msg.Ack()
	}

	log.Infow(
		"bridgr-worker: job discovery run received",
		"run_uuid", payload.RunUUID,
		"user_id", payload.UserID,
		"status", row.Status,
		"run_id", row.ID,
	)

	now := pgtype.Timestamp{Time: time.Now().UTC(), Valid: true}
	if row.Status == "pending" || row.Status == "queued" {
		if _, err := p.repo.UpdateJobSearchDiscoveryRunStarted(ctx, p.q, sqlc.UpdateJobSearchDiscoveryRunStartedParams{
			ID:           row.ID,
			Status:       "running",
			StartedAt:    now,
			SqsMessageID: row.SqsMessageID,
		}); err != nil {
			msg.Nack()
			return fmt.Errorf("mark discovery started: %w", err)
		}
	}

	if p.jobSearch == nil {
		log.Warnw("bridgr-worker: RADAR_ADDR not configured on worker; discovery completes without Radar FindJobs")
		stubMeta, _ := json.Marshal(map[string]interface{}{"stub": true, "worker": "bridgr", "radar_skipped": true})
		return p.finishDiscoveryRun(ctx, msg, row.ID, now, 0, 0, stubMeta)
	}

	fjReq, sourceBoard, err := findJobsRequestFromDiscoveryParams(row.RequestParams)
	if err != nil {
		log.Errorw("bridgr-worker: discovery request_params", "error", err)
		return p.finishDiscoveryRunFailed(ctx, msg, row, now, "bad_request", err.Error())
	}

	log.Infow(
		"bridgr-worker: calling Radar JobSearchService.FindJobs",
		"run_uuid", payload.RunUUID,
		"user_id", payload.UserID,
		"search_query", fjReq.GetSearchQuery(),
		"location", fjReq.GetLocation(),
		"max_results", fjReq.GetMaxResults(),
		"source_board", sourceBoard,
	)

	resp, err := p.jobSearch.FindJobs(ctx, fjReq)
	if err != nil {
		log.Errorw("bridgr-worker: FindJobs failed", "error", err)
		return p.finishDiscoveryRunFailed(ctx, msg, row, now, "radar_find_jobs", err.Error())
	}
	jobs := resp.GetJobs()
	if jobs == nil {
		jobs = []*jobsearchv1.JobInfo{}
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
			"title":         j.GetTitle(),
			"company":       j.GetCompany(),
			"requirements":  reqs,
			"discovery_run": payload.RunUUID,
		})
		jdInline := jdTextFromJobInfo(j)
		suffix := ""
		if i > 0 {
			suffix = fmt.Sprintf("-%d", i)
		}
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

	meta, _ := json.Marshal(map[string]interface{}{
		"radar":               true,
		"find_jobs":           true,
		"jobs_returned":       len(jobs),
		"candidates_upserted": ingested,
		"source_board":        sourceBoard,
		"search_query":        fjReq.GetSearchQuery(),
		"location":            fjReq.GetLocation(),
	})
	return p.finishDiscoveryRun(ctx, msg, row.ID, now, rawCount, ingested, meta)
}

func (p *Processor) finishDiscoveryRunFailed(
	ctx context.Context,
	msg *Message,
	row *sqlc.BridgrJobSearchDiscoveryRun,
	now pgtype.Timestamp,
	code, detail string,
) error {
	log := logger.Get(ctx)
	if _, err := p.repo.UpdateJobSearchDiscoveryRunFinished(ctx, p.q, sqlc.UpdateJobSearchDiscoveryRunFinishedParams{
		ID:                row.ID,
		Status:            "failed",
		CompletedAt:       now,
		RawCandidateCount: row.RawCandidateCount,
		NewCandidateCount: row.NewCandidateCount,
		RadarMeta:         row.RadarMeta,
		ErrorCode:         pgtype.Text{String: code, Valid: true},
		ErrorDetail:       pgtype.Text{String: detail, Valid: true},
	}); err != nil {
		log.Errorw("bridgr-worker: discovery failed state persist error", "error", err)
		msg.Nack()
		return fmt.Errorf("persist discovery failure: %w", err)
	}
	return msg.Ack()
}

func (p *Processor) finishDiscoveryRun(
	ctx context.Context,
	msg *Message,
	runRowID int64,
	now pgtype.Timestamp,
	rawCount, newCount int32,
	radarMeta []byte,
) error {
	log := logger.Get(ctx)
	if _, err := p.repo.UpdateJobSearchDiscoveryRunFinished(ctx, p.q, sqlc.UpdateJobSearchDiscoveryRunFinishedParams{
		ID:                runRowID,
		Status:            "completed",
		CompletedAt:       now,
		RawCandidateCount: rawCount,
		NewCandidateCount: newCount,
		RadarMeta:         radarMeta,
		ErrorCode:         pgtype.Text{},
		ErrorDetail:       pgtype.Text{},
	}); err != nil {
		log.Errorw("bridgr-worker: discovery finish failed", "run_id", runRowID, "error", err)
		if _, uerr := p.repo.UpdateJobSearchDiscoveryRunFinished(ctx, p.q, sqlc.UpdateJobSearchDiscoveryRunFinishedParams{
			ID:                runRowID,
			Status:            "failed",
			CompletedAt:       now,
			RawCandidateCount: 0,
			NewCandidateCount: 0,
			RadarMeta:         radarMeta,
			ErrorCode:         pgtype.Text{String: "worker_error", Valid: true},
			ErrorDetail:       pgtype.Text{String: err.Error(), Valid: true},
		}); uerr != nil {
			msg.Nack()
			return fmt.Errorf("%w: persist discovery failure: %w", err, uerr)
		}
		return msg.Ack()
	}
	return msg.Ack()
}

func (p *Processor) runStubPipeline(ctx context.Context, row *sqlc.BridgrSkillGapAnalysis) error {
	id := row.ID

	if _, err := p.repo.UpdateSkillGapAnalysisStatus(ctx, p.q, sqlc.UpdateSkillGapAnalysisStatusParams{
		ID: id, Status: "extracting",
	}); err != nil {
		return fmt.Errorf("status extracting: %w", err)
	}

	stubExtraction, err := json.Marshal(map[string]interface{}{
		"stub":    true,
		"stage":   "extract",
		"version": 1,
	})
	if err != nil {
		return fmt.Errorf("marshal extraction stub: %w", err)
	}
	if _, err := p.repo.UpdateSkillGapAnalysisSummary(ctx, p.q, sqlc.UpdateSkillGapAnalysisSummaryParams{
		ID:             id,
		GapSummary:     stubExtraction,
		MermaidDiagram: pgtype.Text{},
	}); err != nil {
		return fmt.Errorf("summary after extract: %w", err)
	}

	if _, err := p.repo.UpdateSkillGapAnalysisStatus(ctx, p.q, sqlc.UpdateSkillGapAnalysisStatusParams{
		ID: id, Status: "graphed",
	}); err != nil {
		return fmt.Errorf("status graphed: %w", err)
	}

	stubGraph, err := json.Marshal(map[string]interface{}{
		"stub":  true,
		"stage": "graph",
		"nodes": []string{},
	})
	if err != nil {
		return fmt.Errorf("marshal graph stub: %w", err)
	}
	if _, err := p.repo.UpdateSkillGapAnalysisSummary(ctx, p.q, sqlc.UpdateSkillGapAnalysisSummaryParams{
		ID:             id,
		GapSummary:     stubGraph,
		MermaidDiagram: pgtype.Text{String: "graph TD\n  A[stub]", Valid: true},
	}); err != nil {
		return fmt.Errorf("summary after graph: %w", err)
	}

	if _, err := p.repo.UpdateSkillGapAnalysisStatus(ctx, p.q, sqlc.UpdateSkillGapAnalysisStatusParams{
		ID: id, Status: "pathed",
	}); err != nil {
		return fmt.Errorf("status pathed: %w", err)
	}

	if _, err := p.repo.UpdateSkillGapAnalysisStatus(ctx, p.q, sqlc.UpdateSkillGapAnalysisStatusParams{
		ID: id, Status: "completed",
	}); err != nil {
		return fmt.Errorf("status completed: %w", err)
	}

	return nil
}
