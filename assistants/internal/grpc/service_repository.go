// File: search/internal/grpc/service_repository.go
package grpc

import (
	"context"
	"fmt"
	"middleman/assistants/internal/domain"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"middleman/assistants/internal/models"
	"middleman/internal/auth"
	"middleman/internal/rpc"
	"middleman/services/servicespb"
)

// ServiceRepository calls the remote services service (gRPC) as a fallback.
type ServiceRepository struct {
	endpoint string
	auth     *auth.Auth // Optional auth for authenticated calls
}

var _ domain.ServiceRepository = (*ServiceRepository)(nil)

// NewServiceRepository instantiates the gRPC-based fallback repo.

// NewServiceRepositoryWithAuth creates a new ServiceRepository with JWT authentication support
func NewServiceRepository(endpoint string, authInstance *auth.Auth) ServiceRepository {
	return ServiceRepository{
		endpoint: endpoint,
		auth:     authInstance,
	}
}

func (r ServiceRepository) SearchServicesByName(ctx context.Context, name string) ([]*models.Service, error) {
	log.Printf("SearchServicesByName: searching services with name=%s via gRPC", name)

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("SearchServicesByName: failed to dial gRPC: %v", err)
		return nil, err
	}
	defer conn.Close()

	client := servicespb.NewServicesServiceClient(conn)

	// Use GetServicesWithFilter with SearchText filter
	resp, err := client.GetServicesWithFilter(ctx, &servicespb.GetServicesWithFilterRequest{
		SearchText: name,
		Page:       1,
		PageSize:   20,
	})
	if err != nil {
		return nil, fmt.Errorf("GetServicesWithFilter RPC failed: %w", err)
	}

	services := make([]*models.Service, 0, len(resp.GetServices()))
	for _, pbService := range resp.GetServices() {
		services = append(services, r.serviceToDomain(pbService))
	}

	return services, nil
}

func (r ServiceRepository) GetServiceSuggestions(ctx context.Context, name string) ([]*models.Service, error) {
	// Use SearchServicesByName as a fallback for suggestions
	return r.SearchServicesByName(ctx, name)
}

// GetServiceByID retrieves a service by ID from the services microservice (via gRPC).
func (r ServiceRepository) GetServiceByID(ctx context.Context, serviceID string) (*models.Service, error) {
	log.Printf("GetServiceByID: retrieving service with ID=%s via gRPC", serviceID)

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("GetServiceByID: failed to dial gRPC: %v", err)
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

// SearchServicesAdvanced calls the remote microservice's GetServicesWithFilters.
func (r ServiceRepository) SearchServicesAdvanced(ctx context.Context, categoryID string, categorySlug string, serviceType string, userID string, status models.Status, searchText string, minPrice int64, maxPrice int64, availableFrom time.Time, availableTo time.Time, hasVariants bool, negotiable bool, middlemanService bool, userType models.UserType, tags []string, qualifications []string, offset int64, limit int64, lat float64, lng float64, radius int64, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Service, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := servicespb.NewServicesServiceClient(conn)
	resp, err := client.GetServicesWithFilter(ctx, &servicespb.GetServicesWithFilterRequest{
		CategoryId:     categoryID,
		CategorySlug:   categorySlug,
		ServiceType:    serviceType,
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

// GetServicesByCategorySlug calls the remote microservice to get services by category slug.
func (r ServiceRepository) GetServicesByCategorySlug(
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

// GetServicesByCategoryID calls the remote microservice to get services by category ID.
func (r ServiceRepository) GetServicesByCategoryID(
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

// CreateService calls AddService in the remote gRPC service.
func (r ServiceRepository) CreateService(ctx context.Context, name string, description string, serviceType string, basePrice int64, pricing []string, availability string, providerName string, categoryID string, categorySlug string, descriptionShort string, descriptionLong string, qualifications []string, contact string, faq string, tags []string, status models.Status, userType models.UserType, shippingCost int64, negotiable bool, hasVariants bool, middlemanService bool, attributes []string, options []string, thumbnail string, lat float64, long float64) error {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Convert string arrays to proper protobuf types
	var pbAttributes []*servicespb.Attribute
	for i := 0; i < len(attributes); i += 2 {
		if i+1 < len(attributes) {
			pbAttributes = append(pbAttributes, &servicespb.Attribute{
				Key:   attributes[i],
				Value: attributes[i+1],
			})
		}
	}

	var pbOptions []*servicespb.Option
	for i := 0; i < len(options); i += 3 {
		if i+2 < len(options) {
			price := int64(0)
			if val, err := strconv.ParseInt(options[i+2], 10, 64); err == nil {
				price = val
			}
			pbOptions = append(pbOptions, &servicespb.Option{
				Name:  options[i],
				Value: options[i+1],
				Price: price,
			})
		}
	}

	client := servicespb.NewServicesServiceClient(conn)
	_, err = client.AddService(ctx, &servicespb.AddServiceRequest{
		Name:             name,
		Description:      description,
		BasePrice:        basePrice,
		ServiceType:      serviceType,
		CategoryId:       categoryID,
		Pricing:          pricing,
		Availability:     availability,
		ProviderName:     providerName,
		CategorySlug:     categorySlug,
		DescriptionShort: descriptionShort,
		DescriptionLong:  descriptionLong,
		Qualifications:   qualifications,
		Contact:          contact,
		Faq:              faq,
		Tags:             tags,
		Status:           status.String(),
		UserType:         userType.String(),
		ShippingCost:     shippingCost,
		Negotiable:       negotiable,
		HasVariants:      hasVariants,
		MiddlemanService: middlemanService,
		Attributes:       pbAttributes,
		Options:          pbOptions,
		Thumbnail:        thumbnail,
		Lat:              float32(lat),
		Lng:              float32(long),
	})
	return err
}

// UpdateServicePriceOnly calls UpdateService with only price change (legacy method for backward compatibility)
func (r ServiceRepository) UpdateServicePriceOnly(ctx context.Context, serviceID string, price int64) error {
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

// DeleteService calls RemoveService in the remote gRPC service.
func (r ServiceRepository) DeleteService(ctx context.Context, serviceID, userID string) error {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := servicespb.NewServicesServiceClient(conn)
	_, err = client.RemoveService(ctx, &servicespb.RemoveServiceRequest{
		Id: serviceID,
	})
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
// dial sets up a gRPC connection with the microservice endpoint
func (r ServiceRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

func (r ServiceRepository) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}

// GetAllServices retrieves services with pagination
func (r ServiceRepository) GetAllServices(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*models.Service, int64, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer conn.Close()

	client := servicespb.NewServicesServiceClient(conn)
	resp, err := client.GetServices(ctx, &servicespb.GetServicesRequest{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("GetServices RPC failed: %w", err)
	}

	services := make([]*models.Service, 0, len(resp.GetServices()))
	for _, pbService := range resp.GetServices() {
		services = append(services, r.serviceToDomain(pbService))
	}

	return services, resp.GetTotalCount(), nil
}

// UpdateServiceDetails updates a service
func (r ServiceRepository) UpdateServiceDetails(ctx context.Context, serviceID string, service *models.Service) (string, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := servicespb.NewServicesServiceClient(conn)

	// Convert attributes and options to protobuf format
	var pbAttributes []*servicespb.Attribute
	var pbOptions []*servicespb.Option

	resp, err := client.UpdateService(ctx, &servicespb.UpdateServiceRequest{
		Id:               serviceID,
		Name:             service.Name,
		Description:      service.Description,
		ServiceType:      service.ServiceType,
		BasePrice:        service.BasePrice,
		Pricing:          service.Pricing,
		Availability:     service.Availability,
		ProviderName:     service.ProviderName,
		CategoryId:       service.CategoryID,
		CategorySlug:     service.CategorySlug,
		DescriptionShort: service.DescriptionShort,
		DescriptionLong:  service.DescriptionLong,
		Qualifications:   service.Qualifications,
		Contact:          service.Contact,
		Faq:              service.Faq,
		Tags:             service.Tags,
		Status:           service.Status,
		UserType:         service.UserType,
		ShippingCost:     service.ShippingCost,
		HasVariants:      service.HasVariants,
		MiddlemanService: service.MiddlemanService,
		Negotiable:       service.Negotiable,
		Attributes:       pbAttributes,
		Options:          pbOptions,
		Thumbnail:        service.Thumbnail,
		Lat:              float32(service.Lat),
		Lng:              float32(service.Lng),
	})
	if err != nil {
		return "", fmt.Errorf("UpdateService RPC failed: %w", err)
	}

	return resp.GetId(), nil
}

// GetServiceCatalog retrieves user's service catalog
func (r ServiceRepository) GetServiceCatalog(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*models.Service, int64, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer conn.Close()

	client := servicespb.NewServicesServiceClient(conn)
	resp, err := client.GetCatalog(ctx, &servicespb.GetCatalogRequest{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("GetCatalog RPC failed: %w", err)
	}

	services := make([]*models.Service, 0, len(resp.GetServices()))
	for _, pbService := range resp.GetServices() {
		services = append(services, r.serviceToDomain(pbService))
	}

	return services, resp.GetTotalCount(), nil
}

// GetPublicServiceCatalog retrieves user's public service catalog
func (r ServiceRepository) GetPublicServiceCatalog(ctx context.Context, userID string, page, pageSize int64, sortBy, sortOrder string) ([]*models.Service, int64, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer conn.Close()

	client := servicespb.NewServicesServiceClient(conn)
	resp, err := client.GetPublicCatalog(ctx, &servicespb.GetPublicCatalogRequest{
		UserId:    userID,
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("GetPublicCatalog RPC failed: %w", err)
	}

	services := make([]*models.Service, 0, len(resp.GetServices()))
	for _, pbService := range resp.GetServices() {
		services = append(services, r.serviceToDomain(pbService))
	}

	return services, resp.GetTotalCount(), nil
}

// UpdateServicePriceForUser updates service price
func (r ServiceRepository) UpdateServicePriceForUser(ctx context.Context, serviceID string, newPrice, oldPrice int64) error {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := servicespb.NewServicesServiceClient(conn)
	_, err = client.UpdateServicePrice(ctx, &servicespb.UpdateServicePriceRequest{

		Id:       serviceID,
		NewPrice: newPrice,
		OldPrice: oldPrice,
	})
	if err != nil {
		return fmt.Errorf("UpdateServicePrice RPC failed: %w", err)
	}

	return nil
}

// IncreaseServicePriceBy increases service price
func (r ServiceRepository) IncreaseServicePriceBy(ctx context.Context, serviceID string, increaseAmount int64) (int64, int64, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return 0, 0, err
	}
	defer conn.Close()

	client := servicespb.NewServicesServiceClient(conn)
	resp, err := client.IncreaseServicePrice(ctx, &servicespb.IncreaseServicePriceRequest{
		ServiceId: serviceID,
		Price:     increaseAmount,
	})
	if err != nil {
		return 0, 0, fmt.Errorf("IncreaseServicePrice RPC failed: %w", err)
	}

	return resp.GetOldPrice(), resp.GetNewPrice(), nil
}

// DecreaseServicePriceTo decreases service price
func (r ServiceRepository) DecreaseServicePriceTo(ctx context.Context, serviceID string, newPrice int64) (int64, int64, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return 0, 0, err
	}
	defer conn.Close()

	client := servicespb.NewServicesServiceClient(conn)
	resp, err := client.DecreaseServicePrice(ctx, &servicespb.DecreaseServicePriceRequest{

		ServiceId: serviceID,
		NewPrice:  newPrice,
	})
	if err != nil {
		return 0, 0, fmt.Errorf("DecreaseServicePrice RPC failed: %w", err)
	}

	return resp.GetOldPrice(), resp.GetNewPrice(), nil
}

// UpdateServiceBranding rebrands a service
func (r ServiceRepository) UpdateServiceBranding(ctx context.Context, serviceID string, service *models.Service) error {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := servicespb.NewServicesServiceClient(conn)

	// Convert attributes and options to protobuf format
	var pbAttributes []*servicespb.Attribute
	var pbOptions []*servicespb.Option

	_, err = client.RebrandService(ctx, &servicespb.RebrandServiceRequest{
		Id:               serviceID,
		Name:             service.Name,
		Description:      service.Description,
		BasePrice:        service.BasePrice,
		Pricing:          service.Pricing,
		Availability:     service.Availability,
		ProviderName:     service.ProviderName,
		CategoryId:       service.CategoryID,
		CategorySlug:     service.CategorySlug,
		DescriptionShort: service.DescriptionShort,
		DescriptionLong:  service.DescriptionLong,
		Qualifications:   service.Qualifications,
		Contact:          service.Contact,
		Faq:              service.Faq,
		Tags:             service.Tags,
		Status:           service.Status,
		UserType:         service.UserType,
		ShippingCost:     service.ShippingCost,
		HasVariants:      service.HasVariants,
		MiddlemanService: service.MiddlemanService,
		Negotiable:       service.Negotiable,
		Attributes:       pbAttributes,
		Options:          pbOptions,
		Thumbnail:        service.Thumbnail,
		Lat:              float32(service.Lat),
		Lng:              float32(service.Lng),
	})
	if err != nil {
		return fmt.Errorf("RebrandService RPC failed: %w", err)
	}

	return nil
}

// AdjustServiceInventory adjusts service stock
func (r ServiceRepository) AdjustServiceInventory(ctx context.Context, serviceID string, newStock int64) (int64, int64, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return 0, 0, err
	}
	defer conn.Close()

	client := servicespb.NewServicesServiceClient(conn)
	resp, err := client.AdjustServiceStock(ctx, &servicespb.AdjustServiceStockRequest{

		ServiceId: serviceID,
		NewStock:  newStock,
	})
	if err != nil {
		return 0, 0, fmt.Errorf("AdjustServiceStock RPC failed: %w", err)
	}

	return resp.GetOldStock(), resp.GetNewStock(), nil
}

// ArchiveUserService archives a service
func (r ServiceRepository) ArchiveUserService(ctx context.Context, serviceID string) (bool, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	client := servicespb.NewServicesServiceClient(conn)
	resp, err := client.ArchiveService(ctx, &servicespb.ArchiveServiceRequest{

		ServiceId: serviceID,
	})
	if err != nil {
		return false, fmt.Errorf("ArchiveService RPC failed: %w", err)
	}

	return resp.GetArchived(), nil
}

// MarkServiceAsSold marks a service as sold
func (r ServiceRepository) MarkServiceAsSold(ctx context.Context, serviceID string) (string, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := servicespb.NewServicesServiceClient(conn)
	resp, err := client.MarkServiceSold(ctx, &servicespb.MarkServiceSoldRequest{

		ServiceId: serviceID,
	})
	if err != nil {
		return "", fmt.Errorf("MarkServiceSold RPC failed: %w", err)
	}

	return resp.GetStatus(), nil
}

// MarkServiceAsLeased marks a service as leased
func (r ServiceRepository) MarkServiceAsLeased(ctx context.Context, serviceID string, monthlyPrice, leaseTermMonths int64) (string, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := servicespb.NewServicesServiceClient(conn)
	resp, err := client.MarkServiceLeased(ctx, &servicespb.MarkServiceLeasedRequest{

		ServiceId:       serviceID,
		MonthlyPrice:    monthlyPrice,
		LeaseTermMonths: leaseTermMonths,
	})
	if err != nil {
		return "", fmt.Errorf("MarkServiceLeased RPC failed: %w", err)
	}

	return resp.GetStatus(), nil
}
