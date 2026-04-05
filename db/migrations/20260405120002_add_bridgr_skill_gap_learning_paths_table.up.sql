-- ═══════════════════════════════════════════════════════════════════════════
-- BRIDGR SKILL GAP NAVIGATOR — LEARNING PATHS
-- Generated learning path for an analysis (ordered steps to close gaps).
-- analysis_uuid app-enforced FK to bridgr_skill_gap_analyses.uuid.
-- No foreign keys (app-enforced).
-- ═══════════════════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS hskip_users.bridgr_skill_gap_learning_paths;

CREATE TABLE hskip_users.bridgr_skill_gap_learning_paths (
    uuid          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id            BIGSERIAL NOT NULL UNIQUE,
    analysis_uuid UUID NOT NULL,
    path_version  INTEGER NOT NULL DEFAULT 1,
    algorithm     VARCHAR(64),
    title         TEXT,
    path_metadata JSONB DEFAULT '{}',
    created_at    TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at    TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    UNIQUE(analysis_uuid, path_version)
);

COMMENT ON TABLE hskip_users.bridgr_skill_gap_learning_paths IS 'Learning plan / path for closing skill gaps from one analysis.';
COMMENT ON COLUMN hskip_users.bridgr_skill_gap_learning_paths.analysis_uuid IS 'App-enforced FK to bridgr_skill_gap_analyses.uuid';
COMMENT ON COLUMN hskip_users.bridgr_skill_gap_learning_paths.path_version IS 'Increment when regenerating path for same analysis';

CREATE INDEX idx_bridgr_skill_gap_learning_paths_analysis_uuid ON hskip_users.bridgr_skill_gap_learning_paths(analysis_uuid);

CREATE TRIGGER tr_bridgr_skill_gap_learning_paths_control_time
    BEFORE INSERT OR UPDATE ON hskip_users.bridgr_skill_gap_learning_paths
    FOR EACH ROW EXECUTE FUNCTION hskip_users.tr_control_time();
