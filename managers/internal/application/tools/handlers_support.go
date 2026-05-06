package tools

import (
	"context"
	"fmt"
)

// ==================== SUPPORT HANDLERS ====================
func (r *ComprehensiveToolRegistry) initializeSupportHandlers() {
	r.handlers["support_start"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		if userID == "" {
			return nil, fmt.Errorf("user_id is required")
		}
		return reg.supportRepo.InitiateSupportSessionForUser(ctx, userID)
	}

	r.handlers["support_create_ticket"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		supportID := getStringParam(params, "support_id")
		title := getStringParam(params, "title")
		description := getStringParam(params, "description")
		if supportID == "" || title == "" || description == "" {
			return nil, fmt.Errorf("support_id, title, and description are required")
		}
		return reg.supportRepo.CreateNewSupportTicket(ctx, supportID, title, description)
	}

	r.handlers["support_get_tickets"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		supportID := getStringParam(params, "support_id")
		if supportID == "" {
			return nil, fmt.Errorf("support_id is required")
		}
		return reg.supportRepo.GetAllTicketsForSupport(ctx, supportID)
	}

	r.handlers["support_get_ticket"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		ticketID := getStringParam(params, "ticket_id")
		if ticketID == "" {
			return nil, fmt.Errorf("ticket_id is required")
		}
		return reg.supportRepo.GetSupportTicketByID(ctx, ticketID)
	}

	r.handlers["support_update_ticket"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		ticketID := getStringParam(params, "ticket_id")
		status := getStringParam(params, "status")
		assignedTo := getStringParam(params, "assigned_to")
		if ticketID == "" || status == "" {
			return nil, fmt.Errorf("ticket_id and status are required")
		}
		return reg.supportRepo.ModifySupportTicketStatus(ctx, ticketID, status, assignedTo)
	}

	r.handlers["support_close_ticket"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		ticketID := getStringParam(params, "ticket_id")
		reason := getStringParam(params, "reason")
		if ticketID == "" {
			return nil, fmt.Errorf("ticket_id is required")
		}
		return reg.supportRepo.CloseSupportTicketWithReason(ctx, ticketID, reason)
	}

	r.handlers["support_get_user_support"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		if userID == "" {
			return nil, fmt.Errorf("user_id is required")
		}
		return reg.supportRepo.GetUserSupportInformation(ctx, userID)
	}

	r.handlers["support_get_by_status"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		status := getStringParam(params, "status")
		limit := getInt64Param(params, "limit", 50)
		if status == "" {
			return nil, fmt.Errorf("status is required")
		}
		return reg.supportRepo.FilterTicketsByCurrentStatus(ctx, status, limit)
	}

	r.handlers["support_search_tickets"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		query := getStringParam(params, "query")
		limit := getInt64Param(params, "limit", 50)
		if query == "" {
			return nil, fmt.Errorf("query is required")
		}
		return reg.supportRepo.SearchTicketsByKeyword(ctx, query, limit)
	}

	// Additional handlers using methods from the repository interface
	r.handlers["support_find_session"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		supportID := getStringParam(params, "support_id")
		if supportID == "" {
			return nil, fmt.Errorf("support_id is required")
		}
		return reg.supportRepo.FindSupportSessionByID(ctx, supportID)
	}

	r.handlers["support_get_user_history"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		limit := getInt64Param(params, "limit", 10)
		if userID == "" {
			return nil, fmt.Errorf("user_id is required")
		}
		return reg.supportRepo.GetUserSupportSessionHistory(ctx, userID, limit)
	}

	r.handlers["support_get_open_tickets"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		limit := getInt64Param(params, "limit", 50)
		return reg.supportRepo.GetAllOpenSupportTickets(ctx, limit)
	}

	r.handlers["support_get_resolved_tickets"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		limit := getInt64Param(params, "limit", 50)
		return reg.supportRepo.GetAllResolvedSupportTickets(ctx, limit)
	}

	r.handlers["support_assign_ticket"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		ticketID := getStringParam(params, "ticket_id")
		assignedTo := getStringParam(params, "assigned_to")
		if ticketID == "" || assignedTo == "" {
			return nil, fmt.Errorf("ticket_id and assigned_to are required")
		}
		return nil, reg.supportRepo.AssignTicketToSupportAgent(ctx, ticketID, assignedTo)
	}

	r.handlers["support_add_comment"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		ticketID := getStringParam(params, "ticket_id")
		comment := getStringParam(params, "comment")
		authorID := getStringParam(params, "author_id")
		if ticketID == "" || comment == "" || authorID == "" {
			return nil, fmt.Errorf("ticket_id, comment, and author_id are required")
		}
		return nil, reg.supportRepo.AddCommentToSupportTicket(ctx, ticketID, comment, authorID)
	}

	r.handlers["support_get_comments"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		ticketID := getStringParam(params, "ticket_id")
		if ticketID == "" {
			return nil, fmt.Errorf("ticket_id is required")
		}
		return reg.supportRepo.GetAllCommentsForTicket(ctx, ticketID)
	}

	r.handlers["support_escalate_ticket"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		ticketID := getStringParam(params, "ticket_id")
		reason := getStringParam(params, "reason")
		if ticketID == "" || reason == "" {
			return nil, fmt.Errorf("ticket_id and reason are required")
		}
		return nil, reg.supportRepo.EscalateTicketToPriority(ctx, ticketID, reason)
	}

	r.handlers["support_get_ticket_history"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		ticketID := getStringParam(params, "ticket_id")
		if ticketID == "" {
			return nil, fmt.Errorf("ticket_id is required")
		}
		return reg.supportRepo.GetTicketActivityHistory(ctx, ticketID)
	}
}