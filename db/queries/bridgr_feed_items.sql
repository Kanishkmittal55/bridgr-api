-- =============================================================================
-- BRIDGR — FEED ITEMS
-- Denormalized feed rows per (user, job_candidate). Unique (job_candidate_uuid, user_id).
-- =============================================================================

-- name: CreateFeedItem :one
INSERT INTO bridgr.feed_items (
    user_id,
    job_candidate_uuid,
    score_uuid,
    verification_uuid,
    composite_score,
    gap_severity,
    title,
    company,
    location,
    job_url,
    match_summary,
    gap_summary,
    feed_status,
    surfaced_at,
    seen_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
)
RETURNING *;

-- name: UpsertFeedItem :one
INSERT INTO bridgr.feed_items (
    user_id,
    job_candidate_uuid,
    score_uuid,
    verification_uuid,
    composite_score,
    gap_severity,
    title,
    company,
    location,
    job_url,
    match_summary,
    gap_summary,
    feed_status,
    surfaced_at,
    seen_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
)
ON CONFLICT (job_candidate_uuid, user_id) DO UPDATE SET
    score_uuid = EXCLUDED.score_uuid,
    verification_uuid = EXCLUDED.verification_uuid,
    composite_score = EXCLUDED.composite_score,
    gap_severity = EXCLUDED.gap_severity,
    title = EXCLUDED.title,
    company = EXCLUDED.company,
    location = EXCLUDED.location,
    job_url = EXCLUDED.job_url,
    match_summary = EXCLUDED.match_summary,
    gap_summary = EXCLUDED.gap_summary,
    feed_status = EXCLUDED.feed_status,
    surfaced_at = EXCLUDED.surfaced_at,
    seen_at = EXCLUDED.seen_at
RETURNING *;

-- Upsert used by discovery worker: refresh score/summaries on retry but preserve user feed state and first surfaced time.
-- name: UpsertFeedItemFromDiscovery :one
INSERT INTO bridgr.feed_items (
    user_id,
    job_candidate_uuid,
    score_uuid,
    verification_uuid,
    composite_score,
    gap_severity,
    title,
    company,
    location,
    job_url,
    match_summary,
    gap_summary,
    feed_status,
    surfaced_at,
    seen_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
)
ON CONFLICT (job_candidate_uuid, user_id) DO UPDATE SET
    score_uuid = EXCLUDED.score_uuid,
    verification_uuid = COALESCE(bridgr.feed_items.verification_uuid, EXCLUDED.verification_uuid),
    composite_score = EXCLUDED.composite_score,
    gap_severity = EXCLUDED.gap_severity,
    title = EXCLUDED.title,
    company = EXCLUDED.company,
    location = EXCLUDED.location,
    job_url = EXCLUDED.job_url,
    match_summary = EXCLUDED.match_summary,
    gap_summary = EXCLUDED.gap_summary,
    feed_status = bridgr.feed_items.feed_status,
    surfaced_at = bridgr.feed_items.surfaced_at,
    seen_at = bridgr.feed_items.seen_at
RETURNING *;

-- name: GetFeedItemByUUID :one
SELECT * FROM bridgr.feed_items
WHERE uuid = $1;

-- name: GetFeedItemByID :one
SELECT * FROM bridgr.feed_items
WHERE id = $1;

-- name: GetFeedItemByUserAndCandidateUUID :one
SELECT * FROM bridgr.feed_items
WHERE user_id = $1 AND job_candidate_uuid = $2;

-- name: ListFeedItemsByUser :many
SELECT * FROM bridgr.feed_items
WHERE user_id = $1
ORDER BY surfaced_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateFeedItemByID :one
UPDATE bridgr.feed_items
SET
    score_uuid = $2,
    verification_uuid = $3,
    composite_score = $4,
    gap_severity = $5,
    title = $6,
    company = $7,
    location = $8,
    job_url = $9,
    match_summary = $10,
    gap_summary = $11,
    feed_status = $12,
    surfaced_at = $13,
    seen_at = $14
WHERE id = $1
RETURNING *;

-- name: DeleteFeedItemByUUID :exec
DELETE FROM bridgr.feed_items
WHERE uuid = $1;

-- name: DeleteFeedItemByID :exec
DELETE FROM bridgr.feed_items
WHERE id = $1;
