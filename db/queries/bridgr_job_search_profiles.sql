-- =============================================================================
-- BRIDGR JOB DISCOVERY — JOB SEARCH PROFILES
-- One row = one canonical Radar search slice. user_id app-enforced FK to users.id.
-- =============================================================================

-- name: CreateJobSearchProfile :one
INSERT INTO bridgr.job_search_profiles (
    user_id,
    target_role,
    location,
    source_board,
    career_switch,
    company_stage,
    seniority_goal,
    compensation_goal,
    software_stack_must_have,
    canonical_cv_analysis_uuid
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: ListJobSearchProfilesByUserID :many
-- $1 = user_id. Returns every saved job search profile row for that user (zero or more rows).
SELECT * FROM bridgr.job_search_profiles
WHERE user_id = $1
ORDER BY id ASC;

-- name: GetJobSearchProfileByUserIDAndUUID :one
SELECT * FROM bridgr.job_search_profiles
WHERE user_id = $1 AND uuid = $2;

-- name: GetJobSearchProfileByUUID :one
SELECT * FROM bridgr.job_search_profiles
WHERE uuid = $1;

-- name: GetJobSearchProfileByID :one
SELECT * FROM bridgr.job_search_profiles
WHERE id = $1;

-- name: UpdateJobSearchProfile :one
UPDATE bridgr.job_search_profiles
SET
    target_role = $2,
    location = $3,
    source_board = $4,
    career_switch = $5,
    company_stage = $6,
    seniority_goal = $7,
    compensation_goal = $8,
    software_stack_must_have = $9,
    canonical_cv_analysis_uuid = $10
WHERE id = $1
RETURNING *;

-- name: UpdateJobSearchProfileByUserID :one
UPDATE bridgr.job_search_profiles
SET
    target_role = $2,
    location = $3,
    source_board = $4,
    career_switch = $5,
    company_stage = $6,
    seniority_goal = $7,
    compensation_goal = $8,
    software_stack_must_have = $9,
    canonical_cv_analysis_uuid = $10
WHERE user_id = $1
RETURNING *;

-- name: UpdateJobSearchProfileByUserIDAndUUID :one
UPDATE bridgr.job_search_profiles
SET
    target_role = $3,
    location = $4,
    source_board = $5,
    career_switch = $6,
    company_stage = $7,
    seniority_goal = $8,
    compensation_goal = $9,
    software_stack_must_have = $10,
    canonical_cv_analysis_uuid = $11
WHERE user_id = $1 AND uuid = $2
RETURNING *;

-- name: DeleteJobSearchProfileByUserID :exec
DELETE FROM bridgr.job_search_profiles
WHERE user_id = $1;

-- name: DeleteJobSearchProfileByID :exec
DELETE FROM bridgr.job_search_profiles
WHERE id = $1;

-- name: DeleteJobSearchProfileByUserIDAndUUID :exec
DELETE FROM bridgr.job_search_profiles
WHERE user_id = $1 AND uuid = $2;
