package domain

import "context"

// -----------------------------------------------------------------------------
// 2) CatalogFilter represents a minimal "query model" for filters
// -----------------------------------------------------------------------------
type CatalogFilter struct {
	ID         string
	CategoryID string
	Name       string
	FilterType FilterType
	Values     []string
	IsActive   bool
}

// -----------------------------------------------------------------------------
// CatalogFilterRepository defines the persistence interface for your "Filter" data.
// -----------------------------------------------------------------------------
type CatalogFilterRepository interface {
	// AddFilter: create a new filter record
	AddFilter(ctx context.Context,
		filterID string,
		categoryID string,
		name string,
		filterType FilterType,
		values []string,
		isActive bool,
	) error

	// UpdateFilter: update filter (rename, change type, etc.)
	UpdateFilter(ctx context.Context,
		filterID string,
		name string,
		filterType FilterType,
		values []string,
	) error

	// RemoveFilter: remove or mark as removed
	RemoveFilter(ctx context.Context, filterID string, userID string) error

	// ArchiveFilter: sets isActive=false or similar
	ArchiveFilter(ctx context.Context, filterID string) error

	// RebrandFilter: specialized method to rename or rebrand
	RebrandFilter(ctx context.Context, filterID string, newName string) error

	// FindFilter: get single filter by ID
	FindFilter(ctx context.Context, filterID string) (*CatalogFilter, error)

	// GetFilters: generic listing
	GetFilters(ctx context.Context,
		page, pageSize int64,
		sortBy, sortOrder string,
	) ([]*CatalogFilter, int64, error)

	// GetFiltersByCategory: filters for a specific category
	GetFiltersByCategory(ctx context.Context,
		categoryID string,
		page, pageSize int64,
		sortBy, sortOrder string,
	) ([]*CatalogFilter, int64, error)
}
