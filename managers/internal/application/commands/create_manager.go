package commands

import (
	"context"

	"github.com/google/uuid"
	"middleman/internal/ddd"
	"middleman/managers/internal/domain"
)

// CreateManager defines the contract for the CreateManager use case
type CreateManager struct {
	ID           string
	Name         string
	Description  string
	UserID       string
	Type         domain.ManagerType
	Capabilities []domain.ManagerCapability
	Temperature  float64
	MaxTokens    int
	SystemPrompt string
}

// CreateManagerHandler defines the contract for handling CreateManager commands
type CreateManagerHandler interface {
	CreateManager(ctx context.Context, cmd CreateManager) error
}

// createManagerHandler is the concrete implementation
type createManagerHandler struct {
	managers  domain.ManagerRepository
	publisher ddd.EventPublisher[ddd.Event]
}

// NewCreateManagerHandler constructs a new CreateManagerHandler
func NewCreateManagerHandler(
	managers domain.ManagerRepository,
	publisher ddd.EventPublisher[ddd.Event],
) CreateManagerHandler {
	return &createManagerHandler{
		managers:  managers,
		publisher: publisher,
	}
}

// CreateManager handles the creation of a new manager
func (h *createManagerHandler) CreateManager(ctx context.Context, cmd CreateManager) error {
	// Generate ID if not provided
	if cmd.ID == "" {
		cmd.ID = uuid.New().String()
	}

	// Load or create the manager aggregate
	manager, err := h.managers.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	// Create the manager - the domain CreateManager method already has all defaults
	event, err := manager.CreateManager(
		cmd.ID,
		cmd.Name,
		cmd.Description,
		cmd.UserID,
		cmd.Type,
		cmd.Capabilities,
		cmd.Temperature,
		cmd.MaxTokens,
		cmd.SystemPrompt,
	)
	if err != nil {
		return err
	}

	// Save the manager
	if err = h.managers.Save(ctx, manager); err != nil {
		return err
	}

	// Publish the event
	return h.publisher.Publish(ctx, event)
}