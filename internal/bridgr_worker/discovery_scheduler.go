package bridgr_worker

import (
	"context"
	"time"

	"github.com/Kanishkmittal55/bridgr-api/internal/bridgr_worker/config"
	"github.com/Kanishkmittal55/bridgr-api/internal/bridgr_worker/dependencies"
	"github.com/Kanishkmittal55/bridgr-api/internal/logger"
)

// RunScheduler ticks until ctx is cancelled. Same ctx as the SQS worker so SIGTERM stops both.
// When implemented, this will list due harvest schedules and enqueue discovery work.
func RunScheduler(ctx context.Context, d *dependencies.Deps, opts config.WorkerOpts) {
	_ = d // reserved for DB + enqueue
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
			log.Infow("Triggering job discovery")
		}
	}
}
