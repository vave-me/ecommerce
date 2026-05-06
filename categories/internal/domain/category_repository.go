package domain

import "context"

type CategoryRepository interface {
	Load(ctx context.Context, id string) (*Category, error)
	Save(ctx context.Context, category *Category) error
}
