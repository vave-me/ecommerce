package domain

import (
	"context"
	"time"
)

type MiddlemanVideo struct {
	ID           string
	MediaID      string
	DisplayOrder int
	IsMain       bool
	Url          string
	Metadata     string
	Thumbnail    string
	UserID       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type MiddlemanVideoRepository interface {
	AddVideo(ctx context.Context, id, mediaID string, displayOrder int, isMain bool, url, metadata, thumbnail, userID string) error
	FindVideo(ctx context.Context, id string) (*MiddlemanVideo, error)
	FindAllItemVideos(ctx context.Context, itemId string) ([]*MiddlemanVideo, error)
	FindAllVideos(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*MiddlemanVideo, int64, error)
	FindAllVideosByAuthor(ctx context.Context, userID string, page, pageSize int64, sortBy, sortOrder string) ([]*MiddlemanVideo, int64, error)
	FindAllMediaVideos(ctx context.Context, mediaId string) ([]*MiddlemanVideo, error)
	FindAllItemTypeVideos(ctx context.Context, itemId string) ([]*MiddlemanVideo, error)
	RemoveVideo(ctx context.Context, id string) error
}
