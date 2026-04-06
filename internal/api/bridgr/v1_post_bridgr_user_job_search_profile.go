package bridgr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	apierrors "github.com/Kanishkmittal55/bridgr-api/internal/errors"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository/sqlc"
	types "github.com/Kanishkmittal55/bridgr-api/pkg/types"
	"github.com/jackc/pgx/v5/pgtype"
)

// V1CreateBridgrUserJobSearchProfile handles POST /v1/bridgr/users/{userID}/job-search-profiles
func (s *server) V1CreateBridgrUserJobSearchProfile(w http.ResponseWriter, r *http.Request, userID int32, params types.V1CreateBridgrUserJobSearchProfileParams) {
	ctx := r.Context()
	var body types.CreateBridgrJobSearchProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, nil, fmt.Errorf("%w: invalid JSON body: %v", apierrors.ErrBadRequest, err))
		return
	}
	resp, err := s.v1CreateBridgrUserJobSearchProfile(ctx, userID, body)
	if err != nil {
		s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, nil, err)
		return
	}
	s.writeCreated(w, r, resp)
}

func (s *server) v1CreateBridgrUserJobSearchProfile(ctx context.Context, userID int32, body types.CreateBridgrJobSearchProfileRequest) (*types.BridgrJobSearchProfile, error) {
	if err := s.requireStore(); err != nil {
		return nil, err
	}
	cp, err := profileToCreateParams(userID, body)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err)
	}
	if err := ensureUniqueJobSearchProfile(ctx, s.querier(), userID, nil, naturalKeyFromCreate(cp)); err != nil {
		return nil, err
	}
	row, err := s.deps.Repo.CreateJobSearchProfile(ctx, s.querier(), cp)
	if err != nil {
		return nil, fmt.Errorf("%w: create profile: %w", apierrors.ErrInternal, err)
	}
	out, err := jobSearchProfileFromRow(row)
	if err != nil {
		return nil, fmt.Errorf("%w: map profile: %w", apierrors.ErrInternal, err)
	}
	return &out, nil
}

// ensureUniqueJobSearchProfile returns ErrConflict if another profile for this user
// already matches the same natural key (role, location, board, goals, stack, CV link).
// excludeProfileUUID skips that row (use when updating so the row being edited does not match itself); nil for create.
func ensureUniqueJobSearchProfile(
	ctx context.Context,
	q sqlc.Querier,
	userID int32,
	excludeProfileUUID *pgtype.UUID,
	candidate jobSearchProfileNaturalKey,
) error {
	rows, err := q.ListJobSearchProfilesByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("%w: list profiles: %w", apierrors.ErrInternal, err)
	}
	for i := range rows {
		if excludeProfileUUID != nil && excludeProfileUUID.Valid && rows[i].Uuid.Valid &&
			bytes.Equal(rows[i].Uuid.Bytes[:], excludeProfileUUID.Bytes[:]) {
			continue
		}
		if naturalKeyFromRow(&rows[i]).equals(candidate) {
			return fmt.Errorf("%w: a profile with the same target role, location, source board, goals, and CV analysis already exists", apierrors.ErrConflict)
		}
	}
	return nil
}
