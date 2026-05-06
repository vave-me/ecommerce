package grpc

import (
	"context"
	"fmt"
	"log"
	"middleman/assistants/internal/domain"
	"middleman/assistants/internal/models"
	"middleman/comments/commentspb"
	"middleman/internal/auth"
	"middleman/internal/rpc"
	"time"

	"google.golang.org/grpc"
)

// CommentRepository calls the remote comments service (gRPC) as a fallback.
type CommentRepository struct {
	endpoint string
	auth     *auth.Auth // Optional auth for authenticated calls
}

var _ domain.CommentRepository = (*CommentRepository)(nil)

// NewCommentRepositoryWithAuth creates a new comment repository with JWT authentication support
func NewCommentRepository(endpoint string, authProvider *auth.Auth) CommentRepository {
	return CommentRepository{
		endpoint: endpoint,
		auth:     authProvider,
	}
}

// Core protobuf methods

func (r CommentRepository) CreateNewComment(ctx context.Context, itemID, itemType, content, categoryID, parentID string) (*models.AddCommentResponse, error) {
	log.Printf("[COMMENT_GRPC] CreateNewComment: adding comment for item=%s via gRPC", itemID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		log.Printf("[COMMENT_GRPC] CreateNewComment: failed to dial gRPC: %v", err)
		return nil, err
	}
	defer conn.Close()

	client := commentspb.NewCommentsServiceClient(conn)

	resp, err := client.AddComment(ctx, &commentspb.AddCommentRequest{
		ItemId:     itemID,
		ItemType:   itemType,
		Content:    content,
		CategoryId: categoryID,
		ParentId:   parentID,
	})
	if err != nil {
		log.Printf("[COMMENT_GRPC] CreateNewComment: RPC failed: %v", err)
		return nil, fmt.Errorf("AddComment RPC failed: %w", err)
	}

	return &models.AddCommentResponse{
		ID: resp.GetId(),
	}, nil
}

func (r CommentRepository) ApproveCommentByID(ctx context.Context, id string) (*models.ApproveCommentResponse, error) {
	log.Printf("[COMMENT_GRPC] ApproveCommentByID: approving comment ID=%s via gRPC", id)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := commentspb.NewCommentsServiceClient(conn)
	resp, err := client.ApproveComment(ctx, &commentspb.ApproveCommentRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("ApproveComment RPC failed: %w", err)
	}

	var approvedAt time.Time
	if resp.GetApprovedAt() != nil {
		approvedAt = resp.GetApprovedAt().AsTime()
	}

	return &models.ApproveCommentResponse{
		ID:            resp.GetId(),
		CommentStatus: resp.GetCommentStatus(),
		ApprovedAt:    approvedAt,
	}, nil
}

func (r CommentRepository) RejectCommentByID(ctx context.Context, id string) (*models.RejectCommentResponse, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := commentspb.NewCommentsServiceClient(conn)
	resp, err := client.RejectComment(ctx, &commentspb.RejectCommentRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("RejectComment RPC failed: %w", err)
	}

	return &models.RejectCommentResponse{
		ID:            resp.GetId(),
		CommentStatus: resp.GetCommentStatus(),
	}, nil
}

func (r CommentRepository) FlagCommentAsInappropriate(ctx context.Context, id string) (*models.FlagCommentResponse, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := commentspb.NewCommentsServiceClient(conn)
	resp, err := client.FlagComment(ctx, &commentspb.FlagCommentRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("FlagComment RPC failed: %w", err)
	}

	return &models.FlagCommentResponse{
		ID:      resp.GetId(),
		Flagged: resp.GetFlagged(),
	}, nil
}

func (r CommentRepository) EditCommentContent(ctx context.Context, id, content string) (*models.EditCommentResponse, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := commentspb.NewCommentsServiceClient(conn)
	resp, err := client.EditComment(ctx, &commentspb.EditCommentRequest{
		Id:      id,
		Content: content,
	})
	if err != nil {
		return nil, fmt.Errorf("EditComment RPC failed: %w", err)
	}

	return &models.EditCommentResponse{
		ID:      resp.GetId(),
		Content: resp.GetContent(),
	}, nil
}

func (r CommentRepository) DeleteCommentByID(ctx context.Context, id string) (*models.RemoveCommentResponse, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := commentspb.NewCommentsServiceClient(conn)
	_, err = client.RemoveComment(ctx, &commentspb.RemoveCommentRequest{Id: id})
	if err != nil {
		return &models.RemoveCommentResponse{
			Success: false,
			Message: fmt.Sprintf("RemoveComment RPC failed: %v", err),
		}, err
	}

	return &models.RemoveCommentResponse{
		Success: true,
		Message: "Comment removed successfully",
	}, nil
}

func (r CommentRepository) GetCommentByID(ctx context.Context, id string) (*models.GetCommentResponse, error) {
	log.Printf("[COMMENT_GRPC] GetCommentByID: retrieving comment ID=%s via gRPC", id)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := commentspb.NewCommentsServiceClient(conn)
	resp, err := client.GetComment(ctx, &commentspb.GetCommentRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("GetComment RPC failed: %w", err)
	}

	return &models.GetCommentResponse{
		Comment: r.commentToDomain(resp.GetComment()),
	}, nil
}

func (r CommentRepository) GetAllCommentsForItem(ctx context.Context, itemID string) (*models.GetCommentsResponse, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := commentspb.NewCommentsServiceClient(conn)
	resp, err := client.GetComments(ctx, &commentspb.GetCommentsRequest{ItemId: itemID})
	if err != nil {
		return nil, fmt.Errorf("GetComments RPC failed: %w", err)
	}

	comments := make([]*models.Comment, 0, len(resp.GetComments()))
	for _, pbComment := range resp.GetComments() {
		if domainComment := r.commentToDomain(pbComment); domainComment != nil {
			comments = append(comments, domainComment)
		}
	}

	return &models.GetCommentsResponse{
		Comments: comments,
	}, nil
}

func (r CommentRepository) GetMostCommentedItems(ctx context.Context) (*models.GetMostCommentedResponse, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := commentspb.NewCommentsServiceClient(conn)
	resp, err := client.GetMostCommented(ctx, &commentspb.GetMostCommentedRequest{})
	if err != nil {
		return nil, fmt.Errorf("GetMostCommented RPC failed: %w", err)
	}

	counts := make([]*models.ItemCommentCount, 0, len(resp.GetItemCommentCount()))
	for _, pbCount := range resp.GetItemCommentCount() {
		counts = append(counts, &models.ItemCommentCount{
			ItemID:        pbCount.GetItemId(),
			CategoryID:    pbCount.GetCategoryId(),
			CommentsCount: pbCount.GetCommentsCount(),
		})
	}

	return &models.GetMostCommentedResponse{
		ItemCommentCounts: counts,
	}, nil
}

func (r CommentRepository) GetMostCommentedItemsByCategory(ctx context.Context, categoryID string, offset, limit int64) (*models.GetMostCommentedByCategoryResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := commentspb.NewCommentsServiceClient(conn)
	resp, err := client.GetMostCommentedByCategory(ctx, &commentspb.GetMostCommentedByCategoryRequest{
		CategoryId: categoryID,
		Offset:     offset,
		Limit:      limit,
	})
	if err != nil {
		return nil, fmt.Errorf("GetMostCommentedByCategory RPC failed: %w", err)
	}

	counts := make([]*models.ItemCommentCount, 0, len(resp.GetItemCommentCount()))
	for _, pbCount := range resp.GetItemCommentCount() {
		counts = append(counts, &models.ItemCommentCount{
			ItemID:        pbCount.GetItemId(),
			CategoryID:    pbCount.GetCategoryId(),
			CommentsCount: pbCount.GetCommentsCount(),
		})
	}

	return &models.GetMostCommentedByCategoryResponse{
		ItemCommentCounts: counts,
	}, nil
}

func (r CommentRepository) GetCommentsBySenderUser(ctx context.Context) (*models.GetCommentsBySenderResponse, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := commentspb.NewCommentsServiceClient(conn)
	resp, err := client.GetCommentsBySender(ctx, &commentspb.GetCommentsBySenderRequest{})
	if err != nil {
		return nil, fmt.Errorf("GetCommentsBySender RPC failed: %w", err)
	}

	comments := make([]*models.Comment, 0, len(resp.GetComments()))
	for _, pbComment := range resp.GetComments() {
		if domainComment := r.commentToDomain(pbComment); domainComment != nil {
			comments = append(comments, domainComment)
		}
	}

	return &models.GetCommentsBySenderResponse{
		Comments: comments,
	}, nil
}

func (r CommentRepository) GetAllApprovedComments(ctx context.Context) (*models.GetApprovedCommentsResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := commentspb.NewCommentsServiceClient(conn)
	resp, err := client.GetApprovedComments(ctx, &commentspb.GetApprovedCommentsRequest{})
	if err != nil {
		return nil, fmt.Errorf("GetApprovedComments RPC failed: %w", err)
	}

	comments := make([]*models.Comment, 0, len(resp.GetComments()))
	for _, pbComment := range resp.GetComments() {
		if domainComment := r.commentToDomain(pbComment); domainComment != nil {
			comments = append(comments, domainComment)
		}
	}

	return &models.GetApprovedCommentsResponse{
		Comments: comments,
	}, nil
}

// Additional query methods for AI tooling
func (r CommentRepository) GetPaginatedCommentList(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*models.Comment, error) {
	// This would require additional RPC method
	return nil, fmt.Errorf("GetCommentsWithPagination not implemented - requires additional RPC method")
}

func (r CommentRepository) SearchCommentsByKeyword(ctx context.Context, term string) ([]*models.Comment, error) {
	// This would require additional RPC method
	return nil, fmt.Errorf("SearchCommentsWithTerm not implemented - requires additional RPC method")
}

func (r CommentRepository) GetCommentsByApprovalStatus(ctx context.Context, status string, page, pageSize int64) ([]*models.Comment, error) {
	// This would require additional RPC method
	return nil, fmt.Errorf("GetCommentsByStatus not implemented - requires additional RPC method")
}

func (r CommentRepository) GetCommentsByCategoryID(ctx context.Context, categoryID string, page, pageSize int64) ([]*models.Comment, error) {
	// This would require additional RPC method
	return nil, fmt.Errorf("GetCommentsByCategory not implemented - requires additional RPC method")
}

func (r CommentRepository) GetCommentStatistics(ctx context.Context, itemID string) (*models.CommentStatsResponse, error) {
	// This would require additional RPC method
	return nil, fmt.Errorf("GetCommentStats not implemented - requires additional RPC method")
}

func (r CommentRepository) GetRecentlyPostedComments(ctx context.Context, page, pageSize int64) ([]*models.Comment, error) {
	// This would require additional RPC method
	return nil, fmt.Errorf("GetRecentComments not implemented - requires additional RPC method")
}

func (r CommentRepository) GetChildCommentsForParent(ctx context.Context, parentID string) ([]*models.Comment, error) {
	// This would require additional RPC method
	return nil, fmt.Errorf("GetCommentThread not implemented - requires additional RPC method")
}

// Legacy methods removed - use standard methods instead

func (r CommentRepository) FindCommentByItemAndID(ctx context.Context, commentID, itemID string) (*models.Comment, error) {
	resp, err := r.GetCommentByID(ctx, commentID)
	if err != nil {
		return nil, err
	}
	return resp.Comment, nil
}

func (r CommentRepository) GetAllCommentsByItemID(ctx context.Context, itemID string) ([]*models.Comment, error) {
	resp, err := r.GetAllCommentsForItem(ctx, itemID)
	if err != nil {
		return nil, err
	}
	return resp.Comments, nil
}

func (r CommentRepository) FindAllCommentsBySenderUserID(ctx context.Context) ([]*models.Comment, error) {
	resp, err := r.GetCommentsBySenderUser(ctx)
	if err != nil {
		return nil, err
	}
	return resp.Comments, nil
}

// RemoveCommentLegacy removed - use RemoveComment instead

// Helper methods
func (r CommentRepository) commentToDomain(pb *commentspb.Comment) *models.Comment {
	if pb == nil {
		return nil
	}

	return &models.Comment{
		ID:            pb.GetId(),
		SenderID:      pb.GetSenderId(),
		ItemID:        pb.GetItemId(),
		ItemType:      pb.GetItemType(),
		Content:       pb.GetContent(),
		CategoryID:    pb.GetCategoryId(),
		ParentID:      pb.GetParentId(),
		CommentStatus: pb.GetCommentStatus(),
		Flagged:       pb.GetFlagged(),
		CreatedAt:     time.Now(), // Protobuf doesn't include timestamps in this spec
		UpdatedAt:     time.Now(),
		Approved:      pb.GetCommentStatus() == models.CommentStatusApproved,
	}
}

// dial sets up a gRPC connection with the microservice endpoint
func (r CommentRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

func (r CommentRepository) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}
