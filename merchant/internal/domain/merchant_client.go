package domain

import (
	"context"
	"google.golang.org/api/content/v2.1"
)

// MerchantClient defines the interface for Google Merchant Center operations
type MerchantClient interface {
	InsertProduct(ctx context.Context, product *content.Product) error
	UpdateProduct(ctx context.Context, product *content.Product) error
	GetProduct(ctx context.Context, productID string) (*content.Product, error)
	DeleteProduct(ctx context.Context, productID string) error
	ListProducts(ctx context.Context, pageSize int64, pageToken string) ([]*content.Product, string, error)
	IsNotFoundErr(err error) bool
	MerchantID() uint64
}