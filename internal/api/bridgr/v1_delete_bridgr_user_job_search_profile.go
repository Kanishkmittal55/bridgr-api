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

// V1DeleteBridgrUserJobSearchProfile handles DELETE /v1/bridgr/users/{userID}/job-search-profiles/{profileUUID}
func (s *server) V1DeleteBridgrUserJobSearchProfile(w http.ResponseWriter, r *http.Request, userID int32, profileUUID openapi_types.UUID, params types.V1DeleteBridgrUserJobSearchProfileParams) {
	ctx := r.Context()
	err := s.v1DeleteBridgrUserJobSearchProfile(ctx, userID, profileUUID)
	if err != nil {
		s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, nil, err)
		return
	}
	s.deps.ResponseWriter.WriteNoContentResponse(ctx, w, r, nil)
}

func (s *server) v1DeleteBridgrUserJobSearchProfile(ctx context.Context, userID int32, profileUUID openapi_types.UUID) error {
	if err := s.requireStore(); err != nil {
		return err
	}
	pgid, err := uuid.ConvertOapiUUIDToPgUUID(profileUUID)
	if err != nil {
		return fmt.Errorf("%w: invalid profile uuid: %w", apierrors.ErrBadRequest, err)
	}
	_, err = s.deps.Repo.GetJobSearchProfileByUserIDAndUUID(ctx, s.querier(), sqlc.GetJobSearchProfileByUserIDAndUUIDParams{
		UserID: userID,
		Uuid:   pgid,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%w: job search profile not found", apierrors.ErrNotFound)
		}
		return fmt.Errorf("%w: load profile: %w", apierrors.ErrInternal, err)
	}
	if err := s.deps.Repo.DeleteJobSearchProfileByUserIDAndUUID(ctx, s.querier(), sqlc.DeleteJobSearchProfileByUserIDAndUUIDParams{
		UserID: userID,
		Uuid:   pgid,
	}); err != nil {
		return fmt.Errorf("%w: delete profile: %w", apierrors.ErrInternal, err)
	}
	return nil
}
