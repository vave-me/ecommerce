package queries

import (
	"context"
	"middleman/services/internal/domain"
	"time"
)

// FilterCriteria represents the complete set of filter options
type GetServicesWithFilter struct {
	// Basic Filters
	CategoryID       string
	CategorySlug     string
	ServiceType      string
	UserID           string
	Status           domain.ServiceStatus
	SearchText       string
	MinPrice         int64
	MaxPrice         int64
	Latitude         float64
	Longitude        float64
	Radius           float64 // in meters
	AvailableFrom    time.Time
	AvailableTo      time.Time
	HasVariants      bool
	Negotiable       bool
	MiddlemanService bool
	UserType         domain.UserType
	Tags             []string
	Qualifications   []string
	Page             int64
	PageSize         int64
	SortBy           string
	SortOrder        string
}

// Validate validates the filter criteria
func (f *GetServicesWithFilter) Validate() error {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = 10
	}
	if f.PageSize > 100 {
		f.PageSize = 100
	}
	if f.SortBy == "" {
		f.SortBy = "name"
	}
	if f.SortOrder != "asc" && f.SortOrder != "desc" {
		f.SortOrder = "asc"
	}
	return nil
}

type GetServicesWithFilterHandler struct {
	catalog domain.CatalogRepository
}

func NewGetServicesWithFilterHandler(catalog domain.CatalogRepository) GetServicesWithFilterHandler {
	return GetServicesWithFilterHandler{catalog: catalog}
}

func (h GetServicesWithFilterHandler) GetServicesWithFilter(ctx context.Context, query GetServicesWithFilter) ([]*domain.CatalogService, int64, error) {
	// Validate filter criteria
	if err := query.Validate(); err != nil {
		return nil, 0, err
	}

	// Execute the filtered query
	return h.catalog.GetServicesWithFilter(
		ctx,
		query.CategoryID,
		query.CategorySlug,
		query.ServiceType,
		query.UserID,
		query.Status,
		query.SearchText,
		query.MinPrice,
		query.MaxPrice,
		query.Latitude,
		query.Longitude,
		query.Radius,
		query.AvailableFrom,
		query.AvailableTo,
		query.HasVariants,
		query.Negotiable,
		query.MiddlemanService,
		query.UserType,
		query.Tags,
		query.Qualifications,
		query.Page,
		query.PageSize,
		query.SortBy,
		query.SortOrder,
	)
}
