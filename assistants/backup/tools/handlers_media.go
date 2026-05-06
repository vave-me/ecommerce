package tools

import (
	"context"
	"fmt"
	"middleman/assistants/internal/models"
)

// ==================== MEDIA HANDLERS ====================
func (r *ComprehensiveToolRegistry) initializeMediaHandlers() {
	r.handlers["media_create"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		itemID := getStringParam(params, "item_id")
		itemType := getStringParam(params, "item_type")
		userID := getStringParam(params, "user_id")
		status := getStringParam(params, "status", "active")

		// Validate required parameters
		if err := ValidateIDParam("item_id", itemID); err != nil {
			return nil, fmt.Errorf("invalid item_id: %w", err)
		}
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		if itemType == "" {
			return nil, fmt.Errorf("item_type is required")
		}

		return reg.mediaRepo.CreateMedia(ctx, itemID, itemType, userID, status)
	}

	r.handlers["media_update"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		id := getStringParam(params, "id")
		itemID := getStringParam(params, "item_id")
		itemType := getStringParam(params, "item_type")
		userID := getStringParam(params, "user_id")
		status := getStringParam(params, "status")

		// Validate required parameters
		if err := ValidateIDParam("id", id); err != nil {
			return nil, fmt.Errorf("invalid media id: %w", err)
		}
		if err := ValidateIDParam("item_id", itemID); err != nil {
			return nil, fmt.Errorf("invalid item_id: %w", err)
		}
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}

		return reg.mediaRepo.UpdateMedia(ctx, id, itemID, itemType, userID, status)
	}

	r.handlers["media_get"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		mediaID := getStringParam(params, "media_id")
		if err := ValidateIDParam("media_id", mediaID); err != nil {
			return nil, fmt.Errorf("invalid media_id: %w", err)
		}
		return reg.mediaRepo.GetMedia(ctx, mediaID)
	}

	r.handlers["media_get_by_item"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		itemID := getStringParam(params, "item_id")
		if err := ValidateIDParam("item_id", itemID); err != nil {
			return nil, fmt.Errorf("invalid item_id: %w", err)
		}
		return reg.mediaRepo.GetMediaByItem(ctx, itemID)
	}

	r.handlers["media_remove"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		mediaID := getStringParam(params, "media_id")
		itemID := getStringParam(params, "item_id")
		if err := ValidateIDParam("media_id", mediaID); err != nil {
			return nil, fmt.Errorf("invalid media_id: %w", err)
		}
		if err := ValidateIDParam("item_id", itemID); err != nil {
			return nil, fmt.Errorf("invalid item_id: %w", err)
		}
		return reg.mediaRepo.RemoveMedia(ctx, mediaID, itemID)
	}

	r.handlers["media_add_image"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		mediaID := getStringParam(params, "media_id")
		displayOrder := int32(getInt64Param(params, "display_order", 0))
		isMain := getBoolParam(params, "is_main", false)
		url := getStringParam(params, "url")
		metadata := getStringParam(params, "metadata")
		fileType := getStringParam(params, "file_type")
		thumbnail := getStringParam(params, "thumbnail")
		userID := getStringParam(params, "user_id")

		// Validate required parameters
		if err := ValidateIDParam("media_id", mediaID); err != nil {
			return nil, fmt.Errorf("invalid media_id: %w", err)
		}
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		if url == "" {
			return nil, fmt.Errorf("url is required")
		}

		// Sanitize string inputs
		metadata = SanitizeString(metadata)

		return reg.mediaRepo.AddImage(ctx, mediaID, displayOrder, isMain, url, metadata, fileType, thumbnail, userID)
	}

	r.handlers["media_remove_image"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		mediaID := getStringParam(params, "media_id")
		imageID := getStringParam(params, "image_id")
		if err := ValidateIDParam("media_id", mediaID); err != nil {
			return nil, fmt.Errorf("invalid media_id: %w", err)
		}
		if err := ValidateIDParam("image_id", imageID); err != nil {
			return nil, fmt.Errorf("invalid image_id: %w", err)
		}
		return reg.mediaRepo.RemoveImage(ctx, mediaID, imageID)
	}

	r.handlers["media_get_all_item_images"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		itemID := getStringParam(params, "item_id")
		if err := ValidateIDParam("item_id", itemID); err != nil {
			return nil, fmt.Errorf("invalid item_id: %w", err)
		}
		return reg.mediaRepo.GetAllItemImages(ctx, itemID)
	}

	r.handlers["media_get_all_media_images"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		mediaID := getStringParam(params, "media_id")
		if err := ValidateIDParam("media_id", mediaID); err != nil {
			return nil, fmt.Errorf("invalid media_id: %w", err)
		}
		return reg.mediaRepo.GetAllMediaImages(ctx, mediaID)
	}

	r.handlers["media_get_all_images_by_author"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		sortBy := getStringParam(params, "sort_by", "created_at")
		sortOrder := getStringParam(params, "sort_order", "desc")
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		return reg.mediaRepo.GetAllImagesByAuthor(ctx, userID, page, pageSize, sortBy, sortOrder)
	}

	r.handlers["media_add_video"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		mediaID := getStringParam(params, "media_id")
		displayOrder := int32(getInt64Param(params, "display_order", 0))
		isMain := getBoolParam(params, "is_main", false)
		url := getStringParam(params, "url")
		metadata := getStringParam(params, "metadata")
		fileType := getStringParam(params, "file_type")
		thumbnail := getStringParam(params, "thumbnail")
		userID := getStringParam(params, "user_id")

		// Validate required parameters
		if err := ValidateIDParam("media_id", mediaID); err != nil {
			return nil, fmt.Errorf("invalid media_id: %w", err)
		}
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		if url == "" {
			return nil, fmt.Errorf("url is required")
		}

		// Sanitize string inputs
		metadata = SanitizeString(metadata)

		return reg.mediaRepo.AddVideo(ctx, mediaID, displayOrder, isMain, url, metadata, fileType, thumbnail, userID)
	}

	r.handlers["media_remove_video"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		mediaID := getStringParam(params, "media_id")
		videoID := getStringParam(params, "video_id")
		if err := ValidateIDParam("media_id", mediaID); err != nil {
			return nil, fmt.Errorf("invalid media_id: %w", err)
		}
		if err := ValidateIDParam("video_id", videoID); err != nil {
			return nil, fmt.Errorf("invalid video_id: %w", err)
		}
		return reg.mediaRepo.RemoveVideo(ctx, mediaID, videoID)
	}

	r.handlers["media_get_all_item_videos"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		itemID := getStringParam(params, "item_id")
		if err := ValidateIDParam("item_id", itemID); err != nil {
			return nil, fmt.Errorf("invalid item_id: %w", err)
		}
		return reg.mediaRepo.GetAllItemVideos(ctx, itemID)
	}

	r.handlers["media_get_all_media_videos"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		mediaID := getStringParam(params, "media_id")
		if err := ValidateIDParam("media_id", mediaID); err != nil {
			return nil, fmt.Errorf("invalid media_id: %w", err)
		}
		return reg.mediaRepo.GetAllMediaVideos(ctx, mediaID)
	}

	r.handlers["media_get_all_videos"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		sortBy := getStringParam(params, "sort_by", "created_at")
		sortOrder := getStringParam(params, "sort_order", "desc")
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		return reg.mediaRepo.GetAllVideos(ctx, page, pageSize, sortBy, sortOrder)
	}

	r.handlers["media_get_all_videos_by_author"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		sortBy := getStringParam(params, "sort_by", "created_at")
		sortOrder := getStringParam(params, "sort_order", "desc")
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		return reg.mediaRepo.GetAllVideosByAuthor(ctx, userID, page, pageSize, sortBy, sortOrder)
	}

	r.handlers["media_add_media"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		id := getStringParam(params, "id")
		itemID := getStringParam(params, "item_id")
		itemType := models.EntityType(getStringParam(params, "item_type"))
		userID := getStringParam(params, "user_id")
		status := getStringParam(params, "status", "active")

		// Validate required parameters
		if err := ValidateIDParam("id", id); err != nil {
			return nil, fmt.Errorf("invalid id: %w", err)
		}
		if err := ValidateIDParam("item_id", itemID); err != nil {
			return nil, fmt.Errorf("invalid item_id: %w", err)
		}
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		if string(itemType) == "" {
			return nil, fmt.Errorf("item_type is required")
		}

		return nil, reg.mediaRepo.AddMedia(ctx, id, itemID, itemType, userID, status)
	}

	r.handlers["media_add_media_item_order"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		id := getStringParam(params, "id")
		mediaItemId := getStringParam(params, "media_item_id")
		mediaItemUrl := getStringParam(params, "media_item_url")
		displayOrder := int(getInt64Param(params, "display_order", 0))

		// Validate required parameters
		if err := ValidateIDParam("id", id); err != nil {
			return nil, fmt.Errorf("invalid id: %w", err)
		}
		if err := ValidateIDParam("media_item_id", mediaItemId); err != nil {
			return nil, fmt.Errorf("invalid media_item_id: %w", err)
		}
		if mediaItemUrl == "" {
			return nil, fmt.Errorf("media_item_url is required")
		}

		return nil, reg.mediaRepo.AddMediaItemOrder(ctx, id, mediaItemId, mediaItemUrl, displayOrder)
	}

	r.handlers["media_get_all_user_media"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		return reg.mediaRepo.GetAllUserMedia(ctx, userID)
	}

	r.handlers["media_get_item_media"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		id := getStringParam(params, "id")
		if err := ValidateIDParam("id", id); err != nil {
			return nil, fmt.Errorf("invalid id: %w", err)
		}
		return reg.mediaRepo.GetItemMedia(ctx, id)
	}
}