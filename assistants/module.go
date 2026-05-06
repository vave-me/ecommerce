package assistants

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"middleman/assistants/assistantspb"
	"middleman/assistants/internal/application"
	"middleman/assistants/internal/application/processor"
	"middleman/assistants/internal/application/services"
	"middleman/assistants/internal/application/tools"
	"middleman/assistants/internal/constants"
	"middleman/assistants/internal/domain"
	"middleman/assistants/internal/grpc"
	"middleman/assistants/internal/handlers"

	"middleman/assistants/internal/postgres"
	"middleman/assistants/internal/rest"
	ai2 "middleman/internal/ai"
	"middleman/internal/am"
	"middleman/internal/amotel"
	"middleman/internal/amprom"
	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/internal/es"
	"middleman/internal/jetstream"
	pg "middleman/internal/postgres"
	"middleman/internal/postgresotel"
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
	"middleman/internal/system"
	"middleman/internal/tm"
	"middleman/users/userspb"
	"time"

	"github.com/rs/zerolog"
)

type Module struct{}

// simplePromptProvider implements domain.SystemPromptProvider
type simplePromptProvider struct{}

func (s *simplePromptProvider) GetCompleteSystemPrompt() string {
	return constants.SystemPrompt
}

func (m *Module) Startup(ctx context.Context, mono system.AssistantsService) (err error) {
	return Root(ctx, mono)
}

func Root(ctx context.Context, svc system.AssistantsService) (err error) {

	container := di.New()
	// setup Driven adapters
	container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.New()
		if err := registrations(reg); err != nil {
			return nil, err
		}
		if err := assistantspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := userspb.Registrations(reg); err != nil {
			return nil, err
		}
		return reg, nil
	})

	anthropicClient, anthropicErr := ai2.NewAnthropicClient(svc.AssistantsConfig().AnthropicAPIKey, svc.AssistantsConfig().AnthropicBaseURL, svc.AssistantsConfig().AnthropicModel)
	container.AddSingleton(constants.AnthropicClient, func(c di.Container) (any, error) {
		return anthropicClient, anthropicErr
	})

	openAiClient, openAiErr := ai2.NewOpenAIClient(svc.AssistantsConfig().OpenAIAPIKey, svc.AssistantsConfig().OpenAIBaseURL, svc.AssistantsConfig().OpenAIBaseModel)
	container.AddSingleton(constants.OpenAIClient, func(c di.Container) (any, error) {
		return openAiClient, openAiErr
	})

	deepSeekClient, deepSeekErr := ai2.NewDeepSeekClient(svc.AssistantsConfig().DeepSeekAPIKey, svc.AssistantsConfig().DeepSeekBaseURL, svc.AssistantsConfig().DeepSeekModel)

	container.AddSingleton(constants.DeepSeekAiClient, func(c di.Container) (any, error) {
		return deepSeekClient, deepSeekErr
	})

	container.AddSingleton(constants.SecurityValidator, func(c di.Container) (any, error) {
		validator := application.NewSecurityValidator()
		return validator, nil
	})

	// Add AIClientProvider
	container.AddSingleton(constants.AIClientProvider, func(c di.Container) (any, error) {
		// Create a map of AI clients for the provider
		clientMap := map[string]ai2.EnhancedAIService{
			"openai":    c.Get(constants.OpenAIClient).(ai2.EnhancedAIService),
			"anthropic": c.Get(constants.AnthropicClient).(ai2.EnhancedAIService),
			"deepseek":  c.Get(constants.DeepSeekAiClient).(ai2.EnhancedAIService),
		}

		// Define fallback order
		fallbackOrder := []string{"openai", "anthropic", "deepseek"}

		// Create circuit breaker config
		cbConfig := application.CircuitBreakerConfig{
			Name:             "ai-client-provider",
			MaxFailures:      5,
			ResetTimeout:     time.Minute * 2,
			SuccessThreshold: 3,
			Timeout:          time.Minute * 5,
		}

		clientProvider, err := application.NewAIClientProvider(clientMap, "openai", fallbackOrder, cbConfig)
		return clientProvider, err
	})

	// Add AI Model Selector
	container.AddSingleton(constants.AIModelSelector, func(c di.Container) (any, error) {
		modelSelector := application.NewAIModelSelector()
		return modelSelector, nil
	})

	stream := jetstream.NewStream(svc.Config().Nats.Stream, svc.JS(), svc.Logger())
	container.AddSingleton(constants.DomainDispatcherKey, func(c di.Container) (any, error) {
		return ddd.NewEventDispatcher[ddd.Event](), nil
	})

	container.AddScoped(constants.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return svc.DB().Begin()
	})

	sentCounter := amprom.SentMessagesCounter(constants.ServiceName)

	container.AddScoped(constants.MessagePublisherKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx))
		outboxStore := pg.NewOutboxStore(constants.OutboxTableName, tx)
		return am.NewMessagePublisher(
			stream,
			svc.Logger(),
			amotel.OtelMessageContextInjector(),
			sentCounter,
			tm.OutboxPublisher(outboxStore),
		), nil
	})

	container.AddSingleton(constants.MessageSubscriberKey, func(c di.Container) (any, error) {
		return am.NewMessageSubscriber(
			stream,
			svc.Logger(),
			amotel.OtelMessageContextExtractor(),
			amprom.ReceivedMessagesCounter(constants.ServiceName),
		), nil
	})

	container.AddScoped(constants.EventPublisherKey, func(c di.Container) (any, error) {
		return am.NewEventPublisher(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.MessagePublisherKey).(am.MessagePublisher),
			svc.Logger(),
		), nil
	})

	container.AddScoped(constants.InboxStoreKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx))
		return pg.NewInboxStore(constants.InboxTableName, tx), nil
	})

	container.AddScoped(constants.AggregateStoreKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx))
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		return es.AggregateStoreWithMiddleware(
			pg.NewEventStore(constants.EventsTableName, tx, reg),
			pg.NewSnapshotStore(constants.SnapshotsTableName, tx, reg),
		), nil
	})

	// Add Assistant Repository
	container.AddScoped(constants.AssistantsRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Assistant](
			domain.AssistantAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})

	container.AddScoped(constants.ConversationsRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Conversation](
			domain.ConversationAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})

	// === AUTHENTICATED REPOSITORIES (using rpc.DialWithAuth) ===
	container.AddScoped(constants.UserRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewUserRepository(svc.Config().Rpc.Service(constants.UserServiceName), svc.Auth()), nil
	})

	container.AddScoped(constants.PostRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewPostRepository(svc.Config().Rpc.Service(constants.PostsServiceName), svc.Auth()), nil
	})

	container.AddScoped(constants.ProductRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewProductRepository(svc.Config().Rpc.Service(constants.ProductsServiceName), svc.Auth()), nil
	})
	container.AddScoped(constants.WishlistRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewWishlistRepository(svc.Config().Rpc.Service(constants.WishlistsServiceName), svc.Auth()), nil
	})
	container.AddScoped(constants.FollowingRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewFollowingRepository(svc.Config().Rpc.Service(constants.FollowingServiceName), svc.Auth()), nil
	})
	container.AddScoped(constants.CommentRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewCommentRepository(svc.Config().Rpc.Service(constants.CommentsServiceName), svc.Auth()), nil
	})
	container.AddScoped(constants.SupportRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewSupportRepository(svc.Config().Rpc.Service(constants.SupportsServiceName), svc.Auth()), nil
	})
	container.AddScoped(constants.ShippingRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewShippingRepository(svc.Config().Rpc.Service(constants.ShippingServiceName), svc.Auth()), nil
	})
	container.AddScoped(constants.PaymentRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewPaymentRepository(svc.Config().Rpc.Service(constants.PaymentsServiceName), svc.Auth()), nil
	})
	container.AddScoped(constants.OrderRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewOrderRepository(svc.Config().Rpc.Service(constants.OrderingServiceName), svc.Auth()), nil
	})
	container.AddScoped(constants.OfferRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewOfferRepository(svc.Config().Rpc.Service(constants.OffersServiceName), svc.Auth()), nil
	})
	container.AddScoped(constants.NotificationRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewNotificationRepository(svc.Config().Rpc.Service(constants.NotificationsServiceName), svc.Auth()), nil
	})
	container.AddScoped(constants.NewsletterRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewNewsletterRepository(svc.Config().Rpc.Service(constants.NewsServiceName), svc.Auth()), nil
	})
	container.AddScoped(constants.MailerRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewMailerRepository(svc.Config().Rpc.Service(constants.MailerServiceName), svc.Auth()), nil
	})
	container.AddScoped(constants.GeocodingRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewGeocodingRepository(svc.Config().Rpc.Service(constants.GeocodingServiceName), svc.Auth()), nil
	})
	container.AddScoped(constants.CategoryRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewCategoryRepository(svc.Config().Rpc.Service(constants.CategoriesServiceName), svc.Auth()), nil
	})
	container.AddScoped(constants.MediaRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewMiddlemanMediaRepository(svc.Config().Rpc.Service(constants.MediaServiceName), svc.Auth()), nil
	})
	container.AddScoped(constants.ReviewRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewReviewRepository(svc.Config().Rpc.Service(constants.ReviewsServiceName), svc.Auth()), nil
	})
	container.AddScoped(constants.MessagesRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewMessagesRepository(svc.Config().Rpc.Service(constants.MessagesServiceName), svc.Auth()), nil
	})
	container.AddScoped(constants.ActivityRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewActivityRepository(svc.Config().Rpc.Service(constants.ActivityServiceName)), nil
	})
	container.AddScoped(constants.BasketRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewBasketRepository(svc.Config().Rpc.Service(constants.BasketsServiceName)), nil
	})
	container.AddScoped(constants.VariantRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewVariantRepository(svc.Config().Rpc.Service(constants.ProductsServiceName), svc.Auth()), nil
	})
	container.AddScoped(constants.MetricRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewMetricRepository(svc.Config().Rpc.Service(constants.MetricsServiceName)), nil
	})

	container.AddScoped(constants.CatalogRepositoryKey, func(c di.Container) (any, error) {
		return postgres.NewCatalogRepository(
			constants.CatalogTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	// Add Vector Repository
	container.AddScoped(constants.VectorRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewVectorRepository(
			svc.Config().Rpc.Service("vectors"),
			nil, // auth instance can be nil for now
		), nil
	})
	container.AddScoped(constants.ServiceRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewServiceRepository(svc.Config().Rpc.Service(constants.ServicesServiceName), svc.Auth()), nil
	})
	container.AddScoped(constants.ReadConversationRepositoryKey, func(c di.Container) (any, error) {
		return postgres.NewReadConversationRepository(
			constants.ReadConversationTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	container.AddScoped(constants.ReadMessagesRepositoryKey, func(c di.Container) (any, error) {
		return postgres.NewReadMessagesRepository(
			constants.ReadMessagesTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	// ToolRegistry - The simplified tool system with direct repository access
	container.AddScoped(constants.ProductionToolRegistry, func(c di.Container) (any, error) {
		// Use the new simplified ToolRegistry
		toolRegistry := tools.NewToolRegistry(
			c.Get(constants.ActivityRepositoryKey).(domain.ActivityRepository),
			c.Get(constants.BasketRepositoryKey).(domain.BasketRepository),
			c.Get(constants.CategoryRepositoryKey).(domain.CategoryRepository),
			c.Get(constants.CommentRepositoryKey).(domain.CommentRepository),
			c.Get(constants.GeocodingRepositoryKey).(domain.GeocodingRepository),
			c.Get(constants.MailerRepositoryKey).(domain.MailerRepository),
			c.Get(constants.MediaRepositoryKey).(domain.MiddlemanMediaRepository),
			c.Get(constants.MessagesRepositoryKey).(domain.MessagesRepository),
			c.Get(constants.MetricRepositoryKey).(domain.MetricRepository),
			c.Get(constants.NewsletterRepositoryKey).(domain.NewsletterRepository),
			c.Get(constants.NotificationRepositoryKey).(domain.NotificationRepository),
			c.Get(constants.OfferRepositoryKey).(domain.OfferRepository),
			c.Get(constants.OrderRepositoryKey).(domain.OrderRepository),
			c.Get(constants.PaymentRepositoryKey).(domain.PaymentRepository),
			c.Get(constants.PostRepositoryKey).(domain.PostRepository),
			c.Get(constants.ProductRepositoryKey).(domain.ProductRepository),
			c.Get(constants.ReviewRepositoryKey).(domain.ReviewRepository),
			c.Get(constants.ServiceRepositoryKey).(domain.ServiceRepository),
			c.Get(constants.ShippingRepositoryKey).(domain.ShippingRepository),
			c.Get(constants.SupportRepositoryKey).(domain.SupportRepository),
			c.Get(constants.UserRepositoryKey).(domain.UserRepository),
			c.Get(constants.VariantRepositoryKey).(domain.VariantRepository),
			c.Get(constants.VectorRepositoryKey).(domain.VectorRepository),
			c.Get(constants.WishlistRepositoryKey).(domain.WishlistRepository),
			c.Get(constants.FollowingRepositoryKey).(domain.FollowingRepository),
		)

		return toolRegistry, nil
	})

	// setup speech processor
	container.AddScoped(constants.SpeechProcessorKey, func(c di.Container) (any, error) {
		aiClient := c.Get(constants.AIClientProvider).(services.AIClientProvider)
		defaultClient, err := aiClient.GetDefaultClient(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to get default AI client for speech processor: %w", err)
		}
		return processor.NewOpenAISpeechProcessor(defaultClient, nil), nil
	})

	// setup vision processor
	container.AddScoped(constants.VisionProcessorKey, func(c di.Container) (any, error) {
		aiClient := c.Get(constants.AIClientProvider).(services.AIClientProvider)
		defaultClient, err := aiClient.GetDefaultClient(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to get default AI client for vision processor: %w", err)
		}
		return processor.NewOpenAIVisionProcessor(defaultClient, nil), nil
	})

	// setup document processor
	container.AddScoped(constants.DocumentProcessorKey, func(c di.Container) (any, error) {
		aiClient := c.Get(constants.AIClientProvider).(services.AIClientProvider)
		defaultClient, err := aiClient.GetDefaultClient(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to get default AI client for document processor: %w", err)
		}
		return processor.NewOpenAIDocumentProcessor(defaultClient, nil), nil
	})

	// setup data processor
	container.AddScoped(constants.DataProcessorKey, func(c di.Container) (any, error) {
		aiClient := c.Get(constants.AIClientProvider).(services.AIClientProvider)
		defaultClient, err := aiClient.GetDefaultClient(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to get default AI client for data processor: %w", err)
		}
		return processor.NewOpenAIDataProcessor(defaultClient, nil), nil
	})

	// Add Performance Optimizer
	container.AddSingleton(constants.PerformanceOptimizer, func(c di.Container) (any, error) {
		return application.NewPerformanceOptimizer(), nil
	})

	// Add AI Repository Language Service
	container.AddSingleton(constants.AIRepositoryLanguageService, func(c di.Container) (any, error) {
		clientProvider := c.Get(constants.AIClientProvider).(services.AIClientProvider)
		return services.NewAIRepositoryLanguageService(clientProvider), nil
	})

	// Unified Schema Service has been replaced by ProductionToolService
	// The system prompt is now provided by constants.SchemaAwareSystemPrompt

	// Add SystemPromptProvider - returns the production schema-aware prompt
	container.AddSingleton(constants.SystemPromptProvider, func(c di.Container) (any, error) {
		// Return a simple provider that implements domain.SystemPromptProvider
		return &simplePromptProvider{}, nil
	})

	// Add simplified LLM processor
	container.AddScoped(constants.LLMProcessor, func(c di.Container) (any, error) {
		clientProvider := c.Get(constants.AIClientProvider).(services.AIClientProvider)
		toolRegistry := c.Get(constants.ProductionToolRegistry).(*tools.ToolRegistry)

		// Get default AI client
		defaultClient, err := clientProvider.GetDefaultClient(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to get AI client: %w", err)
		}

		// Create simplified LLM processor
		llmProcessor := processor.NewLLMProcessor(
			defaultClient,
			clientProvider,
			toolRegistry,
		)

		return llmProcessor, nil
	})

	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		assistants := c.Get(constants.AssistantsRepoKey).(domain.AssistantRepository)
		conversations := c.Get(constants.ConversationsRepoKey).(domain.ConversationRepository)
		readConversations := c.Get(constants.ReadConversationRepositoryKey).(domain.ReadConversationRepository)
		readMessages := c.Get(constants.ReadMessagesRepositoryKey).(domain.ReadMessagesRepository)
		llmProcessor := c.Get(constants.LLMProcessor).(services.LLMProcessor)
		speechProcessor := c.Get(constants.SpeechProcessorKey).(services.SpeechProcessor)
		visionProcessor := c.Get(constants.VisionProcessorKey).(services.VisionProcessor)
		documentProcessor := c.Get(constants.DocumentProcessorKey).(services.DocumentProcessor)
		dataProcessor := c.Get(constants.DataProcessorKey).(services.DataProcessor)
		promptProvider := c.Get(constants.SystemPromptProvider).(domain.SystemPromptProvider)

		// Use the simplified ToolRegistry
		toolRegistry := c.Get(constants.ProductionToolRegistry).(*tools.ToolRegistry)

		// Create application config
		config := &application.Config{
			LLMProcessor:      llmProcessor,
			SpeechProcessor:   speechProcessor,
			VisionProcessor:   visionProcessor,
			DocumentProcessor: documentProcessor,
			DataProcessor:     dataProcessor,
			PromptProvider:    promptProvider,
			ToolConfig: &application.ToolConfig{
				MaxConcurrentTools:   20,
				ToolExecutionTimeout: 30 * time.Minute, // Increased timeout
				EnableMetrics:        true,
			},
			StreamingConfig: &application.StreamingConfig{
				MaxConcurrentTools:    20,
				ToolExecutionTimeout:  30 * time.Minute, // Increased timeout
				StreamBufferSize:      200,
				EnableProgressUpdates: true,
				ChunkSize:             100,
			},
		}

		return application.New(
			assistants,
			conversations,
			readConversations,
			readMessages,
			c.Get(constants.CatalogRepositoryKey).(domain.CatalogRepository),
			c.Get(constants.VectorRepositoryKey).(domain.VectorRepository),
			toolRegistry,
			config,
			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.Event]),
		), nil
	})

	// Unified Repository Interface - Use Application directly
	container.AddScoped(constants.UnifiedRepositoryInterface, func(c di.Container) (any, error) {
		app := c.Get(constants.ApplicationKey).(*application.Application)
		return app, nil
	})

	// Register read conversation handlers
	container.AddScoped(constants.ReadConversationHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewReadConversationHandlers(
			c.Get(constants.ReadConversationRepositoryKey).(domain.ReadConversationRepository),
		), nil
	})
	// Register read conversation handlers
	container.AddScoped(constants.ReadMessagesHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewReadMessagesHandlers(
			c.Get(constants.ReadMessagesRepositoryKey).(domain.ReadMessagesRepository),
		), nil
	})
	// Register assistant catalog handlers
	container.AddScoped(constants.AssistantCatalogHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewAssistantCatalogHandlers(c.Get(constants.CatalogRepositoryKey).(domain.CatalogRepository)), nil
	})

	container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewDomainEventHandlers(c.Get(constants.EventPublisherKey).(am.EventPublisher)), nil
	})

	container.AddScoped(constants.IntegrationEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewIntegrationEventHandlers(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.ApplicationKey).(application.App),
			tm.InboxHandler(c.Get(constants.InboxStoreKey).(tm.InboxStore)),
		), nil
	})

	outboxProcessor := tm.NewOutboxProcessor(
		stream,
		pg.NewOutboxStore(constants.OutboxTableName, svc.DB()),
	)

	if err = grpc.RegisterServerTx(container, svc.RPC()); err != nil {
		return err
	}
	if err = rest.RegisterGateway(ctx, svc.Mux(), svc.Config().Rpc.Address()); err != nil {
		return err
	}
	if err = rest.RegisterSwagger(svc.Mux()); err != nil {
		return err
	}

	handlers.RegisterAssistantCatalogHandlersTx(container)
	handlers.RegisterReadConversationHandlersTx(container)
	handlers.RegisterReadMessagesHandlersTx(container)
	if err = handlers.RegisterIntegrationEventHandlersTx(container); err != nil {
		return err
	}
	handlers.RegisterDomainEventHandlersTx(container)

	startOutboxProcessor(ctx, outboxProcessor, svc.Logger())

	// Eagerly initialize ToolServiceRegistry to trigger debug test
	log.Println("Eagerly initializing ToolServiceRegistry...")
	app := container.Get(constants.ApplicationKey).(*application.Application)
	log.Printf("Application with tool registry eagerly initialized: %v", app != nil)

	return
}

func registrations(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	// Store Assistant aggregate
	if err = serde.Register(domain.Assistant{}, func(v any) error {
		assistant := v.(*domain.Assistant)
		assistant.Aggregate = es.NewAggregate("", domain.AssistantAggregate)
		return nil
	}); err != nil {
		return
	}
	// store assistant events
	if err = serde.Register(domain.AssistantCreated{}); err != nil {
		return
	}
	if err = serde.Register(domain.AssistantActivated{}); err != nil {
		return
	}
	if err = serde.Register(domain.AssistantDeactivated{}); err != nil {
		return
	}
	if err = serde.Register(domain.AssistantConfigurationUpdated{}); err != nil {
		return
	}
	if err = serde.Register(domain.AssistantRequestProcessed{}); err != nil {
		return
	}
	// store assistant snapshots
	if err = serde.RegisterKey(domain.AssistantV1{}.SnapshotName(), domain.AssistantV1{}); err != nil {
		return
	}

	// Store Conversation aggregate
	if err = serde.Register(domain.Conversation{}, func(v any) error {
		conversation := v.(*domain.Conversation)
		conversation.Aggregate = es.NewAggregate("", domain.ConversationAggregate)
		return nil
	}); err != nil {
		return
	}
	// store conversation events
	if err = serde.Register(domain.ConversationCreated{}); err != nil {
		return
	}
	if err = serde.Register(domain.MessageAdded{}); err != nil {
		return
	}
	if err = serde.Register(domain.ConversationContextUpdated{}); err != nil {
		return
	}
	if err = serde.Register(domain.ConversationArchived{}); err != nil {
		return
	}
	// store conversation snapshots
	if err = serde.RegisterKey(domain.ConversationV1{}.SnapshotName(), domain.ConversationV1{}); err != nil {
		return
	}

	return
}

func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger zerolog.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("outbox processor error")
		}
	}()
}
