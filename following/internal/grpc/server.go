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
	"middleman/following/followingpb"
	"middleman/following/internal/application"
	"middleman/following/internal/application/commands"
	"middleman/following/internal/application/queries"
	"middleman/following/internal/domain"
	"middleman/internal/auth"
	"middleman/internal/errorsotel"
)

type server struct {
	app application.App
	followingpb.UnimplementedFollowingServiceServer
}

var _ followingpb.FollowingServiceServer = (*server)(nil)

func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	followingpb.RegisterFollowingServiceServer(registrar, server{app: app})
	return nil
}
func (s server) AddFollow(ctx context.Context, request *followingpb.AddFollowRequest) (*followingpb.AddFollowResponse, error) {

	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	span := trace.SpanFromContext(ctx)

	followID := uuid.New().String()

	span.SetAttributes(
		attribute.String("FollowID", followID),
	)
	fmt.Println("UserID", userID)
	fmt.Println("FollowedUserID", request.FollowedUserId)
	fmt.Println("FollowedUserType", request.FollowedUserType)
	fmt.Println("Content", request.Content)
	fmt.Println("CategoryID", request.CategoryId)
	fmt.Println("ParentID", request.ParentId)

	err := s.app.AddFollow(ctx, commands.AddFollow{
		ID:               followID,
		UserID:           userID,
		FollowedUserID:   request.GetFollowedUserId(),
		FollowedUserType: domain.ToFollowedUserType(request.GetFollowedUserType()),
		Content:          request.GetContent(),
		CategoryID:       request.GetCategoryId(),
		ParentID:         request.GetParentId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &followingpb.AddFollowResponse{
		Id: followID,
	}, nil
}
func (s server) ApproveFollow(ctx context.Context, request *followingpb.ApproveFollowRequest) (*followingpb.ApproveFollowResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("FollowID", request.GetId()),
	)

	err := s.app.ApproveFollow(ctx, commands.ApproveFollow{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &followingpb.ApproveFollowResponse{}, nil
}
func (s server) RejectFollow(ctx context.Context, request *followingpb.RejectFollowRequest) (*followingpb.RejectFollowResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("FollowID", request.GetId()),
	)

	err := s.app.RejectFollow(ctx, commands.RejectFollow{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &followingpb.RejectFollowResponse{}, nil
}
func (s server) FlagFollow(ctx context.Context, request *followingpb.FlagFollowRequest) (*followingpb.FlagFollowResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("FollowID", request.GetId()),
	)

	err := s.app.FlagFollow(ctx, commands.FlagFollow{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &followingpb.FlagFollowResponse{}, nil
}

func (s server) RemoveFollow(ctx context.Context, request *followingpb.RemoveFollowRequest) (*followingpb.RemoveFollowResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("FollowID", request.GetId()),
	)

	err := s.app.RemoveFollow(ctx, commands.RemoveFollow{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &followingpb.RemoveFollowResponse{}, nil
}
func (s server) GetFollowing(ctx context.Context, request *followingpb.GetFollowingRequest) (*followingpb.GetFollowingResponse, error) {
	span := trace.SpanFromContext(ctx)

	following, err := s.app.GetFollowing(ctx, queries.GetFollowing{
		FollowedUserID: request.FollowedUserId,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoFollowing := []*followingpb.Follow{}
	for _, follow := range following {
		protoFollowing = append(protoFollowing, s.followFromDomain(follow))
	}

	return &followingpb.GetFollowingResponse{
		Following: protoFollowing,
	}, nil
}

func (s server) GetFollowingBySender(ctx context.Context, request *followingpb.GetFollowingBySenderRequest) (*followingpb.GetFollowingBySenderResponse, error) {
	span := trace.SpanFromContext(ctx)
	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	following, err := s.app.GetFollowingBySender(ctx, queries.GetFollowingBySender{
		UserID: userID,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoFollowing := []*followingpb.Follow{}
	for _, follow := range following {
		protoFollowing = append(protoFollowing, s.followFromDomain(follow))
	}

	return &followingpb.GetFollowingBySenderResponse{
		Following: protoFollowing,
	}, nil
}
func (s server) GetApprovedFollowing(ctx context.Context, request *followingpb.GetApprovedFollowingRequest) (*followingpb.GetApprovedFollowingResponse, error) {
	span := trace.SpanFromContext(ctx)

	following, err := s.app.GetApprovedFollowing(ctx, queries.GetApprovedFollowing{})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoFollowing := []*followingpb.Follow{}
	for _, follow := range following {
		protoFollowing = append(protoFollowing, s.followFromDomain(follow))
	}

	return &followingpb.GetApprovedFollowingResponse{
		Following: protoFollowing,
	}, nil
}
func (s server) GetFollow(ctx context.Context, request *followingpb.GetFollowRequest) (*followingpb.GetFollowResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("FollowID", request.GetId()),
	)

	follow, err := s.app.GetFollow(ctx, queries.GetFollow{ID: request.GetId()})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &followingpb.GetFollowResponse{Follow: s.followFromDomain(follow)}, nil
}

func (s server) followFromDomain(follow *domain.MiddlemanFollow) *followingpb.Follow {
	return &followingpb.Follow{
		Id:               follow.ID,
		UserId:           follow.UserID,
		FollowedUserId:   follow.FollowedUserID,
		FollowedUserType: follow.FollowedUserType,
		Content:          follow.Content,
		CategoryId:       follow.CategoryID,
		ParentId:         follow.ParentID,
	}
}
