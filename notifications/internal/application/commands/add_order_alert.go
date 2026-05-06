package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/notifications/internal/domain"
)

type AddOrderAlert struct {
	ID             string
	UserID         string
	ProductID      string
	OrderID        string
	UserCustomerID string
	Message        string
}

type AddOrderAlertHandler struct {
	Alerts    domain.AlertRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAddOrderAlertHandler(Alerts domain.AlertRepository, publisher ddd.EventPublisher[ddd.Event]) AddOrderAlertHandler {
	return AddOrderAlertHandler{
		Alerts:    Alerts,
		publisher: publisher,
	}
}

func (h AddOrderAlertHandler) AddOrderAlert(ctx context.Context, cmd AddOrderAlert) error {
	Alerts, err := h.Alerts.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"OrderID": cmd.ProductID,
	}
	event, err := Alerts.AddOrderAlert(cmd.UserID, domain.OrderType, payload)
	if err != nil {
		return err
	}
	if err = h.Alerts.Save(ctx, Alerts); err != nil {
		return err
	}
	return h.publisher.Publish(ctx, event)
}
