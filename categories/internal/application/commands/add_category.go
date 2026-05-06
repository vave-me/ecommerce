package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/categories/internal/domain"
	"middleman/internal/ddd"
)

type AddCategory struct {
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

type AddCategoryHandler struct {
	categories domain.CategoryRepository
	publisher  ddd.EventPublisher[ddd.Event]
}

func NewAddCategoryHandler(
	categories domain.CategoryRepository, publisher ddd.EventPublisher[ddd.Event]) AddCategoryHandler {

	return AddCategoryHandler{
		categories: categories,
		publisher:  publisher,
	}
}

func (h AddCategoryHandler) AddCategory(ctx context.Context, cmd AddCategory) error {
	category, err := h.categories.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error adding category")
	}

	event, err := category.InitCategory(cmd.ID, cmd.Description, cmd.ParentID, cmd.GoogleCategoryID, cmd.Tags, cmd.IsActive, cmd.Slug, cmd.SeoTitle, cmd.SeoKeywords, cmd.SeoDesc, cmd.CategoryType, cmd.Lang)
	if err != nil {
		return errors.Wrap(err, "initializing category")
	}

	err = h.categories.Save(ctx, category)
	if err != nil {
		return errors.Wrap(err, "error adding category")
	}

	return errors.Wrap(h.publisher.Publish(ctx, event), "publishing domain event")
}
