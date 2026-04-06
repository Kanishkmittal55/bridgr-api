package seeds

import (
	"context"
	"fmt"
	"os"

	"github.com/jszwec/csvutil"
)

type jobCandidateSeed struct {
	UUID             string   `csv:"uuid"`
	ID               int64    `csv:"id"`
	UserID           int      `csv:"user_id"`
	DiscoveryRunUUID string   `csv:"discovery_run_uuid"`
	SourceBoard      string   `csv:"source_board"`
	SourceJobID      string   `csv:"source_job_id"`
	JobURL           string   `csv:"job_url"`
	URLHash          string   `csv:"url_hash"`
	ContentHash      string   `csv:"content_hash"`
	Title            string   `csv:"title"`
	Company          string   `csv:"company"`
	Location         string   `csv:"location"`
	JdText           string   `csv:"jd_text"`
	JdS3URI          string   `csv:"jd_s3_uri"`
	FetchedAt        SeedTime `csv:"fetched_at"`
	IngestionStatus  string   `csv:"ingestion_status"`
	RadarPayload     string   `csv:"radar_payload"`
	ApplicationURL   string   `csv:"application_url"`
	CreatedAt        SeedTime `csv:"created_at"`
	UpdatedAt        SeedTime `csv:"updated_at"`
}

func (s Seed) SeedJobCandidates(ctx context.Context) {
	csvPath := "./internal/seeds/seeddata/job_candidates.csv"
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		fmt.Println("job_candidates.csv not found, skipping")
		return
	}
	b, err := os.ReadFile(csvPath)
	if err != nil {
		panic(err)
	}
	rows := []jobCandidateSeed{}
	if err := csvutil.Unmarshal(b, &rows); err != nil {
		panic(err)
	}
	if len(rows) == 0 {
		fmt.Println("No job_candidates to seed")
		return
	}
	var count int
	if err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM bridgr.job_candidates`).Scan(&count); err != nil {
		panic(err)
	}
	if count > 0 {
		fmt.Printf("job_candidates already seeded (%d rows), skipping...\n", count)
		return
	}
	const q = `
		INSERT INTO bridgr.job_candidates (
			uuid, id, user_id, discovery_run_uuid, source_board, source_job_id, job_url, url_hash,
			content_hash, title, company, location, jd_text, jd_s3_uri, fetched_at, ingestion_status,
			radar_payload, application_url, created_at, updated_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17::jsonb,$18,$19,$20
		)`
	for _, r := range rows {
		rp := r.RadarPayload
		if rp == "" {
			rp = "{}"
		}
		var fetched interface{}
		if !r.FetchedAt.isNil {
			fetched = r.FetchedAt.Time
		}
		if _, err := s.db.Exec(ctx, q,
			r.UUID, r.ID, r.UserID, strOrNil(r.DiscoveryRunUUID), r.SourceBoard,
			strOrNil(r.SourceJobID), r.JobURL, r.URLHash, strOrNil(r.ContentHash),
			strOrNil(r.Title), strOrNil(r.Company), strOrNil(r.Location),
			strOrNil(r.JdText), strOrNil(r.JdS3URI), fetched, r.IngestionStatus,
			rp, strOrNil(r.ApplicationURL), r.CreatedAt, r.UpdatedAt,
		); err != nil {
			panic(fmt.Errorf("seed job_candidates: %w", err))
		}
	}
	fmt.Printf("✓ Seeded %d job_candidates\n", len(rows))
}
