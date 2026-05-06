package domain

const (
	// Media Entity Lifecycle Events
	MediaCreatedEvent       = "media.MediaCreated"
	MediaUpdatedEvent       = "media.MediaUpdated"
	MediaDeletedEvent       = "media.MediaDeleted"
	MediaStatusChangedEvent = "media.MediaStatusChanged"
)

type MediaCreated struct {
	ItemID   string
	ItemType ItemType
	UserID   string
	Status   MediaStatus
}

// Key implements registry.Registerable
func (MediaCreated) Key() string { return MediaCreatedEvent }

type MediaUpdated struct {
	ItemID   string
	ItemType ItemType
	UserID   string
	Status   MediaStatus
}

// Key implements registry.Registerable
func (MediaUpdated) Key() string { return MediaUpdatedEvent }

type MediaDeleted struct {
	ID     string
	UserID string
}

// Key implements registry.Registerable
func (MediaDeleted) Key() string { return MediaDeletedEvent }

type MediaStatusChanged struct {
}

// Key implements registry.Registerable
func (MediaStatusChanged) Key() string { return MediaStatusChangedEvent }
