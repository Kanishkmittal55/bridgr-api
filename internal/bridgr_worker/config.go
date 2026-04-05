package bridgr_worker

import (
	"os"
	"strconv"
)

// WorkerOpts tunes SQS polling (env-overridable for tests).
type WorkerOpts struct {
	QueueURL            string
	PollingEnabled      bool
	MaxMessages         int32
	WaitTimeSeconds     int32
	VisibilityTimeout   int32
	PollErrorBackoffSec int
}

const (
	envBridgrPollingEnabled      = "BRIDGR_WORKER_POLLING_ENABLED"
	envBridgrMaxMessages         = "BRIDGR_WORKER_SQS_MAX_MESSAGES"
	envBridgrWaitSeconds         = "BRIDGR_WORKER_SQS_WAIT_SECONDS"
	envBridgrVisibility          = "BRIDGR_WORKER_SQS_VISIBILITY_TIMEOUT"
	envBridgrPollErrorBackoffSec = "BRIDGR_WORKER_POLL_ERROR_BACKOFF_SEC"
)

// WorkerOptsFromEnv returns defaults merged with environment overrides.
func WorkerOptsFromEnv(queueURL string) WorkerOpts {
	o := WorkerOpts{
		QueueURL:            queueURL,
		PollingEnabled:      true,
		MaxMessages:         5,
		WaitTimeSeconds:     20,
		VisibilityTimeout:   300,
		PollErrorBackoffSec: 5,
	}
	if v := os.Getenv(envBridgrPollingEnabled); v != "" {
		o.PollingEnabled = v == "1" || v == "true" || v == "TRUE"
	}
	if v := os.Getenv(envBridgrMaxMessages); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			o.MaxMessages = int32(n)
		}
	}
	if v := os.Getenv(envBridgrWaitSeconds); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			o.WaitTimeSeconds = int32(n)
		}
	}
	if v := os.Getenv(envBridgrVisibility); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			o.VisibilityTimeout = int32(n)
		}
	}
	if v := os.Getenv(envBridgrPollErrorBackoffSec); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			o.PollErrorBackoffSec = n
		}
	}
	return o
}
