package postgres

import (
	"middleman/geocoding/internal/domain"
	"middleman/internal/postgres"
)

type CatalogLocationRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.CatalogRepository = (*CatalogRepository)(nil)

func NewCatalogLocationRepository(tableName string, db postgres.DB) CatalogRepository {
	return CatalogRepository{
		tableName: tableName,
		db:        db,
	}
}
