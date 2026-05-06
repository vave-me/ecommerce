package queries

import (
	"context"
	"middleman/shipping/internal/domain"
)

type (
	ListShipments struct {
		Limit     int
		Offset    int
		Status    string
		ProductID string
		OrderID   string
		CarrierID string
		StartDate string
		EndDate   string
	}

	ListShipmentsHandler struct {
		catalog domain.ShippingCatalogRepository
	}
)

func NewListShipmentsHandler(catalog domain.ShippingCatalogRepository) ListShipmentsHandler {
	return ListShipmentsHandler{
		catalog: catalog,
	}
}

func (h ListShipmentsHandler) ListShipments(ctx context.Context, query ListShipments) ([]*domain.CatalogShipment, error) {
	filters := make(map[string]interface{})
	
	if query.Status != "" {
		filters["status"] = query.Status
	}
	if query.ProductID != "" {
		return h.catalog.GetByProductID(ctx, query.ProductID)
	}
	if query.OrderID != "" {
		return h.catalog.GetByOrderID(ctx, query.OrderID)
	}
	if query.CarrierID != "" {
		filters["carrier_id"] = query.CarrierID
	}
	if query.StartDate != "" {
		filters["date_from"] = query.StartDate
	}
	if query.EndDate != "" {
		filters["date_to"] = query.EndDate
	}
	
	limit := query.Limit
	if limit == 0 {
		limit = 20 // Default limit
	}
	
	return h.catalog.SearchShipments(ctx, filters, limit, query.Offset)
}
