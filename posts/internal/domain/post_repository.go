package domain

import "context"

type PostRepository interface {
	Load(ctx context.Context, id string) (*Post, error)
	Save(ctx context.Context, post *Post) error
}
