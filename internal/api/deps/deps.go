package deps

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/Kanishkmittal55/bridgr-api/internal/auth"
	"github.com/Kanishkmittal55/bridgr-api/internal/cloud"
	"github.com/Kanishkmittal55/bridgr-api/internal/config"
	"github.com/Kanishkmittal55/bridgr-api/internal/httpx"
	"github.com/Kanishkmittal55/bridgr-api/internal/logger"
	"github.com/Kanishkmittal55/bridgr-api/internal/rdbms"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository/sqlc"
)

// Deps are Bridgr API runtime dependencies (narrow slice of former monolith dependencies).
type Deps struct {
	ResponseWriter  *httpx.ResponseWriter
	HsQuerier       sqlc.Querier
	Repo            *repository.Repo
	SQSClient       *sqs.Client
	AccessToApiKeys map[auth.Access][]string
	S3              cloud.Interface
}

// New wires DB, SQS, S3, API key map, and HTTP response helper. cleanup closes the DB pool and must be called once when the server stops.
func New(ctx context.Context, cfg *config.Config) (*Deps, func(), error) {
	pool, err := rdbms.NewConn(rdbms.ConnStr(cfg))
	if err != nil {
		return nil, nil, fmt.Errorf("database: %w", err)
	}
	cleanup := func() { pool.Close() }

	hsQ := sqlc.New(pool)
	repo := repository.New()
	rw := httpx.NewResponseWriterWithCapture(logger.Get(ctx), cfg.CaptureTestOutput)

	sqsClient, err := cloud.NewSQSClient(ctx, cfg)
	if err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("sqs client: %w", err)
	}

	d := &Deps{
		ResponseWriter:  rw,
		HsQuerier:       hsQ,
		Repo:            repo,
		SQSClient:       sqsClient,
		AccessToApiKeys: buildAccessMap(cfg),
		S3:              cloud.NewClient(cfg),
	}
	return d, cleanup, nil
}

func buildAccessMap(cfg *config.Config) map[auth.Access][]string {
	m := map[auth.Access][]string{}
	add := func(a auth.Access, keys ...string) {
		for _, k := range keys {
			if k == "" {
				continue
			}
			m[a] = append(m[a], k)
		}
	}
	add(auth.AccessRead, cfg.ReadAPIKey, cfg.AllAccessAPIKey)
	add(auth.AccessWrite, cfg.WriteAPIKey, cfg.AllAccessAPIKey)
	return m
}
