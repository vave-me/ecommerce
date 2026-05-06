package domain

type VideoV1 struct {
	MediaID      string
	DisplayOrder int
	IsMain       bool
	URL          string
	MetaData     string
	FileType     string
	Thumbnail    string
}

func (VideoV1) SnapshotName() string { return "media.VideoV1" }
