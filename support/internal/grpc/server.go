package grpc

import (
	"context"
	"time"
	"middleman/internal/errorsotel"
	"middleman/support/internal/application"
	"middleman/support/internal/application/commands"
	"middleman/support/internal/application/queries"
	"middleman/support/internal/domain"
	"middleman/support/supportpb"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	app application.App
	supportpb.UnimplementedSupportServiceServer
}

var _ supportpb.SupportServiceServer = (*server)(nil)

func RegisterServer(_ context.Context, app application.App, registrar grpc.ServiceRegistrar) error {
	supportpb.RegisterSupportServiceServer(registrar, server{app: app})
	return nil
}

// Support Channel Methods
func (s server) CreateSupportChannel(ctx context.Context, request *supportpb.CreateSupportChannelRequest) (*supportpb.CreateSupportChannelResponse, error) {
	span := trace.SpanFromContext(ctx)
	channelID := uuid.New().String()
	
	span.SetAttributes(
		attribute.String("channelID", channelID),
		attribute.String("userID", request.GetUserId()),
	)
	
	settings := domain.SupportChannelSettings{
		EmailNotifications: request.GetSettings().GetEmailNotifications(),
		SMSNotifications:   request.GetSettings().GetSmsNotifications(),
		AutoAssignTickets:  request.GetSettings().GetAutoAssignTickets(),
		PreferredLanguage:  request.GetSettings().GetPreferredLanguage(),
		Timezone:           request.GetSettings().GetTimezone(),
		NotificationEmails: request.GetSettings().GetNotificationEmails(),
		SLASettings: domain.SLASettings{
			FirstResponseMinutes: int(request.GetSettings().GetSlaSettings().GetFirstResponseMinutes()),
			ResolutionHours:      int(request.GetSettings().GetSlaSettings().GetResolutionHours()),
			PriorityResponseTimes: convertInt32MapToIntMap(request.GetSettings().GetSlaSettings().GetPriorityResponseTimes()),
		},
	}
	
	err := s.app.CreateSupportChannel(ctx, commands.CreateSupportChannel{
		ID:          channelID,
		UserID:      request.GetUserId(),
		BusinessID:  request.GetBusinessId(),
		ChannelType: mapChannelType(request.GetChannelType()),
		Settings:    settings,
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	return &supportpb.CreateSupportChannelResponse{
		Id: channelID,
	}, nil
}

func (s server) GetSupportChannel(ctx context.Context, request *supportpb.GetSupportChannelRequest) (*supportpb.GetSupportChannelResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("channelID", request.GetId()),
	)
	
	channel, err := s.app.GetSupportChannel(ctx, queries.GetSupportChannel{
		ID: request.GetId(),
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	return &supportpb.GetSupportChannelResponse{
		Channel: toProtoChannel(channel),
	}, nil
}

func (s server) GetUserSupportChannels(ctx context.Context, request *supportpb.GetUserSupportChannelsRequest) (*supportpb.GetUserSupportChannelsResponse, error) {
	channels, err := s.app.GetUserSupportChannels(ctx, queries.GetUserSupportChannels{
		UserID:     request.GetUserId(),
		ActiveOnly: request.GetActiveOnly(),
		Page:       int(request.GetPage()),
		Limit:      int(request.GetLimit()),
	})
	
	if err != nil {
		return nil, err
	}
	
	protoChannels := make([]*supportpb.SupportChannel, len(channels))
	for i, channel := range channels {
		protoChannels[i] = toProtoChannel(channel)
	}
	
	return &supportpb.GetUserSupportChannelsResponse{
		Channels:   protoChannels,
		TotalCount: int32(len(channels)),
	}, nil
}

// Ticket Methods
func (s server) CreateTicket(ctx context.Context, request *supportpb.CreateTicketRequest) (*supportpb.CreateTicketResponse, error) {
	span := trace.SpanFromContext(ctx)
	ticketID := uuid.New().String()
	
	span.SetAttributes(
		attribute.String("ticketID", ticketID),
		attribute.String("channelID", request.GetChannelId()),
	)
	
	// Get authenticated user ID from context
	userID := getUserIDFromContext(ctx)
	
	err := s.app.CreateTicket(ctx, commands.CreateTicket{
		ID:          ticketID,
		ChannelID:   request.GetChannelId(),
		Title:       request.GetTitle(),
		Description: request.GetDescription(),
		Category:    mapTicketCategory(request.GetCategory()),
		Priority:    mapTicketPriority(request.GetPriority()),
		Tags:        request.GetTags(),
		Metadata:    request.GetMetadata(),
		CreatedBy:   userID,
		Attachments: mapAttachments(request.GetAttachments()),
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	return &supportpb.CreateTicketResponse{
		Id: ticketID,
	}, nil
}

func (s server) GetTicket(ctx context.Context, request *supportpb.GetTicketRequest) (*supportpb.GetTicketResponse, error) {
	ticket, err := s.app.GetTicket(ctx, queries.GetTicket{
		ID: request.GetId(),
	})
	
	if err != nil {
		return nil, err
	}
	
	return &supportpb.GetTicketResponse{
		Ticket: toProtoTicket(ticket),
	}, nil
}

func (s server) GetTickets(ctx context.Context, request *supportpb.GetTicketsRequest) (*supportpb.GetTicketsResponse, error) {
	tickets, err := s.app.GetTickets(ctx, queries.GetTickets{
		IDs: request.GetIds(),
	})
	
	if err != nil {
		return nil, err
	}
	
	protoTickets := make([]*supportpb.Ticket, len(tickets))
	for i, ticket := range tickets {
		protoTickets[i] = toProtoTicket(ticket)
	}
	
	return &supportpb.GetTicketsResponse{
		Tickets: protoTickets,
	}, nil
}

func (s server) GetChannelTickets(ctx context.Context, request *supportpb.GetChannelTicketsRequest) (*supportpb.GetChannelTicketsResponse, error) {
	var statusFilter *string
	if request.StatusFilter != supportpb.TicketStatus_DRAFT {
		status := request.GetStatusFilter().String()
		statusFilter = &status
	}
	
	tickets, err := s.app.GetChannelTickets(ctx, queries.GetChannelTickets{
		ChannelID:    request.GetChannelId(),
		StatusFilter: statusFilter,
		Page:         int(request.GetPage()),
		Limit:        int(request.GetLimit()),
		SortBy:       request.GetSortBy(),
		Descending:   request.GetDescending(),
	})
	
	if err != nil {
		return nil, err
	}
	
	protoTickets := make([]*supportpb.Ticket, len(tickets))
	for i, ticket := range tickets {
		protoTickets[i] = toProtoTicket(ticket)
	}
	
	return &supportpb.GetChannelTicketsResponse{
		Tickets:    protoTickets,
		TotalCount: int32(len(tickets)),
	}, nil
}

func (s server) ListTickets(ctx context.Context, request *supportpb.ListTicketsRequest) (*supportpb.ListTicketsResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("supportID", request.GetSupportId()),
		attribute.Int("page", int(request.GetPage())),
		attribute.Int("limit", int(request.GetLimit())),
	)
	
	limit := int(request.GetLimit())
	if limit == 0 {
		limit = 20
	}
	
	page := int(request.GetPage())
	if page == 0 {
		page = 1
	}
	
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	
	// Convert proto filters to domain filters
	filters := make(map[string]interface{})
	for k, v := range request.GetFilters() {
		filters[k] = v
	}
	
	// Add support_id/channel_id to filters if provided
	if request.GetSupportId() != "" {
		filters["channel_id"] = request.GetSupportId()
	}
	
	// Get tickets from application layer
	tickets, err := s.app.SearchTickets(ctx, queries.SearchTickets{
		Query:      request.GetSearchQuery(),
		Filters:    filters,
		Limit:      limit,
		Offset:     offset,
		SortBy:     request.GetSortBy(),
		Descending: request.GetDescending(),
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	// Count total tickets
	totalCount, err := s.app.CountTickets(ctx, queries.CountTickets{
		Query:   request.GetSearchQuery(),
		Filters: filters,
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	// Convert to proto tickets
	protoTickets := make([]*supportpb.Ticket, len(tickets))
	for i, t := range tickets {
		protoTickets[i] = toProtoTicket(t)
	}
	
	totalPages := (totalCount + limit - 1) / limit
	
	return &supportpb.ListTicketsResponse{
		Tickets:     protoTickets,
		TotalCount:  int32(totalCount),
		Page:        int32(page),
		TotalPages:  int32(totalPages),
	}, nil
}

// Communication Methods
func (s server) AddTicketReply(ctx context.Context, request *supportpb.AddTicketReplyRequest) (*supportpb.AddTicketReplyResponse, error) {
	replyID := uuid.New().String()
	
	err := s.app.AddTicketReply(ctx, commands.AddTicketReply{
		ID:          replyID,
		TicketID:    request.GetTicketId(),
		AuthorID:    request.GetAuthorId(),
		AuthorType:  mapAuthorType(request.GetAuthorType()),
		Content:     request.GetContent(),
		Attachments: mapAttachments(request.GetAttachments()),
		IsPublic:    request.GetIsPublic(),
	})
	
	if err != nil {
		return nil, err
	}
	
	return &supportpb.AddTicketReplyResponse{
		Id: replyID,
	}, nil
}

func (s server) GetTicketCommunications(ctx context.Context, request *supportpb.GetTicketCommunicationsRequest) (*supportpb.GetTicketCommunicationsResponse, error) {
	comms, err := s.app.GetTicketCommunications(ctx, queries.GetTicketCommunications{
		TicketID:        request.GetTicketId(),
		IncludeInternal: request.GetIncludeInternal(),
		Page:            int(request.GetPage()),
		Limit:           int(request.GetLimit()),
	})
	
	if err != nil {
		return nil, err
	}
	
	protoComms := make([]*supportpb.Communication, len(comms))
	for i, comm := range comms {
		protoComms[i] = toProtoCommunication(comm)
	}
	
	return &supportpb.GetTicketCommunicationsResponse{
		Communications: protoComms,
		TotalCount:     int32(len(comms)),
	}, nil
}

func (s server) UpdateSupportChannelSettings(ctx context.Context, request *supportpb.UpdateSupportChannelSettingsRequest) (*supportpb.UpdateSupportChannelSettingsResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("channelID", request.GetId()),
	)
	
	settings := domain.SupportChannelSettings{
		EmailNotifications: request.GetSettings().GetEmailNotifications(),
		SMSNotifications:   request.GetSettings().GetSmsNotifications(),
		AutoAssignTickets:  request.GetSettings().GetAutoAssignTickets(),
		PreferredLanguage:  request.GetSettings().GetPreferredLanguage(),
		Timezone:           request.GetSettings().GetTimezone(),
		NotificationEmails: request.GetSettings().GetNotificationEmails(),
		SLASettings: domain.SLASettings{
			FirstResponseMinutes: int(request.GetSettings().GetSlaSettings().GetFirstResponseMinutes()),
			ResolutionHours:      int(request.GetSettings().GetSlaSettings().GetResolutionHours()),
			PriorityResponseTimes: convertInt32MapToIntMap(request.GetSettings().GetSlaSettings().GetPriorityResponseTimes()),
		},
	}
	
	err := s.app.UpdateSupportChannelSettings(ctx, commands.UpdateSupportChannelSettings{
		ID:       request.GetId(),
		Settings: settings,
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	return &supportpb.UpdateSupportChannelSettingsResponse{}, nil
}

func (s server) CloseSupportChannel(ctx context.Context, request *supportpb.CloseSupportChannelRequest) (*supportpb.CloseSupportChannelResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("channelID", request.GetId()),
	)
	
	err := s.app.CloseSupportChannel(ctx, commands.CloseSupportChannel{
		ID:       request.GetId(),
		ClosedBy: getUserIDFromContext(ctx),
		Reason:   request.GetReason(),
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	return &supportpb.CloseSupportChannelResponse{}, nil
}

func (s server) ReactivateSupportChannel(ctx context.Context, request *supportpb.ReactivateSupportChannelRequest) (*supportpb.ReactivateSupportChannelResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("channelID", request.GetId()),
	)
	
	err := s.app.ReactivateSupportChannel(ctx, commands.ReactivateSupportChannel{
		ID:            request.GetId(),
		ReactivatedBy: getUserIDFromContext(ctx),
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	return &supportpb.ReactivateSupportChannelResponse{}, nil
}

func (s server) UpdateTicket(ctx context.Context, request *supportpb.UpdateTicketRequest) (*supportpb.UpdateTicketResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("ticketID", request.GetId()),
	)
	
	var category *domain.TicketCategory
	if request.Category != nil {
		c := mapTicketCategory(*request.Category)
		category = &c
	}
	
	err := s.app.UpdateTicket(ctx, commands.UpdateTicket{
		ID:          request.GetId(),
		Title:       request.Title,
		Description: request.Description,
		Category:    category,
		Tags:        request.GetTags(),
		Metadata:    request.GetMetadata(),
		UpdatedBy:   getUserIDFromContext(ctx),
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	return &supportpb.UpdateTicketResponse{}, nil
}

func (s server) AssignTicket(ctx context.Context, request *supportpb.AssignTicketRequest) (*supportpb.AssignTicketResponse, error) {
	err := s.app.AssignTicket(ctx, commands.AssignTicket{
		ID:               request.GetId(),
		AssigneeID:       request.GetAssigneeId(),
		AssigneeType:     mapAssigneeType(request.GetAssigneeType()),
		AssignedBy:       getUserIDFromContext(ctx),
		AssignmentReason: request.GetAssignmentReason(),
	})
	
	return &supportpb.AssignTicketResponse{}, err
}

func (s server) UpdateTicketPriority(ctx context.Context, request *supportpb.UpdateTicketPriorityRequest) (*supportpb.UpdateTicketPriorityResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("ticketID", request.GetId()),
		attribute.String("priority", request.GetPriority().String()),
	)
	
	err := s.app.UpdateTicketPriority(ctx, commands.UpdateTicketPriority{
		ID:       request.GetId(),
		Priority: mapTicketPriority(request.GetPriority()),
		Reason:   request.GetReason(),
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	return &supportpb.UpdateTicketPriorityResponse{}, nil
}

func (s server) EscalateTicket(ctx context.Context, request *supportpb.EscalateTicketRequest) (*supportpb.EscalateTicketResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("ticketID", request.GetId()),
		attribute.String("toTier", request.GetEscalationTier().String()),
	)
	
	err := s.app.EscalateTicket(ctx, commands.EscalateTicket{
		ID:               request.GetId(),
		EscalationTier:   mapSupportTier(request.GetEscalationTier()),
		EscalatedBy:      getUserIDFromContext(ctx),
		EscalationReason: request.GetEscalationReason(),
		EscalationNotes:  request.GetEscalationNotes(),
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	return &supportpb.EscalateTicketResponse{}, nil
}

func (s server) ResolveTicket(ctx context.Context, request *supportpb.ResolveTicketRequest) (*supportpb.ResolveTicketResponse, error) {
	err := s.app.ResolveTicket(ctx, commands.ResolveTicket{
		ID:               request.GetId(),
		ResolvedBy:       getUserIDFromContext(ctx),
		Resolution:       request.GetResolution(),
		AppliedSolutions: request.GetAppliedSolutions(),
	})
	
	return &supportpb.ResolveTicketResponse{}, err
}

func (s server) ReopenTicket(ctx context.Context, request *supportpb.ReopenTicketRequest) (*supportpb.ReopenTicketResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("ticketID", request.GetId()),
	)
	
	err := s.app.ReopenTicket(ctx, commands.ReopenTicket{
		ID:           request.GetId(),
		ReopenedBy:   getUserIDFromContext(ctx),
		ReopenReason: request.GetReopenReason(),
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	return &supportpb.ReopenTicketResponse{}, nil
}

func (s server) CloseTicket(ctx context.Context, request *supportpb.CloseTicketRequest) (*supportpb.CloseTicketResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("ticketID", request.GetId()),
	)
	
	var satisfaction *domain.CustomerSatisfaction
	if request.GetSatisfaction() != supportpb.CustomerSatisfaction_VERY_DISSATISFIED {
		s := mapCustomerSatisfaction(request.GetSatisfaction())
		satisfaction = &s
	}
	
	err := s.app.CloseTicket(ctx, commands.CloseTicket{
		ID:                 request.GetId(),
		ClosedBy:           getUserIDFromContext(ctx),
		ClosureNotes:       request.GetClosureNotes(),
		SatisfactionRating: satisfaction,
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	return &supportpb.CloseTicketResponse{}, nil
}

func (s server) MergeTickets(ctx context.Context, request *supportpb.MergeTicketsRequest) (*supportpb.MergeTicketsResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("primaryTicketID", request.GetPrimaryTicketId()),
		attribute.Int("secondaryCount", len(request.GetSecondaryTicketIds())),
	)
	
	err := s.app.MergeTickets(ctx, commands.MergeTickets{
		PrimaryTicketID:    request.GetPrimaryTicketId(),
		SecondaryTicketIDs: request.GetSecondaryTicketIds(),
		MergedBy:           getUserIDFromContext(ctx),
		MergeReason:        request.GetMergeReason(),
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	return &supportpb.MergeTicketsResponse{
		MergedTicketId: request.GetPrimaryTicketId(),
	}, nil
}

func (s server) LinkTickets(ctx context.Context, request *supportpb.LinkTicketsRequest) (*supportpb.LinkTicketsResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("ticketID", request.GetTicketId()),
		attribute.Int("relatedCount", len(request.GetRelatedTicketIds())),
	)
	
	err := s.app.LinkTickets(ctx, commands.LinkTickets{
		TicketID:         request.GetTicketId(),
		RelatedTicketIDs: request.GetRelatedTicketIds(),
		LinkedBy:         getUserIDFromContext(ctx),
		RelationshipType: request.GetRelationshipType(),
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	return &supportpb.LinkTicketsResponse{}, nil
}

func (s server) AddInternalNote(ctx context.Context, request *supportpb.AddInternalNoteRequest) (*supportpb.AddInternalNoteResponse, error) {
	span := trace.SpanFromContext(ctx)
	noteID := uuid.New().String()
	
	span.SetAttributes(
		attribute.String("noteID", noteID),
		attribute.String("ticketID", request.GetTicketId()),
	)
	
	err := s.app.AddInternalNote(ctx, commands.AddInternalNote{
		ID:             noteID,
		TicketID:       request.GetTicketId(),
		AuthorID:       getUserIDFromContext(ctx),
		Content:        request.GetContent(),
		MentionedUsers: request.GetMentionedUsers(),
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	return &supportpb.AddInternalNoteResponse{
		Id: noteID,
	}, nil
}

func (s server) EnableAISupport(ctx context.Context, request *supportpb.EnableAISupportRequest) (*supportpb.EnableAISupportResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("channelID", request.GetChannelId()),
		attribute.String("assistantID", request.GetAssistantId()),
	)
	
	// Map configuration
	config := domain.AIConfiguration{
		ChannelID:              request.GetChannelId(),
		AssistantID:            request.GetAssistantId(),
		AllowedActions:         request.GetConfiguration().GetAllowedActions(),
		KnowledgeBaseAccess:    request.GetConfiguration().GetKnowledgeBaseAccess(),
		MaxHandlingTier:        mapSupportTier(request.GetConfiguration().GetMaxHandlingTier()),
		CanCloseTickets:        request.GetConfiguration().GetCanCloseTickets(),
		CanIssueRefunds:        request.GetConfiguration().GetCanIssueRefunds(),
		ConfidenceThreshold:    request.GetConfiguration().GetConfidenceThreshold(),
		AutoResponseCategories: request.GetConfiguration().GetAutoResponseCategories(),
		MaxTokens:              int(request.GetConfiguration().GetMaxTokens()),
		Temperature:            request.GetConfiguration().GetTemperature(),
		Active:                 true,
	}
	
	err := s.app.EnableAISupport(ctx, commands.EnableAISupport{
		ChannelID:    request.GetChannelId(),
		AssistantID:  request.GetAssistantId(),
		Configuration: config,
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...)) 
		return nil, err
	}
	
	return &supportpb.EnableAISupportResponse{}, nil
}

func (s server) ConfigureAIAssistant(ctx context.Context, request *supportpb.ConfigureAIAssistantRequest) (*supportpb.ConfigureAIAssistantResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("channelID", request.GetChannelId()),
	)
	
	// Map configuration
	config := domain.AIConfiguration{
		ChannelID:              request.GetChannelId(),
		AssistantID:            request.GetConfiguration().GetAssistantId(),
		AllowedActions:         request.GetConfiguration().GetAllowedActions(),
		KnowledgeBaseAccess:    request.GetConfiguration().GetKnowledgeBaseAccess(),
		MaxHandlingTier:        mapSupportTier(request.GetConfiguration().GetMaxHandlingTier()),
		CanCloseTickets:        request.GetConfiguration().GetCanCloseTickets(),
		CanIssueRefunds:        request.GetConfiguration().GetCanIssueRefunds(),
		ConfidenceThreshold:    request.GetConfiguration().GetConfidenceThreshold(),
		AutoResponseCategories: request.GetConfiguration().GetAutoResponseCategories(),
		MaxTokens:              int(request.GetConfiguration().GetMaxTokens()),
		Temperature:            request.GetConfiguration().GetTemperature(),
		Active:                 true,
	}
	
	err := s.app.ConfigureAIAssistant(ctx, commands.ConfigureAIAssistant{
		ChannelID:     request.GetChannelId(),
		Configuration: config,
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...)) 
		return nil, err
	}
	
	return &supportpb.ConfigureAIAssistantResponse{}, nil
}

func (s server) HandoffToHuman(ctx context.Context, request *supportpb.HandoffToHumanRequest) (*supportpb.HandoffToHumanResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("ticketID", request.GetTicketId()),
		attribute.String("reason", request.GetReason()),
	)
	
	agentID, err := s.app.HandoffToHuman(ctx, commands.HandoffToHuman{
		TicketID:    request.GetTicketId(),
		Reason:      request.GetReason(),
		Context:     request.GetContext(),
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	return &supportpb.HandoffToHumanResponse{
		AssignedAgentId: agentID,
	}, nil
}

func (s server) HandoffToAI(ctx context.Context, request *supportpb.HandoffToAIRequest) (*supportpb.HandoffToAIResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("ticketID", request.GetTicketId()),
	)
	
	err := s.app.HandoffToAI(ctx, commands.HandoffToAI{
		TicketID:    request.GetTicketId(),
		Context:     request.GetContext(),
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	return &supportpb.HandoffToAIResponse{}, nil
}

func (s server) GetAISuggestions(ctx context.Context, request *supportpb.GetAISuggestionsRequest) (*supportpb.GetAISuggestionsResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("ticketID", request.GetTicketId()),
		attribute.String("suggestionType", request.GetSuggestionType().String()),
	)
	
	suggestions, err := s.app.GetAISuggestions(ctx, queries.GetAISuggestions{
		TicketID:       request.GetTicketId(),
		SuggestionType: mapSuggestionType(request.GetSuggestionType()),
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	// Convert domain suggestions to proto
	protoSuggestions := make([]*supportpb.AISuggestion, len(suggestions))
	for i, s := range suggestions {
		protoSuggestions[i] = &supportpb.AISuggestion{
			Id:         s.ID,
			Type:       mapToProtoSuggestionType(s.Type),
			Content:    s.Content,
			Confidence: s.Confidence,
			Reasoning:  s.Reasoning,
			Metadata:   s.Metadata,
		}
	}
	
	return &supportpb.GetAISuggestionsResponse{
		Suggestions: protoSuggestions,
	}, nil
}

func (s server) CreateKnowledgeArticle(ctx context.Context, request *supportpb.CreateKnowledgeArticleRequest) (*supportpb.CreateKnowledgeArticleResponse, error) {
	span := trace.SpanFromContext(ctx)
	articleID := uuid.New().String()
	
	span.SetAttributes(
		attribute.String("articleID", articleID),
		attribute.String("title", request.GetTitle()),
	)
	
	err := s.app.CreateKnowledgeArticle(ctx, commands.CreateKnowledgeArticle{
		ID:          articleID,
		Title:       request.GetTitle(),
		Content:     request.GetContent(),
		Categories:  request.GetCategories(),
		Tags:        request.GetTags(),
		Public:      request.GetPublic(),
		CreatedBy:   getUserIDFromContext(ctx),
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	return &supportpb.CreateKnowledgeArticleResponse{
		Id: articleID,
	}, nil
}

func (s server) GetKnowledgeArticle(ctx context.Context, request *supportpb.GetKnowledgeArticleRequest) (*supportpb.GetKnowledgeArticleResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("articleID", request.GetId()),
	)
	
	article, err := s.app.GetKnowledgeArticle(ctx, queries.GetKnowledgeArticle{
		ID: request.GetId(),
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	return &supportpb.GetKnowledgeArticleResponse{
		Article: &supportpb.KnowledgeArticle{
			Id:            article.ID,
			Title:         article.Title,
			Content:       "", // Not available in catalog
			Categories:    article.Categories,
			Tags:          []string{}, // Not available in catalog
			Public:        article.Public,
			ViewCount:     int32(article.ViewCount),
			AverageRating: article.AverageRating,
			RatingCount:   0, // Not available in catalog
			CreatedAt:     timestamppb.New(article.CreatedAt),
			UpdatedAt:     timestamppb.New(article.CreatedAt), // Use created_at as updated_at
			CreatedBy:     "", // Not available in catalog
			RelatedTicketIds: []string{}, // Not available in catalog
		},
	}, nil
}

func (s server) SearchKnowledgeBase(ctx context.Context, request *supportpb.SearchKnowledgeBaseRequest) (*supportpb.SearchKnowledgeBaseResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("query", request.GetQuery()),
		attribute.Int("limit", int(request.GetLimit())),
	)
	
	limit := int(request.GetLimit())
	if limit == 0 {
		limit = 10 // Default limit
	}
	
	articles, err := s.app.SearchKnowledgeBase(ctx, queries.SearchKnowledgeBase{
		Query:      request.GetQuery(),
		Categories: request.GetCategories(),
		Limit:      limit,
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	// Convert domain articles to proto
	protoArticles := make([]*supportpb.KnowledgeArticle, len(articles))
	for i, a := range articles {
		protoArticles[i] = &supportpb.KnowledgeArticle{
			Id:            a.ID,
			Title:         a.Title,
			Content:       "", // Not available in catalog
			Categories:    a.Categories,
			Tags:          []string{}, // Not available in catalog
			Public:        a.Public,
			ViewCount:     int32(a.ViewCount),
			AverageRating: a.AverageRating,
			RatingCount:   0, // Not available in catalog
			CreatedAt:     timestamppb.New(a.CreatedAt),
			UpdatedAt:     timestamppb.New(a.CreatedAt), // Use created_at as updated_at
			CreatedBy:     "", // Not available in catalog
			RelatedTicketIds: []string{}, // Not available in catalog
		}
	}
	
	return &supportpb.SearchKnowledgeBaseResponse{
		Articles: protoArticles,
	}, nil
}

func (s server) LinkArticleToTicket(ctx context.Context, request *supportpb.LinkArticleToTicketRequest) (*supportpb.LinkArticleToTicketResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("ticketID", request.GetTicketId()),
		attribute.String("articleID", request.GetArticleId()),
		attribute.Bool("helpedResolve", request.GetHelpedResolve()),
	)
	
	err := s.app.LinkArticleToTicket(ctx, commands.LinkArticleToTicket{
		TicketID:      request.GetTicketId(),
		ArticleID:     request.GetArticleId(),
		HelpedResolve: request.GetHelpedResolve(),
		LinkedBy:      getUserIDFromContext(ctx),
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	return &supportpb.LinkArticleToTicketResponse{}, nil
}

func (s server) RateArticle(ctx context.Context, request *supportpb.RateArticleRequest) (*supportpb.RateArticleResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("articleID", request.GetArticleId()),
		attribute.Int("rating", int(request.GetRating())),
	)
	
	err := s.app.RateArticle(ctx, commands.RateArticle{
		ArticleID: request.GetArticleId(),
		Rating:    int(request.GetRating()),
		Feedback:  request.GetFeedback(),
		RatedBy:   getUserIDFromContext(ctx),
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	return &supportpb.RateArticleResponse{}, nil
}

func (s server) GetSupportMetrics(ctx context.Context, request *supportpb.GetSupportMetricsRequest) (*supportpb.GetSupportMetricsResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	// Set default time range if not provided (last 30 days)
	startTime := request.GetStartTime()
	endTime := request.GetEndTime()
	if startTime == nil {
		defaultStart := timestamppb.New(time.Now().AddDate(0, 0, -30))
		startTime = defaultStart
	}
	if endTime == nil {
		defaultEnd := timestamppb.New(time.Now())
		endTime = defaultEnd
	}
	
	metrics, err := s.app.GetSupportMetrics(ctx, queries.GetSupportMetrics{
		StartTime: startTime.AsTime(),
		EndTime:   endTime.AsTime(),
		ChannelIDs: request.GetChannelIds(),
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	// Convert domain metrics to proto
	ticketsByCategory := make(map[string]int32)
	for k, v := range metrics.TicketsByCategory {
		ticketsByCategory[k] = int32(v)
	}
	
	ticketsByPriority := make(map[string]int32)
	for k, v := range metrics.TicketsByPriority {
		ticketsByPriority[k] = int32(v)
	}
	
	ticketsByStatus := make(map[string]int32)
	for k, v := range metrics.TicketsByStatus {
		ticketsByStatus[k] = int32(v)
	}
	
	return &supportpb.GetSupportMetricsResponse{
		Metrics: &supportpb.SupportMetrics{
			TotalTickets:                    int32(metrics.TotalTickets),
			OpenTickets:                     int32(metrics.OpenTickets),
			ResolvedTickets:                 int32(metrics.ResolvedTickets),
			EscalatedTickets:                int32(metrics.EscalatedTickets),
			AverageResolutionTimeHours:      metrics.AverageResolutionTimeHours,
			AverageFirstResponseTimeMinutes: metrics.AverageFirstResponseTimeMinutes,
			CustomerSatisfactionScore:       metrics.CustomerSatisfactionScore,
			TicketsByCategory:               ticketsByCategory,
			TicketsByPriority:               ticketsByPriority,
			TicketsByStatus:                 ticketsByStatus,
			AiResolvedTickets:               int32(metrics.AIResolvedTickets),
			AiResolutionRate:                metrics.AIResolutionRate,
		},
	}, nil
}

func (s server) GetAgentPerformance(ctx context.Context, request *supportpb.GetAgentPerformanceRequest) (*supportpb.GetAgentPerformanceResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("agentID", request.GetAgentId()),
	)
	
	// Set default time range if not provided (last 30 days)
	startTime := request.GetStartTime()
	endTime := request.GetEndTime()
	if startTime == nil {
		defaultStart := timestamppb.New(time.Now().AddDate(0, 0, -30))
		startTime = defaultStart
	}
	if endTime == nil {
		defaultEnd := timestamppb.New(time.Now())
		endTime = defaultEnd
	}
	
	performance, err := s.app.GetAgentPerformance(ctx, queries.GetAgentPerformance{
		AgentID:   request.GetAgentId(),
		StartTime: startTime.AsTime(),
		EndTime:   endTime.AsTime(),
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	// Convert tickets by category map
	ticketsByCategory := make(map[string]int32)
	for k, v := range performance.TicketsByCategory {
		ticketsByCategory[k] = int32(v)
	}
	
	return &supportpb.GetAgentPerformanceResponse{
		Performance: &supportpb.AgentPerformance{
			AgentId:                         performance.AgentID,
			AgentName:                       performance.AgentName,
			TicketsHandled:                  int32(performance.TicketsHandled),
			TicketsResolved:                 int32(performance.TicketsResolved),
			AverageResolutionTimeHours:      performance.AverageResolutionTimeHours,
			AverageFirstResponseTimeMinutes: performance.AverageFirstResponseTimeMinutes,
			CustomerSatisfactionScore:       performance.CustomerSatisfactionScore,
			EscalationsReceived:             int32(performance.EscalationsReceived),
			EscalationsSent:                 int32(performance.EscalationsSent),
			TicketsByCategory:               ticketsByCategory,
			PeriodStart:                     timestamppb.New(performance.PeriodStart),
			PeriodEnd:                       timestamppb.New(performance.PeriodEnd),
		},
	}, nil
}

func (s server) GetTicketAnalytics(ctx context.Context, request *supportpb.GetTicketAnalyticsRequest) (*supportpb.GetTicketAnalyticsResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.Int("ticketCount", len(request.GetTicketIds())),
	)
	
	analytics, err := s.app.GetTicketAnalytics(ctx, queries.GetTicketAnalytics{
		TicketIDs: request.GetTicketIds(),
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, err
	}
	
	// Convert time in status map
	timeInStatus := make(map[string]float64)
	for k, v := range analytics.TimeInStatus {
		timeInStatus[k] = v
	}
	
	return &supportpb.GetTicketAnalyticsResponse{
		Analytics: &supportpb.TicketAnalytics{
			TicketId:                    analytics.TicketID,
			TimeToFirstResponseMinutes:  analytics.TimeToFirstResponseMinutes,
			TotalResolutionTimeHours:    analytics.TotalResolutionTimeHours,
			TotalCommunications:         int32(analytics.TotalCommunications),
			AgentCommunications:         int32(analytics.AgentCommunications),
			CustomerCommunications:      int32(analytics.CustomerCommunications),
			AiCommunications:            int32(analytics.AICommunications),
			EscalationCount:             int32(analytics.EscalationCount),
			ReopenCount:                 int32(analytics.ReopenCount),
			AgentsInvolved:              analytics.AgentsInvolved,
			TimeInStatus:                timeInStatus,
		},
	}, nil
}

// Helper functions
func getUserIDFromContext(ctx context.Context) string {
	// TODO: Extract from auth context
	return "user-id"
}

func mapChannelType(t supportpb.SupportChannelType) domain.SupportChannelType {
	switch t {
	case supportpb.SupportChannelType_TECHNICAL:
		return domain.ChannelTypeTechnical
	case supportpb.SupportChannelType_BILLING:
		return domain.ChannelTypeBilling
	case supportpb.SupportChannelType_SALES:
		return domain.ChannelTypeSales
	case supportpb.SupportChannelType_PRODUCT:
		return domain.ChannelTypeProduct
	default:
		return domain.ChannelTypeGeneral
	}
}

func mapTicketCategory(c supportpb.TicketCategory) domain.TicketCategory {
	switch c {
	case supportpb.TicketCategory_TECHNICAL_ISSUE:
		return domain.CategoryTechnicalIssue
	case supportpb.TicketCategory_BILLING_ISSUE:
		return domain.CategoryBillingIssue
	case supportpb.TicketCategory_ACCOUNT_ISSUE:
		return domain.CategoryAccountIssue
	case supportpb.TicketCategory_PRODUCT_QUESTION:
		return domain.CategoryProductQuestion
	case supportpb.TicketCategory_FEATURE_REQUEST:
		return domain.CategoryFeatureRequest
	case supportpb.TicketCategory_COMPLAINT:
		return domain.CategoryComplaint
	case supportpb.TicketCategory_REFUND_REQUEST:
		return domain.CategoryRefundRequest
	case supportpb.TicketCategory_ORDER_ISSUE:
		return domain.CategoryOrderIssue
	case supportpb.TicketCategory_SHIPPING_ISSUE:
		return domain.CategoryShippingIssue
	default:
		return domain.CategoryGeneralInquiry
	}
}

func mapTicketPriority(p supportpb.TicketPriority) domain.TicketPriority {
	switch p {
	case supportpb.TicketPriority_HIGH:
		return domain.PriorityHigh
	case supportpb.TicketPriority_URGENT:
		return domain.PriorityUrgent
	case supportpb.TicketPriority_CRITICAL:
		return domain.PriorityCritical
	case supportpb.TicketPriority_LOW:
		return domain.PriorityLow
	default:
		return domain.PriorityMedium
	}
}

func mapAuthorType(t supportpb.AuthorType) domain.AuthorType {
	switch t {
	case supportpb.AuthorType_AGENT:
		return domain.AuthorTypeAgent
	case supportpb.AuthorType_AI:
		return domain.AuthorTypeAI
	case supportpb.AuthorType_SYSTEM:
		return domain.AuthorTypeSystem
	default:
		return domain.AuthorTypeCustomer
	}
}

func mapAssigneeType(t supportpb.AssigneeType) domain.AssigneeType {
	switch t {
	case supportpb.AssigneeType_AI_ASSISTANT:
		return domain.AssigneeTypeAI
	case supportpb.AssigneeType_TEAM:
		return domain.AssigneeTypeTeam
	default:
		return domain.AssigneeTypeHuman
	}
}

func mapAttachments(attachments []*supportpb.Attachment) []domain.Attachment {
	result := make([]domain.Attachment, len(attachments))
	for i, a := range attachments {
		result[i] = domain.Attachment{
			ID:          a.GetId(),
			Filename:    a.GetFilename(),
			ContentType: a.GetContentType(),
			SizeBytes:   a.GetSizeBytes(),
			URL:         a.GetUrl(),
			UploadedAt:  a.GetUploadedAt().AsTime(),
		}
	}
	return result
}

func mapSuggestionType(t supportpb.SuggestionType) domain.SuggestionType {
	switch t {
	case supportpb.SuggestionType_RESPONSE:
		return domain.SuggestionTypeResponse
	case supportpb.SuggestionType_KNOWLEDGE_ARTICLE:
		return domain.SuggestionTypeKnowledgeArticle
	case supportpb.SuggestionType_SIMILAR_TICKET:
		return domain.SuggestionTypeSimilarTicket
	case supportpb.SuggestionType_ESCALATION:
		return domain.SuggestionTypeEscalation
	default:
		return domain.SuggestionTypeResponse
	}
}

func mapToProtoSuggestionType(t domain.SuggestionType) supportpb.SuggestionType {
	switch t {
	case domain.SuggestionTypeResponse:
		return supportpb.SuggestionType_RESPONSE
	case domain.SuggestionTypeKnowledgeArticle:
		return supportpb.SuggestionType_KNOWLEDGE_ARTICLE
	case domain.SuggestionTypeSimilarTicket:
		return supportpb.SuggestionType_SIMILAR_TICKET
	case domain.SuggestionTypeEscalation:
		return supportpb.SuggestionType_ESCALATION
	default:
		return supportpb.SuggestionType_RESPONSE
	}
}


func toProtoChannel(c *domain.SupportChannelCatalog) *supportpb.SupportChannel {
	return &supportpb.SupportChannel{
		Id:           c.ID,
		UserId:       c.UserID,
		BusinessId:   c.BusinessID,
		ChannelType:  supportpb.SupportChannelType(supportpb.SupportChannelType_value[c.ChannelType]),
		Active:       c.Active,
		Settings:     &supportpb.SupportChannelSettings{},
		CreatedAt:    timestamppb.New(c.CreatedAt),
		UpdatedAt:    timestamppb.New(c.UpdatedAt),
		OpenTickets:  int32(c.OpenTickets),
		TotalTickets: int32(c.TotalTickets),
	}
}

func toProtoTicket(t *domain.TicketCatalog) *supportpb.Ticket {
	ticket := &supportpb.Ticket{
		Id:          t.ID,
		ChannelId:   t.ChannelID,
		Title:       t.Title,
		Description: "", // Not in catalog
		Status:      supportpb.TicketStatus(supportpb.TicketStatus_value[t.Status]),
		Priority:    supportpb.TicketPriority(supportpb.TicketPriority_value[t.Priority]),
		Category:    supportpb.TicketCategory(supportpb.TicketCategory_value[t.Category]),
		CreatedBy:   t.CreatedBy,
		CreatedAt:   timestamppb.New(t.CreatedAt),
		UpdatedAt:   timestamppb.New(t.UpdatedAt),
	}
	
	if t.AssigneeID != nil {
		ticket.AssigneeId = *t.AssigneeID
	}
	if t.AssigneeType != nil {
		ticket.AssigneeType = supportpb.AssigneeType(supportpb.AssigneeType_value[*t.AssigneeType])
	}
	
	return ticket
}

func toProtoCommunication(c *domain.Communication) *supportpb.Communication {
	return &supportpb.Communication{
		Id:             c.ID,
		TicketId:       c.TicketID,
		AuthorId:       c.AuthorID,
		AuthorType:     supportpb.AuthorType(supportpb.AuthorType_value[string(c.AuthorType)]),
		Content:        c.Content,
		IsPublic:       c.IsPublic,
		Attachments:    toProtoAttachments(c.Attachments),
		CreatedAt:      timestamppb.New(c.CreatedAt),
		MentionedUsers: c.MentionedUsers,
		Metadata:       c.Metadata,
	}
}

func toProtoAttachments(attachments []domain.Attachment) []*supportpb.Attachment {
	result := make([]*supportpb.Attachment, len(attachments))
	for i, a := range attachments {
		result[i] = &supportpb.Attachment{
			Id:          a.ID,
			Filename:    a.Filename,
			ContentType: a.ContentType,
			SizeBytes:   a.SizeBytes,
			Url:         a.URL,
			UploadedAt:  timestamppb.New(a.UploadedAt),
		}
	}
	return result
}

func mapSupportTier(t supportpb.SupportTier) domain.SupportTier {
	switch t {
	case supportpb.SupportTier_TIER_2:
		return domain.TierTwo
	case supportpb.SupportTier_TIER_3:
		return domain.TierThree
	case supportpb.SupportTier_MANAGEMENT:
		return domain.TierManagement
	default:
		return domain.TierOne
	}
}

func mapCustomerSatisfaction(s supportpb.CustomerSatisfaction) domain.CustomerSatisfaction {
	switch s {
	case supportpb.CustomerSatisfaction_VERY_DISSATISFIED:
		return domain.SatisfactionVeryDissatisfied
	case supportpb.CustomerSatisfaction_DISSATISFIED:
		return domain.SatisfactionDissatisfied
	case supportpb.CustomerSatisfaction_NEUTRAL:
		return domain.SatisfactionNeutral
	case supportpb.CustomerSatisfaction_SATISFIED:
		return domain.SatisfactionSatisfied
	case supportpb.CustomerSatisfaction_VERY_SATISFIED:
		return domain.SatisfactionVerySatisfied
	default:
		return domain.SatisfactionNeutral
	}
}

func convertInt32MapToIntMap(m map[string]int32) map[string]int {
	result := make(map[string]int)
	for k, v := range m {
		result[k] = int(v)
	}
	return result
}