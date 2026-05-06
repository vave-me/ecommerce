package domain

// Support Channel Events
const (
	SupportChannelCreatedEvent          = "support.SupportChannelCreated"
	SupportChannelSettingsUpdatedEvent  = "support.SupportChannelSettingsUpdated"
	SupportChannelClosedEvent           = "support.SupportChannelClosed"
	SupportChannelReactivatedEvent      = "support.SupportChannelReactivated"
)

type SupportChannelCreated struct {
	UserID      string
	BusinessID  string
	ChannelType SupportChannelType
	Settings    SupportChannelSettings
}

func (SupportChannelCreated) Key() string { return SupportChannelCreatedEvent }

type SupportChannelSettingsUpdated struct {
	Settings SupportChannelSettings
}

func (SupportChannelSettingsUpdated) Key() string { return SupportChannelSettingsUpdatedEvent }

type SupportChannelClosed struct {
	ClosedBy string
	Reason   string
}

func (SupportChannelClosed) Key() string { return SupportChannelClosedEvent }

type SupportChannelReactivated struct {
	ReactivatedBy string
}

func (SupportChannelReactivated) Key() string { return SupportChannelReactivatedEvent }

// Ticket Events
const (
	TicketCreatedEvent         = "support.TicketCreated"
	TicketUpdatedEvent         = "support.TicketUpdated"
	TicketAssignedEvent        = "support.TicketAssigned"
	TicketPriorityUpdatedEvent = "support.TicketPriorityUpdated"
	TicketEscalatedEvent       = "support.TicketEscalated"
	TicketResolvedEvent        = "support.TicketResolved"
	TicketReopenedEvent        = "support.TicketReopened"
	TicketClosedEvent          = "support.TicketClosed"
	TicketsMergedEvent         = "support.TicketsMerged"
	TicketsLinkedEvent         = "support.TicketsLinked"
)

type TicketCreated struct {
	ChannelID   string
	Title       string
	Description string
	Category    TicketCategory
	Priority    TicketPriority
	Tags        []string
	Metadata    map[string]string
	CreatedBy   string
	Attachments []Attachment
}

func (TicketCreated) Key() string { return TicketCreatedEvent }

type TicketUpdated struct {
	Title       *string
	Description *string
	Category    *TicketCategory
	Tags        []string
	Metadata    map[string]string
	UpdatedBy   string
}

func (TicketUpdated) Key() string { return TicketUpdatedEvent }

type TicketAssigned struct {
	AssigneeID       string
	AssigneeType     AssigneeType
	AssignedBy       string
	AssignmentReason string
}

func (TicketAssigned) Key() string { return TicketAssignedEvent }

type TicketPriorityUpdated struct {
	OldPriority TicketPriority
	NewPriority TicketPriority
	UpdatedBy   string
	Reason      string
}

func (TicketPriorityUpdated) Key() string { return TicketPriorityUpdatedEvent }

type TicketEscalated struct {
	FromTier         SupportTier
	ToTier           SupportTier
	EscalatedBy      string
	EscalationReason string
	EscalationNotes  string
}

func (TicketEscalated) Key() string { return TicketEscalatedEvent }

type TicketResolved struct {
	ResolvedBy       string
	Resolution       string
	AppliedSolutions []string
}

func (TicketResolved) Key() string { return TicketResolvedEvent }

type TicketReopened struct {
	ReopenedBy   string
	ReopenReason string
	ReopenCount  int
}

func (TicketReopened) Key() string { return TicketReopenedEvent }

type TicketClosed struct {
	ClosedBy           string
	ClosureNotes       string
	SatisfactionRating *CustomerSatisfaction
}

func (TicketClosed) Key() string { return TicketClosedEvent }

type TicketsMerged struct {
	PrimaryTicketID    string
	SecondaryTicketIDs []string
	MergedBy           string
	MergeReason        string
}

func (TicketsMerged) Key() string { return TicketsMergedEvent }

type TicketsLinked struct {
	TicketID          string
	RelatedTicketIDs  []string
	LinkedBy          string
	RelationshipType  string
}

func (TicketsLinked) Key() string { return TicketsLinkedEvent }

// Communication Events
const (
	TicketReplyAddedEvent   = "support.TicketReplyAdded"
	InternalNoteAddedEvent  = "support.InternalNoteAdded"
)

type TicketReplyAdded struct {
	ID          string
	TicketID    string
	AuthorID    string
	AuthorType  AuthorType
	Content     string
	Attachments []Attachment
	IsPublic    bool
}

func (TicketReplyAdded) Key() string { return TicketReplyAddedEvent }

type InternalNoteAdded struct {
	ID              string
	TicketID        string
	AuthorID        string
	Content         string
	MentionedUsers  []string
}

func (InternalNoteAdded) Key() string { return InternalNoteAddedEvent }

// Author Types
type AuthorType string

const (
	AuthorTypeCustomer AuthorType = "CUSTOMER"
	AuthorTypeAgent    AuthorType = "AGENT"
	AuthorTypeAI       AuthorType = "AI"
	AuthorTypeSystem   AuthorType = "SYSTEM"
)