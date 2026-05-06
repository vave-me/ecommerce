package sap

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/SAP/go-hdb/driver"
	"github.com/rs/zerolog/log"
)

// HANAClient handles direct connections to SAP HANA database
type HANAClient struct {
	db *sql.DB
}

// NewHANAClient creates a new HANA database client
func NewHANAClient(host, port, user, password string, useTLS bool) (*HANAClient, error) {
	var dsn string
	
	if useTLS {
		// For HANA Cloud with TLS
		dsn = fmt.Sprintf("hdb://%s:%s@%s:%s?TLSServerName=%s", user, password, host, port, host)
	} else {
		// For on-premise HANA without TLS
		dsn = fmt.Sprintf("hdb://%s:%s@%s:%s", user, password, host, port)
	}
	
	db, err := sql.Open("hdb", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening HANA connection: %w", err)
	}
	
	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	
	// Test connection
	if err := db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("pinging HANA database: %w", err)
	}
	
	return &HANAClient{db: db}, nil
}

// Close closes the database connection
func (c *HANAClient) Close() error {
	return c.db.Close()
}

// GetProductMasterData retrieves product master data directly from HANA
func (c *HANAClient) GetProductMasterData(ctx context.Context, productIDs []string) ([]*ProductMasterData, error) {
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
		WHERE MATNR IN (` + generatePlaceholders(len(productIDs)) + `)`
	
	args := make([]interface{}, len(productIDs))
	for i, id := range productIDs {
		args[i] = id
	}
	
	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying product master data: %w", err)
	}
	defer rows.Close()
	
	var products []*ProductMasterData
	for rows.Next() {
		var p ProductMasterData
		err := rows.Scan(
			&p.MaterialNumber,
			&p.Description,
			&p.BaseUnit,
			&p.MaterialType,
			&p.MaterialGroup,
			&p.GrossWeight,
			&p.NetWeight,
			&p.WeightUnit,
			&p.Length,
			&p.Width,
			&p.Height,
			&p.DimensionUnit,
			&p.ChangedDate,
			&p.ChangedTime,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning product data: %w", err)
		}
		products = append(products, &p)
	}
	
	return products, nil
}

// GetStockLevels retrieves current stock levels from HANA
func (c *HANAClient) GetStockLevels(ctx context.Context, productIDs []string) ([]*StockData, error) {
	query := `
		SELECT 
			MATNR,
			WERKS,
			LGORT,
			LABST,
			UMLME,
			EINME,
			SPEME,
			RETME
		FROM MARD
		WHERE MATNR IN (` + generatePlaceholders(len(productIDs)) + `)`
	
	args := make([]interface{}, len(productIDs))
	for i, id := range productIDs {
		args[i] = id
	}
	
	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying stock levels: %w", err)
	}
	defer rows.Close()
	
	var stocks []*StockData
	for rows.Next() {
		var s StockData
		err := rows.Scan(
			&s.MaterialNumber,
			&s.Plant,
			&s.StorageLocation,
			&s.UnrestrictedStock,
			&s.StockInTransfer,
			&s.StockInQuality,
			&s.BlockedStock,
			&s.ReturnsStock,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning stock data: %w", err)
		}
		stocks = append(stocks, &s)
	}
	
	return stocks, nil
}

// GetPricingConditions retrieves pricing conditions from HANA
func (c *HANAClient) GetPricingConditions(ctx context.Context, productIDs []string, conditionType string) ([]*PricingData, error) {
	query := `
		SELECT 
			a.MATNR,
			a.KSCHL,
			a.DATAB,
			a.DATBI,
			b.KBETR,
			b.KONWA,
			b.KPEIN,
			b.KMEIN
		FROM KONH a
		INNER JOIN KONP b ON a.KNUMH = b.KNUMH
		WHERE a.MATNR IN (` + generatePlaceholders(len(productIDs)) + `)
		AND a.KSCHL = ?
		AND a.DATAB <= CURRENT_DATE
		AND a.DATBI >= CURRENT_DATE`
	
	args := make([]interface{}, len(productIDs)+1)
	for i, id := range productIDs {
		args[i] = id
	}
	args[len(productIDs)] = conditionType
	
	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying pricing conditions: %w", err)
	}
	defer rows.Close()
	
	var prices []*PricingData
	for rows.Next() {
		var p PricingData
		err := rows.Scan(
			&p.MaterialNumber,
			&p.ConditionType,
			&p.ValidFrom,
			&p.ValidTo,
			&p.Rate,
			&p.Currency,
			&p.PricingUnit,
			&p.UnitOfMeasure,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning pricing data: %w", err)
		}
		prices = append(prices, &p)
	}
	
	return prices, nil
}

// ExecuteStoredProcedure executes a stored procedure in HANA
func (c *HANAClient) ExecuteStoredProcedure(ctx context.Context, procedure string, args ...interface{}) (*sql.Rows, error) {
	stmt := fmt.Sprintf("CALL %s(%s)", procedure, generatePlaceholders(len(args)))
	return c.db.QueryContext(ctx, stmt, args...)
}

// BeginTransaction starts a new transaction
func (c *HANAClient) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	return c.db.BeginTx(ctx, nil)
}

// Helper function to generate SQL placeholders
func generatePlaceholders(n int) string {
	if n <= 0 {
		return ""
	}
	placeholders := "?"
	for i := 1; i < n; i++ {
		placeholders += ", ?"
	}
	return placeholders
}

// Data structures for HANA queries

type ProductMasterData struct {
	MaterialNumber  string
	Description     string
	BaseUnit        string
	MaterialType    string
	MaterialGroup   string
	GrossWeight     sql.NullFloat64
	NetWeight       sql.NullFloat64
	WeightUnit      string
	Length          sql.NullFloat64
	Width           sql.NullFloat64
	Height          sql.NullFloat64
	DimensionUnit   string
	ChangedDate     string
	ChangedTime     string
}

type StockData struct {
	MaterialNumber    string
	Plant             string
	StorageLocation   string
	UnrestrictedStock float64
	StockInTransfer   float64
	StockInQuality    float64
	BlockedStock      float64
	ReturnsStock      float64
}

type PricingData struct {
	MaterialNumber string
	ConditionType  string
	ValidFrom      time.Time
	ValidTo        time.Time
	Rate           float64
	Currency       string
	PricingUnit    float64
	UnitOfMeasure  string
}