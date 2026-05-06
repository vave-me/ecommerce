package commands

import (
	"context"
	"middleman/activity/internal/domain"
	"middleman/internal/ddd"
)

type (
	RemoveInteraction struct {
		ID         string
		ActivityID string
		ItemID     string
		ItemType   string
		ActionType string
	}

	RemoveInteractionHandler struct {
		interactions domain.InteractionRepository
		publisher    ddd.EventPublisher[ddd.Event]
	}
)

func NewRemoveInteractionHandler(interactions domain.InteractionRepository, publisher ddd.EventPublisher[ddd.Event]) RemoveInteractionHandler {
	return RemoveInteractionHandler{
		interactions: interactions,
		publisher:    publisher,
	}
}

func (h RemoveInteractionHandler) RemoveInteraction(ctx context.Context, cmd RemoveInteraction) error {
	interaction, err := h.interactions.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := interaction.Remove(cmd.ActivityID, cmd.ItemID)
	if err != nil {
		return err
	}

	err = h.interactions.Save(ctx, interaction)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
