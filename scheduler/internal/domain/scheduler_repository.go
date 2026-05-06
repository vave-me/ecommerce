package domain

import "context"

type SchedulerRepository interface {
	Load(ctx context.Context, schedulerID string) (*Scheduler, error)
	Save(ctx context.Context, schedulerID *Scheduler) error
}
