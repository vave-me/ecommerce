package tools

import (
	"context"
	"fmt"
)

// ==================== POST HANDLERS ====================
func (r *ComprehensiveToolRegistry) initializePostHandlers() {
	r.handlers["post_create"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		postID := getStringParam(params, "post_id")
		name := getStringParam(params, "name")
		description := getStringParam(params, "description")
		tags := getStringArrayParam(params, "tags")
		status := getStringParam(params, "status")
		lat := getFloat64Param(params, "lat", 0)
		lng := getFloat64Param(params, "lng", 0)
		thumbnail := getStringParam(params, "thumbnail")

		// Validate required parameters
		if err := ValidateIDParam("post_id", postID); err != nil {
			return nil, fmt.Errorf("invalid post_id: %w", err)
		}
		if name == "" {
			return nil, fmt.Errorf("post name is required")
		}
		if description == "" {
			return nil, fmt.Errorf("post description is required")
		}

		// Validate coordinates if provided
		v := NewValidator()
		if lat != 0 {
			v.ValidateLatitude("lat", lat)
		}
		if lng != 0 {
			v.ValidateLongitude("lng", lng)
		}
		if err := v.GetError(); err != nil {
			return nil, fmt.Errorf("invalid coordinates: %w", err)
		}

		// Sanitize string inputs
		name = SanitizeString(name)
		description = SanitizeString(description)

		return nil, reg.postRepo.CreatePost(ctx, postID, name, description, tags, status, lat, lng, thumbnail)
	}

	r.handlers["post_get_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		postID := getStringParam(params, "post_id")
		if err := ValidateIDParam("post_id", postID); err != nil {
			return nil, fmt.Errorf("invalid post_id: %w", err)
		}
		return reg.postRepo.GetPostByID(ctx, postID)
	}

	r.handlers["post_update_details"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		postID := getStringParam(params, "post_id")
		name := getStringParam(params, "name")
		description := getStringParam(params, "description")
		tags := getStringArrayParam(params, "tags")
		status := getStringParam(params, "status")
		thumbnail := getStringParam(params, "thumbnail")

		// Validate required parameters
		if err := ValidateIDParam("post_id", postID); err != nil {
			return nil, fmt.Errorf("invalid post_id: %w", err)
		}

		// Sanitize string inputs
		name = SanitizeString(name)
		description = SanitizeString(description)

		return nil, reg.postRepo.UpdatePostDetails(ctx, postID, name, description, tags, status, thumbnail)
	}

	r.handlers["post_delete"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		postID := getStringParam(params, "post_id")
		if err := ValidateIDParam("post_id", postID); err != nil {
			return nil, fmt.Errorf("invalid post_id: %w", err)
		}
		return nil, reg.postRepo.DeletePost(ctx, postID)
	}

	r.handlers["post_search_by_term"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		searchTerm := getStringParam(params, "search_term")
		if searchTerm == "" {
			return nil, fmt.Errorf("search_term is required")
		}
		searchTerm = SanitizeString(searchTerm)
		return reg.postRepo.SearchPostsByTerm(ctx, searchTerm)
	}

	r.handlers["post_get_suggestions"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		searchTerm := getStringParam(params, "search_term")
		if searchTerm == "" {
			return nil, fmt.Errorf("search_term is required")
		}
		searchTerm = SanitizeString(searchTerm)
		return reg.postRepo.GetPostSuggestions(ctx, searchTerm)
	}

	r.handlers["post_search_advanced"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return reg.postRepo.SearchPostsAdvanced(ctx,
			getStringParam(params, "name"),
			getStringParam(params, "description"),
			getStringArrayParam(params, "tags"),
			getStringParam(params, "status"),
			getStringParam(params, "user_type"),
			getStringParam(params, "thumbnail"),
			getInt64Param(params, "offset", 0),
			getInt64Param(params, "limit", 20),
			getFloat64Param(params, "lat", 0),
			getFloat64Param(params, "lng", 0),
			getInt64Param(params, "radius", 0),
			getInt64Param(params, "page", 1),
			getInt64Param(params, "page_size", 20),
			getStringParam(params, "sort_by"),
			getStringParam(params, "sort_order"))
	}

	r.handlers["post_get_by_category_slug"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		categorySlug := getStringParam(params, "category_slug")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		sortBy := getStringParam(params, "sort_by")
		sortOrder := getStringParam(params, "sort_order")
		if categorySlug == "" {
			return nil, fmt.Errorf("category_slug is required")
		}
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		return reg.postRepo.GetPostsByCategorySlug(ctx, categorySlug, page, pageSize, sortBy, sortOrder)
	}

	r.handlers["post_get_by_category_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		categoryID := getStringParam(params, "category_id")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		sortBy := getStringParam(params, "sort_by")
		sortOrder := getStringParam(params, "sort_order")
		if err := ValidateIDParam("category_id", categoryID); err != nil {
			return nil, fmt.Errorf("invalid category_id: %w", err)
		}
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		return reg.postRepo.GetPostsByCategoryID(ctx, categoryID, page, pageSize, sortBy, sortOrder)
	}

	r.handlers["post_get_all"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		sortBy := getStringParam(params, "sort_by")
		sortOrder := getStringParam(params, "sort_order")
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		return reg.postRepo.GetAllPosts(ctx, page, pageSize, sortBy, sortOrder)
	}

	r.handlers["post_get_by_user_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		sortBy := getStringParam(params, "sort_by")
		sortOrder := getStringParam(params, "sort_order")
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		return reg.postRepo.GetPostsByUserID(ctx, userID, page, pageSize, sortBy, sortOrder)
	}

	r.handlers["post_get_public_catalog"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		sortBy := getStringParam(params, "sort_by")
		sortOrder := getStringParam(params, "sort_order")
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		return reg.postRepo.GetPublicPostCatalog(ctx, page, pageSize, sortBy, sortOrder)
	}

	r.handlers["post_archive_user_post"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		postID := getStringParam(params, "post_id")
		if err := ValidateIDParam("post_id", postID); err != nil {
			return nil, fmt.Errorf("invalid post_id: %w", err)
		}
		return nil, reg.postRepo.ArchiveUserPost(ctx, postID)
	}

	r.handlers["post_update_thumbnail"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		postID := getStringParam(params, "post_id")
		thumbnail := getStringParam(params, "thumbnail")
		if err := ValidateIDParam("post_id", postID); err != nil {
			return nil, fmt.Errorf("invalid post_id: %w", err)
		}
		if thumbnail == "" {
			return nil, fmt.Errorf("thumbnail is required")
		}
		return nil, reg.postRepo.UpdatePostThumbnail(ctx, postID, thumbnail)
	}

	r.handlers["post_add_thumbnail"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		postID := getStringParam(params, "post_id")
		thumbnail := getStringParam(params, "thumbnail")
		if err := ValidateIDParam("post_id", postID); err != nil {
			return nil, fmt.Errorf("invalid post_id: %w", err)
		}
		if thumbnail == "" {
			return nil, fmt.Errorf("thumbnail is required")
		}
		return nil, reg.postRepo.AddThumbnailToPost(ctx, postID, thumbnail)
	}

	r.handlers["post_get_by_location"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		lat := float32(getFloat64Param(params, "lat", 0))
		lng := float32(getFloat64Param(params, "lng", 0))
		radius := getInt64Param(params, "radius", 10)
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)

		// Validate coordinates
		v := NewValidator()
		v.ValidateLatitude("lat", float64(lat))
		v.ValidateLongitude("lng", float64(lng))
		if err := v.GetError(); err != nil {
			return nil, fmt.Errorf("invalid coordinates: %w", err)
		}
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}

		posts, total, err := reg.postRepo.GetPostsByLocation(ctx, lat, lng, radius, page, pageSize)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"posts": posts,
			"total": total,
		}, nil
	}

	r.handlers["post_get_by_tags"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		tags := getStringArrayParam(params, "tags")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		if len(tags) == 0 {
			return nil, fmt.Errorf("at least one tag is required")
		}
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		posts, total, err := reg.postRepo.GetPostsByTags(ctx, tags, page, pageSize)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"posts": posts,
			"total": total,
		}, nil
	}

	r.handlers["post_get_popular"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		posts, total, err := reg.postRepo.GetPopularPosts(ctx, page, pageSize)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"posts": posts,
			"total": total,
		}, nil
	}

	r.handlers["post_get_recent"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		posts, total, err := reg.postRepo.GetRecentPosts(ctx, page, pageSize)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"posts": posts,
			"total": total,
		}, nil
	}

	r.handlers["post_get_trending"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		posts, total, err := reg.postRepo.GetTrendingPosts(ctx, page, pageSize)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"posts": posts,
			"total": total,
		}, nil
	}

	r.handlers["post_add_like"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		postID := getStringParam(params, "post_id")
		if err := ValidateIDParam("post_id", postID); err != nil {
			return nil, fmt.Errorf("invalid post_id: %w", err)
		}
		return nil, reg.postRepo.AddLikeToPost(ctx, postID)
	}

	r.handlers["post_remove_like"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		postID := getStringParam(params, "post_id")
		if err := ValidateIDParam("post_id", postID); err != nil {
			return nil, fmt.Errorf("invalid post_id: %w", err)
		}
		return nil, reg.postRepo.RemoveLikeFromPost(ctx, postID)
	}

	r.handlers["post_share"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		postID := getStringParam(params, "post_id")
		if err := ValidateIDParam("post_id", postID); err != nil {
			return nil, fmt.Errorf("invalid post_id: %w", err)
		}
		return nil, reg.postRepo.SharePostWithUsers(ctx, postID)
	}

	r.handlers["post_report"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		postID := getStringParam(params, "post_id")
		reason := getStringParam(params, "reason")
		if err := ValidateIDParam("post_id", postID); err != nil {
			return nil, fmt.Errorf("invalid post_id: %w", err)
		}
		if reason == "" {
			return nil, fmt.Errorf("report reason is required")
		}
		reason = SanitizeString(reason)
		return nil, reg.postRepo.ReportInappropriatePost(ctx, postID, reason)
	}
}