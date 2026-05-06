package commands

import "middleman/support/internal/domain"

// Support Channel Commands
type UpdateSupportChannelSettings struct {
	ID       string
	Settings domain.SupportChannelSettings
}

type CloseSupportChannel struct {
	ID       string
	ClosedBy string
	Reason   string
}

type ReactivateSupportChannel struct {
	ID            string
	ReactivatedBy string
}

// Ticket Commands
type UpdateTicket struct {
	ID          string
	Title       *string
	Description *string
	Category    *domain.TicketCategory
	Tags        []string
	Metadata    map[string]string
	UpdatedBy   string
}

type UpdateTicketPriority struct {
	ID       string
	Priority domain.TicketPriority
	Reason   string
}

type EscalateTicket struct {
	ID               string
	EscalationTier   domain.SupportTier
	EscalatedBy      string
	EscalationReason string
	EscalationNotes  string
}

type ReopenTicket struct {
	ID           string
	ReopenedBy   string
	ReopenReason string
}

type CloseTicket struct {
	ID                 string
	ClosedBy           string
	ClosureNotes       string
	SatisfactionRating *domain.CustomerSatisfaction
}

type MergeTickets struct {
	PrimaryTicketID    string
	SecondaryTicketIDs []string
	MergedBy           string
	MergeReason        string
}

type LinkTickets struct {
	TicketID         string
	RelatedTicketIDs []string
	LinkedBy         string
	RelationshipType string
}

// Communication Commands
type AddInternalNote struct {
	ID             string
	TicketID       string
	AuthorID       string
	Content        string
	MentionedUsers []string
}

// AI Integration Commands
type EnableAISupport struct {
	ChannelID     string
	AssistantID   string
	Configuration domain.AIConfiguration
}

type ConfigureAIAssistant struct {
	ChannelID     string
	Configuration domain.AIConfiguration
}

type HandoffToHuman struct {
	TicketID string
	Reason   string
	Context  map[string]string
}

type HandoffToAI struct {
	TicketID string
	Context  map[string]string
}

// Knowledge Base Commands
type CreateKnowledgeArticle struct {
	ID         string
	Title      string
	Content    string
	Categories []string
	Tags       []string
	Public     bool
	CreatedBy  string
}

type LinkArticleToTicket struct {
	TicketID      string
	ArticleID     string
	HelpedResolve bool
	LinkedBy      string
}

type RateArticle struct {
	ArticleID string
	RatedBy   string
	Rating    int
	Feedback  string
}