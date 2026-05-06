package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/notifications/internal/domain"
)

type AddInteractionAlert struct {
	ID            string
	Alert         string
	InteractionID string
	UserSenderID  string
	ProductID     string
	Message       string
}

type AddInteractionAlertHandler struct {
	Alerts    domain.AlertRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAddInteractionAlertHandler(Alerts domain.AlertRepository, publisher ddd.EventPublisher[ddd.Event]) AddInteractionAlertHandler {
	return AddInteractionAlertHandler{
		Alerts:    Alerts,
		publisher: publisher,
	}
}

func (h AddInteractionAlertHandler) AddInteractionAlert(ctx context.Context, cmd AddInteractionAlert) error {
	Alert, err := h.Alerts.Load(ctx, cmd.InteractionID)
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"ProductID": cmd.ProductID,
	}
	event, err := Alert.AddInteractionAlert(cmd.InteractionID, cmd.UserSenderID, domain.InteractionType, payload)
	if err != nil {
		return err
	}
	if err = h.Alerts.Save(ctx, Alert); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
