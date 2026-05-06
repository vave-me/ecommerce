package tools

import (
	"context"
	"fmt"

	"time"

	"middleman/managers/internal/domain"
)

// ActivityToolService handles activity tracking operations
type ActivityToolService struct {
	activityRepo domain.ActivityRepository
	config       *ToolStreamConfig
}

// NewActivityToolService creates a new activity tool service
func NewActivityToolService(activityRepo domain.ActivityRepository) *ActivityToolService {
	return &ActivityToolService{
		activityRepo: activityRepo,
		config: &ToolStreamConfig{
			BufferSize:       100,
			ProgressInterval: 500 * time.Millisecond,
			EnableMetrics:    true,
			MaxRetries:       3,
		},
	}
}

// ExecuteOperation routes activity operations to appropriate handlers
func (s *ActivityToolService) ExecuteOperation(
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
			ToolName: "activity_operation",
			Status:   "started",
			Progress: 0,
			Metadata: map[string]interface{}{
				"operation": operation,
				"service":   "ActivityToolService",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	var result interface{}
	var err error

	switch operation {
	case "find", "get_interaction":
		result, err = s.findInteraction(ctx, parameters, streamChan, toolID)
	case "create_activity":
		result, err = s.createActivity(ctx, parameters, streamChan, toolID)
	case "find_activity":
		result, err = s.findActivity(ctx, parameters, streamChan, toolID)
	case "add_like":
		result, err = s.addLike(ctx, parameters, streamChan, toolID)
	case "add_dislike":
		result, err = s.addDislike(ctx, parameters, streamChan, toolID)
	case "update_interaction":
		result, err = s.updateInteraction(ctx, parameters, streamChan, toolID)
	case "remove_interaction":
		result, err = s.removeInteraction(ctx, parameters, streamChan, toolID)
	case "get_interactions":
		result, err = s.getInteractions(ctx, parameters, streamChan, toolID)
	case "get_most_liked":
		result, err = s.getMostLiked(ctx, parameters, streamChan, toolID)
	case "get_most_disliked":
		result, err = s.getMostDisliked(ctx, parameters, streamChan, toolID)
	case "archive_activity":
		result, err = s.archiveActivity(ctx, parameters, streamChan, toolID)
	case "restore_activity":
		result, err = s.restoreActivity(ctx, parameters, streamChan, toolID)
	default:
		err = fmt.Errorf("unsupported activity operation: %s", operation)
	}

	// Send completion status
	if streamChan != nil {
		status := "completed"
		if err != nil {
			status = "error"
		}

		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "activity_operation",
			Status:   status,
			Progress: 100,
			Result:   result,
			Error:    s.getErrorString(err),
			Metadata: map[string]interface{}{
				"operation": operation,
				"service":   "ActivityToolService",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	return result, err
}

func (s *ActivityToolService) createActivity(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	userID := s.getStringParam(params, "user_id", "")
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	s.sendProgress(streamChan, toolID, "creating_activity", 50)

	activityID, err := s.activityRepo.CreateActivity(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create activity: %w", err)
	}

	return map[string]interface{}{
		"activity_id": activityID,
		"user_id":     userID,
		"created":     true,
	}, nil
}

func (s *ActivityToolService) findActivity(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	userID := s.getStringParam(params, "user_id", "")
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	s.sendProgress(streamChan, toolID, "finding_activity", 50)

	activity, err := s.activityRepo.FindActivityId(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find activity: %w", err)
	}

	return map[string]interface{}{
		"activity": activity,
		"user_id":  userID,
	}, nil
}

func (s *ActivityToolService) findInteraction(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	interactionID := s.getStringParam(params, "interaction_id", "")
	if interactionID == "" {
		interactionID = s.getStringParam(params, "id", "")
	}
	if interactionID == "" {
		return nil, fmt.Errorf("interaction_id is required")
	}

	s.sendProgress(streamChan, toolID, "finding_interaction", 50)

	interaction, err := s.activityRepo.Find(ctx, interactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find interaction: %w", err)
	}

	return map[string]interface{}{
		"interaction":    interaction,
		"interaction_id": interactionID,
	}, nil
}

func (s *ActivityToolService) addLike(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return s.addLikeOrDislike(ctx, params, streamChan, toolID, "like")
}

func (s *ActivityToolService) addDislike(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return s.addLikeOrDislike(ctx, params, streamChan, toolID, "dislike")
}

func (s *ActivityToolService) addLikeOrDislike(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string, actionType string) (interface{}, error) {
	interactionID := s.getStringParam(params, "interaction_id", "")
	activityID := s.getStringParam(params, "activity_id", "")
	itemID := s.getStringParam(params, "item_id", "")
	itemType := s.getStringParam(params, "item_type", "")

	if interactionID == "" || activityID == "" || itemID == "" || itemType == "" {
		return nil, fmt.Errorf("interaction_id, activity_id, item_id, and item_type are required")
	}

	s.sendProgress(streamChan, toolID, fmt.Sprintf("adding_%s", actionType), 50)

	err := s.activityRepo.AddLikeOrDislike(ctx, interactionID, activityID, itemID, itemType, actionType)
	if err != nil {
		return nil, fmt.Errorf("failed to add %s: %w", actionType, err)
	}

	return map[string]interface{}{
		"interaction_id": interactionID,
		"activity_id":    activityID,
		"item_id":        itemID,
		"item_type":      itemType,
		"action_type":    actionType,
		"success":        true,
	}, nil
}

func (s *ActivityToolService) updateInteraction(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	interactionID := s.getStringParam(params, "interaction_id", "")
	actionType := s.getStringParam(params, "action_type", "")

	if interactionID == "" || actionType == "" {
		return nil, fmt.Errorf("interaction_id and action_type are required")
	}

	s.sendProgress(streamChan, toolID, "updating_interaction", 50)

	err := s.activityRepo.UpdateLikeOrDislike(ctx, interactionID, actionType)
	if err != nil {
		return nil, fmt.Errorf("failed to update interaction: %w", err)
	}

	return map[string]interface{}{
		"interaction_id": interactionID,
		"action_type":    actionType,
		"updated":        true,
	}, nil
}

func (s *ActivityToolService) removeInteraction(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	interactionID := s.getStringParam(params, "interaction_id", "")
	if interactionID == "" {
		interactionID = s.getStringParam(params, "id", "")
	}
	if interactionID == "" {
		return nil, fmt.Errorf("interaction_id is required")
	}

	s.sendProgress(streamChan, toolID, "removing_interaction", 50)

	err := s.activityRepo.RemoveInteraction(ctx, interactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to remove interaction: %w", err)
	}

	return map[string]interface{}{
		"interaction_id": interactionID,
		"removed":        true,
	}, nil
}

func (s *ActivityToolService) getInteractions(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	activityID := s.getStringParam(params, "activity_id", "")
	if activityID == "" {
		return nil, fmt.Errorf("activity_id is required")
	}

	s.sendProgress(streamChan, toolID, "getting_interactions", 50)

	interactions, err := s.activityRepo.AllInteraction(ctx, activityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get interactions: %w", err)
	}

	return map[string]interface{}{
		"interactions": interactions,
		"activity_id":  activityID,
		"count":        len(interactions),
	}, nil
}

func (s *ActivityToolService) getMostLiked(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	itemType := s.getStringParam(params, "item_type", "")
	limit := s.getInt64Param(params, "limit", 10)

	if itemType == "" {
		return nil, fmt.Errorf("item_type is required")
	}

	s.sendProgress(streamChan, toolID, "getting_most_liked", 50)

	results, err := s.activityRepo.GetMostLiked(ctx, itemType, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get most liked: %w", err)
	}

	return map[string]interface{}{
		"most_liked": results,
		"item_type":  itemType,
		"limit":      limit,
		"count":      len(results),
	}, nil
}

func (s *ActivityToolService) getMostDisliked(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	itemType := s.getStringParam(params, "item_type", "")
	limit := s.getInt64Param(params, "limit", 10)

	if itemType == "" {
		return nil, fmt.Errorf("item_type is required")
	}

	s.sendProgress(streamChan, toolID, "getting_most_disliked", 50)

	results, err := s.activityRepo.GetMostDisliked(ctx, itemType, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get most disliked: %w", err)
	}

	return map[string]interface{}{
		"most_disliked": results,
		"item_type":     itemType,
		"limit":         limit,
		"count":         len(results),
	}, nil
}

func (s *ActivityToolService) archiveActivity(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	activityID := s.getStringParam(params, "activity_id", "")
	reason := s.getStringParam(params, "reason", "")

	if activityID == "" {
		return nil, fmt.Errorf("activity_id is required")
	}

	s.sendProgress(streamChan, toolID, "archiving_activity", 50)

	err := s.activityRepo.ArchiveActivity(ctx, activityID, reason)
	if err != nil {
		return nil, fmt.Errorf("failed to archive activity: %w", err)
	}

	return map[string]interface{}{
		"activity_id": activityID,
		"reason":      reason,
		"archived":    true,
	}, nil
}

func (s *ActivityToolService) restoreActivity(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	activityID := s.getStringParam(params, "activity_id", "")
	reason := s.getStringParam(params, "reason", "")

	if activityID == "" {
		return nil, fmt.Errorf("activity_id is required")
	}

	s.sendProgress(streamChan, toolID, "restoring_activity", 50)

	err := s.activityRepo.RestoreActivity(ctx, activityID, reason)
	if err != nil {
		return nil, fmt.Errorf("failed to restore activity: %w", err)
	}

	return map[string]interface{}{
		"activity_id": activityID,
		"reason":      reason,
		"restored":    true,
	}, nil
}

// Helper functions
func (s *ActivityToolService) getErrorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func (s *ActivityToolService) sendProgress(streamChan chan<- ToolExecutionStream, toolID string, step string, progress float64) {
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "activity_operation",
			Status:   "progress",
			Progress: progress,
			Metadata: map[string]interface{}{
				"step": step,
			},
			Timestamp: time.Now().Unix(),
		}
	}
}

func (s *ActivityToolService) getStringParam(params map[string]interface{}, key, defaultValue string) string {
	if val, ok := params[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

func (s *ActivityToolService) getInt64Param(params map[string]interface{}, key string, defaultValue int64) int64 {
	if val, ok := params[key]; ok {
		switch v := val.(type) {
		case int:
			return int64(v)
		case int64:
			return v
		case float64:
			return int64(v)
		}
	}
	return defaultValue
}
