package processor

import (
	"sync"
	"time"
)

// ToolPerformanceMetrics tracks tool success rates and patterns
type ToolPerformanceMetrics struct {
	SuccessCount    int64
	FailureCount    int64
	TotalDuration   time.Duration
	AverageDuration time.Duration
	LastUsed        time.Time
	CommonPatterns  map[string]int
	mutex           sync.RWMutex
}

// NewToolPerformanceMetrics creates a new performance metrics instance
func NewToolPerformanceMetrics() *ToolPerformanceMetrics {
	return &ToolPerformanceMetrics{
		CommonPatterns: make(map[string]int),
	}
}

// RecordSuccess records a successful tool execution
func (m *ToolPerformanceMetrics) RecordSuccess(duration time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.SuccessCount++
	m.TotalDuration += duration
	m.LastUsed = time.Now()
	
	// Update average duration
	totalCount := m.SuccessCount + m.FailureCount
	if totalCount > 0 {
		m.AverageDuration = m.TotalDuration / time.Duration(totalCount)
	}
}

// RecordFailure records a failed tool execution
func (m *ToolPerformanceMetrics) RecordFailure(duration time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.FailureCount++
	m.TotalDuration += duration
	m.LastUsed = time.Now()
	
	// Update average duration
	totalCount := m.SuccessCount + m.FailureCount
	if totalCount > 0 {
		m.AverageDuration = m.TotalDuration / time.Duration(totalCount)
	}
}

// GetSuccessRate returns the success rate as a percentage
func (m *ToolPerformanceMetrics) GetSuccessRate() float64 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	totalCount := m.SuccessCount + m.FailureCount
	if totalCount == 0 {
		return 0.0
	}
	
	return float64(m.SuccessCount) / float64(totalCount) * 100.0
}

// RecordPattern records a common usage pattern
func (m *ToolPerformanceMetrics) RecordPattern(pattern string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.CommonPatterns[pattern]++
}