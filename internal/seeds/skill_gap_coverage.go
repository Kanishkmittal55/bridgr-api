package seeds

import (
	"context"
	"fmt"
	"os"

	"github.com/jszwec/csvutil"
)

type skillGapCoverageSeed struct {
	UUID              string   `csv:"uuid"`
	ID                int64    `csv:"id"`
	AnalysisUUID      string   `csv:"analysis_uuid"`
	CoverageKind      string   `csv:"coverage_kind"`
	RoleSkillKey      string   `csv:"role_skill_key"`
	CandidateSkillKey string   `csv:"candidate_skill_key"`
	MatchStatus       string   `csv:"match_status"`
	Summary           string   `csv:"summary"`
	Metrics           string   `csv:"metrics"`
	CreatedAt         SeedTime `csv:"created_at"`
	UpdatedAt         SeedTime `csv:"updated_at"`
}

func (s Seed) SeedSkillGapCoverage(ctx context.Context) {
	csvPath := "./internal/seeds/seeddata/skill_gap_coverage.csv"
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		fmt.Println("skill_gap_coverage.csv not found, skipping")
		return
	}
	b, err := os.ReadFile(csvPath)
	if err != nil {
		panic(err)
	}
	rows := []skillGapCoverageSeed{}
	if err := csvutil.Unmarshal(b, &rows); err != nil {
		panic(err)
	}
	if len(rows) == 0 {
		fmt.Println("No skill_gap_coverage to seed")
		return
	}
	var count int
	if err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM bridgr.skill_gap_coverage`).Scan(&count); err != nil {
		panic(err)
	}
	if count > 0 {
		fmt.Printf("skill_gap_coverage already seeded (%d rows), skipping...\n", count)
		return
	}
	const q = `
		INSERT INTO bridgr.skill_gap_coverage (
			uuid, id, analysis_uuid, coverage_kind, role_skill_key, candidate_skill_key,
			match_status, summary, metrics, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9::jsonb,$10,$11)`
	for _, r := range rows {
		metrics := r.Metrics
		if metrics == "" {
			metrics = "{}"
		}
		if _, err := s.db.Exec(ctx, q,
			r.UUID, r.ID, r.AnalysisUUID, r.CoverageKind,
			strOrNil(r.RoleSkillKey), strOrNil(r.CandidateSkillKey),
			r.MatchStatus, strOrNil(r.Summary), metrics, r.CreatedAt, r.UpdatedAt,
		); err != nil {
			panic(fmt.Errorf("seed skill_gap_coverage: %w", err))
		}
	}
	fmt.Printf("✓ Seeded %d skill_gap_coverage\n", len(rows))
}
