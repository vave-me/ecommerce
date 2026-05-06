package commands

import (
	"context"
	"github.com/google/uuid"
	"github.com/stackus/errors"
	"middleman/internal/auth"
	"middleman/internal/ddd"
	"middleman/newsletters/internal/domain"
)

type CreateTemplate struct {
	Name         string
	Description  string
	HTMLTemplate string
	TextTemplate string
	Variables    map[string]string
	PreviewData  map[string]string
	IsPublic     bool
}

type CreateTemplateHandler struct {
	templates domain.TemplateRepository
	catalog   domain.TemplateCatalogRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewCreateTemplateHandler(
	templates domain.TemplateRepository,
	catalog domain.TemplateCatalogRepository,
	publisher ddd.EventPublisher[ddd.Event],
) CreateTemplateHandler {
	return CreateTemplateHandler{
		templates: templates,
		catalog:   catalog,
		publisher: publisher,
	}
}

func (h CreateTemplateHandler) CreateTemplate(ctx context.Context, cmd CreateTemplate) (string, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return "", errors.ErrUnauthenticated.Msg("authentication required")
	}
	userID := claims.Subject
	templateID := uuid.New().String()
	
	template := domain.NewTemplate(templateID)
	
	event, err := template.CreateTemplate(
		userID,
		cmd.Name,
		cmd.Description,
		cmd.HTMLTemplate,
		cmd.TextTemplate,
		cmd.Variables,
		cmd.PreviewData,
		cmd.IsPublic,
	)
	if err != nil {
		return "", err
	}

	err = h.templates.Save(ctx, template)
	if err != nil {
		return "", err
	}

	// Add to catalog
	catalogTemplate := &domain.CatalogTemplate{
		ID:           templateID,
		UserID:       userID,
		Name:         cmd.Name,
		Description:  cmd.Description,
		HTMLTemplate: cmd.HTMLTemplate,
		TextTemplate: cmd.TextTemplate,
		Variables:    cmd.Variables,
		PreviewData:  cmd.PreviewData,
		IsPublic:     cmd.IsPublic,
		CreatedAt:    template.CreatedAt,
		UpdatedAt:    template.UpdatedAt,
	}
	
	err = h.catalog.Add(ctx, catalogTemplate)
	if err != nil {
		return "", err
	}

	err = h.publisher.Publish(ctx, event)
	if err != nil {
		return "", err
	}

	return templateID, nil
}