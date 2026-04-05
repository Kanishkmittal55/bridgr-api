package bridgr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	guuid "github.com/gofrs/uuid/v5"
	"github.com/hassleskip/bridgr-api/internal/bridgr_worker"
	"github.com/hassleskip/bridgr-api/internal/config"
	"github.com/hassleskip/bridgr-api/internal/repository/sqlc"
	"github.com/hassleskip/bridgr-api/internal/uuid"
	types "github.com/hassleskip/bridgr-api/pkg/types"
	hserr "github.com/hassleskip/hassle-go/pkg/errors"
	"github.com/jackc/pgx/v5/pgtype"
)

// V1PostBridgrAnalyses handles POST /v1/bridgr/analyses
func (s *server) V1PostBridgrAnalyses(w http.ResponseWriter, r *http.Request, _ types.V1PostBridgrAnalysesParams) {
	ctx := r.Context()
	var payload types.CreateBridgrSkillGapAnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, nil, fmt.Errorf("%w: invalid JSON body: %v", hserr.ErrBadRequest, err))
		return
	}
	resp, err := s.v1PostBridgrAnalyses(ctx, payload)
	if err != nil {
		s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, nil, err)
		return
	}
	s.writeCreated(w, r, resp)
}

func (s *server) v1PostBridgrAnalyses(ctx context.Context, payload types.CreateBridgrSkillGapAnalysisRequest) (*types.BridgrSkillGapAnalysis, error) {
	if err := s.requireStore(); err != nil {
		return nil, err
	}
	if payload.UserId == 0 {
		return nil, fmt.Errorf("%w: user_id is required", hserr.ErrBadRequest)
	}
	params := sqlc.CreateSkillGapAnalysisParams{
		UserID: payload.UserId,
		Title:  pgtype.Text{},
		Status: string(types.BridgrSkillGapPending),
	}
	if payload.Title != nil {
		params.Title = pgtype.Text{String: *payload.Title, Valid: true}
	}
	if payload.CvAssetUri != nil {
		params.CvAssetUri = pgtype.Text{String: *payload.CvAssetUri, Valid: true}
	}
	if payload.JdAssetUri != nil {
		params.JdAssetUri = pgtype.Text{String: *payload.JdAssetUri, Valid: true}
	}
	if payload.CvFingerprint != nil {
		params.CvFingerprint = pgtype.Text{String: *payload.CvFingerprint, Valid: true}
	}
	if payload.JdFingerprint != nil {
		params.JdFingerprint = pgtype.Text{String: *payload.JdFingerprint, Valid: true}
	}
	if payload.LlmModel != nil {
		params.LlmModel = pgtype.Text{String: *payload.LlmModel, Valid: true}
	}
	if payload.PromptVersion != nil {
		params.PromptVersion = pgtype.Text{String: *payload.PromptVersion, Valid: true}
	}
	if payload.FounderPersonaUuid != nil {
		pg, err := uuid.ConvertOapiUUIDToPgUUID(*payload.FounderPersonaUuid)
		if err != nil {
			return nil, fmt.Errorf("%w: founder_persona_uuid: %v", hserr.ErrBadRequest, err)
		}
		params.FounderPersonaUuid = pg
	}
	if payload.PursuitUuid != nil {
		pg, err := uuid.ConvertOapiUUIDToPgUUID(*payload.PursuitUuid)
		if err != nil {
			return nil, fmt.Errorf("%w: pursuit_uuid: %v", hserr.ErrBadRequest, err)
		}
		params.PursuitUuid = pg
	}
	row, err := s.deps.Repo.CreateSkillGapAnalysis(ctx, s.querier(), params)
	if err != nil {
		return nil, fmt.Errorf("%w: create analysis: %w", hserr.ErrInternal, err)
	}

	cfg := config.Get()
	if cfg.BridgrQueueURL != "" && s.deps.SQSClient != nil {
		uid, uerr := guuid.FromBytes(row.Uuid.Bytes[:])
		if uerr != nil {
			return nil, fmt.Errorf("%w: analysis uuid: %w", hserr.ErrInternal, uerr)
		}
		if err := bridgr_worker.EnqueueSkillGapAnalysis(ctx, s.deps.SQSClient, cfg.BridgrQueueURL, uid); err != nil {
			return nil, fmt.Errorf("%w: enqueue skill-gap job: %w", hserr.ErrInternal, err)
		}
	}

	out, err := analysisFromRow(row)
	if err != nil {
		return nil, fmt.Errorf("%w: map response: %w", hserr.ErrInternal, err)
	}
	return &out, nil
}
