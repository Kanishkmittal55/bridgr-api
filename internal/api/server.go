package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Kanishkmittal55/bridgr-api/internal/api/deps"
	"github.com/Kanishkmittal55/bridgr-api/internal/config"
	"github.com/Kanishkmittal55/bridgr-api/internal/logger"
)

// Run boots the HTTP server and blocks until ctx is cancelled, then shuts down gracefully
// (same role as users/internal/api/server.go).
func Run(ctx context.Context) error {
	cfg := config.Load()
	log := logger.Get(ctx)

	d, cleanup, err := deps.New(ctx, cfg)
	if err != nil {
		return err
	}
	defer cleanup()

	addr := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{Addr: addr, Handler: Routes(d)}

	errCh := make(chan error, 1)
	go func() {
		log.Infow("bridgr-api listening", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown: %w", err)
		}
		if err := <-errCh; err != nil {
			return err
		}
		log.Info("server stopped")
		return nil
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("listen: %w", err)
		}
		return nil
	}
}
