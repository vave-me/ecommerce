package domain

import (
	"context"
)

// LiveStreamRepository defines the interface for live stream persistence
type LiveStreamRepository interface {
	Find(ctx context.Context, streamID string) (*LiveStream, error)
	FindByStatus(ctx context.Context, status LiveStreamStatus) ([]*LiveStream, error)
	FindByScheduledTime(ctx context.Context, startTime, endTime string) ([]*LiveStream, error)
	Save(ctx context.Context, stream *LiveStream) error
	Update(ctx context.Context, stream *LiveStream) error
	Delete(ctx context.Context, streamID string) error
}