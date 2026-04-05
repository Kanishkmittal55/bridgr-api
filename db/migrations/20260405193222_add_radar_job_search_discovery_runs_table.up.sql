-- ═══════════════════════════════════════════════════════════════════════════
-- RADAR / BRIDGR — JOB SEARCH DISCOVERY RUNS
-- One logical sweep: Radar fetch + ingest metadata for a user.
-- user_id app-enforced FK to users.id.
-- No foreign keys (app-enforced).
-- ═══════════════════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS bridgr.job_search_discovery_runs;

CREATE TABLE bridgr.job_search_discovery_runs (
    uuid                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id                   BIGSERIAL NOT NULL UNIQUE,
    user_id              INTEGER NOT NULL,
    status               VARCHAR(30) NOT NULL DEFAULT 'pending',
    request_params       JSONB NOT NULL DEFAULT '{}',
    radar_meta           JSONB NOT NULL DEFAULT '{}',
    raw_candidate_count  INTEGER NOT NULL DEFAULT 0,
    new_candidate_count  INTEGER NOT NULL DEFAULT 0,
    started_at           TIMESTAMP WITHOUT TIME ZONE,
    completed_at         TIMESTAMP WITHOUT TIME ZONE,
    error_code           VARCHAR(80),
    error_detail         TEXT,
    sqs_message_id       VARCHAR(200),
    created_at           TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at           TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,

    CONSTRAINT chk_job_search_discovery_runs_status
        CHECK (status IN ('pending', 'queued', 'running', 'completed', 'failed', 'cancelled'))
);

COMMENT ON TABLE bridgr.job_search_discovery_runs IS 'Bridgr job discovery: audit trail for one Radar-backed discovery execution.';
COMMENT ON COLUMN bridgr.job_search_discovery_runs.user_id IS 'App-enforced FK to users.id';
COMMENT ON COLUMN bridgr.job_search_discovery_runs.status IS 'pending, queued, running, completed, failed, cancelled';
COMMENT ON COLUMN bridgr.job_search_discovery_runs.request_params IS 'Resolved query/location/boards caps sent to worker/Radar';
COMMENT ON COLUMN bridgr.job_search_discovery_runs.radar_meta IS 'Timings, per-board errors, debug payload';
COMMENT ON COLUMN bridgr.job_search_discovery_runs.sqs_message_id IS 'Optional queue id for async job discovery';

CREATE INDEX idx_job_search_discovery_runs_user_created
    ON bridgr.job_search_discovery_runs(user_id, created_at DESC);
CREATE INDEX idx_job_search_discovery_runs_status
    ON bridgr.job_search_discovery_runs(user_id, status);

CREATE TRIGGER tr_job_search_discovery_runs_control_time
    BEFORE INSERT OR UPDATE ON bridgr.job_search_discovery_runs
    FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();
