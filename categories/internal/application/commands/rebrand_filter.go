package commands

import (
	"context"

	"middleman/categories/internal/domain"
	"middleman/internal/ddd"
)

type RebrandFilter struct {
	ID          string
	Name        string
	Description string
}

type RebrandFilterHandler struct {
	filters   domain.FilterRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRebrandFilterHandler(
	filters domain.FilterRepository,
	publisher ddd.EventPublisher[ddd.Event],
) RebrandFilterHandler {
	return RebrandFilterHandler{
		filters:   filters,
		publisher: publisher,
	}
}

func (h RebrandFilterHandler) RebrandFilter(ctx context.Context, cmd RebrandFilter) error {
	filter, err := h.filters.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := filter.Rebrand(cmd.Name)
	// domain: filter.Rebrand name
	if err != nil {
		return err
	}

	if err = h.filters.Save(ctx, filter); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
