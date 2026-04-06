package bridgr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Kanishkmittal55/bridgr-api/internal/bridgr_worker"
	"github.com/Kanishkmittal55/bridgr-api/internal/cloud"
	"github.com/Kanishkmittal55/bridgr-api/internal/config"
	apierrors "github.com/Kanishkmittal55/bridgr-api/internal/errors"
	"github.com/Kanishkmittal55/bridgr-api/internal/logger"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository/sqlc"
	"github.com/Kanishkmittal55/bridgr-api/internal/uuid"
	types "github.com/Kanishkmittal55/bridgr-api/pkg/types"
	guuid "github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// V1PostBridgrUserJobDiscoveryRuns handles POST /v1/bridgr/users/{userID}/job-discovery/runs
func (s *server) V1PostBridgrUserJobDiscoveryRuns(w http.ResponseWriter, r *http.Request, userID int32, params types.V1PostBridgrUserJobDiscoveryRunsParams) {
	ctx := r.Context()
	var body types.CreateBridgrJobSearchDiscoveryRunRequest
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil && !errors.Is(err, io.EOF) {
			s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, nil, fmt.Errorf("%w: invalid JSON body: %v", apierrors.ErrBadRequest, err))
			return
		}
	}
	resp, err := s.v1PostBridgrUserJobDiscoveryRuns(ctx, userID, body)
	if err != nil {
		s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, nil, err)
		return
	}
	s.writeCreated(w, r, resp)
}

func (s *server) v1PostBridgrUserJobDiscoveryRuns(ctx context.Context, userID int32, body types.CreateBridgrJobSearchDiscoveryRunRequest) (*types.BridgrJobSearchDiscoveryRun, error) {
	if err := s.requireStore(); err != nil {
		return nil, err
	}
	cfg := config.Get()

	profiles, err := s.deps.Repo.ListJobSearchProfilesByUserID(ctx, s.querier(), userID)
	if err != nil {
		return nil, fmt.Errorf("%w: list profiles: %w", apierrors.ErrInternal, err)
	}
	var profPtr *sqlc.BridgrJobSearchProfile
	if len(profiles) > 0 {
		profPtr = &profiles[0]
	}
	var overrides map[string]interface{}
	if body.RequestParams != nil {
		overrides = *body.RequestParams
	}
	reqParams, err := bridgr_worker.BuildDiscoveryRequestParams(userID, profPtr, overrides)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err)
	}

	runPtr, err := s.deps.Repo.CreateJobSearchDiscoveryRun(ctx, s.querier(), sqlc.CreateJobSearchDiscoveryRunParams{
		UserID:            userID,
		Status:            "pending",
		RequestParams:     reqParams,
		RadarMeta:         []byte("{}"),
		RawCandidateCount: 0,
		NewCandidateCount: 0,
		StartedAt:         pgtype.Timestamp{},
		CompletedAt:       pgtype.Timestamp{},
		ErrorCode:         pgtype.Text{},
		ErrorDetail:       pgtype.Text{},
		SqsMessageID:      pgtype.Text{},
	})
	if err != nil {
		return nil, fmt.Errorf("%w: create run: %w", apierrors.ErrInternal, err)
	}

	runUUID, uerr := guuid.FromBytes(runPtr.Uuid.Bytes[:])
	if uerr != nil {
		return nil, fmt.Errorf("%w: run uuid: %w", apierrors.ErrInternal, uerr)
	}

	if cfg.BridgrQueueURL != "" && s.deps.SQSClient != nil {
		msgID, err := cloud.EnqueueJobDiscovery(ctx, s.deps.SQSClient, cfg.BridgrQueueURL, runUUID, userID)
		if err != nil {
			return nil, fmt.Errorf("%w: enqueue job discovery: %w", apierrors.ErrInternal, err)
		}
		updated, err := s.deps.Repo.PatchJobSearchDiscoveryRun(ctx, s.querier(), sqlc.PatchJobSearchDiscoveryRunParams{
			ID:                runPtr.ID,
			Status:            "queued",
			StartedAt:         runPtr.StartedAt,
			CompletedAt:       runPtr.CompletedAt,
			RawCandidateCount: runPtr.RawCandidateCount,
			NewCandidateCount: runPtr.NewCandidateCount,
			RadarMeta:         runPtr.RadarMeta,
			ErrorCode:         runPtr.ErrorCode,
			ErrorDetail:       runPtr.ErrorDetail,
			SqsMessageID:      pgtype.Text{String: msgID, Valid: msgID != ""},
		})
		if err != nil {
			return nil, fmt.Errorf("%w: persist sqs id: %w", apierrors.ErrInternal, err)
		}
		runPtr = updated
	}

	logger.Get(ctx).Infow(
		"bridgr job-discovery run ready",
		"user_id", userID,
		"run_uuid", runUUID.String(),
		"status", runPtr.Status,
	)

	out, err := discoveryRunFromRow(runPtr)
	if err != nil {
		return nil, fmt.Errorf("%w: map run: %w", apierrors.ErrInternal, err)
	}
	return &out, nil
}

// V1GetBridgrUserJobDiscoveryRuns handles GET /v1/bridgr/users/{userID}/job-discovery/runs
func (s *server) V1GetBridgrUserJobDiscoveryRuns(w http.ResponseWriter, r *http.Request, userID int32, params types.V1GetBridgrUserJobDiscoveryRunsParams) {
	ctx := r.Context()
	resp, err := s.v1GetBridgrUserJobDiscoveryRuns(ctx, userID, params)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

func (s *server) v1GetBridgrUserJobDiscoveryRuns(ctx context.Context, userID int32, params types.V1GetBridgrUserJobDiscoveryRunsParams) (*types.BridgrJobSearchDiscoveryRunListResponse, error) {
	if err := s.requireStore(); err != nil {
		return nil, err
	}
	limit, offset := bridgrUserLimitOffset(params.Limit, params.Offset)
	rows, err := s.deps.Repo.ListJobSearchDiscoveryRunsByUser(ctx, s.querier(), sqlc.ListJobSearchDiscoveryRunsByUserParams{
		UserID: userID, Limit: limit, Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: list runs: %w", apierrors.ErrInternal, err)
	}
	out := types.BridgrJobSearchDiscoveryRunListResponse{Runs: make([]types.BridgrJobSearchDiscoveryRun, 0, len(rows))}
	for i := range rows {
		item, err := discoveryRunFromRow(&rows[i])
		if err != nil {
			return nil, fmt.Errorf("%w: map run: %w", apierrors.ErrInternal, err)
		}
		out.Runs = append(out.Runs, item)
	}
	return &out, nil
}

// V1GetBridgrJobDiscoveryRun handles GET /v1/bridgr/job-discovery/runs/{runUUID}
func (s *server) V1GetBridgrJobDiscoveryRun(w http.ResponseWriter, r *http.Request, runUUID openapi_types.UUID, params types.V1GetBridgrJobDiscoveryRunParams) {
	ctx := r.Context()
	resp, err := s.v1GetBridgrJobDiscoveryRun(ctx, runUUID)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

func (s *server) v1GetBridgrJobDiscoveryRun(ctx context.Context, runUUID openapi_types.UUID) (*types.BridgrJobSearchDiscoveryRun, error) {
	if err := s.requireStore(); err != nil {
		return nil, err
	}
	pgid, err := uuid.ConvertOapiUUIDToPgUUID(runUUID)
	if err != nil {
		return nil, fmt.Errorf("%w: run uuid: %w", apierrors.ErrBadRequest, err)
	}
	row, err := s.deps.Repo.GetJobSearchDiscoveryRunByUUID(ctx, s.querier(), pgid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: run not found", apierrors.ErrNotFound)
		}
		return nil, fmt.Errorf("%w: load run: %w", apierrors.ErrInternal, err)
	}
	out, err := discoveryRunFromRow(row)
	if err != nil {
		return nil, fmt.Errorf("%w: map run: %w", apierrors.ErrInternal, err)
	}
	return &out, nil
}

// V1GetBridgrUserJobCandidates handles GET /v1/bridgr/users/{userID}/job-candidates
func (s *server) V1GetBridgrUserJobCandidates(w http.ResponseWriter, r *http.Request, userID int32, params types.V1GetBridgrUserJobCandidatesParams) {
	ctx := r.Context()
	resp, err := s.v1GetBridgrUserJobCandidates(ctx, userID, params)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

func (s *server) v1GetBridgrUserJobCandidates(ctx context.Context, userID int32, params types.V1GetBridgrUserJobCandidatesParams) (*types.BridgrJobCandidateListResponse, error) {
	if err := s.requireStore(); err != nil {
		return nil, err
	}
	limit, offset := bridgrUserLimitOffset(params.Limit, params.Offset)
	f := sqlc.ListJobCandidatesByUserFilteredParams{
		UserID:           userID,
		Limit:            limit,
		Offset:           offset,
		DiscoveryRunUuid: pgtype.UUID{Valid: false},
		IngestionStatus:  pgtype.Text{Valid: false},
		SourceBoard:      pgtype.Text{Valid: false},
	}
	if params.DiscoveryRunUuid != nil {
		dr, err := uuid.ConvertOapiUUIDToPgUUID(*params.DiscoveryRunUuid)
		if err != nil {
			return nil, fmt.Errorf("%w: discovery_run_uuid: %w", apierrors.ErrBadRequest, err)
		}
		f.DiscoveryRunUuid = dr
	}
	if params.IngestionStatus != nil && *params.IngestionStatus != "" {
		f.IngestionStatus = pgtype.Text{String: *params.IngestionStatus, Valid: true}
	}
	if params.SourceBoard != nil && *params.SourceBoard != "" {
		f.SourceBoard = pgtype.Text{String: *params.SourceBoard, Valid: true}
	}
	rows, err := s.deps.Repo.ListJobCandidatesByUserFiltered(ctx, s.querier(), f)
	if err != nil {
		return nil, fmt.Errorf("%w: list candidates: %w", apierrors.ErrInternal, err)
	}
	out := types.BridgrJobCandidateListResponse{Candidates: make([]types.BridgrJobCandidate, 0, len(rows))}
	for i := range rows {
		item, err := jobCandidateFromRow(&rows[i])
		if err != nil {
			return nil, fmt.Errorf("%w: map candidate: %w", apierrors.ErrInternal, err)
		}
		out.Candidates = append(out.Candidates, item)
	}
	return &out, nil
}

// V1GetBridgrJobCandidate handles GET /v1/bridgr/job-candidates/{candidateUUID}
func (s *server) V1GetBridgrJobCandidate(w http.ResponseWriter, r *http.Request, candidateUUID openapi_types.UUID, params types.V1GetBridgrJobCandidateParams) {
	ctx := r.Context()
	resp, err := s.v1GetBridgrJobCandidate(ctx, candidateUUID)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

func (s *server) v1GetBridgrJobCandidate(ctx context.Context, candidateUUID openapi_types.UUID) (*types.BridgrJobCandidate, error) {
	if err := s.requireStore(); err != nil {
		return nil, err
	}
	pgid, err := uuid.ConvertOapiUUIDToPgUUID(candidateUUID)
	if err != nil {
		return nil, fmt.Errorf("%w: candidate uuid: %w", apierrors.ErrBadRequest, err)
	}
	row, err := s.deps.Repo.GetJobCandidateByUUID(ctx, s.querier(), pgid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: candidate not found", apierrors.ErrNotFound)
		}
		return nil, fmt.Errorf("%w: load candidate: %w", apierrors.ErrInternal, err)
	}
	out, err := jobCandidateFromRow(row)
	if err != nil {
		return nil, fmt.Errorf("%w: map candidate: %w", apierrors.ErrInternal, err)
	}
	return &out, nil
}

// V1GetBridgrUserJobNotifications handles GET /v1/bridgr/users/{userID}/job-notifications
func (s *server) V1GetBridgrUserJobNotifications(w http.ResponseWriter, r *http.Request, userID int32, params types.V1GetBridgrUserJobNotificationsParams) {
	ctx := r.Context()
	resp, err := s.v1GetBridgrUserJobNotifications(ctx, userID, params)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

func (s *server) v1GetBridgrUserJobNotifications(ctx context.Context, userID int32, params types.V1GetBridgrUserJobNotificationsParams) (*types.BridgrJobNotificationListResponse, error) {
	if err := s.requireStore(); err != nil {
		return nil, err
	}
	limit, offset := bridgrUserLimitOffset(params.Limit, params.Offset)
	f := sqlc.ListJobNotificationsByUserFilteredParams{
		UserID: userID, Limit: limit, Offset: offset,
		Since:  pgtype.Timestamp{Valid: false},
		Status: pgtype.Text{Valid: false},
	}
	if params.Since != nil {
		f.Since = pgtype.Timestamp{Time: params.Since.UTC(), Valid: true}
	}
	if params.Status != nil && *params.Status != "" {
		f.Status = pgtype.Text{String: *params.Status, Valid: true}
	}
	rows, err := s.deps.Repo.ListJobNotificationsByUserFiltered(ctx, s.querier(), f)
	if err != nil {
		return nil, fmt.Errorf("%w: list notifications: %w", apierrors.ErrInternal, err)
	}
	out := types.BridgrJobNotificationListResponse{Notifications: make([]types.BridgrJobNotification, 0, len(rows))}
	for i := range rows {
		item, err := jobNotificationFromRow(&rows[i])
		if err != nil {
			return nil, fmt.Errorf("%w: map notification: %w", apierrors.ErrInternal, err)
		}
		out.Notifications = append(out.Notifications, item)
	}
	return &out, nil
}

// V1PatchBridgrJobNotification handles PATCH /v1/bridgr/job-notifications/{notificationUUID}
func (s *server) V1PatchBridgrJobNotification(w http.ResponseWriter, r *http.Request, notificationUUID openapi_types.UUID, params types.V1PatchBridgrJobNotificationParams) {
	ctx := r.Context()
	var body types.PatchBridgrJobNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, nil, fmt.Errorf("%w: invalid JSON body: %v", apierrors.ErrBadRequest, err))
		return
	}
	resp, err := s.v1PatchBridgrJobNotification(ctx, notificationUUID, body)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

func (s *server) v1PatchBridgrJobNotification(ctx context.Context, notificationUUID openapi_types.UUID, body types.PatchBridgrJobNotificationRequest) (*types.BridgrJobNotification, error) {
	if err := s.requireStore(); err != nil {
		return nil, err
	}
	if body.MarkSeen == nil && body.Status == nil && body.ErrorDetail == nil {
		return nil, fmt.Errorf("%w: at least one of mark_seen, status, error_detail is required", apierrors.ErrBadRequest)
	}
	pgid, err := uuid.ConvertOapiUUIDToPgUUID(notificationUUID)
	if err != nil {
		return nil, fmt.Errorf("%w: notification uuid: %w", apierrors.ErrBadRequest, err)
	}
	row, err := s.deps.Repo.GetJobNotificationByUUID(ctx, s.querier(), pgid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: notification not found", apierrors.ErrNotFound)
		}
		return nil, fmt.Errorf("%w: load notification: %w", apierrors.ErrInternal, err)
	}

	now := pgtype.Timestamp{Time: time.Now().UTC(), Valid: true}

	if body.MarkSeen != nil && *body.MarkSeen {
		updated, err := s.deps.Repo.UpdateJobNotificationSeen(ctx, s.querier(), sqlc.UpdateJobNotificationSeenParams{
			ID:     row.ID,
			Status: "seen",
			SeenAt: now,
		})
		if err != nil {
			return nil, fmt.Errorf("%w: mark seen: %w", apierrors.ErrInternal, err)
		}
		row = updated
	} else {
		if body.Status != nil {
			ed := pgtype.Text{Valid: false}
			if body.ErrorDetail != nil {
				ed = pgtype.Text{String: *body.ErrorDetail, Valid: true}
			}
			updated, err := s.deps.Repo.UpdateJobNotificationStatus(ctx, s.querier(), sqlc.UpdateJobNotificationStatusParams{
				ID:          row.ID,
				Status:      *body.Status,
				SentAt:      pgtype.Timestamp{Valid: false},
				ErrorDetail: ed,
			})
			if err != nil {
				return nil, fmt.Errorf("%w: update status: %w", apierrors.ErrInternal, err)
			}
			row = updated
		} else if body.ErrorDetail != nil {
			updated, err := s.deps.Repo.UpdateJobNotificationStatus(ctx, s.querier(), sqlc.UpdateJobNotificationStatusParams{
				ID:          row.ID,
				Status:      row.Status,
				SentAt:      pgtype.Timestamp{Valid: false},
				ErrorDetail: pgtype.Text{String: *body.ErrorDetail, Valid: true},
			})
			if err != nil {
				return nil, fmt.Errorf("%w: update error_detail: %w", apierrors.ErrInternal, err)
			}
			row = updated
		}
	}

	out, err := jobNotificationFromRow(row)
	if err != nil {
		return nil, fmt.Errorf("%w: map notification: %w", apierrors.ErrInternal, err)
	}
	return &out, nil
}
