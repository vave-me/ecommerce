package domain

import "time"

type ActionVi struct {
	SchedulerID         string
	NaturalLanguageTask string
	ExecutionTime       time.Time
	Status              string
	CreatedAt           time.Time
	ExecutedAt          *time.Time
	Result              string
	ErrorMessage        string
}

func (ActionVi) SnapshotName() string { return "scheduler.ActionV1" }
