-- ═══════════════════════════════════════════════════════════════════════════
-- BRIDGR — SUPPORTED JOB BOARDS
-- Curated catalog of job boards Radar can target (replaces supported_boards.yaml).
-- board_id is the stable slug used as source_site when crawling and in job_search_profiles.boards_enabled.
-- No foreign keys (app-enforced references from user prefs and discovery runs).
-- ═══════════════════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS bridgr.supported_job_boards;

CREATE TABLE bridgr.supported_job_boards (
    uuid          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id            BIGSERIAL NOT NULL UNIQUE,
    board_id      VARCHAR(64) NOT NULL,
    display_name  VARCHAR(128) NOT NULL,
    engine        VARCHAR(64) NOT NULL,
    site_type     VARCHAR(32) NOT NULL,
    region        VARCHAR(32) NOT NULL DEFAULT 'global',
    is_active     BOOLEAN NOT NULL DEFAULT true,
    config        JSONB NOT NULL DEFAULT '{}',
    sort_order    INTEGER NOT NULL DEFAULT 0,
    created_at    TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at    TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,

    CONSTRAINT uq_supported_job_boards_board_id UNIQUE (board_id)
);

COMMENT ON TABLE bridgr.supported_job_boards IS 'Bridgr/Radar: canonical list of supported job boards (CRUD via API; seeds from former supported_boards.yaml).';
COMMENT ON COLUMN bridgr.supported_job_boards.board_id IS 'Stable slug (e.g. linkedin, indeed); matches Discovery source_site / JobSpy site key';
COMMENT ON COLUMN bridgr.supported_job_boards.display_name IS 'Human-readable label for UI';
COMMENT ON COLUMN bridgr.supported_job_boards.engine IS 'Radar engine: jobspy, smartextract, crawl4ai, workday, indeed, etc.';
COMMENT ON COLUMN bridgr.supported_job_boards.site_type IS 'How the board is crawled: search (query+location URL), static (fixed career pages), etc.';
COMMENT ON COLUMN bridgr.supported_job_boards.region IS 'Primary region: global, us, uk, eu, …';
COMMENT ON COLUMN bridgr.supported_job_boards.is_active IS 'When false, board is hidden from selection and not used for new crawls';
COMMENT ON COLUMN bridgr.supported_job_boards.config IS 'Engine-specific options (e.g. fetch_description, proxy hints); extensible JSON';
COMMENT ON COLUMN bridgr.supported_job_boards.sort_order IS 'Lower values list first in admin and user pickers';

CREATE INDEX idx_supported_job_boards_active_sort ON bridgr.supported_job_boards(is_active, sort_order, display_name);
CREATE INDEX idx_supported_job_boards_engine ON bridgr.supported_job_boards(engine)
    WHERE is_active = true;

CREATE TRIGGER tr_supported_job_boards_control_time
    BEFORE INSERT OR UPDATE ON bridgr.supported_job_boards
    FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();