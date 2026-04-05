// Package repository contains database access functions.
package repository

import (
	"context"
	"fmt"

	"github.com/hassleskip/bridgr-api/internal/repository/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// =============================================================================
// Bridgr Skill Gap — analyses (CV/JD skill gap runs)
// =============================================================================

// CreateSkillGapAnalysis inserts a new skill gap analysis run.
func (r *Repo) CreateSkillGapAnalysis(ctx context.Context, querier sqlc.Querier, params sqlc.CreateSkillGapAnalysisParams) (*sqlc.HskipUsersBridgrSkillGapAnalysis, error) {
	row, err := querier.CreateSkillGapAnalysis(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("CreateSkillGapAnalysis: %w", err)
	}
	return &row, nil
}

// GetSkillGapAnalysis returns an analysis by BIGSERIAL id.
func (r *Repo) GetSkillGapAnalysis(ctx context.Context, querier sqlc.Querier, id int64) (*sqlc.HskipUsersBridgrSkillGapAnalysis, error) {
	row, err := querier.GetSkillGapAnalysis(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetSkillGapAnalysis: %w", err)
	}
	return &row, nil
}

// GetSkillGapAnalysisByUUID returns an analysis by primary UUID.
func (r *Repo) GetSkillGapAnalysisByUUID(ctx context.Context, querier sqlc.Querier, uuid pgtype.UUID) (*sqlc.HskipUsersBridgrSkillGapAnalysis, error) {
	row, err := querier.GetSkillGapAnalysisByUUID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("GetSkillGapAnalysisByUUID: %w", err)
	}
	return &row, nil
}

// GetSkillGapAnalysisByUser lists recent analyses for a user (limit only).
func (r *Repo) GetSkillGapAnalysisByUser(ctx context.Context, querier sqlc.Querier, params sqlc.GetSkillGapAnalysisByUserParams) ([]sqlc.HskipUsersBridgrSkillGapAnalysis, error) {
	rows, err := querier.GetSkillGapAnalysisByUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("GetSkillGapAnalysisByUser: %w", err)
	}
	return rows, nil
}

// GetSkillGapAnalysisByFingerprint finds a run by user + CV/JD fingerprints (idempotency).
func (r *Repo) GetSkillGapAnalysisByFingerprint(ctx context.Context, querier sqlc.Querier, params sqlc.GetSkillGapAnalysisByFingerprintParams) (*sqlc.HskipUsersBridgrSkillGapAnalysis, error) {
	row, err := querier.GetSkillGapAnalysisByFingerprint(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("GetSkillGapAnalysisByFingerprint: %w", err)
	}
	return &row, nil
}

// UpdateSkillGapAnalysisStatus updates pipeline status by analysis id.
func (r *Repo) UpdateSkillGapAnalysisStatus(ctx context.Context, querier sqlc.Querier, params sqlc.UpdateSkillGapAnalysisStatusParams) (*sqlc.HskipUsersBridgrSkillGapAnalysis, error) {
	row, err := querier.UpdateSkillGapAnalysisStatus(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateSkillGapAnalysisStatus: %w", err)
	}
	return &row, nil
}

// UpdateSkillGapAnalysisSummary updates gap_summary and mermaid_diagram.
func (r *Repo) UpdateSkillGapAnalysisSummary(ctx context.Context, querier sqlc.Querier, params sqlc.UpdateSkillGapAnalysisSummaryParams) (*sqlc.HskipUsersBridgrSkillGapAnalysis, error) {
	row, err := querier.UpdateSkillGapAnalysisSummary(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateSkillGapAnalysisSummary: %w", err)
	}
	return &row, nil
}

// UpdateSkillGapAnalysisError sets error fields and status failed.
func (r *Repo) UpdateSkillGapAnalysisError(ctx context.Context, querier sqlc.Querier, params sqlc.UpdateSkillGapAnalysisErrorParams) (*sqlc.HskipUsersBridgrSkillGapAnalysis, error) {
	row, err := querier.UpdateSkillGapAnalysisError(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateSkillGapAnalysisError: %w", err)
	}
	return &row, nil
}

// DeleteSkillGapAnalysis hard-deletes an analysis by id.
func (r *Repo) DeleteSkillGapAnalysis(ctx context.Context, querier sqlc.Querier, id int64) error {
	if err := querier.DeleteSkillGapAnalysis(ctx, id); err != nil {
		return fmt.Errorf("DeleteSkillGapAnalysis: %w", err)
	}
	return nil
}

// DeleteSkillGapAnalysisByUUID hard-deletes an analysis by UUID.
func (r *Repo) DeleteSkillGapAnalysisByUUID(ctx context.Context, querier sqlc.Querier, uuid pgtype.UUID) error {
	if err := querier.DeleteSkillGapAnalysisByUUID(ctx, uuid); err != nil {
		return fmt.Errorf("DeleteSkillGapAnalysisByUUID: %w", err)
	}
	return nil
}

// ListSkillGapAnalysesByUser paginates analyses for a user.
func (r *Repo) ListSkillGapAnalysesByUser(ctx context.Context, querier sqlc.Querier, params sqlc.ListSkillGapAnalysesByUserParams) ([]sqlc.HskipUsersBridgrSkillGapAnalysis, error) {
	rows, err := querier.ListSkillGapAnalysesByUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("ListSkillGapAnalysesByUser: %w", err)
	}
	return rows, nil
}

// =============================================================================
// Bridgr Skill Gap — graphs
// =============================================================================

// CreateSkillGapGraph inserts a candidate or role_requirement graph row.
func (r *Repo) CreateSkillGapGraph(ctx context.Context, querier sqlc.Querier, params sqlc.CreateSkillGapGraphParams) (*sqlc.HskipUsersBridgrSkillGapGraph, error) {
	row, err := querier.CreateSkillGapGraph(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("CreateSkillGapGraph: %w", err)
	}
	return &row, nil
}

// GetSkillGapGraph returns a graph by BIGSERIAL id.
func (r *Repo) GetSkillGapGraph(ctx context.Context, querier sqlc.Querier, id int64) (*sqlc.HskipUsersBridgrSkillGapGraph, error) {
	row, err := querier.GetSkillGapGraph(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetSkillGapGraph: %w", err)
	}
	return &row, nil
}

// GetSkillGapGraphByUUID returns a graph by primary UUID.
func (r *Repo) GetSkillGapGraphByUUID(ctx context.Context, querier sqlc.Querier, uuid pgtype.UUID) (*sqlc.HskipUsersBridgrSkillGapGraph, error) {
	row, err := querier.GetSkillGapGraphByUUID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("GetSkillGapGraphByUUID: %w", err)
	}
	return &row, nil
}

// GetSkillGapGraphsByAnalysis returns all graphs for an analysis (typically candidate + role).
func (r *Repo) GetSkillGapGraphsByAnalysis(ctx context.Context, querier sqlc.Querier, analysisUuid pgtype.UUID) ([]sqlc.HskipUsersBridgrSkillGapGraph, error) {
	rows, err := querier.GetSkillGapGraphsByAnalysis(ctx, analysisUuid)
	if err != nil {
		return nil, fmt.Errorf("GetSkillGapGraphsByAnalysis: %w", err)
	}
	return rows, nil
}

// GetSkillGapGraphByKind returns the graph for an analysis and kind.
func (r *Repo) GetSkillGapGraphByKind(ctx context.Context, querier sqlc.Querier, params sqlc.GetSkillGapGraphByKindParams) (*sqlc.HskipUsersBridgrSkillGapGraph, error) {
	row, err := querier.GetSkillGapGraphByKind(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("GetSkillGapGraphByKind: %w", err)
	}
	return &row, nil
}

// DeleteSkillGapGraphsByAnalysis removes all graph rows for an analysis.
func (r *Repo) DeleteSkillGapGraphsByAnalysis(ctx context.Context, querier sqlc.Querier, analysisUuid pgtype.UUID) error {
	if err := querier.DeleteSkillGapGraphsByAnalysis(ctx, analysisUuid); err != nil {
		return fmt.Errorf("DeleteSkillGapGraphsByAnalysis: %w", err)
	}
	return nil
}

// =============================================================================
// Bridgr Skill Gap — nodes
// =============================================================================

// CreateSkillGapNode inserts a single graph node.
func (r *Repo) CreateSkillGapNode(ctx context.Context, querier sqlc.Querier, params sqlc.CreateSkillGapNodeParams) (*sqlc.HskipUsersBridgrSkillGapNode, error) {
	row, err := querier.CreateSkillGapNode(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("CreateSkillGapNode: %w", err)
	}
	return &row, nil
}

// BulkCreateSkillGapNodes batch-inserts nodes (CopyFrom).
func (r *Repo) BulkCreateSkillGapNodes(ctx context.Context, querier sqlc.Querier, rows []sqlc.BulkCreateSkillGapNodesParams) (int64, error) {
	n, err := querier.BulkCreateSkillGapNodes(ctx, rows)
	if err != nil {
		return 0, fmt.Errorf("BulkCreateSkillGapNodes: %w", err)
	}
	return n, nil
}

// GetSkillGapNode returns a node by BIGSERIAL id.
func (r *Repo) GetSkillGapNode(ctx context.Context, querier sqlc.Querier, id int64) (*sqlc.HskipUsersBridgrSkillGapNode, error) {
	row, err := querier.GetSkillGapNode(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetSkillGapNode: %w", err)
	}
	return &row, nil
}

// GetSkillGapNodeByKey returns a node by graph uuid and node_key.
func (r *Repo) GetSkillGapNodeByKey(ctx context.Context, querier sqlc.Querier, params sqlc.GetSkillGapNodeByKeyParams) (*sqlc.HskipUsersBridgrSkillGapNode, error) {
	row, err := querier.GetSkillGapNodeByKey(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("GetSkillGapNodeByKey: %w", err)
	}
	return &row, nil
}

// ListSkillGapNodesByGraph lists nodes for one graph.
func (r *Repo) ListSkillGapNodesByGraph(ctx context.Context, querier sqlc.Querier, graphUuid pgtype.UUID) ([]sqlc.HskipUsersBridgrSkillGapNode, error) {
	rows, err := querier.ListSkillGapNodesByGraph(ctx, graphUuid)
	if err != nil {
		return nil, fmt.Errorf("ListSkillGapNodesByGraph: %w", err)
	}
	return rows, nil
}

// ListMatchedNodes returns role nodes whose node_key exists on the candidate graph for the analysis.
func (r *Repo) ListMatchedNodes(ctx context.Context, querier sqlc.Querier, analysisUuid pgtype.UUID) ([]sqlc.HskipUsersBridgrSkillGapNode, error) {
	rows, err := querier.ListMatchedNodes(ctx, analysisUuid)
	if err != nil {
		return nil, fmt.Errorf("ListMatchedNodes: %w", err)
	}
	return rows, nil
}

// ListUnmatchedNodes returns role nodes with no candidate match (gaps).
func (r *Repo) ListUnmatchedNodes(ctx context.Context, querier sqlc.Querier, analysisUuid pgtype.UUID) ([]sqlc.HskipUsersBridgrSkillGapNode, error) {
	rows, err := querier.ListUnmatchedNodes(ctx, analysisUuid)
	if err != nil {
		return nil, fmt.Errorf("ListUnmatchedNodes: %w", err)
	}
	return rows, nil
}

// DeleteSkillGapNodesByGraph deletes all nodes for a graph.
func (r *Repo) DeleteSkillGapNodesByGraph(ctx context.Context, querier sqlc.Querier, graphUuid pgtype.UUID) error {
	if err := querier.DeleteSkillGapNodesByGraph(ctx, graphUuid); err != nil {
		return fmt.Errorf("DeleteSkillGapNodesByGraph: %w", err)
	}
	return nil
}

// =============================================================================
// Bridgr Skill Gap — edges
// =============================================================================

// CreateSkillGapEdge inserts a directed edge between two nodes in a graph.
func (r *Repo) CreateSkillGapEdge(ctx context.Context, querier sqlc.Querier, params sqlc.CreateSkillGapEdgeParams) (*sqlc.HskipUsersBridgrSkillGapEdge, error) {
	row, err := querier.CreateSkillGapEdge(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("CreateSkillGapEdge: %w", err)
	}
	return &row, nil
}

// BulkCreateSkillGapEdges batch-inserts edges (CopyFrom).
func (r *Repo) BulkCreateSkillGapEdges(ctx context.Context, querier sqlc.Querier, rows []sqlc.BulkCreateSkillGapEdgesParams) (int64, error) {
	n, err := querier.BulkCreateSkillGapEdges(ctx, rows)
	if err != nil {
		return 0, fmt.Errorf("BulkCreateSkillGapEdges: %w", err)
	}
	return n, nil
}

// ListSkillGapEdgesByGraph lists edges for one graph.
func (r *Repo) ListSkillGapEdgesByGraph(ctx context.Context, querier sqlc.Querier, graphUuid pgtype.UUID) ([]sqlc.HskipUsersBridgrSkillGapEdge, error) {
	rows, err := querier.ListSkillGapEdgesByGraph(ctx, graphUuid)
	if err != nil {
		return nil, fmt.Errorf("ListSkillGapEdgesByGraph: %w", err)
	}
	return rows, nil
}

// ListSkillGapEdgesByNode lists edges incident to a node (from or to) within a graph.
func (r *Repo) ListSkillGapEdgesByNode(ctx context.Context, querier sqlc.Querier, params sqlc.ListSkillGapEdgesByNodeParams) ([]sqlc.HskipUsersBridgrSkillGapEdge, error) {
	rows, err := querier.ListSkillGapEdgesByNode(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("ListSkillGapEdgesByNode: %w", err)
	}
	return rows, nil
}

// DeleteSkillGapEdgesByGraph deletes all edges for a graph.
func (r *Repo) DeleteSkillGapEdgesByGraph(ctx context.Context, querier sqlc.Querier, graphUuid pgtype.UUID) error {
	if err := querier.DeleteSkillGapEdgesByGraph(ctx, graphUuid); err != nil {
		return fmt.Errorf("DeleteSkillGapEdgesByGraph: %w", err)
	}
	return nil
}

// =============================================================================
// Bridgr Skill Gap — learning paths
// =============================================================================

// CreateSkillGapLearningPath inserts a learning path row for an analysis.
func (r *Repo) CreateSkillGapLearningPath(ctx context.Context, querier sqlc.Querier, params sqlc.CreateSkillGapLearningPathParams) (*sqlc.HskipUsersBridgrSkillGapLearningPath, error) {
	row, err := querier.CreateSkillGapLearningPath(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("CreateSkillGapLearningPath: %w", err)
	}
	return &row, nil
}

// GetSkillGapLearningPath returns a path by BIGSERIAL id.
func (r *Repo) GetSkillGapLearningPath(ctx context.Context, querier sqlc.Querier, id int64) (*sqlc.HskipUsersBridgrSkillGapLearningPath, error) {
	row, err := querier.GetSkillGapLearningPath(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetSkillGapLearningPath: %w", err)
	}
	return &row, nil
}

// GetSkillGapLearningPathByUUID returns a learning path by primary UUID.
func (r *Repo) GetSkillGapLearningPathByUUID(ctx context.Context, querier sqlc.Querier, uuid pgtype.UUID) (*sqlc.HskipUsersBridgrSkillGapLearningPath, error) {
	row, err := querier.GetSkillGapLearningPathByUUID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("GetSkillGapLearningPathByUUID: %w", err)
	}
	return &row, nil
}

// GetSkillGapLearningPathByAnalysis returns the latest path_version for an analysis.
func (r *Repo) GetSkillGapLearningPathByAnalysis(ctx context.Context, querier sqlc.Querier, analysisUuid pgtype.UUID) (*sqlc.HskipUsersBridgrSkillGapLearningPath, error) {
	row, err := querier.GetSkillGapLearningPathByAnalysis(ctx, analysisUuid)
	if err != nil {
		return nil, fmt.Errorf("GetSkillGapLearningPathByAnalysis: %w", err)
	}
	return &row, nil
}

// UpdateSkillGapLearningPathMetadata updates path_metadata JSON.
func (r *Repo) UpdateSkillGapLearningPathMetadata(ctx context.Context, querier sqlc.Querier, params sqlc.UpdateSkillGapLearningPathMetadataParams) (*sqlc.HskipUsersBridgrSkillGapLearningPath, error) {
	row, err := querier.UpdateSkillGapLearningPathMetadata(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateSkillGapLearningPathMetadata: %w", err)
	}
	return &row, nil
}

// DeleteSkillGapLearningPathByAnalysis deletes all path versions for an analysis.
func (r *Repo) DeleteSkillGapLearningPathByAnalysis(ctx context.Context, querier sqlc.Querier, analysisUuid pgtype.UUID) error {
	if err := querier.DeleteSkillGapLearningPathByAnalysis(ctx, analysisUuid); err != nil {
		return fmt.Errorf("DeleteSkillGapLearningPathByAnalysis: %w", err)
	}
	return nil
}

// =============================================================================
// Bridgr Skill Gap — path steps
// =============================================================================

// CreateSkillGapPathStep inserts one step.
func (r *Repo) CreateSkillGapPathStep(ctx context.Context, querier sqlc.Querier, params sqlc.CreateSkillGapPathStepParams) (*sqlc.HskipUsersBridgrSkillGapPathStep, error) {
	row, err := querier.CreateSkillGapPathStep(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("CreateSkillGapPathStep: %w", err)
	}
	return &row, nil
}

// BulkCreateSkillGapPathSteps batch-inserts steps (CopyFrom).
func (r *Repo) BulkCreateSkillGapPathSteps(ctx context.Context, querier sqlc.Querier, rows []sqlc.BulkCreateSkillGapPathStepsParams) (int64, error) {
	n, err := querier.BulkCreateSkillGapPathSteps(ctx, rows)
	if err != nil {
		return 0, fmt.Errorf("BulkCreateSkillGapPathSteps: %w", err)
	}
	return n, nil
}

// GetSkillGapPathStep returns a step by BIGSERIAL id.
func (r *Repo) GetSkillGapPathStep(ctx context.Context, querier sqlc.Querier, id int64) (*sqlc.HskipUsersBridgrSkillGapPathStep, error) {
	row, err := querier.GetSkillGapPathStep(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetSkillGapPathStep: %w", err)
	}
	return &row, nil
}

// ListSkillGapPathStepsByPath lists steps ordered by step_index.
func (r *Repo) ListSkillGapPathStepsByPath(ctx context.Context, querier sqlc.Querier, pathUuid pgtype.UUID) ([]sqlc.HskipUsersBridgrSkillGapPathStep, error) {
	rows, err := querier.ListSkillGapPathStepsByPath(ctx, pathUuid)
	if err != nil {
		return nil, fmt.Errorf("ListSkillGapPathStepsByPath: %w", err)
	}
	return rows, nil
}

// UpdateSkillGapPathStep attaches resource_uri / resource_kind to a step.
func (r *Repo) UpdateSkillGapPathStep(ctx context.Context, querier sqlc.Querier, params sqlc.UpdateSkillGapPathStepParams) (*sqlc.HskipUsersBridgrSkillGapPathStep, error) {
	row, err := querier.UpdateSkillGapPathStep(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateSkillGapPathStep: %w", err)
	}
	return &row, nil
}

// DeleteSkillGapPathStepsByPath deletes all steps for a path.
func (r *Repo) DeleteSkillGapPathStepsByPath(ctx context.Context, querier sqlc.Querier, pathUuid pgtype.UUID) error {
	if err := querier.DeleteSkillGapPathStepsByPath(ctx, pathUuid); err != nil {
		return fmt.Errorf("DeleteSkillGapPathStepsByPath: %w", err)
	}
	return nil
}

// =============================================================================
// Bridgr Skill Gap — path step dependencies (DAG edges between steps)
// =============================================================================

// CreateSkillGapPathStepDep inserts one dependency edge.
func (r *Repo) CreateSkillGapPathStepDep(ctx context.Context, querier sqlc.Querier, params sqlc.CreateSkillGapPathStepDepParams) (*sqlc.HskipUsersBridgrSkillGapPathStepDep, error) {
	row, err := querier.CreateSkillGapPathStepDep(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("CreateSkillGapPathStepDep: %w", err)
	}
	return &row, nil
}

// BulkCreateSkillGapPathStepDeps batch-inserts deps (CopyFrom).
func (r *Repo) BulkCreateSkillGapPathStepDeps(ctx context.Context, querier sqlc.Querier, rows []sqlc.BulkCreateSkillGapPathStepDepsParams) (int64, error) {
	n, err := querier.BulkCreateSkillGapPathStepDeps(ctx, rows)
	if err != nil {
		return 0, fmt.Errorf("BulkCreateSkillGapPathStepDeps: %w", err)
	}
	return n, nil
}

// ListDependenciesByStep lists prerequisite steps for step_uuid.
func (r *Repo) ListDependenciesByStep(ctx context.Context, querier sqlc.Querier, params sqlc.ListDependenciesByStepParams) ([]sqlc.HskipUsersBridgrSkillGapPathStepDep, error) {
	rows, err := querier.ListDependenciesByStep(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("ListDependenciesByStep: %w", err)
	}
	return rows, nil
}

// ListDependentsByStep lists steps that depend on depends_on_step_uuid.
func (r *Repo) ListDependentsByStep(ctx context.Context, querier sqlc.Querier, params sqlc.ListDependentsByStepParams) ([]sqlc.HskipUsersBridgrSkillGapPathStepDep, error) {
	rows, err := querier.ListDependentsByStep(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("ListDependentsByStep: %w", err)
	}
	return rows, nil
}

// ListAllDepsByPath returns full DAG edge list for Mermaid / topo sort.
func (r *Repo) ListAllDepsByPath(ctx context.Context, querier sqlc.Querier, pathUuid pgtype.UUID) ([]sqlc.HskipUsersBridgrSkillGapPathStepDep, error) {
	rows, err := querier.ListAllDepsByPath(ctx, pathUuid)
	if err != nil {
		return nil, fmt.Errorf("ListAllDepsByPath: %w", err)
	}
	return rows, nil
}

// DeleteSkillGapPathStepDepsByPath deletes all deps for a path.
func (r *Repo) DeleteSkillGapPathStepDepsByPath(ctx context.Context, querier sqlc.Querier, pathUuid pgtype.UUID) error {
	if err := querier.DeleteSkillGapPathStepDepsByPath(ctx, pathUuid); err != nil {
		return fmt.Errorf("DeleteSkillGapPathStepDepsByPath: %w", err)
	}
	return nil
}

// =============================================================================
// Bridgr Skill Gap — coverage aggregates (derived from graphs)
// =============================================================================

// GetSkillGapCoverage computes candidate vs required counts and coverage % for one analysis.
func (r *Repo) GetSkillGapCoverage(ctx context.Context, querier sqlc.Querier, analysisUuid pgtype.UUID) (*sqlc.GetSkillGapCoverageRow, error) {
	row, err := querier.GetSkillGapCoverage(ctx, analysisUuid)
	if err != nil {
		return nil, fmt.Errorf("GetSkillGapCoverage: %w", err)
	}
	return &row, nil
}

// GetSkillGapCoverageByUser lists analyses for a user with coverage metrics (paginated).
func (r *Repo) GetSkillGapCoverageByUser(ctx context.Context, querier sqlc.Querier, params sqlc.GetSkillGapCoverageByUserParams) ([]sqlc.GetSkillGapCoverageByUserRow, error) {
	rows, err := querier.GetSkillGapCoverageByUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("GetSkillGapCoverageByUser: %w", err)
	}
	return rows, nil
}
