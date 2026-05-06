package grpc

import (
	"context"
	"fmt"
	"log"
	"middleman/assistants/internal/domain"
	"middleman/assistants/internal/models"
	"middleman/categories/categoriespb"
	"middleman/internal/auth"
	"middleman/internal/rpc"
	"time"

	"google.golang.org/grpc"
)

type CategoryRepository struct {
	endpoint string
	auth     *auth.Auth
}

var _ domain.CategoryRepository = (*CategoryRepository)(nil)

func NewCategoryRepository(endpoint string, auth *auth.Auth) CategoryRepository {
	// Removed log to reduce noise - this is created per request
	return CategoryRepository{
		endpoint: endpoint,
		auth:     auth,
	}
}

// Core category operations from protobuf
func (r CategoryRepository) CreateNewCategory(ctx context.Context, req *models.AddCategoryRequest) (*models.AddCategoryResponse, error) {

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		log.Printf("[CATEGORY_GRPC] Failed to connect to categories service: %v", err)
		return nil, err
	}
	defer conn.Close()

	client := categoriespb.NewCategoriesServiceClient(conn)
	resp, err := client.AddCategory(ctx, &categoriespb.AddCategoryRequest{
		Description:      req.Description,
		ParentId:         req.ParentID,
		GoogleCategoryId: req.GoogleCategoryID,
		Tags:             req.Tags,
		Slug:             req.Slug,
		IsActive:         req.IsActive,
		SeoTitle:         req.SeoTitle,
		SeoKeywords:      req.SeoKeywords,
		SeoDesc:          req.SeoDesc,
		CategoryType:     req.CategoryType,
		Lang:             req.Lang,
	})
	if err != nil {
		log.Printf("[CATEGORY_GRPC] AddCategory RPC failed: %v", err)
		return nil, fmt.Errorf("AddCategory RPC failed: %w", err)
	}

	log.Printf("[CATEGORY_GRPC] AddCategory RPC successful, ID: %s", resp.GetId())
	return &models.AddCategoryResponse{
		ID: resp.GetId(),
	}, nil
}

func (r CategoryRepository) GetCategoryByID(ctx context.Context, id, lang, userID string) (*models.GetCategoryResponse, error) {
	log.Printf("[CATEGORY_GRPC] GetCategory called for ID: %s", id)

	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := categoriespb.NewCategoriesServiceClient(conn)
	resp, err := client.GetCategory(ctx, &categoriespb.GetCategoryRequest{
		Id:     id,
		Lang:   lang,
		UserId: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("GetCategory RPC failed: %w", err)
	}

	return &models.GetCategoryResponse{
		Category: r.categoryToDomain(resp.GetCategory()),
	}, nil
}

func (r CategoryRepository) FindCategoryBySlugName(ctx context.Context, slug, lang, userID string) (*models.GetCategoryBySlugResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := categoriespb.NewCategoriesServiceClient(conn)
	resp, err := client.GetCategoryBySlug(ctx, &categoriespb.GetCategoryBySlugRequest{
		Slug:   slug,
		Lang:   lang,
		UserId: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("GetCategoryBySlug RPC failed: %w", err)
	}

	return &models.GetCategoryBySlugResponse{
		Category: r.categoryToDomain(resp.GetCategory()),
	}, nil
}

func (r CategoryRepository) GetPaginatedCategoryList(ctx context.Context, page int64, lang, categoryType string, pageSize int64, sortBy, sortOrder string) (*models.GetCategoriesResponse, error) {
	log.Printf("[CATEGORY_REPO] GetCategories START - page: %d, lang: %s, type: %s, pageSize: %d",
		page, lang, categoryType, pageSize)

	log.Printf("[CATEGORY_REPO] Dialing categories service at: %s", r.endpoint)
	deadline, hasDeadline := ctx.Deadline()
	if hasDeadline {
		log.Printf("[CATEGORY_REPO] Context deadline: %v", deadline)
	} else {
		log.Printf("[CATEGORY_REPO] Context has NO deadline")
	}

	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		log.Printf("[CATEGORY_REPO] ERROR: Context already cancelled before dial: %v", ctx.Err())
		return nil, ctx.Err()
	default:
		log.Printf("[CATEGORY_REPO] Context is active, proceeding with dial")
	}

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[CATEGORY_REPO] ERROR: Failed to dial: %v", err)
		return nil, err
	}
	defer conn.Close()
	log.Printf("[CATEGORY_REPO] Successfully connected to categories service")

	client := categoriespb.NewCategoriesServiceClient(conn)

	log.Printf("[CATEGORY_REPO] Creating GetCategories request")
	req := &categoriespb.GetCategoriesRequest{
		Page:         page,
		Lang:         lang,
		CategoryType: categoryType,
		PageSize:     pageSize,
		SortBy:       sortBy,
		SortOrder:    sortOrder,
	}

	log.Printf("[CATEGORY_REPO] Calling GetCategories RPC")
	startTime := time.Now()
	resp, err := client.GetCategories(ctx, req)
	elapsed := time.Since(startTime)
	log.Printf("[CATEGORY_REPO] GetCategories RPC completed in %v", elapsed)

	if err != nil {
		log.Printf("[CATEGORY_REPO] ERROR: GetCategories RPC failed: %v", err)
		return nil, fmt.Errorf("GetCategories RPC failed: %w", err)
	}

	log.Printf("[CATEGORY_REPO] Response received - Categories count: %d, TotalCount: %d",
		len(resp.GetCategories()), resp.GetTotalCount())

	categories := make([]*models.Category, 0, len(resp.GetCategories()))
	for i, pbCategory := range resp.GetCategories() {
		if pbCategory != nil {
			log.Printf("[CATEGORY_REPO] Converting category[%d]: ID=%s, Description=%s",
				i, pbCategory.GetId(), pbCategory.GetDescription())
		}
		categories = append(categories, r.categoryToDomain(pbCategory))
	}

	result := &models.GetCategoriesResponse{
		Categories:  categories,
		TotalCount:  resp.GetTotalCount(),
		TotalPages:  resp.GetTotalPages(),
		CurrentPage: resp.GetCurrentPage(),
	}

	log.Printf("[CATEGORY_REPO] GetCategories END - Returning %d categories", len(categories))
	return result, nil
}

func (r CategoryRepository) GetMainParentCategories(ctx context.Context, page int64, lang, categoryType string, pageSize int64, sortBy, sortOrder string) (*models.GetMainCategoriesResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := categoriespb.NewCategoriesServiceClient(conn)
	resp, err := client.GetMainCategories(ctx, &categoriespb.GetMainCategoriesRequest{
		Page:         page,
		Lang:         lang,
		CategoryType: categoryType,
		PageSize:     pageSize,
		SortBy:       sortBy,
		SortOrder:    sortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("GetMainCategories RPC failed: %w", err)
	}

	categories := make([]*models.Category, 0, len(resp.GetCategories()))
	for _, pbCategory := range resp.GetCategories() {
		categories = append(categories, r.categoryToDomain(pbCategory))
	}

	return &models.GetMainCategoriesResponse{
		Categories:  categories,
		TotalCount:  resp.GetTotalCount(),
		TotalPages:  resp.GetTotalPages(),
		CurrentPage: resp.GetCurrentPage(),
	}, nil
}

func (r CategoryRepository) GetAllTopLevelCategories(ctx context.Context, page int64, lang string, pageSize int64, sortBy, sortOrder string) (*models.GetAllMainCategoriesResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := categoriespb.NewCategoriesServiceClient(conn)
	resp, err := client.GetAllMainCategories(ctx, &categoriespb.GetAllMainCategoriesRequest{
		Page:      page,
		Lang:      lang,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("GetAllMainCategories RPC failed: %w", err)
	}

	categories := make([]*models.Category, 0, len(resp.GetCategories()))
	for _, pbCategory := range resp.GetCategories() {
		categories = append(categories, r.categoryToDomain(pbCategory))
	}

	return &models.GetAllMainCategoriesResponse{
		Categories:  categories,
		TotalCount:  resp.GetTotalCount(),
		TotalPages:  resp.GetTotalPages(),
		CurrentPage: resp.GetCurrentPage(),
	}, nil
}

func (r CategoryRepository) GetChildCategoriesForParent(ctx context.Context, parentCategoryID, lang string, page, pageSize int64, sortBy, sortOrder string) (*models.GetSubCategoriesResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := categoriespb.NewCategoriesServiceClient(conn)
	resp, err := client.GetSubCategories(ctx, &categoriespb.GetSubCategoriesRequest{
		ParentCategoryId: parentCategoryID,
		Lang:             lang,
		Page:             page,
		PageSize:         pageSize,
		SortBy:           sortBy,
		SortOrder:        sortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("GetSubCategories RPC failed: %w", err)
	}

	categories := make([]*models.Category, 0, len(resp.GetCategories()))
	for _, pbCategory := range resp.GetCategories() {
		categories = append(categories, r.categoryToDomain(pbCategory))
	}

	return &models.GetSubCategoriesResponse{
		Categories:  categories,
		TotalCount:  resp.GetTotalCount(),
		TotalPages:  resp.GetTotalPages(),
		CurrentPage: resp.GetCurrentPage(),
	}, nil
}

func (r CategoryRepository) DeleteCategoryByID(ctx context.Context, id, userID string) (*models.RemoveCategoryResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := categoriespb.NewCategoriesServiceClient(conn)
	resp, err := client.RemoveCategory(ctx, &categoriespb.RemoveCategoryRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("RemoveCategory RPC failed: %w", err)
	}

	return &models.RemoveCategoryResponse{
		ID: resp.GetId(),
	}, nil
}

func (r CategoryRepository) RenameCategoryBrand(ctx context.Context, req *models.RebrandCategoryRequest) (*models.RebrandCategoryResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := categoriespb.NewCategoriesServiceClient(conn)
	_, err = client.RebrandCategory(ctx, &categoriespb.RebrandCategoryRequest{
		Id:      req.ID,
		NewSlug: req.NewSlug,
		NewDesc: req.NewDesc,
		UserId:  req.UserID,
	})
	if err != nil {
		return nil, fmt.Errorf("RebrandCategory RPC failed: %w", err)
	}

	return &models.RebrandCategoryResponse{
		Success: true,
	}, nil
}

func (r CategoryRepository) ModifyCategoryDetails(ctx context.Context, req *models.UpdateCategoryRequest) (*models.UpdateCategoryResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := categoriespb.NewCategoriesServiceClient(conn)
	resp, err := client.UpdateCategory(ctx, &categoriespb.UpdateCategoryRequest{
		Id:               req.ID,
		Description:      req.Description,
		ParentId:         req.ParentID,
		GoogleCategoryId: req.GoogleCategoryID,
		Tags:             req.Tags,
		Slug:             req.Slug,
		IsActive:         req.IsActive,
		SeoTitle:         req.SeoTitle,
		SeoKeywords:      req.SeoKeywords,
		SeoDesc:          req.SeoDesc,
		CategoryType:     req.CategoryType,
		Lang:             req.Lang,
	})
	if err != nil {
		return nil, fmt.Errorf("UpdateCategory RPC failed: %w", err)
	}

	return &models.UpdateCategoryResponse{
		ID: resp.GetId(),
	}, nil
}

func (r CategoryRepository) MoveCategoryToArchive(ctx context.Context, id, userID string) (*models.ArchiveCategoryResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := categoriespb.NewCategoriesServiceClient(conn)
	resp, err := client.ArchiveCategory(ctx, &categoriespb.ArchiveCategoryRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("ArchiveCategory RPC failed: %w", err)
	}

	return &models.ArchiveCategoryResponse{
		CategoryID: resp.GetCategoryId(),
		Archived:   resp.GetArchived(),
	}, nil
}

// Filter operations from protobuf
func (r CategoryRepository) CreateCategoryFilter(ctx context.Context, req *models.AddFilterRequest) (*models.AddFilterResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := categoriespb.NewCategoriesServiceClient(conn)
	resp, err := client.AddFilter(ctx, &categoriespb.AddFilterRequest{
		CategoryId: req.CategoryID,
		Name:       req.Name,
		FilterType: req.FilterType,
		Values:     req.Values,
		UserId:     req.UserID,
		IsActive:   req.IsActive,
	})
	if err != nil {
		return nil, fmt.Errorf("AddFilter RPC failed: %w", err)
	}

	return &models.AddFilterResponse{
		FilterID: resp.GetFilterId(),
	}, nil
}

func (r CategoryRepository) GetFilterByID(ctx context.Context, filterID, userID string) (*models.GetFilterResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := categoriespb.NewCategoriesServiceClient(conn)
	resp, err := client.GetFilter(ctx, &categoriespb.GetFilterRequest{
		FilterId: filterID,
		UserId:   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("GetFilter RPC failed: %w", err)
	}

	return &models.GetFilterResponse{
		Filter: r.filterToDomain(resp.GetFilter()),
	}, nil
}

func (r CategoryRepository) GetCategoryFiltersList(ctx context.Context, categoryID string, page, pageSize int64, sortBy, sortOrder, userID string) (*models.GetFiltersResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := categoriespb.NewCategoriesServiceClient(conn)
	resp, err := client.GetFilters(ctx, &categoriespb.GetFiltersRequest{
		CategoryId: categoryID,
		Page:       page,
		PageSize:   pageSize,
		SortBy:     sortBy,
		SortOrder:  sortOrder,
		UserId:     userID,
	})
	if err != nil {
		return nil, fmt.Errorf("GetFilters RPC failed: %w", err)
	}

	filters := make([]*models.Filter, 0, len(resp.GetFilters()))
	for _, pbFilter := range resp.GetFilters() {
		filters = append(filters, r.filterToDomain(pbFilter))
	}

	return &models.GetFiltersResponse{
		Filters:     filters,
		TotalCount:  resp.GetTotalCount(),
		TotalPages:  resp.GetTotalPages(),
		CurrentPage: resp.GetCurrentPage(),
	}, nil
}

func (r CategoryRepository) MoveFilterToArchive(ctx context.Context, filterID, userID string) (*models.ArchiveFilterResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := categoriespb.NewCategoriesServiceClient(conn)
	resp, err := client.ArchiveFilter(ctx, &categoriespb.ArchiveFilterRequest{
		FilterId: filterID,
		UserId:   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("ArchiveFilter RPC failed: %w", err)
	}

	return &models.ArchiveFilterResponse{
		FilterID: resp.GetFilterId(),
		Archived: resp.GetArchived(),
	}, nil
}

func (r CategoryRepository) DeleteFilterByID(ctx context.Context, filterID, userID string) (*models.RemoveFilterResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := categoriespb.NewCategoriesServiceClient(conn)
	resp, err := client.RemoveFilter(ctx, &categoriespb.RemoveFilterRequest{
		FilterId: filterID,
		UserId:   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("RemoveFilter RPC failed: %w", err)
	}

	return &models.RemoveFilterResponse{
		FilterID: resp.GetFilterId(),
	}, nil
}

// Additional query methods for AI tooling
func (r CategoryRepository) GetPaginatedCategoriesWithSorting(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*models.Category, error) {
	resp, err := r.GetPaginatedCategoryList(ctx, page, "", "", pageSize, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}
	return resp.Categories, nil
}

func (r CategoryRepository) SearchCategoriesByKeyword(ctx context.Context, term string) ([]*models.Category, error) {
	log.Printf("SearchCategoriesWithTerm called with term: %s (mock implementation)", term)
	return []*models.Category{}, nil
}

func (r CategoryRepository) GetCategoriesFilteredByType(ctx context.Context, categoryType string, page, pageSize int64) ([]*models.Category, error) {
	resp, err := r.GetPaginatedCategoryList(ctx, page, "", categoryType, pageSize, "", "")
	if err != nil {
		return nil, err
	}
	return resp.Categories, nil
}

func (r CategoryRepository) GetActiveCategoriesInLanguage(ctx context.Context, lang string, page, pageSize int64) ([]*models.Category, error) {
	resp, err := r.GetPaginatedCategoryList(ctx, page, lang, "", pageSize, "", "")
	if err != nil {
		return nil, err
	}
	var activeCategories []*models.Category
	for _, cat := range resp.Categories {
		if cat.IsActive {
			activeCategories = append(activeCategories, cat)
		}
	}
	return activeCategories, nil
}

func (r CategoryRepository) GetCategoryStatisticsReport(ctx context.Context, categoryID string) (*models.CategoryStatsResponse, error) {
	log.Printf("GetCategoryStats called for category %s (mock implementation)", categoryID)
	return &models.CategoryStatsResponse{CategoryID: categoryID}, nil
}

func (r CategoryRepository) GetAllFiltersForCategory(ctx context.Context, categoryID string) ([]*models.Filter, error) {
	resp, err := r.GetCategoryFiltersList(ctx, categoryID, 1, 100, "", "", "")
	if err != nil {
		return nil, err
	}
	return resp.Filters, nil
}

func (r CategoryRepository) GetChildCategoriesByParentID(ctx context.Context, parentID string, page, pageSize int64) ([]*models.Category, error) {
	resp, err := r.GetChildCategoriesForParent(ctx, parentID, "", page, pageSize, "", "")
	if err != nil {
		return nil, err
	}
	return resp.Categories, nil
}

func (r CategoryRepository) SearchFiltersMatchingKeyword(ctx context.Context, term string) ([]*models.Filter, error) {
	log.Printf("SearchFiltersWithTerm called with term: %s (mock implementation)", term)
	return []*models.Filter{}, nil
}

func (r CategoryRepository) GetMostPopularCategories(ctx context.Context, limit int64) ([]*models.Category, error) {
	resp, err := r.GetPaginatedCategoryList(ctx, 1, "", "", limit, "view_count", "desc")
	if err != nil {
		return nil, err
	}
	return resp.Categories, nil
}

func (r CategoryRepository) GetCompleteCategoryHierarchy(ctx context.Context, categoryID string) (*models.CategoryHierarchyResponse, error) {
	log.Printf("GetCategoryHierarchy called for category %s (mock implementation)", categoryID)
	return &models.CategoryHierarchyResponse{}, nil
}

// Legacy methods - minimal implementations for interface compatibility
func (r CategoryRepository) AllCategories(ctx context.Context) ([]*models.Category, error) {
	return r.GetPaginatedCategoriesWithSorting(ctx, 1, 1000, "", "")
}

func (r CategoryRepository) FindBySlug(ctx context.Context, slug string) (*models.Category, error) {
	resp, err := r.FindCategoryBySlugName(ctx, slug, "", "")
	if err != nil {
		return nil, err
	}
	return resp.Category, nil
}

// FindCategory implements the missing interface method
func (r CategoryRepository) FindCategory(ctx context.Context, categoryID string) (*models.Category, error) {
	resp, err := r.GetCategoryByID(ctx, categoryID, "", "")
	if err != nil {
		return nil, err
	}
	return resp.Category, nil
}

// GetCategoryFilters implements the missing interface method
func (r CategoryRepository) GetCategoryFilters(ctx context.Context, categoryID string) ([]*models.Filter, error) {
	return r.GetAllFiltersForCategory(ctx, categoryID)
}

// Helper methods for domain conversion
func (r CategoryRepository) categoryToDomain(pb *categoriespb.Category) *models.Category {
	if pb == nil {
		return nil
	}
	return &models.Category{
		ID:               pb.GetId(),
		Description:      pb.GetDescription(),
		ParentID:         pb.GetParentId(),
		GoogleCategoryID: pb.GetGoogleCategoryId(),
		Tags:             pb.GetTags(),
		Slug:             pb.GetSlug(),
		IsActive:         pb.GetIsActive(),
		SeoTitle:         pb.GetSeoTitle(),
		SeoKeywords:      pb.GetSeoKeywords(),
		SeoDesc:          pb.GetSeoDesc(),
		CategoryType:     pb.GetCategoryType(),
		Lang:             pb.GetLang(),
		IsPublic:         true,
		IsFeatured:       false,
	}
}

func (r CategoryRepository) filterToDomain(pb *categoriespb.Filter) *models.Filter {
	if pb == nil {
		return nil
	}
	return &models.Filter{
		ID:         pb.GetId(),
		CategoryID: pb.GetCategoryId(),
		Name:       pb.GetName(),
		FilterType: pb.GetFilterType(),
		Values:     pb.GetValues(),
		IsActive:   pb.GetIsActive(),
		IsRequired: false,
		SortOrder:  0,
	}
}

// dial sets up a gRPC connection with the microservice endpoint
func (r CategoryRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	log.Printf("[CATEGORY_REPO] dial: Attempting to connect to %s", r.endpoint)
	log.Printf("[CATEGORY_REPO] dial: Context state: %v", ctx.Err())

	// Add timeout logging
	deadline, hasDeadline := ctx.Deadline()
	if hasDeadline {
		log.Printf("[CATEGORY_REPO] dial: Context has deadline: %v (in %v)", deadline, time.Until(deadline))
	} else {
		log.Printf("[CATEGORY_REPO] dial: Context has NO deadline")
	}

	log.Printf("[CATEGORY_REPO] dial: About to call rpc.Dial")
	conn, err := rpc.Dial(ctx, r.endpoint)
	log.Printf("[CATEGORY_REPO] dial: rpc.Dial returned - err: %v", err)

	if err != nil {
		log.Printf("[CATEGORY_REPO] dial: ERROR - Failed to connect: %v", err)
		return nil, err
	}
	log.Printf("[CATEGORY_REPO] dial: Successfully connected")
	return conn, nil
}
func (r CategoryRepository) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}
