package grpc

import (
	"context"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpc_code "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"middleman/activity/activitypb"
	"middleman/activity/internal/application"
	"middleman/activity/internal/application/commands"
	"middleman/activity/internal/application/queries"
	"middleman/activity/internal/domain"
	"middleman/internal/auth"
	"middleman/internal/errorsotel"
)

type server struct {
	app application.App
	activitypb.UnimplementedActivityServiceServer
}

var _ activitypb.ActivityServiceServer = (*server)(nil)

func RegisterServer(_ context.Context, app application.App, registrar grpc.ServiceRegistrar) error {
	activitypb.RegisterActivityServiceServer(registrar, server{app: app})
	return nil
}

func (s server) CreateActivity(ctx context.Context, request *activitypb.CreateActivityRequest) (*activitypb.CreateActivityResponse, error) {
	span := trace.SpanFromContext(ctx)
	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	activityID := uuid.New().String()

	span.SetAttributes(
		attribute.String("ActivityID", activityID),
	)

	err := s.app.CreateActivity(ctx, commands.CreateActivity{
		ID:     activityID,
		UserID: userID,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &activitypb.CreateActivityResponse{
		Id:     activityID,
		UserId: userID,
	}, nil
}
func (s server) AddInteraction(ctx context.Context, request *activitypb.AddInteractionRequest) (*activitypb.AddInteractionResponse, error) {
	span := trace.SpanFromContext(ctx)
	id := uuid.New().String()

	_, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	activityID := request.GetActivityId()
	actionType := request.GetActionType()
	itemID := request.GetItemId()
	itemType := request.GetItemType()

	// Log the incoming request details
	log.Printf("AddInteraction: Received request with ActivityID=%s, ActionType=%s, ItemID=%s, ItemType=%s", activityID, actionType, itemID, itemType)

	span.SetAttributes(
		attribute.String("ActivityID", activityID),
		attribute.String("ActionType", actionType),
		attribute.String("ItemID", itemID),
		attribute.String("ItemType", itemType),
		attribute.String("RequestID", id),
	)

	// Execute the AddInteraction command
	log.Printf("AddInteraction: Adding interaction with ID=%s", id)
	err := s.app.AddInteraction(ctx, commands.AddInteraction{
		ID:         id,
		ActivityID: activityID,
		ActionType: actionType,
		ItemID:     itemID,
		ItemType:   itemType,
	})
	if err != nil {
		// Log the error details
		log.Printf("AddInteraction: Error adding interaction with ID=%s, ActivityID=%s: %v", id, activityID, err)

		// Record error in tracing span
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Log success
	log.Printf("AddInteraction: Successfully added interaction with ID=%s, ActivityID=%s", id, activityID)

	return &activitypb.AddInteractionResponse{
		Id: id,
	}, nil
}

func (s server) GetActivity(ctx context.Context, request *activitypb.GetActivityRequest) (*activitypb.GetActivityResponse, error) {
	span := trace.SpanFromContext(ctx)
	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	span.SetAttributes(
		attribute.String("UserID", userID),
	)
	// TODO add check with claims if the userID is the same from response and request
	activity, err := s.app.GetActivity(ctx, queries.GetActivity{UserID: userID})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &activitypb.GetActivityResponse{ActivityId: activity.ID, UserId: activity.UserID}, nil
}

func (s server) GetInteraction(ctx context.Context, request *activitypb.GetInteractionRequest) (*activitypb.GetInteractionResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("InteractionID", request.GetId()),
	)

	interaction, err := s.app.GetInteraction(ctx, queries.GetInteraction{ID: request.GetId()})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &activitypb.GetInteractionResponse{Interaction: s.interactionFromDomain(interaction)}, nil
}

func (s server) RemoveInteraction(ctx context.Context, request *activitypb.RemoveInteractionRequest) (*activitypb.RemoveInteractionResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("InteractionID", request.GetId()),
	)

	err := s.app.RemoveInteraction(ctx, commands.RemoveInteraction{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
	}

	return &activitypb.RemoveInteractionResponse{}, err
}

func (s server) GetInteractions(ctx context.Context, request *activitypb.GetInteractionsRequest) (*activitypb.GetInteractionsResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("ActivityID", request.GetActivityId()),
	)

	interactions, err := s.app.GetInteractions(ctx, queries.GetInteractions{ActivityID: request.GetActivityId()})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoInteractions := make([]*activitypb.Interaction, len(interactions))
	for i, interaction := range interactions {
		protoInteractions[i] = s.interactionFromDomain(interaction)
	}

	return &activitypb.GetInteractionsResponse{
		Interactions: protoInteractions,
	}, nil
}

func (s server) GetMostLiked(ctx context.Context, request *activitypb.GetMostLikedRequest) (*activitypb.GetMostLikedResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("Item Type", request.GetItemType()),
	)

	interactions, err := s.app.GetMostLiked(ctx, queries.GetMostLiked{ItemType: request.GetItemType(), Limit: request.GetLimit()})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoInteractions := make([]*activitypb.MostReactionResult, len(interactions))
	for i, interaction := range interactions {
		protoInteractions[i] = s.itemInteractionCountFromDomain(interaction)
	}

	return &activitypb.GetMostLikedResponse{
		Results: protoInteractions,
	}, nil
}
func (s server) GetMostDisliked(ctx context.Context, request *activitypb.GetMostDislikedRequest) (*activitypb.GetMostDislikedResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("Item Type", request.GetItemType()),
	)

	interactions, err := s.app.GetMostLiked(ctx, queries.GetMostLiked{ItemType: request.GetItemType(), Limit: request.GetLimit()})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoInteractions := make([]*activitypb.MostReactionResult, len(interactions))
	for i, interaction := range interactions {
		protoInteractions[i] = s.itemInteractionCountFromDomain(interaction)
	}

	return &activitypb.GetMostDislikedResponse{
		Results: protoInteractions,
	}, nil
}
func (s server) interactionFromDomain(interaction *domain.MiddlemanInteraction) *activitypb.Interaction {
	return &activitypb.Interaction{
		Id:         interaction.ID,
		ActivityId: interaction.ActivityID,
		ItemId:     interaction.ItemID,
		ItemType:   interaction.ItemType,
		ActionType: interaction.ActionType,
	}
}

func (s server) itemInteractionCountFromDomain(interaction *domain.MostReactionResult) *activitypb.MostReactionResult {
	return &activitypb.MostReactionResult{
		ItemId:   interaction.ItemID,
		Action:   interaction.Action,
		ItemType: interaction.ItemType,
	}
}
func (s server) activityFromDomain(activity *domain.MiddlemanActivity) *activitypb.Activity {
	return &activitypb.Activity{
		Id:     activity.ID,
		UserId: activity.UserID,
	}
}
