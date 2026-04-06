package bridgr_worker

import (
	"context"
	"time"

	"github.com/Kanishkmittal55/bridgr-api/internal/bridgr_worker/config"
	"github.com/Kanishkmittal55/bridgr-api/internal/bridgr_worker/dependencies"
	"github.com/Kanishkmittal55/bridgr-api/internal/cloud"
	apicfg "github.com/Kanishkmittal55/bridgr-api/internal/config"
	"github.com/Kanishkmittal55/bridgr-api/internal/logger"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// Worker polls the Bridgr job-discovery SQS queue.
type Worker struct {
	sqsClient *sqs.Client
	opts      config.WorkerOpts
	proc      *Processor
}

// NewWorker constructs a worker. d.JobSearch may be nil if RADAR_ADDR is unset (discovery stub only).
func NewWorker(d *dependencies.Deps, opts config.WorkerOpts) *Worker {
	global := apicfg.Load()
	return &Worker{
		sqsClient: d.SQSClient,
		opts:      opts,
		proc:      NewProcessor(d.Repo, d.HsQuerier, d.JobSearch, d.Radar, d.S3, global.HassleSkipS3Bucket, opts, d.OpenAIAPIKey),
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

	log.Infow("bridgr-worker: polling", "queue_url", w.opts.QueueURL, "max_messages", w.opts.MaxMessages, "wait_seconds", w.opts.WaitTimeSeconds)

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
	msgs, err := cloud.ReceiveQueueMessages(ctx, w.sqsClient, w.opts.QueueURL, w.opts.MaxMessages, w.opts.WaitTimeSeconds, w.opts.VisibilityTimeout)
	if err != nil {
		return err
	}
	if len(msgs) == 0 {
		return nil
	}

	log := logger.Get(ctx)
	log.Infow("bridgr-worker: received messages", "count", len(msgs))

	for _, msg := range msgs {
		if err := w.proc.Process(ctx, msg); err != nil {
			log.Errorw("bridgr-worker: process failed", "message_id", msg.ID, "error", err)
		}
	}
	return nil
}
