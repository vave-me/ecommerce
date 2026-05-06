package domain

type ShipmentV1 struct {
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
	CarrierID       string
	CarrierName     string
}

func (ShipmentV1) SnapshotName() string { return "shipping.ShipmentV1" }
