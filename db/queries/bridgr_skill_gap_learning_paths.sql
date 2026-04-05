-- =============================================================================
-- BRIDGR SKILL GAP — LEARNING PATHS
-- Paths per analysis (v1 often one active path; rows versioned by path_version).
-- =============================================================================

-- name: CreateSkillGapLearningPath :one
INSERT INTO hskip_users.bridgr_skill_gap_learning_paths (
    analysis_uuid,
    path_version,
    algorithm,
    title,
    path_metadata
) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetSkillGapLearningPath :one
SELECT * FROM hskip_users.bridgr_skill_gap_learning_paths
WHERE id = $1;

-- name: GetSkillGapLearningPathByUUID :one
SELECT * FROM hskip_users.bridgr_skill_gap_learning_paths
WHERE uuid = $1;

-- name: GetSkillGapLearningPathByAnalysis :one
SELECT * FROM hskip_users.bridgr_skill_gap_learning_paths
WHERE analysis_uuid = $1
ORDER BY path_version DESC
LIMIT 1;

-- name: UpdateSkillGapLearningPathMetadata :one
UPDATE hskip_users.bridgr_skill_gap_learning_paths
SET path_metadata = $2
WHERE id = $1
RETURNING *;

-- name: DeleteSkillGapLearningPathByAnalysis :exec
DELETE FROM hskip_users.bridgr_skill_gap_learning_paths
WHERE analysis_uuid = $1;
