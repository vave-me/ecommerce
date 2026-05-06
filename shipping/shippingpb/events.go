package shippingpb

import (
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

const (
	ShipmentAggregateChannel   = "middleman.Shipment.events.Shipment"
	ShipmentCreatedEvent       = "shippingapi.ShipmentCreated"
	ShipmentDeletedEvent       = "shippingapi.ShipmentDeleted"
	CarrierAssignedEvent       = "shippingapi.CarrierAssigned"
	ShipmentStartedEvent       = "shippingapi.ShipmentStarted"
	ShipmentStatusUpdatedEvent = "shippingapi.ShipmentStatusUpdated"
	ShipmentCancelledEvent     = "shippingapi.ShipmentCancelled"
	ShipmentDeliveredEvent     = "shippingapi.ShipmentDelivered"
	PickupScheduledEvent       = "shippingapi.PickupScheduled"
	ShipmentReturnedEvent      = "shippingapi.ShipmentReturned"
	ImageAggregateChannel      = "middleman.Shipment.events.Image"
	ImageAddedEvent            = "shippingapi.ImageAdded"
	VideoAggregateChannel      = "middleman.Shipment.events.Video"
	CommandChannel             = "middleman.shipping.commands"
	CreateShipmentCommand      = "shippingapi.CreateShipment"
)

func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) error {
	// Store events
	if err := serde.Register(&ShipmentCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&CarrierAssigned{}); err != nil {
		return err
	}
	if err := serde.Register(&ShipmentStarted{}); err != nil {
		return err
	}
	if err := serde.Register(&ShipmentStatusUpdated{}); err != nil {
		return err
	}
	if err := serde.Register(&ShipmentCancelled{}); err != nil {
		return err
	}
	if err := serde.Register(&ShipmentDelivered{}); err != nil {
		return err
	}
	if err := serde.Register(&PickupScheduled{}); err != nil {
		return err
	}
	if err := serde.Register(&ShipmentReturned{}); err != nil {
		return err
	}
	if err := serde.Register(&CreateShipment{}); err != nil {
		return err
	}
	return nil
}

func (*ShipmentCreated) Key() string       { return ShipmentCreatedEvent }
func (*CarrierAssigned) Key() string       { return CarrierAssignedEvent }
func (*ShipmentStarted) Key() string       { return ShipmentStartedEvent }
func (*ShipmentStatusUpdated) Key() string { return ShipmentStatusUpdatedEvent }
func (*ShipmentCancelled) Key() string     { return ShipmentCancelledEvent }
func (*ShipmentDelivered) Key() string     { return ShipmentDeliveredEvent }
func (*PickupScheduled) Key() string       { return PickupScheduledEvent }
func (*ShipmentReturned) Key() string      { return ShipmentReturnedEvent }
func (*CreateShipment) Key() string        { return CreateShipmentCommand }
