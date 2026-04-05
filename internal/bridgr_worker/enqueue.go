package bridgr_worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	guuid "github.com/gofrs/uuid/v5"
)

// EnqueueSkillGapAnalysis sends one job message. No-op when client or queueURL is nil/empty.
func EnqueueSkillGapAnalysis(ctx context.Context, client *sqs.Client, queueURL string, analysisUUID guuid.UUID) error {
	if client == nil || queueURL == "" {
		return nil
	}
	body, err := json.Marshal(AnalysisJobPayload{AnalysisUUID: analysisUUID.String()})
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	_, err = client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(body)),
	})
	if err != nil {
		return fmt.Errorf("sqs send: %w", err)
	}
	return nil
}
