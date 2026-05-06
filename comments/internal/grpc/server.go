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
	"middleman/comments/commentspb"
	"middleman/comments/internal/application"
	"middleman/comments/internal/application/commands"
	"middleman/comments/internal/application/queries"
	"middleman/comments/internal/domain"
	"middleman/internal/auth"
	"middleman/internal/errorsotel"
)

type server struct {
	app application.App
	commentspb.UnimplementedCommentsServiceServer
}

var _ commentspb.CommentsServiceServer = (*server)(nil)

func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	commentspb.RegisterCommentsServiceServer(registrar, server{app: app})
	return nil
}
func (s server) AddComment(ctx context.Context, request *commentspb.AddCommentRequest) (*commentspb.AddCommentResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	span := trace.SpanFromContext(ctx)
	commentID := uuid.New().String()

	span.SetAttributes(
		attribute.String("CommentID", commentID),
	)
	fmt.Println("SenderID", userID)
	fmt.Println("ItemID", request.ItemId)
	fmt.Println("ItemType", request.ItemType)
	fmt.Println("Content", request.Content)
	fmt.Println("CategoryID", request.CategoryId)
	fmt.Println("ParentID", request.ParentId)

	err := s.app.AddComment(ctx, commands.AddComment{
		ID:         commentID,
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

	return &commentspb.AddCommentResponse{
		Id: commentID,
	}, nil
}
func (s server) ApproveComment(ctx context.Context, request *commentspb.ApproveCommentRequest) (*commentspb.ApproveCommentResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("CommentID", request.GetId()),
	)
	_, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	err := s.app.ApproveComment(ctx, commands.ApproveComment{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &commentspb.ApproveCommentResponse{}, nil
}
func (s server) RejectComment(ctx context.Context, request *commentspb.RejectCommentRequest) (*commentspb.RejectCommentResponse, error) {
	_, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("CommentID", request.GetId()),
	)

	err := s.app.RejectComment(ctx, commands.RejectComment{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &commentspb.RejectCommentResponse{}, nil
}
func (s server) FlagComment(ctx context.Context, request *commentspb.FlagCommentRequest) (*commentspb.FlagCommentResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("CommentID", request.GetId()),
	)

	err := s.app.FlagComment(ctx, commands.FlagComment{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &commentspb.FlagCommentResponse{}, nil
}
func (s server) EditComment(ctx context.Context, request *commentspb.EditCommentRequest) (*commentspb.EditCommentResponse, error) {
	_, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("CommentID", request.GetId()),
	)

	err := s.app.EditComment(ctx, commands.EditComment{
		ID:      request.GetId(),
		Content: request.GetContent(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &commentspb.EditCommentResponse{}, nil
}
func (s server) RemoveComment(ctx context.Context, request *commentspb.RemoveCommentRequest) (*commentspb.RemoveCommentResponse, error) {
	_, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("CommentID", request.GetId()),
	)

	err := s.app.RemoveComment(ctx, commands.RemoveComment{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &commentspb.RemoveCommentResponse{}, nil
}
func (s server) GetComments(ctx context.Context, request *commentspb.GetCommentsRequest) (*commentspb.GetCommentsResponse, error) {
	span := trace.SpanFromContext(ctx)

	comments, err := s.app.GetComments(ctx, queries.GetComments{
		ItemID: request.ItemId,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoComments := []*commentspb.Comment{}
	for _, comment := range comments {
		protoComments = append(protoComments, s.commentFromDomain(comment))
	}

	return &commentspb.GetCommentsResponse{
		Comments: protoComments,
	}, nil
}

func (s server) GetCommentsBySender(ctx context.Context, request *commentspb.GetCommentsBySenderRequest) (*commentspb.GetCommentsBySenderResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	span := trace.SpanFromContext(ctx)
	comments, err := s.app.GetCommentsBySender(ctx, queries.GetCommentsBySender{
		SenderID: userID,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoComments := []*commentspb.Comment{}
	for _, comment := range comments {
		protoComments = append(protoComments, s.commentFromDomain(comment))
	}

	return &commentspb.GetCommentsBySenderResponse{
		Comments: protoComments,
	}, nil
}
func (s server) GetApprovedComments(ctx context.Context, request *commentspb.GetApprovedCommentsRequest) (*commentspb.GetApprovedCommentsResponse, error) {
	span := trace.SpanFromContext(ctx)
	comments, err := s.app.GetApprovedComments(ctx, queries.GetApprovedComments{})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoComments := []*commentspb.Comment{}
	for _, comment := range comments {
		protoComments = append(protoComments, s.commentFromDomain(comment))
	}

	return &commentspb.GetApprovedCommentsResponse{
		Comments: protoComments,
	}, nil
}
func (s server) GetComment(ctx context.Context, request *commentspb.GetCommentRequest) (*commentspb.GetCommentResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("CommentID", request.GetId()),
	)

	comment, err := s.app.GetComment(ctx, queries.GetComment{ID: request.GetId()})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &commentspb.GetCommentResponse{Comment: s.commentFromDomain(comment)}, nil
}

func (s server) commentFromDomain(comment *domain.MiddlemanComment) *commentspb.Comment {
	return &commentspb.Comment{
		Id:         comment.ID,
		SenderId:   comment.SenderID,
		ItemId:     comment.ItemID,
		ItemType:   comment.ItemType,
		Content:    comment.Content,
		CategoryId: comment.CategoryID,
		ParentId:   comment.ParentID,
	}
}
