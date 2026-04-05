package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kanishkmittal55/bridgr-api/internal/bridgr_worker"
	"github.com/Kanishkmittal55/bridgr-api/internal/config"
	"github.com/Kanishkmittal55/bridgr-api/internal/logger"
	"github.com/Kanishkmittal55/bridgr-api/internal/radar"
	"github.com/Kanishkmittal55/bridgr-api/internal/rdbms"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository/sqlc"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := logger.Get(ctx)
	log.Infow("starting bridgr-worker")

	cfg := config.Load()
	queueURL := cfg.BridgrQueueURL
	if queueURL == "" {
		log.Fatalw("bridgr-worker: set SQS_BRIDGR_QUEUE_URL")
	}

	pool, err := rdbms.NewConn(rdbms.ConnStr(cfg))
	if err != nil {
		log.Fatalw("bridgr-worker: db connect", "error", err)
	}
	defer pool.Close()

	q := sqlc.New(pool)
	repo := repository.New()

	sqsClient, err := createSQSClient(ctx, cfg)
	if err != nil {
		log.Fatalw("bridgr-worker: sqs client", "error", err)
	}

	var jobSearch *radar.JobSearchClient
	if cfg.RadarAddr != "" {
		js, err := radar.NewJobSearchClient(radar.JobSearchConfig{Addr: cfg.RadarAddr})
		if err != nil {
			log.Warnw("bridgr-worker: Radar JobSearch client disabled", "addr", cfg.RadarAddr, "error", err)
		} else {
			jobSearch = js
			defer func() {
				if cErr := js.Close(); cErr != nil {
					log.Warnw("bridgr-worker: Radar client close", "error", cErr)
				}
			}()
			log.Infow("bridgr-worker: Radar JobSearch client enabled", "addr", cfg.RadarAddr)
		}
	} else {
		log.Warnw("bridgr-worker: RADAR_ADDR empty; job discovery will not call FindJobs")
	}

	opts := bridgr_worker.WorkerOptsFromEnv(queueURL)
	w := bridgr_worker.NewWorker(sqsClient, repo, q, opts, jobSearch)

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		log.Infow("bridgr-worker: shutdown signal", "signal", sig)
		cancel()
	}()

	log.Infow("bridgr-worker configured",
		"polling_enabled", opts.PollingEnabled,
		"queue_url", opts.QueueURL,
		"max_messages", opts.MaxMessages,
	)

	if err := w.Run(ctx); err != nil && err != context.Canceled {
		log.Fatalw("bridgr-worker failed", "error", err)
	}
	log.Infow("bridgr-worker stopped")
}

func createSQSClient(ctx context.Context, cfg *config.Config) (*sqs.Client, error) {
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
