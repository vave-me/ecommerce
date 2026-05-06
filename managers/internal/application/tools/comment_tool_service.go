package tools

import (
	"context"
	"fmt"
	"log"
	"time"

	"middleman/managers/internal/domain"
)

// CommentToolService handles all comment-related operations with streaming execution
type CommentToolService struct {
	commentRepository domain.CommentRepository
}

// NewCommentToolService creates a new comment tool service instance
func NewCommentToolService(commentRepository domain.CommentRepository) *CommentToolService {
	return &CommentToolService{
		commentRepository: commentRepository,
	}
}

// ExecuteOperation executes comment operations with streaming progress updates
func (s *CommentToolService) ExecuteOperation(
	ctx context.Context,
	operation string,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (*ToolOperationResult, error) {
	startTime := time.Now()

	// Send initial progress
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "comment_operation",
		Status:   "progress",
		Progress: 25.0,
		Metadata: map[string]interface{}{
			"step":      "executing_comment_operation",
			"operation": operation,
		},
		Timestamp: time.Now().Unix(),
	}

	var result interface{}
	var err error

	switch operation {
	case "add", "create":
		result, err = s.handleAddComment(ctx, parameters, streamChan, toolID)
	case "get", "find":
		result, err = s.handleGetComment(ctx, parameters, streamChan, toolID)
	case "get_comments", "list":
		result, err = s.handleGetComments(ctx, parameters, streamChan, toolID)
	case "approve":
		result, err = s.handleApproveComment(ctx, parameters, streamChan, toolID)
	case "reject":
		result, err = s.handleRejectComment(ctx, parameters, streamChan, toolID)
	case "flag":
		result, err = s.handleFlagComment(ctx, parameters, streamChan, toolID)
	case "edit":
		result, err = s.handleEditComment(ctx, parameters, streamChan, toolID)
	case "delete":
		result, err = s.handleDeleteComment(ctx, parameters, streamChan, toolID)
	default:
		err = fmt.Errorf("unsupported comment operation: %s", operation)
	}

	duration := time.Since(startTime)

	if err != nil {
		// Send error update
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "comment_operation",
			Status:   "error",
			Error:    err.Error(),
			Metadata: map[string]interface{}{
				"operation": operation,
			},
			Timestamp: time.Now().Unix(),
		}
		return &ToolOperationResult{
			EntityType: "comments",
			Operation:  operation,
			Success:    false,
			Error:      err.Error(),
			Duration:   duration,
		}, err
	}

	// Send completion update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "comment_operation",
		Status:   "completed",
		Progress: 100,
		Result:   result,
		Metadata: map[string]interface{}{
			"operation": operation,
			"success":   true,
		},
		Timestamp: time.Now().Unix(),
	}

	return &ToolOperationResult{
		EntityType: "comments",
		Operation:  operation,
		Success:    true,
		Result:     result,
		Duration:   duration,
	}, nil
}

// handleAddComment adds a new comment
func (s *CommentToolService) handleAddComment(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	senderID := getCommentStringParam(parameters, "sender_id", "")
	if senderID == "" {
		senderID = getCommentStringParam(parameters, "user_id", "")
	}
	itemID := getCommentStringParam(parameters, "item_id", "")
	itemType := getCommentStringParam(parameters, "item_type", "")
	content := getCommentStringParam(parameters, "content", "")
	categoryID := getCommentStringParam(parameters, "category_id", "")
	parentID := getCommentStringParam(parameters, "parent_id", "")

	if senderID == "" || itemID == "" || content == "" {
		return nil, fmt.Errorf("sender_id, item_id, and content are required for adding comment")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "comment_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":      "adding_comment",
			"sender_id": senderID,
			"item_id":   itemID,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("CommentToolService: Adding comment from user %s to item %s", senderID, itemID)
	resp, err := s.commentRepository.CreateNewComment(ctx, itemID, itemType, content, categoryID, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to add comment: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "comments",
		"operation":   "add",
		"comment_id":  resp.ID,
		"success":     true,
		"result":      resp,
	}, nil
}

// handleGetComment retrieves a specific comment
func (s *CommentToolService) handleGetComment(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	commentID := getCommentStringParam(parameters, "id", "")
	if commentID == "" {
		commentID = getCommentStringParam(parameters, "comment_id", "")
	}

	if commentID == "" {
		return nil, fmt.Errorf("id or comment_id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "comment_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":       "getting_comment",
			"comment_id": commentID,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("CommentToolService: Getting comment with ID: %s", commentID)
	resp, err := s.commentRepository.GetCommentByID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "comments",
		"operation":   "get",
		"result":      resp.Comment,
	}, nil
}

// handleGetComments retrieves comments for an item
func (s *CommentToolService) handleGetComments(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	itemID := getCommentStringParam(parameters, "item_id", "")
	if itemID == "" {
		return nil, fmt.Errorf("item_id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "comment_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":    "getting_comments",
			"item_id": itemID,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("CommentToolService: Getting comments for item: %s", itemID)
	resp, err := s.commentRepository.GetAllCommentsForItem(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "comments",
		"operation":   "get_comments",
		"results":     resp.Comments,
		"count":       len(resp.Comments),
	}, nil
}

// handleApproveComment approves a comment
func (s *CommentToolService) handleApproveComment(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	commentID := getCommentStringParam(parameters, "id", "")
	if commentID == "" {
		commentID = getCommentStringParam(parameters, "comment_id", "")
	}

	if commentID == "" {
		return nil, fmt.Errorf("id or comment_id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "comment_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":       "approving_comment",
			"comment_id": commentID,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("CommentToolService: Approving comment with ID: %s", commentID)
	resp, err := s.commentRepository.ApproveCommentByID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to approve comment: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "comments",
		"operation":   "approve",
		"result":      resp,
		"success":     true,
	}, nil
}

// handleRejectComment rejects a comment
func (s *CommentToolService) handleRejectComment(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	commentID := getCommentStringParam(parameters, "id", "")
	if commentID == "" {
		commentID = getCommentStringParam(parameters, "comment_id", "")
	}

	if commentID == "" {
		return nil, fmt.Errorf("id or comment_id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "comment_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":       "rejecting_comment",
			"comment_id": commentID,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("CommentToolService: Rejecting comment with ID: %s", commentID)
	resp, err := s.commentRepository.RejectCommentByID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to reject comment: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "comments",
		"operation":   "reject",
		"result":      resp,
		"success":     true,
	}, nil
}

// handleFlagComment flags a comment
func (s *CommentToolService) handleFlagComment(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	commentID := getCommentStringParam(parameters, "id", "")
	if commentID == "" {
		commentID = getCommentStringParam(parameters, "comment_id", "")
	}

	if commentID == "" {
		return nil, fmt.Errorf("id or comment_id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "comment_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":       "flagging_comment",
			"comment_id": commentID,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("CommentToolService: Flagging comment with ID: %s", commentID)
	resp, err := s.commentRepository.FlagCommentAsInappropriate(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to flag comment: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "comments",
		"operation":   "flag",
		"result":      resp,
		"success":     true,
	}, nil
}

// handleEditComment edits a comment
func (s *CommentToolService) handleEditComment(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	commentID := getCommentStringParam(parameters, "id", "")
	if commentID == "" {
		commentID = getCommentStringParam(parameters, "comment_id", "")
	}
	content := getCommentStringParam(parameters, "content", "")

	if commentID == "" || content == "" {
		return nil, fmt.Errorf("id/comment_id and content parameters required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "comment_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":       "editing_comment",
			"comment_id": commentID,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("CommentToolService: Editing comment with ID: %s", commentID)
	resp, err := s.commentRepository.EditCommentContent(ctx, commentID, content)
	if err != nil {
		return nil, fmt.Errorf("failed to edit comment: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "comments",
		"operation":   "edit",
		"result":      resp,
		"success":     true,
	}, nil
}

// handleDeleteComment deletes a comment
func (s *CommentToolService) handleDeleteComment(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	commentID := getCommentStringParam(parameters, "id", "")
	if commentID == "" {
		commentID = getCommentStringParam(parameters, "comment_id", "")
	}

	if commentID == "" {
		return nil, fmt.Errorf("id or comment_id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "comment_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":       "deleting_comment",
			"comment_id": commentID,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("CommentToolService: Deleting comment with ID: %s", commentID)
	resp, err := s.commentRepository.DeleteCommentByID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete comment: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "comments",
		"operation":   "delete",
		"result":      resp,
		"success":     true,
	}, nil
}

// Helper function for parameter extraction
func getCommentStringParam(params map[string]interface{}, key, defaultValue string) string {
	if val, ok := params[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}
