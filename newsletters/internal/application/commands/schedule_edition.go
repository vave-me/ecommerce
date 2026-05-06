package commands

import (
	"context"
	"time"
	"github.com/stackus/errors"
	"middleman/internal/auth"
	"middleman/internal/ddd"
	"middleman/newsletters/internal/domain"
)

type ScheduleEdition struct {
	ID          string
	ScheduledAt time.Time
}

type ScheduleEditionHandler struct {
	editions          domain.EditionRepository
	catalog           domain.EditionCatalogRepository
	newsletterCatalog domain.NewsletterCatalogRepository
	publisher         ddd.EventPublisher[ddd.Event]
}

func NewScheduleEditionHandler(
	editions domain.EditionRepository,
	catalog domain.EditionCatalogRepository,
	newsletterCatalog domain.NewsletterCatalogRepository,
	publisher ddd.EventPublisher[ddd.Event],
) ScheduleEditionHandler {
	return ScheduleEditionHandler{
		editions:          editions,
		catalog:           catalog,
		newsletterCatalog: newsletterCatalog,
		publisher:         publisher,
	}
}

func (h ScheduleEditionHandler) ScheduleEdition(ctx context.Context, cmd ScheduleEdition) error {
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
		return errors.ErrForbidden.Msg("not authorized to schedule editions for this newsletter")
	}

	event, err := edition.ScheduleEdition(cmd.ScheduledAt)
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

	catalogEdition.ScheduledAt = &cmd.ScheduledAt
	catalogEdition.Status = domain.ScheduledStatus.String()
	catalogEdition.UpdatedAt = edition.UpdatedAt

	err = h.catalog.Update(ctx, catalogEdition)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}