package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/support/internal/domain"
)

type MergeTicketsHandler struct {
	tickets   domain.TicketRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewMergeTicketsHandler(
	tickets domain.TicketRepository,
	publisher ddd.EventPublisher[ddd.Event],
) MergeTicketsHandler {
	return MergeTicketsHandler{
		tickets:   tickets,
		publisher: publisher,
	}
}

func (h MergeTicketsHandler) MergeTickets(ctx context.Context, cmd MergeTickets) error {
	// Verify primary ticket exists
	_, err := h.tickets.Load(ctx, cmd.PrimaryTicketID)
	if err != nil {
		return err
	}

	// Mark secondary tickets as merged
	for _, secondaryID := range cmd.SecondaryTicketIDs {
		secondaryTicket, err := h.tickets.Load(ctx, secondaryID)
		if err != nil {
			return err
		}

		// Close the secondary ticket with a merge note
		mergeNote := "Merged into ticket " + cmd.PrimaryTicketID + ": " + cmd.MergeReason
		event, err := secondaryTicket.CloseTicket(cmd.MergedBy, mergeNote, nil)
		if err != nil {
			return err
		}
		h.publisher.Publish(ctx, event)

		if err := h.tickets.Save(ctx, secondaryTicket); err != nil {
			return err
		}
	}

	// Publish merge event - this is not an aggregate event, just an integration event
	mergeEvent := ddd.NewEvent(domain.TicketsMergedEvent, &domain.TicketsMerged{
		PrimaryTicketID:    cmd.PrimaryTicketID,
		SecondaryTicketIDs: cmd.SecondaryTicketIDs,
		MergedBy:           cmd.MergedBy,
		MergeReason:        cmd.MergeReason,
	})
	h.publisher.Publish(ctx, mergeEvent)

	return nil
}