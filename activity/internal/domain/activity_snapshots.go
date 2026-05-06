package domain

type ActivityVi struct {
	UserID string
}

func (ActivityVi) SnapshotName() string { return "activity.ActivityV1" }
