-- =============================================================================
-- BRIDGR JOB DISCOVERY — ANALYSIS JOB LINK
-- Junction: skill_gap_analyses.uuid ↔ job_candidates.uuid. App-enforced FKs.
-- =============================================================================

-- name: CreateAnalysisJobLink :one
INSERT INTO bridgr.analysis_job_link (
    user_id,
    analysis_uuid,
    job_candidate_uuid,
    link_kind
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetAnalysisJobLinkByUUID :one
SELECT * FROM bridgr.analysis_job_link
WHERE uuid = $1;

-- name: GetAnalysisJobLinkByID :one
SELECT * FROM bridgr.analysis_job_link
WHERE id = $1;

-- name: GetAnalysisJobLinkByPair :one
SELECT * FROM bridgr.analysis_job_link
WHERE analysis_uuid = $1 AND job_candidate_uuid = $2;

-- name: ListAnalysisJobLinksByUser :many
SELECT * FROM bridgr.analysis_job_link
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListAnalysisJobLinksByAnalysisUUID :many
SELECT * FROM bridgr.analysis_job_link
WHERE analysis_uuid = $1
ORDER BY created_at DESC;

-- name: ListAnalysisJobLinksByCandidateUUID :many
SELECT * FROM bridgr.analysis_job_link
WHERE job_candidate_uuid = $1
ORDER BY created_at DESC;

-- name: UpdateAnalysisJobLinkKind :one
UPDATE bridgr.analysis_job_link
SET link_kind = $2
WHERE id = $1
RETURNING *;

-- name: DeleteAnalysisJobLinkByUUID :exec
DELETE FROM bridgr.analysis_job_link
WHERE uuid = $1;

-- name: DeleteAnalysisJobLinkByID :exec
DELETE FROM bridgr.analysis_job_link
WHERE id = $1;

-- name: DeleteAnalysisJobLinkByPair :exec
DELETE FROM bridgr.analysis_job_link
WHERE analysis_uuid = $1 AND job_candidate_uuid = $2;
