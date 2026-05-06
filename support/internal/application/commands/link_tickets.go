package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/support/internal/domain"
)

type LinkTicketsHandler struct {
	tickets   domain.TicketRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewLinkTicketsHandler(
	tickets domain.TicketRepository,
	publisher ddd.EventPublisher[ddd.Event],
) LinkTicketsHandler {
	return LinkTicketsHandler{
		tickets:   tickets,
		publisher: publisher,
	}
}

func (h LinkTicketsHandler) LinkTickets(ctx context.Context, cmd LinkTickets) error {
	// Verify all tickets exist
	_, err := h.tickets.Load(ctx, cmd.TicketID)
	if err != nil {
		return err
	}

	for _, relatedID := range cmd.RelatedTicketIDs {
		_, err := h.tickets.Load(ctx, relatedID)
		if err != nil {
			return err
		}
	}

	// Publish link event - this is not an aggregate event, just an integration event
	linkEvent := ddd.NewEvent(domain.TicketsLinkedEvent, &domain.TicketsLinked{
		TicketID:         cmd.TicketID,
		RelatedTicketIDs: cmd.RelatedTicketIDs,
		LinkedBy:         cmd.LinkedBy,
		RelationshipType: cmd.RelationshipType,
	})
	
	return h.publisher.Publish(ctx, linkEvent)
}