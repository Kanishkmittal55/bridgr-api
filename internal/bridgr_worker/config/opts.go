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
}
