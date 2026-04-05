-- =============================================================================
-- BRIDGR JOB DISCOVERY — NOTIFICATIONS
-- Outbox for in-app / email / push. user_id and job_candidate_uuid app-enforced.
-- =============================================================================

-- name: CreateJobNotification :one
INSERT INTO bridgr.job_notifications (
    user_id,
    job_candidate_uuid,
    channel,
    status,
    payload,
    sent_at,
    seen_at,
    error_detail
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetJobNotificationByUUID :one
SELECT * FROM bridgr.job_notifications
WHERE uuid = $1;

-- name: GetJobNotificationByID :one
SELECT * FROM bridgr.job_notifications
WHERE id = $1;

-- name: ListJobNotificationsByUser :many
SELECT * FROM bridgr.job_notifications
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListJobNotificationsByUserFiltered :many
SELECT * FROM bridgr.job_notifications AS n
WHERE n.user_id = sqlc.arg('user_id')
  AND (sqlc.narg('since')::timestamp IS NULL OR n.created_at > sqlc.narg('since')::timestamp)
  AND (sqlc.narg('status')::text IS NULL OR sqlc.narg('status')::text = '' OR n.status = sqlc.narg('status')::text)
ORDER BY n.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ListJobNotificationsByUserAndStatus :many
SELECT * FROM bridgr.job_notifications
WHERE user_id = $1 AND status = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: ListPendingJobNotificationsByUser :many
SELECT * FROM bridgr.job_notifications
WHERE user_id = $1 AND status = 'pending'
ORDER BY created_at ASC
LIMIT $2;

-- name: UpdateJobNotificationStatus :one
UPDATE bridgr.job_notifications
SET
    status = $2,
    sent_at = $3,
    error_detail = $4
WHERE id = $1
RETURNING *;

-- name: UpdateJobNotificationSeen :one
UPDATE bridgr.job_notifications
SET
    status = $2,
    seen_at = $3
WHERE id = $1
RETURNING *;

-- name: DeleteJobNotificationByID :exec
DELETE FROM bridgr.job_notifications
WHERE id = $1;

-- name: DeleteJobNotificationByUUID :exec
DELETE FROM bridgr.job_notifications
WHERE uuid = $1;
