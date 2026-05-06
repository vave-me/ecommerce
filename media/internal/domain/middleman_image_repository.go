package domain

import (
	"context"
	"time"
)

type MiddlemanImage struct {
	ID           string
	MediaID      string
	DisplayOrder int
	IsMain       bool
	URL          string
	Metadata     string
	Thumbnail    string
	UserID       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type MiddlemanImageRepository interface {
	AddImage(ctx context.Context, id, mediaID string, displayOrder int, isMain bool, url, metadata, thumbnail, userID string) error
	FindImage(ctx context.Context, id string) (*MiddlemanImage, error)
	FindAllItemImages(ctx context.Context, itemId string) ([]*MiddlemanImage, error)
	FindAllImagesByAuthor(ctx context.Context, userID string, page, pageSize int64, sortBy, sortOrder string) ([]*MiddlemanImage, int64, error)
	FindAllMediaImages(ctx context.Context, mediaId string) ([]*MiddlemanImage, error)
	FindAllItemTypeImages(ctx context.Context, itemId string) ([]*MiddlemanImage, error)
	RemoveImage(ctx context.Context, id string) error
}
