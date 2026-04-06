package bridgr

import (
	"context"
	"fmt"
	"net/http"

	apierrors "github.com/Kanishkmittal55/bridgr-api/internal/errors"
	types "github.com/Kanishkmittal55/bridgr-api/pkg/types"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// V1GetBridgrPathSteps handles GET /v1/bridgr/paths/{pathUUID}/steps
func (s *server) V1GetBridgrPathSteps(w http.ResponseWriter, r *http.Request, pathUUID openapi_types.UUID, _ types.V1GetBridgrPathStepsParams) {
	ctx := r.Context()
	resp, err := s.v1GetBridgrPathSteps(ctx, pathUUID)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

func (s *server) v1GetBridgrPathSteps(ctx context.Context, pathUUID openapi_types.UUID) (*types.BridgrSkillGapPathStepListResponse, error) {
	pgid, _, err := s.loadPath(ctx, pathUUID)
	if err != nil {
		return nil, err
	}
	rows, err := s.deps.Repo.ListSkillGapPathStepsByPath(ctx, s.querier(), pgid)
	if err != nil {
		return nil, fmt.Errorf("%w: list steps: %w", apierrors.ErrInternal, err)
	}
	out := types.BridgrSkillGapPathStepListResponse{Steps: make([]types.BridgrSkillGapPathStep, 0, len(rows))}
	for i := range rows {
		st, err := pathStepFromRow(&rows[i])
		if err != nil {
			return nil, fmt.Errorf("%w: map step: %w", apierrors.ErrInternal, err)
		}
		out.Steps = append(out.Steps, st)
	}
	return &out, nil
}
