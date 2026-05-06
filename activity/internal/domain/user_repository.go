package domain

import "context"

type UserRepository interface {
	Find(ctx context.Context, userID string) (*User, error)
}
