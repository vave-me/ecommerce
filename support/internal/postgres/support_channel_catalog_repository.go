package postgres

import (
	"context"
	"database/sql"
	"fmt"
	
	"github.com/stackus/errors"
	"middleman/support/internal/domain"
)

type SupportChannelCatalogRepository struct {
	tableName string
	db        *sql.DB
}

var _ domain.SupportChannelCatalogRepository = (*SupportChannelCatalogRepository)(nil)

func NewSupportChannelCatalogRepository(tableName string, db *sql.DB) SupportChannelCatalogRepository {
	return SupportChannelCatalogRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r SupportChannelCatalogRepository) Add(ctx context.Context, channel *domain.SupportChannelCatalog) error {
	const query = `
		INSERT INTO %s (id, user_id, business_id, channel_type, active, open_tickets, total_tickets, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	
	_, err := r.db.ExecContext(ctx, r.table(query),
		channel.ID,
		channel.UserID,
		channel.BusinessID,
		channel.ChannelType,
		channel.Active,
		channel.OpenTickets,
		channel.TotalTickets,
		channel.CreatedAt,
		channel.UpdatedAt,
	)
	
	return err
}

func (r SupportChannelCatalogRepository) Update(ctx context.Context, channel *domain.SupportChannelCatalog) error {
	const query = `
		UPDATE %s 
		SET channel_type = $2, active = $3, open_tickets = $4, total_tickets = $5, updated_at = $6
		WHERE id = $1
	`
	
	_, err := r.db.ExecContext(ctx, r.table(query),
		channel.ID,
		channel.ChannelType,
		channel.Active,
		channel.OpenTickets,
		channel.TotalTickets,
		channel.UpdatedAt,
	)
	
	return err
}

func (r SupportChannelCatalogRepository) Delete(ctx context.Context, channelID string) error {
	const query = `DELETE FROM %s WHERE id = $1`
	
	_, err := r.db.ExecContext(ctx, r.table(query), channelID)
	
	return err
}

func (r SupportChannelCatalogRepository) Find(ctx context.Context, channelID string) (*domain.SupportChannelCatalog, error) {
	const query = `
		SELECT id, user_id, business_id, channel_type, active, open_tickets, total_tickets, created_at, updated_at
		FROM %s
		WHERE id = $1
	`
	
	channel := &domain.SupportChannelCatalog{}
	
	err := r.db.QueryRowContext(ctx, r.table(query), channelID).Scan(
		&channel.ID,
		&channel.UserID,
		&channel.BusinessID,
		&channel.ChannelType,
		&channel.Active,
		&channel.OpenTickets,
		&channel.TotalTickets,
		&channel.CreatedAt,
		&channel.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(errors.ErrNotFound, "support channel not found")
		}
		return nil, err
	}
	
	return channel, nil
}

func (r SupportChannelCatalogRepository) GetByUserID(ctx context.Context, userID string, activeOnly bool, limit, offset int) ([]*domain.SupportChannelCatalog, error) {
	query := `
		SELECT id, user_id, business_id, channel_type, active, open_tickets, total_tickets, created_at, updated_at
		FROM %s
		WHERE user_id = $1
	`
	
	args := []interface{}{userID}
	argCount := 1
	
	if activeOnly {
		query += fmt.Sprintf(" AND active = $%d", argCount+1)
		args = append(args, true)
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
	
	var channels []*domain.SupportChannelCatalog
	
	for rows.Next() {
		channel := &domain.SupportChannelCatalog{}
		err := rows.Scan(
			&channel.ID,
			&channel.UserID,
			&channel.BusinessID,
			&channel.ChannelType,
			&channel.Active,
			&channel.OpenTickets,
			&channel.TotalTickets,
			&channel.CreatedAt,
			&channel.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		channels = append(channels, channel)
	}
	
	return channels, rows.Err()
}

func (r SupportChannelCatalogRepository) GetByBusinessID(ctx context.Context, businessID string, activeOnly bool, limit, offset int) ([]*domain.SupportChannelCatalog, error) {
	query := `
		SELECT id, user_id, business_id, channel_type, active, open_tickets, total_tickets, created_at, updated_at
		FROM %s
		WHERE business_id = $1
	`
	
	args := []interface{}{businessID}
	argCount := 1
	
	if activeOnly {
		query += fmt.Sprintf(" AND active = $%d", argCount+1)
		args = append(args, true)
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
	
	var channels []*domain.SupportChannelCatalog
	
	for rows.Next() {
		channel := &domain.SupportChannelCatalog{}
		err := rows.Scan(
			&channel.ID,
			&channel.UserID,
			&channel.BusinessID,
			&channel.ChannelType,
			&channel.Active,
			&channel.OpenTickets,
			&channel.TotalTickets,
			&channel.CreatedAt,
			&channel.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		channels = append(channels, channel)
	}
	
	return channels, rows.Err()
}

func (r SupportChannelCatalogRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}