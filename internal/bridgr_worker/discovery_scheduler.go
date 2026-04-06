package bridgr_worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Kanishkmittal55/bridgr-api/internal/bridgr_worker/config"
	"github.com/Kanishkmittal55/bridgr-api/internal/bridgr_worker/dependencies"
	"github.com/Kanishkmittal55/bridgr-api/internal/cloud"
	"github.com/Kanishkmittal55/bridgr-api/internal/logger"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository/sqlc"
	"github.com/Kanishkmittal55/bridgr-api/internal/uuid"
	guuid "github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// discoverySchedulerUserID is the user whose job_search_profiles rows the scheduler loads, dummy till we get a real user.
const discoverySchedulerUserID int32 = 1

// RunScheduler ticks until ctx is cancelled.
func RunScheduler(ctx context.Context, d *dependencies.Deps, opts config.WorkerOpts) {
	log := logger.Get(ctx)
	tick := time.Duration(opts.DiscoverySchedulerTriggerInterval) * time.Second
	if tick <= 0 {
		tick = 60 * time.Second
	}
	t := time.NewTicker(tick)
	defer t.Stop()
	log.Infow("bridgr-worker: discovery scheduler started", "tick", tick)
	for {
		select {
		case <-ctx.Done():
			log.Infow("bridgr-worker: discovery scheduler stopped")
			return
		case <-t.C:
			runScheduledDiscoveryTick(ctx, d, opts)
		}
	}
}

func runScheduledDiscoveryTick(ctx context.Context, d *dependencies.Deps, opts config.WorkerOpts) {
	log := logger.Get(ctx)
	uid := discoverySchedulerUserID
	profiles, err := d.Repo.ListJobSearchProfilesByUserID(ctx, d.HsQuerier, uid)
	if err != nil {
		log.Errorw("bridgr-worker: discovery scheduler list profiles by user", "user_id", uid, "error", err)
		return
	}
	if len(profiles) == 0 {
		log.Debugw("bridgr-worker: discovery scheduler tick (no profiles)", "user_id", uid)
		return
	}

	var okCount int
	for i := range profiles {
		if enqueueScheduledDiscoveryForProfile(ctx, d, opts, &profiles[i]) {
			okCount++
		}
	}
	log.Infow("bridgr-worker: discovery scheduler tick done", "user_id", uid, "profiles", len(profiles), "enqueued", okCount)
}

func enqueueScheduledDiscoveryForProfile(ctx context.Context, d *dependencies.Deps, opts config.WorkerOpts, prof *sqlc.BridgrJobSearchProfile) bool {
	log := logger.Get(ctx)
	profileUUID := ""
	if prof.Uuid.Valid {
		if s, err := uuid.ToString(prof.Uuid.Bytes); err == nil {
			profileUUID = s
		}
	}

	board := prof.SourceBoard
	if board == "" {
		board = "indeed"
	}

	reqParams, err := BuildDiscoveryRequestParams(prof.UserID, prof, nil)
	if err != nil {
		log.Errorw("bridgr-worker: discovery scheduler request_params", "user_id", prof.UserID, "job_search_profile_uuid", profileUUID, "error", err)
		return false
	}

	// Creating an empty run
	runPtr, err := d.Repo.CreateJobSearchDiscoveryRun(ctx, d.HsQuerier, sqlc.CreateJobSearchDiscoveryRunParams{
		UserID:            prof.UserID,
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
		log.Errorw("bridgr-worker: discovery scheduler create run", "user_id", prof.UserID, "job_search_profile_uuid", profileUUID, "error", err)
		return false
	}

	runUUID, uerr := guuid.FromBytes(runPtr.Uuid.Bytes[:])
	if uerr != nil {
		log.Errorw("bridgr-worker: discovery scheduler run uuid", "user_id", prof.UserID, "job_search_profile_uuid", profileUUID, "error", uerr)
		return false
	}

	msgID, err := cloud.EnqueueJobDiscovery(ctx, d.SQSClient, opts.QueueURL, runUUID, prof.UserID)
	if err != nil {
		log.Errorw("bridgr-worker: discovery scheduler enqueue", "user_id", prof.UserID, "job_search_profile_uuid", profileUUID, "run_uuid", runUUID.String(), "error", err)
		return false
	}

	if _, err := d.Repo.SetJobSearchDiscoveryRunStatus(ctx, d.HsQuerier, sqlc.SetJobSearchDiscoveryRunStatusParams{
		ID:     runPtr.ID,
		Status: "queued",
	}); err != nil {
		log.Errorw("bridgr-worker: discovery scheduler mark queued", "user_id", prof.UserID, "job_search_profile_uuid", profileUUID, "error", err)
		return false
	}

	log.Infow("bridgr-worker: discovery scheduler enqueued", "user_id", prof.UserID, "job_search_profile_uuid", profileUUID, "run_uuid", runUUID.String(), "source_board", board, "sqs_message_id", msgID)

	return true
}

// ###################
// Helpers
// ###################

func BuildDiscoveryRequestParams(userID int32, prof *sqlc.BridgrJobSearchProfile, overrides map[string]interface{}) ([]byte, error) {
	merged := map[string]interface{}{"user_id": userID}
	if prof != nil {
		if prof.Uuid.Valid {
			if u, err := uuid.ToString(prof.Uuid.Bytes); err == nil {
				merged["job_search_profile_uuid"] = u
			}
		}
		merged["target_role"] = prof.TargetRole
		merged["location"] = prof.Location
		merged["source_board"] = prof.SourceBoard
		merged["career_switch"] = prof.CareerSwitch
		merged["company_stage"] = prof.CompanyStage
		merged["seniority_goal"] = prof.SeniorityGoal
		merged["compensation_goal"] = prof.CompensationGoal
		merged["software_stack_must_have"] = prof.SoftwareStackMustHave
		if prof.CanonicalCvAnalysisUuid.Valid {
			u, _ := uuid.ToString(prof.CanonicalCvAnalysisUuid.Bytes)
			merged["canonical_cv_analysis_uuid"] = u
		}
	}
	for k, v := range overrides {
		merged[k] = v
	}
	b, err := json.Marshal(merged)
	if err != nil {
		return nil, fmt.Errorf("request_params: %w", err)
	}
	return b, nil
}
