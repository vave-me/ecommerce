package domain

import "time"

type ImporterV1 struct {
	ExternalSystemID   string
	ExternalSystemType string
	TotalImages        int32
	ProcessedImages    int32
	FailedImages       int32
	Status             ImportStatus
	StartedAt          time.Time
	CompletedAt        time.Time
	Metadata           map[string]string
	BatchCount         int32
}

func (ImporterV1) SnapshotName() string { return "media.ImporterV1" }
