package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/support/internal/domain"
)

type UpdateTicketPriorityHandler struct {
	tickets   domain.TicketRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewUpdateTicketPriorityHandler(
	tickets domain.TicketRepository,
	publisher ddd.EventPublisher[ddd.Event],
) UpdateTicketPriorityHandler {
	return UpdateTicketPriorityHandler{
		tickets:   tickets,
		publisher: publisher,
	}
}

func (h UpdateTicketPriorityHandler) UpdateTicketPriority(ctx context.Context, cmd UpdateTicketPriority) error {
	ticket, err := h.tickets.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	// Get user ID from context
	userID := getUserIDFromContext(ctx)
	
	event, err := ticket.UpdatePriority(cmd.Priority, userID, cmd.Reason)
	if err != nil {
		return err
	}
	h.publisher.Publish(ctx, event)

	return h.tickets.Save(ctx, ticket)
}

func getUserIDFromContext(ctx context.Context) string {
	// TODO: Extract from auth context
	return "user-id"
}