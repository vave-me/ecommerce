package domain

import (
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const VideoAggregate = "media.Video"

var (
	ErrMediaIdIsBlank         = errors.Wrap(errors.ErrBadRequest, "the product name cannot be blank")
	ErrDisplayOrderIsNegative = errors.Wrap(errors.ErrBadRequest, "the product price cannot be negative")
	ErrIsMain                 = errors.Wrap(errors.ErrBadRequest, "the price change would be a decrease")
	ErrNotMetaData            = errors.Wrap(errors.ErrBadRequest, "the price change would be an increase")
	ErrURl                    = errors.Wrap(errors.ErrBadRequest, "the price change would be an increase")
)

type Video struct {
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
} = (*Video)(nil)

func NewVideo(id string) *Video {
	return &Video{
		Aggregate: es.NewAggregate(id, VideoAggregate),
	}
}

func (v *Video) InitVideo(id, mediaID string, displayOrder int, isMain bool, url string, metaData, fileType, userID string) (ddd.Event, error) {
	if mediaID == "" {
		return nil, ErrMediaIdIsBlank
	}

	if displayOrder < 0 {
		return nil, ErrDisplayOrderIsNegative
	}

	v.AddEvent(VideoAddedEvent, &VideoAdded{
		ID:           id,
		MediaID:      mediaID,
		URL:          url,
		DisplayOrder: displayOrder,
		IsMain:       isMain,
		MetaData:     metaData,
		FileType:     fileType,
		UserID:       userID,
	})

	return ddd.NewEvent(VideoAddedEvent, v), nil
}
func (v *Video) Remove(id, mediaID string) (ddd.Event, error) {
	v.AddEvent(VideoRemovedEvent, &VideoRemoved{
		ID:      id,
		MediaID: mediaID,
	})
	return ddd.NewEvent(VideoRemovedEvent, v), nil
}
func (Video) Key() string { return VideoAggregate }
func (i *Video) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *VideoAdded:
		i.MediaID = payload.MediaID
		i.DisplayOrder = payload.DisplayOrder
		i.IsMain = payload.IsMain
		i.URL = payload.URL
		i.MetaData = payload.MetaData
		i.FileType = payload.FileType

	case *VideoDisplayOrderChanged:
		i.DisplayOrder = payload.DisplayOrder

	case *MainVideoSet:
		i.IsMain = payload.IsMain

	case *VideoMetadataUpdated:
		i.MetaData = payload.MetaData

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", i, event.EventName(), payload)
	}

	return nil
}
func (i *Video) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *VideoV1:
		i.MediaID = ss.MediaID
		i.DisplayOrder = ss.DisplayOrder
		i.IsMain = ss.IsMain
		i.URL = ss.URL
		i.MetaData = ss.MetaData
		i.FileType = ss.FileType

	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", i, snapshot)
	}

	return nil
}

func (i Video) ToSnapshot() es.Snapshot {
	return VideoV1{
		MediaID:      i.MediaID,
		DisplayOrder: i.DisplayOrder,
		IsMain:       i.IsMain,
		URL:          i.URL,
		MetaData:     i.MetaData,
		FileType:     i.FileType,
	}
}
