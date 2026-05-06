package domain

import "context"

type FilterRepository interface {
	Load(ctx context.Context, id string) (*Filter, error)
	Save(ctx context.Context, category *Filter) error
}
