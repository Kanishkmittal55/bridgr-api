package bridgr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	apierrors "github.com/Kanishkmittal55/bridgr-api/internal/apierrors"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository/sqlc"
	types "github.com/Kanishkmittal55/bridgr-api/pkg/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// V1GetBridgrAnalysisLearningPath handles GET /v1/bridgr/analyses/{analysisUUID}/learning-path
func (s *server) V1GetBridgrAnalysisLearningPath(w http.ResponseWriter, r *http.Request, analysisUUID openapi_types.UUID, _ types.V1GetBridgrAnalysisLearningPathParams) {
	ctx := r.Context()
	resp, err := s.v1GetBridgrAnalysisLearningPath(ctx, analysisUUID)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

func (s *server) v1GetBridgrAnalysisLearningPath(ctx context.Context, analysisUUID openapi_types.UUID) (*types.BridgrSkillGapLearningPath, error) {
	pgid, _, err := s.loadAnalysis(ctx, analysisUUID)
	if err != nil {
		return nil, err
	}
	row, err := s.deps.Repo.GetSkillGapLearningPathByAnalysis(ctx, s.querier(), pgid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: learning path not found", apierrors.ErrNotFound)
		}
		return nil, fmt.Errorf("%w: %w", apierrors.ErrInternal, err)
	}
	out, err := learningPathFromRow(row)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", apierrors.ErrInternal, err)
	}
	return &out, nil
}

// V1PostBridgrAnalysisLearningPath handles POST /v1/bridgr/analyses/{analysisUUID}/learning-path
func (s *server) V1PostBridgrAnalysisLearningPath(w http.ResponseWriter, r *http.Request, analysisUUID openapi_types.UUID, _ types.V1PostBridgrAnalysisLearningPathParams) {
	ctx := r.Context()
	var body types.CreateBridgrSkillGapLearningPathRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, nil, fmt.Errorf("%w: invalid JSON body: %v", apierrors.ErrBadRequest, err))
		return
	}
	resp, err := s.v1PostBridgrAnalysisLearningPath(ctx, analysisUUID, body)
	if err != nil {
		s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, nil, err)
		return
	}
	s.writeCreated(w, r, resp)
}

func (s *server) v1PostBridgrAnalysisLearningPath(ctx context.Context, analysisUUID openapi_types.UUID, body types.CreateBridgrSkillGapLearningPathRequest) (*types.BridgrSkillGapLearningPath, error) {
	pgid, _, err := s.loadAnalysis(ctx, analysisUUID)
	if err != nil {
		return nil, err
	}
	pathVer := int32(1)
	if body.PathVersion != nil {
		pathVer = *body.PathVersion
	}
	meta, err := mapToJSONBytes(body.PathMetadata)
	if err != nil {
		return nil, fmt.Errorf("%w: path_metadata: %v", apierrors.ErrBadRequest, err)
	}
	params := sqlc.CreateSkillGapLearningPathParams{
		AnalysisUuid: pgid,
		PathVersion:  pathVer,
		PathMetadata: meta,
		Algorithm:    pgtype.Text{},
		Title:        pgtype.Text{},
	}
	if body.Algorithm != nil {
		params.Algorithm = pgtype.Text{String: *body.Algorithm, Valid: true}
	}
	if body.Title != nil {
		params.Title = pgtype.Text{String: *body.Title, Valid: true}
	}
	row, err := s.deps.Repo.CreateSkillGapLearningPath(ctx, s.querier(), params)
	if err != nil {
		return nil, fmt.Errorf("%w: create learning path: %w", apierrors.ErrInternal, err)
	}
	out, err := learningPathFromRow(row)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", apierrors.ErrInternal, err)
	}
	return &out, nil
}

// V1DeleteBridgrAnalysisLearningPath handles DELETE /v1/bridgr/analyses/{analysisUUID}/learning-path
func (s *server) V1DeleteBridgrAnalysisLearningPath(w http.ResponseWriter, r *http.Request, analysisUUID openapi_types.UUID, _ types.V1DeleteBridgrAnalysisLearningPathParams) {
	ctx := r.Context()
	err := s.v1DeleteBridgrAnalysisLearningPath(ctx, analysisUUID)
	if err != nil {
		s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, nil, err)
		return
	}
	s.deps.ResponseWriter.WriteNoContentResponse(ctx, w, r, nil)
}

func (s *server) v1DeleteBridgrAnalysisLearningPath(ctx context.Context, analysisUUID openapi_types.UUID) error {
	pgid, _, err := s.loadAnalysis(ctx, analysisUUID)
	if err != nil {
		return err
	}
	if err := s.deps.Repo.DeleteSkillGapLearningPathByAnalysis(ctx, s.querier(), pgid); err != nil {
		return fmt.Errorf("%w: delete learning paths: %w", apierrors.ErrInternal, err)
	}
	return nil
}

// V1GetBridgrAnalysisNodesMatched handles GET /v1/bridgr/analyses/{analysisUUID}/nodes/matched
func (s *server) V1GetBridgrAnalysisNodesMatched(w http.ResponseWriter, r *http.Request, analysisUUID openapi_types.UUID, _ types.V1GetBridgrAnalysisNodesMatchedParams) {
	ctx := r.Context()
	resp, err := s.v1ListAnalysisGapNodes(ctx, analysisUUID, true)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

// V1GetBridgrAnalysisNodesUnmatched handles GET /v1/bridgr/analyses/{analysisUUID}/nodes/unmatched
func (s *server) V1GetBridgrAnalysisNodesUnmatched(w http.ResponseWriter, r *http.Request, analysisUUID openapi_types.UUID, _ types.V1GetBridgrAnalysisNodesUnmatchedParams) {
	ctx := r.Context()
	resp, err := s.v1ListAnalysisGapNodes(ctx, analysisUUID, false)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

func (s *server) v1ListAnalysisGapNodes(ctx context.Context, analysisUUID openapi_types.UUID, matched bool) (*types.BridgrSkillGapNodeListResponse, error) {
	pgid, _, err := s.loadAnalysis(ctx, analysisUUID)
	if err != nil {
		return nil, err
	}
	var rows []sqlc.HskipUsersBridgrSkillGapNode
	if matched {
		rows, err = s.deps.Repo.ListMatchedNodes(ctx, s.querier(), pgid)
	} else {
		rows, err = s.deps.Repo.ListUnmatchedNodes(ctx, s.querier(), pgid)
	}
	if err != nil {
		return nil, fmt.Errorf("%w: list nodes: %w", apierrors.ErrInternal, err)
	}
	out := types.BridgrSkillGapNodeListResponse{Nodes: make([]types.BridgrSkillGapNode, 0, len(rows))}
	for i := range rows {
		n, err := nodeFromRow(&rows[i])
		if err != nil {
			return nil, fmt.Errorf("%w: map node: %w", apierrors.ErrInternal, err)
		}
		out.Nodes = append(out.Nodes, n)
	}
	return &out, nil
}
