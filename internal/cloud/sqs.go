package cloud

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	guuid "github.com/gofrs/uuid/v5"

	"github.com/Kanishkmittal55/bridgr-api/internal/config"
)

// QueueMessage is one SQS ReceiveMessage result. Ack deletes it from the queue.
type QueueMessage struct {
	ID            string
	Body          []byte
	ReceiptHandle string

	client   *sqs.Client
	queueURL string
	acked    bool
}

// NewSQSClient builds an AWS SDK v2 SQS client from config (LocalStack endpoint or default AWS).
func NewSQSClient(ctx context.Context, cfg *config.Config) (*sqs.Client, error) {
	if cfg.SQSEndpoint != "" {
		awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
			awsconfig.WithRegion(cfg.AWSRegion),
			awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
		)
		if err != nil {
			return nil, err
		}
		return sqs.NewFromConfig(awsCfg, func(o *sqs.Options) {
			o.BaseEndpoint = aws.String(cfg.SQSEndpoint)
		}), nil
	}
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(cfg.AWSRegion))
	if err != nil {
		return nil, err
	}
	return sqs.NewFromConfig(awsCfg), nil
}

// QueuePayload is the JSON body for Bridgr job-discovery SQS messages.
type QueuePayload struct {
	RunUUID string `json:"run_uuid,omitempty"`
}

// EnqueueJobDiscovery sends a job-discovery run message. Returns the SQS message id when send succeeds.
func EnqueueJobDiscovery(ctx context.Context, client *sqs.Client, queueURL string, runUUID guuid.UUID, userID int32) (string, error) {
	return enqueueJSON(ctx, client, queueURL, QueuePayload{
		RunUUID: runUUID.String(),
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

// NewLocalQueueMessage is for in-process runs (tests, tooling): Ack succeeds without calling AWS.
func NewLocalQueueMessage(body []byte) *QueueMessage {
	return &QueueMessage{Body: body}
}

// Ack deletes the message from the queue.
func (m *QueueMessage) Ack() error {
	if m.acked {
		return nil
	}
	if m.client == nil {
		m.acked = true
		return nil
	}
	_, err := m.client.DeleteMessage(context.Background(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(m.queueURL),
		ReceiptHandle: aws.String(m.ReceiptHandle),
	})
	if err == nil {
		m.acked = true
	}
	return err
}

// Nack is a no-op; messages become visible again after the visibility timeout.
func (m *QueueMessage) Nack() {}

// ReceiveQueueMessages long-polls SQS and returns wrapped messages (empty slice if none).
func ReceiveQueueMessages(ctx context.Context, client *sqs.Client, queueURL string, maxMessages, waitTimeSeconds, visibilityTimeout int32) ([]*QueueMessage, error) {
	if client == nil || queueURL == "" {
		return nil, nil
	}
	out, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: maxMessages,
		WaitTimeSeconds:     waitTimeSeconds,
		VisibilityTimeout:   visibilityTimeout,
	})
	if err != nil {
		return nil, err
	}
	var msgs []*QueueMessage
	for _, sm := range out.Messages {
		if sm.Body == nil || sm.ReceiptHandle == nil {
			continue
		}
		id := ""
		if sm.MessageId != nil {
			id = *sm.MessageId
		}
		msgs = append(msgs, &QueueMessage{
			ID:            id,
			Body:          []byte(*sm.Body),
			ReceiptHandle: *sm.ReceiptHandle,
			client:        client,
			queueURL:      queueURL,
		})
	}
	return msgs, nil
}
