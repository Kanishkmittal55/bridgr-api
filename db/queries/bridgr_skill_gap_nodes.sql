-- =============================================================================
-- BRIDGR SKILL GAP — NODES
-- Skill/concept nodes per graph. graph_uuid app-enforced FK.
-- =============================================================================

-- name: CreateSkillGapNode :one
INSERT INTO hskip_users.bridgr_skill_gap_nodes (
    graph_uuid,
    node_key,
    display_name,
    description,
    proficiency_hint,
    source,
    evidence,
    metadata,
    position_x,
    position_y
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: BulkCreateSkillGapNodes :copyfrom
INSERT INTO hskip_users.bridgr_skill_gap_nodes (
    graph_uuid,
    node_key,
    display_name,
    description,
    proficiency_hint,
    source,
    evidence,
    metadata,
    position_x,
    position_y
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
);

-- name: GetSkillGapNode :one
SELECT * FROM hskip_users.bridgr_skill_gap_nodes
WHERE id = $1;

-- name: GetSkillGapNodeByKey :one
SELECT * FROM hskip_users.bridgr_skill_gap_nodes
WHERE graph_uuid = $1 AND node_key = $2;

-- name: ListSkillGapNodesByGraph :many
SELECT * FROM hskip_users.bridgr_skill_gap_nodes
WHERE graph_uuid = $1
ORDER BY node_key ASC;

-- name: ListMatchedNodes :many
-- Role-requirement nodes whose node_key also exists on the candidate graph for the same analysis.
SELECT rn.* FROM hskip_users.bridgr_skill_gap_nodes rn
INNER JOIN hskip_users.bridgr_skill_gap_graphs rg
    ON rg.uuid = rn.graph_uuid AND rg.kind = 'role_requirement'
INNER JOIN hskip_users.bridgr_skill_gap_graphs cg
    ON cg.analysis_uuid = rg.analysis_uuid AND cg.kind = 'candidate'
INNER JOIN hskip_users.bridgr_skill_gap_nodes cn
    ON cn.graph_uuid = cg.uuid AND cn.node_key = rn.node_key
WHERE rg.analysis_uuid = $1
ORDER BY rn.node_key ASC;

-- name: ListUnmatchedNodes :many
-- Role-requirement nodes with no candidate node sharing the same node_key (skill gaps).
SELECT rn.* FROM hskip_users.bridgr_skill_gap_nodes rn
INNER JOIN hskip_users.bridgr_skill_gap_graphs rg
    ON rg.uuid = rn.graph_uuid AND rg.kind = 'role_requirement'
WHERE rg.analysis_uuid = $1
  AND NOT EXISTS (
    SELECT 1 FROM hskip_users.bridgr_skill_gap_graphs cg
    INNER JOIN hskip_users.bridgr_skill_gap_nodes cn ON cn.graph_uuid = cg.uuid
    WHERE cg.analysis_uuid = rg.analysis_uuid
      AND cg.kind = 'candidate'
      AND cn.node_key = rn.node_key
  )
ORDER BY rn.node_key ASC;

-- name: DeleteSkillGapNodesByGraph :exec
DELETE FROM hskip_users.bridgr_skill_gap_nodes
WHERE graph_uuid = $1;
