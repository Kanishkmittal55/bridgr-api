-- =============================================================================
-- BRIDGR SKILL GAP — COVERAGE (aggregated reads)
-- Metrics derived from graphs + node overlap (not only skill_gap_coverage rows).
-- =============================================================================

-- name: GetSkillGapCoverage :one
SELECT
    s.candidate_skills,
    s.required_skills,
    s.matched_skills,
    CASE
        WHEN s.required_skills > 0 THEN
            ROUND(100.0 * s.matched_skills::numeric / s.required_skills::numeric, 2)
        ELSE NULL::numeric
    END AS coverage_pct
FROM (
    SELECT
        COALESCE(
            (
                SELECT COUNT(*)::bigint
                FROM bridgr.skill_gap_nodes n
                INNER JOIN bridgr.skill_gap_graphs g
                    ON g.uuid = n.graph_uuid AND g.kind = 'candidate'
                WHERE g.analysis_uuid = $1
            ),
            0
        ) AS candidate_skills,
        COALESCE(
            (
                SELECT COUNT(*)::bigint
                FROM bridgr.skill_gap_nodes n
                INNER JOIN bridgr.skill_gap_graphs g
                    ON g.uuid = n.graph_uuid AND g.kind = 'role_requirement'
                WHERE g.analysis_uuid = $1
            ),
            0
        ) AS required_skills,
        COALESCE(
            (
                SELECT COUNT(DISTINCT rn.node_key)::bigint
                FROM bridgr.skill_gap_nodes rn
                INNER JOIN bridgr.skill_gap_graphs rg
                    ON rg.uuid = rn.graph_uuid AND rg.kind = 'role_requirement'
                WHERE rg.analysis_uuid = $1
                  AND EXISTS (
                      SELECT 1
                      FROM bridgr.skill_gap_nodes cn
                      INNER JOIN bridgr.skill_gap_graphs cg
                          ON cg.uuid = cn.graph_uuid AND cg.kind = 'candidate'
                      WHERE cg.analysis_uuid = $1 AND cn.node_key = rn.node_key
                  )
            ),
            0
        ) AS matched_skills
) s;

-- name: GetSkillGapCoverageByUser :many
SELECT
    sub.analysis_id,
    sub.analysis_uuid,
    sub.user_id,
    sub.title,
    sub.status,
    sub.created_at,
    sub.candidate_skills,
    sub.required_skills,
    sub.matched_skills,
    CASE
        WHEN sub.required_skills > 0 THEN
            ROUND(100.0 * sub.matched_skills::numeric / sub.required_skills::numeric, 2)
        ELSE NULL::numeric
    END AS coverage_pct
FROM (
    SELECT
        a.id AS analysis_id,
        a.uuid AS analysis_uuid,
        a.user_id,
        a.title,
        a.status,
        a.created_at,
        COALESCE(
            (
                SELECT COUNT(*)::bigint
                FROM bridgr.skill_gap_nodes n
                INNER JOIN bridgr.skill_gap_graphs g
                    ON g.uuid = n.graph_uuid AND g.kind = 'candidate'
                WHERE g.analysis_uuid = a.uuid
            ),
            0
        ) AS candidate_skills,
        COALESCE(
            (
                SELECT COUNT(*)::bigint
                FROM bridgr.skill_gap_nodes n
                INNER JOIN bridgr.skill_gap_graphs g
                    ON g.uuid = n.graph_uuid AND g.kind = 'role_requirement'
                WHERE g.analysis_uuid = a.uuid
            ),
            0
        ) AS required_skills,
        COALESCE(
            (
                SELECT COUNT(DISTINCT rn.node_key)::bigint
                FROM bridgr.skill_gap_nodes rn
                INNER JOIN bridgr.skill_gap_graphs rg
                    ON rg.uuid = rn.graph_uuid AND rg.kind = 'role_requirement'
                WHERE rg.analysis_uuid = a.uuid
                  AND EXISTS (
                      SELECT 1
                      FROM bridgr.skill_gap_nodes cn
                      INNER JOIN bridgr.skill_gap_graphs cg
                          ON cg.uuid = cn.graph_uuid AND cg.kind = 'candidate'
                      WHERE cg.analysis_uuid = a.uuid AND cn.node_key = rn.node_key
                  )
            ),
            0
        ) AS matched_skills
    FROM bridgr.skill_gap_analyses a
    WHERE a.user_id = $1
) sub
ORDER BY sub.created_at DESC
LIMIT $2 OFFSET $3;
