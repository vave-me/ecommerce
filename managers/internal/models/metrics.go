package models

import "time"

// SalesMetrics represents sales metrics for a specific date range
type SalesMetrics struct {
	TotalSales     float64   `json:"total_sales"`
	TotalOrders    int64     `json:"total_orders"`
	AverageOrder   float64   `json:"average_order"`
	TopProducts    []Product `json:"top_products"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	ConversionRate float64   `json:"conversion_rate"`
}

// DashboardMetrics represents comprehensive dashboard metrics
type DashboardMetrics struct {
	TotalUsers         int64         `json:"total_users"`
	ActiveUsers        int64         `json:"active_users"`
	TotalProducts      int64         `json:"total_products"`
	TotalOrders        int64         `json:"total_orders"`
	TotalRevenue       float64       `json:"total_revenue"`
	PendingOrders      int64         `json:"pending_orders"`
	CompletedOrders    int64         `json:"completed_orders"`
	CancelledOrders    int64         `json:"cancelled_orders"`
	SalesMetrics       SalesMetrics  `json:"sales_metrics"`
	TopCategories      []Category    `json:"top_categories"`
	RecentActivity     []interface{} `json:"recent_activity"`
	UpdatedAt          time.Time     `json:"updated_at"`
}