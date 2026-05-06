package domain

import "context"

type VideoRepository interface {
	Load(ctx context.Context, id string) (*Video, error)
	Save(ctx context.Context, video *Video) error
}
