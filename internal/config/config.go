package config

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Kanishkmittal55/bridgr-api/internal/ctxlog"
	"github.com/Kanishkmittal55/bridgr-api/internal/env"
)

// Config holds Bridgr API / worker settings (env-first).
type Config struct {
	Env               env.Env
	LogLevel          ctxlog.Level
	Port              int
	Version           string
	CaptureTestOutput bool

	PostgresUser            string
	PostgresPassword        string
	PostgresHost            string
	PostgresPort            string
	PostgresDb              string
	PostgresSslModeDisabled bool
	PostgresMinIdleConn     int
	PostgresMaxOpenConn     int
	PostgresMaxConnLifetime time.Duration
	PostgresMaxConnIdleTime time.Duration

	S3Url              string
	S3ExternalUrl      string
	S3User             string
	S3Password         string
	HassleSkipS3Bucket string

	AWSRegion      string
	SQSEndpoint    string
	BridgrQueueURL string

	ReadAPIKey      string
	WriteAPIKey     string
	AllAccessAPIKey string

	// Job discovery API / worker
	JobDiscoveryMaxRunsPerHour int
	JobDiscoverySyncInDev      bool

	// RadarAddr gRPC address for the Python radar service (e.g. radar:50051 in Docker).
	RadarAddr string
}

var (
	cfg  *Config
	once sync.Once
)

// Load reads configuration once from environment variables.
func Load() *Config {
	once.Do(func() {
		cfg = loadFromEnv()
	})
	c := *cfg
	return &c
}

func loadFromEnv() *Config {
	c := &Config{
		Env:                        env.ResolveEnvOrDie(),
		LogLevel:                   ctxlog.InfoLevel,
		Port:                       8080,
		PostgresHost:               "localhost",
		PostgresPort:               "5432",
		PostgresUser:               "bridgr",
		PostgresPassword:           "bridgr",
		PostgresDb:                 "bridgr",
		PostgresSslModeDisabled:    true,
		PostgresMinIdleConn:        2,
		PostgresMaxOpenConn:        10,
		PostgresMaxConnLifetime:    time.Hour,
		PostgresMaxConnIdleTime:    30 * time.Minute,
		S3ExternalUrl:              "http://localhost:9000",
		S3User:                     "minioadmin",
		S3Password:                 "minioadmin",
		HassleSkipS3Bucket:         "bridgr",
		AWSRegion:                  "us-east-1",
		AllAccessAPIKey:            "test-all-access-key",
		JobDiscoveryMaxRunsPerHour: 20,
	}

	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		c.LogLevel = ctxlog.DebugLevel
	case "warn":
		c.LogLevel = ctxlog.WarnLevel
	case "error":
		c.LogLevel = ctxlog.ErrorLevel
	default:
		c.LogLevel = ctxlog.InfoLevel
	}
	if v := os.Getenv("PORT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			c.Port = n
		}
	}
	c.Version = os.Getenv("VERSION")
	if os.Getenv("CAPTURE_TEST_OUTPUT") == "true" {
		c.CaptureTestOutput = true
	}

	setStr(&c.PostgresUser, "POSTGRES_USER")
	setStr(&c.PostgresPassword, "POSTGRES_PASSWORD")
	setStr(&c.PostgresHost, "POSTGRES_HOST")
	setStr(&c.PostgresPort, "POSTGRES_PORT")
	setStr(&c.PostgresDb, "POSTGRES_DB")
	if os.Getenv("POSTGRES_SSL_MODE_DISABLED") == "false" {
		c.PostgresSslModeDisabled = false
	}

	setStr(&c.S3Url, "S3_URL")
	setStr(&c.S3ExternalUrl, "S3_EXTERNAL_URL")
	setStr(&c.S3User, "S3_USER")
	setStr(&c.S3Password, "S3_PASSWORD")
	setStr(&c.HassleSkipS3Bucket, "HASSLE_SKIP_S3_BUCKET")

	setStr(&c.AWSRegion, "AWS_REGION")
	setStr(&c.SQSEndpoint, "SQS_ENDPOINT")
	setStr(&c.BridgrQueueURL, "SQS_BRIDGR_QUEUE_URL")

	setStr(&c.ReadAPIKey, "BRIDGR_READ_API_KEY")
	setStr(&c.WriteAPIKey, "BRIDGR_WRITE_API_KEY")
	if v := os.Getenv("BRIDGR_ALL_ACCESS_API_KEY"); v != "" {
		c.AllAccessAPIKey = v
	}
	if c.ReadAPIKey == "" {
		c.ReadAPIKey = os.Getenv("VITE_BRIDGR_API_READ_KEY")
	}
	if c.WriteAPIKey == "" {
		c.WriteAPIKey = os.Getenv("VITE_BRIDGR_API_WRITE_KEY")
	}
	if c.ReadAPIKey == "" && c.AllAccessAPIKey != "" {
		c.ReadAPIKey = c.AllAccessAPIKey
	}
	if c.WriteAPIKey == "" && c.AllAccessAPIKey != "" {
		c.WriteAPIKey = c.AllAccessAPIKey
	}

	if v := os.Getenv("BRIDGR_JOB_DISCOVERY_MAX_RUNS_PER_HOUR"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			c.JobDiscoveryMaxRunsPerHour = n
		}
	}
	c.JobDiscoverySyncInDev = os.Getenv("BRIDGR_JOB_DISCOVERY_SYNC_IN_DEV") == "true"

	setStr(&c.RadarAddr, "RADAR_ADDR")

	return c
}

func setStr(dest *string, key string) {
	if v := os.Getenv(key); v != "" {
		*dest = v
	}
}

// Get is an alias for Load (copy-on-read snapshot).
func Get() *Config {
	return Load()
}
