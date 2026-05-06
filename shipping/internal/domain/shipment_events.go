package domain

const (
	// Shipment Entity Lifecycle Events
	ShipmentCreatedEvent       = "shipping.ShipmentCreated"
	ShipmentDeletedEvent       = "shipping.ShipmentDeleted"
	ShipmentStatusUpdatedEvent = "shipping.ShipmentStatusUpdated"
	CarrierAssignedEvent       = "shipping.CarrierAssigned"
	ShipmentStartedEvent       = "shipping.ShipmentStarted"
	ShipmentCancelledEvent     = "shipping.ShipmentCancelled"
	ShipmentDeliveredEvent     = "shipping.ShipmentDelivered"
	PickupScheduledEvent       = "shipping.PickupScheduled"
	ShipmentReturnedEvent      = "shipping.ShipmentReturned"
)

type ShipmentCreated struct {
	ID              string
	ProductID       string
	OrderID         string
	BasketID        string
	TrackingNumber  string
	LabelUrl        string
	SenderName      string
	SenderAddress   string
	ReceiverName    string
	ReceiverAddress string
	Weight          string
	Dimensions      string
	ServiceType     string
	Status          string
}

// Key implements registry.Registerable
func (ShipmentCreated) Key() string { return ShipmentCreatedEvent }

type ShipmentDeleted struct{}

// Key implements registry.Registerable
func (ShipmentDeleted) Key() string { return ShipmentDeletedEvent }

type ShipmentStatusUpdated struct {
	ShipmentID string
	Status     string
	Location   string
	Notes      string
}

// Key implements registry.Registerable
func (ShipmentStatusUpdated) Key() string { return ShipmentStatusUpdatedEvent }

type CarrierAssigned struct {
	ShipmentID  string
	CarrierID   string
	CarrierName string
}

// Key implements registry.Registerable
func (CarrierAssigned) Key() string { return CarrierAssignedEvent }

type ShipmentStarted struct {
	ShipmentID string
}

// Key implements registry.Registerable
func (ShipmentStarted) Key() string { return ShipmentStartedEvent }

type ShipmentCancelled struct {
	ShipmentID string
	Reason     string
}

// Key implements registry.Registerable
func (ShipmentCancelled) Key() string { return ShipmentCancelledEvent }

type ShipmentDelivered struct {
	ShipmentID         string
	SignedBy           string
	ProofOfDeliveryURL string
}

// Key implements registry.Registerable
func (ShipmentDelivered) Key() string { return ShipmentDeliveredEvent }

type PickupScheduled struct {
	ShipmentID   string
	PickupTime   string
	Instructions string
}

// Key implements registry.Registerable
func (PickupScheduled) Key() string { return PickupScheduledEvent }

type ShipmentReturned struct {
	ShipmentID           string
	Reason               string
	ReturnTrackingNumber string
}

// Key implements registry.Registerable
func (ShipmentReturned) Key() string { return ShipmentReturnedEvent }
