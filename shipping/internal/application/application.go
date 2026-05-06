package application

import (
	"context"
	"middleman/internal/ddd"
	"middleman/shipping/internal/application/commands"
	"middleman/shipping/internal/application/queries"
	"middleman/shipping/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}
	Commands interface {
		CreateShipping(ctx context.Context, cmd commands.CreateShipping) error
		CancelShipment(ctx context.Context, cmd commands.CancelShipment) error
		SchedulePickup(ctx context.Context, cmd commands.SchedulePickup) error
		StartShipment(ctx context.Context, cmd commands.StartShipment) error
		UpdateShipmentStatus(ctx context.Context, cmd commands.UpdateShipmentStatus) error
		AssignCarrier(ctx context.Context, cmd commands.AssignCarrier) error
		MarkShipmentAsDelivered(ctx context.Context, cmd commands.MarkShipmentAsDelivered) error
		ReturnShipment(ctx context.Context, cmd commands.ReturnShipment) error
	}
	Queries interface {
		GetShipment(ctx context.Context, query queries.GetShipment) (*domain.CatalogShipment, error)
		ListShipments(ctx context.Context, query queries.ListShipments) ([]*domain.CatalogShipment, error)
		TrackShipment(ctx context.Context, query queries.TrackShipment) (*domain.CatalogShipment, error)
		GetShipmentHistory(ctx context.Context, query queries.GetShipmentHistory) ([]*domain.ShipmentEvent, error)
	}

	Application struct {
		appCommands
		appQueries
	}
	appCommands struct {
		commands.CreateShippingHandler
		commands.CancelShipmentHandler
		commands.SchedulePickupHandler
		commands.StartShipmentHandler
		commands.UpdateShipmentStatusHandler
		commands.AssignCarrierHandler
		commands.MarkShipmentAsDeliveredHandler
		commands.ReturnShipmentHandler
	}
	appQueries struct {
		queries.GetShipmentHandler
		queries.ListShipmentsHandler
		queries.TrackShipmentHandler
		queries.GetShipmentHistoryHandler
	}
)

var _ App = (*Application)(nil)

func New(
	shippingRepo domain.ShippingRepository,
	catalogRepo domain.ShippingCatalogRepository,
	publisher ddd.EventPublisher[ddd.Event],
) *Application {
	return &Application{
		appCommands: appCommands{
			CreateShippingHandler:           commands.NewCreateShippingHandler(shippingRepo, publisher),
			CancelShipmentHandler:           commands.NewCancelShipmentHandler(shippingRepo, publisher),
			SchedulePickupHandler:           commands.NewSchedulePickupHandler(shippingRepo, publisher),
			StartShipmentHandler:            commands.NewStartShipmentHandler(shippingRepo, publisher),
			UpdateShipmentStatusHandler:     commands.NewUpdateShipmentStatusHandler(shippingRepo, publisher),
			AssignCarrierHandler:            commands.NewAssignCarrierHandler(shippingRepo, publisher),
			MarkShipmentAsDeliveredHandler:  commands.NewMarkShipmentAsDeliveredHandler(shippingRepo, publisher),
			ReturnShipmentHandler:           commands.NewReturnShipmentHandler(shippingRepo, publisher),
		},
		appQueries: appQueries{
			GetShipmentHandler:        queries.NewGetShipmentHandler(catalogRepo),
			ListShipmentsHandler:      queries.NewListShipmentsHandler(catalogRepo),
			TrackShipmentHandler:      queries.NewTrackShipmentHandler(catalogRepo),
			GetShipmentHistoryHandler: queries.NewGetShipmentHistoryHandler(shippingRepo),
		},
	}
}
