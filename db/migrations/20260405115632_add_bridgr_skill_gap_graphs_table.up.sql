-- ═══════════════════════════════════════════════════════════════════════════
-- BRIDGR SKILL GAP GRAPHS — one row per graph (candidate CV vs role JD) per analysis
-- Child of bridgr_skill_gap_analyses. Nodes/edges attach via graph uuid in follow-on tables.
-- No foreign keys (app-enforced).
-- ═══════════════════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS hskip_users.bridgr_skill_gap_graphs;

CREATE TABLE hskip_users.bridgr_skill_gap_graphs (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id BIGSERIAL NOT NULL UNIQUE,

    analysis_uuid UUID NOT NULL,           -- application enforced FK to bridgr_skill_gap_analyses.uuid

    kind VARCHAR(32) NOT NULL,
    -- candidate | role_requirement

    metadata JSONB,

    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,

    CONSTRAINT uq_bridgr_skill_gap_graphs_analysis_kind UNIQUE (analysis_uuid, kind)
);

COMMENT ON TABLE hskip_users.bridgr_skill_gap_graphs IS 'Skill graphs for Bridgr: typically two per analysis (candidate skills vs role requirements).';
COMMENT ON COLUMN hskip_users.bridgr_skill_gap_graphs.analysis_uuid IS 'App-enforced FK to bridgr_skill_gap_analyses.uuid';
COMMENT ON COLUMN hskip_users.bridgr_skill_gap_graphs.kind IS 'candidate = graph from CV; role_requirement = graph from job description';
COMMENT ON COLUMN hskip_users.bridgr_skill_gap_graphs.metadata IS 'Optional graph-level notes, model version, or layout hints';

CREATE INDEX idx_bridgr_skill_gap_graphs_analysis
    ON hskip_users.bridgr_skill_gap_graphs(analysis_uuid);

CREATE TRIGGER tr_bridgr_skill_gap_graphs_control_time
    BEFORE INSERT OR UPDATE ON hskip_users.bridgr_skill_gap_graphs
    FOR EACH ROW EXECUTE FUNCTION hskip_users.tr_control_time();
