package domain

import (
	"context"
	"middleman/assistants/internal/models"
)

type CategoryRepository interface {
	// Core category operations from protobuf
	CreateNewCategory(ctx context.Context, req *models.AddCategoryRequest) (*models.AddCategoryResponse, error)
	GetCategoryByID(ctx context.Context, id, lang, userID string) (*models.GetCategoryResponse, error)
	FindCategoryBySlugName(ctx context.Context, slug, lang, userID string) (*models.GetCategoryBySlugResponse, error)
	GetPaginatedCategoryList(ctx context.Context, page int64, lang, categoryType string, pageSize int64, sortBy, sortOrder string) (*models.GetCategoriesResponse, error)
	GetMainParentCategories(ctx context.Context, page int64, lang, categoryType string, pageSize int64, sortBy, sortOrder string) (*models.GetMainCategoriesResponse, error)
	GetAllTopLevelCategories(ctx context.Context, page int64, lang string, pageSize int64, sortBy, sortOrder string) (*models.GetAllMainCategoriesResponse, error)
	GetChildCategoriesForParent(ctx context.Context, parentCategoryID, lang string, page, pageSize int64, sortBy, sortOrder string) (*models.GetSubCategoriesResponse, error)
	DeleteCategoryByID(ctx context.Context, id, userID string) (*models.RemoveCategoryResponse, error)
	RenameCategoryBrand(ctx context.Context, req *models.RebrandCategoryRequest) (*models.RebrandCategoryResponse, error)
	ModifyCategoryDetails(ctx context.Context, req *models.UpdateCategoryRequest) (*models.UpdateCategoryResponse, error)
	MoveCategoryToArchive(ctx context.Context, id, userID string) (*models.ArchiveCategoryResponse, error)

	// Filter operations from protobuf
	CreateCategoryFilter(ctx context.Context, req *models.AddFilterRequest) (*models.AddFilterResponse, error)
	GetFilterByID(ctx context.Context, filterID, userID string) (*models.GetFilterResponse, error)
	GetCategoryFiltersList(ctx context.Context, categoryID string, page, pageSize int64, sortBy, sortOrder, userID string) (*models.GetFiltersResponse, error)
	MoveFilterToArchive(ctx context.Context, filterID, userID string) (*models.ArchiveFilterResponse, error)
	DeleteFilterByID(ctx context.Context, filterID, userID string) (*models.RemoveFilterResponse, error)

	// Additional query methods for AI tooling and repository pattern compatibility
	GetPaginatedCategoriesWithSorting(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*models.Category, error)
	SearchCategoriesByKeyword(ctx context.Context, term string) ([]*models.Category, error)
	GetCategoriesFilteredByType(ctx context.Context, categoryType string, page, pageSize int64) ([]*models.Category, error)
	GetActiveCategoriesInLanguage(ctx context.Context, lang string, page, pageSize int64) ([]*models.Category, error)
	GetCategoryStatisticsReport(ctx context.Context, categoryID string) (*models.CategoryStatsResponse, error)
	GetAllFiltersForCategory(ctx context.Context, categoryID string) ([]*models.Filter, error)
	GetChildCategoriesByParentID(ctx context.Context, parentID string, page, pageSize int64) ([]*models.Category, error)
	SearchFiltersMatchingKeyword(ctx context.Context, term string) ([]*models.Filter, error)
	GetMostPopularCategories(ctx context.Context, limit int64) ([]*models.Category, error)
	GetCompleteCategoryHierarchy(ctx context.Context, categoryID string) (*models.CategoryHierarchyResponse, error)

	// Legacy methods for backward compatibility
	FindCategory(ctx context.Context, categoryID string) (*models.Category, error)
	AllCategories(ctx context.Context) ([]*models.Category, error)
	FindBySlug(ctx context.Context, slug string) (*models.Category, error)
	GetCategoryFilters(ctx context.Context, categoryID string) ([]*models.Filter, error)
}
