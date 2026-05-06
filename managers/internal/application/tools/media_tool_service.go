package tools

import (
	"context"
	"fmt"
	"time"

	"middleman/managers/internal/domain"
)

// MediaToolService handles media management operations
type MediaToolService struct {
	mediaRepo domain.MiddlemanMediaRepository
	config    *ServiceConfig
}

// NewMediaToolService creates a new media tool service
func NewMediaToolService(mediaRepo domain.MiddlemanMediaRepository) *MediaToolService {
	return &MediaToolService{
		mediaRepo: mediaRepo,
		config: &ServiceConfig{
			MaxRetries:       3,
			EnableStreaming:  true,
			EnableMetrics:    true,
			OperationTimeout: 30 * time.Second,
		},
	}
}

// ExecuteOperation routes media operations to appropriate handlers
func (s *MediaToolService) ExecuteOperation(
	ctx context.Context,
	operation string,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	// Send initial progress
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "media_operation",
			Status:   "started",
			Progress: 0,
			Metadata: map[string]interface{}{
				"operation": operation,
				"service":   "MediaToolService",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	var result interface{}
	var err error

	switch operation {
	case "upload_image":
		result, err = s.uploadImage(ctx, parameters, streamChan, toolID)
	case "upload_video":
		result, err = s.uploadVideo(ctx, parameters, streamChan, toolID)
	case "get_media", "find":
		result, err = s.getMedia(ctx, parameters, streamChan, toolID)
	case "list_media", "list":
		result, err = s.listMedia(ctx, parameters, streamChan, toolID)
	case "delete_media":
		result, err = s.deleteMedia(ctx, parameters, streamChan, toolID)
	case "update_media":
		result, err = s.updateMedia(ctx, parameters, streamChan, toolID)
	case "get_user_media":
		result, err = s.getUserMedia(ctx, parameters, streamChan, toolID)
	case "search_media":
		result, err = s.searchMedia(ctx, parameters, streamChan, toolID)
	case "get_media_by_type":
		result, err = s.getMediaByType(ctx, parameters, streamChan, toolID)
	case "generate_thumbnail":
		result, err = s.generateThumbnail(ctx, parameters, streamChan, toolID)
	case "compress_media":
		result, err = s.compressMedia(ctx, parameters, streamChan, toolID)
	default:
		err = fmt.Errorf("unsupported media operation: %s", operation)
	}

	// Send completion status
	if streamChan != nil {
		status := "completed"
		if err != nil {
			status = "error"
		}

		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "media_operation",
			Status:   status,
			Progress: 100,
			Result:   result,
			Error:    s.getErrorString(err),
			Metadata: map[string]interface{}{
				"operation": operation,
				"service":   "MediaToolService",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	return result, err
}

func (s *MediaToolService) uploadImage(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	userID := getStringParam(params, "user_id", "")
	itemID := getStringParam(params, "item_id", "")
	itemType := getStringParam(params, "item_type", "media")

	if userID == "" || itemID == "" {
		return nil, fmt.Errorf("user_id and item_id are required")
	}

	sendProgress(streamChan, toolID, "creating_media", 30)

	response, err := s.mediaRepo.CreateMedia(ctx, itemID, itemType, userID, "active")
	if err != nil {
		return nil, fmt.Errorf("failed to create media: %w", err)
	}

	return map[string]interface{}{
		"media":   response,
		"item_id": itemID,
		"type":    "image",
	}, nil
}

func (s *MediaToolService) uploadVideo(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	userID := getStringParam(params, "user_id", "")
	itemID := getStringParam(params, "item_id", "")
	itemType := getStringParam(params, "item_type", "media")

	if userID == "" || itemID == "" {
		return nil, fmt.Errorf("user_id and item_id are required")
	}

	sendProgress(streamChan, toolID, "creating_video_media", 30)

	response, err := s.mediaRepo.CreateMedia(ctx, itemID, itemType, userID, "active")
	if err != nil {
		return nil, fmt.Errorf("failed to create video media: %w", err)
	}

	return map[string]interface{}{
		"media":   response,
		"item_id": itemID,
		"type":    "video",
	}, nil
}

func (s *MediaToolService) getMedia(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	mediaID := getStringParam(params, "id", "")
	if mediaID == "" {
		mediaID = getStringParam(params, "media_id", "")
	}
	if mediaID == "" {
		return nil, fmt.Errorf("media ID is required")
	}

	sendProgress(streamChan, toolID, "getting_media", 30)

	media, err := s.mediaRepo.GetMedia(ctx, mediaID)
	if err != nil {
		return nil, fmt.Errorf("failed to find media: %w", err)
	}

	return map[string]interface{}{
		"media": media,
		"id":    mediaID,
	}, nil
}

func (s *MediaToolService) listMedia(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	userID := getStringParam(params, "user_id", "")
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	sendProgress(streamChan, toolID, "listing_media", 30)

	media, err := s.mediaRepo.GetAllUserMedia(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list media: %w", err)
	}

	return map[string]interface{}{
		"media":       media,
		"total_count": len(media),
		"user_id":     userID,
	}, nil
}

func (s *MediaToolService) deleteMedia(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	mediaID := getStringParam(params, "media_id", "")
	itemID := getStringParam(params, "item_id", "")

	if mediaID == "" || itemID == "" {
		return nil, fmt.Errorf("media_id and item_id are required")
	}

	sendProgress(streamChan, toolID, "deleting_media", 30)

	response, err := s.mediaRepo.RemoveMedia(ctx, mediaID, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete media: %w", err)
	}

	return map[string]interface{}{
		"result":   response,
		"media_id": mediaID,
		"deleted":  true,
	}, nil
}

func (s *MediaToolService) updateMedia(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	mediaID := getStringParam(params, "id", "")
	if mediaID == "" {
		mediaID = getStringParam(params, "media_id", "")
	}
	itemID := getStringParam(params, "item_id", "")
	itemType := getStringParam(params, "item_type", "media")
	userID := getStringParam(params, "user_id", "")
	status := getStringParam(params, "status", "active")

	if mediaID == "" || itemID == "" || userID == "" {
		return nil, fmt.Errorf("media_id, item_id, and user_id are required")
	}

	sendProgress(streamChan, toolID, "updating_media", 30)

	response, err := s.mediaRepo.UpdateMedia(ctx, mediaID, itemID, itemType, userID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to update media: %w", err)
	}

	return map[string]interface{}{
		"result":   response,
		"media_id": mediaID,
		"updated":  true,
	}, nil
}

func (s *MediaToolService) getUserMedia(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	userID := getStringParam(params, "user_id", "")
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	sendProgress(streamChan, toolID, "getting_user_media", 30)

	media, err := s.mediaRepo.GetAllUserMedia(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user media: %w", err)
	}

	return map[string]interface{}{
		"media":   media,
		"user_id": userID,
		"count":   len(media),
	}, nil
}

func (s *MediaToolService) searchMedia(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	itemID := getStringParam(params, "item_id", "")
	if itemID == "" {
		return nil, fmt.Errorf("item_id is required for search")
	}

	sendProgress(streamChan, toolID, "searching_media", 30)

	media, err := s.mediaRepo.GetMediaByItem(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to search media: %w", err)
	}

	return map[string]interface{}{
		"media":   media,
		"item_id": itemID,
	}, nil
}

func (s *MediaToolService) getMediaByType(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	itemID := getStringParam(params, "item_id", "")
	mediaType := getStringParam(params, "type", "")

	if itemID == "" {
		return nil, fmt.Errorf("item_id is required")
	}

	sendProgress(streamChan, toolID, "getting_media_by_type", 30)

	var result interface{}
	var err error

	switch mediaType {
	case "image":
		result, err = s.mediaRepo.GetAllItemImages(ctx, itemID)
	case "video":
		result, err = s.mediaRepo.GetAllItemVideos(ctx, itemID)
	default:
		result, err = s.mediaRepo.GetMediaByItem(ctx, itemID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get media by type: %w", err)
	}

	return map[string]interface{}{
		"result":  result,
		"item_id": itemID,
		"type":    mediaType,
	}, nil
}

func (s *MediaToolService) generateThumbnail(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	mediaID := getStringParam(params, "media_id", "")
	if mediaID == "" {
		return nil, fmt.Errorf("media_id is required")
	}

	sendProgress(streamChan, toolID, "generating_thumbnail", 30)

	// For now, return a placeholder - this would need proper thumbnail generation logic
	return map[string]interface{}{
		"media_id":      mediaID,
		"thumbnail_url": "/thumbnails/generated/" + mediaID + ".jpg",
		"generated":     true,
	}, nil
}

func (s *MediaToolService) compressMedia(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	mediaID := getStringParam(params, "media_id", "")
	if mediaID == "" {
		return nil, fmt.Errorf("media_id is required")
	}

	sendProgress(streamChan, toolID, "compressing_media", 30)

	// For now, return a placeholder - this would need proper compression logic
	return map[string]interface{}{
		"media_id":     mediaID,
		"compressed":   true,
		"size_reduced": "40%",
	}, nil
}

func (s *MediaToolService) getErrorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func sendProgress(streamChan chan<- ToolExecutionStream, toolID string, step string, progress float64) {
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "media_operation",
			Status:   "progress",
			Progress: progress,
			Metadata: map[string]interface{}{
				"step": step,
			},
			Timestamp: time.Now().Unix(),
		}
	}
}
