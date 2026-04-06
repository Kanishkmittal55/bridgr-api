package bridgr

import (
	"context"
	"encoding/json"
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

// V1PutBridgrUserJobSearchProfileByUUID handles PUT /v1/bridgr/users/{userID}/job-search-profiles/{profileUUID}
func (s *server) V1PutBridgrUserJobSearchProfileByUUID(w http.ResponseWriter, r *http.Request, userID int32, profileUUID openapi_types.UUID, params types.V1PutBridgrUserJobSearchProfileByUUIDParams) {
	ctx := r.Context()
	var body types.UpsertBridgrJobSearchProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, nil, fmt.Errorf("%w: invalid JSON body: %v", apierrors.ErrBadRequest, err))
		return
	}
	resp, err := s.v1PutBridgrUserJobSearchProfileByUUID(ctx, userID, profileUUID, body)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

func (s *server) v1PutBridgrUserJobSearchProfileByUUID(ctx context.Context, userID int32, profileUUID openapi_types.UUID, body types.UpsertBridgrJobSearchProfileRequest) (*types.BridgrJobSearchProfile, error) {
	if err := s.requireStore(); err != nil {
		return nil, err
	}
	pgid, err := uuid.ConvertOapiUUIDToPgUUID(profileUUID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid profile uuid: %w", apierrors.ErrBadRequest, err)
	}
	existing, err := s.deps.Repo.GetJobSearchProfileByUserIDAndUUID(ctx, s.querier(), sqlc.GetJobSearchProfileByUserIDAndUUIDParams{
		UserID: userID,
		Uuid:   pgid,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: job search profile not found", apierrors.ErrNotFound)
		}
		return nil, fmt.Errorf("%w: load profile: %w", apierrors.ErrInternal, err)
	}
	up, err := profileToUpdateByUserIDAndUUIDParams(existing, body)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err)
	}
	if err := ensureUniqueJobSearchProfile(ctx, s.querier(), userID, &existing.Uuid, naturalKeyFromUpdate(up)); err != nil {
		return nil, err
	}
	row, err := s.deps.Repo.UpdateJobSearchProfileByUserIDAndUUID(ctx, s.querier(), up)
	if err != nil {
		return nil, fmt.Errorf("%w: update profile: %w", apierrors.ErrInternal, err)
	}
	out, err := jobSearchProfileFromRow(row)
	if err != nil {
		return nil, fmt.Errorf("%w: map profile: %w", apierrors.ErrInternal, err)
	}
	return &out, nil
}
