package logger

import (
	"context"
	"sync"

	"github.com/Kanishkmittal55/bridgr-api/internal/config"
	"github.com/Kanishkmittal55/bridgr-api/internal/ctxlog"
	"github.com/Kanishkmittal55/bridgr-api/internal/env"
)

var (
	base *ctxlog.ContextualLogger
	once sync.Once
)

// Get returns the process logger, optionally scoped to context.
func Get(ctx ...context.Context) *ctxlog.ContextualLogger {
	once.Do(func() {
		cfg := config.Load()
		opts := []ctxlog.Option{ctxlog.MinLevel(cfg.LogLevel)}
		if !env.IsNonDevelopment(cfg.Env) {
			opts = append(opts, ctxlog.DevMode())
		}
		l, err := ctxlog.New(opts...)
		if err != nil {
			panic(err)
		}
		base = l
	})
	if len(ctx) == 0 || ctx[0] == nil {
		return base
	}
	return base.AddFromCtx(ctx[0])
}
