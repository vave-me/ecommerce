package grpc

import (
	"context"
	"fmt"
	"middleman/baskets/basketspb"
	"middleman/baskets/internal/application"
	"middleman/baskets/internal/domain"
	"middleman/internal/auth"
	"middleman/internal/errorsotel"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpc_code "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	app application.App
	basketspb.UnimplementedBasketServiceServer
}

var _ basketspb.BasketServiceServer = (*server)(nil)

func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	basketspb.RegisterBasketServiceServer(registrar, server{app: app})
	return nil
}

func (s server) StartBasket(ctx context.Context, request *basketspb.StartBasketRequest) (*basketspb.StartBasketResponse, error) {
	// Retrieve the user claims from the context
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	//TODO check if user from request is same user who send jwt
	fmt.Printf("user ID %s", userID)

	// Generate a new basket ID
	basketID := uuid.New().String()

	// Call the application layer to start a basket
	err := s.app.StartBasket(ctx, application.StartBasket{
		ID:             basketID,
		UserCustomerID: userID,
	})
	if err != nil {
		// Handle the error appropriately
		return nil, status.Error(grpc_code.Internal, err.Error())
	}

	// Return the response with the new basket ID
	return &basketspb.StartBasketResponse{Id: basketID}, nil
}

func (s server) CancelBasket(ctx context.Context, request *basketspb.CancelBasketRequest,
) (*basketspb.CancelBasketResponse, error) {
	err := s.app.CancelBasket(ctx, application.CancelBasket{
		ID: request.GetBasketId(),
	})

	return &basketspb.CancelBasketResponse{}, err
}

func (s server) CheckoutBasket(ctx context.Context, request *basketspb.CheckoutBasketRequest) (*basketspb.CheckoutBasketResponse, error) {

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	err := s.app.CheckoutBasket(ctx, application.CheckoutBasket{
		ID:              request.GetBasketId(),
		UserCustomerID:  userID,
		PaymentIntentID: request.GetPaymentIntentId(), // Include payment intent ID from frontend
	})

	if err != nil {
		return nil, err
	}

	// Return the response with the basket ID that was checked out
	return &basketspb.CheckoutBasketResponse{
		BasketId: request.GetBasketId(),
	}, nil
}

func (s server) ReopenBasket(ctx context.Context, request *basketspb.ReopenBasketRequest) (*basketspb.ReopenBasketResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("BasketID", request.GetBasketId()),
		attribute.String("Reason", request.GetReason()),
	)
	
	err := s.app.ReopenBasket(ctx, application.ReopenBasket{
		ID:     request.GetBasketId(),
		Reason: request.GetReason(),
	})
	
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	
	return &basketspb.ReopenBasketResponse{
		BasketId:     request.GetBasketId(),
		BasketStatus: "open",
	}, nil
}

func (s server) AddItem(ctx context.Context, request *basketspb.AddItemRequest) (*basketspb.AddItemResponse, error) {
	span := trace.SpanFromContext(ctx)
	basketID := request.GetBasketId()
	span.SetAttributes(
		attribute.String("BasketID", basketID),
	)
	err := s.app.AddItem(ctx, application.AddItem{
		ID:        basketID,
		ProductID: request.GetProductId(),
		Quantity:  request.GetQuantity(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return &basketspb.AddItemResponse{}, err
}

func (s server) RemoveItem(ctx context.Context, request *basketspb.RemoveItemRequest) (*basketspb.RemoveItemResponse,
	error,
) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("ItemID", request.GetItemId()),
	)
	//TODO ADD CHECK TOKEN
	err := s.app.RemoveItem(ctx, application.RemoveItem{
		ID:        request.GetItemId(),
		ProductID: request.GetItemId(),
		Quantity:  request.GetQuantity(),
	})

	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
	}

	return &basketspb.RemoveItemResponse{}, err
}

func (s server) GetBasket(ctx context.Context, request *basketspb.GetBasketRequest) (*basketspb.GetBasketResponse,
	error,
) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	basket, err := s.app.GetBasket(ctx, application.GetBasket{
		ID:             request.GetBasketId(),
		UserCustomerID: userID,
	})
	if err != nil {
		return nil, err
	}

	return &basketspb.GetBasketResponse{
		Basket: s.basketFromDomain(basket),
	}, nil
}

func (s server) GetCurrentBasket(ctx context.Context, request *basketspb.GetCurrentBasketRequest) (*basketspb.GetCurrentBasketResponse,
	error,
) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	basket, err := s.app.GetCurrentBasket(ctx, application.GetCurrentBasket{
		UserCustomerID: userID,
	})
	if err != nil {
		return nil, err
	}

	return &basketspb.GetCurrentBasketResponse{
		BasketId:     basket.ID,
		BasketStatus: string(basket.Status),
	}, nil
}

func (s server) GetTotalBasket(ctx context.Context, request *basketspb.GetTotalBasketAmountRequest) (*basketspb.GetTotalBasketAmountResponse,
	error,
) {
	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	amount, err := s.app.GetTotalBasket(ctx, application.GetTotalBasket{
		ID: request.GetBasketId(),
	})
	if err != nil {
		return nil, err
	}

	return &basketspb.GetTotalBasketAmountResponse{
		Amount: amount,
	}, nil
}

func (s server) basketFromDomain(basket *domain.Basket) *basketspb.Basket {
	protoBasket := &basketspb.Basket{
		Id: basket.ID(),
	}

	protoBasket.Items = make([]*basketspb.Item, 0, len(basket.Items))

	for _, item := range basket.Items {
		protoBasket.Items = append(protoBasket.Items, &basketspb.Item{
			UserSellerId:   item.UserSellerID,
			UserSellerName: item.UserSellerName,
			ProductId:      item.ProductID,
			ProductName:    item.ProductName,
			ProductPrice:   item.ProductPrice,
			Quantity:       item.Quantity,
		})
	}

	return protoBasket
}
