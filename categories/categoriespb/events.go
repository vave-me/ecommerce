package categoriespb

import (
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

// Channels (unchanged or simplified)
const (
	CategoryAggregateChannel = "middleman.categories.events.Category"
	FilterAggregateChannel   = "middleman.categories.events.Filter"
	CategoryAddedEvent       = "categoriesapi.CategoryAdded"
	CategoryUpdatedEvent     = "categoriesapi.CategoryUpdated"
	CategoryRebrandedEvent   = "categoriesapi.CategoryRebranded"
	CategoryRemovedEvent     = "categoriesapi.CategoryRemoved"
	CategoryArchivedEvent    = "categoriesapi.CategoryArchived"
	FilterAddedEvent         = "categoriesapi.FilterAdded"
	FilterRebrandedEvent     = "categoriesapi.FilterRebranded"
	FilterArchivedEvent      = "categoriesapi.FilterArchived"
	FilterRemovedEvent       = "categoriesapi.FilterRemoved"
)

func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) error {
	// Category events
	if err := serde.Register(&CategoryAdded{}); err != nil {
		return err
	}

	if err := serde.Register(&CategoryRebranded{}); err != nil {
		return err
	}
	if err := serde.Register(&CategoryRemoved{}); err != nil {
		return err
	}
	if err := serde.Register(&CategoryArchived{}); err != nil {
		return err
	}

	// Filter events
	if err := serde.Register(&FilterAdded{}); err != nil {
		return err
	}

	if err := serde.Register(&FilterArchived{}); err != nil {
		return err
	}
	if err := serde.Register(&FilterRemoved{}); err != nil {
		return err
	}

	return nil
}
func (*CategoryAdded) Key() string { return CategoryAddedEvent }

func (*CategoryRebranded) Key() string { return CategoryRebrandedEvent }

func (*CategoryRemoved) Key() string { return CategoryRemovedEvent }

func (*CategoryArchived) Key() string { return CategoryArchivedEvent }

func (*FilterAdded) Key() string { return FilterAddedEvent }

func (*FilterArchived) Key() string { return FilterArchivedEvent }

func (*FilterRemoved) Key() string { return FilterRemovedEvent }
