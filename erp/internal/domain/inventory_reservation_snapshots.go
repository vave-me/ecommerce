package domain

import "time"

type InventoryReservationV1 struct {
	OrderID             string
	SKU                 string
	WarehouseID         string
	Quantity            int
	Status              ReservationStatus
	ExpiresAt           *time.Time
	CreatedAt           time.Time
	ReleasedAt          *time.Time
	FulfilledAt         *time.Time
	TransferredAt       *time.Time
	TransferToWarehouse string
	ConnectorID         string
}

func (InventoryReservationV1) SnapshotName() string { return "erp.InventoryReservationV1" }