package seeds

import (
	"context"
	"fmt"
	"os"

	"github.com/jszwec/csvutil"
)

type skillGapNodeSeed struct {
	UUID            string   `csv:"uuid"`
	ID              int64    `csv:"id"`
	GraphUUID       string   `csv:"graph_uuid"`
	NodeKey         string   `csv:"node_key"`
	DisplayName     string   `csv:"display_name"`
	Description     string   `csv:"description"`
	ProficiencyHint string   `csv:"proficiency_hint"`
	Source          string   `csv:"source"`
	Evidence        string   `csv:"evidence"`
	Metadata        string   `csv:"metadata"`
	PositionX       int      `csv:"position_x"`
	PositionY       int      `csv:"position_y"`
	CreatedAt       SeedTime `csv:"created_at"`
	UpdatedAt       SeedTime `csv:"updated_at"`
}

func (s Seed) SeedSkillGapNodes(ctx context.Context) {
	csvPath := "./internal/seeds/seeddata/skill_gap_nodes.csv"
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		fmt.Println("skill_gap_nodes.csv not found, skipping")
		return
	}
	b, err := os.ReadFile(csvPath)
	if err != nil {
		panic(err)
	}
	rows := []skillGapNodeSeed{}
	if err := csvutil.Unmarshal(b, &rows); err != nil {
		panic(err)
	}
	if len(rows) == 0 {
		fmt.Println("No skill_gap_nodes to seed")
		return
	}
	var count int
	if err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM bridgr.skill_gap_nodes`).Scan(&count); err != nil {
		panic(err)
	}
	if count > 0 {
		fmt.Printf("skill_gap_nodes already seeded (%d rows), skipping...\n", count)
		return
	}
	const q = `
		INSERT INTO bridgr.skill_gap_nodes (
			uuid, id, graph_uuid, node_key, display_name, description, proficiency_hint, source,
			evidence, metadata, position_x, position_y, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9::jsonb,$10::jsonb,$11,$12,$13,$14)`
	for _, r := range rows {
		ev := r.Evidence
		if ev == "" {
			ev = "{}"
		}
		md := r.Metadata
		if md == "" {
			md = "{}"
		}
		if _, err := s.db.Exec(ctx, q,
			r.UUID, r.ID, r.GraphUUID, r.NodeKey, r.DisplayName,
			strOrNil(r.Description), strOrNil(r.ProficiencyHint), strOrNil(r.Source),
			ev, md, r.PositionX, r.PositionY, r.CreatedAt, r.UpdatedAt,
		); err != nil {
			panic(fmt.Errorf("seed skill_gap_nodes: %w", err))
		}
	}
	fmt.Printf("✓ Seeded %d skill_gap_nodes\n", len(rows))
}
