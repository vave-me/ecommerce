package redis

import (
	"fmt"
	"middleman/vectors/internal/models"
	"strings"
	"time"
)

// QueryBuilder provides fluent interface for building Redis Search queries
// with Redis expert optimizations for high-performance vector
type QueryBuilder struct {
	entityType      models.EntityType
	filters         []string
	sortField       string
	sortDesc        bool
	offset          int
	limit           int
	fields          []string
	timeout         time.Duration
	verbatim        bool // Skip stemming for exact matches
	noStopWords     bool // Skip stop word filtering
	withScores      bool // Include relevance scores
	maxPrefixLength int  // Limit prefix expansion
}

// RedisSearchQuery represents a Redis Search query
type RedisSearchQuery struct {
	Query        string
	Offset       int
	Limit        int
	SortBy       string
	SortDesc     bool
	ReturnFields []string
}

// escapeRedisSearchString escapes special characters for Redis Search
func escapeRedisSearchString(s string) string {
	// Escape special Redis Search characters
	replacer := strings.NewReplacer(
		"@", "\\@",
		"(", "\\(",
		")", "\\)",
		"{", "\\{",
		"}", "\\}",
		"[", "\\[",
		"]", "\\]",
		"|", "\\|",
		"*", "\\*",
		"?", "\\?",
		"~", "\\~",
		"\"", "\\\"",
	)
	return replacer.Replace(s)
}

// QueryMetrics holds performance metrics for monitoring
type QueryMetrics struct {
	QueryTime      time.Duration
	ResultCount    int
	CacheHit       bool
	IndexScanned   int64
	EstimatedTotal int64
}

// NewQueryBuilder creates a new optimized query builder
func NewQueryBuilder(entityType models.EntityType) *QueryBuilder {
	return &QueryBuilder{
		entityType:      entityType,
		filters:         make([]string, 0),
		limit:           50,              // Reasonable default
		timeout:         3 * time.Second, // Prevent slow queries
		maxPrefixLength: 3,               // Limit prefix expansion
	}
}

// WithPerformanceMode enables high-performance query optimizations
func (qb *QueryBuilder) WithPerformanceMode() *QueryBuilder {
	qb.verbatim = true           // Skip expensive stemming
	qb.noStopWords = true        // Skip stop word processing
	qb.maxPrefixLength = 2       // Aggressive prefix limiting
	qb.timeout = 1 * time.Second // Faster timeout
	return qb
}

// WithTimeout sets query execution timeout
func (qb *QueryBuilder) WithTimeout(timeout time.Duration) *QueryBuilder {
	qb.timeout = timeout
	return qb
}

// WithScores enables relevance score calculation
func (qb *QueryBuilder) WithScores() *QueryBuilder {
	qb.withScores = true
	return qb
}

// WithTimeConstraint adds a time-based filter to avoid full collection scans
func (qb *QueryBuilder) WithTimeConstraint(maxAge time.Duration) *QueryBuilder {
	if maxAge > 0 {
		minTimestamp := time.Now().Add(-maxAge).Unix()
		qb.filters = append(qb.filters, fmt.Sprintf("@created_at:[%d +inf]", minTimestamp))
	}
	return qb
}

// WithStatus filters entities by their status fields
func (qb *QueryBuilder) WithStatus(statuses ...string) *QueryBuilder {
	if len(statuses) > 0 {
		statusFilters := make([]string, len(statuses))
		for i, status := range statuses {
			statusFilters[i] = escapeRedisSearchString(status)
		}
		qb.filters = append(qb.filters, fmt.Sprintf("@status:{%s}", strings.Join(statusFilters, "|")))
	}
	return qb
}

// WithNameFilter adds an optimized name filter with prefix limiting
func (qb *QueryBuilder) WithNameFilter(name string) *QueryBuilder {
	if name != "" {
		escaped := escapeRedisSearchString(name)
		// For very short prefixes, use exact matching to avoid explosion
		if len(name) <= qb.maxPrefixLength {
			qb.filters = append(qb.filters, fmt.Sprintf("@name:{%s}", escaped))
		} else {
			qb.filters = append(qb.filters, fmt.Sprintf("@name:(%s)", escaped))
		}
	}
	return qb
}

// WithGeoFilter adds a geographical filter using lat, lng and radius
func (qb *QueryBuilder) WithGeoFilter(lat, lng float64, radiusKm int64) *QueryBuilder {
	if lat != 0 && lng != 0 && radiusKm > 0 {
		qb.filters = append(qb.filters, fmt.Sprintf("@location:[%.6f %.6f %d km]", lng, lat, radiusKm))
	}
	return qb
}

// WithPriceRange adds an optimized price range filter
func (qb *QueryBuilder) WithPriceRange(min, max int64) *QueryBuilder {
	// Only add price filter if we have meaningful constraints
	// Don't add filter for very wide ranges that essentially mean "all prices"
	if min > 0 || (max > 0 && max < 99999999) {
		minStr := "-inf"
		if min > 0 {
			minStr = fmt.Sprintf("%d", min)
		}

		maxStr := "+inf"
		if max > 0 && max < 99999999 {
			maxStr = fmt.Sprintf("%d", max)
		}

		qb.filters = append(qb.filters, fmt.Sprintf("@base_price:[%s %s]", minStr, maxStr))
	}
	return qb
}

// WithSorting sets the sort field and direction
func (qb *QueryBuilder) WithSorting(field string, desc bool) *QueryBuilder {
	qb.sortField = field
	qb.sortDesc = desc
	return qb
}

// WithPagination sets offset and limit with safety bounds
func (qb *QueryBuilder) WithPagination(offset, limit int) *QueryBuilder {
	// Enforce reasonable limits to prevent memory issues
	qb.offset = offset
	if limit > 0 && limit <= 1000 { // Max 1000 results per query
		qb.limit = limit
	} else if limit > 1000 {
		qb.limit = 1000
	}
	return qb
}

// WithReturnFields specifies which fields should be returned
func (qb *QueryBuilder) WithReturnFields(fields ...string) *QueryBuilder {
	qb.fields = fields
	return qb
}

// WithCustomFilter adds a custom filter expression
func (qb *QueryBuilder) WithCustomFilter(filter string) *QueryBuilder {
	if filter != "" {
		qb.filters = append(qb.filters, filter)
	}
	return qb
}

// Build constructs the final optimized query with entity type filtering
// Returns the query string and the RedisSearchQuery object
func (qb *QueryBuilder) Build() (string, *RedisSearchQuery) {
	// Always include entity type filter to prevent cross-contamination
	baseQuery := fmt.Sprintf("@entity_type:{%s}", qb.entityType.String())
	var queryStr string

	if len(qb.filters) > 0 {
		queryStr = fmt.Sprintf("%s %s", baseQuery, strings.Join(qb.filters, " "))
	} else {
		queryStr = baseQuery
	}

	// Create the optimized query
	query := &RedisSearchQuery{
		Query:        queryStr,
		Offset:       qb.offset,
		Limit:        qb.limit,
		SortBy:       qb.sortField,
		SortDesc:     qb.sortDesc,
		ReturnFields: qb.fields,
	}

	return queryStr, query
}

// BuildWithMetrics executes the query and returns performance metrics
// This function is currently disabled as it requires integration with Redis Search client
// func (qb *QueryBuilder) BuildWithMetrics(ctx context.Context, client redis.Client) (*QueryMetrics, []interface{}, error) {
//     // TODO: Implement with proper Redis Search client
//     return nil, nil, fmt.Errorf("BuildWithMetrics not implemented")
// }
