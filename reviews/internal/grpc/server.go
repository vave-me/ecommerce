package grpc

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpc_code "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"middleman/internal/auth"
	"middleman/internal/errorsotel"
	"middleman/reviews/internal/application"
	"middleman/reviews/internal/application/commands"
	"middleman/reviews/internal/application/queries"
	"middleman/reviews/internal/domain"
	"middleman/reviews/reviewspb"
)

type server struct {
	app application.App
	reviewspb.UnimplementedReviewsServiceServer
}

var _ reviewspb.ReviewsServiceServer = (*server)(nil)

func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	reviewspb.RegisterReviewsServiceServer(registrar, server{app: app})
	return nil
}
func (s server) AddReview(ctx context.Context, request *reviewspb.AddReviewRequest) (*reviewspb.AddReviewResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	span := trace.SpanFromContext(ctx)
	reviewID := uuid.New().String()

	span.SetAttributes(
		attribute.String("ReviewID", reviewID),
	)
	fmt.Println("SenderID", userID)
	fmt.Println("ItemID", request.ItemId)
	fmt.Println("ItemType", request.ItemType)
	fmt.Println("Content", request.Content)
	fmt.Println("CategoryID", request.CategoryId)
	fmt.Println("ParentID", request.ParentId)

	err := s.app.AddReview(ctx, commands.AddReview{
		ID:         reviewID,
		SenderID:   userID,
		ItemID:     request.GetItemId(),
		ItemType:   domain.ToItemType(request.GetItemType()),
		Content:    request.GetContent(),
		CategoryID: request.GetCategoryId(),
		ParentID:   request.GetParentId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &reviewspb.AddReviewResponse{
		Id: reviewID,
	}, nil
}
func (s server) ApproveReview(ctx context.Context, request *reviewspb.ApproveReviewRequest) (*reviewspb.ApproveReviewResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("ReviewID", request.GetId()),
	)

	err := s.app.ApproveReview(ctx, commands.ApproveReview{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &reviewspb.ApproveReviewResponse{}, nil
}
func (s server) RejectReview(ctx context.Context, request *reviewspb.RejectReviewRequest) (*reviewspb.RejectReviewResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("ReviewID", request.GetId()),
	)

	err := s.app.RejectReview(ctx, commands.RejectReview{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &reviewspb.RejectReviewResponse{}, nil
}
func (s server) FlagReview(ctx context.Context, request *reviewspb.FlagReviewRequest) (*reviewspb.FlagReviewResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("ReviewID", request.GetId()),
	)

	err := s.app.FlagReview(ctx, commands.FlagReview{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &reviewspb.FlagReviewResponse{}, nil
}
func (s server) EditReview(ctx context.Context, request *reviewspb.EditReviewRequest) (*reviewspb.EditReviewResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("ReviewID", request.GetId()),
	)

	err := s.app.EditReview(ctx, commands.EditReview{
		ID:      request.GetId(),
		Content: request.GetContent(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &reviewspb.EditReviewResponse{}, nil
}
func (s server) RemoveReview(ctx context.Context, request *reviewspb.RemoveReviewRequest) (*reviewspb.RemoveReviewResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("ReviewID", request.GetId()),
	)

	err := s.app.RemoveReview(ctx, commands.RemoveReview{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &reviewspb.RemoveReviewResponse{}, nil
}
func (s server) GetReviews(ctx context.Context, request *reviewspb.GetReviewsRequest) (*reviewspb.GetReviewsResponse, error) {
	span := trace.SpanFromContext(ctx)

	reviews, err := s.app.GetReviews(ctx, queries.GetReviews{
		ItemID: request.ItemId,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoReviews := []*reviewspb.Review{}
	for _, review := range reviews {
		protoReviews = append(protoReviews, s.reviewFromDomain(review))
	}

	return &reviewspb.GetReviewsResponse{
		Reviews: protoReviews,
	}, nil
}

func (s server) GetReviewsBySender(ctx context.Context, request *reviewspb.GetReviewsBySenderRequest) (*reviewspb.GetReviewsBySenderResponse, error) {
	span := trace.SpanFromContext(ctx)

	reviews, err := s.app.GetReviewsBySender(ctx, queries.GetReviewsBySender{
		SenderID: request.SenderId,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoReviews := []*reviewspb.Review{}
	for _, review := range reviews {
		protoReviews = append(protoReviews, s.reviewFromDomain(review))
	}

	return &reviewspb.GetReviewsBySenderResponse{
		Reviews: protoReviews,
	}, nil
}
func (s server) GetApprovedReviews(ctx context.Context, request *reviewspb.GetApprovedReviewsRequest) (*reviewspb.GetApprovedReviewsResponse, error) {
	span := trace.SpanFromContext(ctx)

	reviews, err := s.app.GetApprovedReviews(ctx, queries.GetApprovedReviews{})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoReviews := []*reviewspb.Review{}
	for _, review := range reviews {
		protoReviews = append(protoReviews, s.reviewFromDomain(review))
	}

	return &reviewspb.GetApprovedReviewsResponse{
		Reviews: protoReviews,
	}, nil
}
func (s server) GetReview(ctx context.Context, request *reviewspb.GetReviewRequest) (*reviewspb.GetReviewResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("ReviewID", request.GetId()),
	)

	review, err := s.app.GetReview(ctx, queries.GetReview{ID: request.GetId()})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &reviewspb.GetReviewResponse{Review: s.reviewFromDomain(review)}, nil
}

func (s server) reviewFromDomain(review *domain.MiddlemanReview) *reviewspb.Review {
	return &reviewspb.Review{
		Id:         review.ID,
		SenderId:   review.SenderID,
		ItemId:     review.ItemID,
		ItemType:   review.ItemType,
		Content:    review.Content,
		CategoryId: review.CategoryID,
		ParentId:   review.ParentID,
	}
}
