package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/notifications/internal/domain"
)

type AddWishlistAlert struct {
	ID              string
	UserID          string
	ProductID       string
	WishlistID      string
	WishlistAdderID string
	Message         string
}

type AddWishlistAlertHandler struct {
	alerts    domain.AlertRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAddWishlistAlertHandler(alerts domain.AlertRepository, publisher ddd.EventPublisher[ddd.Event]) AddWishlistAlertHandler {
	return AddWishlistAlertHandler{
		alerts:    alerts,
		publisher: publisher,
	}
}

func (h AddWishlistAlertHandler) AddWishlistAlert(ctx context.Context, cmd AddWishlistAlert) error {
	Alerts, err := h.alerts.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"ProductID": cmd.ProductID,
	}
	event, err := Alerts.AddWishlistAlert(cmd.UserID, domain.WishlistType, payload)
	if err != nil {
		return err
	}
	if err = h.alerts.Save(ctx, Alerts); err != nil {
		return err
	}
	return h.publisher.Publish(ctx, event)
}
