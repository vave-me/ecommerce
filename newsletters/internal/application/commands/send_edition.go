package commands

import (
	"context"
	"time"
	"github.com/stackus/errors"
	"middleman/internal/auth"
	"middleman/internal/ddd"
	"middleman/newsletters/internal/domain"
)

type SendEdition struct {
	ID       string
	TestMode bool // Send only to owner for testing
}

type SendEditionHandler struct {
	editions          domain.EditionRepository
	catalog           domain.EditionCatalogRepository
	newsletterCatalog domain.NewsletterCatalogRepository
	subCatalog        domain.SubscriptionCatalogRepository
	publisher         ddd.EventPublisher[ddd.Event]
}

func NewSendEditionHandler(
	editions domain.EditionRepository,
	catalog domain.EditionCatalogRepository,
	newsletterCatalog domain.NewsletterCatalogRepository,
	subCatalog domain.SubscriptionCatalogRepository,
	publisher ddd.EventPublisher[ddd.Event],
) SendEditionHandler {
	return SendEditionHandler{
		editions:          editions,
		catalog:           catalog,
		newsletterCatalog: newsletterCatalog,
		subCatalog:        subCatalog,
		publisher:         publisher,
	}
}

func (h SendEditionHandler) SendEdition(ctx context.Context, cmd SendEdition) (int, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return 0, errors.ErrUnauthenticated.Msg("authentication required")
	}
	userID := claims.Subject
	
	edition, err := h.editions.Load(ctx, cmd.ID)
	if err != nil {
		return 0, err
	}

	// Verify newsletter exists and user owns it
	newsletter, err := h.newsletterCatalog.Find(ctx, edition.NewsletterID)
	if err != nil {
		return 0, err
	}
	if newsletter.UserID != userID {
		return 0, errors.ErrForbidden.Msg("not authorized to send editions for this newsletter")
	}

	// Mark as sending
	sendingEvent, err := edition.MarkSending()
	if err != nil {
		return 0, err
	}

	// Get recipient count
	recipientCount := 1 // Owner in test mode
	if !cmd.TestMode {
		count, err := h.subCatalog.CountActiveByNewsletter(ctx, edition.NewsletterID)
		if err != nil {
			return 0, err
		}
		recipientCount = count
	}

	// Mark as sent
	event, err := edition.SendEdition(recipientCount)
	if err != nil {
		return 0, err
	}

	err = h.editions.Save(ctx, edition)
	if err != nil {
		return 0, err
	}

	// Update catalog
	catalogEdition, err := h.catalog.Find(ctx, cmd.ID)
	if err != nil {
		return 0, err
	}

	catalogEdition.Status = domain.SentStatus.String()
	now := time.Now()
	catalogEdition.SentAt = &now
	catalogEdition.RecipientCount = recipientCount
	catalogEdition.UpdatedAt = edition.UpdatedAt

	err = h.catalog.Update(ctx, catalogEdition)
	if err != nil {
		return 0, err
	}

	// Publish events
	err = h.publisher.Publish(ctx, sendingEvent)
	if err != nil {
		return 0, err
	}

	err = h.publisher.Publish(ctx, event)
	if err != nil {
		return 0, err
	}

	// TODO: Queue actual email sending to subscribers

	return recipientCount, nil
}