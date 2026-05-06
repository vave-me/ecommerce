package domain

const (
	ImageAddedEvent               = "media.ImageAdded"
	ImageRemovedEvent             = "media.ImageRemoved"
	ImageDisplayOrderChangedEvent = "media.ImageDisplayOrderChanged"
	ImageMetadataUpdatedEvent     = "media.ImageMetadataUpdated"
	MainImageSetEvent             = "media.MainImageSet"
)

type ImageAdded struct {
	ID           string
	MediaID      string
	DisplayOrder int
	IsMain       bool
	URL          string
	MetaData     string
	FileType     string
	Thumbnail    string
	UserID       string
}

// Key implements registry.Registerable
func (ImageAdded) Key() string { return ImageAddedEvent }

type ImageRemoved struct {
	ID      string
	MediaID string
}

// Key implements registry.Registerable
func (ImageRemoved) Key() string { return ImageRemovedEvent }

type ImageDisplayOrderChanged struct {
	ID           string
	MediaID      string
	DisplayOrder int
	IsMain       bool
}

// Key implements registry.Registerable
func (ImageDisplayOrderChanged) Key() string { return ImageDisplayOrderChangedEvent }

type ImageMetadataUpdated struct {
	MediaID  string
	MetaData string
}

func (ImageMetadataUpdated) Key() string { return ImageMetadataUpdatedEvent }

type MainImageSet struct {
	ID      string
	MediaID string
	IsMain  bool
}

func (MainImageSet) Key() string { return MainImageSetEvent }
