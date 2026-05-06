package domain

import (
	"context"
	"middleman/assistants/internal/models"
)

type FollowingRepository interface {
	// Core following operations from protobuf
	CreateNewFollowRelationship(ctx context.Context, userID, followedUserID, followedUserType, content, categoryID, parentID string) (*models.AddFollowResponse, error)
	ApproveFollowRequestByID(ctx context.Context, id string) (*models.ApproveFollowResponse, error)
	RejectFollowRequestByID(ctx context.Context, id string) (*models.RejectFollowResponse, error)
	FlagFollowAsInappropriate(ctx context.Context, id string) (*models.FlagFollowResponse, error)
	EditFollowDescription(ctx context.Context, id, content string) (*models.EditFollowResponse, error)
	DeleteFollowRelationship(ctx context.Context, id string) (*models.RemoveFollowResponse, error)
	GetFollowByID(ctx context.Context, id string) (*models.GetFollowResponse, error)
	GetAllFollowersForUser(ctx context.Context, followedUserID string) (*models.GetFollowingResponse, error)
	GetMostFollowedUsers(ctx context.Context) (*models.GetMostFollowedResponse, error)
	GetMostFollowedUsersByCategory(ctx context.Context, categoryID string, offset, limit int64) (*models.GetMostFollowedByCategoryResponse, error)
	GetUserFollowingList(ctx context.Context, userID string) (*models.GetFollowingBySenderResponse, error)
	GetAllApprovedFollowRelationships(ctx context.Context) (*models.GetApprovedFollowingResponse, error)

	// Additional query methods for AI tooling and repository pattern compatibility
	// These would require additional RPC methods to be added to the protobuf
	GetPaginatedFollowList(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*models.Follow, error)
	SearchFollowsByKeyword(ctx context.Context, term string) ([]*models.Follow, error)
	GetUserFollowingHistoryPaginated(ctx context.Context, userID string, page, pageSize int64) ([]*models.Follow, error)
	GetFollowsByApprovalStatus(ctx context.Context, status string, page, pageSize int64) ([]*models.Follow, error)
	GetFollowingStatistics(ctx context.Context, userID string) (*models.FollowingStatsResponse, error)
	GetMutualFollowingBetweenUsers(ctx context.Context, userID1, userID2 string) ([]*models.Follow, error)
	GetFollowsByCategoryID(ctx context.Context, categoryID string, page, pageSize int64) ([]*models.Follow, error)

	// Legacy methods for backward compatibility
	CreateFollowRelationshipLegacy(ctx context.Context, followID, userID, followedUserID string, followedUserType models.UserType, content, categoryID, parentID string) error
	FindFollowByIDAndUser(ctx context.Context, followID, followedUserID string) (*models.Follow, error)
	GetAllFollowersList(ctx context.Context, followedUserID string) ([]*models.Follow, error)
	FindFollowsByUserID(ctx context.Context, userID string) ([]*models.Follow, error)
	RemoveFollowRelationship(ctx context.Context, followID string) error

	// Methods needed by wishlist tool service
	GetFollowersWithLimit(ctx context.Context, userID string, limit int32) ([]*models.Follow, error)
	CheckIfUserIsFollowing(ctx context.Context, userID, followedUserID string) (bool, error)
	GetTotalFollowerCount(ctx context.Context, userID string) (int, error)
	GetTotalFollowingCount(ctx context.Context, userID string) (int, error)
}
