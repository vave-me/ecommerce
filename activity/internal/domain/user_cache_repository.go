package domain

import (
	"context"
)

type UserCacheRepository interface {
	Add(ctx context.Context, userID, email, username, location string, enabled bool) error
	Rename(ctx context.Context, userID, name string) error
	UserRepository
}
