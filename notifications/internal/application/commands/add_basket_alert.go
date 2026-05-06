package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/notifications/internal/domain"
)

// AddBasketAlert command struct remains unchanged
type AddBasketAlert struct {
	ID             string
	UserID         string
	BasketID       string
	ProductID      string
	UserCustomerID string
	UserSellerID   string
	Message        string
}

// AddBasketAlertHandler handles sending basket Alerts
type AddBasketAlertHandler struct {
	Alerts    domain.AlertRepository
	publisher ddd.EventPublisher[ddd.Event]
}

// NewAddBasketAlertHandler creates a new handler for AddBasketAlert
func NewAddBasketAlertHandler(Alerts domain.AlertRepository, publisher ddd.EventPublisher[ddd.Event]) AddBasketAlertHandler {
	return AddBasketAlertHandler{
		Alerts:    Alerts,
		publisher: publisher,
	}
}

// AddBasketAlert handles the logic for sending a basket Alert
func (h AddBasketAlertHandler) AddBasketAlert(ctx context.Context, cmd AddBasketAlert) error {
	// Load the basket entity to get more information (optional, depending on use case)
	Alert, err := h.Alerts.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"BasketID":       cmd.BasketID,
		"ProductID":      cmd.ProductID,
		"UserCustomerID": cmd.UserCustomerID,
	}

	event, err := Alert.AddBasketAlert(cmd.ID, cmd.UserID, domain.BasketType, payload)
	if err != nil {
		return errors.Wrap(err, "create basket Alert command")
	}

	if err = h.Alerts.Save(ctx, Alert); err != nil {
		return errors.Wrap(err, "basket Alerts creation")
	}
	// Publish the Alert event
	return h.publisher.Publish(ctx, event)
}
