package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/support/internal/domain"
)

type ResolveTicket struct {
	ID               string
	ResolvedBy       string
	Resolution       string
	AppliedSolutions []string
}

type ResolveTicketHandler struct {
	tickets   domain.TicketRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewResolveTicketHandler(
	tickets domain.TicketRepository,
	publisher ddd.EventPublisher[ddd.Event],
) ResolveTicketHandler {
	return ResolveTicketHandler{
		tickets:   tickets,
		publisher: publisher,
	}
}

func (h ResolveTicketHandler) ResolveTicket(ctx context.Context, cmd ResolveTicket) error {
	ticket, err := h.tickets.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := ticket.ResolveTicket(cmd.ResolvedBy, cmd.Resolution, cmd.AppliedSolutions)
	if err != nil {
		return err
	}

	err = h.tickets.Save(ctx, ticket)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}