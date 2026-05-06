package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/stackus/errors"
	"middleman/internal/postgres"
	"middleman/managers/internal/domain"
)

type LLMJournalRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.LLMJournalRepository = (*LLMJournalRepository)(nil)

func NewLLMJournalRepository(tableName string, db postgres.DB) LLMJournalRepository {
	return LLMJournalRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r LLMJournalRepository) Save(entry *domain.LLMJournalEntry) error {
	const query = `
		INSERT INTO %s (
			id, manager_id, conversation_id, message_id, user_id,
			request_type, request_content, request_context,
			response_content, response_metadata, tool_calls,
			detected_patterns, learning_insights, confidence_score,
			processing_time_ms, tokens_used, model_used, provider,
			created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
		)`

	// Convert complex types to JSON
	requestContextJSON, _ := json.Marshal(entry.RequestContext)
	responseMetadataJSON, _ := json.Marshal(entry.ResponseMetadata)
	toolCallsJSON, _ := json.Marshal(entry.ToolCalls)
	patternsJSON, _ := json.Marshal(entry.DetectedPatterns)
	insightsJSON, _ := json.Marshal(entry.LearningInsights)

	_, err := r.db.ExecContext(
		context.Background(),
		r.table(query),
		entry.ID,
		entry.ManagerID,
		entry.ConversationID,
		entry.MessageID,
		entry.UserID,
		entry.RequestType,
		entry.RequestContent,
		requestContextJSON,
		entry.ResponseContent,
		responseMetadataJSON,
		toolCallsJSON,
		patternsJSON,
		insightsJSON,
		entry.ConfidenceScore,
		entry.ProcessingTimeMs,
		entry.TokensUsed,
		entry.ModelUsed,
		entry.Provider,
		entry.CreatedAt,
	)

	return err
}

func (r LLMJournalRepository) FindByID(id string) (*domain.LLMJournalEntry, error) {
	const query = `
		SELECT 
			id, manager_id, conversation_id, message_id, user_id,
			request_type, request_content, request_context,
			response_content, response_metadata, tool_calls,
			detected_patterns, learning_insights, confidence_score,
			processing_time_ms, tokens_used, model_used, provider,
			created_at
		FROM %s WHERE id = $1`

	entry := &domain.LLMJournalEntry{}
	var (
		requestContextJSON   []byte
		responseMetadataJSON []byte
		toolCallsJSON        []byte
		patternsJSON         []byte
		insightsJSON         []byte
	)

	err := r.db.QueryRowContext(
		context.Background(),
		r.table(query),
		id,
	).Scan(
		&entry.ID,
		&entry.ManagerID,
		&entry.ConversationID,
		&entry.MessageID,
		&entry.UserID,
		&entry.RequestType,
		&entry.RequestContent,
		&requestContextJSON,
		&entry.ResponseContent,
		&responseMetadataJSON,
		&toolCallsJSON,
		&patternsJSON,
		&insightsJSON,
		&entry.ConfidenceScore,
		&entry.ProcessingTimeMs,
		&entry.TokensUsed,
		&entry.ModelUsed,
		&entry.Provider,
		&entry.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(errors.ErrNotFound, "journal entry not found")
		}
		return nil, err
	}

	// Unmarshal JSON fields
	json.Unmarshal(requestContextJSON, &entry.RequestContext)
	json.Unmarshal(responseMetadataJSON, &entry.ResponseMetadata)
	json.Unmarshal(toolCallsJSON, &entry.ToolCalls)
	json.Unmarshal(patternsJSON, &entry.DetectedPatterns)
	json.Unmarshal(insightsJSON, &entry.LearningInsights)

	return entry, nil
}

func (r LLMJournalRepository) FindByManagerID(managerID string, limit int, offset int) ([]*domain.LLMJournalEntry, error) {
	const query = `
		SELECT 
			id, manager_id, conversation_id, message_id, user_id,
			request_type, request_content, request_context,
			response_content, response_metadata, tool_calls,
			detected_patterns, learning_insights, confidence_score,
			processing_time_ms, tokens_used, model_used, provider,
			created_at
		FROM %s 
		WHERE manager_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	return r.queryEntries(query, managerID, limit, offset)
}

func (r LLMJournalRepository) FindByUserID(userID string, limit int, offset int) ([]*domain.LLMJournalEntry, error) {
	const query = `
		SELECT 
			id, manager_id, conversation_id, message_id, user_id,
			request_type, request_content, request_context,
			response_content, response_metadata, tool_calls,
			detected_patterns, learning_insights, confidence_score,
			processing_time_ms, tokens_used, model_used, provider,
			created_at
		FROM %s 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	return r.queryEntries(query, userID, limit, offset)
}

func (r LLMJournalRepository) FindByConversationID(conversationID string) ([]*domain.LLMJournalEntry, error) {
	const query = `
		SELECT 
			id, manager_id, conversation_id, message_id, user_id,
			request_type, request_content, request_context,
			response_content, response_metadata, tool_calls,
			detected_patterns, learning_insights, confidence_score,
			processing_time_ms, tokens_used, model_used, provider,
			created_at
		FROM %s 
		WHERE conversation_id = $1
		ORDER BY created_at ASC`

	rows, err := r.db.QueryContext(
		context.Background(),
		r.table(query),
		conversationID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanEntries(rows)
}

func (r LLMJournalRepository) FindPatterns(patternType string, since time.Time) ([]*domain.LLMJournalEntry, error) {
	const query = `
		SELECT 
			id, manager_id, conversation_id, message_id, user_id,
			request_type, request_content, request_context,
			response_content, response_metadata, tool_calls,
			detected_patterns, learning_insights, confidence_score,
			processing_time_ms, tokens_used, model_used, provider,
			created_at
		FROM %s 
		WHERE created_at >= $1
		AND detected_patterns @> $2
		ORDER BY created_at DESC`

	patternFilter, _ := json.Marshal([]map[string]interface{}{
		{"pattern_type": patternType},
	})

	rows, err := r.db.QueryContext(
		context.Background(),
		r.table(query),
		since,
		patternFilter,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanEntries(rows)
}

func (r LLMJournalRepository) GetInsights(userID string, insightType string) ([]domain.LearningInsight, error) {
	const query = `
		SELECT DISTINCT learning_insights
		FROM %s 
		WHERE user_id = $1
		AND learning_insights IS NOT NULL`

	rows, err := r.db.QueryContext(
		context.Background(),
		r.table(query),
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allInsights []domain.LearningInsight
	for rows.Next() {
		var insightsJSON []byte
		if err := rows.Scan(&insightsJSON); err != nil {
			continue
		}

		var insights []domain.LearningInsight
		if err := json.Unmarshal(insightsJSON, &insights); err != nil {
			continue
		}

		// Filter by insight type if specified
		for _, insight := range insights {
			if insightType == "" || insight.InsightType == insightType {
				allInsights = append(allInsights, insight)
			}
		}
	}

	return allInsights, nil
}

func (r LLMJournalRepository) GetPerformanceMetrics(managerID string, since time.Time) (map[string]interface{}, error) {
	const query = `
		SELECT 
			COUNT(*) as total_responses,
			AVG(processing_time_ms) as avg_processing_time,
			AVG(tokens_used) as avg_tokens,
			AVG(confidence_score) as avg_confidence,
			MAX(processing_time_ms) as max_processing_time,
			MIN(processing_time_ms) as min_processing_time,
			COUNT(DISTINCT user_id) as unique_users,
			COUNT(DISTINCT conversation_id) as total_conversations
		FROM %s 
		WHERE manager_id = $1 AND created_at >= $2`

	var metrics struct {
		TotalResponses      int64
		AvgProcessingTime   sql.NullFloat64
		AvgTokens           sql.NullFloat64
		AvgConfidence       sql.NullFloat64
		MaxProcessingTime   sql.NullInt64
		MinProcessingTime   sql.NullInt64
		UniqueUsers         int64
		TotalConversations  int64
	}

	err := r.db.QueryRowContext(
		context.Background(),
		r.table(query),
		managerID,
		since,
	).Scan(
		&metrics.TotalResponses,
		&metrics.AvgProcessingTime,
		&metrics.AvgTokens,
		&metrics.AvgConfidence,
		&metrics.MaxProcessingTime,
		&metrics.MinProcessingTime,
		&metrics.UniqueUsers,
		&metrics.TotalConversations,
	)

	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"total_responses":      metrics.TotalResponses,
		"unique_users":         metrics.UniqueUsers,
		"total_conversations":  metrics.TotalConversations,
	}

	if metrics.AvgProcessingTime.Valid {
		result["avg_processing_time_ms"] = metrics.AvgProcessingTime.Float64
	}
	if metrics.AvgTokens.Valid {
		result["avg_tokens_used"] = metrics.AvgTokens.Float64
	}
	if metrics.AvgConfidence.Valid {
		result["avg_confidence_score"] = metrics.AvgConfidence.Float64
	}
	if metrics.MaxProcessingTime.Valid {
		result["max_processing_time_ms"] = metrics.MaxProcessingTime.Int64
	}
	if metrics.MinProcessingTime.Valid {
		result["min_processing_time_ms"] = metrics.MinProcessingTime.Int64
	}

	return result, nil
}

// Helper methods

func (r LLMJournalRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}

func (r LLMJournalRepository) queryEntries(query string, args ...interface{}) ([]*domain.LLMJournalEntry, error) {
	rows, err := r.db.QueryContext(
		context.Background(),
		r.table(query),
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanEntries(rows)
}

func (r LLMJournalRepository) scanEntries(rows *sql.Rows) ([]*domain.LLMJournalEntry, error) {
	var entries []*domain.LLMJournalEntry

	for rows.Next() {
		entry := &domain.LLMJournalEntry{}
		var (
			requestContextJSON   []byte
			responseMetadataJSON []byte
			toolCallsJSON        []byte
			patternsJSON         []byte
			insightsJSON         []byte
		)

		err := rows.Scan(
			&entry.ID,
			&entry.ManagerID,
			&entry.ConversationID,
			&entry.MessageID,
			&entry.UserID,
			&entry.RequestType,
			&entry.RequestContent,
			&requestContextJSON,
			&entry.ResponseContent,
			&responseMetadataJSON,
			&toolCallsJSON,
			&patternsJSON,
			&insightsJSON,
			&entry.ConfidenceScore,
			&entry.ProcessingTimeMs,
			&entry.TokensUsed,
			&entry.ModelUsed,
			&entry.Provider,
			&entry.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal JSON fields
		json.Unmarshal(requestContextJSON, &entry.RequestContext)
		json.Unmarshal(responseMetadataJSON, &entry.ResponseMetadata)
		json.Unmarshal(toolCallsJSON, &entry.ToolCalls)
		json.Unmarshal(patternsJSON, &entry.DetectedPatterns)
		json.Unmarshal(insightsJSON, &entry.LearningInsights)

		entries = append(entries, entry)
	}

	return entries, rows.Err()
}