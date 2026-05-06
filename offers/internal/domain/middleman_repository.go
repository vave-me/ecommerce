package domain

import "context"

type MiddlemanOffer struct {
	ID             string
	UserSellerID   string
	UserCustomerID string
	ProductID      string
	Price          int64
}

type MiddlemanRepository interface {
	Add(ctx context.Context, offerID, userSellerID, userCustomerID, productID string, price int64) error
	Find(ctx context.Context, offerID string) (*MiddlemanOffer, error)
	All(ctx context.Context, userID string) ([]*MiddlemanOffer, error)
}
