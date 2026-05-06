// File: vector/internal/redis/pagination_utils.go
package redis

// calculatePagination computes offset and limit from page and pageSize parameters
func calculatePagination(page, pageSize int64) (offset, limit int64) {
	if pageSize <= 0 {
		pageSize = 20 // Default page size
	}
	limit = pageSize
	if page <= 0 {
		page = 1 // Default to first page
	}
	offset = (page - 1) * pageSize
	return
}

// calculatePaginationWithParams computes offset and limit with fallback to offsetParam/limitParam
func calculatePaginationWithParams(page, pageSize, offsetParam, limitParam int64) (offset, limit int64) {
	offset = offsetParam
	limit = limitParam

	// If page/pageSize are provided, they take precedence
	if pageSize > 0 {
		limit = pageSize
		if page < 1 {
			page = 1
		}
		offset = (page - 1) * pageSize
	}

	// Apply defaults if needed
	if limit <= 0 {
		limit = 50 // Default limit
	}
	if offset < 0 {
		offset = 0 // Never allow negative offset
	}
	return
}
