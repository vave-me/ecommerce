package domain

import (
	"context"
	"middleman/managers/internal/models"
)

type SupportRepository interface {
	// Core support operations from gRPC service
	StartSupport(ctx context.Context, userID string) (*models.StartSupportResponse, error)
	CreateTicket(ctx context.Context, supportID, title, description string) (*models.CreateTicketResponse, error)
	UpdateTicket(ctx context.Context, ticketID, status, assignedTo string) (*models.UpdateTicketResponse, error)
	CloseTicket(ctx context.Context, ticketID, reason string) (*models.CloseTicketResponse, error)
	GetTicket(ctx context.Context, ticketID string) (*models.GetTicketResponse, error)
	GetTickets(ctx context.Context, supportID string) (*models.GetTicketsResponse, error)
	GetUserSupport(ctx context.Context, userID string) (*models.GetUserSupportResponse, error)

	// Additional query methods for AI tooling
	GetSupportByID(ctx context.Context, supportID string) (*models.Support, error)
	GetSupportsByUser(ctx context.Context, userID string, limit int64) ([]*models.Support, error)
	GetTicketsByStatus(ctx context.Context, status string, limit int64) ([]*models.Ticket, error)
	SearchTickets(ctx context.Context, query string, limit int64) ([]*models.Ticket, error)
	GetOpenTickets(ctx context.Context, limit int64) ([]*models.Ticket, error)
	GetClosedTickets(ctx context.Context, limit int64) ([]*models.Ticket, error)
	AssignTicket(ctx context.Context, ticketID, assignedTo string) error
	AddTicketComment(ctx context.Context, ticketID, comment, authorID string) error
	GetTicketComments(ctx context.Context, ticketID string) ([]*models.TicketComment, error)
	EscalateTicket(ctx context.Context, ticketID, reason string) error
	GetTicketHistory(ctx context.Context, ticketID string) ([]*models.TicketHistoryEntry, error)
}
