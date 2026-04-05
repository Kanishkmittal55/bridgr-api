-- ═══════════════════════════════════════════════════════════════════════════
-- BRIDGR — ANALYSIS JOB LINK
-- Optional junction: skill-gap analysis created from or tied to a job candidate.
-- analysis_uuid app-enforced FK to skill_gap_analyses.uuid.
-- job_candidate_uuid app-enforced FK to job_candidates.uuid.
-- No foreign keys (app-enforced).
-- ═══════════════════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS bridgr.analysis_job_link;

CREATE TABLE bridgr.analysis_job_link (
    uuid                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id                   BIGSERIAL NOT NULL UNIQUE,
    user_id              INTEGER NOT NULL,
    analysis_uuid        UUID NOT NULL,
    job_candidate_uuid   UUID NOT NULL,
    link_kind            VARCHAR(32) NOT NULL DEFAULT 'from_job_feed',
    created_at           TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at           TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,

    CONSTRAINT uq_analysis_job_link_pair UNIQUE (analysis_uuid, job_candidate_uuid),
    CONSTRAINT chk_analysis_job_link_kind
        CHECK (link_kind IN ('from_job_feed', 'manual', 'import'))
);

COMMENT ON TABLE bridgr.analysis_job_link IS 'Bridgr: associates a skill_gap_analysis with a discovered job_candidates row.';
COMMENT ON COLUMN bridgr.analysis_job_link.user_id IS 'App-enforced FK to users.id; denormalized for listing';
COMMENT ON COLUMN bridgr.analysis_job_link.analysis_uuid IS 'App-enforced FK to skill_gap_analyses.uuid';
COMMENT ON COLUMN bridgr.analysis_job_link.job_candidate_uuid IS 'App-enforced FK to job_candidates.uuid';
COMMENT ON COLUMN bridgr.analysis_job_link.link_kind IS 'from_job_feed, manual, import';

CREATE INDEX idx_analysis_job_link_user ON bridgr.analysis_job_link(user_id, created_at DESC);
CREATE INDEX idx_analysis_job_link_analysis ON bridgr.analysis_job_link(analysis_uuid);
CREATE INDEX idx_analysis_job_link_candidate ON bridgr.analysis_job_link(job_candidate_uuid);

CREATE TRIGGER tr_analysis_job_link_control_time
    BEFORE INSERT OR UPDATE ON bridgr.analysis_job_link
    FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();
