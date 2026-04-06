package bridgr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	apierrors "github.com/Kanishkmittal55/bridgr-api/internal/errors"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository/sqlc"
	types "github.com/Kanishkmittal55/bridgr-api/pkg/types"
	"github.com/jackc/pgx/v5"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// V1GetBridgrAnalysisCoverage handles GET /v1/bridgr/analyses/{analysisUUID}/coverage
func (s *server) V1GetBridgrAnalysisCoverage(w http.ResponseWriter, r *http.Request, analysisUUID openapi_types.UUID, _ types.V1GetBridgrAnalysisCoverageParams) {
	ctx := r.Context()
	resp, err := s.v1GetBridgrAnalysisCoverage(ctx, analysisUUID)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

func (s *server) v1GetBridgrAnalysisCoverage(ctx context.Context, analysisUUID openapi_types.UUID) (*types.BridgrSkillGapCoverageSnapshot, error) {
	pgid, _, err := s.loadAnalysis(ctx, analysisUUID)
	if err != nil {
		return nil, err
	}
	row, err := s.deps.Repo.GetSkillGapCoverage(ctx, s.querier(), pgid)
	if err != nil {
		return nil, fmt.Errorf("%w: coverage: %w", apierrors.ErrInternal, err)
	}
	out := types.BridgrSkillGapCoverageSnapshot{
		CandidateSkills: countToInterfacePtr(row.CandidateSkills),
		RequiredSkills:  countToInterfacePtr(row.RequiredSkills),
		MatchedSkills:   countToInterfacePtr(row.MatchedSkills),
		CoveragePct:     numericToFloat32Ptr(row.CoveragePct),
	}
	return &out, nil
}

// V1GetBridgrAnalysisGraphs handles GET /v1/bridgr/analyses/{analysisUUID}/graphs
func (s *server) V1GetBridgrAnalysisGraphs(w http.ResponseWriter, r *http.Request, analysisUUID openapi_types.UUID, _ types.V1GetBridgrAnalysisGraphsParams) {
	ctx := r.Context()
	resp, err := s.v1GetBridgrAnalysisGraphs(ctx, analysisUUID)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

func (s *server) v1GetBridgrAnalysisGraphs(ctx context.Context, analysisUUID openapi_types.UUID) (*types.BridgrSkillGapGraphListResponse, error) {
	pgid, _, err := s.loadAnalysis(ctx, analysisUUID)
	if err != nil {
		return nil, err
	}
	rows, err := s.deps.Repo.GetSkillGapGraphsByAnalysis(ctx, s.querier(), pgid)
	if err != nil {
		return nil, fmt.Errorf("%w: list graphs: %w", apierrors.ErrInternal, err)
	}
	out := types.BridgrSkillGapGraphListResponse{Graphs: make([]types.BridgrSkillGapGraph, 0, len(rows))}
	for i := range rows {
		g, err := graphFromRow(&rows[i])
		if err != nil {
			return nil, fmt.Errorf("%w: map graph: %w", apierrors.ErrInternal, err)
		}
		out.Graphs = append(out.Graphs, g)
	}
	return &out, nil
}

// V1GetBridgrAnalysisGraphByKind handles GET /v1/bridgr/analyses/{analysisUUID}/graphs/{kind}
func (s *server) V1GetBridgrAnalysisGraphByKind(w http.ResponseWriter, r *http.Request, analysisUUID openapi_types.UUID, kind types.BridgrSkillGapGraphKind, _ types.V1GetBridgrAnalysisGraphByKindParams) {
	ctx := r.Context()
	resp, err := s.v1GetBridgrAnalysisGraphByKind(ctx, analysisUUID, kind)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

func (s *server) v1GetBridgrAnalysisGraphByKind(ctx context.Context, analysisUUID openapi_types.UUID, kind types.BridgrSkillGapGraphKind) (*types.BridgrSkillGapGraph, error) {
	pgid, _, err := s.loadAnalysis(ctx, analysisUUID)
	if err != nil {
		return nil, err
	}
	row, err := s.deps.Repo.GetSkillGapGraphByKind(ctx, s.querier(), sqlc.GetSkillGapGraphByKindParams{
		AnalysisUuid: pgid,
		Kind:         string(kind),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: graph not found for kind", apierrors.ErrNotFound)
		}
		return nil, fmt.Errorf("%w: %w", apierrors.ErrInternal, err)
	}
	out, err := graphFromRow(row)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", apierrors.ErrInternal, err)
	}
	return &out, nil
}

// V1PostBridgrAnalysisGraphs handles POST /v1/bridgr/analyses/{analysisUUID}/graphs
func (s *server) V1PostBridgrAnalysisGraphs(w http.ResponseWriter, r *http.Request, analysisUUID openapi_types.UUID, _ types.V1PostBridgrAnalysisGraphsParams) {
	ctx := r.Context()
	var body types.CreateBridgrSkillGapGraphRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, nil, fmt.Errorf("%w: invalid JSON body: %v", apierrors.ErrBadRequest, err))
		return
	}
	resp, err := s.v1PostBridgrAnalysisGraphs(ctx, analysisUUID, body)
	if err != nil {
		s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, nil, err)
		return
	}
	s.writeCreated(w, r, resp)
}

func (s *server) v1PostBridgrAnalysisGraphs(ctx context.Context, analysisUUID openapi_types.UUID, body types.CreateBridgrSkillGapGraphRequest) (*types.BridgrSkillGapGraph, error) {
	pgid, _, err := s.loadAnalysis(ctx, analysisUUID)
	if err != nil {
		return nil, err
	}
	meta, err := mapToJSONBytes(body.Metadata)
	if err != nil {
		return nil, fmt.Errorf("%w: metadata: %v", apierrors.ErrBadRequest, err)
	}
	row, err := s.deps.Repo.CreateSkillGapGraph(ctx, s.querier(), sqlc.CreateSkillGapGraphParams{
		AnalysisUuid: pgid,
		Kind:         string(body.Kind),
		Metadata:     meta,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: create graph: %w", apierrors.ErrInternal, err)
	}
	out, err := graphFromRow(row)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", apierrors.ErrInternal, err)
	}
	return &out, nil
}

// V1DeleteBridgrAnalysisGraphs handles DELETE /v1/bridgr/analyses/{analysisUUID}/graphs
func (s *server) V1DeleteBridgrAnalysisGraphs(w http.ResponseWriter, r *http.Request, analysisUUID openapi_types.UUID, _ types.V1DeleteBridgrAnalysisGraphsParams) {
	ctx := r.Context()
	err := s.v1DeleteBridgrAnalysisGraphs(ctx, analysisUUID)
	if err != nil {
		s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, nil, err)
		return
	}
	s.deps.ResponseWriter.WriteNoContentResponse(ctx, w, r, nil)
}

func (s *server) v1DeleteBridgrAnalysisGraphs(ctx context.Context, analysisUUID openapi_types.UUID) error {
	pgid, _, err := s.loadAnalysis(ctx, analysisUUID)
	if err != nil {
		return err
	}
	if err := s.deps.Repo.DeleteSkillGapGraphsByAnalysis(ctx, s.querier(), pgid); err != nil {
		return fmt.Errorf("%w: delete graphs: %w", apierrors.ErrInternal, err)
	}
	return nil
}
