package domain

import "context"

type ReviewRepository interface {
	Load(ctx context.Context, reviewID string) (*Review, error)
	Save(ctx context.Context, review *Review) error
}
