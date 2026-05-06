package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"middleman/assistants/internal/domain"
	"middleman/internal/postgres"

	"github.com/lib/pq"
	"github.com/stackus/errors"
)

type CatalogRepository struct {
	tableName string
	db        postgres.DB
}

func (r CatalogRepository) Find(ctx context.Context, id string) (*domain.CatalogAssistant, error) {
	const query = `
		SELECT id, assistant_name, description, user_id, type, capabilities, enabled, 
		       temperature, max_tokens, system_prompt, created_at, updated_at
		FROM %s
		WHERE id = $1
	`

	var assistant domain.CatalogAssistant
	var capabilityStrings pq.StringArray
	var name sql.NullString
	var userID sql.NullString
	var assistantType string

	err := r.db.QueryRowContext(ctx, r.table(query), id).Scan(
		&assistant.ID,
		&name,
		&assistant.Description,
		&userID,
		&assistantType,
		&capabilityStrings,
		&assistant.Active,
		&assistant.Temperature,
		&assistant.MaxTokens,
		&assistant.SystemPrompt,
		&assistant.CreatedAt,
		&assistant.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(errors.ErrNotFound, "assistant not found")
		}
		return nil, errors.Wrap(err, "failed to find assistant")
	}

	// Handle nullable name field
	if name.Valid {
		assistant.Name = name.String
	}

	// Handle nullable userID field
	if userID.Valid {
		assistant.UserID = userID.String
	}

	// Convert string to AssistantType
	assistant.Type = domain.AssistantType(assistantType)

	// Convert string array to capabilities
	assistant.Capabilities = make([]domain.AssistantCapability, len(capabilityStrings))
	for i, capStr := range capabilityStrings {
		assistant.Capabilities[i] = domain.AssistantCapability(capStr)
	}

	return &assistant, nil
}

var _ domain.CatalogRepository = (*CatalogRepository)(nil)

func NewCatalogRepository(tableName string, db postgres.DB) CatalogRepository {
	return CatalogRepository{
		tableName: tableName,
		db:        db,
	}
}

// Add creates a new assistant catalog entry
func (r CatalogRepository) Add(
	ctx context.Context,
	id string,
	name string,
	description string,
	userID string,
	assistantType domain.AssistantType,
	capabilities []domain.AssistantCapability,
	active bool,
	temperature float64,
	maxTokens int,
	systemPrompt string,
) error {
	const query = `
        INSERT INTO %s (
            id,
            assistant_name,
            description,
            user_id,
            type,
            capabilities,
            enabled,
            temperature,
            max_tokens,
            system_prompt,
            created_at,
            updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
    `

	// Convert capabilities to string array for PostgreSQL
	capabilityStrings := make([]string, len(capabilities))
	for i, cap := range capabilities {
		capabilityStrings[i] = string(cap)
	}

	now := time.Now()

	_, err := r.db.ExecContext(ctx, r.table(query),
		id,
		name,
		description,
		userID,
		string(assistantType),
		pq.Array(capabilityStrings),
		active,
		temperature,
		maxTokens,
		systemPrompt,
		now,
		now,
	)
	if err != nil {
		log.Printf("Error inserting assistant catalog entry: %v", err)
		return errors.Wrap(err, "failed to add assistant to catalog")
	}

	return nil
}

// Remove deletes an assistant from the catalog
func (r CatalogRepository) Remove(ctx context.Context, id string) error {
	const query = "DELETE FROM %s WHERE id = $1"

	result, err := r.db.ExecContext(ctx, r.table(query), id)
	if err != nil {
		return errors.Wrap(err, "failed to remove assistant from catalog")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.Wrap(errors.ErrNotFound, "assistant not found in catalog")
	}

	return nil
}

// Update modifies an existing assistant in the catalog
func (r CatalogRepository) Update(ctx context.Context, assistant *domain.CatalogAssistant) error {
	const query = `
        UPDATE %s SET
            assistant_name = $2,
            description = $3,
            user_id = $4,
            type = $5,
            capabilities = $6,
            enabled = $7,
            temperature = $8,
            max_tokens = $9,
            system_prompt = $10,
            updated_at = $11
        WHERE id = $1
    `

	// Convert capabilities to string array for PostgreSQL
	capabilityStrings := make([]string, len(assistant.Capabilities))
	for i, cap := range assistant.Capabilities {
		capabilityStrings[i] = string(cap)
	}

	result, err := r.db.ExecContext(ctx, r.table(query),
		assistant.ID,
		assistant.Name,
		assistant.Description,
		assistant.UserID,
		string(assistant.Type),
		pq.Array(capabilityStrings),
		assistant.Active,
		assistant.Temperature,
		assistant.MaxTokens,
		assistant.SystemPrompt,
		time.Now(),
	)
	if err != nil {
		return errors.Wrap(err, "failed to update assistant in catalog")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.Wrap(errors.ErrNotFound, "assistant not found for update")
	}

	return nil
}

// UpdateActiveStatus updates only the active status and timestamp
func (r CatalogRepository) UpdateActiveStatus(ctx context.Context, id string, active bool, updatedAt time.Time) error {
	const query = `
        UPDATE %s SET
            enabled = $2,
            updated_at = $3
        WHERE id = $1
    `

	result, err := r.db.ExecContext(ctx, r.table(query),
		id,
		active,
		updatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "failed to update assistant active status")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.Wrap(errors.ErrNotFound, "assistant not found for active status update")
	}

	return nil
}

// UpdateConfiguration updates only the configuration fields and timestamp
func (r CatalogRepository) UpdateConfiguration(ctx context.Context, id string, temperature float64, maxTokens int, systemPrompt string, updatedAt time.Time) error {
	const query = `
        UPDATE %s SET
            temperature = $2,
            max_tokens = $3,
            system_prompt = $4,
            updated_at = $5
        WHERE id = $1
    `

	result, err := r.db.ExecContext(ctx, r.table(query),
		id,
		temperature,
		maxTokens,
		systemPrompt,
		updatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "failed to update assistant configuration")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.Wrap(errors.ErrNotFound, "assistant not found for configuration update")
	}

	return nil
}

// UpdateConfigurationWithCapabilities updates configuration including capabilities
func (r CatalogRepository) UpdateConfigurationWithCapabilities(ctx context.Context, id string, temperature float64, maxTokens int, systemPrompt string, capabilities []domain.AssistantCapability, updatedAt time.Time) error {
	const query = `
        UPDATE %s SET
            temperature = $2,
            max_tokens = $3,
            system_prompt = $4,
            capabilities = $5,
            updated_at = $6
        WHERE id = $1
    `

	// Convert capabilities to string array
	capabilityStrings := make([]string, len(capabilities))
	for i, cap := range capabilities {
		capabilityStrings[i] = string(cap)
	}

	result, err := r.db.ExecContext(ctx, r.table(query),
		id,
		temperature,
		maxTokens,
		systemPrompt,
		pq.StringArray(capabilityStrings),
		updatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "failed to update assistant configuration with capabilities")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.Wrap(errors.ErrNotFound, "assistant not found for configuration update")
	}

	return nil
}

// FindAll retrieves all assistants from the catalog
func (r CatalogRepository) FindAll(ctx context.Context, userId string) ([]*domain.CatalogAssistant, error) {
	const query = `
        SELECT id, assistant_name, description, user_id, type, capabilities, enabled, 
               temperature, max_tokens, system_prompt, created_at, updated_at
        FROM %s
        ORDER BY created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, r.table(query))
	if err != nil {
		return nil, errors.Wrap(err, "failed to query all assistants")
	}
	defer rows.Close()

	var assistants []*domain.CatalogAssistant
	for rows.Next() {
		assistant, err := r.scanAssistant(rows)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan assistant")
		}
		assistants = append(assistants, assistant)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating assistant rows")
	}

	return assistants, nil
}

// FindActiveByUser retrieves only active assistants for a specific user from the catalog
func (r CatalogRepository) FindActiveByUser(ctx context.Context, userID string) ([]*domain.CatalogAssistant, error) {
	const query = `
        SELECT id, assistant_name, description, user_id, type, capabilities, enabled, 
               temperature, max_tokens, system_prompt, created_at, updated_at
        FROM %s
        WHERE enabled = true AND user_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, r.table(query), userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query active assistants for user")
	}
	defer rows.Close()

	var assistants []*domain.CatalogAssistant
	for rows.Next() {
		assistant, err := r.scanAssistant(rows)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan active assistant for user")
		}
		assistants = append(assistants, assistant)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating active assistant rows for user")
	}

	return assistants, nil
}

// Helper method to scan a row into a CatalogAssistant
func (r CatalogRepository) scanAssistant(rows *sql.Rows) (*domain.CatalogAssistant, error) {
	var assistant domain.CatalogAssistant
	var capabilityStrings pq.StringArray
	var name sql.NullString
	var userID sql.NullString
	var assistantType string

	err := rows.Scan(
		&assistant.ID,
		&name,
		&assistant.Description,
		&userID,
		&assistantType,
		&capabilityStrings,
		&assistant.Active,
		&assistant.Temperature,
		&assistant.MaxTokens,
		&assistant.SystemPrompt,
		&assistant.CreatedAt,
		&assistant.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Handle nullable name field
	if name.Valid {
		assistant.Name = name.String
	}

	// Handle nullable userID field
	if userID.Valid {
		assistant.UserID = userID.String
	}

	// Convert string to AssistantType
	assistant.Type = domain.AssistantType(assistantType)

	// Convert string array to capabilities
	assistant.Capabilities = make([]domain.AssistantCapability, len(capabilityStrings))
	for i, capStr := range capabilityStrings {
		assistant.Capabilities[i] = domain.AssistantCapability(capStr)
	}

	return &assistant, nil
}

// Helper method to format table name in queries
func (r CatalogRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
