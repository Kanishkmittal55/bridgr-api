-- =============================================================================
-- BRIDGR JOB DISCOVERY — JOB SEARCH PROFILES
-- Per-user Radar/search preferences. user_id app-enforced FK to users.id.
-- =============================================================================

-- name: CreateJobSearchProfile :one
INSERT INTO bridgr.job_search_profiles (
    user_id,
    target_roles,
    locations,
    boards_enabled,
    matching,
    canonical_cv_analysis_uuid,
    max_surfaced_jobs
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetJobSearchProfileByUserID :one
SELECT * FROM bridgr.job_search_profiles
WHERE user_id = $1;

-- name: GetJobSearchProfileByUUID :one
SELECT * FROM bridgr.job_search_profiles
WHERE uuid = $1;

-- name: GetJobSearchProfileByID :one
SELECT * FROM bridgr.job_search_profiles
WHERE id = $1;

-- name: UpdateJobSearchProfile :one
UPDATE bridgr.job_search_profiles
SET
    target_roles = $2,
    locations = $3,
    boards_enabled = $4,
    matching = $5,
    canonical_cv_analysis_uuid = $6,
    max_surfaced_jobs = $7
WHERE id = $1
RETURNING *;

-- name: UpdateJobSearchProfileByUserID :one
UPDATE bridgr.job_search_profiles
SET
    target_roles = $2,
    locations = $3,
    boards_enabled = $4,
    matching = $5,
    canonical_cv_analysis_uuid = $6,
    max_surfaced_jobs = $7
WHERE user_id = $1
RETURNING *;

-- name: DeleteJobSearchProfileByUserID :exec
DELETE FROM bridgr.job_search_profiles
WHERE user_id = $1;

-- name: DeleteJobSearchProfileByID :exec
DELETE FROM bridgr.job_search_profiles
WHERE id = $1;
