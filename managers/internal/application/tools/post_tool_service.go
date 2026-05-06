package tools

import (
	"context"
	"fmt"
	"log"
	"time"

	"middleman/managers/internal/domain"
)

// PostToolService handles all post-related operations
type PostToolService struct {
	postRepo domain.PostRepository
}

// NewPostToolService creates a new post tool service
func NewPostToolService(postRepo domain.PostRepository) *PostToolService {
	return &PostToolService{
		postRepo: postRepo,
	}
}

// ExecuteOperation handles post-related operations with streaming support
func (s *PostToolService) ExecuteOperation(ctx context.Context, operation string, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	log.Printf("PostToolService.ExecuteOperation: Executing post operation: %s", operation)

	// Send initial progress
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "post_operation",
		Status:   "started",
		Progress: 0,
		Metadata: map[string]interface{}{
			"operation": operation,
			"message":   fmt.Sprintf("Starting post operation: %s", operation),
		},
		Timestamp: time.Now().Unix(),
	}

	// Progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "post_operation",
		Status:   "progress",
		Progress: 20,
		Metadata: map[string]interface{}{
			"operation": operation,
			"step":      "processing_parameters",
		},
		Timestamp: time.Now().Unix(),
	}

	var result interface{}
	var err error

	switch operation {
	case "find", "get", "get_post":
		result, err = s.handleFindPost(ctx, parameters, streamChan, toolID)
	case "search", "search_posts":
		result, err = s.handleSearchPosts(ctx, parameters, streamChan, toolID)
	case "suggest", "suggest_posts":
		result, err = s.handleSuggestPosts(ctx, parameters, streamChan, toolID)
	case "get_posts", "list", "get_posts_list":
		result, err = s.handleGetPosts(ctx, parameters, streamChan, toolID)
	case "search_by_category", "category_search":
		result, err = s.handleSearchPostsByCategory(ctx, parameters, streamChan, toolID)
	case "search_by_category_slug", "category_slug_search":
		result, err = s.handleSearchPostsByCategorySlug(ctx, parameters, streamChan, toolID)
	case "filter", "filter_posts":
		result, err = s.handleFilterPosts(ctx, parameters, streamChan, toolID)
	case "add", "create", "add_post", "create_post":
		result, err = s.handleAddPost(ctx, parameters, streamChan, toolID)
	case "update", "update_post":
		result, err = s.handleUpdatePost(ctx, parameters, streamChan, toolID)
	case "remove", "delete", "remove_post", "delete_post":
		result, err = s.handleRemovePost(ctx, parameters, streamChan, toolID)
	case "archive", "archive_post":
		result, err = s.handleArchivePost(ctx, parameters, streamChan, toolID)
	case "get_posts_by_user", "user_posts":
		result, err = s.handleGetPostsByUser(ctx, parameters, streamChan, toolID)
	case "get_posts_by_location", "location_posts":
		result, err = s.handleGetPostsByLocation(ctx, parameters, streamChan, toolID)
	case "get_posts_by_tags", "tagged_posts":
		result, err = s.handleGetPostsByTags(ctx, parameters, streamChan, toolID)
	case "get_popular_posts", "popular_posts":
		result, err = s.handleGetPopularPosts(ctx, parameters, streamChan, toolID)
	case "get_recent_posts", "recent_posts":
		result, err = s.handleGetRecentPosts(ctx, parameters, streamChan, toolID)
	case "get_trending_posts", "trending_posts":
		result, err = s.handleGetTrendingPosts(ctx, parameters, streamChan, toolID)
	case "like_post", "like":
		result, err = s.handleLikePost(ctx, parameters, streamChan, toolID)
	case "unlike_post", "unlike":
		result, err = s.handleUnlikePost(ctx, parameters, streamChan, toolID)
	case "share_post", "share":
		result, err = s.handleSharePost(ctx, parameters, streamChan, toolID)
	case "report_post", "report":
		result, err = s.handleReportPost(ctx, parameters, streamChan, toolID)
	default:
		err = fmt.Errorf("unsupported post operation: %s", operation)
	}

	// Handle result
	if err != nil {
		streamChan <- ToolExecutionStream{
			ID:        toolID,
			ToolName:  "post_operation",
			Status:    "error",
			Progress:  100,
			Error:     err.Error(),
			Timestamp: time.Now().Unix(),
		}
		return nil, err
	}

	// Send success
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "post_operation",
		Status:   "completed",
		Progress: 100,
		Result:   result,
		Metadata: map[string]interface{}{
			"operation": operation,
			"message":   fmt.Sprintf("Post operation %s completed successfully", operation),
		},
		Timestamp: time.Now().Unix(),
	}

	return result, nil
}

func (s *PostToolService) handleFindPost(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	postID := getStringParam(parameters, "id", "")
	if postID == "" {
		postID = getStringParam(parameters, "post_id", "")
	}
	if postID == "" {
		return nil, fmt.Errorf("id or post_id parameter required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "post_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":    "finding_post",
			"post_id": postID,
		},
		Timestamp: time.Now().Unix(),
	}

	post, err := s.postRepo.Find(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("post find failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "posts",
		"operation":   "find",
		"result":      post,
		"post_id":     postID,
	}, nil
}

func (s *PostToolService) handleSearchPosts(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	searchTerm := getStringParam(parameters, "search_term", "")
	if searchTerm == "" {
		searchTerm = getStringParam(parameters, "term", "")
	}
	if searchTerm == "" {
		// If no search term, use filter with minimal parameters for a general list
		name := getStringParam(parameters, "name", "")
		description := getStringParam(parameters, "description", "")
		tags := getStringSliceParam(parameters, "tags")
		status := getStringParam(parameters, "status", "")
		userType := getStringParam(parameters, "user_type", "")
		thumbnail := getStringParam(parameters, "thumbnail", "")
		offset := getInt64Param(parameters, "offset", 0)
		limit := getInt64Param(parameters, "limit", 20)
		lat := getFloat32Param(parameters, "lat", 0.0)
		lng := getFloat32Param(parameters, "lng", 0.0)
		radius := getInt64Param(parameters, "radius", 0)
		page := getInt64Param(parameters, "page", 1)
		pageSize := getInt64Param(parameters, "page_size", 20)
		sortBy := getStringParam(parameters, "sort_by", "created_at")
		sortOrder := getStringParam(parameters, "sort_order", "desc")

		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "post_operation",
			Status:   "progress",
			Progress: 50,
			Metadata: map[string]interface{}{
				"step": "searching_posts_with_filters",
			},
			Timestamp: time.Now().Unix(),
		}

		posts, err := s.postRepo.SearchPostsWithFilters(
			ctx, name, description, tags, status, userType, thumbnail,
			offset, limit, float64(lat), float64(lng), radius, page, pageSize, sortBy, sortOrder,
		)
		if err != nil {
			return nil, fmt.Errorf("post search failed: %w", err)
		}

		return map[string]interface{}{
			"entity_type": "posts",
			"operation":   "search",
			"results":     posts,

			"page":      page,
			"page_size": pageSize,
		}, nil
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "post_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":        "searching_posts",
			"search_term": searchTerm,
		},
		Timestamp: time.Now().Unix(),
	}

	posts, err := s.postRepo.SearchWithTerm(ctx, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("post search failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "posts",
		"operation":   "search",
		"results":     posts,
		"count":       len(posts),
		"search_term": searchTerm,
	}, nil
}

func (s *PostToolService) handleSuggestPosts(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	name := getStringParam(parameters, "name", "")
	if name == "" {
		name = getStringParam(parameters, "title", "")
	}
	if name == "" {
		return nil, fmt.Errorf("name or title parameter required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "post_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step": "suggesting_posts",
			"name": name,
		},
		Timestamp: time.Now().Unix(),
	}

	posts, err := s.postRepo.SuggestPosts(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("post suggest failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "posts",
		"operation":   "suggest",
		"results":     posts,
		"count":       len(posts),
		"name":        name,
	}, nil
}

func (s *PostToolService) handleGetPosts(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	page := getInt64Param(parameters, "page", 1)
	pageSize := getInt64Param(parameters, "page_size", 20)
	sortBy := getStringParam(parameters, "sort_by", "created_at")
	sortOrder := getStringParam(parameters, "sort_order", "desc")

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "post_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":       "getting_posts",
			"page":       page,
			"page_size":  pageSize,
			"sort_by":    sortBy,
			"sort_order": sortOrder,
		},
		Timestamp: time.Now().Unix(),
	}

	posts, err := s.postRepo.GetPosts(ctx, page, pageSize, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("get posts failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "posts",
		"operation":   "get_posts",
		"results":     posts,
		"page":        page,
		"page_size":   pageSize,
		"sort_by":     sortBy,
		"sort_order":  sortOrder,
	}, nil
}

func (s *PostToolService) handleSearchPostsByCategory(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	categoryID := getStringParam(parameters, "category_id", "")
	if categoryID == "" {
		return nil, fmt.Errorf("category_id parameter required")
	}

	page := getInt64Param(parameters, "page", 1)
	pageSize := getInt64Param(parameters, "page_size", 20)
	sortBy := getStringParam(parameters, "sort_by", "created_at")
	sortOrder := getStringParam(parameters, "sort_order", "desc")

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "post_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":        "searching_posts_by_category",
			"category_id": categoryID,
		},
		Timestamp: time.Now().Unix(),
	}

	posts, err := s.postRepo.SearchPostsWithCategory(ctx, categoryID, page, pageSize, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("search posts by category failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "posts",
		"operation":   "search_by_category",
		"results":     posts,
		"category_id": categoryID,
		"page":        page,
		"page_size":   pageSize,
	}, nil
}

func (s *PostToolService) handleSearchPostsByCategorySlug(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	categorySlug := getStringParam(parameters, "category_slug", "")
	if categorySlug == "" {
		return nil, fmt.Errorf("category_slug parameter required")
	}

	page := getInt64Param(parameters, "page", 1)
	pageSize := getInt64Param(parameters, "page_size", 20)
	sortBy := getStringParam(parameters, "sort_by", "created_at")
	sortOrder := getStringParam(parameters, "sort_order", "desc")

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "post_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":          "searching_posts_by_category_slug",
			"category_slug": categorySlug,
		},
		Timestamp: time.Now().Unix(),
	}

	posts, err := s.postRepo.SearchPostsWithCategorySlug(ctx, categorySlug, page, pageSize, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("search posts by category slug failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "posts",
		"operation":   "search_by_category_slug",
		"results":     posts,

		"category_slug": categorySlug,
		"page":          page,
		"page_size":     pageSize,
	}, nil
}

func (s *PostToolService) handleFilterPosts(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	// Extract filter parameters
	name := getStringParam(parameters, "name", "")
	description := getStringParam(parameters, "description", "")
	tags := getStringSliceParam(parameters, "tags")
	status := getStringParam(parameters, "status", "")
	userType := getStringParam(parameters, "user_type", "")
	thumbnail := getStringParam(parameters, "thumbnail", "")
	offset := getInt64Param(parameters, "offset", 0)
	limit := getInt64Param(parameters, "limit", 20)
	lat := getFloat32Param(parameters, "lat", 0.0)
	lng := getFloat32Param(parameters, "lng", 0.0)
	radius := getInt64Param(parameters, "radius", 0)
	page := getInt64Param(parameters, "page", 1)
	pageSize := getInt64Param(parameters, "page_size", 20)
	sortBy := getStringParam(parameters, "sort_by", "created_at")
	sortOrder := getStringParam(parameters, "sort_order", "desc")

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "post_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":   "filtering_posts",
			"status": status,
			"tags":   tags,
		},
		Timestamp: time.Now().Unix(),
	}

	posts, err := s.postRepo.SearchPostsWithFilters(
		ctx, name, description, tags, status, userType, thumbnail,
		offset, limit, float64(lat), float64(lng), radius, page, pageSize, sortBy, sortOrder,
	)
	if err != nil {
		return nil, fmt.Errorf("filter posts failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "posts",
		"operation":   "filter",
		"results":     posts,

		"filters": map[string]interface{}{
			"name":        name,
			"description": description,
			"tags":        tags,
			"status":      status,
			"user_type":   userType,
			"thumbnail":   thumbnail,
			"lat":         lat,
			"lng":         lng,
			"radius":      radius,
		},
		"pagination": map[string]interface{}{
			"page":       page,
			"page_size":  pageSize,
			"sort_by":    sortBy,
			"sort_order": sortOrder,
		},
	}, nil
}

func (s *PostToolService) handleAddPost(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	// Extract required parameters
	name := getStringParam(parameters, "name", "")
	description := getStringParam(parameters, "description", "")
	userID := getStringParam(parameters, "user_id", "")

	if name == "" || description == "" || userID == "" {
		return nil, fmt.Errorf("name, description, and user_id are required")
	}

	// Extract optional parameters
	tags := getStringSliceParam(parameters, "tags")
	status := getStringParam(parameters, "status", "active")
	lat := getFloat32Param(parameters, "lat", 0.0)
	lng := getFloat32Param(parameters, "lng", 0.0)
	thumbnail := getStringParam(parameters, "thumbnail", "")

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "post_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step": "adding_post",
			"name": name,
		},
		Timestamp: time.Now().Unix(),
	}

	err := s.postRepo.Add(ctx, "", name, description, tags, status, float64(lat), float64(lng), thumbnail)
	if err != nil {
		return nil, fmt.Errorf("add post failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "posts",
		"operation":   "add",
		"result":      "Post added successfully",
		"name":        name,
	}, nil
}

func (s *PostToolService) handleUpdatePost(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	postID := getStringParam(parameters, "id", "")
	if postID == "" {
		postID = getStringParam(parameters, "post_id", "")
	}
	if postID == "" {
		return nil, fmt.Errorf("id or post_id parameter required")
	}

	// Extract update parameters
	name := getStringParam(parameters, "name", "")
	description := getStringParam(parameters, "description", "")
	userID := getStringParam(parameters, "user_id", "")
	tags := getStringSliceParam(parameters, "tags")
	status := getStringParam(parameters, "status", "")

	if userID == "" {
		return nil, fmt.Errorf("user_id parameter required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "post_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":    "updating_post",
			"post_id": postID,
		},
		Timestamp: time.Now().Unix(),
	}

	err := s.postRepo.UpdatePost(ctx, postID, name, description, tags, userID, status)
	if err != nil {
		return nil, fmt.Errorf("update post failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "posts",
		"operation":   "update",
		"result":      "Post updated successfully",
		"post_id":     postID,
	}, nil
}

func (s *PostToolService) handleRemovePost(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	postID := getStringParam(parameters, "id", "")
	if postID == "" {
		postID = getStringParam(parameters, "post_id", "")
	}
	if postID == "" {
		return nil, fmt.Errorf("id or post_id parameter required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "post_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":    "removing_post",
			"post_id": postID,
		},
		Timestamp: time.Now().Unix(),
	}

	err := s.postRepo.Remove(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("remove post failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "posts",
		"operation":   "remove",
		"result":      "Post removed successfully",
		"post_id":     postID,
	}, nil
}

func (s *PostToolService) handleArchivePost(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	postID := getStringParam(parameters, "id", "")
	if postID == "" {
		postID = getStringParam(parameters, "post_id", "")
	}
	userID := getStringParam(parameters, "user_id", "")

	if postID == "" || userID == "" {
		return nil, fmt.Errorf("post_id and user_id parameters required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "post_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":    "archiving_post",
			"post_id": postID,
		},
		Timestamp: time.Now().Unix(),
	}

	err := s.postRepo.ArchivePost(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("archive post failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "posts",
		"operation":   "archive",
		"post_id":     postID,
	}, nil
}

func (s *PostToolService) handleGetPostsByUser(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	userID := getStringParam(parameters, "user_id", "")
	if userID == "" {
		return nil, fmt.Errorf("user_id parameter required")
	}

	page := getInt64Param(parameters, "page", 1)
	pageSize := getInt64Param(parameters, "page_size", 20)

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "post_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":    "getting_posts_by_user",
			"user_id": userID,
		},
		Timestamp: time.Now().Unix(),
	}

	posts, err := s.postRepo.GetPublicCatalog(ctx, page, pageSize, "created_at", "desc")
	if err != nil {
		return nil, fmt.Errorf("get posts by user failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "posts",
		"operation":   "get_posts_by_user",
		"results":     posts,

		"user_id":   userID,
		"page":      page,
		"page_size": pageSize,
	}, nil
}

func (s *PostToolService) handleGetPostsByLocation(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	lat := getFloat32Param(parameters, "lat", 0.0)
	lng := getFloat32Param(parameters, "lng", 0.0)
	radius := getInt64Param(parameters, "radius", 10)

	if lat == 0.0 && lng == 0.0 {
		return nil, fmt.Errorf("lat and lng parameters required")
	}

	page := getInt64Param(parameters, "page", 1)
	pageSize := getInt64Param(parameters, "page_size", 20)

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "post_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":   "getting_posts_by_location",
			"lat":    lat,
			"lng":    lng,
			"radius": radius,
		},
		Timestamp: time.Now().Unix(),
	}

	posts, totalCount, err := s.postRepo.GetPostsByLocation(ctx, lat, lng, radius, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("get posts by location failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "posts",
		"operation":   "get_posts_by_location",
		"results":     posts,
		"total_count": totalCount,
		"lat":         lat,
		"lng":         lng,
		"radius":      radius,
		"page":        page,
		"page_size":   pageSize,
	}, nil
}

func (s *PostToolService) handleGetPostsByTags(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	tags := getStringSliceParam(parameters, "tags")
	if len(tags) == 0 {
		return nil, fmt.Errorf("tags parameter required")
	}

	page := getInt64Param(parameters, "page", 1)
	pageSize := getInt64Param(parameters, "page_size", 20)

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "post_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step": "getting_posts_by_tags",
			"tags": tags,
		},
		Timestamp: time.Now().Unix(),
	}

	posts, totalCount, err := s.postRepo.GetPostsByTags(ctx, tags, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("get posts by tags failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "posts",
		"operation":   "get_posts_by_tags",
		"results":     posts,
		"total_count": totalCount,
		"tags":        tags,
		"page":        page,
		"page_size":   pageSize,
	}, nil
}

func (s *PostToolService) handleGetPopularPosts(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	page := getInt64Param(parameters, "page", 1)
	pageSize := getInt64Param(parameters, "page_size", 20)

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "post_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step": "getting_popular_posts",
		},
		Timestamp: time.Now().Unix(),
	}

	posts, totalCount, err := s.postRepo.GetPopularPosts(ctx, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("get popular posts failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "posts",
		"operation":   "get_popular_posts",
		"results":     posts,
		"total_count": totalCount,
		"page":        page,
		"page_size":   pageSize,
	}, nil
}

func (s *PostToolService) handleGetRecentPosts(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	page := getInt64Param(parameters, "page", 1)
	pageSize := getInt64Param(parameters, "page_size", 20)

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "post_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step": "getting_recent_posts",
		},
		Timestamp: time.Now().Unix(),
	}

	posts, totalCount, err := s.postRepo.GetRecentPosts(ctx, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("get recent posts failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "posts",
		"operation":   "get_recent_posts",
		"results":     posts,
		"total_count": totalCount,
		"page":        page,
		"page_size":   pageSize,
	}, nil
}

func (s *PostToolService) handleGetTrendingPosts(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	page := getInt64Param(parameters, "page", 1)
	pageSize := getInt64Param(parameters, "page_size", 20)

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "post_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step": "getting_trending_posts",
		},
		Timestamp: time.Now().Unix(),
	}

	posts, totalCount, err := s.postRepo.GetTrendingPosts(ctx, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("get trending posts failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "posts",
		"operation":   "get_trending_posts",
		"results":     posts,
		"total_count": totalCount,
		"page":        page,
		"page_size":   pageSize,
	}, nil
}

func (s *PostToolService) handleLikePost(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	postID := getStringParam(parameters, "post_id", "")
	userID := getStringParam(parameters, "user_id", "")

	if postID == "" || userID == "" {
		return nil, fmt.Errorf("post_id and user_id parameters required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "post_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":    "liking_post",
			"post_id": postID,
		},
		Timestamp: time.Now().Unix(),
	}

	err := s.postRepo.LikePost(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("like post failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "posts",
		"operation":   "like_post",
		"result":      "Post liked successfully",
		"post_id":     postID,
		"user_id":     userID,
	}, nil
}

func (s *PostToolService) handleUnlikePost(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	postID := getStringParam(parameters, "post_id", "")
	userID := getStringParam(parameters, "user_id", "")

	if postID == "" || userID == "" {
		return nil, fmt.Errorf("post_id and user_id parameters required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "post_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":    "unliking_post",
			"post_id": postID,
		},
		Timestamp: time.Now().Unix(),
	}

	err := s.postRepo.UnlikePost(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("unlike post failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "posts",
		"operation":   "unlike_post",
		"result":      "Post unliked successfully",
		"post_id":     postID,
		"user_id":     userID,
	}, nil
}

func (s *PostToolService) handleSharePost(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	postID := getStringParam(parameters, "post_id", "")
	userID := getStringParam(parameters, "user_id", "")

	if postID == "" || userID == "" {
		return nil, fmt.Errorf("post_id and user_id parameters required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "post_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":    "sharing_post",
			"post_id": postID,
		},
		Timestamp: time.Now().Unix(),
	}

	err := s.postRepo.SharePost(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("share post failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "posts",
		"operation":   "share_post",
		"result":      "Post shared successfully",
		"post_id":     postID,
		"user_id":     userID,
	}, nil
}

func (s *PostToolService) handleReportPost(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	postID := getStringParam(parameters, "post_id", "")
	userID := getStringParam(parameters, "user_id", "")
	reason := getStringParam(parameters, "reason", "")

	if postID == "" || userID == "" || reason == "" {
		return nil, fmt.Errorf("post_id, user_id, and reason parameters required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "post_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":    "reporting_post",
			"post_id": postID,
			"reason":  reason,
		},
		Timestamp: time.Now().Unix(),
	}

	err := s.postRepo.ReportPost(ctx, postID, reason)
	if err != nil {
		return nil, fmt.Errorf("report post failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "posts",
		"operation":   "report_post",
		"result":      "Post reported successfully",
		"post_id":     postID,
		"user_id":     userID,
		"reason":      reason,
	}, nil
}
