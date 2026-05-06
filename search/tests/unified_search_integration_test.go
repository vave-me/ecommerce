package search

import (
	"context"
	"log"
	"strings"
	"testing"
	"time"

	"middleman/internal/di"
	"middleman/search/internal/application"
	"middleman/search/internal/constants"
	mockgrpc "middleman/search/internal/grpc"
	"middleman/search/internal/models"
	redisrepo "middleman/search/internal/redis"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"
)

const (
	unifiedTestRedisAddr = "localhost:6379"
	unifiedTestIndexName = "unified_search_integration_test"
	unifiedTestPassword  = "YourStrongPasswordHere"
)

// Minimal mock order repository for testing
type mockOrderRepository struct{}

func (m *mockOrderRepository) Add(ctx context.Context, order *models.Order) error {
	return nil
}

func (m *mockOrderRepository) UpdateStatus(ctx context.Context, orderID, status string) error {
	return nil
}

func (m *mockOrderRepository) Get(ctx context.Context, orderID string) (*models.Order, error) {
	return &models.Order{OrderID: orderID}, nil
}

// Minimal mock user cache repository for testing
type mockUserCacheRepository struct{}

func (m *mockUserCacheRepository) Add(ctx context.Context, userID, email, username, firstName, lastName, location string, enabled bool) error {
	return nil
}

func (m *mockUserCacheRepository) Rename(ctx context.Context, userID string, firstName string) error {
	return nil
}

func (m *mockUserCacheRepository) Update(ctx context.Context, user *models.User) error {
	return nil
}

func (m *mockUserCacheRepository) Find(ctx context.Context, userID string) (*models.User, error) {
	return &models.User{ID: userID}, nil
}

// Minimal mock item metric repository for testing
type mockItemMetricRepository struct{}

func (m *mockItemMetricRepository) GetItemMetric(ctx context.Context, itemId string) (*models.ItemMetric, error) {
	return nil, nil // Return nil to avoid processing
}

func (m *mockItemMetricRepository) GetItemsMetric(ctx context.Context, itemIds []string) ([]*models.ItemMetric, error) {
	return []*models.ItemMetric{}, nil
}

func (m *mockItemMetricRepository) GetHighestMetricsByType(ctx context.Context, metricType string, req application.MetricSortRequest) ([]*models.ItemMetric, error) {
	return []*models.ItemMetric{}, nil
}

func (m *mockItemMetricRepository) GetLowestMetricsByType(ctx context.Context, metricType string, req application.MetricSortRequest) ([]*models.ItemMetric, error) {
	return []*models.ItemMetric{}, nil
}

func setupUnifiedSearchIntegration(t *testing.T) (*redis.Pool, *redisearch.Client, context.Context, application.Application) {
	// Setup Redis pool
	pool := &redis.Pool{
		MaxIdle:   10,
		MaxActive: 20,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", unifiedTestRedisAddr)
			if err != nil {
				return nil, err
			}

			_, err = conn.Do("AUTH", unifiedTestPassword)
			if err != nil {
				conn.Close()
				return nil, err
			}

			return conn, nil
		},
	}

	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("PING")
	if err != nil {
		t.Fatalf("Cannot connect to Redis: %v", err)
	}

	// Create RediSearch client
	client := redisearch.NewClientFromPool(pool, unifiedTestIndexName)

	// Create the complete unified index schema exactly like SearchSystem.initRedisearch
	log.Printf("🔧 Creating complete unified RediSearch index for integration test")

	schema := redisearch.NewSchema(redisearch.DefaultOptions).
		// Core entity identification
		AddField(redisearch.NewTagFieldOptions("entity_type", redisearch.TagFieldOptions{Sortable: true})).

		// Entity-specific ID fields (TAG for exact matching)
		AddField(redisearch.NewTagFieldOptions("id", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("product_id", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("post_id", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("vehicle_id", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("job_id", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("service_id", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("deal_id", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("variant_id", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("property_id", redisearch.TagFieldOptions{Sortable: true})).

		// User identification fields (TAG for exact matching)
		AddField(redisearch.NewTagFieldOptions("user_id", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("user_seller_id", redisearch.TagFieldOptions{Sortable: true})).

		// Category fields (TAG for exact matching)
		AddField(redisearch.NewTagFieldOptions("category_id", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("category_slug", redisearch.TagFieldOptions{Sortable: true})).

		// Common text fields (TEXT with SORTABLE for search)
		AddField(redisearch.NewTextFieldOptions("name", redisearch.TextFieldOptions{Sortable: true, Weight: 2.0})).
		AddField(redisearch.NewTextFieldOptions("title", redisearch.TextFieldOptions{Sortable: true, Weight: 2.0})).
		AddField(redisearch.NewTextFieldOptions("description", redisearch.TextFieldOptions{Sortable: true, Weight: 1.0})).
		AddField(redisearch.NewTextFieldOptions("content", redisearch.TextFieldOptions{Sortable: true, Weight: 1.0})).

		// Pricing fields (NUMERIC for range filtering)
		AddField(redisearch.NewNumericFieldOptions("price", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("base_price", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("listing_price", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("min_price", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("max_price", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("deal_price", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("variant_price", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("shipping_cost", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("salary", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("salary_min", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("salary_max", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("square_footage", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("bedrooms", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("bathrooms", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("year_built", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("mileage", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("year", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("stock", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("weight", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("height", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("width", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("depth", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("created_at", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("updated_at", redisearch.NumericFieldOptions{Sortable: true})).

		// Tag fields for filtering
		AddField(redisearch.NewTagFieldOptions("status", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("condition", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("property_type", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("user_type", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("brand", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("model", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("fuel_type", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("transmission_type", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("employment_type", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("seniority_level", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("service_type", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("deal_type", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("negotiable", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("has_variants", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("featured", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("manage_stock", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("middleman_service", redisearch.TagFieldOptions{Sortable: true})).

		// Geographic field
		AddField(redisearch.NewGeoField("location")).
		AddField(redisearch.NewTextFieldOptions("city", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("state", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("country", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("address", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("company_name", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("provider_name", redisearch.TextFieldOptions{Sortable: true})).

		// Text fields for content
		AddField(redisearch.NewTextFieldOptions("tags", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("attributes", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("thumbnail", redisearch.TextFieldOptions{Sortable: true})).

		// Property-specific fields
		AddField(redisearch.NewTagFieldOptions("type_of_property", redisearch.TagFieldOptions{Sortable: true})).

		// Tag fields for filtering
		AddField(redisearch.NewTagFieldOptions("status", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("condition", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("property_type", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("user_type", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("brand", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("model", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("fuel_type", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("transmission_type", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("employment_type", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("seniority_level", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("service_type", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("deal_type", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("negotiable", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("has_variants", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("featured", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("manage_stock", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("middleman_service", redisearch.TagFieldOptions{Sortable: true})).

		// Geographic field
		AddField(redisearch.NewGeoField("location")).
		AddField(redisearch.NewTextFieldOptions("city", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("state", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("country", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("address", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("company_name", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("provider_name", redisearch.TextFieldOptions{Sortable: true})).

		// Text fields for content
		AddField(redisearch.NewTextFieldOptions("tags", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("attributes", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("thumbnail", redisearch.TextFieldOptions{Sortable: true}))

	// Create the index
	err = client.CreateIndex(schema)
	if err != nil {
		if strings.Contains(err.Error(), "Index already exists") {
			log.Printf("✅ Unified index already exists")
		} else {
			t.Fatalf("❌ Failed to create unified index: %v", err)
		}
	} else {
		log.Printf("✅ Unified index created successfully")
	}

	// Set up DI container with all dependencies
	container := di.New()
	ctx := container.Scoped(context.Background())

	// Add Redis and RediSearch to DI
	container.AddScoped(constants.RedisTransactionKey, func(c di.Container) (any, error) {
		return pool.Get(), nil
	})
	container.AddScoped(constants.RedisearchClientKey, func(c di.Container) (any, error) {
		return *client, nil
	})

	// Create mock fallback repositories
	productFallback := mockgrpc.NewProductRepository("localhost:9001")
	postFallback := mockgrpc.NewPostRepository("localhost:9002")
	vehicleFallback := mockgrpc.NewVehicleRepository("localhost:9003")
	propertyFallback := mockgrpc.NewPropertyRepository("localhost:9004")
	serviceFallback := mockgrpc.NewServiceRepository("localhost:9005")
	dealFallback := mockgrpc.NewDealRepository("localhost:9006")
	jobFallback := mockgrpc.NewJobRepository("localhost:9007")
	variantFallback := mockgrpc.NewVariantRepository("localhost:9008")

	// Create cache repositories
	productCacheRepo := redisrepo.NewProductCacheRepository("idx:products", productFallback)
	postCacheRepo := redisrepo.NewPostCacheRepository("idx:posts", postFallback)
	vehicleCacheRepo := redisrepo.NewVehicleCacheRepository("idx:vehicles", vehicleFallback)
	propertyCacheRepo := redisrepo.NewPropertyCacheRepository("idx:properties", propertyFallback)
	serviceCacheRepo := redisrepo.NewServiceCacheRepository("idx:services", serviceFallback)
	dealCacheRepo := redisrepo.NewDealCacheRepository("idx:deals", dealFallback)
	jobCacheRepo := redisrepo.NewJobCacheRepository("idx:jobs", jobFallback)
	variantCacheRepo := redisrepo.NewVariantCacheRepository("idx:variants", variantFallback)

	// Create mock repositories
	mockOrderRepo := &mockOrderRepository{}
	mockUserCacheRepo := &mockUserCacheRepository{}

	// Create mock metric repository with proper implementations
	mockMetricRepo := &mockItemMetricRepository{}

	// Create application instance
	app := application.New(
		mockOrderRepo,     // orders
		productCacheRepo,  // products
		variantCacheRepo,  // variants
		postCacheRepo,     // posts
		mockUserCacheRepo, // users
		dealCacheRepo,     // deals
		jobCacheRepo,      // jobs
		propertyCacheRepo, // properties
		vehicleCacheRepo,  // vehicles
		serviceCacheRepo,  // services
		mockMetricRepo,    // metrics - mock to avoid null pointer
	)

	return pool, client, ctx, app
}

func indexComprehensiveUnifiedTestData(client *redisearch.Client) {
	log.Printf("📝 Indexing comprehensive test data for UnifiedSearch and UnifiedFeed...")

	currentTime := time.Now().Unix()

	// Product test data with comprehensive fields
	productDoc := redisearch.NewDocument("unified-product-001", 1.0)
	productDoc.Set("entity_type", "product")
	productDoc.Set("product_id", "unified-product-001")
	productDoc.Set("name", "Unified Product iPhone 15 Pro Max àáâãäå 🚀")
	productDoc.Set("description", "Comprehensive product for unified search testing with special characters")
	productDoc.Set("base_price", 125000) // $1,250.00
	productDoc.Set("status", "active")
	productDoc.Set("condition", "new")
	productDoc.Set("brand", "Apple")
	productDoc.Set("model", "iPhone-15-Pro-Max")
	productDoc.Set("category_id", "cat-electronics")
	productDoc.Set("category_slug", "smartphones")
	productDoc.Set("stock", 50)
	productDoc.Set("weight", 240)
	productDoc.Set("height", 160)
	productDoc.Set("width", 77)
	productDoc.Set("depth", 8)
	productDoc.Set("user_seller_id", "seller-tech-001")
	productDoc.Set("user_type", "business")
	productDoc.Set("location", "-74.0060,40.7128") // New York
	productDoc.Set("city", "New York")
	productDoc.Set("state", "NY")
	productDoc.Set("country", "USA")
	productDoc.Set("tags", "smartphone,apple,5g,pro,camera")
	productDoc.Set("featured", "true")
	productDoc.Set("negotiable", "false")
	productDoc.Set("has_variants", "true")
	productDoc.Set("created_at", currentTime)
	productDoc.Set("updated_at", currentTime)
	client.IndexOptions(redisearch.DefaultIndexingOptions, productDoc)

	// Vehicle test data with comprehensive fields
	vehicleDoc := redisearch.NewDocument("unified-vehicle-001", 1.0)
	vehicleDoc.Set("entity_type", "vehicle")
	vehicleDoc.Set("vehicle_id", "unified-vehicle-001")
	vehicleDoc.Set("name", "Unified Vehicle BMW X5 M Competition àáâãäå 🚗")
	vehicleDoc.Set("description", "Comprehensive vehicle for unified search testing")
	vehicleDoc.Set("base_price", 11500000) // $115,000.00
	vehicleDoc.Set("status", "active")
	vehicleDoc.Set("condition", "excellent")
	vehicleDoc.Set("brand", "BMW")
	vehicleDoc.Set("model", "X5-M-Competition")
	vehicleDoc.Set("category_id", "cat-vehicles")
	vehicleDoc.Set("category_slug", "luxury-suv")
	vehicleDoc.Set("year", 2024)
	vehicleDoc.Set("mileage", 1200)
	vehicleDoc.Set("fuel_type", "gasoline")
	vehicleDoc.Set("transmission_type", "automatic")
	vehicleDoc.Set("user_seller_id", "seller-auto-001")
	vehicleDoc.Set("user_type", "dealer")
	vehicleDoc.Set("location", "-87.6298,41.8781") // Chicago
	vehicleDoc.Set("city", "Chicago")
	vehicleDoc.Set("state", "IL")
	vehicleDoc.Set("country", "USA")
	vehicleDoc.Set("tags", "luxury,suv,performance,awd,leather")
	vehicleDoc.Set("featured", "true")
	vehicleDoc.Set("negotiable", "true")
	vehicleDoc.Set("created_at", currentTime)
	vehicleDoc.Set("updated_at", currentTime)
	client.IndexOptions(redisearch.DefaultIndexingOptions, vehicleDoc)

	// Property test data with comprehensive fields
	propertyDoc := redisearch.NewDocument("unified-property-001", 1.0)
	propertyDoc.Set("entity_type", "property")
	propertyDoc.Set("property_id", "unified-property-001")
	propertyDoc.Set("name", "Unified Property Modern Villa àáâãäå 🏠")
	propertyDoc.Set("description", "Comprehensive property for unified search testing")
	propertyDoc.Set("base_price", 185000000) // $1,850,000.00
	propertyDoc.Set("status", "active")
	propertyDoc.Set("condition", "excellent")
	propertyDoc.Set("type_of_property", "residential")
	propertyDoc.Set("category_id", "cat-real-estate")
	propertyDoc.Set("category_slug", "luxury-homes")
	propertyDoc.Set("bedrooms", 5)
	propertyDoc.Set("bathrooms", 4)
	propertyDoc.Set("square_footage", 4500)
	propertyDoc.Set("year_built", 2022)
	propertyDoc.Set("user_seller_id", "seller-realty-001")
	propertyDoc.Set("user_type", "agent")
	propertyDoc.Set("location", "-122.4194,37.7749") // San Francisco
	propertyDoc.Set("city", "San Francisco")
	propertyDoc.Set("state", "CA")
	propertyDoc.Set("country", "USA")
	propertyDoc.Set("address", "123 Pacific Heights Dr, San Francisco, CA")
	propertyDoc.Set("tags", "luxury,modern,ocean-view,pool,garden")
	propertyDoc.Set("featured", "true")
	propertyDoc.Set("negotiable", "true")
	propertyDoc.Set("created_at", currentTime)
	propertyDoc.Set("updated_at", currentTime)
	client.IndexOptions(redisearch.DefaultIndexingOptions, propertyDoc)

	// Service test data with comprehensive fields
	serviceDoc := redisearch.NewDocument("unified-service-001", 1.0)
	serviceDoc.Set("entity_type", "service")
	serviceDoc.Set("service_id", "unified-service-001")
	serviceDoc.Set("name", "Unified Service Cloud Architecture Consulting àáâãäå 💻")
	serviceDoc.Set("description", "Comprehensive service for unified search testing")
	serviceDoc.Set("base_price", 25000000) // $250,000.00
	serviceDoc.Set("status", "active")
	serviceDoc.Set("service_type", "consulting")
	serviceDoc.Set("category_id", "cat-technology")
	serviceDoc.Set("category_slug", "cloud-services")
	serviceDoc.Set("provider_name", "TechCorp Solutions")
	serviceDoc.Set("user_seller_id", "seller-service-001")
	serviceDoc.Set("user_type", "business")
	serviceDoc.Set("location", "-118.2437,34.0522") // Los Angeles
	serviceDoc.Set("city", "Los Angeles")
	serviceDoc.Set("state", "CA")
	serviceDoc.Set("country", "USA")
	serviceDoc.Set("tags", "cloud,aws,azure,architecture,enterprise")
	serviceDoc.Set("featured", "true")
	serviceDoc.Set("negotiable", "true")
	serviceDoc.Set("created_at", currentTime)
	serviceDoc.Set("updated_at", currentTime)
	client.IndexOptions(redisearch.DefaultIndexingOptions, serviceDoc)

	// Deal test data with comprehensive fields
	dealDoc := redisearch.NewDocument("unified-deal-001", 1.0)
	dealDoc.Set("entity_type", "deal")
	dealDoc.Set("deal_id", "unified-deal-001")
	dealDoc.Set("name", "Unified Deal Black Friday Electronics àáâãäå 🎯")
	dealDoc.Set("description", "Comprehensive deal for unified search testing")
	dealDoc.Set("base_price", 79999) // $799.99
	dealDoc.Set("status", "active")
	dealDoc.Set("condition", "new")
	dealDoc.Set("brand", "Samsung")
	dealDoc.Set("model", "Galaxy-S24-Ultra")
	dealDoc.Set("deal_type", "flash_sale")
	dealDoc.Set("category_id", "cat-electronics")
	dealDoc.Set("category_slug", "smartphones")
	dealDoc.Set("user_seller_id", "seller-deals-001")
	dealDoc.Set("user_type", "business")
	dealDoc.Set("location", "-104.9903,39.7392") // Denver
	dealDoc.Set("city", "Denver")
	dealDoc.Set("state", "CO")
	dealDoc.Set("country", "USA")
	dealDoc.Set("tags", "black-friday,electronics,discount,limited-time")
	dealDoc.Set("featured", "true")
	dealDoc.Set("negotiable", "false")
	dealDoc.Set("created_at", currentTime)
	dealDoc.Set("updated_at", currentTime)
	client.IndexOptions(redisearch.DefaultIndexingOptions, dealDoc)

	// Job test data with comprehensive fields
	jobDoc := redisearch.NewDocument("unified-job-001", 1.0)
	jobDoc.Set("entity_type", "job")
	jobDoc.Set("job_id", "unified-job-001")
	jobDoc.Set("name", "Unified Job Senior Full Stack Developer àáâãäå 💼")
	jobDoc.Set("description", "Comprehensive job for unified search testing")
	jobDoc.Set("salary", 18000000)     // $180,000.00
	jobDoc.Set("salary_min", 16000000) // $160,000.00
	jobDoc.Set("salary_max", 20000000) // $200,000.00
	jobDoc.Set("status", "active")
	jobDoc.Set("employment_type", "full_time")
	jobDoc.Set("seniority_level", "senior")
	jobDoc.Set("category_id", "cat-technology")
	jobDoc.Set("category_slug", "software-development")
	jobDoc.Set("company_name", "InnovateTech Inc")
	jobDoc.Set("user_seller_id", "employer-tech-001")
	jobDoc.Set("user_type", "business")
	jobDoc.Set("location", "-122.3321,47.6062") // Seattle
	jobDoc.Set("city", "Seattle")
	jobDoc.Set("state", "WA")
	jobDoc.Set("country", "USA")
	jobDoc.Set("tags", "react,node,typescript,aws,kubernetes")
	jobDoc.Set("featured", "true")
	jobDoc.Set("created_at", currentTime)
	jobDoc.Set("updated_at", currentTime)
	client.IndexOptions(redisearch.DefaultIndexingOptions, jobDoc)

	// Post test data with comprehensive fields
	postDoc := redisearch.NewDocument("unified-post-001", 1.0)
	postDoc.Set("entity_type", "post")
	postDoc.Set("post_id", "unified-post-001")
	postDoc.Set("name", "Unified Post Tech Discussion àáâãäå 💬")
	postDoc.Set("description", "Comprehensive post for unified search testing")
	postDoc.Set("content", "Discussing the latest trends in technology and AI")
	postDoc.Set("status", "active")
	postDoc.Set("category_id", "cat-discussion")
	postDoc.Set("category_slug", "technology")
	postDoc.Set("user_id", "user-community-001")
	postDoc.Set("user_type", "individual")
	postDoc.Set("location", "-73.9857,40.7489") // Manhattan
	postDoc.Set("city", "Manhattan")
	postDoc.Set("state", "NY")
	postDoc.Set("country", "USA")
	postDoc.Set("tags", "ai,technology,discussion,innovation,future")
	postDoc.Set("featured", "false")
	postDoc.Set("created_at", currentTime)
	postDoc.Set("updated_at", currentTime)
	client.IndexOptions(redisearch.DefaultIndexingOptions, postDoc)

	log.Printf("✅ Comprehensive unified test data indexed successfully")
}

func TestUnifiedSearchIntegration(t *testing.T) {
	log.Printf("\n🚀 COMPREHENSIVE UNIFIED SEARCH INTEGRATION TEST")
	log.Printf("🚀 Testing UnifiedSearch and UnifiedFeed with complete Redis setup")

	pool, client, ctx, app := setupUnifiedSearchIntegration(t)
	defer pool.Close()

	// Clean existing data
	log.Printf("\n=== CLEANUP EXISTING DATA ===")
	allQuery := redisearch.NewQuery("*").Limit(0, 100)
	docs, _, _ := client.Search(allQuery)
	for _, doc := range docs {
		client.DeleteDocument(doc.Id)
	}
	log.Printf("Cleaned up %d existing documents", len(docs))

	// Index comprehensive test data
	indexComprehensiveUnifiedTestData(client)
	time.Sleep(300 * time.Millisecond) // Allow indexing to complete

	// Verify data was indexed
	log.Printf("\n=== VERIFYING INDEXED DATA ===")
	allQuery = redisearch.NewQuery("*").Limit(0, 50)
	docs, total, err := client.Search(allQuery)
	if err != nil {
		t.Fatalf("Failed to verify indexed data: %v", err)
	}
	log.Printf("✅ Verified %d documents indexed, %d total", len(docs), total)

	// Log all indexed documents with their entity types
	entityTypeCounts := make(map[string]int)
	for _, doc := range docs {
		if entityType, ok := doc.Properties["entity_type"]; ok {
			entityTypeStr := entityType.(string)
			entityTypeCounts[entityTypeStr]++
			log.Printf("  📄 Doc ID: %s, entity_type: %s, name: %v",
				doc.Id, entityTypeStr, doc.Properties["name"])
		}
	}

	log.Printf("📊 Entity type distribution:")
	for entityType, count := range entityTypeCounts {
		log.Printf("    %s: %d documents", entityType, count)
	}

	// TEST 1: UNIFIED SEARCH - Global search across all entity types
	log.Printf("\n=== TEST 1: UNIFIED SEARCH - GLOBAL SEARCH ===")

	unifiedSearchParams := application.UnifiedSearchParams{
		SearchTerm:   "Unified",
		EntityTypes:  []string{}, // Empty = search all types
		MinPrice:     0,
		MaxPrice:     0,
		Lat:          0,
		Lng:          0,
		Radius:       0,
		CategorySlug: "",
		UserType:     "",
		Negotiable:   false,
		SortBy:       "",
		SortOrder:    "",
		Page:         1,
		PageSize:     20,
	}

	unifiedResults, err := app.UnifiedSearch(ctx, unifiedSearchParams)
	if err != nil {
		log.Printf("❌ UnifiedSearch failed: %v", err)
		t.Errorf("UnifiedSearch failed: %v", err)
	} else {
		log.Printf("✅ UnifiedSearch success: found %d results, %d total",
			len(unifiedResults.Results), unifiedResults.TotalCount)

		log.Printf("📊 Results by entity type:")
		for entityType, count := range unifiedResults.CountsByType {
			log.Printf("    %s: %d results", entityType, count)
		}

		log.Printf("📋 Top results:")
		for i, result := range unifiedResults.Results {
			if i >= 5 {
				break
			} // Show only top 5

			var name string
			var price int64

			// Extract name and price based on entity type
			switch result.EntityType {
			case "product":
				if result.Product != nil {
					name = result.Product.Name
					price = result.Product.BasePrice
				}
			case "vehicle":
				if result.Vehicle != nil {
					name = result.Vehicle.Name
					price = result.Vehicle.BasePrice
				}
			case "property":
				if result.Property != nil {
					name = result.Property.Name
					price = result.Property.ListingPrice
				}
			case "service":
				if result.Service != nil {
					name = result.Service.Name
					price = result.Service.BasePrice
				}
			case "deal":
				if result.Deal != nil {
					name = result.Deal.Name
					price = result.Deal.BasePrice
				}
			case "job":
				if result.Job != nil {
					name = result.Job.Name
					price = result.Job.Salary
				}
			case "post":
				if result.Post != nil {
					name = result.Post.Name
					price = 0 // Posts don't have price
				}
			}

			log.Printf("  %d. [%s] %s (Score: %.2f, Price: $%.2f)",
				i+1, result.EntityType, name, result.RelevanceScore, float64(price)/100)
		}
	}

	// TEST 2: UNIFIED SEARCH - Specific entity types
	log.Printf("\n=== TEST 2: UNIFIED SEARCH - SPECIFIC ENTITY TYPES ===")

	specificSearchParams := application.UnifiedSearchParams{
		SearchTerm:   "BMW",
		EntityTypes:  []string{"product", "vehicle", "deal"},
		MinPrice:     0,
		MaxPrice:     0,
		Lat:          0,
		Lng:          0,
		Radius:       0,
		CategorySlug: "",
		UserType:     "",
		Negotiable:   false,
		SortBy:       "",
		SortOrder:    "",
		Page:         1,
		PageSize:     10,
	}

	specificResults, err := app.UnifiedSearch(ctx, specificSearchParams)
	if err != nil {
		log.Printf("❌ Specific entity search failed: %v", err)
	} else {
		log.Printf("✅ Specific entity search success: found %d results", len(specificResults.Results))

		for i, result := range specificResults.Results {
			var name string
			var price int64

			switch result.EntityType {
			case "product":
				if result.Product != nil {
					name = result.Product.Name
					price = result.Product.BasePrice
				}
			case "vehicle":
				if result.Vehicle != nil {
					name = result.Vehicle.Name
					price = result.Vehicle.BasePrice
				}
			case "deal":
				if result.Deal != nil {
					name = result.Deal.Name
					price = result.Deal.BasePrice
				}
			}

			log.Printf("  %d. [%s] %s - $%.2f",
				i+1, result.EntityType, name, float64(price)/100)
		}
	}

	// TEST 3: UNIFIED FEED - Get latest content across all types
	log.Printf("\n=== TEST 3: UNIFIED FEED - LATEST CONTENT ===")

	unifiedFeedParams := application.UnifiedFeedParams{
		EntityTypes: []string{}, // All types
		FeedType:    "latest",
		Page:        1,
		PageSize:    15,
		Lat:         0,
		Lng:         0,
		Radius:      0,
		UserID:      "",
	}

	feedResults, err := app.UnifiedFeed(ctx, unifiedFeedParams)
	if err != nil {
		log.Printf("❌ UnifiedFeed failed: %v", err)
	} else {
		log.Printf("✅ UnifiedFeed success: found %d results, %d total",
			len(feedResults.Items), feedResults.TotalCount)

		log.Printf("📋 Latest items in feed:")
		for i, result := range feedResults.Items {
			if i >= 10 {
				break
			} // Show only top 10

			var name string
			var price int64

			switch result.EntityType {
			case "product":
				if result.Product != nil {
					name = result.Product.Name
					price = result.Product.BasePrice
				}
			case "vehicle":
				if result.Vehicle != nil {
					name = result.Vehicle.Name
					price = result.Vehicle.BasePrice
				}
			case "property":
				if result.Property != nil {
					name = result.Property.Name
					price = result.Property.ListingPrice
				}
			case "service":
				if result.Service != nil {
					name = result.Service.Name
					price = result.Service.BasePrice
				}
			case "deal":
				if result.Deal != nil {
					name = result.Deal.Name
					price = result.Deal.BasePrice
				}
			case "job":
				if result.Job != nil {
					name = result.Job.Name
					price = result.Job.Salary
				}
			case "post":
				if result.Post != nil {
					name = result.Post.Name
					price = 0
				}
			}

			log.Printf("  %d. [%s] %s - $%.2f",
				i+1, result.EntityType, name, float64(price)/100)
		}
	}

	// TEST 4: UNIFIED FEED - Trending content
	log.Printf("\n=== TEST 4: UNIFIED FEED - TRENDING CONTENT ===")

	trendingFeedParams := application.UnifiedFeedParams{
		EntityTypes: []string{"vehicle", "property", "job"},
		FeedType:    "trending",
		Page:        1,
		PageSize:    10,
		Lat:         0,
		Lng:         0,
		Radius:      0,
		UserID:      "",
	}

	trendingResults, err := app.UnifiedFeed(ctx, trendingFeedParams)
	if err != nil {
		log.Printf("❌ Trending feed failed: %v", err)
	} else {
		log.Printf("✅ Trending feed success: found %d results", len(trendingResults.Items))

		for i, result := range trendingResults.Items {
			var name string
			var price int64

			switch result.EntityType {
			case "vehicle":
				if result.Vehicle != nil {
					name = result.Vehicle.Name
					price = result.Vehicle.BasePrice
				}
			case "property":
				if result.Property != nil {
					name = result.Property.Name
					price = result.Property.ListingPrice
				}
			case "job":
				if result.Job != nil {
					name = result.Job.Name
					price = result.Job.Salary
				}
			}

			log.Printf("  %d. [%s] %s - $%.2f",
				i+1, result.EntityType, name, float64(price)/100)
		}
	}

	// Cleanup
	log.Printf("\n=== CLEANUP ===")
	testIds := []string{
		"unified-product-001", "unified-vehicle-001", "unified-property-001",
		"unified-service-001", "unified-deal-001", "unified-job-001", "unified-post-001",
	}
	for _, id := range testIds {
		err := client.DeleteDocument(id)
		if err != nil {
			log.Printf("⚠️ Failed to cleanup %s: %v", id, err)
		} else {
			log.Printf("🗑️ Cleaned up %s", id)
		}
	}

	log.Printf("\n🏁 COMPREHENSIVE UNIFIED SEARCH INTEGRATION TEST COMPLETED")
}
