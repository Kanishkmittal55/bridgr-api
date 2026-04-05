-- =============================================================================
-- BRIDGR SKILL GAP — GRAPHS
-- Candidate vs role_requirement graphs per analysis. analysis_uuid app-enforced.
-- =============================================================================

-- name: CreateSkillGapGraph :one
INSERT INTO hskip_users.bridgr_skill_gap_graphs (analysis_uuid, kind, metadata)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSkillGapGraph :one
SELECT * FROM hskip_users.bridgr_skill_gap_graphs
WHERE id = $1;

-- name: GetSkillGapGraphByUUID :one
SELECT * FROM hskip_users.bridgr_skill_gap_graphs
WHERE uuid = $1;

-- name: GetSkillGapGraphsByAnalysis :many
SELECT * FROM hskip_users.bridgr_skill_gap_graphs
WHERE analysis_uuid = $1
ORDER BY kind ASC;

-- name: GetSkillGapGraphByKind :one
SELECT * FROM hskip_users.bridgr_skill_gap_graphs
WHERE analysis_uuid = $1 AND kind = $2;

-- name: DeleteSkillGapGraphsByAnalysis :exec
DELETE FROM hskip_users.bridgr_skill_gap_graphs
WHERE analysis_uuid = $1;
