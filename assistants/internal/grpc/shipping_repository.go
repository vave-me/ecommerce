package grpc

import (
	"context"
	"middleman/assistants/internal/domain"
	"middleman/assistants/internal/models"
	"middleman/internal/auth"
	"middleman/internal/rpc"
	"middleman/shipping/shippingpb"

	"google.golang.org/grpc"
)

type ShippingRepository struct {
	endpoint string
	auth     *auth.Auth // Optional auth for authenticated calls
}

var _ domain.ShippingRepository = (*ShippingRepository)(nil)

// NewShippingRepositoryWithAuth creates a new ShippingRepository with JWT authentication support
func NewShippingRepository(endpoint string, authInstance *auth.Auth) ShippingRepository {
	return ShippingRepository{
		endpoint: endpoint,
		auth:     authInstance,
	}
}

// CreateNewShipmentWithDetails creates a new shipping record
func (r ShippingRepository) CreateNewShipmentWithDetails(ctx context.Context, productID, trackingNumber, labelURL, senderName, senderAddress, receiverName, receiverAddress, weight, dimensions, serviceTypes string) (*models.CreateShippingResponse, error) {
	var conn *grpc.ClientConn
	conn, err := r.dial(ctx)

	if err != nil {
		return nil, err
	}

	resp, err := shippingpb.NewShippingServiceClient(conn).CreateShipping(ctx, &shippingpb.CreateShippingRequest{
		ProductId:       productID,
		TrackingNumber:  trackingNumber,
		LabelUrl:        labelURL,
		SenderName:      senderName,
		SenderAddress:   senderAddress,
		ReceiverName:    receiverName,
		ReceiverAddress: receiverAddress,
		Weight:          weight,
		Dimensions:      dimensions,
		ServiceTypes:    serviceTypes,
	})
	if err != nil {
		return nil, err
	}

	return &models.CreateShippingResponse{
		ID: resp.GetId(),
	}, nil
}

// TrackShipmentByTrackingNumber tracks a shipment by tracking number
func (r ShippingRepository) TrackShipmentByTrackingNumber(ctx context.Context, trackingNumber string) (*models.TrackShippingResponse, error) {
	var conn *grpc.ClientConn
	conn, err := r.dial(ctx)

	if err != nil {
		return nil, err
	}

	resp, err := shippingpb.NewShippingServiceClient(conn).TrackShipping(ctx, &shippingpb.TrackShippingRequest{
		TrackingNumber: trackingNumber,
	})
	if err != nil {
		return nil, err
	}

	return &models.TrackShippingResponse{
		Shipping: r.shippingToDomain(resp.GetShipping()),
	}, nil
}

// Domain conversion methods

func (r ShippingRepository) shippingToDomain(shipping *shippingpb.Shipping) *models.Shipping {
	if shipping == nil {
		return nil
	}

	return &models.Shipping{
		ID:              shipping.GetId(),
		ProductID:       shipping.GetProductId(),
		TrackingNumber:  shipping.GetTrackingNumber(),
		LabelURL:        shipping.GetLabelUrl(),
		SenderName:      shipping.GetSenderName(),
		SenderAddress:   shipping.GetSenderAddress(),
		ReceiverName:    shipping.GetReceiverName(),
		ReceiverAddress: shipping.GetReceiverAddress(),
		Weight:          shipping.GetWeight(),
		Dimensions:      shipping.GetDimensions(),
		ServiceTypes:    shipping.GetServiceType(), // Changed from GetServiceTypes to GetServiceType
	}
}

// dial sets up a gRPC connection with the microservice endpoint
func (r ShippingRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

func (r ShippingRepository) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}
