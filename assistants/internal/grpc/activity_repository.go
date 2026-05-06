package grpc

import (
	"context"
	"fmt"
	"middleman/assistants/internal/domain"
	"middleman/assistants/internal/models"
	"middleman/internal/auth"
	"middleman/internal/rpc"

	"google.golang.org/grpc"
)

// ActivityRepository calls the remote activity service (gRPC) as a fallback.
type ActivityRepository struct {
	endpoint string
	auth     *auth.Auth // Optional auth for authenticated calls
}

var _ domain.ActivityRepository = (*ActivityRepository)(nil)

// NewActivityRepository instantiates the gRPC-based fallback repo.
func NewActivityRepository(endpoint string) ActivityRepository {
	return ActivityRepository{
		endpoint: endpoint,
		auth:     nil, // No auth by default for backwards compatibility
	}
}

// NewActivityRepositoryWithAuth instantiates the gRPC-based repo with auth.
func NewActivityRepositoryWithAuth(endpoint string, authProvider *auth.Auth) ActivityRepository {
	return ActivityRepository{
		endpoint: endpoint,
		auth:     authProvider,
	}
}

// CreateUserActivity creates a new activity for a user
func (r ActivityRepository) CreateUserActivity(ctx context.Context, userID string) (string, error) {
	// TODO: Implement when activitypb is available
	return fmt.Sprintf("mock_activity_%s", userID), nil
}

// DeleteAllUserActivities removes all activities via the activity microservice
func (r ActivityRepository) DeleteAllUserActivities(ctx context.Context, activityID string) error {
	// TODO: Implement when activitypb is available
	return fmt.Errorf("DeleteAllUserActivities not implemented - activity service protobuf not available")
}

// ArchiveUserActivity archives an activity
func (r ActivityRepository) ArchiveUserActivity(ctx context.Context, activityID, reason string) error {
	// TODO: Implement when activitypb is available
	return fmt.Errorf("ArchiveUserActivity not implemented - activity service protobuf not available")
}

// RestoreUserActivity restores an archived activity
func (r ActivityRepository) RestoreUserActivity(ctx context.Context, activityID, reason string) error {
	// TODO: Implement when activitypb is available
	return fmt.Errorf("RestoreUserActivity not implemented - activity service protobuf not available")
}

// GetActivityByUserID finds an activity by user ID via the activity microservice
func (r ActivityRepository) GetActivityByUserID(ctx context.Context, userID string) (*models.Activity, error) {
	// TODO: Implement when activitypb is available
	return &models.Activity{
		ID:       fmt.Sprintf("mock_activity_%s", userID),
		UserID:   userID,
		Archived: false,
	}, nil
}

// AddUserInteraction adds a like or dislike interaction via the activity microservice
func (r ActivityRepository) AddUserInteraction(ctx context.Context, interactionID, activityID, itemID, itemType, actionType string) error {
	// TODO: Implement when activitypb is available
	return fmt.Errorf("AddUserInteraction not implemented - activity service protobuf not available")
}

// UpdateUserInteraction updates a like or dislike interaction via the activity microservice
func (r ActivityRepository) UpdateUserInteraction(ctx context.Context, interactionID string, actionType string) error {
	// TODO: Implement when activitypb is available
	return fmt.Errorf("UpdateUserInteraction not implemented - activity service protobuf not available")
}

// DeleteUserInteraction removes an interaction via the activity microservice
func (r ActivityRepository) DeleteUserInteraction(ctx context.Context, interactionID string) error {
	// TODO: Implement when activitypb is available
	return fmt.Errorf("DeleteUserInteraction not implemented - activity service protobuf not available")
}

// GetInteractionByID retrieves an interaction by ID from the activity microservice
func (r ActivityRepository) GetInteractionByID(ctx context.Context, interactionID string) (*models.Interaction, error) {
	// TODO: Implement when activitypb is available
	return &models.Interaction{
		ID:         interactionID,
		ActivityID: "mock_activity",
		ItemID:     "mock_item",
		ItemType:   "mock_type",
		ActionType: "like",
	}, nil
}

// GetAllActivityInteractions retrieves all interactions for an activity from the activity microservice
func (r ActivityRepository) GetAllActivityInteractions(ctx context.Context, activityID string) ([]*models.Interaction, error) {
	// TODO: Implement when activitypb is available
	return []*models.Interaction{
		{
			ID:         "mock_interaction_1",
			ActivityID: activityID,
			ItemID:     "mock_item_1",
			ItemType:   "product",
			ActionType: "like",
		},
	}, nil
}

// GetMostLikedItems retrieves the most liked items of a specific type
func (r ActivityRepository) GetMostLikedItems(ctx context.Context, itemType string, limit int64) ([]*models.MostReactionResult, error) {
	// TODO: Implement when activitypb is available
	return []*models.MostReactionResult{
		{
			ItemID:   "popular_item_1",
			ItemType: itemType,
			Action:   "like",
			Count:    100,
		},
	}, nil
}

// GetMostDislikedItems retrieves the most disliked items of a specific type
func (r ActivityRepository) GetMostDislikedItems(ctx context.Context, itemType string, limit int64) ([]*models.MostReactionResult, error) {
	// TODO: Implement when activitypb is available
	return []*models.MostReactionResult{
		{
			ItemID:   "unpopular_item_1",
			ItemType: itemType,
			Action:   "dislike",
			Count:    50,
		},
	}, nil
}

// dial establishes a gRPC connection to the activity service
func (r ActivityRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

func (r ActivityRepository) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}
