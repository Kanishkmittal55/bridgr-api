package logger

import (
	"context"
	"sync"

	"github.com/hassleskip/bridgr-api/internal/config"
	"github.com/hassleskip/hassle-go/pkg/ctxlogger"
	"github.com/hassleskip/hassle-go/pkg/env"
)

var (
	base *ctxlogger.ContextualLogger
	once sync.Once
)

// Get returns the process logger, optionally scoped to context.
func Get(ctx ...context.Context) *ctxlogger.ContextualLogger {
	once.Do(func() {
		cfg := config.Load()
		opts := []ctxlogger.Option{ctxlogger.MinLevel(cfg.LogLevel)}
		if !env.IsNonDevelopment(cfg.Env) {
			opts = append(opts, ctxlogger.DevMode())
		}
		l, err := ctxlogger.New(opts...)
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
