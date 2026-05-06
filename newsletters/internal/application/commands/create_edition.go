package commands

import (
	"context"
	"time"
	"github.com/google/uuid"
	"github.com/stackus/errors"
	"middleman/internal/auth"
	"middleman/internal/ddd"
	"middleman/newsletters/internal/domain"
)

type CreateEdition struct {
	NewsletterID string
	Subject      string
	ContentHTML  string
	ContentText  string
	TemplateData map[string]string
	ScheduledAt  *time.Time
}

type CreateEditionHandler struct {
	editions          domain.EditionRepository
	catalog           domain.EditionCatalogRepository
	newsletterCatalog domain.NewsletterCatalogRepository
	publisher         ddd.EventPublisher[ddd.Event]
}

func NewCreateEditionHandler(
	editions domain.EditionRepository,
	catalog domain.EditionCatalogRepository,
	newsletterCatalog domain.NewsletterCatalogRepository,
	publisher ddd.EventPublisher[ddd.Event],
) CreateEditionHandler {
	return CreateEditionHandler{
		editions:          editions,
		catalog:           catalog,
		newsletterCatalog: newsletterCatalog,
		publisher:         publisher,
	}
}

func (h CreateEditionHandler) CreateEdition(ctx context.Context, cmd CreateEdition) (string, error) {
	// Get user ID from context
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return "", errors.ErrUnauthenticated.Msg("authentication required")
	}
	userID := claims.Subject
	
	// Verify newsletter exists and user owns it
	newsletter, err := h.newsletterCatalog.Find(ctx, cmd.NewsletterID)
	if err != nil {
		return "", err
	}
	if newsletter.UserID != userID {
		return "", errors.ErrForbidden.Msg("not authorized to create editions for this newsletter")
	}

	editionID := uuid.New().String()
	edition := domain.NewEdition(editionID)

	event, err := edition.CreateEdition(
		cmd.NewsletterID,
		cmd.Subject,
		cmd.ContentHTML,
		cmd.ContentText,
		userID,
		cmd.TemplateData,
		cmd.ScheduledAt,
	)
	if err != nil {
		return "", err
	}

	err = h.editions.Save(ctx, edition)
	if err != nil {
		return "", err
	}

	// Add to catalog
	catalogEdition := &domain.CatalogEdition{
		ID:           editionID,
		NewsletterID: cmd.NewsletterID,
		Subject:      cmd.Subject,
		ContentHTML:  cmd.ContentHTML,
		ContentText:  cmd.ContentText,
		TemplateData: cmd.TemplateData,
		ScheduledAt:  cmd.ScheduledAt,
		Status:       edition.Status.String(),
		CreatedBy:    userID,
		CreatedAt:    edition.CreatedAt,
		UpdatedAt:    edition.UpdatedAt,
	}

	err = h.catalog.Add(ctx, catalogEdition)
	if err != nil {
		return "", err
	}

	err = h.publisher.Publish(ctx, event)
	if err != nil {
		return "", err
	}

	return editionID, nil
}