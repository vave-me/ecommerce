package search

import (
	"context"
	"log"
	"strings"
	"testing"
	"time"

	"middleman/internal/di"
	"middleman/search/internal/constants"
	redisrepo "middleman/search/internal/redis"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"
)

const (
	allTestRedisAddr = "localhost:6379"
	allTestIndexName = "unified_search_all_test"
	allTestPassword  = "YourStrongPasswordHere"
)

func setupAllMethodsTest(t *testing.T) (*redis.Pool, *redisearch.Client, context.Context) {
	pool := &redis.Pool{
		MaxIdle:   10,
		MaxActive: 20,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", allTestRedisAddr)
			if err != nil {
				return nil, err
			}

			_, err = conn.Do("AUTH", allTestPassword)
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

	client := redisearch.NewClientFromPool(pool, allTestIndexName)

	// Create the unified index with the COMPLETE schema from SearchSystem.initRedisearch
	log.Printf("🔧 Creating complete unified RediSearch index")

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

		// Common text fields (TEXT with SORTABLE for threaded search)
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
		AddField(redisearch.NewNumericFieldOptions("rabatt", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("square_footage", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("bedrooms", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("bathrooms", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("year_built", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("mileage", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("performance_hp", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("number_of_owners", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("year", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("salary", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("salary_min", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("salary_max", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("positions_open", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("stock", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("weight", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("height", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("width", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("depth", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("deal_duration", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("status", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("condition", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("type_of_property", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("listing_type", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("user_type", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("brand", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("make", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("model", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("fuel_type", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("transmission", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("transmission_type", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("employment_type", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("seniority_level", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("service_type", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("availability", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("deal_type", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("sku", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("barcode", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("currency_code", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("vin", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewGeoField("location")).
		AddField(redisearch.NewTextFieldOptions("city", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("state", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("country", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("postal_code", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("address", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("company_name", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("provider_name", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("merchant_name", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("created_at", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("updated_at", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("expires_at", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("available_from", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("negotiable", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("has_variants", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("featured", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("verified", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("manage_stock", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("middleman_service", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("accident_free", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("relocation_support", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("third_party_agency", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("is_available", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("has_options", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("tags", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("attributes", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("options", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("images", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("pricing", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("qualifications", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("thumbnail", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("video_url", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("external_url", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("deal_url", redisearch.TextFieldOptions{Sortable: true})).

		// Service-specific text fields
		AddField(redisearch.NewTextFieldOptions("contact", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("faq", redisearch.TextFieldOptions{Sortable: true}))

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

	// Set up DI container with Redis and Redisearch clients
	container := di.New()
	ctx := container.Scoped(context.Background())

	// Add Redis connection to DI
	container.AddScoped(constants.RedisTransactionKey, func(c di.Container) (any, error) {
		return pool.Get(), nil
	})
	container.AddScoped(constants.RedisearchClientKey, func(c di.Container) (any, error) {
		return *client, nil
	})

	return pool, client, ctx
}

func indexTestData(client *redisearch.Client) {
	log.Printf("📝 Indexing comprehensive test data for all entity types...")

	// Product test data
	productDoc := redisearch.NewDocument("test-product-comprehensive", 1.0)
	productDoc.Set("entity_type", "product")
	productDoc.Set("product_id", "test-product-comprehensive")
	productDoc.Set("name", "Test Product Comprehensive àáâãäå 🚀")
	productDoc.Set("description", "Comprehensive product for testing all methods")
	productDoc.Set("base_price", 12345)
	productDoc.Set("status", "active")
	productDoc.Set("condition", "new")
	productDoc.Set("brand", "TestBrand")
	productDoc.Set("model", "TestModel")
	productDoc.Set("category_id", "cat123")
	productDoc.Set("category_slug", "electronics")
	productDoc.Set("stock", 50)
	productDoc.Set("weight", 1000)
	productDoc.Set("height", 100)
	productDoc.Set("width", 200)
	productDoc.Set("depth", 50)
	productDoc.Set("user_seller_id", "seller123")
	productDoc.Set("location", "-74.0060,40.7128")
	productDoc.Set("tags", "electronics,test,gadget")
	client.IndexOptions(redisearch.DefaultIndexingOptions, productDoc)

	// Post test data
	postDoc := redisearch.NewDocument("test-post-comprehensive", 1.0)
	postDoc.Set("entity_type", "post")
	postDoc.Set("post_id", "test-post-comprehensive")
	postDoc.Set("name", "Test Post Comprehensive àáâãäå 💬")
	postDoc.Set("description", "Comprehensive post for testing all methods")
	postDoc.Set("status", "active")
	postDoc.Set("category_id", "cat123")
	postDoc.Set("category_slug", "discussion")
	postDoc.Set("user_id", "user123")
	postDoc.Set("location", "-73.9857,40.7489")
	postDoc.Set("tags", "discussion,community,test")
	client.IndexOptions(redisearch.DefaultIndexingOptions, postDoc)

	// Vehicle test data
	vehicleDoc := redisearch.NewDocument("test-vehicle-comprehensive", 1.0)
	vehicleDoc.Set("entity_type", "vehicle")
	vehicleDoc.Set("vehicle_id", "test-vehicle-comprehensive")
	vehicleDoc.Set("name", "Test Vehicle Comprehensive BMW àáâãäå 🚗")
	vehicleDoc.Set("description", "Comprehensive vehicle for testing all methods")
	vehicleDoc.Set("base_price", 45000)
	vehicleDoc.Set("status", "active")
	vehicleDoc.Set("condition", "used")
	vehicleDoc.Set("brand", "BMW")
	vehicleDoc.Set("model", "X5")
	vehicleDoc.Set("category_id", "cat123")
	vehicleDoc.Set("category_slug", "suv")
	vehicleDoc.Set("year", 2020)
	vehicleDoc.Set("mileage", 25000)
	vehicleDoc.Set("location", "-87.6298,41.8781")
	vehicleDoc.Set("tags", "luxury,suv,awd")
	client.IndexOptions(redisearch.DefaultIndexingOptions, vehicleDoc)

	// Property test data
	propertyDoc := redisearch.NewDocument("test-property-comprehensive", 1.0)
	propertyDoc.Set("entity_type", "property")
	propertyDoc.Set("property_id", "test-property-comprehensive")
	propertyDoc.Set("name", "Test Property Comprehensive Villa àáâãäå 🏠")
	propertyDoc.Set("description", "Comprehensive property for testing all methods")
	propertyDoc.Set("base_price", 850000)
	propertyDoc.Set("status", "active")
	propertyDoc.Set("condition", "excellent")
	propertyDoc.Set("category_id", "cat123")
	propertyDoc.Set("category_slug", "residential")
	propertyDoc.Set("bedrooms", 4)
	propertyDoc.Set("bathrooms", 3)
	propertyDoc.Set("square_footage", 2500)
	propertyDoc.Set("year_built", 2015)
	propertyDoc.Set("location", "-122.4194,37.7749")
	propertyDoc.Set("tags", "luxury,modern,pool")
	client.IndexOptions(redisearch.DefaultIndexingOptions, propertyDoc)

	// Service test data
	serviceDoc := redisearch.NewDocument("test-service-comprehensive", 1.0)
	serviceDoc.Set("entity_type", "service")
	serviceDoc.Set("service_id", "test-service-comprehensive")
	serviceDoc.Set("name", "Test Service Comprehensive IT àáâãäå 💻")
	serviceDoc.Set("description", "Comprehensive service for testing all methods")
	serviceDoc.Set("base_price", 15000)
	serviceDoc.Set("status", "active")
	serviceDoc.Set("category_id", "cat123")
	serviceDoc.Set("category_slug", "technology")
	serviceDoc.Set("service_type", "consulting")
	serviceDoc.Set("provider_name", "TechCorp")
	serviceDoc.Set("location", "-118.2437,34.0522")
	serviceDoc.Set("tags", "IT,consulting,cloud")
	client.IndexOptions(redisearch.DefaultIndexingOptions, serviceDoc)

	// Deal test data
	dealDoc := redisearch.NewDocument("test-deal-comprehensive", 1.0)
	dealDoc.Set("entity_type", "deal")
	dealDoc.Set("deal_id", "test-deal-comprehensive")
	dealDoc.Set("name", "Test Deal Comprehensive Black Friday àáâãäå 🎯")
	dealDoc.Set("description", "Comprehensive deal for testing all methods")
	dealDoc.Set("base_price", 9999)
	dealDoc.Set("status", "active")
	dealDoc.Set("condition", "new")
	dealDoc.Set("brand", "DealBrand")
	dealDoc.Set("category_id", "cat123")
	dealDoc.Set("category_slug", "electronics")
	dealDoc.Set("deal_type", "flash_sale")
	dealDoc.Set("location", "-104.9903,39.7392")
	dealDoc.Set("tags", "black-friday,electronics,discount")
	client.IndexOptions(redisearch.DefaultIndexingOptions, dealDoc)

	// Job test data
	jobDoc := redisearch.NewDocument("test-job-comprehensive", 1.0)
	jobDoc.Set("entity_type", "job")
	jobDoc.Set("job_id", "test-job-comprehensive")
	jobDoc.Set("name", "Test Job Comprehensive Senior Developer àáâãäå 💼")
	jobDoc.Set("description", "Comprehensive job for testing all methods")
	jobDoc.Set("salary", 120000)
	jobDoc.Set("status", "active")
	jobDoc.Set("category_id", "cat123")
	jobDoc.Set("category_slug", "technology")
	jobDoc.Set("employment_type", "full_time")
	jobDoc.Set("seniority_level", "senior")
	jobDoc.Set("company_name", "TechCompany")
	jobDoc.Set("location", "-122.3321,47.6062")
	jobDoc.Set("tags", "software,development,remote")
	client.IndexOptions(redisearch.DefaultIndexingOptions, jobDoc)

	log.Printf("✅ Test data indexed successfully")
}

func TestAllRepositoryMethods(t *testing.T) {
	log.Printf("\n🚀 TESTING ALL REPOSITORY METHODS COMPREHENSIVELY")

	pool, client, ctx := setupAllMethodsTest(t)
	defer pool.Close()

	// Clean existing data
	allQuery := redisearch.NewQuery("*").Limit(0, 100)
	docs, _, _ := client.Search(allQuery)
	for _, doc := range docs {
		client.DeleteDocument(doc.Id)
	}

	// Index comprehensive test data
	indexTestData(client)
	time.Sleep(200 * time.Millisecond) // Allow indexing to complete

	// Initialize all repositories with nil fallbacks for testing
	productRepo := redisrepo.NewProductCacheRepository("idx:products", nil)
	postRepo := redisrepo.NewPostCacheRepository("idx:posts", nil)
	vehicleRepo := redisrepo.NewVehicleCacheRepository("idx:vehicles", nil)
	propertyRepo := redisrepo.NewPropertyCacheRepository("idx:properties", nil)
	serviceRepo := redisrepo.NewServiceCacheRepository("idx:services", nil)
	dealRepo := redisrepo.NewDealCacheRepository("idx:deals", nil)
	jobRepo := redisrepo.NewJobCacheRepository("idx:jobs", nil)
	variantRepo := redisrepo.NewVariantCacheRepository("idx:variants", nil)

	// TEST PRODUCT REPOSITORY METHODS
	log.Printf("\n=== TESTING PRODUCT REPOSITORY METHODS ===")

	log.Printf("\n--- ProductCacheRepository.SuggestProducts ---")
	products, err := productRepo.SuggestProducts(ctx, "Test")
	if err != nil {
		log.Printf("❌ SuggestProducts failed: %v", err)
	} else {
		log.Printf("✅ SuggestProducts: found %d products", len(products))
	}

	log.Printf("\n--- ProductCacheRepository.SearchWithTerm ---")
	products, err = productRepo.SearchWithTerm(ctx, "Comprehensive")
	if err != nil {
		log.Printf("❌ SearchWithTerm failed: %v", err)
	} else {
		log.Printf("✅ SearchWithTerm: found %d products", len(products))
	}

	log.Printf("\n--- ProductCacheRepository.SearchWithFilters ---")
	products, err = productRepo.SearchWithFilters(ctx, "Test", "electronics", int64(0), int64(99999),
		"TestBrand", "new", "TestModel", []string{"electronics"}, false, int64(0), int64(0), "",
		"active", false, "", false, false, int64(0), int64(0), int64(0), int64(0), int64(0), int64(0), int64(0), int64(0),
		int64(0), int64(50), int64(0), float64(0), float64(0), int64(0), int64(1), int64(20), "", "")
	if err != nil {
		log.Printf("❌ SearchWithFilters failed: %v", err)
	} else {
		log.Printf("✅ SearchWithFilters: found %d products", len(products))
	}

	log.Printf("\n--- ProductCacheRepository.SearchProductsWithCategorySlug ---")
	products, err = productRepo.SearchProductsWithCategorySlug(ctx, "electronics", 1, 20, "", "")
	if err != nil {
		log.Printf("❌ SearchProductsWithCategorySlug failed: %v", err)
	} else {
		log.Printf("✅ SearchProductsWithCategorySlug: found %d products", len(products))
	}

	log.Printf("\n--- ProductCacheRepository.SearchProductsWithCategory ---")
	products, err = productRepo.SearchProductsWithCategory(ctx, "cat123", 1, 20, "", "")
	if err != nil {
		log.Printf("❌ SearchProductsWithCategory failed: %v", err)
	} else {
		log.Printf("✅ SearchProductsWithCategory: found %d products", len(products))
	}

	// TEST POST REPOSITORY METHODS
	log.Printf("\n=== TESTING POST REPOSITORY METHODS ===")

	log.Printf("\n--- PostCacheRepository.SuggestPosts ---")
	posts, err := postRepo.SuggestPosts(ctx, "Test")
	if err != nil {
		log.Printf("❌ SuggestPosts failed: %v", err)
	} else {
		log.Printf("✅ SuggestPosts: found %d posts", len(posts))
	}

	log.Printf("\n--- PostCacheRepository.SearchWithTerm ---")
	posts, err = postRepo.SearchWithTerm(ctx, "Comprehensive")
	if err != nil {
		log.Printf("❌ SearchWithTerm failed: %v", err)
	} else {
		log.Printf("✅ SearchWithTerm: found %d posts", len(posts))
	}

	log.Printf("\n--- PostCacheRepository.SearchPostsWithFilters ---")
	posts, err = postRepo.SearchPostsWithFilters(ctx, "Test", "Comprehensive", []string{"discussion"},
		"active", "", "", 0, 50, 0, 0, 0, 1, 20, "", "")
	if err != nil {
		log.Printf("❌ SearchPostsWithFilters failed: %v", err)
	} else {
		log.Printf("✅ SearchPostsWithFilters: found %d posts", len(posts))
	}

	log.Printf("\n--- PostCacheRepository.SearchPostsWithCategorySlug ---")
	posts, err = postRepo.SearchPostsWithCategorySlug(ctx, "discussion", 1, 20, "", "")
	if err != nil {
		log.Printf("❌ SearchPostsWithCategorySlug failed: %v", err)
	} else {
		log.Printf("✅ SearchPostsWithCategorySlug: found %d posts", len(posts))
	}

	log.Printf("\n--- PostCacheRepository.SearchPostsWithCategory ---")
	posts, err = postRepo.SearchPostsWithCategory(ctx, "cat123", 1, 20, "", "")
	if err != nil {
		log.Printf("❌ SearchPostsWithCategory failed: %v", err)
	} else {
		log.Printf("✅ SearchPostsWithCategory: found %d posts", len(posts))
	}

	// TEST VEHICLE REPOSITORY METHODS
	log.Printf("\n=== TESTING VEHICLE REPOSITORY METHODS ===")

	log.Printf("\n--- VehicleCacheRepository.SuggestVehicles ---")
	vehicles, err := vehicleRepo.SuggestVehicles(ctx, "Test")
	if err != nil {
		log.Printf("❌ SuggestVehicles failed: %v", err)
	} else {
		log.Printf("✅ SuggestVehicles: found %d vehicles", len(vehicles))
	}

	log.Printf("\n--- VehicleCacheRepository.SearchWithTerm ---")
	vehicles, err = vehicleRepo.SearchWithTerm(ctx, "BMW")
	if err != nil {
		log.Printf("❌ SearchWithTerm failed: %v", err)
	} else {
		log.Printf("✅ SearchWithTerm: found %d vehicles", len(vehicles))
	}

	log.Printf("\n--- VehicleCacheRepository.SearchVehiclesWithCategorySlug ---")
	vehicles, err = vehicleRepo.SearchVehiclesWithCategorySlug(ctx, "suv", 1, 20, "", "")
	if err != nil {
		log.Printf("❌ SearchVehiclesWithCategorySlug failed: %v", err)
	} else {
		log.Printf("✅ SearchVehiclesWithCategorySlug: found %d vehicles", len(vehicles))
	}

	// TEST PROPERTY REPOSITORY METHODS
	log.Printf("\n=== TESTING PROPERTY REPOSITORY METHODS ===")

	log.Printf("\n--- PropertyCacheRepository.SuggestProperties ---")
	properties, err := propertyRepo.SuggestProperties(ctx, "Test")
	if err != nil {
		log.Printf("❌ SuggestProperties failed: %v", err)
	} else {
		log.Printf("✅ SuggestProperties: found %d properties", len(properties))
	}

	log.Printf("\n--- PropertyCacheRepository.SearchWithTerm ---")
	properties, err = propertyRepo.SearchWithTerm(ctx, "Villa")
	if err != nil {
		log.Printf("❌ SearchWithTerm failed: %v", err)
	} else {
		log.Printf("✅ SearchWithTerm: found %d properties", len(properties))
	}

	// TEST SERVICE REPOSITORY METHODS
	log.Printf("\n=== TESTING SERVICE REPOSITORY METHODS ===")

	log.Printf("\n--- ServiceCacheRepository.SuggestServices ---")
	services, err := serviceRepo.SuggestServices(ctx, "Test")
	if err != nil {
		log.Printf("❌ SuggestServices failed: %v", err)
	} else {
		log.Printf("✅ SuggestServices: found %d services", len(services))
	}

	log.Printf("\n--- ServiceCacheRepository.SearchWithTerm ---")
	services, err = serviceRepo.SearchWithTerm(ctx, "IT")
	if err != nil {
		log.Printf("❌ SearchWithTerm failed: %v", err)
	} else {
		log.Printf("✅ SearchWithTerm: found %d services", len(services))
	}

	log.Printf("\n--- ServiceCacheRepository.SearchServicesWithCategorySlug ---")
	services, err = serviceRepo.SearchServicesWithCategorySlug(ctx, "technology", 1, 20, "", "")
	if err != nil {
		log.Printf("❌ SearchServicesWithCategorySlug failed: %v", err)
	} else {
		log.Printf("✅ SearchServicesWithCategorySlug: found %d services", len(services))
	}

	log.Printf("\n--- ServiceCacheRepository.SearchServicesWithCategory ---")
	services, err = serviceRepo.SearchServicesWithCategory(ctx, "cat123", 1, 20, "", "")
	if err != nil {
		log.Printf("❌ SearchServicesWithCategory failed: %v", err)
	} else {
		log.Printf("✅ SearchServicesWithCategory: found %d services", len(services))
	}

	// TEST DEAL REPOSITORY METHODS
	log.Printf("\n=== TESTING DEAL REPOSITORY METHODS ===")

	log.Printf("\n--- DealCacheRepository.SuggestDeals ---")
	deals, err := dealRepo.SuggestDeals(ctx, "Test")
	if err != nil {
		log.Printf("❌ SuggestDeals failed: %v", err)
	} else {
		log.Printf("✅ SuggestDeals: found %d deals", len(deals))
	}

	log.Printf("\n--- DealCacheRepository.SearchWithTerm ---")
	deals, err = dealRepo.SearchWithTerm(ctx, "Black Friday")
	if err != nil {
		log.Printf("❌ SearchWithTerm failed: %v", err)
	} else {
		log.Printf("✅ SearchWithTerm: found %d deals", len(deals))
	}

	log.Printf("\n--- DealCacheRepository.SearchDealsWithCategorySlug ---")
	deals, err = dealRepo.SearchDealsWithCategorySlug(ctx, "electronics", 1, 20, "", "")
	if err != nil {
		log.Printf("❌ SearchDealsWithCategorySlug failed: %v", err)
	} else {
		log.Printf("✅ SearchDealsWithCategorySlug: found %d deals", len(deals))
	}

	log.Printf("\n--- DealCacheRepository.SearchDealsWithCategory ---")
	deals, err = dealRepo.SearchDealsWithCategory(ctx, "cat123", 1, 20, "", "")
	if err != nil {
		log.Printf("❌ SearchDealsWithCategory failed: %v", err)
	} else {
		log.Printf("✅ SearchDealsWithCategory: found %d deals", len(deals))
	}

	// TEST JOB REPOSITORY METHODS
	log.Printf("\n=== TESTING JOB REPOSITORY METHODS ===")

	log.Printf("\n--- JobCacheRepository.SuggestJobs ---")
	jobs, err := jobRepo.SuggestJobs(ctx, "Test")
	if err != nil {
		log.Printf("❌ SuggestJobs failed: %v", err)
	} else {
		log.Printf("✅ SuggestJobs: found %d jobs", len(jobs))
	}

	log.Printf("\n--- JobCacheRepository.SearchJobsWithCategory ---")
	jobs, err = jobRepo.SearchJobsWithCategory(ctx, "cat123", 1, 20, "", "")
	if err != nil {
		log.Printf("❌ SearchJobsWithCategory failed: %v", err)
	} else {
		log.Printf("✅ SearchJobsWithCategory: found %d jobs", len(jobs))
	}

	// TEST VARIANT REPOSITORY METHODS
	log.Printf("\n=== TESTING VARIANT REPOSITORY METHODS ===")

	log.Printf("\n--- VariantCacheRepository.SuggestVariants ---")
	variants, err := variantRepo.SuggestVariants(ctx, "Test")
	if err != nil {
		log.Printf("❌ SuggestVariants failed: %v", err)
	} else {
		log.Printf("✅ SuggestVariants: found %d variants", len(variants))
	}

	log.Printf("\n--- VariantCacheRepository.SearchWithTerm ---")
	variants, err = variantRepo.SearchWithTerm(ctx, "Comprehensive")
	if err != nil {
		log.Printf("❌ SearchWithTerm failed: %v", err)
	} else {
		log.Printf("✅ SearchWithTerm: found %d variants", len(variants))
	}

	// Cleanup
	log.Printf("\n=== CLEANUP ===")
	testIds := []string{
		"test-product-comprehensive", "test-post-comprehensive", "test-vehicle-comprehensive",
		"test-property-comprehensive", "test-service-comprehensive", "test-deal-comprehensive",
		"test-job-comprehensive",
	}
	for _, id := range testIds {
		err := client.DeleteDocument(id)
		if err != nil {
			log.Printf("⚠️ Failed to cleanup %s: %v", id, err)
		} else {
			log.Printf("🗑️ Cleaned up %s", id)
		}
	}

	log.Printf("\n🏁 ALL REPOSITORY METHODS TESTED")
}
