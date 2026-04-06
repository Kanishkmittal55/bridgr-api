package bridgr_worker

import (
	"context"
	"fmt"

	"github.com/Kanishkmittal55/bridgr-api/internal/cloud"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository/sqlc"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	guuid "github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// EnqueueJobDiscoveryRunForProfile is the single path used by the discovery scheduler and
// POST /v1/bridgr/users/{userID}/job-discovery/runs: build request_params from the profile,
// insert a run, enqueue to SQS when configured, then patch status to queued with sqs_message_id.
func EnqueueJobDiscoveryRunForProfile(
	ctx context.Context,
	repo *repository.Repo,
	q sqlc.Querier,
	sqsClient *sqs.Client,
	queueURL string,
	userID int32,
	prof *sqlc.BridgrJobSearchProfile,
	overrides map[string]interface{},
) (*sqlc.BridgrJobSearchDiscoveryRun, error) {
	reqParams, err := BuildDiscoveryRequestParams(userID, prof, overrides)
	if err != nil {
		return nil, err
	}

	runPtr, err := repo.CreateJobSearchDiscoveryRun(ctx, q, sqlc.CreateJobSearchDiscoveryRunParams{
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
		return nil, fmt.Errorf("create discovery run: %w", err)
	}

	if queueURL != "" && sqsClient != nil {
		runUUID, uerr := guuid.FromBytes(runPtr.Uuid.Bytes[:])
		if uerr != nil {
			return nil, fmt.Errorf("run uuid: %w", uerr)
		}
		msgID, err := cloud.EnqueueJobDiscovery(ctx, sqsClient, queueURL, runUUID, userID)
		if err != nil {
			return nil, fmt.Errorf("enqueue job discovery: %w", err)
		}
		updated, err := repo.PatchJobSearchDiscoveryRun(ctx, q, sqlc.PatchJobSearchDiscoveryRunParams{
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
			return nil, fmt.Errorf("persist sqs id: %w", err)
		}
		runPtr = updated
	}

	return runPtr, nil
}
