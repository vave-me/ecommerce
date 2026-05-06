package postgres

import (
	"context"
	"middleman/geocoding/internal/domain"
	"middleman/internal/postgres"
)

type CatalogRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.CatalogRepository = (*CatalogRepository)(nil)

func NewCatalogRepository(tableName string, db postgres.DB) CatalogRepository {
	return CatalogRepository{
		tableName: tableName,
		db:        db,
	}
}
func (r CatalogRepository) GetAddressForCoordinates(ctx context.Context, lat, lng float64) (*domain.CatalogAddress, error) {
	return nil, nil
}
func (r CatalogRepository) GetCoordinatesForAddress(ctx context.Context, address string) (*domain.CatalogAddress, error) {
	return nil, nil
}
func (r CatalogRepository) GetGeocodingCache(ctx context.Context, address string) (*domain.CatalogAddress, error) {
	return nil, nil
}
func (r CatalogRepository) GetGeocodingDetails(ctx context.Context, address string) (*domain.CatalogAddress, error) {
	return nil, nil
}
func (r CatalogRepository) SuggestAddress(ctx context.Context, address string) ([]*domain.CatalogAddress, error) {
	return nil, nil
}
func (r CatalogRepository) SuggestCity(ctx context.Context, address string) ([]*domain.CatalogAddress, error) {
	return nil, nil
}
