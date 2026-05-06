package models

import "time"

// Core category entities
type Category struct {
	ID               string    `json:"id"`
	Description      string    `json:"description"`
	ParentID         string    `json:"parent_id"`
	GoogleCategoryID string    `json:"google_category_id"`
	Tags             []string  `json:"tags"`
	Slug             string    `json:"slug"`
	IsActive         bool      `json:"is_active"`
	SeoTitle         string    `json:"seo_title"`
	SeoKeywords      []string  `json:"seo_keywords"`
	SeoDesc          string    `json:"seo_desc"`
	CategoryType     string    `json:"category_type"`
	Lang             string    `json:"lang"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	ArchivedAt       time.Time `json:"archived_at,omitempty"`

	// Extended fields for comprehensive functionality
	ItemCount   int64  `json:"item_count,omitempty"`
	Level       int32  `json:"level,omitempty"`
	SortOrder   int32  `json:"sort_order,omitempty"`
	IconURL     string `json:"icon_url,omitempty"`
	BannerURL   string `json:"banner_url,omitempty"`
	MetaTitle   string `json:"meta_title,omitempty"`
	MetaDesc    string `json:"meta_desc,omitempty"`
	IsPublic    bool   `json:"is_public"`
	IsFeatured  bool   `json:"is_featured"`
	ViewCount   int64  `json:"view_count,omitempty"`
	SearchCount int64  `json:"search_count,omitempty"`
}

type Filter struct {
	ID         string    `json:"id"`
	CategoryID string    `json:"category_id"`
	Name       string    `json:"name"`
	FilterType string    `json:"filter_type"`
	Values     []string  `json:"values"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	ArchivedAt time.Time `json:"archived_at,omitempty"`

	// Extended fields
	IsRequired   bool   `json:"is_required"`
	SortOrder    int32  `json:"sort_order"`
	Placeholder  string `json:"placeholder,omitempty"`
	HelpText     string `json:"help_text,omitempty"`
	DefaultValue string `json:"default_value,omitempty"`
}

// Category protobuf request/response types

// AddCategory
type AddCategoryRequest struct {
	Description      string   `json:"description"`
	ParentID         string   `json:"parent_id"`
	GoogleCategoryID string   `json:"google_category_id"`
	Tags             []string `json:"tags"`
	Slug             string   `json:"slug"`
	IsActive         bool     `json:"is_active"`
	SeoTitle         string   `json:"seo_title"`
	SeoKeywords      []string `json:"seo_keywords"`
	SeoDesc          string   `json:"seo_desc"`
	CategoryType     string   `json:"category_type"`
	Lang             string   `json:"lang"`
}

type AddCategoryResponse struct {
	ID string `json:"id"`
}

// GetCategory
type GetCategoryRequest struct {
	ID     string `json:"id"`
	Lang   string `json:"lang"`
	UserID string `json:"user_id"`
}

type GetCategoryResponse struct {
	Category *Category `json:"category"`
}

// GetCategoryBySlug
type GetCategoryBySlugRequest struct {
	Slug   string `json:"slug"`
	Lang   string `json:"lang"`
	UserID string `json:"user_id"`
}

type GetCategoryBySlugResponse struct {
	Category *Category `json:"category"`
}

// GetCategories
type GetCategoriesRequest struct {
	Page         int64  `json:"page"`
	Lang         string `json:"lang"`
	CategoryType string `json:"category_type"`
	PageSize     int64  `json:"page_size"`
	SortBy       string `json:"sort_by"`
	SortOrder    string `json:"sort_order"`
}

type GetCategoriesResponse struct {
	Categories  []*Category `json:"categories"`
	TotalCount  int64       `json:"total_count"`
	TotalPages  int64       `json:"total_pages"`
	CurrentPage int64       `json:"current_page"`
}

// GetMainCategories
type GetMainCategoriesRequest struct {
	Page         int64  `json:"page"`
	Lang         string `json:"lang"`
	CategoryType string `json:"category_type"`
	PageSize     int64  `json:"page_size"`
	SortBy       string `json:"sort_by"`
	SortOrder    string `json:"sort_order"`
}

type GetMainCategoriesResponse struct {
	Categories  []*Category `json:"categories"`
	TotalCount  int64       `json:"total_count"`
	TotalPages  int64       `json:"total_pages"`
	CurrentPage int64       `json:"current_page"`
}

// GetAllMainCategories
type GetAllMainCategoriesRequest struct {
	Page      int64  `json:"page"`
	Lang      string `json:"lang"`
	PageSize  int64  `json:"page_size"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}

type GetAllMainCategoriesResponse struct {
	Categories  []*Category `json:"categories"`
	TotalCount  int64       `json:"total_count"`
	TotalPages  int64       `json:"total_pages"`
	CurrentPage int64       `json:"current_page"`
}

// GetSubCategories
type GetSubCategoriesRequest struct {
	ParentCategoryID string `json:"parent_category_id"`
	Lang             string `json:"lang"`
	Page             int64  `json:"page"`
	PageSize         int64  `json:"page_size"`
	SortBy           string `json:"sort_by"`
	SortOrder        string `json:"sort_order"`
}

type GetSubCategoriesResponse struct {
	Categories  []*Category `json:"categories"`
	TotalCount  int64       `json:"total_count"`
	TotalPages  int64       `json:"total_pages"`
	CurrentPage int64       `json:"current_page"`
}

// RemoveCategory
type RemoveCategoryRequest struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

type RemoveCategoryResponse struct {
	ID string `json:"id"`
}

// RebrandCategory
type RebrandCategoryRequest struct {
	ID      string `json:"id"`
	NewSlug string `json:"new_slug"`
	NewDesc string `json:"new_desc"`
	UserID  string `json:"user_id"`
}

type RebrandCategoryResponse struct {
	Success bool `json:"success"`
}

// UpdateCategory
type UpdateCategoryRequest struct {
	ID               string   `json:"id"`
	Description      string   `json:"description"`
	ParentID         string   `json:"parent_id"`
	GoogleCategoryID string   `json:"google_category_id"`
	Tags             []string `json:"tags"`
	Slug             string   `json:"slug"`
	IsActive         bool     `json:"is_active"`
	SeoTitle         string   `json:"seo_title"`
	SeoKeywords      []string `json:"seo_keywords"`
	SeoDesc          string   `json:"seo_desc"`
	CategoryType     string   `json:"category_type"`
	Lang             string   `json:"lang"`
}

type UpdateCategoryResponse struct {
	ID string `json:"id"`
}

// ArchiveCategory
type ArchiveCategoryRequest struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

type ArchiveCategoryResponse struct {
	CategoryID string `json:"category_id"`
	Archived   bool   `json:"archived"`
}

// Filter protobuf request/response types

// AddFilter
type AddFilterRequest struct {
	CategoryID string   `json:"category_id"`
	Name       string   `json:"name"`
	FilterType string   `json:"filter_type"`
	Values     []string `json:"values"`
	UserID     string   `json:"user_id"`
	IsActive   bool     `json:"is_active"`
}

type AddFilterResponse struct {
	FilterID string `json:"filter_id"`
}

// GetFilter
type GetFilterRequest struct {
	FilterID string `json:"filter_id"`
	UserID   string `json:"user_id"`
}

type GetFilterResponse struct {
	Filter *Filter `json:"filter"`
}

// GetFilters
type GetFiltersRequest struct {
	CategoryID string `json:"category_id"`
	Page       int64  `json:"page"`
	PageSize   int64  `json:"page_size"`
	SortBy     string `json:"sort_by"`
	SortOrder  string `json:"sort_order"`
	UserID     string `json:"user_id"`
}

type GetFiltersResponse struct {
	Filters     []*Filter `json:"filters"`
	TotalCount  int64     `json:"total_count"`
	TotalPages  int64     `json:"total_pages"`
	CurrentPage int64     `json:"current_page"`
}

// ArchiveFilter
type ArchiveFilterRequest struct {
	FilterID string `json:"filter_id"`
	UserID   string `json:"user_id"`
}

type ArchiveFilterResponse struct {
	FilterID string `json:"filter_id"`
	Archived bool   `json:"archived"`
}

// RemoveFilter
type RemoveFilterRequest struct {
	FilterID string `json:"filter_id"`
	UserID   string `json:"user_id"`
}

type RemoveFilterResponse struct {
	FilterID string `json:"filter_id"`
}

// Extended response types for AI tooling and comprehensive functionality

type CategoryStatsResponse struct {
	CategoryID       string `json:"category_id"`
	ItemCount        int64  `json:"item_count"`
	SubCategoryCount int64  `json:"sub_category_count"`
	FilterCount      int64  `json:"filter_count"`
	ViewCount        int64  `json:"view_count"`
	SearchCount      int64  `json:"search_count"`
	LastActivityAt   string `json:"last_activity_at"`
}

type CategoryHierarchyResponse struct {
	Category    *Category   `json:"category"`
	Parent      *Category   `json:"parent,omitempty"`
	Children    []*Category `json:"children"`
	Breadcrumbs []*Category `json:"breadcrumbs"`
	Level       int32       `json:"level"`
	TotalLevels int32       `json:"total_levels"`
}

// Filter types constants
const (
	FilterTypeText     = "text"
	FilterTypeNumber   = "number"
	FilterTypeSelect   = "select"
	FilterTypeCheckbox = "checkbox"
	FilterTypeRange    = "range"
	FilterTypeDate     = "date"
	FilterTypeColor    = "color"
	FilterTypeBoolean  = "boolean"
)

// Category type constants
const (
	CategoryTypeProduct  = "product"
	CategoryTypeService  = "service"
	CategoryTypeJob      = "job"
	CategoryTypeProperty = "property"
	CategoryTypeVehicle  = "vehicle"
	CategoryTypeDeal     = "deal"
	CategoryTypePost     = "post"
	CategoryTypeGeneral  = "general"
)

// Category status constants
const (
	CategoryStatusActive   = "active"
	CategoryStatusInactive = "inactive"
	CategoryStatusArchived = "archived"
	CategoryStatusDraft    = "draft"
	CategoryStatusPending  = "pending"
)
