package config

import (
	"testing"
)

func TestLoad_DevelopmentDefaults(t *testing.T) {
	t.Setenv(EnvVar, "development")
	opts, err := Load("https://sqs.example/queue")
	if err != nil {
		t.Fatal(err)
	}
	if opts.QueueURL != "https://sqs.example/queue" {
		t.Fatalf("QueueURL: got %q", opts.QueueURL)
	}
	if !opts.PollingEnabled {
		t.Fatal("expected polling enabled in development.yaml")
	}
	if opts.MaxMessages != 5 {
		t.Fatalf("max_messages: got %d", opts.MaxMessages)
	}
	if opts.DiscoverySchedulerTriggerInterval != 60 {
		t.Fatalf("discovery_scheduler_trigger_interval: got %d", opts.DiscoverySchedulerTriggerInterval)
	}
}

func TestLoad_DiscoverySchedulerEnvOverrides(t *testing.T) {
	t.Setenv(EnvVar, "development")
	t.Setenv("BRIDGR_DISCOVERY_SCHEDULER_ENABLED", "true")
	t.Setenv("BRIDGR_DISCOVERY_SCHEDULER_TRIGGER_INTERVAL", "120")

	opts, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if !opts.DiscoverySchedulerEnabled {
		t.Fatal("expected discovery scheduler enabled")
	}
	if opts.DiscoverySchedulerTriggerInterval != 120 {
		t.Fatalf("discovery_scheduler_trigger_interval: got %d", opts.DiscoverySchedulerTriggerInterval)
	}
}
