package domain

import (
	"context"
	"time"
)

// Read models for queries
type CatalogNewsletter struct {
	ID              string
	UserID          string
	Name            string
	Description     string
	Frequency       string
	Category        string
	TemplateID      string
	IsActive        bool
	SubscriberCount int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type CatalogSubscription struct {
	ID                string
	UserID            string
	NewsletterID      string
	Status            string
	FrequencyOverride string
	Topics            []string
	Format            string
	SubscribedAt      time.Time
	UnsubscribedAt    *time.Time
	Newsletter        *CatalogNewsletter // Populated in queries
}

type CatalogEdition struct {
	ID             string
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

type CatalogTemplate struct {
	ID           string
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

// Repository interfaces for read operations
type NewsletterCatalogRepository interface {
	Find(ctx context.Context, id string) (*CatalogNewsletter, error)
	FindByUser(ctx context.Context, userID string, limit, offset int) ([]*CatalogNewsletter, int, error)
	FindByCategory(ctx context.Context, category string, activeOnly bool, limit, offset int) ([]*CatalogNewsletter, int, error)
	FindAll(ctx context.Context, activeOnly bool, limit, offset int) ([]*CatalogNewsletter, int, error)
	Add(ctx context.Context, newsletter *CatalogNewsletter) error
	Update(ctx context.Context, newsletter *CatalogNewsletter) error
	Delete(ctx context.Context, id string) error
}

type SubscriptionCatalogRepository interface {
	Find(ctx context.Context, id string) (*CatalogSubscription, error)
	FindByUser(ctx context.Context, userID string, status string, limit, offset int) ([]*CatalogSubscription, int, error)
	FindByNewsletter(ctx context.Context, newsletterID string, status string, limit, offset int) ([]*CatalogSubscription, int, error)
	FindByUserAndNewsletter(ctx context.Context, userID, newsletterID string) (*CatalogSubscription, error)
	Add(ctx context.Context, subscription *CatalogSubscription) error
	Update(ctx context.Context, subscription *CatalogSubscription) error
	CountActiveByNewsletter(ctx context.Context, newsletterID string) (int, error)
}

type EditionCatalogRepository interface {
	Find(ctx context.Context, id string) (*CatalogEdition, error)
	FindByNewsletter(ctx context.Context, newsletterID string, status string, limit, offset int) ([]*CatalogEdition, int, error)
	FindScheduled(ctx context.Context, before time.Time) ([]*CatalogEdition, error)
	Add(ctx context.Context, edition *CatalogEdition) error
	Update(ctx context.Context, edition *CatalogEdition) error
}

type TemplateCatalogRepository interface {
	Find(ctx context.Context, id string) (*CatalogTemplate, error)
	FindByUser(ctx context.Context, userID string, limit, offset int) ([]*CatalogTemplate, int, error)
	FindPublic(ctx context.Context, limit, offset int) ([]*CatalogTemplate, int, error)
	Add(ctx context.Context, template *CatalogTemplate) error
	Update(ctx context.Context, template *CatalogTemplate) error
	Delete(ctx context.Context, id string) error
}