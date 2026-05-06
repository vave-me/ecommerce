package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/notifications/internal/domain"
)

type AddSupportAlert struct {
	ID       string
	UserID   string
	TicketID string
	Message  string
}

type AddSupportAlertHandler struct {
	alerts    domain.AlertRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAddSupportAlertHandler(alerts domain.AlertRepository, publisher ddd.EventPublisher[ddd.Event]) AddSupportAlertHandler {
	return AddSupportAlertHandler{
		alerts:    alerts,
		publisher: publisher,
	}
}

func (h AddSupportAlertHandler) AddSupportAlert(ctx context.Context, cmd AddSupportAlert) error {
	Alert, err := h.alerts.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"TicketID": cmd.TicketID,
	}
	event, err := Alert.AddSupportAlert(cmd.UserID, domain.SupportType, payload)
	if err != nil {
		return err
	}
	if err = h.alerts.Save(ctx, Alert); err != nil {
		return err
	}
	return h.publisher.Publish(ctx, event)
}
