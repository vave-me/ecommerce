package domain

import "time"

const (
	NewsletterV1Version   = "newsletters.NewsletterV1"
	SubscriptionV1Version = "newsletters.SubscriptionV1"
	EditionV1Version      = "newsletters.EditionV1"
	TemplateV1Version     = "newsletters.TemplateV1"
)

// Newsletter Snapshots
type NewsletterV1 struct {
	UserID      string
	Name        string
	Description string
	Frequency   string
	Category    string
	TemplateID  string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (NewsletterV1) SnapshotName() string { return NewsletterV1Version }

// Subscription Snapshots
type SubscriptionV1 struct {
	UserID            string
	NewsletterID      string
	Status            string
	FrequencyOverride string
	Topics            []string
	Format            string
	SubscribedAt      time.Time
	UnsubscribedAt    *time.Time
}

func (SubscriptionV1) SnapshotName() string { return SubscriptionV1Version }

// Edition Snapshots
type EditionV1 struct {
	NewsletterID   string
	Subject        string
	ContentHTML    string
	ContentText    string
	TemplateData   map[string]string
	ScheduledAt    *time.Time
	SentAt         *time.Time
	Status         string
	CreatedBy      string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	RecipientCount int
}

func (EditionV1) SnapshotName() string { return EditionV1Version }

// Template Snapshots
type TemplateV1 struct {
	UserID       string
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

func (TemplateV1) SnapshotName() string { return TemplateV1Version }