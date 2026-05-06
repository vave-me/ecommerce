package domain

import "context"

type AlertRepository interface {
	Load(ctx context.Context, alertID string) (*Alert, error)
	Save(ctx context.Context, alert *Alert) error
}
