package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/notifications/internal/domain"
)

type AddMessageAlert struct {
	ID              string
	UserID          string
	ProductID       string
	MessageID       string
	MessageSenderID string
	Message         string
}

type AddMessageAlertHandler struct {
	Alerts    domain.AlertRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAddMessageAlertHandler(Alerts domain.AlertRepository, publisher ddd.EventPublisher[ddd.Event]) AddMessageAlertHandler {
	return AddMessageAlertHandler{
		Alerts:    Alerts,
		publisher: publisher,
	}
}

func (h AddMessageAlertHandler) AddMessageAlert(ctx context.Context, cmd AddMessageAlert) error {
	Alert, err := h.Alerts.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"MessageID":       cmd.MessageID,
		"MessageSenderID": cmd.MessageSenderID,
		"ProductID":       cmd.ProductID,
		"Message":         cmd.Message,
	}
	event, err := Alert.AddMessageAlert(cmd.UserID, domain.MessageType, payload)
	if err != nil {
		return err
	}
	if err = h.Alerts.Save(ctx, Alert); err != nil {
		return err
	}
	return h.publisher.Publish(ctx, event)
}
