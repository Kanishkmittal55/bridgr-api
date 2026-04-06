package config

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// EnvVar is the environment variable used to pick YAML (development | staging | production).
const EnvVar = "ENV"

// GetEnvironment resolves ENV to a config file stem (development.yaml, etc.).
func GetEnvironment() string {
	env := strings.ToLower(strings.TrimSpace(os.Getenv(EnvVar)))
	switch env {
	case "production", "prod":
		return "production"
	case "staging", "stage":
		return "staging"
	default:
		return "development"
	}
}

// Load reads embedded YAML for the current ENV, merges BRIDGR_* overrides, validates, and sets QueueURL.
func Load(queueURL string) (WorkerOpts, error) {
	env := GetEnvironment()
	filename := env + ".yaml"
	data, err := configFiles.ReadFile(filename)
	if err != nil {
		return WorkerOpts{}, fmt.Errorf("bridgr_worker/config: read %s: %w", filename, err)
	}

	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		return WorkerOpts{}, fmt.Errorf("bridgr_worker/config: parse %s: %w", filename, err)
	}

	// Defaults if a key is missing from YAML (defensive).
	v.SetDefault("polling_enabled", true)
	v.SetDefault("max_messages", int32(5))
	v.SetDefault("wait_time_seconds", int32(20))
	v.SetDefault("visibility_timeout", int32(300))
	v.SetDefault("poll_error_backoff_sec", 5)
	v.SetDefault("discovery_scheduler_enabled", false)
	v.SetDefault("discovery_scheduler_trigger_interval", 60)

	v.SetEnvPrefix("BRIDGR_WORKER")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Historical env names (SQS_* infix) kept for backward compatibility.
	_ = v.BindEnv("max_messages", "BRIDGR_WORKER_SQS_MAX_MESSAGES")
	_ = v.BindEnv("wait_time_seconds", "BRIDGR_WORKER_SQS_WAIT_SECONDS")
	_ = v.BindEnv("visibility_timeout", "BRIDGR_WORKER_SQS_VISIBILITY_TIMEOUT")
	_ = v.BindEnv("poll_error_backoff_sec", "BRIDGR_WORKER_POLL_ERROR_BACKOFF_SEC")

	// Discovery scheduler: short names without BRIDGR_WORKER_ prefix.
	_ = v.BindEnv("discovery_scheduler_enabled", "BRIDGR_DISCOVERY_SCHEDULER_ENABLED")
	_ = v.BindEnv("discovery_scheduler_trigger_interval", "BRIDGR_DISCOVERY_SCHEDULER_TRIGGER_INTERVAL")

	var opts WorkerOpts
	if err := v.Unmarshal(&opts); err != nil {
		return WorkerOpts{}, fmt.Errorf("bridgr_worker/config: unmarshal: %w", err)
	}
	opts.QueueURL = queueURL

	// Env-only booleans often use "1"/"true" — viper may not parse those for bool without cast.
	if s := strings.TrimSpace(os.Getenv("BRIDGR_WORKER_POLLING_ENABLED")); s != "" {
		opts.PollingEnabled = parseBool(s)
	}
	if s := strings.TrimSpace(os.Getenv("BRIDGR_DISCOVERY_SCHEDULER_ENABLED")); s != "" {
		opts.DiscoverySchedulerEnabled = parseBool(s)
	}

	if opts.DiscoverySchedulerTriggerInterval < 1 {
		opts.DiscoverySchedulerTriggerInterval = 60
	}

	if err := validate(&opts); err != nil {
		return WorkerOpts{}, err
	}
	return opts, nil
}

// MustLoad calls Load and panics on error.
func MustLoad(queueURL string) WorkerOpts {
	o, err := Load(queueURL)
	if err != nil {
		panic(err)
	}
	return o
}

func parseBool(s string) bool {
	switch strings.ToLower(s) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
