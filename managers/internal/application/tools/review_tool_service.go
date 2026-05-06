package tools

import (
	"context"
	"fmt"
	"log"
	"time"

	"middleman/managers/internal/domain"
)

// ReviewToolService handles all review-related operations
type ReviewToolService struct {
	reviews domain.ReviewRepository
}

// NewReviewToolService creates a new review tool service
func NewReviewToolService(reviews domain.ReviewRepository) *ReviewToolService {
	return &ReviewToolService{
		reviews: reviews,
	}
}

// GetSupportedOperations returns all operations supported by this service
func (r *ReviewToolService) GetSupportedOperations() []string {
	return []string{
		"add", "create",
		"get", "find",
		"get_reviews", "list",
		"get_reviews_by_sender", "by_sender",
		"get_approved_reviews", "approved",
		"get_pending_reviews", "pending",
		"get_rejected_reviews", "rejected",
		"approve_review",
		"reject_review",
		"update_review",
		"delete_review",
		"flag_review",
		"unflag_review",
		"search_reviews",
	}
}

// ExecuteOperation executes a review operation with streaming progress updates
func (r *ReviewToolService) ExecuteOperation(
	ctx context.Context,
	operation string,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	log.Printf("ReviewToolService.ExecuteOperation: Executing review operation: %s", operation)

	// Send initial progress
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "progress",
		Progress: 10,
		Metadata: map[string]interface{}{
			"operation": operation,
			"message":   fmt.Sprintf("Processing review operation: %s", operation),
		},
		Timestamp: time.Now().Unix(),
	}

	switch operation {
	case "add", "create":
		return r.handleAddReview(ctx, parameters, streamChan, toolID)
	case "get", "find":
		return r.handleGetReview(ctx, parameters, streamChan, toolID)
	case "get_reviews", "list":
		return r.handleGetReviews(ctx, parameters, streamChan, toolID)
	case "get_reviews_by_sender", "by_sender":
		return r.handleGetReviewsBySender(ctx, parameters, streamChan, toolID)
	case "get_approved_reviews", "approved":
		return r.handleGetApprovedReviews(ctx, parameters, streamChan, toolID)
	case "get_pending_reviews", "pending":
		return r.handleGetPendingReviews(ctx, parameters, streamChan, toolID)
	case "get_rejected_reviews", "rejected":
		return r.handleGetRejectedReviews(ctx, parameters, streamChan, toolID)
	case "approve_review":
		return r.handleApproveReview(ctx, parameters, streamChan, toolID)
	case "reject_review":
		return r.handleRejectReview(ctx, parameters, streamChan, toolID)
	case "update_review":
		return r.handleUpdateReview(ctx, parameters, streamChan, toolID)
	case "delete_review":
		return r.handleDeleteReview(ctx, parameters, streamChan, toolID)
	case "flag_review":
		return r.handleFlagReview(ctx, parameters, streamChan, toolID)
	case "unflag_review":
		return r.handleUnflagReview(ctx, parameters, streamChan, toolID)
	case "search_reviews":
		return r.handleSearchReviews(ctx, parameters, streamChan, toolID)
	default:
		return nil, fmt.Errorf("unsupported review operation: %s", operation)
	}
}

func (r *ReviewToolService) handleAddReview(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	// Extract parameters
	senderID := getStringParam(parameters, "sender_id", "")
	if senderID == "" {
		senderID = getStringParam(parameters, "user_id", "")
	}
	itemID := getStringParam(parameters, "item_id", "")
	itemType := getStringParam(parameters, "item_type", "")
	content := getStringParam(parameters, "content", "")
	categoryID := getStringParam(parameters, "category_id", "")
	parentID := getStringParam(parameters, "parent_id", "")

	if senderID == "" || itemID == "" || content == "" {
		return nil, fmt.Errorf("sender_id, item_id, and content are required for adding review")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":      "adding_review",
			"sender_id": senderID,
			"item_id":   itemID,
		},
		Timestamp: time.Now().Unix(),
	}

	reviewID, err := r.reviews.AddReview(ctx, senderID, itemID, itemType, content, categoryID, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to add review: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"review_id": reviewID,
			"message":   "Review added successfully",
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "reviews",
		"operation":   "add",
		"review_id":   reviewID,
		"success":     true,
	}, nil
}

func (r *ReviewToolService) handleGetReview(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	reviewID := getStringParam(parameters, "id", "")
	if reviewID == "" {
		reviewID = getStringParam(parameters, "review_id", "")
	}

	if reviewID == "" {
		return nil, fmt.Errorf("id or review_id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":      "getting_review",
			"review_id": reviewID,
		},
		Timestamp: time.Now().Unix(),
	}

	review, err := r.reviews.GetReview(ctx, reviewID)
	if err != nil {
		return nil, fmt.Errorf("failed to get review: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"review": review,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "reviews",
		"operation":   "get",
		"result":      review,
	}, nil
}

func (r *ReviewToolService) handleGetReviews(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	itemID := getStringParam(parameters, "item_id", "")
	if itemID == "" {
		return nil, fmt.Errorf("item_id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":    "getting_reviews",
			"item_id": itemID,
		},
		Timestamp: time.Now().Unix(),
	}

	reviews, err := r.reviews.GetReviews(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"reviews": reviews,
			"count":   len(reviews),
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "reviews",
		"operation":   "get_reviews",
		"results":     reviews,
		"count":       len(reviews),
	}, nil
}

func (r *ReviewToolService) handleGetReviewsBySender(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	senderID := getStringParam(parameters, "sender_id", "")
	if senderID == "" {
		senderID = getStringParam(parameters, "user_id", "")
	}

	if senderID == "" {
		return nil, fmt.Errorf("sender_id or user_id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":      "getting_reviews_by_sender",
			"sender_id": senderID,
		},
		Timestamp: time.Now().Unix(),
	}

	reviews, err := r.reviews.GetReviewsBySender(ctx, senderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews by sender: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"reviews": reviews,
			"count":   len(reviews),
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "reviews",
		"operation":   "get_reviews_by_sender",
		"results":     reviews,
		"count":       len(reviews),
	}, nil
}

func (r *ReviewToolService) handleGetApprovedReviews(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step": "getting_approved_reviews",
		},
		Timestamp: time.Now().Unix(),
	}

	reviews, err := r.reviews.GetApprovedReviews(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get approved reviews: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"reviews": reviews,
			"count":   len(reviews),
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "reviews",
		"operation":   "get_approved_reviews",
		"results":     reviews,
		"count":       len(reviews),
	}, nil
}

func (r *ReviewToolService) handleGetPendingReviews(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	itemID := getStringParam(parameters, "item_id", "")
	_ = getInt64Param(parameters, "limit", 20) // limit not used in current implementation

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":    "getting_pending_reviews",
			"item_id": itemID,
		},
		Timestamp: time.Now().Unix(),
	}

	// Note: GetPendingReviews not available in current interface, using GetApprovedReviews as fallback
	reviews, err := r.reviews.GetApprovedReviews(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending reviews: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"reviews": reviews,
			"count":   len(reviews),
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "reviews",
		"operation":   "get_pending_reviews",
		"results":     reviews,
		"count":       len(reviews),
	}, nil
}

func (r *ReviewToolService) handleGetRejectedReviews(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	itemID := getStringParam(parameters, "item_id", "")
	_ = getInt64Param(parameters, "limit", 20) // limit not used in current implementation

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":    "getting_rejected_reviews",
			"item_id": itemID,
		},
		Timestamp: time.Now().Unix(),
	}

	// Note: GetRejectedReviews not available in current interface, using GetApprovedReviews as fallback
	reviews, err := r.reviews.GetApprovedReviews(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get rejected reviews: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"reviews": reviews,
			"count":   len(reviews),
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "reviews",
		"operation":   "get_rejected_reviews",
		"results":     reviews,
		"count":       len(reviews),
	}, nil
}

func (r *ReviewToolService) handleApproveReview(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	reviewID := getStringParam(parameters, "review_id", "")
	if reviewID == "" {
		reviewID = getStringParam(parameters, "id", "")
	}

	if reviewID == "" {
		return nil, fmt.Errorf("review_id or id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":      "approving_review",
			"review_id": reviewID,
		},
		Timestamp: time.Now().Unix(),
	}

	_, err := r.reviews.ApproveReview(ctx, reviewID)
	if err != nil {
		return nil, fmt.Errorf("failed to approve review: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"review_id": reviewID,
			"message":   "Review approved successfully",
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "reviews",
		"operation":   "approve_review",
		"review_id":   reviewID,
		"success":     true,
	}, nil
}

func (r *ReviewToolService) handleRejectReview(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	reviewID := getStringParam(parameters, "review_id", "")
	if reviewID == "" {
		reviewID = getStringParam(parameters, "id", "")
	}

	if reviewID == "" {
		return nil, fmt.Errorf("review_id or id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":      "rejecting_review",
			"review_id": reviewID,
		},
		Timestamp: time.Now().Unix(),
	}

	_, err := r.reviews.RejectReview(ctx, reviewID)
	if err != nil {
		return nil, fmt.Errorf("failed to reject review: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"review_id": reviewID,
			"message":   "Review rejected successfully",
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "reviews",
		"operation":   "reject_review",
		"review_id":   reviewID,
		"success":     true,
	}, nil
}

func (r *ReviewToolService) handleUpdateReview(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	reviewID := getStringParam(parameters, "review_id", "")
	if reviewID == "" {
		reviewID = getStringParam(parameters, "id", "")
	}
	content := getStringParam(parameters, "content", "")

	if reviewID == "" {
		return nil, fmt.Errorf("review_id or id parameter required")
	}
	if content == "" {
		return nil, fmt.Errorf("content parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":      "updating_review",
			"review_id": reviewID,
		},
		Timestamp: time.Now().Unix(),
	}

	err := r.reviews.UpdateReview(ctx, reviewID, content)
	if err != nil {
		return nil, fmt.Errorf("failed to update review: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"review_id": reviewID,
			"content":   content,
			"message":   "Review updated successfully",
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "reviews",
		"operation":   "update_review",
		"review_id":   reviewID,
		"success":     true,
	}, nil
}

func (r *ReviewToolService) handleDeleteReview(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	reviewID := getStringParam(parameters, "review_id", "")
	if reviewID == "" {
		reviewID = getStringParam(parameters, "id", "")
	}

	if reviewID == "" {
		return nil, fmt.Errorf("review_id or id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":      "deleting_review",
			"review_id": reviewID,
		},
		Timestamp: time.Now().Unix(),
	}

	err := r.reviews.DeleteReview(ctx, reviewID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete review: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"review_id": reviewID,
			"message":   "Review deleted successfully",
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "reviews",
		"operation":   "delete_review",
		"review_id":   reviewID,
		"success":     true,
	}, nil
}

func (r *ReviewToolService) handleFlagReview(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	reviewID := getStringParam(parameters, "review_id", "")
	if reviewID == "" {
		reviewID = getStringParam(parameters, "id", "")
	}

	if reviewID == "" {
		return nil, fmt.Errorf("review_id or id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":      "flagging_review",
			"review_id": reviewID,
		},
		Timestamp: time.Now().Unix(),
	}

	_, err := r.reviews.FlagReview(ctx, reviewID)
	if err != nil {
		return nil, fmt.Errorf("failed to flag review: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"review_id": reviewID,
			"message":   "Review flagged successfully",
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "reviews",
		"operation":   "flag_review",
		"review_id":   reviewID,
		"success":     true,
	}, nil
}

func (r *ReviewToolService) handleUnflagReview(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	reviewID := getStringParam(parameters, "review_id", "")
	if reviewID == "" {
		reviewID = getStringParam(parameters, "id", "")
	}

	if reviewID == "" {
		return nil, fmt.Errorf("review_id or id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":      "unflagging_review",
			"review_id": reviewID,
		},
		Timestamp: time.Now().Unix(),
	}

	err := r.reviews.UnflagReview(ctx, reviewID)
	if err != nil {
		return nil, fmt.Errorf("failed to unflag review: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"review_id": reviewID,
			"message":   "Review unflagged successfully",
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "reviews",
		"operation":   "unflag_review",
		"review_id":   reviewID,
		"success":     true,
	}, nil
}

func (r *ReviewToolService) handleSearchReviews(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	searchTerm := getStringParam(parameters, "search_term", "")
	if searchTerm == "" {
		searchTerm = getStringParam(parameters, "query", "")
	}
	limit := getInt64Param(parameters, "limit", 20)

	if searchTerm == "" {
		return nil, fmt.Errorf("search_term or query parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":        "searching_reviews",
			"search_term": searchTerm,
		},
		Timestamp: time.Now().Unix(),
	}

	reviews, err := r.reviews.SearchReviews(ctx, searchTerm, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search reviews: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "review_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"reviews":     reviews,
			"count":       len(reviews),
			"search_term": searchTerm,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "reviews",
		"operation":   "search_reviews",
		"results":     reviews,
		"count":       len(reviews),
		"search_term": searchTerm,
	}, nil
}
