-- =============================================================================
-- BRIDGR SKILL GAP — ANALYSES
-- One CV+JD skill-gap analysis run. user_id / persona / pursuit are app-enforced FKs.
-- =============================================================================

-- name: CreateSkillGapAnalysis :one
INSERT INTO bridgr.skill_gap_analyses (
    user_id,
    founder_persona_uuid,
    pursuit_uuid,
    title,
    status,
    cv_asset_uri,
    jd_asset_uri,
    cv_fingerprint,
    jd_fingerprint,
    llm_model,
    prompt_version,
    extraction_payload,
    gap_summary,
    mermaid_diagram,
    error_code,
    error_detail
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
)
RETURNING *;

-- name: GetSkillGapAnalysis :one
SELECT * FROM bridgr.skill_gap_analyses
WHERE id = $1;

-- name: GetSkillGapAnalysisByUUID :one
SELECT * FROM bridgr.skill_gap_analyses
WHERE uuid = $1;

-- name: GetSkillGapAnalysisByUser :many
SELECT * FROM bridgr.skill_gap_analyses
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: GetSkillGapAnalysisByFingerprint :one
SELECT * FROM bridgr.skill_gap_analyses
WHERE user_id = $1
  AND cv_fingerprint = $2
  AND jd_fingerprint = $3
ORDER BY created_at DESC
LIMIT 1;

-- name: UpdateSkillGapAnalysisStatus :one
UPDATE bridgr.skill_gap_analyses
SET status = $2
WHERE id = $1
RETURNING *;

-- name: UpdateSkillGapAnalysisSummary :one
UPDATE bridgr.skill_gap_analyses
SET gap_summary = $2,
    mermaid_diagram = $3
WHERE id = $1
RETURNING *;

-- name: UpdateSkillGapAnalysisError :one
UPDATE bridgr.skill_gap_analyses
SET error_code = $2,
    error_detail = $3,
    status = 'failed'
WHERE id = $1
RETURNING *;

-- name: DeleteSkillGapAnalysis :exec
DELETE FROM bridgr.skill_gap_analyses
WHERE id = $1;

-- name: DeleteSkillGapAnalysisByUUID :exec
DELETE FROM bridgr.skill_gap_analyses
WHERE uuid = $1;

-- name: ListSkillGapAnalysesByUser :many
SELECT * FROM bridgr.skill_gap_analyses
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
