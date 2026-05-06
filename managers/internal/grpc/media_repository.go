package grpc

import (
	"context"
	"fmt"
	"middleman/internal/auth"
	"middleman/internal/rpc"
	"middleman/managers/internal/domain"
	"middleman/managers/internal/models"
	"middleman/media/mediapb"
	"time"

	"google.golang.org/grpc"
)

// MiddlemanMediaRepository calls the remote media service (gRPC) as a fallback.
type MiddlemanMediaRepository struct {
	endpoint string
	auth     *auth.Auth // Optional auth for authenticated calls
}

var _ domain.MiddlemanMediaRepository = (*MiddlemanMediaRepository)(nil)

// NewMiddlemanMediaRepositoryWithAuth creates a new MediaRepository with JWT authentication support
func NewMiddlemanMediaRepository(endpoint string, authInstance *auth.Auth) MiddlemanMediaRepository {
	return MiddlemanMediaRepository{
		endpoint: endpoint,
		auth:     authInstance,
	}
}

// Media Management Methods

// CreateMedia creates a new media record
func (r MiddlemanMediaRepository) CreateMedia(ctx context.Context, itemID, itemType, userID, status string) (*models.CreateMediaResponse, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := mediapb.NewMediaServiceClient(conn)
	resp, err := client.CreateMedia(ctx, &mediapb.CreateMediaRequest{
		ItemId:   itemID,
		ItemType: itemType,
		Status:   status,
	})
	if err != nil {
		return nil, fmt.Errorf("CreateMedia RPC failed: %w", err)
	}

	return &models.CreateMediaResponse{
		ID: resp.GetId(),
	}, nil
}

// UpdateMedia updates an existing media record
func (r MiddlemanMediaRepository) UpdateMedia(ctx context.Context, id, itemID, itemType, userID, status string) (*models.UpdateMediaResponse, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := mediapb.NewMediaServiceClient(conn)
	resp, err := client.UpdateMedia(ctx, &mediapb.UpdateMediaRequest{
		Id:       id,
		ItemId:   itemID,
		ItemType: itemType,
		UserId:   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("UpdateMedia RPC failed: %w", err)
	}

	return &models.UpdateMediaResponse{
		ID: resp.GetId(),
	}, nil
}

// GetMedia retrieves media by ID from the media microservice
func (r MiddlemanMediaRepository) GetMedia(ctx context.Context, mediaID string) (*models.Media, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := mediapb.NewMediaServiceClient(conn)
	resp, err := client.GetMedia(ctx, &mediapb.GetMediaRequest{
		MediaId: mediaID,
	})
	if err != nil {
		return nil, fmt.Errorf("GetMedia RPC failed: %w", err)
	}

	return r.convertMediaFromPb(resp.GetMedia()), nil
}

// GetMediaByItem retrieves media by item ID from the media microservice
func (r MiddlemanMediaRepository) GetMediaByItem(ctx context.Context, itemID string) (*models.Media, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := mediapb.NewMediaServiceClient(conn)
	resp, err := client.GetMediaByItem(ctx, &mediapb.GetMediaByItemRequest{
		ItemId: itemID,
	})
	if err != nil {
		return nil, fmt.Errorf("GetMediaByItem RPC failed: %w", err)
	}

	return r.convertMediaFromPb(resp.GetMedia()), nil
}

// RemoveMedia removes media by ID from the media microservice
func (r MiddlemanMediaRepository) RemoveMedia(ctx context.Context, mediaID, itemID string) (*models.RemoveMediaResponse, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := mediapb.NewMediaServiceClient(conn)
	resp, err := client.RemoveMedia(ctx, &mediapb.RemoveMediaRequest{
		MediaId: mediaID,
		ItemId:  itemID,
	})
	if err != nil {
		return nil, fmt.Errorf("RemoveMedia RPC failed: %w", err)
	}

	return &models.RemoveMediaResponse{
		ID: resp.GetId(),
	}, nil
}

// Image Management Methods

// AddImage adds an image to a media record
func (r MiddlemanMediaRepository) AddImage(ctx context.Context, mediaID string, displayOrder int32, isMain bool, url, metadata, fileType, thumbnail, userID string) (*models.AddImageResponse, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := mediapb.NewMediaServiceClient(conn)
	resp, err := client.AddImage(ctx, &mediapb.AddImageRequest{
		MediaId:      mediaID,
		DisplayOrder: displayOrder,
		IsMain:       isMain,
		Url:          url,
		Metadata:     metadata,
		FileType:     fileType,
		Thumbnail:    thumbnail,
	})
	if err != nil {
		return nil, fmt.Errorf("AddImage RPC failed: %w", err)
	}

	return &models.AddImageResponse{
		URL:     resp.GetUrl(),
		ViewURL: resp.GetViewUrl(),
	}, nil
}

// RemoveImage removes an image from a media record
func (r MiddlemanMediaRepository) RemoveImage(ctx context.Context, mediaID, imageID string) (*models.RemoveImageResponse, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := mediapb.NewMediaServiceClient(conn)
	resp, err := client.RemoveImage(ctx, &mediapb.RemoveImageRequest{
		MediaId: mediaID,
		ImageId: imageID,
	})
	if err != nil {
		return nil, fmt.Errorf("RemoveImage RPC failed: %w", err)
	}

	return &models.RemoveImageResponse{
		ID: resp.GetId(),
	}, nil
}

// GetAllItemImages retrieves all images for an item
func (r MiddlemanMediaRepository) GetAllItemImages(ctx context.Context, itemID string) (*models.GetAllItemImagesResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := mediapb.NewMediaServiceClient(conn)
	resp, err := client.GetAllItemImages(ctx, &mediapb.GetAllItemImagesRequest{
		ItemId: itemID,
	})
	if err != nil {
		return nil, fmt.Errorf("GetAllItemImages RPC failed: %w", err)
	}

	images := make([]models.Image, len(resp.GetImages()))
	for i, img := range resp.GetImages() {
		images[i] = *r.convertImageFromPb(img)
	}

	return &models.GetAllItemImagesResponse{
		Images: images,
	}, nil
}

// GetAllMediaImages retrieves all images for a media record
func (r MiddlemanMediaRepository) GetAllMediaImages(ctx context.Context, mediaID string) (*models.GetAllMediaImagesResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := mediapb.NewMediaServiceClient(conn)
	resp, err := client.GetAllMediaImages(ctx, &mediapb.GetAllMediaImagesRequest{
		MediaId: mediaID,
	})
	if err != nil {
		return nil, fmt.Errorf("GetAllMediaImages RPC failed: %w", err)
	}

	images := make([]models.Image, len(resp.GetImages()))
	for i, img := range resp.GetImages() {
		images[i] = *r.convertImageFromPb(img)
	}

	return &models.GetAllMediaImagesResponse{
		Images: images,
	}, nil
}

// GetAllImagesByAuthor retrieves all images by an author with pagination
func (r MiddlemanMediaRepository) GetAllImagesByAuthor(ctx context.Context, userID string, page, pageSize int64, sortBy, sortOrder string) (*models.GetAllImagesByAuthorResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := mediapb.NewMediaServiceClient(conn)
	resp, err := client.GetAllImagesByAuthor(ctx, &mediapb.GetAllImagesByAuthorRequest{
		UserId:    userID,
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("GetAllImagesByAuthor RPC failed: %w", err)
	}

	images := make([]models.Image, len(resp.GetImages()))
	for i, img := range resp.GetImages() {
		images[i] = *r.convertImageFromPb(img)
	}

	return &models.GetAllImagesByAuthorResponse{
		Images:      images,
		TotalCount:  resp.GetTotalCount(),
		CurrentPage: resp.GetCurrentPage(),
		TotalPages:  resp.GetTotalPages(),
	}, nil
}

// Video Management Methods

// AddVideo adds a video to a media record
func (r MiddlemanMediaRepository) AddVideo(ctx context.Context, mediaID string, displayOrder int32, isMain bool, url, metadata, fileType, thumbnail, userID string) (*models.AddVideoResponse, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := mediapb.NewMediaServiceClient(conn)
	resp, err := client.AddVideo(ctx, &mediapb.AddVideoRequest{
		MediaId:      mediaID,
		DisplayOrder: displayOrder,
		IsMain:       isMain,
		Url:          url,
		Metadata:     metadata,
		FileType:     fileType,
		Thumbnail:    thumbnail,
	})
	if err != nil {
		return nil, fmt.Errorf("AddVideo RPC failed: %w", err)
	}

	return &models.AddVideoResponse{
		URL:     resp.GetUrl(),
		ViewURL: resp.GetViewUrl(),
	}, nil
}

// RemoveVideo removes a video from a media record
func (r MiddlemanMediaRepository) RemoveVideo(ctx context.Context, mediaID, videoID string) (*models.RemoveVideoResponse, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := mediapb.NewMediaServiceClient(conn)
	resp, err := client.RemoveVideo(ctx, &mediapb.RemoveVideoRequest{
		MediaId: mediaID,
		VideoId: videoID,
	})
	if err != nil {
		return nil, fmt.Errorf("RemoveVideo RPC failed: %w", err)
	}

	return &models.RemoveVideoResponse{
		ID: resp.GetId(),
	}, nil
}

// GetAllItemVideos retrieves all videos for an item
func (r MiddlemanMediaRepository) GetAllItemVideos(ctx context.Context, itemID string) (*models.GetAllItemVideosResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := mediapb.NewMediaServiceClient(conn)
	resp, err := client.GetAllItemVideos(ctx, &mediapb.GetAllItemVideosRequest{
		ItemId: itemID,
	})
	if err != nil {
		return nil, fmt.Errorf("GetAllItemVideos RPC failed: %w", err)
	}

	videos := make([]models.Video, len(resp.GetVideos()))
	for i, vid := range resp.GetVideos() {
		videos[i] = *r.convertVideoFromPb(vid)
	}

	return &models.GetAllItemVideosResponse{
		Videos: videos,
	}, nil
}

// GetAllMediaVideos retrieves all videos for a media record
func (r MiddlemanMediaRepository) GetAllMediaVideos(ctx context.Context, mediaID string) (*models.GetAllMediaVideosResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := mediapb.NewMediaServiceClient(conn)
	resp, err := client.GetAllMediaVideos(ctx, &mediapb.GetAllMediaVideosRequest{
		MediaId: mediaID,
	})
	if err != nil {
		return nil, fmt.Errorf("GetAllMediaVideos RPC failed: %w", err)
	}

	videos := make([]models.Video, len(resp.GetVideos()))
	for i, vid := range resp.GetVideos() {
		videos[i] = *r.convertVideoFromPb(vid)
	}

	return &models.GetAllMediaVideosResponse{
		Videos: videos,
	}, nil
}

// GetAllVideos retrieves all videos with pagination
func (r MiddlemanMediaRepository) GetAllVideos(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) (*models.GetAllVideosResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := mediapb.NewMediaServiceClient(conn)
	resp, err := client.GetAllVideos(ctx, &mediapb.GetAllVideosRequest{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("GetAllVideos RPC failed: %w", err)
	}

	videos := make([]models.Video, len(resp.GetVideos()))
	for i, vid := range resp.GetVideos() {
		videos[i] = *r.convertVideoFromPb(vid)
	}

	return &models.GetAllVideosResponse{
		Videos:      videos,
		TotalCount:  resp.GetTotalCount(),
		CurrentPage: resp.GetCurrentPage(),
		TotalPages:  resp.GetTotalPages(),
	}, nil
}

// GetAllVideosByAuthor retrieves all videos by an author with pagination
func (r MiddlemanMediaRepository) GetAllVideosByAuthor(ctx context.Context, userID string, page, pageSize int64, sortBy, sortOrder string) (*models.GetAllVideosByAuthorResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := mediapb.NewMediaServiceClient(conn)
	resp, err := client.GetAllVideosByAuthor(ctx, &mediapb.GetAllVideosByAuthorRequest{
		UserId:    userID,
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("GetAllVideosByAuthor RPC failed: %w", err)
	}

	videos := make([]models.Video, len(resp.GetVideos()))
	for i, vid := range resp.GetVideos() {
		videos[i] = *r.convertVideoFromPb(vid)
	}

	return &models.GetAllVideosByAuthorResponse{
		Videos:      videos,
		TotalCount:  resp.GetTotalCount(),
		CurrentPage: resp.GetCurrentPage(),
		TotalPages:  resp.GetTotalPages(),
	}, nil
}

// Legacy Methods (for backward compatibility)

// AddMedia adds media via the media microservice (legacy method)
func (r MiddlemanMediaRepository) AddMedia(ctx context.Context, id, itemID string, itemType models.EntityType, userID string, status string) error {
	_, err := r.CreateMedia(ctx, itemID, string(itemType), userID, status)
	return err
}

// AddMediaItemOrder adds media item order via the media microservice (legacy method)
func (r MiddlemanMediaRepository) AddMediaItemOrder(ctx context.Context, id, mediaItemId, mediaItemUrl string, displayOrder int) error {
	// This would need to be implemented separately or integrated into the media creation
	return fmt.Errorf("AddMediaItemOrder not directly supported in new API - use CreateMedia and AddImage/AddVideo instead")
}

// GetAllUserMedia retrieves all media for a user from the media microservice (legacy method)
func (r MiddlemanMediaRepository) GetAllUserMedia(ctx context.Context, userID string) ([]*models.Media, error) {
	// This would need a new gRPC method or could be implemented using existing methods
	return nil, fmt.Errorf("GetAllUserMedia not directly supported in new API - use GetAllImagesByAuthor and GetAllVideosByAuthor instead")
}

// GetItemMedia retrieves item media by ID from the media microservice (legacy method)
func (r MiddlemanMediaRepository) GetItemMedia(ctx context.Context, id string) (*models.Media, error) {
	return r.GetMedia(ctx, id)
}

// Helper Methods for protobuf conversion

// convertMediaFromPb converts protobuf Media to domain model
func (r MiddlemanMediaRepository) convertMediaFromPb(pbMedia *mediapb.Media) *models.Media {
	if pbMedia == nil {
		return nil
	}

	mediaOrder := make([]models.MediaOrderItem, len(pbMedia.GetMediaOrder()))
	for i, order := range pbMedia.GetMediaOrder() {
		mediaOrder[i] = models.MediaOrderItem{
			MediaItemID: order.GetMediaItemId(),
			URL:         order.GetUrl(),
		}
	}

	return &models.Media{
		ID:         pbMedia.GetId(),
		ItemID:     pbMedia.GetItemId(),
		ItemType:   pbMedia.GetItemType(),
		UserID:     pbMedia.GetUserId(),
		FileType:   pbMedia.GetFileType(),
		MediaOrder: mediaOrder,
	}
}

// convertImageFromPb converts protobuf Image to domain model
func (r MiddlemanMediaRepository) convertImageFromPb(pbImage *mediapb.Image) *models.Image {
	if pbImage == nil {
		return nil
	}

	return &models.Image{
		ID:           pbImage.GetId(),
		MediaID:      pbImage.GetMediaId(),
		DisplayOrder: pbImage.GetDisplayOrder(),
		IsMain:       pbImage.GetIsMain(),
		URL:          pbImage.GetUrl(),
		Metadata:     pbImage.GetMetadata(),
		FileType:     pbImage.GetFileType(),
		Thumbnail:    pbImage.GetThumbnail(),
		UserID:       pbImage.GetUserId(),
		CreatedAt:    time.Now(), // Would be set from protobuf timestamp if available
		UpdatedAt:    time.Now(),
	}
}

// convertVideoFromPb converts protobuf Video to domain model
func (r MiddlemanMediaRepository) convertVideoFromPb(pbVideo *mediapb.Video) *models.Video {
	if pbVideo == nil {
		return nil
	}

	return &models.Video{
		ID:           pbVideo.GetId(),
		MediaID:      pbVideo.GetMediaId(),
		DisplayOrder: pbVideo.GetDisplayOrder(),
		IsMain:       pbVideo.GetIsMain(),
		URL:          pbVideo.GetUrl(),
		Metadata:     pbVideo.GetMetadata(),
		FileType:     pbVideo.GetFileType(),
		Thumbnail:    pbVideo.GetThumbnail(),
		UserID:       pbVideo.GetUserId(),
		CreatedAt:    time.Now(), // Would be set from protobuf timestamp if available
		UpdatedAt:    time.Now(),
	}
}

// dial establishes a gRPC connection to the media service
// dial sets up a gRPC connection with the microservice endpoint
func (r MiddlemanMediaRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

func (r MiddlemanMediaRepository) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}
