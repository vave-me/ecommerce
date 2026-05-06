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
	"middleman/internal/auth"
	"middleman/internal/errorsotel"
	"middleman/wishlists/internal/application"
	"middleman/wishlists/internal/application/commands"
	"middleman/wishlists/internal/application/queries"
	"middleman/wishlists/internal/domain"
	"middleman/wishlists/wishlistspb"
)

type server struct {
	app application.App
	wishlistspb.UnimplementedWishlistServiceServer
}

var _ wishlistspb.WishlistServiceServer = (*server)(nil)

func RegisterServer(_ context.Context, app application.App, registrar grpc.ServiceRegistrar) error {
	wishlistspb.RegisterWishlistServiceServer(registrar, server{app: app})
	return nil
}

func (s server) CreateWishlist(ctx context.Context, request *wishlistspb.CreateWishlistRequest) (*wishlistspb.CreateWishlistResponse, error) {
	span := trace.SpanFromContext(ctx)

	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject

	wishlistID := uuid.New().String()
	span.SetAttributes(
		attribute.String("Wishlist", wishlistID),
	)

	err := s.app.CreateWishlist(ctx, commands.CreateWishlist{
		ID:     wishlistID,
		UserID: userID,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &wishlistspb.CreateWishlistResponse{
		Id: wishlistID,
	}, nil
}
func (s server) RemoveWishlist(ctx context.Context, request *wishlistspb.RemoveWishlistRequest) (*wishlistspb.RemoveWishlistResponse, error) {
	span := trace.SpanFromContext(ctx)

	_, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	wishlistID := uuid.New().String()
	span.SetAttributes(
		attribute.String("Wishlist", wishlistID),
	)

	err := s.app.RemoveWishlist(ctx, commands.RemoveWishlist{
		ID: wishlistID,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &wishlistspb.RemoveWishlistResponse{}, nil
}

func (s server) AddWishlistItem(ctx context.Context, request *wishlistspb.AddWishlistItemRequest) (*wishlistspb.AddWishlistItemResponse, error) {

	_, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	span := trace.SpanFromContext(ctx)
	id := uuid.New().String()

	span.SetAttributes(
		attribute.String("WishlistItemID", id),
	)

	err := s.app.AddWishlistItem(ctx, commands.AddWishlistItem{
		ID:         id,
		WishlistID: request.GetWishlistId(),
		ItemID:     request.GetItemId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &wishlistspb.AddWishlistItemResponse{Id: id}, nil
}

func (s server) RemoveWishlistItem(ctx context.Context, request *wishlistspb.RemoveWishlistItemRequest) (*wishlistspb.RemoveWishlistItemResponse, error) {

	_, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	span := trace.SpanFromContext(ctx)

	// Retrieve the user claims from the context

	span.SetAttributes(
		attribute.String("WishlistItem", request.GetId()),
	)

	err := s.app.RemoveWishlistItem(ctx, commands.RemoveWishlistItem{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
	}

	return &wishlistspb.RemoveWishlistItemResponse{}, err
}
func (s server) GetWishlist(ctx context.Context, request *wishlistspb.GetWishlistRequest) (*wishlistspb.GetWishlistResponse, error) {
	span := trace.SpanFromContext(ctx)

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject

	span.SetAttributes(
		attribute.String("UserID", userID),
	)

	wishlist, err := s.app.GetWishlist(ctx, queries.GetWishlist{
		UserID: userID,
		Name:   request.GetName()})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &wishlistspb.GetWishlistResponse{WishlistId: wishlist.ID}, nil
}

func (s server) GetWishlists(ctx context.Context, request *wishlistspb.GetWishlistsRequest) (*wishlistspb.GetWishlistsResponse, error) {
	span := trace.SpanFromContext(ctx)

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject

	span.SetAttributes(
		attribute.String("UserID", userID),
	)

	wishlists, err := s.app.GetWishlists(ctx, queries.GetWishlists{UserID: userID})

	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	protoWishlists := make([]*wishlistspb.Wishlist, len(wishlists))
	for i, wishlist := range wishlists {
		protoWishlists[i] = s.wishlistFromDomain(wishlist)
	}

	return &wishlistspb.GetWishlistsResponse{Wishlists: protoWishlists}, nil
}
func (s server) wishlistFromDomain(wishlist *domain.MiddlemanWishlist) *wishlistspb.Wishlist {
	return &wishlistspb.Wishlist{
		Id:     wishlist.ID,
		UserId: wishlist.UserID,
	}
}
func (s server) GetWishlistItem(ctx context.Context, request *wishlistspb.GetWishlistItemRequest) (*wishlistspb.GetWishlistItemResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("WishlistItemID", request.GetId()),
	)

	item, err := s.app.GetWishlistItem(ctx, queries.GetWishlistItem{
		WishlistItemID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &wishlistspb.GetWishlistItemResponse{Item: s.wishlistItemFromDomain(item)}, nil
}

func (s server) GetWishlistItems(ctx context.Context, request *wishlistspb.GetWishlistItemsRequest) (*wishlistspb.GetWishlistItemsResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("WishlistID", request.GetWishlistId()),
	)

	items, err := s.app.GetWishlistItems(ctx, queries.GetWishlistItems{WishlistID: request.GetWishlistId()})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoItems := make([]*wishlistspb.WishlistItem, len(items))
	for i, item := range items {
		protoItems[i] = s.wishlistItemFromDomain(item)
	}

	return &wishlistspb.GetWishlistItemsResponse{
		Items: protoItems,
	}, nil
}

func (s server) wishlistItemFromDomain(wishlistItem *domain.CatalogWishlistItem) *wishlistspb.WishlistItem {
	return &wishlistspb.WishlistItem{
		Id:         wishlistItem.ID,
		WishlistId: wishlistItem.WishlistID,
		ItemId:     wishlistItem.ItemID,
		EntityType: wishlistItem.EntityType,
	}
}
