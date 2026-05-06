package domain

type OrderStatus string

const (
	OrderUnknown     OrderStatus = ""
	OrderIsPending   OrderStatus = "PENDING"  // after creation
	OrderIsApproved  OrderStatus = "APPROVED" // payment confirmed
	OrderIsRejected  OrderStatus = "REJECTED"
	OrderIsCanceled  OrderStatus = "CANCELED"
	OrderIsReady     OrderStatus = "READY"     // ready for shipping
	OrderIsShipped   OrderStatus = "SHIPPED"   // in transit
	OrderIsDelivered OrderStatus = "DELIVERED" // physically delivered
	OrderIsCompleted OrderStatus = "COMPLETED" // final business closure
)

func (s OrderStatus) String() string {
	switch s {
	case OrderIsPending, OrderIsRejected, OrderIsApproved, OrderIsCanceled, OrderIsDelivered, OrderIsReady, OrderIsCompleted, OrderIsShipped:
		return string(s)
	default:
		return ""
	}
}

func ToOrderStatus(status string) OrderStatus {
	switch status {
	case OrderIsPending.String():
		return OrderIsPending
	case OrderIsRejected.String():
		return OrderIsRejected
	case OrderIsApproved.String():
		return OrderIsApproved
	case OrderIsCanceled.String():
		return OrderIsCanceled
	case OrderIsReady.String():
		return OrderIsReady
	case OrderIsDelivered.String():
		return OrderIsDelivered
	case OrderIsCompleted.String():
		return OrderIsCompleted
	case OrderIsShipped.String():
		return OrderIsShipped
	default:
		return OrderUnknown
	}
}
