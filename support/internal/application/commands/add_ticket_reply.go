package commands

import (
	"context"
	"time"
	"middleman/internal/ddd"
	"middleman/support/internal/domain"
)

type AddTicketReply struct {
	ID          string
	TicketID    string
	AuthorID    string
	AuthorType  domain.AuthorType
	Content     string
	Attachments []domain.Attachment
	IsPublic    bool
}

type AddTicketReplyHandler struct {
	tickets         domain.TicketRepository
	communications  domain.CommunicationRepository
	publisher       ddd.EventPublisher[ddd.Event]
}

func NewAddTicketReplyHandler(
	tickets domain.TicketRepository,
	communications domain.CommunicationRepository,
	publisher ddd.EventPublisher[ddd.Event],
) AddTicketReplyHandler {
	return AddTicketReplyHandler{
		tickets:        tickets,
		communications: communications,
		publisher:      publisher,
	}
}

func (h AddTicketReplyHandler) AddTicketReply(ctx context.Context, cmd AddTicketReply) error {
	// Verify ticket exists
	ticket, err := h.tickets.Load(ctx, cmd.TicketID)
	if err != nil {
		return err
	}

	// Create communication record
	comm := &domain.Communication{
		ID:          cmd.ID,
		TicketID:    cmd.TicketID,
		AuthorID:    cmd.AuthorID,
		AuthorType:  cmd.AuthorType,
		Content:     cmd.Content,
		IsPublic:    cmd.IsPublic,
		Attachments: cmd.Attachments,
		Metadata:    make(map[string]string),
		CreatedAt:   time.Now(),
	}

	err = h.communications.Add(ctx, comm)
	if err != nil {
		return err
	}

	// Update ticket's first response time if this is the first agent/AI response
	if ticket.FirstResponseAt == nil && 
		(cmd.AuthorType == domain.AuthorTypeAgent || cmd.AuthorType == domain.AuthorTypeAI) {
		now := time.Now()
		ticket.FirstResponseAt = &now
	}

	// Update response count
	ticket.ResponseCount++

	// Update ticket status if needed
	if ticket.Status == domain.TicketStatusPendingCustomer && cmd.AuthorType == domain.AuthorTypeCustomer {
		ticket.Status = domain.TicketStatusInProgress
	} else if ticket.Status == domain.TicketStatusAssigned && 
		(cmd.AuthorType == domain.AuthorTypeAgent || cmd.AuthorType == domain.AuthorTypeAI) {
		ticket.Status = domain.TicketStatusInProgress
	}

	ticket.UpdatedAt = time.Now()
	err = h.tickets.Save(ctx, ticket)
	if err != nil {
		return err
	}

	// Publish event
	event := ddd.NewEvent(domain.TicketReplyAddedEvent, &domain.TicketReplyAdded{
		ID:          cmd.ID,
		TicketID:    cmd.TicketID,
		AuthorID:    cmd.AuthorID,
		AuthorType:  cmd.AuthorType,
		Content:     cmd.Content,
		Attachments: cmd.Attachments,
		IsPublic:    cmd.IsPublic,
	})

	return h.publisher.Publish(ctx, event)
}