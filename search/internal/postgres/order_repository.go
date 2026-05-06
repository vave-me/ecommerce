package postgres

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/stackus/errors"
	"middleman/internal/postgres"
	"middleman/search/internal/application"
	"middleman/search/internal/models"
	"strings"
)

type OrderRepository struct {
	tableName string
	db        postgres.DB
}

var _ application.OrderRepository = (*OrderRepository)(nil)

func NewOrderRepository(tableName string, db postgres.DB) OrderRepository {
	return OrderRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r OrderRepository) Add(ctx context.Context, order *models.Order) error {
	const query = `INSERT INTO %s (
order_id, user_customer_id, user_customer_name,
items, status, product_ids, user_seller_ids,
created_at VALUES (
$1, $2, $3,
$4, $5, $6, $7
$8
)`
	items, err := json.Marshal(order.Items)
	if err != nil {
		return err
	}

	productIDs := make(IDArray, len(order.Items))
	userSellerMap := make(map[string]struct{})
	for i, item := range order.Items {
		productIDs[i] = item.ProductID
		userSellerMap[item.UserSellerID] = struct{}{}
	}
	userSellerIDs := make(IDArray, 0, len(userSellerMap))
	for userSellerID, _ := range userSellerMap {
		userSellerIDs = append(userSellerIDs, userSellerID)
	}

	_, err = r.db.ExecContext(ctx, r.table(query),
		order.OrderID, order.UserCustomerID, order.UserCustomerName,
		items, order.Status, productIDs, userSellerIDs,
		order.CreatedAt,
	)
	return err
}

func (r OrderRepository) UpdateStatus(ctx context.Context, orderID, status string) error {
	const query = `UPDATE %s SET status = $2 WHERE order_id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), orderID, status)
	return err
}

//
//func (r OrderRepository) Search(ctx context.Context, search application.SearchOrders) ([]*models.Order, error) {
//	const baseQuery = `SELECT order_id, customer_id, customer_name, items, status, total, created_at FROM %s`
//
//	var (
//		conditions []string
//		args       []interface{}
//		argID      = 1
//	)
//
//	// Build conditions based on filters
//	if search.Filters.UserCustomerID != "" {
//		conditions = append(conditions, fmt.Sprintf("customer_id = $%d", argID))
//		args = append(args, search.Filters.UserCustomerID)
//		argID++
//	}
//
//	if !search.Filters.After.IsZero() {
//		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argID))
//		args = append(args, search.Filters.After)
//		argID++
//	}
//
//	if !search.Filters.Before.IsZero() {
//		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argID))
//		args = append(args, search.Filters.Before)
//		argID++
//	}
//
//	if len(search.Filters.UserSellerIDs) > 0 {
//		conditions = append(conditions, fmt.Sprintf("user_seller_ids && $%d", argID))
//		args = append(args, IDArray(search.Filters.UserSellerIDs))
//		argID++
//	}
//
//	if len(search.Filters.ProductIDs) > 0 {
//		conditions = append(conditions, fmt.Sprintf("product_ids && $%d", argID))
//		args = append(args, IDArray(search.Filters.ProductIDs))
//		argID++
//	}
//
//	if search.Filters.MinTotal > 0 {
//		conditions = append(conditions, fmt.Sprintf("total >= $%d", argID))
//		args = append(args, search.Filters.MinTotal)
//		argID++
//	}
//
//	if search.Filters.MaxTotal > 0 {
//		conditions = append(conditions, fmt.Sprintf("total <= $%d", argID))
//		args = append(args, search.Filters.MaxTotal)
//		argID++
//	}
//
//	if search.Filters.Status != "" {
//		conditions = append(conditions, fmt.Sprintf("status = $%d", argID))
//		args = append(args, search.Filters.Status)
//		argID++
//	}
//
//	// Build the final query
//	query := r.table(baseQuery)
//
//	if len(conditions) > 0 {
//		query += " WHERE " + strings.Join(conditions, " AND ")
//	}
//
//	query += " ORDER BY created_at DESC"
//
//	if search.Limit > 0 {
//		query += fmt.Sprintf(" LIMIT $%d", argID)
//		args = append(args, search.Limit)
//		argID++
//	}
//
//	if search.Offset > 0 {
//		query += fmt.Sprintf(" OFFSET $%d", argID)
//		args = append(args, search.Offset)
//		argID++
//	}
//
//	// Execute the query
//	rows, err := r.db.QueryContext(ctx, query, args...)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var orders []*models.Order
//
//	for rows.Next() {
//		order := &models.Order{}
//		var itemData []byte
//		err := rows.Scan(&order.OrderID, &order.UserCustomerID, &order.UserCustomerName, &itemData, &order.Status, &order.Total, &order.CreatedAt)
//		if err != nil {
//			return nil, err
//		}
//
//		// Deserialize items from JSON
//		err = json.Unmarshal(itemData, &order.Items)
//		if err != nil {
//			return nil, err
//		}
//
//		orders = append(orders, order)
//	}
//
//	if err := rows.Err(); err != nil {
//		return nil, err
//	}
//
//	return orders, nil
//}

func (r OrderRepository) Get(ctx context.Context, orderID string) (*models.Order, error) {
	const query = `SELECT customer_id, customer_name, items, status, created_at FROM %s WHERE order_id = $1`

	order := &models.Order{
		OrderID: orderID,
	}

	var itemData []byte
	err := r.db.QueryRowContext(ctx, r.table(query)).Scan(&order.UserCustomerID, &order.UserCustomerName, &itemData, &order.Status, &order.CreatedAt)
	if err != nil {
		return nil, err
	}

	var items []models.Item
	err = json.Unmarshal(itemData, &items)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return order, nil
}

func (r OrderRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}

type IDArray []string

func (a *IDArray) Scan(src any) error {
	var sep = []byte(",")

	var data []byte
	switch v := src.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return errors.ErrInvalidArgument.Msgf("IDArray: unsupported type: %T", src)
	}

	ids := make([]string, bytes.Count(data, sep))
	for i, id := range bytes.Split(bytes.Trim(data, "{}"), sep) {
		ids[i] = string(id)
	}

	*a = ids

	return nil
}

func (a IDArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	if len(a) == 0 {
		return "{}", nil
	}
	// unsafe way to do this; assumption is all ids are UUIDs
	return fmt.Sprintf("{%s}", strings.Join(a, ",")), nil
}
