package seeds

import (
	"context"
	"fmt"
	"os"

	"github.com/jszwec/csvutil"
)

type skillGapEdgeSeed struct {
	UUID      string   `csv:"uuid"`
	ID        int64    `csv:"id"`
	GraphUUID string   `csv:"graph_uuid"`
	FromNode  string   `csv:"from_node_uuid"`
	ToNode    string   `csv:"to_node_uuid"`
	Relation  string   `csv:"relation"`
	Weight    float64  `csv:"weight"`
	Metadata  string   `csv:"metadata"`
	CreatedAt SeedTime `csv:"created_at"`
	UpdatedAt SeedTime `csv:"updated_at"`
}

func (s Seed) SeedSkillGapEdges(ctx context.Context) {
	csvPath := "./internal/seeds/seeddata/skill_gap_edges.csv"
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		fmt.Println("skill_gap_edges.csv not found, skipping")
		return
	}
	b, err := os.ReadFile(csvPath)
	if err != nil {
		panic(err)
	}
	rows := []skillGapEdgeSeed{}
	if err := csvutil.Unmarshal(b, &rows); err != nil {
		panic(err)
	}
	if len(rows) == 0 {
		fmt.Println("No skill_gap_edges to seed")
		return
	}
	var count int
	if err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM bridgr.skill_gap_edges`).Scan(&count); err != nil {
		panic(err)
	}
	if count > 0 {
		fmt.Printf("skill_gap_edges already seeded (%d rows), skipping...\n", count)
		return
	}
	const q = `
		INSERT INTO bridgr.skill_gap_edges (
			uuid, id, graph_uuid, from_node_uuid, to_node_uuid, relation, weight, metadata, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8::jsonb,$9,$10)`
	for _, r := range rows {
		md := r.Metadata
		if md == "" {
			md = "{}"
		}
		if _, err := s.db.Exec(ctx, q,
			r.UUID, r.ID, r.GraphUUID, r.FromNode, r.ToNode, r.Relation, r.Weight, md, r.CreatedAt, r.UpdatedAt,
		); err != nil {
			panic(fmt.Errorf("seed skill_gap_edges: %w", err))
		}
	}
	fmt.Printf("✓ Seeded %d skill_gap_edges\n", len(rows))
}
