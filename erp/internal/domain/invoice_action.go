package domain

type InvoiceAction string

const (
	InvoiceActionCreate  InvoiceAction = "create"
	InvoiceActionUpdate  InvoiceAction = "update"
	InvoiceActionApprove InvoiceAction = "approve"
	InvoiceActionVoid    InvoiceAction = "void"
	InvoiceActionSend    InvoiceAction = "send"
	InvoiceActionPayment InvoiceAction = "payment"
)
