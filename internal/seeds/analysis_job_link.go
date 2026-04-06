package seeds

import (
	"context"
	"fmt"
	"os"

	"github.com/jszwec/csvutil"
)

type analysisJobLinkSeed struct {
	UUID             string   `csv:"uuid"`
	ID               int64    `csv:"id"`
	UserID           int      `csv:"user_id"`
	AnalysisUUID     string   `csv:"analysis_uuid"`
	JobCandidateUUID string   `csv:"job_candidate_uuid"`
	LinkKind         string   `csv:"link_kind"`
	CreatedAt        SeedTime `csv:"created_at"`
	UpdatedAt        SeedTime `csv:"updated_at"`
}

func (s Seed) SeedAnalysisJobLink(ctx context.Context) {
	csvPath := "./internal/seeds/seeddata/analysis_job_link.csv"
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		fmt.Println("analysis_job_link.csv not found, skipping")
		return
	}
	b, err := os.ReadFile(csvPath)
	if err != nil {
		panic(err)
	}
	rows := []analysisJobLinkSeed{}
	if err := csvutil.Unmarshal(b, &rows); err != nil {
		panic(err)
	}
	if len(rows) == 0 {
		fmt.Println("No analysis_job_link to seed")
		return
	}
	var count int
	if err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM bridgr.analysis_job_link`).Scan(&count); err != nil {
		panic(err)
	}
	if count > 0 {
		fmt.Printf("analysis_job_link already seeded (%d rows), skipping...\n", count)
		return
	}
	const q = `
		INSERT INTO bridgr.analysis_job_link (
			uuid, id, user_id, analysis_uuid, job_candidate_uuid, link_kind, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	for _, r := range rows {
		if _, err := s.db.Exec(ctx, q,
			r.UUID, r.ID, r.UserID, r.AnalysisUUID, r.JobCandidateUUID, r.LinkKind, r.CreatedAt, r.UpdatedAt,
		); err != nil {
			panic(fmt.Errorf("seed analysis_job_link: %w", err))
		}
	}
	fmt.Printf("✓ Seeded %d analysis_job_link\n", len(rows))
}
