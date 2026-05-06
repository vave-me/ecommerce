package grpc

import (
	"context"
	"fmt"
	"middleman/internal/auth"
	"middleman/internal/rpc"
	"middleman/managers/internal/domain"
	"middleman/managers/internal/models"

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

// CreateActivity creates a new activity for a user
func (r ActivityRepository) CreateActivity(ctx context.Context, userID string) (string, error) {
	// TODO: Implement when activitypb is available
	return fmt.Sprintf("mock_activity_%s", userID), nil
}

// RemoveAllActivities removes all activities via the activity microservice
func (r ActivityRepository) RemoveAllActivities(ctx context.Context, activityID string) error {
	// TODO: Implement when activitypb is available
	return fmt.Errorf("RemoveAllActivities not implemented - activity service protobuf not available")
}

// ArchiveActivity archives an activity
func (r ActivityRepository) ArchiveActivity(ctx context.Context, activityID, reason string) error {
	// TODO: Implement when activitypb is available
	return fmt.Errorf("ArchiveActivity not implemented - activity service protobuf not available")
}

// RestoreActivity restores an archived activity
func (r ActivityRepository) RestoreActivity(ctx context.Context, activityID, reason string) error {
	// TODO: Implement when activitypb is available
	return fmt.Errorf("RestoreActivity not implemented - activity service protobuf not available")
}

// FindActivityId finds an activity by user ID via the activity microservice
func (r ActivityRepository) FindActivityId(ctx context.Context, userID string) (*models.Activity, error) {
	// TODO: Implement when activitypb is available
	return &models.Activity{
		ID:       fmt.Sprintf("mock_activity_%s", userID),
		UserID:   userID,
		Archived: false,
	}, nil
}

// AddLikeOrDislike adds a like or dislike interaction via the activity microservice
func (r ActivityRepository) AddLikeOrDislike(ctx context.Context, interactionID, activityID, itemID, itemType, actionType string) error {
	// TODO: Implement when activitypb is available
	return fmt.Errorf("AddLikeOrDislike not implemented - activity service protobuf not available")
}

// UpdateLikeOrDislike updates a like or dislike interaction via the activity microservice
func (r ActivityRepository) UpdateLikeOrDislike(ctx context.Context, interactionID string, actionType string) error {
	// TODO: Implement when activitypb is available
	return fmt.Errorf("UpdateLikeOrDislike not implemented - activity service protobuf not available")
}

// RemoveInteraction removes an interaction via the activity microservice
func (r ActivityRepository) RemoveInteraction(ctx context.Context, interactionID string) error {
	// TODO: Implement when activitypb is available
	return fmt.Errorf("RemoveInteraction not implemented - activity service protobuf not available")
}

// Find retrieves an interaction by ID from the activity microservice
func (r ActivityRepository) Find(ctx context.Context, interactionID string) (*models.Interaction, error) {
	// TODO: Implement when activitypb is available
	return &models.Interaction{
		ID:         interactionID,
		ActivityID: "mock_activity",
		ItemID:     "mock_item",
		ItemType:   "mock_type",
		ActionType: "like",
	}, nil
}

// AllInteraction retrieves all interactions for an activity from the activity microservice
func (r ActivityRepository) AllInteraction(ctx context.Context, activityID string) ([]*models.Interaction, error) {
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

// GetMostLiked retrieves the most liked items of a specific type
func (r ActivityRepository) GetMostLiked(ctx context.Context, itemType string, limit int64) ([]*models.MostReactionResult, error) {
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

// GetMostDisliked retrieves the most disliked items of a specific type
func (r ActivityRepository) GetMostDisliked(ctx context.Context, itemType string, limit int64) ([]*models.MostReactionResult, error) {
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
