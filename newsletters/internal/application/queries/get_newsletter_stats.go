package queries

import (
	"context"
	"time"
	"middleman/newsletters/internal/domain"
)

type GetNewsletterStats struct {
	NewsletterID string
	StartDate    time.Time
	EndDate      time.Time
}

type NewsletterStats struct {
	TotalSubscribers  int
	NewSubscribers    int
	Unsubscribes      int
	EditionsSent      int
	AverageOpenRate   float32
	AverageClickRate  float32
}

type GetNewsletterStatsHandler struct {
	subCatalog     domain.SubscriptionCatalogRepository
	editionCatalog domain.EditionCatalogRepository
}

func NewGetNewsletterStatsHandler(
	subCatalog domain.SubscriptionCatalogRepository,
	editionCatalog domain.EditionCatalogRepository,
) GetNewsletterStatsHandler {
	return GetNewsletterStatsHandler{
		subCatalog:     subCatalog,
		editionCatalog: editionCatalog,
	}
}

func (h GetNewsletterStatsHandler) GetNewsletterStats(ctx context.Context, query GetNewsletterStats) (*NewsletterStats, error) {
	// Get current active subscribers
	totalSubscribers, err := h.subCatalog.CountActiveByNewsletter(ctx, query.NewsletterID)
	if err != nil {
		return nil, err
	}

	// Get sent editions in date range
	editions, _, err := h.editionCatalog.FindByNewsletter(ctx, query.NewsletterID, domain.SentStatus.String(), 1000, 0)
	if err != nil {
		return nil, err
	}

	editionsSent := 0
	for _, edition := range editions {
		if edition.SentAt != nil && 
		   edition.SentAt.After(query.StartDate) && 
		   edition.SentAt.Before(query.EndDate) {
			editionsSent++
		}
	}

	// TODO: Calculate new subscribers, unsubscribes, open/click rates from send logs

	return &NewsletterStats{
		TotalSubscribers: totalSubscribers,
		NewSubscribers:   0, // TODO: Implement
		Unsubscribes:     0, // TODO: Implement
		EditionsSent:     editionsSent,
		AverageOpenRate:  0, // TODO: Implement from send logs
		AverageClickRate: 0, // TODO: Implement from send logs
	}, nil
}