package application

import (
	"context"
	"middleman/internal/ddd"
	"middleman/newsletters/internal/application/commands"
	"middleman/newsletters/internal/application/queries"
	"middleman/newsletters/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}

	Commands interface {
		// Newsletter Management
		CreateNewsletter(ctx context.Context, cmd commands.CreateNewsletter) (string, error)
		UpdateNewsletter(ctx context.Context, cmd commands.UpdateNewsletter) error
		DeleteNewsletter(ctx context.Context, cmd commands.DeleteNewsletter) error

		// Subscription Management
		Subscribe(ctx context.Context, cmd commands.Subscribe) (string, error)
		Unsubscribe(ctx context.Context, cmd commands.Unsubscribe) error
		UpdateSubscription(ctx context.Context, cmd commands.UpdateSubscription) error

		// Edition Management
		CreateEdition(ctx context.Context, cmd commands.CreateEdition) (string, error)
		UpdateEdition(ctx context.Context, cmd commands.UpdateEdition) error
		ScheduleEdition(ctx context.Context, cmd commands.ScheduleEdition) error
		SendEdition(ctx context.Context, cmd commands.SendEdition) (int, error)

		// Template Management
		CreateTemplate(ctx context.Context, cmd commands.CreateTemplate) (string, error)
		UpdateTemplate(ctx context.Context, cmd commands.UpdateTemplate) error
		DeleteTemplate(ctx context.Context, cmd commands.DeleteTemplate) error
	}

	Queries interface {
		// Newsletter Queries
		GetNewsletter(ctx context.Context, query queries.GetNewsletter) (*domain.CatalogNewsletter, error)
		ListNewsletters(ctx context.Context, query queries.ListNewsletters) ([]*domain.CatalogNewsletter, int, error)

		// Subscription Queries
		GetSubscription(ctx context.Context, query queries.GetSubscription) (*domain.CatalogSubscription, error)
		ListSubscriptions(ctx context.Context, query queries.ListSubscriptions) ([]*domain.CatalogSubscription, int, error)

		// Edition Queries
		GetEdition(ctx context.Context, query queries.GetEdition) (*domain.CatalogEdition, error)
		ListEditions(ctx context.Context, query queries.ListEditions) ([]*domain.CatalogEdition, int, error)

		// Template Queries
		GetTemplate(ctx context.Context, query queries.GetTemplate) (*domain.CatalogTemplate, error)
		ListTemplates(ctx context.Context, query queries.ListTemplates) ([]*domain.CatalogTemplate, int, error)

		// Analytics
		GetNewsletterStats(ctx context.Context, query queries.GetNewsletterStats) (*queries.NewsletterStats, error)
		GetEditionStats(ctx context.Context, query queries.GetEditionStats) (*queries.EditionStats, error)
	}

	Application struct {
		appCommands
		appQueries
	}

	appCommands struct {
		commands.CreateNewsletterHandler
		commands.UpdateNewsletterHandler
		commands.DeleteNewsletterHandler

		commands.SubscribeHandler
		commands.UnsubscribeHandler
		commands.UpdateSubscriptionHandler

		commands.CreateEditionHandler
		commands.UpdateEditionHandler
		commands.ScheduleEditionHandler
		commands.SendEditionHandler

		commands.CreateTemplateHandler
		commands.UpdateTemplateHandler
		commands.DeleteTemplateHandler
	}

	appQueries struct {
		queries.GetNewsletterHandler
		queries.ListNewslettersHandler

		queries.GetSubscriptionHandler
		queries.ListSubscriptionsHandler

		queries.GetEditionHandler
		queries.ListEditionsHandler

		queries.GetTemplateHandler
		queries.ListTemplatesHandler

		queries.GetNewsletterStatsHandler
		queries.GetEditionStatsHandler
	}
)

var _ App = (*Application)(nil)

func New(
	newsletters domain.NewsletterRepository,
	subscriptions domain.SubscriptionRepository,
	editions domain.EditionRepository,
	templates domain.TemplateRepository,
	newsletterCatalog domain.NewsletterCatalogRepository,
	subCatalog domain.SubscriptionCatalogRepository,
	editionCatalog domain.EditionCatalogRepository,
	templateCatalog domain.TemplateCatalogRepository,
	publisher ddd.EventPublisher[ddd.Event],
) *Application {
	return &Application{
		appCommands: appCommands{
			CreateNewsletterHandler: commands.NewCreateNewsletterHandler(newsletters, newsletterCatalog, publisher),
			UpdateNewsletterHandler: commands.NewUpdateNewsletterHandler(newsletters, newsletterCatalog, publisher),
			DeleteNewsletterHandler: commands.NewDeleteNewsletterHandler(newsletters, newsletterCatalog, publisher),

			SubscribeHandler:       commands.NewSubscribeHandler(subscriptions, subCatalog, newsletterCatalog, publisher),
			UnsubscribeHandler:     commands.NewUnsubscribeHandler(subscriptions, subCatalog, newsletterCatalog, publisher),
			UpdateSubscriptionHandler: commands.NewUpdateSubscriptionHandler(subscriptions, subCatalog, publisher),

			CreateEditionHandler:  commands.NewCreateEditionHandler(editions, editionCatalog, newsletterCatalog, publisher),
			UpdateEditionHandler:  commands.NewUpdateEditionHandler(editions, editionCatalog, newsletterCatalog, publisher),
			ScheduleEditionHandler: commands.NewScheduleEditionHandler(editions, editionCatalog, newsletterCatalog, publisher),
			SendEditionHandler:    commands.NewSendEditionHandler(editions, editionCatalog, newsletterCatalog, subCatalog, publisher),

			CreateTemplateHandler: commands.NewCreateTemplateHandler(templates, templateCatalog, publisher),
			UpdateTemplateHandler: commands.NewUpdateTemplateHandler(templates, templateCatalog, publisher),
			DeleteTemplateHandler: commands.NewDeleteTemplateHandler(templates, templateCatalog, publisher),
		},
		appQueries: appQueries{
			GetNewsletterHandler:  queries.NewGetNewsletterHandler(newsletterCatalog),
			ListNewslettersHandler: queries.NewListNewslettersHandler(newsletterCatalog),

			GetSubscriptionHandler:  queries.NewGetSubscriptionHandler(subCatalog, newsletterCatalog),
			ListSubscriptionsHandler: queries.NewListSubscriptionsHandler(subCatalog, newsletterCatalog),

			GetEditionHandler:  queries.NewGetEditionHandler(editionCatalog),
			ListEditionsHandler: queries.NewListEditionsHandler(editionCatalog),

			GetTemplateHandler:  queries.NewGetTemplateHandler(templateCatalog),
			ListTemplatesHandler: queries.NewListTemplatesHandler(templateCatalog),

			GetNewsletterStatsHandler: queries.NewGetNewsletterStatsHandler(subCatalog, editionCatalog),
			GetEditionStatsHandler:    queries.NewGetEditionStatsHandler(editionCatalog),
		},
	}
}