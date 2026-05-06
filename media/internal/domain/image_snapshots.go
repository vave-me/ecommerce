package domain

type ImageV1 struct {
	MediaID      string
	DisplayOrder int
	IsMain       bool
	URL          string
	MetaData     string
	FileType     string
	Thumbnail    string
	UserID       string
}

func (ImageV1) SnapshotName() string { return "media.ImageV1" }
