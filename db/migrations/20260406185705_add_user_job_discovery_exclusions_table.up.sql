-- ═══════════════════════════════════════════════════════════════════════════
-- BRIDGR — USER JOB DISCOVERY EXCLUSIONS
-- Slim per-user registry of job URL hashes already evaluated or rejected during
-- discovery so subsequent FindJobs batches can skip them (batch lookup only).
-- user_id app-enforced FK to users.id. No foreign keys on job_candidates.
-- ═══════════════════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS bridgr.user_job_discovery_exclusions;

CREATE TABLE bridgr.user_job_discovery_exclusions (
    uuid                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id                   BIGSERIAL NOT NULL UNIQUE,
    user_id              INTEGER NOT NULL,
    url_hash             VARCHAR(64) NOT NULL,
    source_board         VARCHAR(64),
    reason               VARCHAR(32),
    discovery_run_uuid   UUID,
    created_at           TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at           TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,

    CONSTRAINT uq_user_job_discovery_exclusions_user_url UNIQUE (user_id, url_hash)
);

COMMENT ON TABLE bridgr.user_job_discovery_exclusions IS 'Per-user URL hashes excluded from re-processing in job discovery (prefilter reject, low score, etc.); supports batch EXISTS lookups.';
COMMENT ON COLUMN bridgr.user_job_discovery_exclusions.user_id IS 'App-enforced FK to users.id';
COMMENT ON COLUMN bridgr.user_job_discovery_exclusions.url_hash IS 'Hex SHA-256 of normalized job URL (same canonicalization as job_candidates.url_hash)';
COMMENT ON COLUMN bridgr.user_job_discovery_exclusions.source_board IS 'Optional board slug when exclusion is board-scoped in app logic';
COMMENT ON COLUMN bridgr.user_job_discovery_exclusions.reason IS 'Short code e.g. prefilter, low_score, duplicate';
COMMENT ON COLUMN bridgr.user_job_discovery_exclusions.discovery_run_uuid IS 'Optional app-enforced FK to job_search_discovery_runs.uuid for traceability';

CREATE TRIGGER tr_user_job_discovery_exclusions_control_time
    BEFORE INSERT OR UPDATE ON bridgr.user_job_discovery_exclusions
    FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();
