package seeds

import (
	"context"
	"fmt"
	"os"

	"github.com/jszwec/csvutil"
)

type jobSearchDiscoveryRunSeed struct {
	UUID              string   `csv:"uuid"`
	ID                int64    `csv:"id"`
	UserID            int      `csv:"user_id"`
	Status            string   `csv:"status"`
	RequestParams     string   `csv:"request_params"`
	RadarMeta         string   `csv:"radar_meta"`
	RawCandidateCount int      `csv:"raw_candidate_count"`
	NewCandidateCount int      `csv:"new_candidate_count"`
	StartedAt         SeedTime `csv:"started_at"`
	CompletedAt       SeedTime `csv:"completed_at"`
	ErrorCode         string   `csv:"error_code"`
	ErrorDetail       string   `csv:"error_detail"`
	SqsMessageID      string   `csv:"sqs_message_id"`
	CreatedAt         SeedTime `csv:"created_at"`
	UpdatedAt         SeedTime `csv:"updated_at"`
}

func (s Seed) SeedJobSearchDiscoveryRuns(ctx context.Context) {
	csvPath := "./internal/seeds/seeddata/job_search_discovery_runs.csv"
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		fmt.Println("job_search_discovery_runs.csv not found, skipping")
		return
	}
	b, err := os.ReadFile(csvPath)
	if err != nil {
		panic(err)
	}
	rows := []jobSearchDiscoveryRunSeed{}
	if err := csvutil.Unmarshal(b, &rows); err != nil {
		panic(err)
	}
	if len(rows) == 0 {
		fmt.Println("No job_search_discovery_runs to seed")
		return
	}
	var count int
	if err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM bridgr.job_search_discovery_runs`).Scan(&count); err != nil {
		panic(err)
	}
	if count > 0 {
		fmt.Printf("job_search_discovery_runs already seeded (%d rows), skipping...\n", count)
		return
	}
	const q = `
		INSERT INTO bridgr.job_search_discovery_runs (
			uuid, id, user_id, status, request_params, radar_meta,
			raw_candidate_count, new_candidate_count, started_at, completed_at,
			error_code, error_detail, sqs_message_id, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5::jsonb,$6::jsonb,$7,$8,$9,$10,$11,$12,$13,$14,$15)`
	for _, r := range rows {
		rp := r.RequestParams
		if rp == "" {
			rp = "{}"
		}
		rm := r.RadarMeta
		if rm == "" {
			rm = "{}"
		}
		var started, completed interface{}
		if !r.StartedAt.isNil {
			started = r.StartedAt.Time
		}
		if !r.CompletedAt.isNil {
			completed = r.CompletedAt.Time
		}
		if _, err := s.db.Exec(ctx, q,
			r.UUID, r.ID, r.UserID, r.Status, rp, rm,
			r.RawCandidateCount, r.NewCandidateCount, started, completed,
			strOrNil(r.ErrorCode), strOrNil(r.ErrorDetail), strOrNil(r.SqsMessageID),
			r.CreatedAt, r.UpdatedAt,
		); err != nil {
			panic(fmt.Errorf("seed job_search_discovery_runs: %w", err))
		}
	}
	fmt.Printf("✓ Seeded %d job_search_discovery_runs\n", len(rows))
}
