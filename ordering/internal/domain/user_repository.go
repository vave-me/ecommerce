package domain

import "context"

type UserRepository interface {
	Authorize(ctx context.Context, userCustomerID string) error
}
