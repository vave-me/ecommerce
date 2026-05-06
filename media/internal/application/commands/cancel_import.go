package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/media/internal/domain"
	"time"
)

type (
	CancelImport struct {
		SessionID string
		Reason    string
		UserID    string
	}

	CancelImportHandler struct {
		imports   domain.ImportSessionRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

func NewCancelImportHandler(imports domain.ImportSessionRepository, publisher ddd.EventPublisher[ddd.Event]) CancelImportHandler {
	return CancelImportHandler{
		imports:   imports,
		publisher: publisher,
	}
}

func (h CancelImportHandler) CancelImport(ctx context.Context, cmd CancelImport) error {
	session, err := h.imports.Get(ctx, cmd.SessionID)
	if err != nil {
		return err
	}

	if session.Status == domain.ImportStatusCompleted || session.Status == domain.ImportStatusCancelled {
		return domain.ErrImportAlreadyFinished
	}

	// Update session status
	session.Status = domain.ImportStatusCancelled
	session.CompletedAt = time.Now()
	
	if err := h.imports.Update(ctx, session); err != nil {
		return err
	}

	// Publish cancellation event
	event := &domain.BulkImportCancelled{
		SessionID:                cmd.SessionID,
		Reason:                   cmd.Reason,
		ProcessedAtCancellation:  session.ProcessedImages,
	}

	return h.publisher.Publish(ctx, ddd.NewEvent(event))
}