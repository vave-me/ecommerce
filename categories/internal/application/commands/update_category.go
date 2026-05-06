package commands

import (
	"context"
	"middleman/categories/internal/domain"
	"middleman/internal/ddd"
)

type UpdateCategory struct {
	ID               string
	Description      string
	ParentID         string
	GoogleCategoryID string
	Tags             []string // optional set of tags
	IsActive         bool
	Slug             string
	SeoTitle         string
	SeoKeywords      []string
	SeoDesc          string
	CategoryType     string
	Lang             string
}

type UpdateCategoryHandler struct {
	categories domain.CategoryRepository
	publisher  ddd.EventPublisher[ddd.Event]
}

func NewUpdateCategoryHandler(categories domain.CategoryRepository, publisher ddd.EventPublisher[ddd.Event]) UpdateCategoryHandler {
	return UpdateCategoryHandler{
		categories: categories,
		publisher:  publisher,
	}
}

func (h UpdateCategoryHandler) UpdateCategory(ctx context.Context, cmd UpdateCategory) error {
	category, err := h.categories.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := category.Update(cmd.Description, cmd.ParentID, cmd.GoogleCategoryID, cmd.Tags, cmd.IsActive, cmd.Slug, cmd.SeoTitle, cmd.SeoKeywords, cmd.SeoDesc, cmd.CategoryType, cmd.Lang)
	if err != nil {
		return err
	}

	err = h.categories.Save(ctx, category)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
