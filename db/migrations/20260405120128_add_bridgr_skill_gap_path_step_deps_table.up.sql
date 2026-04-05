-- ═══════════════════════════════════════════════════════════════════════════
-- BRIDGR SKILL GAP NAVIGATOR — PATH STEP DEPENDENCIES
-- DAG edges: step_uuid must be completed after depends_on_step_uuid.
-- path_uuid, step_uuid, depends_on_step_uuid app-enforced FKs.
-- No foreign keys (app-enforced).
-- ═══════════════════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS bridgr.skill_gap_path_step_deps;

CREATE TABLE bridgr.skill_gap_path_step_deps (
    uuid                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id                     BIGSERIAL NOT NULL UNIQUE,
    path_uuid              UUID NOT NULL,
    step_uuid              UUID NOT NULL,
    depends_on_step_uuid   UUID NOT NULL,
    created_at             TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at             TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    CONSTRAINT skill_gap_path_step_deps_no_self CHECK (step_uuid != depends_on_step_uuid),
    UNIQUE(path_uuid, step_uuid, depends_on_step_uuid)
);

COMMENT ON TABLE bridgr.skill_gap_path_step_deps IS 'Prerequisite edges between steps in a learning path (DAG).';
COMMENT ON COLUMN bridgr.skill_gap_path_step_deps.path_uuid IS 'App-enforced FK to skill_gap_learning_paths.uuid';
COMMENT ON COLUMN bridgr.skill_gap_path_step_deps.step_uuid IS 'App-enforced FK to skill_gap_path_steps.uuid — dependent step';
COMMENT ON COLUMN bridgr.skill_gap_path_step_deps.depends_on_step_uuid IS 'App-enforced FK to skill_gap_path_steps.uuid — must be done first';

CREATE INDEX idx_skill_gap_path_step_deps_path ON bridgr.skill_gap_path_step_deps(path_uuid);
CREATE INDEX idx_skill_gap_path_step_deps_step ON bridgr.skill_gap_path_step_deps(step_uuid);
CREATE INDEX idx_skill_gap_path_step_deps_depends ON bridgr.skill_gap_path_step_deps(depends_on_step_uuid);

CREATE TRIGGER tr_skill_gap_path_step_deps_control_time
    BEFORE INSERT OR UPDATE ON bridgr.skill_gap_path_step_deps
    FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();
