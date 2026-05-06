package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/auth"
	"middleman/internal/ddd"
	"middleman/newsletters/internal/domain"
)

type UpdateTemplate struct {
	ID           string
	Name         string
	Description  string
	HTMLTemplate string
	TextTemplate string
	Variables    map[string]string
	PreviewData  map[string]string
	IsPublic     bool
}

type UpdateTemplateHandler struct {
	templates domain.TemplateRepository
	catalog   domain.TemplateCatalogRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewUpdateTemplateHandler(
	templates domain.TemplateRepository,
	catalog domain.TemplateCatalogRepository,
	publisher ddd.EventPublisher[ddd.Event],
) UpdateTemplateHandler {
	return UpdateTemplateHandler{
		templates: templates,
		catalog:   catalog,
		publisher: publisher,
	}
}

func (h UpdateTemplateHandler) UpdateTemplate(ctx context.Context, cmd UpdateTemplate) error {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return errors.ErrUnauthenticated.Msg("authentication required")
	}
	userID := claims.Subject
	
	template, err := h.templates.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	// Only owner can update template
	if template.UserID != userID {
		return errors.ErrForbidden.Msg("not authorized to update this template")
	}

	event, err := template.UpdateTemplate(
		cmd.Name,
		cmd.Description,
		cmd.HTMLTemplate,
		cmd.TextTemplate,
		cmd.Variables,
		cmd.PreviewData,
		cmd.IsPublic,
	)
	if err != nil {
		return err
	}

	err = h.templates.Save(ctx, template)
	if err != nil {
		return err
	}

	// Update catalog
	catalogTemplate, err := h.catalog.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	catalogTemplate.Name = cmd.Name
	catalogTemplate.Description = cmd.Description
	catalogTemplate.HTMLTemplate = cmd.HTMLTemplate
	catalogTemplate.TextTemplate = cmd.TextTemplate
	catalogTemplate.Variables = cmd.Variables
	catalogTemplate.PreviewData = cmd.PreviewData
	catalogTemplate.IsPublic = cmd.IsPublic
	catalogTemplate.UpdatedAt = template.UpdatedAt

	err = h.catalog.Update(ctx, catalogTemplate)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}