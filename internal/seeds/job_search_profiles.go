package seeds

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jszwec/csvutil"
)

// CSV columns match bridgr.job_search_profiles (migration 20260405193142):
// uuid, id, user_id, target_role, location, source_board,
// career_switch, company_stage, seniority_goal, compensation_goal, software_stack (semicolon-separated tags),
// canonical_cv_analysis_uuid, created_at, updated_at.
// Rows must be unique on the API natural key (role, location, board, goals, stack, CV) per ensureUniqueJobSearchProfile.

type jobSearchProfileSeed struct {
	UUID                    string   `csv:"uuid"`
	ID                      int64    `csv:"id"`
	UserID                  int      `csv:"user_id"`
	TargetRole              string   `csv:"target_role"`
	Location                string   `csv:"location"`
	SourceBoard             string   `csv:"source_board"`
	CareerSwitch            bool     `csv:"career_switch"`
	CompanyStage            string   `csv:"company_stage"`
	SeniorityGoal           string   `csv:"seniority_goal"`
	CompensationGoal        string   `csv:"compensation_goal"`
	SoftwareStack           string   `csv:"software_stack"`
	CanonicalCvAnalysisUUID string   `csv:"canonical_cv_analysis_uuid"`
	CreatedAt               SeedTime `csv:"created_at"`
	UpdatedAt               SeedTime `csv:"updated_at"`
}

func parseSoftwareStackCSV(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ";")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func (s Seed) SeedJobSearchProfiles(ctx context.Context) {
	csvPath := "./internal/seeds/seeddata/job_search_profiles.csv"
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		fmt.Println("job_search_profiles.csv not found, skipping")
		return
	}
	b, err := os.ReadFile(csvPath)
	if err != nil {
		panic(err)
	}
	rows := []jobSearchProfileSeed{}
	if err := csvutil.Unmarshal(b, &rows); err != nil {
		panic(err)
	}
	if len(rows) == 0 {
		fmt.Println("No job_search_profiles to seed")
		return
	}
	var count int
	if err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM bridgr.job_search_profiles`).Scan(&count); err != nil {
		panic(err)
	}
	if count > 0 {
		fmt.Printf("job_search_profiles already seeded (%d rows), skipping...\n", count)
		return
	}
	const q = `
		INSERT INTO bridgr.job_search_profiles (
			uuid, id, user_id, target_role, location, source_board,
			career_switch, company_stage, seniority_goal, compensation_goal, software_stack_must_have,
			canonical_cv_analysis_uuid, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`
	for _, r := range rows {
		sb := r.SourceBoard
		if sb == "" {
			sb = "indeed"
		}
		cs := r.CompanyStage
		if cs == "" {
			cs = "any"
		}
		sg := r.SeniorityGoal
		if sg == "" {
			sg = "any"
		}
		cg := r.CompensationGoal
		if cg == "" {
			cg = "any"
		}
		stack := parseSoftwareStackCSV(r.SoftwareStack)
		if _, err := s.db.Exec(ctx, q,
			r.UUID, r.ID, r.UserID, r.TargetRole, r.Location, sb,
			r.CareerSwitch, cs, sg, cg, stack,
			strOrNil(r.CanonicalCvAnalysisUUID), r.CreatedAt, r.UpdatedAt,
		); err != nil {
			panic(fmt.Errorf("seed job_search_profiles: %w", err))
		}
	}
	fmt.Printf("✓ Seeded %d job_search_profiles\n", len(rows))
}
