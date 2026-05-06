package tools

import (
	"context"
	"fmt"
	"time"

	"middleman/managers/internal/domain"
)

// VariantToolService handles product variant operations
type VariantToolService struct {
	variantRepo domain.VariantRepository
	config      *ServiceConfig
}

// NewVariantToolService creates a new variant tool service
func NewVariantToolService(variantRepo domain.VariantRepository) *VariantToolService {
	return &VariantToolService{
		variantRepo: variantRepo,
		config: &ServiceConfig{
			MaxRetries:      3,
			EnableStreaming: true,
			EnableMetrics:   true,
		},
	}
}

// ExecuteOperation routes variant operations to appropriate handlers
func (s *VariantToolService) ExecuteOperation(
	ctx context.Context,
	operation string,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	// Send initial progress
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "variant_operation",
			Status:   "started",
			Progress: 0,
			Metadata: map[string]interface{}{
				"operation": operation,
				"service":   "VariantToolService",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	var result interface{}
	var err error

	switch operation {
	case "create_variant", "create":
		result, err = s.createVariant(ctx, parameters, streamChan, toolID)
	case "get_variant", "find":
		result, err = s.getVariant(ctx, parameters, streamChan, toolID)
	case "update_variant", "update":
		result, err = s.updateVariant(ctx, parameters, streamChan, toolID)
	case "delete_variant", "delete":
		result, err = s.deleteVariant(ctx, parameters, streamChan, toolID)
	case "list_variants", "list":
		result, err = s.listVariants(ctx, parameters, streamChan, toolID)
	case "get_product_variants":
		result, err = s.getProductVariants(ctx, parameters, streamChan, toolID)
	case "search_variants":
		result, err = s.searchVariants(ctx, parameters, streamChan, toolID)
	case "update_inventory":
		result, err = s.updateInventory(ctx, parameters, streamChan, toolID)
	case "get_inventory":
		result, err = s.getInventory(ctx, parameters, streamChan, toolID)
	default:
		err = fmt.Errorf("unsupported variant operation: %s", operation)
	}

	// Send completion status
	if streamChan != nil {
		status := "completed"
		errorStr := ""
		if err != nil {
			status = "error"
			errorStr = err.Error()
		}

		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "variant_operation",
			Status:   status,
			Progress: 100,
			Result:   result,
			Error:    errorStr,
			Metadata: map[string]interface{}{
				"operation": operation,
				"service":   "VariantToolService",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	return result, err
}

func (s *VariantToolService) createVariant(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	productID := getStringParam(params, "product_id", "")
	name := getStringParam(params, "name", "")
	sku := getStringParam(params, "sku", "")
	price := getInt64Param(params, "price", 0)

	if productID == "" || name == "" {
		return nil, fmt.Errorf("product_id and name are required")
	}

	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "variant_operation",
			Status:   "progress",
			Progress: 30,
			Metadata: map[string]interface{}{
				"step": "creating_variant",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	variantID, err := s.variantRepo.CreateVariant(ctx, productID, name, sku, price)
	if err != nil {
		return nil, fmt.Errorf("failed to create variant: %w", err)
	}

	return map[string]interface{}{
		"variant_id": variantID,
		"product_id": productID,
		"name":       name,
		"sku":        sku,
		"price":      price,
	}, nil
}

func (s *VariantToolService) getVariant(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	variantID := getStringParam(params, "id", "")
	if variantID == "" {
		variantID = getStringParam(params, "variant_id", "")
	}
	if variantID == "" {
		return nil, fmt.Errorf("variant ID is required")
	}

	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "variant_operation",
			Status:   "progress",
			Progress: 30,
			Metadata: map[string]interface{}{
				"step": "getting_variant",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	variant, err := s.variantRepo.Find(ctx, variantID)
	if err != nil {
		return nil, fmt.Errorf("failed to find variant: %w", err)
	}

	return map[string]interface{}{
		"variant": variant,
		"id":      variantID,
	}, nil
}

func (s *VariantToolService) updateVariant(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	variantID := getStringParam(params, "id", "")
	if variantID == "" {
		variantID = getStringParam(params, "variant_id", "")
	}
	if variantID == "" {
		return nil, fmt.Errorf("variant ID is required")
	}

	name := getStringParam(params, "name", "")
	sku := getStringParam(params, "sku", "")
	price := getInt64Param(params, "price", 0)

	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "variant_operation",
			Status:   "progress",
			Progress: 30,
			Metadata: map[string]interface{}{
				"step": "updating_variant",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	err := s.variantRepo.UpdateVariant(ctx, variantID, name, sku, price)
	if err != nil {
		return nil, fmt.Errorf("failed to update variant: %w", err)
	}

	return map[string]interface{}{
		"variant_id": variantID,
		"name":       name,
		"sku":        sku,
		"price":      price,
		"updated":    true,
	}, nil
}

func (s *VariantToolService) deleteVariant(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	variantID := getStringParam(params, "id", "")
	if variantID == "" {
		variantID = getStringParam(params, "variant_id", "")
	}
	if variantID == "" {
		return nil, fmt.Errorf("variant ID is required")
	}

	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "variant_operation",
			Status:   "progress",
			Progress: 30,
			Metadata: map[string]interface{}{
				"step": "deleting_variant",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	err := s.variantRepo.DeleteVariant(ctx, variantID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete variant: %w", err)
	}

	return map[string]interface{}{
		"variant_id": variantID,
		"deleted":    true,
	}, nil
}

func (s *VariantToolService) listVariants(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	page := getInt64Param(params, "page", 1)
	limit := getInt64Param(params, "limit", 20)
	sortBy := getStringParam(params, "sort_by", "created_at")
	sortOrder := getStringParam(params, "sort_order", "desc")

	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "variant_operation",
			Status:   "progress",
			Progress: 30,
			Metadata: map[string]interface{}{
				"step": "listing_variants",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	variants, totalCount, err := s.variantRepo.GetVariants(ctx, page, limit, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to list variants: %w", err)
	}

	return map[string]interface{}{
		"variants":    variants,
		"total_count": totalCount,
		"page":        page,
		"limit":       limit,
		"sort_by":     sortBy,
		"sort_order":  sortOrder,
	}, nil
}

func (s *VariantToolService) getProductVariants(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	productID := getStringParam(params, "product_id", "")
	if productID == "" {
		return nil, fmt.Errorf("product ID is required")
	}

	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "variant_operation",
			Status:   "progress",
			Progress: 30,
			Metadata: map[string]interface{}{
				"step": "getting_product_variants",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	variants, err := s.variantRepo.GetProductVariants(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product variants: %w", err)
	}

	return map[string]interface{}{
		"product_id": productID,
		"variants":   variants,
	}, nil
}

func (s *VariantToolService) searchVariants(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	query := getStringParam(params, "query", "")
	if query == "" {
		query = getStringParam(params, "search_term", "")
	}
	if query == "" {
		return nil, fmt.Errorf("search query is required")
	}

	page := getInt64Param(params, "page", 1)
	limit := getInt64Param(params, "limit", 20)

	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "variant_operation",
			Status:   "progress",
			Progress: 30,
			Metadata: map[string]interface{}{
				"step": "searching_variants",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	variants, totalCount, err := s.variantRepo.SearchVariants(ctx, query, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search variants: %w", err)
	}

	return map[string]interface{}{
		"query":       query,
		"variants":    variants,
		"total_count": totalCount,
		"page":        page,
		"limit":       limit,
	}, nil
}

func (s *VariantToolService) updateInventory(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	variantID := getStringParam(params, "variant_id", "")
	quantity := getInt64Param(params, "quantity", 0)

	if variantID == "" {
		return nil, fmt.Errorf("variant ID is required")
	}

	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "variant_operation",
			Status:   "progress",
			Progress: 30,
			Metadata: map[string]interface{}{
				"step": "updating_inventory",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	err := s.variantRepo.UpdateInventory(ctx, variantID, int(quantity))
	if err != nil {
		return nil, fmt.Errorf("failed to update inventory: %w", err)
	}

	return map[string]interface{}{
		"variant_id": variantID,
		"quantity":   quantity,
		"updated":    true,
	}, nil
}

func (s *VariantToolService) getInventory(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	variantID := getStringParam(params, "variant_id", "")
	if variantID == "" {
		return nil, fmt.Errorf("variant ID is required")
	}

	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "variant_operation",
			Status:   "progress",
			Progress: 30,
			Metadata: map[string]interface{}{
				"step": "getting_inventory",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	quantity, err := s.variantRepo.GetInventory(ctx, variantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}

	return map[string]interface{}{
		"variant_id": variantID,
		"quantity":   quantity,
	}, nil
}
