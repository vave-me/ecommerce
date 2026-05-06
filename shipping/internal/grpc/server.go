// grpc/shipping_service_server.go
package grpc

import (
	"context"
	"fmt"
	"log"
	"time"
	"middleman/shipping/internal/dhl"
	"middleman/internal/di"
	"middleman/internal/errorsotel"
	"middleman/shipping/internal/application"
	"middleman/shipping/internal/application/commands"
	"middleman/shipping/internal/application/queries"
	"middleman/shipping/internal/constants"
	"middleman/shipping/internal/domain"
	"middleman/shipping/shippingpb"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// server implements shippingpb.ShippingServiceServer
type server struct {
	app application.App

	shippingpb.UnimplementedShippingServiceServer
}

// Ensure server implements shippingpb.ShippingServiceServer
var _ shippingpb.ShippingServiceServer = (*server)(nil)

// RegisterServer registers the ShippingServiceServer with the provided gRPC registrar
func RegisterServer(ctx context.Context, app application.App, registrar grpc.ServiceRegistrar) error {
	shippingpb.RegisterShippingServiceServer(registrar, &server{
		app: app,
	})
	log.Println("Shipping service server registered successfully")
	return nil
}

// CreateShipping handles the creation of a new Shipping entity
func (s *server) CreateShipping(ctx context.Context, request *shippingpb.CreateShippingRequest) (*shippingpb.CreateShippingResponse, error) {
	span := trace.SpanFromContext(ctx)

	shippingID := uuid.New().String()

	span.SetAttributes(
		attribute.String("ShippingID", shippingID),
	)

	client := di.Get(ctx, constants.DHLClient).(*dhl.Client)

	shipmentData := map[string]interface{}{
		"sender": map[string]interface{}{
			"name":    request.SenderName,
			"address": request.SenderAddress,
		},
		"receiver": map[string]interface{}{
			"name":    request.ReceiverName,
			"address": request.ReceiverAddress,
		},
		"package": map[string]interface{}{
			"weight": request.Weight,
			"dimensions": map[string]float64{
				"length": 12.00,
				"width":  12.00,
				"height": 12.00,
			},
		},
		"serviceType": request.ServiceTypes,
	}
	response, err := client.CreateShipment(ctx, shipmentData)
	if err != nil {
		return nil, err
	}
	err = s.app.CreateShipping(ctx, commands.CreateShipping{
		ID:              shippingID,
		OrderID:         request.OrderId,
		BasketID:        request.BasketId,
		TrackingNumber:  response.ShipmentTrackingNumber,
		LabelUrl:        getLabelURL(response),
		ProductID:       request.ProductId,
		SenderName:      request.SenderName,
		SenderAddress:   request.SenderAddress,
		ReceiverName:    request.ReceiverName,
		ReceiverAddress: request.ReceiverAddress,
		Weight:          request.Weight,
		Dimensions:      request.Dimensions,
		ServiceType:     request.ServiceTypes,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &shippingpb.CreateShippingResponse{
		Id: shippingID,
	}, nil
}

// TrackShipping handles tracking request for a shipment
func (s *server) TrackShipping(ctx context.Context, request *shippingpb.TrackShippingRequest) (*shippingpb.TrackShippingResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("TrackingNumber", request.TrackingNumber),
	)
	
	// First try to get shipment from our database
	shipment, err := s.app.TrackShipment(ctx, queries.TrackShipment{
		TrackingNumber: request.TrackingNumber,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	
	// Optionally, also get real-time tracking from DHL
	client := di.Get(ctx, constants.DHLClient).(*dhl.Client)
	trackingResp, err := client.TrackShipment(ctx, request.TrackingNumber)
	if err == nil && len(trackingResp.Shipments) > 0 {
		// Update status if DHL has newer information
		// This could trigger an update command if needed
	}
	
	return &shippingpb.TrackShippingResponse{
		Shipping: toProtoShipping(shipment),
	}, nil
}

// CancelShipment cancels an existing shipment
func (s *server) CancelShipment(ctx context.Context, request *shippingpb.CancelShipmentRequest) (*shippingpb.CancelShipmentResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("ShipmentID", request.Id),
		attribute.String("Reason", request.Reason),
	)
	
	err := s.app.CancelShipment(ctx, commands.CancelShipment{
		ID:     request.Id,
		Reason: request.Reason,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	
	return &shippingpb.CancelShipmentResponse{
		Success: true,
	}, nil
}

// SchedulePickup schedules a pickup for the shipment
func (s *server) SchedulePickup(ctx context.Context, request *shippingpb.SchedulePickupRequest) (*shippingpb.SchedulePickupResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("ShipmentID", request.Id),
	)
	
	err := s.app.SchedulePickup(ctx, commands.SchedulePickup{
		ID:           request.Id,
		PickupTime:   request.PickupTime.AsTime().Format(time.RFC3339),
		Instructions: request.Instructions,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	
	// TODO: In production, integrate with DHL pickup API
	confirmationNumber := fmt.Sprintf("PU-%s-%d", request.Id[:8], time.Now().Unix())
	
	return &shippingpb.SchedulePickupResponse{
		PickupConfirmationNumber: confirmationNumber,
	}, nil
}

// StartShipment starts the shipment journey
func (s *server) StartShipment(ctx context.Context, request *shippingpb.StartShipmentRequest) (*shippingpb.StartShipmentResponse, error) {
	err := s.app.StartShipment(ctx, commands.StartShipment{
		ID: request.Id,
	})
	if err != nil {
		return nil, err
	}
	
	return &shippingpb.StartShipmentResponse{
		Success: true,
	}, nil
}

// UpdateShipmentStatus updates the shipment status
func (s *server) UpdateShipmentStatus(ctx context.Context, request *shippingpb.UpdateShipmentStatusRequest) (*shippingpb.UpdateShipmentStatusResponse, error) {
	err := s.app.UpdateShipmentStatus(ctx, commands.UpdateShipmentStatus{
		ID:       request.Id,
		Status:   request.Status,
		Location: request.Location,
		Notes:    request.Notes,
	})
	if err != nil {
		return nil, err
	}
	
	return &shippingpb.UpdateShipmentStatusResponse{
		Success: true,
	}, nil
}

// AssignCarrier assigns a carrier to the shipment
func (s *server) AssignCarrier(ctx context.Context, request *shippingpb.AssignCarrierRequest) (*shippingpb.AssignCarrierResponse, error) {
	err := s.app.AssignCarrier(ctx, commands.AssignCarrier{
		ID:          request.Id,
		CarrierID:   request.CarrierId,
		CarrierName: request.CarrierName,
	})
	if err != nil {
		return nil, err
	}
	
	return &shippingpb.AssignCarrierResponse{
		Success: true,
	}, nil
}

// MarkShipmentAsDelivered marks the shipment as delivered
func (s *server) MarkShipmentAsDelivered(ctx context.Context, request *shippingpb.MarkShipmentAsDeliveredRequest) (*shippingpb.MarkShipmentAsDeliveredResponse, error) {
	err := s.app.MarkShipmentAsDelivered(ctx, commands.MarkShipmentAsDelivered{
		ID:                 request.Id,
		SignedBy:           request.SignedBy,
		ProofOfDeliveryURL: request.ProofOfDeliveryUrl,
	})
	if err != nil {
		return nil, err
	}
	
	return &shippingpb.MarkShipmentAsDeliveredResponse{
		Success: true,
	}, nil
}

// ReturnShipment initiates a return for the shipment
func (s *server) ReturnShipment(ctx context.Context, request *shippingpb.ReturnShipmentRequest) (*shippingpb.ReturnShipmentResponse, error) {
	err := s.app.ReturnShipment(ctx, commands.ReturnShipment{
		ID:                   request.Id,
		Reason:               request.Reason,
		ReturnTrackingNumber: request.ReturnTrackingNumber,
	})
	if err != nil {
		return nil, err
	}
	
	// TODO: In production, create actual return shipment
	returnShipmentID := fmt.Sprintf("RET-%s", request.Id)
	
	return &shippingpb.ReturnShipmentResponse{
		ReturnShipmentId: returnShipmentID,
	}, nil
}

// GetShipment retrieves a shipment by ID
func (s *server) GetShipment(ctx context.Context, request *shippingpb.GetShipmentRequest) (*shippingpb.GetShipmentResponse, error) {
	shipment, err := s.app.GetShipment(ctx, queries.GetShipment{
		ID: request.Id,
	})
	if err != nil {
		return nil, err
	}
	
	return &shippingpb.GetShipmentResponse{
		Shipping: toProtoShipping(shipment),
	}, nil
}

// ListShipments lists shipments with filters
func (s *server) ListShipments(ctx context.Context, request *shippingpb.ListShipmentsRequest) (*shippingpb.ListShipmentsResponse, error) {
	query := queries.ListShipments{
		Limit:     int(request.Limit),
		Offset:    int(request.Offset),
		Status:    request.Status,
		ProductID: request.ProductId,
		OrderID:   request.OrderId,
		CarrierID: request.CarrierId,
	}
	
	if request.StartDate != nil {
		query.StartDate = request.StartDate.AsTime().Format(time.RFC3339)
	}
	if request.EndDate != nil {
		query.EndDate = request.EndDate.AsTime().Format(time.RFC3339)
	}
	
	shipments, err := s.app.ListShipments(ctx, query)
	if err != nil {
		return nil, err
	}
	
	protoShipments := make([]*shippingpb.Shipping, len(shipments))
	for i, shipment := range shipments {
		protoShipments[i] = toProtoShipping(shipment)
	}
	
	return &shippingpb.ListShipmentsResponse{
		Shipments: protoShipments,
		Total:     int32(len(shipments)), // TODO: Get total count from DB
	}, nil
}

// GetShipmentHistory retrieves the history of a shipment
func (s *server) GetShipmentHistory(ctx context.Context, request *shippingpb.GetShipmentHistoryRequest) (*shippingpb.GetShipmentHistoryResponse, error) {
	events, err := s.app.GetShipmentHistory(ctx, queries.GetShipmentHistory{
		ID: request.Id,
	})
	if err != nil {
		return nil, err
	}
	
	protoEvents := make([]*shippingpb.ShipmentEvent, len(events))
	for i, event := range events {
		timestamp, _ := time.Parse(time.RFC3339, event.Timestamp)
		protoEvents[i] = &shippingpb.ShipmentEvent{
			Id:          event.ID,
			ShipmentId:  event.ShipmentID,
			EventType:   event.EventType,
			Status:      event.Status,
			Location:    event.Location,
			Description: event.Description,
			Timestamp:   timestamppb.New(timestamp),
		}
	}
	
	return &shippingpb.GetShipmentHistoryResponse{
		Events: protoEvents,
	}, nil
}

// GetLabel retrieves the shipping label
func (s *server) GetLabel(ctx context.Context, request *shippingpb.GetLabelRequest) (*shippingpb.GetLabelResponse, error) {
	// TODO: In production, retrieve from storage or DHL API
	// For now, return mock data
	mockPDF := []byte("%PDF-1.4 Mock Label Content")
	
	contentType := "application/pdf"
	if request.Format == "png" {
		contentType = "image/png"
		mockPDF = []byte("PNG Mock Label Content")
	} else if request.Format == "zpl" {
		contentType = "text/plain"
		mockPDF = []byte("^XA^FO50,50^ADN,36,20^FDMock Label^FS^XZ")
	}
	
	return &shippingpb.GetLabelResponse{
		LabelData:   mockPDF,
		ContentType: contentType,
	}, nil
}

// GetRates calculates shipping rates
func (s *server) GetRates(ctx context.Context, request *shippingpb.GetRatesRequest) (*shippingpb.GetRatesResponse, error) {
	client := di.Get(ctx, constants.DHLClient).(*dhl.Client)
	
	rateReq := dhl.RateRequest{
		CustomerDetails: dhl.CustomerDetails{
			ShipperDetails: dhl.ShipperReceiverDetails{
				PostalAddress: dhl.PostalAddress{
					PostalCode:  request.SenderPostalCode,
					CountryCode: request.SenderCountryCode,
				},
			},
			ReceiverDetails: dhl.ShipperReceiverDetails{
				PostalAddress: dhl.PostalAddress{
					PostalCode:  request.ReceiverPostalCode,
					CountryCode: request.ReceiverCountryCode,
				},
			},
		},
		PlannedShippingDateAndTime: time.Now().Format("2006-01-02T15:04:05 MST"),
		UnitOfMeasurement:          "metric",
		IsCustomsDeclarable:        false,
		Packages: []dhl.Package{{
			Weight: float64(request.Weight),
			Dimensions: dhl.Dimensions{
				Length: float64(request.Length),
				Width:  float64(request.Width),
				Height: float64(request.Height),
			},
		}},
	}
	
	rateResp, err := client.GetRates(ctx, rateReq)
	if err != nil {
		return nil, err
	}
	
	var rates []*shippingpb.Rate
	for _, product := range rateResp.Products {
		for _, price := range product.TotalPrice {
			deliveryTime := time.Now().Add(time.Duration(product.DeliveryCapabilities.TotalTransitDays) * 24 * time.Hour)
			rates = append(rates, &shippingpb.Rate{
				ServiceType:       product.ProductCode,
				ServiceName:       product.ProductName,
				Amount:            float32(price.Price),
				Currency:          price.PriceCurrency,
				EstimatedDays:     int32(product.DeliveryCapabilities.TotalTransitDays),
				EstimatedDelivery: timestamppb.New(deliveryTime),
			})
		}
	}
	
	return &shippingpb.GetRatesResponse{
		Rates: rates,
	}, nil
}

// Helper functions
func getLabelURL(response *dhl.CreateShipmentResponse) string {
	if response.Documents != nil && len(response.Documents) > 0 {
		for _, doc := range response.Documents {
			if doc.TypeCode == "label" {
				// In production, you might want to store the PDF content and return a URL
				return fmt.Sprintf("/labels/%s.pdf", response.ShipmentTrackingNumber)
			}
		}
	}
	return ""
}

func toProtoShipping(shipment *domain.CatalogShipment) *shippingpb.Shipping {
	shipping := &shippingpb.Shipping{
		Id:              shipment.ID,
		ProductId:       shipment.ProductID,
		OrderId:         shipment.OrderID,
		BasketId:        shipment.BasketID,
		TrackingNumber:  shipment.TrackingNumber,
		LabelUrl:        shipment.LabelURL,
		SenderName:      shipment.SenderName,
		SenderAddress:   shipment.SenderAddress,
		ReceiverName:    shipment.ReceiverName,
		ReceiverAddress: shipment.ReceiverAddress,
		Weight:          shipment.Weight,
		Dimensions:      shipment.Dimensions,
		ServiceType:     string(shipment.ServiceType),
		Status:          string(shipment.Status),
		CarrierId:       shipment.CarrierID,
		CarrierName:     shipment.CarrierName,
		CreatedAt:       timestamppb.New(shipment.CreatedAt),
		UpdatedAt:       timestamppb.New(shipment.UpdatedAt),
	}
	
	if shipment.PickupScheduledAt != nil {
		shipping.PickupScheduledAt = timestamppb.New(*shipment.PickupScheduledAt)
	}
	if shipment.DeliveredAt != nil {
		shipping.DeliveredAt = timestamppb.New(*shipment.DeliveredAt)
	}
	if shipment.CancelledAt != nil {
		shipping.CancelledAt = timestamppb.New(*shipment.CancelledAt)
	}
	
	return shipping
}
