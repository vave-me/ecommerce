package vector

import (
	"context"
	"fmt"
	"log"

	"github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// VectorService provides pure vector database operations for LLM consumption
type VectorService struct {
	client       qdrant.CollectionsClient
	pointsClient qdrant.PointsClient
	conn         *grpc.ClientConn
	config       Config
}

// Config holds vector service configuration
type Config struct {
	QdrantHost     string
	QdrantPort     string
	CollectionName string
	VectorSize     int
}

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
}

// SearchResponse contains search results
type SearchResponse struct {
	Points []VectorPoint `json:"points"`
	Total  int64         `json:"total"`
}

// NewVectorService creates a new vector service instance
func NewVectorService(config Config) (*VectorService, error) {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", config.QdrantHost, config.QdrantPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Qdrant: %w", err)
	}

	service := &VectorService{
		client:       qdrant.NewCollectionsClient(conn),
		pointsClient: qdrant.NewPointsClient(conn),
		conn:         conn,
		config:       config,
	}

	if err := service.initializeCollection(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to initialize collection: %w", err)
	}

	return service, nil
}

// Close closes the vector service connection
func (vs *VectorService) Close() error {
	if vs.conn != nil {
		return vs.conn.Close()
	}
	return nil
}

// initializeCollection creates the collection if it doesn't exist
func (vs *VectorService) initializeCollection(ctx context.Context) error {
	_, err := vs.client.Get(ctx, &qdrant.GetCollectionInfoRequest{
		CollectionName: vs.config.CollectionName,
	})

	if err != nil {
		log.Printf("Creating collection: %s", vs.config.CollectionName)
		_, err = vs.client.Create(ctx, &qdrant.CreateCollection{
			CollectionName: vs.config.CollectionName,
			VectorsConfig: &qdrant.VectorsConfig{
				Config: &qdrant.VectorsConfig_Params{
					Params: &qdrant.VectorParams{
						Size:     uint64(vs.config.VectorSize),
						Distance: qdrant.Distance_Cosine,
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to create collection: %w", err)
		}
	}
	return nil
}

// IndexVector inserts or updates a vector point
func (vs *VectorService) IndexVector(ctx context.Context, point VectorPoint) error {
	payload := vs.buildPayload(point.Metadata)

	qdrantPoint := &qdrant.PointStruct{
		Id: &qdrant.PointId{
			PointIdOptions: &qdrant.PointId_Uuid{Uuid: point.ID},
		},
		Vectors: &qdrant.Vectors{
			VectorsOptions: &qdrant.Vectors_Vector{
				Vector: &qdrant.Vector{Data: point.Vector},
			},
		},
		Payload: payload,
	}

	_, err := vs.pointsClient.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: vs.config.CollectionName,
		Points:         []*qdrant.PointStruct{qdrantPoint},
	})

	return err
}

// BatchIndexVectors inserts or updates multiple vector points
func (vs *VectorService) BatchIndexVectors(ctx context.Context, points []VectorPoint) error {
	if len(points) == 0 {
		return nil
	}

	qdrantPoints := make([]*qdrant.PointStruct, 0, len(points))
	for _, point := range points {
		qdrantPoints = append(qdrantPoints, &qdrant.PointStruct{
			Id: &qdrant.PointId{
				PointIdOptions: &qdrant.PointId_Uuid{Uuid: point.ID},
			},
			Vectors: &qdrant.Vectors{
				VectorsOptions: &qdrant.Vectors_Vector{
					Vector: &qdrant.Vector{Data: point.Vector},
				},
			},
			Payload: vs.buildPayload(point.Metadata),
		})
	}

	_, err := vs.pointsClient.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: vs.config.CollectionName,
		Points:         qdrantPoints,
	})

	return err
}

// SearchSimilar performs vector similarity search
func (vs *VectorService) SearchSimilar(ctx context.Context, req SearchRequest) (*SearchResponse, error) {
	searchPoints := &qdrant.SearchPoints{
		CollectionName: vs.config.CollectionName,
		Vector:         req.Vector,
		Limit:          uint64(req.TopK),
		WithPayload:    &qdrant.WithPayloadSelector{SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: true}},
		WithVectors:    &qdrant.WithVectorsSelector{SelectorOptions: &qdrant.WithVectorsSelector_Enable{Enable: req.WithVector}},
	}

	if req.ScoreThreshold > 0 {
		searchPoints.ScoreThreshold = &req.ScoreThreshold
	}

	if len(req.Filter) > 0 {
		searchPoints.Filter = vs.buildFilter(req.Filter)
	}

	response, err := vs.pointsClient.Search(ctx, searchPoints)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	points := make([]VectorPoint, 0, len(response.Result))
	for _, result := range response.Result {
		point := VectorPoint{
			ID:       result.Id.GetUuid(),
			Score:    result.Score,
			Metadata: vs.extractMetadata(result.Payload),
		}

		if req.WithVector && result.Vectors != nil {
			if vector := result.Vectors.GetVector(); vector != nil {
				point.Vector = vector.Data
			}
		}

		points = append(points, point)
	}

	return &SearchResponse{
		Points: points,
		Total:  int64(len(points)),
	}, nil
}

// GetVector retrieves a specific vector point by ID
func (vs *VectorService) GetVector(ctx context.Context, id string, withVector bool) (*VectorPoint, error) {
	response, err := vs.pointsClient.Get(ctx, &qdrant.GetPoints{
		CollectionName: vs.config.CollectionName,
		Ids: []*qdrant.PointId{{
			PointIdOptions: &qdrant.PointId_Uuid{Uuid: id},
		}},
		WithPayload: &qdrant.WithPayloadSelector{SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: true}},
		WithVectors: &qdrant.WithVectorsSelector{SelectorOptions: &qdrant.WithVectorsSelector_Enable{Enable: withVector}},
	})

	if err != nil {
		return nil, fmt.Errorf("get vector failed: %w", err)
	}

	if len(response.Result) == 0 {
		return nil, fmt.Errorf("vector not found")
	}

	result := response.Result[0]
	point := &VectorPoint{
		ID:       result.Id.GetUuid(),
		Metadata: vs.extractMetadata(result.Payload),
	}

	if withVector && result.Vectors != nil {
		if vector := result.Vectors.GetVector(); vector != nil {
			point.Vector = vector.Data
		}
	}

	return point, nil
}

// DeleteVector removes a vector point by ID
func (vs *VectorService) DeleteVector(ctx context.Context, id string) error {
	_, err := vs.pointsClient.Delete(ctx, &qdrant.DeletePoints{
		CollectionName: vs.config.CollectionName,
		Points: &qdrant.PointsSelector{
			PointsSelectorOneOf: &qdrant.PointsSelector_Points{
				Points: &qdrant.PointsIdsList{
					Ids: []*qdrant.PointId{{
						PointIdOptions: &qdrant.PointId_Uuid{Uuid: id},
					}},
				},
			},
		},
	})

	return err
}

// BatchDeleteVectors removes multiple vector points by IDs
func (vs *VectorService) BatchDeleteVectors(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	pointIds := make([]*qdrant.PointId, 0, len(ids))
	for _, id := range ids {
		pointIds = append(pointIds, &qdrant.PointId{
			PointIdOptions: &qdrant.PointId_Uuid{Uuid: id},
		})
	}

	_, err := vs.pointsClient.Delete(ctx, &qdrant.DeletePoints{
		CollectionName: vs.config.CollectionName,
		Points: &qdrant.PointsSelector{
			PointsSelectorOneOf: &qdrant.PointsSelector_Points{
				Points: &qdrant.PointsIdsList{Ids: pointIds},
			},
		},
	})

	return err
}

// Helper methods

func (vs *VectorService) buildPayload(metadata map[string]interface{}) map[string]*qdrant.Value {
	payload := make(map[string]*qdrant.Value)
	for key, value := range metadata {
		if qdrantValue := vs.convertToQdrantValue(value); qdrantValue != nil {
			payload[key] = qdrantValue
		}
	}
	return payload
}

func (vs *VectorService) buildFilter(filters map[string]interface{}) *qdrant.Filter {
	if len(filters) == 0 {
		return nil
	}

	conditions := make([]*qdrant.Condition, 0, len(filters))
	for key, value := range filters {
		if condition := vs.createCondition(key, value); condition != nil {
			conditions = append(conditions, condition)
		}
	}

	if len(conditions) == 0 {
		return nil
	}

	return &qdrant.Filter{
		Must: conditions,
	}
}

func (vs *VectorService) createCondition(key string, value interface{}) *qdrant.Condition {
	// TODO: Implement proper condition creation based on Qdrant client version
	// For now, return nil to disable complex filtering
	return nil
}

func (vs *VectorService) convertToQdrantValue(value interface{}) *qdrant.Value {
	switch v := value.(type) {
	case string:
		return &qdrant.Value{Kind: &qdrant.Value_StringValue{StringValue: v}}
	case int:
		return &qdrant.Value{Kind: &qdrant.Value_IntegerValue{IntegerValue: int64(v)}}
	case int64:
		return &qdrant.Value{Kind: &qdrant.Value_IntegerValue{IntegerValue: v}}
	case float32:
		return &qdrant.Value{Kind: &qdrant.Value_DoubleValue{DoubleValue: float64(v)}}
	case float64:
		return &qdrant.Value{Kind: &qdrant.Value_DoubleValue{DoubleValue: v}}
	case bool:
		return &qdrant.Value{Kind: &qdrant.Value_BoolValue{BoolValue: v}}
	case []string:
		list := &qdrant.ListValue{}
		for _, item := range v {
			list.Values = append(list.Values, &qdrant.Value{
				Kind: &qdrant.Value_StringValue{StringValue: item},
			})
		}
		return &qdrant.Value{Kind: &qdrant.Value_ListValue{ListValue: list}}
	default:
		return nil
	}
}

func (vs *VectorService) extractMetadata(payload map[string]*qdrant.Value) map[string]interface{} {
	metadata := make(map[string]interface{})
	for key, value := range payload {
		if converted := vs.convertFromQdrantValue(value); converted != nil {
			metadata[key] = converted
		}
	}
	return metadata
}

func (vs *VectorService) convertFromQdrantValue(value *qdrant.Value) interface{} {
	if value == nil {
		return nil
	}

	switch v := value.Kind.(type) {
	case *qdrant.Value_StringValue:
		return v.StringValue
	case *qdrant.Value_IntegerValue:
		return v.IntegerValue
	case *qdrant.Value_DoubleValue:
		return v.DoubleValue
	case *qdrant.Value_BoolValue:
		return v.BoolValue
	case *qdrant.Value_ListValue:
		result := make([]interface{}, 0, len(v.ListValue.Values))
		for _, item := range v.ListValue.Values {
			if converted := vs.convertFromQdrantValue(item); converted != nil {
				result = append(result, converted)
			}
		}
		return result
	default:
		return nil
	}
}
