package domain

import (
	"context"
)

type UserCacheRepository interface {
	Add(ctx context.Context, userSellerID, email, username, firstName, lastName, location string, enabled bool) error
	Rename(ctx context.Context, userSellerID, name string) error
	UserRepository
}
