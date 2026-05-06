package queries

import (
	"context"
	"middleman/products/internal/domain"
)

type GetProductsWithFilters struct {
	Name             string
	Category         string
	MinPrice         int64
	MaxPrice         int64
	Brand            string
	Condition        string
	Model            string
	Tags             []string
	ManageStock      bool
	MinStock         int64
	MaxStock         int64
	SKU              string
	Status           string
	Negotiable       bool
	UserType         string
	MiddlemanService bool
	HasVariants      bool
	ShippingCost     int64
	MinWeight        int64
	MaxWeight        int64
	MinHeight        int64
	MaxHeight        int64
	MinWidth         int64
	MaxWidth         int64
	MinDepth         int64
	MaxDepth         int64
	Offset           int64
	Limit            int64
	Lat              float64
	Lng              float64
	Radius           int64
	Page             int64
	PageSize         int64
	SortBy           string
	SortOrder        string
}
type GetProductsWithFiltersHandler struct {
	catalog domain.CatalogRepository
}

func NewGetProductsWithFiltersHandler(catalog domain.CatalogRepository) GetProductsWithFiltersHandler {
	return GetProductsWithFiltersHandler{catalog: catalog}
}

func (h GetProductsWithFiltersHandler) GetProductsWithFilters(ctx context.Context, query GetProductsWithFilters) ([]*domain.CatalogProduct, int64, error) {
	// Set defaults for pagination and sorting
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 10
	}
	if query.SortBy == "" {
		query.SortBy = "createdAt" // Default sort field
	}
	if query.SortOrder != "asc" && query.SortOrder != "desc" {
		query.SortOrder = "desc" // Default sort order
	}

	return h.catalog.GetProductsWithFilters(ctx, query.Name, query.Category, query.MinPrice, query.MaxPrice, query.Brand, query.Condition, query.Model, query.Tags,
		query.ManageStock, query.MinStock, query.MaxStock, query.SKU, query.Status, query.Negotiable, query.UserType, query.MiddlemanService, query.HasVariants,
		query.ShippingCost, query.MinWeight, query.MaxWeight, query.MinHeight, query.MaxHeight, query.MinWidth, query.MaxWidth, query.MinDepth, query.MaxDepth, query.Offset, query.Limit, query.Lat, query.Lng, query.Radius, query.Page, query.PageSize, query.SortBy, query.SortOrder)
}
