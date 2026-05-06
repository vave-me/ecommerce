package domain

import (
	"context"
	"middleman/assistants/internal/models"
	"time"
)

type ServiceRepository interface {
	// Basic CRUD operations
	CreateService(ctx context.Context,
		name, description, serviceType string,
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
	GetServiceByID(ctx context.Context, serviceID string) (*models.Service, error)
	DeleteService(ctx context.Context, serviceID, userID string) error
	UpdateServiceDetails(ctx context.Context, serviceID string, service *models.Service) (string, error)

	// Search operations
	GetAllServices(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*models.Service, int64, error)
	SearchServicesByName(ctx context.Context, name string) ([]*models.Service, error)
	GetServiceSuggestions(ctx context.Context, name string) ([]*models.Service, error)
	SearchServicesAdvanced(
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
	GetServicesByCategorySlug(ctx context.Context, categorySlug string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Service, error)
	GetServicesByCategoryID(ctx context.Context, categoryID string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Service, error)

	// Catalog operations
	GetServiceCatalog(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*models.Service, int64, error)
	GetPublicServiceCatalog(ctx context.Context, userID string, page, pageSize int64, sortBy, sortOrder string) ([]*models.Service, int64, error)

	// Price operations
	UpdateServicePriceForUser(ctx context.Context, serviceID string, newPrice, oldPrice int64) error
	IncreaseServicePriceBy(ctx context.Context, serviceID string, increaseAmount int64) (int64, int64, error) // returns old, new price
	DecreaseServicePriceTo(ctx context.Context, serviceID string, newPrice int64) (int64, int64, error)       // returns old, new price

	// Service management operations
	UpdateServiceBranding(ctx context.Context, serviceID string, service *models.Service) error
	AdjustServiceInventory(ctx context.Context, serviceID string, newStock int64) (int64, int64, error) // returns old, new stock
	ArchiveUserService(ctx context.Context, serviceID string) (bool, error)
	MarkServiceAsSold(ctx context.Context, serviceID string) (string, error)                                        // returns status
	MarkServiceAsLeased(ctx context.Context, serviceID string, monthlyPrice, leaseTermMonths int64) (string, error) // returns status
}
