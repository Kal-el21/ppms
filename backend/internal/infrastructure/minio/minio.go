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
	if err := c.ensureBucket(ctx); err != nil {
		return err
	}

	_, err := c.client.PutObject(ctx, c.bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (c *Client) ensureBucket(ctx context.Context) error {
	exists, err := c.client.BucketExists(ctx, c.bucketName)
	if err != nil {
		return fmt.Errorf("failed to check minio bucket %q: %w", c.bucketName, err)
	}
	if exists {
		return nil
	}

	if err := c.client.MakeBucket(ctx, c.bucketName, minio.MakeBucketOptions{}); err != nil {
		// Another request/process may have created it between BucketExists and MakeBucket.
		exists, checkErr := c.client.BucketExists(ctx, c.bucketName)
		if checkErr == nil && exists {
			return nil
		}
		return fmt.Errorf("failed to create minio bucket %q: %w", c.bucketName, err)
	}

	return nil
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
