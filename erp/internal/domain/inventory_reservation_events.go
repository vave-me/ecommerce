package domain

import (
	"time"
)

// Inventory reservation domain event types
const (
	ReservationCreatedEvent     = "erp.ReservationCreated"
	ReservationReleasedEvent    = "erp.ReservationReleased"
	ReservationFulfilledEvent   = "erp.ReservationFulfilled"
	ReservationTransferredEvent = "erp.ReservationTransferred"
	ReservationExpiredEvent     = "erp.ReservationExpired"
)

// ReservationCreated event is raised when an inventory reservation is created
type ReservationCreated struct {
	OrderID     string
	SKU         string
	WarehouseID string
	Quantity    int
	ExpiresAt   *time.Time
	CreatedAt   time.Time
	ConnectorID string
}

// Key implements event key
func (ReservationCreated) Key() string { return ReservationCreatedEvent }

// ReservationReleased event is raised when a reservation is released
type ReservationReleased struct {
	ReleasedAt time.Time
	Reason     string
}

// Key implements event key
func (ReservationReleased) Key() string { return ReservationReleasedEvent }

// ReservationFulfilled event is raised when a reservation is fulfilled
type ReservationFulfilled struct {
	FulfilledAt time.Time
}

// Key implements event key
func (ReservationFulfilled) Key() string { return ReservationFulfilledEvent }

// ReservationTransferred event is raised when a reservation is transferred to another warehouse
type ReservationTransferred struct {
	FromWarehouseID string
	ToWarehouseID   string
	TransferredAt   time.Time
}

// Key implements event key
func (ReservationTransferred) Key() string { return ReservationTransferredEvent }

// ReservationExpired event is raised when a reservation expires
type ReservationExpired struct {
	ExpiredAt time.Time
}

// Key implements event key
func (ReservationExpired) Key() string { return ReservationExpiredEvent }