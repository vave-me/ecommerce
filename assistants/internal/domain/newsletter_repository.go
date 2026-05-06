package domain

import (
	"context"
	"middleman/assistants/internal/models"
)

type NewsletterRepository interface {
	// Core newsletter operations from gRPC service
	SubscribeNewsletter(ctx context.Context, userID, newsletterID, subscriptionPreferences string) (*models.SubscribeNewsletterResponse, error)
	UnsubscribeNewsletter(ctx context.Context, subscriptionID string) error
	GetSubscription(ctx context.Context, subscriptionID string) (*models.GetSubscriptionResponse, error)
	ListSubscriptions(ctx context.Context, userID, subscriptionStatus string, page, limit int32) (*models.ListSubscriptionsResponse, error)
	UpdateSubscription(ctx context.Context, subscriptionID, subscriptionPreferences, subscriptionStatus string) (*models.UpdateSubscriptionResponse, error)
	SendNewsletter(ctx context.Context, newsletterID, content string) (*models.SendNewsletterResponse, error)

	// Additional query methods for AI tooling
	GetSubscriptionByID(ctx context.Context, subscriptionID string) (*models.Subscription, error)
	GetSubscriptionsByUser(ctx context.Context, userID string, limit int32) ([]*models.Subscription, error)
	GetSubscriptionsByNewsletter(ctx context.Context, newsletterID string, limit int32) ([]*models.Subscription, error)
	GetActiveSubscriptions(ctx context.Context, userID string, limit int32) ([]*models.Subscription, error)
	GetInactiveSubscriptions(ctx context.Context, userID string, limit int32) ([]*models.Subscription, error)
	SearchSubscriptions(ctx context.Context, userID, query string, limit int32) ([]*models.Subscription, error)
	CountSubscriptions(ctx context.Context, userID, subscriptionStatus string) (int32, error)
	CountActiveSubscriptions(ctx context.Context, userID string) (int32, error)
	CountSubscriptionsByNewsletter(ctx context.Context, newsletterID string) (int32, error)

	// Newsletter management methods
	GetNewsletter(ctx context.Context, newsletterID string) (*models.Newsletter, error)
	GetNewsletters(ctx context.Context, status string, limit int32) ([]*models.Newsletter, error)
	SearchNewsletters(ctx context.Context, query string, limit int32) ([]*models.Newsletter, error)
}
