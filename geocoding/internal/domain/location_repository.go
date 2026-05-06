package domain

import (
	"context"
)

type LocationRepository interface {
	Load(ctx context.Context, id string) (*Location, error)
	Save(ctx context.Context, location *Location) error
}
