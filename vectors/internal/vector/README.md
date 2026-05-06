# Vector Service - Native Vector-Oriented Architecture

This directory contains a pure vector-oriented implementation designed specifically for LLM consumption and semantic search capabilities.

## Architecture Overview

The vector service has been redesigned to be:
- **Native Vector-Oriented**: Core operations focus on pure vector mathematics and similarity search
- **LLM-Friendly**: Accepts raw `[]float32` vectors directly from any LLM
- **Clean API**: Simple, focused interfaces without legacy complexity
- **High Performance**: Optimized for batch operations and real-time search

## Core Components

### VectorService (`service.go`)
The main vector database interface providing:

#### Core Operations
- `IndexVector(ctx, VectorPoint)` - Index a single vector point
- `BatchIndexVectors(ctx, []VectorPoint)` - Efficient batch indexing
- `SearchSimilar(ctx, SearchRequest)` - Semantic similarity search
- `GetVector(ctx, id, withVector)` - Retrieve specific vector by ID
- `DeleteVector(ctx, id)` - Remove vector by ID
- `BatchDeleteVectors(ctx, []id)` - Batch deletion

#### Key Features
- Accepts raw `[]float32` vectors from LLMs
- Flexible metadata storage
- Advanced filtering capabilities
- Cosine similarity search
- Production-ready error handling

### EmbeddingService (`embedding_service.go`)
Generates vector embeddings from text and structured data:

#### Methods
- `GenerateEmbedding(ctx, text)` - Convert text to vector
- `GenerateBatchEmbeddings(ctx, []text)` - Batch text processing
- `GenerateEntityEmbedding(ctx, entityData)` - Structured data to vector

#### Features
- 384-dimensional vectors (configurable)
- Deterministic embedding generation
- Vector normalization
- Priority field extraction from entities

### IntegrationService (`integration.go`)
Handles entity indexing with semantic-aware processing:

#### Entity Indexing
- `IndexProduct(ctx, product)` - Index product entities
- `IndexPost(ctx, post)` - Index post entities  
- `IndexDeal(ctx, deal)` - Index deal entities
- `IndexJob(ctx, job)` - Index job entities
- `IndexProperty(ctx, property)` - Index property entities
- `IndexVehicle(ctx, vehicle)` - Index vehicle entities
- `IndexService(ctx, service)` - Index service entities

#### Generic Operations
- `IndexEntity(ctx, id, type, data)` - Index any entity type
- `BatchIndexEntities(ctx, []entities)` - Efficient batch processing
- `RemoveEntity(ctx, id)` - Remove entity from index

## Data Types

### VectorPoint
```go
type VectorPoint struct {
    ID       string                 // Unique identifier
    Vector   []float32              // Raw vector embedding
    Metadata map[string]interface{} // Searchable metadata
    Score    float32                // Similarity score (search results)
}
```

### SearchRequest
```go
type SearchRequest struct {
    Vector         []float32              // Query vector from LLM
    TopK           int64                  // Number of results
    ScoreThreshold float32                // Minimum similarity score
    Filter         map[string]interface{} // Metadata filters
    WithVector     bool                   // Include vectors in response
}
```

## LLM Integration Pattern

### 1. LLM Generates Vector
```go
// LLM produces embedding
queryVector := []float32{0.1, 0.2, 0.3, ...} // From LLM
```

### 2. Search Vector Database
```go
request := SearchRequest{
    Vector: queryVector,
    TopK:   10,
    Filter: map[string]interface{}{
        "entity_type": "product",
        "price": map[string]interface{}{"lte": 1000},
    },
}

results, err := vectorService.SearchSimilar(ctx, request)
```

### 3. Process Results
```go
for _, point := range results.Points {
    fmt.Printf("ID: %s, Score: %.3f\n", point.ID, point.Score)
    // Access metadata: point.Metadata["name"], point.Metadata["price"], etc.
}
```

## Event-Driven Indexing

The system automatically indexes entities when they're created/updated through integration events:

```go
// Automatic indexing on entity events
productAdded → IndexProduct() → Vector Database
dealUpdated  → IndexDeal()   → Vector Database
postRemoved  → RemoveEntity() → Vector Database
```

## Performance Features

- **Batch Operations**: Efficient bulk indexing and deletion
- **Async Processing**: Non-blocking vector indexing
- **Filtered Search**: Metadata-based result filtering  
- **Normalized Vectors**: Optimized similarity calculations
- **Connection Pooling**: Efficient Qdrant client management

## Configuration

```go
config := vector.Config{
    QdrantHost:     "localhost",
    QdrantPort:     "6334", 
    CollectionName: "entities",
    VectorSize:     384,
}

embeddingConfig := vector.EmbeddingConfig{
    Model:      "all-MiniLM-L6-v2",
    Dimensions: 384,
}
```

## Usage Examples

### Direct Vector Search
```go
// LLM generates query vector
queryVector := llm.GenerateEmbedding("Find red sports cars under $50k")

// Search similar entities
results, err := vectorService.SearchSimilar(ctx, SearchRequest{
    Vector: queryVector,
    TopK:   5,
    Filter: map[string]interface{}{
        "entity_type": "vehicle",
        "price": map[string]interface{}{"lte": 50000},
    },
})
```

### Entity Indexing
```go
// Index a new product
product := &models.Product{
    ProductID:   "prod-123",
    Name:        "Red Sports Car",
    Description: "Fast and beautiful",
    BasePrice:   45000,
    // ... other fields
}

err := integrationService.IndexProduct(ctx, product)
```

### Batch Processing
```go
// Index multiple entities efficiently
entities := []vector.EntityIndexRequest{
    {ID: "1", Type: "product", Data: productData1},
    {ID: "2", Type: "service", Data: serviceData1},
}

err := integrationService.BatchIndexEntities(ctx, entities)
```

## Design Principles

1. **Vector-First**: All operations center around vector similarity
2. **LLM-Native**: Direct integration with LLM-generated embeddings
3. **Metadata-Rich**: Flexible filtering without sacrificing performance
4. **Event-Driven**: Automatic synchronization with entity changes
5. **Production-Ready**: Comprehensive error handling and observability
6. **Clean Interfaces**: Simple, focused APIs without legacy complexity

This architecture provides a solid foundation for LLM-powered semantic search while maintaining high performance and clean separation of concerns. 