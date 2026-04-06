package dependencies

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/Kanishkmittal55/bridgr-api/internal/cloud"
	"github.com/Kanishkmittal55/bridgr-api/internal/config"
	"github.com/Kanishkmittal55/bridgr-api/internal/logger"
	"github.com/Kanishkmittal55/bridgr-api/internal/radar"
	"github.com/Kanishkmittal55/bridgr-api/internal/rdbms"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository/sqlc"
)

// Deps are bridgr-worker runtime dependencies (DB, SQS consumer, optional Radar JobSearch).
type Deps struct {
	Repo      *repository.Repo
	HsQuerier sqlc.Querier
	SQSClient *sqs.Client
	JobSearch *radar.JobSearchClient
}

// New wires the worker pool, repository, SQS client, and optional Radar gRPC client.
// cleanup closes Radar (if any) then the DB pool; call once on shutdown.
func New(ctx context.Context, cfg *config.Config) (*Deps, func(), error) {
	pool, err := rdbms.NewConn(rdbms.ConnStr(cfg))
	if err != nil {
		return nil, nil, fmt.Errorf("database: %w", err)
	}

	q := sqlc.New(pool)
	repo := repository.New()

	sqsClient, err := cloud.NewSQSClient(ctx, cfg)
	if err != nil {
		pool.Close()
		return nil, nil, fmt.Errorf("sqs client: %w", err)
	}

	var jobSearch *radar.JobSearchClient
	if cfg.RadarAddr != "" {
		js, jerr := radar.NewJobSearchClient(radar.JobSearchConfig{Addr: cfg.RadarAddr})
		if jerr != nil {
			logger.Get(ctx).Warnw("bridgr-worker: Radar JobSearch client disabled", "addr", cfg.RadarAddr, "error", jerr)
		} else {
			jobSearch = js
			logger.Get(ctx).Infow("bridgr-worker: Radar JobSearch client enabled", "addr", cfg.RadarAddr)
		}
	} else {
		logger.Get(ctx).Warnw("bridgr-worker: RADAR_ADDR empty; job discovery will not call FindJobs")
	}

	d := &Deps{
		Repo:      repo,
		HsQuerier: q,
		SQSClient: sqsClient,
		JobSearch: jobSearch,
	}

	cleanup := func() {
		if jobSearch != nil {
			if cErr := jobSearch.Close(); cErr != nil {
				logger.Get(ctx).Warnw("bridgr-worker: Radar client close", "error", cErr)
			}
		}
		pool.Close()
	}

	return d, cleanup, nil
}
