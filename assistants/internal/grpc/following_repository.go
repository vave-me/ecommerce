package grpc

import (
	"context"
	"fmt"
	"log"
	"middleman/assistants/internal/domain"
	"middleman/assistants/internal/models"
	"middleman/following/followingpb"
	"middleman/internal/auth"
	"middleman/internal/rpc"
	"time"

	"google.golang.org/grpc"
)

// FollowingRepository calls the remote following service (gRPC) as a fallback.
type FollowingRepository struct {
	endpoint string
	auth     *auth.Auth // Optional auth for authenticated calls
}

var _ domain.FollowingRepository = (*FollowingRepository)(nil)

func NewFollowingRepository(endpoint string, authInstance *auth.Auth) FollowingRepository {
	return FollowingRepository{
		endpoint: endpoint,
		auth:     authInstance,
	}
}

// CreateNewFollowRelationship creates a new follow relationship
func (r FollowingRepository) CreateNewFollowRelationship(ctx context.Context, userID, followedUserID, followedUserType, content, categoryID, parentID string) (*models.AddFollowResponse, error) {
	log.Printf("[FOLLOWING_GRPC] CreateNewFollowRelationship called for userID=%s, followedUserID=%s", userID, followedUserID)

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[FOLLOWING_GRPC] Failed to connect to following service: %v", err)
		return nil, fmt.Errorf("failed to connect to following service: %w", err)
	}
	defer conn.Close()

	client := followingpb.NewFollowingServiceClient(conn)
	resp, err := client.AddFollow(ctx, &followingpb.AddFollowRequest{
		FollowedUserId:   followedUserID,
		FollowedUserType: followedUserType,
		Content:          content,
		CategoryId:       categoryID,
		ParentId:         parentID,
	})
	if err != nil {
		log.Printf("[FOLLOWING_GRPC] CreateNewFollowRelationship RPC failed: %v", err)
		return nil, fmt.Errorf("AddFollow RPC failed: %w", err)
	}

	log.Printf("[FOLLOWING_GRPC] CreateNewFollowRelationship RPC successful, created follow with ID: %s", resp.GetId())
	return &models.AddFollowResponse{
		ID: resp.GetId(),
	}, nil
}

// ApproveFollowRequestByID approves a follow request
func (r FollowingRepository) ApproveFollowRequestByID(ctx context.Context, id string) (*models.ApproveFollowResponse, error) {
	log.Printf("[FOLLOWING_GRPC] ApproveFollow called for ID: %s", id)

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[FOLLOWING_GRPC] Failed to connect to following service: %v", err)
		return nil, fmt.Errorf("failed to connect to following service: %w", err)
	}
	defer conn.Close()

	client := followingpb.NewFollowingServiceClient(conn)
	resp, err := client.ApproveFollow(ctx, &followingpb.ApproveFollowRequest{
		Id: id,
	})
	if err != nil {
		log.Printf("[FOLLOWING_GRPC] ApproveFollow RPC failed: %v", err)
		return nil, fmt.Errorf("ApproveFollow RPC failed: %w", err)
	}

	var approvedAt time.Time
	if resp.GetApprovedAt() != nil {
		approvedAt = resp.GetApprovedAt().AsTime()
	}

	log.Printf("[FOLLOWING_GRPC] ApproveFollow RPC successful for ID: %s", id)
	return &models.ApproveFollowResponse{
		ID:           resp.GetId(),
		FollowStatus: resp.GetFollowStatus(),
		ApprovedAt:   approvedAt,
	}, nil
}

// RejectFollowRequestByID rejects a follow request
func (r FollowingRepository) RejectFollowRequestByID(ctx context.Context, id string) (*models.RejectFollowResponse, error) {
	log.Printf("[FOLLOWING_GRPC] RejectFollow called for ID: %s", id)

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[FOLLOWING_GRPC] Failed to connect to following service: %v", err)
		return nil, fmt.Errorf("failed to connect to following service: %w", err)
	}
	defer conn.Close()

	client := followingpb.NewFollowingServiceClient(conn)
	resp, err := client.RejectFollow(ctx, &followingpb.RejectFollowRequest{
		Id: id,
	})
	if err != nil {
		log.Printf("[FOLLOWING_GRPC] RejectFollow RPC failed: %v", err)
		return nil, fmt.Errorf("RejectFollow RPC failed: %w", err)
	}

	log.Printf("[FOLLOWING_GRPC] RejectFollow RPC successful for ID: %s", id)
	return &models.RejectFollowResponse{
		ID:           resp.GetId(),
		FollowStatus: resp.GetFollowStatus(),
	}, nil
}

// FlagFollowAsInappropriate flags a follow for review
func (r FollowingRepository) FlagFollowAsInappropriate(ctx context.Context, id string) (*models.FlagFollowResponse, error) {
	log.Printf("[FOLLOWING_GRPC] FlagFollow called for ID: %s", id)

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[FOLLOWING_GRPC] Failed to connect to following service: %v", err)
		return nil, fmt.Errorf("failed to connect to following service: %w", err)
	}
	defer conn.Close()

	client := followingpb.NewFollowingServiceClient(conn)
	resp, err := client.FlagFollow(ctx, &followingpb.FlagFollowRequest{
		Id: id,
	})
	if err != nil {
		log.Printf("[FOLLOWING_GRPC] FlagFollow RPC failed: %v", err)
		return nil, fmt.Errorf("FlagFollow RPC failed: %w", err)
	}

	log.Printf("[FOLLOWING_GRPC] FlagFollow RPC successful for ID: %s", id)
	return &models.FlagFollowResponse{
		ID:      resp.GetId(),
		Flagged: resp.GetFlagged(),
	}, nil
}

// EditFollowDescription edits the content of a follow
func (r FollowingRepository) EditFollowDescription(ctx context.Context, id, content string) (*models.EditFollowResponse, error) {
	log.Printf("[FOLLOWING_GRPC] EditFollow called for ID: %s", id)

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[FOLLOWING_GRPC] Failed to connect to following service: %v", err)
		return nil, fmt.Errorf("failed to connect to following service: %w", err)
	}
	defer conn.Close()

	client := followingpb.NewFollowingServiceClient(conn)
	resp, err := client.EditFollow(ctx, &followingpb.EditFollowRequest{
		Id:      id,
		Content: content,
	})
	if err != nil {
		log.Printf("[FOLLOWING_GRPC] EditFollow RPC failed: %v", err)
		return nil, fmt.Errorf("EditFollow RPC failed: %w", err)
	}

	log.Printf("[FOLLOWING_GRPC] EditFollow RPC successful for ID: %s", id)
	return &models.EditFollowResponse{
		ID:      resp.GetId(),
		Content: resp.GetContent(),
	}, nil
}

// DeleteFollowRelationship removes a follow relationship
func (r FollowingRepository) DeleteFollowRelationship(ctx context.Context, id string) (*models.RemoveFollowResponse, error) {
	log.Printf("[FOLLOWING_GRPC] RemoveFollow called for ID: %s", id)

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[FOLLOWING_GRPC] Failed to connect to following service: %v", err)
		return nil, fmt.Errorf("failed to connect to following service: %w", err)
	}
	defer conn.Close()

	client := followingpb.NewFollowingServiceClient(conn)
	_, err = client.RemoveFollow(ctx, &followingpb.RemoveFollowRequest{
		Id: id,
	})
	if err != nil {
		log.Printf("[FOLLOWING_GRPC] RemoveFollow RPC failed: %v", err)
		return nil, fmt.Errorf("RemoveFollow RPC failed: %w", err)
	}

	log.Printf("[FOLLOWING_GRPC] RemoveFollow RPC successful for ID: %s", id)
	return &models.RemoveFollowResponse{
		Success: true,
		Message: "Follow removed successfully",
	}, nil
}

// GetFollowByID retrieves a specific follow by ID
func (r FollowingRepository) GetFollowByID(ctx context.Context, id string) (*models.GetFollowResponse, error) {
	log.Printf("[FOLLOWING_GRPC] GetFollow called for ID: %s", id)

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[FOLLOWING_GRPC] Failed to connect to following service: %v", err)
		return nil, fmt.Errorf("failed to connect to following service: %w", err)
	}
	defer conn.Close()

	client := followingpb.NewFollowingServiceClient(conn)
	resp, err := client.GetFollow(ctx, &followingpb.GetFollowRequest{
		Id: id,
	})
	if err != nil {
		log.Printf("[FOLLOWING_GRPC] GetFollow RPC failed: %v", err)
		return nil, fmt.Errorf("GetFollow RPC failed: %w", err)
	}

	follow := r.followToDomain(resp.GetFollow())

	log.Printf("[FOLLOWING_GRPC] GetFollow RPC successful for ID: %s", id)
	return &models.GetFollowResponse{
		Follow: follow,
	}, nil
}

// GetAllFollowersForUser retrieves all followers for a user
func (r FollowingRepository) GetAllFollowersForUser(ctx context.Context, followedUserID string) (*models.GetFollowingResponse, error) {
	log.Printf("[FOLLOWING_GRPC] GetFollowing called for followedUserID: %s", followedUserID)

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[FOLLOWING_GRPC] Failed to connect to following service: %v", err)
		return nil, fmt.Errorf("failed to connect to following service: %w", err)
	}
	defer conn.Close()

	client := followingpb.NewFollowingServiceClient(conn)
	resp, err := client.GetFollowing(ctx, &followingpb.GetFollowingRequest{
		FollowedUserId: followedUserID,
	})
	if err != nil {
		log.Printf("[FOLLOWING_GRPC] GetFollowing RPC failed: %v", err)
		return nil, fmt.Errorf("GetFollowing RPC failed: %w", err)
	}

	follows := make([]*models.Follow, len(resp.GetFollowing()))
	for i, pbFollow := range resp.GetFollowing() {
		follows[i] = r.followToDomain(pbFollow)
	}

	log.Printf("[FOLLOWING_GRPC] GetFollowing RPC successful, returned %d follows", len(follows))
	return &models.GetFollowingResponse{
		Following: follows,
		Count:     int64(len(follows)),
	}, nil
}

// GetMostFollowedUsers retrieves the most followed items
func (r FollowingRepository) GetMostFollowedUsers(ctx context.Context) (*models.GetMostFollowedResponse, error) {
	log.Printf("[FOLLOWING_GRPC] GetMostFollowed called")

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[FOLLOWING_GRPC] Failed to connect to following service: %v", err)
		return nil, fmt.Errorf("failed to connect to following service: %w", err)
	}
	defer conn.Close()

	client := followingpb.NewFollowingServiceClient(conn)
	resp, err := client.GetMostFollowed(ctx, &followingpb.GetMostFollowedRequest{})
	if err != nil {
		log.Printf("[FOLLOWING_GRPC] GetMostFollowed RPC failed: %v", err)
		return nil, fmt.Errorf("GetMostFollowed RPC failed: %w", err)
	}

	counts := make([]*models.ItemFollowCount, len(resp.GetItemFollowCount()))
	for i, pbCount := range resp.GetItemFollowCount() {
		counts[i] = &models.ItemFollowCount{
			FollowedUserID: pbCount.GetFollowedUserId(),
			CategoryID:     pbCount.GetCategoryId(),
			FollowingCount: pbCount.GetFollowingCount(),
		}
	}

	log.Printf("[FOLLOWING_GRPC] GetMostFollowed RPC successful, returned %d counts", len(counts))
	return &models.GetMostFollowedResponse{
		ItemFollowCounts: counts,
		Count:            int64(len(counts)),
	}, nil
}

// GetMostFollowedUsersByCategory retrieves the most followed items in a category
func (r FollowingRepository) GetMostFollowedUsersByCategory(ctx context.Context, categoryID string, offset, limit int64) (*models.GetMostFollowedByCategoryResponse, error) {
	log.Printf("[FOLLOWING_GRPC] GetMostFollowedByCategory called for categoryID: %s", categoryID)

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[FOLLOWING_GRPC] Failed to connect to following service: %v", err)
		return nil, fmt.Errorf("failed to connect to following service: %w", err)
	}
	defer conn.Close()

	client := followingpb.NewFollowingServiceClient(conn)
	resp, err := client.GetMostFollowedByCategory(ctx, &followingpb.GetMostFollowedByCategoryRequest{
		CategoryId: categoryID,
		Offset:     offset,
		Limit:      limit,
	})
	if err != nil {
		log.Printf("[FOLLOWING_GRPC] GetMostFollowedByCategory RPC failed: %v", err)
		return nil, fmt.Errorf("GetMostFollowedByCategory RPC failed: %w", err)
	}

	counts := make([]*models.ItemFollowCount, len(resp.GetItemFollowCount()))
	for i, pbCount := range resp.GetItemFollowCount() {
		counts[i] = &models.ItemFollowCount{
			FollowedUserID: pbCount.GetFollowedUserId(),
			CategoryID:     pbCount.GetCategoryId(),
			FollowingCount: pbCount.GetFollowingCount(),
		}
	}

	log.Printf("[FOLLOWING_GRPC] GetMostFollowedByCategory RPC successful, returned %d counts", len(counts))
	return &models.GetMostFollowedByCategoryResponse{
		ItemFollowCounts: counts,
		Count:            int64(len(counts)),
	}, nil
}

// GetUserFollowingList retrieves all follows by a specific sender
func (r FollowingRepository) GetUserFollowingList(ctx context.Context, userID string) (*models.GetFollowingBySenderResponse, error) {
	log.Printf("[FOLLOWING_GRPC] GetFollowingBySender called for userID: %s", userID)

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[FOLLOWING_GRPC] Failed to connect to following service: %v", err)
		return nil, fmt.Errorf("failed to connect to following service: %w", err)
	}
	defer conn.Close()

	client := followingpb.NewFollowingServiceClient(conn)
	resp, err := client.GetFollowingBySender(ctx, &followingpb.GetFollowingBySenderRequest{})
	if err != nil {
		log.Printf("[FOLLOWING_GRPC] GetFollowingBySender RPC failed: %v", err)
		return nil, fmt.Errorf("GetFollowingBySender RPC failed: %w", err)
	}

	follows := make([]*models.Follow, len(resp.GetFollowing()))
	for i, pbFollow := range resp.GetFollowing() {
		follows[i] = r.followToDomain(pbFollow)
	}

	log.Printf("[FOLLOWING_GRPC] GetFollowingBySender RPC successful, returned %d follows", len(follows))
	return &models.GetFollowingBySenderResponse{
		Following: follows,
		Count:     int64(len(follows)),
	}, nil
}

// GetAllApprovedFollowRelationships retrieves all approved follows
func (r FollowingRepository) GetAllApprovedFollowRelationships(ctx context.Context) (*models.GetApprovedFollowingResponse, error) {
	log.Printf("[FOLLOWING_GRPC] GetApprovedFollowing called")

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[FOLLOWING_GRPC] Failed to connect to following service: %v", err)
		return nil, fmt.Errorf("failed to connect to following service: %w", err)
	}
	defer conn.Close()

	client := followingpb.NewFollowingServiceClient(conn)
	resp, err := client.GetApprovedFollowing(ctx, &followingpb.GetApprovedFollowingRequest{})
	if err != nil {
		log.Printf("[FOLLOWING_GRPC] GetApprovedFollowing RPC failed: %v", err)
		return nil, fmt.Errorf("GetApprovedFollowing RPC failed: %w", err)
	}

	follows := make([]*models.Follow, len(resp.GetFollowing()))
	for i, pbFollow := range resp.GetFollowing() {
		follows[i] = r.followToDomain(pbFollow)
	}

	log.Printf("[FOLLOWING_GRPC] GetApprovedFollowing RPC successful, returned %d follows", len(follows))
	return &models.GetApprovedFollowingResponse{
		Following: follows,
		Count:     int64(len(follows)),
	}, nil
}

// Additional query methods for AI tooling - these would require extending the protobuf
// For now, they return "not implemented" errors since the protobuf doesn't have these methods

func (r FollowingRepository) GetPaginatedFollowList(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*models.Follow, error) {
	return nil, fmt.Errorf("GetFollowsWithPagination not implemented - requires additional RPC method in protobuf")
}

func (r FollowingRepository) SearchFollowsByKeyword(ctx context.Context, term string) ([]*models.Follow, error) {
	return nil, fmt.Errorf("SearchFollowsWithTerm not implemented - requires additional RPC method in protobuf")
}

func (r FollowingRepository) GetUserFollowingHistoryPaginated(ctx context.Context, userID string, page, pageSize int64) ([]*models.Follow, error) {
	return nil, fmt.Errorf("GetUserFollowingHistory not implemented - requires additional RPC method in protobuf")
}

func (r FollowingRepository) GetFollowsByApprovalStatus(ctx context.Context, status string, page, pageSize int64) ([]*models.Follow, error) {
	return nil, fmt.Errorf("GetFollowsByStatus not implemented - requires additional RPC method in protobuf")
}

func (r FollowingRepository) GetFollowingStatistics(ctx context.Context, userID string) (*models.FollowingStatsResponse, error) {
	return nil, fmt.Errorf("GetFollowingStats not implemented - requires additional RPC method in protobuf")
}

func (r FollowingRepository) GetMutualFollowingBetweenUsers(ctx context.Context, userID1, userID2 string) ([]*models.Follow, error) {
	return nil, fmt.Errorf("GetMutualFollowing not implemented - requires additional RPC method in protobuf")
}

func (r FollowingRepository) GetFollowsByCategoryID(ctx context.Context, categoryID string, page, pageSize int64) ([]*models.Follow, error) {
	return nil, fmt.Errorf("GetFollowingByCategory not implemented - requires additional RPC method in protobuf")
}

// Legacy methods for backward compatibility
func (r FollowingRepository) CreateFollowRelationshipLegacy(ctx context.Context, followID, userID, followedUserID string, followedUserType models.UserType, content, categoryID, parentID string) error {
	_, err := r.CreateNewFollowRelationship(ctx, userID, followedUserID, followedUserType.String(), content, categoryID, parentID)
	return err
}

func (r FollowingRepository) FindFollowByIDAndUser(ctx context.Context, followID, followedUserID string) (*models.Follow, error) {
	resp, err := r.GetFollowByID(ctx, followID)
	if err != nil {
		return nil, err
	}
	return resp.Follow, nil
}

func (r FollowingRepository) GetAllFollowersList(ctx context.Context, followedUserID string) ([]*models.Follow, error) {
	resp, err := r.GetAllFollowersForUser(ctx, followedUserID)
	if err != nil {
		return nil, err
	}
	return resp.Following, nil
}

func (r FollowingRepository) FindFollowsByUserID(ctx context.Context, userID string) ([]*models.Follow, error) {
	resp, err := r.GetUserFollowingList(ctx, userID)
	if err != nil {
		return nil, err
	}
	return resp.Following, nil
}

func (r FollowingRepository) RemoveFollowRelationship(ctx context.Context, followID string) error {
	log.Printf("[FOLLOWING_GRPC] RemoveFollowRelationship called for followID: %s", followID)
	// Implementation would call the unfollow RPC
	return nil
}

// GetTotalFollowerCount gets the count of followers for a user
func (r FollowingRepository) GetTotalFollowerCount(ctx context.Context, userID string) (int, error) {
	log.Printf("[FOLLOWING_GRPC] GetTotalFollowerCount called for userID: %s", userID)

	// TODO: Implement when gRPC method is available
	// For now, return a placeholder count
	log.Printf("[FOLLOWING_GRPC] GetTotalFollowerCount placeholder implementation for userID: %s", userID)
	return 0, nil
}

// GetFollowersWithLimit gets the followers for a user
func (r FollowingRepository) GetFollowersWithLimit(ctx context.Context, userID string, limit int32) ([]*models.Follow, error) {
	log.Printf("[FOLLOWING_GRPC] GetFollowersWithLimit called for userID: %s, limit: %d", userID, limit)

	// TODO: Implement when gRPC method is available
	// For now, return an empty slice
	log.Printf("[FOLLOWING_GRPC] GetFollowersWithLimit placeholder implementation for userID: %s", userID)
	return []*models.Follow{}, nil
}

// GetTotalFollowingCount gets the count of users being followed by a user
func (r FollowingRepository) GetTotalFollowingCount(ctx context.Context, userID string) (int, error) {
	log.Printf("[FOLLOWING_GRPC] GetTotalFollowingCount called for userID: %s", userID)

	// TODO: Implement when gRPC method is available
	// For now, return a placeholder count
	log.Printf("[FOLLOWING_GRPC] GetTotalFollowingCount placeholder implementation for userID: %s", userID)
	return 0, nil
}

// CheckIfUserIsFollowing checks if a user is following another user
func (r FollowingRepository) CheckIfUserIsFollowing(ctx context.Context, userID, followedUserID string) (bool, error) {
	log.Printf("[FOLLOWING_GRPC] CheckIfUserIsFollowing called for userID: %s, followedUserID: %s", userID, followedUserID)

	// TODO: Implement when gRPC method is available
	// For now, return false
	log.Printf("[FOLLOWING_GRPC] CheckIfUserIsFollowing placeholder implementation for userID: %s, followedUserID: %s", userID, followedUserID)
	return false, nil
}

// Helper methods

// followToDomain converts a protobuf Follow to domain Follow
func (r FollowingRepository) followToDomain(pbFollow *followingpb.Follow) *models.Follow {
	if pbFollow == nil {
		return nil
	}

	return &models.Follow{
		ID:               pbFollow.GetId(),
		UserID:           pbFollow.GetUserId(),
		FollowedUserID:   pbFollow.GetFollowedUserId(),
		FollowedUserType: pbFollow.GetFollowedUserType(),
		Content:          pbFollow.GetContent(),
		CategoryID:       pbFollow.GetCategoryId(),
		ParentID:         pbFollow.GetParentId(),
		FollowStatus:     pbFollow.GetFollowStatus(),
		Flagged:          pbFollow.GetFlagged(),
		// Set legacy Approved field based on status
		Approved: pbFollow.GetFollowStatus() == models.FollowStatusApproved,
		// Note: CreatedAt, UpdatedAt, ApprovedAt would need to be added to protobuf
		CreatedAt:  time.Now(),  // Placeholder
		UpdatedAt:  time.Now(),  // Placeholder
		ApprovedAt: time.Time{}, // Placeholder
	}
}

// dial establishes a gRPC connection to the following service
// dial sets up a gRPC connection with the microservice endpoint
func (r FollowingRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

func (r FollowingRepository) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}
