package domain

import (
	"middleman/internal/ddd"
	"middleman/internal/es"
	"time"

	"github.com/stackus/errors"
)

const TicketAggregate = "support.Ticket"

var (
	ErrTicketNotFound          = errors.Wrap(errors.ErrNotFound, "ticket not found")
	ErrInvalidTicketStatus     = errors.Wrap(errors.ErrBadRequest, "invalid ticket status transition")
	ErrTicketAlreadyAssigned   = errors.Wrap(errors.ErrBadRequest, "ticket already assigned")
	ErrTicketAlreadyResolved   = errors.Wrap(errors.ErrBadRequest, "ticket already resolved")
	ErrTicketAlreadyClosed     = errors.Wrap(errors.ErrBadRequest, "ticket already closed")
	ErrCannotReopenClosedTicket = errors.Wrap(errors.ErrBadRequest, "cannot reopen a closed ticket")
)

type Ticket struct {
	es.Aggregate
	ChannelID           string
	Title               string
	Description         string
	Status              TicketStatus
	Priority            TicketPriority
	Category            TicketCategory
	Tags                []string
	Metadata            map[string]string
	AssigneeID          string
	AssigneeType        AssigneeType
	CreatedBy           string
	CurrentTier         SupportTier
	ResponseCount       int
	ReopenCount         int
	SatisfactionRating  *CustomerSatisfaction
	LinkedTicketIDs     []string
	MergedTicketIDs     []string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	ResolvedAt          *time.Time
	ClosedAt            *time.Time
	FirstResponseAt     *time.Time
	Attachments         []Attachment
}

type TicketStatus string

const (
	TicketStatusDraft           TicketStatus = "DRAFT"
	TicketStatusSubmitted       TicketStatus = "SUBMITTED"
	TicketStatusAssigned        TicketStatus = "ASSIGNED"
	TicketStatusInProgress      TicketStatus = "IN_PROGRESS"
	TicketStatusPendingCustomer TicketStatus = "PENDING_CUSTOMER"
	TicketStatusResolved        TicketStatus = "RESOLVED"
	TicketStatusClosed          TicketStatus = "CLOSED"
	TicketStatusEscalated       TicketStatus = "ESCALATED"
	TicketStatusReopened        TicketStatus = "REOPENED"
)

type TicketPriority string

const (
	PriorityLow      TicketPriority = "LOW"
	PriorityMedium   TicketPriority = "MEDIUM"
	PriorityHigh     TicketPriority = "HIGH"
	PriorityUrgent   TicketPriority = "URGENT"
	PriorityCritical TicketPriority = "CRITICAL"
)

type TicketCategory string

const (
	CategoryGeneralInquiry  TicketCategory = "GENERAL_INQUIRY"
	CategoryTechnicalIssue  TicketCategory = "TECHNICAL_ISSUE"
	CategoryBillingIssue    TicketCategory = "BILLING_ISSUE"
	CategoryAccountIssue    TicketCategory = "ACCOUNT_ISSUE"
	CategoryProductQuestion TicketCategory = "PRODUCT_QUESTION"
	CategoryFeatureRequest  TicketCategory = "FEATURE_REQUEST"
	CategoryComplaint       TicketCategory = "COMPLAINT"
	CategoryRefundRequest   TicketCategory = "REFUND_REQUEST"
	CategoryOrderIssue      TicketCategory = "ORDER_ISSUE"
	CategoryShippingIssue   TicketCategory = "SHIPPING_ISSUE"
)

type SupportTier string

const (
	TierOne        SupportTier = "TIER_1"
	TierTwo        SupportTier = "TIER_2"
	TierThree      SupportTier = "TIER_3"
	TierManagement SupportTier = "MANAGEMENT"
)

type AssigneeType string

const (
	AssigneeTypeHuman AssigneeType = "HUMAN_AGENT"
	AssigneeTypeAI    AssigneeType = "AI_ASSISTANT"
	AssigneeTypeTeam  AssigneeType = "TEAM"
)

type CustomerSatisfaction string

const (
	SatisfactionVeryDissatisfied CustomerSatisfaction = "VERY_DISSATISFIED"
	SatisfactionDissatisfied     CustomerSatisfaction = "DISSATISFIED"
	SatisfactionNeutral          CustomerSatisfaction = "NEUTRAL"
	SatisfactionSatisfied        CustomerSatisfaction = "SATISFIED"
	SatisfactionVerySatisfied    CustomerSatisfaction = "VERY_SATISFIED"
)

type Attachment struct {
	ID          string    `json:"id"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"content_type"`
	SizeBytes   int64     `json:"size_bytes"`
	URL         string    `json:"url"`
	UploadedAt  time.Time `json:"uploaded_at"`
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Ticket)(nil)

func NewTicket(id string) *Ticket {
	return &Ticket{
		Aggregate: es.NewAggregate(id, TicketAggregate),
	}
}

func (Ticket) Key() string { return TicketAggregate }

func (t *Ticket) InitTicket(
	channelID, title, description string,
	category TicketCategory,
	priority TicketPriority,
	tags []string,
	metadata map[string]string,
	createdBy string,
	attachments []Attachment,
) (ddd.Event, error) {
	if t.ChannelID != "" {
		return nil, errors.Wrap(errors.ErrBadRequest, "ticket already initialized")
	}

	if title == "" {
		return nil, errors.Wrap(errors.ErrBadRequest, "ticket title cannot be empty")
	}

	t.AddEvent(TicketCreatedEvent, &TicketCreated{
		ChannelID:   channelID,
		Title:       title,
		Description: description,
		Category:    category,
		Priority:    priority,
		Tags:        tags,
		Metadata:    metadata,
		CreatedBy:   createdBy,
		Attachments: attachments,
	})

	return ddd.NewEvent(TicketCreatedEvent, t), nil
}

func (t *Ticket) UpdateTicket(
	title, description *string,
	category *TicketCategory,
	tags []string,
	metadata map[string]string,
	updatedBy string,
) (ddd.Event, error) {
	if t.Status == TicketStatusClosed {
		return nil, errors.Wrap(errors.ErrBadRequest, "cannot update closed ticket")
	}

	t.AddEvent(TicketUpdatedEvent, &TicketUpdated{
		Title:       title,
		Description: description,
		Category:    category,
		Tags:        tags,
		Metadata:    metadata,
		UpdatedBy:   updatedBy,
	})

	return ddd.NewEvent(TicketUpdatedEvent, t), nil
}

func (t *Ticket) AssignTicket(assigneeID string, assigneeType AssigneeType, assignedBy, reason string) (ddd.Event, error) {
	if t.Status == TicketStatusClosed {
		return nil, errors.Wrap(errors.ErrBadRequest, "cannot assign closed ticket")
	}

	t.AddEvent(TicketAssignedEvent, &TicketAssigned{
		AssigneeID:        assigneeID,
		AssigneeType:      assigneeType,
		AssignedBy:        assignedBy,
		AssignmentReason:  reason,
	})

	return ddd.NewEvent(TicketAssignedEvent, t), nil
}

func (t *Ticket) UpdatePriority(newPriority TicketPriority, updatedBy, reason string) (ddd.Event, error) {
	if t.Priority == newPriority {
		return nil, errors.Wrap(errors.ErrBadRequest, "priority unchanged")
	}

	t.AddEvent(TicketPriorityUpdatedEvent, &TicketPriorityUpdated{
		OldPriority: t.Priority,
		NewPriority: newPriority,
		UpdatedBy:   updatedBy,
		Reason:      reason,
	})

	return ddd.NewEvent(TicketPriorityUpdatedEvent, t), nil
}

func (t *Ticket) EscalateTicket(toTier SupportTier, escalatedBy, reason, notes string) (ddd.Event, error) {
	if t.Status == TicketStatusClosed {
		return nil, errors.Wrap(errors.ErrBadRequest, "cannot escalate closed ticket")
	}

	t.AddEvent(TicketEscalatedEvent, &TicketEscalated{
		FromTier:         t.CurrentTier,
		ToTier:           toTier,
		EscalatedBy:      escalatedBy,
		EscalationReason: reason,
		EscalationNotes:  notes,
	})

	return ddd.NewEvent(TicketEscalatedEvent, t), nil
}

func (t *Ticket) ResolveTicket(resolvedBy, resolution string, appliedSolutions []string) (ddd.Event, error) {
	if t.Status == TicketStatusResolved || t.Status == TicketStatusClosed {
		return nil, ErrTicketAlreadyResolved
	}

	t.AddEvent(TicketResolvedEvent, &TicketResolved{
		ResolvedBy:       resolvedBy,
		Resolution:       resolution,
		AppliedSolutions: appliedSolutions,
	})

	return ddd.NewEvent(TicketResolvedEvent, t), nil
}

func (t *Ticket) ReopenTicket(reopenedBy, reason string) (ddd.Event, error) {
	if t.Status == TicketStatusClosed {
		return nil, ErrCannotReopenClosedTicket
	}

	if t.Status != TicketStatusResolved {
		return nil, errors.Wrap(errors.ErrBadRequest, "can only reopen resolved tickets")
	}

	t.AddEvent(TicketReopenedEvent, &TicketReopened{
		ReopenedBy:   reopenedBy,
		ReopenReason: reason,
		ReopenCount:  t.ReopenCount + 1,
	})

	return ddd.NewEvent(TicketReopenedEvent, t), nil
}

func (t *Ticket) CloseTicket(closedBy, notes string, satisfaction *CustomerSatisfaction) (ddd.Event, error) {
	if t.Status == TicketStatusClosed {
		return nil, ErrTicketAlreadyClosed
	}

	if t.Status != TicketStatusResolved {
		return nil, errors.Wrap(errors.ErrBadRequest, "can only close resolved tickets")
	}

	t.AddEvent(TicketClosedEvent, &TicketClosed{
		ClosedBy:           closedBy,
		ClosureNotes:       notes,
		SatisfactionRating: satisfaction,
	})

	return ddd.NewEvent(TicketClosedEvent, t), nil
}

// ApplyEvent implements es.EventApplier
func (t *Ticket) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *TicketCreated:
		t.ChannelID = payload.ChannelID
		t.Title = payload.Title
		t.Description = payload.Description
		t.Category = payload.Category
		t.Priority = payload.Priority
		t.Tags = payload.Tags
		t.Metadata = payload.Metadata
		t.CreatedBy = payload.CreatedBy
		t.Attachments = payload.Attachments
		t.Status = TicketStatusSubmitted
		t.CurrentTier = TierOne
		t.ResponseCount = 0
		t.ReopenCount = 0
		t.CreatedAt = event.OccurredAt()
		t.UpdatedAt = event.OccurredAt()

	case *TicketUpdated:
		if payload.Title != nil {
			t.Title = *payload.Title
		}
		if payload.Description != nil {
			t.Description = *payload.Description
		}
		if payload.Category != nil {
			t.Category = *payload.Category
		}
		if payload.Tags != nil {
			t.Tags = payload.Tags
		}
		if payload.Metadata != nil {
			t.Metadata = payload.Metadata
		}
		t.UpdatedAt = event.OccurredAt()

	case *TicketAssigned:
		t.AssigneeID = payload.AssigneeID
		t.AssigneeType = payload.AssigneeType
		if t.Status == TicketStatusSubmitted {
			t.Status = TicketStatusAssigned
		}
		t.UpdatedAt = event.OccurredAt()

	case *TicketPriorityUpdated:
		t.Priority = payload.NewPriority
		t.UpdatedAt = event.OccurredAt()

	case *TicketEscalated:
		t.CurrentTier = payload.ToTier
		t.Status = TicketStatusEscalated
		t.UpdatedAt = event.OccurredAt()

	case *TicketResolved:
		t.Status = TicketStatusResolved
		resolvedAt := event.OccurredAt()
		t.ResolvedAt = &resolvedAt
		t.UpdatedAt = event.OccurredAt()

	case *TicketReopened:
		t.Status = TicketStatusReopened
		t.ReopenCount = payload.ReopenCount
		t.ResolvedAt = nil
		t.UpdatedAt = event.OccurredAt()

	case *TicketClosed:
		t.Status = TicketStatusClosed
		t.SatisfactionRating = payload.SatisfactionRating
		closedAt := event.OccurredAt()
		t.ClosedAt = &closedAt
		t.UpdatedAt = event.OccurredAt()

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", t, event.EventName(), payload)
	}

	return nil
}

// ApplySnapshot implements es.Snapshotter
func (t *Ticket) ApplySnapshot(snapshot es.Snapshot) error {
	switch snap := snapshot.(type) {
	case *TicketV1:
		t.ChannelID = snap.ChannelID
		t.Title = snap.Title
		t.Description = snap.Description
		t.Status = snap.Status
		t.Priority = snap.Priority
		t.Category = snap.Category
		t.Tags = snap.Tags
		t.Metadata = snap.Metadata
		t.AssigneeID = snap.AssigneeID
		t.AssigneeType = snap.AssigneeType
		t.CreatedBy = snap.CreatedBy
		t.CurrentTier = snap.CurrentTier
		t.ResponseCount = snap.ResponseCount
		t.ReopenCount = snap.ReopenCount
		t.SatisfactionRating = snap.SatisfactionRating
		t.LinkedTicketIDs = snap.LinkedTicketIDs
		t.MergedTicketIDs = snap.MergedTicketIDs
		t.CreatedAt = snap.CreatedAt
		t.UpdatedAt = snap.UpdatedAt
		t.ResolvedAt = snap.ResolvedAt
		t.ClosedAt = snap.ClosedAt
		t.FirstResponseAt = snap.FirstResponseAt
		t.Attachments = snap.Attachments

	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", t, snapshot)
	}

	return nil
}

// ToSnapshot implements es.Snapshotter
func (t Ticket) ToSnapshot() es.Snapshot {
	return &TicketV1{
		ChannelID:          t.ChannelID,
		Title:              t.Title,
		Description:        t.Description,
		Status:             t.Status,
		Priority:           t.Priority,
		Category:           t.Category,
		Tags:               t.Tags,
		Metadata:           t.Metadata,
		AssigneeID:         t.AssigneeID,
		AssigneeType:       t.AssigneeType,
		CreatedBy:          t.CreatedBy,
		CurrentTier:        t.CurrentTier,
		ResponseCount:      t.ResponseCount,
		ReopenCount:        t.ReopenCount,
		SatisfactionRating: t.SatisfactionRating,
		LinkedTicketIDs:    t.LinkedTicketIDs,
		MergedTicketIDs:    t.MergedTicketIDs,
		CreatedAt:          t.CreatedAt,
		UpdatedAt:          t.UpdatedAt,
		ResolvedAt:         t.ResolvedAt,
		ClosedAt:           t.ClosedAt,
		FirstResponseAt:    t.FirstResponseAt,
		Attachments:        t.Attachments,
	}
}