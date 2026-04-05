-- ═══════════════════════════════════════════════════════════════════════════
-- RADAR / BRIDGR — JOB CANDIDATES
-- Deduplicated postings per user from discovery (URL/content hash).
-- user_id app-enforced FK to users.id; discovery_run_uuid app-enforced to job_search_discovery_runs.uuid.
-- No foreign keys (app-enforced).
-- ═══════════════════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS bridgr.job_candidates;

CREATE TABLE bridgr.job_candidates (
    uuid                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id                   BIGSERIAL NOT NULL UNIQUE,
    user_id              INTEGER NOT NULL,
    discovery_run_uuid   UUID,
    source_board         VARCHAR(80) NOT NULL DEFAULT '',
    source_job_id        TEXT,
    job_url              TEXT NOT NULL,
    url_hash             VARCHAR(128) NOT NULL,
    content_hash         VARCHAR(128),
    title                TEXT,
    company              TEXT,
    location             TEXT,
    jd_text              TEXT,
    jd_s3_uri            TEXT,
    fetched_at           TIMESTAMP WITHOUT TIME ZONE,
    ingestion_status     VARCHAR(30) NOT NULL DEFAULT 'pending',
    radar_payload        JSONB NOT NULL DEFAULT '{}',
    application_url      TEXT,
    created_at           TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at           TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,

    CONSTRAINT chk_job_candidates_ingestion_status
        CHECK (ingestion_status IN ('pending', 'fetched', 'fetch_failed', 'partial'))
);

COMMENT ON TABLE bridgr.job_candidates IS 'Bridgr job discovery: one row per unique posting per user (dedupe by url_hash).';
COMMENT ON COLUMN bridgr.job_candidates.user_id IS 'App-enforced FK to users.id';
COMMENT ON COLUMN bridgr.job_candidates.discovery_run_uuid IS 'App-enforced FK to job_search_discovery_runs.uuid';
COMMENT ON COLUMN bridgr.job_candidates.url_hash IS 'Stable hash of normalized URL for deduplication';
COMMENT ON COLUMN bridgr.job_candidates.content_hash IS 'Optional hash of JD body for refresh/dedup';
COMMENT ON COLUMN bridgr.job_candidates.jd_text IS 'Inline job description text when small enough';
COMMENT ON COLUMN bridgr.job_candidates.jd_s3_uri IS 'Pointer to stored JD blob when large';
COMMENT ON COLUMN bridgr.job_candidates.radar_payload IS 'Raw or normalized Radar response slice';
COMMENT ON COLUMN bridgr.job_candidates.application_url IS 'Direct apply URL when distinct from job_url';

CREATE UNIQUE INDEX uq_job_candidates_user_url_hash ON bridgr.job_candidates(user_id, url_hash);
CREATE INDEX idx_job_candidates_user_created ON bridgr.job_candidates(user_id, created_at DESC);
CREATE INDEX idx_job_candidates_discovery_run ON bridgr.job_candidates(discovery_run_uuid);

CREATE TRIGGER tr_job_candidates_control_time
    BEFORE INSERT OR UPDATE ON bridgr.job_candidates
    FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();
