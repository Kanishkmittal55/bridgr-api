package config

import "embed"

//go:embed development.yaml staging.yaml production.yaml
var configFiles embed.FS
