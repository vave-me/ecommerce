package postgres

import (
	"context"
	"database/sql"
	"fmt"
	
	"github.com/stackus/errors"
	"middleman/support/internal/domain"
)

type TicketCatalogRepository struct {
	tableName string
	db        *sql.DB
}

var _ domain.TicketCatalogRepository = (*TicketCatalogRepository)(nil)

func NewTicketCatalogRepository(tableName string, db *sql.DB) TicketCatalogRepository {
	return TicketCatalogRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r TicketCatalogRepository) Add(ctx context.Context, ticket *domain.TicketCatalog) error {
	const query = `
		INSERT INTO %s (id, channel_id, title, status, priority, category, assignee_id, assignee_type, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	
	_, err := r.db.ExecContext(ctx, r.table(query),
		ticket.ID,
		ticket.ChannelID,
		ticket.Title,
		ticket.Status,
		ticket.Priority,
		ticket.Category,
		ticket.AssigneeID,
		ticket.AssigneeType,
		ticket.CreatedBy,
		ticket.CreatedAt,
		ticket.UpdatedAt,
	)
	
	return err
}

func (r TicketCatalogRepository) Update(ctx context.Context, ticket *domain.TicketCatalog) error {
	const query = `
		UPDATE %s 
		SET title = $2, status = $3, priority = $4, category = $5, assignee_id = $6, assignee_type = $7, updated_at = $8
		WHERE id = $1
	`
	
	_, err := r.db.ExecContext(ctx, r.table(query),
		ticket.ID,
		ticket.Title,
		ticket.Status,
		ticket.Priority,
		ticket.Category,
		ticket.AssigneeID,
		ticket.AssigneeType,
		ticket.UpdatedAt,
	)
	
	return err
}

func (r TicketCatalogRepository) Delete(ctx context.Context, ticketID string) error {
	const query = `DELETE FROM %s WHERE id = $1`
	
	_, err := r.db.ExecContext(ctx, r.table(query), ticketID)
	
	return err
}

func (r TicketCatalogRepository) Find(ctx context.Context, ticketID string) (*domain.TicketCatalog, error) {
	const query = `
		SELECT id, channel_id, title, status, priority, category, assignee_id, assignee_type, created_by, created_at, updated_at
		FROM %s
		WHERE id = $1
	`
	
	ticket := &domain.TicketCatalog{}
	
	err := r.db.QueryRowContext(ctx, r.table(query), ticketID).Scan(
		&ticket.ID,
		&ticket.ChannelID,
		&ticket.Title,
		&ticket.Status,
		&ticket.Priority,
		&ticket.Category,
		&ticket.AssigneeID,
		&ticket.AssigneeType,
		&ticket.CreatedBy,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(errors.ErrNotFound, "ticket not found")
		}
		return nil, err
	}
	
	return ticket, nil
}

func (r TicketCatalogRepository) GetByChannelID(ctx context.Context, channelID string, status *string, limit, offset int) ([]*domain.TicketCatalog, error) {
	query := `
		SELECT id, channel_id, title, status, priority, category, assignee_id, assignee_type, created_by, created_at, updated_at
		FROM %s
		WHERE channel_id = $1
	`
	
	args := []interface{}{channelID}
	argCount := 1
	
	if status != nil {
		query += fmt.Sprintf(" AND status = $%d", argCount+1)
		args = append(args, *status)
		argCount++
	}
	
	query += " ORDER BY created_at DESC"
	
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount+1)
		args = append(args, limit)
		argCount++
	}
	
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount+1)
		args = append(args, offset)
	}
	
	rows, err := r.db.QueryContext(ctx, r.table(query), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return r.scanTickets(rows)
}

func (r TicketCatalogRepository) GetByAssigneeID(ctx context.Context, assigneeID string, status *string, limit, offset int) ([]*domain.TicketCatalog, error) {
	query := `
		SELECT id, channel_id, title, status, priority, category, assignee_id, assignee_type, created_by, created_at, updated_at
		FROM %s
		WHERE assignee_id = $1
	`
	
	args := []interface{}{assigneeID}
	argCount := 1
	
	if status != nil {
		query += fmt.Sprintf(" AND status = $%d", argCount+1)
		args = append(args, *status)
		argCount++
	}
	
	query += " ORDER BY priority DESC, created_at DESC"
	
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount+1)
		args = append(args, limit)
		argCount++
	}
	
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount+1)
		args = append(args, offset)
	}
	
	rows, err := r.db.QueryContext(ctx, r.table(query), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return r.scanTickets(rows)
}

func (r TicketCatalogRepository) Search(ctx context.Context, searchQuery string, filters map[string]interface{}, limit, offset int) ([]*domain.TicketCatalog, error) {
	query := `
		SELECT id, channel_id, title, status, priority, category, assignee_id, assignee_type, created_by, created_at, updated_at
		FROM %s
		WHERE 1=1
	`
	
	args := []interface{}{}
	argCount := 0
	
	// Add search query if provided
	if searchQuery != "" {
		argCount++
		query += fmt.Sprintf(" AND (title ILIKE $%d OR id = $%d)", argCount, argCount)
		args = append(args, "%"+searchQuery+"%")
	}
	
	// Add filters
	for key, value := range filters {
		if value == nil {
			continue
		}
		argCount++
		switch key {
		case "channel_id", "status", "priority", "category", "assignee_id", "created_by":
			query += fmt.Sprintf(" AND %s = $%d", key, argCount)
			args = append(args, value)
		}
	}
	
	query += " ORDER BY created_at DESC"
	
	if limit > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, limit)
	}
	
	if offset > 0 {
		argCount++
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, offset)
	}
	
	rows, err := r.db.QueryContext(ctx, r.table(query), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return r.scanTickets(rows)
}

func (r TicketCatalogRepository) Count(ctx context.Context, channelID string, status *string) (int, error) {
	query := `SELECT COUNT(*) FROM %s WHERE channel_id = $1`
	args := []interface{}{channelID}
	
	if status != nil {
		query += " AND status = $2"
		args = append(args, *status)
	}
	
	var count int
	err := r.db.QueryRowContext(ctx, r.table(query), args...).Scan(&count)
	
	return count, err
}

func (r TicketCatalogRepository) scanTickets(rows *sql.Rows) ([]*domain.TicketCatalog, error) {
	var tickets []*domain.TicketCatalog
	
	for rows.Next() {
		ticket := &domain.TicketCatalog{}
		err := rows.Scan(
			&ticket.ID,
			&ticket.ChannelID,
			&ticket.Title,
			&ticket.Status,
			&ticket.Priority,
			&ticket.Category,
			&ticket.AssigneeID,
			&ticket.AssigneeType,
			&ticket.CreatedBy,
			&ticket.CreatedAt,
			&ticket.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}
	
	return tickets, rows.Err()
}

func (r TicketCatalogRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}