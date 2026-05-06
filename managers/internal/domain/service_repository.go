package domain

import (
	"context"
	"middleman/managers/internal/models"
	"time"
)

type ServiceRepository interface {
	// Basic CRUD operations
	Add(ctx context.Context,
		id, name, description, serviceType string,
		basePrice int64, pricing []string, availability string,
		providerName, categoryID, categorySlug string,
		descriptionShort, descriptionLong string,
		qualifications []string,
		contact, faq string,
		tags []string,
		status models.Status,
		userType models.UserType,
		shippingCost int64,
		negotiable, hasVariants, middlemanService bool,
		attributes []string,
		options []string,
		thumbnail string,
		lat, long float64) error
	Find(ctx context.Context, serviceID string) (*models.Service, error)
	Remove(ctx context.Context, serviceID, userID string) error
	Update(ctx context.Context, serviceID string, service *models.Service) (string, error)

	// Search operations
	GetServices(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*models.Service, int64, error)
	SearchWithTerm(ctx context.Context, term string) ([]*models.Service, error)
	SuggestServices(ctx context.Context, name string) ([]*models.Service, error)
	SearchServicesWithFilter(
		ctx context.Context,
		categoryID, categorySlug, serviceType string,
		userID string,
		status models.Status,
		searchText string,
		minPrice, maxPrice int64,
		availableFrom, availableTo time.Time,
		hasVariants, negotiable, middlemanService bool,
		userType models.UserType,
		tags, qualifications []string,
		offset int64,
		limit int64,
		lat, lng float64,
		radius int64,
		page, pageSize int64,
		sortBy, sortOrder string,
	) ([]*models.Service, error)
	SearchServicesWithCategorySlug(ctx context.Context, categorySlug string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Service, error)
	SearchServicesWithCategory(ctx context.Context, categoryID string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Service, error)

	// Catalog operations
	GetCatalog(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*models.Service, int64, error)
	GetPublicCatalog(ctx context.Context, userID string, page, pageSize int64, sortBy, sortOrder string) ([]*models.Service, int64, error)

	// Price operations
	UpdateServicePrice(ctx context.Context, serviceID string, newPrice, oldPrice int64) error
	IncreaseServicePrice(ctx context.Context, serviceID string, price int64) (int64, int64, error)    // returns old, new price
	DecreaseServicePrice(ctx context.Context, serviceID string, newPrice int64) (int64, int64, error) // returns old, new price

	// Service management operations
	RebrandService(ctx context.Context, serviceID string, service *models.Service) error
	AdjustServiceStock(ctx context.Context, serviceID string, newStock int64) (int64, int64, error) // returns old, new stock
	ArchiveService(ctx context.Context, serviceID string) (bool, error)
	MarkServiceSold(ctx context.Context, serviceID string) (string, error)                                        // returns status
	MarkServiceLeased(ctx context.Context, serviceID string, monthlyPrice, leaseTermMonths int64) (string, error) // returns status
}
