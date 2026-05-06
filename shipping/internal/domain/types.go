package domain

// ServiceType represents the type of shipping service
type ServiceType string

const (
	ServiceTypeStandard  ServiceType = "standard"
	ServiceTypeExpress   ServiceType = "express"
	ServiceTypeOvernight ServiceType = "overnight"
	ServiceTypeSameDay   ServiceType = "same_day"
	ServiceTypeEconomy   ServiceType = "economy"
)

// ShipmentStatus represents the status of a shipment
type ShipmentStatus string

const (
	ShipmentStatusCreated          ShipmentStatus = "created"
	ShipmentStatusAssigned         ShipmentStatus = "assigned"
	ShipmentStatusPickupScheduled  ShipmentStatus = "pickup_scheduled"
	ShipmentStatusInTransit        ShipmentStatus = "in_transit"
	ShipmentStatusOutForDelivery   ShipmentStatus = "out_for_delivery"
	ShipmentStatusDelivered        ShipmentStatus = "delivered"
	ShipmentStatusCancelled        ShipmentStatus = "cancelled"
	ShipmentStatusReturned         ShipmentStatus = "returned"
	ShipmentStatusFailed           ShipmentStatus = "failed"
)

// String returns the string representation of ServiceType
func (s ServiceType) String() string {
	return string(s)
}

// String returns the string representation of ShipmentStatus
func (s ShipmentStatus) String() string {
	return string(s)
}