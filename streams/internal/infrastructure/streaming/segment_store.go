package streaming

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// SegmentStore manages video segment storage
type SegmentStore struct {
	basePath      string
	segmentIndex  map[string]*SegmentInfo
	mu            sync.RWMutex
}

// SegmentInfo contains segment metadata
type SegmentInfo struct {
	StreamID    string
	Quality     string
	SegmentID   string
	FilePath    string
	Duration    float64
	Size        int64
	CreatedAt   time.Time
	AccessedAt  time.Time
	UploadedToCDN bool
}

// NewSegmentStore creates a new segment store
func NewSegmentStore(basePath string) SegmentStore {
	return SegmentStore{
		basePath:     basePath,
		segmentIndex: make(map[string]*SegmentInfo),
	}
}

// StoreSegment stores segment information
func (ss *SegmentStore) StoreSegment(info *SegmentInfo) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	key := fmt.Sprintf("%s/%s/%s", info.StreamID, info.Quality, info.SegmentID)
	ss.segmentIndex[key] = info
	return nil
}

// GetSegment retrieves segment information
func (ss *SegmentStore) GetSegment(streamID, quality, segmentID string) (*SegmentInfo, error) {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	key := fmt.Sprintf("%s/%s/%s", streamID, quality, segmentID)
	info, exists := ss.segmentIndex[key]
	if !exists {
		return nil, fmt.Errorf("segment not found")
	}

	// Update access time
	info.AccessedAt = time.Now()
	return info, nil
}

// GetNewSegments finds segments that haven't been uploaded to CDN
func (ss *SegmentStore) GetNewSegments(outputDir string) ([]string, error) {
	files, err := ioutil.ReadDir(outputDir)
	if err != nil {
		return nil, err
	}

	var newSegments []string
	ss.mu.Lock()
	defer ss.mu.Unlock()

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".ts") {
			continue
		}

		filePath := filepath.Join(outputDir, file.Name())
		
		// Check if we've already processed this segment
		streamID := filepath.Base(filepath.Dir(filepath.Dir(filePath)))
		quality := filepath.Base(filepath.Dir(filePath))
		segmentID := file.Name()
		
		key := fmt.Sprintf("%s/%s/%s", streamID, quality, segmentID)
		if info, exists := ss.segmentIndex[key]; exists && info.UploadedToCDN {
			continue
		}

		// New segment found
		newSegments = append(newSegments, filePath)
		
		// Store segment info
		ss.segmentIndex[key] = &SegmentInfo{
			StreamID:    streamID,
			Quality:     quality,
			SegmentID:   segmentID,
			FilePath:    filePath,
			Size:        file.Size(),
			CreatedAt:   file.ModTime(),
			AccessedAt:  time.Now(),
			UploadedToCDN: false,
		}
	}

	return newSegments, nil
}

// MarkUploaded marks segments as uploaded to CDN
func (ss *SegmentStore) MarkUploaded(streamID, quality, segmentID string) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	key := fmt.Sprintf("%s/%s/%s", streamID, quality, segmentID)
	if info, exists := ss.segmentIndex[key]; exists {
		info.UploadedToCDN = true
		return nil
	}

	return fmt.Errorf("segment not found")
}

// GetStreamSegments returns all segments for a stream
func (ss *SegmentStore) GetStreamSegments(streamID string) ([]*SegmentInfo, error) {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	var segments []*SegmentInfo
	for key, info := range ss.segmentIndex {
		if strings.HasPrefix(key, streamID+"/") {
			segments = append(segments, info)
		}
	}

	// Sort by creation time
	sort.Slice(segments, func(i, j int) bool {
		return segments[i].CreatedAt.Before(segments[j].CreatedAt)
	})

	return segments, nil
}

// CleanupOldSegments removes segments older than DVR window
func (ss *SegmentStore) CleanupOldSegments(dvrWindowMinutes int) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	cutoffTime := time.Now().Add(-time.Duration(dvrWindowMinutes) * time.Minute)
	var toDelete []string

	for key, info := range ss.segmentIndex {
		if info.CreatedAt.Before(cutoffTime) && info.UploadedToCDN {
			toDelete = append(toDelete, key)
			
			// Delete physical file
			if err := os.Remove(info.FilePath); err != nil && !os.IsNotExist(err) {
				fmt.Printf("Failed to delete segment file %s: %v\n", info.FilePath, err)
			}
		}
	}

	// Remove from index
	for _, key := range toDelete {
		delete(ss.segmentIndex, key)
	}

	return nil
}

// GetStorageStats returns storage statistics
func (ss *SegmentStore) GetStorageStats() map[string]interface{} {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	totalSize := int64(0)
	totalSegments := len(ss.segmentIndex)
	uploadedCount := 0
	streamCount := make(map[string]bool)

	for _, info := range ss.segmentIndex {
		totalSize += info.Size
		if info.UploadedToCDN {
			uploadedCount++
		}
		streamCount[info.StreamID] = true
	}

	return map[string]interface{}{
		"totalSegments":   totalSegments,
		"totalSize":       totalSize,
		"totalSizeMB":     float64(totalSize) / 1024 / 1024,
		"uploadedCount":   uploadedCount,
		"pendingUploads":  totalSegments - uploadedCount,
		"activeStreams":   len(streamCount),
		"averageSegmentSize": func() int64 {
			if totalSegments > 0 {
				return totalSize / int64(totalSegments)
			}
			return 0
		}(),
	}
}

// PruneOrphanedFiles removes files not in the index
func (ss *SegmentStore) PruneOrphanedFiles() error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	// Build set of known files
	knownFiles := make(map[string]bool)
	for _, info := range ss.segmentIndex {
		knownFiles[info.FilePath] = true
	}

	// Walk storage directory
	orphanedCount := 0
	err := filepath.Walk(ss.basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Check if file is a segment
		if strings.HasSuffix(path, ".ts") || strings.HasSuffix(path, ".m4s") {
			if !knownFiles[path] {
				// Orphaned file found
				if err := os.Remove(path); err != nil {
					fmt.Printf("Failed to remove orphaned file %s: %v\n", path, err)
				} else {
					orphanedCount++
				}
			}
		}

		return nil
	})

	if orphanedCount > 0 {
		fmt.Printf("Pruned %d orphaned segment files\n", orphanedCount)
	}

	return err
}