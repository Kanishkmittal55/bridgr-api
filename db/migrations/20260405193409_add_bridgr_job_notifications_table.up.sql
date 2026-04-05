-- ═══════════════════════════════════════════════════════════════════════════
-- BRIDGR — JOB NOTIFICATIONS
-- Outbox for surfacing alerts (in-app, email, etc.).
-- user_id app-enforced FK to users.id; job_candidate_uuid app-enforced to job_candidates.uuid.
-- No foreign keys (app-enforced).
-- ═══════════════════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS bridgr.job_notifications;

CREATE TABLE bridgr.job_notifications (
    uuid                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id                   BIGSERIAL NOT NULL UNIQUE,
    user_id              INTEGER NOT NULL,
    job_candidate_uuid   UUID NOT NULL,
    channel              VARCHAR(32) NOT NULL DEFAULT 'in_app',
    status               VARCHAR(30) NOT NULL DEFAULT 'pending',
    payload              JSONB NOT NULL DEFAULT '{}',
    sent_at              TIMESTAMP WITHOUT TIME ZONE,
    seen_at              TIMESTAMP WITHOUT TIME ZONE,
    error_detail         TEXT,
    created_at           TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at           TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,

    CONSTRAINT chk_job_notifications_channel
        CHECK (channel IN ('in_app', 'email', 'push')),
    CONSTRAINT chk_job_notifications_status
        CHECK (status IN ('pending', 'sent', 'failed', 'seen', 'skipped'))
);

COMMENT ON TABLE bridgr.job_notifications IS 'Bridgr job discovery: notification rows for new surfaced or relevant jobs.';
COMMENT ON COLUMN bridgr.job_notifications.user_id IS 'App-enforced FK to users.id';
COMMENT ON COLUMN bridgr.job_notifications.job_candidate_uuid IS 'App-enforced FK to job_candidates.uuid';
COMMENT ON COLUMN bridgr.job_notifications.channel IS 'in_app, email, push';
COMMENT ON COLUMN bridgr.job_notifications.status IS 'pending, sent, failed, seen, skipped';
COMMENT ON COLUMN bridgr.job_notifications.payload IS 'Title snippet, deep link hints, template vars';

CREATE INDEX idx_job_notifications_user_status ON bridgr.job_notifications(user_id, status, created_at DESC);
CREATE INDEX idx_job_notifications_candidate ON bridgr.job_notifications(job_candidate_uuid);

CREATE TRIGGER tr_job_notifications_control_time
    BEFORE INSERT OR UPDATE ON bridgr.job_notifications
    FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();
