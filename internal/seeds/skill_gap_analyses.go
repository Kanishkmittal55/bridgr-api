package seeds

import (
	"context"
	"fmt"
	"os"

	"github.com/jszwec/csvutil"
)

type skillGapAnalysisSeed struct {
	UUID               string   `csv:"uuid"`
	ID                 int64    `csv:"id"`
	UserID             int      `csv:"user_id"`
	FounderPersonaUUID string   `csv:"founder_persona_uuid"`
	PursuitUUID        string   `csv:"pursuit_uuid"`
	Title              string   `csv:"title"`
	Status             string   `csv:"status"`
	CvAssetURI         string   `csv:"cv_asset_uri"`
	JdAssetURI         string   `csv:"jd_asset_uri"`
	CvFingerprint      string   `csv:"cv_fingerprint"`
	JdFingerprint      string   `csv:"jd_fingerprint"`
	LlmModel           string   `csv:"llm_model"`
	PromptVersion      string   `csv:"prompt_version"`
	ExtractionPayload  string   `csv:"extraction_payload"`
	GapSummary         string   `csv:"gap_summary"`
	MermaidDiagram     string   `csv:"mermaid_diagram"`
	CreatedAt          SeedTime `csv:"created_at"`
	UpdatedAt          SeedTime `csv:"updated_at"`
}

func (s Seed) SeedSkillGapAnalyses(ctx context.Context) {
	csvPath := "./internal/seeds/seeddata/skill_gap_analyses.csv"
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		fmt.Println("skill_gap_analyses.csv not found, skipping")
		return
	}
	b, err := os.ReadFile(csvPath)
	if err != nil {
		panic(err)
	}
	rows := []skillGapAnalysisSeed{}
	if err := csvutil.Unmarshal(b, &rows); err != nil {
		panic(err)
	}
	if len(rows) == 0 {
		fmt.Println("No skill_gap_analyses to seed")
		return
	}
	var count int
	if err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM bridgr.skill_gap_analyses`).Scan(&count); err != nil {
		panic(err)
	}
	if count > 0 {
		fmt.Printf("skill_gap_analyses already seeded (%d rows), skipping...\n", count)
		return
	}
	const q = `
		INSERT INTO bridgr.skill_gap_analyses (
			uuid, id, user_id, founder_persona_uuid, pursuit_uuid, title, status,
			cv_asset_uri, jd_asset_uri, cv_fingerprint, jd_fingerprint,
			llm_model, prompt_version, extraction_payload, gap_summary, mermaid_diagram,
			created_at, updated_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14::jsonb,$15::jsonb,$16,$17,$18
		)`
	for _, r := range rows {
		ext := r.ExtractionPayload
		if ext == "" {
			ext = "{}"
		}
		gap := r.GapSummary
		if gap == "" {
			gap = "{}"
		}
		if _, err := s.db.Exec(ctx, q,
			r.UUID, r.ID, r.UserID,
			strOrNil(r.FounderPersonaUUID),
			strOrNil(r.PursuitUUID),
			strOrNil(r.Title),
			r.Status,
			strOrNil(r.CvAssetURI),
			strOrNil(r.JdAssetURI),
			strOrNil(r.CvFingerprint),
			strOrNil(r.JdFingerprint),
			strOrNil(r.LlmModel),
			strOrNil(r.PromptVersion),
			ext, gap,
			strOrNil(r.MermaidDiagram),
			r.CreatedAt, r.UpdatedAt,
		); err != nil {
			panic(fmt.Errorf("seed skill_gap_analyses: %w", err))
		}
	}
	fmt.Printf("✓ Seeded %d skill_gap_analyses\n", len(rows))
}
