package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	hshttp "github.com/hassleskip/hassle-go/pkg/http"
	authmiddleware "github.com/hassleskip/hassle-go/pkg/middleware/auth"

	"github.com/hassleskip/bridgr-api/internal/api/deps"
	"github.com/hassleskip/bridgr-api/internal/cloud"
	"github.com/hassleskip/bridgr-api/internal/config"
	"github.com/hassleskip/bridgr-api/internal/logger"
	"github.com/hassleskip/bridgr-api/internal/rdbms"
	"github.com/hassleskip/bridgr-api/internal/repository"
	"github.com/hassleskip/bridgr-api/internal/repository/sqlc"
)

// Run boots the HTTP server and blocks until ctx is cancelled, then shuts down gracefully
// (same role as users/internal/api/server.go).
func Run(ctx context.Context) error {
	cfg := config.Load()
	log := logger.Get(ctx)

	pool, err := rdbms.NewConn(rdbms.ConnStr(cfg))
	if err != nil {
		return fmt.Errorf("database: %w", err)
	}
	defer pool.Close()

	hsQ := sqlc.New(pool)
	repo := repository.New()
	rw := hshttp.NewResponseWriterWithCapture(logger.Get(ctx), cfg.CaptureTestOutput)

	sqsClient, err := createSQSClient(ctx, cfg)
	if err != nil {
		return fmt.Errorf("sqs client: %w", err)
	}

	d := &deps.Deps{
		ResponseWriter:  rw,
		HsQuerier:       hsQ,
		Repo:            repo,
		SQSClient:       sqsClient,
		AccessToApiKeys: buildAccessMap(cfg),
		S3:              cloud.NewClient(cfg),
	}

	addr := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{Addr: addr, Handler: Routes(d)}

	errCh := make(chan error, 1)
	go func() {
		log.Infow("bridgr-api listening", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown: %w", err)
		}
		if err := <-errCh; err != nil {
			return err
		}
		log.Info("server stopped")
		return nil
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("listen: %w", err)
		}
		return nil
	}
}

func buildAccessMap(cfg *config.Config) map[authmiddleware.Access][]string {
	m := map[authmiddleware.Access][]string{}
	add := func(a authmiddleware.Access, keys ...string) {
		for _, k := range keys {
			if k == "" {
				continue
			}
			m[a] = append(m[a], k)
		}
	}
	add(authmiddleware.AccessRead, cfg.ReadAPIKey, cfg.AllAccessAPIKey)
	add(authmiddleware.AccessWrite, cfg.WriteAPIKey, cfg.AllAccessAPIKey)
	return m
}

func createSQSClient(ctx context.Context, cfg *config.Config) (*sqs.Client, error) {
	if cfg.SQSEndpoint != "" {
		awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
			awsconfig.WithRegion(cfg.AWSRegion),
			awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
		)
		if err != nil {
			return nil, err
		}
		return sqs.NewFromConfig(awsCfg, func(o *sqs.Options) {
			o.BaseEndpoint = aws.String(cfg.SQSEndpoint)
		}), nil
	}
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(cfg.AWSRegion))
	if err != nil {
		return nil, err
	}
	return sqs.NewFromConfig(awsCfg), nil
}
