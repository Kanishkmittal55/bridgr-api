package cloud

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/Kanishkmittal55/bridgr-api/internal/config"
	"github.com/Kanishkmittal55/bridgr-api/internal/env"
)

// Interface is the S3 surface Bridgr needs (presigned uploads, reads).
type Interface interface {
	Upload(ctx context.Context, bucket string, key string, body io.Reader) (string, error)
	Download(ctx context.Context, bucket string, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, bucket string, key string) error
	CopyS3Folder(ctx context.Context, sourceBucket, sourcePrefix, destinationBucket, destinationPrefix string) error
	GetPresignedUploadURL(ctx context.Context, bucket, key, contentType string, expiry time.Duration) (string, error)
}

type client struct {
	s3        *s3.Client
	presigner *s3.PresignClient
}

// NewClient builds an S3 client (MinIO in development; AWS in non-development when S3Url is empty).
func NewClient(cfg *config.Config) Interface {
	if env.IsNonDevelopment(cfg.Env) && cfg.S3Url == "" {
		awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion(cfg.AWSRegion))
		if err != nil {
			log.Printf("cloud: aws config: %v", err)
			return nil
		}
		s3Client := s3.NewFromConfig(awsCfg)
		return &client{s3: s3Client, presigner: s3.NewPresignClient(s3Client)}
	}
	if cfg.S3Url == "" {
		return nil
	}
	s3Client := s3.New(s3.Options{
		BaseEndpoint: aws.String(cfg.S3Url),
		Credentials:  credentials.NewStaticCredentialsProvider(cfg.S3User, cfg.S3Password, ""),
		Region:       cfg.AWSRegion,
		UsePathStyle: true,
	})
	external := cfg.S3ExternalUrl

	presignS3Client := s3.New(s3.Options{
		BaseEndpoint: aws.String(external),
		Credentials:  credentials.NewStaticCredentialsProvider(cfg.S3User, cfg.S3Password, ""),
		Region:       cfg.AWSRegion,
		UsePathStyle: true,
	})
	c := &client{s3: s3Client, presigner: s3.NewPresignClient(presignS3Client)}
	if cfg.HassleSkipS3Bucket != "" {
		if err := c.ensureBucket(context.Background(), cfg.HassleSkipS3Bucket); err != nil {
			log.Printf("Warning: ensure bucket %s: %v", cfg.HassleSkipS3Bucket, err)
		}
	}
	return c
}

func (c *client) Upload(ctx context.Context, bucket string, key string, body io.Reader) (string, error) {
	_, err := c.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   body,
	})
	if err != nil {
		return "", fmt.Errorf("put object: %w", err)
	}
	return URI(bucket, key), nil
}

func (c *client) Download(ctx context.Context, bucket string, key string) (io.ReadCloser, error) {
	resp, err := c.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("get object: %w", err)
	}
	return resp.Body, nil
}

func (c *client) Delete(ctx context.Context, bucket string, key string) error {
	_, err := c.s3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("delete object: %w", err)
	}
	return nil
}

func (c *client) CopyS3Folder(ctx context.Context, sourceBucket, sourcePrefix, destinationBucket, destinationPrefix string) error {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(sourceBucket),
		Prefix: aws.String(sourcePrefix),
	}
	paginator := s3.NewListObjectsV2Paginator(c.s3, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("list objects: %w", err)
		}
		for _, object := range page.Contents {
			sourceKey := *object.Key
			destinationKey := destinationPrefix + sourceKey[len(sourcePrefix):]
			_, err := c.s3.CopyObject(ctx, &s3.CopyObjectInput{
				Bucket:     aws.String(destinationBucket),
				CopySource: aws.String(sourceBucket + "/" + sourceKey),
				Key:        aws.String(destinationKey),
			})
			if err != nil {
				return fmt.Errorf("copy %s -> %s: %w", sourceKey, destinationKey, err)
			}
		}
	}
	return nil
}

func URI(bucket, key string) string {
	return fmt.Sprintf("s3://%s/%s", bucket, key)
}

func (c *client) ensureBucket(ctx context.Context, bucket string) error {
	_, err := c.s3.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(bucket)})
	if err == nil {
		return nil
	}
	var notFound *types.NotFound
	if !errors.As(err, &notFound) {
		log.Printf("HeadBucket %s: %v, trying create", bucket, err)
	}
	_, err = c.s3.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(bucket)})
	if err != nil {
		var bucketExists *types.BucketAlreadyExists
		var bucketOwned *types.BucketAlreadyOwnedByYou
		if errors.As(err, &bucketExists) || errors.As(err, &bucketOwned) {
			return nil
		}
		return fmt.Errorf("create bucket: %w", err)
	}
	return nil
}

func (c *client) GetPresignedUploadURL(ctx context.Context, bucket, key, contentType string, expiry time.Duration) (string, error) {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}
	presignedReq, err := c.presigner.PresignPutObject(ctx, input, func(opts *s3.PresignOptions) {
		opts.Expires = expiry
	})
	if err != nil {
		return "", fmt.Errorf("presign put: %w", err)
	}
	return presignedReq.URL, nil
}
