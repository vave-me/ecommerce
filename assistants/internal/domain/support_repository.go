package domain

import (
	"context"
	"middleman/assistants/internal/models"
)

type SupportRepository interface {
	// Core support operations from gRPC service
	InitiateSupportSessionForUser(ctx context.Context, userID string) (*models.StartSupportResponse, error)
	CreateNewSupportTicket(ctx context.Context, supportID, title, description string) (*models.CreateTicketResponse, error)
	ModifySupportTicketStatus(ctx context.Context, ticketID, status, assignedTo string) (*models.UpdateTicketResponse, error)
	CloseSupportTicketWithReason(ctx context.Context, ticketID, reason string) (*models.CloseTicketResponse, error)
	GetSupportTicketByID(ctx context.Context, ticketID string) (*models.GetTicketResponse, error)
	GetAllTicketsForSupport(ctx context.Context, supportID string) (*models.GetTicketsResponse, error)
	GetUserSupportInformation(ctx context.Context, userID string) (*models.GetUserSupportResponse, error)

	// Additional query methods for AI tooling
	FindSupportSessionByID(ctx context.Context, supportID string) (*models.Support, error)
	GetUserSupportSessionHistory(ctx context.Context, userID string, limit int64) ([]*models.Support, error)
	FilterTicketsByCurrentStatus(ctx context.Context, status string, limit int64) ([]*models.Ticket, error)
	SearchTicketsByKeyword(ctx context.Context, query string, limit int64) ([]*models.Ticket, error)
	GetAllOpenSupportTickets(ctx context.Context, limit int64) ([]*models.Ticket, error)
	GetAllResolvedSupportTickets(ctx context.Context, limit int64) ([]*models.Ticket, error)
	AssignTicketToSupportAgent(ctx context.Context, ticketID, assignedTo string) error
	AddCommentToSupportTicket(ctx context.Context, ticketID, comment, authorID string) error
	GetAllCommentsForTicket(ctx context.Context, ticketID string) ([]*models.TicketComment, error)
	EscalateTicketToPriority(ctx context.Context, ticketID, reason string) error
	GetTicketActivityHistory(ctx context.Context, ticketID string) ([]*models.TicketHistoryEntry, error)
}
