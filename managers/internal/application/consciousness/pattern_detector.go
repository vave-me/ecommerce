package consciousness

import (
	"context"
	"middleman/internal/ddd"
)

// PatternDetector manages multiple pattern detection functions
type PatternDetector struct {
	detectors []PatternDetectorFunc
	memoryCore *MemoryCore
}

// NewPatternDetector creates a new pattern detector with the given detector functions
func NewPatternDetector(detectors ...PatternDetectorFunc) *PatternDetector {
	return &PatternDetector{
		detectors: detectors,
	}
}

// SetMemoryCore sets the memory core for accessing recent events
func (pd *PatternDetector) SetMemoryCore(memoryCore *MemoryCore) {
	pd.memoryCore = memoryCore
}

// DetectPattern runs all pattern detectors and returns the first pattern found
func (pd *PatternDetector) DetectPattern(ctx context.Context, event ddd.Event) *Pattern {
	if pd.memoryCore == nil {
		return nil
	}
	
	// Get recent events from memory
	recentEvents := pd.memoryCore.GetRecentEvents(1000)
	
	// Run each detector
	for _, detector := range pd.detectors {
		if pattern := detector(event, recentEvents); pattern != nil {
			return pattern
		}
	}
	
	return nil
}