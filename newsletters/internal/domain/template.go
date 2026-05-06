package domain

import (
	"time"

	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const TemplateAggregate = "newsletters.Template"

type Template struct {
	es.Aggregate
	UserID       string // Empty for system templates
	Name         string
	Description  string
	HTMLTemplate string
	TextTemplate string
	Variables    map[string]string
	PreviewData  map[string]string
	IsPublic     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Template)(nil)

func NewTemplate(id string) *Template {
	return &Template{
		Aggregate: es.NewAggregate(id, TemplateAggregate),
	}
}

func (t *Template) CreateTemplate(userID, name, description, htmlTemplate, textTemplate string, variables, previewData map[string]string, isPublic bool) (ddd.Event, error) {
	if name == "" {
		return nil, errors.ErrBadRequest.Msg("template name is required")
	}
	if htmlTemplate == "" && textTemplate == "" {
		return nil, errors.ErrBadRequest.Msg("at least one template format is required")
	}

	t.AddEvent(TemplateCreatedEvent, &TemplateCreated{
		UserID:       userID,
		Name:         name,
		Description:  description,
		HTMLTemplate: htmlTemplate,
		TextTemplate: textTemplate,
		Variables:    variables,
		PreviewData:  previewData,
		IsPublic:     isPublic,
	})

	return ddd.NewEvent(TemplateCreatedEvent, t), nil
}

func (t *Template) UpdateTemplate(name, description, htmlTemplate, textTemplate string, variables, previewData map[string]string, isPublic bool) (ddd.Event, error) {
	if name == "" {
		return nil, errors.ErrBadRequest.Msg("template name is required")
	}
	if htmlTemplate == "" && textTemplate == "" {
		return nil, errors.ErrBadRequest.Msg("at least one template format is required")
	}

	t.AddEvent(TemplateUpdatedEvent, &TemplateUpdated{
		Name:         name,
		Description:  description,
		HTMLTemplate: htmlTemplate,
		TextTemplate: textTemplate,
		Variables:    variables,
		PreviewData:  previewData,
		IsPublic:     isPublic,
	})

	return ddd.NewEvent(TemplateUpdatedEvent, t), nil
}

func (t *Template) DeleteTemplate() (ddd.Event, error) {
	t.AddEvent(TemplateDeletedEvent, &TemplateDeleted{})
	return ddd.NewEvent(TemplateDeletedEvent, t), nil
}

// Key implements registry.Registerable
func (Template) Key() string { return TemplateAggregate }

// ApplyEvent implements es.EventApplier
func (t *Template) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *TemplateCreated:
		t.UserID = payload.UserID
		t.Name = payload.Name
		t.Description = payload.Description
		t.HTMLTemplate = payload.HTMLTemplate
		t.TextTemplate = payload.TextTemplate
		t.Variables = payload.Variables
		t.PreviewData = payload.PreviewData
		t.IsPublic = payload.IsPublic
		t.CreatedAt = event.OccurredAt()
		t.UpdatedAt = event.OccurredAt()

	case *TemplateUpdated:
		t.Name = payload.Name
		t.Description = payload.Description
		t.HTMLTemplate = payload.HTMLTemplate
		t.TextTemplate = payload.TextTemplate
		t.Variables = payload.Variables
		t.PreviewData = payload.PreviewData
		t.IsPublic = payload.IsPublic
		t.UpdatedAt = event.OccurredAt()

	case *TemplateDeleted:
		t.UpdatedAt = event.OccurredAt()

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", t, event.EventName(), payload)
	}

	return nil
}

// ApplySnapshot implements es.Snapshotter
func (t *Template) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *TemplateV1:
		t.UserID = ss.UserID
		t.Name = ss.Name
		t.Description = ss.Description
		t.HTMLTemplate = ss.HTMLTemplate
		t.TextTemplate = ss.TextTemplate
		t.Variables = ss.Variables
		t.PreviewData = ss.PreviewData
		t.IsPublic = ss.IsPublic
		t.CreatedAt = ss.CreatedAt
		t.UpdatedAt = ss.UpdatedAt

	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", t, snapshot)
	}

	return nil
}

// ToSnapshot implements es.Snapshotter
func (t Template) ToSnapshot() es.Snapshot {
	return TemplateV1{
		UserID:       t.UserID,
		Name:         t.Name,
		Description:  t.Description,
		HTMLTemplate: t.HTMLTemplate,
		TextTemplate: t.TextTemplate,
		Variables:    t.Variables,
		PreviewData:  t.PreviewData,
		IsPublic:     t.IsPublic,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}
}