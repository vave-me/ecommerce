package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"middleman/internal/postgres"
	"middleman/managers/internal/domain"

	"github.com/lib/pq"
	"github.com/stackus/errors"
)

type CatalogRepository struct {
	tableName string
	db        postgres.DB
}

func (r CatalogRepository) Find(ctx context.Context, id string) (*domain.CatalogManager, error) {
	const query = `
		SELECT id, manager_name, description, user_id, type, capabilities, enabled, 
		       temperature, max_tokens, system_prompt, created_at, updated_at
		FROM %s
		WHERE id = $1
	`

	var manager domain.CatalogManager
	var capabilityStrings pq.StringArray
	var name sql.NullString
	var userID sql.NullString
	var managerType string

	err := r.db.QueryRowContext(ctx, r.table(query), id).Scan(
		&manager.ID,
		&name,
		&manager.Description,
		&userID,
		&managerType,
		&capabilityStrings,
		&manager.Active,
		&manager.Temperature,
		&manager.MaxTokens,
		&manager.SystemPrompt,
		&manager.CreatedAt,
		&manager.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(errors.ErrNotFound, "manager not found")
		}
		return nil, errors.Wrap(err, "failed to find manager")
	}

	// Handle nullable name field
	if name.Valid {
		manager.Name = name.String
	}

	// Handle nullable userID field
	if userID.Valid {
		manager.UserID = userID.String
	}

	// Convert string to ManagerType
	manager.Type = domain.ManagerType(managerType)

	// Convert string array to capabilities
	manager.Capabilities = make([]domain.ManagerCapability, len(capabilityStrings))
	for i, capStr := range capabilityStrings {
		manager.Capabilities[i] = domain.ManagerCapability(capStr)
	}

	return &manager, nil
}

var _ domain.CatalogRepository = (*CatalogRepository)(nil)

func NewCatalogRepository(tableName string, db postgres.DB) CatalogRepository {
	return CatalogRepository{
		tableName: tableName,
		db:        db,
	}
}

// Add creates a new manager catalog entry
func (r CatalogRepository) Add(
	ctx context.Context,
	id string,
	name string,
	description string,
	userID string,
	managerType domain.ManagerType,
	capabilities []domain.ManagerCapability,
	active bool,
	temperature float64,
	maxTokens int,
	systemPrompt string,
) error {
	const query = `
        INSERT INTO %s (
            id,
            manager_name,
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
		string(managerType),
		pq.Array(capabilityStrings),
		active,
		temperature,
		maxTokens,
		systemPrompt,
		now,
		now,
	)
	if err != nil {
		log.Printf("Error inserting manager catalog entry: %v", err)
		return errors.Wrap(err, "failed to add manager to catalog")
	}

	return nil
}

// Remove deletes an manager from the catalog
func (r CatalogRepository) Remove(ctx context.Context, id string) error {
	const query = "DELETE FROM %s WHERE id = $1"

	result, err := r.db.ExecContext(ctx, r.table(query), id)
	if err != nil {
		return errors.Wrap(err, "failed to remove manager from catalog")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.Wrap(errors.ErrNotFound, "manager not found in catalog")
	}

	return nil
}

// Update modifies an existing manager in the catalog
func (r CatalogRepository) Update(ctx context.Context, manager *domain.CatalogManager) error {
	const query = `
        UPDATE %s SET
            manager_name = $2,
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
	capabilityStrings := make([]string, len(manager.Capabilities))
	for i, cap := range manager.Capabilities {
		capabilityStrings[i] = string(cap)
	}

	result, err := r.db.ExecContext(ctx, r.table(query),
		manager.ID,
		manager.Name,
		manager.Description,
		manager.UserID,
		string(manager.Type),
		pq.Array(capabilityStrings),
		manager.Active,
		manager.Temperature,
		manager.MaxTokens,
		manager.SystemPrompt,
		time.Now(),
	)
	if err != nil {
		return errors.Wrap(err, "failed to update manager in catalog")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.Wrap(errors.ErrNotFound, "manager not found for update")
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
		return errors.Wrap(err, "failed to update manager active status")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.Wrap(errors.ErrNotFound, "manager not found for active status update")
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
		return errors.Wrap(err, "failed to update manager configuration")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.Wrap(errors.ErrNotFound, "manager not found for configuration update")
	}

	return nil
}

// UpdateConfigurationWithCapabilities updates configuration including capabilities
func (r CatalogRepository) UpdateConfigurationWithCapabilities(ctx context.Context, id string, temperature float64, maxTokens int, systemPrompt string, capabilities []domain.ManagerCapability, updatedAt time.Time) error {
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
		return errors.Wrap(err, "failed to update manager configuration with capabilities")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.Wrap(errors.ErrNotFound, "manager not found for configuration update")
	}

	return nil
}

// FindAll retrieves all managers from the catalog
func (r CatalogRepository) FindAll(ctx context.Context, userId string) ([]*domain.CatalogManager, error) {
	const query = `
        SELECT id, manager_name, description, user_id, type, capabilities, enabled, 
               temperature, max_tokens, system_prompt, created_at, updated_at
        FROM %s
        ORDER BY created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, r.table(query))
	if err != nil {
		return nil, errors.Wrap(err, "failed to query all managers")
	}
	defer rows.Close()

	var managers []*domain.CatalogManager
	for rows.Next() {
		manager, err := r.scanManager(rows)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan manager")
		}
		managers = append(managers, manager)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating manager rows")
	}

	return managers, nil
}

// FindActiveByUser retrieves only active managers for a specific user from the catalog
func (r CatalogRepository) FindActiveByUser(ctx context.Context, userID string) ([]*domain.CatalogManager, error) {
	const query = `
        SELECT id, manager_name, description, user_id, type, capabilities, enabled, 
               temperature, max_tokens, system_prompt, created_at, updated_at
        FROM %s
        WHERE enabled = true AND user_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, r.table(query), userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query active managers for user")
	}
	defer rows.Close()

	var managers []*domain.CatalogManager
	for rows.Next() {
		manager, err := r.scanManager(rows)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan active manager for user")
		}
		managers = append(managers, manager)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating active manager rows for user")
	}

	return managers, nil
}

// Helper method to scan a row into a CatalogManager
func (r CatalogRepository) scanManager(rows *sql.Rows) (*domain.CatalogManager, error) {
	var manager domain.CatalogManager
	var capabilityStrings pq.StringArray
	var name sql.NullString
	var userID sql.NullString
	var managerType string

	err := rows.Scan(
		&manager.ID,
		&name,
		&manager.Description,
		&userID,
		&managerType,
		&capabilityStrings,
		&manager.Active,
		&manager.Temperature,
		&manager.MaxTokens,
		&manager.SystemPrompt,
		&manager.CreatedAt,
		&manager.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Handle nullable name field
	if name.Valid {
		manager.Name = name.String
	}

	// Handle nullable userID field
	if userID.Valid {
		manager.UserID = userID.String
	}

	// Convert string to ManagerType
	manager.Type = domain.ManagerType(managerType)

	// Convert string array to capabilities
	manager.Capabilities = make([]domain.ManagerCapability, len(capabilityStrings))
	for i, capStr := range capabilityStrings {
		manager.Capabilities[i] = domain.ManagerCapability(capStr)
	}

	return &manager, nil
}

// Helper method to format table name in queries
func (r CatalogRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
