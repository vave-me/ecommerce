package grpc

import (
	"context"
	"middleman/products/productspb"
	"middleman/wishlists/internal/domain"

	"github.com/stackus/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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

func (r ProductRepository) Find(ctx context.Context, itemID string) (product *domain.Product, err error) {
	var conn *grpc.ClientConn
	conn, err = r.dial(ctx)
	if err != nil {
		return nil, err
	}

	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	resp, err := productspb.NewProductsServiceClient(conn).GetProduct(ctx, &productspb.GetProductRequest{
		Id: itemID,
	})
	if err != nil {
		if errors.GRPCCode(err) == codes.NotFound {
			return nil, errors.ErrNotFound.Msg("product was not located")
		}
		return nil, errors.Wrap(err, "requesting product")
	}

	return r.productToDomain(resp.Product), nil
}

func (r ProductRepository) productToDomain(product *productspb.Product) *domain.Product {
	return &domain.Product{
		ID:           product.GetId(),
		UserSellerID: product.GetUserSellerId(),
		Name:         product.GetName(),
		Price:        product.GetBasePrice(),
	}
}

func (r ProductRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}
