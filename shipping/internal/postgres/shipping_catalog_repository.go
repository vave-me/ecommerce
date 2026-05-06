package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/stackus/errors"
	"middleman/internal/postgres"
	"middleman/shipping/internal/domain"
)

type ShippingCatalogRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.ShippingCatalogRepository = (*ShippingCatalogRepository)(nil)

func NewShippingCatalogRepository(tableName string, db postgres.DB) ShippingCatalogRepository {
	return ShippingCatalogRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r ShippingCatalogRepository) Find(ctx context.Context, shipmentID string) (*domain.CatalogShipment, error) {
	const query = `
		SELECT id, product_id, order_id, basket_id, tracking_number, label_url, 
		       sender_name, sender_address, receiver_name, receiver_address,
		       weight, dimensions, service_type, status, carrier_id, carrier_name,
		       created_at, updated_at, pickup_scheduled_at, delivered_at, cancelled_at
		FROM %s
		WHERE id = $1
	`

	shipment := &domain.CatalogShipment{}
	var orderID, basketID, labelURL, carrierID, carrierName sql.NullString
	var pickupScheduledAt, deliveredAt, cancelledAt sql.NullTime

	err := r.db.QueryRowContext(ctx, r.table(query), shipmentID).Scan(
		&shipment.ID,
		&shipment.ProductID,
		&orderID,
		&basketID,
		&shipment.TrackingNumber,
		&labelURL,
		&shipment.SenderName,
		&shipment.SenderAddress,
		&shipment.ReceiverName,
		&shipment.ReceiverAddress,
		&shipment.Weight,
		&shipment.Dimensions,
		&shipment.ServiceType,
		&shipment.Status,
		&carrierID,
		&carrierName,
		&shipment.CreatedAt,
		&shipment.UpdatedAt,
		&pickupScheduledAt,
		&deliveredAt,
		&cancelledAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(errors.ErrNotFound, "shipment not found")
		}
		return nil, errors.Wrap(err, "scanning shipment")
	}

	// Handle nullable fields
	if orderID.Valid {
		shipment.OrderID = orderID.String
	}
	if basketID.Valid {
		shipment.BasketID = basketID.String
	}
	if labelURL.Valid {
		shipment.LabelURL = labelURL.String
	}
	if carrierID.Valid {
		shipment.CarrierID = carrierID.String
	}
	if carrierName.Valid {
		shipment.CarrierName = carrierName.String
	}
	if pickupScheduledAt.Valid {
		shipment.PickupScheduledAt = &pickupScheduledAt.Time
	}
	if deliveredAt.Valid {
		shipment.DeliveredAt = &deliveredAt.Time
	}
	if cancelledAt.Valid {
		shipment.CancelledAt = &cancelledAt.Time
	}

	return shipment, nil
}

func (r ShippingCatalogRepository) GetByOrderID(ctx context.Context, orderID string) ([]*domain.CatalogShipment, error) {
	const query = `
		SELECT id, product_id, order_id, basket_id, tracking_number, label_url, 
		       sender_name, sender_address, receiver_name, receiver_address,
		       weight, dimensions, service_type, status, carrier_id, carrier_name,
		       created_at, updated_at, pickup_scheduled_at, delivered_at, cancelled_at
		FROM %s
		WHERE order_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, r.table(query), orderID)
	if err != nil {
		return nil, errors.Wrap(err, "querying shipments by order ID")
	}
	defer rows.Close()

	return r.scanShipments(rows)
}

func (r ShippingCatalogRepository) GetByBasketID(ctx context.Context, basketID string) ([]*domain.CatalogShipment, error) {
	const query = `
		SELECT id, product_id, order_id, basket_id, tracking_number, label_url, 
		       sender_name, sender_address, receiver_name, receiver_address,
		       weight, dimensions, service_type, status, carrier_id, carrier_name,
		       created_at, updated_at, pickup_scheduled_at, delivered_at, cancelled_at
		FROM %s
		WHERE basket_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, r.table(query), basketID)
	if err != nil {
		return nil, errors.Wrap(err, "querying shipments by basket ID")
	}
	defer rows.Close()

	return r.scanShipments(rows)
}

func (r ShippingCatalogRepository) GetByProductID(ctx context.Context, productID string) ([]*domain.CatalogShipment, error) {
	const query = `
		SELECT id, product_id, order_id, basket_id, tracking_number, label_url, 
		       sender_name, sender_address, receiver_name, receiver_address,
		       weight, dimensions, service_type, status, carrier_id, carrier_name,
		       created_at, updated_at, pickup_scheduled_at, delivered_at, cancelled_at
		FROM %s
		WHERE product_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, r.table(query), productID)
	if err != nil {
		return nil, errors.Wrap(err, "querying shipments by product ID")
	}
	defer rows.Close()

	return r.scanShipments(rows)
}

func (r ShippingCatalogRepository) GetByTrackingNumber(ctx context.Context, trackingNumber string) (*domain.CatalogShipment, error) {
	const query = `
		SELECT id, product_id, order_id, basket_id, tracking_number, label_url, 
		       sender_name, sender_address, receiver_name, receiver_address,
		       weight, dimensions, service_type, status, carrier_id, carrier_name,
		       created_at, updated_at, pickup_scheduled_at, delivered_at, cancelled_at
		FROM %s
		WHERE tracking_number = $1
	`

	shipment := &domain.CatalogShipment{}
	var orderID, basketID, labelURL, carrierID, carrierName sql.NullString
	var pickupScheduledAt, deliveredAt, cancelledAt sql.NullTime

	err := r.db.QueryRowContext(ctx, r.table(query), trackingNumber).Scan(
		&shipment.ID,
		&shipment.ProductID,
		&orderID,
		&basketID,
		&shipment.TrackingNumber,
		&labelURL,
		&shipment.SenderName,
		&shipment.SenderAddress,
		&shipment.ReceiverName,
		&shipment.ReceiverAddress,
		&shipment.Weight,
		&shipment.Dimensions,
		&shipment.ServiceType,
		&shipment.Status,
		&carrierID,
		&carrierName,
		&shipment.CreatedAt,
		&shipment.UpdatedAt,
		&pickupScheduledAt,
		&deliveredAt,
		&cancelledAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(errors.ErrNotFound, "shipment not found")
		}
		return nil, errors.Wrap(err, "scanning shipment")
	}

	// Handle nullable fields
	if orderID.Valid {
		shipment.OrderID = orderID.String
	}
	if basketID.Valid {
		shipment.BasketID = basketID.String
	}
	if labelURL.Valid {
		shipment.LabelURL = labelURL.String
	}
	if carrierID.Valid {
		shipment.CarrierID = carrierID.String
	}
	if carrierName.Valid {
		shipment.CarrierName = carrierName.String
	}
	if pickupScheduledAt.Valid {
		shipment.PickupScheduledAt = &pickupScheduledAt.Time
	}
	if deliveredAt.Valid {
		shipment.DeliveredAt = &deliveredAt.Time
	}
	if cancelledAt.Valid {
		shipment.CancelledAt = &cancelledAt.Time
	}

	return shipment, nil
}

func (r ShippingCatalogRepository) GetShipmentsByStatus(ctx context.Context, status domain.ShipmentStatus, limit, offset int) ([]*domain.CatalogShipment, error) {
	const query = `
		SELECT id, product_id, order_id, basket_id, tracking_number, label_url, 
		       sender_name, sender_address, receiver_name, receiver_address,
		       weight, dimensions, service_type, status, carrier_id, carrier_name,
		       created_at, updated_at, pickup_scheduled_at, delivered_at, cancelled_at
		FROM %s
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, r.table(query), status, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "querying shipments by status")
	}
	defer rows.Close()

	return r.scanShipments(rows)
}

func (r ShippingCatalogRepository) GetShipmentsByCarrier(ctx context.Context, carrierID string, limit, offset int) ([]*domain.CatalogShipment, error) {
	const query = `
		SELECT id, product_id, order_id, basket_id, tracking_number, label_url, 
		       sender_name, sender_address, receiver_name, receiver_address,
		       weight, dimensions, service_type, status, carrier_id, carrier_name,
		       created_at, updated_at, pickup_scheduled_at, delivered_at, cancelled_at
		FROM %s
		WHERE carrier_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, r.table(query), carrierID, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "querying shipments by carrier")
	}
	defer rows.Close()

	return r.scanShipments(rows)
}

func (r ShippingCatalogRepository) GetPendingPickups(ctx context.Context, limit int) ([]*domain.CatalogShipment, error) {
	const query = `
		SELECT id, product_id, order_id, basket_id, tracking_number, label_url, 
		       sender_name, sender_address, receiver_name, receiver_address,
		       weight, dimensions, service_type, status, carrier_id, carrier_name,
		       created_at, updated_at, pickup_scheduled_at, delivered_at, cancelled_at
		FROM %s
		WHERE status = $1 
		  AND pickup_scheduled_at IS NOT NULL 
		  AND pickup_scheduled_at > CURRENT_TIMESTAMP
		ORDER BY pickup_scheduled_at ASC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, r.table(query), domain.ShipmentStatusPickupScheduled, limit)
	if err != nil {
		return nil, errors.Wrap(err, "querying pending pickups")
	}
	defer rows.Close()

	return r.scanShipments(rows)
}

func (r ShippingCatalogRepository) SearchShipments(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*domain.CatalogShipment, error) {
	var conditions []string
	var args []interface{}
	argCount := 1

	query := `
		SELECT id, product_id, order_id, basket_id, tracking_number, label_url, 
		       sender_name, sender_address, receiver_name, receiver_address,
		       weight, dimensions, service_type, status, carrier_id, carrier_name,
		       created_at, updated_at, pickup_scheduled_at, delivered_at, cancelled_at
		FROM %s
		WHERE 1=1
	`

	// Build dynamic conditions based on filters
	for key, value := range filters {
		switch key {
		case "status":
			conditions = append(conditions, fmt.Sprintf("status = $%d", argCount))
			args = append(args, value)
			argCount++
		case "carrier_id":
			conditions = append(conditions, fmt.Sprintf("carrier_id = $%d", argCount))
			args = append(args, value)
			argCount++
		case "service_type":
			conditions = append(conditions, fmt.Sprintf("service_type = $%d", argCount))
			args = append(args, value)
			argCount++
		case "sender_name":
			conditions = append(conditions, fmt.Sprintf("sender_name ILIKE $%d", argCount))
			args = append(args, "%"+value.(string)+"%")
			argCount++
		case "receiver_name":
			conditions = append(conditions, fmt.Sprintf("receiver_name ILIKE $%d", argCount))
			args = append(args, "%"+value.(string)+"%")
			argCount++
		case "date_from":
			conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argCount))
			args = append(args, value)
			argCount++
		case "date_to":
			conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argCount))
			args = append(args, value)
			argCount++
		}
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, r.table(query), args...)
	if err != nil {
		return nil, errors.Wrap(err, "searching shipments")
	}
	defer rows.Close()

	return r.scanShipments(rows)
}

func (r ShippingCatalogRepository) AddShipment(ctx context.Context, shipment *domain.CatalogShipment) error {
	const query = `
		INSERT INTO %s (
			id, product_id, order_id, basket_id, tracking_number, label_url,
			sender_name, sender_address, receiver_name, receiver_address,
			weight, dimensions, service_type, status, carrier_id, carrier_name,
			created_at, updated_at, pickup_scheduled_at, delivered_at, cancelled_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21
		)
	`

	_, err := r.db.ExecContext(ctx, r.table(query),
		shipment.ID,
		shipment.ProductID,
		nullString(shipment.OrderID),
		nullString(shipment.BasketID),
		shipment.TrackingNumber,
		nullString(shipment.LabelURL),
		shipment.SenderName,
		shipment.SenderAddress,
		shipment.ReceiverName,
		shipment.ReceiverAddress,
		shipment.Weight,
		shipment.Dimensions,
		shipment.ServiceType,
		shipment.Status,
		nullString(shipment.CarrierID),
		nullString(shipment.CarrierName),
		shipment.CreatedAt,
		shipment.UpdatedAt,
		nullTime(shipment.PickupScheduledAt),
		nullTime(shipment.DeliveredAt),
		nullTime(shipment.CancelledAt),
	)

	return err
}

func (r ShippingCatalogRepository) UpdateShipmentStatus(ctx context.Context, shipmentID string, status domain.ShipmentStatus, updatedAt time.Time) error {
	const query = `
		UPDATE %s
		SET status = $2, updated_at = $3
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, r.table(query), shipmentID, status, updatedAt)
	return err
}

func (r ShippingCatalogRepository) UpdateTrackingInfo(ctx context.Context, shipmentID, trackingNumber, labelURL string) error {
	const query = `
		UPDATE %s
		SET tracking_number = $2, label_url = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, r.table(query), shipmentID, trackingNumber, nullString(labelURL))
	return err
}

func (r ShippingCatalogRepository) UpdateDeliveryInfo(ctx context.Context, shipmentID string, deliveredAt time.Time) error {
	const query = `
		UPDATE %s
		SET delivered_at = $2, status = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, r.table(query), shipmentID, deliveredAt, domain.ShipmentStatusDelivered)
	return err
}

func (r ShippingCatalogRepository) UpdatePickupInfo(ctx context.Context, shipmentID string, pickupScheduledAt time.Time) error {
	const query = `
		UPDATE %s
		SET pickup_scheduled_at = $2, status = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, r.table(query), shipmentID, pickupScheduledAt, domain.ShipmentStatusPickupScheduled)
	return err
}

func (r ShippingCatalogRepository) CancelShipment(ctx context.Context, shipmentID string, cancelledAt time.Time) error {
	const query = `
		UPDATE %s
		SET cancelled_at = $2, status = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, r.table(query), shipmentID, cancelledAt, domain.ShipmentStatusCancelled)
	return err
}

// Helper function to scan multiple shipments
func (r ShippingCatalogRepository) scanShipments(rows *sql.Rows) ([]*domain.CatalogShipment, error) {
	var shipments []*domain.CatalogShipment

	for rows.Next() {
		shipment := &domain.CatalogShipment{}
		var orderID, basketID, labelURL, carrierID, carrierName sql.NullString
		var pickupScheduledAt, deliveredAt, cancelledAt sql.NullTime

		err := rows.Scan(
			&shipment.ID,
			&shipment.ProductID,
			&orderID,
			&basketID,
			&shipment.TrackingNumber,
			&labelURL,
			&shipment.SenderName,
			&shipment.SenderAddress,
			&shipment.ReceiverName,
			&shipment.ReceiverAddress,
			&shipment.Weight,
			&shipment.Dimensions,
			&shipment.ServiceType,
			&shipment.Status,
			&carrierID,
			&carrierName,
			&shipment.CreatedAt,
			&shipment.UpdatedAt,
			&pickupScheduledAt,
			&deliveredAt,
			&cancelledAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "scanning shipment")
		}

		// Handle nullable fields
		if orderID.Valid {
			shipment.OrderID = orderID.String
		}
		if basketID.Valid {
			shipment.BasketID = basketID.String
		}
		if labelURL.Valid {
			shipment.LabelURL = labelURL.String
		}
		if carrierID.Valid {
			shipment.CarrierID = carrierID.String
		}
		if carrierName.Valid {
			shipment.CarrierName = carrierName.String
		}
		if pickupScheduledAt.Valid {
			shipment.PickupScheduledAt = &pickupScheduledAt.Time
		}
		if deliveredAt.Valid {
			shipment.DeliveredAt = &deliveredAt.Time
		}
		if cancelledAt.Valid {
			shipment.CancelledAt = &cancelledAt.Time
		}

		shipments = append(shipments, shipment)
	}

	return shipments, rows.Err()
}

func (r ShippingCatalogRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}

func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func nullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}