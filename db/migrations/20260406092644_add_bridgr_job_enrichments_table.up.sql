-- ═══════════════════════════════════════════════════════════════════════════
-- BRIDGR — JOB ENRICHMENTS
-- Structured extraction / deep-dive results for a job candidate (JD parsing, skills, etc.).
-- user_id app-enforced FK to users.id; job_candidate_uuid app-enforced to job_candidates.uuid.
-- No foreign keys (app-enforced).
-- ═══════════════════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS bridgr.job_enrichments;

CREATE TABLE bridgr.job_enrichments (
    uuid                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id                     BIGSERIAL NOT NULL UNIQUE,
    user_id                INTEGER NOT NULL,
    job_candidate_uuid     UUID NOT NULL,
    status                 VARCHAR(30) NOT NULL DEFAULT 'pending',
    required_skills        JSONB,
    experience_range       JSONB,
    salary_range           JSONB,
    remote_policy          VARCHAR(64),
    visa_sponsorship       BOOLEAN,
    company_size           VARCHAR(64),
    industry               VARCHAR(128),
    application_deadline   TIMESTAMP WITHOUT TIME ZONE,
    structured_jd          JSONB,
    llm_model              VARCHAR(128),
    prompt_version         VARCHAR(64),
    error_detail           TEXT,
    created_at             TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at             TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    
    CONSTRAINT uq_job_enrichments_candidate_user UNIQUE (job_candidate_uuid, user_id)
);

COMMENT ON TABLE bridgr.job_enrichments IS 'Bridgr feed: LLM- or rules-based structured fields extracted from a job posting for matching.';
COMMENT ON COLUMN bridgr.job_enrichments.user_id IS 'App-enforced FK to users.id';
COMMENT ON COLUMN bridgr.job_enrichments.job_candidate_uuid IS 'App-enforced FK to job_candidates.uuid';
COMMENT ON COLUMN bridgr.job_enrichments.required_skills IS 'JSON array of {skill, level, mandatory} or similar';
COMMENT ON COLUMN bridgr.job_enrichments.experience_range IS 'JSON {min_years, max_years, seniority}';
COMMENT ON COLUMN bridgr.job_enrichments.salary_range IS 'JSON {min, max, currency, period}';
COMMENT ON COLUMN bridgr.job_enrichments.structured_jd IS 'Full normalized job description blob for scoring';

CREATE INDEX idx_job_enrichments_user_status ON bridgr.job_enrichments(user_id, status, updated_at DESC);
CREATE INDEX idx_job_enrichments_candidate ON bridgr.job_enrichments(job_candidate_uuid);

CREATE TRIGGER tr_job_enrichments_control_time
    BEFORE INSERT OR UPDATE ON bridgr.job_enrichments
    FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();
