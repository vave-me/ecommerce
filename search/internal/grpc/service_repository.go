// File: search/internal/grpc/service_repository.go
package grpc

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"time"

	"middleman/internal/rpc"
	"middleman/search/internal/application"
	"middleman/search/internal/models"
	"middleman/services/servicespb"
)

// ServiceRepository calls the remote services service (gRPC) as a fallback.
type ServiceRepository struct {
	endpoint string
}

var _ application.ServiceRepository = (*ServiceRepository)(nil)

// NewServiceRepository instantiates the gRPC-based fallback repo.
func NewServiceRepository(endpoint string) ServiceRepository {
	return ServiceRepository{
		endpoint: endpoint,
	}
}

// Find retrieves a service by ID from the services microservice (via gRPC).
func (r ServiceRepository) Find(ctx context.Context, serviceID string) (*models.Service, error) {
	log.Printf("Find: retrieving service with ID=%s via gRPC fallback", serviceID)

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("Find: failed to dial gRPC: %v", err)
		return nil, err
	}
	defer conn.Close()

	client := servicespb.NewServicesServiceClient(conn)
	resp, err := client.GetService(ctx, &servicespb.GetServiceRequest{Id: serviceID})
	if err != nil {
		return nil, fmt.Errorf("GetService RPC failed: %w", err)
	}
	return r.serviceToDomain(resp.GetService()), nil
}
func (r ServiceRepository) GetCatalog(ctx context.Context, userID string) ([]*models.Service, error) {
	log.Printf("Find: retrieving deal with ID=%s via gRPC fallback", userID)

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("Find: failed to dial gRPC: %v", err)
		return nil, err
	}
	defer conn.Close()

	client := servicespb.NewServicesServiceClient(conn)
	resp, err := client.GetCatalog(ctx, &servicespb.GetCatalogRequest{UserId: userID})
	if err != nil {
		return nil, fmt.Errorf("GetDeal RPC failed: %w", err)
	}

	var results []*models.Service
	for _, p := range resp.GetServices() {
		domainProd := r.serviceToDomain(p)
		if domainProd != nil {
			results = append(results, domainProd)
		}
	}
	return results, nil
}

// SearchWithFilters calls the remote microservice's GetServicesWithFilters (or similar).
func (r ServiceRepository) SearchServicesWithFilter(ctx context.Context, categoryID string, categorySlug string, serviceType string, userID string, status models.Status, searchText string, minPrice int64, maxPrice int64, availableFrom time.Time, availableTo time.Time, hasVariants bool, negotiable bool, middlemanService bool, userType models.UserType, tags []string, qualifications []string, offset int64, limit int64, lat float64, lng float64, radius int64, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Service, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := servicespb.NewServicesServiceClient(conn)
	resp, err := client.GetServicesWithFilter(ctx, &servicespb.GetServicesWithFilterRequest{
		CategoryId:     categoryID,
		CategorySlug:   categorySlug,
		ServiceType:    serviceType,
		UserId:         userID,
		SearchText:     searchText,
		Qualifications: qualifications,
		//AvailableFrom:    availableFrom,
		//AvailableTo: availableTo,
		//AvailableFrom: availableFrom,
		//AvailableTo: availableTo,

		MinPrice:         minPrice,
		MaxPrice:         maxPrice,
		Tags:             tags,
		Status:           status.String(),
		Negotiable:       negotiable,
		UserType:         userType.String(),
		MiddlemanService: middlemanService,
		HasVariants:      hasVariants,
		Lat:              float32(lat),
		Lng:              float32(lng),
		Radius:           radius,
		Page:             page,
		PageSize:         pageSize,
		SortBy:           sortBy,
		SortOrder:        sortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("GetServicesWithFilters RPC failed: %w", err)
	}

	var results []*models.Service
	for _, p := range resp.GetServices() {
		domainProd := r.serviceToDomain(p)
		if domainProd != nil {
			results = append(results, domainProd)
		}
	}
	return results, nil
}

// SearchWithFilters calls the remote microservice's GetServicesWithFilters (or similar).
func (r ServiceRepository) SearchServicesWithCategorySlug(
	ctx context.Context,
	categorySlug string,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*models.Service, error) {

	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := servicespb.NewServicesServiceClient(conn)
	resp, err := client.GetServicesByCategorySlug(ctx, &servicespb.GetServicesByCategorySlugRequest{
		CategorySlug: categorySlug,
		Page:         page,
		PageSize:     pageSize,
		SortBy:       sortBy,
		SortOrder:    sortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("GetServicesWithFilters RPC failed: %w", err)
	}

	var results []*models.Service
	for _, p := range resp.GetServices() {
		domainProd := r.serviceToDomain(p)
		if domainProd != nil {
			results = append(results, domainProd)
		}
	}
	return results, nil
}

// SearchWithFilters calls the remote microservice's GetServicesWithFilters (or similar).
func (r ServiceRepository) SearchServicesWithCategory(
	ctx context.Context,
	categoryId string,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*models.Service, error) {

	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := servicespb.NewServicesServiceClient(conn)
	resp, err := client.GetServicesByCategory(ctx, &servicespb.GetServicesByCategoryRequest{
		CategoryId: categoryId,
		Page:       page,
		PageSize:   pageSize,
		SortBy:     sortBy,
		SortOrder:  sortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("GetServicesWithFilters RPC failed: %w", err)
	}

	var results []*models.Service
	for _, p := range resp.GetServices() {
		domainProd := r.serviceToDomain(p)
		if domainProd != nil {
			results = append(results, domainProd)
		}
	}
	return results, nil
}

// Add calls AddService in the remote gRPC service.
func (r ServiceRepository) Add(ctx context.Context, serviceID string, name string, description string, basePrice int64, userSellerID string, categoryID string, brand string, condition string, model string, tags []string) error {
	conn, err := r.dial(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := servicespb.NewServicesServiceClient(conn)
	_, err = client.AddService(ctx, &servicespb.AddServiceRequest{
		Name:        name,
		Description: description,
		BasePrice:   basePrice,
		CategoryId:  categoryID,
		Tags:        tags,
	})
	return err
}

// Update calls UpdateService in the remote gRPC service (or partial, etc.).
func (r ServiceRepository) Update(ctx context.Context, serviceID string, price int64) error {
	conn, err := r.dial(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := servicespb.NewServicesServiceClient(conn)
	_, err = client.UpdateService(ctx, &servicespb.UpdateServiceRequest{
		Id:        serviceID,
		BasePrice: price,
	})
	return err
}

// Remove calls RemoveService in the remote gRPC service.
func (r ServiceRepository) Remove(ctx context.Context, serviceID string) error {
	conn, err := r.dial(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := servicespb.NewServicesServiceClient(conn)
	_, err = client.RemoveService(ctx, &servicespb.RemoveServiceRequest{Id: serviceID})
	return err
}

// serviceToDomain converts a servicespb.Service into our internal models.Service.
func (r ServiceRepository) serviceToDomain(pb *servicespb.Service) *models.Service {
	if pb == nil {
		return nil
	}
	attrs := make([]models.Attribute, len(pb.GetAttributes()))
	for i, a := range pb.GetAttributes() {
		attrs[i] = models.Attribute{
			Key:   a.GetKey(),
			Value: a.GetValue(),
		}
	}
	opts := make([]models.Option, len(pb.GetOptions()))
	for i, o := range pb.GetOptions() {
		opts[i] = models.Option{
			Name:  o.GetName(),
			Value: o.GetValue(),
			Price: float64(o.GetPrice()),
		}
	}

	return &models.Service{
		ID:               pb.GetId(),
		Name:             pb.GetName(),
		Description:      pb.GetDescription(),
		BasePrice:        pb.GetBasePrice(),
		UserID:           pb.GetUserId(),
		CategoryID:       pb.GetCategoryId(),
		CategorySlug:     pb.GetCategorySlug(),
		Tags:             pb.GetTags(),
		Attributes:       attrs,
		Status:           pb.GetStatus(),
		Negotiable:       pb.GetNegotiable(),
		MiddlemanService: pb.GetMiddlemanService(),
		UserType:         pb.GetUserType(),
		ShippingCost:     pb.GetShippingCost(),
		HasVariants:      pb.GetHasVariants(),
		Options:          opts,
		Lat:              float64(pb.GetLat()),
		Lng:              float64(pb.GetLng()),
		Thumbnail:        pb.GetThumbnail(),
	}
}

// dial sets up a gRPC connection with the microservice endpoint.
func (r ServiceRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}
