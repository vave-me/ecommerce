package application

import (
	"context"
	"middleman/internal/ddd"
	"middleman/notifications/internal/application/commands"
	"middleman/notifications/internal/application/queries"
	"middleman/notifications/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}

	Commands interface {
		AddBasketAlert(ctx context.Context, cmd commands.AddBasketAlert) error
		AddWishlistAlert(ctx context.Context, cmd commands.AddWishlistAlert) error
		AddMessageAlert(ctx context.Context, cmd commands.AddMessageAlert) error
		AddCommentAlert(ctx context.Context, cmd commands.AddCommentAlert) error
		AddProductAlert(ctx context.Context, cmd commands.AddProductAlert) error
		AddInteractionAlert(ctx context.Context, cmd commands.AddInteractionAlert) error
		AddOfferAlert(ctx context.Context, cmd commands.AddOfferAlert) error
		AddSupportAlert(ctx context.Context, cmd commands.AddSupportAlert) error
		AddOrderAlert(ctx context.Context, cmd commands.AddOrderAlert) error
		AddReviewAlert(ctx context.Context, cmd commands.AddReviewAlert) error
		AddPaymentAlert(ctx context.Context, cmd commands.AddPaymentAlert) error
		AddFollowingAlert(ctx context.Context, cmd commands.AddFollowingAlert) error
		ReadAlert(ctx context.Context, cmd commands.ReadAlert) error
		DeleteAlert(ctx context.Context, cmd commands.DeleteAlert) error
		UpdatePreferences(ctx context.Context, cmd commands.UpdatePreferences) error
	}

	Queries interface {
		ListAlerts(ctx context.Context, query queries.ListAlerts) ([]*domain.MiddlemanAlert, error)
		GetAlertsByType(ctx context.Context, query queries.GetAlertsByType) ([]*domain.MiddlemanAlert, error)
		GetPreferences(ctx context.Context, query queries.GetPreferences) (*domain.UserPreferences, error)
	}

	Application struct {
		appCommands
		appQueries
	}

	appCommands struct {
		commands.AddBasketAlertHandler
		commands.AddWishlistAlertHandler
		commands.AddMessageAlertHandler
		commands.AddCommentAlertHandler
		commands.AddProductAlertHandler
		commands.AddInteractionAlertHandler
		commands.AddOfferAlertHandler
		commands.AddSupportAlertHandler
		commands.AddOrderAlertHandler
		commands.AddReviewAlertHandler
		commands.AddPaymentAlertHandler
		commands.AddFollowingAlertHandler
		commands.ReadAlertHandler
		commands.DeleteAlertHandler
		commands.UpdatePreferencesHandler
	}

	appQueries struct {
		queries.GetAlertsByTypeHandler
		queries.ListAlertsHandler
		queries.GetPreferencesHandler
	}
)

var _ App = (*Application)(nil)

func New(alerts domain.AlertRepository, catalogAlerts domain.CatalogRepository, prefsRepo domain.PreferencesRepository, publisher ddd.EventPublisher[ddd.Event]) *Application {
	return &Application{
		appCommands: appCommands{

			AddBasketAlertHandler:      commands.NewAddBasketAlertHandler(alerts, publisher),
			AddCommentAlertHandler:     commands.NewAddCommentAlertHandler(alerts, publisher),
			AddWishlistAlertHandler:    commands.NewAddWishlistAlertHandler(alerts, publisher),
			AddMessageAlertHandler:     commands.NewAddMessageAlertHandler(alerts, publisher),
			AddProductAlertHandler:     commands.NewAddProductAlertHandler(alerts, publisher),
			AddInteractionAlertHandler: commands.NewAddInteractionAlertHandler(alerts, publisher),
			AddOfferAlertHandler:       commands.NewAddOfferAlertHandler(alerts, publisher),
			AddSupportAlertHandler:     commands.NewAddSupportAlertHandler(alerts, publisher),
			AddOrderAlertHandler:       commands.NewAddOrderAlertHandler(alerts, publisher),
			AddReviewAlertHandler:      commands.NewAddReviewAlertHandler(alerts, publisher),
			AddPaymentAlertHandler:     commands.NewAddPaymentAlertHandler(alerts, publisher),
			AddFollowingAlertHandler:   commands.NewAddFollowingAlertHandler(alerts, publisher),
			ReadAlertHandler:           commands.NewReadAlertHandler(alerts, publisher),
			DeleteAlertHandler:         commands.NewDeleteAlertHandler(alerts, catalogAlerts, publisher),
			UpdatePreferencesHandler:   commands.NewUpdatePreferencesHandler(prefsRepo),
		},
		//	TODO implement all kind of queries which will use different repositories, it will later be decided by user which
		//  notifications he wants to receive,
		appQueries: appQueries{
			ListAlertsHandler:      queries.NewListAlertsHandler(catalogAlerts),
			GetAlertsByTypeHandler: queries.NewGetAlertsByTypeHandler(catalogAlerts),
			GetPreferencesHandler:  queries.NewGetPreferencesHandler(prefsRepo),
		},
	}
}
