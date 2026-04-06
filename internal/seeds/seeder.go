package seeds

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Seed holds a pgx pool for Bridgr schema seeding.
type Seed struct {
	db *pgxpool.Pool
}

// Execute runs all CSV seeders in FK-safe order. Paths are relative to the repo root
// (same convention as github.com/hassleskip/users: run from project root).
func Execute(db *pgxpool.Pool) {
	s := Seed{db: db}
	ctx := context.Background()

	s.SeedSupportedJobBoards(ctx)
	s.SeedSkillGapAnalyses(ctx)
	s.SeedSkillGapGraphs(ctx)
	s.SeedSkillGapNodes(ctx)
	s.SeedSkillGapEdges(ctx)
	s.SeedSkillGapCoverage(ctx)
	s.SeedSkillGapLearningPaths(ctx)
	s.SeedSkillGapPathSteps(ctx)
	s.SeedSkillGapPathStepDeps(ctx)

	// Job search profiles (Radar slices) — must run before job_harvest_schedules, which references profile_uuid.
	s.SeedJobSearchProfiles(ctx)
	s.SeedJobHarvestSchedules(ctx)
	s.SeedJobSearchDiscoveryRuns(ctx)
	s.SeedJobCandidates(ctx)
	s.SeedJobEnrichments(ctx)
	s.SeedJobScores(ctx)
	s.SeedFeedItems(ctx)
	s.SeedJobNotifications(ctx)
	s.SeedAnalysisJobLink(ctx)

	s.ResetSequences(ctx)
}

// ResetSequences sets SERIAL/BIGSERIAL sequences to max(id) after explicit IDs from CSV.
func (s Seed) ResetSequences(ctx context.Context) {
	tables := []struct {
		table    string
		sequence string
	}{
		{"bridgr.supported_job_boards", "bridgr.supported_job_boards_id_seq"},
		{"bridgr.skill_gap_analyses", "bridgr.skill_gap_analyses_id_seq"},
		{"bridgr.skill_gap_graphs", "bridgr.skill_gap_graphs_id_seq"},
		{"bridgr.skill_gap_nodes", "bridgr.skill_gap_nodes_id_seq"},
		{"bridgr.skill_gap_edges", "bridgr.skill_gap_edges_id_seq"},
		{"bridgr.skill_gap_coverage", "bridgr.skill_gap_coverage_id_seq"},
		{"bridgr.skill_gap_learning_paths", "bridgr.skill_gap_learning_paths_id_seq"},
		{"bridgr.skill_gap_path_steps", "bridgr.skill_gap_path_steps_id_seq"},
		{"bridgr.skill_gap_path_step_deps", "bridgr.skill_gap_path_step_deps_id_seq"},
		{"bridgr.job_search_profiles", "bridgr.job_search_profiles_id_seq"},
		{"bridgr.job_harvest_schedules", "bridgr.job_harvest_schedules_id_seq"},
		{"bridgr.job_search_discovery_runs", "bridgr.job_search_discovery_runs_id_seq"},
		{"bridgr.job_candidates", "bridgr.job_candidates_id_seq"},
		{"bridgr.job_enrichments", "bridgr.job_enrichments_id_seq"},
		{"bridgr.job_scores", "bridgr.job_scores_id_seq"},
		{"bridgr.feed_items", "bridgr.feed_items_id_seq"},
		{"bridgr.job_notifications", "bridgr.job_notifications_id_seq"},
		{"bridgr.analysis_job_link", "bridgr.analysis_job_link_id_seq"},
	}

	for _, t := range tables {
		var maxID *int64
		q := fmt.Sprintf("SELECT MAX(id) FROM %s", t.table)
		err := s.db.QueryRow(ctx, q).Scan(&maxID)
		if err != nil || maxID == nil {
			continue
		}
		_, err = s.db.Exec(ctx, fmt.Sprintf("SELECT setval('%s', %d, true)", t.sequence, *maxID))
		if err != nil {
			fmt.Printf("⚠ Warning: failed to reset sequence %s: %v\n", t.sequence, err)
		}
	}
	fmt.Println("✓ Reset Bridgr sequences to max(id)")
}
