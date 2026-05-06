package domain

import "context"

type MediaRepository interface {
	Load(ctx context.Context, id string) (*Media, error)
	Save(ctx context.Context, media *Media) error
}
