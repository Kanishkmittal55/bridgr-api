package bridgr_worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	guuid "github.com/gofrs/uuid/v5"
	"github.com/hassleskip/bridgr-api/internal/logger"
	"github.com/hassleskip/bridgr-api/internal/repository"
	"github.com/hassleskip/bridgr-api/internal/repository/sqlc"
	"github.com/hassleskip/bridgr-api/internal/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Processor runs Bridgr skill-gap pipeline stages for one analysis.
// Stubs LLM/extraction/graph construction until wired to real services; still drives status + summary fields.
type Processor struct {
	repo *repository.Repo
	q    sqlc.Querier
}

// NewProcessor builds a processor.
func NewProcessor(repo *repository.Repo, q sqlc.Querier) *Processor {
	return &Processor{repo: repo, q: q}
}

// Process handles one SQS message. Ack on success or unrecoverable bad payload; Nack on transient DB errors.
func (p *Processor) Process(ctx context.Context, msg *Message) error {
	log := logger.Get(ctx)

	var payload AnalysisJobPayload
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		log.Errorw("bridgr-worker: invalid json", "error", err)
		_ = msg.Ack()
		return fmt.Errorf("unmarshal: %w", err)
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

func (p *Processor) runStubPipeline(ctx context.Context, row *sqlc.HskipUsersBridgrSkillGapAnalysis) error {
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
