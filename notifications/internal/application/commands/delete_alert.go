package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/notifications/internal/domain"
)

type DeleteAlert struct {
	AlertID string
	UserID  string // To verify ownership
}

type DeleteAlertHandler struct {
	alerts    domain.AlertRepository
	catalog   domain.CatalogRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewDeleteAlertHandler(alerts domain.AlertRepository, catalog domain.CatalogRepository, publisher ddd.EventPublisher[ddd.Event]) DeleteAlertHandler {
	return DeleteAlertHandler{
		alerts:    alerts,
		catalog:   catalog,
		publisher: publisher,
	}
}

func (h DeleteAlertHandler) DeleteAlert(ctx context.Context, cmd DeleteAlert) error {
	// First verify the alert exists and belongs to the user
	alert, err := h.catalog.Find(ctx, cmd.AlertID)
	if err != nil {
		return errors.Wrap(err, "finding alert")
	}
	
	if alert == nil {
		return domain.ErrAlertNotFound
	}
	
	// Verify ownership
	if alert.UserID != cmd.UserID {
		return errors.Wrap(errors.ErrUnauthorized, "alert does not belong to user")
	}
	
	// Remove from catalog (read model)
	if err := h.catalog.Remove(ctx, cmd.AlertID); err != nil {
		return errors.Wrap(err, "removing alert from catalog")
	}
	
	// TODO: Consider if we want to delete from event store or just mark as deleted
	// For now, just remove from the read model
	
	return nil
}