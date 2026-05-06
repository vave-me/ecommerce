package infrastructure

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"middleman/streams/internal/domain"
	
	"go.uber.org/zap"
)

// WebhookClient handles webhook delivery
type WebhookClient struct {
	httpClient *http.Client
	logger     *zap.Logger
}

// NewWebhookClient creates a new webhook client
func NewWebhookClient(logger *zap.Logger) *WebhookClient {
	// Configure HTTP client with reasonable timeouts
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   10,
		},
	}

	return &WebhookClient{
		httpClient: httpClient,
		logger:     logger,
	}
}

// DeliverWebhook delivers a webhook event to a subscription
func (c *WebhookClient) DeliverWebhook(ctx context.Context, subscription *domain.WebhookSubscription, event *domain.WebhookEvent) (*WebhookDeliveryResult, error) {
	// Validate URL
	if err := c.validateURL(subscription.URL); err != nil {
		return nil, fmt.Errorf("invalid webhook URL: %w", err)
	}

	// Marshal payload
	payload, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	// Generate signature
	signature := c.generateSignature(payload, subscription.Secret)

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, subscription.URL, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create webhook request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "StreamingService-Webhook/1.0")
	req.Header.Set("X-Webhook-ID", event.ID)
	req.Header.Set("X-Webhook-Event", string(event.Type))
	req.Header.Set("X-Webhook-Signature", signature)
	req.Header.Set("X-Webhook-Timestamp", event.Timestamp.Format(time.RFC3339))

	// Add custom headers from subscription
	for key, value := range subscription.Headers {
		req.Header.Set(key, value)
	}

	// Log delivery attempt
	c.logger.Info("Delivering webhook",
		zap.String("subscription_id", subscription.ID),
		zap.String("event_type", event.Type),
		zap.String("url", subscription.URL))

	// Execute request
	startTime := time.Now()
	resp, err := c.httpClient.Do(req)
	duration := time.Since(startTime)

	result := &WebhookDeliveryResult{
		Duration: duration,
	}

	if err != nil {
		result.Error = err.Error()
		c.logger.Error("Webhook delivery failed",
			zap.String("subscription_id", subscription.ID),
			zap.Error(err),
			zap.Duration("duration", duration))
		return result, err
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode

	// Read response body (limited to prevent memory issues)
	bodyReader := io.LimitReader(resp.Body, 1024*1024) // 1MB limit
	responseBody, err := io.ReadAll(bodyReader)
	if err != nil {
		c.logger.Warn("Failed to read webhook response body",
			zap.String("subscription_id", subscription.ID),
			zap.Error(err))
	} else {
		result.ResponseBody = string(responseBody)
	}

	// Check status code
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		c.logger.Info("Webhook delivered successfully",
			zap.String("subscription_id", subscription.ID),
			zap.Int("status_code", resp.StatusCode),
			zap.Duration("duration", duration))
		return result, nil
	}

	// Non-2xx status code
	result.Error = fmt.Sprintf("webhook returned status %d", resp.StatusCode)
	c.logger.Warn("Webhook delivery returned non-success status",
		zap.String("subscription_id", subscription.ID),
		zap.Int("status_code", resp.StatusCode),
		zap.String("response_body", result.ResponseBody),
		zap.Duration("duration", duration))

	return result, fmt.Errorf("webhook returned status %d", resp.StatusCode)
}

// generateSignature generates HMAC-SHA256 signature for webhook payload
func (c *WebhookClient) generateSignature(payload []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	return "sha256=" + hex.EncodeToString(h.Sum(nil))
}

// validateURL validates webhook URL
func (c *WebhookClient) validateURL(webhookURL string) error {
	u, err := url.Parse(webhookURL)
	if err != nil {
		return err
	}

	// Enforce HTTPS in production
	if os.Getenv("ENVIRONMENT") != "development" && u.Scheme != "https" {
		return fmt.Errorf("HTTPS required for webhooks in production")
	}

	// Validate scheme
	if u.Scheme != "https" && u.Scheme != "http" {
		return fmt.Errorf("invalid scheme: %s", u.Scheme)
	}

	// Parse IP to check if it's private
	host := u.Hostname()
	if ip := net.ParseIP(host); ip != nil {
		// Block private IP ranges
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			return fmt.Errorf("webhook URL points to private/local IP address")
		}
		
		// Block multicast and unspecified addresses
		if ip.IsMulticast() || ip.IsUnspecified() {
			return fmt.Errorf("webhook URL points to invalid IP address")
		}
	}

	// Block common private hostnames
	privateHostnames := []string{
		"localhost", "127.0.0.1", "0.0.0.0", "::1",
		"metadata.google.internal", // GCP metadata
		"169.254.169.254", // AWS metadata
		"metadata.azure.internal", // Azure metadata
	}
	
	for _, private := range privateHostnames {
		if host == private {
			return fmt.Errorf("webhook URL points to private hostname: %s", host)
		}
	}
	
	// Allow localhost only in development
	if os.Getenv("ENVIRONMENT") == "development" {
		if host == "localhost" || host == "127.0.0.1" || host == "::1" {
			return nil
		}
	}

	return nil
}

// TestWebhook sends a test webhook to verify configuration
func (c *WebhookClient) TestWebhook(ctx context.Context, subscription *domain.WebhookSubscription) error {
	testEvent := &domain.WebhookEvent{
		ID:        "test_" + fmt.Sprintf("%d", time.Now().Unix()),
		Type:      "webhook.test",
		StreamID:  "test_stream",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"message": "This is a test webhook delivery",
			"subscription_id": subscription.ID,
		},
	}

	_, err := c.DeliverWebhook(ctx, subscription, testEvent)
	return err
}

// WebhookDeliveryResult contains the result of a webhook delivery attempt
type WebhookDeliveryResult struct {
	StatusCode   int
	ResponseBody string
	Error        string
	Duration     time.Duration
}