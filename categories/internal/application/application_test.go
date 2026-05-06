package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"middleman/categories/internal/application"
	"middleman/categories/internal/application/commands"
	"middleman/categories/internal/application/queries"
	"middleman/categories/internal/domain"
	"middleman/internal/ddd"
)

// Example test suite for application.Application

func TestApplication_AddCategory(t *testing.T) {
	// GIVEN: A mock CategoryRepository and event publisher
	catRepo := domain.NewMockCategoryRepository(t)
	publisher := ddd.NewMockEventPublisher[ddd.Event](t)

	// You may want to set up catRepo or publisher expectations.
	// For example:
	// catRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(nil)
	// publisher.On("Publish", mock.Anything, mock.AnythingOfType("ddd.event")).Return(nil)

	// We also need the other repositories for the New() constructor
	// but we can stub them out with minimal mocks:
	catalogRepo := domain.NewMockCatalogRepository(t)
	cacheRepo := domain.NewMockCatalogRepository(t)
	filterRepo := domain.NewMockFilterRepository(t)
	filterCatRepo := domain.NewMockCatalogFilterRepository(t)

	// WHEN: We build the real Application
	app := application.New(catRepo, catalogRepo, cacheRepo, filterRepo, filterCatRepo, publisher)

	// THEN: We call AddCategory
	err := app.AddCategory(context.Background(), commands.AddCategory{
		ID:               "Test Category",
		Description:      "Test Category",
		ParentCategoryID: "Main Category",
		GoogleCategoryID: "Google",
		Slug:             "/test/category",
		// fill other fields as needed
	})

	// ASSERT: For now, we simply check that no error occurred
	assert.NoError(t, err)

	// Optionally verify any expectations on catRepo or publisher.
	catRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestApplication_RebrandCategory(t *testing.T) {
	catRepo := domain.NewMockCategoryRepository(t)
	publisher := ddd.NewMockEventPublisher[ddd.Event](t)

	catalogRepo := domain.NewMockCatalogRepository(t)
	cacheRepo := domain.NewMockCatalogRepository(t)
	filterRepo := domain.NewMockFilterRepository(t)
	filterCatRepo := domain.NewMockCatalogFilterRepository(t)

	app := application.New(catRepo, catalogRepo, cacheRepo, filterRepo, filterCatRepo, publisher)

	// Example: Return an error from the category repo to see if it's handled
	catRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(errors.New("update failed"))

	err := app.RebrandCategory(context.Background(), commands.RebrandCategory{
		ID:          "cat-123",
		Description: "New Name",
	})

	assert.Error(t, err, "Should return error if underlying update fails")
	catRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestApplication_RemoveCategory(t *testing.T) {
	catRepo := domain.NewMockCategoryRepository(t)
	publisher := ddd.NewMockEventPublisher[ddd.Event](t)

	app := application.New(
		catRepo,
		domain.NewMockCatalogRepository(t),
		domain.NewMockCatalogRepository(t),
		domain.NewMockFilterRepository(t),
		domain.NewMockCatalogFilterRepository(t),
		publisher,
	)

	// Suppose removing a category is always successful in the repo
	catRepo.On("Remove", mock.Anything, "cat-999").Return(nil)

	err := app.RemoveCategory(context.Background(), commands.RemoveCategory{
		ID: "cat-999",
	})
	assert.NoError(t, err)

	catRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestApplication_ArchiveCategory(t *testing.T) {
	catRepo := domain.NewMockCategoryRepository(t)
	publisher := ddd.NewMockEventPublisher[ddd.Event](t)

	app := application.New(
		catRepo,
		domain.NewMockCatalogRepository(t),
		domain.NewMockCatalogRepository(t),
		domain.NewMockFilterRepository(t),
		domain.NewMockCatalogFilterRepository(t),
		publisher,
	)

	catRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Category")).
		Return(nil)

	err := app.ArchiveCategory(context.Background(), commands.ArchiveCategory{
		ID: "cat-321",
	})
	assert.NoError(t, err)

	catRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestApplication_AddFilter(t *testing.T) {
	filterRepo := domain.NewMockFilterRepository(t)
	publisher := ddd.NewMockEventPublisher[ddd.Event](t)

	app := application.New(
		domain.NewMockCategoryRepository(t),
		domain.NewMockCatalogRepository(t),
		domain.NewMockCatalogRepository(t),
		filterRepo,
		domain.NewMockCatalogFilterRepository(t),
		publisher,
	)

	filterRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.Filter")).
		Return(nil)

	err := app.AddFilter(context.Background(), commands.AddFilter{
		Name: "Size",
	})
	assert.NoError(t, err)

	filterRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestApplication_RebrandFilter(t *testing.T) {
	filterRepo := domain.NewMockFilterRepository(t)
	publisher := ddd.NewMockEventPublisher[ddd.Event](t)

	app := application.New(
		domain.NewMockCategoryRepository(t),
		domain.NewMockCatalogRepository(t),
		domain.NewMockCatalogRepository(t),
		filterRepo,
		domain.NewMockCatalogFilterRepository(t),
		publisher,
	)

	filterRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Filter")).
		Return(nil)

	err := app.RebrandFilter(context.Background(), commands.RebrandFilter{
		ID:   "filter-123",
		Name: "Updated Name",
	})
	assert.NoError(t, err)

	filterRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestApplication_ArchiveFilter(t *testing.T) {
	filterRepo := domain.NewMockFilterRepository(t)
	publisher := ddd.NewMockEventPublisher[ddd.Event](t)

	app := application.New(
		domain.NewMockCategoryRepository(t),
		domain.NewMockCatalogRepository(t),
		domain.NewMockCatalogRepository(t),
		filterRepo,
		domain.NewMockCatalogFilterRepository(t),
		publisher,
	)

	filterRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Filter")).
		Return(nil)

	err := app.ArchiveFilter(context.Background(), commands.ArchiveFilter{
		ID: "filter-321",
	})
	assert.NoError(t, err)

	filterRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

func TestApplication_RemoveFilter(t *testing.T) {
	filterRepo := domain.NewMockFilterRepository(t)
	publisher := ddd.NewMockEventPublisher[ddd.Event](t)

	app := application.New(
		domain.NewMockCategoryRepository(t),
		domain.NewMockCatalogRepository(t),
		domain.NewMockCatalogRepository(t),
		filterRepo,
		domain.NewMockCatalogFilterRepository(t),
		publisher,
	)

	filterRepo.On("Remove", mock.Anything, "filter-001").Return(nil)

	err := app.RemoveFilter(context.Background(), commands.RemoveFilter{
		ID: "filter-001",
	})
	assert.NoError(t, err)

	filterRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

// -------------------------------------------------
// Queries
// -------------------------------------------------

func TestApplication_GetCatalog(t *testing.T) {
	catRepo := domain.NewMockCatalogRepository(t)
	app := application.New(
		domain.NewMockCategoryRepository(t),
		catRepo,
		domain.NewMockCatalogRepository(t),
		domain.NewMockFilterRepository(t),
		domain.NewMockCatalogFilterRepository(t),
		ddd.NewMockEventPublisher[ddd.Event](t),
	)

	// Suppose the catalog repo returns two categories
	catRepo.On("FindAll", mock.Anything, mock.AnythingOfType("queries.GetCatalog")).
		Return([]*domain.CatalogCategory{
			{ID: "cat-1"}, {ID: "cat-2"},
		}, int64(2), nil).
		Once()

	cats, total, err := app.GetCatalog(context.Background(), queries.GetCatalog{
		// possible filter/pagination
	})
	assert.NoError(t, err)
	assert.Len(t, cats, 2)
	assert.EqualValues(t, 2, total)

	catRepo.AssertExpectations(t)
}

func TestApplication_GetCategory(t *testing.T) {
	catRepo := domain.NewMockCatalogRepository(t)
	app := application.New(
		domain.NewMockCategoryRepository(t),
		catRepo,
		domain.NewMockCatalogRepository(t),
		domain.NewMockFilterRepository(t),
		domain.NewMockCatalogFilterRepository(t),
		ddd.NewMockEventPublisher[ddd.Event](t),
	)

	expectedCat := &domain.CatalogCategory{ID: "cat-xyz"}

	// The "GetCategoryHandler" calls something like catRepo.FindOne(...) or so
	catRepo.On("FindOne", mock.Anything, mock.AnythingOfType("queries.GetCategory")).
		Return(expectedCat, nil).
		Once()

	cat, err := app.GetCategory(context.Background(), queries.GetCategory{
		ID: "cat-xyz",
	})
	assert.NoError(t, err)
	assert.Equal(t, "cat-xyz", cat.ID)
}

func TestApplication_GetCategoryBySlug(t *testing.T) {
	catRepo := domain.NewMockCatalogRepository(t)
	app := application.New(
		domain.NewMockCategoryRepository(t),
		catRepo,
		domain.NewMockCatalogRepository(t),
		domain.NewMockFilterRepository(t),
		domain.NewMockCatalogFilterRepository(t),
		ddd.NewMockEventPublisher[ddd.Event](t),
	)

	expectedCat := &domain.CatalogCategory{ID: "cat-abc", Slug: "shoes"}

	catRepo.On("FindBySlug", mock.Anything, "shoes").
		Return(expectedCat, nil)

	cat, err := app.GetCategoryBySlug(context.Background(), queries.GetCategoryBySlug{
		Slug: "shoes",
	})
	assert.NoError(t, err)
	assert.Equal(t, "cat-abc", cat.ID)
	assert.Equal(t, "shoes", cat.Slug)
}

func TestApplication_GetCategories(t *testing.T) {
	catRepo := domain.NewMockCatalogRepository(t)
	app := application.New(
		domain.NewMockCategoryRepository(t),
		domain.NewMockCatalogRepository(t),
		catRepo,
		domain.NewMockFilterRepository(t),
		domain.NewMockCatalogFilterRepository(t),
		ddd.NewMockEventPublisher[ddd.Event](t),
	)

	catRepo.On("FindAll", mock.Anything, mock.AnythingOfType("queries.GetCategories")).
		Return([]*domain.CatalogCategory{{ID: "cat-A"}}, int64(1), nil)

	cats, total, err := app.GetCategories(context.Background(), queries.GetCategories{})
	assert.NoError(t, err)
	assert.Len(t, cats, 1)
	assert.EqualValues(t, 1, total)
}

func TestApplication_GetMainCategories(t *testing.T) {
	catRepo := domain.NewMockCatalogRepository(t)
	app := application.New(
		domain.NewMockCategoryRepository(t),
		domain.NewMockCatalogRepository(t),
		catRepo,
		domain.NewMockFilterRepository(t),
		domain.NewMockCatalogFilterRepository(t),
		ddd.NewMockEventPublisher[ddd.Event](t),
	)

	catRepo.On("FindMain", mock.Anything, mock.AnythingOfType("queries.GetMainCategories")).
		Return([]*domain.CatalogCategory{{ID: "main-1"}}, int64(1), nil)

	cats, total, err := app.GetMainCategories(context.Background(), queries.GetMainCategories{})
	assert.NoError(t, err)
	assert.Len(t, cats, 1)
	assert.EqualValues(t, 1, total)
}

func TestApplication_GetSubCategories(t *testing.T) {
	catRepo := domain.NewMockCatalogRepository(t)
	app := application.New(
		domain.NewMockCategoryRepository(t),
		domain.NewMockCatalogRepository(t),
		catRepo,
		domain.NewMockFilterRepository(t),
		domain.NewMockCatalogFilterRepository(t),
		ddd.NewMockEventPublisher[ddd.Event](t),
	)

	catRepo.On("FindSubs", mock.Anything, mock.AnythingOfType("queries.GetSubCategories")).
		Return([]*domain.CatalogCategory{
			{ID: "sub-1"}, {ID: "sub-2"},
		}, int64(2), nil)

	subs, total, err := app.GetSubCategories(context.Background(), queries.GetSubCategories{
		ParentCategoryID: "cat-parent",
	})
	assert.NoError(t, err)
	assert.Len(t, subs, 2)
	assert.EqualValues(t, 2, total)
}
