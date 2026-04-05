package radar

import (
	"context"
	"time"

	jobsearchv1 "github.com/Kanishkmittal55/bridgr-api/internal/gen/radar/services/job_search/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	// DefaultJobSearchTimeout caps a single FindJobs RPC when no custom timeout is set.
	DefaultJobSearchTimeout = 10 * time.Minute
)

// JobSearchClient is a thin gRPC client for Radar JobSearchService (FindJobs, etc.).
type JobSearchClient struct {
	conn    *grpc.ClientConn
	client  jobsearchv1.JobSearchServiceClient
	timeout time.Duration
}

// JobSearchConfig configures JobSearchClient.
type JobSearchConfig struct {
	Addr    string
	Timeout time.Duration // per-RPC; defaults to DefaultJobSearchTimeout
}

// NewJobSearchClient dials Radar at addr (e.g. RADAR_ADDR or localhost:50051).
func NewJobSearchClient(cfg JobSearchConfig) (*JobSearchClient, error) {
	addr := cfg.Addr
	if addr == "" {
		addr = "localhost:50051"
	}
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = DefaultJobSearchTimeout
	}
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &JobSearchClient{
		conn:    conn,
		client:  jobsearchv1.NewJobSearchServiceClient(conn),
		timeout: timeout,
	}, nil
}

// Close releases the underlying connection.
func (c *JobSearchClient) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

// FindJobs calls JobSearchService.FindJobs with the configured timeout.
func (c *JobSearchClient) FindJobs(ctx context.Context, req *jobsearchv1.FindJobsRequest) (*jobsearchv1.FindJobsResponse, error) {
	if c == nil || c.client == nil {
		return nil, nil
	}
	callCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.FindJobs(callCtx, req)
}
