-- ═══════════════════════════════════════════════════════════════════════════
-- RADAR / BRIDGR — JOB SEARCH PROFILES
-- One row = one canonical Radar search slice (single query, location, board).
-- user_id app-enforced FK to users.id.
-- No foreign keys (app-enforced).
-- ═══════════════════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS bridgr.job_search_profiles;

CREATE TABLE bridgr.job_search_profiles (
    uuid                       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id                         BIGSERIAL NOT NULL UNIQUE,
    user_id                    INTEGER NOT NULL,
    target_role                TEXT    NOT NULL DEFAULT '',
    location                   TEXT    NOT NULL DEFAULT '',
    source_board               TEXT    NOT NULL DEFAULT 'indeed',
    career_switch              BOOLEAN NOT NULL DEFAULT false,
    company_stage              TEXT    NOT NULL DEFAULT 'any',
    seniority_goal             TEXT    NOT NULL DEFAULT 'any',
    compensation_goal          TEXT    NOT NULL DEFAULT 'any',
    software_stack_must_have   TEXT[]  NOT NULL DEFAULT '{}',
    canonical_cv_analysis_uuid UUID,
    created_at                 TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at                 TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL
);

COMMENT ON TABLE bridgr.job_search_profiles IS 'Bridgr job discovery: one row = one Radar search slice (query, location, board, goals).';
COMMENT ON COLUMN bridgr.job_search_profiles.user_id IS 'App-enforced FK to users.id';
COMMENT ON COLUMN bridgr.job_search_profiles.target_role IS 'Single search string passed to Radar (e.g. primary target role or fixed query).';
COMMENT ON COLUMN bridgr.job_search_profiles.location IS 'Single location string for this slice (city, remote label, etc.).';
COMMENT ON COLUMN bridgr.job_search_profiles.source_board IS 'Single board id for this slice (e.g. linkedin, indeed).';
COMMENT ON COLUMN bridgr.job_search_profiles.career_switch IS 'True when the user is seeking a career change / new domain.';
COMMENT ON COLUMN bridgr.job_search_profiles.company_stage IS 'Employer stage filter; product-defined values (e.g. startup, growth, enterprise, any).';
COMMENT ON COLUMN bridgr.job_search_profiles.seniority_goal IS 'Allowed: up, lateral, down, any.';
COMMENT ON COLUMN bridgr.job_search_profiles.compensation_goal IS 'Allowed: up, lateral, down, any.';
COMMENT ON COLUMN bridgr.job_search_profiles.software_stack_must_have IS 'Tags for required technologies in this search slice.';
COMMENT ON COLUMN bridgr.job_search_profiles.canonical_cv_analysis_uuid IS 'Optional app-enforced FK to skill_gap_analyses.uuid for CV fingerprint context';

CREATE TRIGGER tr_job_search_profiles_control_time
    BEFORE INSERT OR UPDATE ON bridgr.job_search_profiles
    FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();
