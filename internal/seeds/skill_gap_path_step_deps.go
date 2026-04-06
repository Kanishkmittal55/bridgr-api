package seeds

import (
	"context"
	"fmt"
	"os"

	"github.com/jszwec/csvutil"
)

type skillGapPathStepDepSeed struct {
	UUID              string   `csv:"uuid"`
	ID                int64    `csv:"id"`
	PathUUID          string   `csv:"path_uuid"`
	StepUUID          string   `csv:"step_uuid"`
	DependsOnStepUUID string   `csv:"depends_on_step_uuid"`
	CreatedAt         SeedTime `csv:"created_at"`
	UpdatedAt         SeedTime `csv:"updated_at"`
}

func (s Seed) SeedSkillGapPathStepDeps(ctx context.Context) {
	csvPath := "./internal/seeds/seeddata/skill_gap_path_step_deps.csv"
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		fmt.Println("skill_gap_path_step_deps.csv not found, skipping")
		return
	}
	b, err := os.ReadFile(csvPath)
	if err != nil {
		panic(err)
	}
	rows := []skillGapPathStepDepSeed{}
	if err := csvutil.Unmarshal(b, &rows); err != nil {
		panic(err)
	}
	if len(rows) == 0 {
		fmt.Println("No skill_gap_path_step_deps to seed")
		return
	}
	var count int
	if err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM bridgr.skill_gap_path_step_deps`).Scan(&count); err != nil {
		panic(err)
	}
	if count > 0 {
		fmt.Printf("skill_gap_path_step_deps already seeded (%d rows), skipping...\n", count)
		return
	}
	const q = `
		INSERT INTO bridgr.skill_gap_path_step_deps (
			uuid, id, path_uuid, step_uuid, depends_on_step_uuid, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	for _, r := range rows {
		if _, err := s.db.Exec(ctx, q,
			r.UUID, r.ID, r.PathUUID, r.StepUUID, r.DependsOnStepUUID, r.CreatedAt, r.UpdatedAt,
		); err != nil {
			panic(fmt.Errorf("seed skill_gap_path_step_deps: %w", err))
		}
	}
	fmt.Printf("✓ Seeded %d skill_gap_path_step_deps\n", len(rows))
}
