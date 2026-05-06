package domain

import (
	"context"
)

type CatalogAddress struct {
	Address string
	Lat     float64
	Lng     float64
}

type CatalogRepository interface {
	GetAddressForCoordinates(ctx context.Context, lat, lng float64) (*CatalogAddress, error)
	GetCoordinatesForAddress(ctx context.Context, address string) (*CatalogAddress, error)
	GetGeocodingCache(ctx context.Context, address string) (*CatalogAddress, error)
	GetGeocodingDetails(ctx context.Context, address string) (*CatalogAddress, error)
	SuggestAddress(ctx context.Context, address string) ([]*CatalogAddress, error)
	SuggestCity(ctx context.Context, address string) ([]*CatalogAddress, error)
}
