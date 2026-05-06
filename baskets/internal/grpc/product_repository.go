package grpc

import (
	"context"
	"middleman/products/productspb"

	"github.com/stackus/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"middleman/baskets/internal/domain"
	"middleman/internal/rpc"
)

type ProductRepository struct {
	endpoint string
}

var _ domain.ProductRepository = (*ProductRepository)(nil)

func NewProductRepository(endpoint string) ProductRepository {
	return ProductRepository{
		endpoint: endpoint,
	}
}

func (r ProductRepository) Find(ctx context.Context, productID string) (*domain.Product, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)
	// Extract metadata from the incoming context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	// Create a new outgoing context with the extracted metadata
	outCtx := metadata.NewOutgoingContext(ctx, md)

	// Make the gRPC call using the outgoing context
	resp, err := productspb.NewProductsServiceClient(conn).GetProduct(outCtx, &productspb.GetProductRequest{
		Id: productID,
	})
	if err != nil {
		if errors.GRPCCode(err) == codes.NotFound {
			return nil, errors.ErrNotFound.Msg("product was not located")
		}
		if errors.GRPCCode(err) == codes.Unauthenticated {
			return nil, errors.ErrUnauthorized.Msg("unauthorized access to product service")
		}
		return nil, errors.Wrap(err, "requesting product")
	}

	return r.productToDomain(resp.Product), nil
}

// productToDomain maps a productspb.Product to our domain.Product struct.
func (r ProductRepository) productToDomain(product *productspb.Product) *domain.Product {
	// 1. Convert productspb.Attribute => domain.Attribute
	domainAttrs := make([]domain.Attribute, len(product.Attributes))
	for i, pbAttr := range product.Attributes {
		domainAttrs[i] = domain.Attribute{
			Key:   pbAttr.Key,
			Value: pbAttr.Value,
		}
	}

	// 2. Convert productspb.Option => domain.Option
	domainOpts := make([]domain.Option, len(product.Options))
	for i, pbOpt := range product.Options {
		domainOpts[i] = domain.Option{
			Name:  pbOpt.Name,
			Value: pbOpt.Value,
			// If domain.Option uses float64, cast from int64
			Price: float64(pbOpt.Price),
		}
	}

	// 3. Return a domain.Product with all fields mapped
	return &domain.Product{
		ID:               product.GetId(),
		Name:             product.GetName(),
		Description:      product.GetDescription(),
		BasePrice:        product.GetBasePrice(),
		UserSellerID:     product.GetUserSellerId(),
		CategoryID:       product.GetCategoryId(),
		Brand:            product.GetBrand(),
		Condition:        domain.ProductCondition(product.GetCondition()),
		Model:            product.GetModel(),
		Tags:             product.GetTags(),
		ManageStock:      product.GetManageStocks(),
		Stock:            product.GetStock(),
		SKU:              product.GetSku(),
		Attributes:       domainAttrs,
		Weight:           product.GetWeight(),
		Height:           product.GetHeight(),
		Width:            product.GetWidth(),
		Depth:            product.GetDepth(),
		Status:           domain.ProductStatus(product.GetStatus()),
		Negotiable:       product.GetNegotiable(),
		UserType:         domain.UserType(product.GetUserType()),
		MiddlemanService: product.GetMiddlemanService(),
		ShippingCost:     product.GetShippingCost(),
		HasVariants:      product.GetHasVariants(),
		Options:          domainOpts,
		Thumbnail:        product.GetThumbnail(),
		Lat:              float64(product.GetLat()),
		Lng:              float64(product.GetLng()),
	}
}

func (r ProductRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}
