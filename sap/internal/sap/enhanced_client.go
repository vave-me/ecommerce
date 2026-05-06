package sap

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	_ "github.com/rs/zerolog/log"
)

// EnhancedSAPClient combines all SAP integration capabilities
type EnhancedSAPClient struct {
	*SAPClient
	hanaClient     *HANAClient
	securityClient *SecurityClient
	config         *Config
}

// Config holds comprehensive SAP configuration
type Config struct {
	// API Configuration
	BaseURL       string
	APIKey        string
	WebhookSecret string

	// OAuth2 Configuration
	ClientID     string
	ClientSecret string
	TokenURL     string

	// HANA Database Configuration
	HANAHost     string
	HANAPort     string
	HANAUser     string
	HANAPassword string
	HANAUseTLS   bool

	// Security Configuration
	IASInstanceName string
	Issuer          string
	Audience        []string

	// Integration Configuration
	UseDirectHANA  bool // Use direct HANA queries instead of APIs
	EnableSecurity bool // Enable IAS security validation
}

// NewEnhancedSAPClient creates a new enhanced SAP client with all capabilities
func NewEnhancedSAPClient(config *Config) (*EnhancedSAPClient, error) {
	// Create base SAP client
	baseClient := NewSAPClient(
		config.BaseURL,
		config.APIKey,
		config.WebhookSecret,
		config.ClientID,
		config.ClientSecret,
	)

	client := &EnhancedSAPClient{
		SAPClient: baseClient,
		config:    config,
	}

	// Initialize HANA client if configured
	if config.UseDirectHANA && config.HANAHost != "" {
		hanaClient, err := NewHANAClient(
			config.HANAHost,
			config.HANAPort,
			config.HANAUser,
			config.HANAPassword,
			config.HANAUseTLS,
		)
		if err != nil {
			return nil, fmt.Errorf("creating HANA client: %w", err)
		}
		client.hanaClient = hanaClient
	}

	// Initialize security client if configured
	if config.EnableSecurity {
		securityConfig := &SecurityConfig{
			IASInstanceName: config.IASInstanceName,
			ClientID:        config.ClientID,
			ClientSecret:    config.ClientSecret,
			TokenURL:        config.TokenURL,
			Issuer:          config.Issuer,
			Audience:        config.Audience,
		}

		securityClient, err := NewSecurityClient(securityConfig)
		if err != nil {
			return nil, fmt.Errorf("creating security client: %w", err)
		}
		client.securityClient = securityClient

		// Update HTTP client with authentication
		authClient, err := securityClient.CreateAuthenticatedHTTPClient(context.Background())
		if err != nil {
			return nil, fmt.Errorf("creating authenticated HTTP client: %w", err)
		}
		client.httpClient = authClient
	}

	return client, nil
}

// Close closes all client connections
func (c *EnhancedSAPClient) Close() error {
	if c.hanaClient != nil {
		if err := c.hanaClient.Close(); err != nil {
			return fmt.Errorf("closing HANA client: %w", err)
		}
	}
	return nil
}

// GetProductChangesEnhanced retrieves product changes using the most efficient method
func (c *EnhancedSAPClient) GetProductChangesEnhanced(ctx context.Context, since time.Time) ([]*ProductChange, error) {
	// If direct HANA access is available and enabled, use it
	if c.config.UseDirectHANA && c.hanaClient != nil {
		return c.getProductChangesFromHANA(ctx, since)
	}

	// Otherwise, fall back to API
	return c.GetProductChanges(ctx, since)
}

// getProductChangesFromHANA retrieves product changes directly from HANA
func (c *EnhancedSAPClient) getProductChangesFromHANA(ctx context.Context, since time.Time) ([]*ProductChange, error) {
	// Query for changed products
	query := `
		SELECT 
			MATNR,
			MAKTX,
			MEINS,
			MTART,
			MATKL,
			BRGEW,
			NTGEW,
			GEWEI,
			LAENG,
			BREIT,
			HOEHE,
			MEABM,
			AEDAT,
			AETIM
		FROM MARA
		WHERE AEDAT >= ? OR (AEDAT = ? AND AETIM >= ?)
		ORDER BY AEDAT, AETIM`

	dateStr := since.Format("20060102")
	timeStr := since.Format("150405")

	rows, err := c.hanaClient.db.QueryContext(ctx, query, dateStr, dateStr, timeStr)
	if err != nil {
		return nil, fmt.Errorf("querying changed products: %w", err)
	}
	defer rows.Close()

	var changes []*ProductChange
	for rows.Next() {
		var data ProductMasterData
		err := rows.Scan(
			&data.MaterialNumber,
			&data.Description,
			&data.BaseUnit,
			&data.MaterialType,
			&data.MaterialGroup,
			&data.GrossWeight,
			&data.NetWeight,
			&data.WeightUnit,
			&data.Length,
			&data.Width,
			&data.Height,
			&data.DimensionUnit,
			&data.ChangedDate,
			&data.ChangedTime,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning product data: %w", err)
		}

		// Convert to ProductChange
		change := &ProductChange{
			ProductID:  data.MaterialNumber,
			SKU:        data.MaterialNumber,
			Name:       data.Description,
			Category:   data.MaterialGroup,
			Weight:     data.GrossWeight.Float64,
			ChangeType: "UPDATE",
		}

		if data.Length.Valid && data.Width.Valid && data.Height.Valid {
			change.Dimensions = Dimensions{
				Length: data.Length.Float64,
				Width:  data.Width.Float64,
				Height: data.Height.Float64,
				Unit:   data.DimensionUnit,
			}
		}

		// Parse change timestamp
		if changedAt, err := ParseSAPDateTime(data.ChangedDate, data.ChangedTime); err == nil {
			change.ChangedAt = changedAt
		}

		changes = append(changes, change)
	}

	return changes, nil
}

// ValidateWebhookWithSecurity validates webhook with enhanced security
func (c *EnhancedSAPClient) ValidateWebhookWithSecurity(r *http.Request) error {
	// First, validate using security client if enabled
	if c.securityClient != nil {
		if err := c.securityClient.AuthenticateWebhook(r); err != nil {
			return fmt.Errorf("security authentication failed: %w", err)
		}
	}

	// Then validate signature
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("reading request body: %w", err)
	}
	r.Body = io.NopCloser(bytes.NewReader(body))

	signature := r.Header.Get("X-SAP-Signature")
	if signature == "" {
		signature = r.Header.Get("SAP-Signature")
	}

	return c.ValidateSignature(body, signature)
}

// GetStockLevelsEnhanced retrieves stock levels using the most efficient method
func (c *EnhancedSAPClient) GetStockLevelsEnhanced(ctx context.Context, productIDs []string) ([]*StockLevel, error) {
	// If direct HANA access is available, use it
	if c.config.UseDirectHANA && c.hanaClient != nil {
		stockData, err := c.hanaClient.GetStockLevels(ctx, productIDs)
		if err != nil {
			return nil, err
		}

		// Convert HANA stock data to API format
		var levels []*StockLevel
		for _, data := range stockData {
			level := &StockLevel{
				ProductID:    data.MaterialNumber,
				SKU:          data.MaterialNumber,
				WarehouseID:  fmt.Sprintf("%s-%s", data.Plant, data.StorageLocation),
				Quantity:     int(data.UnrestrictedStock),
				AvailableQty: int(data.UnrestrictedStock),
				UpdatedAt:    time.Now(),
			}
			levels = append(levels, level)
		}

		return levels, nil
	}

	// Otherwise, use API
	return c.GetStockLevels(ctx, productIDs)
}

// Middleware returns authentication middleware for HTTP handlers
func (c *EnhancedSAPClient) Middleware() func(http.Handler) http.Handler {
	if c.securityClient != nil {
		return c.securityClient.Middleware()
	}

	// Return a no-op middleware if security is not enabled
	return func(next http.Handler) http.Handler {
		return next
	}
}
