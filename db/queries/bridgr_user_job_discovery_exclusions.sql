-- =============================================================================
-- BRIDGR — USER JOB DISCOVERY EXCLUSIONS
-- Per-user url_hash registry for discovery batch skip. Unique (user_id, url_hash).
-- =============================================================================

-- name: ListUserJobDiscoveryExclusionsByURLHashes :many
SELECT url_hash
FROM bridgr.user_job_discovery_exclusions
WHERE user_id = sqlc.arg('user_id')
  AND url_hash = ANY(sqlc.arg('url_hashes')::text[]);

-- name: AddUserJobDiscoveryExclusion :exec
INSERT INTO bridgr.user_job_discovery_exclusions (
    user_id,
    url_hash,
    source_board,
    reason,
    discovery_run_uuid
) VALUES (
    sqlc.arg('user_id'),
    sqlc.arg('url_hash'),
    sqlc.arg('source_board'),
    sqlc.arg('reason'),
    sqlc.arg('discovery_run_uuid')
)
ON CONFLICT (user_id, url_hash) DO NOTHING;

-- name: GetUserJobDiscoveryExclusionByUserAndURLHash :one
SELECT *
FROM bridgr.user_job_discovery_exclusions
WHERE user_id = sqlc.arg('user_id') AND url_hash = sqlc.arg('url_hash');

-- name: DeleteUserJobDiscoveryExclusionByUserAndURLHash :exec
DELETE FROM bridgr.user_job_discovery_exclusions
WHERE user_id = sqlc.arg('user_id') AND url_hash = sqlc.arg('url_hash');
