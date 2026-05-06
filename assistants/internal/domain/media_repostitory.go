package domain

import (
	"context"
	"middleman/assistants/internal/models"
)

type MiddlemanMediaRepository interface {
	// Media Management
	CreateMedia(ctx context.Context, itemID, itemType, userID, status string) (*models.CreateMediaResponse, error)
	UpdateMedia(ctx context.Context, id, itemID, itemType, userID, status string) (*models.UpdateMediaResponse, error)
	GetMedia(ctx context.Context, mediaID string) (*models.Media, error)
	GetMediaByItem(ctx context.Context, itemID string) (*models.Media, error)
	RemoveMedia(ctx context.Context, mediaID, itemID string) (*models.RemoveMediaResponse, error)
	AddImage(ctx context.Context, mediaID string, displayOrder int32, isMain bool, url, metadata, fileType, thumbnail, userID string) (*models.AddImageResponse, error)
	RemoveImage(ctx context.Context, mediaID, imageID string) (*models.RemoveImageResponse, error)
	GetAllItemImages(ctx context.Context, itemID string) (*models.GetAllItemImagesResponse, error)
	GetAllMediaImages(ctx context.Context, mediaID string) (*models.GetAllMediaImagesResponse, error)
	GetAllImagesByAuthor(ctx context.Context, userID string, page, pageSize int64, sortBy, sortOrder string) (*models.GetAllImagesByAuthorResponse, error)
	AddVideo(ctx context.Context, mediaID string, displayOrder int32, isMain bool, url, metadata, fileType, thumbnail, userID string) (*models.AddVideoResponse, error)
	RemoveVideo(ctx context.Context, mediaID, videoID string) (*models.RemoveVideoResponse, error)
	GetAllItemVideos(ctx context.Context, itemID string) (*models.GetAllItemVideosResponse, error)
	GetAllMediaVideos(ctx context.Context, mediaID string) (*models.GetAllMediaVideosResponse, error)
	GetAllVideos(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) (*models.GetAllVideosResponse, error)
	GetAllVideosByAuthor(ctx context.Context, userID string, page, pageSize int64, sortBy, sortOrder string) (*models.GetAllVideosByAuthorResponse, error)
	AddMedia(ctx context.Context, id, itemID string, itemType models.EntityType, userID string, status string) error
	AddMediaItemOrder(ctx context.Context, id, mediaItemId, mediaItemUrl string, displayOrder int) error
	GetAllUserMedia(ctx context.Context, userID string) ([]*models.Media, error)
	GetItemMedia(ctx context.Context, id string) (*models.Media, error)
}
