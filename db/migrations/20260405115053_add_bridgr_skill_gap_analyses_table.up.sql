-- ═══════════════════════════════════════════════════════════════════════════
-- BRIDGR SKILL GAP ANALYSES — one CV+JD skill-gap analysis run per user/session
-- Stores run header, asset pointers, LLM extraction metadata, gap summary, Mermaid.
-- No foreign keys (app-enforced).
-- ═══════════════════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS bridgr.skill_gap_analyses;

CREATE TABLE bridgr.skill_gap_analyses (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id BIGSERIAL NOT NULL UNIQUE,

    user_id INTEGER NOT NULL,              -- application enforced FK to users.id
    founder_persona_uuid UUID,             -- application enforced FK to founder_persona.uuid
    pursuit_uuid UUID,                     -- application enforced FK to founder_pursuits.uuid

    title TEXT,

    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    -- pending, extracting, graphed, pathed, completed, failed

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

    CONSTRAINT chk_skill_gap_analyses_status
        CHECK (status IN ('pending', 'extracting', 'graphed', 'pathed', 'completed', 'failed'))
);

COMMENT ON TABLE bridgr.skill_gap_analyses IS 'Bridgr Skill Gap Navigator: analysis run linking user (and optional founder context) to CV/JD extraction and learning-path outputs.';
COMMENT ON COLUMN bridgr.skill_gap_analyses.status IS 'pending, extracting, graphed, pathed, completed, failed';
COMMENT ON COLUMN bridgr.skill_gap_analyses.cv_asset_uri IS 'Pointer to stored CV (e.g. S3); optional if inline processing only';
COMMENT ON COLUMN bridgr.skill_gap_analyses.jd_asset_uri IS 'Pointer to stored job description';
COMMENT ON COLUMN bridgr.skill_gap_analyses.cv_fingerprint IS 'Hash for dedup / idempotency';
COMMENT ON COLUMN bridgr.skill_gap_analyses.jd_fingerprint IS 'Hash for dedup / idempotency';
COMMENT ON COLUMN bridgr.skill_gap_analyses.extraction_payload IS 'Validated structured LLM output (skills graph extraction)';
COMMENT ON COLUMN bridgr.skill_gap_analyses.gap_summary IS 'Rollup metrics and narrative gap summary';
COMMENT ON COLUMN bridgr.skill_gap_analyses.mermaid_diagram IS 'Optional v1 DAG visualization text';
COMMENT ON COLUMN bridgr.skill_gap_analyses.sqs_message_id IS 'AWS SQS MessageId after successful enqueue (debugging / dedupe)';

CREATE INDEX idx_skill_gap_analyses_user_created
    ON bridgr.skill_gap_analyses(user_id, created_at DESC);
CREATE INDEX idx_skill_gap_analyses_status
    ON bridgr.skill_gap_analyses(status);
CREATE INDEX idx_skill_gap_analyses_persona
    ON bridgr.skill_gap_analyses(founder_persona_uuid);
CREATE INDEX idx_skill_gap_analyses_pursuit
    ON bridgr.skill_gap_analyses(pursuit_uuid);

CREATE TRIGGER tr_skill_gap_analyses_control_time
    BEFORE INSERT OR UPDATE ON bridgr.skill_gap_analyses
    FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();
