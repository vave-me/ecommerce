package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/newsletters/internal/domain"
)

type DeleteNewsletter struct {
	ID string
}

type DeleteNewsletterHandler struct {
	newsletters domain.NewsletterRepository
	catalog     domain.NewsletterCatalogRepository
	publisher   ddd.EventPublisher[ddd.Event]
}

func NewDeleteNewsletterHandler(
	newsletters domain.NewsletterRepository,
	catalog domain.NewsletterCatalogRepository,
	publisher ddd.EventPublisher[ddd.Event],
) DeleteNewsletterHandler {
	return DeleteNewsletterHandler{
		newsletters: newsletters,
		catalog:     catalog,
		publisher:   publisher,
	}
}

func (h DeleteNewsletterHandler) DeleteNewsletter(ctx context.Context, cmd DeleteNewsletter) error {
	newsletter, err := h.newsletters.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := newsletter.DeleteNewsletter()
	if err != nil {
		return err
	}

	err = h.newsletters.Save(ctx, newsletter)
	if err != nil {
		return err
	}

	// Remove from catalog
	err = h.catalog.Delete(ctx, cmd.ID)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}