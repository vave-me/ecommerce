package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/newsletters/internal/domain"
)

type UpdateSubscription struct {
	ID          string
	Status      string
	Preferences domain.SubscriptionPreferences
}

type UpdateSubscriptionHandler struct {
	subscriptions domain.SubscriptionRepository
	catalog       domain.SubscriptionCatalogRepository
	publisher     ddd.EventPublisher[ddd.Event]
}

func NewUpdateSubscriptionHandler(
	subscriptions domain.SubscriptionRepository,
	catalog domain.SubscriptionCatalogRepository,
	publisher ddd.EventPublisher[ddd.Event],
) UpdateSubscriptionHandler {
	return UpdateSubscriptionHandler{
		subscriptions: subscriptions,
		catalog:       catalog,
		publisher:     publisher,
	}
}

func (h UpdateSubscriptionHandler) UpdateSubscription(ctx context.Context, cmd UpdateSubscription) error {
	subscription, err := h.subscriptions.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	var event ddd.Event

	// Handle status changes
	if cmd.Status != "" && cmd.Status != subscription.Status.String() {
		switch cmd.Status {
		case domain.PausedStatus.String():
			event, err = subscription.Pause()
		case domain.ActiveStatus.String():
			event, err = subscription.Resume()
		}
		if err != nil {
			return err
		}
	}

	// Update preferences
	prefEvent, err := subscription.UpdatePreferences(cmd.Preferences)
	if err != nil {
		return err
	}

	err = h.subscriptions.Save(ctx, subscription)
	if err != nil {
		return err
	}

	// Update catalog
	catalogSub, err := h.catalog.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if cmd.Status != "" {
		catalogSub.Status = cmd.Status
	}
	
	catalogSub.Topics = cmd.Preferences.Topics
	catalogSub.Format = cmd.Preferences.Format.String()
	if cmd.Preferences.FrequencyOverride != nil {
		catalogSub.FrequencyOverride = cmd.Preferences.FrequencyOverride.String()
	} else {
		catalogSub.FrequencyOverride = ""
	}

	err = h.catalog.Update(ctx, catalogSub)
	if err != nil {
		return err
	}

	// Publish events
	if event != nil {
		err = h.publisher.Publish(ctx, event)
		if err != nil {
			return err
		}
	}

	return h.publisher.Publish(ctx, prefEvent)
}