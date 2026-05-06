package domain

import "context"

type VariantRepository interface {
	Load(ctx context.Context, id string) (*Variant, error)
	Save(ctx context.Context, product *Variant) error
}
