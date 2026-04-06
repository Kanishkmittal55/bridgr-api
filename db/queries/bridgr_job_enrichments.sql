-- =============================================================================
-- BRIDGR — JOB ENRICHMENTS
-- Structured JD extraction per (user, job_candidate). Unique (job_candidate_uuid, user_id).
-- =============================================================================

-- name: CreateJobEnrichment :one
INSERT INTO bridgr.job_enrichments (
    user_id,
    job_candidate_uuid,
    status,
    required_skills,
    experience_range,
    salary_range,
    remote_policy,
    visa_sponsorship,
    company_size,
    industry,
    application_deadline,
    structured_jd,
    llm_model,
    prompt_version,
    error_detail
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
)
RETURNING *;

-- name: UpsertJobEnrichment :one
INSERT INTO bridgr.job_enrichments (
    user_id,
    job_candidate_uuid,
    status,
    required_skills,
    experience_range,
    salary_range,
    remote_policy,
    visa_sponsorship,
    company_size,
    industry,
    application_deadline,
    structured_jd,
    llm_model,
    prompt_version,
    error_detail
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
)
ON CONFLICT (job_candidate_uuid, user_id) DO UPDATE SET
    status = EXCLUDED.status,
    required_skills = EXCLUDED.required_skills,
    experience_range = EXCLUDED.experience_range,
    salary_range = EXCLUDED.salary_range,
    remote_policy = EXCLUDED.remote_policy,
    visa_sponsorship = EXCLUDED.visa_sponsorship,
    company_size = EXCLUDED.company_size,
    industry = EXCLUDED.industry,
    application_deadline = EXCLUDED.application_deadline,
    structured_jd = EXCLUDED.structured_jd,
    llm_model = EXCLUDED.llm_model,
    prompt_version = EXCLUDED.prompt_version,
    error_detail = EXCLUDED.error_detail
RETURNING *;

-- name: GetJobEnrichmentByUUID :one
SELECT * FROM bridgr.job_enrichments
WHERE uuid = $1;

-- name: GetJobEnrichmentByID :one
SELECT * FROM bridgr.job_enrichments
WHERE id = $1;

-- name: GetJobEnrichmentByUserAndCandidateUUID :one
SELECT * FROM bridgr.job_enrichments
WHERE user_id = $1 AND job_candidate_uuid = $2;

-- name: ListJobEnrichmentsByUser :many
SELECT * FROM bridgr.job_enrichments
WHERE user_id = $1
ORDER BY updated_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateJobEnrichmentByID :one
UPDATE bridgr.job_enrichments
SET
    status = $2,
    required_skills = $3,
    experience_range = $4,
    salary_range = $5,
    remote_policy = $6,
    visa_sponsorship = $7,
    company_size = $8,
    industry = $9,
    application_deadline = $10,
    structured_jd = $11,
    llm_model = $12,
    prompt_version = $13,
    error_detail = $14
WHERE id = $1
RETURNING *;

-- name: DeleteJobEnrichmentByUUID :exec
DELETE FROM bridgr.job_enrichments
WHERE uuid = $1;

-- name: DeleteJobEnrichmentByID :exec
DELETE FROM bridgr.job_enrichments
WHERE id = $1;
