package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/notifications/internal/domain"
)

type AddPaymentAlert struct {
	ID        string
	UserID    string
	PaymentID string
	OrderID   string
	Message   string
}

type AddPaymentAlertHandler struct {
	Alerts    domain.AlertRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAddPaymentAlertHandler(Alerts domain.AlertRepository, publisher ddd.EventPublisher[ddd.Event]) AddPaymentAlertHandler {
	return AddPaymentAlertHandler{
		Alerts:    Alerts,
		publisher: publisher,
	}
}

func (h AddPaymentAlertHandler) AddPaymentAlert(ctx context.Context, cmd AddPaymentAlert) error {
	Alert, err := h.Alerts.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"PaymentID": cmd.PaymentID,
		"OrderID":   cmd.OrderID,
		"Message":   cmd.Message,
	}
	event, err := Alert.AddPaymentAlert(cmd.UserID, domain.PaymentType, payload)
	if err != nil {
		return err
	}
	if err = h.Alerts.Save(ctx, Alert); err != nil {
		return err
	}
	return h.publisher.Publish(ctx, event)
}