package commands

import (
	"context"
	"middleman/categories/internal/domain"
	"middleman/internal/ddd"
)

type RebrandCategory struct {
	ID          string
	Slug        string
	Description string
}

type RebrandCategoryHandler struct {
	categories domain.CategoryRepository
	publisher  ddd.EventPublisher[ddd.Event]
}

func NewRebrandCategoryHandler(categories domain.CategoryRepository, publisher ddd.EventPublisher[ddd.Event]) RebrandCategoryHandler {
	return RebrandCategoryHandler{
		categories: categories,
		publisher:  publisher,
	}
}

func (h RebrandCategoryHandler) RebrandCategory(ctx context.Context, cmd RebrandCategory) error {
	category, err := h.categories.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := category.Rebrand(cmd.Slug, cmd.Description)
	if err != nil {
		return err
	}

	err = h.categories.Save(ctx, category)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
