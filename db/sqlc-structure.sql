-- Bridgr API — sqlc schema mirror (see db/migrations for versioned DDL).
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE SCHEMA IF NOT EXISTS hskip_users;

CREATE OR REPLACE FUNCTION hskip_users.tr_control_time()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        IF NEW.created_at IS NULL THEN
            NEW.created_at := NOW();
        END IF;
        IF NEW.updated_at IS NULL THEN
            NEW.updated_at := NEW.created_at;
        END IF;
    ELSIF TG_OP = 'UPDATE' THEN
        NEW.created_at = OLD.created_at;
        NEW.updated_at = NOW();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE hskip_users.bridgr_skill_gap_analyses (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id BIGSERIAL NOT NULL UNIQUE,
    user_id INTEGER NOT NULL,
    founder_persona_uuid UUID,
    pursuit_uuid UUID,
    title TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    cv_asset_uri TEXT,
    jd_asset_uri TEXT,
    cv_fingerprint VARCHAR(128),
    jd_fingerprint VARCHAR(128),
    llm_model VARCHAR(120),
    prompt_version VARCHAR(80),
    extraction_payload JSONB,
    gap_summary JSONB,
    mermaid_diagram TEXT,
    error_code VARCHAR(80),
    error_detail TEXT,
    sqs_message_id VARCHAR(200),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    CONSTRAINT chk_bridgr_skill_gap_analyses_status
        CHECK (status IN ('pending', 'extracting', 'graphed', 'pathed', 'completed', 'failed'))
);

CREATE INDEX idx_bridgr_skill_gap_analyses_user_created
    ON hskip_users.bridgr_skill_gap_analyses(user_id, created_at DESC);
CREATE INDEX idx_bridgr_skill_gap_analyses_status
    ON hskip_users.bridgr_skill_gap_analyses(status);
CREATE INDEX idx_bridgr_skill_gap_analyses_persona
    ON hskip_users.bridgr_skill_gap_analyses(founder_persona_uuid);
CREATE INDEX idx_bridgr_skill_gap_analyses_pursuit
    ON hskip_users.bridgr_skill_gap_analyses(pursuit_uuid);

CREATE TRIGGER tr_bridgr_skill_gap_analyses_control_time
    BEFORE INSERT OR UPDATE ON hskip_users.bridgr_skill_gap_analyses
    FOR EACH ROW EXECUTE FUNCTION hskip_users.tr_control_time();

CREATE TABLE hskip_users.bridgr_skill_gap_graphs (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id BIGSERIAL NOT NULL UNIQUE,
    analysis_uuid UUID NOT NULL,
    kind VARCHAR(32) NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    CONSTRAINT uq_bridgr_skill_gap_graphs_analysis_kind UNIQUE (analysis_uuid, kind)
);

CREATE INDEX idx_bridgr_skill_gap_graphs_analysis
    ON hskip_users.bridgr_skill_gap_graphs(analysis_uuid);

CREATE TRIGGER tr_bridgr_skill_gap_graphs_control_time
    BEFORE INSERT OR UPDATE ON hskip_users.bridgr_skill_gap_graphs
    FOR EACH ROW EXECUTE FUNCTION hskip_users.tr_control_time();

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

CREATE INDEX idx_bridgr_skill_gap_nodes_graph_uuid ON hskip_users.bridgr_skill_gap_nodes(graph_uuid);
CREATE INDEX idx_bridgr_skill_gap_nodes_node_key ON hskip_users.bridgr_skill_gap_nodes(node_key);

CREATE TRIGGER tr_bridgr_skill_gap_nodes_control_time
    BEFORE INSERT OR UPDATE ON hskip_users.bridgr_skill_gap_nodes
    FOR EACH ROW EXECUTE FUNCTION hskip_users.tr_control_time();

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

CREATE INDEX idx_bridgr_skill_gap_edges_graph_uuid ON hskip_users.bridgr_skill_gap_edges(graph_uuid);
CREATE INDEX idx_bridgr_skill_gap_edges_from ON hskip_users.bridgr_skill_gap_edges(from_node_uuid);
CREATE INDEX idx_bridgr_skill_gap_edges_to ON hskip_users.bridgr_skill_gap_edges(to_node_uuid);
CREATE INDEX idx_bridgr_skill_gap_edges_relation ON hskip_users.bridgr_skill_gap_edges(relation);

CREATE TRIGGER tr_bridgr_skill_gap_edges_control_time
    BEFORE INSERT OR UPDATE ON hskip_users.bridgr_skill_gap_edges
    FOR EACH ROW EXECUTE FUNCTION hskip_users.tr_control_time();

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

CREATE INDEX idx_bridgr_skill_gap_learning_paths_analysis_uuid ON hskip_users.bridgr_skill_gap_learning_paths(analysis_uuid);

CREATE TRIGGER tr_bridgr_skill_gap_learning_paths_control_time
    BEFORE INSERT OR UPDATE ON hskip_users.bridgr_skill_gap_learning_paths
    FOR EACH ROW EXECUTE FUNCTION hskip_users.tr_control_time();

CREATE TABLE hskip_users.bridgr_skill_gap_path_steps (
    uuid                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id                        BIGSERIAL NOT NULL UNIQUE,
    path_uuid                 UUID NOT NULL,
    step_index                INTEGER NOT NULL,
    title                     TEXT NOT NULL,
    rationale                 TEXT,
    estimated_hours           NUMERIC(10,2),
    resource_uri              TEXT,
    resource_kind             VARCHAR(64),
    founder_learning_item_uuid UUID,
    course_lesson_uuid        UUID,
    linked_node_keys          JSONB DEFAULT '[]',
    metadata                  JSONB DEFAULT '{}',
    created_at                TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at                TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    UNIQUE(path_uuid, step_index)
);

CREATE INDEX idx_bridgr_skill_gap_path_steps_path_uuid ON hskip_users.bridgr_skill_gap_path_steps(path_uuid);

CREATE TRIGGER tr_bridgr_skill_gap_path_steps_control_time
    BEFORE INSERT OR UPDATE ON hskip_users.bridgr_skill_gap_path_steps
    FOR EACH ROW EXECUTE FUNCTION hskip_users.tr_control_time();

CREATE TABLE hskip_users.bridgr_skill_gap_path_step_deps (
    uuid                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id                     BIGSERIAL NOT NULL UNIQUE,
    path_uuid              UUID NOT NULL,
    step_uuid              UUID NOT NULL,
    depends_on_step_uuid   UUID NOT NULL,
    created_at             TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at             TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    CONSTRAINT bridgr_skill_gap_path_step_deps_no_self CHECK (step_uuid != depends_on_step_uuid),
    UNIQUE(path_uuid, step_uuid, depends_on_step_uuid)
);

CREATE INDEX idx_bridgr_skill_gap_path_step_deps_path ON hskip_users.bridgr_skill_gap_path_step_deps(path_uuid);
CREATE INDEX idx_bridgr_skill_gap_path_step_deps_step ON hskip_users.bridgr_skill_gap_path_step_deps(step_uuid);
CREATE INDEX idx_bridgr_skill_gap_path_step_deps_depends ON hskip_users.bridgr_skill_gap_path_step_deps(depends_on_step_uuid);

CREATE TRIGGER tr_bridgr_skill_gap_path_step_deps_control_time
    BEFORE INSERT OR UPDATE ON hskip_users.bridgr_skill_gap_path_step_deps
    FOR EACH ROW EXECUTE FUNCTION hskip_users.tr_control_time();

CREATE TABLE hskip_users.bridgr_skill_gap_coverage (
    uuid                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id                  BIGSERIAL NOT NULL UNIQUE,
    analysis_uuid       UUID NOT NULL,
    coverage_kind       VARCHAR(32) NOT NULL DEFAULT 'role_skill',
    role_skill_key      VARCHAR(256),
    candidate_skill_key VARCHAR(256),
    match_status        VARCHAR(32) NOT NULL DEFAULT 'unknown',
    summary             TEXT,
    metrics             JSONB DEFAULT '{}',
    created_at          TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at          TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    CONSTRAINT bridgr_skill_gap_coverage_kind_chk CHECK (coverage_kind IN ('role_skill', 'summary', 'aggregate'))
);

CREATE INDEX idx_bridgr_skill_gap_coverage_analysis_uuid ON hskip_users.bridgr_skill_gap_coverage(analysis_uuid);
CREATE INDEX idx_bridgr_skill_gap_coverage_kind ON hskip_users.bridgr_skill_gap_coverage(analysis_uuid, coverage_kind);

CREATE TRIGGER tr_bridgr_skill_gap_coverage_control_time
    BEFORE INSERT OR UPDATE ON hskip_users.bridgr_skill_gap_coverage
    FOR EACH ROW EXECUTE FUNCTION hskip_users.tr_control_time();
