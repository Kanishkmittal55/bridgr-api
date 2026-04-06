package bridgr

import (
	"context"
	"fmt"
	"net/http"

	apierrors "github.com/Kanishkmittal55/bridgr-api/internal/errors"
	types "github.com/Kanishkmittal55/bridgr-api/pkg/types"
)

// V1ListBridgrUserJobSearchProfiles handles GET /v1/bridgr/users/{userID}/job-search-profiles
func (s *server) V1ListBridgrUserJobSearchProfiles(w http.ResponseWriter, r *http.Request, userID int32, params types.V1ListBridgrUserJobSearchProfilesParams) {
	ctx := r.Context()
	resp, err := s.v1ListBridgrUserJobSearchProfiles(ctx, userID)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

func (s *server) v1ListBridgrUserJobSearchProfiles(ctx context.Context, userID int32) (*types.BridgrJobSearchProfileListResponse, error) {
	if err := s.requireStore(); err != nil {
		return nil, err
	}
	rows, err := s.deps.Repo.ListJobSearchProfilesByUserID(ctx, s.querier(), userID)
	if err != nil {
		return nil, fmt.Errorf("%w: list profiles: %w", apierrors.ErrInternal, err)
	}
	out := types.BridgrJobSearchProfileListResponse{Items: make([]types.BridgrJobSearchProfile, 0, len(rows))}
	for i := range rows {
		item, err := jobSearchProfileFromRow(&rows[i])
		if err != nil {
			return nil, fmt.Errorf("%w: map profile: %w", apierrors.ErrInternal, err)
		}
		out.Items = append(out.Items, item)
	}
	return &out, nil
}
