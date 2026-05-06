package domain

import "context"

type ActivityRepository interface {
	Load(ctx context.Context, activityID string) (*Activity, error)
	Save(ctx context.Context, activityID *Activity) error
}
