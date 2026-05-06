package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/support/internal/domain"
	"time"
)

type AddInternalNoteHandler struct {
	tickets        domain.TicketRepository
	communications domain.CommunicationRepository
	publisher      ddd.EventPublisher[ddd.Event]
}

func NewAddInternalNoteHandler(
	tickets domain.TicketRepository,
	communications domain.CommunicationRepository,
	publisher ddd.EventPublisher[ddd.Event],
) AddInternalNoteHandler {
	return AddInternalNoteHandler{
		tickets:        tickets,
		communications: communications,
		publisher:      publisher,
	}
}

func (h AddInternalNoteHandler) AddInternalNote(ctx context.Context, cmd AddInternalNote) error {
	// Verify ticket exists
	_, err := h.tickets.Load(ctx, cmd.TicketID)
	if err != nil {
		return err
	}

	// Create internal note (non-public communication)
	comm := &domain.Communication{
		ID:             cmd.ID,
		TicketID:       cmd.TicketID,
		AuthorID:       cmd.AuthorID,
		AuthorType:     domain.AuthorTypeAgent,
		Content:        cmd.Content,
		IsPublic:       false, // Internal notes are never public
		CreatedAt:      time.Now(),
		MentionedUsers: cmd.MentionedUsers,
		Metadata:       map[string]string{"type": "internal_note"},
	}

	if err := h.communications.Add(ctx, comm); err != nil {
		return err
	}

	// Publish event - this is not an aggregate event, just an integration event
	event := ddd.NewEvent(domain.InternalNoteAddedEvent, &domain.InternalNoteAdded{
		ID:             cmd.ID,
		TicketID:       cmd.TicketID,
		AuthorID:       cmd.AuthorID,
		Content:        cmd.Content,
		MentionedUsers: cmd.MentionedUsers,
	})
	
	return h.publisher.Publish(ctx, event)
}