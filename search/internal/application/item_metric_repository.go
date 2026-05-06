package application

// This file is deprecated - use MetricRepository instead
// Keeping for backward compatibility but should be removed in next major version

// ItemMetricRepository is deprecated - use MetricRepository
// Deprecated: Use MetricRepository instead
type ItemMetricRepository interface {
	MetricRepository
}

// ItemMetricCacheRepository is deprecated - use MetricCacheRepository
// Deprecated: Use MetricCacheRepository instead
type ItemMetricCacheRepository interface {
	MetricCacheRepository
}
