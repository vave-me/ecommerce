package domain

import "context"

type MiddlemanNewsletter struct {
	ID             string
	UserSellerID   string
	UserCustomerID string
	ProductID      string
	Price          int64
}

type MiddlemanRepository interface {
	Add(ctx context.Context, newsletterID, userSellerID, userCustomerID, productID string, price int64) error
	Find(ctx context.Context, newsletterID string) (*MiddlemanNewsletter, error)
	All(ctx context.Context, userID string) ([]*MiddlemanNewsletter, error)
}
