package commands

import (
	"context"
	"middleman/activity/internal/domain"
	"middleman/internal/ddd"
)

type (
	AddInteraction struct {
		ID         string
		ActivityID string
		ItemID     string
		ItemType   string
		ActionType string
	}

	AddInteractionHandler struct {
		interactions domain.InteractionRepository
		publisher    ddd.EventPublisher[ddd.Event]
	}
)

func NewAddInteractionHandler(interactions domain.InteractionRepository, publisher ddd.EventPublisher[ddd.Event]) AddInteractionHandler {
	return AddInteractionHandler{
		interactions: interactions,
		publisher:    publisher,
	}
}

func (h AddInteractionHandler) AddInteraction(ctx context.Context, cmd AddInteraction) error {
	interaction, err := h.interactions.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := interaction.AddInteraction(cmd.ActivityID, cmd.ItemID, cmd.ItemType, cmd.ActionType)
	if err != nil {
		return err
	}

	err = h.interactions.Save(ctx, interaction)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
