-- ═══════════════════════════════════════════════════════════════════════════
-- BRIDGR SKILL GAP NAVIGATOR — NODES
-- Skill/concept nodes for a graph (candidate CV graph or role requirement graph).
-- graph_uuid app-enforced FK to bridgr_skill_gap_graphs.uuid.
-- No foreign keys (app-enforced).
-- ═══════════════════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS hskip_users.bridgr_skill_gap_nodes;

CREATE TABLE hskip_users.bridgr_skill_gap_nodes (
    uuid              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id                BIGSERIAL NOT NULL UNIQUE,
    graph_uuid        UUID NOT NULL,
    node_key          VARCHAR(256) NOT NULL,
    display_name      TEXT NOT NULL,
    description       TEXT,
    proficiency_hint  TEXT,
    source              VARCHAR(32),
    evidence            JSONB DEFAULT '{}',
    metadata            JSONB DEFAULT '{}',
    position_x          INTEGER DEFAULT 0,
    position_y          INTEGER DEFAULT 0,
    created_at          TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at          TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    UNIQUE(graph_uuid, node_key)
);

COMMENT ON TABLE hskip_users.bridgr_skill_gap_nodes IS 'Skill/concept nodes in a skill-gap graph (candidate or role requirement).';
COMMENT ON COLUMN hskip_users.bridgr_skill_gap_nodes.graph_uuid IS 'App-enforced FK to bridgr_skill_gap_graphs.uuid';
COMMENT ON COLUMN hskip_users.bridgr_skill_gap_nodes.node_key IS 'Stable key for this node within the graph (e.g. normalized skill id)';
COMMENT ON COLUMN hskip_users.bridgr_skill_gap_nodes.display_name IS 'Human-readable label';
COMMENT ON COLUMN hskip_users.bridgr_skill_gap_nodes.source IS 'cv, jd, inferred, merged, etc.';
COMMENT ON COLUMN hskip_users.bridgr_skill_gap_nodes.evidence IS 'Spans, quotes, or LLM citations supporting the node';

CREATE INDEX idx_bridgr_skill_gap_nodes_graph_uuid ON hskip_users.bridgr_skill_gap_nodes(graph_uuid);
CREATE INDEX idx_bridgr_skill_gap_nodes_node_key ON hskip_users.bridgr_skill_gap_nodes(node_key);

CREATE TRIGGER tr_bridgr_skill_gap_nodes_control_time
    BEFORE INSERT OR UPDATE ON hskip_users.bridgr_skill_gap_nodes
    FOR EACH ROW EXECUTE FUNCTION hskip_users.tr_control_time();
