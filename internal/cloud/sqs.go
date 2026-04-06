package cloud

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/Kanishkmittal55/bridgr-api/internal/config"
)

// NewSQSClient builds an AWS SDK v2 SQS client from config (LocalStack endpoint or default AWS).
func NewSQSClient(ctx context.Context, cfg *config.Config) (*sqs.Client, error) {
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
