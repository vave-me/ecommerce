package domain

import (
	"time"

	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const ImporterAggregate = "media.Importer"

var (
	ErrImporterAlreadyStarted   = errors.Wrap(errors.ErrBadRequest, "importer already started")
	ErrImporterNotActive        = errors.Wrap(errors.ErrBadRequest, "importer is not active")
	ErrImporterAlreadyCompleted = errors.Wrap(errors.ErrBadRequest, "importer already completed")
	ErrExternalSystemIDBlank    = errors.Wrap(errors.ErrBadRequest, "external system id cannot be blank")
	ErrInvalidBatchSize         = errors.Wrap(errors.ErrBadRequest, "invalid batch size")
)

type Importer struct {
	es.Aggregate
	ExternalSystemID   string
	ExternalSystemType string
	TotalImages        int32
	ProcessedImages    int32
	FailedImages       int32
	Status             ImportStatus
	StartedAt          time.Time
	CompletedAt        time.Time
	Metadata           map[string]string
	BatchCount         int32
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Importer)(nil)

func NewImporter(id string) *Importer {
	return &Importer{
		Aggregate: es.NewAggregate(id, ImporterAggregate),
		Metadata:  make(map[string]string),
	}
}

// StartImport initializes a new import session
func (i *Importer) StartImport(externalSystemID, externalSystemType string, estimatedCount int32, userID string, options map[string]string) (ddd.Event, error) {
	if i.Status != "" {
		return nil, ErrImporterAlreadyStarted
	}

	if externalSystemID == "" {
		return nil, ErrExternalSystemIDBlank
	}

	i.AddEvent(BulkImportStartedEvent, &BulkImportStarted{
		SessionID:          i.ID(),
		ExternalSystemID:   externalSystemID,
		ExternalSystemType: externalSystemType,
		EstimatedCount:     estimatedCount,
		UserID:             userID,
	})

	// Store options in metadata
	for k, v := range options {
		i.Metadata[k] = v
	}

	return ddd.NewEvent(BulkImportStartedEvent, i), nil
}

// AddBatch records a batch of items being added to the import queue
func (i *Importer) AddBatch(batchSize int32, batchNumber int32) (ddd.Event, error) {
	if i.Status != ImportStatusPending && i.Status != ImportStatusProcessing {
		return nil, ErrImporterNotActive
	}

	if batchSize <= 0 {
		return nil, ErrInvalidBatchSize
	}

	// Update status to processing if still pending
	if i.Status == ImportStatusPending {
		i.Status = ImportStatusProcessing
	}

	i.AddEvent(ImportBatchAddedEvent, &ImportBatchAdded{
		SessionID:   i.ID(),
		BatchSize:   batchSize,
		BatchNumber: batchNumber,
	})

	return ddd.NewEvent(ImportBatchAddedEvent, i), nil
}

// ProcessItem records successful processing of an import item
func (i *Importer) ProcessItem(externalID, mediaID, imageID string) (ddd.Event, error) {
	if i.Status != ImportStatusProcessing {
		return nil, ErrImporterNotActive
	}

	i.AddEvent(ImportItemProcessedEvent, &ImportItemProcessed{
		SessionID:  i.ID(),
		ExternalID: externalID,
		MediaID:    mediaID,
		ImageID:    imageID,
		Status:     "completed",
	})

	return ddd.NewEvent(ImportItemProcessedEvent, i), nil
}

// FailItem records failed processing of an import item
func (i *Importer) FailItem(externalID, errorCode, errorMessage string, retryCount int32) (ddd.Event, error) {
	if i.Status != ImportStatusProcessing {
		return nil, ErrImporterNotActive
	}

	i.AddEvent(ImportItemFailedEvent, &ImportItemFailed{
		SessionID:    i.ID(),
		ExternalID:   externalID,
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
		RetryCount:   retryCount,
	})

	return ddd.NewEvent(ImportItemFailedEvent, i), nil
}

// Complete marks the import as successfully completed
func (i *Importer) Complete() (ddd.Event, error) {
	if i.Status == ImportStatusCompleted || i.Status == ImportStatusCancelled {
		return nil, ErrImporterAlreadyCompleted
	}

	duration := time.Since(i.StartedAt).Milliseconds()

	i.AddEvent(BulkImportCompletedEvent, &BulkImportCompleted{
		SessionID:         i.ID(),
		TotalProcessed:    i.ProcessedImages + i.FailedImages,
		SuccessfulImports: i.ProcessedImages,
		FailedImports:     i.FailedImages,
		DurationMs:        duration,
	})

	return ddd.NewEvent(BulkImportCompletedEvent, i), nil
}

// Cancel cancels an active import
func (i *Importer) Cancel(reason string) (ddd.Event, error) {
	if i.Status == ImportStatusCompleted || i.Status == ImportStatusCancelled {
		return nil, ErrImporterAlreadyCompleted
	}

	i.AddEvent(BulkImportCancelledEvent, &BulkImportCancelled{
		SessionID:               i.ID(),
		Reason:                  reason,
		ProcessedAtCancellation: i.ProcessedImages,
	})

	return ddd.NewEvent(BulkImportCancelledEvent, i), nil
}

// Key implements registry.Registerable
func (Importer) Key() string { return ImporterAggregate }

// ApplyEvent implements es.EventApplier
func (i *Importer) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *BulkImportStarted:
		i.ExternalSystemID = payload.ExternalSystemID
		i.ExternalSystemType = payload.ExternalSystemType
		i.TotalImages = payload.EstimatedCount
		i.ProcessedImages = 0
		i.FailedImages = 0
		i.Status = ImportStatusPending
		i.StartedAt = time.Now()
		i.BatchCount = 0

	case *ImportBatchAdded:
		i.BatchCount++
		i.TotalImages += payload.BatchSize
		if i.Status == ImportStatusPending {
			i.Status = ImportStatusProcessing
		}

	case *ImportItemProcessed:
		i.ProcessedImages++

	case *ImportItemFailed:
		i.FailedImages++

	case *BulkImportCompleted:
		i.Status = ImportStatusCompleted
		i.CompletedAt = time.Now()

	case *BulkImportCancelled:
		i.Status = ImportStatusCancelled
		i.CompletedAt = time.Now()

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", i, event.EventName(), payload)
	}

	return nil
}

// ApplySnapshot implements es.Snapshotter
func (i *Importer) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *ImporterV1:
		i.ExternalSystemID = ss.ExternalSystemID
		i.ExternalSystemType = ss.ExternalSystemType
		i.TotalImages = ss.TotalImages
		i.ProcessedImages = ss.ProcessedImages
		i.FailedImages = ss.FailedImages
		i.Status = ss.Status
		i.StartedAt = ss.StartedAt
		i.CompletedAt = ss.CompletedAt
		i.Metadata = ss.Metadata
		i.BatchCount = ss.BatchCount

	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", i, snapshot)
	}

	return nil
}

// ToSnapshot implements es.Snapshotter
func (i Importer) ToSnapshot() es.Snapshot {
	return ImporterV1{
		ExternalSystemID:   i.ExternalSystemID,
		ExternalSystemType: i.ExternalSystemType,
		TotalImages:        i.TotalImages,
		ProcessedImages:    i.ProcessedImages,
		FailedImages:       i.FailedImages,
		Status:             i.Status,
		StartedAt:          i.StartedAt,
		CompletedAt:        i.CompletedAt,
		Metadata:           i.Metadata,
		BatchCount:         i.BatchCount,
	}
}