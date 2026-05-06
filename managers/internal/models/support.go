package models

import "time"

// Support represents a support channel associated with a user
type Support struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

// Ticket represents an individual support request or issue within a support channel
type Ticket struct {
	ID           string       `json:"id"`
	SupportID    string       `json:"support_id"`
	Title        string       `json:"title"`
	Description  string       `json:"description"`
	TicketStatus TicketStatus `json:"ticket_status"`
	CreatedAt    time.Time    `json:"created_at,omitempty"`
	UpdatedAt    time.Time    `json:"updated_at,omitempty"`
}

// TicketStatus represents the possible statuses of a support ticket
type TicketStatus string

const (
	TicketStatusOpen       TicketStatus = "open"
	TicketStatusInProgress TicketStatus = "in_progress"
	TicketStatusResolved   TicketStatus = "resolved"
	TicketStatusClosed     TicketStatus = "closed"
	TicketStatusCancelled  TicketStatus = "cancelled"
	TicketStatusUnknown    TicketStatus = ""
)

func (s TicketStatus) String() string {
	switch s {
	case TicketStatusOpen, TicketStatusInProgress, TicketStatusResolved, TicketStatusClosed, TicketStatusCancelled:
		return string(s)
	default:
		return ""
	}
}

func ToTicketStatus(s string) TicketStatus {
	switch s {
	case TicketStatusOpen.String():
		return TicketStatusOpen
	case TicketStatusInProgress.String():
		return TicketStatusInProgress
	case TicketStatusResolved.String():
		return TicketStatusResolved
	case TicketStatusClosed.String():
		return TicketStatusClosed
	case TicketStatusCancelled.String():
		return TicketStatusCancelled
	default:
		return TicketStatusUnknown
	}
}

// StartSupportResponse represents the response after starting a support channel
type StartSupportResponse struct {
	ID string `json:"id"`
}

// CreateTicketResponse represents the response after creating a ticket
type CreateTicketResponse struct {
	ID string `json:"id"`
}

// ListTicketsResponse represents the response for listing tickets
type ListTicketsResponse struct {
	Tickets []*Ticket `json:"tickets"`
	Total   int64     `json:"total"`
	Page    int64     `json:"page"`
	Limit   int64     `json:"limit"`
}

// GetTicketResponse represents the response for getting a specific ticket
type GetTicketResponse struct {
	Ticket *Ticket `json:"ticket"`
}

// UpdateTicketResponse represents the response after updating a ticket
type UpdateTicketResponse struct {
	Ticket *Ticket `json:"ticket"`
}

// DeleteTicketResponse represents the response after deleting a ticket
type DeleteTicketResponse struct {
}

// CloseTicketResponse represents the response after closing a ticket
type CloseTicketResponse struct {
	Ticket *Ticket `json:"ticket"`
}

// GetTicketsResponse represents the response for getting tickets
type GetTicketsResponse struct {
	Tickets []*Ticket `json:"tickets"`
	Total   int64     `json:"total"`
}

// GetUserSupportResponse represents the response for getting user support
type GetUserSupportResponse struct {
	Support *Support `json:"support"`
}

// TicketComment represents a comment on a support ticket
type TicketComment struct {
	ID        string    `json:"id"`
	TicketID  string    `json:"ticket_id"`
	AuthorID  string    `json:"author_id"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

// TicketHistoryEntry represents a history entry for a support ticket
type TicketHistoryEntry struct {
	ID        string    `json:"id"`
	TicketID  string    `json:"ticket_id"`
	Action    string    `json:"action"`
	Details   string    `json:"details"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
