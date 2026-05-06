package tools

import (
	"context"
	"fmt"
	"log"
	"time"

	"middleman/managers/internal/domain"
)

// CategoryToolService handles category management operations
type CategoryToolService struct {
	categories domain.CategoryRepository
}

// NewCategoryToolService creates a new category tool service
func NewCategoryToolService(categoryRepo domain.CategoryRepository) *CategoryToolService {
	return &CategoryToolService{
		categories: categoryRepo,
	}
}

// ExecuteOperation executes a category operation with streaming progress
func (s *CategoryToolService) ExecuteOperation(
	ctx context.Context,
	operation string,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	// Send initial progress
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "category_operation",
		Status:   "progress",
		Progress: 25.0,
		Metadata: map[string]interface{}{
			"step":      "initializing_category_operation",
			"operation": operation,
		},
		Timestamp: time.Now().Unix(),
	}

	switch operation {
	case "add", "create":
		return s.addCategory(ctx, parameters, streamChan, toolID)
	case "get", "find":
		return s.getCategory(ctx, parameters, streamChan, toolID)
	case "get_by_slug":
		return s.getCategoryBySlug(ctx, parameters, streamChan, toolID)
	case "list", "get_all":
		return s.getCategories(ctx, parameters, streamChan, toolID)
	case "get_main":
		return s.getMainCategories(ctx, parameters, streamChan, toolID)
	case "get_sub":
		return s.getSubCategories(ctx, parameters, streamChan, toolID)
	case "update":
		return s.updateCategory(ctx, parameters, streamChan, toolID)
	case "remove", "delete":
		return s.removeCategory(ctx, parameters, streamChan, toolID)
	case "archive":
		return s.archiveCategory(ctx, parameters, streamChan, toolID)
	case "rebrand":
		return s.rebrandCategory(ctx, parameters, streamChan, toolID)
	case "search":
		return s.searchCategories(ctx, parameters, streamChan, toolID)
	default:
		return s.handleUnsupportedOperation(operation, streamChan, toolID)
	}
}

func (s *CategoryToolService) addCategory(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "category_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step": "adding_category",
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("CategoryToolService: Adding category with parameters: %+v", parameters)

	// Send completion (mock implementation)
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "category_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"category_id": "mock_category_123",
			"message":     "Category added successfully",
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"success":     true,
		"operation":   "add_category",
		"category_id": "mock_category_123",
		"message":     "Category added successfully (mock implementation)",
	}, nil
}

func (s *CategoryToolService) getCategory(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	categoryID := getStringParam(parameters, "id", "")

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "category_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":        "getting_category",
			"category_id": categoryID,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("CategoryToolService: Getting category with ID: %s", categoryID)

	// Send completion (mock implementation)
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "category_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"category": map[string]interface{}{
				"id":          categoryID,
				"description": "Mock Category",
				"slug":        "mock-category",
				"is_active":   true,
			},
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"success":   true,
		"operation": "get_category",
		"category": map[string]interface{}{
			"id":          categoryID,
			"description": "Mock Category",
			"slug":        "mock-category",
			"is_active":   true,
		},
		"message": "Category retrieved successfully (mock implementation)",
	}, nil
}

func (s *CategoryToolService) getCategoryBySlug(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	slug := getStringParam(parameters, "slug", "")

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "category_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step": "getting_category_by_slug",
			"slug": slug,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("CategoryToolService: Getting category with slug: %s", slug)

	// Send completion (mock implementation)
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "category_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"category": map[string]interface{}{
				"id":          "mock_category_123",
				"description": "Mock Category",
				"slug":        slug,
				"is_active":   true,
			},
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"success":   true,
		"operation": "get_category_by_slug",
		"category": map[string]interface{}{
			"id":          "mock_category_123",
			"description": "Mock Category",
			"slug":        slug,
			"is_active":   true,
		},
		"message": "Category retrieved by slug successfully (mock implementation)",
	}, nil
}

func (s *CategoryToolService) getCategories(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	page := getInt64Param(parameters, "page", 1)
	pageSize := getInt64Param(parameters, "page_size", 20)

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "category_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":      "getting_categories",
			"page":      page,
			"page_size": pageSize,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("CategoryToolService: Getting categories with page: %d, page_size: %d", page, pageSize)

	// Send completion (mock implementation)
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "category_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"categories": []map[string]interface{}{
				{
					"id":          "cat_1",
					"description": "Electronics",
					"slug":        "electronics",
					"is_active":   true,
				},
				{
					"id":          "cat_2",
					"description": "Clothing",
					"slug":        "clothing",
					"is_active":   true,
				},
			},
			"total_count": 2,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"success":   true,
		"operation": "get_categories",
		"categories": []map[string]interface{}{
			{
				"id":          "cat_1",
				"description": "Electronics",
				"slug":        "electronics",
				"is_active":   true,
			},
			{
				"id":          "cat_2",
				"description": "Clothing",
				"slug":        "clothing",
				"is_active":   true,
			},
		},
		"total_count":  2,
		"current_page": page,
		"message":      "Categories retrieved successfully (mock implementation)",
	}, nil
}

func (s *CategoryToolService) getMainCategories(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "category_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step": "getting_main_categories",
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("CategoryToolService: Getting main categories")

	// Send completion (mock implementation)
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "category_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"main_categories": []map[string]interface{}{
				{
					"id":          "main_1",
					"description": "Electronics",
					"slug":        "electronics",
					"is_main":     true,
				},
			},
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"success":   true,
		"operation": "get_main_categories",
		"main_categories": []map[string]interface{}{
			{
				"id":          "main_1",
				"description": "Electronics",
				"slug":        "electronics",
				"is_main":     true,
			},
		},
		"message": "Main categories retrieved successfully (mock implementation)",
	}, nil
}

func (s *CategoryToolService) getSubCategories(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	parentID := getStringParam(parameters, "parent_id", "")

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "category_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":      "getting_sub_categories",
			"parent_id": parentID,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("CategoryToolService: Getting sub categories for parent: %s", parentID)

	// Send completion (mock implementation)
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "category_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"sub_categories": []map[string]interface{}{
				{
					"id":          "sub_1",
					"description": "Smartphones",
					"slug":        "smartphones",
					"parent_id":   parentID,
				},
			},
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"success":   true,
		"operation": "get_sub_categories",
		"sub_categories": []map[string]interface{}{
			{
				"id":          "sub_1",
				"description": "Smartphones",
				"slug":        "smartphones",
				"parent_id":   parentID,
			},
		},
		"message": "Sub categories retrieved successfully (mock implementation)",
	}, nil
}

func (s *CategoryToolService) updateCategory(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	categoryID := getStringParam(parameters, "id", "")

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "category_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":        "updating_category",
			"category_id": categoryID,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("CategoryToolService: Updating category: %s", categoryID)

	// Send completion (mock implementation)
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "category_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"category_id": categoryID,
			"message":     "Category updated successfully",
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"success":     true,
		"operation":   "update_category",
		"category_id": categoryID,
		"message":     "Category updated successfully (mock implementation)",
	}, nil
}

func (s *CategoryToolService) removeCategory(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	categoryID := getStringParam(parameters, "id", "")

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "category_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":        "removing_category",
			"category_id": categoryID,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("CategoryToolService: Removing category: %s", categoryID)

	// Send completion (mock implementation)
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "category_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"category_id": categoryID,
			"message":     "Category removed successfully",
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"success":     true,
		"operation":   "remove_category",
		"category_id": categoryID,
		"message":     "Category removed successfully (mock implementation)",
	}, nil
}

func (s *CategoryToolService) archiveCategory(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	categoryID := getStringParam(parameters, "id", "")

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "category_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":        "archiving_category",
			"category_id": categoryID,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("CategoryToolService: Archiving category: %s", categoryID)

	// Send completion (mock implementation)
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "category_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"category_id": categoryID,
			"message":     "Category archived successfully",
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"success":     true,
		"operation":   "archive_category",
		"category_id": categoryID,
		"message":     "Category archived successfully (mock implementation)",
	}, nil
}

func (s *CategoryToolService) rebrandCategory(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	categoryID := getStringParam(parameters, "id", "")
	newName := getStringParam(parameters, "new_name", "")

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "category_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":        "rebranding_category",
			"category_id": categoryID,
			"new_name":    newName,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("CategoryToolService: Rebranding category: %s to %s", categoryID, newName)

	// Send completion (mock implementation)
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "category_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"category_id": categoryID,
			"new_name":    newName,
			"message":     "Category rebranded successfully",
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"success":     true,
		"operation":   "rebrand_category",
		"category_id": categoryID,
		"new_name":    newName,
		"message":     "Category rebranded successfully (mock implementation)",
	}, nil
}

func (s *CategoryToolService) searchCategories(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	query := getStringParam(parameters, "query", "")
	limit := getInt64Param(parameters, "limit", 10)

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "category_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":  "searching_categories",
			"query": query,
			"limit": limit,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("CategoryToolService: Searching categories with query: %s", query)

	// Send completion (mock implementation)
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "category_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"categories": []map[string]interface{}{
				{
					"id":          "search_1",
					"description": "Electronics",
					"slug":        "electronics",
					"relevance":   0.95,
				},
			},
			"total_found": 1,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"success":   true,
		"operation": "search_categories",
		"categories": []map[string]interface{}{
			{
				"id":          "search_1",
				"description": "Electronics",
				"slug":        "electronics",
				"relevance":   0.95,
			},
		},
		"total_found": 1,
		"query":       query,
		"message":     "Category search completed successfully (mock implementation)",
	}, nil
}

func (s *CategoryToolService) handleUnsupportedOperation(
	operation string,
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	streamChan <- ToolExecutionStream{
		ID:        toolID,
		ToolName:  "category_operation",
		Status:    "error",
		Progress:  100.0,
		Error:     fmt.Sprintf("Category operation '%s' not supported", operation),
		Timestamp: time.Now().Unix(),
	}

	return nil, fmt.Errorf("unsupported category operation: %s", operation)
}
