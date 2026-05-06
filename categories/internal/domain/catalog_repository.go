package domain

import "context"

type CatalogCategory struct {
	ID               string
	Description      string
	ParentID         string
	GoogleCategoryID string
	Tags             []string // optional set of tags
	IsActive         bool
	Slug             string
	SeoTitle         string
	SeoKeywords      []string
	SeoDesc          string
	CategoryType     string
	Lang             string
}

type CatalogRepository interface {
	AddCategory(ctx context.Context, id, description, parentID, googleCategoryID string, tags []string, isActive bool, slug, seoTitle string, seoKeywords []string, seoDesc, categoryType, lang string) error
	UpdateCategory(ctx context.Context, categoryID, description, parentID, googleCategoryID string, tags []string, isActive bool, slug, seoTitle string, seoKeywords []string, seoDesc, categoryType, land string) error
	RemoveCategory(ctx context.Context, categoryID, userID string) error
	ArchiveCategory(ctx context.Context, categoryID, userID string) error
	RebrandCategory(ctx context.Context, categoryID, newDescription, newSlug string) error
	Find(ctx context.Context, categoryID string) (*CatalogCategory, error)
	GetCategories(ctx context.Context, categoryType, lang string, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogCategory, int64, error)
	GetMainCategories(ctx context.Context, categoryType, lang string, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogCategory, int64, error)
	GetAllMainCategories(ctx context.Context, lang string, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogCategory, int64, error)
	GetSubCategories(ctx context.Context, lang string, parentCategoryID string, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogCategory, int64, error)
	GetCatalog(ctx context.Context, lang string, userID string, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogCategory, int64, error)
}
