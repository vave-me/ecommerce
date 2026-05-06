package domain

import "time"

// AISuggestion represents a suggestion from AI
type AISuggestion struct {
	ID         string
	Type       SuggestionType
	Content    string
	Confidence float64
	Reasoning  string
	Metadata   map[string]string
}

type SuggestionType string

const (
	SuggestionTypeResponse         SuggestionType = "RESPONSE"
	SuggestionTypeKnowledgeArticle SuggestionType = "KNOWLEDGE_ARTICLE"
	SuggestionTypeSimilarTicket    SuggestionType = "SIMILAR_TICKET"
	SuggestionTypeEscalation       SuggestionType = "ESCALATION"
)

// SupportMetrics represents aggregated support metrics
type SupportMetrics struct {
	TotalTickets                     int
	OpenTickets                      int
	ResolvedTickets                  int
	EscalatedTickets                 int
	AverageResolutionTimeHours       float64
	AverageFirstResponseTimeMinutes  float64
	CustomerSatisfactionScore        float64
	TicketsByCategory                map[string]int
	TicketsByPriority                map[string]int
	TicketsByStatus                  map[string]int
	AIResolvedTickets                int
	AIResolutionRate                 float64
}

// AgentPerformance represents performance metrics for an agent
type AgentPerformance struct {
	AgentID                          string
	AgentName                        string
	TicketsHandled                   int
	TicketsResolved                  int
	AverageResolutionTimeHours       float64
	AverageFirstResponseTimeMinutes  float64
	CustomerSatisfactionScore        float64
	EscalationsReceived              int
	EscalationsSent                  int
	TicketsByCategory                map[string]int
	PeriodStart                      time.Time
	PeriodEnd                        time.Time
}

// TicketAnalytics represents analytics for a specific ticket
type TicketAnalytics struct {
	TicketID                      string
	TimeToFirstResponseMinutes    float64
	TotalResolutionTimeHours      float64
	TotalCommunications           int
	AgentCommunications           int
	CustomerCommunications        int
	AICommunications              int
	EscalationCount               int
	ReopenCount                   int
	AgentsInvolved                []string
	TimeInStatus                  map[string]float64
}