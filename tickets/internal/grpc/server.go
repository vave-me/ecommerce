package grpc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpc_code "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"middleman/internal/auth"
	"middleman/internal/errorsotel"
	"middleman/tickets/internal/application"
	"middleman/tickets/internal/application/commands"
	"middleman/tickets/internal/application/queries"
	"middleman/tickets/internal/domain"
	"middleman/tickets/ticketspb"
)

type server struct {
	app application.App
	ticketspb.UnimplementedTicketsServiceServer
}

var _ ticketspb.TicketsServiceServer = (*server)(nil)

// RegisterServer registers the gRPC server implementation
func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	ticketspb.RegisterTicketsServiceServer(registrar, server{app: app})
	return nil
}

// Match Management

func (s server) CreateMatch(ctx context.Context, req *ticketspb.CreateMatchRequest) (*ticketspb.CreateMatchResponse, error) {
	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	matchID := uuid.New().String()

	// Convert sector pricing
	sectorPricing := make(map[string]domain.SectorPricing)
	for _, sp := range req.GetSectorPricing() {
		sectorPricing[sp.GetSectorId()] = domain.SectorPricing{
			BasePrice:      sp.GetBasePrice(),
			DynamicPricing: sp.GetDynamicPricing(),
		}
	}

	cmd := commands.CreateMatch{
		ID:              matchID,
		HomeTeamID:      req.GetHomeTeamId(),
		AwayTeamID:      req.GetAwayTeamId(),
		StadiumID:       req.GetStadiumId(),
		MatchDate:       req.GetMatchDate().AsTime(),
		CompetitionType: domain.CompetitionType(req.GetCompetitionType().String()),
		CompetitionName: req.GetCompetitionName(),
		Round:           req.GetRound(),
		Season:          req.GetSeason(),
		SectorPricing:   sectorPricing,
	}

	if err := s.app.CreateMatch(ctx, cmd); err != nil {
		return nil, err
	}

	return &ticketspb.CreateMatchResponse{MatchId: matchID}, nil
}

func (s server) GetMatch(ctx context.Context, req *ticketspb.GetMatchRequest) (*ticketspb.GetMatchResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("match_id", req.GetMatchId()))

	match, err := s.app.GetMatch(ctx, queries.GetMatch{
		MatchID: req.GetMatchId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &ticketspb.GetMatchResponse{
		Match: s.matchFromDomain(match),
	}, nil
}

func (s server) UpdateMatch(ctx context.Context, req *ticketspb.UpdateMatchRequest) (*ticketspb.UpdateMatchResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("match_id", req.GetMatchId()))

	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	cmd := commands.UpdateMatch{
		ID:        req.GetMatchId(),
		MatchDate: req.GetMatchDate().AsTime(),
		Referee:   req.GetReferee(),
		Officials: req.GetOfficials(),
	}

	if err := s.app.UpdateMatch(ctx, cmd); err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &ticketspb.UpdateMatchResponse{Success: true}, nil
}

func (s server) OpenMatchSales(ctx context.Context, req *ticketspb.OpenMatchSalesRequest) (*ticketspb.OpenMatchSalesResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("match_id", req.GetMatchId()))

	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	var salesStartAt, salesEndAt time.Time
	if req.GetSalesStartAt() != nil {
		salesStartAt = req.GetSalesStartAt().AsTime()
	} else {
		salesStartAt = time.Now()
	}
	if req.GetSalesEndAt() != nil {
		salesEndAt = req.GetSalesEndAt().AsTime()
	}

	totalTickets, err := s.app.OpenMatchSales(ctx, commands.OpenMatchSales{
		MatchID:      req.GetMatchId(),
		SalesStartAt: salesStartAt,
		SalesEndAt:   salesEndAt,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &ticketspb.OpenMatchSalesResponse{
		Success:      true,
		TotalTickets: totalTickets,
	}, nil
}

func (s server) CloseMatchSales(ctx context.Context, req *ticketspb.CloseMatchSalesRequest) (*ticketspb.CloseMatchSalesResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("match_id", req.GetMatchId()))

	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	result, err := s.app.CloseMatchSales(ctx, commands.CloseMatchSales{
		MatchID: req.GetMatchId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &ticketspb.CloseMatchSalesResponse{
		Success:          true,
		TicketsSold:      result.TicketsSold,
		TicketsRemaining: result.TicketsRemaining,
	}, nil
}

func (s server) CancelMatch(ctx context.Context, req *ticketspb.CancelMatchRequest) (*ticketspb.CancelMatchResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("match_id", req.GetMatchId()))

	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	ticketsAffected, err := s.app.CancelMatch(ctx, commands.CancelMatch{
		MatchID:      req.GetMatchId(),
		Reason:       req.GetReason(),
		IssueRefunds: req.GetIssueRefunds(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &ticketspb.CancelMatchResponse{
		Success:         true,
		TicketsAffected: ticketsAffected,
	}, nil
}

func (s server) PostponeMatch(ctx context.Context, req *ticketspb.PostponeMatchRequest) (*ticketspb.PostponeMatchResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("match_id", req.GetMatchId()))

	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	cmd := commands.PostponeMatch{
		MatchID: req.GetMatchId(),
		NewDate: req.GetNewDate().AsTime(),
		Reason:  req.GetReason(),
	}

	if err := s.app.PostponeMatch(ctx, cmd); err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &ticketspb.PostponeMatchResponse{Success: true}, nil
}

// Sector Management

func (s server) UpdateSectorPricing(ctx context.Context, req *ticketspb.UpdateSectorPricingRequest) (*ticketspb.UpdateSectorPricingResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("match_id", req.GetMatchId()),
		attribute.String("sector_id", req.GetSectorId()),
	)

	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	oldPrice, newPrice, err := s.app.UpdateSectorPricing(ctx, commands.UpdateSectorPricing{
		MatchID:              req.GetMatchId(),
		SectorID:             req.GetSectorId(),
		NewPrice:             req.GetNewPrice(),
		EnableDynamicPricing: req.GetEnableDynamicPricing(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &ticketspb.UpdateSectorPricingResponse{
		Success:  true,
		OldPrice: oldPrice,
		NewPrice: newPrice,
	}, nil
}

func (s server) GetSectorAvailability(ctx context.Context, req *ticketspb.GetSectorAvailabilityRequest) (*ticketspb.GetSectorAvailabilityResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("match_id", req.GetMatchId()),
		attribute.String("sector_id", req.GetSectorId()),
	)

	availability, err := s.app.GetSectorAvailability(ctx, queries.GetSectorAvailability{
		MatchID:  req.GetMatchId(),
		SectorID: req.GetSectorId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Convert available seats
	availableSeats := make([]*ticketspb.AvailableSeat, len(availability.AvailableSeats))
	for i, seat := range availability.AvailableSeats {
		availableSeats[i] = &ticketspb.AvailableSeat{
			RowNumber:              seat.RowNumber,
			SeatNumber:             seat.SeatNumber,
			Price:                  seat.Price,
			HasRestrictedView:      seat.HasRestrictedView,
			IsWheelchairAccessible: seat.IsWheelchairAccessible,
		}
	}

	return &ticketspb.GetSectorAvailabilityResponse{
		SectorId:            availability.SectorID,
		TotalSeats:          availability.TotalSeats,
		AvailableSeats:      availability.AvailableSeats,
		AvailableSeatDetails: availableSeats,
	}, nil
}

// Ticket Operations

func (s server) PurchaseTickets(ctx context.Context, req *ticketspb.PurchaseTicketsRequest) (*ticketspb.PurchaseTicketsResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("match_id", req.GetMatchId()),
		attribute.String("user_id", req.GetUserId()),
		attribute.Int("seat_count", len(req.GetSeats())),
	)

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Ensure the user can only purchase for themselves
	if req.GetUserId() != claims.Subject {
		return nil, status.Error(grpc_code.PermissionDenied, "cannot purchase tickets for another user")
	}

	// Convert seat selections
	seats := make([]domain.SeatSelection, len(req.GetSeats()))
	for i, seat := range req.GetSeats() {
		seats[i] = domain.SeatSelection{
			SectorID:   seat.GetSectorId(),
			RowNumber:  seat.GetRowNumber(),
			SeatNumber: seat.GetSeatNumber(),
		}
	}

	cmd := commands.PurchaseTickets{
		MatchID:       req.GetMatchId(),
		UserID:        req.GetUserId(),
		Seats:         seats,
		PaymentMethod: s.paymentMethodFromProto(req.GetPaymentMethod()),
		PromoCode:     req.GetPromoCode(),
	}

	result, err := s.app.PurchaseTickets(ctx, cmd)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &ticketspb.PurchaseTicketsResponse{
		Success:     true,
		TicketIds:   result.TicketIDs,
		TotalAmount: result.TotalAmount,
		PaymentId:   result.PaymentID,
	}, nil
}

func (s server) ReserveTickets(ctx context.Context, req *ticketspb.ReserveTicketsRequest) (*ticketspb.ReserveTicketsResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("match_id", req.GetMatchId()),
		attribute.String("user_id", req.GetUserId()),
		attribute.Int("seat_count", len(req.GetSeats())),
	)

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Ensure the user can only reserve for themselves
	if req.GetUserId() != claims.Subject {
		return nil, status.Error(grpc_code.PermissionDenied, "cannot reserve tickets for another user")
	}

	// Convert seat selections
	seats := make([]domain.SeatSelection, len(req.GetSeats()))
	for i, seat := range req.GetSeats() {
		seats[i] = domain.SeatSelection{
			SectorID:   seat.GetSectorId(),
			RowNumber:  seat.GetRowNumber(),
			SeatNumber: seat.GetSeatNumber(),
		}
	}

	duration := req.GetReservationDurationMinutes()
	if duration <= 0 {
		duration = 15 // Default 15 minutes
	}

	result, err := s.app.ReserveTickets(ctx, commands.ReserveTickets{
		MatchID:                   req.GetMatchId(),
		UserID:                    req.GetUserId(),
		Seats:                     seats,
		ReservationDurationMinutes: duration,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &ticketspb.ReserveTicketsResponse{
		Success:   true,
		TicketIds: result.TicketIDs,
		ExpiresAt: timestamppb.New(result.ExpiresAt),
	}, nil
}

func (s server) GetTicket(ctx context.Context, req *ticketspb.GetTicketRequest) (*ticketspb.GetTicketResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ticket_id", req.GetTicketId()))

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	ticket, match, err := s.app.GetTicket(ctx, queries.GetTicket{
		TicketID:      req.GetTicketId(),
		UserID:        claims.Subject,
		IncludeQRCode: req.GetIncludeQrCode(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &ticketspb.GetTicketResponse{
		Ticket:    s.ticketFromDomain(ticket),
		MatchInfo: s.matchFromDomain(match),
	}, nil
}

func (s server) TransferTicket(ctx context.Context, req *ticketspb.TransferTicketRequest) (*ticketspb.TransferTicketResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ticket_id", req.GetTicketId()),
		attribute.String("from_user_id", req.GetFromUserId()),
		attribute.String("to_user_id", req.GetToUserId()),
	)

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Ensure the user owns the ticket
	if req.GetFromUserId() != claims.Subject {
		return nil, status.Error(grpc_code.PermissionDenied, "can only transfer your own tickets")
	}

	transferredAt, err := s.app.TransferTicket(ctx, commands.TransferTicket{
		TicketID:   req.GetTicketId(),
		FromUserID: req.GetFromUserId(),
		ToUserID:   req.GetToUserId(),
		ToEmail:    req.GetToEmail(),
		Reason:     req.GetReason(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &ticketspb.TransferTicketResponse{
		Success:       true,
		TransferredAt: timestamppb.New(transferredAt),
	}, nil
}

func (s server) ValidateTicket(ctx context.Context, req *ticketspb.ValidateTicketRequest) (*ticketspb.ValidateTicketResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("qr_code", req.GetQrCode()),
		attribute.String("gate_number", req.GetGateNumber()),
	)

	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	result, err := s.app.ValidateTicket(ctx, commands.ValidateTicket{
		QRCode:      req.GetQrCode(),
		GateNumber:  req.GetGateNumber(),
		ValidatorID: req.GetValidatorId(),
	})
	if err != nil {
		// Don't log full error for security reasons
		span.SetStatus(codes.Error, "validation failed")
		return &ticketspb.ValidateTicketResponse{
			Valid:   false,
			Message: "Invalid ticket",
		}, nil
	}

	return &ticketspb.ValidateTicketResponse{
		Valid:      true,
		TicketId:   result.TicketID,
		Sector:     result.Sector,
		Row:        result.Row,
		Seat:       result.Seat,
		HolderName: result.HolderName,
		Message:    "Ticket validated successfully",
	}, nil
}

func (s server) CancelTicket(ctx context.Context, req *ticketspb.CancelTicketRequest) (*ticketspb.CancelTicketResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ticket_id", req.GetTicketId()))

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	err := s.app.CancelTicket(ctx, commands.CancelTicket{
		TicketID: req.GetTicketId(),
		UserID:   claims.Subject,
		Reason:   req.GetReason(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &ticketspb.CancelTicketResponse{Success: true}, nil
}

func (s server) RefundTicket(ctx context.Context, req *ticketspb.RefundTicketRequest) (*ticketspb.RefundTicketResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ticket_id", req.GetTicketId()))

	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	refundAmount, refundID, err := s.app.RefundTicket(ctx, commands.RefundTicket{
		TicketID:     req.GetTicketId(),
		RefundAmount: req.GetRefundAmount(),
		Reason:       req.GetReason(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &ticketspb.RefundTicketResponse{
		Success:      true,
		RefundAmount: refundAmount,
		RefundId:     refundID,
	}, nil
}

// User Tickets

func (s server) GetUserTickets(ctx context.Context, req *ticketspb.GetUserTicketsRequest) (*ticketspb.GetUserTicketsResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("user_id", req.GetUserId()))

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Users can only get their own tickets
	if req.GetUserId() != claims.Subject {
		return nil, status.Error(grpc_code.PermissionDenied, "cannot access other user's tickets")
	}

	page := req.GetPage()
	if page <= 0 {
		page = 1
	}

	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		pageSize = 20
	}

	tickets, totalCount, err := s.app.GetUserTickets(ctx, queries.GetUserTickets{
		UserID:      req.GetUserId(),
		IncludePast: req.GetIncludePast(),
		Page:        page,
		PageSize:    pageSize,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Convert to summaries
	summaries := make([]*ticketspb.TicketSummary, len(tickets))
	for i, ticket := range tickets {
		summaries[i] = s.ticketSummaryFromDomain(ticket)
	}

	return &ticketspb.GetUserTicketsResponse{
		Tickets:    summaries,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s server) GetUserUpcomingMatches(ctx context.Context, req *ticketspb.GetUserUpcomingMatchesRequest) (*ticketspb.GetUserUpcomingMatchesResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("user_id", req.GetUserId()))

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Users can only get their own upcoming matches
	if req.GetUserId() != claims.Subject {
		return nil, status.Error(grpc_code.PermissionDenied, "cannot access other user's matches")
	}

	matches, err := s.app.GetUserUpcomingMatches(ctx, queries.GetUserUpcomingMatches{
		UserID: req.GetUserId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	pbMatches := make([]*ticketspb.Match, len(matches))
	for i, match := range matches {
		pbMatches[i] = s.matchFromDomain(match)
	}

	return &ticketspb.GetUserUpcomingMatchesResponse{Matches: pbMatches}, nil
}

// Match Discovery

func (s server) GetUpcomingMatches(ctx context.Context, req *ticketspb.GetUpcomingMatchesRequest) (*ticketspb.GetUpcomingMatchesResponse, error) {
	span := trace.SpanFromContext(ctx)

	page := req.GetPage()
	if page <= 0 {
		page = 1
	}

	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		pageSize = 20
	}

	daysAhead := req.GetDaysAhead()
	if daysAhead <= 0 {
		daysAhead = 30 // Default to 30 days
	}

	matches, totalCount, err := s.app.GetUpcomingMatches(ctx, queries.GetUpcomingMatches{
		TeamID:          req.GetTeamId(),
		StadiumID:       req.GetStadiumId(),
		CompetitionType: domain.CompetitionType(req.GetCompetitionType().String()),
		DaysAhead:       daysAhead,
		Page:            page,
		PageSize:        pageSize,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	pbMatches := make([]*ticketspb.Match, len(matches))
	for i, match := range matches {
		pbMatches[i] = s.matchFromDomain(match)
	}

	return &ticketspb.GetUpcomingMatchesResponse{
		Matches:    pbMatches,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s server) SearchMatches(ctx context.Context, req *ticketspb.SearchMatchesRequest) (*ticketspb.SearchMatchesResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("query", req.GetQuery()))

	var filters *domain.MatchFilters
	if req.GetFilters() != nil {
		f := req.GetFilters()
		filters = &domain.MatchFilters{
			TeamIDs:         f.GetTeamIds(),
			StadiumIDs:      f.GetStadiumIds(),
			CompetitionType: domain.CompetitionType(f.GetCompetitionType().String()),
			DateFrom:        f.GetDateFrom().AsTime(),
			DateTo:          f.GetDateTo().AsTime(),
			OnlyAvailable:   f.GetOnlyAvailable(),
			MaxPrice:        f.GetMaxPrice(),
			SortBy:          f.GetSortBy(),
			SortOrder:       f.GetSortOrder(),
		}
	}

	matches, totalCount, err := s.app.SearchMatches(ctx, queries.SearchMatches{
		Query:   req.GetQuery(),
		Filters: filters,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	pbMatches := make([]*ticketspb.Match, len(matches))
	for i, match := range matches {
		pbMatches[i] = s.matchFromDomain(match)
	}

	return &ticketspb.SearchMatchesResponse{
		Matches:    pbMatches,
		TotalCount: totalCount,
	}, nil
}

// Season Tickets

func (s server) CreateSeasonTicket(ctx context.Context, req *ticketspb.CreateSeasonTicketRequest) (*ticketspb.CreateSeasonTicketResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Ensure the user can only create season tickets for themselves
	if req.GetUserId() != claims.Subject {
		return nil, status.Error(grpc_code.PermissionDenied, "cannot create season ticket for another user")
	}

	seasonTicketID := uuid.New().String()

	cmd := commands.CreateSeasonTicket{
		ID:                   seasonTicketID,
		UserID:               req.GetUserId(),
		TeamID:               req.GetTeamId(),
		SectorID:             req.GetSectorId(),
		RowNumber:            req.GetRowNumber(),
		SeatNumber:           req.GetSeatNumber(),
		Season:               req.GetSeason(),
		IncludedCompetitions: req.GetIncludedCompetitions(),
		Price:                req.GetPrice(),
	}

	if err := s.app.CreateSeasonTicket(ctx, cmd); err != nil {
		return nil, err
	}

	return &ticketspb.CreateSeasonTicketResponse{SeasonTicketId: seasonTicketID}, nil
}

func (s server) RenewSeasonTicket(ctx context.Context, req *ticketspb.RenewSeasonTicketRequest) (*ticketspb.RenewSeasonTicketResponse, error) {
	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	result, err := s.app.RenewSeasonTicket(ctx, commands.RenewSeasonTicket{
		SeasonTicketID: req.GetSeasonTicketId(),
		NewSeason:      req.GetNewSeason(),
		PaymentMethod:  s.paymentMethodFromProto(req.GetPaymentMethod()),
	})
	if err != nil {
		return nil, err
	}

	return &ticketspb.RenewSeasonTicketResponse{
		Success:      true,
		RenewalPrice: result.RenewalPrice,
		ValidUntil:   timestamppb.New(result.ValidUntil),
	}, nil
}

// Stadium Management

func (s server) CreateStadium(ctx context.Context, req *ticketspb.CreateStadiumRequest) (*ticketspb.CreateStadiumResponse, error) {
	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	stadiumID := uuid.New().String()

	// Convert sections
	sections := make([]domain.StadiumSection, len(req.GetSections()))
	for i, section := range req.GetSections() {
		// Convert rows
		rows := make([]domain.Row, len(section.GetRows()))
		for j, row := range section.GetRows() {
			rows[j] = domain.Row{
				Number:    row.GetNumber(),
				SeatCount: row.GetSeatCount(),
			}
		}

		sections[i] = domain.StadiumSection{
			Name:     section.GetName(),
			Type:     domain.SectorType(section.GetType().String()),
			Capacity: section.GetCapacity(),
			Rows:     rows,
		}
	}

	cmd := commands.CreateStadium{
		ID:       stadiumID,
		Name:     req.GetName(),
		City:     req.GetCity(),
		Country:  req.GetCountry(),
		Capacity: req.GetCapacity(),
		Location: domain.Location{
			Latitude:   req.GetLocation().GetLatitude(),
			Longitude:  req.GetLocation().GetLongitude(),
			Address:    req.GetLocation().GetAddress(),
			PostalCode: req.GetLocation().GetPostalCode(),
		},
		Sections:   sections,
		ImageURL:   req.GetImageUrl(),
		Facilities: req.GetFacilities(),
	}

	if err := s.app.CreateStadium(ctx, cmd); err != nil {
		return nil, err
	}

	return &ticketspb.CreateStadiumResponse{StadiumId: stadiumID}, nil
}

func (s server) GetStadium(ctx context.Context, req *ticketspb.GetStadiumRequest) (*ticketspb.GetStadiumResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("stadium_id", req.GetStadiumId()))

	stadium, err := s.app.GetStadium(ctx, queries.GetStadium{
		StadiumID: req.GetStadiumId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &ticketspb.GetStadiumResponse{
		Stadium: s.stadiumFromDomain(stadium),
	}, nil
}

// Helper methods

func (s server) matchFromDomain(match *domain.Match) *ticketspb.Match {
	// Convert sectors
	sectors := make([]*ticketspb.SectorInfo, len(match.Sectors))
	i := 0
	for sectorID, sector := range match.Sectors {
		sectors[i] = &ticketspb.SectorInfo{
			SectorId:       sectorID,
			Name:           sector.Name,
			Type:           s.sectorTypeToProto(sector.Type),
			BasePrice:      sector.BasePrice,
			CurrentPrice:   sector.CurrentPrice,
			TotalSeats:     sector.TotalSeats,
			AvailableSeats: sector.AvailableSeats,
			DynamicPricing: sector.DynamicPricing,
		}
		i++
	}

	return &ticketspb.Match{
		Id:               match.ID,
		HomeTeam:         s.teamFromDomain(match.HomeTeam),
		AwayTeam:         s.teamFromDomain(match.AwayTeam),
		Stadium:          s.stadiumFromDomain(match.Stadium),
		MatchDate:        timestamppb.New(match.MatchDate),
		Status:           s.matchStatusToProto(match.Status),
		CompetitionType:  s.competitionTypeToProto(match.CompetitionType),
		CompetitionName:  match.CompetitionName,
		Round:            match.Round,
		Season:           match.Season,
		TotalCapacity:    match.TotalCapacity,
		TicketsSold:      match.TicketsSold,
		TicketsAvailable: match.TicketsAvailable,
		Sectors:          sectors,
		Referee:          match.Referee,
		Officials:        match.Officials,
		IsHighDemand:     match.IsHighDemand,
		CreatedAt:        timestamppb.New(match.CreatedAt),
		UpdatedAt:        timestamppb.New(match.UpdatedAt),
		SalesStartAt:     timestamppb.New(match.SalesStartAt),
		SalesEndAt:       timestamppb.New(match.SalesEndAt),
	}
}

func (s server) teamFromDomain(team domain.Team) *ticketspb.Team {
	return &ticketspb.Team{
		Id:        team.ID,
		Name:      team.Name,
		Code:      team.Code,
		LogoUrl:   team.LogoURL,
		StadiumId: team.StadiumID,
		City:      team.City,
		Country:   team.Country,
	}
}

func (s server) stadiumFromDomain(stadium *domain.Stadium) *ticketspb.Stadium {
	// Convert sections
	sections := make([]*ticketspb.StadiumSection, len(stadium.Sections))
	for i, section := range stadium.Sections {
		// Convert rows
		rows := make([]*ticketspb.Row, len(section.Rows))
		for j, row := range section.Rows {
			// Convert seats
			seats := make([]*ticketspb.Seat, len(row.Seats))
			for k, seat := range row.Seats {
				seats[k] = &ticketspb.Seat{
					Number:                 seat.Number,
					IsAvailable:            seat.IsAvailable,
					HasRestrictedView:      seat.HasRestrictedView,
					IsWheelchairAccessible: seat.IsWheelchairAccessible,
				}
			}

			rows[j] = &ticketspb.Row{
				Id:        row.ID,
				Number:    row.Number,
				SeatCount: row.SeatCount,
				Seats:     seats,
			}
		}

		sections[i] = &ticketspb.StadiumSection{
			Id:       section.ID,
			Name:     section.Name,
			Type:     s.sectorTypeToProto(section.Type),
			Capacity: section.Capacity,
			Rows:     rows,
		}
	}

	var location *ticketspb.Location
	if stadium.Location != nil {
		location = &ticketspb.Location{
			Latitude:   stadium.Location.Latitude,
			Longitude:  stadium.Location.Longitude,
			Address:    stadium.Location.Address,
			PostalCode: stadium.Location.PostalCode,
		}
	}

	return &ticketspb.Stadium{
		Id:         stadium.ID,
		Name:       stadium.Name,
		City:       stadium.City,
		Country:    stadium.Country,
		Capacity:   stadium.Capacity,
		Location:   location,
		Sections:   sections,
		ImageUrl:   stadium.ImageURL,
		Facilities: stadium.Facilities,
	}
}

func (s server) ticketFromDomain(ticket *domain.Ticket) *ticketspb.Ticket {
	// Convert transfers
	transfers := make([]*ticketspb.Transfer, len(ticket.Transfers))
	for i, transfer := range ticket.Transfers {
		transfers[i] = &ticketspb.Transfer{
			FromUserId:    transfer.FromUserID,
			ToUserId:      transfer.ToUserID,
			TransferredAt: timestamppb.New(transfer.TransferredAt),
			Reason:        transfer.Reason,
		}
	}

	var validatedAt, expiresAt *timestamppb.Timestamp
	if !ticket.ValidatedAt.IsZero() {
		validatedAt = timestamppb.New(ticket.ValidatedAt)
	}
	if !ticket.ExpiresAt.IsZero() {
		expiresAt = timestamppb.New(ticket.ExpiresAt)
	}

	return &ticketspb.Ticket{
		Id:             ticket.ID,
		MatchId:        ticket.MatchID,
		SectorId:       ticket.SectorID,
		RowNumber:      ticket.RowNumber,
		SeatNumber:     ticket.SeatNumber,
		Status:         s.ticketStatusToProto(ticket.Status),
		OwnerId:        ticket.OwnerID,
		OwnerName:      ticket.OwnerName,
		OwnerEmail:     ticket.OwnerEmail,
		OriginalPrice:  ticket.OriginalPrice,
		PaidPrice:      ticket.PaidPrice,
		ResalePrice:    ticket.ResalePrice,
		QrCode:         ticket.QRCode,
		Barcode:        ticket.Barcode,
		SecurityCode:   ticket.SecurityCode,
		Transfers:      transfers,
		IsTransferable: ticket.IsTransferable,
		IsResellable:   ticket.IsResellable,
		ValidatedAt:    validatedAt,
		ValidatedBy:    ticket.ValidatedBy,
		GateNumber:     ticket.GateNumber,
		CreatedAt:      timestamppb.New(ticket.CreatedAt),
		PurchasedAt:    timestamppb.New(ticket.PurchasedAt),
		ExpiresAt:      expiresAt,
	}
}

func (s server) ticketSummaryFromDomain(ticket *domain.TicketSummary) *ticketspb.TicketSummary {
	return &ticketspb.TicketSummary{
		Id:        ticket.ID,
		MatchId:   ticket.MatchID,
		Section:   ticket.SectionName,
		Row:       ticket.Row,
		Seat:      ticket.Seat,
		Price:     ticket.Price,
		Status:    s.ticketStatusToProto(ticket.Status),
		MatchDate: timestamppb.New(ticket.MatchDate),
		HomeTeam:  ticket.HomeTeam,
		AwayTeam:  ticket.AwayTeam,
	}
}

func (s server) paymentMethodFromProto(pm *ticketspb.PaymentMethod) domain.PaymentMethod {
	if pm == nil {
		return domain.PaymentMethod{}
	}
	return domain.PaymentMethod{
		Type:  pm.GetType(),
		Token: pm.GetToken(),
	}
}

// Proto conversion helpers

func (s server) matchStatusToProto(status domain.MatchStatus) ticketspb.MatchStatus {
	switch status {
	case domain.MatchStatusScheduled:
		return ticketspb.MatchStatus_MATCH_STATUS_SCHEDULED
	case domain.MatchStatusOnSale:
		return ticketspb.MatchStatus_MATCH_STATUS_ON_SALE
	case domain.MatchStatusSoldOut:
		return ticketspb.MatchStatus_MATCH_STATUS_SOLD_OUT
	case domain.MatchStatusInProgress:
		return ticketspb.MatchStatus_MATCH_STATUS_IN_PROGRESS
	case domain.MatchStatusCompleted:
		return ticketspb.MatchStatus_MATCH_STATUS_COMPLETED
	case domain.MatchStatusCancelled:
		return ticketspb.MatchStatus_MATCH_STATUS_CANCELLED
	case domain.MatchStatusPostponed:
		return ticketspb.MatchStatus_MATCH_STATUS_POSTPONED
	default:
		return ticketspb.MatchStatus_MATCH_STATUS_UNSPECIFIED
	}
}

func (s server) ticketStatusToProto(status domain.TicketStatus) ticketspb.TicketStatus {
	switch status {
	case domain.TicketStatusAvailable:
		return ticketspb.TicketStatus_TICKET_STATUS_AVAILABLE
	case domain.TicketStatusReserved:
		return ticketspb.TicketStatus_TICKET_STATUS_RESERVED
	case domain.TicketStatusSold:
		return ticketspb.TicketStatus_TICKET_STATUS_SOLD
	case domain.TicketStatusTransferred:
		return ticketspb.TicketStatus_TICKET_STATUS_TRANSFERRED
	case domain.TicketStatusUsed:
		return ticketspb.TicketStatus_TICKET_STATUS_USED
	case domain.TicketStatusCancelled:
		return ticketspb.TicketStatus_TICKET_STATUS_CANCELLED
	case domain.TicketStatusRefunded:
		return ticketspb.TicketStatus_TICKET_STATUS_REFUNDED
	default:
		return ticketspb.TicketStatus_TICKET_STATUS_UNSPECIFIED
	}
}

func (s server) sectorTypeToProto(st domain.SectorType) ticketspb.SectorType {
	switch st {
	case domain.SectorTypeVIP:
		return ticketspb.SectorType_SECTOR_TYPE_VIP
	case domain.SectorTypePremium:
		return ticketspb.SectorType_SECTOR_TYPE_PREMIUM
	case domain.SectorTypeStandard:
		return ticketspb.SectorType_SECTOR_TYPE_STANDARD
	case domain.SectorTypeFamily:
		return ticketspb.SectorType_SECTOR_TYPE_FAMILY
	case domain.SectorTypeAway:
		return ticketspb.SectorType_SECTOR_TYPE_AWAY
	case domain.SectorTypeHome:
		return ticketspb.SectorType_SECTOR_TYPE_HOME
	case domain.SectorTypeNeutral:
		return ticketspb.SectorType_SECTOR_TYPE_NEUTRAL
	default:
		return ticketspb.SectorType_SECTOR_TYPE_UNSPECIFIED
	}
}

func (s server) competitionTypeToProto(ct domain.CompetitionType) ticketspb.CompetitionType {
	switch ct {
	case domain.CompetitionTypeLeague:
		return ticketspb.CompetitionType_COMPETITION_TYPE_LEAGUE
	case domain.CompetitionTypeCup:
		return ticketspb.CompetitionType_COMPETITION_TYPE_CUP
	case domain.CompetitionTypeFriendly:
		return ticketspb.CompetitionType_COMPETITION_TYPE_FRIENDLY
	case domain.CompetitionTypeChampionsLeague:
		return ticketspb.CompetitionType_COMPETITION_TYPE_CHAMPIONS_LEAGUE
	case domain.CompetitionTypeEuropaLeague:
		return ticketspb.CompetitionType_COMPETITION_TYPE_EUROPA_LEAGUE
	case domain.CompetitionTypeInternational:
		return ticketspb.CompetitionType_COMPETITION_TYPE_INTERNATIONAL
	default:
		return ticketspb.CompetitionType_COMPETITION_TYPE_UNSPECIFIED
	}
}