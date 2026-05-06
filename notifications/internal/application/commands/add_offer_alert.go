package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/notifications/internal/domain"
)

type AddOfferAlert struct {
	ID            string
	UserID        string
	ProductID     string
	OfferID       string
	OfferSenderID string
	Message       string
}

type AddOfferAlertHandler struct {
	Alerts    domain.AlertRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAddOfferAlertHandler(Alerts domain.AlertRepository, publisher ddd.EventPublisher[ddd.Event]) AddOfferAlertHandler {
	return AddOfferAlertHandler{
		Alerts:    Alerts,
		publisher: publisher,
	}
}

func (h AddOfferAlertHandler) AddOfferAlert(ctx context.Context, cmd AddOfferAlert) error {
	Alert, err := h.Alerts.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"ProductID": cmd.ProductID,
	}
	event, err := Alert.AddOfferAlert(cmd.UserID, domain.OfferType, payload)
	if err != nil {
		return err
	}
	if err = h.Alerts.Save(ctx, Alert); err != nil {
		return err
	}
	return h.publisher.Publish(ctx, event)
}
