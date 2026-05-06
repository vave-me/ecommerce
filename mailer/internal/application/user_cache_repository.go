package application

import (
	"context"
)

type UserCacheRepository interface {
	Add(ctx context.Context, userID, email, firstname, lastname string, enabled bool) error
	UserRepository
}
