package domain

type MediaV1 struct {
	ItemID   string
	ItemType ItemType
	UserID   string
	Status   MediaStatus
}

func (MediaV1) SnapshotName() string { return "media.MediaV1" }
