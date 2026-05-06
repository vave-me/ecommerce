package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/notifications/internal/domain"
)

type AddReviewAlert struct {
	ID        string
	UserID    string
	ReviewID  string
	ProductID string
	Message   string
}

type AddReviewAlertHandler struct {
	Alerts    domain.AlertRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAddReviewAlertHandler(Alerts domain.AlertRepository, publisher ddd.EventPublisher[ddd.Event]) AddReviewAlertHandler {
	return AddReviewAlertHandler{
		Alerts:    Alerts,
		publisher: publisher,
	}
}

func (h AddReviewAlertHandler) AddReviewAlert(ctx context.Context, cmd AddReviewAlert) error {
	Alert, err := h.Alerts.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"ReviewID":  cmd.ReviewID,
		"ProductID": cmd.ProductID,
		"Message":   cmd.Message,
	}
	event, err := Alert.AddReviewAlert(cmd.UserID, domain.ReviewType, payload)
	if err != nil {
		return err
	}
	if err = h.Alerts.Save(ctx, Alert); err != nil {
		return err
	}
	return h.publisher.Publish(ctx, event)
}