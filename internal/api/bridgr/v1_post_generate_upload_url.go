package bridgr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/hassleskip/bridgr-api/internal/config"
	"github.com/hassleskip/bridgr-api/internal/uuid"
	types "github.com/hassleskip/bridgr-api/pkg/types"
	hserr "github.com/hassleskip/hassle-go/pkg/errors"
)

const presignedURLExpiry = 5 * time.Minute

// V1PostGenerateUploadUrl handles POST /v1/autograph/upload-url
func (s *server) V1PostGenerateUploadUrl(w http.ResponseWriter, r *http.Request, _ types.V1PostGenerateUploadUrlParams) {
	ctx := r.Context()
	var payload types.PresignedUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, nil, fmt.Errorf("%w: invalid JSON: %v", hserr.ErrBadRequest, err))
		return
	}
	resp, err := s.postGenerateUploadURL(ctx, payload)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

func (s *server) postGenerateUploadURL(ctx context.Context, params types.PresignedUploadRequest) (*types.PresignedUploadResponse, error) {
	if params.UserId == 0 {
		return nil, fmt.Errorf("%w: user_id is required", hserr.ErrBadRequest)
	}
	if params.Filename == "" {
		return nil, fmt.Errorf("%w: filename is required", hserr.ErrBadRequest)
	}
	if s.deps.S3 == nil {
		return nil, fmt.Errorf("%w: S3 not configured (set S3_URL)", hserr.ErrServiceUnavailable)
	}

	userID := params.UserId
	fileUUID, err := uuid.NewDBUuid()
	if err != nil {
		return nil, fmt.Errorf("%w: file uuid: %w", hserr.ErrInternal, err)
	}
	filename := params.Filename
	s3Key := fmt.Sprintf("inputs/%d/%s/%s", userID, fileUUID.String(), filename)

	contentType := "application/octet-stream"
	if params.ContentType != nil && *params.ContentType != "" {
		contentType = *params.ContentType
	} else {
		switch filepath.Ext(filename) {
		case ".pdf":
			contentType = "application/pdf"
		case ".txt":
			contentType = "text/plain"
		case ".doc", ".docx":
			contentType = "application/msword"
		}
	}

	cfg := config.Get()
	bucket := cfg.HassleSkipS3Bucket
	if bucket == "" {
		return nil, fmt.Errorf("%w: HASSLE_SKIP_S3_BUCKET is required", hserr.ErrInternal)
	}

	uploadURL, err := s.deps.S3.GetPresignedUploadURL(ctx, bucket, s3Key, contentType, presignedURLExpiry)
	if err != nil {
		return nil, fmt.Errorf("%w: presign: %w", hserr.ErrInternal, err)
	}
	expiresAt := time.Now().Add(presignedURLExpiry)
	s3URI := fmt.Sprintf("s3://%s/%s", bucket, s3Key)
	return &types.PresignedUploadResponse{
		UploadUrl: &uploadURL,
		S3Key:     &s3Key,
		S3Uri:     &s3URI,
		ExpiresAt: &expiresAt,
	}, nil
}
