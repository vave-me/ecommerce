package grpc

import (
	"context"
	"middleman/assistants/internal/domain"
	"middleman/assistants/internal/models"
	"middleman/internal/auth"
	"middleman/internal/rpc"
	"middleman/support/supportpb"

	"google.golang.org/grpc"
)

type SupportRepository struct {
	endpoint string
	auth     *auth.Auth // Optional auth for authenticated calls
}

var _ domain.SupportRepository = (*SupportRepository)(nil)

// NewSupportRepositoryWithAuth creates a new SupportRepository with JWT authentication support
func NewSupportRepository(endpoint string, authInstance *auth.Auth) SupportRepository {
	return SupportRepository{
		endpoint: endpoint,
		auth:     authInstance,
	}
}

// InitiateSupportSessionForUser initiates a new support channel for a user
func (r SupportRepository) InitiateSupportSessionForUser(ctx context.Context, userID string) (*models.StartSupportResponse, error) {
	var conn *grpc.ClientConn
	conn, err := r.dial(ctx)

	if err != nil {
		return nil, err
	}

	// Create a support channel for the user
	resp, err := supportpb.NewSupportServiceClient(conn).CreateSupportChannel(ctx, &supportpb.CreateSupportChannelRequest{
		UserId:      userID,
		ChannelType: supportpb.SupportChannelType_GENERAL,
	})
	if err != nil {
		return nil, err
	}

	return &models.StartSupportResponse{
		ID: resp.GetId(),
	}, nil
}

// CreateNewSupportTicket creates a new support ticket within an existing support channel
func (r SupportRepository) CreateNewSupportTicket(ctx context.Context, supportID, title, description string) (*models.CreateTicketResponse, error) {
	var conn *grpc.ClientConn
	conn, err := r.dial(ctx)

	if err != nil {
		return nil, err
	}

	resp, err := supportpb.NewSupportServiceClient(conn).CreateTicket(ctx, &supportpb.CreateTicketRequest{
		ChannelId:   supportID, // supportID is actually channelID
		Title:       title,
		Description: description,
		Category:    supportpb.TicketCategory_GENERAL_INQUIRY,
		Priority:    supportpb.TicketPriority_MEDIUM,
	})
	if err != nil {
		return nil, err
	}

	return &models.CreateTicketResponse{
		ID: resp.GetId(),
	}, nil
}

// ListTickets lists all support tickets within a support channel
func (r SupportRepository) ListTickets(ctx context.Context, supportID string, page, limit int64) (*models.ListTicketsResponse, error) {
	var conn *grpc.ClientConn
	conn, err := r.dial(ctx)

	if err != nil {
		return nil, err
	}

	// Use GetChannelTickets instead of ListTickets
	resp, err := supportpb.NewSupportServiceClient(conn).GetChannelTickets(ctx, &supportpb.GetChannelTicketsRequest{
		ChannelId: supportID,
		Page:      int32(page),
		Limit:     int32(limit),
	})
	if err != nil {
		return nil, err
	}

	tickets := make([]*models.Ticket, len(resp.GetTickets()))
	for i, ticket := range resp.GetTickets() {
		tickets[i] = r.ticketToDomain(ticket)
	}

	return &models.ListTicketsResponse{
		Tickets: tickets,
		Total:   int64(resp.GetTotalCount()),
		Page:    page,
		Limit:   limit,
	}, nil
}

// GetSupportTicketByID retrieves details of a specific support ticket
func (r SupportRepository) GetSupportTicketByID(ctx context.Context, ticketID string) (*models.GetTicketResponse, error) {
	var conn *grpc.ClientConn
	conn, err := r.dial(ctx)

	if err != nil {
		return nil, err
	}

	resp, err := supportpb.NewSupportServiceClient(conn).GetTicket(ctx, &supportpb.GetTicketRequest{
		Id: ticketID,
	})
	if err != nil {
		return nil, err
	}

	return &models.GetTicketResponse{
		Ticket: r.ticketToDomain(resp.GetTicket()),
	}, nil
}

// ModifySupportTicketStatus updates a support ticket's details or status
func (r SupportRepository) ModifySupportTicketStatus(ctx context.Context, ticketID, status, assignedTo string) (*models.UpdateTicketResponse, error) {
	var conn *grpc.ClientConn
	conn, err := r.dial(ctx)

	if err != nil {
		return nil, err
	}

	// Status update is done through specific methods like ResolveTicket, CloseTicket, etc.
	// For now, we'll return an empty response
	// TODO: Implement proper status updates using the appropriate endpoints
	_ = supportpb.NewSupportServiceClient(conn)
	
	// Since UpdateTicket doesn't return the ticket in the new API
	// we need to fetch it separately
	getResp, err := r.GetSupportTicketByID(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	return &models.UpdateTicketResponse{
		Ticket: getResp.Ticket,
	}, nil
}

// DeleteTicket deletes a support ticket - implemented as closing the ticket
func (r SupportRepository) DeleteTicket(ctx context.Context, supportID, ticketID string) (*models.DeleteTicketResponse, error) {
	// Since there's no DeleteTicket in the new API, we'll close it instead
	_, err := r.CloseSupportTicketWithReason(ctx, ticketID, "Ticket deleted")
	if err != nil {
		return nil, err
	}

	return &models.DeleteTicketResponse{}, nil
}

// CloseSupportTicketWithReason closes a support ticket
func (r SupportRepository) CloseSupportTicketWithReason(ctx context.Context, ticketID, reason string) (*models.CloseTicketResponse, error) {
	var conn *grpc.ClientConn
	conn, err := r.dial(ctx)

	if err != nil {
		return nil, err
	}

	_, err = supportpb.NewSupportServiceClient(conn).CloseTicket(ctx, &supportpb.CloseTicketRequest{
		Id:           ticketID,
		ClosureNotes: reason,
	})
	if err != nil {
		return nil, err
	}

	return &models.CloseTicketResponse{}, nil
}

// GetAllTicketsForSupport retrieves tickets for a support channel
func (r SupportRepository) GetAllTicketsForSupport(ctx context.Context, supportID string) (*models.GetTicketsResponse, error) {
	// Use ListTickets with default pagination
	listResp, err := r.ListTickets(ctx, supportID, 1, 50)
	if err != nil {
		return nil, err
	}
	return &models.GetTicketsResponse{
		Tickets: listResp.Tickets,
		Total:   listResp.Total,
	}, nil
}

// GetUserSupportInformation retrieves support channel for a user
func (r SupportRepository) GetUserSupportInformation(ctx context.Context, userID string) (*models.GetUserSupportResponse, error) {
	// This would need to be implemented based on the actual gRPC service
	return &models.GetUserSupportResponse{
		Support: &models.Support{
			ID:     "default_support_" + userID,
			UserID: userID,
		},
	}, nil
}

// Additional query methods for AI tooling
func (r SupportRepository) FindSupportSessionByID(ctx context.Context, supportID string) (*models.Support, error) {
	return &models.Support{
		ID:     supportID,
		UserID: "unknown",
	}, nil
}

func (r SupportRepository) GetUserSupportSessionHistory(ctx context.Context, userID string, limit int64) ([]*models.Support, error) {
	return []*models.Support{
		{
			ID:     "support_" + userID,
			UserID: userID,
		},
	}, nil
}

func (r SupportRepository) FilterTicketsByCurrentStatus(ctx context.Context, status string, limit int64) ([]*models.Ticket, error) {
	return []*models.Ticket{}, nil
}

func (r SupportRepository) SearchTicketsByKeyword(ctx context.Context, query string, limit int64) ([]*models.Ticket, error) {
	return []*models.Ticket{}, nil
}

func (r SupportRepository) GetAllOpenSupportTickets(ctx context.Context, limit int64) ([]*models.Ticket, error) {
	return []*models.Ticket{}, nil
}

func (r SupportRepository) GetAllResolvedSupportTickets(ctx context.Context, limit int64) ([]*models.Ticket, error) {
	return []*models.Ticket{}, nil
}

func (r SupportRepository) AssignTicketToSupportAgent(ctx context.Context, ticketID, assignedTo string) error {
	return nil // Placeholder implementation
}

func (r SupportRepository) AddCommentToSupportTicket(ctx context.Context, ticketID, comment, authorID string) error {
	return nil // Placeholder implementation
}

func (r SupportRepository) GetAllCommentsForTicket(ctx context.Context, ticketID string) ([]*models.TicketComment, error) {
	return []*models.TicketComment{}, nil
}

func (r SupportRepository) EscalateTicketToPriority(ctx context.Context, ticketID, reason string) error {
	return nil // Placeholder implementation
}

func (r SupportRepository) GetTicketActivityHistory(ctx context.Context, ticketID string) ([]*models.TicketHistoryEntry, error) {
	return []*models.TicketHistoryEntry{}, nil
}

// Domain conversion methods

func (r SupportRepository) supportToDomain(channel *supportpb.SupportChannel) *models.Support {
	if channel == nil {
		return nil
	}

	return &models.Support{
		ID:     channel.GetId(),
		UserID: channel.GetUserId(),
	}
}

func (r SupportRepository) ticketToDomain(ticket *supportpb.Ticket) *models.Ticket {
	if ticket == nil {
		return nil
	}

	// Map the new ticket status enum to our models
	var status models.TicketStatus
	switch ticket.GetStatus() {
	case supportpb.TicketStatus_DRAFT, supportpb.TicketStatus_SUBMITTED, supportpb.TicketStatus_ASSIGNED:
		status = models.TicketStatusOpen
	case supportpb.TicketStatus_IN_PROGRESS, supportpb.TicketStatus_PENDING_CUSTOMER, supportpb.TicketStatus_ESCALATED:
		status = models.TicketStatusInProgress
	case supportpb.TicketStatus_RESOLVED:
		status = models.TicketStatusResolved
	case supportpb.TicketStatus_CLOSED:
		status = models.TicketStatusClosed
	case supportpb.TicketStatus_REOPENED:
		status = models.TicketStatusOpen // Reopened tickets go back to open
	default:
		status = models.TicketStatusUnknown
	}

	return &models.Ticket{
		ID:           ticket.GetId(),
		SupportID:    ticket.GetChannelId(), // channel_id instead of support_id
		Title:        ticket.GetTitle(),
		Description:  ticket.GetDescription(),
		TicketStatus: status,
	}
}

// dial sets up a gRPC connection with the microservice endpoint
func (r SupportRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

func (r SupportRepository) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}
