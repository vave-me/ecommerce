package domain

import (
	"context"
	"time"
)

type CatalogService struct {
	ID               string
	Name             string
	Description      string
	ServiceType      string
	BasePrice        int64
	Pricing          []string
	Availability     string
	ProviderName     string
	UserID           string
	CategoryID       string
	CategorySlug     string
	DescriptionShort string
	DescriptionLong  string
	Qualifications   []string
	Contact          string
	Faq              string
	Tags             []string
	Status           ServiceStatus
	UserType         UserType
	ShippingCost     int64
	Negotiable       bool
	HasVariants      bool
	MiddlemanService bool
	Attributes       []Attribute
	Options          []Option
	Thumbnail        string
	Lat              float64
	Lng              float64
	Distance         float64 // Distance from search point in meters
}

type CatalogRepository interface {
	AddService(ctx context.Context,
		id, name, description, serviceType string,
		basePrice int64, pricing []string, availability string,
		providerName, userID, categoryID, categorySlug string,
		descriptionShort, descriptionLong string,
		qualifications []string,
		contact, faq string,
		tags []string,
		status ServiceStatus,
		userType UserType,
		shippingCost int64,
		negotiable, hasVariants, middlemanService bool,
		attributes []Attribute,
		options []Option,
		thumbnail string,
		lat, long float64) error
	UpdatePrice(ctx context.Context, serviceID string, oldPrice, newPrice int64) error
	RemoveService(ctx context.Context, serviceID string, userID string) error
	Find(ctx context.Context, serviceID string) (*CatalogService, error)
	GetCatalog(ctx context.Context, userID string, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogService, int64, error)
	GetPublicCatalog(ctx context.Context, userID string, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogService, int64, error)
	GetServices(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogService, int64, error)
	GetServicesByCategory(ctx context.Context, categoryID string, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogService, int64, error)
	GetServicesByCategorySlug(ctx context.Context, categorySlug string, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogService, int64, error)
	RebrandService(ctx context.Context, serviceID, name, description string, tags, qualifications []string, faq string) error
	ArchiveService(ctx context.Context, serviceID string, userID string) error
	MarkServiceSold(ctx context.Context, serviceID string, userID string, finalPrice int64) error
	MarkServiceLeased(ctx context.Context, serviceID string, userID string) error
	ToggleNegotiable(ctx context.Context, serviceID string, userID string, currentValue bool) error
	UpdateService(ctx context.Context,
		id, name, description, serviceType string,
		basePrice int64, pricing []string, availability string,
		providerName, userID, categoryID, categorySlug string,
		descriptionShort, descriptionLong string,
		qualifications []string,
		contact, faq string,
		tags []string,
		status ServiceStatus,
		userType UserType,
		shippingCost int64,
		negotiable, hasVariants, middlemanService bool,
		attributes []Attribute,
		options []Option,
		thumbnail string,
		lat, long float64) error
	FindByLocation(ctx context.Context, lat, lng float64, radiusMeters float64, limit int) ([]*CatalogService, error)
	GetServicesWithFilter(ctx context.Context,
		categoryID, categorySlug, serviceType string,
		userID string,
		status ServiceStatus,
		searchText string,
		minPrice, maxPrice int64,
		lat, lng, radius float64,
		availableFrom, availableTo time.Time,
		hasVariants, negotiable, middlemanService bool,
		userType UserType,
		tags, qualifications []string,
		page, pageSize int64,
		sortBy, sortOrder string) ([]*CatalogService, int64, error)
}
