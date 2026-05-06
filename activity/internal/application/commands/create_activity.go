package commands

import (
	"context"
	"fmt"
	"middleman/activity/internal/domain"
	"middleman/internal/ddd"
)

type (
	CreateActivity struct {
		ID     string
		UserID string
	}

	CreateActivityHandler struct {
		activities domain.ActivityRepository
		publisher  ddd.EventPublisher[ddd.Event]
	}
)

func NewCreateActivityHandler(activities domain.ActivityRepository, publisher ddd.EventPublisher[ddd.Event]) CreateActivityHandler {
	return CreateActivityHandler{
		activities: activities,
		publisher:  publisher,
	}
}

func (h CreateActivityHandler) CreateActivity(ctx context.Context, cmd CreateActivity) error {
	interaction, err := h.activities.Load(ctx, cmd.ID)

	if err != nil {
		fmt.Println("Error loading activities")
		return err
	}

	event, err := interaction.InitActivity(cmd.UserID)
	if err != nil {

		fmt.Println("Error initializing activity")
		return err
	}

	err = h.activities.Save(ctx, interaction)
	if err != nil {

		fmt.Println("Error saving activities")
		return err
	}
	fmt.Println("Activities created")
	return h.publisher.Publish(ctx, event)
}
