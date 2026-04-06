package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kanishkmittal55/bridgr-api/internal/bridgr_worker"
	"github.com/Kanishkmittal55/bridgr-api/internal/bridgr_worker/config"
	"github.com/Kanishkmittal55/bridgr-api/internal/bridgr_worker/dependencies"
	apicfg "github.com/Kanishkmittal55/bridgr-api/internal/config"
	"github.com/Kanishkmittal55/bridgr-api/internal/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := logger.Get(ctx)
	log.Infow("starting bridgr-worker")

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		log.Infow("bridgr-worker: shutdown signal", "signal", sig)
		cancel()
	}()

	cfg := apicfg.Load()
	if cfg.BridgrQueueURL == "" {
		log.Fatalw("bridgr-worker: set SQS_BRIDGR_QUEUE_URL")
	}

	d, cleanup, err := dependencies.New(ctx, cfg)
	if err != nil {
		log.Fatalw("bridgr-worker: deps", "error", err)
	}
	defer cleanup()

	opts, err := config.Load(cfg.BridgrQueueURL)
	if err != nil {
		log.Fatalw("bridgr-worker: worker config", "error", err)
	}

	// Discovery scheduler started as its own go routine
	if opts.DiscoverySchedulerEnabled {
		go bridgr_worker.RunScheduler(ctx, d, opts)
	}

	w := bridgr_worker.NewWorker(d, opts)

	log.Infow("bridgr-worker configured",
		"env", config.GetEnvironment(),
		"polling_enabled", opts.PollingEnabled,
		"queue_url", opts.QueueURL,
		"max_messages", opts.MaxMessages,
		"discovery_scheduler_enabled", opts.DiscoverySchedulerEnabled,
		"discovery_scheduler_trigger_interval", opts.DiscoverySchedulerTriggerInterval,
	)

	if err := w.Run(ctx); err != nil && err != context.Canceled {
		log.Fatalw("bridgr-worker failed", "error", err)
	}
	log.Infow("bridgr-worker stopped")
}
