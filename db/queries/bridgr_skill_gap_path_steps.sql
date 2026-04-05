-- =============================================================================
-- BRIDGR SKILL GAP — PATH STEPS
-- Ordered steps inside a learning path.
-- =============================================================================

-- name: CreateSkillGapPathStep :one
INSERT INTO hskip_users.bridgr_skill_gap_path_steps (
    path_uuid,
    step_index,
    title,
    rationale,
    estimated_hours,
    resource_uri,
    resource_kind,
    founder_learning_item_uuid,
    course_lesson_uuid,
    linked_node_keys,
    metadata
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: BulkCreateSkillGapPathSteps :copyfrom
INSERT INTO hskip_users.bridgr_skill_gap_path_steps (
    path_uuid,
    step_index,
    title,
    rationale,
    estimated_hours,
    resource_uri,
    resource_kind,
    founder_learning_item_uuid,
    course_lesson_uuid,
    linked_node_keys,
    metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
);

-- name: GetSkillGapPathStep :one
SELECT * FROM hskip_users.bridgr_skill_gap_path_steps
WHERE id = $1;

-- name: ListSkillGapPathStepsByPath :many
SELECT * FROM hskip_users.bridgr_skill_gap_path_steps
WHERE path_uuid = $1
ORDER BY step_index ASC;

-- name: UpdateSkillGapPathStep :one
UPDATE hskip_users.bridgr_skill_gap_path_steps
SET resource_uri = $2,
    resource_kind = $3
WHERE id = $1
RETURNING *;

-- name: DeleteSkillGapPathStepsByPath :exec
DELETE FROM hskip_users.bridgr_skill_gap_path_steps
WHERE path_uuid = $1;
