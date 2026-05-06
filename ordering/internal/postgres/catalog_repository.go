package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/stackus/errors"
	"middleman/internal/postgres"
	"middleman/ordering/internal/domain"
)

type CatalogRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.CatalogRepository = (*CatalogRepository)(nil)

func NewCatalogRepository(tableName string, db postgres.DB) CatalogRepository {
	return CatalogRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r CatalogRepository) AddOrder(ctx context.Context, order *domain.Order) error {
	const query = `
		INSERT INTO %s (id, user_customer_id, payment_method_id, status, total, item_count)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	total := order.GetTotal()
	itemCount := len(order.Items)

	_, err := r.db.ExecContext(ctx, r.table(query), order.ID(), order.UserCustomerID, order.PaymentMethodID,
		order.Status.String(), total, itemCount)

	return errors.Wrap(err, "inserting order into catalog")
}

func (r CatalogRepository) UpdateOrder(ctx context.Context, order *domain.Order) error {
	const query = `
		UPDATE %s 
		SET user_customer_id = $2, payment_method_id = $3, status = $4, total = $5, item_count = $6, updated_at = NOW()
		WHERE id = $1
	`

	total := order.GetTotal()
	itemCount := len(order.Items)

	_, err := r.db.ExecContext(ctx, r.table(query), order.ID(), order.UserCustomerID, order.PaymentMethodID,
		order.Status.String(), total, itemCount)

	return errors.Wrap(err, "updating order in catalog")
}

func (r CatalogRepository) RemoveOrder(ctx context.Context, orderID string) error {
	const query = `DELETE FROM %s WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), orderID)

	return errors.Wrap(err, "removing order from catalog")
}

func (r CatalogRepository) Find(ctx context.Context, orderID string) (*domain.OrderCatalog, error) {
	const query = `
		SELECT id, user_customer_id, payment_method_id, status, total, item_count, created_at, updated_at
		FROM %s
		WHERE id = $1
	`

	order := &domain.OrderCatalog{}

	err := r.db.QueryRowContext(ctx, r.table(query), orderID).Scan(
		&order.ID,
		&order.UserCustomerID,
		&order.PaymentMethodID,
		&order.Status,
		&order.Total,
		&order.ItemCount,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(errors.ErrNotFound, "order not found in catalog")
		}
		return nil, errors.Wrap(err, "finding order in catalog")
	}

	return order, nil
}

func (r CatalogRepository) ListOrders(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*domain.OrderCatalog, int64, error) {
	// Validate and set defaults
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	// Validate sortBy field
	validSortFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"status":     true,
		"total":      true,
	}
	if !validSortFields[sortBy] {
		sortBy = "created_at"
	}

	offset := (page - 1) * pageSize

	// Count total
	var total int64
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, r.tableName)
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, errors.Wrap(err, "counting orders")
	}

	// Get orders
	query := fmt.Sprintf(`
		SELECT id, user_customer_id, payment_method_id, status, total, item_count, created_at, updated_at
		FROM %s
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "listing orders")
	}
	defer rows.Close()

	orders := make([]*domain.OrderCatalog, 0, pageSize)
	for rows.Next() {
		order := &domain.OrderCatalog{}
		err := rows.Scan(
			&order.ID,
			&order.UserCustomerID,
			&order.PaymentMethodID,
			&order.Status,
			&order.Total,
			&order.ItemCount,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scanning order")
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "iterating orders")
	}

	return orders, total, nil
}

func (r CatalogRepository) GetOrdersByCustomer(ctx context.Context, userCustomerID string, page, pageSize int64, sortBy, sortOrder string) ([]*domain.OrderCatalog, int64, error) {
	// Validate and set defaults
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	// Validate sortBy field
	validSortFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"status":     true,
		"total":      true,
	}
	if !validSortFields[sortBy] {
		sortBy = "created_at"
	}

	offset := (page - 1) * pageSize

	// Count total
	var total int64
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE user_customer_id = $1`, r.tableName)
	err := r.db.QueryRowContext(ctx, countQuery, userCustomerID).Scan(&total)
	if err != nil {
		return nil, 0, errors.Wrap(err, "counting customer orders")
	}

	// Get orders
	query := fmt.Sprintf(`
		SELECT id, user_customer_id, payment_method_id, status, total, item_count, created_at, updated_at
		FROM %s
		WHERE user_customer_id = $1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, userCustomerID, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "getting customer orders")
	}
	defer rows.Close()

	orders := make([]*domain.OrderCatalog, 0, pageSize)
	for rows.Next() {
		order := &domain.OrderCatalog{}
		err := rows.Scan(
			&order.ID,
			&order.UserCustomerID,
			&order.PaymentMethodID,
			&order.Status,
			&order.Total,
			&order.ItemCount,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scanning customer order")
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "iterating customer orders")
	}

	return orders, total, nil
}

func (r CatalogRepository) GetOrdersByStatus(ctx context.Context, status domain.OrderStatus, page, pageSize int64, sortBy, sortOrder string) ([]*domain.OrderCatalog, int64, error) {
	// Validate and set defaults
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	// Validate sortBy field
	validSortFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"total":      true,
	}
	if !validSortFields[sortBy] {
		sortBy = "created_at"
	}

	offset := (page - 1) * pageSize

	// Count total
	var total int64
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE status = $1`, r.tableName)
	err := r.db.QueryRowContext(ctx, countQuery, status.String()).Scan(&total)
	if err != nil {
		return nil, 0, errors.Wrap(err, "counting orders by status")
	}

	// Get orders
	query := fmt.Sprintf(`
		SELECT id, user_customer_id, payment_method_id, status, total, item_count, created_at, updated_at
		FROM %s
		WHERE status = $1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, status.String(), pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "getting orders by status")
	}
	defer rows.Close()

	orders := make([]*domain.OrderCatalog, 0, pageSize)
	for rows.Next() {
		order := &domain.OrderCatalog{}
		err := rows.Scan(
			&order.ID,
			&order.UserCustomerID,
			&order.PaymentMethodID,
			&order.Status,
			&order.Total,
			&order.ItemCount,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scanning order by status")
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "iterating orders by status")
	}

	return orders, total, nil
}

func (r CatalogRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
