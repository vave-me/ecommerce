package odoo

import (
	"context"
	"fmt"
	"middleman/internal/erp"
	"strconv"
	"time"
)

// JSONRPCRequest represents a JSON-RPC request to Odoo
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      interface{} `json:"id"`
}

// JSONRPCResponse represents a JSON-RPC response from Odoo
type JSONRPCResponse struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      interface{}            `json:"id"`
	Result  interface{}            `json:"result"`
	Error   map[string]interface{} `json:"error,omitempty"`
}

// OdooWebhook represents an incoming webhook from Odoo
type OdooWebhook struct {
	ID        string                 `json:"id"`
	Event     string                 `json:"event"`
	Model     string                 `json:"model"`
	RecordID  int                    `json:"record_id"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// OdooProduct represents a product in Odoo
type OdooProduct struct {
	ID              int           `json:"id"`
	Name            string        `json:"name"`
	DefaultCode     string        `json:"default_code"`
	Description     string        `json:"description"`
	DescriptionSale string        `json:"description_sale"`
	CategoryID      []interface{} `json:"categ_id"` // [id, name]
	ListPrice       float64       `json:"list_price"`
	StandardPrice   float64       `json:"standard_price"`
	Weight          float64       `json:"weight"`
	Volume          float64       `json:"volume"`
	Barcode         string        `json:"barcode"`
	Type            string        `json:"type"`
	Active          bool          `json:"active"`
	UOMID           []interface{} `json:"uom_id"`
	TaxesID         []int         `json:"taxes_id"`
	CreateDate      time.Time     `json:"create_date"`
	WriteDate       time.Time     `json:"write_date"`
}

// OdooStock represents stock/inventory in Odoo
type OdooStock struct {
	ID                int           `json:"id"`
	ProductID         []interface{} `json:"product_id"`  // [id, name]
	LocationID        []interface{} `json:"location_id"` // [id, name]
	Quantity          float64       `json:"quantity"`
	ReservedQuantity  float64       `json:"reserved_quantity"`
	AvailableQuantity float64       `json:"available_quantity"`
	InDate            *time.Time    `json:"in_date"`
	WriteDate         time.Time     `json:"write_date"`
}

// OdooPriceListItem represents a price list item in Odoo
type OdooPriceListItem struct {
	ID            int           `json:"id"`
	PricelistID   []interface{} `json:"pricelist_id"` // [id, name]
	ProductID     []interface{} `json:"product_id"`   // [id, name]
	ProductTmplID []interface{} `json:"product_tmpl_id"`
	MinQuantity   float64       `json:"min_quantity"`
	FixedPrice    float64       `json:"fixed_price"`
	Price         float64       `json:"price"`
	DateStart     *time.Time    `json:"date_start"`
	DateEnd       *time.Time    `json:"date_end"`
	CurrencyID    []interface{} `json:"currency_id"` // [id, code]
}

// OdooOrder represents a sales order in Odoo
type OdooOrder struct {
	ID                int           `json:"id"`
	Name              string        `json:"name"`
	PartnerID         []interface{} `json:"partner_id"`          // [id, name]
	PartnerInvoiceID  []interface{} `json:"partner_invoice_id"`  // [id, name]
	PartnerShippingID []interface{} `json:"partner_shipping_id"` // [id, name]
	DateOrder         time.Time     `json:"date_order"`
	CreateDate        time.Time     `json:"create_date"`
	State             string        `json:"state"`
	AmountTotal       float64       `json:"amount_total"`
	AmountUntaxed     float64       `json:"amount_untaxed"`
	CurrencyID        []interface{} `json:"currency_id"`
	OrderLine         []int         `json:"order_line"` // Line IDs
}

// OdooOrderLine represents a sales order line
type OdooOrderLine struct {
	ID            int           `json:"id"`
	OrderID       []interface{} `json:"order_id"`
	ProductID     []interface{} `json:"product_id"`
	Name          string        `json:"name"`
	ProductUOMQty float64       `json:"product_uom_qty"`
	PriceUnit     float64       `json:"price_unit"`
	PriceSubtotal float64       `json:"price_subtotal"`
	PriceTotal    float64       `json:"price_total"`
	Discount      float64       `json:"discount"`
	TaxID         []int         `json:"tax_id"`
}

// OdooPartner represents a customer/partner in Odoo
type OdooPartner struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	Email        string        `json:"email"`
	Phone        string        `json:"phone"`
	Mobile       string        `json:"mobile"`
	Website      string        `json:"website"`
	Street       string        `json:"street"`
	Street2      string        `json:"street2"`
	City         string        `json:"city"`
	StateID      []interface{} `json:"state_id"`   // [id, name]
	CountryID    []interface{} `json:"country_id"` // [id, name]
	Zip          string        `json:"zip"`
	VAT          string        `json:"vat"`
	Ref          string        `json:"ref"`
	Lang         string        `json:"lang"`
	IsCompany    bool          `json:"is_company"`
	CustomerRank int           `json:"customer_rank"`
	SupplierRank int           `json:"supplier_rank"`
	Active       bool          `json:"active"`
	CreateDate   time.Time     `json:"create_date"`
	WriteDate    time.Time     `json:"write_date"`
}

// Helper functions for parsing Odoo responses

func (c *Connector) parseProductsResponse(result interface{}) ([]*erp.ProductPayload, error) {
	records, ok := result.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected products response format")
	}

	products := make([]*erp.ProductPayload, 0, len(records))
	for _, record := range records {
		data, ok := record.(map[string]interface{})
		if !ok {
			continue
		}

		product := &erp.ProductPayload{
			SKU:         getStringValue(data, "default_code"),
			Name:        getStringValue(data, "name"),
			Description: getStringValue(data, "description_sale", "description"),
			Category:    getRelationName(data["categ_id"]),
			Weight:      getFloatValue(data, "weight"),
			Attributes: map[string]interface{}{
				"odoo_id":        data["id"],
				"barcode":        getStringValue(data, "barcode"),
				"type":           getStringValue(data, "type"),
				"list_price":     getFloatValue(data, "list_price"),
				"standard_price": getFloatValue(data, "standard_price"),
				"volume":         getFloatValue(data, "volume"),
				"active":         getBoolValue(data, "active"),
			},
		}

		// Add UOM if present
		if uom := getRelationName(data["uom_id"]); uom != "" {
			product.Attributes["unit_of_measure"] = uom
		}

		products = append(products, product)
	}

	return products, nil
}

func (c *Connector) parseStockResponse(result interface{}) ([]*erp.StockPayload, error) {
	records, ok := result.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected stock response format")
	}

	// Aggregate stock by SKU and location
	stockMap := make(map[string]*erp.StockPayload)

	for _, record := range records {
		data, ok := record.(map[string]interface{})
		if !ok {
			continue
		}

		// Extract SKU from product relation
		productInfo := data["product_id"]
		sku := ""
		if prod, ok := productInfo.([]interface{}); ok && len(prod) > 1 {
			// Get product details to find SKU
			sku = c.getProductSKU(prod[0])
		}

		if sku == "" {
			continue
		}

		locationID := fmt.Sprintf("%v", getRelationID(data["location_id"]))
		key := fmt.Sprintf("%s-%s", sku, locationID)

		stock := &erp.StockPayload{
			SKU:        sku,
			LocationID: locationID,
			Quantity:   int(getFloatValue(data, "quantity")),
			Available:  int(getFloatValue(data, "available_quantity")),
			Reserved:   int(getFloatValue(data, "reserved_quantity")),
			StockType:  "physical",
			UpdatedAt:  time.Now(),
		}

		// Aggregate quantities if multiple records for same SKU/location
		if existing, ok := stockMap[key]; ok {
			existing.Quantity += stock.Quantity
			existing.Available += stock.Available
			existing.Reserved += stock.Reserved
		} else {
			stockMap[key] = stock
		}
	}

	// Convert map to slice
	stocks := make([]*erp.StockPayload, 0, len(stockMap))
	for _, stock := range stockMap {
		stocks = append(stocks, stock)
	}

	return stocks, nil
}

func (c *Connector) parsePriceResponse(result interface{}, productIDs []string, priceListID string) ([]*erp.PricePayload, error) {
	priceMap, ok := result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected price response format")
	}

	prices := make([]*erp.PricePayload, 0, len(priceMap))

	for productIDStr, price := range priceMap {
		priceValue, ok := price.(float64)
		if !ok {
			continue
		}

		// Map back to original product ID
		var originalID string
		for _, pid := range productIDs {
			if pid == productIDStr {
				originalID = pid
				break
			}
		}
		
		if originalID == "" {
			continue
		}
		
		// Get SKU for product
		sku := c.getSKUByProductID(parseIntOrZero(productIDStr))

		pricePayload := &erp.PricePayload{
			SKU:         sku,
			PriceListID: priceListID,
			Currency:    "EUR", // Default, should be fetched from pricelist
			Price:       priceValue,
			ValidFrom:   time.Now(),
		}

		prices = append(prices, pricePayload)
	}

	return prices, nil
}

func (c *Connector) parseOrdersResponse(ctx context.Context, result interface{}) ([]*erp.OrderPayload, error) {
	records, ok := result.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected orders response format")
	}

	orders := make([]*erp.OrderPayload, 0, len(records))

	for _, record := range records {
		data, ok := record.(map[string]interface{})
		if !ok {
			continue
		}

		// Get order line details
		lineIDs := getIntSlice(data["order_line"])
		orderItems, err := c.fetchOrderLines(ctx, lineIDs)
		if err != nil {
			c.Logger.Error().Err(err).Msg("Failed to fetch order lines")
			continue
		}

		order := &erp.OrderPayload{
			OrderID:     getStringValue(data, "name"),
			CustomerID:  fmt.Sprintf("%v", getRelationID(data["partner_id"])),
			Items:       orderItems,
			TotalAmount: getFloatValue(data, "amount_total"),
			Currency:    getRelationCode(data["currency_id"]),
			Status:      getStringValue(data, "state"),
			CreatedAt:   getTimeValue(data, "date_order", "create_date"),
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (c *Connector) parseCustomersResponse(result interface{}) ([]*erp.CustomerPayload, error) {
	records, ok := result.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected customers response format")
	}

	customers := make([]*erp.CustomerPayload, 0, len(records))

	for _, record := range records {
		data, ok := record.(map[string]interface{})
		if !ok {
			continue
		}

		customer := &erp.CustomerPayload{
			CustomerID: fmt.Sprintf("%v", data["id"]),
			Email:      getStringValue(data, "email"),
			Name:       getStringValue(data, "name"),
			Phone:      getStringValue(data, "phone", "mobile"),
			Address: &erp.Address{
				Street:     getStringValue(data, "street"),
				City:       getStringValue(data, "city"),
				State:      getRelationName(data["state_id"]),
				PostalCode: getStringValue(data, "zip"),
				Country:    getRelationName(data["country_id"]),
			},
			Attributes: map[string]interface{}{
				"vat":           getStringValue(data, "vat"),
				"ref":           getStringValue(data, "ref"),
				"lang":          getStringValue(data, "lang"),
				"is_company":    getBoolValue(data, "is_company"),
				"customer_rank": getIntValue(data, "customer_rank"),
			},
		}

		// Add street2 if present
		if street2 := getStringValue(data, "street2"); street2 != "" {
			customer.Address.Street += ", " + street2
		}

		customers = append(customers, customer)
	}

	return customers, nil
}


// mapOdooEventType maps Odoo event types to canonical event types
func mapOdooEventType(odooEvent string) (erp.EventType, error) {
	eventMap := map[string]erp.EventType{
		"product.template.write":        erp.EventTypeProductMasterUpdated,
		"product.product.write":         erp.EventTypeProductMasterUpdated,
		"product.template.create":       erp.EventTypeProductCreated,
		"product.product.create":        erp.EventTypeProductCreated,
		"product.template.unlink":       erp.EventTypeProductDeleted,
		"product.product.unlink":        erp.EventTypeProductDeleted,
		"stock.quant.write":             erp.EventTypeStockLevelUpdated,
		"stock.move.write":              erp.EventTypeStockLevelUpdated,
		"product.pricelist.item.write":  erp.EventTypePriceUpdated,
		"product.pricelist.item.create": erp.EventTypePriceUpdated,
		"sale.order.create":             erp.EventTypeOrderCreated,
		"sale.order.write":              erp.EventTypeOrderUpdated,
		"res.partner.create":            erp.EventTypeCustomerCreated,
		"res.partner.write":             erp.EventTypeCustomerUpdated,
	}

	if eventType, ok := eventMap[odooEvent]; ok {
		return eventType, nil
	}

	return "", fmt.Errorf("unknown Odoo event type: %s", odooEvent)
}

// getRelationID extracts the ID from a relation field
func getRelationID(relation interface{}) interface{} {
	if rel, ok := relation.([]interface{}); ok && len(rel) > 0 {
		return rel[0]
	}
	return 0
}

// getRelationName extracts the name from a relation field
func getRelationName(relation interface{}) string {
	if rel, ok := relation.([]interface{}); ok && len(rel) > 1 {
		if name, ok := rel[1].(string); ok {
			return name
		}
	}
	return ""
}

// getRelationCode extracts the code from a relation field (e.g., currency)
func getRelationCode(relation interface{}) string {
	// For currency, the code is often the second element
	if rel, ok := relation.([]interface{}); ok && len(rel) > 1 {
		if code, ok := rel[1].(string); ok {
			// Extract currency code from name like "EUR (€)"
			if len(code) >= 3 {
				return code[:3]
			}
			return code
		}
	}
	return "EUR" // Default currency
}

// getStringValue safely extracts a string value from a map
func getStringValue(data map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if val, ok := data[key].(string); ok && val != "" {
			return val
		}
	}
	return ""
}

// getFloatValue safely extracts a float value from a map
func getFloatValue(data map[string]interface{}, key string) float64 {
	if val, ok := data[key].(float64); ok {
		return val
	}
	return 0
}

// getIntValue safely extracts an int value from a map
func getIntValue(data map[string]interface{}, key string) int {
	if val, ok := data[key].(float64); ok {
		return int(val)
	}
	return 0
}

// getBoolValue safely extracts a bool value from a map
func getBoolValue(data map[string]interface{}, key string) bool {
	if val, ok := data[key].(bool); ok {
		return val
	}
	return false
}

// getTimeValue safely extracts a time value from a map
func getTimeValue(data map[string]interface{}, keys ...string) time.Time {
	for _, key := range keys {
		if val, ok := data[key].(string); ok {
			if t, err := time.Parse("2006-01-02 15:04:05", val); err == nil {
				return t
			}
			if t, err := time.Parse("2006-01-02", val); err == nil {
				return t
			}
		}
	}
	return time.Now()
}

// getIntSlice converts an interface{} to []int
func getIntSlice(data interface{}) []int {
	if slice, ok := data.([]interface{}); ok {
		result := make([]int, 0, len(slice))
		for _, v := range slice {
			if id, ok := v.(float64); ok {
				result = append(result, int(id))
			}
		}
		return result
	}
	return []int{}
}

// parseIntOrZero safely parses a string to int
func parseIntOrZero(s string) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return 0
}

// Additional helper methods for the connector

func (c *Connector) transformWebhookPayload(webhook *OdooWebhook, eventType erp.EventType) (interface{}, error) {
	// Transform based on event type
	switch eventType {
	case erp.EventTypeProductMasterUpdated, erp.EventTypeProductCreated:
		return c.transformProductData(webhook.Data), nil
	case erp.EventTypeStockLevelUpdated:
		return c.transformStockData(webhook.Data), nil
	case erp.EventTypePriceUpdated:
		return c.transformPriceData(webhook.Data), nil
	case erp.EventTypeOrderCreated, erp.EventTypeOrderUpdated:
		return c.transformOrderData(webhook.Data), nil
	case erp.EventTypeCustomerCreated, erp.EventTypeCustomerUpdated:
		return c.transformCustomerData(webhook.Data), nil
	default:
		return webhook.Data, nil
	}
}

func (c *Connector) transformProductData(data map[string]interface{}) *erp.ProductPayload {
	return &erp.ProductPayload{
		SKU:         getStringValue(data, "default_code"),
		Name:        getStringValue(data, "name"),
		Description: getStringValue(data, "description_sale", "description"),
		Category:    getStringValue(data, "categ_id"),
		Weight:      getFloatValue(data, "weight"),
		Attributes:  data,
	}
}

func (c *Connector) transformStockData(data map[string]interface{}) *erp.StockPayload {
	return &erp.StockPayload{
		SKU:        getStringValue(data, "product_code"),
		LocationID: fmt.Sprintf("%v", data["location_id"]),
		Quantity:   int(getFloatValue(data, "quantity")),
		Available:  int(getFloatValue(data, "available_quantity")),
		Reserved:   int(getFloatValue(data, "reserved_quantity")),
		StockType:  "physical",
		UpdatedAt:  time.Now(),
	}
}

func (c *Connector) transformPriceData(data map[string]interface{}) *erp.PricePayload {
	return &erp.PricePayload{
		SKU:         getStringValue(data, "product_code"),
		PriceListID: fmt.Sprintf("%v", data["pricelist_id"]),
		Currency:    getStringValue(data, "currency"),
		Price:       getFloatValue(data, "price"),
		ValidFrom:   time.Now(),
	}
}

func (c *Connector) transformOrderData(data map[string]interface{}) *erp.OrderPayload {
	return &erp.OrderPayload{
		OrderID:     getStringValue(data, "name"),
		CustomerID:  fmt.Sprintf("%v", data["partner_id"]),
		TotalAmount: getFloatValue(data, "amount_total"),
		Currency:    getStringValue(data, "currency"),
		Status:      getStringValue(data, "state"),
		CreatedAt:   getTimeValue(data, "date_order"),
	}
}

func (c *Connector) transformCustomerData(data map[string]interface{}) *erp.CustomerPayload {
	return &erp.CustomerPayload{
		CustomerID: fmt.Sprintf("%v", data["id"]),
		Email:      getStringValue(data, "email"),
		Name:       getStringValue(data, "name"),
		Phone:      getStringValue(data, "phone", "mobile"),
		Attributes: data,
	}
}

// Stub methods that would need implementation

func (c *Connector) getProductSKU(productID interface{}) string {
	// Convert to int
	var id int
	switch v := productID.(type) {
	case float64:
		id = int(v)
	case int:
		id = v
	default:
		return ""
	}

	// Read product to get SKU
	req := &JSONRPCRequest{
		JSONRPC: c.apiVersion,
		Method:  "call",
		ID:      generateRequestID(),
		Params: c.buildObjectRequest("product.product", "read", []interface{}{
			[]int{id},
		}, map[string]interface{}{
			"fields": []string{"default_code"},
		}),
	}

	resp, err := c.sendJSONRPC(context.Background(), "/jsonrpc", req)
	if err != nil || resp.Error != nil {
		return ""
	}

	if records, ok := resp.Result.([]interface{}); ok && len(records) > 0 {
		if record, ok := records[0].(map[string]interface{}); ok {
			return getStringValue(record, "default_code")
		}
	}

	return ""
}

func (c *Connector) getSKUByProductID(productID int) string {
	return c.getProductSKU(productID)
}

func (c *Connector) getProductIDsBySKUs(ctx context.Context, skus []string) []int {
	if len(skus) == 0 {
		return []int{}
	}

	// Search for products by SKUs
	req := &JSONRPCRequest{
		JSONRPC: c.apiVersion,
		Method:  "call",
		ID:      generateRequestID(),
		Params: c.buildObjectRequest("product.product", "search", []interface{}{
			[]interface{}{[]interface{}{"default_code", "in", skus}},
		}, nil),
	}

	resp, err := c.sendJSONRPC(ctx, "/jsonrpc", req)
	if err != nil || resp.Error != nil {
		c.Logger.Error().Err(err).Msg("Failed to get product IDs by SKUs")
		return []int{}
	}

	if ids, ok := resp.Result.([]interface{}); ok {
		result := make([]int, 0, len(ids))
		for _, id := range ids {
			if idFloat, ok := id.(float64); ok {
				result = append(result, int(idFloat))
			}
		}
		return result
	}

	return []int{}
}

func (c *Connector) fetchOrderLines(ctx context.Context, lineIDs []int) ([]erp.OrderItem, error) {
	if len(lineIDs) == 0 {
		return []erp.OrderItem{}, nil
	}

	// Read order lines
	req := &JSONRPCRequest{
		JSONRPC: c.apiVersion,
		Method:  "call",
		ID:      generateRequestID(),
		Params: c.buildObjectRequest("sale.order.line", "read", []interface{}{
			lineIDs,
		}, map[string]interface{}{
			"fields": []string{
				"product_id", "product_uom_qty", "price_unit",
				"price_subtotal", "discount", "name",
			},
		}),
	}

	resp, err := c.sendJSONRPC(ctx, "/jsonrpc", req)
	if err != nil {
		return nil, fmt.Errorf("fetching order lines: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("API error: %v", resp.Error)
	}

	// Parse order lines
	items := make([]erp.OrderItem, 0, len(lineIDs))
	if lines, ok := resp.Result.([]interface{}); ok {
		for _, line := range lines {
			if lineData, ok := line.(map[string]interface{}); ok {
				// Get SKU from product
				sku := ""
				if productID := getRelationID(lineData["product_id"]); productID != nil {
					sku = c.getProductSKU(productID)
				}

				item := erp.OrderItem{
					SKU:         sku,
					Quantity:    int(getFloatValue(lineData, "product_uom_qty")),
					Price:       getFloatValue(lineData, "price_unit"),
					Total:       getFloatValue(lineData, "price_subtotal"),
					Discount:    getFloatValue(lineData, "discount"),
					Description: getStringValue(lineData, "name"),
				}
				items = append(items, item)
			}
		}
	}

	return items, nil
}

func (c *Connector) transformOrderToOdoo(ctx context.Context, order *erp.OrderPayload) (map[string]interface{}, error) {
	// Convert customer ID to int
	partnerID, err := strconv.Atoi(order.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("invalid customer ID: %w", err)
	}

	// Build order lines
	orderLines := make([]interface{}, 0, len(order.Items))
	for _, item := range order.Items {
		// Get product ID from SKU
		productID := c.getProductIDBySKU(ctx, item.SKU)
		if productID == 0 {
			c.Logger.Warn().Str("sku", item.SKU).Msg("Product not found for SKU")
			continue
		}

		// Create order line
		line := []interface{}{
			0, // Command for create
			0, // ID (0 for new)
			map[string]interface{}{
				"product_id":      productID,
				"product_uom_qty": item.Quantity,
				"price_unit":      item.Price,
			},
		}
		orderLines = append(orderLines, line)
	}

	// Build order data
	odooOrder := map[string]interface{}{
		"partner_id":     partnerID,
		"order_line":     orderLines,
		"client_order_ref": order.OrderID, // Store original order ID as reference
	}

	// Add shipping address if different
	if order.ShippingAddress != nil {
		// You might need to create/find shipping partner
		odooOrder["partner_shipping_id"] = partnerID // Simplified: same as billing
	}

	return odooOrder, nil
}

func (c *Connector) confirmOrder(ctx context.Context, orderID interface{}) error {
	// Convert order ID to int
	var id int
	switch v := orderID.(type) {
	case float64:
		id = int(v)
	case int:
		id = v
	default:
		return fmt.Errorf("invalid order ID type")
	}

	// Confirm the order
	req := &JSONRPCRequest{
		JSONRPC: c.apiVersion,
		Method:  "call",
		ID:      generateRequestID(),
		Params: c.buildObjectRequest("sale.order", "action_confirm", []interface{}{
			[]int{id},
		}, nil),
	}

	resp, err := c.sendJSONRPC(ctx, "/jsonrpc", req)
	if err != nil {
		return fmt.Errorf("confirming order: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("order confirmation error: %v", resp.Error["message"])
	}

	return nil
}

// parseOdooError extracts meaningful error information from Odoo JSON-RPC errors
func parseOdooError(errorData map[string]interface{}) error {
	if errorData == nil {
		return nil
	}

	code := errorData["code"]
	message := errorData["message"]

	// Extract detailed error data if available
	if data, ok := errorData["data"].(map[string]interface{}); ok {
		if debugMsg, ok := data["debug"].(string); ok {
			return fmt.Errorf("odoo error %v: %v (debug: %s)", code, message, debugMsg)
		}
		if exceptionMsg, ok := data["message"].(string); ok {
			return fmt.Errorf("odoo error %v: %v (exception: %s)", code, message, exceptionMsg)
		}
	}

	return fmt.Errorf("odoo error %v: %v", code, message)
}
