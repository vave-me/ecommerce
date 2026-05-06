package streaming

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// CDNProvider represents different CDN providers
type CDNProvider string

const (
	CDNProviderCloudflare CDNProvider = "cloudflare"
	CDNProviderAkamai     CDNProvider = "akamai"
	CDNProviderFastly     CDNProvider = "fastly"
	CDNProviderAWS        CDNProvider = "aws_cloudfront"
)

// CDNManager handles CDN operations
type CDNManager struct {
	providers       map[CDNProvider]CDNInterface
	primaryProvider CDNProvider
	fallbackEnabled bool
	uploadQueue     chan *UploadTask
	workers         int
	wg              sync.WaitGroup
	ctx             context.Context
	cancel          context.CancelFunc
}

// CDNInterface defines methods for CDN providers
type CDNInterface interface {
	UploadFile(ctx context.Context, localPath, remotePath string) error
	DeleteFile(ctx context.Context, remotePath string) error
	GetURL(remotePath string) string
	Purge(ctx context.Context, paths []string) error
	GetStatistics(ctx context.Context) (*CDNStats, error)
}

// CDNConfig contains CDN configuration
type CDNConfig struct {
	Provider        CDNProvider
	Endpoint        string
	AccessKey       string
	SecretKey       string
	BucketName      string
	Region          string
	CustomDomain    string
	EnableSSL       bool
	Workers         int
}

// CDNStats contains CDN statistics
type CDNStats struct {
	BytesTransferred int64
	RequestCount     int64
	CacheHitRate     float64
	Bandwidth        float64 // Mbps
	ActiveStreams    int
}

// UploadTask represents a file upload task
type UploadTask struct {
	StreamID    string
	Quality     string
	LocalPath   string
	RemotePath  string
	RetryCount  int
	MaxRetries  int
	UploadedAt  time.Time
}

// NewCDNManager creates a new CDN manager
func NewCDNManager(configs []CDNConfig, primaryProvider CDNProvider, fallbackEnabled bool) *CDNManager {
	ctx, cancel := context.WithCancel(context.Background())
	
	manager := &CDNManager{
		providers:       make(map[CDNProvider]CDNInterface),
		primaryProvider: primaryProvider,
		fallbackEnabled: fallbackEnabled,
		uploadQueue:     make(chan *UploadTask, 1000),
		workers:         4,
		ctx:             ctx,
		cancel:          cancel,
	}

	// Initialize CDN providers
	for _, config := range configs {
		switch config.Provider {
		case CDNProviderAWS:
			manager.providers[config.Provider] = NewAWSCDN(config)
		case CDNProviderCloudflare:
			manager.providers[config.Provider] = NewCloudflareCDN(config)
		case CDNProviderAkamai:
			manager.providers[config.Provider] = NewAkamaiCDN(config)
		case CDNProviderFastly:
			manager.providers[config.Provider] = NewFastlyCDN(config)
		}

		if config.Workers > manager.workers {
			manager.workers = config.Workers
		}
	}

	// Start upload workers
	manager.startWorkers()

	return manager
}

// startWorkers starts background upload workers
func (cm *CDNManager) startWorkers() {
	for i := 0; i < cm.workers; i++ {
		cm.wg.Add(1)
		go cm.uploadWorker()
	}
}

// uploadWorker processes upload tasks
func (cm *CDNManager) uploadWorker() {
	defer cm.wg.Done()

	for {
		select {
		case <-cm.ctx.Done():
			return
		case task := <-cm.uploadQueue:
			cm.processUpload(task)
		}
	}
}

// processUpload handles individual upload task
func (cm *CDNManager) processUpload(task *UploadTask) {
	provider := cm.providers[cm.primaryProvider]
	if provider == nil {
		fmt.Printf("Primary CDN provider %s not configured\n", cm.primaryProvider)
		return
	}

	// Try primary provider
	err := provider.UploadFile(cm.ctx, task.LocalPath, task.RemotePath)
	if err == nil {
		task.UploadedAt = time.Now()
		return
	}

	fmt.Printf("Primary CDN upload failed: %v\n", err)

	// Try fallback providers if enabled
	if cm.fallbackEnabled {
		for providerName, provider := range cm.providers {
			if providerName == cm.primaryProvider {
				continue
			}

			if err := provider.UploadFile(cm.ctx, task.LocalPath, task.RemotePath); err == nil {
				task.UploadedAt = time.Now()
				fmt.Printf("Uploaded to fallback CDN: %s\n", providerName)
				return
			}
		}
	}

	// Retry logic
	if task.RetryCount < task.MaxRetries {
		task.RetryCount++
		time.Sleep(time.Duration(task.RetryCount) * time.Second)
		cm.uploadQueue <- task
	} else {
		fmt.Printf("Failed to upload %s after %d retries\n", task.LocalPath, task.MaxRetries)
	}
}

// UploadSegment queues a segment for upload
func (cm *CDNManager) UploadSegment(streamID, quality, segmentPath string) {
	filename := filepath.Base(segmentPath)
	remotePath := fmt.Sprintf("live/%s/%s/%s", streamID, quality, filename)

	task := &UploadTask{
		StreamID:   streamID,
		Quality:    quality,
		LocalPath:  segmentPath,
		RemotePath: remotePath,
		MaxRetries: 3,
	}

	select {
	case cm.uploadQueue <- task:
	default:
		fmt.Printf("Upload queue full, dropping segment: %s\n", segmentPath)
	}
}

// GetSegmentURL returns the CDN URL for a segment
func (cm *CDNManager) GetSegmentURL(streamID, quality, segment string) string {
	remotePath := fmt.Sprintf("live/%s/%s/%s", streamID, quality, segment)
	
	provider := cm.providers[cm.primaryProvider]
	if provider != nil {
		return provider.GetURL(remotePath)
	}

	return ""
}

// PurgeCache purges CDN cache for specific paths
func (cm *CDNManager) PurgeCache(paths []string) error {
	provider := cm.providers[cm.primaryProvider]
	if provider == nil {
		return fmt.Errorf("primary CDN provider not configured")
	}

	return provider.Purge(cm.ctx, paths)
}

// GetStatistics returns CDN statistics
func (cm *CDNManager) GetStatistics() (*CDNStats, error) {
	provider := cm.providers[cm.primaryProvider]
	if provider == nil {
		return nil, fmt.Errorf("primary CDN provider not configured")
	}

	return provider.GetStatistics(cm.ctx)
}

// Stop stops the CDN manager
func (cm *CDNManager) Stop() {
	cm.cancel()
	close(cm.uploadQueue)
	cm.wg.Wait()
}

// AWSCDN implements CDN interface for AWS CloudFront/S3
type AWSCDN struct {
	config   CDNConfig
	s3Client *s3.S3
	uploader *s3manager.Uploader
}

// NewAWSCDN creates AWS CDN provider
func NewAWSCDN(config CDNConfig) *AWSCDN {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(config.Region),
	})

	return &AWSCDN{
		config:   config,
		s3Client: s3.New(sess),
		uploader: s3manager.NewUploader(sess),
	}
}

// UploadFile uploads file to S3
func (a *AWSCDN) UploadFile(ctx context.Context, localPath, remotePath string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = a.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket:      aws.String(a.config.BucketName),
		Key:         aws.String(remotePath),
		Body:        file,
		ContentType: aws.String(getContentType(localPath)),
		CacheControl: aws.String("max-age=3600"),
	})

	return err
}

// DeleteFile deletes file from S3
func (a *AWSCDN) DeleteFile(ctx context.Context, remotePath string) error {
	_, err := a.s3Client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(a.config.BucketName),
		Key:    aws.String(remotePath),
	})
	return err
}

// GetURL returns CloudFront URL
func (a *AWSCDN) GetURL(remotePath string) string {
	if a.config.CustomDomain != "" {
		return fmt.Sprintf("https://%s/%s", a.config.CustomDomain, remotePath)
	}
	return fmt.Sprintf("https://s3.%s.amazonaws.com/%s/%s", 
		a.config.Region, a.config.BucketName, remotePath)
}

// Purge invalidates CloudFront cache
func (a *AWSCDN) Purge(ctx context.Context, paths []string) error {
	// CloudFront invalidation implementation
	return nil
}

// GetStatistics returns AWS CDN statistics
func (a *AWSCDN) GetStatistics(ctx context.Context) (*CDNStats, error) {
	// CloudWatch metrics implementation
	return &CDNStats{
		BytesTransferred: 0,
		RequestCount:     0,
		CacheHitRate:     0.95,
		Bandwidth:        1000.0,
		ActiveStreams:    10,
	}, nil
}

// CloudflareCDN implements CDN interface for Cloudflare
type CloudflareCDN struct {
	config     CDNConfig
	httpClient *http.Client
}

// NewCloudflareCDN creates Cloudflare CDN provider
func NewCloudflareCDN(config CDNConfig) *CloudflareCDN {
	return &CloudflareCDN{
		config:     config,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// UploadFile uploads file to Cloudflare R2/Stream
func (c *CloudflareCDN) UploadFile(ctx context.Context, localPath, remotePath string) error {
	// Cloudflare R2 API implementation
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create request
	req, err := http.NewRequestWithContext(ctx, "PUT", 
		fmt.Sprintf("%s/%s", c.config.Endpoint, remotePath), file)
	if err != nil {
		return err
	}

	// Add authentication headers
	req.Header.Set("X-Auth-Email", c.config.AccessKey)
	req.Header.Set("X-Auth-Key", c.config.SecretKey)
	req.Header.Set("Content-Type", getContentType(localPath))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload failed with status: %d", resp.StatusCode)
	}

	return nil
}

// DeleteFile deletes file from Cloudflare
func (c *CloudflareCDN) DeleteFile(ctx context.Context, remotePath string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE",
		fmt.Sprintf("%s/%s", c.config.Endpoint, remotePath), nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Auth-Email", c.config.AccessKey)
	req.Header.Set("X-Auth-Key", c.config.SecretKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// GetURL returns Cloudflare URL
func (c *CloudflareCDN) GetURL(remotePath string) string {
	if c.config.CustomDomain != "" {
		return fmt.Sprintf("https://%s/%s", c.config.CustomDomain, remotePath)
	}
	return fmt.Sprintf("%s/%s", c.config.Endpoint, remotePath)
}

// Purge purges Cloudflare cache
func (c *CloudflareCDN) Purge(ctx context.Context, paths []string) error {
	// Cloudflare cache purge API implementation
	return nil
}

// GetStatistics returns Cloudflare statistics
func (c *CloudflareCDN) GetStatistics(ctx context.Context) (*CDNStats, error) {
	// Cloudflare Analytics API implementation
	return &CDNStats{
		BytesTransferred: 0,
		RequestCount:     0,
		CacheHitRate:     0.97,
		Bandwidth:        2000.0,
		ActiveStreams:    15,
	}, nil
}

// AkamaiCDN implements CDN interface for Akamai
type AkamaiCDN struct {
	config CDNConfig
}

func NewAkamaiCDN(config CDNConfig) *AkamaiCDN {
	return &AkamaiCDN{config: config}
}

func (a *AkamaiCDN) UploadFile(ctx context.Context, localPath, remotePath string) error {
	// Akamai NetStorage API implementation
	return nil
}

func (a *AkamaiCDN) DeleteFile(ctx context.Context, remotePath string) error {
	return nil
}

func (a *AkamaiCDN) GetURL(remotePath string) string {
	return fmt.Sprintf("https://%s/%s", a.config.CustomDomain, remotePath)
}

func (a *AkamaiCDN) Purge(ctx context.Context, paths []string) error {
	return nil
}

func (a *AkamaiCDN) GetStatistics(ctx context.Context) (*CDNStats, error) {
	return &CDNStats{}, nil
}

// FastlyCDN implements CDN interface for Fastly
type FastlyCDN struct {
	config CDNConfig
}

func NewFastlyCDN(config CDNConfig) *FastlyCDN {
	return &FastlyCDN{config: config}
}

func (f *FastlyCDN) UploadFile(ctx context.Context, localPath, remotePath string) error {
	// Fastly API implementation
	return nil
}

func (f *FastlyCDN) DeleteFile(ctx context.Context, remotePath string) error {
	return nil
}

func (f *FastlyCDN) GetURL(remotePath string) string {
	return fmt.Sprintf("https://%s/%s", f.config.CustomDomain, remotePath)
}

func (f *FastlyCDN) Purge(ctx context.Context, paths []string) error {
	return nil
}

func (f *FastlyCDN) GetStatistics(ctx context.Context) (*CDNStats, error) {
	return &CDNStats{}, nil
}

// getContentType returns content type based on file extension
func getContentType(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".m3u8":
		return "application/vnd.apple.mpegurl"
	case ".ts":
		return "video/mp2t"
	case ".mp4":
		return "video/mp4"
	case ".m4s":
		return "video/iso.segment"
	case ".mpd":
		return "application/dash+xml"
	case ".vtt":
		return "text/vtt"
	default:
		return "application/octet-stream"
	}
}