package vector

import (
	"context"
	"fmt"
	"log"
	"time"
	"middleman/internal/ai"
	"middleman/internal/am"
	"middleman/internal/amotel"
	"middleman/internal/amprom"
	"middleman/internal/di"
	"middleman/internal/jetstream"
	"middleman/internal/registry"
	"middleman/internal/system"
	"middleman/ordering/orderingpb"
	"middleman/posts/postspb"
	"middleman/products/productspb"
	"middleman/services/servicespb"
	"middleman/users/userspb"
	"middleman/vectors/internal/application"
	"middleman/vectors/internal/constants"
	"middleman/vectors/internal/grpc"
	"middleman/vectors/internal/handlers"
	"middleman/vectors/internal/infra"
	"middleman/vectors/internal/ports"
	"middleman/vectors/internal/rest"
	"middleman/vectors/internal/vector"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono system.VectorsService) (err error) {
	return Root(ctx, mono)
}

func Root(ctx context.Context, svc system.VectorsService) (err error) {
	container := di.New()

	// setup Driven adapters
	container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.New()
		if err := orderingpb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := userspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := productspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := postspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := servicespb.Registrations(reg); err != nil {
			return nil, err
		}
		return reg, nil
	})

	stream := jetstream.NewStream(svc.Config().Nats.Stream, svc.JS(), svc.Logger())

	container.AddScoped(constants.RedisTransactionKey, func(c di.Container) (any, error) { return *svc.Redis(), nil })

	// Embedding Client Provider with auto-discovery
	container.AddSingleton("embeddingProvider", func(c di.Container) (any, error) {
		// Create embedding provider
		embeddingProvider := infra.NewEmbeddingClientProvider()
		
		// Try to register OpenAI client with config
		openAIKey := svc.VectorsConfig().OpenAIAPIKey
		if openAIKey != "" {
			openAIClient, err := ai.NewOpenAIClient(openAIKey, svc.VectorsConfig().OpenAIBaseURL, "text-embedding-3-small")
			if err != nil {
				log.Printf("ERROR: Failed to create OpenAI client: %v", err)
			} else {
				config := infra.EmbeddingClientConfig{
					Model:         "text-embedding-3-small",
					Dimensions:    1536,
					PromptEnabled: true,
					MaxRetries:    3,
					Timeout:       30 * time.Second,
				}
				embeddingClient := infra.NewOpenAIEmbeddingClient(openAIClient, config)
				if err := embeddingProvider.RegisterClient("openai", embeddingClient); err != nil {
					log.Printf("ERROR: Failed to register OpenAI embedding client: %v", err)
				} else {
					log.Println("Registered OpenAI embedding client with model: text-embedding-3-small")
				}
			}
		} else {
			log.Println("WARNING: No OpenAI API key configured in VectorsConfig")
		}
		
		// Register simple client as fallback
		simpleClient := infra.NewSimpleEmbeddingClient("text-embedding-3-small", 1536)
		if err := embeddingProvider.RegisterClient("simple", simpleClient); err != nil {
			log.Printf("ERROR: Failed to register simple embedding client: %v", err)
		} else {
			log.Println("Registered simple embedding client as fallback")
		}
		
		return embeddingProvider, nil
	})

	// Message Subscriber
	container.AddSingleton(constants.MessageSubscriberKey, func(c di.Container) (any, error) {
		return am.NewMessageSubscriber(
			stream, svc.Logger(),
			amotel.OtelMessageContextExtractor(),
			amprom.ReceivedMessagesCounter(constants.ServiceName),
		), nil
	})

	// ===============================
	// QDRANT + REDIS INFRASTRUCTURE (Real implementations)
	// ===============================

	// Vector Service (Qdrant-based)
	container.AddSingleton("vectorService", func(c di.Container) (any, error) {
		// Use environment variables with Docker-friendly defaults
		qdrantHost := "qdrant" // Docker service name
		qdrantPort := "6334"   // Default gRPC port

		// Allow override via config if available
		vectorsConfig := svc.VectorsConfig()
		if vectorsConfig.QdrantHost != "" {
			qdrantHost = vectorsConfig.QdrantHost
		}
		if vectorsConfig.QdrantPort != "" {
			qdrantPort = vectorsConfig.QdrantPort
		}

		config := vector.Config{
			QdrantHost:     qdrantHost,
			QdrantPort:     qdrantPort,
			CollectionName: "vectors",
			VectorSize:     1536, // OpenAI text-embedding-3-small dimensions
		}
		return vector.NewVectorService(config)
	})

	// Embedding Service (Simple implementation for Qdrant + Redis)
	container.AddSingleton("embeddingService", func(c di.Container) (any, error) {
		config := vector.EmbeddingConfig{
			Model:         "text-embedding-3-small",
			Dimensions:    1536,
			PromptEnabled: false, // Simplified for now
			LLMClient:     nil,   // Not needed for simple implementation
		}
		return vector.NewEmbeddingService(config), nil
	})

	// Vector Repository (Real implementation using Qdrant + Redis)
	container.AddScoped("vectorRepository", func(c di.Container) (any, error) {
		vectorService := c.Get("vectorService").(*vector.VectorService)
		embeddingService := c.Get("embeddingService").(*vector.EmbeddingService)
		embeddingProvider := c.Get("embeddingProvider").(*infra.EmbeddingClientProvider)

		// Create vector provider with Qdrant client
		vectorProvider := infra.NewVectorClientProvider()
		qdrantClient := infra.NewQdrantVectorClient(vectorService)
		vectorProvider.RegisterClient("qdrant", qdrantClient)

		return vector.NewVectorRepository(embeddingProvider, vectorProvider, vectorService, embeddingService), nil
	})

	// Repository Collection (Qdrant + Redis only)
	container.AddScoped("repositoryCollection", func(c di.Container) (any, error) {
		return &repositoryCollection{
			vectorRepo: c.Get("vectorRepository").(*vector.VectorRepository),
		}, nil
	})

	// Embedding Interface implementation
	container.AddScoped("embeddingInterface", func(c di.Container) (any, error) {
		embeddingProvider := c.Get("embeddingProvider").(*infra.EmbeddingClientProvider)
		return &embeddingWrapper{provider: embeddingProvider}, nil
	})

	// ===============================
	// APPLICATION LAYER
	// ===============================

	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		vectorRepo := c.Get("vectorRepository").(*vector.VectorRepository)
		repos := c.Get("repositoryCollection").(ports.RepositoryCollection)
		embedding := c.Get("embeddingInterface").(application.EmbeddingInterface)

		return application.New(vectorRepo, repos, embedding), nil
	})

	// Integration Event Handlers
	container.AddScoped(constants.IntegrationEventHandlersKey, func(c di.Container) (any, error) {
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		app := c.Get(constants.ApplicationKey).(application.Application)

		return handlers.NewVectorIntegrationEventHandlers(reg, app), nil
	})

	// ===============================
	// GRPC + REST SETUP
	// ===============================

	// Register GRPC server
	if err = grpc.RegisterServerTx(container, svc.RPC()); err != nil {
		return err
	}

	// Register REST gateway
	if err = rest.RegisterGateway(ctx, svc.Mux(), svc.Config().Rpc.Address()); err != nil {
		return err
	}

	// Register Swagger UI
	if err = rest.RegisterSwagger(svc.Mux()); err != nil {
		return err
	}

	// Register event handlers
	if err = handlers.RegisterIntegrationEventHandlersTx(container); err != nil {
		return err
	}

	// start service
	return nil
}

// ===============================
// REPOSITORY COLLECTION (Qdrant + Redis only)
// ===============================

type repositoryCollection struct {
	vectorRepo *vector.VectorRepository
}

func (r *repositoryCollection) Products() ports.ProductCacheRepositoryPort {
	return nil // Redis cache can be added later
}

func (r *repositoryCollection) Posts() ports.PostCacheRepositoryPort {
	return nil // Redis cache can be added later
}

func (r *repositoryCollection) Users() ports.UserCacheRepositoryPort {
	return nil // Redis cache can be added later
}

func (r *repositoryCollection) Variants() ports.VariantCacheRepositoryPort {
	return nil // Redis cache can be added later
}

func (r *repositoryCollection) Orders() ports.OrderRepositoryPort {
	return nil // Redis cache can be added later
}

func (r *repositoryCollection) ItemMetrics() ports.ItemMetricRepositoryPort {
	return nil // Redis cache can be added later
}
func (r *repositoryCollection) Services() ports.ServiceCacheRepositoryPort {
	return nil // Redis cache can be added later
}
func (r *repositoryCollection) Vectors() ports.VectorRepositoryPort {
	return r.vectorRepo
}

// ===============================
// EMBEDDING INTERFACE IMPLEMENTATION
// ===============================

type embeddingWrapper struct {
	provider *infra.EmbeddingClientProvider
}

func (e *embeddingWrapper) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	client, err := e.provider.GetDefaultClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get embedding client: %w", err)
	}
	return client.GenerateEmbedding(ctx, text)
}

func (e *embeddingWrapper) GenerateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	client, err := e.provider.GetDefaultClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get embedding client: %w", err)
	}
	return client.GenerateBatchEmbeddings(ctx, texts)
}

func (e *embeddingWrapper) GenerateEntityEmbedding(ctx context.Context, entityData map[string]interface{}) ([]float32, error) {
	client, err := e.provider.GetDefaultClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get embedding client: %w", err)
	}
	return client.GenerateEntityEmbedding(ctx, entityData)
}

func (e *embeddingWrapper) GenerateEntityEmbeddingWithPrompt(ctx context.Context, entityType string, entityData map[string]interface{}, strategy application.TransformationStrategy) ([]float32, error) {
	client, _, err := e.provider.GetOptimalProvider(ctx, entityType)
	if err != nil {
		return nil, fmt.Errorf("failed to get optimal embedding client: %w", err)
	}
	return client.GenerateOptimizedEmbedding(ctx, entityType, entityData, string(strategy))
}

func (e *embeddingWrapper) GenerateOptimizedEmbedding(ctx context.Context, entityType string, entityData map[string]interface{}, optimization string) ([]float32, error) {
	client, _, err := e.provider.GetOptimalProvider(ctx, entityType)
	if err != nil {
		return nil, fmt.Errorf("failed to get optimal embedding client: %w", err)
	}
	return client.GenerateOptimizedEmbedding(ctx, entityType, entityData, optimization)
}

func (e *embeddingWrapper) GetDimensions() int {
	client, err := e.provider.GetDefaultClient(context.Background())
	if err != nil {
		return 1536 // Default dimensions
	}
	return client.GetDimensions()
}

func (e *embeddingWrapper) GetModel() string {
	client, err := e.provider.GetDefaultClient(context.Background())
	if err != nil {
		return "unknown"
	}
	return client.GetModel()
}

func (e *embeddingWrapper) IsPromptEnabled() bool {
	client, err := e.provider.GetDefaultClient(context.Background())
	if err != nil {
		return false
	}
	return client.IsPromptEnabled()
}

// ===============================
// SIMPLE LLM CLIENT (Not needed for simple implementation)
// ===============================

type simpleLLMWrapper struct{}

func (l *simpleLLMWrapper) Transform(ctx context.Context, text string, prompt string) (string, error) {
	// Simple implementation just returns the original text
	return text, nil
}

func (l *simpleLLMWrapper) IsAvailable() bool {
	return false // Not available in simple implementation
}
