package seeds

import (
	"context"
	"fmt"
	"os"

	"github.com/jszwec/csvutil"
)

type jobHarvestScheduleSeed struct {
	UUID           string   `csv:"uuid"`
	ID             int64    `csv:"id"`
	UserID         int      `csv:"user_id"`
	ProfileUUID    string   `csv:"profile_uuid"`
	Enabled        bool     `csv:"enabled"`
	CadenceMinutes int      `csv:"cadence_minutes"`
	BoardsRotation string   `csv:"boards_rotation"`
	LastRunAt      SeedTime `csv:"last_run_at"`
	NextRunAt      SeedTime `csv:"next_run_at"`
	CreatedAt      SeedTime `csv:"created_at"`
	UpdatedAt      SeedTime `csv:"updated_at"`
}

func (s Seed) SeedJobHarvestSchedules(ctx context.Context) {
	csvPath := "./internal/seeds/seeddata/job_harvest_schedules.csv"
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		fmt.Println("job_harvest_schedules.csv not found, skipping")
		return
	}
	b, err := os.ReadFile(csvPath)
	if err != nil {
		panic(err)
	}
	rows := []jobHarvestScheduleSeed{}
	if err := csvutil.Unmarshal(b, &rows); err != nil {
		panic(err)
	}
	if len(rows) == 0 {
		fmt.Println("No job_harvest_schedules to seed")
		return
	}
	var count int
	if err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM bridgr.job_harvest_schedules`).Scan(&count); err != nil {
		panic(err)
	}
	if count > 0 {
		fmt.Printf("job_harvest_schedules already seeded (%d rows), skipping...\n", count)
		return
	}
	const q = `
		INSERT INTO bridgr.job_harvest_schedules (
			uuid, id, user_id, profile_uuid, enabled, cadence_minutes, boards_rotation,
			last_run_at, next_run_at, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7::jsonb,$8,$9,$10,$11)`
	for _, r := range rows {
		br := r.BoardsRotation
		if br == "" {
			br = "[]"
		}
		var last, next interface{}
		if !r.LastRunAt.isNil {
			last = r.LastRunAt.Time
		}
		if !r.NextRunAt.isNil {
			next = r.NextRunAt.Time
		}
		if _, err := s.db.Exec(ctx, q,
			r.UUID, r.ID, r.UserID, r.ProfileUUID, r.Enabled, r.CadenceMinutes, br,
			last, next, r.CreatedAt, r.UpdatedAt,
		); err != nil {
			panic(fmt.Errorf("seed job_harvest_schedules: %w", err))
		}
	}
	fmt.Printf("✓ Seeded %d job_harvest_schedules\n", len(rows))
}
