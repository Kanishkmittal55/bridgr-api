-- ═══════════════════════════════════════════════════════════════════════════
-- BRIDGR — JOB HARVEST SCHEDULES
-- Drives periodic / continuous job discovery per user profile (cadence, board rotation).
-- user_id app-enforced FK to users.id; profile_uuid app-enforced to job_search_profiles.uuid.
-- No foreign keys (app-enforced).
-- ═══════════════════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS bridgr.job_harvest_schedules;

CREATE TABLE bridgr.job_harvest_schedules (
    uuid               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id                 BIGSERIAL NOT NULL UNIQUE,
    user_id            INTEGER NOT NULL,
    profile_uuid       UUID NOT NULL,
    enabled            BOOLEAN NOT NULL DEFAULT true,
    cadence_minutes    INTEGER NOT NULL DEFAULT 360,
    boards_rotation    JSONB NOT NULL DEFAULT '[]',
    last_run_at        TIMESTAMP WITHOUT TIME ZONE,
    next_run_at        TIMESTAMP WITHOUT TIME ZONE,
    created_at         TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at         TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,

    CONSTRAINT uq_job_harvest_schedules_user_profile UNIQUE (user_id, profile_uuid)
);

COMMENT ON TABLE bridgr.job_harvest_schedules IS 'Bridgr feed: when and how to re-run discovery for a job_search_profile.';
COMMENT ON COLUMN bridgr.job_harvest_schedules.user_id IS 'App-enforced FK to users.id';
COMMENT ON COLUMN bridgr.job_harvest_schedules.profile_uuid IS 'App-enforced FK to job_search_profiles.uuid';
COMMENT ON COLUMN bridgr.job_harvest_schedules.cadence_minutes IS 'Minimum interval between automated harvest runs';
COMMENT ON COLUMN bridgr.job_harvest_schedules.boards_rotation IS 'JSON array of board ids to rotate through on successive runs';
COMMENT ON COLUMN bridgr.job_harvest_schedules.last_run_at IS 'When the last scheduled harvest started or completed';
COMMENT ON COLUMN bridgr.job_harvest_schedules.next_run_at IS 'When the scheduler should enqueue the next harvest';

CREATE INDEX idx_job_harvest_schedules_user ON bridgr.job_harvest_schedules(user_id, enabled);
CREATE INDEX idx_job_harvest_schedules_next_run ON bridgr.job_harvest_schedules(next_run_at)
    WHERE enabled = true;

CREATE TRIGGER tr_job_harvest_schedules_control_time
    BEFORE INSERT OR UPDATE ON bridgr.job_harvest_schedules
    FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();
