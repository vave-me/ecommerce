package domain

import "context"

type MiddlemanSupport struct {
	ID             string
	UserSellerID   string
	UserCustomerID string
	ProductID      string
	Price          int64
}

type MiddlemanRepository interface {
	Add(ctx context.Context, supportID, userSellerID, userCustomerID, productID string, price int64) error
	Find(ctx context.Context, supportID string) (*MiddlemanSupport, error)
	All(ctx context.Context, userID string) ([]*MiddlemanSupport, error)
}
