package bridgr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	apierrors "github.com/Kanishkmittal55/bridgr-api/internal/apierrors"
	"github.com/Kanishkmittal55/bridgr-api/internal/bridgr_worker"
	"github.com/Kanishkmittal55/bridgr-api/internal/config"
	"github.com/Kanishkmittal55/bridgr-api/internal/logger"
	"github.com/Kanishkmittal55/bridgr-api/internal/radar"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository/sqlc"
	"github.com/Kanishkmittal55/bridgr-api/internal/uuid"
	types "github.com/Kanishkmittal55/bridgr-api/pkg/types"
	guuid "github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// V1GetBridgrUserJobSearchProfile handles GET /v1/bridgr/users/{userID}/job-search-profile
func (s *server) V1GetBridgrUserJobSearchProfile(w http.ResponseWriter, r *http.Request, userID int32, params types.V1GetBridgrUserJobSearchProfileParams) {
	ctx := r.Context()
	resp, err := s.v1GetBridgrUserJobSearchProfile(ctx, userID)
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, err)
}

func (s *server) v1GetBridgrUserJobSearchProfile(ctx context.Context, userID int32) (*types.BridgrJobSearchProfile, error) {
	if err := s.requireStore(); err != nil {
		return nil, err
	}
	row, err := s.deps.Repo.GetJobSearchProfileByUserID(ctx, s.querier(), userID)
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

// V1PutBridgrUserJobSearchProfile handles PUT /v1/bridgr/users/{userID}/job-search-profile
func (s *server) V1PutBridgrUserJobSearchProfile(w http.ResponseWriter, r *http.Request, userID int32, params types.V1PutBridgrUserJobSearchProfileParams) {
	ctx := r.Context()
	var body types.UpsertBridgrJobSearchProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, nil, fmt.Errorf("%w: invalid JSON body: %v", apierrors.ErrBadRequest, err))
		return
	}
	resp, created, err := s.v1PutBridgrUserJobSearchProfile(ctx, userID, body)
	if err != nil {
		s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, nil, err)
		return
	}
	if created {
		s.writeCreated(w, r, resp)
		return
	}
	s.deps.ResponseWriter.WriteOkResponse(ctx, w, r, resp, nil)
}

func (s *server) v1PutBridgrUserJobSearchProfile(ctx context.Context, userID int32, body types.UpsertBridgrJobSearchProfileRequest) (*types.BridgrJobSearchProfile, bool, error) {
	if err := s.requireStore(); err != nil {
		return nil, false, err
	}
	existing, err := s.deps.Repo.GetJobSearchProfileByUserID(ctx, s.querier(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			cp, perr := profileToCreateParams(userID, body)
			if perr != nil {
				return nil, false, fmt.Errorf("%w: %v", apierrors.ErrBadRequest, perr)
			}
			row, cerr := s.deps.Repo.CreateJobSearchProfile(ctx, s.querier(), cp)
			if cerr != nil {
				return nil, false, fmt.Errorf("%w: create profile: %w", apierrors.ErrInternal, cerr)
			}
			out, merr := jobSearchProfileFromRow(row)
			if merr != nil {
				return nil, true, fmt.Errorf("%w: map profile: %w", apierrors.ErrInternal, merr)
			}
			return &out, true, nil
		}
		return nil, false, fmt.Errorf("%w: load profile: %w", apierrors.ErrInternal, err)
	}
	up, err := profileToUpdateByUserParams(existing, body)
	if err != nil {
		return nil, false, fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err)
	}
	row, err := s.deps.Repo.UpdateJobSearchProfileByUserID(ctx, s.querier(), up)
	if err != nil {
		return nil, false, fmt.Errorf("%w: update profile: %w", apierrors.ErrInternal, err)
	}
	out, err := jobSearchProfileFromRow(row)
	if err != nil {
		return nil, false, fmt.Errorf("%w: map profile: %w", apierrors.ErrInternal, err)
	}
	return &out, false, nil
}

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
	if cfg.JobDiscoveryMaxRunsPerHour > 0 {
		n, err := s.deps.Repo.CountJobSearchDiscoveryRunsByUserSince(ctx, s.querier(), sqlc.CountJobSearchDiscoveryRunsByUserSinceParams{
			UserID:    userID,
			CreatedAt: discoveryRateWindowStart(),
		})
		if err != nil {
			return nil, fmt.Errorf("%w: rate check: %w", apierrors.ErrInternal, err)
		}
		if n >= int64(cfg.JobDiscoveryMaxRunsPerHour) {
			return nil, fmt.Errorf("%w: too many job discovery runs in the last hour for this user", apierrors.ErrTooManyRequests)
		}
	}

	reqParams, err := s.buildDiscoveryRequestParams(ctx, userID, body.RequestParams)
	if err != nil {
		return nil, err
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

	// Sync-in-dev runs FindJobs inside the API process (RADAR_ADDR). When BRIDGR_JOB_DISCOVERY_SYNC_IN_DEV
	// is set, it wins over SQS so local Docker does not require bridgr-worker for discovery; otherwise a
	// queued run is only processed when the worker drains SQS_BRIDGR_QUEUE_URL.
	if cfg.JobDiscoverySyncInDev {
		syncBody, err := json.Marshal(bridgr_worker.QueuePayload{
			Kind:    bridgr_worker.KindJobDiscovery,
			RunUUID: runUUID.String(),
			UserID:  userID,
		})
		if err != nil {
			return nil, fmt.Errorf("%w: marshal sync payload: %w", apierrors.ErrInternal, err)
		}
		var jobSearch *radar.JobSearchClient
		if addr := config.Get().RadarAddr; addr != "" {
			js, jerr := radar.NewJobSearchClient(radar.JobSearchConfig{Addr: addr})
			if jerr != nil {
				logger.Get(ctx).Warnw("sync job discovery: Radar client disabled", "addr", addr, "error", jerr)
			} else {
				jobSearch = js
				defer func() { _ = js.Close() }()
			}
		}
		if err := bridgr_worker.ProcessLocalJSON(ctx, s.deps.Repo, s.querier(), jobSearch, syncBody); err != nil {
			return nil, fmt.Errorf("%w: sync job discovery: %w", apierrors.ErrInternal, err)
		}
		ref, err := s.deps.Repo.GetJobSearchDiscoveryRunByUUID(ctx, s.querier(), runPtr.Uuid)
		if err != nil {
			return nil, fmt.Errorf("%w: reload run: %w", apierrors.ErrInternal, err)
		}
		runPtr = ref
	} else if cfg.BridgrQueueURL != "" && s.deps.SQSClient != nil {
		msgID, err := bridgr_worker.EnqueueJobDiscovery(ctx, s.deps.SQSClient, cfg.BridgrQueueURL, runUUID, userID)
		if err != nil {
			return nil, fmt.Errorf("%w: enqueue job discovery: %w", apierrors.ErrInternal, err)
		}
		updated, err := s.deps.Repo.UpdateJobSearchDiscoveryRunStarted(ctx, s.querier(), sqlc.UpdateJobSearchDiscoveryRunStartedParams{
			ID:           runPtr.ID,
			Status:       "queued",
			StartedAt:    pgtype.Timestamp{Valid: false},
			SqsMessageID: pgtype.Text{String: msgID, Valid: msgID != ""},
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

func (s *server) buildDiscoveryRequestParams(ctx context.Context, userID int32, overrides *map[string]interface{}) ([]byte, error) {
	merged := map[string]interface{}{"user_id": userID}
	prof, err := s.deps.Repo.GetJobSearchProfileByUserID(ctx, s.querier(), userID)
	if err == nil && prof != nil {
		var roles interface{}
		_ = json.Unmarshal(prof.TargetRoles, &roles)
		var locs interface{}
		_ = json.Unmarshal(prof.Locations, &locs)
		var boards interface{}
		_ = json.Unmarshal(prof.BoardsEnabled, &boards)
		var match interface{}
		_ = json.Unmarshal(prof.Matching, &match)
		merged["target_roles"] = roles
		merged["locations"] = locs
		merged["boards_enabled"] = boards
		merged["matching"] = match
		merged["max_surfaced_jobs"] = prof.MaxSurfacedJobs
		if prof.CanonicalCvAnalysisUuid.Valid {
			u, _ := uuid.ToString(prof.CanonicalCvAnalysisUuid.Bytes)
			merged["canonical_cv_analysis_uuid"] = u
		}
	}
	if overrides != nil {
		for k, v := range *overrides {
			merged[k] = v
		}
	}
	b, err := json.Marshal(merged)
	if err != nil {
		return nil, fmt.Errorf("%w: request_params: %w", apierrors.ErrBadRequest, err)
	}
	return b, nil
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
