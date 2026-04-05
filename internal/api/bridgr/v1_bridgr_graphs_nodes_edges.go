package bridgr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/hassleskip/bridgr-api/internal/repository/sqlc"
	"github.com/hassleskip/bridgr-api/internal/uuid"
	types "github.com/hassleskip/bridgr-api/pkg/types"
	hserr "github.com/hassleskip/hassle-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// V1GetBridgrGraphNodes handles GET /v1/bridgr/graphs/{graphUUID}/nodes
func (s *server) V1GetBridgrGraphNodes(w http.ResponseWriter, r *http.Request, graphUUID openapi_types.UUID, _ types.V1GetBridgrGraphNodesParams) {
	ctx := r.Context()
	resp, err := s.v1GetBridgrGraphNodes(ctx, graphUUID)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

func (s *server) v1GetBridgrGraphNodes(ctx context.Context, graphUUID openapi_types.UUID) (*types.BridgrSkillGapNodeListResponse, error) {
	pgid, _, err := s.loadGraph(ctx, graphUUID)
	if err != nil {
		return nil, err
	}
	rows, err := s.deps.Repo.ListSkillGapNodesByGraph(ctx, s.querier(), pgid)
	if err != nil {
		return nil, fmt.Errorf("%w: list nodes: %w", hserr.ErrInternal, err)
	}
	out := types.BridgrSkillGapNodeListResponse{Nodes: make([]types.BridgrSkillGapNode, 0, len(rows))}
	for i := range rows {
		n, err := nodeFromRow(&rows[i])
		if err != nil {
			return nil, fmt.Errorf("%w: map node: %w", hserr.ErrInternal, err)
		}
		out.Nodes = append(out.Nodes, n)
	}
	return &out, nil
}

// V1GetBridgrGraphEdges handles GET /v1/bridgr/graphs/{graphUUID}/edges
func (s *server) V1GetBridgrGraphEdges(w http.ResponseWriter, r *http.Request, graphUUID openapi_types.UUID, _ types.V1GetBridgrGraphEdgesParams) {
	ctx := r.Context()
	resp, err := s.v1GetBridgrGraphEdges(ctx, graphUUID)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

func (s *server) v1GetBridgrGraphEdges(ctx context.Context, graphUUID openapi_types.UUID) (*types.BridgrSkillGapEdgeListResponse, error) {
	pgid, _, err := s.loadGraph(ctx, graphUUID)
	if err != nil {
		return nil, err
	}
	rows, err := s.deps.Repo.ListSkillGapEdgesByGraph(ctx, s.querier(), pgid)
	if err != nil {
		return nil, fmt.Errorf("%w: list edges: %w", hserr.ErrInternal, err)
	}
	out := types.BridgrSkillGapEdgeListResponse{Edges: make([]types.BridgrSkillGapEdge, 0, len(rows))}
	for i := range rows {
		e, err := edgeFromRow(&rows[i])
		if err != nil {
			return nil, fmt.Errorf("%w: map edge: %w", hserr.ErrInternal, err)
		}
		out.Edges = append(out.Edges, e)
	}
	return &out, nil
}

// V1GetBridgrGraphNodeByKey handles GET /v1/bridgr/graphs/{graphUUID}/nodes/by-key/{nodeKey}
func (s *server) V1GetBridgrGraphNodeByKey(w http.ResponseWriter, r *http.Request, graphUUID openapi_types.UUID, nodeKey string, _ types.V1GetBridgrGraphNodeByKeyParams) {
	ctx := r.Context()
	resp, err := s.v1GetBridgrGraphNodeByKey(ctx, graphUUID, nodeKey)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

func (s *server) v1GetBridgrGraphNodeByKey(ctx context.Context, graphUUID openapi_types.UUID, nodeKey string) (*types.BridgrSkillGapNode, error) {
	pgid, _, err := s.loadGraph(ctx, graphUUID)
	if err != nil {
		return nil, err
	}
	row, err := s.deps.Repo.GetSkillGapNodeByKey(ctx, s.querier(), sqlc.GetSkillGapNodeByKeyParams{
		GraphUuid: pgid,
		NodeKey:   nodeKey,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: node not found", hserr.ErrNotFound)
		}
		return nil, fmt.Errorf("%w: %w", hserr.ErrInternal, err)
	}
	out, err := nodeFromRow(row)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", hserr.ErrInternal, err)
	}
	return &out, nil
}

// V1PostBridgrGraphNodes handles POST /v1/bridgr/graphs/{graphUUID}/nodes
func (s *server) V1PostBridgrGraphNodes(w http.ResponseWriter, r *http.Request, graphUUID openapi_types.UUID, _ types.V1PostBridgrGraphNodesParams) {
	ctx := r.Context()
	var body types.CreateBridgrSkillGapNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, nil, fmt.Errorf("%w: invalid JSON body: %v", hserr.ErrBadRequest, err))
		return
	}
	resp, err := s.v1PostBridgrGraphNodes(ctx, graphUUID, body)
	if err != nil {
		s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, nil, err)
		return
	}
	s.writeCreated(w, r, resp)
}

func (s *server) v1PostBridgrGraphNodes(ctx context.Context, graphUUID openapi_types.UUID, body types.CreateBridgrSkillGapNodeRequest) (*types.BridgrSkillGapNode, error) {
	if body.DisplayName == "" || body.NodeKey == "" {
		return nil, fmt.Errorf("%w: display_name and node_key are required", hserr.ErrBadRequest)
	}
	pgid, _, err := s.loadGraph(ctx, graphUUID)
	if err != nil {
		return nil, err
	}
	evidence, err := mapToJSONBytes(body.Evidence)
	if err != nil {
		return nil, fmt.Errorf("%w: evidence: %v", hserr.ErrBadRequest, err)
	}
	meta, err := mapToJSONBytes(body.Metadata)
	if err != nil {
		return nil, fmt.Errorf("%w: metadata: %v", hserr.ErrBadRequest, err)
	}
	params := sqlc.CreateSkillGapNodeParams{
		GraphUuid:       pgid,
		NodeKey:         body.NodeKey,
		DisplayName:     body.DisplayName,
		Evidence:        evidence,
		Metadata:        meta,
		Description:     pgtype.Text{},
		ProficiencyHint: pgtype.Text{},
		Source:          pgtype.Text{},
		PositionX:       pgtype.Int4{},
		PositionY:       pgtype.Int4{},
	}
	if body.Description != nil {
		params.Description = pgtype.Text{String: *body.Description, Valid: true}
	}
	if body.ProficiencyHint != nil {
		params.ProficiencyHint = pgtype.Text{String: *body.ProficiencyHint, Valid: true}
	}
	if body.Source != nil {
		params.Source = pgtype.Text{String: *body.Source, Valid: true}
	}
	if body.PositionX != nil {
		params.PositionX = pgtype.Int4{Int32: *body.PositionX, Valid: true}
	}
	if body.PositionY != nil {
		params.PositionY = pgtype.Int4{Int32: *body.PositionY, Valid: true}
	}
	row, err := s.deps.Repo.CreateSkillGapNode(ctx, s.querier(), params)
	if err != nil {
		if bridgrFKOrNotFound(err) {
			return nil, fmt.Errorf("%w: graph not found or invalid reference", hserr.ErrNotFound)
		}
		return nil, fmt.Errorf("%w: create node: %w", hserr.ErrInternal, err)
	}
	out, err := nodeFromRow(row)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", hserr.ErrInternal, err)
	}
	return &out, nil
}

// V1PostBridgrGraphEdges handles POST /v1/bridgr/graphs/{graphUUID}/edges
func (s *server) V1PostBridgrGraphEdges(w http.ResponseWriter, r *http.Request, graphUUID openapi_types.UUID, _ types.V1PostBridgrGraphEdgesParams) {
	ctx := r.Context()
	var body types.CreateBridgrSkillGapEdgeRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, nil, fmt.Errorf("%w: invalid JSON body: %v", hserr.ErrBadRequest, err))
		return
	}
	resp, err := s.v1PostBridgrGraphEdges(ctx, graphUUID, body)
	if err != nil {
		s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, nil, err)
		return
	}
	s.writeCreated(w, r, resp)
}

func (s *server) v1PostBridgrGraphEdges(ctx context.Context, graphUUID openapi_types.UUID, body types.CreateBridgrSkillGapEdgeRequest) (*types.BridgrSkillGapEdge, error) {
	if body.Relation == "" {
		return nil, fmt.Errorf("%w: relation is required", hserr.ErrBadRequest)
	}
	pgid, _, err := s.loadGraph(ctx, graphUUID)
	if err != nil {
		return nil, err
	}
	fromPg, err := uuid.ConvertOapiUUIDToPgUUID(body.FromNodeUuid)
	if err != nil {
		return nil, fmt.Errorf("%w: from_node_uuid: %v", hserr.ErrBadRequest, err)
	}
	toPg, err := uuid.ConvertOapiUUIDToPgUUID(body.ToNodeUuid)
	if err != nil {
		return nil, fmt.Errorf("%w: to_node_uuid: %v", hserr.ErrBadRequest, err)
	}
	meta, err := mapToJSONBytes(body.Metadata)
	if err != nil {
		return nil, fmt.Errorf("%w: metadata: %v", hserr.ErrBadRequest, err)
	}
	row, err := s.deps.Repo.CreateSkillGapEdge(ctx, s.querier(), sqlc.CreateSkillGapEdgeParams{
		GraphUuid:    pgid,
		FromNodeUuid: fromPg,
		ToNodeUuid:   toPg,
		Relation:     body.Relation,
		Weight:       float32ToNumericPtr(body.Weight),
		Metadata:     meta,
	})
	if err != nil {
		if bridgrFKOrNotFound(err) {
			return nil, fmt.Errorf("%w: graph or node reference not found", hserr.ErrNotFound)
		}
		return nil, fmt.Errorf("%w: create edge: %w", hserr.ErrInternal, err)
	}
	out, err := edgeFromRow(row)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", hserr.ErrInternal, err)
	}
	return &out, nil
}

func bridgrFKOrNotFound(err error) bool {
	var pe *pgconn.PgError
	if errors.As(err, &pe) {
		if pe.Code == "23503" {
			return true
		}
	}
	return errors.Is(err, pgx.ErrNoRows)
}
