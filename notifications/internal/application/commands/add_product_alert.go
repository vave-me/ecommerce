package commands

import (
	"context"
	"fmt"
	"middleman/internal/ddd"
	"middleman/notifications/internal/domain"
)

type AddProductAlert struct {
	ID        string
	ProductID string
	UserID    string
	Message   string
}

type AddProductAlertHandler struct {
	alerts    domain.AlertRepository
	catalog   domain.CatalogRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAddProductAlertHandler(alerts domain.AlertRepository, publisher ddd.EventPublisher[ddd.Event]) AddProductAlertHandler {
	return AddProductAlertHandler{
		alerts:    alerts,
		publisher: publisher,
	}
}

func (h AddProductAlertHandler) AddProductAlert(ctx context.Context, cmd AddProductAlert) error {
	Alert, err := h.alerts.Load(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("loading alert %w", err)
	}

	payload := map[string]interface{}{
		"ProductID": cmd.ProductID,
		"Message":   cmd.Message,
	}
	event, err := Alert.AddProductAlert(cmd.UserID, domain.ProductType, payload)
	if err != nil {
		return err
	}
	if err = h.alerts.Save(ctx, Alert); err != nil {
		return err
	}
	return h.publisher.Publish(ctx, event)
}
