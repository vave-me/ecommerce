package domain

import "context"

type ShippingRepository interface {
	Load(ctx context.Context, id string) (*Shipment, error)
	Save(ctx context.Context, shipping *Shipment) error
}


type ListFilters struct {
	Limit      int
	Offset     int
	Status     string
	ProductID  string
	OrderID    string
	CarrierID  string
	StartDate  string
	EndDate    string
}

type ShipmentEvent struct {
	ID          string
	ShipmentID  string
	EventType   string
	Status      string
	Location    string
	Description string
	Timestamp   string
}
