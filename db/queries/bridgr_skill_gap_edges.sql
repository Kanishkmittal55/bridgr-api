-- =============================================================================
-- BRIDGR SKILL GAP — EDGES
-- Directed edges between nodes in one graph. All uuids app-enforced FKs.
-- =============================================================================

-- name: CreateSkillGapEdge :one
INSERT INTO bridgr.skill_gap_edges (
    graph_uuid,
    from_node_uuid,
    to_node_uuid,
    relation,
    weight,
    metadata
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: BulkCreateSkillGapEdges :copyfrom
INSERT INTO bridgr.skill_gap_edges (
    graph_uuid,
    from_node_uuid,
    to_node_uuid,
    relation,
    weight,
    metadata
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: ListSkillGapEdgesByGraph :many
SELECT * FROM bridgr.skill_gap_edges
WHERE graph_uuid = $1
ORDER BY id ASC;

-- name: ListSkillGapEdgesByNode :many
SELECT * FROM bridgr.skill_gap_edges
WHERE graph_uuid = $1
  AND (from_node_uuid = $2 OR to_node_uuid = $2)
ORDER BY id ASC;

-- name: DeleteSkillGapEdgesByGraph :exec
DELETE FROM bridgr.skill_gap_edges
WHERE graph_uuid = $1;
