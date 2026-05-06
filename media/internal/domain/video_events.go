package domain

const (
	VideoAddedEvent               = "media.VideoAdded"
	VideoRemovedEvent             = "media.VideoRemoved"
	VideoDisplayOrderChangedEvent = "media.VideoDisplayOrderChanged"
	VideoMetadataUpdatedEvent     = "media.VideoMetadataUpdated"
	MainVideoSetEvent             = "media.MainVideoSet"
)

type VideoAdded struct {
	ID           string
	MediaID      string
	DisplayOrder int
	IsMain       bool
	URL          string
	MetaData     string
	Thumbnail    string
	FileType     string
	UserID       string
}

// Key implements registry.Registerable
func (VideoAdded) Key() string { return VideoAddedEvent }

type VideoRemoved struct {
	ID      string
	MediaID string
}

// Key implements registry.Registerable
func (VideoRemoved) Key() string { return VideoRemovedEvent }

type VideoDisplayOrderChanged struct {
	ID           string
	MediaID      string
	DisplayOrder int
	IsMain       bool
}

// Key implements registry.Registerable
func (VideoDisplayOrderChanged) Key() string { return VideoDisplayOrderChangedEvent }

type VideoMetadataUpdated struct {
	ID       string
	MediaID  string
	MetaData string
}

func (VideoMetadataUpdated) Key() string { return VideoMetadataUpdatedEvent }

type MainVideoSet struct {
	ID      string
	MediaID string
	IsMain  bool
}

func (MainVideoSet) Key() string { return MainVideoSetEvent }
