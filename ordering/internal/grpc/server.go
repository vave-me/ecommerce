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
	"middleman/ordering/internal/application"
	"middleman/ordering/internal/application/commands"
	"middleman/ordering/internal/application/queries"
	"middleman/ordering/internal/domain"
	"middleman/ordering/orderingpb"
)

type server struct {
	app application.App
	orderingpb.UnimplementedOrderingServiceServer
}

var _ orderingpb.OrderingServiceServer = (*server)(nil)

func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	orderingpb.RegisterOrderingServiceServer(registrar, server{app: app})
	return nil
}

func (s server) CreateOrder(ctx context.Context, request *orderingpb.CreateOrderRequest) (*orderingpb.CreateOrderResponse, error) {
	span := trace.SpanFromContext(ctx)

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	id := uuid.New().String()

	span.SetAttributes(
		attribute.String("OrderID", id),
		attribute.String("UserCustomerID", userID),
	)

	items := make([]domain.Item, len(request.Items))
	for i, item := range request.Items {
		items[i] = s.itemToDomain(item)
	}

	err := s.app.CreateOrder(ctx, commands.CreateOrder{
		ID:             id,
		UserCustomerID: userID,
		BasketID:       request.GetBasketId(),
		PaymentIntent:  request.GetPaymentIntent(),
		Items:          items,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
	}

	return &orderingpb.CreateOrderResponse{Id: id}, err
}

func (s server) CancelOrder(ctx context.Context, request *orderingpb.CancelOrderRequest) (*orderingpb.CancelOrderResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("OrderID", request.GetId()),
	)

	err := s.app.CancelOrder(ctx, commands.CancelOrder{ID: request.GetId()})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
	}

	return &orderingpb.CancelOrderResponse{}, err
}

func (s server) ReadyOrder(ctx context.Context, request *orderingpb.ReadyOrderRequest) (*orderingpb.ReadyOrderResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("OrderID", request.GetId()),
	)

	err := s.app.ReadyOrder(ctx, commands.ReadyOrder{ID: request.GetId()})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
	}

	return &orderingpb.ReadyOrderResponse{}, err
}
func (s server) DeliverOrder(ctx context.Context, req *orderingpb.DeliverOrderRequest) (*orderingpb.DeliverOrderResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("OrderID", req.GetId()))

	err := s.app.DeliverOrder(ctx, commands.DeliverOrder{ID: req.GetId()})
	if err != nil {
		recordSpanError(span, err)
		return nil, err
	}

	return &orderingpb.DeliverOrderResponse{}, nil
}
func (s server) CompleteOrder(ctx context.Context, request *orderingpb.CompleteOrderRequest) (*orderingpb.CompleteOrderResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("OrderID", request.GetId()),
	)

	err := s.app.CompleteOrder(ctx, commands.CompleteOrder{ID: request.GetId()})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
	}

	return &orderingpb.CompleteOrderResponse{}, err
}
func (s server) ApproveOrder(ctx context.Context, req *orderingpb.ApproveOrderRequest) (*orderingpb.ApproveOrderResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("OrderID", req.GetId()), attribute.String("ShoppingID", req.GetShoppingId()))

	err := s.app.ApproveOrder(ctx, commands.ApproveOrder{
		ID:         req.GetId(),
		ShoppingID: req.GetShoppingId(),
	})
	if err != nil {
		recordSpanError(span, err)
		return nil, err
	}

	return &orderingpb.ApproveOrderResponse{}, nil
}

func (s server) RejectOrder(ctx context.Context, req *orderingpb.RejectOrderRequest) (*orderingpb.RejectOrderResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("OrderID", req.GetId()))

	err := s.app.RejectOrder(ctx, commands.RejectOrder{ID: req.GetId()})
	if err != nil {
		recordSpanError(span, err)
		return nil, err
	}

	return &orderingpb.RejectOrderResponse{}, nil
}
func (s server) ShipOrder(ctx context.Context, req *orderingpb.ShipOrderRequest) (*orderingpb.ShipOrderResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("OrderID", req.GetId()))

	err := s.app.ShipOrder(ctx, commands.ShipOrder{ID: req.GetId()})
	if err != nil {
		recordSpanError(span, err)
		return nil, err
	}

	return &orderingpb.ShipOrderResponse{}, nil
}
func (s server) GetOrder(ctx context.Context, request *orderingpb.GetOrderRequest) (*orderingpb.GetOrderResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("OrderID", request.GetId()),
	)

	order, err := s.app.GetOrder(ctx, queries.GetOrder{ID: request.GetId()})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &orderingpb.GetOrderResponse{
		Order: s.orderFromDomain(order),
	}, nil
}

func (s server) ListOrders(ctx context.Context, request *orderingpb.ListOrdersRequest) (*orderingpb.ListOrdersResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.Int64("Page", request.GetPage()),
		attribute.Int64("PageSize", request.GetPageSize()),
	)

	orders, totalCount, err := s.app.ListOrders(ctx, queries.ListOrders{
		Page:      request.GetPage(),
		PageSize:  request.GetPageSize(),
		SortBy:    request.GetSortBy(),
		SortOrder: request.GetSortOrder(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Calculate total pages
	pageSize := request.GetPageSize()
	if pageSize <= 0 {
		pageSize = 10
	}
	totalPages := (totalCount + pageSize - 1) / pageSize

	return &orderingpb.ListOrdersResponse{
		Orders:      s.ordersFromCatalog(orders),
		TotalCount:  totalCount,
		TotalPages:  totalPages,
		CurrentPage: request.GetPage(),
	}, nil
}

func (s server) GetOrdersByCustomer(ctx context.Context, request *orderingpb.GetOrdersByCustomerRequest) (*orderingpb.GetOrdersByCustomerResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("UserCustomerID", request.GetUserCustomerId()),
		attribute.Int64("Page", request.GetPage()),
		attribute.Int64("PageSize", request.GetPageSize()),
	)

	orders, totalCount, err := s.app.GetOrdersByCustomer(ctx, queries.GetOrdersByCustomer{
		UserCustomerID: request.GetUserCustomerId(),
		Page:           request.GetPage(),
		PageSize:       request.GetPageSize(),
		SortBy:         request.GetSortBy(),
		SortOrder:      request.GetSortOrder(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Calculate total pages
	pageSize := request.GetPageSize()
	if pageSize <= 0 {
		pageSize = 10
	}
	totalPages := (totalCount + pageSize - 1) / pageSize

	return &orderingpb.GetOrdersByCustomerResponse{
		Orders:      s.ordersFromCatalog(orders),
		TotalCount:  totalCount,
		TotalPages:  totalPages,
		CurrentPage: request.GetPage(),
	}, nil
}

func (s server) GetOrdersByStatus(ctx context.Context, request *orderingpb.GetOrdersByStatusRequest) (*orderingpb.GetOrdersByStatusResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("Status", request.GetStatus()),
		attribute.Int64("Page", request.GetPage()),
		attribute.Int64("PageSize", request.GetPageSize()),
	)

	orders, totalCount, err := s.app.GetOrdersByStatus(ctx, queries.GetOrdersByStatus{
		Status:    request.GetStatus(),
		Page:      request.GetPage(),
		PageSize:  request.GetPageSize(),
		SortBy:    request.GetSortBy(),
		SortOrder: request.GetSortOrder(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Calculate total pages
	pageSize := request.GetPageSize()
	if pageSize <= 0 {
		pageSize = 10
	}
	totalPages := (totalCount + pageSize - 1) / pageSize

	return &orderingpb.GetOrdersByStatusResponse{
		Orders:      s.ordersFromCatalog(orders),
		TotalCount:  totalCount,
		TotalPages:  totalPages,
		CurrentPage: request.GetPage(),
	}, nil
}

func (s server) UpdateOrderStatus(ctx context.Context, request *orderingpb.UpdateOrderStatusRequest) (*orderingpb.UpdateOrderStatusResponse, error) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("OrderID", request.GetId()),
		attribute.String("Status", request.GetStatus()),
	)

	err := s.app.UpdateOrderStatus(ctx, commands.UpdateOrderStatus{
		ID:     request.GetId(),
		Status: request.GetStatus(),
		Reason: request.GetReason(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &orderingpb.UpdateOrderStatusResponse{
		Id:     request.GetId(),
		Status: request.GetStatus(),
	}, nil
}

func (s server) orderFromDomain(order *domain.Order) *orderingpb.Order {
	items := make([]*orderingpb.Item, len(order.Items))
	for i, item := range order.Items {
		items[i] = s.itemFromDomain(item)
	}

	return &orderingpb.Order{
		Id:              order.ID(),
		UserCustomerId:  order.UserCustomerID,
		PaymentMethodId: order.PaymentMethodID,
		Items:           items,
		Status:          order.Status.String(),
	}
}

func (s server) itemToDomain(item *orderingpb.Item) domain.Item {
	return domain.Item{
		ProductID:      item.GetProductId(),
		UserSellerID:   item.GetUserSellerId(),
		UserSellerName: item.GetUserSellerName(),
		ProductName:    item.GetProductName(),
		Price:          item.GetPrice(),
		Quantity:       item.GetQuantity(),
	}
}

func (s server) itemFromDomain(item domain.Item) *orderingpb.Item {
	return &orderingpb.Item{
		UserSellerId:   item.UserSellerID,
		ProductId:      item.ProductID,
		UserSellerName: item.UserSellerName,
		ProductName:    item.ProductName,
		Price:          item.Price,
		Quantity:       item.Quantity,
	}
}

func (s server) ordersFromCatalog(orders []*domain.OrderCatalog) []*orderingpb.Order {
	result := make([]*orderingpb.Order, len(orders))
	for i, order := range orders {
		result[i] = s.orderFromCatalog(order)
	}
	return result
}

func (s server) orderFromCatalog(order *domain.OrderCatalog) *orderingpb.Order {
	return &orderingpb.Order{
		Id:              order.ID,
		UserCustomerId:  order.UserCustomerID,
		PaymentMethodId: order.PaymentMethodID,
		Status:          order.Status.String(),
		// Note: Items are not stored in catalog, would need to fetch from main order if needed
		Items: []*orderingpb.Item{},
	}
}

func recordSpanError(span trace.Span, err error) {
	span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
	span.SetStatus(codes.Error, err.Error())
}
