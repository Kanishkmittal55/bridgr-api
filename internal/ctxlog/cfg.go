package ctxlog

type cfg struct {
	devMode  bool
	minLevel Level
}

func defaultCfg() *cfg {
	return &cfg{
		devMode:  false,
		minLevel: InfoLevel,
	}
}
