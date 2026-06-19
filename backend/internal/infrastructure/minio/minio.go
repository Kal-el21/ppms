package minio

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/Kal-el21/backend/configs"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	client     *minio.Client
	bucketName string
}

func NewClient(cfg *configs.Config) (*Client, error) {
	client, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	return &Client{client: client, bucketName: cfg.MinioBucket}, nil
}

func (c *Client) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error {
	_, err := c.client.PutObject(ctx, c.bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// GetPresignedDownloadURL menghasilkan URL sementara (default 15 menit) untuk download
// langsung dari MinIO, menghindari proxy file besar lewat backend Go.
func (c *Client) GetPresignedDownloadURL(ctx context.Context, objectName string) (string, error) {
	url, err := c.client.PresignedGetObject(ctx, c.bucketName, objectName, 15*time.Minute, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (c *Client) Delete(ctx context.Context, objectName string) error {
	return c.client.RemoveObject(ctx, c.bucketName, objectName, minio.RemoveObjectOptions{})
}
