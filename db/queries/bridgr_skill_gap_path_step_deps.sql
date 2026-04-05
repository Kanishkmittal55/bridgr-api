-- =============================================================================
-- BRIDGR SKILL GAP — PATH STEP DEPENDENCIES
-- DAG edges between path steps (prerequisites).
-- =============================================================================

-- name: CreateSkillGapPathStepDep :one
INSERT INTO hskip_users.bridgr_skill_gap_path_step_deps (
    path_uuid,
    step_uuid,
    depends_on_step_uuid
) VALUES ($1, $2, $3)
RETURNING *;

-- name: BulkCreateSkillGapPathStepDeps :copyfrom
INSERT INTO hskip_users.bridgr_skill_gap_path_step_deps (
    path_uuid,
    step_uuid,
    depends_on_step_uuid
) VALUES (
    $1, $2, $3
);

-- name: ListDependenciesByStep :many
-- Prerequisite steps that step_uuid depends on within this path.
SELECT * FROM hskip_users.bridgr_skill_gap_path_step_deps
WHERE path_uuid = $1 AND step_uuid = $2
ORDER BY id ASC;

-- name: ListDependentsByStep :many
-- Steps in this path that depend on depends_on_step_uuid.
SELECT * FROM hskip_users.bridgr_skill_gap_path_step_deps
WHERE path_uuid = $1 AND depends_on_step_uuid = $2
ORDER BY id ASC;

-- name: ListAllDepsByPath :many
SELECT * FROM hskip_users.bridgr_skill_gap_path_step_deps
WHERE path_uuid = $1
ORDER BY id ASC;

-- name: DeleteSkillGapPathStepDepsByPath :exec
DELETE FROM hskip_users.bridgr_skill_gap_path_step_deps
WHERE path_uuid = $1;
