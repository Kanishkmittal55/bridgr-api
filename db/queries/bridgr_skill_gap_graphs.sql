-- =============================================================================
-- BRIDGR SKILL GAP — GRAPHS
-- Candidate vs role_requirement graphs per analysis. analysis_uuid app-enforced.
-- =============================================================================

-- name: CreateSkillGapGraph :one
INSERT INTO bridgr.skill_gap_graphs (analysis_uuid, kind, metadata)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSkillGapGraph :one
SELECT * FROM bridgr.skill_gap_graphs
WHERE id = $1;

-- name: GetSkillGapGraphByUUID :one
SELECT * FROM bridgr.skill_gap_graphs
WHERE uuid = $1;

-- name: GetSkillGapGraphsByAnalysis :many
SELECT * FROM bridgr.skill_gap_graphs
WHERE analysis_uuid = $1
ORDER BY kind ASC;

-- name: GetSkillGapGraphByKind :one
SELECT * FROM bridgr.skill_gap_graphs
WHERE analysis_uuid = $1 AND kind = $2;

-- name: DeleteSkillGapGraphsByAnalysis :exec
DELETE FROM bridgr.skill_gap_graphs
WHERE analysis_uuid = $1;
