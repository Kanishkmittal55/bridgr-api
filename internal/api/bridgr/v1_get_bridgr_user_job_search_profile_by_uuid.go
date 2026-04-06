package bridgr

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	apierrors "github.com/Kanishkmittal55/bridgr-api/internal/errors"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository/sqlc"
	"github.com/Kanishkmittal55/bridgr-api/internal/uuid"
	types "github.com/Kanishkmittal55/bridgr-api/pkg/types"
	"github.com/jackc/pgx/v5"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// V1GetBridgrUserJobSearchProfileByUUID handles GET /v1/bridgr/users/{userID}/job-search-profiles/{profileUUID}
func (s *server) V1GetBridgrUserJobSearchProfileByUUID(w http.ResponseWriter, r *http.Request, userID int32, profileUUID openapi_types.UUID, params types.V1GetBridgrUserJobSearchProfileByUUIDParams) {
	ctx := r.Context()
	resp, err := s.v1GetBridgrUserJobSearchProfileByUUID(ctx, userID, profileUUID)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

func (s *server) v1GetBridgrUserJobSearchProfileByUUID(ctx context.Context, userID int32, profileUUID openapi_types.UUID) (*types.BridgrJobSearchProfile, error) {
	if err := s.requireStore(); err != nil {
		return nil, err
	}
	pgid, err := uuid.ConvertOapiUUIDToPgUUID(profileUUID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid profile uuid: %w", apierrors.ErrBadRequest, err)
	}
	row, err := s.deps.Repo.GetJobSearchProfileByUserIDAndUUID(ctx, s.querier(), sqlc.GetJobSearchProfileByUserIDAndUUIDParams{
		UserID: userID,
		Uuid:   pgid,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: job search profile not found", apierrors.ErrNotFound)
		}
		return nil, fmt.Errorf("%w: load profile: %w", apierrors.ErrInternal, err)
	}
	out, err := jobSearchProfileFromRow(row)
	if err != nil {
		return nil, fmt.Errorf("%w: map profile: %w", apierrors.ErrInternal, err)
	}
	return &out, nil
}
