package domain

import "context"

type ImageRepository interface {
	Load(ctx context.Context, id string) (*Image, error)
	Save(ctx context.Context, video *Image) error
}
