package bridgr

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	apierrors "github.com/Kanishkmittal55/bridgr-api/internal/apierrors"
	"github.com/Kanishkmittal55/bridgr-api/internal/cloud"
	"github.com/Kanishkmittal55/bridgr-api/internal/config"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository/sqlc"
	types "github.com/Kanishkmittal55/bridgr-api/pkg/types"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

const presignedAssetReadExpiry = 15 * time.Minute

// V1GetBridgrAnalysisCvAssetReadUrl handles GET /v1/bridgr/analyses/{analysisUUID}/assets/cv/read-url
func (s *server) V1GetBridgrAnalysisCvAssetReadUrl(w http.ResponseWriter, r *http.Request, analysisUUID openapi_types.UUID, _ types.V1GetBridgrAnalysisCvAssetReadUrlParams) {
	ctx := r.Context()
	resp, err := s.v1GetAnalysisAssetReadURL(ctx, analysisUUID, false)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

// V1GetBridgrAnalysisJdAssetReadUrl handles GET /v1/bridgr/analyses/{analysisUUID}/assets/jd/read-url
func (s *server) V1GetBridgrAnalysisJdAssetReadUrl(w http.ResponseWriter, r *http.Request, analysisUUID openapi_types.UUID, _ types.V1GetBridgrAnalysisJdAssetReadUrlParams) {
	ctx := r.Context()
	resp, err := s.v1GetAnalysisAssetReadURL(ctx, analysisUUID, true)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

func (s *server) v1GetAnalysisAssetReadURL(ctx context.Context, analysisUUID openapi_types.UUID, jd bool) (*types.BridgrAnalysisAssetReadUrlResponse, error) {
	_, row, err := s.loadAnalysis(ctx, analysisUUID)
	if err != nil {
		return nil, err
	}
	raw := assetURIFromRow(row, jd)
	if strings.TrimSpace(raw) == "" {
		return nil, fmt.Errorf("%w: %s", apierrors.ErrNotFound, assetMissingMsg(jd))
	}
	return s.resolveAssetReadURL(ctx, raw)
}

func assetMissingMsg(jd bool) string {
	logicalName, abbreviation := "CV", "cv"
	if jd {
		logicalName, abbreviation = "Job description", "jd"
	}
	// Keeps error text greppable for clients and humans
	return fmt.Sprintf("%s asset is not set for this analysis (empty %s_asset_uri)", logicalName, abbreviation)
}

func assetURIFromRow(row *sqlc.BridgrSkillGapAnalysis, jd bool) string {
	if jd {
		if row.JdAssetUri.Valid {
			return row.JdAssetUri.String
		}
		return ""
	}
	if row.CvAssetUri.Valid {
		return row.CvAssetUri.String
	}
	return ""
}

func (s *server) resolveAssetReadURL(ctx context.Context, raw string) (*types.BridgrAnalysisAssetReadUrlResponse, error) {
	uri := strings.TrimSpace(raw)
	if uri == "" {
		return nil, fmt.Errorf("%w: empty asset uri", apierrors.ErrNotFound)
	}

	if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
		parsed, err := url.Parse(uri)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
			return nil, fmt.Errorf("%w: invalid HTTP(S) asset URL", apierrors.ErrBadRequest)
		}
		assetCopy := uri
		ct := contentTypeHintForURL(uri)
		out := &types.BridgrAnalysisAssetReadUrlResponse{
			Access:   types.Direct,
			Url:      uri,
			AssetUri: &assetCopy,
		}
		if ct != "" {
			out.ContentType = &ct
		}
		return out, nil
	}

	if !strings.HasPrefix(uri, "s3://") {
		return nil, fmt.Errorf("%w: asset URI must be s3:// or http(s)://", apierrors.ErrBadRequest)
	}

	bucket, key, err := cloud.ParseS3URI(uri)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err)
	}

	cfg := config.Get()
	if cfg.HassleSkipS3Bucket != "" && bucket != cfg.HassleSkipS3Bucket {
		return nil, fmt.Errorf("%w: asset bucket is not allowed", apierrors.ErrBadRequest)
	}

	if s.deps.S3 == nil {
		return nil, fmt.Errorf("%w: S3 not configured", apierrors.ErrServiceUnavailable)
	}

	ct := cloud.ContentTypeByExtension(key)
	signed, err := s.deps.S3.GetPresignedDownloadURL(ctx, bucket, key, presignedAssetReadExpiry, ct)
	if err != nil {
		return nil, fmt.Errorf("%w: presign asset read: %w", apierrors.ErrInternal, err)
	}

	exp := time.Now().Add(presignedAssetReadExpiry)
	assetCopy := uri
	return &types.BridgrAnalysisAssetReadUrlResponse{
		Access:      types.Presigned,
		Url:         signed,
		AssetUri:    &assetCopy,
		ContentType: &ct,
		ExpiresAt:   &exp,
	}, nil
}

func contentTypeHintForURL(uri string) string {
	u := strings.ToLower(uri)
	switch {
	case strings.HasSuffix(u, ".pdf"):
		return "application/pdf"
	case strings.HasSuffix(u, ".txt"):
		return "text/plain"
	case strings.HasSuffix(u, ".html"), strings.HasSuffix(u, ".htm"):
		return "text/html"
	default:
		return ""
	}
}
