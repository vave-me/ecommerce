package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/categories/internal/domain"
	"middleman/internal/ddd"
)

type AddFilter struct {
	ID         string
	CategoryID string
	Name       string
	FilterType domain.FilterType
	Values     []string
}

type AddFilterHandler struct {
	filters   domain.FilterRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAddFilterHandler(filters domain.FilterRepository, publisher ddd.EventPublisher[ddd.Event]) AddFilterHandler {

	return AddFilterHandler{
		filters:   filters,
		publisher: publisher,
	}
}

func (h AddFilterHandler) AddFilter(ctx context.Context, cmd AddFilter) error {
	filter, err := h.filters.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error adding category")
	}

	event, err := filter.InitFilter(cmd.ID, cmd.CategoryID, cmd.Name, cmd.FilterType, cmd.Values)
	if err != nil {
		return errors.Wrap(err, "initializing category")
	}

	err = h.filters.Save(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "error adding category")
	}

	return errors.Wrap(h.publisher.Publish(ctx, event), "publishing domain event")
}
