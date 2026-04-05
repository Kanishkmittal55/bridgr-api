package bridgr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/hassleskip/bridgr-api/internal/repository/sqlc"
	types "github.com/hassleskip/bridgr-api/pkg/types"
	hserr "github.com/hassleskip/hassle-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// V1PatchBridgrAnalysisStatus handles PATCH /v1/bridgr/analyses/{analysisUUID}/status
func (s *server) V1PatchBridgrAnalysisStatus(w http.ResponseWriter, r *http.Request, analysisUUID openapi_types.UUID, _ types.V1PatchBridgrAnalysisStatusParams) {
	ctx := r.Context()
	var body types.UpdateBridgrSkillGapAnalysisStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, nil, fmt.Errorf("%w: invalid JSON body: %v", hserr.ErrBadRequest, err))
		return
	}
	resp, err := s.v1PatchBridgrAnalysisStatus(ctx, analysisUUID, body)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

func (s *server) v1PatchBridgrAnalysisStatus(ctx context.Context, analysisUUID openapi_types.UUID, body types.UpdateBridgrSkillGapAnalysisStatusRequest) (*types.BridgrSkillGapAnalysis, error) {
	_, row, err := s.loadAnalysis(ctx, analysisUUID)
	if err != nil {
		return nil, err
	}
	updated, err := s.deps.Repo.UpdateSkillGapAnalysisStatus(ctx, s.querier(), sqlc.UpdateSkillGapAnalysisStatusParams{
		ID:     row.ID,
		Status: string(body.Status),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: analysis not found", hserr.ErrNotFound)
		}
		return nil, fmt.Errorf("%w: update status: %w", hserr.ErrInternal, err)
	}
	out, err := analysisFromRow(updated)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", hserr.ErrInternal, err)
	}
	return &out, nil
}
