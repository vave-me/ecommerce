package grpc

import (
	"context"
	"database/sql"
	"middleman/internal/di"
	"middleman/offers/internal/application"
	"middleman/offers/internal/constants"
	"middleman/offers/offerspb"

	"google.golang.org/grpc"
)

type serverTx struct {
	c di.Container
	offerspb.UnimplementedOffersServiceServer
}

var _ offerspb.OffersServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	offerspb.RegisterOffersServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

func (s serverTx) CreateOffer(ctx context.Context, request *offerspb.CreateOfferRequest) (resp *offerspb.CreateOfferResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.CreateOffer(ctx, request)
}
func (s serverTx) StartLease(ctx context.Context, request *offerspb.StartLeaseRequest) (resp *offerspb.StartLeaseResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.StartLease(ctx, request)
}
func (s serverTx) RequestReservationNegotiation(ctx context.Context, request *offerspb.RequestReservationNegotiationRequest) (resp *offerspb.RequestReservationNegotiationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RequestReservationNegotiation(ctx, request)
}
func (s serverTx) RequestBuyNowNegotiation(ctx context.Context, request *offerspb.RequestBuyNowNegotiationRequest) (resp *offerspb.RequestBuyNowNegotiationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RequestBuyNowNegotiation(ctx, request)
}
func (s serverTx) RequestBuyBackNegotiation(ctx context.Context, request *offerspb.RequestBuyBackNegotiationRequest) (resp *offerspb.RequestBuyBackNegotiationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RequestBuyBackNegotiation(ctx, request)
}
func (s serverTx) RequestLeaseNegotiation(ctx context.Context, request *offerspb.RequestLeaseNegotiationRequest) (resp *offerspb.RequestLeaseNegotiationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RequestLeaseNegotiation(ctx, request)
}
func (s serverTx) RedeemReservation(ctx context.Context, request *offerspb.RedeemReservationRequest) (resp *offerspb.RedeemReservationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RedeemReservation(ctx, request)
}
func (s serverTx) RedeemBuyBack(ctx context.Context, request *offerspb.RedeemBuyBackRequest) (resp *offerspb.RedeemBuyBackResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RedeemBuyBack(ctx, request)
}
func (s serverTx) MakeLeasePayment(ctx context.Context, request *offerspb.MakeLeasePaymentRequest) (resp *offerspb.MakeLeasePaymentResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.MakeLeasePayment(ctx, request)
}
func (s serverTx) ExpireReservation(ctx context.Context, request *offerspb.ExpireReservationRequest) (resp *offerspb.ExpireReservationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ExpireReservation(ctx, request)
}
func (s serverTx) ExpireBuyBack(ctx context.Context, request *offerspb.ExpireBuyBackRequest) (resp *offerspb.ExpireBuyBackResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ExpireBuyBack(ctx, request)
}
func (s serverTx) ExecuteLeaseBuyout(ctx context.Context, request *offerspb.ExecuteLeaseBuyoutRequest) (resp *offerspb.ExecuteLeaseBuyoutResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ExecuteLeaseBuyout(ctx, request)
}
func (s serverTx) EndLease(ctx context.Context, request *offerspb.EndLeaseRequest) (resp *offerspb.EndLeaseResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.EndLease(ctx, request)
}
func (s serverTx) DefaultLease(ctx context.Context, request *offerspb.DefaultLeaseRequest) (resp *offerspb.DefaultLeaseResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.DefaultLease(ctx, request)
}
func (s serverTx) DeclineReservationNegotiation(ctx context.Context, request *offerspb.DeclineReservationNegotiationRequest) (resp *offerspb.DeclineReservationNegotiationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.DeclineReservationNegotiation(ctx, request)
}
func (s serverTx) DeclineLeaseNegotiation(ctx context.Context, request *offerspb.DeclineLeaseNegotiationRequest) (resp *offerspb.DeclineLeaseNegotiationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.DeclineLeaseNegotiation(ctx, request)
}
func (s serverTx) DeclineBuyNowNegotiation(ctx context.Context, request *offerspb.DeclineBuyNowNegotiationRequest) (resp *offerspb.DeclineBuyNowNegotiationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.DeclineBuyNowNegotiation(ctx, request)
}
func (s serverTx) DeclineBuyBackNegotiation(ctx context.Context, request *offerspb.DeclineBuyBackNegotiationRequest) (resp *offerspb.DeclineBuyBackNegotiationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.DeclineBuyBackNegotiation(ctx, request)
}

func (s serverTx) CreateReservation(ctx context.Context, request *offerspb.CreateReservationRequest) (resp *offerspb.CreateReservationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.CreateReservation(ctx, request)
}
func (s serverTx) CreateLease(ctx context.Context, request *offerspb.CreateLeaseRequest) (resp *offerspb.CreateLeaseResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.CreateLease(ctx, request)
}
func (s serverTx) CreateBuyNow(ctx context.Context, request *offerspb.CreateBuyNowRequest) (resp *offerspb.CreateBuyNowResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.CreateBuyNow(ctx, request)
}
func (s serverTx) CreateBuyBack(ctx context.Context, request *offerspb.CreateBuyBackRequest) (resp *offerspb.CreateBuyBackResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.CreateBuyBack(ctx, request)
}
func (s serverTx) ConfirmBuyNow(ctx context.Context, request *offerspb.ConfirmBuyNowRequest) (resp *offerspb.ConfirmBuyNowResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ConfirmBuyNow(ctx, request)
}
func (s serverTx) CloseOffer(ctx context.Context, request *offerspb.CloseOfferRequest) (resp *offerspb.CloseOfferResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.CloseOffer(ctx, request)
}
func (s serverTx) CancelReservation(ctx context.Context, request *offerspb.CancelReservationRequest) (resp *offerspb.CancelReservationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.CancelReservation(ctx, request)
}

func (s serverTx) AcceptBuyBackNegotiation(ctx context.Context, request *offerspb.AcceptBuyBackNegotiationRequest) (resp *offerspb.AcceptBuyBackNegotiationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.AcceptBuyBackNegotiation(ctx, request)
}

func (s serverTx) AcceptBuyNowNegotiation(ctx context.Context, request *offerspb.AcceptBuyNowNegotiationRequest) (resp *offerspb.AcceptBuyNowNegotiationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.AcceptBuyNowNegotiation(ctx, request)
}
func (s serverTx) AcceptLeaseNegotiation(ctx context.Context, request *offerspb.AcceptLeaseNegotiationRequest) (resp *offerspb.AcceptLeaseNegotiationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.AcceptLeaseNegotiation(ctx, request)
}

func (s serverTx) AcceptOffer(ctx context.Context, request *offerspb.AcceptOfferRequest) (resp *offerspb.AcceptOfferResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.AcceptOffer(ctx, request)
}

func (s serverTx) AcceptReservationNegotiation(ctx context.Context, request *offerspb.AcceptReservationNegotiationRequest) (resp *offerspb.AcceptReservationNegotiationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.AcceptReservationNegotiation(ctx, request)
}

func (s serverTx) ActivateOffer(ctx context.Context, request *offerspb.ActivateOfferRequest) (resp *offerspb.ActivateOfferResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ActivateOffer(ctx, request)
}

func (s serverTx) CancelBuyBack(ctx context.Context, request *offerspb.CancelBuyBackRequest) (resp *offerspb.CancelBuyBackResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.CancelBuyBack(ctx, request)
}
func (s serverTx) CancelBuyNow(ctx context.Context, request *offerspb.CancelBuyNowRequest) (resp *offerspb.CancelBuyNowResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.CancelBuyNow(ctx, request)
}
func (s serverTx) CancelLease(ctx context.Context, request *offerspb.CancelLeaseRequest) (resp *offerspb.CancelLeaseResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.CancelLease(ctx, request)
}

func (s serverTx) GetOffer(ctx context.Context, request *offerspb.GetOfferRequest) (resp *offerspb.GetOfferResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetOffer(ctx, request)
}

func (s serverTx) ListOffers(ctx context.Context, request *offerspb.ListOffersRequest) (resp *offerspb.ListOffersResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ListOffers(ctx, request)
}

func (s serverTx) closeTx(tx *sql.Tx, err error) error {
	if p := recover(); p != nil {
		_ = tx.Rollback()
		panic(p)
	} else if err != nil {
		_ = tx.Rollback()
		return err
	} else {
		return tx.Commit()
	}
}
