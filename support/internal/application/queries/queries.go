package queries

import (
	"context"
	"middleman/support/internal/domain"
	"time"
)

// Support Channel Queries
type GetUserSupportChannels struct {
	UserID     string
	ActiveOnly bool
	Page       int
	Limit      int
}

type GetUserSupportChannelsHandler struct {
	catalog domain.SupportChannelCatalogRepository
}

func NewGetUserSupportChannelsHandler(catalog domain.SupportChannelCatalogRepository) GetUserSupportChannelsHandler {
	return GetUserSupportChannelsHandler{catalog: catalog}
}

func (h GetUserSupportChannelsHandler) GetUserSupportChannels(ctx context.Context, query GetUserSupportChannels) ([]*domain.SupportChannelCatalog, error) {
	limit := query.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := (query.Page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	return h.catalog.GetByUserID(ctx, query.UserID, query.ActiveOnly, limit, offset)
}

// Ticket Queries
type GetTickets struct {
	IDs []string
}

type GetTicketsHandler struct {
	catalog domain.TicketCatalogRepository
}

func NewGetTicketsHandler(catalog domain.TicketCatalogRepository) GetTicketsHandler {
	return GetTicketsHandler{catalog: catalog}
}

func (h GetTicketsHandler) GetTickets(ctx context.Context, query GetTickets) ([]*domain.TicketCatalog, error) {
	tickets := make([]*domain.TicketCatalog, 0, len(query.IDs))
	for _, id := range query.IDs {
		ticket, err := h.catalog.Find(ctx, id)
		if err != nil {
			continue // Skip if not found
		}
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}

type GetChannelTickets struct {
	ChannelID    string
	StatusFilter *string
	Page         int
	Limit        int
	SortBy       string
	Descending   bool
}

type GetChannelTicketsHandler struct {
	catalog domain.TicketCatalogRepository
}

func NewGetChannelTicketsHandler(catalog domain.TicketCatalogRepository) GetChannelTicketsHandler {
	return GetChannelTicketsHandler{catalog: catalog}
}

func (h GetChannelTicketsHandler) GetChannelTickets(ctx context.Context, query GetChannelTickets) ([]*domain.TicketCatalog, error) {
	limit := query.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := (query.Page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	return h.catalog.GetByChannelID(ctx, query.ChannelID, query.StatusFilter, limit, offset)
}

type GetTicketCommunications struct {
	TicketID        string
	IncludeInternal bool
	Page            int
	Limit           int
}

type GetTicketCommunicationsHandler struct {
	communications domain.CommunicationRepository
}

func NewGetTicketCommunicationsHandler(communications domain.CommunicationRepository) GetTicketCommunicationsHandler {
	return GetTicketCommunicationsHandler{communications: communications}
}

func (h GetTicketCommunicationsHandler) GetTicketCommunications(ctx context.Context, query GetTicketCommunications) ([]*domain.Communication, error) {
	limit := query.Limit
	if limit <= 0 {
		limit = 50
	}
	offset := (query.Page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	return h.communications.GetByTicketID(ctx, query.TicketID, query.IncludeInternal, limit, offset)
}

// AI Queries
type GetAISuggestions struct {
	TicketID       string
	SuggestionType domain.SuggestionType
}

type GetAISuggestionsHandler struct {
	// This would integrate with AI service
}

func NewGetAISuggestionsHandler() GetAISuggestionsHandler {
	return GetAISuggestionsHandler{}
}

func (h GetAISuggestionsHandler) GetAISuggestions(ctx context.Context, query GetAISuggestions) ([]*domain.AISuggestion, error) {
	// TODO: Implement AI integration
	return []*domain.AISuggestion{}, nil
}

// Knowledge Base Queries
type GetKnowledgeArticle struct {
	ID string
}

type GetKnowledgeArticleHandler struct {
	catalog domain.KnowledgeArticleCatalogRepository
}

func NewGetKnowledgeArticleHandler(catalog domain.KnowledgeArticleCatalogRepository) GetKnowledgeArticleHandler {
	return GetKnowledgeArticleHandler{catalog: catalog}
}

func (h GetKnowledgeArticleHandler) GetKnowledgeArticle(ctx context.Context, query GetKnowledgeArticle) (*domain.KnowledgeArticleCatalog, error) {
	return h.catalog.Find(ctx, query.ID)
}

type SearchKnowledgeBase struct {
	Query      string
	Categories []string
	Limit      int
}

type SearchKnowledgeBaseHandler struct {
	catalog domain.KnowledgeArticleCatalogRepository
}

func NewSearchKnowledgeBaseHandler(catalog domain.KnowledgeArticleCatalogRepository) SearchKnowledgeBaseHandler {
	return SearchKnowledgeBaseHandler{catalog: catalog}
}

func (h SearchKnowledgeBaseHandler) SearchKnowledgeBase(ctx context.Context, query SearchKnowledgeBase) ([]*domain.KnowledgeArticleCatalog, error) {
	limit := query.Limit
	if limit <= 0 {
		limit = 10
	}
	return h.catalog.Search(ctx, query.Query, query.Categories, true, limit, 0)
}

// Analytics Queries
type GetSupportMetrics struct {
	StartTime  time.Time
	EndTime    time.Time
	ChannelIDs []string
}

type GetSupportMetricsHandler struct {
	ticketCatalog domain.TicketCatalogRepository
}

func NewGetSupportMetricsHandler(ticketCatalog domain.TicketCatalogRepository) GetSupportMetricsHandler {
	return GetSupportMetricsHandler{ticketCatalog: ticketCatalog}
}

func (h GetSupportMetricsHandler) GetSupportMetrics(ctx context.Context, query GetSupportMetrics) (*domain.SupportMetrics, error) {
	// TODO: Implement metrics aggregation
	return &domain.SupportMetrics{}, nil
}

type GetAgentPerformance struct {
	AgentID   string
	StartTime time.Time
	EndTime   time.Time
}

type GetAgentPerformanceHandler struct {
	ticketCatalog domain.TicketCatalogRepository
}

func NewGetAgentPerformanceHandler(ticketCatalog domain.TicketCatalogRepository) GetAgentPerformanceHandler {
	return GetAgentPerformanceHandler{ticketCatalog: ticketCatalog}
}

func (h GetAgentPerformanceHandler) GetAgentPerformance(ctx context.Context, query GetAgentPerformance) (*domain.AgentPerformance, error) {
	// TODO: Implement performance tracking
	return &domain.AgentPerformance{
		AgentID: query.AgentID,
		PeriodStart: query.StartTime,
		PeriodEnd: query.EndTime,
	}, nil
}

type GetTicketAnalytics struct {
	TicketIDs []string
}

type GetTicketAnalyticsHandler struct {
	ticketCatalog  domain.TicketCatalogRepository
	communications domain.CommunicationRepository
}

func NewGetTicketAnalyticsHandler(
	ticketCatalog domain.TicketCatalogRepository,
	communications domain.CommunicationRepository,
) GetTicketAnalyticsHandler {
	return GetTicketAnalyticsHandler{
		ticketCatalog:  ticketCatalog,
		communications: communications,
	}
}

func (h GetTicketAnalyticsHandler) GetTicketAnalytics(ctx context.Context, query GetTicketAnalytics) (*domain.TicketAnalytics, error) {
	// TODO: Implement analytics aggregation
	return &domain.TicketAnalytics{}, nil
}