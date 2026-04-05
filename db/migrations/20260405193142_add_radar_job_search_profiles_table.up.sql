-- ═══════════════════════════════════════════════════════════════════════════
-- RADAR / BRIDGR — JOB SEARCH PROFILES
-- Per-user search preferences (targets, locations, boards, matching policy).
-- user_id app-enforced FK to users.id.
-- No foreign keys (app-enforced).
-- ═══════════════════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS bridgr.job_search_profiles;

CREATE TABLE bridgr.job_search_profiles (
    uuid                       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id                         BIGSERIAL NOT NULL UNIQUE,
    user_id                    INTEGER NOT NULL,
    target_roles               JSONB   NOT NULL DEFAULT '[]',
    locations                  JSONB   NOT NULL DEFAULT '[]',
    boards_enabled             JSONB   NOT NULL DEFAULT '[]',
    matching                   JSONB   NOT NULL DEFAULT '{}',
    canonical_cv_analysis_uuid UUID,
    max_surfaced_jobs          INTEGER NOT NULL DEFAULT 3,
    created_at                 TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at                 TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,

    CONSTRAINT uq_job_search_profiles_user UNIQUE (user_id),
    CONSTRAINT chk_job_search_profiles_max_surfaced CHECK (max_surfaced_jobs > 0 AND max_surfaced_jobs <= 20)
);

COMMENT ON TABLE bridgr.job_search_profiles IS 'Bridgr job discovery: user preferences for Radar search and surfacing.';
COMMENT ON COLUMN bridgr.job_search_profiles.user_id IS 'App-enforced FK to users.id';
COMMENT ON COLUMN bridgr.job_search_profiles.target_roles IS 'JSON array of role strings or objects';
COMMENT ON COLUMN bridgr.job_search_profiles.locations IS 'JSON array of location descriptors (city, remote flags, etc.)';
COMMENT ON COLUMN bridgr.job_search_profiles.boards_enabled IS 'JSON array of board ids (e.g. linkedin, indeed)';
COMMENT ON COLUMN bridgr.job_search_profiles.matching IS 'Policy blob: strict_no_skill_gaps, tiers, etc.';
COMMENT ON COLUMN bridgr.job_search_profiles.canonical_cv_analysis_uuid IS 'Optional app-enforced FK to skill_gap_analyses.uuid for CV fingerprint context';

CREATE TRIGGER tr_job_search_profiles_control_time
    BEFORE INSERT OR UPDATE ON bridgr.job_search_profiles
    FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();
