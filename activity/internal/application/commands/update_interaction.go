package commands

import (
	"context"
	"middleman/activity/internal/domain"
	"middleman/internal/ddd"
)

type (
	UpdateInteraction struct {
		ID         string
		ActivityID string
		ItemID     string
		ItemType   string
		ActionType string
	}

	UpdateInteractionHandler struct {
		interactions domain.InteractionRepository
		publisher    ddd.EventPublisher[ddd.Event]
	}
)

func NewUpdateInteractionHandler(interactions domain.InteractionRepository, publisher ddd.EventPublisher[ddd.Event]) UpdateInteractionHandler {
	return UpdateInteractionHandler{
		interactions: interactions,
		publisher:    publisher,
	}
}

func (h UpdateInteractionHandler) UpdateInteraction(ctx context.Context, cmd UpdateInteraction) error {
	interaction, err := h.interactions.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}
	event, err := interaction.Update(cmd.ActivityID, cmd.ItemID, cmd.ActionType)
	if err != nil {
		return err
	}

	err = h.interactions.Save(ctx, interaction)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
