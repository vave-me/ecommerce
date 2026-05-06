package application

import (
	"context"
	"middleman/search/internal/models"
	"time"
)

type ServiceRepository interface {
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
	Find(ctx context.Context, serviceID string) (*models.Service, error)
	SearchServicesWithCategorySlug(ctx context.Context, categorySlug string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Service, error)
	SearchServicesWithCategory(ctx context.Context, categoryID string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Service, error)
	GetCatalog(ctx context.Context, userId string) ([]*models.Service, error)
}

type ServiceCacheRepository interface {
	ServiceRepository
	Add(ctx context.Context,
		id, name, description, serviceType string,
		basePrice int64, pricing []string, availability string,
		providerName, userID, categoryID, categorySlug string,
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
		lat, long float64, entityType models.EntityType) error
	Remove(ctx context.Context, serviceID string) error
	UpdateService(ctx context.Context,
		id, name, description, serviceType string,
		basePrice int64, pricing []string, availability string,
		providerName, userID, categoryID, categorySlug string,
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
		lat, long float64, entityType models.EntityType) error
	SearchWithTerm(ctx context.Context, term string) ([]*models.Service, error)
	SuggestServices(ctx context.Context, name string) ([]*models.Service, error)
	FindBatch(ctx context.Context, serviceIDs []string) (map[string]*models.Service, error)
}
