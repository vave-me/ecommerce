package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/scheduler/internal/domain"
)

type (
	CreateScheduler struct {
		ID     string
		UserID string
	}

	CreateSchedulerHandler struct {
		activities domain.SchedulerRepository
		publisher  ddd.EventPublisher[ddd.Event]
	}
)

func NewCreateSchedulerHandler(activities domain.SchedulerRepository, publisher ddd.EventPublisher[ddd.Event]) CreateSchedulerHandler {
	return CreateSchedulerHandler{
		activities: activities,
		publisher:  publisher,
	}
}

func (h CreateSchedulerHandler) CreateScheduler(ctx context.Context, cmd CreateScheduler) error {
	scheduler, err := h.activities.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "loading scheduler")
	}

	event, err := scheduler.InitScheduler(cmd.UserID)
	if err != nil {
		return errors.Wrap(err, "initializing scheduler")
	}

	err = h.activities.Save(ctx, scheduler)
	if err != nil {
		return errors.Wrap(err, "saving scheduler")
	}

	return h.publisher.Publish(ctx, event)
}
