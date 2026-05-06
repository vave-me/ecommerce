package domain

import (
	"time"

	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const NewsletterAggregate = "newsletters.Newsletter"

type Newsletter struct {
	es.Aggregate
	UserID      string
	Name        string
	Description string
	Frequency   NewsletterFrequency
	Category    string
	TemplateID  string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Newsletter)(nil)

func NewNewsletter(id string) *Newsletter {
	return &Newsletter{
		Aggregate: es.NewAggregate(id, NewsletterAggregate),
	}
}

func (n *Newsletter) CreateNewsletter(userID, name, description, frequency, category, templateID string) (ddd.Event, error) {
	if userID == "" {
		return nil, errors.ErrBadRequest.Msg("user ID is required")
	}
	if name == "" {
		return nil, errors.ErrBadRequest.Msg("newsletter name is required")
	}
	
	freq := ToNewsletterFrequency(frequency)
	if freq == UnknownFrequency {
		return nil, errors.ErrBadRequest.Msg("invalid frequency")
	}

	n.AddEvent(NewsletterCreatedEvent, &NewsletterCreated{
		UserID:      userID,
		Name:        name,
		Description: description,
		Frequency:   frequency,
		Category:    category,
		TemplateID:  templateID,
	})

	return ddd.NewEvent(NewsletterCreatedEvent, n), nil
}

func (n *Newsletter) UpdateNewsletter(name, description, frequency, category, templateID string) (ddd.Event, error) {
	if !n.IsActive {
		return nil, errors.ErrBadRequest.Msg("cannot update inactive newsletter")
	}

	freq := ToNewsletterFrequency(frequency)
	if freq == UnknownFrequency {
		return nil, errors.ErrBadRequest.Msg("invalid frequency")
	}

	n.AddEvent(NewsletterUpdatedEvent, &NewsletterUpdated{
		Name:        name,
		Description: description,
		Frequency:   frequency,
		Category:    category,
		TemplateID:  templateID,
	})

	return ddd.NewEvent(NewsletterUpdatedEvent, n), nil
}

func (n *Newsletter) ActivateNewsletter() (ddd.Event, error) {
	if n.IsActive {
		return nil, errors.ErrBadRequest.Msg("newsletter is already active")
	}

	n.AddEvent(NewsletterActivatedEvent, &NewsletterActivated{})
	return ddd.NewEvent(NewsletterActivatedEvent, n), nil
}

func (n *Newsletter) DeactivateNewsletter() (ddd.Event, error) {
	if !n.IsActive {
		return nil, errors.ErrBadRequest.Msg("newsletter is already inactive")
	}

	n.AddEvent(NewsletterDeactivatedEvent, &NewsletterDeactivated{})
	return ddd.NewEvent(NewsletterDeactivatedEvent, n), nil
}

func (n *Newsletter) DeleteNewsletter() (ddd.Event, error) {
	n.AddEvent(NewsletterDeletedEvent, &NewsletterDeleted{})
	return ddd.NewEvent(NewsletterDeletedEvent, n), nil
}

// Key implements registry.Registerable
func (Newsletter) Key() string { return NewsletterAggregate }

// ApplyEvent implements es.EventApplier
func (n *Newsletter) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *NewsletterCreated:
		n.UserID = payload.UserID
		n.Name = payload.Name
		n.Description = payload.Description
		n.Frequency = ToNewsletterFrequency(payload.Frequency)
		n.Category = payload.Category
		n.TemplateID = payload.TemplateID
		n.IsActive = true
		n.CreatedAt = event.OccurredAt()
		n.UpdatedAt = event.OccurredAt()

	case *NewsletterUpdated:
		n.Name = payload.Name
		n.Description = payload.Description
		n.Frequency = ToNewsletterFrequency(payload.Frequency)
		n.Category = payload.Category
		n.TemplateID = payload.TemplateID
		n.UpdatedAt = event.OccurredAt()

	case *NewsletterActivated:
		n.IsActive = true
		n.UpdatedAt = event.OccurredAt()

	case *NewsletterDeactivated:
		n.IsActive = false
		n.UpdatedAt = event.OccurredAt()

	case *NewsletterDeleted:
		n.IsActive = false
		n.UpdatedAt = event.OccurredAt()

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", n, event.EventName(), payload)
	}

	return nil
}

// ApplySnapshot implements es.Snapshotter
func (n *Newsletter) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *NewsletterV1:
		n.UserID = ss.UserID
		n.Name = ss.Name
		n.Description = ss.Description
		n.Frequency = ToNewsletterFrequency(ss.Frequency)
		n.Category = ss.Category
		n.TemplateID = ss.TemplateID
		n.IsActive = ss.IsActive
		n.CreatedAt = ss.CreatedAt
		n.UpdatedAt = ss.UpdatedAt

	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", n, snapshot)
	}

	return nil
}

// ToSnapshot implements es.Snapshotter
func (n Newsletter) ToSnapshot() es.Snapshot {
	return NewsletterV1{
		UserID:      n.UserID,
		Name:        n.Name,
		Description: n.Description,
		Frequency:   n.Frequency.String(),
		Category:    n.Category,
		TemplateID:  n.TemplateID,
		IsActive:    n.IsActive,
		CreatedAt:   n.CreatedAt,
		UpdatedAt:   n.UpdatedAt,
	}
}