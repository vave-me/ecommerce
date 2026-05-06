package minio

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioStorageClient encapsulates interactions with a MinIO (S3-compatible) server.
type MinioStorageClient struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	BucketName string
	UseSSL     bool
	Client     *minio.Client
}

// NewMinioStorageClient initializes a new MinioStorageClient by reading
// environment variables or defaulting them. This is an alternative approach
// if you want a single self-contained client.
func NewMinioStorageClient() (*MinioStorageClient, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT") // e.g. "localhost:9000"
	if endpoint == "" {
		endpoint = "localhost:9096"
	}

	accessKey := os.Getenv("MINIO_ACCESS_KEY_ID")     // or MINIO_ROOT_USER
	secretKey := os.Getenv("MINIO_SECRET_ACCESS_KEY") // or MINIO_ROOT_PASSWORD
	useSSL := false
	if os.Getenv("S3_USE_SSL") == "true" {
		useSSL = true
	}

	bucketName := os.Getenv("MINIO_BUCKET_NAME")
	if bucketName == "" {
		bucketName = "classified"
	}

	c := &MinioStorageClient{
		Endpoint:   endpoint,
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		BucketName: bucketName,
		UseSSL:     useSSL,
	}

	// Initialize actual minio.Client
	minioClient, err := minio.New(c.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.AccessKey, c.SecretKey, ""),
		Secure: c.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	c.Client = minioClient
	return c, nil
}

// EnsureBucketExists ensures the BucketName is present or creates it.
func (c *MinioStorageClient) EnsureBucketExists(ctx context.Context) error {
	exists, err := c.Client.BucketExists(ctx, c.BucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		fmt.Printf("Bucket %s does not exist; creating...\n", c.BucketName)
		err = c.Client.MakeBucket(ctx, c.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket %q: %w", c.BucketName, err)
		}
		fmt.Printf("Created bucket: %s\n", c.BucketName)
	}
	return nil
}

// PresignedPutObjectURL returns a signed URL valid for 15 minutes.
func (c *MinioStorageClient) PresignedPutObjectURL(ctx context.Context, objectName string, expires time.Duration) (string, error) {
	if expires <= 0 {
		expires = 15 * time.Minute
	}
	url, err := c.Client.PresignedPutObject(ctx, c.BucketName, objectName, expires)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned PUT url for %q: %w", objectName, err)
	}
	return url.String(), nil
}

// RemoveObject deletes a single object.
func (c *MinioStorageClient) RemoveObject(ctx context.Context, objectKey string) error {
	err := c.Client.RemoveObject(ctx, c.BucketName, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to remove object %q: %w", objectKey, err)
	}
	return nil
}

// RemovePrefix recursively deletes all objects matching a prefix.
func (c *MinioStorageClient) RemovePrefix(ctx context.Context, prefix string) error {
	objectCh := c.Client.ListObjects(ctx, c.BucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})
	for obj := range objectCh {
		if obj.Err != nil {
			return fmt.Errorf("error listing object with prefix %q: %w", prefix, obj.Err)
		}
		if err := c.Client.RemoveObject(ctx, c.BucketName, obj.Key, minio.RemoveObjectOptions{}); err != nil {
			return fmt.Errorf("failed removing object %s: %w", obj.Key, err)
		}
	}
	return nil
}
