package domain

import (
	"context"
)

type UserCacheRepository interface {
	Add(ctx context.Context, userID, firstName, lastName, email string) error
	Rename(ctx context.Context, userID, name string) error
	UserRepository
}
