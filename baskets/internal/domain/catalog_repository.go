package domain

import "context"

type CatalogBasket struct {
	ID     string
	UserID string
	Status BasketStatus
}

type CatalogRepository interface {
	Add(ctx context.Context, basketID, userID string, status BasketStatus) error
	Remove(ctx context.Context, basketID string) error
	Find(ctx context.Context, userID string) (*CatalogBasket, error)
	All(ctx context.Context, activityID string) ([]*CatalogBasket, error)
}
