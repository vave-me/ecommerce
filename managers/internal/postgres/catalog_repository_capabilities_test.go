package postgres_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"middleman/managers/internal/domain"
	"middleman/managers/internal/postgres"
)

func TestCatalogRepository_UpdateConfigurationWithCapabilities(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := postgres.NewCatalogRepository("managers", db)
	ctx := context.Background()

	managerID := "test-manager-id"
	temperature := 0.8
	maxTokens := 2000
	systemPrompt := "Updated system prompt"
	capabilities := []domain.ManagerCapability{
		domain.CapabilityUserInteraction,
		domain.CapabilityDataAnalysis,
		domain.CapabilityLocationServices,
		domain.CapabilityAuthentication,
	}
	updatedAt := time.Now()

	// Expected capability strings
	capabilityStrings := []string{
		"user_interaction",
		"data_analysis",
		"location_services",
		"authentication",
	}

	// Set up expectation
	mock.ExpectExec("UPDATE managers SET").
		WithArgs(
			managerID,
			temperature,
			maxTokens,
			systemPrompt,
			pq.StringArray(capabilityStrings),
			updatedAt,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Execute
	err = repo.UpdateConfigurationWithCapabilities(
		ctx,
		managerID,
		temperature,
		maxTokens,
		systemPrompt,
		capabilities,
		updatedAt,
	)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCatalogRepository_UpdateConfigurationWithCapabilities_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := postgres.NewCatalogRepository("managers", db)
	ctx := context.Background()

	managerID := "non-existent-id"
	temperature := 0.8
	maxTokens := 2000
	systemPrompt := "Updated system prompt"
	capabilities := []domain.ManagerCapability{
		domain.CapabilityUserInteraction,
	}
	updatedAt := time.Now()

	// Set up expectation - no rows affected
	mock.ExpectExec("UPDATE managers SET").
		WithArgs(
			managerID,
			temperature,
			maxTokens,
			systemPrompt,
			pq.StringArray([]string{"user_interaction"}),
			updatedAt,
		).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Execute
	err = repo.UpdateConfigurationWithCapabilities(
		ctx,
		managerID,
		temperature,
		maxTokens,
		systemPrompt,
		capabilities,
		updatedAt,
	)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCatalogRepository_UpdateConfigurationWithCapabilities_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := postgres.NewCatalogRepository("managers", db)
	ctx := context.Background()

	managerID := "test-manager-id"
	temperature := 0.8
	maxTokens := 2000
	systemPrompt := "Updated system prompt"
	capabilities := []domain.ManagerCapability{
		domain.CapabilityUserInteraction,
	}
	updatedAt := time.Now()

	// Set up expectation - database error
	mock.ExpectExec("UPDATE managers SET").
		WithArgs(
			managerID,
			temperature,
			maxTokens,
			systemPrompt,
			pq.StringArray([]string{"user_interaction"}),
			updatedAt,
		).
		WillReturnError(sql.ErrConnDone)

	// Execute
	err = repo.UpdateConfigurationWithCapabilities(
		ctx,
		managerID,
		temperature,
		maxTokens,
		systemPrompt,
		capabilities,
		updatedAt,
	)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update manager configuration with capabilities")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCatalogRepository_Add_WithDuplicateCapabilities(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := postgres.NewCatalogRepository("managers", db)
	ctx := context.Background()

	// Note: The database trigger should handle deduplication
	capabilities := []domain.ManagerCapability{
		domain.CapabilityUserInteraction,
		domain.CapabilityDataAnalysis,
		domain.CapabilityUserInteraction, // Duplicate
		domain.CapabilityLocationServices,
		domain.CapabilityDataAnalysis, // Duplicate
	}

	// Expected capability strings (with duplicates - trigger will handle)
	capabilityStrings := []string{
		"user_interaction",
		"data_analysis",
		"user_interaction",
		"location_services",
		"data_analysis",
	}

	// Set up expectation
	mock.ExpectExec("INSERT INTO managers").
		WithArgs(
			"test-id",
			"Test Manager",
			"Test Description",
			"user-123",
			"standard",
			pq.StringArray(capabilityStrings),
			true,
			0.7,
			1000,
			"Test prompt",
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Execute
	err = repo.Add(
		ctx,
		"test-id",
		"Test Manager",
		"Test Description",
		"user-123",
		domain.ManagerTypeStandard,
		capabilities,
		true,
		0.7,
		1000,
		"Test prompt",
	)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCatalogRepository_FindAll_WithCapabilities(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := postgres.NewCatalogRepository("managers", db)
	ctx := context.Background()

	// Set up mock rows
	rows := sqlmock.NewRows([]string{
		"id", "manager_name", "description", "user_id", "type",
		"capabilities", "enabled", "temperature", "max_tokens",
		"system_prompt", "created_at", "updated_at",
	}).
		AddRow(
			"manager-1",
			"Manager 1",
			"Description 1",
			"user-123",
			"standard",
			pq.StringArray{"user_interaction", "data_analysis"},
			true,
			0.7,
			1000,
			"Prompt 1",
			time.Now(),
			time.Now(),
		).
		AddRow(
			"manager-2",
			"Manager 2",
			"Description 2",
			"user-123",
			"admin",
			pq.StringArray{"user_interaction", "data_analysis", "authentication", "manager_management"},
			true,
			0.8,
			2000,
			"Prompt 2",
			time.Now(),
			time.Now(),
		)

	mock.ExpectQuery("SELECT .* FROM managers").
		WillReturnRows(rows)

	// Execute
	managers, err := repo.FindAll(ctx, "user-123")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, managers, 2)

	// Check first manager
	assert.Equal(t, "manager-1", managers[0].ID)
	assert.Len(t, managers[0].Capabilities, 2)
	assert.Contains(t, managers[0].Capabilities, domain.CapabilityUserInteraction)
	assert.Contains(t, managers[0].Capabilities, domain.CapabilityDataAnalysis)

	// Check second manager
	assert.Equal(t, "manager-2", managers[1].ID)
	assert.Len(t, managers[1].Capabilities, 4)
	assert.Contains(t, managers[1].Capabilities, domain.CapabilityUserInteraction)
	assert.Contains(t, managers[1].Capabilities, domain.CapabilityDataAnalysis)
	assert.Contains(t, managers[1].Capabilities, domain.CapabilityAuthentication)
	assert.Contains(t, managers[1].Capabilities, domain.CapabilityManagerManagement)

	assert.NoError(t, mock.ExpectationsWereMet())
}
