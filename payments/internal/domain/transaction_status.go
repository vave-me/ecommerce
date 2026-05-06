package domain

type TransactionStatus string

const (
	TransactionStatusUnknown           TransactionStatus = "Unknown"
	TransactionStatusPending           TransactionStatus = "Pending"
	TransactionStatusCleared           TransactionStatus = "Cleared"
	TransactionStatusRefunded          TransactionStatus = "Refunded"
	TransactionStatusPartiallyRefunded TransactionStatus = "PartiallyRefunded"
)

func (t TransactionStatus) String() string {
	switch t {
	case TransactionStatusPending, TransactionStatusCleared, TransactionStatusRefunded, TransactionStatusPartiallyRefunded:
		return string(t)
	default:
		return ""
	}
}

func ToTransactionStatus(s string) TransactionStatus {
	switch s {
	case TransactionStatusPending.String():
		return TransactionStatusPending
	case TransactionStatusCleared.String():
		return TransactionStatusCleared
	case TransactionStatusRefunded.String():
		return TransactionStatusRefunded
	case TransactionStatusPartiallyRefunded.String():
		return TransactionStatusPartiallyRefunded
	default:
		return TransactionStatusUnknown
	}
}
