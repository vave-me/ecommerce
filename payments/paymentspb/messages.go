package paymentspb

import (
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

const (
	InvoiceAggregateChannel = "middleman.payments.events.Invoice"
	PaymentAggregateChannel = "middleman.payments.events.Payment"
	InvoicePaidEvent        = "paymentsapi.InvoicePaid"
	CommandChannel          = "middleman.payments.commands"
	PaymentConfirmedEvent   = "middleman.payments.PaymentConfirmed"
	PaymentAuthorizedEvent  = "middleman.payments.PaymentAuthorized"
	ConfirmPaymentCommand   = "paymentsapi.ConfirmPayment"
	AuthorizePaymentCommand = "paymentsapi.AuthorizePayment"
)

func Registrations(reg registry.Registry) (err error) {
	serde := serdes.NewProtoSerde(reg)

	// Invoice events
	if err = serde.Register(&InvoicePaid{}); err != nil {
		return err
	}
	if err = serde.Register(&PaymentConfirmed{}); err != nil {
		return err
	}
	if err = serde.Register(&PaymentAuthorized{}); err != nil {
		return err
	}

	// commands
	if err = serde.Register(&ConfirmPayment{}); err != nil {
		return
	}
	if err = serde.Register(&AuthorizePaymentRequest{}); err != nil {
		return
	}

	return
}

func (*InvoicePaid) Key() string             { return InvoicePaidEvent }
func (*ConfirmPayment) Key() string          { return ConfirmPaymentCommand }
func (*PaymentConfirmed) Key() string        { return PaymentConfirmedEvent }
func (*PaymentAuthorized) Key() string       { return PaymentAuthorizedEvent }
func (*AuthorizePaymentRequest) Key() string { return AuthorizePaymentCommand }
