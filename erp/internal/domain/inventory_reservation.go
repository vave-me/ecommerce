package domain

import (
	"middleman/internal/ddd"
	"middleman/internal/es"
	"time"

	"github.com/stackus/errors"
)

const InventoryReservationAggregate = "erp.InventoryReservation"

var (
	ErrReservationAlreadyExists = errors.Wrap(errors.ErrBadRequest, "reservation already exists")
	ErrInvalidQuantity          = errors.Wrap(errors.ErrBadRequest, "quantity must be greater than zero")
	ErrMissingSKU               = errors.Wrap(errors.ErrBadRequest, "SKU is required")
	ErrMissingWarehouseID       = errors.Wrap(errors.ErrBadRequest, "warehouse ID is required")
	ErrMissingOrderID           = errors.Wrap(errors.ErrBadRequest, "order ID is required")
	ErrReservationNotActive     = errors.Wrap(errors.ErrBadRequest, "reservation is not active")
)

// ReservationStatus represents the status of an inventory reservation
type ReservationStatus string

const (
	ReservationStatusActive      ReservationStatus = "active"
	ReservationStatusReleased    ReservationStatus = "released"
	ReservationStatusFulfilled   ReservationStatus = "fulfilled"
	ReservationStatusExpired     ReservationStatus = "expired"
	ReservationStatusTransferred ReservationStatus = "transferred"
)

// InventoryReservation represents an inventory reservation aggregate
type InventoryReservation struct {
	es.Aggregate
	OrderID             string
	SKU                 string
	WarehouseID         string
	Quantity            int
	Status              ReservationStatus
	ExpiresAt           *time.Time
	CreatedAt           time.Time
	ReleasedAt          *time.Time
	FulfilledAt         *time.Time
	TransferredAt       *time.Time
	TransferToWarehouse string
	ConnectorID         string
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*InventoryReservation)(nil)

func NewInventoryReservation(id string) *InventoryReservation {
	return &InventoryReservation{
		Aggregate: es.NewAggregate(id, InventoryReservationAggregate),
	}
}

// Key implements registry.Registerable
func (InventoryReservation) Key() string { return InventoryReservationAggregate }

// CreateReservation creates a new inventory reservation
func (r *InventoryReservation) CreateReservation(
	orderID string,
	sku string,
	warehouseID string,
	quantity int,
	expiresAt *time.Time,
	connectorID string,
) (ddd.Event, error) {
	if r.OrderID != "" {
		return nil, ErrReservationAlreadyExists
	}
	if orderID == "" {
		return nil, ErrMissingOrderID
	}
	if sku == "" {
		return nil, ErrMissingSKU
	}
	if warehouseID == "" {
		return nil, ErrMissingWarehouseID
	}
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	r.AddEvent(ReservationCreatedEvent, &ReservationCreated{
		OrderID:     orderID,
		SKU:         sku,
		WarehouseID: warehouseID,
		Quantity:    quantity,
		ExpiresAt:   expiresAt,
		CreatedAt:   time.Now(),
		ConnectorID: connectorID,
	})
	return ddd.NewEvent(ReservationCreatedEvent, r), nil
}

// ReleaseReservation releases the inventory reservation
func (r *InventoryReservation) ReleaseReservation(reason string) (ddd.Event, error) {
	if r.Status != ReservationStatusActive {
		return nil, ErrReservationNotActive
	}

	r.AddEvent(ReservationReleasedEvent, &ReservationReleased{
		ReleasedAt: time.Now(),
		Reason:     reason,
	})
	return ddd.NewEvent(ReservationReleasedEvent, r), nil
}

// FulfillReservation marks the reservation as fulfilled
func (r *InventoryReservation) FulfillReservation() (ddd.Event, error) {
	if r.Status != ReservationStatusActive {
		return nil, ErrReservationNotActive
	}

	r.AddEvent(ReservationFulfilledEvent, &ReservationFulfilled{
		FulfilledAt: time.Now(),
	})
	return ddd.NewEvent(ReservationFulfilledEvent, r), nil
}

// TransferReservation transfers the reservation to another warehouse
func (r *InventoryReservation) TransferReservation(toWarehouseID string) (ddd.Event, error) {
	if r.Status != ReservationStatusActive {
		return nil, ErrReservationNotActive
	}
	if toWarehouseID == "" {
		return nil, ErrMissingWarehouseID
	}

	r.AddEvent(ReservationTransferredEvent, &ReservationTransferred{
		FromWarehouseID: r.WarehouseID,
		ToWarehouseID:   toWarehouseID,
		TransferredAt:   time.Now(),
	})
	return ddd.NewEvent(ReservationTransferredEvent, r), nil
}

// ExpireReservation marks the reservation as expired
func (r *InventoryReservation) ExpireReservation() (ddd.Event, error) {
	if r.Status != ReservationStatusActive {
		return nil, ErrReservationNotActive
	}

	r.AddEvent(ReservationExpiredEvent, &ReservationExpired{
		ExpiredAt: time.Now(),
	})
	return ddd.NewEvent(ReservationExpiredEvent, r), nil
}

// ApplyEvent implements es.EventApplier
func (r *InventoryReservation) ApplyEvent(event ddd.Event) error {
	switch e := event.Payload().(type) {
	case *ReservationCreated:
		r.OrderID = e.OrderID
		r.SKU = e.SKU
		r.WarehouseID = e.WarehouseID
		r.Quantity = e.Quantity
		r.Status = ReservationStatusActive
		r.ExpiresAt = e.ExpiresAt
		r.CreatedAt = e.CreatedAt
		r.ConnectorID = e.ConnectorID

	case *ReservationReleased:
		r.Status = ReservationStatusReleased
		r.ReleasedAt = &e.ReleasedAt

	case *ReservationFulfilled:
		r.Status = ReservationStatusFulfilled
		r.FulfilledAt = &e.FulfilledAt

	case *ReservationTransferred:
		r.Status = ReservationStatusTransferred
		r.WarehouseID = e.ToWarehouseID
		r.TransferToWarehouse = e.ToWarehouseID
		r.TransferredAt = &e.TransferredAt

	case *ReservationExpired:
		r.Status = ReservationStatusExpired

	default:
		return errors.ErrInternal.Msgf(
			"%T received the event %s with unexpected payload %T",
			r, event.EventName(), e)
	}
	return nil
}

// ApplySnapshot implements es.Snapshotter
func (r *InventoryReservation) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *InventoryReservationV1:
		r.OrderID = ss.OrderID
		r.SKU = ss.SKU
		r.WarehouseID = ss.WarehouseID
		r.Quantity = ss.Quantity
		r.Status = ss.Status
		r.ExpiresAt = ss.ExpiresAt
		r.CreatedAt = ss.CreatedAt
		r.ReleasedAt = ss.ReleasedAt
		r.FulfilledAt = ss.FulfilledAt
		r.TransferredAt = ss.TransferredAt
		r.TransferToWarehouse = ss.TransferToWarehouse
		r.ConnectorID = ss.ConnectorID

	default:
		return errors.ErrInternal.Msgf(
			"%T received the unexpected snapshot %T", r, snapshot)
	}
	return nil
}

// ToSnapshot implements es.Snapshotter
func (r InventoryReservation) ToSnapshot() es.Snapshot {
	return InventoryReservationV1{
		OrderID:             r.OrderID,
		SKU:                 r.SKU,
		WarehouseID:         r.WarehouseID,
		Quantity:            r.Quantity,
		Status:              r.Status,
		ExpiresAt:           r.ExpiresAt,
		CreatedAt:           r.CreatedAt,
		ReleasedAt:          r.ReleasedAt,
		FulfilledAt:         r.FulfilledAt,
		TransferredAt:       r.TransferredAt,
		TransferToWarehouse: r.TransferToWarehouse,
		ConnectorID:         r.ConnectorID,
	}
}

// Fulfill fulfills the reservation (convenience method for FulfillReservation)
func (r *InventoryReservation) Fulfill() (ddd.Event, error) {
	return r.FulfillReservation()
}

// Release releases the reservation with a reason (convenience method for ReleaseReservation)
func (r *InventoryReservation) Release(reason string) (ddd.Event, error) {
	return r.ReleaseReservation(reason)
}

// Transfer transfers the reservation to another warehouse (convenience method for TransferReservation)
func (r *InventoryReservation) Transfer(toWarehouseID string) (ddd.Event, error) {
	return r.TransferReservation(toWarehouseID)
}