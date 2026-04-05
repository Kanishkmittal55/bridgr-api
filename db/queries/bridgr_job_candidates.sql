-- =============================================================================
-- BRIDGR JOB DISCOVERY — JOB CANDIDATES
-- Deduped postings per user. (user_id, url_hash) unique.
-- =============================================================================

-- name: InsertJobCandidate :one
INSERT INTO bridgr.job_candidates (
    user_id,
    discovery_run_uuid,
    source_board,
    source_job_id,
    job_url,
    url_hash,
    content_hash,
    title,
    company,
    location,
    jd_text,
    jd_s3_uri,
    fetched_at,
    ingestion_status,
    radar_payload,
    application_url
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
)
RETURNING *;

-- name: UpsertJobCandidate :one
INSERT INTO bridgr.job_candidates (
    user_id,
    discovery_run_uuid,
    source_board,
    source_job_id,
    job_url,
    url_hash,
    content_hash,
    title,
    company,
    location,
    jd_text,
    jd_s3_uri,
    fetched_at,
    ingestion_status,
    radar_payload,
    application_url
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
)
ON CONFLICT (user_id, url_hash) DO UPDATE SET
    discovery_run_uuid = COALESCE(EXCLUDED.discovery_run_uuid, bridgr.job_candidates.discovery_run_uuid),
    source_board = EXCLUDED.source_board,
    source_job_id = EXCLUDED.source_job_id,
    job_url = EXCLUDED.job_url,
    content_hash = COALESCE(EXCLUDED.content_hash, bridgr.job_candidates.content_hash),
    title = COALESCE(EXCLUDED.title, bridgr.job_candidates.title),
    company = COALESCE(EXCLUDED.company, bridgr.job_candidates.company),
    location = COALESCE(EXCLUDED.location, bridgr.job_candidates.location),
    jd_text = COALESCE(EXCLUDED.jd_text, bridgr.job_candidates.jd_text),
    jd_s3_uri = COALESCE(EXCLUDED.jd_s3_uri, bridgr.job_candidates.jd_s3_uri),
    fetched_at = COALESCE(EXCLUDED.fetched_at, bridgr.job_candidates.fetched_at),
    ingestion_status = EXCLUDED.ingestion_status,
    radar_payload = EXCLUDED.radar_payload,
    application_url = COALESCE(EXCLUDED.application_url, bridgr.job_candidates.application_url)
RETURNING *;

-- name: GetJobCandidateByUUID :one
SELECT * FROM bridgr.job_candidates
WHERE uuid = $1;

-- name: GetJobCandidateByID :one
SELECT * FROM bridgr.job_candidates
WHERE id = $1;

-- name: GetJobCandidateByUserAndURLHash :one
SELECT * FROM bridgr.job_candidates
WHERE user_id = $1 AND url_hash = $2;

-- name: ListJobCandidatesByUser :many
SELECT * FROM bridgr.job_candidates
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListJobCandidatesByUserFiltered :many
SELECT * FROM bridgr.job_candidates AS c
WHERE c.user_id = sqlc.arg('user_id')
  AND (sqlc.narg('discovery_run_uuid')::uuid IS NULL OR c.discovery_run_uuid = sqlc.narg('discovery_run_uuid')::uuid)
  AND (sqlc.narg('ingestion_status')::text IS NULL OR sqlc.narg('ingestion_status')::text = '' OR c.ingestion_status = sqlc.narg('ingestion_status')::text)
  AND (sqlc.narg('source_board')::text IS NULL OR sqlc.narg('source_board')::text = '' OR c.source_board = sqlc.narg('source_board')::text)
ORDER BY c.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ListJobCandidatesByDiscoveryRunUUID :many
SELECT * FROM bridgr.job_candidates
WHERE discovery_run_uuid = $1
ORDER BY created_at ASC;

-- name: UpdateJobCandidateIngestion :one
UPDATE bridgr.job_candidates
SET
    jd_text = $2,
    jd_s3_uri = $3,
    content_hash = $4,
    fetched_at = $5,
    ingestion_status = $6,
    radar_payload = $7
WHERE id = $1
RETURNING *;

-- name: UpdateJobCandidateByID :one
UPDATE bridgr.job_candidates
SET
    discovery_run_uuid = $2,
    source_board = $3,
    source_job_id = $4,
    job_url = $5,
    url_hash = $6,
    content_hash = $7,
    title = $8,
    company = $9,
    location = $10,
    jd_text = $11,
    jd_s3_uri = $12,
    fetched_at = $13,
    ingestion_status = $14,
    radar_payload = $15,
    application_url = $16
WHERE id = $1
RETURNING *;

-- name: DeleteJobCandidateByUUID :exec
DELETE FROM bridgr.job_candidates
WHERE uuid = $1;

-- name: DeleteJobCandidateByID :exec
DELETE FROM bridgr.job_candidates
WHERE id = $1;
