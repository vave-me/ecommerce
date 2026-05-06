package grpc

import (
	"context"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"

	"middleman/internal/errorsotel"
	"middleman/offers/internal/application"
	"middleman/offers/internal/application/commands"
	"middleman/offers/internal/application/queries"
	"middleman/offers/offerspb"
)

// server implements offerspb.OffersServiceServer
type server struct {
	app application.App
	offerspb.UnimplementedOffersServiceServer
}

// Ensure compile-time check
var _ offerspb.OffersServiceServer = (*server)(nil)

// RegisterServer is called to register this service with the gRPC server
func RegisterServer(
	ctx context.Context,
	app application.App,
	registrar grpc.ServiceRegistrar,
) error {
	offerspb.RegisterOffersServiceServer(registrar, server{app: app})
	return nil
}

// 1.1 CreateOffer
func (s server) CreateOffer(ctx context.Context, request *offerspb.CreateOfferRequest) (*offerspb.CreateOfferResponse, error) {

	span := trace.SpanFromContext(ctx)
	offerID := uuid.New().String()
	span.SetAttributes(attribute.String("offerID", offerID))

	err := s.app.CreateOffer(ctx, commands.CreateOffer{
		ID:           offerID,
		UserSellerID: request.GetUserSellerId(),
		ProductID:    request.GetProductId(),
		Price:        request.GetPrice(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.CreateOfferResponse{
		Id: offerID,
		// created_at: ... (if needed)
	}, nil
}

// 1.6 ActivateOffer
func (s server) ActivateOffer(ctx context.Context, req *offerspb.ActivateOfferRequest) (*offerspb.ActivateOfferResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("offerID", req.GetOfferId()))

	err := s.app.ActivateOffer(ctx, commands.ActivateOffer{ID: req.GetOfferId()})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.ActivateOfferResponse{
		OfferId:     req.GetOfferId(),
		OfferStatus: "active",
	}, nil
}

func (s server) CloseOffer(ctx context.Context, req *offerspb.CloseOfferRequest) (*offerspb.CloseOfferResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("offerID", req.GetOfferId()))

	err := s.app.CloseOffer(ctx, commands.CloseOffer{
		ID:     req.GetOfferId(),
		Reason: req.GetReason(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.CloseOfferResponse{
		OfferId:     req.GetOfferId(),
		OfferStatus: "closed",
	}, nil
}

func (s server) AcceptOffer(ctx context.Context, req *offerspb.AcceptOfferRequest) (*offerspb.AcceptOfferResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("offerID", req.GetOfferId()))

	err := s.app.AcceptOffer(ctx, commands.AcceptOffer{
		ID:             req.GetOfferId(),
		UserCustomerID: req.GetUserCustomerId(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.AcceptOfferResponse{
		OfferId:     req.GetOfferId(),
		OfferStatus: "accepted",
	}, nil
}

// GetOffer retrieves a single offer by ID
func (s server) GetOffer(ctx context.Context, req *offerspb.GetOfferRequest) (*offerspb.GetOfferResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("offerID", req.GetOfferId()))

	offer, err := s.app.GetOffer(ctx, queries.GetOffer{
		OfferID: req.GetOfferId(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.GetOfferResponse{
		Offer: &offerspb.Offer{
			Id:             offer.ID(),
			UserSellerId:   offer.UserSellerID,
			UserCustomerId: offer.UserCustomerID,
			ProductId:      offer.ProductID,
			Price:          offer.Price,
			OfferStatus:    string(offer.Status),
		},
	}, nil
}

// ListOffers returns a list of offers with optional filtering
func (s server) ListOffers(ctx context.Context, req *offerspb.ListOffersRequest) (*offerspb.ListOffersResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("userSellerId", req.GetUserSellerId()),
		attribute.String("userCustomerId", req.GetUserCustomerId()),
		attribute.String("offerStatus", req.GetOfferStatus()),
	)

	offers, total, err := s.app.ListOffers(ctx, queries.ListOffers{
		UserSellerID:   req.GetUserSellerId(),
		UserCustomerID: req.GetUserCustomerId(),
		OfferStatus:    req.GetOfferStatus(),
		Limit:          int(req.GetLimit()),
		Offset:         int(req.GetPage() * req.GetLimit()), // Convert page to offset
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Convert domain offers to protobuf offers
	pbOffers := make([]*offerspb.Offer, 0, len(offers))
	for _, offer := range offers {
		pbOffers = append(pbOffers, &offerspb.Offer{
			Id:             offer.ID(),
			UserSellerId:   offer.UserSellerID,
			UserCustomerId: offer.UserCustomerID,
			ProductId:      offer.ProductID,
			Price:          offer.Price,
			OfferStatus:    string(offer.Status),
		})
	}

	return &offerspb.ListOffersResponse{
		Offers: pbOffers,
		Total:  total,
		Page:   req.GetPage(),
		Limit:  req.GetLimit(),
	}, nil
}

// -----------------------------------------------------------
// 2) BuyNow aggregator methods
// -----------------------------------------------------------

func (s server) CreateBuyNow(ctx context.Context, req *offerspb.CreateBuyNowRequest) (*offerspb.CreateBuyNowResponse, error) {

	span := trace.SpanFromContext(ctx)
	buyNowID := uuid.New().String()
	span.SetAttributes(attribute.String("buyNowID", buyNowID))

	err := s.app.CreateBuyNow(ctx, commands.CreateBuyNow{
		ID:         buyNowID,
		OfferID:    req.GetOfferId(),
		FinalPrice: req.GetFinalPrice(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.CreateBuyNowResponse{
		BuyNowId: buyNowID,
		// created_at ...
	}, nil
}

func (s server) ConfirmBuyNow(ctx context.Context, req *offerspb.ConfirmBuyNowRequest) (*offerspb.ConfirmBuyNowResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("buyNowID", req.GetBuyNowId()))

	err := s.app.ConfirmBuyNow(ctx, commands.ConfirmBuyNow{
		ID: req.GetBuyNowId(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.ConfirmBuyNowResponse{
		BuyNowId: req.GetBuyNowId(),
		Status:   "confirmed",
		// confirmed_at ...
	}, nil
}

// BuyNow Negotiation
func (s server) RequestBuyNowNegotiation(ctx context.Context, req *offerspb.RequestBuyNowNegotiationRequest) (*offerspb.RequestBuyNowNegotiationResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("buyNowID", req.GetBuyNowId()))

	err := s.app.RequestBuyNowNegotiation(ctx, commands.RequestBuyNowNegotiation{
		ID: req.GetBuyNowId(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.RequestBuyNowNegotiationResponse{
		BuyNowId:          req.GetBuyNowId(),
		NegotiationStatus: "requested",
	}, nil
}

func (s server) AcceptBuyNowNegotiation(ctx context.Context, req *offerspb.AcceptBuyNowNegotiationRequest) (*offerspb.AcceptBuyNowNegotiationResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("buyNowID", req.GetBuyNowId()))

	err := s.app.AcceptBuyNowNegotiation(ctx, commands.AcceptBuyNowNegotiation{
		ID:              req.GetBuyNowId(),
		NegotiatedPrice: req.GetFinalPrice(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.AcceptBuyNowNegotiationResponse{
		BuyNowId:          req.GetBuyNowId(),
		NegotiationStatus: "accepted",
	}, nil
}

func (s server) DeclineBuyNowNegotiation(ctx context.Context, req *offerspb.DeclineBuyNowNegotiationRequest) (*offerspb.DeclineBuyNowNegotiationResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("buyNowID", req.GetBuyNowId()))

	err := s.app.DeclineBuyNowNegotiation(ctx, commands.DeclineBuyNowNegotiation{
		ID:     req.GetBuyNowId(),
		Reason: req.GetReason(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.DeclineBuyNowNegotiationResponse{
		BuyNowId:          req.GetBuyNowId(),
		NegotiationStatus: "declined",
	}, nil
}

func (s server) CreateLease(ctx context.Context, req *offerspb.CreateLeaseRequest) (*offerspb.CreateLeaseResponse, error) {

	span := trace.SpanFromContext(ctx)
	leaseID := uuid.New().String()
	span.SetAttributes(attribute.String("leaseID", leaseID))

	err := s.app.CreateLease(ctx, commands.CreateLease{
		ID:              leaseID,
		OfferID:         req.GetOfferId(),
		MonthlyPrice:    req.GetMonthlyPrice(),
		LeaseTermMonths: req.GetLeaseTermMonths(),
		HasBuyout:       req.GetHasBuyout(),
		BuyoutPrice:     req.GetBuyoutPrice(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.CreateLeaseResponse{
		LeaseId:     leaseID,
		LeaseStatus: "pending",
		// created_at ...
	}, nil
}

func (s server) MakeLeasePayment(ctx context.Context, req *offerspb.MakeLeasePaymentRequest) (*offerspb.MakeLeasePaymentResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("leaseID", req.GetLeaseId()))

	err := s.app.MakeLeasePayment(ctx, commands.MakeLeasePayment{
		ID:     req.GetLeaseId(),
		Amount: req.GetAmount(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.MakeLeasePaymentResponse{
		LeaseId:     req.GetLeaseId(),
		PaymentDate: nil, // or actual date
	}, nil
}

func (s server) ExecuteLeaseBuyout(ctx context.Context, req *offerspb.ExecuteLeaseBuyoutRequest) (*offerspb.ExecuteLeaseBuyoutResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("leaseID", req.GetLeaseId()))

	err := s.app.ExecuteLeaseBuyout(ctx, commands.ExecuteLeaseBuyout{
		ID: req.GetLeaseId(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.ExecuteLeaseBuyoutResponse{
		LeaseId:     req.GetLeaseId(),
		LeaseStatus: "completed",
		// executed_at ...
	}, nil
}

func (s server) EndLease(ctx context.Context, req *offerspb.EndLeaseRequest) (*offerspb.EndLeaseResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("leaseID", req.GetLeaseId()))

	err := s.app.EndLease(ctx, commands.EndLease{
		ID: req.GetLeaseId(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.EndLeaseResponse{
		LeaseId:     req.GetLeaseId(),
		LeaseStatus: "completed",
		// ended_at ...
	}, nil
}

func (s server) DefaultLease(ctx context.Context, req *offerspb.DefaultLeaseRequest) (*offerspb.DefaultLeaseResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("leaseID", req.GetLeaseId()))

	err := s.app.DefaultLease(ctx, commands.DefaultLease{
		ID:     req.GetLeaseId(),
		Reason: req.GetReason(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.DefaultLeaseResponse{
		LeaseId:     req.GetLeaseId(),
		LeaseStatus: "defaulted",
	}, nil
}

// Lease Negotiation
func (s server) RequestLeaseNegotiation(ctx context.Context, req *offerspb.RequestLeaseNegotiationRequest) (*offerspb.RequestLeaseNegotiationResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("leaseID", req.GetLeaseId()))

	// Placeholder domain command
	err := s.app.RequestLeaseNegotiation(ctx, commands.RequestLeaseNegotiation{
		Message: req.GetMessage(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.RequestLeaseNegotiationResponse{
		LeaseId:           req.GetLeaseId(),
		NegotiationStatus: "requested",
	}, nil
}

func (s server) AcceptLeaseNegotiation(ctx context.Context, req *offerspb.AcceptLeaseNegotiationRequest) (*offerspb.AcceptLeaseNegotiationResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("leaseID", req.GetLeaseId()))

	// Placeholder domain command
	err := s.app.AcceptLeaseNegotiation(ctx, commands.AcceptLeaseNegotiation{})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.AcceptLeaseNegotiationResponse{
		LeaseId:           req.GetLeaseId(),
		NegotiationStatus: "accepted",
	}, nil
}

func (s server) DeclineLeaseNegotiation(ctx context.Context, req *offerspb.DeclineLeaseNegotiationRequest) (*offerspb.DeclineLeaseNegotiationResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("leaseID", req.GetLeaseId()))

	err := s.app.DeclineLeaseNegotiation(ctx, commands.DeclineLeaseNegotiation{
		ID:     req.GetLeaseId(),
		Reason: req.GetReason(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.DeclineLeaseNegotiationResponse{
		LeaseId:           req.GetLeaseId(),
		NegotiationStatus: "declined",
	}, nil
}

func (s server) CreateBuyBack(
	ctx context.Context,
	req *offerspb.CreateBuyBackRequest,
) (*offerspb.CreateBuyBackResponse, error) {

	span := trace.SpanFromContext(ctx)
	buyBackID := uuid.New().String()
	span.SetAttributes(attribute.String("buyBackID", buyBackID))

	err := s.app.CreateBuyBack(ctx, commands.CreateBuyBack{
		BuyBackID:        buyBackID,
		OfferID:          req.GetOfferId(),
		LockedPrice:      req.GetLockedPrice(),
		RedemptionFee:    req.GetRedemptionFee(),
		LockDurationDays: int(req.GetLockDurationDays()),
		LockBuyerID:      req.GetLockBuyerId(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.CreateBuyBackResponse{
		BuyBackId:     buyBackID,
		BuyBackStatus: "active",
		CreatedAt:     nil,
	}, nil
}

func (s server) RedeemBuyBack(ctx context.Context, req *offerspb.RedeemBuyBackRequest) (*offerspb.RedeemBuyBackResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("buyBackID", req.GetBuyBackId()))

	err := s.app.RedeemBuyBack(ctx, commands.RedeemBuyBack{
		BuyBackID: req.GetBuyBackId(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.RedeemBuyBackResponse{
		BuyBackId:     req.GetBuyBackId(),
		BuyBackStatus: "redeemed",
	}, nil
}

func (s server) ExpireBuyBack(
	ctx context.Context,
	req *offerspb.ExpireBuyBackRequest,
) (*offerspb.ExpireBuyBackResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("buyBackID", req.GetBuyBackId()))

	err := s.app.ExpireBuyBack(ctx, commands.ExpireBuyBack{})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.ExpireBuyBackResponse{
		BuyBackId:     req.GetBuyBackId(),
		BuyBackStatus: "expired",
	}, nil
}

func (s server) CancelBuyBack(ctx context.Context, req *offerspb.CancelBuyBackRequest) (*offerspb.CancelBuyBackResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("buyBackID", req.GetBuyBackId()))

	err := s.app.CancelBuyBack(ctx, commands.CancelBuyBack{})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.CancelBuyBackResponse{
		BuyBackId:     req.GetBuyBackId(),
		BuyBackStatus: "canceled",
		// canceled_at ...
	}, nil
}

func (s server) RequestBuyBackNegotiation(ctx context.Context, req *offerspb.RequestBuyBackNegotiationRequest) (*offerspb.RequestBuyBackNegotiationResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("buyBackID", req.GetBuyBackId()))

	err := s.app.RequestBuyBackNegotiation(ctx, commands.RequestBuyBackNegotiation{
		BuyBackID: req.GetBuyBackId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.RequestBuyBackNegotiationResponse{
		BuyBackId:         req.GetBuyBackId(),
		NegotiationStatus: "requested",
	}, nil
}

func (s server) AcceptBuyBackNegotiation(ctx context.Context, req *offerspb.AcceptBuyBackNegotiationRequest) (*offerspb.AcceptBuyBackNegotiationResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("buyBackID", req.GetBuyBackId()))

	err := s.app.AcceptBuyBackNegotiation(ctx, commands.AcceptBuyBackNegotiation{
		BuyBackID: req.GetBuyBackId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.AcceptBuyBackNegotiationResponse{
		BuyBackId:         req.GetBuyBackId(),
		NegotiationStatus: "accepted",
	}, nil
}

func (s server) CreateReservation(
	ctx context.Context,
	req *offerspb.CreateReservationRequest,
) (*offerspb.CreateReservationResponse, error) {

	span := trace.SpanFromContext(ctx)
	reservationID := uuid.New().String()
	span.SetAttributes(attribute.String("reservationID", reservationID))

	err := s.app.CreateReservation(ctx, commands.CreateReservation{
		ReservationID:    reservationID,
		OfferID:          req.GetOfferId(),
		LockedPrice:      req.GetLockedPrice(),
		RedemptionFee:    req.GetReservationFee(),
		LockDurationDays: int(req.GetLockDurationDays()),
		LockBuyerID:      req.GetLockBuyerId(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.CreateReservationResponse{
		ReservationId:     reservationID,
		ReservationStatus: "active",
		CreatedAt:         nil,
		IsFree:            req.GetLockDurationDays() == 1 && req.GetReservationFee() == 0,
	}, nil
}

func (s server) RedeemReservation(ctx context.Context, req *offerspb.RedeemReservationRequest) (*offerspb.RedeemReservationResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("reservationID", req.GetReservationId()))

	err := s.app.RedeemReservation(ctx, commands.RedeemReservation{
		ReservationID: req.GetReservationId(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.RedeemReservationResponse{
		ReservationId:     req.GetReservationId(),
		ReservationStatus: "redeemed",
	}, nil
}

func (s server) ExpireReservation(
	ctx context.Context,
	req *offerspb.ExpireReservationRequest,
) (*offerspb.ExpireReservationResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("reservationID", req.GetReservationId()))

	err := s.app.ExpireReservation(ctx, commands.ExpireReservation{})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.ExpireReservationResponse{
		ReservationId:     req.GetReservationId(),
		ReservationStatus: "expired",
	}, nil
}

func (s server) CancelReservation(ctx context.Context, req *offerspb.CancelReservationRequest) (*offerspb.CancelReservationResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("reservationID", req.GetReservationId()))

	err := s.app.CancelReservation(ctx, commands.CancelReservation{})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.CancelReservationResponse{
		ReservationId:     req.GetReservationId(),
		ReservationStatus: "canceled",
		// canceled_at ...
	}, nil
}

func (s server) RequestReservationNegotiation(ctx context.Context, req *offerspb.RequestReservationNegotiationRequest) (*offerspb.RequestReservationNegotiationResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("reservationID", req.GetReservationId()))

	err := s.app.RequestReservationNegotiation(ctx, commands.RequestReservationNegotiation{
		ReservationID: req.GetReservationId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.RequestReservationNegotiationResponse{
		ReservationId:     req.GetReservationId(),
		NegotiationStatus: "requested",
	}, nil
}

func (s server) AcceptReservationNegotiation(ctx context.Context, req *offerspb.AcceptReservationNegotiationRequest) (*offerspb.AcceptReservationNegotiationResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("reservationID", req.GetReservationId()))

	err := s.app.AcceptReservationNegotiation(ctx, commands.AcceptReservationNegotiation{
		ReservationID: req.GetReservationId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.AcceptReservationNegotiationResponse{
		ReservationId:     req.GetReservationId(),
		NegotiationStatus: "accepted",
	}, nil
}

// Missing methods implementation

// CancelBuyNow cancels a buy now offer
func (s server) CancelBuyNow(ctx context.Context, req *offerspb.CancelBuyNowRequest) (*offerspb.CancelBuyNowResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("buyNowID", req.GetBuyNowId()))

	err := s.app.CancelBuyNow(ctx, commands.CancelBuyNow{
		ID: req.GetBuyNowId(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.CancelBuyNowResponse{
		BuyNowId:     req.GetBuyNowId(),
		BuyNowStatus: "canceled",
	}, nil
}

// StartLease starts a lease
func (s server) StartLease(ctx context.Context, req *offerspb.StartLeaseRequest) (*offerspb.StartLeaseResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("leaseID", req.GetLeaseId()))

	err := s.app.StartLease(ctx, commands.StartLease{
		ID: req.GetLeaseId(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.StartLeaseResponse{
		LeaseId:     req.GetLeaseId(),
		LeaseStatus: "active",
	}, nil
}

// CancelLease cancels a lease
func (s server) CancelLease(ctx context.Context, req *offerspb.CancelLeaseRequest) (*offerspb.CancelLeaseResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("leaseID", req.GetLeaseId()))

	err := s.app.CancelLease(ctx, commands.CancelLease{
		LeaseID: req.GetLeaseId(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.CancelLeaseResponse{
		LeaseId:     req.GetLeaseId(),
		LeaseStatus: "canceled",
	}, nil
}

// DeclineBuyBackNegotiation declines a buy back negotiation
func (s server) DeclineBuyBackNegotiation(ctx context.Context, req *offerspb.DeclineBuyBackNegotiationRequest) (*offerspb.DeclineBuyBackNegotiationResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("buyBackID", req.GetBuyBackId()))

	err := s.app.DeclineBuyBackNegotiation(ctx, commands.DeclineBuyBackNegotiation{
		BuyBackID: req.GetBuyBackId(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.DeclineBuyBackNegotiationResponse{
		BuyBackId:         req.GetBuyBackId(),
		NegotiationStatus: "declined",
	}, nil
}

// DeclineReservationNegotiation declines a reservation negotiation
func (s server) DeclineReservationNegotiation(ctx context.Context, req *offerspb.DeclineReservationNegotiationRequest) (*offerspb.DeclineReservationNegotiationResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("reservationID", req.GetReservationId()))

	err := s.app.DeclineReservationNegotiation(ctx, commands.DeclineReservationNegotiation{
		ReservationID: req.GetReservationId(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &offerspb.DeclineReservationNegotiationResponse{
		ReservationId:     req.GetReservationId(),
		NegotiationStatus: "declined",
	}, nil
}
