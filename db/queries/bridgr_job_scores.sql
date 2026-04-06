-- =============================================================================
-- BRIDGR — JOB SCORES
-- Relevance + gap summary per (user, job_candidate). Unique (job_candidate_uuid, user_id).
-- =============================================================================

-- name: CreateJobScore :one
INSERT INTO bridgr.job_scores (
    user_id,
    job_candidate_uuid,
    enrichment_uuid,
    skill_match_score,
    experience_match_score,
    location_match_score,
    recency_score,
    board_quality_score,
    composite_score,
    matched_skills,
    gap_skills,
    gap_severity,
    scoring_model,
    scoring_version
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
)
RETURNING *;

-- name: UpsertJobScore :one
INSERT INTO bridgr.job_scores (
    user_id,
    job_candidate_uuid,
    enrichment_uuid,
    skill_match_score,
    experience_match_score,
    location_match_score,
    recency_score,
    board_quality_score,
    composite_score,
    matched_skills,
    gap_skills,
    gap_severity,
    scoring_model,
    scoring_version
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
)
ON CONFLICT (job_candidate_uuid, user_id) DO UPDATE SET
    enrichment_uuid = EXCLUDED.enrichment_uuid,
    skill_match_score = EXCLUDED.skill_match_score,
    experience_match_score = EXCLUDED.experience_match_score,
    location_match_score = EXCLUDED.location_match_score,
    recency_score = EXCLUDED.recency_score,
    board_quality_score = EXCLUDED.board_quality_score,
    composite_score = EXCLUDED.composite_score,
    matched_skills = EXCLUDED.matched_skills,
    gap_skills = EXCLUDED.gap_skills,
    gap_severity = EXCLUDED.gap_severity,
    scoring_model = EXCLUDED.scoring_model,
    scoring_version = EXCLUDED.scoring_version
RETURNING *;

-- name: GetJobScoreByUUID :one
SELECT * FROM bridgr.job_scores
WHERE uuid = $1;

-- name: GetJobScoreByID :one
SELECT * FROM bridgr.job_scores
WHERE id = $1;

-- name: GetJobScoreByUserAndCandidateUUID :one
SELECT * FROM bridgr.job_scores
WHERE user_id = $1 AND job_candidate_uuid = $2;

-- name: ListJobScoresByUser :many
SELECT * FROM bridgr.job_scores
WHERE user_id = $1
ORDER BY composite_score DESC, updated_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateJobScoreByID :one
UPDATE bridgr.job_scores
SET
    enrichment_uuid = $2,
    skill_match_score = $3,
    experience_match_score = $4,
    location_match_score = $5,
    recency_score = $6,
    board_quality_score = $7,
    composite_score = $8,
    matched_skills = $9,
    gap_skills = $10,
    gap_severity = $11,
    scoring_model = $12,
    scoring_version = $13
WHERE id = $1
RETURNING *;

-- name: DeleteJobScoreByUUID :exec
DELETE FROM bridgr.job_scores
WHERE uuid = $1;

-- name: DeleteJobScoreByID :exec
DELETE FROM bridgr.job_scores
WHERE id = $1;
