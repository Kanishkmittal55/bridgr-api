package bridgr_worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	guuid "github.com/gofrs/uuid/v5"
)

// EnqueueSkillGapAnalysis sends one skill-gap job message. No-op when client or queueURL is nil/empty.
func EnqueueSkillGapAnalysis(ctx context.Context, client *sqs.Client, queueURL string, analysisUUID guuid.UUID) error {
	_, err := enqueueJSON(ctx, client, queueURL, QueuePayload{
		Kind:         KindSkillGapAnalysis,
		AnalysisUUID: analysisUUID.String(),
	})
	return err
}

// EnqueueJobDiscovery sends a job-discovery run message. Returns the SQS message id when send succeeds.
func EnqueueJobDiscovery(ctx context.Context, client *sqs.Client, queueURL string, runUUID guuid.UUID, userID int32) (string, error) {
	return enqueueJSON(ctx, client, queueURL, QueuePayload{
		Kind:    KindJobDiscovery,
		RunUUID: runUUID.String(),
		UserID:  userID,
	})
}

func enqueueJSON(ctx context.Context, client *sqs.Client, queueURL string, payload QueuePayload) (string, error) {
	if client == nil || queueURL == "" {
		return "", nil
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}
	out, err := client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(body)),
	})
	if err != nil {
		return "", fmt.Errorf("sqs send: %w", err)
	}
	if out.MessageId == nil {
		return "", nil
	}
	return *out.MessageId, nil
}
