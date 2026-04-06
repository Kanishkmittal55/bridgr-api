package config

import (
	"fmt"
)

func validate(o *WorkerOpts) error {
	if o.MaxMessages < 1 || o.MaxMessages > 10 {
		return fmt.Errorf("max_messages must be between 1 and 10, got %d", o.MaxMessages)
	}
	if o.WaitTimeSeconds < 0 || o.WaitTimeSeconds > 20 {
		return fmt.Errorf("wait_time_seconds must be between 0 and 20 (SQS long poll cap), got %d", o.WaitTimeSeconds)
	}
	if o.VisibilityTimeout < 1 {
		return fmt.Errorf("visibility_timeout must be positive, got %d", o.VisibilityTimeout)
	}
	if o.PollErrorBackoffSec < 1 {
		return fmt.Errorf("poll_error_backoff_sec must be at least 1, got %d", o.PollErrorBackoffSec)
	}
	if o.DiscoverySchedulerEnabled && o.DiscoverySchedulerTriggerInterval < 1 {
		return fmt.Errorf("discovery_scheduler_trigger_interval must be at least 1 when discovery scheduler is enabled, got %d", o.DiscoverySchedulerTriggerInterval)
	}
	return nil
}
