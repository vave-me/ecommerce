package domain

import (
	"github.com/stackus/errors"

	"middleman/internal/ddd"
	"middleman/internal/es"
)

const FilterAggregate = "categories.Filter"

// Domain errors
var (
	ErrFilterNameBlank = errors.Wrap(errors.ErrBadRequest, "the filter name cannot be blank")
)

type Filter struct {
	es.Aggregate
	CategoryID string
	Name       string
	FilterType FilterType
	Values     []string
	IsActive   bool
}

// Make sure Filter implements the correct interfaces
var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Filter)(nil)

// NewFilter constructor
func NewFilter(id string) *Filter {
	return &Filter{
		Aggregate: es.NewAggregate(id, FilterAggregate),
		IsActive:  true,
	}
}

// Key implements registry.Registerable
func (Filter) Key() string { return FilterAggregate }

// InitFilter sets the initial state for a new filter.
func (f *Filter) InitFilter(id, categoryID, name string, filterType FilterType, values []string) (ddd.Event, error) {
	if name == "" {
		return nil, ErrFilterNameBlank
	}

	f.AddEvent(FilterAddedEvent, &FilterAdded{
		CategoryID: categoryID,
		Name:       name,
		FilterType: filterType,
		Values:     values,
		IsActive:   true,
	})

	return ddd.NewEvent(FilterAddedEvent, f), nil
}

// Update modifies the filter’s name/type/values.
func (f *Filter) Update(
	name string,
	filterType FilterType,
	values []string,
) (ddd.Event, error) {
	if name == "" {
		return nil, ErrFilterNameBlank
	}

	f.AddEvent(FilterUpdatedEvent, &FilterUpdated{
		Name:       name,
		FilterType: filterType,
		Values:     values,
	})
	return ddd.NewEvent(FilterUpdatedEvent, f), nil
}

// Archive sets the filter as inactive.
func (f *Filter) Archive() (ddd.Event, error) {
	if !f.IsActive {
		// Potentially do nothing, or return an error that it’s already inactive
	}
	f.AddEvent(FilterArchivedEvent, &FilterArchived{
		FilterID: f.ID(),
	})
	return ddd.NewEvent(FilterArchivedEvent, f), nil
}

// Remove might mark it as removed from domain or physically delete it
func (f *Filter) Remove() (ddd.Event, error) {
	f.AddEvent(FilterRemovedEvent, &FilterRemoved{
		FilterID: f.ID(),
	})
	return ddd.NewEvent(FilterRemovedEvent, f), nil
}

// Rebrand modifies only the filter name
func (f *Filter) Rebrand(name string) (ddd.Event, error) {
	if name == "" {
		return nil, ErrFilterNameBlank
	}
	f.AddEvent(FilterRebrandedEvent, &FilterRebranded{
		FilterID: f.ID(),
		Name:     name,
	})
	return ddd.NewEvent(FilterRebrandedEvent, f), nil
}

// ApplyEvent applies event changes to the filter.
func (f *Filter) ApplyEvent(event ddd.Event) error {
	switch e := event.Payload().(type) {

	case *FilterAdded:
		f.CategoryID = e.CategoryID
		f.Name = e.Name
		f.FilterType = e.FilterType
		f.Values = e.Values
		f.IsActive = true

	case *FilterUpdated:
		f.Name = e.Name
		f.FilterType = e.FilterType
		f.Values = e.Values

	case *FilterArchived:
		f.IsActive = false

	case *FilterRemoved:
		// Possibly mark a “removed” state if needed
		// or do nothing if ephemeral

	case *FilterRebranded:
		f.Name = e.Name

	default:
		return errors.ErrInternal.Msgf(
			"%T received unexpected event type %s with payload %T",
			f, event.EventName(), e,
		)
	}
	return nil
}

// Snapshotter interface:

// ApplySnapshot hydrates the Filter from a snapshot.
func (f *Filter) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *FilterV1:
		f.CategoryID = ss.CategoryID
		f.Name = ss.Name
		f.FilterType = ss.FilterType
		f.Values = ss.Values
		f.IsActive = ss.IsActive
	default:
		return errors.ErrInternal.Msgf(
			"Filter(%s) got unexpected snapshot type %T", f.ID(), snapshot)
	}
	return nil
}

func (f Filter) ToSnapshot() es.Snapshot {
	return &FilterV1{
		CategoryID: f.CategoryID,
		Name:       f.Name,
		FilterType: f.FilterType,
		Values:     f.Values,
		IsActive:   f.IsActive,
	}
}
