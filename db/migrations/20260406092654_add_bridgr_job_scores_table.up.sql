-- ═══════════════════════════════════════════════════════════════════════════
-- BRIDGR — JOB SCORES
-- Relevance ranking and skill-gap summary for a (user, job_candidate) pair.
-- user_id app-enforced FK to users.id; job_candidate_uuid app-enforced to job_candidates.uuid;
-- enrichment_uuid app-enforced to job_enrichments.uuid (optional).
-- No foreign keys (app-enforced).
-- ═══════════════════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS bridgr.job_scores;

CREATE TABLE bridgr.job_scores (
    uuid                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id                       BIGSERIAL NOT NULL UNIQUE,
    user_id                  INTEGER NOT NULL,
    job_candidate_uuid       UUID NOT NULL,
    enrichment_uuid          UUID,
    skill_match_score        REAL NOT NULL DEFAULT 0,
    experience_match_score   REAL NOT NULL DEFAULT 0,
    location_match_score     REAL NOT NULL DEFAULT 0,
    recency_score            REAL NOT NULL DEFAULT 0,
    board_quality_score      REAL NOT NULL DEFAULT 0,
    composite_score          REAL NOT NULL DEFAULT 0,
    matched_skills           JSONB,
    gap_skills               JSONB,
    gap_severity             VARCHAR(32),
    scoring_model            VARCHAR(128),
    scoring_version          VARCHAR(64),
    created_at               TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at               TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,

    CONSTRAINT uq_job_scores_candidate_user UNIQUE (job_candidate_uuid, user_id)
);

COMMENT ON TABLE bridgr.job_scores IS 'Bridgr feed: per-user relevance and skill-gap scores for a discovered job.';
COMMENT ON COLUMN bridgr.job_scores.user_id IS 'App-enforced FK to users.id';
COMMENT ON COLUMN bridgr.job_scores.job_candidate_uuid IS 'App-enforced FK to job_candidates.uuid';
COMMENT ON COLUMN bridgr.job_scores.enrichment_uuid IS 'Optional app-enforced FK to job_enrichments.uuid used for this score';
COMMENT ON COLUMN bridgr.job_scores.composite_score IS 'Weighted blend of dimension scores; primary feed sort key';
COMMENT ON COLUMN bridgr.job_scores.matched_skills IS 'JSON: skills the user satisfies for this role';
COMMENT ON COLUMN bridgr.job_scores.gap_skills IS 'JSON: required skills still missing for the user';
COMMENT ON COLUMN bridgr.job_scores.gap_severity IS 'none, minor, moderate, major';

CREATE INDEX idx_job_scores_user_composite ON bridgr.job_scores(user_id, composite_score DESC);
CREATE INDEX idx_job_scores_candidate ON bridgr.job_scores(job_candidate_uuid);

CREATE TRIGGER tr_job_scores_control_time
    BEFORE INSERT OR UPDATE ON bridgr.job_scores
    FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();
