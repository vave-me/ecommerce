package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"middleman/internal/postgres"
	"middleman/managers/internal/domain"
	"time"

	"github.com/stackus/errors"
)

type ReadMessagesRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.ReadMessagesRepository = (*ReadMessagesRepository)(nil)

func NewReadMessagesRepository(tableName string, db postgres.DB) ReadMessagesRepository {
	return ReadMessagesRepository{
		tableName: tableName,
		db:        db,
	}
}

// GetConversationMessages retrieves messages with integrated access control
func (r ReadMessagesRepository) GetConversationMessages(ctx context.Context, conversationID string, userID string, limit, offset int) ([]*domain.ReadMessage, error) {
	const query = `
		SELECT 
			m.id, m.conversation_id, m.manager_id, m.role, m.content, m.timestamp, 
			m.metadata, m.actions_taken
		FROM %s m
		INNER JOIN conversations c ON m.conversation_id = c.id
		WHERE m.conversation_id = $1 AND c.user_id = $2
		ORDER BY m.timestamp ASC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, r.table(query), conversationID, userID, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get conversation messages")
	}
	defer rows.Close()

	return r.scanMessages(rows)
}

// GetLatestMessage retrieves latest message with integrated access control
func (r ReadMessagesRepository) GetLatestMessage(ctx context.Context, conversationID string, userID string) (*domain.ReadMessage, error) {
	const query = `
		SELECT 
			m.id, m.conversation_id, m.manager_id, m.role, m.content, m.timestamp, 
			m.metadata, m.actions_taken
		FROM %s m
		INNER JOIN conversations c ON m.conversation_id = c.id
		WHERE m.conversation_id = $1 AND c.user_id = $2
		ORDER BY m.timestamp DESC
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, r.table(query), conversationID, userID)

	var msg conversationMessageDB
	err := row.Scan(
		&msg.ID, &msg.ConversationID, &msg.ManagerID, &msg.Role, &msg.Content, &msg.Timestamp,
		&msg.Metadata, &msg.ActionsTaken,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrConversationNotFound
		}
		return nil, errors.Wrap(err, "failed to get latest message")
	}

	return r.dbToConversationMessage(&msg)
}

// GetUserMessageCount optimized with single query
func (r ReadMessagesRepository) GetUserMessageCount(ctx context.Context, userID string, dateRange string) (int64, error) {
	var whereClause string

	switch dateRange {
	case "today":
		whereClause = "AND DATE(m.timestamp) = CURRENT_DATE"
	case "week":
		whereClause = "AND m.timestamp >= DATE_TRUNC('week', CURRENT_DATE)"
	case "month":
		whereClause = "AND m.timestamp >= DATE_TRUNC('month', CURRENT_DATE)"
	case "year":
		whereClause = "AND m.timestamp >= DATE_TRUNC('year', CURRENT_DATE)"
	default:
		whereClause = ""
	}

	query := `
		SELECT COUNT(*)
		FROM %s m
		INNER JOIN conversations c ON m.conversation_id = c.id
		WHERE c.user_id = $1 AND m.role = 'user' ` + whereClause

	var count int64
	err := r.db.QueryRowContext(ctx, r.table(query), userID).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get user message count")
	}

	return count, nil
}

// AddMessage with proper conflict handling
func (r ReadMessagesRepository) AddMessage(ctx context.Context, conversationID, id string, managerID string, role domain.MessageRole, content string, timestamp time.Time, metadata map[string]interface{}, actionsTaken []domain.ManagerAction) error {
	// Verify conversation exists
	const accessQuery = "SELECT 1 FROM conversations WHERE id = $1"
	var exists int
	err := r.db.QueryRowContext(ctx, accessQuery, conversationID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ErrConversationNotFound
		}
		return errors.Wrap(err, "failed to verify conversation exists")
	}

	const query = `
		INSERT INTO %s (id, conversation_id, manager_id, role, content, timestamp, metadata, actions_taken)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			manager_id = EXCLUDED.manager_id,
			role = EXCLUDED.role,
			content = EXCLUDED.content,
			timestamp = EXCLUDED.timestamp,
			metadata = EXCLUDED.metadata,
			actions_taken = EXCLUDED.actions_taken
	`

	metadataJSON, actionsJSON, err := r.prepareJSONFieldsVars(metadata, actionsTaken)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, r.table(query),
		id, conversationID, managerID, string(role), content,
		timestamp, metadataJSON, actionsJSON,
	)
	if err != nil {
		return errors.Wrap(err, "failed to add message")
	}

	return nil
}

// Database struct for scanning
type conversationMessageDB struct {
	ID             string         `db:"id"`
	ConversationID string         `db:"conversation_id"`
	ManagerID      sql.NullString `db:"manager_id"`
	Role           string         `db:"role"`
	Content        string         `db:"content"`
	Timestamp      time.Time      `db:"timestamp"`
	Metadata       sql.NullString `db:"metadata"`
	ActionsTaken   sql.NullString `db:"actions_taken"`
}

// Helper method to scan messages from rows
func (r ReadMessagesRepository) scanMessages(rows *sql.Rows) ([]*domain.ReadMessage, error) {
	var messages []*domain.ReadMessage
	for rows.Next() {
		var msg conversationMessageDB
		err := rows.Scan(
			&msg.ID, &msg.ConversationID, &msg.ManagerID, &msg.Role, &msg.Content, &msg.Timestamp,
			&msg.Metadata, &msg.ActionsTaken,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan message")
		}

		message, err := r.dbToConversationMessage(&msg)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert message")
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating message rows")
	}

	return messages, nil
}

// Helper method to convert database row to domain message
func (r ReadMessagesRepository) dbToConversationMessage(msg *conversationMessageDB) (*domain.ReadMessage, error) {
	var metadata map[string]interface{}
	if msg.Metadata.Valid && msg.Metadata.String != "" {
		if err := json.Unmarshal([]byte(msg.Metadata.String), &metadata); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal message metadata")
		}
	}

	var actionsTaken []domain.ManagerAction
	if msg.ActionsTaken.Valid && msg.ActionsTaken.String != "" {
		if err := json.Unmarshal([]byte(msg.ActionsTaken.String), &actionsTaken); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal actions taken")
		}
	}

	managerID := ""
	if msg.ManagerID.Valid {
		managerID = msg.ManagerID.String
	}

	return &domain.ReadMessage{
		ID:             msg.ID,
		ConversationID: msg.ConversationID,
		ManagerID:      managerID,
		Role:           domain.MessageRole(msg.Role),
		Content:        msg.Content,
		Timestamp:      msg.Timestamp,
		Metadata:       metadata,
		ActionsTaken:   actionsTaken,
	}, nil
}

// Helper method to prepare JSON fields for database storage
func (r ReadMessagesRepository) prepareJSONFields(message *domain.ReadMessage) (sql.NullString, sql.NullString, error) {
	var metadataJSON, actionsJSON sql.NullString

	if message.Metadata != nil {
		metadata, err := json.Marshal(message.Metadata)
		if err != nil {
			return metadataJSON, actionsJSON, errors.Wrap(err, "failed to marshal metadata")
		}
		metadataJSON = sql.NullString{String: string(metadata), Valid: true}
	}

	if message.ActionsTaken != nil && len(message.ActionsTaken) > 0 {
		actions, err := json.Marshal(message.ActionsTaken)
		if err != nil {
			return metadataJSON, actionsJSON, errors.Wrap(err, "failed to marshal actions taken")
		}
		actionsJSON = sql.NullString{String: string(actions), Valid: true}
	}

	return metadataJSON, actionsJSON, nil
}

// Helper method to prepare JSON fields for database storage
func (r ReadMessagesRepository) prepareJSONFieldsVars(msgMetadata map[string]interface{}, msgActionsTaken []domain.ManagerAction) (sql.NullString, sql.NullString, error) {
	var metadataJSON, actionsJSON sql.NullString

	if msgMetadata != nil {
		metadata, err := json.Marshal(msgMetadata)
		if err != nil {
			return metadataJSON, actionsJSON, errors.Wrap(err, "failed to marshal metadata")
		}
		metadataJSON = sql.NullString{String: string(metadata), Valid: true}
	}

	if msgActionsTaken != nil && len(msgActionsTaken) > 0 {
		actions, err := json.Marshal(msgActionsTaken)
		if err != nil {
			return metadataJSON, actionsJSON, errors.Wrap(err, "failed to marshal actions taken")
		}
		actionsJSON = sql.NullString{String: string(actions), Valid: true}
	}

	return metadataJSON, actionsJSON, nil
}

// Helper method to format table name in queries
func (r ReadMessagesRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
