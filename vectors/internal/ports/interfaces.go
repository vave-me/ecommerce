package ports

import (
	"context"
	"middleman/vectors/internal/models"
	"time"
)

// ===============================
// CORE CLIENT PROVIDER INTERFACES
// ===============================

// VectorClientProvider provides access to vector database clients
type VectorClientProvider interface {
	GetClient(ctx context.Context, provider string) (VectorDatabaseClient, error)
	GetDefaultClient(ctx context.Context) (VectorDatabaseClient, error)
	GetHealthyProvider(ctx context.Context) (VectorDatabaseClient, string, error)
	GetOptimalProvider(ctx context.Context, operationType string) (VectorDatabaseClient, string, error)
}

// EmbeddingClientProvider provides access to embedding generation clients
type EmbeddingClientProvider interface {
	GetClient(ctx context.Context, provider string) (EmbeddingClient, error)
	GetDefaultClient(ctx context.Context) (EmbeddingClient, error)
	GetOptimalProvider(ctx context.Context, entityType string) (EmbeddingClient, string, error)
}

// ExternalServiceClientProvider provides access to external microservice clients
type ExternalServiceClientProvider interface {
	GetServiceClient(ctx context.Context, serviceName string) (ExternalServiceClient, error)
	GetHealthyServiceClient(ctx context.Context, serviceName string) (ExternalServiceClient, error)
	RegisterServiceClient(serviceName string, client ExternalServiceClient) error
}

// ===============================
// VECTOR DATABASE CLIENT INTERFACE
// ===============================

// VectorDatabaseClient defines the interface for vector database operations
type VectorDatabaseClient interface {
	// Core vector operations
	IndexVector(ctx context.Context, point VectorPoint) error
	BatchIndexVectors(ctx context.Context, points []VectorPoint) error
	SearchSimilar(ctx context.Context, req SearchRequest) (*SearchResponse, error)
	GetVector(ctx context.Context, id string, withVector bool) (*VectorPoint, error)
	DeleteVector(ctx context.Context, id string) error
	BatchDeleteVectors(ctx context.Context, ids []string) error

	// Health and diagnostics
	HealthCheck(ctx context.Context) error
	GetStats(ctx context.Context) (*VectorDBStats, error)
	GetProviderInfo() VectorDBProviderInfo

	// Connection management
	Close() error
	IsConnected() bool
}

// ===============================
// EMBEDDING CLIENT INTERFACE
// ===============================

// EmbeddingClient defines the interface for embedding generation
type EmbeddingClient interface {
	// Core embedding operations
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
	GenerateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error)
	GenerateEntityEmbedding(ctx context.Context, entityData map[string]interface{}) ([]float32, error)

	// Advanced embedding operations
	GenerateEmbeddingWithPrompt(ctx context.Context, text string, prompt string) ([]float32, error)
	GenerateOptimizedEmbedding(ctx context.Context, entityType string, entityData map[string]interface{}, optimization string) ([]float32, error)

	// Client information
	GetDimensions() int
	GetModel() string
	IsPromptEnabled() bool
	HealthCheck(ctx context.Context) error
	GetProviderInfo() EmbeddingProviderInfo
}

// ===============================
// EXTERNAL SERVICE CLIENT INTERFACE
// ===============================

// ExternalServiceClient defines the interface for communicating with external microservices
type ExternalServiceClient interface {
	// Entity operations
	GetEntity(ctx context.Context, entityType string, entityID string) (interface{}, error)
	SearchEntities(ctx context.Context, entityType string, filters map[string]interface{}) ([]interface{}, error)
	BatchGetEntities(ctx context.Context, entityType string, entityIDs []string) ([]interface{}, error)

	// Health and diagnostics
	HealthCheck(ctx context.Context) error
	GetServiceInfo() ExternalServiceInfo
	IsConnected() bool
}

// ===============================
// LLM CLIENT INTERFACE
// ===============================

// LLMClient defines the interface for LLM interactions (for text transformation)
type LLMClient interface {
	Transform(ctx context.Context, text string, prompt string) (string, error)
	BatchTransform(ctx context.Context, texts []string, prompt string) ([]string, error)
	IsAvailable() bool
	HealthCheck(ctx context.Context) error
	GetProviderInfo() LLMProviderInfo
}

// ===============================
// REPOSITORY PORT INTERFACES
// ===============================

// VectorRepositoryPort abstracts vector operations for domain entities
type VectorRepositoryPort interface {
	// Entity indexing operations
	IndexEntity(ctx context.Context, entityID string, entityType string, entityData interface{}) error
	BatchIndexEntities(ctx context.Context, entities []EntityIndexRequest) error
	ReindexEntity(ctx context.Context, entityID string, entityType string, entityData interface{}) error

	// Entity removal operations
	RemoveEntity(ctx context.Context, entityID string) error
	BatchRemoveEntities(ctx context.Context, entityIDs []string) error
	RemoveEntitiesByType(ctx context.Context, entityType string) error

	// Search operations
	SearchSimilarEntities(ctx context.Context, params VectorSearchParams) (*VectorSearchResults, error)
	SearchByVector(ctx context.Context, vector []float32, params VectorSearchParams) (*VectorSearchResults, error)
	FindSimilarToEntity(ctx context.Context, entityID string, entityType string, params VectorSearchParams) (*VectorSearchResults, error)

	// Entity vector operations
	GetEntityVector(ctx context.Context, entityID string) ([]float32, error)
	HasEntityVector(ctx context.Context, entityID string) (bool, error)
	GetEntityVectors(ctx context.Context, entityIDs []string) (map[string][]float32, error)

	// Health and diagnostics
	HealthCheck(ctx context.Context) error
	GetVectorStats(ctx context.Context) (*VectorRepositoryStats, error)
	ValidateVectorIndex(ctx context.Context) error
}

// VectorSearchParams defines parameters for vector search operations
type VectorSearchParams struct {
	EntityTypes    []string               `json:"entity_types,omitempty"`
	TopK           int64                  `json:"top_k"`
	ScoreThreshold float64                `json:"score_threshold,omitempty"`
	Filters        map[string]interface{} `json:"filters,omitempty"`
	IncludeVector  bool                   `json:"include_vector"`
	IncludeEntity  bool                   `json:"include_entity"`
}

// VectorSearchResults contains the results of a vector search
type VectorSearchResults struct {
	Results    []VectorSearchResult `json:"results"`
	TotalCount int64                `json:"total_count"`
	SearchTime time.Duration        `json:"search_time"`
	Query      VectorSearchParams   `json:"query"`
}

// VectorSearchResult represents a single search result
type VectorSearchResult struct {
	EntityID    string                 `json:"entity_id"`
	EntityType  string                 `json:"entity_type"`
	Score       float64                `json:"score"`
	Distance    float64                `json:"distance"`
	Vector      []float32              `json:"vector,omitempty"`
	Entity      interface{}            `json:"entity,omitempty"`
	Metadata    map[string]interface{} `json:"metadata"`
	Explanation string                 `json:"explanation,omitempty"`
}

// VectorRepositoryStats provides statistics about the vector repository
type VectorRepositoryStats struct {
	TotalVectors      int64                  `json:"total_vectors"`
	VectorsByType     map[string]int64       `json:"vectors_by_type"`
	IndexSize         int64                  `json:"index_size_bytes"`
	MemoryUsage       int64                  `json:"memory_usage_bytes"`
	AverageQueryTime  time.Duration          `json:"average_query_time"`
	IndexingRate      float64                `json:"indexing_rate_per_second"`
	EmbeddingProvider string                 `json:"embedding_provider"`
	VectorProvider    string                 `json:"vector_provider"`
	HealthStatus      string                 `json:"health_status"`
	LastUpdated       time.Time              `json:"last_updated"`
	ProviderStats     map[string]interface{} `json:"provider_stats"`
}

// ProductCacheRepositoryPort abstracts product cache data access
type ProductCacheRepositoryPort interface {
	Find(ctx context.Context, productID string) (*models.Product, error)
	Search(ctx context.Context, filters map[string]interface{}) ([]*models.Product, error)
	GetCatalog(ctx context.Context, userID string) ([]*models.Product, error)
	BatchFind(ctx context.Context, productIDs []string) ([]*models.Product, error)
}

// PostCacheRepositoryPort abstracts post cache data access
type PostCacheRepositoryPort interface {
	Find(ctx context.Context, postID string) (*models.Post, error)
	Search(ctx context.Context, filters map[string]interface{}) ([]*models.Post, error)
	GetCatalog(ctx context.Context, userID string) ([]*models.Post, error)
	BatchFind(ctx context.Context, postIDs []string) ([]*models.Post, error)
}

// UserCacheRepositoryPort abstracts user cache data access
type UserCacheRepositoryPort interface {
	Find(ctx context.Context, userID string) (*models.User, error)
	Search(ctx context.Context, filters map[string]interface{}) ([]*models.User, error)
	BatchFind(ctx context.Context, userIDs []string) ([]*models.User, error)
}

// ServiceCacheRepositoryPort abstracts service cache data access
type ServiceCacheRepositoryPort interface {
	Find(ctx context.Context, serviceID string) (*models.Service, error)
	Search(ctx context.Context, filters map[string]interface{}) ([]*models.Service, error)
	GetCatalog(ctx context.Context, userID string) ([]*models.Service, error)
	BatchFind(ctx context.Context, serviceIDs []string) ([]*models.Service, error)
}

// VariantCacheRepositoryPort abstracts variant cache data access
type VariantCacheRepositoryPort interface {
	Find(ctx context.Context, variantID string) (*models.Variant, error)
	Search(ctx context.Context, filters map[string]interface{}) ([]*models.Variant, error)
	BatchFind(ctx context.Context, variantIDs []string) ([]*models.Variant, error)
}

// OrderRepositoryPort abstracts order data access
type OrderRepositoryPort interface {
	Find(ctx context.Context, orderID string) (*models.Order, error)
	Search(ctx context.Context, filters map[string]interface{}) ([]*models.Order, error)
	BatchFind(ctx context.Context, orderIDs []string) ([]*models.Order, error)
}

// ItemMetricRepositoryPort abstracts item metric data access
type ItemMetricRepositoryPort interface {
	Find(ctx context.Context, entityType, entityID string) (*models.ItemMetric, error)
	Search(ctx context.Context, filters map[string]interface{}) ([]*models.ItemMetric, error)
	BatchFind(ctx context.Context, requests []MetricRequest) ([]*models.ItemMetric, error)
}

// ===============================
// REPOSITORY COLLECTION INTERFACE
// ===============================

// RepositoryCollection provides access to all repositories through ports
type RepositoryCollection interface {
	Products() ProductCacheRepositoryPort
	Posts() PostCacheRepositoryPort
	Users() UserCacheRepositoryPort
	Services() ServiceCacheRepositoryPort
	Variants() VariantCacheRepositoryPort
	Orders() OrderRepositoryPort
	ItemMetrics() ItemMetricRepositoryPort
	Vectors() VectorRepositoryPort
}

// ===============================
// CIRCUIT BREAKER INTERFACE
// ===============================

// CircuitBreaker defines the circuit breaker interface for fault tolerance
type CircuitBreaker interface {
	Execute(ctx context.Context, operation func() (interface{}, error)) (interface{}, error)
	GetState() string
	GetMetrics() CircuitBreakerMetrics
	Reset() error
	IsHealthy() bool
}

// ===============================
// CACHE MANAGER INTERFACE
// ===============================

// CacheManager defines the caching interface
type CacheManager interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, duration time.Duration) error
	Delete(ctx context.Context, key string) error
	BatchGet(ctx context.Context, keys []string) (map[string]interface{}, error)
	BatchSet(ctx context.Context, items map[string]CacheItem) error
	Clear(ctx context.Context) error
	GetStats(ctx context.Context) (*CacheStats, error)
}

// ===============================
// METRICS COLLECTOR INTERFACE
// ===============================

// MetricsCollector defines the metrics collection interface
type MetricsCollector interface {
	RecordLatency(ctx context.Context, operation string, duration time.Duration, tags map[string]string) error
	IncrementCounter(ctx context.Context, metric string, tags map[string]string) error
	RecordGauge(ctx context.Context, metric string, value float64, tags map[string]string) error
	RecordHistogram(ctx context.Context, metric string, value float64, tags map[string]string) error
	GetMetrics(ctx context.Context, metric string) (interface{}, error)
}

// ===============================
// HEALTH CHECKER INTERFACE
// ===============================

// HealthChecker defines the health checking interface
type HealthChecker interface {
	CheckHealth(ctx context.Context) error
	GetHealthStatus(ctx context.Context) (*HealthStatus, error)
	RegisterHealthCheck(name string, checker func(ctx context.Context) error) error
	GetDependencyHealth(ctx context.Context) (map[string]*DependencyHealth, error)
}

// ===============================
// CONFIGURATION PROVIDER INTERFACE
// ===============================

// ConfigurationProvider defines the configuration interface
type ConfigurationProvider interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetFloat(key string) float64
	GetDuration(key string) time.Duration
	GetStringSlice(key string) []string
	Set(key string, value interface{}) error
	Reload(ctx context.Context) error
	Watch(key string, callback func(interface{})) error
}

// ===============================
// DATA STRUCTURES
// ===============================

// VectorPoint represents a point in the vector space
type VectorPoint struct {
	ID       string                 `json:"id"`
	Vector   []float32              `json:"vector"`
	Metadata map[string]interface{} `json:"metadata"`
	Score    float32                `json:"score,omitempty"`
}

// SearchRequest represents a vector search request
type SearchRequest struct {
	Vector         []float32              `json:"vector"`
	TopK           int64                  `json:"top_k"`
	ScoreThreshold float32                `json:"score_threshold,omitempty"`
	Filter         map[string]interface{} `json:"filter,omitempty"`
	WithVector     bool                   `json:"with_vector"`
	EntityTypes    []string               `json:"entity_types,omitempty"`
}

// SearchResponse contains search results
type SearchResponse struct {
	Points []VectorPoint `json:"points"`
	Total  int64         `json:"total"`
	Took   time.Duration `json:"took"`
}

// VectorDBStats contains vector database statistics
type VectorDBStats struct {
	TotalVectors   int64                  `json:"total_vectors"`
	IndexSize      int64                  `json:"index_size_bytes"`
	MemoryUsage    int64                  `json:"memory_usage_bytes"`
	QueryLatency   time.Duration          `json:"avg_query_latency"`
	IndexingRate   float64                `json:"indexing_rate_per_second"`
	Collections    map[string]interface{} `json:"collections"`
	ConnectionPool *ConnectionPoolStats   `json:"connection_pool"`
	LastUpdated    time.Time              `json:"last_updated"`
}

// VectorDBProviderInfo contains provider information
type VectorDBProviderInfo struct {
	Provider    string    `json:"provider"`
	Version     string    `json:"version"`
	Endpoint    string    `json:"endpoint"`
	Collection  string    `json:"collection"`
	Dimensions  int       `json:"dimensions"`
	IndexType   string    `json:"index_type"`
	ConnectedAt time.Time `json:"connected_at"`
}

// EmbeddingProviderInfo contains embedding provider information
type EmbeddingProviderInfo struct {
	Provider     string    `json:"provider"`
	Model        string    `json:"model"`
	Dimensions   int       `json:"dimensions"`
	MaxTokens    int       `json:"max_tokens"`
	RateLimit    int       `json:"rate_limit_per_minute"`
	CostPerToken float64   `json:"cost_per_token"`
	ConnectedAt  time.Time `json:"connected_at"`
}

// ExternalServiceInfo contains external service information
type ExternalServiceInfo struct {
	ServiceName  string            `json:"service_name"`
	Endpoint     string            `json:"endpoint"`
	Version      string            `json:"version"`
	Status       string            `json:"status"`
	LastSeen     time.Time         `json:"last_seen"`
	Capabilities []string          `json:"capabilities"`
	Metadata     map[string]string `json:"metadata"`
}

// LLMProviderInfo contains LLM provider information
type LLMProviderInfo struct {
	Provider     string    `json:"provider"`
	Model        string    `json:"model"`
	MaxTokens    int       `json:"max_tokens"`
	Temperature  float64   `json:"temperature"`
	RateLimit    int       `json:"rate_limit_per_minute"`
	CostPerToken float64   `json:"cost_per_token"`
	ConnectedAt  time.Time `json:"connected_at"`
}

// CircuitBreakerMetrics contains circuit breaker metrics
type CircuitBreakerMetrics struct {
	TotalRequests        int64         `json:"total_requests"`
	SuccessfulRequests   int64         `json:"successful_requests"`
	FailedRequests       int64         `json:"failed_requests"`
	ConsecutiveFailures  int           `json:"consecutive_failures"`
	ConsecutiveSuccesses int           `json:"consecutive_successes"`
	State                string        `json:"state"`
	LastStateChange      time.Time     `json:"last_state_change"`
	AverageLatency       time.Duration `json:"average_latency"`
}

// CacheItem represents an item to be cached
type CacheItem struct {
	Key        string        `json:"key"`
	Value      interface{}   `json:"value"`
	Expiration time.Duration `json:"expiration"`
	Tags       []string      `json:"tags,omitempty"`
}

// CacheStats contains cache statistics
type CacheStats struct {
	HitRate      float64   `json:"hit_rate"`
	MissRate     float64   `json:"miss_rate"`
	ItemCount    int64     `json:"item_count"`
	MemoryUsage  int64     `json:"memory_usage_bytes"`
	Evictions    int64     `json:"evictions"`
	LastAccessed time.Time `json:"last_accessed"`
}

// ConnectionPoolStats contains connection pool statistics
type ConnectionPoolStats struct {
	ActiveConnections int `json:"active_connections"`
	IdleConnections   int `json:"idle_connections"`
	MaxConnections    int `json:"max_connections"`
	WaitingRequests   int `json:"waiting_requests"`
}

// HealthStatus represents overall health status
type HealthStatus struct {
	Status       string                       `json:"status"`
	Timestamp    time.Time                    `json:"timestamp"`
	Duration     time.Duration                `json:"check_duration"`
	Version      string                       `json:"version"`
	Dependencies map[string]*DependencyHealth `json:"dependencies"`
	Details      map[string]interface{}       `json:"details"`
}

// DependencyHealth represents health of a dependency
type DependencyHealth struct {
	Name      string        `json:"name"`
	Status    string        `json:"status"`
	Latency   time.Duration `json:"latency"`
	Error     string        `json:"error,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
	Essential bool          `json:"essential"`
}

// MetricRequest represents a request for metrics
type MetricRequest struct {
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
}

// ===============================
// PROCESSOR INTERFACES
// ===============================

// VectorProcessor defines the interface for vector processing operations
type VectorProcessor interface {
	ProcessEntityForIndexing(ctx context.Context, entityType string, entityData map[string]interface{}) (*VectorPoint, error)
	ProcessBatchForIndexing(ctx context.Context, entities []EntityIndexRequest) ([]*VectorPoint, error)
	ProcessSearchQuery(ctx context.Context, query VectorSearchQuery) (*SearchRequest, error)
	OptimizeVector(ctx context.Context, vector []float32, optimization string) ([]float32, error)
}

// EntityIndexRequest represents a request to index an entity
type EntityIndexRequest struct {
	EntityID   string                 `json:"entity_id"`
	EntityType string                 `json:"entity_type"`
	EntityData map[string]interface{} `json:"entity_data"`
	Strategy   string                 `json:"strategy,omitempty"`
}

// VectorSearchQuery represents a vector search query
type VectorSearchQuery struct {
	Query           string                 `json:"query"`
	EntityTypes     []string               `json:"entity_types"`
	TopK            int64                  `json:"top_k"`
	Filters         map[string]interface{} `json:"filters"`
	Optimization    string                 `json:"optimization"`
	IncludeVector   bool                   `json:"include_vector"`
	IncludeMetadata bool                   `json:"include_metadata"`
}
