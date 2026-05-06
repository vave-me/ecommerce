package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/stackus/errors"
	"middleman/support/internal/domain"
)

type AIConfigurationRepository struct {
	tableName string
	db        *sql.DB
}

var _ domain.AIConfigurationRepository = (*AIConfigurationRepository)(nil)

func NewAIConfigurationRepository(tableName string, db *sql.DB) AIConfigurationRepository {
	return AIConfigurationRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r AIConfigurationRepository) Add(ctx context.Context, config *domain.AIConfiguration) error {
	const query = `
		INSERT INTO %s (
			id, channel_id, assistant_id, allowed_actions, knowledge_base_access,
			max_handling_tier, can_close_tickets, can_issue_refunds, confidence_threshold,
			auto_response_categories, max_tokens, temperature, active, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err := r.db.ExecContext(ctx, r.table(query),
		config.ID,
		config.ChannelID,
		config.AssistantID,
		pq.Array(config.AllowedActions),
		config.KnowledgeBaseAccess, // Stored as JSONB
		config.MaxHandlingTier,
		config.CanCloseTickets,
		config.CanIssueRefunds,
		config.ConfidenceThreshold,
		pq.Array(config.AutoResponseCategories),
		config.MaxTokens,
		config.Temperature,
		config.Active,
		config.CreatedAt,
		config.UpdatedAt,
	)

	return err
}

func (r AIConfigurationRepository) Update(ctx context.Context, config *domain.AIConfiguration) error {
	const query = `
		UPDATE %s SET
			assistant_id = $2,
			allowed_actions = $3,
			knowledge_base_access = $4,
			max_handling_tier = $5,
			can_close_tickets = $6,
			can_issue_refunds = $7,
			confidence_threshold = $8,
			auto_response_categories = $9,
			max_tokens = $10,
			temperature = $11,
			active = $12,
			updated_at = $13
		WHERE channel_id = $1
	`

	_, err := r.db.ExecContext(ctx, r.table(query),
		config.ChannelID,
		config.AssistantID,
		pq.Array(config.AllowedActions),
		config.KnowledgeBaseAccess,
		config.MaxHandlingTier,
		config.CanCloseTickets,
		config.CanIssueRefunds,
		config.ConfidenceThreshold,
		pq.Array(config.AutoResponseCategories),
		config.MaxTokens,
		config.Temperature,
		config.Active,
		time.Now(),
	)

	return err
}

func (r AIConfigurationRepository) GetByChannelID(ctx context.Context, channelID string) (*domain.AIConfiguration, error) {
	const query = `
		SELECT
			id, channel_id, assistant_id, allowed_actions, knowledge_base_access,
			max_handling_tier, can_close_tickets, can_issue_refunds, confidence_threshold,
			auto_response_categories, max_tokens, temperature, active, created_at, updated_at
		FROM %s
		WHERE channel_id = $1
	`

	config := &domain.AIConfiguration{}
	var kbAccess sql.NullString

	err := r.db.QueryRowContext(ctx, r.table(query), channelID).Scan(
		&config.ID,
		&config.ChannelID,
		&config.AssistantID,
		pq.Array(&config.AllowedActions),
		&kbAccess,
		&config.MaxHandlingTier,
		&config.CanCloseTickets,
		&config.CanIssueRefunds,
		&config.ConfidenceThreshold,
		pq.Array(&config.AutoResponseCategories),
		&config.MaxTokens,
		&config.Temperature,
		&config.Active,
		&config.CreatedAt,
		&config.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(errors.ErrNotFound, "AI configuration not found")
		}
		return nil, err
	}

	// Parse JSONB knowledge base access
	if kbAccess.Valid {
		// Parse JSON into map[string]string
		// This would need proper JSON unmarshaling in production
		config.KnowledgeBaseAccess = make(map[string]string)
	}

	return config, nil
}

func (r AIConfigurationRepository) Delete(ctx context.Context, channelID string) error {
	const query = `DELETE FROM %s WHERE channel_id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), channelID)

	return err
}

func (r AIConfigurationRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}