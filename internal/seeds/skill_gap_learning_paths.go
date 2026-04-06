package seeds

import (
	"context"
	"fmt"
	"os"

	"github.com/jszwec/csvutil"
)

type skillGapLearningPathSeed struct {
	UUID         string   `csv:"uuid"`
	ID           int64    `csv:"id"`
	AnalysisUUID string   `csv:"analysis_uuid"`
	PathVersion  int      `csv:"path_version"`
	Algorithm    string   `csv:"algorithm"`
	Title        string   `csv:"title"`
	PathMetadata string   `csv:"path_metadata"`
	CreatedAt    SeedTime `csv:"created_at"`
	UpdatedAt    SeedTime `csv:"updated_at"`
}

func (s Seed) SeedSkillGapLearningPaths(ctx context.Context) {
	csvPath := "./internal/seeds/seeddata/skill_gap_learning_paths.csv"
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		fmt.Println("skill_gap_learning_paths.csv not found, skipping")
		return
	}
	b, err := os.ReadFile(csvPath)
	if err != nil {
		panic(err)
	}
	rows := []skillGapLearningPathSeed{}
	if err := csvutil.Unmarshal(b, &rows); err != nil {
		panic(err)
	}
	if len(rows) == 0 {
		fmt.Println("No skill_gap_learning_paths to seed")
		return
	}
	var count int
	if err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM bridgr.skill_gap_learning_paths`).Scan(&count); err != nil {
		panic(err)
	}
	if count > 0 {
		fmt.Printf("skill_gap_learning_paths already seeded (%d rows), skipping...\n", count)
		return
	}
	const q = `
		INSERT INTO bridgr.skill_gap_learning_paths (
			uuid, id, analysis_uuid, path_version, algorithm, title, path_metadata, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7::jsonb,$8,$9)`
	for _, r := range rows {
		md := r.PathMetadata
		if md == "" {
			md = "{}"
		}
		if _, err := s.db.Exec(ctx, q,
			r.UUID, r.ID, r.AnalysisUUID, r.PathVersion,
			strOrNil(r.Algorithm), strOrNil(r.Title), md, r.CreatedAt, r.UpdatedAt,
		); err != nil {
			panic(fmt.Errorf("seed skill_gap_learning_paths: %w", err))
		}
	}
	fmt.Printf("✓ Seeded %d skill_gap_learning_paths\n", len(rows))
}
