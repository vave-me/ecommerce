package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/newsletters/internal/domain"
)

type UpdateNewsletter struct {
	ID          string
	Name        string
	Description string
	Frequency   string
	Category    string
	TemplateID  string
	IsActive    bool
}

type UpdateNewsletterHandler struct {
	newsletters domain.NewsletterRepository
	catalog     domain.NewsletterCatalogRepository
	publisher   ddd.EventPublisher[ddd.Event]
}

func NewUpdateNewsletterHandler(
	newsletters domain.NewsletterRepository,
	catalog domain.NewsletterCatalogRepository,
	publisher ddd.EventPublisher[ddd.Event],
) UpdateNewsletterHandler {
	return UpdateNewsletterHandler{
		newsletters: newsletters,
		catalog:     catalog,
		publisher:   publisher,
	}
}

func (h UpdateNewsletterHandler) UpdateNewsletter(ctx context.Context, cmd UpdateNewsletter) error {
	newsletter, err := h.newsletters.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	// Handle activation/deactivation
	var event ddd.Event
	if cmd.IsActive && !newsletter.IsActive {
		event, err = newsletter.ActivateNewsletter()
		if err != nil {
			return err
		}
	} else if !cmd.IsActive && newsletter.IsActive {
		event, err = newsletter.DeactivateNewsletter()
		if err != nil {
			return err
		}
	}

	// Update newsletter details
	updateEvent, err := newsletter.UpdateNewsletter(
		cmd.Name,
		cmd.Description,
		cmd.Frequency,
		cmd.Category,
		cmd.TemplateID,
	)
	if err != nil {
		return err
	}

	err = h.newsletters.Save(ctx, newsletter)
	if err != nil {
		return err
	}

	// Update catalog
	catalogNewsletter, err := h.catalog.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	catalogNewsletter.Name = cmd.Name
	catalogNewsletter.Description = cmd.Description
	catalogNewsletter.Frequency = cmd.Frequency
	catalogNewsletter.Category = cmd.Category
	catalogNewsletter.TemplateID = cmd.TemplateID
	catalogNewsletter.IsActive = cmd.IsActive
	catalogNewsletter.UpdatedAt = newsletter.UpdatedAt

	err = h.catalog.Update(ctx, catalogNewsletter)
	if err != nil {
		return err
	}

	// Publish events
	if event != nil {
		err = h.publisher.Publish(ctx, event)
		if err != nil {
			return err
		}
	}

	return h.publisher.Publish(ctx, updateEvent)
}