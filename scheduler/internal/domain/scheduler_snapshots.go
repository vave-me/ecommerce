package domain

type SchedulerVi struct {
	UserID string
}

func (SchedulerVi) SnapshotName() string { return "scheduler.SchedulerV1" }
