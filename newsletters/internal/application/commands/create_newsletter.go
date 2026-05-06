package commands

import (
	"context"
	"github.com/google/uuid"
	"middleman/internal/ddd"
	"middleman/newsletters/internal/domain"
)

type CreateNewsletter struct {
	UserID      string
	Name        string
	Description string
	Frequency   string
	Category    string
	TemplateID  string
}

type CreateNewsletterHandler struct {
	newsletters domain.NewsletterRepository
	catalog     domain.NewsletterCatalogRepository
	publisher   ddd.EventPublisher[ddd.Event]
}

func NewCreateNewsletterHandler(
	newsletters domain.NewsletterRepository,
	catalog domain.NewsletterCatalogRepository,
	publisher ddd.EventPublisher[ddd.Event],
) CreateNewsletterHandler {
	return CreateNewsletterHandler{
		newsletters: newsletters,
		catalog:     catalog,
		publisher:   publisher,
	}
}

func (h CreateNewsletterHandler) CreateNewsletter(ctx context.Context, cmd CreateNewsletter) (string, error) {
	newsletterID := uuid.New().String()
	
	newsletter := domain.NewNewsletter(newsletterID)
	
	event, err := newsletter.CreateNewsletter(
		cmd.UserID,
		cmd.Name,
		cmd.Description,
		cmd.Frequency,
		cmd.Category,
		cmd.TemplateID,
	)
	if err != nil {
		return "", err
	}

	err = h.newsletters.Save(ctx, newsletter)
	if err != nil {
		return "", err
	}

	// Add to catalog for queries
	catalogNewsletter := &domain.CatalogNewsletter{
		ID:           newsletterID,
		UserID:       cmd.UserID,
		Name:         cmd.Name,
		Description:  cmd.Description,
		Frequency:    cmd.Frequency,
		Category:     cmd.Category,
		TemplateID:   cmd.TemplateID,
		IsActive:     true,
		CreatedAt:    newsletter.CreatedAt,
		UpdatedAt:    newsletter.UpdatedAt,
	}
	
	err = h.catalog.Add(ctx, catalogNewsletter)
	if err != nil {
		return "", err
	}

	err = h.publisher.Publish(ctx, event)
	if err != nil {
		return "", err
	}

	return newsletterID, nil
}