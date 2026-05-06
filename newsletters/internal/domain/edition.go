package domain

import (
	"time"

	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const EditionAggregate = "newsletters.Edition"

type Edition struct {
	es.Aggregate
	NewsletterID   string
	Subject        string
	ContentHTML    string
	ContentText    string
	TemplateData   map[string]string
	ScheduledAt    *time.Time
	SentAt         *time.Time
	Status         EditionStatus
	CreatedBy      string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	RecipientCount int
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Edition)(nil)

func NewEdition(id string) *Edition {
	return &Edition{
		Aggregate: es.NewAggregate(id, EditionAggregate),
	}
}

func (e *Edition) CreateEdition(newsletterID, subject, contentHTML, contentText, createdBy string, templateData map[string]string, scheduledAt *time.Time) (ddd.Event, error) {
	if newsletterID == "" {
		return nil, errors.ErrBadRequest.Msg("newsletter ID is required")
	}
	if subject == "" {
		return nil, errors.ErrBadRequest.Msg("subject is required")
	}
	if contentHTML == "" && contentText == "" {
		return nil, errors.ErrBadRequest.Msg("content is required")
	}
	if createdBy == "" {
		return nil, errors.ErrBadRequest.Msg("creator ID is required")
	}

	var schedAt int64
	if scheduledAt != nil {
		schedAt = scheduledAt.Unix()
	}

	e.AddEvent(EditionCreatedEvent, &EditionCreated{
		NewsletterID: newsletterID,
		Subject:      subject,
		ContentHTML:  contentHTML,
		ContentText:  contentText,
		TemplateData: templateData,
		ScheduledAt:  schedAt,
		CreatedBy:    createdBy,
	})

	return ddd.NewEvent(EditionCreatedEvent, e), nil
}

func (e *Edition) UpdateEdition(subject, contentHTML, contentText string, templateData map[string]string, scheduledAt *time.Time) (ddd.Event, error) {
	if e.Status != DraftStatus {
		return nil, errors.ErrBadRequest.Msg("can only update draft editions")
	}
	if subject == "" {
		return nil, errors.ErrBadRequest.Msg("subject is required")
	}

	var schedAt int64
	if scheduledAt != nil {
		schedAt = scheduledAt.Unix()
	}

	e.AddEvent(EditionUpdatedEvent, &EditionUpdated{
		Subject:      subject,
		ContentHTML:  contentHTML,
		ContentText:  contentText,
		TemplateData: templateData,
		ScheduledAt:  schedAt,
	})

	return ddd.NewEvent(EditionUpdatedEvent, e), nil
}

func (e *Edition) ScheduleEdition(scheduledAt time.Time) (ddd.Event, error) {
	if e.Status != DraftStatus {
		return nil, errors.ErrBadRequest.Msg("can only schedule draft editions")
	}
	if scheduledAt.Before(time.Now()) {
		return nil, errors.ErrBadRequest.Msg("scheduled time must be in the future")
	}

	e.AddEvent(EditionScheduledEvent, &EditionScheduled{
		ScheduledAt: scheduledAt.Unix(),
	})

	return ddd.NewEvent(EditionScheduledEvent, e), nil
}

func (e *Edition) SendEdition(recipientCount int) (ddd.Event, error) {
	if e.Status != DraftStatus && e.Status != ScheduledStatus {
		return nil, errors.ErrBadRequest.Msg("can only send draft or scheduled editions")
	}

	e.AddEvent(EditionSentEvent, &EditionSent{
		EditionID:      e.ID(),
		RecipientCount: recipientCount,
	})

	return ddd.NewEvent(EditionSentEvent, e), nil
}

func (e *Edition) MarkSending() (ddd.Event, error) {
	if e.Status != DraftStatus && e.Status != ScheduledStatus {
		return nil, errors.ErrBadRequest.Msg("invalid status for sending")
	}

	e.AddEvent(EditionSendingEvent, &EditionSending{})
	return ddd.NewEvent(EditionSendingEvent, e), nil
}

// Key implements registry.Registerable
func (Edition) Key() string { return EditionAggregate }

// ApplyEvent implements es.EventApplier
func (e *Edition) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *EditionCreated:
		e.NewsletterID = payload.NewsletterID
		e.Subject = payload.Subject
		e.ContentHTML = payload.ContentHTML
		e.ContentText = payload.ContentText
		e.TemplateData = payload.TemplateData
		if payload.ScheduledAt > 0 {
			schedAt := time.Unix(payload.ScheduledAt, 0)
			e.ScheduledAt = &schedAt
			e.Status = ScheduledStatus
		} else {
			e.Status = DraftStatus
		}
		e.CreatedBy = payload.CreatedBy
		e.CreatedAt = event.OccurredAt()
		e.UpdatedAt = event.OccurredAt()

	case *EditionUpdated:
		e.Subject = payload.Subject
		e.ContentHTML = payload.ContentHTML
		e.ContentText = payload.ContentText
		e.TemplateData = payload.TemplateData
		if payload.ScheduledAt > 0 {
			schedAt := time.Unix(payload.ScheduledAt, 0)
			e.ScheduledAt = &schedAt
			e.Status = ScheduledStatus
		}
		e.UpdatedAt = event.OccurredAt()

	case *EditionScheduled:
		schedAt := time.Unix(payload.ScheduledAt, 0)
		e.ScheduledAt = &schedAt
		e.Status = ScheduledStatus
		e.UpdatedAt = event.OccurredAt()

	case *EditionSending:
		e.Status = SendingStatus
		e.UpdatedAt = event.OccurredAt()

	case *EditionSent:
		e.Status = SentStatus
		sentAt := event.OccurredAt()
		e.SentAt = &sentAt
		e.RecipientCount = payload.RecipientCount
		e.UpdatedAt = event.OccurredAt()

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", e, event.EventName(), payload)
	}

	return nil
}

// ApplySnapshot implements es.Snapshotter
func (e *Edition) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *EditionV1:
		e.NewsletterID = ss.NewsletterID
		e.Subject = ss.Subject
		e.ContentHTML = ss.ContentHTML
		e.ContentText = ss.ContentText
		e.TemplateData = ss.TemplateData
		e.ScheduledAt = ss.ScheduledAt
		e.SentAt = ss.SentAt
		e.Status = ToEditionStatus(ss.Status)
		e.CreatedBy = ss.CreatedBy
		e.CreatedAt = ss.CreatedAt
		e.UpdatedAt = ss.UpdatedAt
		e.RecipientCount = ss.RecipientCount

	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", e, snapshot)
	}

	return nil
}

// ToSnapshot implements es.Snapshotter
func (e Edition) ToSnapshot() es.Snapshot {
	return EditionV1{
		NewsletterID:   e.NewsletterID,
		Subject:        e.Subject,
		ContentHTML:    e.ContentHTML,
		ContentText:    e.ContentText,
		TemplateData:   e.TemplateData,
		ScheduledAt:    e.ScheduledAt,
		SentAt:         e.SentAt,
		Status:         e.Status.String(),
		CreatedBy:      e.CreatedBy,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
		RecipientCount: e.RecipientCount,
	}
}