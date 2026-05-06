package application

import (
	"context"
	"middleman/internal/ddd"
	"middleman/offers/internal/application/commands"
	"middleman/offers/internal/application/queries"
	"middleman/offers/internal/domain"
)

type (
	// App aggregates all possible domain Commands in one interface.
	App interface {
		Commands
		Queries
	}

	// Queries enumerates all possible read operations.
	Queries interface {
		GetOffer(ctx context.Context, query queries.GetOffer) (*domain.Offer, error)
		ListOffers(ctx context.Context, query queries.ListOffers) ([]*domain.Offer, int64, error)
	}

	// Commands enumerates all possible domain actions (including negotiations).
	Commands interface {
		// Existing Offer/Lease/BuyBack
		CreateOffer(ctx context.Context, cmd commands.CreateOffer) error
		ActivateOffer(ctx context.Context, cmd commands.ActivateOffer) error
		CloseOffer(ctx context.Context, cmd commands.CloseOffer) error
		AcceptOffer(ctx context.Context, cmd commands.AcceptOffer) error

		//BuyNow
		CreateBuyNow(ctx context.Context, cmd commands.CreateBuyNow) error
		ConfirmBuyNow(ctx context.Context, cmd commands.ConfirmBuyNow) error
		CancelBuyNow(ctx context.Context, cmd commands.CancelBuyNow) error
		RequestBuyNowNegotiation(ctx context.Context, cmd commands.RequestBuyNowNegotiation) error
		AcceptBuyNowNegotiation(ctx context.Context, cmd commands.AcceptBuyNowNegotiation) error
		DeclineBuyNowNegotiation(ctx context.Context, cmd commands.DeclineBuyNowNegotiation) error

		//Lease
		CreateLease(ctx context.Context, cmd commands.CreateLease) error
		StartLease(ctx context.Context, cmd commands.StartLease) error
		MakeLeasePayment(ctx context.Context, cmd commands.MakeLeasePayment) error
		ExecuteLeaseBuyout(ctx context.Context, cmd commands.ExecuteLeaseBuyout) error
		EndLease(ctx context.Context, cmd commands.EndLease) error
		CancelLease(ctx context.Context, cmd commands.CancelLease) error
		DefaultLease(ctx context.Context, cmd commands.DefaultLease) error
		RequestLeaseNegotiation(ctx context.Context, cmd commands.RequestLeaseNegotiation) error
		AcceptLeaseNegotiation(ctx context.Context, cmd commands.AcceptLeaseNegotiation) error
		DeclineLeaseNegotiation(ctx context.Context, cmd commands.DeclineLeaseNegotiation) error

		//BuyBack
		CreateBuyBack(ctx context.Context, cmd commands.CreateBuyBack) error
		RedeemBuyBack(ctx context.Context, cmd commands.RedeemBuyBack) error
		ExpireBuyBack(ctx context.Context, cmd commands.ExpireBuyBack) error
		CancelBuyBack(ctx context.Context, cmd commands.CancelBuyBack) error
		RequestBuyBackNegotiation(ctx context.Context, cmd commands.RequestBuyBackNegotiation) error
		AcceptBuyBackNegotiation(ctx context.Context, cmd commands.AcceptBuyBackNegotiation) error
		DeclineBuyBackNegotiation(ctx context.Context, cmd commands.DeclineBuyBackNegotiation) error

		//Reservation
		CreateReservation(ctx context.Context, cmd commands.CreateReservation) error
		RedeemReservation(ctx context.Context, cmd commands.RedeemReservation) error
		ExpireReservation(ctx context.Context, cmd commands.ExpireReservation) error
		CancelReservation(ctx context.Context, cmd commands.CancelReservation) error
		RequestReservationNegotiation(ctx context.Context, cmd commands.RequestReservationNegotiation) error
		AcceptReservationNegotiation(ctx context.Context, cmd commands.AcceptReservationNegotiation) error
		DeclineReservationNegotiation(ctx context.Context, cmd commands.DeclineReservationNegotiation) error
	}

	// Application is the concrete struct embedding all command handlers.
	Application struct {
		appCommands
		appQueries
	}

	// appCommands aggregates all the individual command handlers.
	appCommands struct {
		//Offer
		commands.CreateOfferHandler
		commands.ActivateOfferHandler
		commands.CloseOfferHandler
		commands.AcceptOfferHandler

		//BuyNow
		commands.CreateBuyNowHandler
		commands.ConfirmBuyNowHandler
		commands.CancelBuyNowHandler
		commands.RequestBuyNowNegotiationHandler
		commands.AcceptBuyNowNegotiationHandler
		commands.DeclineBuyNowNegotiationHandler

		//Lease
		commands.CreateLeaseHandler
		commands.StartLeaseHandler
		commands.MakeLeasePaymentHandler
		commands.ExecuteLeaseBuyoutHandler
		commands.EndLeaseHandler
		commands.CancelLeaseHandler
		commands.DefaultLeaseHandler
		commands.RequestLeaseNegotiationHandler
		commands.AcceptLeaseNegotiationHandler
		commands.DeclineLeaseNegotiationHandler

		//BuyBack
		commands.CreateBuyBackHandler
		commands.RedeemBuyBackHandler
		commands.ExpireBuyBackHandler
		commands.CancelBuyBackHandler
		commands.RequestBuyBackNegotiationHandler
		commands.AcceptBuyBackNegotiationHandler
		commands.DeclineBuyBackNegotiationHandler

		//reservation
		commands.CreateReservationHandler
		commands.RedeemReservationHandler
		commands.ExpireReservationHandler
		commands.CancelReservationHandler
		commands.RequestReservationNegotiationHandler
		commands.AcceptReservationNegotiationHandler
		commands.DeclineReservationNegotiationHandler
	}

	appQueries struct {
		queries.GetOfferHandler
		queries.ListOffersHandler
	}
)

// Compile-time check ensures Application satisfies App interface.
var _ App = (*Application)(nil)

// New instantiates the Application with all required repositories & publishers.
func New(
	offers domain.OfferRepository,
	leasing domain.LeaseRepository,
	buyBacks domain.BuyBackRepository,
	buynow domain.BuyNowRepository,
	reservations domain.ReservationRepository,
	middleman domain.MiddlemanRepository, // For read queries
	publisher ddd.EventPublisher[ddd.Event],
) *Application {

	return &Application{
		appCommands: appCommands{
			//Offers
			CreateOfferHandler:   commands.NewCreateOfferHandler(offers, publisher),
			ActivateOfferHandler: commands.NewActivateOfferHandler(offers, publisher),
			CloseOfferHandler:    commands.NewCloseOfferHandler(offers, publisher),
			AcceptOfferHandler:   commands.NewAcceptOfferHandler(offers, publisher),

			//BuyNow
			CreateBuyNowHandler:             commands.NewCreateBuyNowHandler(buynow, publisher),
			ConfirmBuyNowHandler:            commands.NewConfirmBuyNowHandler(buynow, publisher),
			CancelBuyNowHandler:             commands.NewCancelBuyNowHandler(buynow, publisher),
			RequestBuyNowNegotiationHandler: commands.NewRequestBuyNowNegotiationHandler(buynow, publisher),
			AcceptBuyNowNegotiationHandler:  commands.NewAcceptBuyNowNegotiationHandler(buynow, publisher),
			DeclineBuyNowNegotiationHandler: commands.NewDeclineBuyNowNegotiationHandler(buynow, publisher),

			//Lease
			CreateLeaseHandler:             commands.NewCreateLeaseHandler(leasing, publisher),
			StartLeaseHandler:              commands.NewStartLeaseHandler(leasing, publisher),
			MakeLeasePaymentHandler:        commands.NewMakeLeasePaymentHandler(leasing, publisher),
			ExecuteLeaseBuyoutHandler:      commands.NewExecuteLeaseBuyoutHandler(leasing, publisher),
			EndLeaseHandler:                commands.NewEndLeaseHandler(leasing, publisher),
			CancelLeaseHandler:             commands.NewCancelLeaseHandler(leasing, publisher),
			DefaultLeaseHandler:            commands.NewDefaultLeaseHandler(leasing, publisher),
			RequestLeaseNegotiationHandler: commands.NewRequestLeaseNegotiationHandler(leasing, publisher),
			AcceptLeaseNegotiationHandler:  commands.NewAcceptLeaseNegotiationHandler(leasing, publisher),
			DeclineLeaseNegotiationHandler: commands.NewDeclineLeaseNegotiationHandler(leasing, publisher),

			//BuyBack
			CreateBuyBackHandler:             commands.NewCreateBuyBackHandler(buyBacks, publisher),
			RedeemBuyBackHandler:             commands.NewRedeemBuyBackHandler(buyBacks, publisher),
			ExpireBuyBackHandler:             commands.NewExpireBuyBackHandler(buyBacks, publisher),
			CancelBuyBackHandler:             commands.NewCancelBuyBackHandler(buyBacks, publisher),
			RequestBuyBackNegotiationHandler: commands.NewRequestBuyBackNegotiationHandler(buyBacks, publisher),
			AcceptBuyBackNegotiationHandler:  commands.NewAcceptBuyBackNegotiationHandler(buyBacks, publisher),
			DeclineBuyBackNegotiationHandler: commands.NewDeclineBuyBackNegotiationHandler(buyBacks, publisher),

			//Reservation
			CreateReservationHandler:             commands.NewCreateReservationHandler(reservations, publisher),
			RedeemReservationHandler:             commands.NewRedeemReservationHandler(reservations, publisher),
			ExpireReservationHandler:             commands.NewExpireReservationHandler(reservations, publisher),
			CancelReservationHandler:             commands.NewCancelReservationHandler(reservations, publisher),
			RequestReservationNegotiationHandler: commands.NewRequestReservationNegotiationHandler(reservations, publisher),
			AcceptReservationNegotiationHandler:  commands.NewAcceptReservationNegotiationHandler(reservations, publisher),
			DeclineReservationNegotiationHandler: commands.NewDeclineReservationNegotiationHandler(reservations, publisher),
		},
		appQueries: appQueries{
			GetOfferHandler:   queries.NewGetOfferHandler(offers),
			ListOffersHandler: queries.NewListOffersHandler(middleman),
		},
	}
}
