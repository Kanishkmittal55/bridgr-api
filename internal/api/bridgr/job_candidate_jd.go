package bridgr

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"strings"

	apierrors "github.com/Kanishkmittal55/bridgr-api/internal/errors"
	"github.com/Kanishkmittal55/bridgr-api/internal/cloud"
	"github.com/Kanishkmittal55/bridgr-api/internal/config"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository/sqlc"
	"github.com/Kanishkmittal55/bridgr-api/internal/uuid"
)

// jdAssetURIFromJobCandidate resolves a stored JD pointer for skill-gap analysis when the client
// omits jd_asset_uri. Order: jd_s3_uri → upload jd_text to S3 → https (or http) job_url.
func (s *server) jdAssetURIFromJobCandidate(ctx context.Context, userID int32, row *sqlc.BridgrJobCandidate) (string, error) {
	if row.UserID != userID {
		return "", fmt.Errorf("%w: job candidate does not belong to this user", apierrors.ErrForbidden)
	}
	if row.JdS3Uri.Valid {
		u := strings.TrimSpace(row.JdS3Uri.String)
		if u != "" {
			if err := validateJdS3URIAgainstConfig(u); err != nil {
				return "", err
			}
			return u, nil
		}
	}
	if row.JdText.Valid {
		t := strings.TrimSpace(row.JdText.String)
		if t != "" {
			return s.uploadJobCandidateJdText(ctx, userID, row, t)
		}
	}
	jobURL := strings.TrimSpace(row.JobUrl)
	if jobURL != "" {
		ok, err := isHTTPOrHTTPSURL(jobURL)
		if err != nil {
			return "", fmt.Errorf("%w: job_url: %v", apierrors.ErrBadRequest, err)
		}
		if ok {
			return jobURL, nil
		}
	}
	return "", fmt.Errorf("%w: job candidate has no usable job description (need jd_s3_uri, non-empty jd_text, or http(s) job_url); configure S3 to store inline jd_text", apierrors.ErrBadRequest)
}

func isHTTPOrHTTPSURL(raw string) (bool, error) {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" || u.Scheme == "" {
		return false, fmt.Errorf("invalid URL")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false, nil
	}
	return true, nil
}

func validateJdS3URIAgainstConfig(raw string) error {
	bucket, _, err := cloud.ParseS3URI(raw)
	if err != nil {
		return fmt.Errorf("%w: jd_s3_uri: %v", apierrors.ErrBadRequest, err)
	}
	cfg := config.Get()
	if cfg.HassleSkipS3Bucket != "" && bucket != cfg.HassleSkipS3Bucket {
		return fmt.Errorf("%w: jd_s3_uri uses a bucket that is not allowed", apierrors.ErrBadRequest)
	}
	return nil
}

func (s *server) uploadJobCandidateJdText(ctx context.Context, userID int32, row *sqlc.BridgrJobCandidate, text string) (string, error) {
	cfg := config.Get()
	bucket := cfg.HassleSkipS3Bucket
	if bucket == "" || s.deps.S3 == nil {
		return "", fmt.Errorf("%w: S3 is required to store inline job description from the feed (configure HASSLE_SKIP_S3_BUCKET and S3)", apierrors.ErrServiceUnavailable)
	}
	if !row.Uuid.Valid {
		return "", fmt.Errorf("%w: candidate uuid", apierrors.ErrInternal)
	}
	uidStr, err := uuid.ToString(row.Uuid.Bytes)
	if err != nil {
		return "", fmt.Errorf("%w: candidate uuid: %w", apierrors.ErrInternal, err)
	}
	key := fmt.Sprintf("inputs/%d/%s/jd-from-feed.txt", userID, uidStr)
	_, err = s.deps.S3.Upload(ctx, bucket, key, bytes.NewReader([]byte(text)))
	if err != nil {
		return "", fmt.Errorf("%w: upload jd: %w", apierrors.ErrInternal, err)
	}
	return cloud.URI(bucket, key), nil
}
