package commands

import (
	"context"

	"middleman/categories/internal/domain"
	"middleman/internal/ddd"
)

type ArchiveFilter struct {
	ID string // filter ID
}

type ArchiveFilterHandler struct {
	filters   domain.FilterRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewArchiveFilterHandler(
	filters domain.FilterRepository,
	publisher ddd.EventPublisher[ddd.Event],
) ArchiveFilterHandler {
	return ArchiveFilterHandler{
		filters:   filters,
		publisher: publisher,
	}
}

func (h ArchiveFilterHandler) ArchiveFilter(ctx context.Context, cmd ArchiveFilter) error {
	filter, err := h.filters.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := filter.Archive() // domain: filter.Archive
	if err != nil {
		return err
	}

	if err = h.filters.Save(ctx, filter); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
