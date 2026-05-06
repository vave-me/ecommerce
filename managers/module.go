package managers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"middleman/erp/erppb"
	"middleman/managers/internal/application"
	"middleman/managers/internal/application/processor"
	"middleman/managers/internal/application/services"
	"middleman/managers/internal/application/tools"
	"middleman/managers/internal/application/consciousness"
	"middleman/managers/internal/constants"
	"middleman/managers/internal/domain"
	"middleman/managers/internal/grpc"
	"middleman/managers/internal/handlers"
	"middleman/managers/managerspb"
	"middleman/newsletters/newsletterspb"

	"middleman/activity/activitypb"
	"middleman/baskets/basketspb"
	"middleman/comments/commentspb"
	"middleman/following/followingpb"
	"middleman/messages/messagespb"
	"middleman/offers/offerspb"
	"middleman/ordering/orderingpb"
	"middleman/payments/paymentspb"
	"middleman/posts/postspb"
	"middleman/products/productspb"
	"middleman/reviews/reviewspb"
	"middleman/services/servicespb"
	"middleman/support/supportpb"
	"middleman/wishlists/wishlistspb"

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
	"middleman/managers/internal/postgres"
	"middleman/managers/internal/rest"
	"middleman/users/userspb"

	"time"

	"github.com/rs/zerolog"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, mono system.ManagersService) (err error) {
	return Root(ctx, mono)
}

func Root(ctx context.Context, svc system.ManagersService) (err error) {

	container := di.New()
	// setup Driven adapters
	container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.New()
		if err := registrations(reg); err != nil {
			return nil, err
		}
		if err := managerspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := userspb.Registrations(reg); err != nil {
			return nil, err
		}
		// Register all protobuf packages for platform event handlers
		if err := activitypb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := erppb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := basketspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := commentspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := followingpb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := messagespb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := newsletterspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := offerspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := orderingpb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := paymentspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := postspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := productspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := reviewspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := servicespb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := supportpb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := wishlistspb.Registrations(reg); err != nil {
			return nil, err
		}
		return reg, nil
	})

	anthropicClient, anthropicErr := ai2.NewAnthropicClient(svc.ManagersConfig().AnthropicAPIKey, svc.ManagersConfig().AnthropicBaseURL, svc.ManagersConfig().AnthropicModel)
	container.AddSingleton(constants.AnthropicClient, func(c di.Container) (any, error) {
		return anthropicClient, anthropicErr
	})

	openAiClient, openAiErr := ai2.NewOpenAIClient(svc.ManagersConfig().OpenAIAPIKey, svc.ManagersConfig().OpenAIBaseURL, svc.ManagersConfig().OpenAIBaseModel)
	container.AddSingleton(constants.OpenAIClient, func(c di.Container) (any, error) {
		return openAiClient, openAiErr
	})

	deepSeekClient, deepSeekErr := ai2.NewDeepSeekClient(svc.ManagersConfig().DeepSeekAPIKey, svc.ManagersConfig().DeepSeekBaseURL, svc.ManagersConfig().DeepSeekModel)

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
		clientProvider := c.Get(constants.AIClientProvider).(services.AIClientProvider)
		modelSelector := application.NewAIModelSelector(clientProvider)
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

	// Add Manager Repository
	container.AddScoped(constants.ManagersRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Manager](
			domain.ManagerAggregate,
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

	// Add LLM Journal Repository
	container.AddScoped(constants.LLMJournalRepositoryKey, func(c di.Container) (any, error) {
		return postgres.NewLLMJournalRepository(
			constants.LLMJournalTableName,
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

	// Add Scheduler Repository following the standard pattern
	container.AddScoped(constants.SchedulerRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewSchedulerRepository(svc.Config().Rpc.Service("scheduler"), svc.Auth()), nil
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

	container.AddScoped(constants.LLMJournalRepositoryKey, func(c di.Container) (any, error) {
		return postgres.NewLLMJournalRepository(
			constants.LLMJournalTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	// ToolServiceRegistry - The modern tool system using application services
	container.AddScoped(constants.ProductionToolRegistry, func(c di.Container) (any, error) {
		// Use the modern ToolServiceRegistry from application layer
		toolRegistry := tools.NewToolServiceRegistry(
			c.Get(constants.ProductRepositoryKey).(domain.ProductRepository),           // productRepo
			c.Get(constants.OrderRepositoryKey).(domain.OrderRepository),               // orderRepo
			c.Get(constants.UserRepositoryKey).(domain.UserRepository),                 // userRepo
			c.Get(constants.PaymentRepositoryKey).(domain.PaymentRepository),           // paymentRepo
			c.Get(constants.CommentRepositoryKey).(domain.CommentRepository),           // commentRepo
			c.Get(constants.ShippingRepositoryKey).(domain.ShippingRepository),         // shippingRepo
			c.Get(constants.NotificationRepositoryKey).(domain.NotificationRepository), // notificationRepo
			c.Get(constants.SupportRepositoryKey).(domain.SupportRepository),           // supportRepo
			c.Get(constants.CategoryRepositoryKey).(domain.CategoryRepository),         // categoryRepo
			c.Get(constants.MetricRepositoryKey).(domain.MetricRepository),             // metricRepo
			c.Get(constants.ReviewRepositoryKey).(domain.ReviewRepository),             // reviewRepo
			c.Get(constants.NewsletterRepositoryKey).(domain.NewsletterRepository),     // newsletterRepo
			c.Get(constants.MailerRepositoryKey).(domain.MailerRepository),             // mailerRepo
			c.Get(constants.WishlistRepositoryKey).(domain.WishlistRepository),         // wishlistRepo
			c.Get(constants.FollowingRepositoryKey).(domain.FollowingRepository),       // followingRepo
			c.Get(constants.PostRepositoryKey).(domain.PostRepository),                 // postRepo
			c.Get(constants.OfferRepositoryKey).(domain.OfferRepository),               // offerRepo
			c.Get(constants.ActivityRepositoryKey).(domain.ActivityRepository),         // activityRepo
			c.Get(constants.MediaRepositoryKey).(domain.MiddlemanMediaRepository),      // mediaRepo
			c.Get(constants.BasketRepositoryKey).(domain.BasketRepository),             // basketRepo
			c.Get(constants.GeocodingRepositoryKey).(domain.GeocodingRepository),       // geocodingRepo
			c.Get(constants.MessagesRepositoryKey).(domain.MessagesRepository),         // messageRepo
			c.Get(constants.VariantRepositoryKey).(domain.VariantRepository),           // variantRepo
			c.Get(constants.VectorRepositoryKey).(domain.VectorRepository),             // vectorRepo
			c.Get(constants.ServiceRepositoryKey).(domain.ServiceRepository),           // serviceRepo
			c.Get(constants.LLMJournalRepositoryKey).(domain.LLMJournalRepository),     // journalRepo
			&tools.ToolStreamConfig{
				BufferSize:       100,
				ProgressInterval: 500 * time.Millisecond,
				EnableMetrics:    true,
				MaxRetries:       3,
			},
		)

		log.Println("ToolServiceRegistry created successfully")
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
		defaultClient, err := clientProvider.GetDefaultClient(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to get default AI client: %w", err)
		}
		return services.NewAIRepositoryLanguageService(defaultClient), nil
	})

	// Add Schema Consultation Service
	container.AddSingleton(constants.SchemaConsultationService, func(c di.Container) (any, error) {
		return processor.NewSchemaConsultationService(), nil
	})

	// Add SystemPromptProvider
	container.AddSingleton(constants.SystemPromptProvider, func(c di.Container) (any, error) {
		return processor.NewEnhancedLLMInterface(), nil
	})


	// Note: Consciousness functionality is now integrated directly into Application
	// No separate factory needed

	// Add enhanced LLM processor with optimized services
	container.AddScoped(constants.LLMProcessor, func(c di.Container) (any, error) {
		clientProvider := c.Get(constants.AIClientProvider).(services.AIClientProvider)
		toolRegistry := c.Get(constants.ProductionToolRegistry).(*tools.ToolServiceRegistry)
		aiLanguageService := c.Get(constants.AIRepositoryLanguageService).(*services.AIRepositoryLanguageService)
		schemaService := c.Get(constants.SchemaConsultationService).(*processor.SchemaConsultationService)

		// Get optimized client with communication features
		defaultClient, err := clientProvider.GetDefaultClient(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to get optimized AI client: %w", err)
		}

		journalRepo := c.Get(constants.LLMJournalRepositoryKey).(domain.LLMJournalRepository)

		return processor.NewLLMProcessor(
			defaultClient,
			clientProvider,
			toolRegistry,
			aiLanguageService,
			schemaService,
			journalRepo,
		), nil
	})

	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		managers := c.Get(constants.ManagersRepoKey).(domain.ManagerRepository)
		conversations := c.Get(constants.ConversationsRepoKey).(domain.ConversationRepository)
		readConversations := c.Get(constants.ReadConversationRepositoryKey).(domain.ReadConversationRepository)
		readMessages := c.Get(constants.ReadMessagesRepositoryKey).(domain.ReadMessagesRepository)
		llmProcessor := c.Get(constants.LLMProcessor).(services.LLMProcessor)
		speechProcessor := c.Get(constants.SpeechProcessorKey).(services.SpeechProcessor)
		visionProcessor := c.Get(constants.VisionProcessorKey).(services.VisionProcessor)
		documentProcessor := c.Get(constants.DocumentProcessorKey).(services.DocumentProcessor)
		dataProcessor := c.Get(constants.DataProcessorKey).(services.DataProcessor)
		promptProvider := c.Get(constants.SystemPromptProvider).(domain.SystemPromptProvider)

		// Use the existing ProductionToolRegistry instead of creating duplicate
		toolRegistry := c.Get(constants.ProductionToolRegistry).(*tools.ToolServiceRegistry)

		// Create AI manager
		factory := ai2.NewClientFactory()
		
		// Register providers
		factory.RegisterProvider(ai2.ProviderOpenAI, ai2.ProviderConfig{
			APIKey:       svc.ManagersConfig().OpenAIAPIKey,
			DefaultModel: ai2.ModelGPT4oMini,
			Enabled:      true,
		})
		
		factory.RegisterProvider(ai2.ProviderAnthropic, ai2.ProviderConfig{
			APIKey:       svc.ManagersConfig().AnthropicAPIKey,
			DefaultModel: ai2.ModelClaude35SonnetLatest,
			Enabled:      true,
		})
		
		factory.RegisterProvider(ai2.ProviderDeepSeek, ai2.ProviderConfig{
			APIKey:       svc.ManagersConfig().DeepSeekAPIKey,
			DefaultModel: ai2.ModelDeepSeekV3,
			Enabled:      true,
		})
		
		aiManager := ai2.NewClientManager(factory)
		aiManager.SetDefaultProvider(ai2.ProviderDeepSeek) // Cost-effective for autonomous operations
		
		// Create application with AI manager
		config := &application.Config{
			LLMProcessor:         llmProcessor,
			SpeechProcessor:      speechProcessor,
			VisionProcessor:      visionProcessor,
			DocumentProcessor:    documentProcessor,
			DataProcessor:        dataProcessor,
			PromptProvider:       promptProvider,
			AIClientManager:      aiManager,
			ToolConfig: &application.ToolConfig{
				MaxConcurrentTools:   10,
				ToolExecutionTimeout: 10 * time.Minute,
				EnableMetrics:        true,
			},
			StreamingConfig: &application.StreamingConfig{
				MaxConcurrentTools:    10,
				ToolExecutionTimeout:  10 * time.Minute,
				StreamBufferSize:      100,
				EnableProgressUpdates: true,
				ChunkSize:             50,
			},
		}

		app, err := application.New(
			managers,
			conversations,
			readConversations,
			readMessages,
			c.Get(constants.CatalogRepositoryKey).(domain.CatalogRepository),
			c.Get(constants.VectorRepositoryKey).(domain.VectorRepository),
			toolRegistry,
			config,
			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.Event]),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create application: %w", err)
		}

		return app, nil
	})
	
	// Add Consciousness Manager
	container.AddScoped(constants.ConsciousnessManagerKey, func(c di.Container) (any, error) {
		app := c.Get(constants.ApplicationKey).(*application.Application)
		logger := svc.Logger()
		
		// Load consciousness configuration
		consciousnessConfig := config.LoadConsciousnessConfig()
		
		// Check if we should use production features
		if consciousnessConfig.Enabled && !consciousnessConfig.DryRunMode {
			// Create production consciousness manager with all features
			productionManager, err := consciousness.NewProductionConsciousnessManager(
				app,
				app.GetMemoryCore(),
				app.GetPatternDetector(),
				app.GetLearningProcessor(),
				app.GetDecisionOrchestrator(),
				app.GetActionExecutor(),
				app.GetAIManager(),
				consciousnessConfig,
				logger,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to create production consciousness manager: %w", err)
			}
			
			// Register health and readiness endpoints
			svc.Mux().HandleFunc("/health/consciousness", productionManager.GetHealthHandler())
			svc.Mux().HandleFunc("/ready/consciousness", productionManager.GetReadinessHandler())
			
			logger.Info().Msg("Production consciousness manager initialized with full features")
			return productionManager, nil
		}
		
		// Fallback to basic consciousness manager
		consciousnessManager := consciousness.NewConsciousnessManager(
			app,
			app.GetMemoryCore(),
			app.GetPatternDetector(),
			app.GetLearningProcessor(),
			app.GetDecisionOrchestrator(),
			app.GetActionExecutor(),
			app.GetAIManager(),
			logger,
		)
		
		logger.Info().Msg("Basic consciousness manager initialized")
		return consciousnessManager, nil
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
	// Register manager catalog handlers
	container.AddScoped(constants.ManagerCatalogHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewManagerCatalogHandlers(c.Get(constants.CatalogRepositoryKey).(domain.CatalogRepository)), nil
	})

	container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewDomainEventHandlers(c.Get(constants.EventPublisherKey).(am.EventPublisher)), nil
	})

	container.AddScoped(constants.IntegrationEventHandlersKey, func(c di.Container) (any, error) {
		// Check if consciousness is enabled
		if svc.ManagersConfig().ConsciousnessEnabled {
			// Get consciousness manager (can be either basic or production)
			var cm interface{ ProcessEvent(context.Context, ddd.Event) error }
			cmInterface := c.Get(constants.ConsciousnessManagerKey)
			
			// Check if it's a production manager
			if prodCM, ok := cmInterface.(*consciousness.ProductionConsciousnessManager); ok {
				cm = prodCM
			} else if basicCM, ok := cmInterface.(*consciousness.ConsciousnessManager); ok {
				cm = basicCM
			} else {
				return nil, fmt.Errorf("invalid consciousness manager type")
			}
			
			return handlers.NewIntegrationEventHandlersWithConsciousness(
				c.Get(constants.RegistryKey).(registry.Registry),
				c.Get(constants.ApplicationKey).(application.App),
				cm,
				tm.InboxHandler(c.Get(constants.InboxStoreKey).(tm.InboxStore)),
			), nil
		}
		
		// Fallback to standard handlers
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

	handlers.RegisterManagerCatalogHandlersTx(container)
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

	// Store Manager aggregate
	if err = serde.Register(domain.Manager{}, func(v any) error {
		manager := v.(*domain.Manager)
		manager.Aggregate = es.NewAggregate("", domain.ManagerAggregate)
		return nil
	}); err != nil {
		return
	}
	// store manager events
	if err = serde.Register(domain.ManagerCreated{}); err != nil {
		return
	}
	if err = serde.Register(domain.ManagerActivated{}); err != nil {
		return
	}
	if err = serde.Register(domain.ManagerDeactivated{}); err != nil {
		return
	}
	if err = serde.Register(domain.ManagerConfigurationUpdated{}); err != nil {
		return
	}
	if err = serde.Register(domain.ManagerRequestProcessed{}); err != nil {
		return
	}
	// store manager snapshots
	if err = serde.RegisterKey(domain.ManagerV1{}.SnapshotName(), domain.ManagerV1{}); err != nil {
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
