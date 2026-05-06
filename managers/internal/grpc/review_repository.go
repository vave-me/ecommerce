package grpc

import (
	"context"
	"fmt"
	"middleman/internal/auth"
	"middleman/internal/rpc"
	"middleman/managers/internal/domain"
	"middleman/managers/internal/models"
	"middleman/reviews/reviewspb"

	"google.golang.org/grpc"
)

// ReviewRepository calls the remote reviews service (gRPC) as a fallback.
type ReviewRepository struct {
	endpoint string
	auth     *auth.Auth // Optional auth for authenticated calls
}

var _ domain.ReviewRepository = (*ReviewRepository)(nil)

// NewReviewRepositoryWithAuth creates a new ReviewRepository with JWT authentication support
func NewReviewRepository(endpoint string, authInstance *auth.Auth) ReviewRepository {
	return ReviewRepository{
		endpoint: endpoint,
		auth:     authInstance,
	}
}

// AddReview adds a review via the reviews microservice
func (r ReviewRepository) AddReview(ctx context.Context, senderID, itemID, itemType, content, categoryID, parentID string) (string, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := reviewspb.NewReviewsServiceClient(conn)
	resp, err := client.AddReview(ctx, &reviewspb.AddReviewRequest{

		ItemId:     itemID,
		ItemType:   itemType,
		Content:    content,
		CategoryId: categoryID,
		ParentId:   parentID,
	})
	if err != nil {
		return "", fmt.Errorf("AddReview RPC failed: %w", err)
	}

	return resp.GetId(), nil
}

// GetReview retrieves a review by ID from the reviews microservice
func (r ReviewRepository) GetReview(ctx context.Context, reviewID string) (*models.Review, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := reviewspb.NewReviewsServiceClient(conn)
	resp, err := client.GetReview(ctx, &reviewspb.GetReviewRequest{
		Id: reviewID,
	})
	if err != nil {
		return nil, fmt.Errorf("GetReview RPC failed: %w", err)
	}

	return r.reviewToDomain(resp.GetReview()), nil
}

// GetReviews retrieves all reviews for an item
func (r ReviewRepository) GetReviews(ctx context.Context, itemID string) ([]*models.Review, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := reviewspb.NewReviewsServiceClient(conn)
	resp, err := client.GetReviews(ctx, &reviewspb.GetReviewsRequest{
		ItemId: itemID,
	})
	if err != nil {
		return nil, fmt.Errorf("GetReviews RPC failed: %w", err)
	}

	reviews := make([]*models.Review, 0, len(resp.GetReviews()))
	for _, pbReview := range resp.GetReviews() {
		reviews = append(reviews, r.reviewToDomain(pbReview))
	}

	return reviews, nil
}

// EditReview edits a review's content
func (r ReviewRepository) EditReview(ctx context.Context, reviewID, content string) error {
	conn, err := r.dial(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := reviewspb.NewReviewsServiceClient(conn)
	_, err = client.EditReview(ctx, &reviewspb.EditReviewRequest{
		Id:      reviewID,
		Content: content,
	})
	if err != nil {
		return fmt.Errorf("EditReview RPC failed: %w", err)
	}

	return nil
}

// RemoveReview removes a review by ID
func (r ReviewRepository) RemoveReview(ctx context.Context, reviewID string) error {
	conn, err := r.dial(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := reviewspb.NewReviewsServiceClient(conn)
	_, err = client.RemoveReview(ctx, &reviewspb.RemoveReviewRequest{
		Id: reviewID,
	})
	if err != nil {
		return fmt.Errorf("RemoveReview RPC failed: %w", err)
	}

	return nil
}

// ApproveReview approves a review
func (r ReviewRepository) ApproveReview(ctx context.Context, reviewID string) (*models.Review, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := reviewspb.NewReviewsServiceClient(conn)
	resp, err := client.ApproveReview(ctx, &reviewspb.ApproveReviewRequest{
		Id: reviewID,
	})
	if err != nil {
		return nil, fmt.Errorf("ApproveReview RPC failed: %w", err)
	}

	// Create a review object with approval info
	review := &models.Review{
		ID:           resp.GetId(),
		ReviewStatus: resp.GetReviewStatus(),
	}

	// Convert timestamp if provided
	if resp.GetApprovedAt() != nil {
		approvedAt := resp.GetApprovedAt().AsTime()
		review.ApprovedAt = &approvedAt
	}

	return review, nil
}

// RejectReview rejects a review
func (r ReviewRepository) RejectReview(ctx context.Context, reviewID string) (*models.Review, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := reviewspb.NewReviewsServiceClient(conn)
	resp, err := client.RejectReview(ctx, &reviewspb.RejectReviewRequest{
		Id: reviewID,
	})
	if err != nil {
		return nil, fmt.Errorf("RejectReview RPC failed: %w", err)
	}

	return &models.Review{
		ID:           resp.GetId(),
		ReviewStatus: resp.GetReviewStatus(),
	}, nil
}

// FlagReview flags a review
func (r ReviewRepository) FlagReview(ctx context.Context, reviewID string) (*models.Review, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := reviewspb.NewReviewsServiceClient(conn)
	resp, err := client.FlagReview(ctx, &reviewspb.FlagReviewRequest{
		Id: reviewID,
	})
	if err != nil {
		return nil, fmt.Errorf("FlagReview RPC failed: %w", err)
	}

	return &models.Review{
		ID:      resp.GetId(),
		Flagged: resp.GetFlagged(),
	}, nil
}

// GetReviewsBySender retrieves all reviews by a sender
func (r ReviewRepository) GetReviewsBySender(ctx context.Context, senderID string) ([]*models.Review, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := reviewspb.NewReviewsServiceClient(conn)
	resp, err := client.GetReviewsBySender(ctx, &reviewspb.GetReviewsBySenderRequest{
		SenderId: senderID,
	})
	if err != nil {
		return nil, fmt.Errorf("GetReviewsBySender RPC failed: %w", err)
	}

	reviews := make([]*models.Review, 0, len(resp.GetReviews()))
	for _, pbReview := range resp.GetReviews() {
		reviews = append(reviews, r.reviewToDomain(pbReview))
	}

	return reviews, nil
}

// GetApprovedReviews retrieves all approved reviews
func (r ReviewRepository) GetApprovedReviews(ctx context.Context) ([]*models.Review, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := reviewspb.NewReviewsServiceClient(conn)
	resp, err := client.GetApprovedReviews(ctx, &reviewspb.GetApprovedReviewsRequest{})
	if err != nil {
		return nil, fmt.Errorf("GetApprovedReviews RPC failed: %w", err)
	}

	reviews := make([]*models.Review, 0, len(resp.GetReviews()))
	for _, pbReview := range resp.GetReviews() {
		reviews = append(reviews, r.reviewToDomain(pbReview))
	}

	return reviews, nil
}

// GetMostReviewed retrieves items with the most reviews
func (r ReviewRepository) GetMostReviewed(ctx context.Context) ([]*models.ItemReviewCount, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := reviewspb.NewReviewsServiceClient(conn)
	resp, err := client.GetMostReviewed(ctx, &reviewspb.GetMostReviewedRequest{})
	if err != nil {
		return nil, fmt.Errorf("GetMostReviewed RPC failed: %w", err)
	}

	counts := make([]*models.ItemReviewCount, 0, len(resp.GetItemReviewCount()))
	for _, pbCount := range resp.GetItemReviewCount() {
		counts = append(counts, &models.ItemReviewCount{
			ItemID:       pbCount.GetItemId(),
			CategoryID:   pbCount.GetCategoryId(),
			ReviewsCount: pbCount.GetReviewsCount(),
		})
	}

	return counts, nil
}

// GetMostReviewedByCategory retrieves items with the most reviews in a category
func (r ReviewRepository) GetMostReviewedByCategory(ctx context.Context, categoryID string, offset, limit int64) ([]*models.ItemReviewCount, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := reviewspb.NewReviewsServiceClient(conn)
	resp, err := client.GetMostReviewedByCategory(ctx, &reviewspb.GetMostReviewedByCategoryRequest{
		CategoryId: categoryID,
		Offset:     offset,
		Limit:      limit,
	})
	if err != nil {
		return nil, fmt.Errorf("GetMostReviewedByCategory RPC failed: %w", err)
	}

	counts := make([]*models.ItemReviewCount, 0, len(resp.GetItemReviewCount()))
	for _, pbCount := range resp.GetItemReviewCount() {
		counts = append(counts, &models.ItemReviewCount{
			ItemID:       pbCount.GetItemId(),
			CategoryID:   pbCount.GetCategoryId(),
			ReviewsCount: pbCount.GetReviewsCount(),
		})
	}

	return counts, nil
}

// Legacy methods removed - use standard methods instead

// reviewToDomain converts a reviewspb.Review into our internal models.Review
func (r ReviewRepository) reviewToDomain(pb *reviewspb.Review) *models.Review {
	if pb == nil {
		return nil
	}

	return &models.Review{
		ID:           pb.GetId(),
		SenderID:     pb.GetSenderId(),
		ItemID:       pb.GetItemId(),
		ItemType:     pb.GetItemType(),
		Content:      pb.GetContent(),
		CategoryID:   pb.GetCategoryId(),
		ParentID:     pb.GetParentId(),
		ReviewStatus: pb.GetReviewStatus(),
		Flagged:      pb.GetFlagged(),
		// Note: CreatedAt, UpdatedAt, and ApprovedAt would need to be added to the protobuf
		// if they're needed for domain operations
	}
}

// dial establishes a gRPC connection to the reviews service
// dial sets up a gRPC connection with the microservice endpoint
func (r ReviewRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

func (r ReviewRepository) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}

// Additional review operations needed by tool service

// UpdateReview updates the content of an existing review
func (r ReviewRepository) UpdateReview(ctx context.Context, reviewID, content string) error {
	// For now, use the existing EditReview method
	return r.EditReview(ctx, reviewID, content)
}

// DeleteReview deletes a review (alias for RemoveReview)
func (r ReviewRepository) DeleteReview(ctx context.Context, reviewID string) error {
	return r.RemoveReview(ctx, reviewID)
}

// UnflagReview removes the flag from a review
func (r ReviewRepository) UnflagReview(ctx context.Context, reviewID string) error {
	// For now, return nil as this would need to be implemented in the reviews service
	// This is a placeholder implementation
	return nil
}

// SearchReviews searches for reviews based on a query
func (r ReviewRepository) SearchReviews(ctx context.Context, query string, limit int64) ([]*models.Review, error) {
	// For now, return empty results as this would need to be implemented in the reviews service
	// This is a placeholder implementation
	return []*models.Review{}, nil
}

// AllReviews implements the missing interface method (legacy compatibility)
func (r ReviewRepository) AllReviews(ctx context.Context, itemID string) ([]*models.Review, error) {
	return r.GetReviews(ctx, itemID)
}

// FindOneReview implements the missing interface method (legacy compatibility)
func (r ReviewRepository) FindOneReview(ctx context.Context, reviewID, itemID string) (*models.Review, error) {
	return r.GetReview(ctx, reviewID)
}

// FindBySenderID implements the missing interface method (legacy compatibility)
func (r ReviewRepository) FindBySenderID(ctx context.Context, senderID string) ([]*models.Review, error) {
	return r.GetReviewsBySender(ctx, senderID)
}
