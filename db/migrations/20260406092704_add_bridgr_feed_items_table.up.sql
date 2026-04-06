-- ═══════════════════════════════════════════════════════════════════════════
-- BRIDGR — FEED ITEMS
-- Curated, denormalized rows surfaced to the user (newsletter / feed UI).
-- user_id app-enforced FK to users.id; job_candidate_uuid app-enforced to job_candidates.uuid;
-- score_uuid app-enforced to job_scores.uuid; verification_uuid reserved for future verification rows.
-- No foreign keys (app-enforced).
-- ═══════════════════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS bridgr.feed_items;

CREATE TABLE bridgr.feed_items (
    uuid                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id                   BIGSERIAL NOT NULL UNIQUE,
    user_id              INTEGER NOT NULL,
    job_candidate_uuid   UUID NOT NULL,
    score_uuid           UUID,
    verification_uuid    UUID,
    composite_score      REAL NOT NULL DEFAULT 0,
    gap_severity         VARCHAR(32),
    title                TEXT,
    company              TEXT,
    location             TEXT,
    job_url              TEXT,
    match_summary        TEXT,
    gap_summary          TEXT,
    feed_status          VARCHAR(30) NOT NULL DEFAULT 'new',
    surfaced_at          TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    seen_at              TIMESTAMP WITHOUT TIME ZONE,
    created_at           TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at           TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,

    CONSTRAINT uq_feed_items_candidate_user UNIQUE (job_candidate_uuid, user_id)
);

COMMENT ON TABLE bridgr.feed_items IS 'Bridgr feed: user-facing rows for high-quality job matches (denormalized for fast reads).';
COMMENT ON COLUMN bridgr.feed_items.user_id IS 'App-enforced FK to users.id';
COMMENT ON COLUMN bridgr.feed_items.job_candidate_uuid IS 'App-enforced FK to job_candidates.uuid';
COMMENT ON COLUMN bridgr.feed_items.score_uuid IS 'App-enforced FK to job_scores.uuid';
COMMENT ON COLUMN bridgr.feed_items.verification_uuid IS 'Optional future link to a verification record; app-enforced when set';
COMMENT ON COLUMN bridgr.feed_items.match_summary IS 'Short human-readable why this job matches';
COMMENT ON COLUMN bridgr.feed_items.gap_summary IS 'Short human-readable skill-gap note';
COMMENT ON COLUMN bridgr.feed_items.feed_status IS 'new, seen, saved, dismissed, applied';
COMMENT ON COLUMN bridgr.feed_items.surfaced_at IS 'When this item entered the feed';

CREATE INDEX idx_feed_items_user_status_surfaced ON bridgr.feed_items(user_id, feed_status, surfaced_at DESC);
CREATE INDEX idx_feed_items_user_score ON bridgr.feed_items(user_id, composite_score DESC);
CREATE INDEX idx_feed_items_candidate ON bridgr.feed_items(job_candidate_uuid);

CREATE TRIGGER tr_feed_items_control_time
    BEFORE INSERT OR UPDATE ON bridgr.feed_items
    FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();
