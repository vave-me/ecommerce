package models

import "time"

// Subscription represents a user's subscription to a newsletter
type Subscription struct {
	SubscriptionID          string    `json:"subscription_id"`
	UserID                  string    `json:"user_id"`
	NewsletterID            string    `json:"newsletter_id"`
	SubscriptionPreferences string    `json:"subscription_preferences"`
	SubscriptionStatus      string    `json:"subscription_status"`
	CreatedAt               time.Time `json:"created_at,omitempty"`
	UpdatedAt               time.Time `json:"updated_at,omitempty"`
}

// SubscribeNewsletterResponse represents the response for subscribing to a newsletter
type SubscribeNewsletterResponse struct {
	SubscriptionID string `json:"subscription_id"`
}

// GetSubscriptionResponse represents the response for getting a subscription
type GetSubscriptionResponse struct {
	Subscription *Subscription `json:"subscription"`
}

// ListSubscriptionsResponse represents the response for listing subscriptions
type ListSubscriptionsResponse struct {
	Subscriptions []*Subscription `json:"subscriptions"`
	Total         int32           `json:"total"`
	Page          int32           `json:"page"`
	Limit         int32           `json:"limit"`
}

// UpdateSubscriptionResponse represents the response for updating a subscription
type UpdateSubscriptionResponse struct {
	Subscription *Subscription `json:"subscription"`
}

// SendNewsletterResponse represents the response for sending a newsletter
type SendNewsletterResponse struct {
	NewsletterID string `json:"newsletter_id"`
	SentCount    int32  `json:"sent_count,omitempty"`
	Message      string `json:"message,omitempty"`
}

// Newsletter represents a newsletter entity
type Newsletter struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Subscription status constants
const (
	SubscriptionStatusActive    = "active"
	SubscriptionStatusInactive  = "inactive"
	SubscriptionStatusPending   = "pending"
	SubscriptionStatusCancelled = "cancelled"
	SubscriptionStatusSuspended = "suspended"
)

// Newsletter status constants
const (
	NewsletterStatusDraft     = "draft"
	NewsletterStatusPublished = "published"
	NewsletterStatusArchived  = "archived"
	NewsletterStatusScheduled = "scheduled"
)

// Subscription preference constants
const (
	PreferenceDaily   = "daily"
	PreferenceWeekly  = "weekly"
	PreferenceMonthly = "monthly"
	PreferenceInstant = "instant"
	PreferenceDigest  = "digest"
)
