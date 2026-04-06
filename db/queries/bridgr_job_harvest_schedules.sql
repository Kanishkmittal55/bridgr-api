-- =============================================================================
-- BRIDGR — JOB HARVEST SCHEDULES
-- Per-user cadence + board rotation for a job_search_profile.
-- =============================================================================

-- name: CreateJobHarvestSchedule :one
INSERT INTO bridgr.job_harvest_schedules (
    user_id,
    profile_uuid,
    enabled,
    cadence_minutes,
    boards_rotation,
    last_run_at,
    next_run_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetJobHarvestScheduleByUUID :one
SELECT * FROM bridgr.job_harvest_schedules
WHERE uuid = $1;

-- name: GetJobHarvestScheduleByID :one
SELECT * FROM bridgr.job_harvest_schedules
WHERE id = $1;

-- name: GetJobHarvestScheduleByUserAndProfileUUID :one
SELECT * FROM bridgr.job_harvest_schedules
WHERE user_id = $1 AND profile_uuid = $2;

-- name: ListJobHarvestSchedulesByUser :many
SELECT * FROM bridgr.job_harvest_schedules
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: ListJobHarvestSchedulesDue :many
-- Scheduler: enabled schedules whose next_run_at is set and not in the future.
SELECT * FROM bridgr.job_harvest_schedules
WHERE enabled = true
  AND next_run_at IS NOT NULL
  AND next_run_at <= now()
ORDER BY next_run_at ASC;

-- name: UpdateJobHarvestScheduleByID :one
UPDATE bridgr.job_harvest_schedules
SET
    enabled = $2,
    cadence_minutes = $3,
    boards_rotation = $4,
    last_run_at = $5,
    next_run_at = $6
WHERE id = $1
RETURNING *;

-- name: DeleteJobHarvestScheduleByUUID :exec
DELETE FROM bridgr.job_harvest_schedules
WHERE uuid = $1;

-- name: DeleteJobHarvestScheduleByID :exec
DELETE FROM bridgr.job_harvest_schedules
WHERE id = $1;

