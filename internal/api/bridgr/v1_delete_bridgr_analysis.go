package bridgr

import (
	"context"
	"fmt"
	"net/http"

	apierrors "github.com/Kanishkmittal55/bridgr-api/internal/apierrors"
	types "github.com/Kanishkmittal55/bridgr-api/pkg/types"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// V1DeleteBridgrAnalysis handles DELETE /v1/bridgr/analyses/{analysisUUID}
func (s *server) V1DeleteBridgrAnalysis(w http.ResponseWriter, r *http.Request, analysisUUID openapi_types.UUID, _ types.V1DeleteBridgrAnalysisParams) {
	ctx := r.Context()
	err := s.v1DeleteBridgrAnalysis(ctx, analysisUUID)
	if err != nil {
		s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, nil, err)
		return
	}
	s.deps.ResponseWriter.WriteNoContentResponse(ctx, w, r, nil)
}

func (s *server) v1DeleteBridgrAnalysis(ctx context.Context, analysisUUID openapi_types.UUID) error {
	pgid, _, err := s.loadAnalysis(ctx, analysisUUID)
	if err != nil {
		return err
	}
	if err := s.deps.Repo.DeleteSkillGapAnalysisByUUID(ctx, s.querier(), pgid); err != nil {
		return fmt.Errorf("%w: delete analysis: %w", apierrors.ErrInternal, err)
	}
	return nil
}
