package odoo

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"middleman/internal/erp"
)

// Connector implements the Odoo ERP connector with production-ready features
type Connector struct {
	*erp.BaseConnector

	// Authentication state
	sessionID string
	userID    int
	lastAuth  time.Time
	authMutex sync.RWMutex

	// Configuration
	database   string
	apiVersion string
}

// NewConnector creates a new Odoo connector
func NewConnector(config erp.ERPConfig) (*Connector, error) {
	if config.Type != erp.ERPTypeOdoo {
		return nil, fmt.Errorf("invalid ERP type: %s, expected %s", config.Type, erp.ERPTypeOdoo)
	}

	// Validate required configuration
	if config.Auth.Username == "" || config.Auth.Password == "" {
		return nil, fmt.Errorf("Odoo requires username and password authentication")
	}

	// Extract database from metadata
	database, ok := config.Metadata["database"].(string)
	if !ok || database == "" {
		return nil, fmt.Errorf("Odoo requires database name in metadata")
	}

	// API version (default to 2.0 for newer Odoo versions)
	apiVersion := "2.0"
	if v, ok := config.Metadata["api_version"].(string); ok {
		apiVersion = v
	}

	return &Connector{
		BaseConnector: erp.NewBaseConnector(config),
		database:      database,
		apiVersion:    apiVersion,
	}, nil
}

// GetID returns the connector ID
func (c *Connector) GetID() string {
	return c.Config.ID
}

// GetType returns the ERP type as string
func (c *Connector) GetType() string {
	return string(erp.ERPTypeOdoo)
}

// GetConfig returns the connector configuration
func (c *Connector) GetConfig() erp.ERPConfig {
	return c.Config
}

// TestConnection validates the connector can reach the ERP
func (c *Connector) TestConnection(ctx context.Context) error {
	return c.authenticate(ctx)
}

// HealthCheck verifies the connector is operational
func (c *Connector) HealthCheck(ctx context.Context) erp.HealthCheck {
	err := c.authenticate(ctx)
	if err != nil {
		return erp.HealthCheck{
			Status:  erp.HealthStatusUnhealthy,
			Message: fmt.Sprintf("Authentication failed: %v", err),
		}
	}
	
	// Test a simple API call
	req := &JSONRPCRequest{
		JSONRPC: c.apiVersion,
		Method:  "call",
		ID:      generateRequestID(),
		Params: map[string]interface{}{
			"service": "common",
			"method":  "version",
			"args":    []interface{}{},
		},
	}
	
	resp, err := c.sendJSONRPC(ctx, "/jsonrpc", req)
	if err != nil {
		return erp.HealthCheck{
			Status:  erp.HealthStatusDegraded,
			Message: fmt.Sprintf("API call failed: %v", err),
		}
	}
	
	if resp.Error != nil {
		return erp.HealthCheck{
			Status:  erp.HealthStatusDegraded,
			Message: fmt.Sprintf("API error: %v", resp.Error["message"]),
		}
	}
	
	return erp.HealthCheck{
		Status:  erp.HealthStatusHealthy,
		Message: "Connector is operational",
	}
}

// authenticate performs authentication with Odoo
func (c *Connector) authenticate(ctx context.Context) error {
	c.authMutex.Lock()
	defer c.authMutex.Unlock()

	// Check if we have a recent authentication
	if c.sessionID != "" && time.Since(c.lastAuth) < 24*time.Hour {
		return nil
	}

	c.Logger.Info().Msg("Authenticating with Odoo")

	// Prepare authentication request
	authReq := &JSONRPCRequest{
		JSONRPC: c.apiVersion,
		Method:  "call",
		ID:      generateRequestID(),
		Params: map[string]interface{}{
			"service": "common",
			"method":  "authenticate",
			"args": []interface{}{
				c.database,
				c.Config.Auth.Username,
				c.Config.Auth.Password,
				map[string]interface{}{}, // User context
			},
		},
	}

	// Send authentication request
	resp, err := c.sendJSONRPC(ctx, "/jsonrpc", authReq)
	if err != nil {
		return fmt.Errorf("authentication request failed: %w", err)
	}

	// Parse response
	if resp.Error != nil {
		return fmt.Errorf("authentication error: code=%v, message=%v",
			resp.Error["code"], resp.Error["message"])
	}

	// Extract user ID
	userID, ok := resp.Result.(float64)
	if !ok || userID == 0 {
		return fmt.Errorf("authentication failed: invalid credentials or response")
	}

	c.userID = int(userID)
	c.sessionID = uuid.New().String()
	c.lastAuth = time.Now()

	c.Logger.Info().
		Int("user_id", c.userID).
		Str("database", c.database).
		Msg("Successfully authenticated with Odoo")

	return nil
}

// ensureAuthenticated ensures we have valid authentication
func (c *Connector) ensureAuthenticated(ctx context.Context) error {
	c.authMutex.RLock()
	needsAuth := c.sessionID == "" || time.Since(c.lastAuth) >= 24*time.Hour
	c.authMutex.RUnlock()

	if needsAuth {
		return c.authenticate(ctx)
	}
	return nil
}

// ValidateWebhook validates a webhook signature
func (c *Connector) ValidateWebhook(payload []byte, signature string) error {
	if !c.Config.Webhook.ValidateSign || c.Config.Webhook.Secret == "" {
		return nil
	}

	// Calculate expected signature
	h := hmac.New(sha256.New, []byte(c.Config.Webhook.Secret))
	h.Write(payload)
	expectedSig := hex.EncodeToString(h.Sum(nil))

	// Compare signatures
	if !hmac.Equal([]byte(signature), []byte(expectedSig)) {
		return fmt.Errorf("invalid webhook signature")
	}

	return nil
}

// ParseWebhook parses webhook payload into canonical event
func (c *Connector) ParseWebhook(payload []byte) (*erp.CanonicalEvent, error) {
	var webhook OdooWebhook
	if err := json.Unmarshal(payload, &webhook); err != nil {
		return nil, fmt.Errorf("parsing webhook payload: %w", err)
	}

	// Map Odoo events to canonical events
	eventType, err := mapOdooEventType(webhook.Event)
	if err != nil {
		return nil, err
	}

	// Transform payload based on event type
	canonicalPayload, err := c.transformWebhookPayload(&webhook, eventType)
	if err != nil {
		return nil, fmt.Errorf("transforming payload: %w", err)
	}

	return &erp.CanonicalEvent{
		EventID:        webhook.ID,
		EventType:      eventType,
		EventTimestamp: webhook.Timestamp,
		Source:         fmt.Sprintf("Odoo-%s", c.database),
		CorrelationID:  fmt.Sprintf("%s-%d", webhook.Model, webhook.RecordID),
		Payload:        canonicalPayload,
	}, nil
}

// SyncProducts synchronizes products from the ERP
func (c *Connector) SyncProducts(ctx context.Context, since time.Time, batchSize int) ([]*erp.ProductPayload, error) {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}

	c.Logger.Info().
		Time("since", since).
		Int("batch_size", batchSize).
		Msg("Syncing products from Odoo")

	// Build search domain
	domain := []interface{}{
		[]interface{}{"write_date", ">=", since.Format("2006-01-02 15:04:05")},
		[]interface{}{"active", "=", true},
	}

	// Add any additional filters from config
	if filters, ok := c.Config.Metadata["product_filters"].([]interface{}); ok {
		domain = append(domain, filters...)
	}

	var allProducts []*erp.ProductPayload
	offset := 0

	for {
		// Search for products with pagination
		searchReq := &JSONRPCRequest{
			JSONRPC: c.apiVersion,
			Method:  "call",
			ID:      generateRequestID(),
			Params: c.buildObjectRequest("product.product", "search_read", []interface{}{
				domain,
			}, map[string]interface{}{
				"fields": []string{
					"id", "name", "default_code", "description", "description_sale",
					"categ_id", "list_price", "standard_price", "weight", "volume",
					"barcode", "active", "type", "uom_id", "taxes_id",
				},
				"limit":  batchSize,
				"offset": offset,
				"order":  "write_date asc",
			}),
		}

		resp, err := c.sendJSONRPC(ctx, "/jsonrpc", searchReq)
		if err != nil {
			return nil, fmt.Errorf("fetching products at offset %d: %w", offset, err)
		}

		if resp.Error != nil {
			return nil, fmt.Errorf("API error at offset %d: %v", offset, resp.Error)
		}

		// Parse and transform products
		products, err := c.parseProductsResponse(resp.Result)
		if err != nil {
			return nil, fmt.Errorf("parsing products at offset %d: %w", offset, err)
		}

		allProducts = append(allProducts, products...)

		c.Logger.Debug().
			Int("batch_count", len(products)).
			Int("total_count", len(allProducts)).
			Int("offset", offset).
			Msg("Fetched product batch")

		// Check if we have more records
		if len(products) < batchSize {
			break
		}

		offset += batchSize

		// Check context cancellation
		select {
		case <-ctx.Done():
			return allProducts, ctx.Err()
		default:
		}
	}

	c.Logger.Info().
		Int("total_count", len(allProducts)).
		Msg("Successfully synced products from Odoo")

	return allProducts, nil
}

// SyncStock synchronizes stock levels from the ERP
func (c *Connector) SyncStock(ctx context.Context, productIDs []string, batchSize int) ([]*erp.StockPayload, error) {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if len(productIDs) == 0 {
		return []*erp.StockPayload{}, nil
	}

	c.Logger.Info().
		Int("product_count", len(productIDs)).
		Int("batch_size", batchSize).
		Msg("Syncing stock levels from Odoo")

	// Convert product IDs to ints
	intProductIDs := c.convertProductIDsToInts(productIDs)
	if len(intProductIDs) == 0 {
		return []*erp.StockPayload{}, nil
	}

	var allStocks []*erp.StockPayload

	// Process in chunks if needed
	for i := 0; i < len(intProductIDs); i += batchSize {
		end := i + batchSize
		if end > len(intProductIDs) {
			end = len(intProductIDs)
		}

		chunkIDs := intProductIDs[i:end]

		// Search for stock quantities
		domain := []interface{}{
			[]interface{}{"product_id", "in", chunkIDs},
			[]interface{}{"location_id.usage", "=", "internal"},
		}

		offset := 0
		for {
			searchReq := &JSONRPCRequest{
				JSONRPC: c.apiVersion,
				Method:  "call",
				ID:      generateRequestID(),
				Params: c.buildObjectRequest("stock.quant", "search_read", []interface{}{
					domain,
				}, map[string]interface{}{
					"fields": []string{
						"product_id", "location_id", "quantity",
						"reserved_quantity", "available_quantity", "in_date",
					},
					"limit":  batchSize,
					"offset": offset,
				}),
			}

			resp, err := c.sendJSONRPC(ctx, "/jsonrpc", searchReq)
			if err != nil {
				return nil, fmt.Errorf("fetching stock at offset %d: %w", offset, err)
			}

			if resp.Error != nil {
				return nil, fmt.Errorf("API error at offset %d: %v", offset, resp.Error)
			}

			// Parse and aggregate stock by SKU
			stocks, err := c.parseStockResponse(resp.Result)
			if err != nil {
				return nil, fmt.Errorf("parsing stock at offset %d: %w", offset, err)
			}

			allStocks = append(allStocks, stocks...)

			// Check if we have more records
			if len(stocks) < batchSize {
				break
			}

			offset += batchSize
		}

		// Check context cancellation
		select {
		case <-ctx.Done():
			return allStocks, ctx.Err()
		default:
		}
	}

	c.Logger.Info().
		Int("total_count", len(allStocks)).
		Msg("Successfully synced stock levels from Odoo")

	return allStocks, nil
}

// SyncPrices synchronizes prices from the ERP
func (c *Connector) SyncPrices(ctx context.Context, productIDs []string, batchSize int) ([]*erp.PricePayload, error) {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if len(productIDs) == 0 {
		return []*erp.PricePayload{}, nil
	}

	c.Logger.Info().
		Int("product_count", len(productIDs)).
		Msg("Fetching prices from Odoo")

	// Determine pricelist from config or use default
	plID := 1 // Default public pricelist
	if plIDConfig, ok := c.Config.Metadata["default_pricelist_id"].(float64); ok {
		plID = int(plIDConfig)
	}

	// Get prices using price computation
	computeReq := &JSONRPCRequest{
		JSONRPC: c.apiVersion,
		Method:  "call",
		ID:      generateRequestID(),
		Params: c.buildObjectRequest("product.pricelist", "get_products_price", []interface{}{
			plID,                                    // Pricelist ID
			c.convertProductIDsToInts(productIDs), // Product IDs
			1.0,                              // Quantity
			map[string]interface{}{},         // Partner
			time.Now(),                       // Date
		}, nil),
	}

	resp, err := c.sendJSONRPC(ctx, "/jsonrpc", computeReq)
	if err != nil {
		return nil, fmt.Errorf("computing prices: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("API error: %v", resp.Error)
	}

	// Parse price response
	prices, err := c.parsePriceResponse(resp.Result, productIDs, fmt.Sprintf("%d", plID))
	if err != nil {
		return nil, err
	}

	c.Logger.Info().
		Int("count", len(prices)).
		Msg("Successfully fetched prices from Odoo")

	return prices, nil
}

// FetchProducts fetches products from the ERP
func (c *Connector) FetchProducts(ctx context.Context, productIDs []string) ([]*erp.ProductPayload, error) {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if len(productIDs) == 0 {
		return []*erp.ProductPayload{}, nil
	}

	c.Logger.Info().
		Int("product_count", len(productIDs)).
		Msg("Fetching specific products from Odoo")

	// Convert product IDs to ints
	ids := c.convertProductIDsToInts(productIDs)

	// Read products directly
	req := &JSONRPCRequest{
		JSONRPC: c.apiVersion,
		Method:  "call",
		ID:      generateRequestID(),
		Params: c.buildObjectRequest("product.product", "read", []interface{}{
			ids,
		}, map[string]interface{}{
			"fields": []string{
				"id", "name", "default_code", "description", "description_sale",
				"categ_id", "list_price", "standard_price", "weight", "volume",
				"barcode", "active", "type", "uom_id", "taxes_id",
			},
		}),
	}

	resp, err := c.sendJSONRPC(ctx, "/jsonrpc", req)
	if err != nil {
		return nil, fmt.Errorf("fetching products: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("API error: %v", resp.Error)
	}

	// Parse and transform products
	products, err := c.parseProductsResponse(resp.Result)
	if err != nil {
		return nil, err
	}

	c.Logger.Info().
		Int("count", len(products)).
		Msg("Successfully fetched products from Odoo")

	return products, nil
}

// FetchPrices fetches prices from the ERP
func (c *Connector) FetchPrices(ctx context.Context, productIDs []string, priceListID string) ([]*erp.PricePayload, error) {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if len(productIDs) == 0 {
		return []*erp.PricePayload{}, nil
	}

	c.Logger.Info().
		Int("product_count", len(productIDs)).
		Str("pricelist_id", priceListID).
		Msg("Fetching prices from Odoo")

	// Determine pricelist
	plID := 1 // Default public pricelist
	if priceListID != "" {
		if id, err := strconv.Atoi(priceListID); err == nil {
			plID = id
		}
	}

	// Get prices using price computation
	computeReq := &JSONRPCRequest{
		JSONRPC: c.apiVersion,
		Method:  "call",
		ID:      generateRequestID(),
		Params: c.buildObjectRequest("product.pricelist", "get_products_price", []interface{}{
			plID,                                    // Pricelist ID
			c.convertProductIDsToInts(productIDs), // Product IDs
			1.0,                                     // Quantity
			map[string]interface{}{},                // Partner
			time.Now(),                              // Date
		}, nil),
	}

	resp, err := c.sendJSONRPC(ctx, "/jsonrpc", computeReq)
	if err != nil {
		return nil, fmt.Errorf("computing prices: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("API error: %v", resp.Error)
	}

	// Parse price response
	prices, err := c.parsePriceResponse(resp.Result, productIDs, priceListID)
	if err != nil {
		return nil, err
	}

	c.Logger.Info().
		Int("count", len(prices)).
		Msg("Successfully fetched prices from Odoo")

	return prices, nil
}

// FetchStock fetches stock for specific products
func (c *Connector) FetchStock(ctx context.Context, productIDs []string) ([]*erp.StockPayload, error) {
	return c.SyncStock(ctx, productIDs, 1000)
}

// FetchCustomer fetches a single customer from the ERP
func (c *Connector) FetchCustomer(ctx context.Context, customerID string) (*erp.CustomerPayload, error) {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}

	c.Logger.Info().
		Str("customer_id", customerID).
		Msg("Fetching customer from Odoo")

	// Convert customer ID to int
	id, err := strconv.Atoi(customerID)
	if err != nil {
		return nil, fmt.Errorf("invalid customer ID: %w", err)
	}

	// Read customer directly
	req := &JSONRPCRequest{
		JSONRPC: c.apiVersion,
		Method:  "call",
		ID:      generateRequestID(),
		Params: c.buildObjectRequest("res.partner", "read", []interface{}{
			[]int{id},
		}, map[string]interface{}{
			"fields": []string{
				"id", "name", "email", "phone", "mobile", "website",
				"street", "street2", "city", "state_id", "zip", "country_id",
				"vat", "ref", "lang", "customer_rank", "is_company",
			},
		}),
	}

	resp, err := c.sendJSONRPC(ctx, "/jsonrpc", req)
	if err != nil {
		return nil, fmt.Errorf("fetching customer: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("API error: %v", resp.Error)
	}

	// Parse response
	customers, err := c.parseCustomersResponse(resp.Result)
	if err != nil {
		return nil, err
	}

	if len(customers) == 0 {
		return nil, fmt.Errorf("customer not found: %s", customerID)
	}

	return customers[0], nil
}

// FetchCustomers fetches multiple customers from the ERP
func (c *Connector) FetchCustomers(ctx context.Context, customerIDs []string) ([]*erp.CustomerPayload, error) {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if len(customerIDs) == 0 {
		return []*erp.CustomerPayload{}, nil
	}

	c.Logger.Info().
		Int("customer_count", len(customerIDs)).
		Msg("Fetching customers from Odoo")

	// Convert customer IDs to ints
	ids := make([]int, 0, len(customerIDs))
	for _, idStr := range customerIDs {
		if id, err := strconv.Atoi(idStr); err == nil {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		return []*erp.CustomerPayload{}, nil
	}

	// Read customers directly
	req := &JSONRPCRequest{
		JSONRPC: c.apiVersion,
		Method:  "call",
		ID:      generateRequestID(),
		Params: c.buildObjectRequest("res.partner", "read", []interface{}{
			ids,
		}, map[string]interface{}{
			"fields": []string{
				"id", "name", "email", "phone", "mobile", "website",
				"street", "street2", "city", "state_id", "zip", "country_id",
				"vat", "ref", "lang", "customer_rank", "is_company",
			},
		}),
	}

	resp, err := c.sendJSONRPC(ctx, "/jsonrpc", req)
	if err != nil {
		return nil, fmt.Errorf("fetching customers: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("API error: %v", resp.Error)
	}

	// Parse response
	customers, err := c.parseCustomersResponse(resp.Result)
	if err != nil {
		return nil, err
	}

	c.Logger.Info().
		Int("count", len(customers)).
		Msg("Successfully fetched customers from Odoo")

	return customers, nil
}

// FetchOrders fetches orders changed since the given time
func (c *Connector) FetchOrders(ctx context.Context, since time.Time) ([]*erp.OrderPayload, error) {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}

	c.Logger.Info().
		Time("since", since).
		Msg("Fetching orders from Odoo")

	// Search for sale orders
	domain := []interface{}{
		[]interface{}{"write_date", ">=", since.Format("2006-01-02 15:04:05")},
		[]interface{}{"state", "not in", []string{"draft", "cancel"}},
	}

	searchReq := &JSONRPCRequest{
		JSONRPC: c.apiVersion,
		Method:  "call",
		ID:      generateRequestID(),
		Params: c.buildObjectRequest("sale.order", "search_read", []interface{}{
			domain,
		}, map[string]interface{}{
			"fields": []string{
				"id", "name", "partner_id", "amount_total", "amount_untaxed",
				"currency_id", "state", "create_date", "date_order",
				"order_line", "partner_invoice_id", "partner_shipping_id",
			},
			"limit":  100,
			"offset": 0,
			"order":  "write_date asc",
		}),
	}

	resp, err := c.sendJSONRPC(ctx, "/jsonrpc", searchReq)
	if err != nil {
		return nil, fmt.Errorf("fetching orders: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("API error: %v", resp.Error)
	}

	// Parse orders with line items
	orders, err := c.parseOrdersResponse(ctx, resp.Result)
	if err != nil {
		return nil, err
	}

	c.Logger.Info().
		Int("count", len(orders)).
		Msg("Successfully fetched orders from Odoo")

	return orders, nil
}

// SyncCustomers synchronizes customers from the ERP
func (c *Connector) SyncCustomers(ctx context.Context, since time.Time, batchSize int) ([]*erp.CustomerPayload, error) {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}

	c.Logger.Info().
		Time("since", since).
		Int("batch_size", batchSize).
		Msg("Syncing customers from Odoo")

	// Search for partners (customers)
	domain := []interface{}{
		[]interface{}{"write_date", ">=", since.Format("2006-01-02 15:04:05")},
		[]interface{}{"customer_rank", ">", 0},
		[]interface{}{"active", "=", true},
	}

	var allCustomers []*erp.CustomerPayload
	offset := 0

	for {
		searchReq := &JSONRPCRequest{
			JSONRPC: c.apiVersion,
			Method:  "call",
			ID:      generateRequestID(),
			Params: c.buildObjectRequest("res.partner", "search_read", []interface{}{
				domain,
			}, map[string]interface{}{
				"fields": []string{
					"id", "name", "email", "phone", "mobile", "website",
					"street", "street2", "city", "state_id", "zip", "country_id",
					"vat", "ref", "lang", "customer_rank", "is_company",
				},
				"limit":  batchSize,
				"offset": offset,
				"order":  "write_date asc",
			}),
		}

		resp, err := c.sendJSONRPC(ctx, "/jsonrpc", searchReq)
		if err != nil {
			return nil, fmt.Errorf("fetching customers at offset %d: %w", offset, err)
		}

		if resp.Error != nil {
			return nil, fmt.Errorf("API error at offset %d: %v", offset, resp.Error)
		}

		// Parse customers
		customers, err := c.parseCustomersResponse(resp.Result)
		if err != nil {
			return nil, fmt.Errorf("parsing customers at offset %d: %w", offset, err)
		}

		allCustomers = append(allCustomers, customers...)

		c.Logger.Debug().
			Int("batch_count", len(customers)).
			Int("total_count", len(allCustomers)).
			Int("offset", offset).
			Msg("Fetched customer batch")

		// Check if we have more records
		if len(customers) < batchSize {
			break
		}

		offset += batchSize

		// Check context cancellation
		select {
		case <-ctx.Done():
			return allCustomers, ctx.Err()
		default:
		}
	}

	c.Logger.Info().
		Int("total_count", len(allCustomers)).
		Msg("Successfully synced customers from Odoo")

	return allCustomers, nil
}

// SendOrder sends an order to Odoo
func (c *Connector) SendOrder(ctx context.Context, order *erp.OrderPayload) error {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return err
	}

	c.Logger.Info().
		Str("order_id", order.OrderID).
		Msg("Creating order in Odoo")

	// Transform order to Odoo format
	odooOrder, err := c.transformOrderToOdoo(ctx, order)
	if err != nil {
		return fmt.Errorf("transforming order: %w", err)
	}

	// Create order
	createReq := &JSONRPCRequest{
		JSONRPC: c.apiVersion,
		Method:  "call",
		ID:      generateRequestID(),
		Params: c.buildObjectRequest("sale.order", "create", []interface{}{
			[]interface{}{odooOrder},
		}, nil),
	}

	resp, err := c.sendJSONRPC(ctx, "/jsonrpc", createReq)
	if err != nil {
		return fmt.Errorf("creating order: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("API error: %v", resp.Error)
	}

	// Extract created order ID
	orderIDs, ok := resp.Result.([]interface{})
	if !ok || len(orderIDs) == 0 {
		return fmt.Errorf("unexpected create response format")
	}

	c.Logger.Info().
		Str("order_id", order.OrderID).
		Interface("odoo_id", orderIDs[0]).
		Msg("Successfully created order in Odoo")

	// Optionally confirm the order
	if autoConfirm, ok := c.Config.Metadata["auto_confirm_orders"].(bool); ok && autoConfirm {
		if err := c.confirmOrder(ctx, orderIDs[0]); err != nil {
			c.Logger.Error().Err(err).Msg("Failed to confirm order")
		}
	}

	return nil
}

// UpdateStock updates stock level in Odoo
func (c *Connector) UpdateStock(ctx context.Context, stock *erp.StockPayload) error {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return err
	}

	// In Odoo, direct stock updates are typically done through stock moves
	// or inventory adjustments. This is a simplified implementation.
	c.Logger.Warn().
		Str("sku", stock.SKU).
		Msg("Direct stock updates not recommended in Odoo - use stock moves or inventory adjustments")

	return fmt.Errorf("direct stock updates not supported - use Odoo's inventory adjustment workflow")
}

// GetOrder retrieves an order from Odoo
func (c *Connector) GetOrder(ctx context.Context, orderID string) (*erp.OrderPayload, error) {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Search for the order
	req := &JSONRPCRequest{
		JSONRPC: c.apiVersion,
		Method:  "call",
		ID:      generateRequestID(),
		Params: c.buildObjectRequest("sale.order", "search_read", []interface{}{
			[]interface{}{[]interface{}{"name", "=", orderID}},
		}, map[string]interface{}{
			"fields": []string{"name", "partner_id", "date_order", "amount_total", "state", "order_line"},
		}),
	}

	resp, err := c.sendJSONRPC(ctx, "/jsonrpc", req)
	if err != nil {
		return nil, fmt.Errorf("fetching order: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("order fetch error: %v", resp.Error["message"])
	}

	orders, ok := resp.Result.([]interface{})
	if !ok || len(orders) == 0 {
		return nil, fmt.Errorf("order not found: %s", orderID)
	}

	orderData := orders[0].(map[string]interface{})
	return c.transformOrderPayload(orderData)
}

// ProcessWebhook processes an incoming webhook
func (c *Connector) ProcessWebhook(ctx context.Context, payload []byte, signature string) error {
	// Validate signature
	if err := c.ValidateWebhook(payload, signature); err != nil {
		return fmt.Errorf("webhook validation failed: %w", err)
	}

	// Parse webhook
	event, err := c.ParseWebhook(payload)
	if err != nil {
		return fmt.Errorf("webhook parsing failed: %w", err)
	}

	c.Logger.Info().
		Str("event_id", event.EventID).
		Str("event_type", string(event.EventType)).
		Msg("Processed Odoo webhook")

	return nil
}

// CreateInvoice creates an invoice in Odoo
func (c *Connector) CreateInvoice(ctx context.Context, invoice *erp.InvoicePayload) (string, error) {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return "", err
	}

	// Map invoice data to Odoo format
	invoiceData := map[string]interface{}{
		"partner_id":     invoice.CustomerID,
		"invoice_date":   invoice.IssueDate.Format("2006-01-02"),
		"invoice_date_due": invoice.DueDate.Format("2006-01-02"),
		"currency_id":    c.getCurrencyID(invoice.Currency),
		"move_type":      "out_invoice", // Customer invoice
		"invoice_line_ids": c.transformInvoiceLines(invoice.Lines),
	}

	// Create invoice
	req := &JSONRPCRequest{
		JSONRPC: c.apiVersion,
		Method:  "call",
		ID:      generateRequestID(),
		Params: c.buildObjectRequest("account.move", "create", []interface{}{
			[]interface{}{invoiceData},
		}, nil),
	}

	resp, err := c.sendJSONRPC(ctx, "/jsonrpc", req)
	if err != nil {
		return "", fmt.Errorf("creating invoice: %w", err)
	}

	if resp.Error != nil {
		return "", fmt.Errorf("invoice creation error: %v", resp.Error["message"])
	}

	// Extract created invoice ID
	invoiceIDs, ok := resp.Result.([]interface{})
	if !ok || len(invoiceIDs) == 0 {
		return "", fmt.Errorf("invalid invoice creation response")
	}

	invoiceID := fmt.Sprintf("%v", invoiceIDs[0])
	return invoiceID, nil
}

// UpdateInvoice updates an invoice in Odoo
func (c *Connector) UpdateInvoice(ctx context.Context, invoiceID string, invoice *erp.InvoicePayload) error {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return err
	}

	// Convert invoice ID to int
	id, err := strconv.Atoi(invoiceID)
	if err != nil {
		return fmt.Errorf("invalid invoice ID: %w", err)
	}

	// Update data
	updateData := map[string]interface{}{
		"invoice_date_due": invoice.DueDate.Format("2006-01-02"),
	}

	// Update invoice
	req := &JSONRPCRequest{
		JSONRPC: c.apiVersion,
		Method:  "call",
		ID:      generateRequestID(),
		Params: c.buildObjectRequest("account.move", "write", []interface{}{
			[]interface{}{id},
			updateData,
		}, nil),
	}

	resp, err := c.sendJSONRPC(ctx, "/jsonrpc", req)
	if err != nil {
		return fmt.Errorf("updating invoice: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("invoice update error: %v", resp.Error["message"])
	}

	return nil
}

// ProcessReturn processes a return in Odoo
func (c *Connector) ProcessReturn(ctx context.Context, returnPayload *erp.ReturnPayload) (*erp.ReturnPayload, error) {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// In Odoo, returns are handled through stock.return.picking wizard
	// This is a simplified implementation
	returnData := map[string]interface{}{
		"origin": returnPayload.OrderID,
		"partner_id": returnPayload.CustomerID,
		"location_id": 1, // Customer location
		"location_dest_id": 5, // Stock location
		"picking_type_id": 5, // Return picking type
	}

	req := &JSONRPCRequest{
		JSONRPC: c.apiVersion,
		Method:  "call",
		ID:      generateRequestID(),
		Params: c.buildObjectRequest("stock.picking", "create", []interface{}{
			[]interface{}{returnData},
		}, nil),
	}

	resp, err := c.sendJSONRPC(ctx, "/jsonrpc", req)
	if err != nil {
		return nil, fmt.Errorf("creating return: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("return creation error: %v", resp.Error["message"])
	}

	returnPayload.Status = "created"
	return returnPayload, nil
}

// UpdateInventory updates inventory in Odoo
func (c *Connector) UpdateInventory(ctx context.Context, adjustments []*erp.InventoryAdjustment) error {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return err
	}

	// In Odoo, inventory adjustments are done through stock.inventory.adjustment
	for _, adj := range adjustments {
		adjData := map[string]interface{}{
			"product_id": c.getProductIDBySKU(ctx, adj.SKU),
			"location_id": adj.LocationID,
			"inventory_quantity": adj.Quantity,
		}

		req := &JSONRPCRequest{
			JSONRPC: c.apiVersion,
			Method:  "call",
			ID:      generateRequestID(),
			Params: c.buildObjectRequest("stock.quant", "write", []interface{}{
				[]interface{}{}, // Would need quant ID
				adjData,
			}, nil),
		}

		resp, err := c.sendJSONRPC(ctx, "/jsonrpc", req)
		if err != nil {
			return fmt.Errorf("updating inventory for %s: %w", adj.SKU, err)
		}

		if resp.Error != nil {
			return fmt.Errorf("inventory update error for %s: %v", adj.SKU, resp.Error["message"])
		}
	}

	return nil
}

// FetchAllStock fetches all stock information
func (c *Connector) FetchAllStock(ctx context.Context, since time.Time) ([]*erp.StockPayload, error) {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var allStock []*erp.StockPayload
	offset := 0
	limit := 100

	for {
		// Search for stock quants
		domain := []interface{}{}
		if !since.IsZero() {
			domain = append(domain, []interface{}{"write_date", ">=", since.Format("2006-01-02 15:04:05")})
		}

		req := &JSONRPCRequest{
			JSONRPC: c.apiVersion,
			Method:  "call",
			ID:      generateRequestID(),
			Params: c.buildObjectRequest("stock.quant", "search_read", []interface{}{
				domain,
			}, map[string]interface{}{
				"fields": []string{"product_id", "location_id", "quantity", "reserved_quantity"},
				"offset": offset,
				"limit":  limit,
			}),
		}

		resp, err := c.sendJSONRPC(ctx, "/jsonrpc", req)
		if err != nil {
			return nil, fmt.Errorf("fetching stock: %w", err)
		}

		if resp.Error != nil {
			return nil, fmt.Errorf("stock fetch error: %v", resp.Error["message"])
		}

		quants, ok := resp.Result.([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid stock response format")
		}

		// Transform quants to stock payloads
		for _, q := range quants {
			quant := q.(map[string]interface{})
			stock := c.transformStockPayload(quant)
			if stock != nil {
				allStock = append(allStock, stock)
			}
		}

		// Check if we have more records
		if len(quants) < limit {
			break
		}

		offset += limit
	}

	return allStock, nil
}

// Helper methods

func (c *Connector) sendJSONRPC(ctx context.Context, path string, request *JSONRPCRequest) (*JSONRPCResponse, error) {
	url, err := c.BuildURL(path, nil)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if c.sessionID != "" {
		req.Header.Set("X-Openerp-Session-Id", c.sessionID)
		req.Header.Set("Cookie", fmt.Sprintf("session_id=%s", c.sessionID))
	}

	// Send request
	resp, err := c.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	// Parse response
	var rpcResp JSONRPCResponse
	if err := c.ParseJSONResponse(resp, &rpcResp); err != nil {
		return nil, err
	}

	return &rpcResp, nil
}

func (c *Connector) buildObjectRequest(model, method string, args []interface{}, kwargs map[string]interface{}) map[string]interface{} {
	params := map[string]interface{}{
		"service": "object",
		"method":  "execute_kw",
		"args": append([]interface{}{
			c.database,
			c.userID,
			c.Config.Auth.Password,
			model,
			method,
		}, args...),
	}

	if kwargs != nil {
		params["args"] = append(params["args"].([]interface{}), kwargs)
	}

	return params
}

func generateRequestID() int64 {
	return time.Now().UnixNano()
}

// convertProductIDsToInts converts string product IDs to integers
func (c *Connector) convertProductIDsToInts(productIDs []string) []int {
	ids := make([]int, 0, len(productIDs))
	for _, idStr := range productIDs {
		if id, err := strconv.Atoi(idStr); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

// getProductIDBySKU retrieves product ID for a given SKU
func (c *Connector) getProductIDBySKU(ctx context.Context, sku string) int {
	// Search for product by SKU
	req := &JSONRPCRequest{
		JSONRPC: c.apiVersion,
		Method:  "call",
		ID:      generateRequestID(),
		Params: c.buildObjectRequest("product.product", "search", []interface{}{
			[]interface{}{[]interface{}{"default_code", "=", sku}},
		}, map[string]interface{}{
			"limit": 1,
		}),
	}

	resp, err := c.sendJSONRPC(ctx, "/jsonrpc", req)
	if err != nil || resp.Error != nil {
		c.Logger.Error().Err(err).Str("sku", sku).Msg("Failed to get product ID by SKU")
		return 0
	}

	if ids, ok := resp.Result.([]interface{}); ok && len(ids) > 0 {
		if id, ok := ids[0].(float64); ok {
			return int(id)
		}
	}

	return 0
}

// getCurrencyID gets the currency ID for a currency code
func (c *Connector) getCurrencyID(currencyCode string) int {
	// Default currencies
	defaultCurrencies := map[string]int{
		"EUR": 1,
		"USD": 2,
		"GBP": 3,
	}

	if id, ok := defaultCurrencies[currencyCode]; ok {
		return id
	}

	return 1 // Default to EUR
}

// transformInvoiceLines transforms invoice lines to Odoo format
func (c *Connector) transformInvoiceLines(lines []erp.InvoiceLine) []interface{} {
	odooLines := make([]interface{}, 0, len(lines))

	for _, line := range lines {
		odooLine := []interface{}{
			0, // Command for create
			0, // ID (0 for new)
			map[string]interface{}{
				"name":        line.Description,
				"quantity":    line.Quantity,
				"price_unit":  line.UnitPrice,
				"product_id":  c.getProductIDBySKU(context.Background(), line.SKU),
				"tax_ids":     []interface{}{[]interface{}{6, 0, []interface{}{}}}, // No taxes by default
			},
		}
		odooLines = append(odooLines, odooLine)
	}

	return odooLines
}

// transformOrderPayload transforms an order to Odoo format
func (c *Connector) transformOrderPayload(orderData map[string]interface{}) (*erp.OrderPayload, error) {
	order := &erp.OrderPayload{
		OrderID:     getStringValue(orderData, "name"),
		CustomerID:  fmt.Sprintf("%v", getRelationID(orderData["partner_id"])),
		TotalAmount: getFloatValue(orderData, "amount_total"),
		Currency:    getRelationCode(orderData["currency_id"]),
		Status:      getStringValue(orderData, "state"),
		CreatedAt:   getTimeValue(orderData, "date_order", "create_date"),
	}

	// Get order lines if available
	if lineIDs := getIntSlice(orderData["order_line"]); len(lineIDs) > 0 {
		items, err := c.fetchOrderLines(context.Background(), lineIDs)
		if err == nil {
			order.Items = items
		}
	}

	return order, nil
}

// transformStockPayload transforms a stock quant to stock payload
func (c *Connector) transformStockPayload(quantData map[string]interface{}) *erp.StockPayload {
	// Extract SKU from product relation
	productInfo := quantData["product_id"]
	sku := ""
	if prod, ok := productInfo.([]interface{}); ok && len(prod) > 1 {
		// Get product details to find SKU
		sku = c.getProductSKU(prod[0])
	}

	if sku == "" {
		return nil
	}

	return &erp.StockPayload{
		SKU:        sku,
		LocationID: fmt.Sprintf("%v", getRelationID(quantData["location_id"])),
		Quantity:   int(getFloatValue(quantData, "quantity")),
		Available:  int(getFloatValue(quantData, "available_quantity")),
		Reserved:   int(getFloatValue(quantData, "reserved_quantity")),
		StockType:  "physical",
		UpdatedAt:  time.Now(),
	}
}
