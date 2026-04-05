package ctxlog

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type CtxLoggerKey string

const (
	ctxKeyLogFields           CtxLoggerKey = "ctxKeyLogFields"
	CtxKeyLoggerLevelOverride CtxLoggerKey = "ctxKeyLoggerLevelOverride"
)

// ContextualLogger is a sugared Zap logger with optional per-request fields in context.
type ContextualLogger struct {
	*zap.SugaredLogger
	*zap.Config
	opts []Option
}

func New(opts ...Option) (*ContextualLogger, error) {
	cfg := defaultCfg()
	for _, opt := range opts {
		opt(cfg)
	}

	var zc zap.Config
	if cfg.devMode {
		zc = zap.NewDevelopmentConfig()
	} else {
		zc = zap.NewProductionConfig()
		zc.Sampling = nil
	}

	switch cfg.minLevel {
	case DebugLevel:
		zc.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case InfoLevel:
		zc.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case WarnLevel:
		zc.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case ErrorLevel:
		zc.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		return nil, fmt.Errorf("unknown log level: %v", cfg.minLevel)
	}

	logger, err := zc.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create zap logger: %w", err)
	}

	return &ContextualLogger{logger.Sugar(), &zc, opts}, nil
}

func NewNop() *ContextualLogger {
	zc := zap.NewDevelopmentConfig()
	return &ContextualLogger{zap.NewNop().Sugar(), &zc, nil}
}

func (l *ContextualLogger) Add(args ...any) *ContextualLogger {
	return &ContextualLogger{l.With(args...), l.Config, l.opts}
}

func (l *ContextualLogger) AddToCtx(ctx context.Context, args ...any) (*ContextualLogger, context.Context) {
	return l.Add(args...), AddToCtx(ctx, args...)
}

func (l *ContextualLogger) AddFromCtx(ctx context.Context) *ContextualLogger {
	if ctx == nil {
		return &ContextualLogger{l.SugaredLogger, l.Config, l.opts}
	}

	args, ok := ctx.Value(ctxKeyLogFields).([]any)
	if !ok {
		return &ContextualLogger{l.SugaredLogger, l.Config, l.opts}
	}

	returnedLogger := l.Add(args...)
	level := levelOverride(ctx, Level(l.SugaredLogger.Level().String()))
	if level != Level(l.SugaredLogger.Level().String()) {
		return returnedLogger.WithLevel(level, args...)
	}

	return returnedLogger
}

func (l *ContextualLogger) AtomicLevel() zap.AtomicLevel {
	return l.Config.Level
}

func AddToCtx(ctx context.Context, args ...any) context.Context {
	contextArgs, ok := ctx.Value(ctxKeyLogFields).([]any)
	if !ok {
		contextArgs = []any{}
	}
	return context.WithValue(ctx, ctxKeyLogFields, append(contextArgs, args...))
}

func levelOverride(ctx context.Context, fallback Level) Level {
	override := ctx.Value(CtxKeyLoggerLevelOverride)
	if override == nil {
		return fallback
	}
	return override.(Level)
}

func (l *ContextualLogger) WithLevel(level Level, args ...any) *ContextualLogger {
	newOpts := append(l.opts, MinLevel(level))
	newLogger, _ := New(newOpts...)
	return newLogger.Add(args...)
}
