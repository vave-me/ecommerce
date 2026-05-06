package domain

import (
	"context"
)

type UserRepository interface {
	Load(ctx context.Context, userID string) (*User, error)
	Save(ctx context.Context, user *User) error
}
