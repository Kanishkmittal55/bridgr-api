package bridgr_worker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path"
	"strings"
	"unicode/utf8"

	"github.com/Kanishkmittal55/bridgr-api/internal/cloud"
	"github.com/Kanishkmittal55/bridgr-api/internal/radar"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository/sqlc"
	"github.com/Kanishkmittal55/bridgr-api/internal/uuid"
	guuid "github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
)

const maxCVPdfBytesForExtraction = 20 << 20 // 20 MiB safety cap for in-memory read

var errNoCanonicalCvURI = errors.New("canonical cv analysis has no cv_asset_uri or extraction text")

// resolveResumeTextForDiscovery loads CV text for FindJobs: prefers stored extraction_payload,
// otherwise downloads cv_asset_uri from S3 and runs Radar PdfExtractionService.ExtractText for PDFs,
// or treats bytes as UTF-8 text for non-PDF assets.
func resolveResumeTextForDiscovery(
	ctx context.Context,
	repo *repository.Repo,
	q sqlc.Querier,
	radarClient *radar.Client,
	s3Client cloud.Interface,
	expectedUserID int32,
	canonicalCvAnalysisUUID string,
) (string, error) {
	canonicalCvAnalysisUUID = strings.TrimSpace(canonicalCvAnalysisUUID)
	if canonicalCvAnalysisUUID == "" {
		return "", fmt.Errorf("canonical_cv_analysis_uuid is empty: %w", errNoCanonicalCvURI)
	}

	u, err := guuid.FromString(canonicalCvAnalysisUUID)
	if err != nil {
		return "", fmt.Errorf("canonical_cv_analysis_uuid: %w", err)
	}
	analysis, err := repo.GetSkillGapAnalysisByUUID(ctx, q, uuid.ToPgUuid(u))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("canonical skill-gap analysis not found: %w", err)
		}
		return "", fmt.Errorf("load canonical analysis: %w", err)
	}
	if analysis.UserID != expectedUserID {
		return "", fmt.Errorf("canonical analysis user_id %d does not match discovery run user_id %d", analysis.UserID, expectedUserID)
	}

	if txt, ok := resumeTextFromExtractionPayload(analysis.ExtractionPayload); ok {
		return strings.TrimSpace(txt), nil
	}

	uri := strings.TrimSpace(analysis.CvAssetUri.String)
	if uri == "" || !analysis.CvAssetUri.Valid {
		return "", errNoCanonicalCvURI
	}
	if s3Client == nil {
		return "", fmt.Errorf("s3 client is nil; cannot download CV from %s", uri)
	}

	bucket, key, err := cloud.ParseS3URI(uri)
	if err != nil {
		return "", fmt.Errorf("cv_asset_uri: %w", err)
	}

	rc, err := s3Client.Download(ctx, bucket, key)
	if err != nil {
		return "", fmt.Errorf("download cv from s3: %w", err)
	}
	defer rc.Close()

	body, err := io.ReadAll(io.LimitReader(rc, maxCVPdfBytesForExtraction+1))
	if err != nil {
		return "", fmt.Errorf("read cv object: %w", err)
	}
	if len(body) > maxCVPdfBytesForExtraction {
		return "", fmt.Errorf("cv object exceeds max size (%d bytes)", maxCVPdfBytesForExtraction)
	}

	if isLikelyPDF(body) {
		if radarClient == nil {
			return "", fmt.Errorf("radar client is nil; cannot extract text from PDF")
		}
		filename := path.Base(key)
		if filename == "" || filename == "." {
			filename = "cv.pdf"
		}
		text, err := radarClient.ExtractText(ctx, body, filename)
		if err != nil {
			return "", fmt.Errorf("radar pdf extract: %w", err)
		}
		if strings.TrimSpace(text) == "" {
			return "", fmt.Errorf("radar pdf extract returned empty text")
		}
		return strings.TrimSpace(text), nil
	}

	if !utf8.Valid(body) {
		return "", fmt.Errorf("cv asset is not valid UTF-8 and not a PDF")
	}
	s := strings.TrimSpace(string(body))
	if s == "" {
		return "", fmt.Errorf("cv text from object is empty")
	}
	return s, nil
}

func resumeTextFromExtractionPayload(payload []byte) (string, bool) {
	if len(payload) == 0 {
		return "", false
	}
	var m map[string]interface{}
	if err := json.Unmarshal(payload, &m); err != nil || len(m) == 0 {
		return "", false
	}
	keys := []string{
		"resume_text",
		"cv_text",
		"extracted_resume_text",
		"full_text",
		"raw_cv_text",
		"raw_text",
	}
	for _, k := range keys {
		v, ok := m[k]
		if !ok || v == nil {
			continue
		}
		switch t := v.(type) {
		case string:
			if s := strings.TrimSpace(t); s != "" {
				return s, true
			}
		}
	}
	return "", false
}

func isLikelyPDF(b []byte) bool {
	return len(b) >= 4 && bytes.Equal(b[:4], []byte("%PDF"))
}
