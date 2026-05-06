package application

import (
	"context"
	"middleman/internal/ddd"
	"middleman/ordering/internal/application/commands"
	"middleman/ordering/internal/application/queries"
	"middleman/ordering/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}

	Commands interface {
		CreateOrder(ctx context.Context, cmd commands.CreateOrder) error
		ApproveOrder(ctx context.Context, cmd commands.ApproveOrder) error
		RejectOrder(ctx context.Context, cmd commands.RejectOrder) error
		CancelOrder(ctx context.Context, cmd commands.CancelOrder) error
		ReadyOrder(ctx context.Context, cmd commands.ReadyOrder) error
		ShipOrder(ctx context.Context, cmd commands.ShipOrder) error
		DeliverOrder(ctx context.Context, cmd commands.DeliverOrder) error
		CompleteOrder(ctx context.Context, cmd commands.CompleteOrder) error
		UpdateOrderStatus(ctx context.Context, cmd commands.UpdateOrderStatus) error
	}

	Queries interface {
		GetOrder(ctx context.Context, q queries.GetOrder) (*domain.Order, error)
		ListOrders(ctx context.Context, q queries.ListOrders) ([]*domain.OrderCatalog, int64, error)
		GetOrdersByCustomer(ctx context.Context, q queries.GetOrdersByCustomer) ([]*domain.OrderCatalog, int64, error)
		GetOrdersByStatus(ctx context.Context, q queries.GetOrdersByStatus) ([]*domain.OrderCatalog, int64, error)
	}
	Application struct {
		appCommands
		appQueries
	}

	// appCommands holds each command handler
	appCommands struct {
		commands.CreateOrderHandler
		commands.ApproveOrderHandler
		commands.RejectOrderHandler
		commands.CancelOrderHandler
		commands.ReadyOrderHandler
		commands.ShipOrderHandler
		commands.DeliverOrderHandler
		commands.CompleteOrderHandler
		commands.UpdateOrderStatusHandler
	}

	// appQueries holds each query handler
	appQueries struct {
		queries.GetOrderHandler
		queries.ListOrdersHandler
		queries.GetOrdersByCustomerHandler
		queries.GetOrdersByStatusHandler
	}
)

var _ App = (*Application)(nil)

func New(orders domain.OrderRepository, catalog domain.CatalogRepository, publisher ddd.EventPublisher[ddd.Event]) *Application {
	return &Application{
		appCommands: appCommands{
			CreateOrderHandler:       commands.NewCreateOrderHandler(orders, publisher),
			ApproveOrderHandler:      commands.NewApproveOrderHandler(orders, publisher),
			RejectOrderHandler:       commands.NewRejectOrderHandler(orders, publisher),
			CancelOrderHandler:       commands.NewCancelOrderHandler(orders, publisher),
			ReadyOrderHandler:        commands.NewReadyOrderHandler(orders, publisher),
			ShipOrderHandler:         commands.NewShipOrderHandler(orders, publisher),
			DeliverOrderHandler:      commands.NewDeliverOrderHandler(orders, publisher),
			CompleteOrderHandler:     commands.NewCompleteOrderHandler(orders, publisher),
			UpdateOrderStatusHandler: commands.NewUpdateOrderStatusHandler(orders, publisher),
		},
		appQueries: appQueries{
			GetOrderHandler:            queries.NewGetOrderHandler(orders),
			ListOrdersHandler:          queries.NewListOrdersHandler(catalog),
			GetOrdersByCustomerHandler: queries.NewGetOrdersByCustomerHandler(catalog),
			GetOrdersByStatusHandler:   queries.NewGetOrdersByStatusHandler(catalog),
		},
	}
}
