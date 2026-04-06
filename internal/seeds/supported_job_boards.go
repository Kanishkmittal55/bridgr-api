package seeds

import (
	"context"
	"fmt"
	"os"

	"github.com/jszwec/csvutil"
)

type supportedJobBoardSeed struct {
	UUID        string   `csv:"uuid"`
	ID          int64    `csv:"id"`
	BoardID     string   `csv:"board_id"`
	DisplayName string   `csv:"display_name"`
	Engine      string   `csv:"engine"`
	SiteType    string   `csv:"site_type"`
	Region      string   `csv:"region"`
	IsActive    bool     `csv:"is_active"`
	Config      string   `csv:"config"`
	SortOrder   int      `csv:"sort_order"`
	CreatedAt   SeedTime `csv:"created_at"`
	UpdatedAt   SeedTime `csv:"updated_at"`
}

func (s Seed) SeedSupportedJobBoards(ctx context.Context) {
	csvPath := "./internal/seeds/seeddata/supported_job_boards.csv"
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		fmt.Println("supported_job_boards.csv not found, skipping")
		return
	}
	b, err := os.ReadFile(csvPath)
	if err != nil {
		panic(err)
	}
	rows := []supportedJobBoardSeed{}
	if err := csvutil.Unmarshal(b, &rows); err != nil {
		panic(err)
	}
	if len(rows) == 0 {
		fmt.Println("No supported_job_boards to seed")
		return
	}
	var count int
	if err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM bridgr.supported_job_boards`).Scan(&count); err != nil {
		panic(err)
	}
	if count > 0 {
		fmt.Printf("supported_job_boards already seeded (%d rows), skipping...\n", count)
		return
	}
	const q = `
		INSERT INTO bridgr.supported_job_boards (
			uuid, id, board_id, display_name, engine, site_type, region, is_active, config, sort_order, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9::jsonb,$10,$11,$12)`
	for _, r := range rows {
		region := r.Region
		if region == "" {
			region = "global"
		}
		cfg := r.Config
		if cfg == "" {
			cfg = "{}"
		}
		if _, err := s.db.Exec(ctx, q, r.UUID, r.ID, r.BoardID, r.DisplayName, r.Engine, r.SiteType, region, r.IsActive, cfg, r.SortOrder, r.CreatedAt, r.UpdatedAt); err != nil {
			panic(fmt.Errorf("seed supported_job_boards %s: %w", r.BoardID, err))
		}
	}
	fmt.Printf("✓ Seeded %d supported_job_boards\n", len(rows))
}
