package domain

import (
	"fmt"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const ImageAggregate = "media.Image"

type Image struct {
	es.Aggregate
	MediaID      string
	DisplayOrder int
	IsMain       bool
	URL          string
	MetaData     string
	FileType     string
	Thumbnail    string
	UserID       string
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Image)(nil)

func NewImage(id string) *Image {
	return &Image{
		Aggregate: es.NewAggregate(id, ImageAggregate),
	}
}

func (i *Image) Remove(id, mediaID string) (ddd.Event, error) {
	i.AddEvent(ImageRemovedEvent, &ImageRemoved{
		ID:      id,
		MediaID: mediaID,
	})
	return ddd.NewEvent(ImageRemovedEvent, i), nil
}
func (i *Image) InitImage(id, mediaID string, displayOrder int, isMain bool, url, metaData, fileType string, thumbnail, userID string) (ddd.Event, error) {
	if mediaID == "" {
		return nil, ErrMediaIdIsBlank
	}

	if displayOrder < 0 {
		return nil, ErrDisplayOrderIsNegative
	}
	fmt.Printf("metadata %s ", metaData)
	i.AddEvent(ImageAddedEvent, &ImageAdded{
		ID:           id,
		MediaID:      mediaID,
		URL:          url,
		DisplayOrder: displayOrder,
		IsMain:       isMain,
		MetaData:     metaData,
		FileType:     fileType,
		Thumbnail:    thumbnail,
		UserID:       userID,
	})

	return ddd.NewEvent(ImageAddedEvent, i), nil
}

// Key implements registry.Registerable
func (Image) Key() string { return ImageAggregate }
func (i *Image) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *ImageAdded:
		i.MediaID = payload.MediaID
		i.DisplayOrder = payload.DisplayOrder
		i.IsMain = payload.IsMain
		i.URL = payload.URL
		i.MetaData = payload.MetaData
		i.FileType = payload.FileType
		i.Thumbnail = payload.Thumbnail
		i.UserID = payload.UserID

	case *ImageDisplayOrderChanged:
		i.DisplayOrder = payload.DisplayOrder

	case *MainImageSet:
		i.IsMain = payload.IsMain

	case *ImageMetadataUpdated:
		i.MetaData = payload.MetaData

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", i, event.EventName(), payload)
	}

	return nil
}
func (i *Image) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *ImageV1:
		i.MediaID = ss.MediaID
		i.DisplayOrder = ss.DisplayOrder
		i.IsMain = ss.IsMain
		i.URL = ss.URL
		i.MetaData = ss.MetaData
		i.FileType = ss.FileType
		i.Thumbnail = ss.Thumbnail
		i.UserID = ss.UserID
	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", i, snapshot)
	}

	return nil
}

func (i Image) ToSnapshot() es.Snapshot {
	return ImageV1{
		MediaID:      i.MediaID,
		DisplayOrder: i.DisplayOrder,
		IsMain:       i.IsMain,
		URL:          i.URL,
		MetaData:     i.MetaData,
		FileType:     i.FileType,
		Thumbnail:    i.Thumbnail,
		UserID:       i.UserID,
	}
}
