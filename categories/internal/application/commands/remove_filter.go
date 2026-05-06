package commands

import (
	"context"

	"middleman/categories/internal/domain"
	"middleman/internal/ddd"
)

type RemoveFilter struct {
	ID     string
	UserID string // optional if needed
}

type RemoveFilterHandler struct {
	filters   domain.FilterRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRemoveFilterHandler(
	filters domain.FilterRepository,
	publisher ddd.EventPublisher[ddd.Event],
) RemoveFilterHandler {
	return RemoveFilterHandler{
		filters:   filters,
		publisher: publisher,
	}
}

func (h RemoveFilterHandler) RemoveFilter(ctx context.Context, cmd RemoveFilter) error {
	filter, err := h.filters.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := filter.Remove() // domain: filter.Remove
	if err != nil {
		return err
	}

	if err = h.filters.Save(ctx, filter); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
