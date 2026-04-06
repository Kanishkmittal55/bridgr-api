package seeds

import (
	"context"
	"fmt"
	"os"

	"github.com/jszwec/csvutil"
)

type feedItemSeed struct {
	UUID             string   `csv:"uuid"`
	ID               int64    `csv:"id"`
	UserID           int      `csv:"user_id"`
	JobCandidateUUID string   `csv:"job_candidate_uuid"`
	ScoreUUID        string   `csv:"score_uuid"`
	VerificationUUID string   `csv:"verification_uuid"`
	CompositeScore   float64  `csv:"composite_score"`
	GapSeverity      string   `csv:"gap_severity"`
	Title            string   `csv:"title"`
	Company          string   `csv:"company"`
	Location         string   `csv:"location"`
	JobURL           string   `csv:"job_url"`
	MatchSummary     string   `csv:"match_summary"`
	GapSummary       string   `csv:"gap_summary"`
	FeedStatus       string   `csv:"feed_status"`
	SurfacedAt       SeedTime `csv:"surfaced_at"`
	SeenAt           SeedTime `csv:"seen_at"`
	CreatedAt        SeedTime `csv:"created_at"`
	UpdatedAt        SeedTime `csv:"updated_at"`
}

func (s Seed) SeedFeedItems(ctx context.Context) {
	csvPath := "./internal/seeds/seeddata/feed_items.csv"
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		fmt.Println("feed_items.csv not found, skipping")
		return
	}
	b, err := os.ReadFile(csvPath)
	if err != nil {
		panic(err)
	}
	rows := []feedItemSeed{}
	if err := csvutil.Unmarshal(b, &rows); err != nil {
		panic(err)
	}
	if len(rows) == 0 {
		fmt.Println("No feed_items to seed")
		return
	}
	var count int
	if err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM bridgr.feed_items`).Scan(&count); err != nil {
		panic(err)
	}
	if count > 0 {
		fmt.Printf("feed_items already seeded (%d rows), skipping...\n", count)
		return
	}
	const q = `
		INSERT INTO bridgr.feed_items (
			uuid, id, user_id, job_candidate_uuid, score_uuid, verification_uuid,
			composite_score, gap_severity, title, company, location, job_url,
			match_summary, gap_summary, feed_status, surfaced_at, seen_at, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)`
	for _, r := range rows {
		var seen interface{}
		if !r.SeenAt.isNil {
			seen = r.SeenAt.Time
		}
		if _, err := s.db.Exec(ctx, q,
			r.UUID, r.ID, r.UserID, r.JobCandidateUUID, strOrNil(r.ScoreUUID), strOrNil(r.VerificationUUID),
			r.CompositeScore, strOrNil(r.GapSeverity), strOrNil(r.Title), strOrNil(r.Company), strOrNil(r.Location), strOrNil(r.JobURL),
			strOrNil(r.MatchSummary), strOrNil(r.GapSummary), r.FeedStatus, r.SurfacedAt, seen, r.CreatedAt, r.UpdatedAt,
		); err != nil {
			panic(fmt.Errorf("seed feed_items: %w", err))
		}
	}
	fmt.Printf("✓ Seeded %d feed_items\n", len(rows))
}
