-- ═══════════════════════════════════════════════════════════════════════════
-- BRIDGR SKILL GAP NAVIGATOR — PATH STEPS
-- Ordered steps in a learning path.
-- path_uuid app-enforced FK to skill_gap_learning_paths.uuid.
-- No foreign keys (app-enforced).
-- ═══════════════════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS bridgr.skill_gap_path_steps;

CREATE TABLE bridgr.skill_gap_path_steps (
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

COMMENT ON TABLE bridgr.skill_gap_path_steps IS 'One step in a generated learning path.';
COMMENT ON COLUMN bridgr.skill_gap_path_steps.path_uuid IS 'App-enforced FK to skill_gap_learning_paths.uuid';
COMMENT ON COLUMN bridgr.skill_gap_path_steps.step_index IS 'Order within the path (0-based or 1-based per app convention)';
COMMENT ON COLUMN bridgr.skill_gap_path_steps.linked_node_keys IS 'JSON array of graph node_key values this step addresses';
COMMENT ON COLUMN bridgr.skill_gap_path_steps.founder_learning_item_uuid IS 'App-enforced FK to founder_learning_items.uuid when linked';
COMMENT ON COLUMN bridgr.skill_gap_path_steps.course_lesson_uuid IS 'App-enforced FK to course lesson when linked';

CREATE INDEX idx_skill_gap_path_steps_path_uuid ON bridgr.skill_gap_path_steps(path_uuid);

CREATE TRIGGER tr_skill_gap_path_steps_control_time
    BEFORE INSERT OR UPDATE ON bridgr.skill_gap_path_steps
    FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();
