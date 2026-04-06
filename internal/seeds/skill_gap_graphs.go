package seeds

import (
	"context"
	"fmt"
	"os"

	"github.com/jszwec/csvutil"
)

type skillGapGraphSeed struct {
	UUID         string   `csv:"uuid"`
	ID           int64    `csv:"id"`
	AnalysisUUID string   `csv:"analysis_uuid"`
	Kind         string   `csv:"kind"`
	Metadata     string   `csv:"metadata"`
	CreatedAt    SeedTime `csv:"created_at"`
	UpdatedAt    SeedTime `csv:"updated_at"`
}

func (s Seed) SeedSkillGapGraphs(ctx context.Context) {
	csvPath := "./internal/seeds/seeddata/skill_gap_graphs.csv"
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		fmt.Println("skill_gap_graphs.csv not found, skipping")
		return
	}
	b, err := os.ReadFile(csvPath)
	if err != nil {
		panic(err)
	}
	rows := []skillGapGraphSeed{}
	if err := csvutil.Unmarshal(b, &rows); err != nil {
		panic(err)
	}
	if len(rows) == 0 {
		fmt.Println("No skill_gap_graphs to seed")
		return
	}
	var count int
	if err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM bridgr.skill_gap_graphs`).Scan(&count); err != nil {
		panic(err)
	}
	if count > 0 {
		fmt.Printf("skill_gap_graphs already seeded (%d rows), skipping...\n", count)
		return
	}
	const q = `
		INSERT INTO bridgr.skill_gap_graphs (uuid, id, analysis_uuid, kind, metadata, created_at, updated_at)
		VALUES ($1,$2,$3,$4,CAST($5 AS jsonb),$6,$7)`
	for _, r := range rows {
		var meta interface{}
		if r.Metadata != "" {
			meta = r.Metadata
		}
		if _, err := s.db.Exec(ctx, q, r.UUID, r.ID, r.AnalysisUUID, r.Kind, meta, r.CreatedAt, r.UpdatedAt); err != nil {
			panic(fmt.Errorf("seed skill_gap_graphs: %w", err))
		}
	}
	fmt.Printf("✓ Seeded %d skill_gap_graphs\n", len(rows))
}
