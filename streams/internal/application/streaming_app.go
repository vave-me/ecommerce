package application

import (
	"context"

	"middleman/internal/ddd"
	"middleman/streams/internal/application/commands"
	"middleman/streams/internal/application/queries"
	"middleman/streams/internal/domain"
	"middleman/streams/internal/handlers"
	"middleman/streams/internal/infrastructure"
	"middleman/streams/internal/infrastructure/streaming"
)

type StreamingApp interface {
	Commands
	Queries
}

type Commands interface {
	CreateLiveStream(ctx context.Context, cmd commands.CreateLiveStream) error
	ConfigureStreaming(ctx context.Context, cmd commands.ConfigureStreaming) error
	StartLiveStream(ctx context.Context, cmd commands.StartLiveStream) error
	
	// Webhook commands
	SubscribeWebhook(ctx context.Context, cmd commands.SubscribeWebhook) error
	UpdateWebhook(ctx context.Context, cmd commands.UpdateWebhook) error
	UnsubscribeWebhook(ctx context.Context, cmd commands.UnsubscribeWebhook) error
	TestWebhook(ctx context.Context, cmd commands.TestWebhook) error
	RetryWebhookDelivery(ctx context.Context, cmd commands.RetryWebhookDelivery) error
}

type Queries interface {
	GetLiveStream(ctx context.Context, query queries.GetLiveStream) (*domain.LiveStream, error)
	
	// Webhook queries
	ListWebhookSubscriptions(ctx context.Context, query queries.ListWebhookSubscriptions) ([]*domain.WebhookSubscription, error)
	GetWebhookSubscription(ctx context.Context, query queries.GetWebhookSubscription) (*domain.WebhookSubscription, error)
	GetWebhookDeliveries(ctx context.Context, query queries.GetWebhookDeliveries) ([]*domain.WebhookDelivery, error)
}

type Application struct {
	appCommands
	appQueries
}

type appCommands struct {
	commands.CreateLiveStreamHandler
	commands.ConfigureStreamingHandler
	commands.StartLiveStreamHandler
	
	// Webhook handlers
	commands.SubscribeWebhookHandler
	commands.UpdateWebhookHandler
	commands.UnsubscribeWebhookHandler
	commands.TestWebhookHandler
	commands.RetryWebhookDeliveryHandler
}

type appQueries struct {
	queries.GetLiveStreamHandler
	
	// Webhook query handlers
	queries.ListWebhookSubscriptionsHandler
	queries.GetWebhookSubscriptionHandler
	queries.GetWebhookDeliveriesHandler
}

var _ StreamingApp = (*Application)(nil)

func New(
	liveStreams domain.LiveStreamRepository,
	webhookSubscriptions domain.WebhookSubscriptionRepository,
	webhookDeliveries domain.WebhookDeliveryRepository,
	publisher ddd.EventPublisher,
	streamingServer *streaming.StreamingServer,
	cdnManager *streaming.CDNManager,
	drmManager *streaming.DRMManager,
	webrtcServer *streaming.WebRTCServer,
	webhookClient *infrastructure.WebhookClient,
	webhookDispatcher *handlers.WebhookDispatcher,
) *Application {
	return &Application{
		appCommands: appCommands{
			CreateLiveStreamHandler: commands.NewCreateLiveStreamHandler(liveStreams, publisher),
			ConfigureStreamingHandler: commands.NewConfigureStreamingHandler(
				liveStreams, publisher, streamingServer, cdnManager, drmManager,
			),
			StartLiveStreamHandler: commands.NewStartLiveStreamHandler(
				liveStreams, publisher, streamingServer,
			),
			// Webhook handlers
			SubscribeWebhookHandler: commands.NewSubscribeWebhookHandler(webhookSubscriptions),
			UpdateWebhookHandler: commands.NewUpdateWebhookHandler(webhookSubscriptions),
			UnsubscribeWebhookHandler: commands.NewUnsubscribeWebhookHandler(webhookSubscriptions),
			TestWebhookHandler: commands.NewTestWebhookHandler(webhookSubscriptions, webhookClient),
			RetryWebhookDeliveryHandler: commands.NewRetryWebhookDeliveryHandler(webhookDeliveries, webhookDispatcher),
		},
		appQueries: appQueries{
			GetLiveStreamHandler: queries.NewGetLiveStreamHandler(liveStreams),
			// Webhook query handlers
			ListWebhookSubscriptionsHandler: queries.NewListWebhookSubscriptionsHandler(webhookSubscriptions),
			GetWebhookSubscriptionHandler: queries.NewGetWebhookSubscriptionHandler(webhookSubscriptions),
			GetWebhookDeliveriesHandler: queries.NewGetWebhookDeliveriesHandler(webhookDeliveries),
		},
	}
}