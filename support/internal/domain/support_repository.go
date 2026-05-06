package domain

import (
	"context"
	"time"
)

// Repository interfaces for event-sourced aggregates
type SupportChannelRepository interface {
	Load(ctx context.Context, channelID string) (*SupportChannel, error)
	Save(ctx context.Context, channel *SupportChannel) error
}

type TicketRepository interface {
	Load(ctx context.Context, ticketID string) (*Ticket, error)
	Save(ctx context.Context, ticket *Ticket) error
}

// Catalog types for read models
type SupportChannelCatalog struct {
	ID           string
	UserID       string
	BusinessID   string
	ChannelType  string
	Active       bool
	OpenTickets  int
	TotalTickets int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type TicketCatalog struct {
	ID           string
	ChannelID    string
	Title        string
	Status       string
	Priority     string
	Category     string
	AssigneeID   *string
	AssigneeType *string
	CreatedBy    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type KnowledgeArticleCatalog struct {
	ID            string
	Title         string
	Categories    []string
	Public        bool
	ViewCount     int
	AverageRating float64
	CreatedAt     time.Time
}

// Domain models for non-event-sourced entities
type KnowledgeArticle struct {
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

type Communication struct {
	ID             string
	TicketID       string
	AuthorID       string
	AuthorType     AuthorType
	Content        string
	IsPublic       bool
	Attachments    []Attachment
	MentionedUsers []string
	Metadata       map[string]string
	CreatedAt      time.Time
}

type AIConfiguration struct {
	ID                     string
	ChannelID              string
	AssistantID            string
	AllowedActions         []string
	KnowledgeBaseAccess    map[string]string
	MaxHandlingTier        SupportTier
	CanCloseTickets        bool
	CanIssueRefunds        bool
	ConfidenceThreshold    float64
	AutoResponseCategories []string
	MaxTokens              int
	Temperature            float64
	Active                 bool
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

// Catalog repository interfaces for read models
type SupportChannelCatalogRepository interface {
	Add(ctx context.Context, channel *SupportChannelCatalog) error
	Update(ctx context.Context, channel *SupportChannelCatalog) error
	Delete(ctx context.Context, channelID string) error
	Find(ctx context.Context, channelID string) (*SupportChannelCatalog, error)
	GetByUserID(ctx context.Context, userID string, activeOnly bool, limit, offset int) ([]*SupportChannelCatalog, error)
	GetByBusinessID(ctx context.Context, businessID string, activeOnly bool, limit, offset int) ([]*SupportChannelCatalog, error)
}

type TicketCatalogRepository interface {
	Add(ctx context.Context, ticket *TicketCatalog) error
	Update(ctx context.Context, ticket *TicketCatalog) error
	Delete(ctx context.Context, ticketID string) error
	Find(ctx context.Context, ticketID string) (*TicketCatalog, error)
	GetByChannelID(ctx context.Context, channelID string, status *string, limit, offset int) ([]*TicketCatalog, error)
	GetByAssigneeID(ctx context.Context, assigneeID string, status *string, limit, offset int) ([]*TicketCatalog, error)
	Search(ctx context.Context, query string, filters map[string]interface{}, limit, offset int) ([]*TicketCatalog, error)
	Count(ctx context.Context, channelID string, status *string) (int, error)
}

type KnowledgeArticleCatalogRepository interface {
	Add(ctx context.Context, article *KnowledgeArticleCatalog) error
	Update(ctx context.Context, article *KnowledgeArticleCatalog) error
	Delete(ctx context.Context, articleID string) error
	Find(ctx context.Context, articleID string) (*KnowledgeArticleCatalog, error)
	Search(ctx context.Context, query string, categories []string, publicOnly bool, limit, offset int) ([]*KnowledgeArticleCatalog, error)
	GetPopular(ctx context.Context, limit int) ([]*KnowledgeArticleCatalog, error)
}

// Communication repository for storing replies and notes
type CommunicationRepository interface {
	Add(ctx context.Context, comm *Communication) error
	GetByTicketID(ctx context.Context, ticketID string, includeInternal bool, limit, offset int) ([]*Communication, error)
	Count(ctx context.Context, ticketID string, includeInternal bool) (int, error)
}

// AI configuration repository
type AIConfigurationRepository interface {
	Add(ctx context.Context, config *AIConfiguration) error
	Update(ctx context.Context, config *AIConfiguration) error
	GetByChannelID(ctx context.Context, channelID string) (*AIConfiguration, error)
	Delete(ctx context.Context, channelID string) error
}