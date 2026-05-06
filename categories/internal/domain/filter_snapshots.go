package domain

type FilterV1 struct {
	ID         string
	CategoryID string
	Name       string
	FilterType FilterType
	Values     []string
	IsActive   bool
}

func (FilterV1) SnapshotName() string { return "categories.FilterV1" }
