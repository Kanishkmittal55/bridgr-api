-- =============================================================================
-- BRIDGR JOB DISCOVERY — DISCOVERY RUNS
-- One Radar sweep audit row per execution. user_id app-enforced FK to users.id.
-- =============================================================================

-- name: CreateJobSearchDiscoveryRun :one
INSERT INTO bridgr.job_search_discovery_runs (
    user_id,
    status,
    request_params,
    radar_meta,
    raw_candidate_count,
    new_candidate_count,
    started_at,
    completed_at,
    error_code,
    error_detail,
    sqs_message_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING *;

-- name: GetJobSearchDiscoveryRunByUUID :one
SELECT * FROM bridgr.job_search_discovery_runs
WHERE uuid = $1;

-- name: GetJobSearchDiscoveryRunByID :one
SELECT * FROM bridgr.job_search_discovery_runs
WHERE id = $1;

-- name: ListJobSearchDiscoveryRunsByUser :many
SELECT * FROM bridgr.job_search_discovery_runs
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountJobSearchDiscoveryRunsByUserSince :one
SELECT COUNT(*)::bigint AS count FROM bridgr.job_search_discovery_runs
WHERE user_id = $1 AND created_at >= $2;

-- name: UpdateJobSearchDiscoveryRunStatus :one
UPDATE bridgr.job_search_discovery_runs
SET status = $2
WHERE id = $1
RETURNING *;

-- name: UpdateJobSearchDiscoveryRunProgress :one
UPDATE bridgr.job_search_discovery_runs
SET
    status = $2,
    raw_candidate_count = $3,
    new_candidate_count = $4,
    radar_meta = $5
WHERE id = $1
RETURNING *;

-- name: UpdateJobSearchDiscoveryRunStarted :one
UPDATE bridgr.job_search_discovery_runs
SET
    status = $2,
    started_at = $3,
    sqs_message_id = $4
WHERE id = $1
RETURNING *;

-- name: UpdateJobSearchDiscoveryRunFinished :one
UPDATE bridgr.job_search_discovery_runs
SET
    status = $2,
    completed_at = $3,
    raw_candidate_count = $4,
    new_candidate_count = $5,
    radar_meta = $6,
    error_code = $7,
    error_detail = $8
WHERE id = $1
RETURNING *;

-- name: DeleteJobSearchDiscoveryRunByID :exec
DELETE FROM bridgr.job_search_discovery_runs
WHERE id = $1;
