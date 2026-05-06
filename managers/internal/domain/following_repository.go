package domain

import (
	"context"
	"middleman/managers/internal/models"
)

type FollowingRepository interface {
	// Core following operations from protobuf
	AddFollow(ctx context.Context, userID, followedUserID, followedUserType, content, categoryID, parentID string) (*models.AddFollowResponse, error)
	ApproveFollow(ctx context.Context, id string) (*models.ApproveFollowResponse, error)
	RejectFollow(ctx context.Context, id string) (*models.RejectFollowResponse, error)
	FlagFollow(ctx context.Context, id string) (*models.FlagFollowResponse, error)
	EditFollow(ctx context.Context, id, content string) (*models.EditFollowResponse, error)
	RemoveFollow(ctx context.Context, id string) (*models.RemoveFollowResponse, error)
	GetFollow(ctx context.Context, id string) (*models.GetFollowResponse, error)
	GetFollowing(ctx context.Context, followedUserID string) (*models.GetFollowingResponse, error)
	GetMostFollowed(ctx context.Context) (*models.GetMostFollowedResponse, error)
	GetMostFollowedByCategory(ctx context.Context, categoryID string, offset, limit int64) (*models.GetMostFollowedByCategoryResponse, error)
	GetFollowingBySender(ctx context.Context, userID string) (*models.GetFollowingBySenderResponse, error)
	GetApprovedFollowing(ctx context.Context) (*models.GetApprovedFollowingResponse, error)

	// Additional query methods for AI tooling and repository pattern compatibility
	// These would require additional RPC methods to be added to the protobuf
	GetFollowsWithPagination(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*models.Follow, error)
	SearchFollowsWithTerm(ctx context.Context, term string) ([]*models.Follow, error)
	GetUserFollowingHistory(ctx context.Context, userID string, page, pageSize int64) ([]*models.Follow, error)
	GetFollowsByStatus(ctx context.Context, status string, page, pageSize int64) ([]*models.Follow, error)
	GetFollowingStats(ctx context.Context, userID string) (*models.FollowingStatsResponse, error)
	GetMutualFollowing(ctx context.Context, userID1, userID2 string) ([]*models.Follow, error)
	GetFollowingByCategory(ctx context.Context, categoryID string, page, pageSize int64) ([]*models.Follow, error)

	// Legacy methods for backward compatibility
	FollowUser(ctx context.Context, followID, userID, followedUserID string, followedUserType models.UserType, content, categoryID, parentID string) error
	FindFollow(ctx context.Context, followID, followedUserID string) (*models.Follow, error)
	AllFollowers(ctx context.Context, followedUserID string) ([]*models.Follow, error)
	FindByUserID(ctx context.Context, userID string) ([]*models.Follow, error)
	UnfollowUser(ctx context.Context, followID string) error

	// Methods needed by wishlist tool service
	GetFollowers(ctx context.Context, userID string, limit int32) ([]*models.Follow, error)
	IsFollowing(ctx context.Context, userID, followedUserID string) (bool, error)
	GetFollowerCount(ctx context.Context, userID string) (int, error)
	GetFollowingCount(ctx context.Context, userID string) (int, error)
}
