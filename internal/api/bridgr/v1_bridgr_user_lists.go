package bridgr

import (
	"context"
	"fmt"
	"net/http"

	apierrors "github.com/Kanishkmittal55/bridgr-api/internal/errors"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository/sqlc"
	types "github.com/Kanishkmittal55/bridgr-api/pkg/types"
)

// V1GetBridgrUserAnalyses handles GET /v1/bridgr/users/{userID}/analyses
func (s *server) V1GetBridgrUserAnalyses(w http.ResponseWriter, r *http.Request, userID int32, params types.V1GetBridgrUserAnalysesParams) {
	ctx := r.Context()
	resp, err := s.v1GetBridgrUserAnalyses(ctx, userID, params)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

func (s *server) v1GetBridgrUserAnalyses(ctx context.Context, userID int32, params types.V1GetBridgrUserAnalysesParams) (*types.BridgrSkillGapAnalysisListResponse, error) {
	if err := s.requireStore(); err != nil {
		return nil, err
	}
	limit, offset := bridgrUserLimitOffset(params.Limit, params.Offset)
	rows, err := s.deps.Repo.ListSkillGapAnalysesByUser(ctx, s.querier(), sqlc.ListSkillGapAnalysesByUserParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: list analyses: %w", apierrors.ErrInternal, err)
	}
	out := types.BridgrSkillGapAnalysisListResponse{Analyses: make([]types.BridgrSkillGapAnalysis, 0, len(rows))}
	for i := range rows {
		a, err := analysisFromRow(&rows[i])
		if err != nil {
			return nil, fmt.Errorf("%w: map analysis: %w", apierrors.ErrInternal, err)
		}
		out.Analyses = append(out.Analyses, a)
	}
	return &out, nil
}

// V1GetBridgrUserCoverage handles GET /v1/bridgr/users/{userID}/coverage
func (s *server) V1GetBridgrUserCoverage(w http.ResponseWriter, r *http.Request, userID int32, params types.V1GetBridgrUserCoverageParams) {
	ctx := r.Context()
	resp, err := s.v1GetBridgrUserCoverage(ctx, userID, params)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

func (s *server) v1GetBridgrUserCoverage(ctx context.Context, userID int32, params types.V1GetBridgrUserCoverageParams) (*types.BridgrSkillGapUserCoverageListResponse, error) {
	if err := s.requireStore(); err != nil {
		return nil, err
	}
	limit, offset := bridgrUserLimitOffset(params.Limit, params.Offset)
	rows, err := s.deps.Repo.GetSkillGapCoverageByUser(ctx, s.querier(), sqlc.GetSkillGapCoverageByUserParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: list coverage: %w", apierrors.ErrInternal, err)
	}
	out := types.BridgrSkillGapUserCoverageListResponse{Items: make([]types.BridgrSkillGapUserCoverageRow, 0, len(rows))}
	for i := range rows {
		item, err := coverageUserRowFromSQL(rows[i])
		if err != nil {
			return nil, fmt.Errorf("%w: map coverage row: %w", apierrors.ErrInternal, err)
		}
		out.Items = append(out.Items, item)
	}
	return &out, nil
}
