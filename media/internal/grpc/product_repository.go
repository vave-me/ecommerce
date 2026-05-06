package grpc

import (
	"context"
	"middleman/products/productspb"

	"github.com/stackus/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"middleman/internal/rpc"
)

type ProductRepository struct {
	endpoint string
}

func NewProductRepository(endpoint string) ProductRepository {
	return ProductRepository{
		endpoint: endpoint,
	}
}

// ValidateSKUs checks if the provided SKUs exist in the products service
// Returns a map of SKU to product ID for valid SKUs and a list of invalid SKUs
func (r ProductRepository) ValidateSKUs(ctx context.Context, skus []string) (map[string]string, []string, error) {
	if len(skus) == 0 {
		return make(map[string]string), []string{}, nil
	}

	conn, err := r.dial(ctx)
	if err != nil {
		return nil, nil, err
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
	resp, err := productspb.NewProductsServiceClient(conn).GetProductsBySKUs(outCtx, &productspb.GetProductsBySKUsRequest{
		Skus: skus,
	})
	if err != nil {
		if errors.GRPCCode(err) == codes.Unauthenticated {
			return nil, nil, errors.ErrUnauthorized.Msg("unauthorized access to product service")
		}
		return nil, nil, errors.Wrap(err, "requesting products by SKUs")
	}

	// Build SKU to product ID map
	skuToProductID := make(map[string]string)
	for _, product := range resp.Products {
		skuToProductID[product.GetSku()] = product.GetId()
	}

	return skuToProductID, resp.NotFoundSkus, nil
}

func (r ProductRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}