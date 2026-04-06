package bridgr_worker

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/Kanishkmittal55/bridgr-api/internal/cloud"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository/sqlc"
	"github.com/Kanishkmittal55/bridgr-api/internal/uuid"
	guuid "github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

const discoveryFeedPromptVersion = "discovery-feed-p4-v1"

func fingerprintHex(parts ...string) string {
	h := sha256.New()
	for _, p := range parts {
		h.Write([]byte(p))
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))
}

func cvAssetURIFromCanonicalAnalysis(ctx context.Context, repo *repository.Repo, q sqlc.Querier, userID int32, canonicalAnalysisUUID string) (string, error) {
	canonicalAnalysisUUID = strings.TrimSpace(canonicalAnalysisUUID)
	if canonicalAnalysisUUID == "" {
		return "", errors.New("empty canonical_cv_analysis_uuid")
	}
	u, err := guuid.FromString(canonicalAnalysisUUID)
	if err != nil {
		return "", fmt.Errorf("canonical_cv_analysis_uuid: %w", err)
	}
	analysis, err := repo.GetSkillGapAnalysisByUUID(ctx, q, uuid.ToPgUuid(u))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("canonical skill-gap analysis not found")
		}
		return "", err
	}
	if analysis.UserID != userID {
		return "", fmt.Errorf("canonical analysis user_id mismatch")
	}
	if !analysis.CvAssetUri.Valid {
		return "", fmt.Errorf("canonical analysis has no cv_asset_uri")
	}
	s := strings.TrimSpace(analysis.CvAssetUri.String)
	if s == "" {
		return "", fmt.Errorf("canonical analysis has empty cv_asset_uri")
	}
	return s, nil
}

func jdAssetURIForWorker(ctx context.Context, s3Client cloud.Interface, bucket string, userID int32, cand *sqlc.BridgrJobCandidate, jdInline, jobURL string) (string, error) {
	if cand == nil {
		return "", errors.New("nil job candidate")
	}
	if cand.JdS3Uri.Valid {
		u := strings.TrimSpace(cand.JdS3Uri.String)
		if u != "" {
			return u, nil
		}
	}
	jdInline = strings.TrimSpace(jdInline)
	if jdInline != "" && bucket != "" && s3Client != nil {
		if !cand.Uuid.Valid {
			return "", errors.New("candidate uuid missing")
		}
		uidStr, err := uuid.ToString(cand.Uuid.Bytes)
		if err != nil {
			return "", fmt.Errorf("candidate uuid: %w", err)
		}
		key := fmt.Sprintf("inputs/%d/%s/jd-from-discovery-feed.txt", userID, uidStr)
		if _, err := s3Client.Upload(ctx, bucket, key, bytes.NewReader([]byte(jdInline))); err != nil {
			return "", fmt.Errorf("upload jd to s3: %w", err)
		}
		return cloud.URI(bucket, key), nil
	}
	jobURL = strings.TrimSpace(jobURL)
	if jobURL != "" {
		u, err := url.Parse(jobURL)
		if err == nil && u.Host != "" && (u.Scheme == "http" || u.Scheme == "https") {
			return jobURL, nil
		}
	}
	return "", errors.New("no usable jd (need jd_text + S3 bucket, jd_s3_uri, or http(s) job_url)")
}

func feedMatchSummaryText(score jobDiscoveryScoreOutput) string {
	var matched []string
	_ = json.Unmarshal(score.MatchedSkillsJSON, &matched)
	if len(matched) > 0 {
		return fmt.Sprintf("Fit score %.2f. Strong overlap on: %s.", score.CompositeScore, strings.Join(matched, ", "))
	}
	return fmt.Sprintf("Fit score %.2f (discovery ranking).", score.CompositeScore)
}

func feedGapSummaryText(score jobDiscoveryScoreOutput) string {
	var gaps []string
	_ = json.Unmarshal(score.GapSkillsJSON, &gaps)
	if len(gaps) == 0 {
		if score.GapSeverity == "none" || score.GapSeverity == "" {
			return "No major skill gaps flagged vs your must-have stack."
		}
		return fmt.Sprintf("Gap severity: %s.", score.GapSeverity)
	}
	return fmt.Sprintf("Focus areas: %s.", strings.Join(gaps, ", "))
}

// runDiscoveryFeedPipeline creates or reuses a completed skill-gap analysis (CV+JD pointers + structured summary),
// links it to the job candidate, and upserts a feed row keyed by (user_id, job_candidate_uuid).
func (p *Processor) runDiscoveryFeedPipeline(
	ctx context.Context,
	userID int32,
	discoveryRunUUID pgtype.UUID,
	canonicalCvAnalysisUUID string,
	cand *sqlc.BridgrJobCandidate,
	scoreRow *sqlc.BridgrJobScore,
	scoreOut jobDiscoveryScoreOutput,
	jdInline, jobURL string,
	feedLocation string,
	surfacedAt pgtype.Timestamp,
) error {
	cvURI, err := cvAssetURIFromCanonicalAnalysis(ctx, p.repo, p.q, userID, canonicalCvAnalysisUUID)
	if err != nil {
		return fmt.Errorf("cv asset uri: %w", err)
	}
	jdURI, err := jdAssetURIForWorker(ctx, p.s3, strings.TrimSpace(p.s3Bucket), userID, cand, jdInline, jobURL)
	if err != nil {
		return fmt.Errorf("jd asset uri: %w", err)
	}

	cvFP := fingerprintHex(cvURI, canonicalCvAnalysisUUID)
	jdFP := fingerprintHex(jdURI, jdInline, cand.UrlHash)

	gapJSON, extractJSON, analysisTitle := discoveryAnalysisPayloads(userID, discoveryRunUUID, cand, scoreOut)

	existing, ferr := p.repo.GetSkillGapAnalysisByFingerprint(ctx, p.q, sqlc.GetSkillGapAnalysisByFingerprintParams{
		UserID:        userID,
		CvFingerprint: pgtype.Text{String: cvFP, Valid: true},
		JdFingerprint: pgtype.Text{String: jdFP, Valid: true},
	})

	var analysis *sqlc.BridgrSkillGapAnalysis
	switch {
	case ferr == nil:
		analysis = existing
		if _, err := p.repo.UpdateSkillGapAnalysisSummary(ctx, p.q, sqlc.UpdateSkillGapAnalysisSummaryParams{
			ID:             analysis.ID,
			GapSummary:     gapJSON,
			MermaidDiagram: pgtype.Text{},
		}); err != nil {
			return fmt.Errorf("update analysis gap summary: %w", err)
		}
		if strings.TrimSpace(analysis.Status) != "completed" {
			if _, err := p.repo.UpdateSkillGapAnalysisStatus(ctx, p.q, sqlc.UpdateSkillGapAnalysisStatusParams{
				ID:     analysis.ID,
				Status: "completed",
			}); err != nil {
				return fmt.Errorf("update analysis status: %w", err)
			}
		}
	case errors.Is(ferr, pgx.ErrNoRows):
		row, cerr := p.repo.CreateSkillGapAnalysis(ctx, p.q, sqlc.CreateSkillGapAnalysisParams{
			UserID:             userID,
			FounderPersonaUuid: pgtype.UUID{Valid: false},
			PursuitUuid:        pgtype.UUID{Valid: false},
			Title:              pgtype.Text{String: analysisTitle, Valid: true},
			Status:             "completed",
			CvAssetUri:         pgtype.Text{String: cvURI, Valid: true},
			JdAssetUri:         pgtype.Text{String: jdURI, Valid: true},
			CvFingerprint:      pgtype.Text{String: cvFP, Valid: true},
			JdFingerprint:      pgtype.Text{String: jdFP, Valid: true},
			LlmModel:           pgtype.Text{String: scoreOut.ScoringModel, Valid: scoreOut.ScoringModel != ""},
			PromptVersion:      pgtype.Text{String: discoveryFeedPromptVersion, Valid: true},
			ExtractionPayload:  extractJSON,
			GapSummary:         gapJSON,
			MermaidDiagram:     pgtype.Text{},
			ErrorCode:          pgtype.Text{},
			ErrorDetail:        pgtype.Text{},
		})
		if cerr != nil {
			return fmt.Errorf("create skill gap analysis: %w", cerr)
		}
		analysis = row
	default:
		return fmt.Errorf("fingerprint lookup: %w", ferr)
	}

	if analysis == nil {
		return errors.New("analysis nil after create/lookup")
	}

	if err := p.tryCreateAnalysisJobLink(ctx, userID, analysis.Uuid, cand.Uuid); err != nil {
		return err
	}

	ver := pgtype.UUID{Valid: false}
	_, uerr := p.repo.UpsertFeedItemFromDiscovery(ctx, p.q, sqlc.UpsertFeedItemFromDiscoveryParams{
		UserID:           userID,
		JobCandidateUuid: cand.Uuid,
		ScoreUuid:        scoreRow.Uuid,
		VerificationUuid: ver,
		CompositeScore:   scoreRow.CompositeScore,
		GapSeverity:      scoreRow.GapSeverity,
		Title:            cand.Title,
		Company:          cand.Company,
		Location:         pgtype.Text{String: strings.TrimSpace(feedLocation), Valid: strings.TrimSpace(feedLocation) != ""},
		JobUrl:           pgtype.Text{String: jobURL, Valid: jobURL != ""},
		MatchSummary:     pgtype.Text{String: feedMatchSummaryText(scoreOut), Valid: true},
		GapSummary:       pgtype.Text{String: feedGapSummaryText(scoreOut), Valid: true},
		FeedStatus:       "new",
		SurfacedAt:       surfacedAt,
		SeenAt:           pgtype.Timestamp{Valid: false},
	})
	if uerr != nil {
		return fmt.Errorf("upsert feed item: %w", uerr)
	}
	return nil
}

func discoveryAnalysisPayloads(userID int32, discoveryRunUUID pgtype.UUID, cand *sqlc.BridgrJobCandidate, scoreOut jobDiscoveryScoreOutput) (gapSummary []byte, extract []byte, title string) {
	if cand == nil {
		return nil, nil, "Discovery match"
	}
	t := strings.TrimSpace(cand.Title.String)
	co := strings.TrimSpace(cand.Company.String)
	if t != "" && co != "" {
		title = fmt.Sprintf("Discovery: %s @ %s", t, co)
	} else if t != "" {
		title = fmt.Sprintf("Discovery: %s", t)
	} else {
		title = "Discovery job match"
	}

	var matched, gaps []string
	_ = json.Unmarshal(scoreOut.MatchedSkillsJSON, &matched)
	_ = json.Unmarshal(scoreOut.GapSkillsJSON, &gaps)

	gapObj := map[string]interface{}{
		"source":          "discovery",
		"composite_score": scoreOut.CompositeScore,
		"gap_severity":    scoreOut.GapSeverity,
		"matched_skills":  matched,
		"gap_skills":      gaps,
		"scoring_model":   scoreOut.ScoringModel,
		"scoring_version": scoreOut.ScoringVersion,
	}
	gapSummary, _ = json.Marshal(gapObj)

	candStr := ""
	if cand.Uuid.Valid {
		if s, err := uuid.ToString(cand.Uuid.Bytes); err == nil {
			candStr = s
		}
	}
	runStr := ""
	if discoveryRunUUID.Valid {
		if s, err := uuid.ToString(discoveryRunUUID.Bytes); err == nil {
			runStr = s
		}
	}
	extObj := map[string]interface{}{
		"discovery_run_uuid":     runStr,
		"job_candidate_uuid":     candStr,
		"user_id":                userID,
		"url_hash":               cand.UrlHash,
		"skill_match_score":      scoreOut.SkillMatchScore,
		"experience_match_score": scoreOut.ExperienceMatchScore,
	}
	extract, _ = json.Marshal(extObj)
	return gapSummary, extract, title
}

func (p *Processor) tryCreateAnalysisJobLink(ctx context.Context, userID int32, analysisUUID, jobCandidateUUID pgtype.UUID) error {
	_, err := p.repo.CreateAnalysisJobLink(ctx, p.q, sqlc.CreateAnalysisJobLinkParams{
		UserID:           userID,
		AnalysisUuid:     analysisUUID,
		JobCandidateUuid: jobCandidateUUID,
		LinkKind:         "from_job_feed",
	})
	if err == nil {
		return nil
	}
	var pe *pgconn.PgError
	if errors.As(err, &pe) && pe.Code == "23505" {
		return nil
	}
	return fmt.Errorf("analysis_job_link: %w", err)
}
