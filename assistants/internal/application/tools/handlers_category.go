package tools

import (
	"context"
	"fmt"
	"middleman/managers/internal/models"
)

// ==================== CATEGORY HANDLERS ====================
func (r *ComprehensiveToolRegistry) initializeCategoryHandlers() {
	r.handlers["category_find"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		categoryID := getStringParam(params, "category_id")
		if err := ValidateIDParam("category_id", categoryID); err != nil {
			return nil, fmt.Errorf("invalid category_id: %w", err)
		}
		return reg.categoryRepo.FindCategory(ctx, categoryID)
	}

	r.handlers["category_add"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		description := getStringParam(params, "description")
		slug := getStringParam(params, "slug")

		// Validate required parameters
		if description == "" {
			return nil, fmt.Errorf("description is required")
		}
		if slug == "" {
			return nil, fmt.Errorf("slug is required")
		}

		// Sanitize string inputs
		description = SanitizeString(description)
		seoTitle := SanitizeString(getStringParam(params, "seo_title"))
		seoDesc := SanitizeString(getStringParam(params, "seo_desc"))

		req := &models.AddCategoryRequest{
			Description:      description,
			ParentID:         getStringParam(params, "parent_id"),
			GoogleCategoryID: getStringParam(params, "google_category_id"),
			Tags:             getStringArrayParam(params, "tags"),
			Slug:             slug,
			IsActive:         getBoolParam(params, "is_active", true),
			SeoTitle:         seoTitle,
			SeoKeywords:      getStringArrayParam(params, "seo_keywords"),
			SeoDesc:          seoDesc,
			CategoryType:     getStringParam(params, "category_type"),
			Lang:             getStringParam(params, "lang"),
		}
		return reg.categoryRepo.CreateNewCategory(ctx, req)
	}

	r.handlers["category_update"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		categoryID := getStringParam(params, "category_id")
		description := getStringParam(params, "description")

		// Validate required parameters
		if err := ValidateIDParam("category_id", categoryID); err != nil {
			return nil, fmt.Errorf("invalid category_id: %w", err)
		}

		// Sanitize string inputs
		description = SanitizeString(description)
		seoTitle := SanitizeString(getStringParam(params, "seo_title"))
		seoDesc := SanitizeString(getStringParam(params, "seo_desc"))

		req := &models.UpdateCategoryRequest{
			ID:               categoryID,
			Description:      description,
			ParentID:         getStringParam(params, "parent_id"),
			GoogleCategoryID: getStringParam(params, "google_category_id"),
			Tags:             getStringArrayParam(params, "tags"),
			Slug:             getStringParam(params, "slug"),
			IsActive:         getBoolParam(params, "is_active", true),
			SeoTitle:         seoTitle,
			SeoKeywords:      getStringArrayParam(params, "seo_keywords"),
			SeoDesc:          seoDesc,
			CategoryType:     getStringParam(params, "category_type"),
			Lang:             getStringParam(params, "lang"),
		}
		return reg.categoryRepo.ModifyCategoryDetails(ctx, req)
	}

	r.handlers["category_remove"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		categoryID := getStringParam(params, "category_id")
		userID := getStringParam(params, "user_id")
		if err := ValidateIDParam("category_id", categoryID); err != nil {
			return nil, fmt.Errorf("invalid category_id: %w", err)
		}
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		return reg.categoryRepo.DeleteCategoryByID(ctx, categoryID, userID)
	}

	r.handlers["category_get_paginated"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		sortBy := getStringParam(params, "sort_by")
		sortOrder := getStringParam(params, "sort_order")
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		return reg.categoryRepo.GetPaginatedCategoriesWithSorting(ctx, page, pageSize, sortBy, sortOrder)
	}

	r.handlers["category_get_all"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return reg.categoryRepo.AllCategories(ctx)
	}

	r.handlers["category_get_by_parent"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		parentID := getStringParam(params, "parent_id")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 100)
		if err := ValidateIDParam("parent_id", parentID); err != nil {
			return nil, fmt.Errorf("invalid parent_id: %w", err)
		}
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		return reg.categoryRepo.GetChildCategoriesByParentID(ctx, parentID, page, pageSize)
	}

	r.handlers["category_get_root"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		page := getInt64Param(params, "page", 1)
		lang := getStringParam(params, "lang")
		pageSize := getInt64Param(params, "page_size", 100)
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		resp, err := reg.categoryRepo.GetAllTopLevelCategories(ctx, page, lang, pageSize, "", "")
		if err != nil {
			return nil, err
		}
		return resp.Categories, nil
	}

	r.handlers["category_search"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		query := getStringParam(params, "query")
		if query == "" {
			return nil, fmt.Errorf("search query is required")
		}
		query = SanitizeString(query)
		return reg.categoryRepo.SearchCategoriesByKeyword(ctx, query)
	}

	r.handlers["category_get_by_slug"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		slug := getStringParam(params, "slug")
		if slug == "" {
			return nil, fmt.Errorf("slug is required")
		}
		return reg.categoryRepo.FindBySlug(ctx, slug)
	}
}