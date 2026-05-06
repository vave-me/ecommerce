package consciousness

import (
	"math"
	"time"

	"github.com/google/uuid"
	"middleman/internal/ddd"
)

const (
	PatternTypeActivitySurge     = "activity_surge"
	PatternTypeAbandonmentRisk   = "abandonment_risk"
	PatternTypeFraudRisk         = "fraud_risk"
	PatternTypeSupportCrisis     = "support_crisis"
	PatternTypeUserSurge         = "user_surge"
	PatternTypeCancellationSpike = "cancellation_spike"
)

type PatternDetectorFunc func(event ddd.Event, recentEvents []EventMemory) *Pattern

// Activity Surge Detector
func NewActivitySurgeDetector() PatternDetectorFunc {
	return func(event ddd.Event, recentEvents []EventMemory) *Pattern {
		threshold := time.Now().Add(-5 * time.Minute)
		recentCount := 0

		for _, e := range recentEvents {
			if e.Timestamp.After(threshold) {
				recentCount++
			}
		}

		// Calculate average events per 5 minutes from older data
		olderCount := 0
		olderThreshold := time.Now().Add(-1 * time.Hour)
		for _, e := range recentEvents {
			if e.Timestamp.Before(threshold) && e.Timestamp.After(olderThreshold) {
				olderCount++
			}
		}

		avgPer5Min := float64(olderCount) / 11.0

		// Detect surge if current rate is 2x normal
		if float64(recentCount) > avgPer5Min*2 && recentCount > 50 {
			return &Pattern{
				ID:         uuid.New().String(),
				Type:       PatternTypeActivitySurge,
				Confidence: math.Min(float64(recentCount)/(avgPer5Min*3), 1.0),
				Data: map[string]interface{}{
					"recent_count":  recentCount,
					"average_count": avgPer5Min,
					"surge_factor":  float64(recentCount) / avgPer5Min,
				},
			}
		}

		return nil
	}
}

// Abandonment Risk Detector
func NewAbandonmentRiskDetector() PatternDetectorFunc {
	return func(event ddd.Event, recentEvents []EventMemory) *Pattern {
		if event.EventName() != "BasketItemAdded" {
			return nil
		}

		threshold := time.Now().Add(-2 * time.Hour)
		basketsStarted := 0
		basketsCompleted := 0

		for _, e := range recentEvents {
			if e.Timestamp.After(threshold) {
				switch e.Type {
				case "BasketStarted", "BasketItemAdded":
					basketsStarted++
				case "BasketCheckedOut":
					basketsCompleted++
				}
			}
		}

		if basketsStarted > 10 {
			abandonmentRate := 1.0 - (float64(basketsCompleted) / float64(basketsStarted))
			if abandonmentRate > 0.7 {
				return &Pattern{
					ID:         uuid.New().String(),
					Type:       PatternTypeAbandonmentRisk,
					Confidence: abandonmentRate,
					Data: map[string]interface{}{
						"baskets_started":   basketsStarted,
						"baskets_completed": basketsCompleted,
						"abandonment_rate":  abandonmentRate,
					},
				}
			}
		}

		return nil
	}
}

// Fraud Risk Detector
func NewFraudRiskDetector() PatternDetectorFunc {
	return func(event ddd.Event, recentEvents []EventMemory) *Pattern {
		switch event.EventName() {
		case "OrderCreated", "PaymentRequested":
			// Try type assertion to get AggregateID if the concrete type supports it
			aggregateID := ""
			if eventWithAggregateID, ok := event.(interface{ AggregateID() string }); ok {
				aggregateID = eventWithAggregateID.AggregateID()
			}
			threshold := time.Now().Add(-1 * time.Hour)
			sameUserEvents := 0

			for _, e := range recentEvents {
				if e.Timestamp.After(threshold) && e.AggregateID == aggregateID {
					if e.Type == "OrderCreated" || e.Type == "PaymentRequested" {
						sameUserEvents++
					}
				}
			}

			// Flag if more than 5 orders in an hour
			if sameUserEvents > 5 {
				return &Pattern{
					ID:         uuid.New().String(),
					Type:       PatternTypeFraudRisk,
					Confidence: math.Min(float64(sameUserEvents)/10.0, 1.0),
					Data: map[string]interface{}{
						"user_events_count": sameUserEvents,
						"time_window":       "1h",
					},
				}
			}
		}

		return nil
	}
}

// Support Crisis Detector
func NewSupportCrisisDetector() PatternDetectorFunc {
	return func(event ddd.Event, recentEvents []EventMemory) *Pattern {
		if event.EventName() != "TicketCreated" {
			return nil
		}

		threshold := time.Now().Add(-1 * time.Hour)
		ticketCount := 0

		for _, e := range recentEvents {
			if e.Timestamp.After(threshold) && e.Type == "TicketCreated" {
				ticketCount++
			}
		}

		// Crisis if more than 20 tickets in an hour
		if ticketCount > 20 {
			return &Pattern{
				ID:         uuid.New().String(),
				Type:       PatternTypeSupportCrisis,
				Confidence: math.Min(float64(ticketCount)/30.0, 1.0),
				Data: map[string]interface{}{
					"ticket_count": ticketCount,
					"time_window":  "1h",
					"threshold":    20,
				},
			}
		}

		return nil
	}
}

// User Surge Detector
func NewUserSurgeDetector() PatternDetectorFunc {
	return func(event ddd.Event, recentEvents []EventMemory) *Pattern {
		if event.EventName() != "UserCreated" {
			return nil
		}

		threshold := time.Now().Add(-30 * time.Minute)
		newUserCount := 0

		for _, e := range recentEvents {
			if e.Timestamp.After(threshold) && e.Type == "UserCreated" {
				newUserCount++
			}
		}

		// Surge if more than 50 new users in 30 minutes
		if newUserCount > 50 {
			return &Pattern{
				ID:         uuid.New().String(),
				Type:       PatternTypeUserSurge,
				Confidence: math.Min(float64(newUserCount)/100.0, 1.0),
				Data: map[string]interface{}{
					"new_user_count": newUserCount,
					"time_window":    "30m",
					"threshold":      50,
				},
			}
		}

		return nil
	}
}

// Cancellation Spike Detector
func NewCancellationSpikeDetector() PatternDetectorFunc {
	return func(event ddd.Event, recentEvents []EventMemory) *Pattern {
		if event.EventName() != "OrderCanceled" {
			return nil
		}

		threshold := time.Now().Add(-2 * time.Hour)
		canceledCount := 0
		completedCount := 0

		for _, e := range recentEvents {
			if e.Timestamp.After(threshold) {
				switch e.Type {
				case "OrderCanceled":
					canceledCount++
				case "OrderCompleted":
					completedCount++
				}
			}
		}

		totalOrders := canceledCount + completedCount
		if totalOrders > 10 {
			cancellationRate := float64(canceledCount) / float64(totalOrders)
			if cancellationRate > 0.3 {
				return &Pattern{
					ID:         uuid.New().String(),
					Type:       PatternTypeCancellationSpike,
					Confidence: cancellationRate,
					Data: map[string]interface{}{
						"canceled_count":    canceledCount,
						"completed_count":   completedCount,
						"cancellation_rate": cancellationRate,
					},
				}
			}
		}

		return nil
	}
}