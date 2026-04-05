package ctxlog

type Option func(*cfg)

func DevMode() Option {
	return func(cfg *cfg) {
		cfg.devMode = true
	}
}

func MinLevel(level Level) Option {
	return func(cfg *cfg) {
		cfg.minLevel = level
	}
}
