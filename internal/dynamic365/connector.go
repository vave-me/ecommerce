package dynamic365

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"middleman/erp/internal/erp"
)

// Connector implements the Dynamics 365 ERP connector with OAuth2 authentication
type Connector struct {
	*erp.BaseConnector

	// OAuth2 state
	accessToken  string
	refreshToken string
	tokenExpiry  time.Time
	tokenMutex   sync.RWMutex

	// Configuration
	tenantID     string
	clientID     string
	clientSecret string
	environment  string // production or sandbox
	companyID    string
}

// NewConnector creates a new Dynamics 365 connector
func NewConnector(config erp.ERPConfig) (*Connector, error) {
	if config.Type != erp.ERPTypeDynamics365 {
		return nil, fmt.Errorf("invalid ERP type: %s, expected %s", config.Type, erp.ERPTypeDynamics365)
	}

	// Validate OAuth2 configuration
	if config.Auth.Type != "oauth2" {
		return nil, fmt.Errorf("Dynamics 365 requires OAuth2 authentication")
	}

	// Extract configuration from metadata
	tenantID, ok := config.Metadata["tenant_id"].(string)
	if !ok || tenantID == "" {
		return nil, fmt.Errorf("Dynamics 365 requires tenant_id in metadata")
	}

	clientID := config.Auth.ClientID
	if clientID == "" {
		return nil, fmt.Errorf("Dynamics 365 requires client_id for OAuth2")
	}

	clientSecret := config.Auth.ClientSecret
	if clientSecret == "" {
		return nil, fmt.Errorf("Dynamics 365 requires client_secret for OAuth2")
	}

	// Optional configurations
	environment := "production"
	if env, ok := config.Metadata["environment"].(string); ok {
		environment = env
	}

	companyID := ""
	if cid, ok := config.Metadata["company_id"].(string); ok {
		companyID = cid
	}

	return &Connector{
		BaseConnector: erp.NewBaseConnector(config),
		tenantID:      tenantID,
		clientID:      clientID,
		clientSecret:  clientSecret,
		environment:   environment,
		companyID:     companyID,
	}, nil
}

// GetType returns the ERP type
func (c *Connector) GetType() erp.ERPType {
	return erp.ERPTypeDynamics365
}

// authenticate performs OAuth2 authentication with Microsoft Identity Platform
func (c *Connector) authenticate(ctx context.Context) error {
	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()

	// Check if we have a valid token
	if c.accessToken != "" && time.Now().Before(c.tokenExpiry.Add(-5*time.Minute)) {
		return nil
	}

	c.Logger.Info().Msg("Authenticating with Dynamics 365")

	// Build token endpoint URL
	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", c.tenantID)

	// Prepare request body
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)
	data.Set("scope", "https://api.businesscentral.dynamics.com/.default")

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("creating auth request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("authentication request failed: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
		Error        string `json:"error"`
		ErrorDesc    string `json:"error_description"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("parsing auth response: %w", err)
	}

	if tokenResp.Error != "" {
		return fmt.Errorf("authentication error: %s - %s", tokenResp.Error, tokenResp.ErrorDesc)
	}

	if tokenResp.AccessToken == "" {
		return fmt.Errorf("no access token received")
	}

	// Update token state
	c.accessToken = tokenResp.AccessToken
	c.refreshToken = tokenResp.RefreshToken
	c.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	c.Logger.Info().
		Str("tenant_id", c.tenantID).
		Time("expires_at", c.tokenExpiry).
		Msg("Successfully authenticated with Dynamics 365")

	return nil
}

// ensureAuthenticated ensures we have valid authentication
func (c *Connector) ensureAuthenticated(ctx context.Context) error {
	c.tokenMutex.RLock()
	needsAuth := c.accessToken == "" || time.Now().After(c.tokenExpiry.Add(-5*time.Minute))
	c.tokenMutex.RUnlock()

	if needsAuth {
		return c.authenticate(ctx)
	}
	return nil
}

// ValidateWebhook validates Dynamics 365 webhook notification
func (c *Connector) ValidateWebhook(ctx context.Context, payload []byte, signature string) error {
	// Dynamics 365 uses notification tokens for validation
	if !c.Config.Webhook.ValidateSign || c.Config.Webhook.ValidationToken == "" {
		return nil
	}

	// Parse the notification
	var notification struct {
		ValidationToken string `json:"validationToken"`
	}

	if err := json.Unmarshal(payload, &notification); err == nil && notification.ValidationToken != "" {
		// This is a validation request
		if notification.ValidationToken != c.Config.Webhook.ValidationToken {
			return fmt.Errorf("invalid validation token")
		}
	}

	// For actual notifications, Dynamics 365 doesn't provide signatures
	// Additional security can be implemented using client certificates
	return nil
}

// ParseWebhook parses Dynamics 365 webhook into canonical event
func (c *Connector) ParseWebhook(ctx context.Context, payload []byte) (*erp.CanonicalEvent, error) {
	var notification D365Notification
	if err := json.Unmarshal(payload, &notification); err != nil {
		return nil, fmt.Errorf("parsing webhook payload: %w", err)
	}

	// Handle validation request
	if notification.ValidationToken != "" {
		return nil, fmt.Errorf("validation request, not an event")
	}

	// Process change notifications
	if len(notification.Value) == 0 {
		return nil, fmt.Errorf("no change notifications in payload")
	}

	// Take the first notification (usually there's only one)
	change := notification.Value[0]

	// Map resource to event type
	eventType, err := mapD365EventType(change.Resource, change.ChangeType)
	if err != nil {
		return nil, err
	}

	// Fetch the actual resource data if needed
	var eventPayload interface{}
	if change.ResourceData != nil {
		eventPayload = change.ResourceData
	} else {
		// Fetch the resource data using the resource URL
		eventPayload, err = c.fetchResourceData(ctx, change.Resource)
		if err != nil {
			c.Logger.Warn().Err(err).Str("resource", change.Resource).Msg("Failed to fetch resource data")
			eventPayload = map[string]interface{}{
				"resource":    change.Resource,
				"change_type": change.ChangeType,
			}
		}
	}

	return &erp.CanonicalEvent{
		EventID:        change.ID,
		EventType:      eventType,
		EventTimestamp: change.SubscriptionExpirationDateTime, // Best approximation
		Source:         fmt.Sprintf("Dynamics365-%s", c.tenantID),
		CorrelationID:  change.ClientState,
		Payload:        eventPayload,
	}, nil
}

// FetchProducts fetches products (items) changed since the given time
func (c *Connector) FetchProducts(ctx context.Context, since time.Time) ([]*erp.ProductPayload, error) {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}

	c.Logger.Info().
		Time("since", since).
		Msg("Fetching products from Dynamics 365")

	// Build query with OData filters
	query := url.Values{}
	query.Set("$filter", fmt.Sprintf("lastModifiedDateTime gt %s", since.Format("2006-01-02T15:04:05Z")))
	query.Set("$top", fmt.Sprint(c.Config.Sync.BatchSize))
	query.Set("$orderby", "lastModifiedDateTime asc")

	// Fetch items
	items, err := c.makeAPIRequest(ctx, "items", query)
	if err != nil {
		return nil, fmt.Errorf("fetching items: %w", err)
	}

	// Transform to canonical format
	products := make([]*erp.ProductPayload, 0)
	if value, ok := items["value"].([]interface{}); ok {
		for _, item := range value {
			if data, ok := item.(map[string]interface{}); ok {
				product := c.transformItemToProduct(data)
				products = append(products, product)
			}
		}
	}

	c.Logger.Info().
		Int("count", len(products)).
		Msg("Successfully fetched products from Dynamics 365")

	return products, nil
}

// FetchStock fetches inventory levels for given SKUs
func (c *Connector) FetchStock(ctx context.Context, skus []string) ([]*erp.StockPayload, error) {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if len(skus) == 0 {
		return []*erp.StockPayload{}, nil
	}

	c.Logger.Info().
		Int("sku_count", len(skus)).
		Msg("Fetching stock levels from Dynamics 365")

	// Build filter for item numbers
	filterParts := make([]string, 0, len(skus))
	for _, sku := range skus {
		filterParts = append(filterParts, fmt.Sprintf("number eq '%s'", sku))
	}

	query := url.Values{}
	query.Set("$filter", strings.Join(filterParts, " or "))
	query.Set("$expand", "itemLedgerEntries")

	// Fetch items with ledger entries
	items, err := c.makeAPIRequest(ctx, "items", query)
	if err != nil {
		return nil, fmt.Errorf("fetching items with stock: %w", err)
	}

	// Transform to stock payloads
	stocks := make([]*erp.StockPayload, 0)
	if value, ok := items["value"].([]interface{}); ok {
		for _, item := range value {
			if data, ok := item.(map[string]interface{}); ok {
				stock := c.transformItemToStock(data)
				stocks = append(stocks, stock)
			}
		}
	}

	c.Logger.Info().
		Int("count", len(stocks)).
		Msg("Successfully fetched stock levels from Dynamics 365")

	return stocks, nil
}

// FetchPrices fetches prices for given SKUs
func (c *Connector) FetchPrices(ctx context.Context, skus []string, priceListID string) ([]*erp.PricePayload, error) {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if len(skus) == 0 {
		return []*erp.PricePayload{}, nil
	}

	c.Logger.Info().
		Int("sku_count", len(skus)).
		Str("pricelist_id", priceListID).
		Msg("Fetching prices from Dynamics 365")

	// In D365 BC, prices can come from multiple sources
	// We'll fetch item prices and sales prices

	// Build filter for items
	filterParts := make([]string, 0, len(skus))
	for _, sku := range skus {
		filterParts = append(filterParts, fmt.Sprintf("itemNo eq '%s'", sku))
	}

	query := url.Values{}
	query.Set("$filter", strings.Join(filterParts, " or "))

	// Fetch sales prices
	salesPrices, err := c.makeAPIRequest(ctx, "salesPrices", query)
	if err != nil {
		return nil, fmt.Errorf("fetching sales prices: %w", err)
	}

	// Transform to price payloads
	prices := make([]*erp.PricePayload, 0)
	if value, ok := salesPrices["value"].([]interface{}); ok {
		for _, price := range value {
			if data, ok := price.(map[string]interface{}); ok {
				pricePayload := c.transformSalesPrice(data, priceListID)
				prices = append(prices, pricePayload)
			}
		}
	}

	// If no sales prices found, fall back to item unit prices
	if len(prices) == 0 {
		itemFilter := make([]string, 0, len(skus))
		for _, sku := range skus {
			itemFilter = append(itemFilter, fmt.Sprintf("number eq '%s'", sku))
		}

		itemQuery := url.Values{}
		itemQuery.Set("$filter", strings.Join(itemFilter, " or "))

		items, err := c.makeAPIRequest(ctx, "items", itemQuery)
		if err != nil {
			return nil, fmt.Errorf("fetching item prices: %w", err)
		}

		if value, ok := items["value"].([]interface{}); ok {
			for _, item := range value {
				if data, ok := item.(map[string]interface{}); ok {
					pricePayload := c.transformItemPrice(data, priceListID)
					prices = append(prices, pricePayload)
				}
			}
		}
	}

	c.Logger.Info().
		Int("count", len(prices)).
		Msg("Successfully fetched prices from Dynamics 365")

	return prices, nil
}

// FetchOrders fetches sales orders changed since the given time
func (c *Connector) FetchOrders(ctx context.Context, since time.Time) ([]*erp.OrderPayload, error) {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}

	c.Logger.Info().
		Time("since", since).
		Msg("Fetching orders from Dynamics 365")

	// Build query
	query := url.Values{}
	query.Set("$filter", fmt.Sprintf("lastModifiedDateTime gt %s", since.Format("2006-01-02T15:04:05Z")))
	query.Set("$expand", "salesOrderLines")
	query.Set("$top", fmt.Sprint(c.Config.Sync.BatchSize))
	query.Set("$orderby", "lastModifiedDateTime asc")

	// Fetch sales orders
	salesOrders, err := c.makeAPIRequest(ctx, "salesOrders", query)
	if err != nil {
		return nil, fmt.Errorf("fetching sales orders: %w", err)
	}

	// Transform to canonical format
	orders := make([]*erp.OrderPayload, 0)
	if value, ok := salesOrders["value"].([]interface{}); ok {
		for _, order := range value {
			if data, ok := order.(map[string]interface{}); ok {
				orderPayload := c.transformSalesOrder(data)
				orders = append(orders, orderPayload)
			}
		}
	}

	c.Logger.Info().
		Int("count", len(orders)).
		Msg("Successfully fetched orders from Dynamics 365")

	return orders, nil
}

// FetchCustomers fetches customers changed since the given time
func (c *Connector) FetchCustomers(ctx context.Context, since time.Time) ([]*erp.CustomerPayload, error) {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}

	c.Logger.Info().
		Time("since", since).
		Msg("Fetching customers from Dynamics 365")

	// Build query
	query := url.Values{}
	query.Set("$filter", fmt.Sprintf("lastModifiedDateTime gt %s", since.Format("2006-01-02T15:04:05Z")))
	query.Set("$top", fmt.Sprint(c.Config.Sync.BatchSize))
	query.Set("$orderby", "lastModifiedDateTime asc")

	// Fetch customers
	customers, err := c.makeAPIRequest(ctx, "customers", query)
	if err != nil {
		return nil, fmt.Errorf("fetching customers: %w", err)
	}

	// Transform to canonical format
	customerPayloads := make([]*erp.CustomerPayload, 0)
	if value, ok := customers["value"].([]interface{}); ok {
		for _, customer := range value {
			if data, ok := customer.(map[string]interface{}); ok {
				customerPayload := c.transformCustomer(data)
				customerPayloads = append(customerPayloads, customerPayload)
			}
		}
	}

	c.Logger.Info().
		Int("count", len(customerPayloads)).
		Msg("Successfully fetched customers from Dynamics 365")

	return customerPayloads, nil
}

// SendOrder creates a sales order in Dynamics 365
func (c *Connector) SendOrder(ctx context.Context, order *erp.OrderPayload) error {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return err
	}

	c.Logger.Info().
		Str("order_id", order.OrderID).
		Msg("Creating order in Dynamics 365")

	// Transform order to D365 format
	d365Order, err := c.transformOrderToD365(ctx, order)
	if err != nil {
		return fmt.Errorf("transforming order: %w", err)
	}

	// Create the sales order
	result, err := c.makeAPIRequest(ctx, "salesOrders", nil, "POST", d365Order)
	if err != nil {
		return fmt.Errorf("creating sales order: %w", err)
	}

	// Extract created order ID
	if orderID, ok := result["id"].(string); ok {
		c.Logger.Info().
			Str("order_id", order.OrderID).
			Str("d365_id", orderID).
			Msg("Successfully created order in Dynamics 365")

		// Create order lines
		if lines, ok := d365Order["lines"].([]interface{}); ok {
			for _, line := range lines {
				lineData := line.(map[string]interface{})
				lineData["documentId"] = orderID

				_, err := c.makeAPIRequest(ctx, fmt.Sprintf("salesOrders(%s)/salesOrderLines", orderID), nil, "POST", lineData)
				if err != nil {
					c.Logger.Error().Err(err).Msg("Failed to create order line")
				}
			}
		}
	}

	return nil
}

// UpdateStock is not directly supported in Dynamics 365
func (c *Connector) UpdateStock(ctx context.Context, stock *erp.StockPayload) error {
	c.Logger.Warn().
		Str("sku", stock.SKU).
		Msg("Direct stock updates not supported in Dynamics 365 - use item journal entries")

	return fmt.Errorf("direct stock updates not supported - use Dynamics 365 item journal entries")
}

// HealthCheck verifies the Dynamics 365 connection
func (c *Connector) HealthCheck(ctx context.Context) error {
	// Try to fetch company information
	companies, err := c.makeAPIRequest(ctx, "companies", nil)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	// Log company info
	if value, ok := companies["value"].([]interface{}); ok && len(value) > 0 {
		if company, ok := value[0].(map[string]interface{}); ok {
			c.Logger.Debug().
				Interface("company", company["name"]).
				Msg("Dynamics 365 health check successful")
		}
	}

	return nil
}

// Helper methods

func (c *Connector) makeAPIRequest(ctx context.Context, resource string, query url.Values, method ...string) (map[string]interface{}, error) {
	// Ensure authentication
	if err := c.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Build URL
	baseURL := fmt.Sprintf("https://api.businesscentral.dynamics.com/v2.0/%s/%s", c.tenantID, c.environment)
	if c.companyID != "" {
		resource = fmt.Sprintf("companies(%s)/%s", c.companyID, resource)
	}

	apiURL, err := url.Parse(fmt.Sprintf("%s/api/v2.0/%s", baseURL, resource))
	if err != nil {
		return nil, err
	}

	if query != nil {
		apiURL.RawQuery = query.Encode()
	}

	// Determine method
	httpMethod := "GET"
	var body []byte
	if len(method) > 0 {
		httpMethod = method[0]
		if len(method) > 1 {
			if data, ok := method[1].(map[string]interface{}); ok {
				body, _ = json.Marshal(data)
			}
		}
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, httpMethod, apiURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Body = http.NoBody
		req.ContentLength = int64(len(body))
		req.GetBody = func() (http.ReadCloser, error) {
			return http.NopCloser(strings.NewReader(string(body))), nil
		}
	}

	// Set headers
	c.tokenMutex.RLock()
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	c.tokenMutex.RUnlock()

	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute request with retry
	resp, err := c.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	// Check for errors
	if errData, ok := result["error"].(map[string]interface{}); ok {
		return nil, fmt.Errorf("API error: %s - %s", errData["code"], errData["message"])
	}

	return result, nil
}

func (c *Connector) fetchResourceData(ctx context.Context, resourceURL string) (interface{}, error) {
	// Parse the resource URL to extract the entity and ID
	parts := strings.Split(resourceURL, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid resource URL: %s", resourceURL)
	}

	// Make request to fetch the resource
	result, err := c.makeAPIRequest(ctx, resourceURL, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Transform functions

func (c *Connector) transformItemToProduct(item map[string]interface{}) *erp.ProductPayload {
	return &erp.ProductPayload{
		SKU:         getStringValue(item, "number"),
		Name:        getStringValue(item, "displayName"),
		Description: getStringValue(item, "description"),
		Category:    getStringValue(item, "itemCategoryCode"),
		Weight:      getFloatValue(item, "grossWeight"),
		Attributes: map[string]interface{}{
			"d365_id":        item["id"],
			"type":           getStringValue(item, "type"),
			"base_unit":      getStringValue(item, "baseUnitOfMeasure"),
			"unit_price":     getFloatValue(item, "unitPrice"),
			"unit_cost":      getFloatValue(item, "unitCost"),
			"tax_group_code": getStringValue(item, "taxGroupCode"),
			"blocked":        getBoolValue(item, "blocked"),
			"gtin":           getStringValue(item, "gtin"),
			"inventory":      getFloatValue(item, "inventory"),
			"last_modified":  getStringValue(item, "lastModifiedDateTime"),
		},
	}
}

func (c *Connector) transformItemToStock(item map[string]interface{}) *erp.StockPayload {
	inventory := int(getFloatValue(item, "inventory"))

	return &erp.StockPayload{
		SKU:        getStringValue(item, "number"),
		LocationID: "MAIN", // Default location
		Quantity:   inventory,
		Available:  inventory, // Simplified - real implementation would calculate
		Reserved:   0,
		StockType:  "physical",
		UpdatedAt:  time.Now(),
	}
}

func (c *Connector) transformSalesPrice(price map[string]interface{}, priceListID string) *erp.PricePayload {
	return &erp.PricePayload{
		SKU:         getStringValue(price, "itemNo"),
		PriceListID: priceListID,
		Currency:    getStringValue(price, "currencyCode"),
		Price:       getFloatValue(price, "unitPrice"),
		MinQuantity: getFloatValue(price, "minimumQuantity"),
		ValidFrom:   parseD365Date(getStringValue(price, "startingDate")),
		ValidTo:     parseD365Date(getStringValue(price, "endingDate")),
	}
}

func (c *Connector) transformItemPrice(item map[string]interface{}, priceListID string) *erp.PricePayload {
	return &erp.PricePayload{
		SKU:         getStringValue(item, "number"),
		PriceListID: priceListID,
		Currency:    "USD", // Default currency
		Price:       getFloatValue(item, "unitPrice"),
		ValidFrom:   time.Now(),
	}
}

func (c *Connector) transformSalesOrder(order map[string]interface{}) *erp.OrderPayload {
	// Extract order lines
	orderItems := make([]erp.OrderItem, 0)
	if lines, ok := order["salesOrderLines"].([]interface{}); ok {
		for _, line := range lines {
			if lineData, ok := line.(map[string]interface{}); ok {
				orderItems = append(orderItems, erp.OrderItem{
					SKU:         getStringValue(lineData, "itemNumber"),
					Name:        getStringValue(lineData, "description"),
					Quantity:    int(getFloatValue(lineData, "quantity")),
					Price:       getFloatValue(lineData, "unitPrice"),
					Discount:    getFloatValue(lineData, "discountPercent"),
					TotalAmount: getFloatValue(lineData, "netAmount"),
				})
			}
		}
	}

	return &erp.OrderPayload{
		OrderID:     getStringValue(order, "number"),
		CustomerID:  getStringValue(order, "customerId"),
		Items:       orderItems,
		TotalAmount: getFloatValue(order, "totalAmountIncludingTax"),
		Currency:    getStringValue(order, "currencyCode"),
		Status:      mapD365OrderStatus(getStringValue(order, "status")),
		CreatedAt:   parseD365Date(getStringValue(order, "orderDate")),
		Attributes: map[string]interface{}{
			"d365_id":              order["id"],
			"external_document_no": getStringValue(order, "externalDocumentNumber"),
			"payment_terms":        getStringValue(order, "paymentTermsId"),
			"shipment_method":      getStringValue(order, "shipmentMethodId"),
			"requested_delivery":   getStringValue(order, "requestedDeliveryDate"),
		},
	}
}

func (c *Connector) transformCustomer(customer map[string]interface{}) *erp.CustomerPayload {
	return &erp.CustomerPayload{
		CustomerID: getStringValue(customer, "number"),
		Email:      getStringValue(customer, "email"),
		Name:       getStringValue(customer, "displayName"),
		Phone:      getStringValue(customer, "phoneNumber"),
		Address: &erp.Address{
			Street:     getStringValue(customer, "addressLine1"),
			Street2:    getStringValue(customer, "addressLine2"),
			City:       getStringValue(customer, "city"),
			State:      getStringValue(customer, "state"),
			PostalCode: getStringValue(customer, "postalCode"),
			Country:    getStringValue(customer, "country"),
		},
		Attributes: map[string]interface{}{
			"d365_id":        customer["id"],
			"type":           getStringValue(customer, "type"),
			"tax_id":         getStringValue(customer, "taxRegistrationNumber"),
			"currency_code":  getStringValue(customer, "currencyCode"),
			"payment_terms":  getStringValue(customer, "paymentTermsId"),
			"payment_method": getStringValue(customer, "paymentMethodId"),
			"credit_limit":   getFloatValue(customer, "creditLimit"),
			"balance":        getFloatValue(customer, "balance"),
			"blocked":        getStringValue(customer, "blocked"),
			"last_modified":  getStringValue(customer, "lastModifiedDateTime"),
		},
	}
}

func (c *Connector) transformOrderToD365(ctx context.Context, order *erp.OrderPayload) (map[string]interface{}, error) {
	// Create sales order structure
	d365Order := map[string]interface{}{
		"customerId":             order.CustomerID,
		"orderDate":              time.Now().Format("2006-01-02"),
		"currencyCode":           order.Currency,
		"externalDocumentNumber": order.OrderID,
	}

	// Prepare order lines
	lines := make([]interface{}, 0, len(order.Items))
	for i, item := range order.Items {
		line := map[string]interface{}{
			"sequence":    (i + 1) * 10000,
			"itemNumber":  item.SKU,
			"description": item.Name,
			"quantity":    item.Quantity,
			"unitPrice":   item.Price,
		}
		lines = append(lines, line)
	}

	d365Order["lines"] = lines

	return d365Order, nil
}

// Utility functions

func getStringValue(data map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if val, ok := data[key].(string); ok && val != "" {
			return val
		}
	}
	return ""
}

func getFloatValue(data map[string]interface{}, key string) float64 {
	if val, ok := data[key].(float64); ok {
		return val
	}
	return 0
}

func getBoolValue(data map[string]interface{}, key string) bool {
	if val, ok := data[key].(bool); ok {
		return val
	}
	return false
}

func parseD365Date(dateStr string) time.Time {
	if dateStr == "" || dateStr == "0001-01-01" {
		return time.Time{}
	}

	// Try different date formats
	formats := []string{
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t
		}
	}

	return time.Now()
}

func mapD365OrderStatus(status string) string {
	// Map D365 status to canonical status
	switch strings.ToLower(status) {
	case "open":
		return "pending"
	case "released":
		return "processing"
	case "partially shipped":
		return "partial"
	case "shipped":
		return "shipped"
	case "invoiced":
		return "completed"
	case "canceled":
		return "cancelled"
	default:
		return status
	}
}

func mapD365EventType(resource string, changeType string) (erp.EventType, error) {
	// Extract entity type from resource
	entity := ""
	if strings.Contains(resource, "/items/") {
		entity = "item"
	} else if strings.Contains(resource, "/customers/") {
		entity = "customer"
	} else if strings.Contains(resource, "/salesOrders/") {
		entity = "salesOrder"
	} else if strings.Contains(resource, "/salesPrices/") {
		entity = "salesPrice"
	}

	// Map to canonical event type
	eventKey := fmt.Sprintf("%s.%s", entity, changeType)
	eventMap := map[string]erp.EventType{
		"item.created":       erp.EventTypeProductCreated,
		"item.updated":       erp.EventTypeProductMasterUpdated,
		"item.deleted":       erp.EventTypeProductDeleted,
		"customer.created":   erp.EventTypeCustomerCreated,
		"customer.updated":   erp.EventTypeCustomerUpdated,
		"salesOrder.created": erp.EventTypeOrderCreated,
		"salesOrder.updated": erp.EventTypeOrderUpdated,
		"salesPrice.created": erp.EventTypePriceUpdated,
		"salesPrice.updated": erp.EventTypePriceUpdated,
	}

	if eventType, ok := eventMap[eventKey]; ok {
		return eventType, nil
	}

	return "", fmt.Errorf("unknown Dynamics 365 event: %s", eventKey)
}
