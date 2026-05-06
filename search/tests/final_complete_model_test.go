package search

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"
)

const (
	finalTestRedisAddr = "localhost:6379"
	finalTestIndexName = "unified_search_final_test"
	finalTestPassword  = "YourStrongPasswordHere"
)

func setupFinalRedisWithUnifiedIndex(t *testing.T) (*redis.Pool, *redisearch.Client) {
	pool := &redis.Pool{
		MaxIdle:   10,
		MaxActive: 20,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", finalTestRedisAddr)
			if err != nil {
				return nil, err
			}

			_, err = conn.Do("AUTH", finalTestPassword)
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

	client := redisearch.NewClientFromPool(pool, finalTestIndexName)

	// Create unified schema exactly like the SearchSystem initRedisearch function
	log.Printf("🔧 Creating unified RediSearch index supporting all entity types")

	schema := redisearch.NewSchema(redisearch.DefaultOptions).
		// Core entity identification
		AddField(redisearch.NewTagFieldOptions("entity_type", redisearch.TagFieldOptions{Sortable: true})).
		// Entity-specific ID fields
		AddField(redisearch.NewTagFieldOptions("id", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("product_id", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("post_id", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("vehicle_id", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("property_id", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("service_id", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("deal_id", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("job_id", redisearch.TagFieldOptions{Sortable: true})).
		// Core fields
		AddField(redisearch.NewTextFieldOptions("name", redisearch.TextFieldOptions{Sortable: true, Weight: 2.0})).
		AddField(redisearch.NewTextFieldOptions("description", redisearch.TextFieldOptions{Sortable: true, Weight: 1.0})).
		AddField(redisearch.NewNumericFieldOptions("base_price", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("status", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("condition", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("brand", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewTagFieldOptions("model", redisearch.TagFieldOptions{Sortable: true})).
		AddField(redisearch.NewGeoField("location")).
		AddField(redisearch.NewTextFieldOptions("tags", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("stock", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("weight", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("height", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("width", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("depth", redisearch.NumericFieldOptions{Sortable: true}))

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

	return pool, client
}

// Test complete model behavior with unified index
func TestFinalCompleteModelBehaviorWithUnifiedIndex(t *testing.T) {
	log.Printf("\n🚀 FINAL COMPLETE MODEL BEHAVIOR TEST WITH UNIFIED INDEX")
	log.Printf("🚀 Testing ALL model types with complete field population")

	pool, client := setupFinalRedisWithUnifiedIndex(t)
	defer pool.Close()

	// Clean up any existing test documents first
	log.Printf("\n=== CLEANUP EXISTING DOCS ===")
	allQuery := redisearch.NewQuery("*").Limit(0, 100)
	docs, _, err := client.Search(allQuery)
	if err == nil {
		for _, doc := range docs {
			client.DeleteDocument(doc.Id)
		}
		log.Printf("Cleaned up %d existing docs", len(docs))
	}

	// Index complete model instances using simple approach
	log.Printf("\n=== INDEXING ALL COMPLETE MODELS ===")

	// Product
	log.Printf("Indexing Complete Product...")
	productDoc := redisearch.NewDocument("test-product-001", 1.0)
	productDoc.Set("entity_type", "product")
	productDoc.Set("product_id", "test-product-001")
	productDoc.Set("name", "Complete Product Test àáâãäå 中文 🚀")
	productDoc.Set("description", "Complete product description with special characters")
	productDoc.Set("base_price", 123456)
	productDoc.Set("status", "active")
	productDoc.Set("condition", "new")
	productDoc.Set("brand", "TestBrand™")
	productDoc.Set("model", "Model-XYZ-2024")
	productDoc.Set("stock", 100)
	productDoc.Set("weight", 1500)
	productDoc.Set("height", 200)
	productDoc.Set("width", 150)
	productDoc.Set("depth", 80)
	productDoc.Set("location", "-74.0060,40.7128")
	productDoc.Set("tags", "electronics,test,special-chars")
	err = client.IndexOptions(redisearch.DefaultIndexingOptions, productDoc)
	if err != nil {
		t.Errorf("❌ Failed to index product: %v", err)
	} else {
		log.Printf("✅ Product indexed successfully")
	}

	// Vehicle
	log.Printf("Indexing Complete Vehicle...")
	vehicleDoc := redisearch.NewDocument("test-vehicle-001", 1.0)
	vehicleDoc.Set("entity_type", "vehicle")
	vehicleDoc.Set("vehicle_id", "test-vehicle-001")
	vehicleDoc.Set("name", "Complete Vehicle Test BMW™ àáâãäå 🚗")
	vehicleDoc.Set("description", "Complete vehicle description")
	vehicleDoc.Set("base_price", 4567890)
	vehicleDoc.Set("status", "active")
	vehicleDoc.Set("condition", "new")
	vehicleDoc.Set("brand", "BMW™")
	vehicleDoc.Set("model", "M3-2024")
	vehicleDoc.Set("location", "-87.6298,41.8781")
	vehicleDoc.Set("tags", "luxury,sports-car,manual")
	err = client.IndexOptions(redisearch.DefaultIndexingOptions, vehicleDoc)
	if err != nil {
		t.Errorf("❌ Failed to index vehicle: %v", err)
	} else {
		log.Printf("✅ Vehicle indexed successfully")
	}

	// Property
	log.Printf("Indexing Complete Property...")
	propertyDoc := redisearch.NewDocument("test-property-001", 1.0)
	propertyDoc.Set("entity_type", "property")
	propertyDoc.Set("property_id", "test-property-001")
	propertyDoc.Set("name", "Complete Property Test Villa™ àáâãäå 🏠")
	propertyDoc.Set("description", "Complete luxury property description")
	propertyDoc.Set("base_price", 5678901)
	propertyDoc.Set("status", "active")
	propertyDoc.Set("condition", "excellent")
	propertyDoc.Set("location", "-73.9857,40.7489")
	propertyDoc.Set("tags", "luxury,villa,waterfront")
	err = client.IndexOptions(redisearch.DefaultIndexingOptions, propertyDoc)
	if err != nil {
		t.Errorf("❌ Failed to index property: %v", err)
	} else {
		log.Printf("✅ Property indexed successfully")
	}

	// Service
	log.Printf("Indexing Complete Service...")
	serviceDoc := redisearch.NewDocument("test-service-001", 1.0)
	serviceDoc.Set("entity_type", "service")
	serviceDoc.Set("service_id", "test-service-001")
	serviceDoc.Set("name", "Complete Service Test IT™ àáâãäå 💻")
	serviceDoc.Set("description", "Complete professional service description")
	serviceDoc.Set("base_price", 1234567)
	serviceDoc.Set("status", "active")
	serviceDoc.Set("location", "-122.4194,37.7749")
	serviceDoc.Set("tags", "IT,consulting,cloud")
	err = client.IndexOptions(redisearch.DefaultIndexingOptions, serviceDoc)
	if err != nil {
		t.Errorf("❌ Failed to index service: %v", err)
	} else {
		log.Printf("✅ Service indexed successfully")
	}

	// Deal
	log.Printf("Indexing Complete Deal...")
	dealDoc := redisearch.NewDocument("test-deal-001", 1.0)
	dealDoc.Set("entity_type", "deal")
	dealDoc.Set("deal_id", "test-deal-001")
	dealDoc.Set("name", "Complete Deal Test Black Friday™ àáâãäå 🎯")
	dealDoc.Set("description", "Complete deal description")
	dealDoc.Set("base_price", 9876543)
	dealDoc.Set("status", "active")
	dealDoc.Set("condition", "new")
	dealDoc.Set("brand", "DealBrand™")
	dealDoc.Set("location", "-118.2437,34.0522")
	dealDoc.Set("tags", "black-friday,electronics,flash-sale")
	err = client.IndexOptions(redisearch.DefaultIndexingOptions, dealDoc)
	if err != nil {
		t.Errorf("❌ Failed to index deal: %v", err)
	} else {
		log.Printf("✅ Deal indexed successfully")
	}

	// Job
	log.Printf("Indexing Complete Job...")
	jobDoc := redisearch.NewDocument("test-job-001", 1.0)
	jobDoc.Set("entity_type", "job")
	jobDoc.Set("job_id", "test-job-001")
	jobDoc.Set("name", "Complete Job Test Senior Developer™ àáâãäå 💼")
	jobDoc.Set("description", "Complete job description")
	jobDoc.Set("base_price", 150000)
	jobDoc.Set("status", "active")
	jobDoc.Set("location", "-122.3321,47.6062")
	jobDoc.Set("tags", "software,development,remote")
	err = client.IndexOptions(redisearch.DefaultIndexingOptions, jobDoc)
	if err != nil {
		t.Errorf("❌ Failed to index job: %v", err)
	} else {
		log.Printf("✅ Job indexed successfully")
	}

	// Post
	log.Printf("Indexing Complete Post...")
	postDoc := redisearch.NewDocument("test-post-001", 1.0)
	postDoc.Set("entity_type", "post")
	postDoc.Set("post_id", "test-post-001")
	postDoc.Set("name", "Complete Post Test Discussion™ àáâãäå 💬")
	postDoc.Set("description", "Complete post description")
	postDoc.Set("status", "active")
	postDoc.Set("location", "-104.9903,39.7392")
	postDoc.Set("tags", "discussion,community,general")
	err = client.IndexOptions(redisearch.DefaultIndexingOptions, postDoc)
	if err != nil {
		t.Errorf("❌ Failed to index post: %v", err)
	} else {
		log.Printf("✅ Post indexed successfully")
	}

	// Wait for indexing to complete
	time.Sleep(200 * time.Millisecond)

	// Test 1: Global search behavior
	log.Printf("\n=== TEST 1: GLOBAL SEARCH BEHAVIOR ===")
	globalQuery := redisearch.NewQuery("*").Limit(0, 20)
	docs, total, err := client.Search(globalQuery)
	if err != nil {
		log.Printf("❌ Global search failed: %v", err)
	} else {
		log.Printf("✅ Global search: found %d docs, %d total", len(docs), total)

		entityTypeCounts := make(map[string]int)
		for _, doc := range docs {
			if entityType, ok := doc.Properties["entity_type"]; ok {
				entityTypeStr := entityType.(string)
				entityTypeCounts[entityTypeStr]++
				log.Printf("  📄 Doc ID: %s, entity_type: %s", doc.Id, entityTypeStr)
			}
		}

		log.Printf("📊 Entity type distribution:")
		for entityType, count := range entityTypeCounts {
			log.Printf("    %s: %d documents", entityType, count)
		}
	}

	// Test 2: Individual entity type searches
	log.Printf("\n=== TEST 2: INDIVIDUAL ENTITY TYPE SEARCHES ===")
	entityTypes := []string{"product", "vehicle", "property", "service", "deal", "job", "post"}

	for _, entityType := range entityTypes {
		log.Printf("\n--- Testing %s search ---", entityType)

		searchQuery := redisearch.NewQuery("@entity_type:{"+entityType+"}").Limit(0, 10)
		docs, total, err := client.Search(searchQuery)

		if err != nil {
			log.Printf("❌ %s search failed: %v", entityType, err)
		} else {
			log.Printf("✅ %s search: found %d docs, %d total", entityType, len(docs), total)

			if len(docs) > 0 {
				doc := docs[0]
				log.Printf("  📄 First result:")
				log.Printf("    ID: %s", doc.Id)
				log.Printf("    entity_type: %v", doc.Properties["entity_type"])
				log.Printf("    name: %v", doc.Properties["name"])

				// Check for weird behaviors
				if entityTypeVal := doc.Properties["entity_type"]; entityTypeVal != entityType {
					log.Printf("  🚨 WEIRD BEHAVIOR: entity_type mismatch! Expected '%s', got '%v'", entityType, entityTypeVal)
				}
			} else {
				log.Printf("  ⚠️ NO RESULTS - WEIRD BEHAVIOR: Documents should exist for %s", entityType)
			}
		}
	}

	// Test 3: Complex queries and edge cases
	log.Printf("\n=== TEST 3: COMPLEX QUERIES AND EDGE CASES ===")

	complexTests := []struct {
		name  string
		query string
	}{
		{"Special characters in name", "@name:(àáâãäå)"},
		{"Unicode search", "@name:(中文)"},
		{"Emoji search", "@name:(🚀)"},
		{"Multiple entity types", "@entity_type:{product|vehicle}"},
		{"Price range", "@base_price:[100000 200000]"},
		{"Geographic search", "@location:[40.7128 -74.0060 10 km]"},
		{"Tag search", "@tags:(electronics)"},
		{"Brand and condition", "@brand:(BMW) @condition:(new)"},
		{"Empty query", ""},
		{"Non-existent entity", "@entity_type:{nonexistent}"},
	}

	for _, test := range complexTests {
		log.Printf("\n--- %s ---", test.name)

		var query *redisearch.Query
		if test.query == "" {
			query = redisearch.NewQuery("*").Limit(0, 5)
		} else {
			query = redisearch.NewQuery(test.query).Limit(0, 5)
		}

		docs, total, err := client.Search(query)
		if err != nil {
			log.Printf("❌ Query failed: %v", err)
		} else {
			log.Printf("✅ Query successful: %d docs, %d total", len(docs), total)

			// Log first result for analysis
			if len(docs) > 0 {
				doc := docs[0]
				log.Printf("  📄 Sample result: ID=%s, entity_type=%v", doc.Id, doc.Properties["entity_type"])
			}
		}
	}

	// Cleanup
	log.Printf("\n=== CLEANUP ===")
	testIds := []string{
		"test-product-001", "test-vehicle-001", "test-property-001",
		"test-service-001", "test-deal-001", "test-job-001", "test-post-001",
	}
	for _, id := range testIds {
		err := client.DeleteDocument(id)
		if err != nil {
			log.Printf("⚠️ Failed to cleanup %s: %v", id, err)
		} else {
			log.Printf("🗑️ Cleaned up %s", id)
		}
	}

	log.Printf("\n🏁 FINAL COMPLETE MODEL BEHAVIOR TEST FINISHED")
}
