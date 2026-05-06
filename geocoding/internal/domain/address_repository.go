package domain

import (
	"context"
)

type AddressRepository interface {
	Load(ctx context.Context, id string) (*Address, error)
	Save(ctx context.Context, address *Address) error
}
