package mediapb

import (
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

const (
	MediaAggregateChannel = "middleman.media.events.Media"
	MediaCreatedEvent     = "mediaapi.MediaCreated"
	MediaUpdatedEvent     = "mediaapi.MediaUpdated"

	ImageAggregateChannel = "middleman.media.events.Image"
	ImageAddedEvent       = "mediaapi.ImageAdded"
	VideoAggregateChannel = "middleman.media.events.Video"
	VideoAddedEvent       = "mediaapi.VideoAdded"

	ImageRemovedEvent = "mediaapi.ImageRemoved"
	MediaRemovedEvent = "mediaapi.MediaRemoved"
	VideoRemovedEvent = "mediaapi.VideoRemoved"

	// Bulk Import events
	ImporterAggregateChannel = "middleman.media.events.Importer"
	BulkImportStartedEvent   = "mediaapi.BulkImportStarted"
	ImportBatchAddedEvent    = "mediaapi.ImportBatchAdded"
	ImportItemProcessedEvent = "mediaapi.ImportItemProcessed"
	ImportItemFailedEvent    = "mediaapi.ImportItemFailed"
	BulkImportCompletedEvent = "mediaapi.BulkImportCompleted"
	BulkImportCancelledEvent = "mediaapi.BulkImportCancelled"
)

func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) error {
	// Store events
	if err := serde.Register(&MediaCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&MediaUpdated{}); err != nil {
		return err
	}
	if err := serde.Register(&MediaRemoved{}); err != nil {
		return err
	}

	if err := serde.Register(&ImageAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&ImageRemoved{}); err != nil {
		return err
	}

	if err := serde.Register(&VideoAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&VideoRemoved{}); err != nil {
		return err
	}

	// Bulk Import events
	if err := serde.Register(&BulkImportStarted{}); err != nil {
		return err
	}
	if err := serde.Register(&ImportBatchAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&ImportItemProcessed{}); err != nil {
		return err
	}
	if err := serde.Register(&ImportItemFailed{}); err != nil {
		return err
	}
	if err := serde.Register(&BulkImportCompleted{}); err != nil {
		return err
	}
	if err := serde.Register(&BulkImportCancelled{}); err != nil {
		return err
	}
	return nil
}

func (*MediaCreated) Key() string { return MediaCreatedEvent }
func (*MediaUpdated) Key() string { return MediaUpdatedEvent }
func (*MediaRemoved) Key() string { return MediaRemovedEvent }

func (*ImageAdded) Key() string   { return ImageAddedEvent }
func (*ImageRemoved) Key() string { return ImageRemovedEvent }

func (*VideoAdded) Key() string   { return VideoAddedEvent }
func (*VideoRemoved) Key() string { return VideoRemovedEvent }

func (*BulkImportStarted) Key() string   { return BulkImportStartedEvent }
func (*ImportBatchAdded) Key() string    { return ImportBatchAddedEvent }
func (*ImportItemProcessed) Key() string { return ImportItemProcessedEvent }
func (*ImportItemFailed) Key() string    { return ImportItemFailedEvent }
func (*BulkImportCompleted) Key() string { return BulkImportCompletedEvent }
func (*BulkImportCancelled) Key() string { return BulkImportCancelledEvent }
