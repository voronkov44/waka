package s3store

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Config struct {
	Endpoint        string // внутренний http://minio:9000 (для API контейнера)
	PresignEndpoint string // внешний  http://localhost:19000 (для ссылок)
	AccessKey       string
	SecretKey       string
	Bucket          string
	Region          string
	PublicBaseURL   string // public режим
}

type Minio struct {
	client        *s3.Client
	presign       *s3.PresignClient
	bucket        string
	publicBaseURL string
}

func NewMinio(ctx context.Context, cfg Config) (*Minio, error) {
	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}
	if cfg.Endpoint == "" || cfg.AccessKey == "" || cfg.SecretKey == "" || cfg.Bucket == "" {
		return nil, fmt.Errorf("minio config invalid")
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		),
	)
	if err != nil {
		return nil, err
	}

	// client для реальных запросов из контейнера
	cli := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(cfg.Endpoint)
	})

	// client для presign
	presignEndpoint := cfg.PresignEndpoint
	if presignEndpoint == "" {
		presignEndpoint = cfg.Endpoint
	}
	presignCli := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(presignEndpoint)
	})

	return &Minio{
		client:        cli,
		presign:       s3.NewPresignClient(presignCli),
		bucket:        cfg.Bucket,
		publicBaseURL: strings.TrimRight(cfg.PublicBaseURL, "/"),
	}, nil
}

func (m *Minio) Put(ctx context.Context, key string, body io.Reader, contentType string) error {
	key = strings.TrimLeft(key, "/")
	_, err := m.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(m.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	return err
}

func (m *Minio) Delete(ctx context.Context, key string) error {
	key = strings.TrimLeft(key, "/")
	_, err := m.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(m.bucket),
		Key:    aws.String(key),
	})
	return err
}

func (m *Minio) PublicURL(key string) string {
	key = strings.TrimLeft(key, "/")
	return fmt.Sprintf("%s/%s/%s", m.publicBaseURL, m.bucket, key)
}

// PresignedGetURL - private bucket: presigned GET
func (m *Minio) PresignedGetURL(ctx context.Context, key string, ttl time.Duration) (string, error) {
	key = strings.TrimLeft(key, "/")
	out, err := m.presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(m.bucket),
		Key:    aws.String(key),
	}, func(po *s3.PresignOptions) {
		po.Expires = ttl
	})
	if err != nil {
		return "", err
	}
	return out.URL, nil
}
