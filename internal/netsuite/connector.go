package netsuite

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"middleman/internal/erp"
)

// Connector implements the NetSuite ERP connector with OAuth 1.0a TBA
type Connector struct {
	*erp.BaseConnector
	
	// OAuth 1.0a TBA credentials
	accountID       string
	consumerKey     string
	consumerSecret  string
	tokenID         string
	tokenSecret     string
	
	// Configuration
	restletURL      string // Base URL for RESTlet endpoints
	apiVersion      string
	realm           string
	
	// Request tracking
	requestMutex    sync.Mutex
	nonceCounter    int64
}

// NewConnector creates a new NetSuite connector
func NewConnector(config erp.ERPConfig) (*Connector, error) {
	if config.Type != erp.ERPTypeNetSuite {
		return nil, fmt.Errorf("invalid ERP type: %s, expected %s", config.Type, erp.ERPTypeNetSuite)
	}

	// Validate OAuth 1.0a TBA configuration
	if config.Auth.ConsumerKey == "" {
		return nil, fmt.Errorf("NetSuite requires ConsumerKey for OAuth 1.0a")
	}
	if config.Auth.ConsumerSecret == "" {
		return nil, fmt.Errorf("NetSuite requires ConsumerSecret for OAuth 1.0a")
	}
	if config.Auth.TokenID == "" {
		return nil, fmt.Errorf("NetSuite requires TokenID for OAuth 1.0a")
	}
	if config.Auth.TokenSecret == "" {
		return nil, fmt.Errorf("NetSuite requires TokenSecret for OAuth 1.0a")
	}

	accountID, ok := config.Metadata["account_id"].(string)
	if !ok || accountID == "" {
		return nil, fmt.Errorf("NetSuite requires account_id in metadata")
	}

	consumerKey := config.Auth.ConsumerKey
	consumerSecret := config.Auth.ConsumerSecret
	tokenID := config.Auth.TokenID
	tokenSecret := config.Auth.TokenSecret

	// Optional configurations
	apiVersion := "2024.1"
	if v, ok := config.Metadata["api_version"].(string); ok {
		apiVersion = v
	}

	// Build RESTlet URL
	datacenter := strings.ToLower(accountID)
	if dc, ok := config.Metadata["datacenter"].(string); ok {
		datacenter = dc
	}
	
	restletURL := fmt.Sprintf("https://%s.restlets.api.netsuite.com", datacenter)
	if url, ok := config.Metadata["restlet_url"].(string); ok {
		restletURL = url
	}

	// Extract realm (account ID with underscores)
	realm := strings.ReplaceAll(strings.ToUpper(accountID), "-", "_")

	return &Connector{
		BaseConnector:  erp.NewBaseConnector(config),
		accountID:      accountID,
		consumerKey:    consumerKey,
		consumerSecret: consumerSecret,
		tokenID:        tokenID,
		tokenSecret:    tokenSecret,
		restletURL:     restletURL,
		apiVersion:     apiVersion,
		realm:          realm,
	}, nil
}

// GetID returns the connector ID
func (c *Connector) GetID() string {
	return c.Config.ID
}

// GetType returns the ERP type as string
func (c *Connector) GetType() string {
	return string(erp.ERPTypeNetSuite)
}

// GetConfig returns the connector configuration
func (c *Connector) GetConfig() erp.ERPConfig {
	return c.Config
}

// TestConnection validates the connector can reach the ERP
func (c *Connector) TestConnection(ctx context.Context) error {
	// Execute a simple query to verify connection
	query := "SELECT 1 as test FROM DUAL"
	_, err := c.executeSuiteQL(ctx, query)
	return err
}

// generateOAuthHeader generates OAuth 1.0a header for NetSuite
func (c *Connector) generateOAuthHeader(method, urlStr string) string {
	c.requestMutex.Lock()
	c.nonceCounter++
	nonce := fmt.Sprintf("%s_%d", uuid.New().String(), c.nonceCounter)
	c.requestMutex.Unlock()

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// OAuth parameters
	oauthParams := map[string]string{
		"oauth_consumer_key":     c.consumerKey,
		"oauth_token":            c.tokenID,
		"oauth_signature_method": "HMAC-SHA256",
		"oauth_timestamp":        timestamp,
		"oauth_nonce":            nonce,
		"oauth_version":          "1.0",
	}

	// Create signature base string
	signatureBase := c.createSignatureBase(method, urlStr, oauthParams)

	// Generate signature
	signingKey := fmt.Sprintf("%s&%s", 
		url.QueryEscape(c.consumerSecret), 
		url.QueryEscape(c.tokenSecret))
	
	h := hmac.New(sha256.New, []byte(signingKey))
	h.Write([]byte(signatureBase))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Build authorization header
	authHeader := fmt.Sprintf(`OAuth realm="%s",`, c.realm)
	authHeader += fmt.Sprintf(`oauth_consumer_key="%s",`, oauthParams["oauth_consumer_key"])
	authHeader += fmt.Sprintf(`oauth_token="%s",`, oauthParams["oauth_token"])
	authHeader += fmt.Sprintf(`oauth_signature_method="%s",`, oauthParams["oauth_signature_method"])
	authHeader += fmt.Sprintf(`oauth_timestamp="%s",`, oauthParams["oauth_timestamp"])
	authHeader += fmt.Sprintf(`oauth_nonce="%s",`, oauthParams["oauth_nonce"])
	authHeader += fmt.Sprintf(`oauth_version="%s",`, oauthParams["oauth_version"])
	authHeader += fmt.Sprintf(`oauth_signature="%s"`, url.QueryEscape(signature))

	return authHeader
}

// createSignatureBase creates the OAuth 1.0a signature base string
func (c *Connector) createSignatureBase(method, urlStr string, params map[string]string) string {
	// Parse URL to get base URL without query params
	u, _ := url.Parse(urlStr)
	baseURL := fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path)

	// Collect and sort all parameters
	allParams := make([]string, 0)
	for k, v := range params {
		allParams = append(allParams, fmt.Sprintf("%s=%s", 
			url.QueryEscape(k), url.QueryEscape(v)))
	}
	
	// Add query parameters if any
	for k, v := range u.Query() {
		allParams = append(allParams, fmt.Sprintf("%s=%s", 
			url.QueryEscape(k), url.QueryEscape(v[0])))
	}

	sort.Strings(allParams)
	paramString := strings.Join(allParams, "&")

	// Create signature base
	signatureBase := fmt.Sprintf("%s&%s&%s",
		strings.ToUpper(method),
		url.QueryEscape(baseURL),
		url.QueryEscape(paramString))

	return signatureBase
}

// ValidateWebhook validates a webhook signature
func (c *Connector) ValidateWebhook(payload []byte, signature string) error {
	// NetSuite webhooks use HMAC-SHA256 with a shared secret
	if !c.Config.Webhook.ValidateSign || c.Config.Webhook.Secret == "" {
		return nil
	}

	// Calculate expected signature
	h := hmac.New(sha256.New, []byte(c.Config.Webhook.Secret))
	h.Write(payload)
	expectedSig := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Compare signatures
	if signature != expectedSig {
		return fmt.Errorf("invalid webhook signature")
	}

	return nil
}

// ParseWebhook parses webhook payload into canonical event
func (c *Connector) ParseWebhook(payload []byte) (*erp.CanonicalEvent, error) {
	var webhook NetSuiteWebhook
	if err := json.Unmarshal(payload, &webhook); err != nil {
		return nil, fmt.Errorf("parsing webhook payload: %w", err)
	}

	// Map NetSuite events to canonical events
	eventType, err := mapNetSuiteEventType(webhook.RecordType, webhook.EventType)
	if err != nil {
		return nil, err
	}

	// Transform payload based on event type
	canonicalPayload, err := c.transformWebhookPayload(&webhook, eventType)
	if err != nil {
		return nil, fmt.Errorf("transforming payload: %w", err)
	}

	return &erp.CanonicalEvent{
		EventID:        webhook.EventID,
		EventType:      eventType,
		EventTimestamp: webhook.Timestamp,
		Source:         fmt.Sprintf("NetSuite-%s", c.accountID),
		CorrelationID:  fmt.Sprintf("%s-%s", webhook.RecordType, webhook.RecordID),
		Payload:        canonicalPayload,
	}, nil
}

// SyncProducts synchronizes products from the ERP
func (c *Connector) SyncProducts(ctx context.Context, since time.Time, batchSize int) ([]*erp.ProductPayload, error) {
	c.Logger.Info().
		Time("since", since).
		Msg("Fetching products from NetSuite")

	// Build SuiteQL query
	query := fmt.Sprintf(`
		SELECT 
			i.id,
			i.itemid,
			i.displayname,
			i.salesdescription,
			i.purchasedescription,
			i.class,
			i.weight,
			i.weightunit,
			i.unitstype,
			i.baseprice,
			i.cost,
			i.costingmethod,
			i.isinactive,
			i.upccode,
			i.vendorname,
			i.lastmodifieddate
		FROM item i
		WHERE i.lastmodifieddate >= TO_DATE('%s', 'YYYY-MM-DD HH24:MI:SS')
		AND i.isinactive = 'F'
		ORDER BY i.lastmodifieddate ASC
		FETCH FIRST %d ROWS ONLY
	`, since.Format("2006-01-02 15:04:05"), c.Config.Sync.BatchSize)

	// Execute query via RESTlet
	result, err := c.executeSuiteQL(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("executing SuiteQL query: %w", err)
	}

	// Transform results
	products := make([]*erp.ProductPayload, 0)
	if items, ok := result["items"].([]interface{}); ok {
		for _, item := range items {
			if data, ok := item.(map[string]interface{}); ok {
				product := c.transformItemToProduct(data)
				products = append(products, product)
			}
		}
	}

	c.Logger.Info().
		Int("count", len(products)).
		Msg("Successfully fetched products from NetSuite")

	return products, nil
}

// SyncStock synchronizes stock levels from the ERP
func (c *Connector) SyncStock(ctx context.Context, productIDs []string, batchSize int) ([]*erp.StockPayload, error) {
	if len(productIDs) == 0 {
		return []*erp.StockPayload{}, nil
	}

	c.Logger.Info().
		Int("product_count", len(productIDs)).
		Msg("Fetching stock levels from NetSuite")

	// Build item list for query
	itemList := make([]string, len(productIDs))
	for i, id := range productIDs {
		itemList[i] = fmt.Sprintf("'%s'", id)
	}

	// Query inventory balance
	query := fmt.Sprintf(`
		SELECT 
			i.itemid,
			l.name as location,
			ib.quantityonhand,
			ib.quantityavailable,
			ib.quantityonorder,
			ib.quantityintransit
		FROM inventorybalance ib
		JOIN item i ON ib.item = i.id
		JOIN location l ON ib.location = l.id
		WHERE i.id IN (%s)
		AND l.isinactive = 'F'
	`, strings.Join(itemList, ","))

	// Execute query
	result, err := c.executeSuiteQL(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("executing inventory query: %w", err)
	}

	// Transform results
	stocks := make([]*erp.StockPayload, 0)
	if items, ok := result["items"].([]interface{}); ok {
		for _, item := range items {
			if data, ok := item.(map[string]interface{}); ok {
				stock := c.transformInventoryToStock(data)
				stocks = append(stocks, stock)
			}
		}
	}

	c.Logger.Info().
		Int("count", len(stocks)).
		Msg("Successfully fetched stock levels from NetSuite")

	return stocks, nil
}

// SyncPrices synchronizes prices from the ERP
func (c *Connector) SyncPrices(ctx context.Context, productIDs []string, batchSize int) ([]*erp.PricePayload, error) {
	if len(productIDs) == 0 {
		return []*erp.PricePayload{}, nil
	}

	c.Logger.Info().
		Int("product_count", len(productIDs)).
		Msg("Fetching prices from NetSuite")

	// Build item list
	itemList := make([]string, len(productIDs))
	for i, id := range productIDs {
		itemList[i] = fmt.Sprintf("'%s'", id)
	}

	// Determine price level from config
	priceLevel := "1" // Base price
	if pl, ok := c.Config.Metadata["default_price_level"].(string); ok {
		priceLevel = pl
	}

	// Query pricing
	query := fmt.Sprintf(`
		SELECT 
			i.itemid,
			i.baseprice,
			ip.pricelevel,
			ip.price,
			c.symbol as currency,
			pl.name as pricelevelname
		FROM item i
		LEFT JOIN itemprice ip ON i.id = ip.item
		LEFT JOIN pricelevel pl ON ip.pricelevel = pl.id
		LEFT JOIN currency c ON ip.currency = c.id
		WHERE i.id IN (%s)
		AND (ip.pricelevel = %s OR ip.pricelevel IS NULL)
	`, strings.Join(itemList, ","), priceLevel)

	// Execute query
	result, err := c.executeSuiteQL(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("executing price query: %w", err)
	}

	// Transform results
	prices := make([]*erp.PricePayload, 0)
	if items, ok := result["items"].([]interface{}); ok {
		for _, item := range items {
			if data, ok := item.(map[string]interface{}); ok {
				price := c.transformPriceData(data, priceLevel)
				prices = append(prices, price)
			}
		}
	}

	c.Logger.Info().
		Int("count", len(prices)).
		Msg("Successfully fetched prices from NetSuite")

	return prices, nil
}

// FetchOrders fetches sales orders changed since the given time
func (c *Connector) FetchOrders(ctx context.Context, since time.Time) ([]*erp.OrderPayload, error) {
	c.Logger.Info().
		Time("since", since).
		Msg("Fetching orders from NetSuite")

	// Query sales orders with line items
	query := fmt.Sprintf(`
		SELECT 
			so.id,
			so.tranid,
			so.entity,
			so.trandate,
			so.status,
			so.total,
			so.currency,
			so.lastmodifieddate,
			c.companyname,
			c.email
		FROM transaction so
		JOIN customer c ON so.entity = c.id
		WHERE so.type = 'SalesOrd'
		AND so.lastmodifieddate >= TO_DATE('%s', 'YYYY-MM-DD HH24:MI:SS')
		ORDER BY so.lastmodifieddate ASC
		FETCH FIRST %d ROWS ONLY
	`, since.Format("2006-01-02 15:04:05"), c.Config.Sync.BatchSize)

	// Execute query
	result, err := c.executeSuiteQL(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("executing order query: %w", err)
	}

	// Transform results and fetch line items
	orders := make([]*erp.OrderPayload, 0)
	if items, ok := result["items"].([]interface{}); ok {
		for _, item := range items {
			if data, ok := item.(map[string]interface{}); ok {
				// Fetch line items for each order
				orderID := fmt.Sprintf("%v", data["id"])
				lineItems, err := c.fetchOrderLines(ctx, orderID)
				if err != nil {
					c.Logger.Error().Err(err).Str("order_id", orderID).Msg("Failed to fetch order lines")
					continue
				}
				
				order := c.transformSalesOrder(data, lineItems)
				orders = append(orders, order)
			}
		}
	}

	c.Logger.Info().
		Int("count", len(orders)).
		Msg("Successfully fetched orders from NetSuite")

	return orders, nil
}

// SyncCustomers synchronizes customers from the ERP
func (c *Connector) SyncCustomers(ctx context.Context, since time.Time, batchSize int) ([]*erp.CustomerPayload, error) {
	c.Logger.Info().
		Time("since", since).
		Msg("Fetching customers from NetSuite")

	// Query customers
	query := fmt.Sprintf(`
		SELECT 
			c.id,
			c.entityid,
			c.companyname,
			c.firstname,
			c.lastname,
			c.email,
			c.phone,
			c.fax,
			c.url,
			c.defaultaddress,
			c.terms,
			c.creditlimit,
			c.balance,
			c.currency,
			c.isinactive,
			c.lastmodifieddate,
			a.addr1,
			a.addr2,
			a.city,
			a.state,
			a.zip,
			a.country
		FROM customer c
		LEFT JOIN customeraddressbook cab ON c.id = cab.entity AND cab.defaultbilling = 'T'
		LEFT JOIN customeraddress a ON cab.addressbookaddress = a.nkey
		WHERE c.lastmodifieddate >= TO_DATE('%s', 'YYYY-MM-DD HH24:MI:SS')
		AND c.isinactive = 'F'
		ORDER BY c.lastmodifieddate ASC
		FETCH FIRST %d ROWS ONLY
	`, since.Format("2006-01-02 15:04:05"), c.Config.Sync.BatchSize)

	// Execute query
	result, err := c.executeSuiteQL(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("executing customer query: %w", err)
	}

	// Transform results
	customers := make([]*erp.CustomerPayload, 0)
	if items, ok := result["items"].([]interface{}); ok {
		for _, item := range items {
			if data, ok := item.(map[string]interface{}); ok {
				customer := c.transformCustomer(data)
				customers = append(customers, customer)
			}
		}
	}

	c.Logger.Info().
		Int("count", len(customers)).
		Msg("Successfully fetched customers from NetSuite")

	return customers, nil
}

// SendOrder creates a sales order in NetSuite
func (c *Connector) SendOrder(ctx context.Context, order *erp.OrderPayload) error {
	c.Logger.Info().
		Str("order_id", order.OrderID).
		Msg("Creating order in NetSuite")

	// Transform order to NetSuite format
	nsOrder, err := c.transformOrderToNetSuite(ctx, order)
	if err != nil {
		return fmt.Errorf("transforming order: %w", err)
	}

	// Create order via RESTlet
	restletPath := "/restlets/create_order"
	if path, ok := c.Config.Metadata["create_order_restlet"].(string); ok {
		restletPath = path
	}

	result, err := c.callRESTlet(ctx, restletPath, nsOrder)
	if err != nil {
		return fmt.Errorf("creating order: %w", err)
	}

	// Check result
	if orderID, ok := result["id"].(string); ok {
		c.Logger.Info().
			Str("order_id", order.OrderID).
			Str("netsuite_id", orderID).
			Msg("Successfully created order in NetSuite")
	}

	return nil
}

// UpdateStock updates inventory in NetSuite
func (c *Connector) UpdateStock(ctx context.Context, stock *erp.StockPayload) error {
	c.Logger.Warn().
		Str("sku", stock.SKU).
		Msg("Direct stock updates not recommended in NetSuite - use inventory adjustment")

	return fmt.Errorf("direct stock updates not supported - use NetSuite inventory adjustment transactions")
}

// ProcessWebhook handles incoming webhook from NetSuite
func (c *Connector) ProcessWebhook(ctx context.Context, payload []byte, signature string) error {
	c.Logger.Info().
		Int("payload_size", len(payload)).
		Msg("Processing NetSuite webhook")

	// Validate webhook signature
	if err := c.ValidateWebhook(payload, signature); err != nil {
		return fmt.Errorf("webhook validation failed: %w", err)
	}

	// Parse webhook into canonical event
	canonicalEvent, err := c.ParseWebhook(payload)
	if err != nil {
		return fmt.Errorf("parsing webhook: %w", err)
	}

	c.Logger.Info().
		Str("event_id", canonicalEvent.EventID).
		Str("event_type", string(canonicalEvent.EventType)).
		Msg("Successfully processed NetSuite webhook")

	// The actual processing of the event (e.g., updating local data)
	// would be handled by the application layer based on the canonical event
	
	return nil
}

// HealthCheck verifies the connector is operational
func (c *Connector) HealthCheck(ctx context.Context) erp.HealthCheck {
	// Execute a simple query to verify connection
	query := "SELECT 1 as test FROM DUAL"
	
	_, err := c.executeSuiteQL(ctx, query)
	if err != nil {
		return erp.HealthCheck{
			Status:  erp.HealthStatusUnhealthy,
			Message: fmt.Sprintf("Connection failed: %v", err),
		}
	}

	c.Logger.Debug().Msg("NetSuite health check successful")
	return erp.HealthCheck{
		Status:  erp.HealthStatusHealthy,
		Message: "Connector is operational",
	}
}

// GetOrder retrieves an order from NetSuite
func (c *Connector) GetOrder(ctx context.Context, orderID string) (*erp.OrderPayload, error) {
	c.Logger.Info().
		Str("order_id", orderID).
		Msg("Fetching order from NetSuite")

	// Build sales order search request
	request := map[string]interface{}{
		"recordtype": "salesorder",
		"filters": []map[string]interface{}{
			{
				"field":    "tranid",
				"operator": "is",
				"value":    orderID,
			},
		},
		"columns": []string{
			"tranid", "entity", "total", "subtotal", "taxtotal",
			"currency", "status", "trandate", "memo", "item",
		},
	}

	// Call search RESTlet
	restletPath := "/restlets/search_records"
	if path, ok := c.Config.Metadata["search_records_restlet"].(string); ok {
		restletPath = path
	}

	result, err := c.callRESTlet(ctx, restletPath, request)
	if err != nil {
		return nil, fmt.Errorf("searching order: %w", err)
	}

	// Parse results
	records, ok := result["records"].([]interface{})
	if !ok || len(records) == 0 {
		return nil, fmt.Errorf("order not found: %s", orderID)
	}

	// Transform the first matching order
	orderRecord := records[0].(map[string]interface{})
	
	// Get line items
	var lineItems []erp.OrderItem
	if items, ok := orderRecord["item"].([]interface{}); ok {
		for _, item := range items {
			itemData := item.(map[string]interface{})
			lineItem := erp.OrderItem{
				Name:        getStringValue(itemData, "item_display"),
				SKU:         getStringValue(itemData, "item"),
				Quantity:    int(getFloatValue(itemData, "quantity")),
				Price:       getFloatValue(itemData, "rate"),
				TotalAmount: getFloatValue(itemData, "amount"),
				TaxRate:     getFloatValue(itemData, "taxrate1"),
			}
			lineItems = append(lineItems, lineItem)
		}
	}

	order := c.transformSalesOrder(orderRecord, lineItems)
	
	c.Logger.Info().
		Str("order_id", orderID).
		Msg("Successfully fetched order from NetSuite")

	return order, nil
}

// Helper methods

func (c *Connector) executeSuiteQL(ctx context.Context, query string) (map[string]interface{}, error) {
	// Call SuiteQL RESTlet
	restletPath := "/restlets/suiteql"
	if path, ok := c.Config.Metadata["suiteql_restlet"].(string); ok {
		restletPath = path
	}

	payload := map[string]interface{}{
		"query": query,
	}

	return c.callRESTlet(ctx, restletPath, payload)
}

func (c *Connector) callRESTlet(ctx context.Context, path string, payload interface{}) (map[string]interface{}, error) {
	// Build URL
	urlStr := fmt.Sprintf("%s%s", c.restletURL, path)

	// Marshal payload
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshaling payload: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", urlStr, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", c.generateOAuthHeader("POST", urlStr))

	// Execute request
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
	if errorMsg, ok := result["error"].(map[string]interface{}); ok {
		return nil, fmt.Errorf("NetSuite error: %s - %s", 
			errorMsg["code"], errorMsg["message"])
	}

	return result, nil
}

func (c *Connector) fetchOrderLines(ctx context.Context, orderID string) ([]erp.OrderItem, error) {
	query := fmt.Sprintf(`
		SELECT 
			tl.item,
			i.itemid,
			i.displayname,
			tl.quantity,
			tl.rate,
			tl.amount,
			tl.taxrate1
		FROM transactionline tl
		JOIN item i ON tl.item = i.id
		WHERE tl.transaction = %s
		AND tl.mainline = 'F'
		ORDER BY tl.linesequencenumber
	`, orderID)

	result, err := c.executeSuiteQL(ctx, query)
	if err != nil {
		return nil, err
	}

	items := make([]erp.OrderItem, 0)
	if lines, ok := result["items"].([]interface{}); ok {
		for _, line := range lines {
			if data, ok := line.(map[string]interface{}); ok {
				items = append(items, erp.OrderItem{
					SKU:         getStringValue(data, "itemid"),
					Name:        getStringValue(data, "displayname"),
					Quantity:    int(getFloatValue(data, "quantity")),
					Price:       getFloatValue(data, "rate"),
					TotalAmount: getFloatValue(data, "amount"),
					TaxRate:     getFloatValue(data, "taxrate1"),
				})
			}
		}
	}

	return items, nil
}

// Transform functions

func (c *Connector) transformWebhookPayload(webhook *NetSuiteWebhook, eventType erp.EventType) (interface{}, error) {
	// Transform based on event type
	switch eventType {
	case erp.EventTypeProductMasterUpdated, erp.EventTypeProductCreated:
		return c.transformProductWebhook(webhook.Data), nil
	case erp.EventTypeStockLevelUpdated:
		return c.transformStockWebhook(webhook.Data), nil
	case erp.EventTypePriceUpdated:
		return c.transformPriceWebhook(webhook.Data), nil
	case erp.EventTypeOrderCreated, erp.EventTypeOrderUpdated:
		return c.transformOrderWebhook(webhook.Data), nil
	case erp.EventTypeCustomerCreated, erp.EventTypeCustomerUpdated:
		return c.transformCustomerWebhook(webhook.Data), nil
	default:
		return webhook.Data, nil
	}
}

func (c *Connector) transformProductWebhook(data map[string]interface{}) *erp.ProductPayload {
	return &erp.ProductPayload{
		SKU:         getStringValue(data, "itemid"),
		Name:        getStringValue(data, "displayname"),
		Description: getStringValue(data, "salesdescription"),
		Category:    getStringValue(data, "class"),
		Weight:      getFloatValue(data, "weight"),
		Attributes:  data,
	}
}

func (c *Connector) transformStockWebhook(data map[string]interface{}) *erp.StockPayload {
	return &erp.StockPayload{
		SKU:        getStringValue(data, "itemid"),
		LocationID: getStringValue(data, "location"),
		Quantity:   int(getFloatValue(data, "quantityonhand")),
		Available:  int(getFloatValue(data, "quantityavailable")),
		StockType:  "physical",
		UpdatedAt:  time.Now(),
	}
}

func (c *Connector) transformPriceWebhook(data map[string]interface{}) *erp.PricePayload {
	return &erp.PricePayload{
		SKU:         getStringValue(data, "itemid"),
		PriceListID: getStringValue(data, "pricelevel"),
		Currency:    getStringValue(data, "currency"),
		Price:       getFloatValue(data, "price"),
		ValidFrom:   time.Now(),
	}
}

func (c *Connector) transformOrderWebhook(data map[string]interface{}) *erp.OrderPayload {
	return &erp.OrderPayload{
		OrderID:     getStringValue(data, "tranid"),
		CustomerID:  getStringValue(data, "entity"),
		TotalAmount: getFloatValue(data, "total"),
		Currency:    getStringValue(data, "currency"),
		Status:      mapNetSuiteOrderStatus(getStringValue(data, "status")),
		CreatedAt:   parseNetSuiteDate(getStringValue(data, "trandate")),
	}
}

func (c *Connector) transformCustomerWebhook(data map[string]interface{}) *erp.CustomerPayload {
	return &erp.CustomerPayload{
		CustomerID: getStringValue(data, "entityid"),
		Email:      getStringValue(data, "email"),
		Name:       getStringValue(data, "companyname", "entityid"),
		Phone:      getStringValue(data, "phone"),
		Attributes: data,
	}
}

func (c *Connector) transformItemToProduct(item map[string]interface{}) *erp.ProductPayload {
	return &erp.ProductPayload{
		SKU:         getStringValue(item, "itemid"),
		Name:        getStringValue(item, "displayname"),
		Description: getStringValue(item, "salesdescription", "purchasedescription"),
		Category:    getStringValue(item, "class"),
		Weight:      getFloatValue(item, "weight"),
		Attributes: map[string]interface{}{
			"netsuite_id":    item["id"],
			"weight_unit":    getStringValue(item, "weightunit"),
			"base_price":     getFloatValue(item, "baseprice"),
			"cost":           getFloatValue(item, "cost"),
			"costing_method": getStringValue(item, "costingmethod"),
			"upc_code":       getStringValue(item, "upccode"),
			"vendor":         getStringValue(item, "vendorname"),
			"inactive":       getBoolValue(item, "isinactive"),
			"last_modified":  getStringValue(item, "lastmodifieddate"),
		},
	}
}

func (c *Connector) transformInventoryToStock(inv map[string]interface{}) *erp.StockPayload {
	onHand := int(getFloatValue(inv, "quantityonhand"))
	available := int(getFloatValue(inv, "quantityavailable"))
	onOrder := int(getFloatValue(inv, "quantityonorder"))
	inTransit := int(getFloatValue(inv, "quantityintransit"))

	return &erp.StockPayload{
		SKU:        getStringValue(inv, "itemid"),
		LocationID: getStringValue(inv, "location"),
		Quantity:   onHand,
		Available:  available,
		Reserved:   onHand - available,
		StockType:  "physical",
		UpdatedAt:  time.Now(),
		Attributes: map[string]interface{}{
			"on_order":   onOrder,
			"in_transit": inTransit,
		},
	}
}

func (c *Connector) transformPriceData(price map[string]interface{}, priceListID string) *erp.PricePayload {
	// Use specific price if available, otherwise base price
	priceValue := getFloatValue(price, "price")
	if priceValue == 0 {
		priceValue = getFloatValue(price, "baseprice")
	}

	return &erp.PricePayload{
		SKU:         getStringValue(price, "itemid"),
		PriceListID: priceListID,
		Currency:    getStringValue(price, "currency", "USD"),
		Price:       priceValue,
		ValidFrom:   time.Now(),
		Attributes: map[string]interface{}{
			"price_level_name": getStringValue(price, "pricelevelname"),
		},
	}
}

func (c *Connector) transformSalesOrder(order map[string]interface{}, lineItems []erp.OrderItem) *erp.OrderPayload {
	return &erp.OrderPayload{
		OrderID:     getStringValue(order, "tranid"),
		CustomerID:  getStringValue(order, "entity"),
		Items:       lineItems,
		TotalAmount: getFloatValue(order, "total"),
		Currency:    getStringValue(order, "currency"),
		Status:      mapNetSuiteOrderStatus(getStringValue(order, "status")),
		CreatedAt:   parseNetSuiteDate(getStringValue(order, "trandate")),
		Attributes: map[string]interface{}{
			"netsuite_id":    order["id"],
			"customer_name":  getStringValue(order, "companyname"),
			"customer_email": getStringValue(order, "email"),
			"last_modified":  getStringValue(order, "lastmodifieddate"),
		},
	}
}

func (c *Connector) transformCustomer(customer map[string]interface{}) *erp.CustomerPayload {
	// Build full name
	name := getStringValue(customer, "companyname")
	if name == "" {
		firstName := getStringValue(customer, "firstname")
		lastName := getStringValue(customer, "lastname")
		if firstName != "" || lastName != "" {
			name = strings.TrimSpace(fmt.Sprintf("%s %s", firstName, lastName))
		}
	}

	return &erp.CustomerPayload{
		CustomerID: getStringValue(customer, "entityid"),
		Email:      getStringValue(customer, "email"),
		Name:       name,
		Phone:      getStringValue(customer, "phone"),
		Address: &erp.Address{
			Street:     getStringValue(customer, "addr1"),
			Street2:    getStringValue(customer, "addr2"),
			City:       getStringValue(customer, "city"),
			State:      getStringValue(customer, "state"),
			PostalCode: getStringValue(customer, "zip"),
			Country:    getStringValue(customer, "country"),
		},
		Attributes: map[string]interface{}{
			"netsuite_id":    customer["id"],
			"fax":            getStringValue(customer, "fax"),
			"website":        getStringValue(customer, "url"),
			"terms":          getStringValue(customer, "terms"),
			"credit_limit":   getFloatValue(customer, "creditlimit"),
			"balance":        getFloatValue(customer, "balance"),
			"currency":       getStringValue(customer, "currency"),
			"inactive":       getBoolValue(customer, "isinactive"),
			"last_modified":  getStringValue(customer, "lastmodifieddate"),
		},
	}
}

func (c *Connector) transformOrderToNetSuite(ctx context.Context, order *erp.OrderPayload) (map[string]interface{}, error) {
	// Create NetSuite order structure
	nsOrder := map[string]interface{}{
		"recordtype": "salesorder",
		"entity":     order.CustomerID,
		"trandate":   time.Now().Format("01/02/2006"),
		"otherrefnum": order.OrderID,
		"currency":   order.Currency,
		"item": make([]map[string]interface{}, 0, len(order.Items)),
	}

	// Add line items
	for _, item := range order.Items {
		nsOrder["item"] = append(nsOrder["item"].([]map[string]interface{}), map[string]interface{}{
			"item":     item.SKU,
			"quantity": item.Quantity,
			"rate":     item.Price,
		})
	}

	return nsOrder, nil
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
	if val, ok := data[key].(string); ok {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return 0
}

func getBoolValue(data map[string]interface{}, key string) bool {
	if val, ok := data[key].(bool); ok {
		return val
	}
	if val, ok := data[key].(string); ok {
		return val == "T" || val == "true" || val == "1"
	}
	return false
}

func parseNetSuiteDate(dateStr string) time.Time {
	// NetSuite dates can be in various formats
	formats := []string{
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"01/02/2006",
		"2006-01-02",
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t
		}
	}
	
	return time.Now()
}

func mapNetSuiteOrderStatus(status string) string {
	// Map NetSuite status to canonical status
	statusMap := map[string]string{
		"pendingApproval":     "pending",
		"pendingFulfillment":  "processing",
		"partiallyFulfilled":  "partial",
		"pendingBilling":      "shipped",
		"fullyBilled":         "completed",
		"closed":              "completed",
		"cancelled":           "cancelled",
	}

	if mapped, ok := statusMap[strings.ToLower(status)]; ok {
		return mapped
	}
	return status
}

func mapNetSuiteEventType(recordType, eventType string) (erp.EventType, error) {
	// Map NetSuite record types and events to canonical events
	eventKey := fmt.Sprintf("%s.%s", recordType, eventType)
	eventMap := map[string]erp.EventType{
		"item.create":              erp.EventTypeProductCreated,
		"item.update":              erp.EventTypeProductMasterUpdated,
		"item.delete":              erp.EventTypeProductDeleted,
		"inventoryitem.create":     erp.EventTypeProductCreated,
		"inventoryitem.update":     erp.EventTypeProductMasterUpdated,
		"inventoryadjustment.create": erp.EventTypeStockLevelUpdated,
		"inventoryadjustment.update": erp.EventTypeStockLevelUpdated,
		"itemfulfillment.create":   erp.EventTypeStockLevelUpdated,
		"itemreceipt.create":       erp.EventTypeStockLevelUpdated,
		"pricebook.update":         erp.EventTypePriceUpdated,
		"salesorder.create":        erp.EventTypeOrderCreated,
		"salesorder.update":        erp.EventTypeOrderUpdated,
		"customer.create":          erp.EventTypeCustomerCreated,
		"customer.update":          erp.EventTypeCustomerUpdated,
	}

	if eventType, ok := eventMap[eventKey]; ok {
		return eventType, nil
	}

	return "", fmt.Errorf("unknown NetSuite event: %s", eventKey)
}// FetchCustomer fetches a specific customer by ID
func (c *Connector) FetchCustomer(ctx context.Context, customerID string) (*erp.CustomerPayload, error) {
	c.Logger.Info().
		Str("customer_id", customerID).
		Msg("Fetching customer from NetSuite")

	// Query specific customer
	query := fmt.Sprintf(`
		SELECT 
			c.id,
			c.entityid,
			c.companyname,
			c.firstname,
			c.lastname,
			c.email,
			c.phone,
			c.fax,
			c.url,
			c.defaultaddress,
			c.terms,
			c.creditlimit,
			c.balance,
			c.currency,
			c.isinactive,
			c.lastmodifieddate,
			a.addr1,
			a.addr2,
			a.city,
			a.state,
			a.zip,
			a.country
		FROM customer c
		LEFT JOIN customeraddressbook cab ON c.id = cab.entity AND cab.defaultbilling = 'T'
		LEFT JOIN customeraddress a ON cab.addressbookaddress = a.nkey
		WHERE c.entityid = '%s'
		AND c.isinactive = 'F'
	`, customerID)

	// Execute query
	result, err := c.executeSuiteQL(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("executing customer query: %w", err)
	}

	// Check if customer found
	if items, ok := result["items"].([]interface{}); ok && len(items) > 0 {
		if data, ok := items[0].(map[string]interface{}); ok {
			customer := c.transformCustomer(data)
			return customer, nil
		}
	}

	return nil, fmt.Errorf("customer not found: %s", customerID)
}

// FetchCustomers fetches multiple customers from NetSuite
func (c *Connector) FetchCustomers(ctx context.Context, customerIDs []string) ([]*erp.CustomerPayload, error) {
	c.Logger.Info().
		Int("count", len(customerIDs)).
		Msg("Fetching multiple customers from NetSuite")

	if len(customerIDs) == 0 {
		return []*erp.CustomerPayload{}, nil
	}

	customers := make([]*erp.CustomerPayload, 0, len(customerIDs))
	
	// Fetch each customer individually
	// Note: NetSuite doesn't have efficient bulk customer fetch, so we fetch one by one
	for _, customerID := range customerIDs {
		customer, err := c.FetchCustomer(ctx, customerID)
		if err != nil {
			c.Logger.Warn().
				Str("customer_id", customerID).
				Err(err).
				Msg("Failed to fetch customer, skipping")
			continue
		}
		customers = append(customers, customer)
	}

	c.Logger.Info().
		Int("fetched", len(customers)).
		Int("requested", len(customerIDs)).
		Msg("Completed fetching customers from NetSuite")

	return customers, nil
}

// FetchPrices fetches prices for specific products and price list
func (c *Connector) FetchPrices(ctx context.Context, productIDs []string, priceListID string) ([]*erp.PricePayload, error) {
	c.Logger.Info().
		Int("product_count", len(productIDs)).
		Str("price_list_id", priceListID).
		Msg("Fetching prices from NetSuite")

	if len(productIDs) == 0 {
		return []*erp.PricePayload{}, nil
	}

	prices := make([]*erp.PricePayload, 0, len(productIDs))
	
	// Build item search request for price list items
	// NetSuite requires specific filters for price list retrieval
	request := map[string]interface{}{
		"recordtype": "pricing",
		"filters": []map[string]interface{}{
			{
				"field":    "pricelist",
				"operator": "is",
				"value":    priceListID,
			},
			{
				"field":    "item",
				"operator": "anyof",
				"value":    productIDs,
			},
		},
		"columns": []string{"item", "pricelevel", "price", "currency", "quantity"},
	}

	// Call search RESTlet
	restletPath := "/restlets/search_records"
	if path, ok := c.Config.Metadata["search_records_restlet"].(string); ok {
		restletPath = path
	}

	result, err := c.callRESTlet(ctx, restletPath, request)
	if err != nil {
		return nil, fmt.Errorf("searching prices: %w", err)
	}

	// Transform results
	records, ok := result["records"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid search response format")
	}

	for _, record := range records {
		priceRecord := record.(map[string]interface{})
		price := c.transformPriceData(priceRecord, priceListID)
		if price != nil {
			prices = append(prices, price)
		}
	}

	c.Logger.Info().
		Int("found", len(prices)).
		Msg("Completed fetching prices from NetSuite")

	return prices, nil
}

// FetchProducts fetches specific products by their IDs
func (c *Connector) FetchProducts(ctx context.Context, productIDs []string) ([]*erp.ProductPayload, error) {
	c.Logger.Info().
		Int("count", len(productIDs)).
		Msg("Fetching products from NetSuite")

	if len(productIDs) == 0 {
		return []*erp.ProductPayload{}, nil
	}

	products := make([]*erp.ProductPayload, 0, len(productIDs))
	
	// Build item search request
	request := map[string]interface{}{
		"recordtype": "item",
		"filters": []map[string]interface{}{
			{
				"field":    "internalid",
				"operator": "anyof",
				"value":    productIDs,
			},
		},
		"columns": []string{
			"itemid", "displayname", "salesdescription", "class",
			"custitem_brand", "weight", "custitem_dimensions",
			"isinactive", "matrix", "type",
		},
	}

	// Call search RESTlet
	restletPath := "/restlets/search_records"
	if path, ok := c.Config.Metadata["search_records_restlet"].(string); ok {
		restletPath = path
	}

	result, err := c.callRESTlet(ctx, restletPath, request)
	if err != nil {
		return nil, fmt.Errorf("searching products: %w", err)
	}

	// Transform results
	records, ok := result["records"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid search response format")
	}

	for _, record := range records {
		itemRecord := record.(map[string]interface{})
		product := c.transformItemToProduct(itemRecord)
		if product != nil {
			products = append(products, product)
		}
	}

	c.Logger.Info().
		Int("found", len(products)).
		Msg("Completed fetching products from NetSuite")

	return products, nil
}

// FetchAllPrices fetches all prices changed since the given time
func (c *Connector) FetchAllPrices(ctx context.Context, since time.Time) ([]*erp.PricePayload, error) {
	c.Logger.Info().
		Time("since", since).
		Msg("Fetching all prices from NetSuite")

	// Query all price changes
	query := fmt.Sprintf(`
		SELECT 
			i.itemid,
			i.baseprice,
			ip.pricelevel,
			ip.price,
			c.symbol as currency,
			pl.name as pricelevelname,
			i.lastmodifieddate
		FROM item i
		LEFT JOIN itemprice ip ON i.id = ip.item
		LEFT JOIN pricelevel pl ON ip.pricelevel = pl.id
		LEFT JOIN currency c ON ip.currency = c.id
		WHERE i.lastmodifieddate >= TO_DATE('%s', 'YYYY-MM-DD HH24:MI:SS')
		ORDER BY i.lastmodifieddate ASC
		FETCH FIRST %d ROWS ONLY
	`, since.Format("2006-01-02 15:04:05"), c.Config.Sync.BatchSize)

	// Execute query
	result, err := c.executeSuiteQL(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("executing price query: %w", err)
	}

	// Transform results
	prices := make([]*erp.PricePayload, 0)
	if items, ok := result["items"].([]interface{}); ok {
		for _, item := range items {
			if data, ok := item.(map[string]interface{}); ok {
				// Create price for each price level
				priceLevel := getStringValue(data, "pricelevel", "1")
				price := c.transformPriceData(data, priceLevel)
				prices = append(prices, price)
			}
		}
	}

	c.Logger.Info().
		Int("count", len(prices)).
		Msg("Successfully fetched all prices from NetSuite")

	return prices, nil
}

// FetchAllStock fetches all stock changes since the given time
func (c *Connector) FetchAllStock(ctx context.Context, since time.Time) ([]*erp.StockPayload, error) {
	c.Logger.Info().
		Time("since", since).
		Msg("Fetching all stock levels from NetSuite")

	// Query inventory changes
	query := fmt.Sprintf(`
		SELECT 
			i.itemid,
			l.name as location,
			ib.quantityonhand,
			ib.quantityavailable,
			ib.quantityonorder,
			ib.quantityintransit,
			ib.lastmodifieddate
		FROM inventorybalance ib
		JOIN item i ON ib.item = i.id
		JOIN location l ON ib.location = l.id
		WHERE ib.lastmodifieddate >= TO_DATE('%s', 'YYYY-MM-DD HH24:MI:SS')
		AND l.isinactive = 'F'
		ORDER BY ib.lastmodifieddate ASC
		FETCH FIRST %d ROWS ONLY
	`, since.Format("2006-01-02 15:04:05"), c.Config.Sync.BatchSize)

	// Execute query
	result, err := c.executeSuiteQL(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("executing inventory query: %w", err)
	}

	// Transform results
	stocks := make([]*erp.StockPayload, 0)
	if items, ok := result["items"].([]interface{}); ok {
		for _, item := range items {
			if data, ok := item.(map[string]interface{}); ok {
				stock := c.transformInventoryToStock(data)
				stocks = append(stocks, stock)
			}
		}
	}

	c.Logger.Info().
		Int("count", len(stocks)).
		Msg("Successfully fetched all stock levels from NetSuite")

	return stocks, nil
}

// FetchStock fetches stock for specific products
func (c *Connector) FetchStock(ctx context.Context, productIDs []string) ([]*erp.StockPayload, error) {
	c.Logger.Info().
		Int("product_count", len(productIDs)).
		Msg("Fetching stock for specific products from NetSuite")

	if len(productIDs) == 0 {
		return []*erp.StockPayload{}, nil
	}

	// Build inventory search request
	request := map[string]interface{}{
		"recordtype": "inventoryitem",
		"filters": []map[string]interface{}{
			{
				"field":    "internalid",
				"operator": "anyof",
				"value":    productIDs,
			},
		},
		"columns": []string{
			"itemid", "quantityavailable", "quantityonhand", 
			"quantityonorder", "reorderpoint", "location",
		},
	}

	// Call search RESTlet
	restletPath := "/restlets/search_records"
	if path, ok := c.Config.Metadata["search_records_restlet"].(string); ok {
		restletPath = path
	}

	result, err := c.callRESTlet(ctx, restletPath, request)
	if err != nil {
		return nil, fmt.Errorf("searching inventory: %w", err)
	}

	// Transform results
	records, ok := result["records"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid search response format")
	}

	stockPayloads := make([]*erp.StockPayload, 0, len(records))
	for _, record := range records {
		invRecord := record.(map[string]interface{})
		stock := c.transformInventoryToStock(invRecord)
		if stock != nil {
			stockPayloads = append(stockPayloads, stock)
		}
	}

	c.Logger.Info().
		Int("found", len(stockPayloads)).
		Msg("Completed fetching stock from NetSuite")

	return stockPayloads, nil
}

// UpdateInventory processes inventory adjustments
func (c *Connector) UpdateInventory(ctx context.Context, adjustments []*erp.InventoryAdjustment) error {
	for _, adjustment := range adjustments {
		c.Logger.Info().
			Str("sku", adjustment.SKU).
			Str("type", string(adjustment.Type)).
			Int("delta", adjustment.QuantityDelta).
			Msg("Processing inventory adjustment in NetSuite")

		// Transform adjustment to NetSuite format
		nsAdjustment := map[string]interface{}{
			"recordtype": "inventoryadjustment",
			"subsidiary": "1", // Default subsidiary
			"account":    "1", // Default adjustment account
			"trandate":   time.Now().Format("01/02/2006"),
			"memo":       fmt.Sprintf("%s - %s", adjustment.ReferenceType, adjustment.Reason),
			"item": []map[string]interface{}{
				{
					"item":           adjustment.SKU,
					"location":       adjustment.WarehouseID,
					"adjustqtyby":    adjustment.QuantityDelta,
					"memo":           adjustment.Reason,
				},
			},
		}

		// Add reference fields
		if adjustment.ReferenceID != "" {
			nsAdjustment["custbody_reference_id"] = adjustment.ReferenceID
			nsAdjustment["custbody_reference_type"] = adjustment.ReferenceType
		}

		// Call inventory adjustment RESTlet
		restletPath := "/restlets/inventory_adjustment"
		if path, ok := c.Config.Metadata["inventory_adjustment_restlet"].(string); ok {
			restletPath = path
		}

		result, err := c.callRESTlet(ctx, restletPath, nsAdjustment)
		if err != nil {
			return fmt.Errorf("creating inventory adjustment: %w", err)
		}

		if adjustmentID, ok := result["id"].(string); ok {
			c.Logger.Info().
				Str("adjustment_id", adjustmentID).
				Str("sku", adjustment.SKU).
				Msg("Successfully created inventory adjustment in NetSuite")
		}
	}

	return nil
}

// ProcessReturn processes a return/RMA
func (c *Connector) ProcessReturn(ctx context.Context, returnData *erp.ReturnPayload) (*erp.ReturnPayload, error) {
	c.Logger.Info().
		Str("return_id", returnData.ReturnID).
		Str("order_id", returnData.OriginalOrderID).
		Msg("Processing return in NetSuite")

	// Transform return to NetSuite RMA format
	nsReturn := map[string]interface{}{
		"recordtype": "returnauthorization",
		"entity":     returnData.CustomerID,
		"trandate":   time.Now().Format("01/02/2006"),
		"memo":       fmt.Sprintf("Return %s - %s", returnData.ReturnID, returnData.Reason),
		"custbody_original_order": returnData.OriginalOrderID,
		"custbody_return_reason":  returnData.Reason,
		"location":   returnData.WarehouseID,
		"item":       make([]map[string]interface{}, 0, len(returnData.Items)),
	}

	// Add line items
	for _, item := range returnData.Items {
		nsReturn["item"] = append(nsReturn["item"].([]map[string]interface{}), map[string]interface{}{
			"item":     item.SKU,
			"quantity": item.Quantity,
			"rate":     item.RefundAmount / float64(item.Quantity),
			"custcol_return_reason": item.Reason,
			"description": item.Notes,
		})
	}

	// Call RMA creation RESTlet
	restletPath := "/restlets/create_return"
	if path, ok := c.Config.Metadata["create_return_restlet"].(string); ok {
		restletPath = path
	}

	result, err := c.callRESTlet(ctx, restletPath, nsReturn)
	if err != nil {
		return nil, fmt.Errorf("creating return authorization: %w", err)
	}

	rmaID, ok := result["id"].(string)
	if !ok {
		return nil, fmt.Errorf("return authorization ID not returned from NetSuite")
	}
	
	c.Logger.Info().
		Str("rma_id", rmaID).
		Str("return_id", returnData.ReturnID).
		Msg("Successfully created return authorization in NetSuite")
	
	// Return the updated return payload
	returnData.Status = "processed"
	returnData.ExternalID = rmaID
	
	return returnData, nil
}

// CreateInvoice creates a new invoice
func (c *Connector) CreateInvoice(ctx context.Context, invoice *erp.InvoicePayload) (string, error) {
	c.Logger.Info().
		Str("invoice_id", invoice.InvoiceID).
		Str("customer_id", invoice.CustomerID).
		Msg("Creating invoice in NetSuite")

	// Transform invoice to NetSuite format
	nsInvoice := map[string]interface{}{
		"recordtype": "invoice",
		"entity":     invoice.CustomerID,
		"trandate":   invoice.IssueDate.Format("01/02/2006"),
		"duedate":    invoice.DueDate.Format("01/02/2006"),
		"currency":   invoice.Currency,
		"terms":      invoice.PaymentTerms,
		"memo":       invoice.Notes,
		"otherrefnum": invoice.InvoiceNumber,
		"item":       make([]map[string]interface{}, 0, len(invoice.Lines)),
	}

	// Add order reference if provided
	if invoice.OrderID != "" {
		nsInvoice["custbody_order_reference"] = invoice.OrderID
	}

	// Add line items
	for _, line := range invoice.Lines {
		lineItem := map[string]interface{}{
			"description": line.Description,
			"quantity":    line.Quantity,
			"rate":        line.UnitPrice,
			"taxrate1":    line.TaxRate,
		}
		
		// Add SKU if provided
		if line.SKU != "" {
			lineItem["item"] = line.SKU
		}
		
		nsInvoice["item"] = append(nsInvoice["item"].([]map[string]interface{}), lineItem)
	}

	// Call invoice creation RESTlet
	restletPath := "/restlets/create_invoice"
	if path, ok := c.Config.Metadata["create_invoice_restlet"].(string); ok {
		restletPath = path
	}

	result, err := c.callRESTlet(ctx, restletPath, nsInvoice)
	if err != nil {
		return "", fmt.Errorf("creating invoice: %w", err)
	}

	invoiceID, ok := result["id"].(string)
	if !ok {
		return "", fmt.Errorf("invoice ID not returned from NetSuite")
	}

	c.Logger.Info().
		Str("netsuite_invoice_id", invoiceID).
		Str("invoice_id", invoice.InvoiceID).
		Msg("Successfully created invoice in NetSuite")

	return invoiceID, nil
}

// UpdateInvoice updates an existing invoice
func (c *Connector) UpdateInvoice(ctx context.Context, invoiceID string, invoice *erp.InvoicePayload) error {
	c.Logger.Info().
		Str("invoice_id", invoice.InvoiceID).
		Msg("Updating invoice in NetSuite")

	// Find NetSuite internal ID
	internalID, err := c.findInvoiceInternalID(ctx, invoice.InvoiceID)
	if err != nil {
		return fmt.Errorf("finding invoice: %w", err)
	}

	// Transform updates
	updates := map[string]interface{}{
		"recordtype": "invoice",
		"id":         internalID,
		"duedate":    invoice.DueDate.Format("01/02/2006"),
		"terms":      invoice.PaymentTerms,
		"memo":       invoice.Notes,
	}

	// Call invoice update RESTlet
	restletPath := "/restlets/update_record"
	if path, ok := c.Config.Metadata["update_record_restlet"].(string); ok {
		restletPath = path
	}

	result, err := c.callRESTlet(ctx, restletPath, updates)
	if err != nil {
		return fmt.Errorf("updating invoice: %w", err)
	}

	c.Logger.Info().
		Str("invoice_id", invoice.InvoiceID).
		Interface("result", result).
		Msg("Successfully updated invoice in NetSuite")

	return nil
}

// ApproveInvoice approves an invoice
func (c *Connector) ApproveInvoice(ctx context.Context, invoiceID string) error {
	c.Logger.Info().
		Str("invoice_id", invoiceID).
		Msg("Approving invoice in NetSuite")

	// Find NetSuite internal ID
	internalID, err := c.findInvoiceInternalID(ctx, invoiceID)
	if err != nil {
		return fmt.Errorf("finding invoice: %w", err)
	}

	// Call approval RESTlet
	restletPath := "/restlets/approve_invoice"
	if path, ok := c.Config.Metadata["approve_invoice_restlet"].(string); ok {
		restletPath = path
	}

	payload := map[string]interface{}{
		"invoiceId": internalID,
	}

	result, err := c.callRESTlet(ctx, restletPath, payload)
	if err != nil {
		return fmt.Errorf("approving invoice: %w", err)
	}

	c.Logger.Info().
		Str("invoice_id", invoiceID).
		Interface("result", result).
		Msg("Successfully approved invoice in NetSuite")

	return nil
}

// VoidInvoice voids an invoice
func (c *Connector) VoidInvoice(ctx context.Context, invoiceID string, reason string) error {
	c.Logger.Info().
		Str("invoice_id", invoiceID).
		Str("reason", reason).
		Msg("Voiding invoice in NetSuite")

	// Find NetSuite internal ID
	internalID, err := c.findInvoiceInternalID(ctx, invoiceID)
	if err != nil {
		return fmt.Errorf("finding invoice: %w", err)
	}

	// Call void RESTlet
	restletPath := "/restlets/void_invoice"
	if path, ok := c.Config.Metadata["void_invoice_restlet"].(string); ok {
		restletPath = path
	}

	payload := map[string]interface{}{
		"invoiceId": internalID,
		"reason":    reason,
	}

	result, err := c.callRESTlet(ctx, restletPath, payload)
	if err != nil {
		return fmt.Errorf("voiding invoice: %w", err)
	}

	c.Logger.Info().
		Str("invoice_id", invoiceID).
		Interface("result", result).
		Msg("Successfully voided invoice in NetSuite")

	return nil
}

// SendInvoice sends an invoice to customer
func (c *Connector) SendInvoice(ctx context.Context, invoiceID string, customerID string) error {
	c.Logger.Info().
		Str("invoice_id", invoiceID).
		Str("customer_id", customerID).
		Msg("Sending invoice to customer via NetSuite")

	// Find NetSuite internal ID
	internalID, err := c.findInvoiceInternalID(ctx, invoiceID)
	if err != nil {
		return fmt.Errorf("finding invoice: %w", err)
	}

	// Call send invoice RESTlet
	restletPath := "/restlets/send_invoice"
	if path, ok := c.Config.Metadata["send_invoice_restlet"].(string); ok {
		restletPath = path
	}

	payload := map[string]interface{}{
		"invoiceId":  internalID,
		"customerId": customerID,
		"emailBody":  "Please find attached your invoice.",
	}

	result, err := c.callRESTlet(ctx, restletPath, payload)
	if err != nil {
		return fmt.Errorf("sending invoice: %w", err)
	}

	c.Logger.Info().
		Str("invoice_id", invoiceID).
		Interface("result", result).
		Msg("Successfully sent invoice via NetSuite")

	return nil
}

// RecordPayment records a payment against an invoice
func (c *Connector) RecordPayment(ctx context.Context, payment *erp.PaymentPayload) error {
	c.Logger.Info().
		Str("payment_id", payment.PaymentID).
		Str("invoice_id", payment.InvoiceID).
		Float64("amount", payment.Amount).
		Msg("Recording payment in NetSuite")

	// Find NetSuite invoice internal ID
	invoiceInternalID, err := c.findInvoiceInternalID(ctx, payment.InvoiceID)
	if err != nil {
		return fmt.Errorf("finding invoice: %w", err)
	}

	// Create customer payment record
	nsPayment := map[string]interface{}{
		"recordtype": "customerpayment",
		"customer":   payment.CustomerID,
		"payment":    payment.Amount,
		"currency":   payment.Currency,
		"trandate":   payment.PaymentDate.Format("01/02/2006"),
		"memo":       payment.Notes,
		"checknum":   payment.Reference,
		"apply": []map[string]interface{}{
			{
				"doc":    invoiceInternalID,
				"apply":  true,
				"amount": payment.Amount,
			},
		},
	}

	// Map payment method
	if paymentMethod, ok := mapPaymentMethod(payment.PaymentMethod); ok {
		nsPayment["paymentmethod"] = paymentMethod
	}

	// Call payment creation RESTlet
	restletPath := "/restlets/create_payment"
	if path, ok := c.Config.Metadata["create_payment_restlet"].(string); ok {
		restletPath = path
	}

	result, err := c.callRESTlet(ctx, restletPath, nsPayment)
	if err != nil {
		return fmt.Errorf("recording payment: %w", err)
	}

	if paymentID, ok := result["id"].(string); ok {
		c.Logger.Info().
			Str("netsuite_payment_id", paymentID).
			Str("payment_id", payment.PaymentID).
			Msg("Successfully recorded payment in NetSuite")
	}

	return nil
}

// Helper method to find invoice internal ID
func (c *Connector) findInvoiceInternalID(ctx context.Context, invoiceID string) (string, error) {
	query := fmt.Sprintf(`
		SELECT id
		FROM transaction
		WHERE type = 'CustInvc'
		AND otherrefnum = '%s'
		FETCH FIRST 1 ROWS ONLY
	`, invoiceID)

	result, err := c.executeSuiteQL(ctx, query)
	if err != nil {
		return "", err
	}

	if items, ok := result["items"].([]interface{}); ok && len(items) > 0 {
		if data, ok := items[0].(map[string]interface{}); ok {
			if id, ok := data["id"].(string); ok {
				return id, nil
			}
		}
	}

	return "", fmt.Errorf("invoice not found: %s", invoiceID)
}

// Helper function to map payment methods
func mapPaymentMethod(method string) (string, bool) {
	methodMap := map[string]string{
		"credit_card": "1",  // Credit Card
		"cash":        "2",  // Cash
		"check":       "3",  // Check
		"wire":        "4",  // Wire Transfer
		"ach":         "5",  // ACH
		"paypal":      "6",  // PayPal
	}

	if nsMethod, ok := methodMap[strings.ToLower(method)]; ok {
		return nsMethod, true
	}
	return "", false
}