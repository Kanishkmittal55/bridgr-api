package bridgr_worker

import (
	"context"
	"time"

	"github.com/Kanishkmittal55/bridgr-api/internal/logger"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository/sqlc"
	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
)

// Worker polls the Bridgr skill-gap SQS queue.
type Worker struct {
	sqsClient *awssqs.Client
	opts      WorkerOpts
	proc      *Processor
}

// NewWorker constructs a worker.
func NewWorker(sqsClient *awssqs.Client, repo *repository.Repo, q sqlc.Querier, opts WorkerOpts) *Worker {
	return &Worker{
		sqsClient: sqsClient,
		opts:      opts,
		proc:      NewProcessor(repo, q),
	}
}

// Run blocks until ctx is cancelled.
func (w *Worker) Run(ctx context.Context) error {
	log := logger.Get(ctx)
	if !w.opts.PollingEnabled {
		log.Infow("bridgr-worker: polling disabled, standby until shutdown")
		<-ctx.Done()
		return ctx.Err()
	}
	if w.opts.QueueURL == "" {
		log.Infow("bridgr-worker: empty queue URL, standby until shutdown")
		<-ctx.Done()
		return ctx.Err()
	}

	log.Infow("bridgr-worker: polling",
		"queue_url", w.opts.QueueURL,
		"max_messages", w.opts.MaxMessages,
		"wait_seconds", w.opts.WaitTimeSeconds,
	)

	backoff := time.Duration(w.opts.PollErrorBackoffSec) * time.Second
	if backoff <= 0 {
		backoff = 5 * time.Second
	}

	for {
		select {
		case <-ctx.Done():
			log.Infow("bridgr-worker: shutting down")
			return ctx.Err()
		default:
			if err := w.pollAndProcess(ctx); err != nil {
				log.Errorw("bridgr-worker: poll cycle", "error", err)
				time.Sleep(backoff)
			}
		}
	}
}

func (w *Worker) pollAndProcess(ctx context.Context) error {
	out, err := w.sqsClient.ReceiveMessage(ctx, &awssqs.ReceiveMessageInput{
		QueueUrl:            &w.opts.QueueURL,
		MaxNumberOfMessages: w.opts.MaxMessages,
		WaitTimeSeconds:     w.opts.WaitTimeSeconds,
		VisibilityTimeout:   w.opts.VisibilityTimeout,
	})
	if err != nil {
		return err
	}
	if len(out.Messages) == 0 {
		return nil
	}

	log := logger.Get(ctx)
	log.Infow("bridgr-worker: received messages", "count", len(out.Messages))

	for _, sqsMsg := range out.Messages {
		msg := &Message{
			ID:            *sqsMsg.MessageId,
			Body:          []byte(*sqsMsg.Body),
			ReceiptHandle: *sqsMsg.ReceiptHandle,
			sqsClient:     w.sqsClient,
			queueURL:      w.opts.QueueURL,
		}
		if err := w.proc.Process(ctx, msg); err != nil {
			log.Errorw("bridgr-worker: process failed", "message_id", msg.ID, "error", err)
		}
	}
	return nil
}

// Message wraps SQS delivery for Ack.
type Message struct {
	ID            string
	Body          []byte
	ReceiptHandle string
	sqsClient     *awssqs.Client
	queueURL      string
	acked         bool
}

// Ack deletes the message from the queue.
func (m *Message) Ack() error {
	if m.acked {
		return nil
	}
	_, err := m.sqsClient.DeleteMessage(context.Background(), &awssqs.DeleteMessageInput{
		QueueUrl:      &m.queueURL,
		ReceiptHandle: &m.ReceiptHandle,
	})
	if err == nil {
		m.acked = true
	}
	return err
}

// Nack leaves the message to retry after visibility timeout.
func (m *Message) Nack() {}
