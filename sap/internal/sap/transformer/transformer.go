package transformer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"middleman/sap/internal/sap"
)

// EventType represents canonical event types
type EventType string

const (
	EventTypeProductMasterUpdated EventType = "ProductMasterUpdated"
	EventTypeStockLevelUpdated    EventType = "StockLevelUpdated"
	EventTypePriceUpdated         EventType = "PriceUpdated"
	EventTypeProductCreated       EventType = "ProductCreated"
	EventTypeProductDeleted       EventType = "ProductDeleted"
)

// CanonicalEvent represents a standardized event format
type CanonicalEvent struct {
	EventID        uuid.UUID   `json:"eventId"`
	EventType      EventType   `json:"eventType"`
	EventTimestamp time.Time   `json:"eventTimestamp"`
	Source         string      `json:"source"`
	CorrelationID  string      `json:"correlationId,omitempty"`
	Payload        interface{} `json:"payload"`
}

// ProductMasterPayload represents product master data
type ProductMasterPayload struct {
	SKU         string                 `json:"sku"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Category    string                 `json:"category"`
	Weight      float64                `json:"weight"`
	Dimensions  Dimensions             `json:"dimensions"`
	Attributes  map[string]interface{} `json:"attributes"`
}

// StockLevelPayload represents stock level data
type StockLevelPayload struct {
	SKU          string    `json:"sku"`
	WarehouseID  string    `json:"warehouseId"`
	Quantity     int       `json:"quantity"`
	StockType    string    `json:"stockType"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// PricePayload represents pricing data
type PricePayload struct {
	SKU         string     `json:"sku"`
	PriceListID string     `json:"priceListId"`
	Currency    string     `json:"currency"`
	Price       float64    `json:"price"`
	ValidFrom   time.Time  `json:"validFrom"`
	ValidTo     *time.Time `json:"validTo,omitempty"`
}

// Dimensions represents product dimensions
type Dimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Unit   string  `json:"unit"`
}

// SAPToCanonicalTransformer transforms SAP events to canonical format
type SAPToCanonicalTransformer struct {
	sourceSystem string
}

// NewSAPToCanonicalTransformer creates a new transformer
func NewSAPToCanonicalTransformer(sourceSystem string) *SAPToCanonicalTransformer {
	return &SAPToCanonicalTransformer{
		sourceSystem: sourceSystem,
	}
}

// TransformSAPEvent transforms a generic SAP event to canonical format
func (t *SAPToCanonicalTransformer) TransformSAPEvent(event *sap.SAPEvent) (*CanonicalEvent, error) {
	switch event.Type {
	case sap.EventTypeProductCreated, sap.EventTypeProductUpdated:
		return t.transformProductEvent(event)
	case sap.EventTypeStockUpdated:
		return t.transformStockEvent(event)
	case sap.EventTypePriceUpdated:
		return t.transformPriceEvent(event)
	case sap.EventTypeProductDeleted:
		return t.transformDeleteEvent(event)
	default:
		return nil, fmt.Errorf("unsupported event type: %s", event.Type)
	}
}

// TransformMATMAS transforms a MATMAS IDoc to canonical event
func (t *SAPToCanonicalTransformer) TransformMATMAS(idoc *sap.MATMAS, correlationID string) (*CanonicalEvent, error) {
	material := idoc.E1MARAM
	
	weight, _ := sap.ParseSAPFloat(material.BRGEW)
	length, _ := sap.ParseSAPFloat(material.LAENG)
	width, _ := sap.ParseSAPFloat(material.BREIT)
	height, _ := sap.ParseSAPFloat(material.HOEHE)
	
	payload := ProductMasterPayload{
		SKU:         material.MATNR,
		Name:        material.MAKTX,
		Description: material.MAKTX_LONG,
		Category:    material.MATKL,
		Weight:      weight,
		Dimensions: Dimensions{
			Length: length,
			Width:  width,
			Height: height,
			Unit:   material.MEABM,
		},
		Attributes: map[string]interface{}{
			"baseUnit":     material.MEINS,
			"materialType": material.MTART,
			"division":     material.SPART,
			"weightUnit":   material.GEWEI,
			"volumeUnit":   material.VOLEH,
		},
	}
	
	return &CanonicalEvent{
		EventID:        uuid.New(),
		EventType:      EventTypeProductMasterUpdated,
		EventTimestamp: time.Now(),
		Source:         t.sourceSystem,
		CorrelationID:  correlationID,
		Payload:        payload,
	}, nil
}

// TransformINVCON transforms an INVCON IDoc to canonical event
func (t *SAPToCanonicalTransformer) TransformINVCON(idoc *sap.INVCON, correlationID string) (*CanonicalEvent, error) {
	inventory := idoc.E1INVCO
	
	quantity, _ := sap.ParseSAPInt(inventory.MENGE)
	updatedAt, _ := sap.ParseSAPDateTime(inventory.AEDAT, inventory.AETIM)
	
	payload := StockLevelPayload{
		SKU:         inventory.MATNR,
		WarehouseID: fmt.Sprintf("%s-%s", inventory.WERKS, inventory.LGORT),
		Quantity:    quantity,
		StockType:   mapStockType(inventory.SOBKZ),
		UpdatedAt:   updatedAt,
	}
	
	return &CanonicalEvent{
		EventID:        uuid.New(),
		EventType:      EventTypeStockLevelUpdated,
		EventTimestamp: time.Now(),
		Source:         t.sourceSystem,
		CorrelationID:  correlationID,
		Payload:        payload,
	}, nil
}

// TransformCOND_A transforms a COND_A IDoc to canonical event
func (t *SAPToCanonicalTransformer) TransformCOND_A(idoc *sap.COND_A, correlationID string) (*CanonicalEvent, error) {
	condition := idoc.E1KOMG
	pricing := idoc.E1KONP
	validity := idoc.E1KONH
	
	price, _ := sap.ParseSAPFloat(pricing.KBETR)
	validFrom, _ := sap.ParseSAPDate(validity.DATAB)
	
	var validTo *time.Time
	if validity.DATBI != "" && validity.DATBI != "99991231" {
		if vt, err := sap.ParseSAPDate(validity.DATBI); err == nil {
			validTo = &vt
		}
	}
	
	payload := PricePayload{
		SKU:         condition.MATNR,
		PriceListID: fmt.Sprintf("%s-%s-%s", condition.VKORG, condition.VTWEG, condition.KSCHL),
		Currency:    pricing.WAERS,
		Price:       price,
		ValidFrom:   validFrom,
		ValidTo:     validTo,
	}
	
	return &CanonicalEvent{
		EventID:        uuid.New(),
		EventType:      EventTypePriceUpdated,
		EventTimestamp: time.Now(),
		Source:         t.sourceSystem,
		CorrelationID:  correlationID,
		Payload:        payload,
	}, nil
}

// Helper functions

func (t *SAPToCanonicalTransformer) transformProductEvent(event *sap.SAPEvent) (*CanonicalEvent, error) {
	var productData sap.ProductEventData
	if err := json.Unmarshal(event.Data, &productData); err != nil {
		return nil, fmt.Errorf("unmarshaling product data: %w", err)
	}
	
	eventType := EventTypeProductCreated
	if event.Type == sap.EventTypeProductUpdated {
		eventType = EventTypeProductMasterUpdated
	}
	
	payload := ProductMasterPayload{
		SKU:         productData.SKU,
		Name:        productData.Name,
		Description: productData.Description,
		Category:    productData.Category,
		Weight:      productData.Weight,
		Dimensions: Dimensions{
			Length: productData.Dimensions.Length,
			Width:  productData.Dimensions.Width,
			Height: productData.Dimensions.Height,
			Unit:   productData.Dimensions.Unit,
		},
		Attributes: productData.Attributes,
	}
	
	return &CanonicalEvent{
		EventID:        uuid.New(),
		EventType:      eventType,
		EventTimestamp: event.Timestamp,
		Source:         t.sourceSystem,
		CorrelationID:  event.CorrelationID,
		Payload:        payload,
	}, nil
}

func (t *SAPToCanonicalTransformer) transformStockEvent(event *sap.SAPEvent) (*CanonicalEvent, error) {
	var stockData sap.StockEventData
	if err := json.Unmarshal(event.Data, &stockData); err != nil {
		return nil, fmt.Errorf("unmarshaling stock data: %w", err)
	}
	
	payload := StockLevelPayload{
		SKU:         stockData.SKU,
		WarehouseID: stockData.WarehouseID,
		Quantity:    stockData.NewQuantity,
		StockType:   "available",
		UpdatedAt:   event.Timestamp,
	}
	
	return &CanonicalEvent{
		EventID:        uuid.New(),
		EventType:      EventTypeStockLevelUpdated,
		EventTimestamp: event.Timestamp,
		Source:         t.sourceSystem,
		CorrelationID:  event.CorrelationID,
		Payload:        payload,
	}, nil
}

func (t *SAPToCanonicalTransformer) transformPriceEvent(event *sap.SAPEvent) (*CanonicalEvent, error) {
	var priceData sap.PriceEventData
	if err := json.Unmarshal(event.Data, &priceData); err != nil {
		return nil, fmt.Errorf("unmarshaling price data: %w", err)
	}
	
	payload := PricePayload{
		SKU:         priceData.SKU,
		PriceListID: priceData.PriceListID,
		Currency:    priceData.Currency,
		Price:       priceData.NewPrice,
		ValidFrom:   priceData.ValidFrom,
		ValidTo:     priceData.ValidTo,
	}
	
	return &CanonicalEvent{
		EventID:        uuid.New(),
		EventType:      EventTypePriceUpdated,
		EventTimestamp: event.Timestamp,
		Source:         t.sourceSystem,
		CorrelationID:  event.CorrelationID,
		Payload:        payload,
	}, nil
}

func (t *SAPToCanonicalTransformer) transformDeleteEvent(event *sap.SAPEvent) (*CanonicalEvent, error) {
	var productData map[string]string
	if err := json.Unmarshal(event.Data, &productData); err != nil {
		return nil, fmt.Errorf("unmarshaling delete data: %w", err)
	}
	
	return &CanonicalEvent{
		EventID:        uuid.New(),
		EventType:      EventTypeProductDeleted,
		EventTimestamp: event.Timestamp,
		Source:         t.sourceSystem,
		CorrelationID:  event.CorrelationID,
		Payload:        productData,
	}, nil
}

// mapStockType maps SAP stock type to canonical stock type
func mapStockType(sapStockType string) string {
	switch sapStockType {
	case "":
		return "unrestricted"
	case "Q":
		return "quality_inspection"
	case "S":
		return "blocked"
	case "E":
		return "reserved"
	default:
		return "other"
	}
}