-- ═══════════════════════════════════════════════════════════════════════════
-- BRIDGR SKILL GAP NAVIGATOR — EDGES
-- Directed edges between skill-gap nodes (prerequisite, similarity, etc.).
-- graph_uuid, from_node_uuid, to_node_uuid app-enforced FKs.
-- No foreign keys (app-enforced).
-- ═══════════════════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS hskip_users.bridgr_skill_gap_edges;

CREATE TABLE hskip_users.bridgr_skill_gap_edges (
    uuid              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id                BIGSERIAL NOT NULL UNIQUE,
    graph_uuid        UUID NOT NULL,
    from_node_uuid    UUID NOT NULL,
    to_node_uuid      UUID NOT NULL,
    relation          VARCHAR(64) NOT NULL DEFAULT 'related',
    weight            NUMERIC(10,4) DEFAULT 1.0,
    metadata          JSONB DEFAULT '{}',
    created_at        TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at        TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    CONSTRAINT bridgr_skill_gap_edges_no_self CHECK (from_node_uuid != to_node_uuid),
    UNIQUE(graph_uuid, from_node_uuid, to_node_uuid, relation)
);

COMMENT ON TABLE hskip_users.bridgr_skill_gap_edges IS 'Directed edges between skill-gap nodes in a graph.';
COMMENT ON COLUMN hskip_users.bridgr_skill_gap_edges.graph_uuid IS 'App-enforced FK to bridgr_skill_gap_graphs.uuid';
COMMENT ON COLUMN hskip_users.bridgr_skill_gap_edges.from_node_uuid IS 'App-enforced FK to bridgr_skill_gap_nodes.uuid';
COMMENT ON COLUMN hskip_users.bridgr_skill_gap_edges.to_node_uuid IS 'App-enforced FK to bridgr_skill_gap_nodes.uuid';
COMMENT ON COLUMN hskip_users.bridgr_skill_gap_edges.relation IS 'prerequisite, similarity, required_by, etc.';

CREATE INDEX idx_bridgr_skill_gap_edges_graph_uuid ON hskip_users.bridgr_skill_gap_edges(graph_uuid);
CREATE INDEX idx_bridgr_skill_gap_edges_from ON hskip_users.bridgr_skill_gap_edges(from_node_uuid);
CREATE INDEX idx_bridgr_skill_gap_edges_to ON hskip_users.bridgr_skill_gap_edges(to_node_uuid);
CREATE INDEX idx_bridgr_skill_gap_edges_relation ON hskip_users.bridgr_skill_gap_edges(relation);

CREATE TRIGGER tr_bridgr_skill_gap_edges_control_time
    BEFORE INSERT OR UPDATE ON hskip_users.bridgr_skill_gap_edges
    FOR EACH ROW EXECUTE FUNCTION hskip_users.tr_control_time();
