package domain

import "context"

type FollowRepository interface {
	Load(ctx context.Context, followID string) (*Follow, error)
	Save(ctx context.Context, follow *Follow) error
}
