package servicespb

import (
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

// Channels
const (
	ServiceAggregateChannel = "middleman.services.events.Service"
)

// Service Event Names
const (
	ServiceAddedEvent             = "servicesapi.ServiceAdded"
	ServiceUpdatedEvent           = "servicesapi.ServiceUpdated"
	ServiceRebrandedEvent         = "servicesapi.ServiceRebranded"
	ServicePriceIncreasedEvent    = "servicesapi.ServicePriceIncreased"
	ServicePriceDecreasedEvent    = "servicesapi.ServicePriceDecreased"
	ServiceStockAdjustedEvent     = "servicesapi.ServiceStockAdjusted"
	ServiceRemovedEvent           = "servicesapi.ServiceRemoved"
	ServiceArchivedEvent          = "servicesapi.ServiceArchived"
	ServiceNegotiableToggledEvent = "servicesapi.ServiceNegotiableToggled"
	ServiceSoldEvent              = "servicesapi.ServiceSold"
	ServiceLeasedEvent            = "servicesapi.ServiceLeased"
	ServicePawnedEvent            = "servicesapi.ServicePawned"
)

// Registrations and Serde
func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) error {
	// Service
	if err := serde.Register(&ServiceAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&ServiceUpdated{}); err != nil {
		return err
	}
	if err := serde.Register(&ServiceRebranded{}); err != nil {
		return err
	}
	if err := serde.Register(&ServicePriceIncreased{}); err != nil {
		return err
	}
	if err := serde.Register(&ServicePriceDecreased{}); err != nil {
		return err
	}
	if err := serde.Register(&ServiceRemoved{}); err != nil {
		return err
	}
	if err := serde.Register(&ServiceStockAdjusted{}); err != nil {
		return err
	}
	if err := serde.Register(&ServiceArchived{}); err != nil {
		return err
	}
	if err := serde.Register(&ServiceNegotiableToggled{}); err != nil {
		return err
	}
	if err := serde.Register(&ServiceSold{}); err != nil {
		return err
	}
	if err := serde.Register(&ServiceLeased{}); err != nil {
		return err
	}
	if err := serde.Register(&ServicePawned{}); err != nil {
		return err
	}

	return nil
}

func (*ServiceAdded) Key() string             { return ServiceAddedEvent }
func (*ServiceUpdated) Key() string           { return ServiceUpdatedEvent }
func (*ServiceRebranded) Key() string         { return ServiceRebrandedEvent }
func (*ServicePriceIncreased) Key() string    { return ServicePriceIncreasedEvent }
func (*ServicePriceDecreased) Key() string    { return ServicePriceDecreasedEvent }
func (*ServiceRemoved) Key() string           { return ServiceRemovedEvent }
func (*ServiceStockAdjusted) Key() string     { return ServiceStockAdjustedEvent }
func (*ServiceArchived) Key() string          { return ServiceArchivedEvent }
func (*ServiceNegotiableToggled) Key() string { return ServiceNegotiableToggledEvent }
func (*ServiceSold) Key() string              { return ServiceSoldEvent }
func (*ServiceLeased) Key() string            { return ServiceLeasedEvent }
func (*ServicePawned) Key() string            { return ServicePawnedEvent }
