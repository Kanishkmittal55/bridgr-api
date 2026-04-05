package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/hassleskip/bridgr-api/internal/api"
	"github.com/hassleskip/bridgr-api/internal/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := api.Run(ctx); err != nil {
		logger.Get(ctx).Fatalw("bridgr-api exited", "error", err)
	}
}
