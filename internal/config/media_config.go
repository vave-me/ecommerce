package config

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MediaConfig holds the configuration for the Media service (MinIO-based).
type MediaConfig struct {
	MinioEndpoint string // e.g. "http://46.4.91.91:9096"
	MinioBucket   string
	ACL           string
	ServerAddress string
	UseSSL        bool
}

// LoadMediaConfig loads configuration for the Media service from environment variables
func LoadMediaConfig() *MediaConfig {
	fmt.Println("Loading environment variables for MediaConfig...")

	// MINIO_ENDPOINT
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	if minioEndpoint == "" {
		minioEndpoint = "minio.sfx-markt.de:9096" // default if not set
		fmt.Println("MINIO_ENDPOINT not set, using default:", minioEndpoint)
	} else {
		fmt.Println("MINIO_ENDPOINT:", minioEndpoint)
	}

	// MINIO_BUCKET
	minioBucket := os.Getenv("MINIO_BUCKET")
	if minioBucket == "" {
		minioBucket = "classified"
		fmt.Println("MINIO_BUCKET not set, using default:", minioBucket)
	} else {
		fmt.Println("MINIO_BUCKET:", minioBucket)
	}

	// S3_ACL (or rename to MINIO_ACL if desired)
	acl := os.Getenv("S3_ACL")
	if acl == "" {
		acl = "private"
		fmt.Println("S3_ACL not set, using default:", acl)
	} else {
		fmt.Println("S3_ACL:", acl)
	}

	// SERVER_ADDRESS
	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = "minio.sfx-markt.de"
		fmt.Println("SERVER_ADDRESS not set, using default:", serverAddress)
	} else {
		fmt.Println("SERVER_ADDRESS:", serverAddress)
	}

	// MINIO_USE_SSL (bool)
	useSslEnv := os.Getenv("MINIO_USE_SSL")
	useSSL := strings.EqualFold(useSslEnv, "true")

	return &MediaConfig{
		MinioEndpoint: minioEndpoint,
		MinioBucket:   minioBucket,
		ACL:           acl,
		ServerAddress: serverAddress,
		UseSSL:        useSSL,
	}
}

// NewMinioClient creates a new MinIO client using minio-go/v7
func NewMinioClient(ctx context.Context, cfg MediaConfig) (*minio.Client, error) {
	fmt.Println("Initializing MinIO client...")

	// Load credentials from environment
	accessKey := os.Getenv("MINIO_ACCESS_KEY_ID") // or MINIO_ROOT_USER
	secretKey := os.Getenv("MINIO_SECRET_ACCESS_KEY")

	if accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("MinIO credentials are not set in env (MINIO_ACCESS_KEY_ID / MINIO_SECRET_ACCESS_KEY)")
	}

	// Create the client
	minioClient, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	// Ensure the bucket exists
	exists, err := minioClient.BucketExists(ctx, cfg.MinioBucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check if bucket %q exists: %w", cfg.MinioBucket, err)
	}
	if !exists {
		fmt.Printf("Bucket %s does not exist, creating...\n", cfg.MinioBucket)
		err = minioClient.MakeBucket(ctx, cfg.MinioBucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket %q: %w", cfg.MinioBucket, err)
		}
		fmt.Printf("Created bucket: %s\n", cfg.MinioBucket)
	}

	fmt.Println("MinIO client initialized successfully.")
	return minioClient, nil
}

// (Optional) If you want a presigned PUT example
func PresignedPutObject(
	ctx context.Context,
	client *minio.Client,
	cfg MediaConfig,
	objectKey string,
) (string, error) {
	expiry := 15 * time.Minute
	url, err := client.PresignedPutObject(ctx, cfg.MinioBucket, objectKey, expiry)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned PUT URL for %q: %w", objectKey, err)
	}
	return url.String(), nil
}
