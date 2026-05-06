package domain

import (
	"time"
)

// Support Channel Snapshots
type SupportChannelV1 struct {
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

func (SupportChannelV1) SnapshotName() string { return "support.SupportChannelV1" }

// Ticket Snapshots
type TicketV1 struct {
	ChannelID          string
	Title              string
	Description        string
	Status             TicketStatus
	Priority           TicketPriority
	Category           TicketCategory
	Tags               []string
	Metadata           map[string]string
	AssigneeID         string
	AssigneeType       AssigneeType
	CreatedBy          string
	CurrentTier        SupportTier
	ResponseCount      int
	ReopenCount        int
	SatisfactionRating *CustomerSatisfaction
	LinkedTicketIDs    []string
	MergedTicketIDs    []string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	ResolvedAt         *time.Time
	ClosedAt           *time.Time
	FirstResponseAt    *time.Time
	Attachments        []Attachment
}

func (TicketV1) SnapshotName() string { return "support.TicketV1" }

// Knowledge Article Snapshots
type KnowledgeArticleV1 struct {
	ID               string
	Title            string
	Content          string
	Categories       []string
	Tags             []string
	Public           bool
	ViewCount        int
	AverageRating    float64
	RatingCount      int
	CreatedBy        string
	RelatedTicketIDs []string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (KnowledgeArticleV1) SnapshotName() string { return "support.KnowledgeArticleV1" }