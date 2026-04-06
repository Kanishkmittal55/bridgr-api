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
	v.SetDefault("discovery_require_canonical_cv", true)
	v.SetDefault("discovery_score_threshold", 0.72)
	v.SetDefault("discovery_prefilter_min_score", 0.5)
	v.SetDefault("discovery_max_iterations", int32(5))
	v.SetDefault("discovery_max_llm_calls_per_run", int32(30))
	v.SetDefault("discovery_iteration_time_budget_sec", 0)
	v.SetDefault("discovery_findjobs_page_cap", int32(20))
	v.SetDefault("discovery_llm_prefilter_enabled", false)
	v.SetDefault("discovery_llm_score_enabled", false)
	v.SetDefault("discovery_llm_query_refine_enabled", false)
	v.SetDefault("discovery_openai_model", "gpt-4o-mini")
	v.SetDefault("discovery_openai_base_url", "")

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
	_ = v.BindEnv("discovery_require_canonical_cv", "BRIDGR_DISCOVERY_REQUIRE_CANONICAL_CV")
	_ = v.BindEnv("discovery_score_threshold", "BRIDGR_DISCOVERY_SCORE_THRESHOLD")
	_ = v.BindEnv("discovery_prefilter_min_score", "BRIDGR_DISCOVERY_PREFILTER_MIN_SCORE")
	_ = v.BindEnv("discovery_max_iterations", "BRIDGR_DISCOVERY_MAX_ITERATIONS")
	_ = v.BindEnv("discovery_max_llm_calls_per_run", "BRIDGR_DISCOVERY_MAX_LLM_CALLS_PER_RUN")
	_ = v.BindEnv("discovery_iteration_time_budget_sec", "BRIDGR_DISCOVERY_ITERATION_TIME_BUDGET_SEC")
	_ = v.BindEnv("discovery_findjobs_page_cap", "BRIDGR_DISCOVERY_FINDJOBS_PAGE_CAP")
	_ = v.BindEnv("discovery_llm_prefilter_enabled", "BRIDGR_DISCOVERY_LLM_PREFILTER_ENABLED")
	_ = v.BindEnv("discovery_llm_score_enabled", "BRIDGR_DISCOVERY_LLM_SCORE_ENABLED")
	_ = v.BindEnv("discovery_llm_query_refine_enabled", "BRIDGR_DISCOVERY_LLM_QUERY_REFINE_ENABLED")
	_ = v.BindEnv("discovery_openai_model", "BRIDGR_DISCOVERY_OPENAI_MODEL")
	_ = v.BindEnv("discovery_openai_base_url", "BRIDGR_DISCOVERY_OPENAI_BASE_URL")

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
	if s := strings.TrimSpace(os.Getenv("BRIDGR_DISCOVERY_REQUIRE_CANONICAL_CV")); s != "" {
		opts.DiscoveryRequireCanonicalCV = parseBool(s)
	}
	if s := strings.TrimSpace(os.Getenv("BRIDGR_WORKER_DISCOVERY_REQUIRE_CANONICAL_CV")); s != "" {
		opts.DiscoveryRequireCanonicalCV = parseBool(s)
	}
	if s := strings.TrimSpace(os.Getenv("BRIDGR_DISCOVERY_LLM_PREFILTER_ENABLED")); s != "" {
		opts.DiscoveryLLMPrefilterEnabled = parseBool(s)
	}
	if s := strings.TrimSpace(os.Getenv("BRIDGR_DISCOVERY_LLM_SCORE_ENABLED")); s != "" {
		opts.DiscoveryLLMScoreEnabled = parseBool(s)
	}
	if s := strings.TrimSpace(os.Getenv("BRIDGR_DISCOVERY_LLM_QUERY_REFINE_ENABLED")); s != "" {
		opts.DiscoveryLLMQueryRefineEnabled = parseBool(s)
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
