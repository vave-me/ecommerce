package commands

import (
	"context"

	"middleman/categories/internal/domain"
	"middleman/internal/ddd"
)

// ArchiveCategory marks a category as archived or inactive.
type ArchiveCategory struct {
	ID string
}

type ArchiveCategoryHandler struct {
	categories domain.CategoryRepository
	publisher  ddd.EventPublisher[ddd.Event]
}

func NewArchiveCategoryHandler(
	categories domain.CategoryRepository,
	publisher ddd.EventPublisher[ddd.Event],
) ArchiveCategoryHandler {
	return ArchiveCategoryHandler{
		categories: categories,
		publisher:  publisher,
	}
}

func (h ArchiveCategoryHandler) ArchiveCategory(ctx context.Context, cmd ArchiveCategory) error {
	category, err := h.categories.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := category.Archive(cmd.ID) // domain: category.Archive(userSellerID)
	if err != nil {
		return err
	}

	if err = h.categories.Save(ctx, category); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
