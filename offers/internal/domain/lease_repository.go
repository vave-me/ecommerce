package domain

import (
	"context"
)

type LeaseRepository interface {
	Load(ctx context.Context, leaseID string) (*Lease, error)
	Save(ctx context.Context, lease *Lease) error
}
