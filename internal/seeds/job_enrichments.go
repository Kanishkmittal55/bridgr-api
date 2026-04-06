package seeds

import (
	"context"
	"fmt"
	"os"

	"github.com/jszwec/csvutil"
)

type jobEnrichmentSeed struct {
	UUID                string   `csv:"uuid"`
	ID                  int64    `csv:"id"`
	UserID              int      `csv:"user_id"`
	JobCandidateUUID    string   `csv:"job_candidate_uuid"`
	Status              string   `csv:"status"`
	RequiredSkills      string   `csv:"required_skills"`
	ExperienceRange     string   `csv:"experience_range"`
	SalaryRange         string   `csv:"salary_range"`
	RemotePolicy        string   `csv:"remote_policy"`
	VisaSponsorship     string   `csv:"visa_sponsorship"`
	CompanySize         string   `csv:"company_size"`
	Industry            string   `csv:"industry"`
	ApplicationDeadline SeedTime `csv:"application_deadline"`
	StructuredJd        string   `csv:"structured_jd"`
	LlmModel            string   `csv:"llm_model"`
	PromptVersion       string   `csv:"prompt_version"`
	ErrorDetail         string   `csv:"error_detail"`
	CreatedAt           SeedTime `csv:"created_at"`
	UpdatedAt           SeedTime `csv:"updated_at"`
}

func (s Seed) SeedJobEnrichments(ctx context.Context) {
	csvPath := "./internal/seeds/seeddata/job_enrichments.csv"
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		fmt.Println("job_enrichments.csv not found, skipping")
		return
	}
	b, err := os.ReadFile(csvPath)
	if err != nil {
		panic(err)
	}
	rows := []jobEnrichmentSeed{}
	if err := csvutil.Unmarshal(b, &rows); err != nil {
		panic(err)
	}
	if len(rows) == 0 {
		fmt.Println("No job_enrichments to seed")
		return
	}
	var count int
	if err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM bridgr.job_enrichments`).Scan(&count); err != nil {
		panic(err)
	}
	if count > 0 {
		fmt.Printf("job_enrichments already seeded (%d rows), skipping...\n", count)
		return
	}
	const q = `
		INSERT INTO bridgr.job_enrichments (
			uuid, id, user_id, job_candidate_uuid, status,
			required_skills, experience_range, salary_range, remote_policy, visa_sponsorship,
			company_size, industry, application_deadline, structured_jd, llm_model, prompt_version,
			error_detail, created_at, updated_at
		) VALUES (
			$1,$2,$3,$4,$5,$6::jsonb,$7::jsonb,$8::jsonb,$9,$10,$11,$12,$13,$14::jsonb,$15,$16,$17,$18,$19
		)`
	for _, r := range rows {
		rs := r.RequiredSkills
		if rs == "" {
			rs = "[]"
		}
		er := r.ExperienceRange
		if er == "" {
			er = "{}"
		}
		sr := r.SalaryRange
		if sr == "" {
			sr = "{}"
		}
		sj := r.StructuredJd
		if sj == "" {
			sj = "{}"
		}
		var visa interface{}
		if r.VisaSponsorship != "" {
			switch r.VisaSponsorship {
			case "true", "1":
				visa = true
			case "false", "0":
				visa = false
			}
		}
		var deadline interface{}
		if !r.ApplicationDeadline.isNil {
			deadline = r.ApplicationDeadline.Time
		}
		if _, err := s.db.Exec(ctx, q,
			r.UUID, r.ID, r.UserID, r.JobCandidateUUID, r.Status,
			rs, er, sr, strOrNil(r.RemotePolicy), visa,
			strOrNil(r.CompanySize), strOrNil(r.Industry), deadline, sj,
			strOrNil(r.LlmModel), strOrNil(r.PromptVersion),
			strOrNil(r.ErrorDetail), r.CreatedAt, r.UpdatedAt,
		); err != nil {
			panic(fmt.Errorf("seed job_enrichments: %w", err))
		}
	}
	fmt.Printf("✓ Seeded %d job_enrichments\n", len(rows))
}
