package domain

import (
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const ShipmentAggregate = "Shipment.Shipment"

var (
	ErrProductIDIsBlank    = errors.Wrap(errors.ErrBadRequest, "product id cannot be blank")
	ErrInvalidStatus       = errors.Wrap(errors.ErrBadRequest, "invalid shipment status")
	ErrShipmentCancelled   = errors.Wrap(errors.ErrBadRequest, "shipment is already cancelled")
	ErrShipmentDelivered   = errors.Wrap(errors.ErrBadRequest, "shipment is already delivered")
	ErrCarrierNotAssigned  = errors.Wrap(errors.ErrBadRequest, "carrier not assigned")
	ErrShipmentNotStarted  = errors.Wrap(errors.ErrBadRequest, "shipment not started")
)

type Shipment struct {
	es.Aggregate
	ProductID       string
	OrderID         string
	BasketID        string
	TrackingNumber  string
	LabelUrl        string
	SenderName      string
	SenderAddress   string
	ReceiverName    string
	ReceiverAddress string
	Weight          string
	Dimensions      string
	ServiceType     string
	Status          string
	CarrierID       string
	CarrierName     string
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Shipment)(nil)

func NewShipment(id string) *Shipment {
	return &Shipment{
		Aggregate: es.NewAggregate(id, ShipmentAggregate),
	}
}

func (m *Shipment) InitShipment(productID, orderID, basketID, trackingNumber, labelUrl, senderName, senderAddress, receiverName, receiverAddress, weight, dimensions, serviceType string) (ddd.Event, error) {
	if productID == "" {
		return nil, ErrProductIDIsBlank
	}

	m.Status = "created"
	m.AddEvent(ShipmentCreatedEvent, &ShipmentCreated{
		ProductID:       productID,
		OrderID:         orderID,
		BasketID:        basketID,
		TrackingNumber:  trackingNumber,
		LabelUrl:        labelUrl,
		SenderName:      senderName,
		SenderAddress:   senderAddress,
		ReceiverName:    receiverName,
		ReceiverAddress: receiverAddress,
		Weight:          weight,
		Dimensions:      dimensions,
		ServiceType:     serviceType,
		Status:          "created",
	})

	return ddd.NewEvent(ShipmentCreatedEvent, m), nil
}

func (m *Shipment) AssignCarrier(carrierID, carrierName string) (ddd.Event, error) {
	if m.Status == "cancelled" {
		return nil, ErrShipmentCancelled
	}
	if m.Status == "delivered" {
		return nil, ErrShipmentDelivered
	}

	m.AddEvent(CarrierAssignedEvent, &CarrierAssigned{
		ShipmentID:  m.ID(),
		CarrierID:   carrierID,
		CarrierName: carrierName,
	})

	return ddd.NewEvent(CarrierAssignedEvent, m), nil
}

func (m *Shipment) StartShipment() (ddd.Event, error) {
	if m.Status == "cancelled" {
		return nil, ErrShipmentCancelled
	}
	if m.Status == "delivered" {
		return nil, ErrShipmentDelivered
	}
	if m.CarrierID == "" {
		return nil, ErrCarrierNotAssigned
	}

	m.AddEvent(ShipmentStartedEvent, &ShipmentStarted{
		ShipmentID: m.ID(),
	})

	return ddd.NewEvent(ShipmentStartedEvent, m), nil
}

func (m *Shipment) UpdateStatus(status, location, notes string) (ddd.Event, error) {
	if m.Status == "cancelled" {
		return nil, ErrShipmentCancelled
	}
	if m.Status == "delivered" {
		return nil, ErrShipmentDelivered
	}

	m.AddEvent(ShipmentStatusUpdatedEvent, &ShipmentStatusUpdated{
		ShipmentID: m.ID(),
		Status:     status,
		Location:   location,
		Notes:      notes,
	})

	return ddd.NewEvent(ShipmentStatusUpdatedEvent, m), nil
}

func (m *Shipment) CancelShipment(reason string) (ddd.Event, error) {
	if m.Status == "cancelled" {
		return nil, ErrShipmentCancelled
	}
	if m.Status == "delivered" {
		return nil, ErrShipmentDelivered
	}

	m.AddEvent(ShipmentCancelledEvent, &ShipmentCancelled{
		ShipmentID: m.ID(),
		Reason:     reason,
	})

	return ddd.NewEvent(ShipmentCancelledEvent, m), nil
}

func (m *Shipment) MarkAsDelivered(signedBy string, proofURL string) (ddd.Event, error) {
	if m.Status == "cancelled" {
		return nil, ErrShipmentCancelled
	}
	if m.Status == "delivered" {
		return nil, ErrShipmentDelivered
	}
	if m.Status != "in_transit" && m.Status != "out_for_delivery" {
		return nil, ErrShipmentNotStarted
	}

	m.AddEvent(ShipmentDeliveredEvent, &ShipmentDelivered{
		ShipmentID:           m.ID(),
		SignedBy:             signedBy,
		ProofOfDeliveryURL:   proofURL,
	})

	return ddd.NewEvent(ShipmentDeliveredEvent, m), nil
}

func (m *Shipment) SchedulePickup(pickupTime string, instructions string) (ddd.Event, error) {
	if m.Status == "cancelled" {
		return nil, ErrShipmentCancelled
	}
	if m.Status == "delivered" {
		return nil, ErrShipmentDelivered
	}

	m.AddEvent(PickupScheduledEvent, &PickupScheduled{
		ShipmentID:   m.ID(),
		PickupTime:   pickupTime,
		Instructions: instructions,
	})

	return ddd.NewEvent(PickupScheduledEvent, m), nil
}

func (m *Shipment) InitiateReturn(reason, returnTrackingNumber string) (ddd.Event, error) {
	if m.Status != "delivered" {
		return nil, errors.Wrap(errors.ErrBadRequest, "can only return delivered shipments")
	}

	m.AddEvent(ShipmentReturnedEvent, &ShipmentReturned{
		ShipmentID:           m.ID(),
		Reason:               reason,
		ReturnTrackingNumber: returnTrackingNumber,
	})

	return ddd.NewEvent(ShipmentReturnedEvent, m), nil
}

// Key implements registry.Registerable
func (Shipment) Key() string { return ShipmentAggregate }

// ApplyEvent implements es.EventApplier
func (m *Shipment) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *ShipmentCreated:
		m.ProductID = payload.ProductID
		m.OrderID = payload.OrderID
		m.BasketID = payload.BasketID
		m.TrackingNumber = payload.TrackingNumber
		m.LabelUrl = payload.LabelUrl
		m.SenderName = payload.SenderName
		m.SenderAddress = payload.SenderAddress
		m.ReceiverName = payload.ReceiverName
		m.ReceiverAddress = payload.ReceiverAddress
		m.Weight = payload.Weight
		m.Dimensions = payload.Dimensions
		m.ServiceType = payload.ServiceType
		m.Status = payload.Status

	case *CarrierAssigned:
		m.CarrierID = payload.CarrierID
		m.CarrierName = payload.CarrierName
		m.Status = "carrier_assigned"

	case *ShipmentStarted:
		m.Status = "in_transit"

	case *ShipmentStatusUpdated:
		m.Status = payload.Status

	case *ShipmentCancelled:
		m.Status = "cancelled"

	case *ShipmentDelivered:
		m.Status = "delivered"

	case *PickupScheduled:
		m.Status = "pickup_scheduled"

	case *ShipmentReturned:
		m.Status = "returned"

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", m, event.EventName(), payload)
	}

	return nil
}

// ApplySnapshot implements es.Snapshotter
func (m *Shipment) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *ShipmentV1:
		m.ProductID = ss.ProductID
		m.OrderID = ss.OrderID
		m.BasketID = ss.BasketID
		m.TrackingNumber = ss.TrackingNumber
		m.LabelUrl = ss.LabelUrl
		m.SenderName = ss.SenderName
		m.SenderAddress = ss.SenderAddress
		m.ReceiverName = ss.ReceiverName
		m.ReceiverAddress = ss.ReceiverAddress
		m.Weight = ss.Weight
		m.Dimensions = ss.Dimensions
		m.ServiceType = ss.ServiceType
		m.Status = ss.Status
		m.CarrierID = ss.CarrierID
		m.CarrierName = ss.CarrierName

	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", m, snapshot)
	}

	return nil
}

// ToSnapshot implements es.Snapshotter
func (s Shipment) ToSnapshot() es.Snapshot {
	return ShipmentV1{
		ProductID:       s.ProductID,
		OrderID:         s.OrderID,
		BasketID:        s.BasketID,
		TrackingNumber:  s.TrackingNumber,
		LabelUrl:        s.LabelUrl,
		SenderName:      s.SenderName,
		SenderAddress:   s.SenderAddress,
		ReceiverName:    s.ReceiverName,
		ReceiverAddress: s.ReceiverAddress,
		Weight:          s.Weight,
		Dimensions:      s.Dimensions,
		ServiceType:     s.ServiceType,
		Status:          s.Status,
		CarrierID:       s.CarrierID,
		CarrierName:     s.CarrierName,
	}
}
