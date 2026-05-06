package domain

import (
	"context"
	"middleman/assistants/internal/models"
)

type CommentRepository interface {
	// Core comment operations from protobuf
	CreateNewComment(ctx context.Context, itemID, itemType, content, categoryID, parentID string) (*models.AddCommentResponse, error)
	ApproveCommentByID(ctx context.Context, id string) (*models.ApproveCommentResponse, error)
	RejectCommentByID(ctx context.Context, id string) (*models.RejectCommentResponse, error)
	FlagCommentAsInappropriate(ctx context.Context, id string) (*models.FlagCommentResponse, error)
	EditCommentContent(ctx context.Context, id, content string) (*models.EditCommentResponse, error)
	DeleteCommentByID(ctx context.Context, id string) (*models.RemoveCommentResponse, error)
	GetCommentByID(ctx context.Context, id string) (*models.GetCommentResponse, error)
	GetAllCommentsForItem(ctx context.Context, itemID string) (*models.GetCommentsResponse, error)
	GetMostCommentedItems(ctx context.Context) (*models.GetMostCommentedResponse, error)
	GetMostCommentedItemsByCategory(ctx context.Context, categoryID string, offset, limit int64) (*models.GetMostCommentedByCategoryResponse, error)
	GetCommentsBySenderUser(ctx context.Context) (*models.GetCommentsBySenderResponse, error)
	GetAllApprovedComments(ctx context.Context) (*models.GetApprovedCommentsResponse, error)

	// Additional query methods for AI tooling and repository pattern compatibility
	// These would require additional RPC methods to be added to the protobuf
	GetPaginatedCommentList(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*models.Comment, error)
	SearchCommentsByKeyword(ctx context.Context, term string) ([]*models.Comment, error)
	GetCommentsByApprovalStatus(ctx context.Context, status string, page, pageSize int64) ([]*models.Comment, error)
	GetCommentsByCategoryID(ctx context.Context, categoryID string, page, pageSize int64) ([]*models.Comment, error)
	GetCommentStatistics(ctx context.Context, itemID string) (*models.CommentStatsResponse, error)
	GetRecentlyPostedComments(ctx context.Context, page, pageSize int64) ([]*models.Comment, error)
	GetChildCommentsForParent(ctx context.Context, parentID string) ([]*models.Comment, error)

	// Simplified methods for backward compatibility
	FindCommentByItemAndID(ctx context.Context, commentID, itemID string) (*models.Comment, error)
	GetAllCommentsByItemID(ctx context.Context, itemID string) ([]*models.Comment, error)
	FindAllCommentsBySenderUserID(ctx context.Context) ([]*models.Comment, error)
}
