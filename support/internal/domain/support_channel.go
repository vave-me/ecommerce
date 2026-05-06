package domain

import (
	"middleman/internal/ddd"
	"middleman/internal/es"
	"time"

	"github.com/stackus/errors"
)

const SupportChannelAggregate = "support.SupportChannel"

var (
	ErrChannelAlreadyExists = errors.Wrap(errors.ErrBadRequest, "support channel already exists")
	ErrChannelNotActive     = errors.Wrap(errors.ErrBadRequest, "support channel is not active")
	ErrChannelAlreadyClosed = errors.Wrap(errors.ErrBadRequest, "support channel is already closed")
)

type SupportChannel struct {
	es.Aggregate
	UserID       string
	BusinessID   string
	ChannelType  SupportChannelType
	Active       bool
	Settings     SupportChannelSettings
	OpenTickets  int
	TotalTickets int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ClosedAt     *time.Time
}

type SupportChannelType string

const (
	ChannelTypeGeneral   SupportChannelType = "GENERAL"
	ChannelTypeTechnical SupportChannelType = "TECHNICAL"
	ChannelTypeBilling   SupportChannelType = "BILLING"
	ChannelTypeSales     SupportChannelType = "SALES"
	ChannelTypeProduct   SupportChannelType = "PRODUCT"
)

type SupportChannelSettings struct {
	EmailNotifications  bool              `json:"email_notifications"`
	SMSNotifications    bool              `json:"sms_notifications"`
	AutoAssignTickets   bool              `json:"auto_assign_tickets"`
	PreferredLanguage   string            `json:"preferred_language"`
	Timezone            string            `json:"timezone"`
	NotificationEmails  []string          `json:"notification_emails"`
	SLASettings         SLASettings       `json:"sla_settings"`
}

type SLASettings struct {
	FirstResponseMinutes  int            `json:"first_response_minutes"`
	ResolutionHours       int            `json:"resolution_hours"`
	PriorityResponseTimes map[string]int `json:"priority_response_times"`
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*SupportChannel)(nil)

func NewSupportChannel(id string) *SupportChannel {
	return &SupportChannel{
		Aggregate: es.NewAggregate(id, SupportChannelAggregate),
	}
}

func (SupportChannel) Key() string { return SupportChannelAggregate }

func (c *SupportChannel) InitChannel(
	userID, businessID string,
	channelType SupportChannelType,
	settings SupportChannelSettings,
) (ddd.Event, error) {
	if c.UserID != "" {
		return nil, ErrChannelAlreadyExists
	}

	c.AddEvent(SupportChannelCreatedEvent, &SupportChannelCreated{
		UserID:      userID,
		BusinessID:  businessID,
		ChannelType: channelType,
		Settings:    settings,
	})

	return ddd.NewEvent(SupportChannelCreatedEvent, c), nil
}

func (c *SupportChannel) UpdateSettings(settings SupportChannelSettings) (ddd.Event, error) {
	if !c.Active {
		return nil, ErrChannelNotActive
	}

	c.AddEvent(SupportChannelSettingsUpdatedEvent, &SupportChannelSettingsUpdated{
		Settings: settings,
	})

	return ddd.NewEvent(SupportChannelSettingsUpdatedEvent, c), nil
}

func (c *SupportChannel) CloseChannel(closedBy, reason string) (ddd.Event, error) {
	if !c.Active {
		return nil, ErrChannelAlreadyClosed
	}

	c.AddEvent(SupportChannelClosedEvent, &SupportChannelClosed{
		ClosedBy: closedBy,
		Reason:   reason,
	})

	return ddd.NewEvent(SupportChannelClosedEvent, c), nil
}

func (c *SupportChannel) ReactivateChannel(reactivatedBy string) (ddd.Event, error) {
	if c.Active {
		return nil, errors.Wrap(errors.ErrBadRequest, "support channel is already active")
	}

	c.AddEvent(SupportChannelReactivatedEvent, &SupportChannelReactivated{
		ReactivatedBy: reactivatedBy,
	})

	return ddd.NewEvent(SupportChannelReactivatedEvent, c), nil
}

func (c *SupportChannel) IncrementTicketCount() {
	c.OpenTickets++
	c.TotalTickets++
}

func (c *SupportChannel) DecrementOpenTickets() {
	if c.OpenTickets > 0 {
		c.OpenTickets--
	}
}

// ApplyEvent implements es.EventApplier
func (c *SupportChannel) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *SupportChannelCreated:
		c.UserID = payload.UserID
		c.BusinessID = payload.BusinessID
		c.ChannelType = payload.ChannelType
		c.Settings = payload.Settings
		c.Active = true
		c.OpenTickets = 0
		c.TotalTickets = 0
		c.CreatedAt = event.OccurredAt()
		c.UpdatedAt = event.OccurredAt()

	case *SupportChannelSettingsUpdated:
		c.Settings = payload.Settings
		c.UpdatedAt = event.OccurredAt()

	case *SupportChannelClosed:
		c.Active = false
		closedAt := event.OccurredAt()
		c.ClosedAt = &closedAt
		c.UpdatedAt = event.OccurredAt()

	case *SupportChannelReactivated:
		c.Active = true
		c.ClosedAt = nil
		c.UpdatedAt = event.OccurredAt()

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", c, event.EventName(), payload)
	}

	return nil
}

// ApplySnapshot implements es.Snapshotter
func (c *SupportChannel) ApplySnapshot(snapshot es.Snapshot) error {
	switch snap := snapshot.(type) {
	case *SupportChannelV1:
		c.UserID = snap.UserID
		c.BusinessID = snap.BusinessID
		c.ChannelType = snap.ChannelType
		c.Active = snap.Active
		c.Settings = snap.Settings
		c.OpenTickets = snap.OpenTickets
		c.TotalTickets = snap.TotalTickets
		c.CreatedAt = snap.CreatedAt
		c.UpdatedAt = snap.UpdatedAt
		c.ClosedAt = snap.ClosedAt

	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", c, snapshot)
	}

	return nil
}

// ToSnapshot implements es.Snapshotter
func (c SupportChannel) ToSnapshot() es.Snapshot {
	return &SupportChannelV1{
		UserID:       c.UserID,
		BusinessID:   c.BusinessID,
		ChannelType:  c.ChannelType,
		Active:       c.Active,
		Settings:     c.Settings,
		OpenTickets:  c.OpenTickets,
		TotalTickets: c.TotalTickets,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
		ClosedAt:     c.ClosedAt,
	}
}