package commands

import (
	"context"
	"time"
	"github.com/stackus/errors"
	"middleman/internal/auth"
	"middleman/internal/ddd"
	"middleman/newsletters/internal/domain"
)

type UpdateEdition struct {
	ID           string
	Subject      string
	ContentHTML  string
	ContentText  string
	TemplateData map[string]string
	ScheduledAt  *time.Time
}

type UpdateEditionHandler struct {
	editions          domain.EditionRepository
	catalog           domain.EditionCatalogRepository
	newsletterCatalog domain.NewsletterCatalogRepository
	publisher         ddd.EventPublisher[ddd.Event]
}

func NewUpdateEditionHandler(
	editions domain.EditionRepository,
	catalog domain.EditionCatalogRepository,
	newsletterCatalog domain.NewsletterCatalogRepository,
	publisher ddd.EventPublisher[ddd.Event],
) UpdateEditionHandler {
	return UpdateEditionHandler{
		editions:          editions,
		catalog:           catalog,
		newsletterCatalog: newsletterCatalog,
		publisher:         publisher,
	}
}

func (h UpdateEditionHandler) UpdateEdition(ctx context.Context, cmd UpdateEdition) error {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return errors.ErrUnauthenticated.Msg("authentication required")
	}
	userID := claims.Subject
	
	edition, err := h.editions.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	// Verify newsletter exists and user owns it
	newsletter, err := h.newsletterCatalog.Find(ctx, edition.NewsletterID)
	if err != nil {
		return err
	}
	if newsletter.UserID != userID {
		return errors.ErrForbidden.Msg("not authorized to update editions for this newsletter")
	}

	event, err := edition.UpdateEdition(
		cmd.Subject,
		cmd.ContentHTML,
		cmd.ContentText,
		cmd.TemplateData,
		cmd.ScheduledAt,
	)
	if err != nil {
		return err
	}

	err = h.editions.Save(ctx, edition)
	if err != nil {
		return err
	}

	// Update catalog
	catalogEdition, err := h.catalog.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	catalogEdition.Subject = cmd.Subject
	catalogEdition.ContentHTML = cmd.ContentHTML
	catalogEdition.ContentText = cmd.ContentText
	catalogEdition.TemplateData = cmd.TemplateData
	catalogEdition.ScheduledAt = cmd.ScheduledAt
	catalogEdition.Status = edition.Status.String()
	catalogEdition.UpdatedAt = edition.UpdatedAt

	err = h.catalog.Update(ctx, catalogEdition)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}