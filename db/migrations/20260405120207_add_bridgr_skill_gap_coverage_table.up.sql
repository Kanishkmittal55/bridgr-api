-- ═══════════════════════════════════════════════════════════════════════════
-- BRIDGR SKILL GAP NAVIGATOR — COVERAGE
-- Per-analysis coverage: role vs candidate alignment, summary rows, or metrics.
-- analysis_uuid app-enforced FK to bridgr_skill_gap_analyses.uuid.
-- No foreign keys (app-enforced).
-- ═══════════════════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS hskip_users.bridgr_skill_gap_coverage;

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

COMMENT ON TABLE hskip_users.bridgr_skill_gap_coverage IS 'Skill coverage / gap rows or summary metrics for an analysis.';
COMMENT ON COLUMN hskip_users.bridgr_skill_gap_coverage.analysis_uuid IS 'App-enforced FK to bridgr_skill_gap_analyses.uuid';
COMMENT ON COLUMN hskip_users.bridgr_skill_gap_coverage.coverage_kind IS 'role_skill: pairing row; summary: one row roll-up; aggregate: optional bucket';
COMMENT ON COLUMN hskip_users.bridgr_skill_gap_coverage.match_status IS 'covered, gap, partial, surplus, not_applicable, unknown';

CREATE INDEX idx_bridgr_skill_gap_coverage_analysis_uuid ON hskip_users.bridgr_skill_gap_coverage(analysis_uuid);
CREATE INDEX idx_bridgr_skill_gap_coverage_kind ON hskip_users.bridgr_skill_gap_coverage(analysis_uuid, coverage_kind);

CREATE TRIGGER tr_bridgr_skill_gap_coverage_control_time
    BEFORE INSERT OR UPDATE ON hskip_users.bridgr_skill_gap_coverage
    FOR EACH ROW EXECUTE FUNCTION hskip_users.tr_control_time();
