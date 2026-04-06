package seeds

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/jszwec/csvutil"
)

type skillGapPathStepSeed struct {
	UUID                    string   `csv:"uuid"`
	ID                      int64    `csv:"id"`
	PathUUID                string   `csv:"path_uuid"`
	StepIndex               int      `csv:"step_index"`
	Title                   string   `csv:"title"`
	Rationale               string   `csv:"rationale"`
	EstimatedHours          string   `csv:"estimated_hours"`
	ResourceURI             string   `csv:"resource_uri"`
	ResourceKind            string   `csv:"resource_kind"`
	FounderLearningItemUUID string   `csv:"founder_learning_item_uuid"`
	CourseLessonUUID        string   `csv:"course_lesson_uuid"`
	LinkedNodeKeys          string   `csv:"linked_node_keys"`
	Metadata                string   `csv:"metadata"`
	CreatedAt               SeedTime `csv:"created_at"`
	UpdatedAt               SeedTime `csv:"updated_at"`
}

func (s Seed) SeedSkillGapPathSteps(ctx context.Context) {
	csvPath := "./internal/seeds/seeddata/skill_gap_path_steps.csv"
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		fmt.Println("skill_gap_path_steps.csv not found, skipping")
		return
	}
	b, err := os.ReadFile(csvPath)
	if err != nil {
		panic(err)
	}
	rows := []skillGapPathStepSeed{}
	if err := csvutil.Unmarshal(b, &rows); err != nil {
		panic(err)
	}
	if len(rows) == 0 {
		fmt.Println("No skill_gap_path_steps to seed")
		return
	}
	var count int
	if err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM bridgr.skill_gap_path_steps`).Scan(&count); err != nil {
		panic(err)
	}
	if count > 0 {
		fmt.Printf("skill_gap_path_steps already seeded (%d rows), skipping...\n", count)
		return
	}
	const q = `
		INSERT INTO bridgr.skill_gap_path_steps (
			uuid, id, path_uuid, step_index, title, rationale, estimated_hours,
			resource_uri, resource_kind, founder_learning_item_uuid, course_lesson_uuid,
			linked_node_keys, metadata, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12::jsonb,$13::jsonb,$14,$15)`
	for _, r := range rows {
		linked := r.LinkedNodeKeys
		if linked == "" {
			linked = "[]"
		}
		md := r.Metadata
		if md == "" {
			md = "{}"
		}
		var est interface{}
		if r.EstimatedHours != "" {
			f, err := strconv.ParseFloat(r.EstimatedHours, 64)
			if err != nil {
				panic(fmt.Errorf("skill_gap_path_steps estimated_hours %q: %w", r.EstimatedHours, err))
			}
			est = f
		}
		if _, err := s.db.Exec(ctx, q,
			r.UUID, r.ID, r.PathUUID, r.StepIndex, r.Title,
			strOrNil(r.Rationale), est,
			strOrNil(r.ResourceURI), strOrNil(r.ResourceKind),
			strOrNil(r.FounderLearningItemUUID), strOrNil(r.CourseLessonUUID),
			linked, md, r.CreatedAt, r.UpdatedAt,
		); err != nil {
			panic(fmt.Errorf("seed skill_gap_path_steps: %w", err))
		}
	}
	fmt.Printf("✓ Seeded %d skill_gap_path_steps\n", len(rows))
}
