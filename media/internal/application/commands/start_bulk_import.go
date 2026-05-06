package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/media/internal/domain"
)

type (
	StartBulkImport struct {
		SessionID            string
		ExternalSystemID     string
		ExternalSystemType   string
		EstimatedCount       int32
		Options              map[string]string
		UserID               string
	}

	StartBulkImportHandler struct {
		importers domain.ImporterRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

func NewStartBulkImportHandler(importers domain.ImporterRepository, publisher ddd.EventPublisher[ddd.Event]) StartBulkImportHandler {
	return StartBulkImportHandler{
		importers: importers,
		publisher: publisher,
	}
}

func (h StartBulkImportHandler) StartBulkImport(ctx context.Context, cmd StartBulkImport) error {
	importer, err := h.importers.Load(ctx, cmd.SessionID)
	if err != nil {
		return err
	}

	event, err := importer.StartImport(
		cmd.ExternalSystemID,
		cmd.ExternalSystemType,
		cmd.EstimatedCount,
		cmd.UserID,
		cmd.Options,
	)
	if err != nil {
		return err
	}

	if err = h.importers.Save(ctx, importer); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}