-- =============================================================================
-- BRIDGR — SUPPORTED JOB BOARDS
-- Canonical Radar board catalog. Unique board_id.
-- =============================================================================

-- name: CreateSupportedJobBoard :one
INSERT INTO bridgr.supported_job_boards (
    board_id,
    display_name,
    engine,
    site_type,
    region,
    is_active,
    config,
    sort_order
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetSupportedJobBoardByUUID :one
SELECT * FROM bridgr.supported_job_boards
WHERE uuid = $1;

-- name: GetSupportedJobBoardByID :one
SELECT * FROM bridgr.supported_job_boards
WHERE id = $1;

-- name: GetSupportedJobBoardByBoardID :one
SELECT * FROM bridgr.supported_job_boards
WHERE board_id = $1;

-- name: ListSupportedJobBoards :many
SELECT * FROM bridgr.supported_job_boards
ORDER BY sort_order ASC, display_name ASC;

-- name: ListActiveSupportedJobBoards :many
SELECT * FROM bridgr.supported_job_boards
WHERE is_active = true
ORDER BY sort_order ASC, display_name ASC;

-- name: UpdateSupportedJobBoardByID :one
UPDATE bridgr.supported_job_boards
SET
    board_id = $2,
    display_name = $3,
    engine = $4,
    site_type = $5,
    region = $6,
    is_active = $7,
    config = $8,
    sort_order = $9
WHERE id = $1
RETURNING *;

-- name: DeleteSupportedJobBoardByUUID :exec
DELETE FROM bridgr.supported_job_boards
WHERE uuid = $1;

-- name: DeleteSupportedJobBoardByID :exec
DELETE FROM bridgr.supported_job_boards
WHERE id = $1;
