package domain

import (
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const MediaAggregate = "media.Media"

var (
	ErrItemIdIsBlank       = errors.Wrap(errors.ErrBadRequest, "the item id cannot be blank")
	ErrItemBlankIsBlank    = errors.Wrap(errors.ErrBadRequest, "the store location cannot be blank")
	ErrUserIDBlank         = errors.Wrap(errors.ErrBadRequest, "the user id cannot be blank")
	ErrStatusParticipating = errors.Wrap(errors.ErrBadRequest, "the store is already not participating")
)

// MediaOrderItem holds info about a single item in the "display order" list/map.
type MediaOrderItem struct {
	MediaItemID string
	URL         string
}
type Media struct {
	es.Aggregate
	ItemID     string
	ItemType   ItemType
	UserID     string
	Status     MediaStatus
	MediaOrder map[int]MediaOrderItem
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Media)(nil)

func NewMedia(id string) *Media {
	return &Media{
		Aggregate: es.NewAggregate(id, MediaAggregate),
	}
}

// Example domain method to add a single item to MediaOrder:
func (m *Media) AddMediaItem(order int, itemID, url string) (ddd.Event, error) {
	// Very simple approach:
	// - If there's already an occupant at "order", you shift them up by 1, etc.
	// - Or just overwrite if that's your domain rule (adjust as needed).

	// For example, detect occupant:
	occupant, found := m.MediaOrder[order]
	if found {
		// SHIFT occupant up by 1, and so on...
		// This is just an example; you might do more advanced shifting logic
		// or reassign occupant. Let's do a simple loop:
		newOrder := order
		for {
			if _, exists := m.MediaOrder[newOrder]; !exists {
				break
			}
			newOrder++
		}
		// Move occupant to 'newOrder'
		m.MediaOrder[newOrder] = occupant
		delete(m.MediaOrder, order)
	}

	// Insert the new item
	m.MediaOrder[order] = MediaOrderItem{
		MediaItemID: itemID,
		URL:         url,
	}

	// Potentially emit an event, e.g. MediaItemAddedEvent, if you want
	// m.AddEvent("media.MediaItemAdded", &MediaItemAdded{...})

	return nil, nil
}
func (m *Media) InitMedia(itemID string, itemType ItemType, userID string, status MediaStatus) (ddd.Event, error) {
	if itemID == "" {
		return nil, ErrItemIdIsBlank
	}

	if userID == "" {
		return nil, ErrUserIDBlank
	}

	m.AddEvent(MediaCreatedEvent, &MediaCreated{
		ItemID:   itemID,
		ItemType: itemType,
		UserID:   userID,
		Status:   status,
	})

	return ddd.NewEvent(MediaCreatedEvent, m), nil
}
func (m *Media) UpdateMedia(itemID string, itemType ItemType, userID string, status MediaStatus) (ddd.Event, error) {
	if itemID == "" {
		return nil, ErrItemIdIsBlank
	}

	if userID == "" {
		return nil, ErrUserIDBlank
	}

	m.AddEvent(MediaUpdatedEvent, &MediaUpdated{
		ItemID:   itemID,
		ItemType: itemType,
		UserID:   userID,
		Status:   status,
	})

	return ddd.NewEvent(MediaUpdatedEvent, m), nil
}
func (m *Media) Delete(mediaId, userID string) (ddd.Event, error) {
	m.AddEvent(MediaDeletedEvent, &MediaDeleted{
		ID:     mediaId,
		UserID: userID,
	})
	return ddd.NewEvent(MediaDeletedEvent, m), nil
}

// Key implements registry.Registerable
func (Media) Key() string { return MediaAggregate }

// ApplyEvent implements es.EventApplier
func (m *Media) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *MediaCreated:
		m.ItemID = payload.ItemID
		m.ItemType = payload.ItemType
		m.UserID = payload.UserID
		m.Status = payload.Status
	case *MediaUpdated:
		m.ItemID = payload.ItemID
		m.ItemType = payload.ItemType
		m.UserID = payload.UserID
		m.Status = payload.Status
	case *MediaDeleted:

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", m, event.EventName(), payload)
	}

	return nil
}

// ApplySnapshot implements es.Snapshotter
func (m *Media) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *MediaV1:
		m.ItemID = ss.ItemID
		m.ItemType = ss.ItemType
		m.UserID = ss.UserID
		m.Status = ss.Status

	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", m, snapshot)
	}

	return nil
}

// ToSnapshot implements es.Snapshotter
func (s Media) ToSnapshot() es.Snapshot {
	return MediaV1{
		ItemID:   s.ItemID,
		ItemType: s.ItemType,
		UserID:   s.UserID,
		Status:   s.Status,
	}
}
