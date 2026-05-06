package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"middleman/assistants/internal/domain"
	"middleman/internal/postgres"
	"strings"
	"time"

	"github.com/stackus/errors"
)

type ReadConversationRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.ReadConversationRepository = (*ReadConversationRepository)(nil)

func NewReadConversationRepository(tableName string, db postgres.DB) ReadConversationRepository {
	return ReadConversationRepository{
		tableName: tableName,
		db:        db,
	}
}

// GetConversation retrieves a conversation by ID with user access check
func (r ReadConversationRepository) GetConversation(ctx context.Context, id string, userID string) (*domain.ReadConversation, error) {
	const query = `
		SELECT 
			id, user_id, assistant_id, created_at, updated_at, 
			active, context
		FROM %s
		WHERE id = $1 AND user_id = $2
	`

	row := r.db.QueryRowContext(ctx, r.table(query), id, userID)

	var c conversationDB
	err := row.Scan(
		&c.ID, &c.UserID, &c.AssistantID, &c.CreatedAt, &c.UpdatedAt,
		&c.Active, &c.Context,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrConversationNotFound
		}
		return nil, errors.Wrap(err, "failed to get conversation")
	}

	return r.dbToConversation(&c)
}

// GetUserConversations retrieves conversations for a user with pagination
func (r ReadConversationRepository) GetUserConversations(ctx context.Context, userID string, activeOnly bool, limit, offset int) ([]*domain.ReadConversation, int, error) {
	var whereClause []string
	var args []interface{}
	argIndex := 1

	whereClause = append(whereClause, fmt.Sprintf("user_id = $%d", argIndex))
	args = append(args, userID)
	argIndex++

	if activeOnly {
		whereClause = append(whereClause, fmt.Sprintf("active = $%d", argIndex))
		args = append(args, true)
		argIndex++
	}

	whereSQL := strings.Join(whereClause, " AND ")

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", r.tableName, whereSQL)
	var totalCount int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get conversation count")
	}

	// Get paginated results
	query := fmt.Sprintf(`
		SELECT 
			id, user_id, assistant_id, created_at, updated_at,
			active, context
		FROM %s
		WHERE %s
		ORDER BY updated_at DESC, created_at DESC
		LIMIT $%d OFFSET $%d
	`, r.tableName, whereSQL, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get user conversations")
	}
	defer rows.Close()

	conversations, err := r.scanConversations(rows)
	if err != nil {
		return nil, 0, err
	}

	return conversations, totalCount, nil
}

// GetAssistantConversations retrieves conversations for a specific assistant
func (r ReadConversationRepository) GetAssistantConversations(ctx context.Context, assistantID string, limit, offset int) ([]*domain.ReadConversation, int, error) {
	// Get total count
	const countQuery = "SELECT COUNT(*) FROM %s WHERE assistant_id = $1"
	var totalCount int
	err := r.db.QueryRowContext(ctx, r.table(countQuery), assistantID).Scan(&totalCount)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get assistant conversation count")
	}

	// Get paginated results
	const query = `
		SELECT 
			id, user_id, assistant_id, created_at, updated_at,
			active, context
		FROM %s
		WHERE assistant_id = $1
		ORDER BY updated_at DESC, created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, r.table(query), assistantID, limit, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get assistant conversations")
	}
	defer rows.Close()

	conversations, err := r.scanConversations(rows)
	if err != nil {
		return nil, 0, err
	}

	return conversations, totalCount, nil
}

// GetRecentConversations retrieves the most recent conversations for a user
func (r ReadConversationRepository) GetRecentConversations(ctx context.Context, userID string, limit int) ([]*domain.ReadConversation, error) {
	const query = `
		SELECT 
			id, user_id, assistant_id, created_at, updated_at,
			active, context
		FROM %s
		WHERE user_id = $1 AND active = true
		ORDER BY updated_at DESC, created_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, r.table(query), userID, limit)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get recent conversations")
	}
	defer rows.Close()

	return r.scanConversations(rows)
}

// GetConversationsByDateRange retrieves conversations within a date range
func (r ReadConversationRepository) GetConversationsByDateRange(ctx context.Context, userID string, startDate, endDate time.Time) ([]*domain.ReadConversation, error) {
	const query = `
		SELECT 
			id, user_id, assistant_id, created_at, updated_at,
			active, context
		FROM %s
		WHERE user_id = $1 AND created_at BETWEEN $2 AND $3
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, r.table(query), userID, startDate, endDate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get conversations by date range")
	}
	defer rows.Close()

	return r.scanConversations(rows)
}

// Database struct for scanning
type conversationDB struct {
	ID          string         `db:"id"`
	UserID      string         `db:"user_id"`
	AssistantID string         `db:"assistant_id"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
	Active      bool           `db:"active"`
	Context     sql.NullString `db:"context"`
}

// Helper method to scan conversations from rows
func (r ReadConversationRepository) scanConversations(rows *sql.Rows) ([]*domain.ReadConversation, error) {
	var conversations []*domain.ReadConversation
	for rows.Next() {
		var c conversationDB
		err := rows.Scan(
			&c.ID, &c.UserID, &c.AssistantID, &c.CreatedAt, &c.UpdatedAt,
			&c.Active, &c.Context,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan conversation")
		}

		conversation, err := r.dbToConversation(&c)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert conversation")
		}

		conversations = append(conversations, conversation)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating conversation rows")
	}

	return conversations, nil
}

// Helper method to convert database row to domain conversation
func (r ReadConversationRepository) dbToConversation(c *conversationDB) (*domain.ReadConversation, error) {
	var context map[string]interface{}
	if c.Context.Valid && c.Context.String != "" {
		if err := json.Unmarshal([]byte(c.Context.String), &context); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal context")
		}
	}

	conversation := &domain.ReadConversation{
		ID:          c.ID,
		UserID:      c.UserID,
		AssistantID: c.AssistantID,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		Active:      c.Active,
		Context:     context,
		Messages:    []domain.ConversationMessage{}, // Initialize empty slice
	}

	// Set the aggregate ID

	return conversation, nil
}

// Helper method to format table name in queries
func (r ReadConversationRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}

// Write operations for read model maintenance

// AddConversation creates a new conversation in the read model
func (r ReadConversationRepository) AddConversation(ctx context.Context, id, userID, assistantID string, createdAt time.Time, context map[string]interface{}) error {
	const query = `
		INSERT INTO %s (id, user_id, assistant_id, created_at, updated_at, active, context)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			user_id = EXCLUDED.user_id,
			assistant_id = EXCLUDED.assistant_id,
			created_at = EXCLUDED.created_at,
			updated_at = EXCLUDED.updated_at,
			active = EXCLUDED.active,
			context = EXCLUDED.context
	`

	contextJSON, err := r.prepareJSONFields(context)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, r.table(query),
		id, userID, assistantID, createdAt, createdAt, true, contextJSON,
	)
	if err != nil {
		return errors.Wrap(err, "failed to add conversation to read model")
	}

	return nil
}

// UpdateConversationContext updates the conversation context
func (r ReadConversationRepository) UpdateConversationContext(ctx context.Context, id string, context map[string]interface{}, updatedAt time.Time) error {
	const query = `
		UPDATE %s SET
			context = $2,
			updated_at = $3
		WHERE id = $1
	`

	contextJSON, err := r.prepareJSONFields(context)
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(ctx, r.table(query), id, contextJSON, updatedAt)
	if err != nil {
		return errors.Wrap(err, "failed to update conversation context")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.Wrap(domain.ErrConversationNotFound, "conversation not found for context update")
	}

	return nil
}

// ArchiveConversation marks a conversation as inactive
func (r ReadConversationRepository) ArchiveConversation(ctx context.Context, id string, archivedAt time.Time) error {
	const query = `
		UPDATE %s SET
			active = false,
			updated_at = $2
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, r.table(query), id, archivedAt)
	if err != nil {
		return errors.Wrap(err, "failed to archive conversation")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.Wrap(domain.ErrConversationNotFound, "conversation not found for archival")
	}

	return nil
}

// UpdateConversationTimestamp updates the conversation's updated_at timestamp
func (r ReadConversationRepository) UpdateConversationTimestamp(ctx context.Context, id string, updatedAt time.Time) error {
	const query = `
		UPDATE %s SET
			updated_at = $2
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, r.table(query), id, updatedAt)
	if err != nil {
		return errors.Wrap(err, "failed to update conversation timestamp")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.Wrap(domain.ErrConversationNotFound, "conversation not found for timestamp update")
	}

	return nil
}

// Helper method to prepare JSON fields for database storage
func (r ReadConversationRepository) prepareJSONFields(context map[string]interface{}) ([]byte, error) {
	if context == nil {
		return nil, nil
	}
	return json.Marshal(context)
}
