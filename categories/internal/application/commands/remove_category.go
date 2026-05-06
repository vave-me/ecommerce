package commands

import (
	"context"
	"middleman/categories/internal/domain"
	"middleman/internal/ddd"
)

type RemoveCategory struct {
	ID     string
	UserID string
}

type RemoveCategoryHandler struct {
	categories domain.CategoryRepository
	publisher  ddd.EventPublisher[ddd.Event]
}

func NewRemoveCategoryHandler(categories domain.CategoryRepository, publisher ddd.EventPublisher[ddd.Event]) RemoveCategoryHandler {
	return RemoveCategoryHandler{
		categories: categories,
		publisher:  publisher,
	}
}

func (h RemoveCategoryHandler) RemoveCategory(ctx context.Context, cmd RemoveCategory) error {
	category, err := h.categories.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := category.Remove(cmd.ID)
	if err != nil {
		return err
	}

	err = h.categories.Save(ctx, category)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
