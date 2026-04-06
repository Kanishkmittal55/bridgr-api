package seeds

import (
	"context"
	"fmt"
	"os"

	"github.com/jszwec/csvutil"
)

type jobNotificationSeed struct {
	UUID             string   `csv:"uuid"`
	ID               int64    `csv:"id"`
	UserID           int      `csv:"user_id"`
	JobCandidateUUID string   `csv:"job_candidate_uuid"`
	Channel          string   `csv:"channel"`
	Status           string   `csv:"status"`
	Payload          string   `csv:"payload"`
	SentAt           SeedTime `csv:"sent_at"`
	SeenAt           SeedTime `csv:"seen_at"`
	ErrorDetail      string   `csv:"error_detail"`
	CreatedAt        SeedTime `csv:"created_at"`
	UpdatedAt        SeedTime `csv:"updated_at"`
}

func (s Seed) SeedJobNotifications(ctx context.Context) {
	csvPath := "./internal/seeds/seeddata/job_notifications.csv"
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		fmt.Println("job_notifications.csv not found, skipping")
		return
	}
	b, err := os.ReadFile(csvPath)
	if err != nil {
		panic(err)
	}
	rows := []jobNotificationSeed{}
	if err := csvutil.Unmarshal(b, &rows); err != nil {
		panic(err)
	}
	if len(rows) == 0 {
		fmt.Println("No job_notifications to seed")
		return
	}
	var count int
	if err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM bridgr.job_notifications`).Scan(&count); err != nil {
		panic(err)
	}
	if count > 0 {
		fmt.Printf("job_notifications already seeded (%d rows), skipping...\n", count)
		return
	}
	const q = `
		INSERT INTO bridgr.job_notifications (
			uuid, id, user_id, job_candidate_uuid, channel, status, payload,
			sent_at, seen_at, error_detail, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7::jsonb,$8,$9,$10,$11,$12)`
	for _, r := range rows {
		pl := r.Payload
		if pl == "" {
			pl = "{}"
		}
		var sent, seen interface{}
		if !r.SentAt.isNil {
			sent = r.SentAt.Time
		}
		if !r.SeenAt.isNil {
			seen = r.SeenAt.Time
		}
		if _, err := s.db.Exec(ctx, q,
			r.UUID, r.ID, r.UserID, r.JobCandidateUUID, r.Channel, r.Status, pl,
			sent, seen, strOrNil(r.ErrorDetail), r.CreatedAt, r.UpdatedAt,
		); err != nil {
			panic(fmt.Errorf("seed job_notifications: %w", err))
		}
	}
	fmt.Printf("✓ Seeded %d job_notifications\n", len(rows))
}
