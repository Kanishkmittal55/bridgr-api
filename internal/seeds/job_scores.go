package seeds

import (
	"context"
	"fmt"
	"os"

	"github.com/jszwec/csvutil"
)

type jobScoreSeed struct {
	UUID                 string   `csv:"uuid"`
	ID                   int64    `csv:"id"`
	UserID               int      `csv:"user_id"`
	JobCandidateUUID     string   `csv:"job_candidate_uuid"`
	EnrichmentUUID       string   `csv:"enrichment_uuid"`
	SkillMatchScore      float64  `csv:"skill_match_score"`
	ExperienceMatchScore float64  `csv:"experience_match_score"`
	LocationMatchScore   float64  `csv:"location_match_score"`
	RecencyScore         float64  `csv:"recency_score"`
	BoardQualityScore    float64  `csv:"board_quality_score"`
	CompositeScore       float64  `csv:"composite_score"`
	MatchedSkills        string   `csv:"matched_skills"`
	GapSkills            string   `csv:"gap_skills"`
	GapSeverity          string   `csv:"gap_severity"`
	ScoringModel         string   `csv:"scoring_model"`
	ScoringVersion       string   `csv:"scoring_version"`
	CreatedAt            SeedTime `csv:"created_at"`
	UpdatedAt            SeedTime `csv:"updated_at"`
}

func (s Seed) SeedJobScores(ctx context.Context) {
	csvPath := "./internal/seeds/seeddata/job_scores.csv"
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		fmt.Println("job_scores.csv not found, skipping")
		return
	}
	b, err := os.ReadFile(csvPath)
	if err != nil {
		panic(err)
	}
	rows := []jobScoreSeed{}
	if err := csvutil.Unmarshal(b, &rows); err != nil {
		panic(err)
	}
	if len(rows) == 0 {
		fmt.Println("No job_scores to seed")
		return
	}
	var count int
	if err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM bridgr.job_scores`).Scan(&count); err != nil {
		panic(err)
	}
	if count > 0 {
		fmt.Printf("job_scores already seeded (%d rows), skipping...\n", count)
		return
	}
	const q = `
		INSERT INTO bridgr.job_scores (
			uuid, id, user_id, job_candidate_uuid, enrichment_uuid,
			skill_match_score, experience_match_score, location_match_score, recency_score, board_quality_score,
			composite_score, matched_skills, gap_skills, gap_severity, scoring_model, scoring_version,
			created_at, updated_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12::jsonb,$13::jsonb,$14,$15,$16,$17,$18
		)`
	for _, r := range rows {
		ms := r.MatchedSkills
		if ms == "" {
			ms = "[]"
		}
		gs := r.GapSkills
		if gs == "" {
			gs = "[]"
		}
		if _, err := s.db.Exec(ctx, q,
			r.UUID, r.ID, r.UserID, r.JobCandidateUUID, strOrNil(r.EnrichmentUUID),
			r.SkillMatchScore, r.ExperienceMatchScore, r.LocationMatchScore, r.RecencyScore, r.BoardQualityScore,
			r.CompositeScore, ms, gs, strOrNil(r.GapSeverity), strOrNil(r.ScoringModel), strOrNil(r.ScoringVersion),
			r.CreatedAt, r.UpdatedAt,
		); err != nil {
			panic(fmt.Errorf("seed job_scores: %w", err))
		}
	}
	fmt.Printf("✓ Seeded %d job_scores\n", len(rows))
}
