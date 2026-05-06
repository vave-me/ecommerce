package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/notifications/internal/domain"
)

// ReadAlert command struct remains unchanged
type ReadAlert struct {
	ID string
}

// ReadAlertHandler handles sending basket Alerts
type ReadAlertHandler struct {
	alerts    domain.AlertRepository
	publisher ddd.EventPublisher[ddd.Event]
}

// NewReadAlertHandler creates a new handler for ReadAlert
func NewReadAlertHandler(alerts domain.AlertRepository, publisher ddd.EventPublisher[ddd.Event]) ReadAlertHandler {
	return ReadAlertHandler{
		alerts:    alerts,
		publisher: publisher,
	}
}

// ReadAlert handles the logic for sending a basket Alert
func (h ReadAlertHandler) ReadAlert(ctx context.Context, cmd ReadAlert) error {
	// Load the basket entity to get more information (optional, depending on use case)
	alert, err := h.alerts.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := alert.Read(cmd.ID)
	if err != nil {
		return errors.Wrap(err, "create basket Alert command")
	}

	if err = h.alerts.Save(ctx, alert); err != nil {
		return errors.Wrap(err, "basket Alerts creation")
	}
	// Publish the Alert event
	return h.publisher.Publish(ctx, event)
}
