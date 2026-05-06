package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/notifications/internal/domain"
)

type AddCommentAlert struct {
	ID          string
	UserID      string
	CommentID   string
	UserAddedID string
	ProductID   string
	Message     string
}

type AddCommentAlertHandler struct {
	Alerts    domain.AlertRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAddCommentAlertHandler(Alerts domain.AlertRepository, publisher ddd.EventPublisher[ddd.Event]) AddCommentAlertHandler {
	return AddCommentAlertHandler{
		Alerts:    Alerts,
		publisher: publisher,
	}
}

func (h AddCommentAlertHandler) AddCommentAlert(ctx context.Context, cmd AddCommentAlert) error {
	Alert, err := h.Alerts.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"CommentID": cmd.CommentID,
		"ProductID": cmd.ProductID,
		"Message":   cmd.Message,
	}
	event, err := Alert.AddCommentAlert(cmd.UserID, domain.CommentType, payload)
	if err != nil {
		return err
	}
	if err = h.Alerts.Save(ctx, Alert); err != nil {
		return err
	}
	return h.publisher.Publish(ctx, event)
}
