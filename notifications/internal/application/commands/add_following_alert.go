package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/notifications/internal/domain"
)

type AddFollowingAlert struct {
	ID         string
	UserID     string
	FollowerID string
	Message    string
}

type AddFollowingAlertHandler struct {
	Alerts    domain.AlertRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAddFollowingAlertHandler(Alerts domain.AlertRepository, publisher ddd.EventPublisher[ddd.Event]) AddFollowingAlertHandler {
	return AddFollowingAlertHandler{
		Alerts:    Alerts,
		publisher: publisher,
	}
}

func (h AddFollowingAlertHandler) AddFollowingAlert(ctx context.Context, cmd AddFollowingAlert) error {
	Alert, err := h.Alerts.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"FollowerID": cmd.FollowerID,
		"Message":    cmd.Message,
	}
	event, err := Alert.AddFollowingAlert(cmd.UserID, domain.FollowingType, payload)
	if err != nil {
		return err
	}
	if err = h.Alerts.Save(ctx, Alert); err != nil {
		return err
	}
	return h.publisher.Publish(ctx, event)
}