package domain

import (
	"context"
	"time"
)

type CatalogShipment struct {
	ID                 string
	ProductID          string
	OrderID            string
	BasketID           string
	TrackingNumber     string
	LabelURL           string
	SenderName         string
	SenderAddress      string
	ReceiverName       string
	ReceiverAddress    string
	Weight             string
	Dimensions         string
	ServiceType        ServiceType
	Status             ShipmentStatus
	CarrierID          string
	CarrierName        string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	PickupScheduledAt  *time.Time
	DeliveredAt        *time.Time
	CancelledAt        *time.Time
}

type ShippingCatalogRepository interface {
	Find(ctx context.Context, shipmentID string) (*CatalogShipment, error)
	GetByOrderID(ctx context.Context, orderID string) ([]*CatalogShipment, error)
	GetByBasketID(ctx context.Context, basketID string) ([]*CatalogShipment, error)
	GetByProductID(ctx context.Context, productID string) ([]*CatalogShipment, error)
	GetByTrackingNumber(ctx context.Context, trackingNumber string) (*CatalogShipment, error)
	GetShipmentsByStatus(ctx context.Context, status ShipmentStatus, limit, offset int) ([]*CatalogShipment, error)
	GetShipmentsByCarrier(ctx context.Context, carrierID string, limit, offset int) ([]*CatalogShipment, error)
	GetPendingPickups(ctx context.Context, limit int) ([]*CatalogShipment, error)
	SearchShipments(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*CatalogShipment, error)
	AddShipment(ctx context.Context, shipment *CatalogShipment) error
	UpdateShipmentStatus(ctx context.Context, shipmentID string, status ShipmentStatus, updatedAt time.Time) error
	UpdateTrackingInfo(ctx context.Context, shipmentID, trackingNumber, labelURL string) error
	UpdateDeliveryInfo(ctx context.Context, shipmentID string, deliveredAt time.Time) error
	UpdatePickupInfo(ctx context.Context, shipmentID string, pickupScheduledAt time.Time) error
	CancelShipment(ctx context.Context, shipmentID string, cancelledAt time.Time) error
}