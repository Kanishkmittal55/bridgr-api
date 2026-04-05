package bridgr

import (
	"context"
	"fmt"
	"net/http"

	apierrors "github.com/Kanishkmittal55/bridgr-api/internal/apierrors"
	types "github.com/Kanishkmittal55/bridgr-api/pkg/types"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// V1GetBridgrAnalysis handles GET /v1/bridgr/analyses/{analysisUUID}
func (s *server) V1GetBridgrAnalysis(w http.ResponseWriter, r *http.Request, analysisUUID openapi_types.UUID, _ types.V1GetBridgrAnalysisParams) {
	ctx := r.Context()
	resp, err := s.v1GetBridgrAnalysis(ctx, analysisUUID)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

func (s *server) v1GetBridgrAnalysis(ctx context.Context, analysisUUID openapi_types.UUID) (*types.BridgrSkillGapAnalysis, error) {
	_, row, err := s.loadAnalysis(ctx, analysisUUID)
	if err != nil {
		return nil, err
	}
	out, err := analysisFromRow(row)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", apierrors.ErrInternal, err)
	}
	return &out, nil
}
