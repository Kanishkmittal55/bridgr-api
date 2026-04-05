package bridgr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	apierrors "github.com/Kanishkmittal55/bridgr-api/internal/apierrors"
	"github.com/Kanishkmittal55/bridgr-api/internal/config"
	"github.com/Kanishkmittal55/bridgr-api/internal/logger"
	"github.com/Kanishkmittal55/bridgr-api/internal/radar"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository/sqlc"
	"github.com/Kanishkmittal55/bridgr-api/internal/uuid"
	types "github.com/Kanishkmittal55/bridgr-api/pkg/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// V1PostBridgrUserJobDiscoveryRunCancel handles POST /v1/bridgr/users/{userID}/job-discovery/runs/{runUUID}/cancel
func (s *server) V1PostBridgrUserJobDiscoveryRunCancel(w http.ResponseWriter, r *http.Request, userID int32, runUUID openapi_types.UUID, params types.V1PostBridgrUserJobDiscoveryRunCancelParams) {
	ctx := r.Context()
	_ = params
	out, cerr := s.cancelJobDiscoveryRun(ctx, userID, runUUID.String())
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, out, cerr)
}

func (s *server) cancelJobDiscoveryRun(ctx context.Context, userID int32, runUUIDStr string) (*types.BridgrJobSearchDiscoveryRun, error) {
	if err := s.requireStore(); err != nil {
		return nil, err
	}
	runUID, err := uuid.FromString(runUUIDStr)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid run UUID", apierrors.ErrBadRequest)
	}
	pgid := uuid.ToPgUuid(runUID)

	row, err := s.deps.Repo.GetJobSearchDiscoveryRunByUUID(ctx, s.querier(), pgid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: discovery run not found", apierrors.ErrNotFound)
		}
		return nil, fmt.Errorf("%w: load run: %w", apierrors.ErrInternal, err)
	}
	if row.UserID != userID {
		return nil, fmt.Errorf("%w: run does not belong to this user", apierrors.ErrForbidden)
	}

	switch row.Status {
	case "completed", "failed", "cancelled":
		out, merr := discoveryRunFromRow(row)
		if merr != nil {
			return nil, fmt.Errorf("%w: map run: %w", apierrors.ErrInternal, merr)
		}
		return &out, nil
	}

	cfg := config.Get()
	if cfg.RadarAddr != "" {
		rc, rerr := radar.NewClient(ctx, radar.Config{Addr: cfg.RadarAddr})
		if rerr != nil {
			logger.Get(ctx).Warnw("job discovery cancel: Radar client", "error", rerr)
		} else {
			defer func() { _ = rc.Close() }()
			ok, cerr := rc.CancelCrawl(ctx, runUUIDStr)
			logger.Get(ctx).Infow("job discovery cancel: CancelCrawl", "run_uuid", runUUIDStr, "signalled", ok, "err", cerr)
		}
	}

	now := pgtype.Timestamp{Time: time.Now().UTC(), Valid: true}
	meta, _ := json.Marshal(map[string]interface{}{
		"cancelled":    true,
		"cancelled_at": now.Time.Format(time.RFC3339Nano),
	})
	updated, uerr := s.deps.Repo.UpdateJobSearchDiscoveryRunFinished(ctx, s.querier(), sqlc.UpdateJobSearchDiscoveryRunFinishedParams{
		ID:                row.ID,
		Status:            "cancelled",
		CompletedAt:       now,
		RawCandidateCount: row.RawCandidateCount,
		NewCandidateCount: row.NewCandidateCount,
		RadarMeta:         meta,
		ErrorCode:         pgtype.Text{String: "cancelled", Valid: true},
		ErrorDetail:       pgtype.Text{String: "cancelled by user", Valid: true},
	})
	if uerr != nil {
		return nil, fmt.Errorf("%w: update run: %w", apierrors.ErrInternal, uerr)
	}
	out, merr := discoveryRunFromRow(updated)
	if merr != nil {
		return nil, fmt.Errorf("%w: map run: %w", apierrors.ErrInternal, merr)
	}
	logger.Get(ctx).Infow("bridgr job-discovery run cancelled", "user_id", userID, "run_uuid", runUUIDStr)
	return &out, nil
}
