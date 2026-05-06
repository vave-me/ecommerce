package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/auth"
	"middleman/internal/ddd"
	"middleman/newsletters/internal/domain"
)

type DeleteTemplate struct {
	ID string
}

type DeleteTemplateHandler struct {
	templates domain.TemplateRepository
	catalog   domain.TemplateCatalogRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewDeleteTemplateHandler(
	templates domain.TemplateRepository,
	catalog domain.TemplateCatalogRepository,
	publisher ddd.EventPublisher[ddd.Event],
) DeleteTemplateHandler {
	return DeleteTemplateHandler{
		templates: templates,
		catalog:   catalog,
		publisher: publisher,
	}
}

func (h DeleteTemplateHandler) DeleteTemplate(ctx context.Context, cmd DeleteTemplate) error {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return errors.ErrUnauthenticated.Msg("authentication required")
	}
	userID := claims.Subject
	
	template, err := h.templates.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	// Only owner can delete template
	if template.UserID != userID {
		return errors.ErrForbidden.Msg("not authorized to delete this template")
	}

	event, err := template.DeleteTemplate()
	if err != nil {
		return err
	}

	err = h.templates.Save(ctx, template)
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