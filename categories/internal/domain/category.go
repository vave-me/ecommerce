package domain

import (
	"github.com/stackus/errors"

	"middleman/internal/ddd"
	"middleman/internal/es"
)

const (
	// CategoryAggregate is the unique identifier for this aggregate.
	CategoryAggregate = "categories.Category"
)

// Domain errors
var (
	ErrCategoryNameIsBlank = errors.Wrap(errors.ErrBadRequest, "the category name cannot be blank")
)

// Category is our aggregate root for category-related behavior.
type Category struct {
	es.Aggregate
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
}

// Ensure Category implements the required interfaces for event sourcing.
var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Category)(nil)

// NewCategory is a constructor for a new category domain aggregate.
func NewCategory(id string) *Category {
	return &Category{
		Aggregate: es.NewAggregate(id, CategoryAggregate),
		IsActive:  true, // by default new categories might be active
	}
}

// Key implements registry.Registerable, for aggregator discovery or similar.
func (Category) Key() string {
	return CategoryAggregate
}

// InitCategory handles the initial creation/adding of a Category.
func (c *Category) InitCategory(
	id, description, parentID, googleCategoryID string, tags []string, isActive bool, slug string, seoTitle string, seoKeywords []string, seoDesc string, categoryType string, lang string) (ddd.Event, error) {
	if slug == "" {
		return nil, ErrCategoryNameIsBlank
	}

	c.AddEvent(CategoryAddedEvent, &CategoryAdded{
		CategoryID:       c.ID(),
		Description:      description,
		Tags:             tags,
		IsActive:         isActive,
		ParentID:         parentID,
		GoogleCategoryID: googleCategoryID,
		Slug:             slug,
		SeoTitle:         seoTitle,
		SeoKeywords:      seoKeywords,
		SeoDesc:          seoDesc,
		CategoryType:     categoryType,
		Lang:             lang,
	})

	return ddd.NewEvent(CategoryAddedEvent, c), nil
}

// Update modifies an existing Category’s basic fields.
func (c *Category) Update(
	description, parentID, googleCategoryID string,
	tags []string, // optional set of tags
	isActive bool,
	slug, seoTitle string,
	seoKeywords []string,
	seoDesc string,
	categoryType string,
	lang string,
) (ddd.Event, error) {
	if slug == "" {
		return nil, ErrCategoryNameIsBlank
	}

	c.AddEvent(CategoryUpdatedEvent, &CategoryUpdated{
		CategoryID:       c.ID(),
		Description:      description,
		Tags:             tags,
		IsActive:         isActive,
		SeoTitle:         seoTitle,
		SeoKeywords:      seoKeywords,
		SeoDesc:          seoDesc,
		ParentID:         parentID,
		GoogleCategoryID: googleCategoryID,
		Slug:             slug,
		CategoryType:     categoryType,
		Lang:             lang,
	})

	return ddd.NewEvent(CategoryUpdatedEvent, c), nil
}

// Rebrand modifies only name/description, if that’s a separate concept.
func (c *Category) Rebrand(slug, description string) (ddd.Event, error) {
	if slug == "" {
		return nil, ErrCategoryNameIsBlank
	}

	c.AddEvent(CategoryRebrandedEvent, &CategoryRebranded{
		Slug:        slug,
		Description: description,
	})

	return ddd.NewEvent(CategoryRebrandedEvent, c), nil
}

// Remove could mark the category for deletion or immediate removal.
func (c *Category) Remove(userID string) (ddd.Event, error) {
	c.AddEvent(CategoryRemovedEvent, &CategoryRemoved{
		CategoryID: c.ID(),
		UserID:     userID,
	})

	return ddd.NewEvent(CategoryRemovedEvent, c), nil
}

// Archive sets isActive=false or a “status = archived”.
func (c *Category) Archive(userID string) (ddd.Event, error) {
	c.AddEvent(CategoryArchivedEvent, &CategoryArchived{
		CategoryID: c.ID(),
	})
	return ddd.NewEvent(CategoryArchivedEvent, c), nil
}

// ApplyEvent changes the Category state based on domain events.
func (c *Category) ApplyEvent(event ddd.Event) error {
	switch e := event.Payload().(type) {

	case *CategoryAdded:
		c.Description = e.Description
		c.ParentID = e.ParentID
		c.GoogleCategoryID = e.GoogleCategoryID
		c.Slug = e.Slug
		c.IsActive = true // new category is active by default

	case *CategoryUpdated:

		c.Description = e.Description
		c.ParentID = e.ParentID
		c.GoogleCategoryID = e.GoogleCategoryID
		c.Slug = e.Slug
		// c.IsActive remains or re-activated if needed

	case *CategoryRebranded:
		c.Slug = e.Slug
		c.Description = e.Description

	case *CategoryRemoved:
		// optionally do c.IsActive = false or track “removed”
		// e.UserID can be used in logs or business logic

	case *CategoryArchived:
		// c.IsActive = false or c.Status = “archived”

	default:
		return errors.ErrInternal.Msgf(
			"Category(%s) received unexpected event payload type %T for event %s",
			c.ID(), e, event.EventName(),
		)
	}
	return nil
}

// Snapshot-related logic:

// ApplySnapshot to restore from a saved snapshot.
func (c *Category) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *CategoryV1:
		c.Description = ss.Description
		c.ParentID = ss.ParentID
		c.GoogleCategoryID = ss.GoogleCategoryID
		c.IsActive = ss.IsActive
		c.Slug = ss.Slug
		c.SeoTitle = ss.SeoTitle
		c.SeoKeywords = ss.SeoKeywords
		c.SeoDesc = ss.SeoDesc

	default:
		return errors.ErrInternal.Msgf(
			"Category(%s) received unexpected snapshot type %T",
			c.ID(), snapshot,
		)
	}
	return nil
}

// ToSnapshot creates a new snapshot from the current Category state.
func (c Category) ToSnapshot() es.Snapshot {
	return &CategoryV1{
		Description:      c.Description,
		ParentID:         c.ParentID,
		GoogleCategoryID: c.GoogleCategoryID,
		IsActive:         c.IsActive,
		Tags:             c.Tags,
		Slug:             c.Slug,
		SeoTitle:         c.SeoTitle,
		SeoKeywords:      c.SeoKeywords,
		SeoDesc:          c.SeoDesc,
	}
}
