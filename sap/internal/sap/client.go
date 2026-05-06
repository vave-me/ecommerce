package sap

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// SAPClient handles communication with SAP system
type SAPClient struct {
	httpClient     *http.Client
	baseURL        string
	apiKey         string
	webhookSecret  string
	clientID       string
	clientSecret   string
}

// NewSAPClient creates a new SAP client instance
func NewSAPClient(baseURL, apiKey, webhookSecret, clientID, clientSecret string) *SAPClient {
	return &SAPClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:       baseURL,
		apiKey:        apiKey,
		webhookSecret: webhookSecret,
		clientID:      clientID,
		clientSecret:  clientSecret,
	}
}

// ValidateSignature validates SAP webhook signature
func (c *SAPClient) ValidateSignature(payload []byte, signature string) error {
	// SAP typically uses HMAC-SHA256 for webhook signatures
	// Implementation depends on SAP's specific webhook security model
	
	if signature == "" {
		return fmt.Errorf("missing signature")
	}
	
	// TODO: Implement actual signature validation based on SAP's webhook security
	// For now, we'll do a simple check
	if signature != c.webhookSecret {
		return fmt.Errorf("invalid signature")
	}
	
	return nil
}

// ParseWebhookEvent parses SAP webhook event from raw payload
func (c *SAPClient) ParseWebhookEvent(payload []byte) (*SAPEvent, error) {
	var event SAPEvent
	
	// Try JSON first
	if err := json.Unmarshal(payload, &event); err == nil {
		return &event, nil
	}
	
	// Try XML (IDoc format)
	if err := xml.Unmarshal(payload, &event); err == nil {
		return &event, nil
	}
	
	return nil, fmt.Errorf("unable to parse webhook payload")
}

// GetProductChanges polls SAP for product changes since last sync
func (c *SAPClient) GetProductChanges(ctx context.Context, since time.Time) ([]*ProductChange, error) {
	url := fmt.Sprintf("%s/api/products/changes?since=%s", c.baseURL, since.Format(time.RFC3339))
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Accept", "application/json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("SAP API error: status %d, body: %s", resp.StatusCode, string(body))
	}
	
	var changes []*ProductChange
	if err := json.NewDecoder(resp.Body).Decode(&changes); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	
	return changes, nil
}

// GetStockLevels retrieves current stock levels from SAP
func (c *SAPClient) GetStockLevels(ctx context.Context, productIDs []string) ([]*StockLevel, error) {
	payload, err := json.Marshal(map[string]interface{}{
		"productIds": productIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}
	
	url := fmt.Sprintf("%s/api/inventory/levels", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("SAP API error: status %d, body: %s", resp.StatusCode, string(body))
	}
	
	var levels []*StockLevel
	if err := json.NewDecoder(resp.Body).Decode(&levels); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	
	return levels, nil
}

// GetPrices retrieves current prices from SAP
func (c *SAPClient) GetPrices(ctx context.Context, productIDs []string, priceListID string) ([]*Price, error) {
	payload, err := json.Marshal(map[string]interface{}{
		"productIds":  productIDs,
		"priceListId": priceListID,
	})
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}
	
	url := fmt.Sprintf("%s/api/pricing/prices", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("SAP API error: status %d, body: %s", resp.StatusCode, string(body))
	}
	
	var prices []*Price
	if err := json.NewDecoder(resp.Body).Decode(&prices); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	
	return prices, nil
}

// SendOrderToSAP sends order information to SAP
func (c *SAPClient) SendOrderToSAP(ctx context.Context, order *OrderData) error {
	payload, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("marshaling order: %w", err)
	}
	
	url := fmt.Sprintf("%s/api/orders", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("SAP API error: status %d, body: %s", resp.StatusCode, string(body))
	}
	
	log.Info().
		Str("orderId", order.OrderID).
		Msg("Successfully sent order to SAP")
	
	return nil
}