package config

// WorkerOpts configures SQS polling and the optional discovery scheduler.
// QueueURL is injected at runtime (from global config / secrets), not from YAML.
type WorkerOpts struct {
	QueueURL                          string `mapstructure:"-"`
	PollingEnabled                    bool   `mapstructure:"polling_enabled"`
	MaxMessages                       int32  `mapstructure:"max_messages"`
	WaitTimeSeconds                   int32  `mapstructure:"wait_time_seconds"`
	VisibilityTimeout                 int32  `mapstructure:"visibility_timeout"`
	PollErrorBackoffSec               int    `mapstructure:"poll_error_backoff_sec"`
	DiscoverySchedulerEnabled         bool   `mapstructure:"discovery_scheduler_enabled"`
	DiscoverySchedulerTriggerInterval int    `mapstructure:"discovery_scheduler_trigger_interval"` // seconds between ticks

	// Discovery pipeline (prefilter / score / LLM phases — placeholders until wired).
	DiscoveryRequireCanonicalCV bool    `mapstructure:"discovery_require_canonical_cv"`
	DiscoveryScoreThreshold     float64 `mapstructure:"discovery_score_threshold"`       // 0–1 composite score gate (placeholder)
	DiscoveryPrefilterMinScore  float64 `mapstructure:"discovery_prefilter_min_score"`   // placeholder for LLM prefilter confidence
	DiscoveryMaxIterations      int32   `mapstructure:"discovery_max_iterations"`        // FindJobs refinement loop cap (placeholder)
	DiscoveryMaxLLMCallsPerRun  int32   `mapstructure:"discovery_max_llm_calls_per_run"` // cost cap (placeholder)
	// DiscoveryIterationTimeBudgetSec caps total wall time for the iterative FindJobs loop (0 = no limit).
	DiscoveryIterationTimeBudgetSec int `mapstructure:"discovery_iteration_time_budget_sec"`
	// DiscoveryFindJobsPageCap max jobs to request per FindJobs call per iteration (0 = default 20).
	DiscoveryFindJobsPageCap int32 `mapstructure:"discovery_findjobs_page_cap"`

	DiscoveryLLMPrefilterEnabled   bool   `mapstructure:"discovery_llm_prefilter_enabled"`
	DiscoveryLLMScoreEnabled       bool   `mapstructure:"discovery_llm_score_enabled"`        // optional refinement + explainability; shares max LLM budget with prefilter
	DiscoveryLLMQueryRefineEnabled bool   `mapstructure:"discovery_llm_query_refine_enabled"` // optional search_query / target_roles refinement between iterations
	DiscoveryOpenAIModel           string `mapstructure:"discovery_openai_model"`
	DiscoveryOpenAIBaseURL         string `mapstructure:"discovery_openai_base_url"` // optional; default OpenAI v1
}
