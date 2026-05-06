package streaming

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/stackus/errors"
)

// StreamingServer handles live stream ingestion and distribution
type StreamingServer struct {
	config            *StreamingConfig
	manifestGenerator *ManifestGenerator
	segmentStore      SegmentStore
	cdnManager        *CDNManager
	qualityProfiles   []QualityProfile
	activeStreams     map[string]*ActiveStream
	mu                sync.RWMutex
	httpServer        *http.Server
}

// StreamingConfig contains server configuration
type StreamingConfig struct {
	IngestPort         int
	StreamingPort      int
	RTMPEnabled        bool
	SRTEnabled         bool
	WebRTCEnabled      bool
	SegmentDuration    int // seconds
	PlaylistSize       int // number of segments
	StoragePath        string
	CDNUploadEnabled   bool
	LowLatencyMode     bool
	DVRWindowMinutes   int
}

// QualityProfile defines encoding settings for different qualities
type QualityProfile struct {
	Name       string
	Width      int
	Height     int
	Bitrate    int    // kbps
	Framerate  int    // fps
	Codec      string // h264, h265, av1
	Profile    string // baseline, main, high
	Level      string // 3.0, 3.1, 4.0, 4.1
}

// ActiveStream represents a currently live stream
type ActiveStream struct {
	StreamID          string
	StreamKey         string
	IngestProtocol    string
	StartTime         time.Time
	CurrentViewers    int64
	PeakViewers       int64
	SegmentsGenerated int64
	LastSegmentTime   time.Time
	Transcoder        *Transcoder
	mu                sync.RWMutex
}

// NewStreamingServer creates a new streaming server
func NewStreamingServer(config *StreamingConfig, cdnManager *CDNManager) *StreamingServer {
	return &StreamingServer{
		config:            config,
		manifestGenerator: NewManifestGenerator(config),
		segmentStore:      NewSegmentStore(config.StoragePath),
		cdnManager:        cdnManager,
		activeStreams:     make(map[string]*ActiveStream),
		qualityProfiles:   getDefaultQualityProfiles(),
	}
}

// Start starts the streaming server
func (s *StreamingServer) Start(ctx context.Context) error {
	// Start RTMP ingest server
	if s.config.RTMPEnabled {
		go s.startRTMPServer(ctx)
	}

	// Start SRT ingest server
	if s.config.SRTEnabled {
		go s.startSRTServer(ctx)
	}

	// Start HTTP server for HLS/DASH delivery
	router := mux.NewRouter()
	s.setupHTTPRoutes(router)

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.StreamingPort),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	// Start segment cleanup routine
	go s.cleanupOldSegments(ctx)

	<-ctx.Done()
	return s.Stop()
}

// Stop stops the streaming server
func (s *StreamingServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}

// StartStream begins processing a new live stream
func (s *StreamingServer) StartStream(streamID, streamKey, ingestProtocol string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.activeStreams[streamID]; exists {
		return errors.Wrap(errors.ErrConflict, "stream already active")
	}

	// Create transcoder for multiple quality profiles
	transcoder := NewTranscoder(s.qualityProfiles, s.config)

	activeStream := &ActiveStream{
		StreamID:       streamID,
		StreamKey:      streamKey,
		IngestProtocol: ingestProtocol,
		StartTime:      time.Now(),
		Transcoder:     transcoder,
	}

	s.activeStreams[streamID] = activeStream

	// Start transcoding pipeline
	go s.processStream(streamID)

	return nil
}

// StopStream stops processing a live stream
func (s *StreamingServer) StopStream(streamID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	stream, exists := s.activeStreams[streamID]
	if !exists {
		return errors.Wrap(errors.ErrNotFound, "stream not found")
	}

	// Stop transcoder
	if stream.Transcoder != nil {
		stream.Transcoder.Stop()
	}

	delete(s.activeStreams, streamID)
	return nil
}

// processStream handles the transcoding and segmentation pipeline
func (s *StreamingServer) processStream(streamID string) {
	stream := s.getActiveStream(streamID)
	if stream == nil {
		return
	}

	// Create output directories
	baseDir := filepath.Join(s.config.StoragePath, streamID)
	os.MkdirAll(baseDir, 0755)

	// Process each quality profile
	var wg sync.WaitGroup
	for _, profile := range s.qualityProfiles {
		wg.Add(1)
		go func(p QualityProfile) {
			defer wg.Done()
			s.processQualityProfile(streamID, p)
		}(profile)
	}

	wg.Wait()
}

// processQualityProfile handles transcoding for a specific quality
func (s *StreamingServer) processQualityProfile(streamID string, profile QualityProfile) {
	stream := s.getActiveStream(streamID)
	if stream == nil {
		return
	}

	outputDir := filepath.Join(s.config.StoragePath, streamID, profile.Name)
	os.MkdirAll(outputDir, 0755)

	// Generate segments
	segmentPattern := filepath.Join(outputDir, "segment_%d.ts")
	playlistPath := filepath.Join(outputDir, "playlist.m3u8")

	// Start transcoding to this quality
	transcodeOpts := TranscodeOptions{
		InputProtocol:   stream.IngestProtocol,
		OutputFormat:    "hls",
		SegmentDuration: s.config.SegmentDuration,
		SegmentPattern:  segmentPattern,
		PlaylistPath:    playlistPath,
		VideoCodec:      profile.Codec,
		VideoBitrate:    profile.Bitrate,
		VideoWidth:      profile.Width,
		VideoHeight:     profile.Height,
		VideoFramerate:  profile.Framerate,
		VideoProfile:    profile.Profile,
		VideoLevel:      profile.Level,
		LowLatency:      s.config.LowLatencyMode,
	}

	stream.Transcoder.TranscodeToQuality(profile.Name, transcodeOpts)

	// Monitor segment generation
	s.monitorSegments(streamID, profile.Name, outputDir)
}

// monitorSegments watches for new segments and uploads to CDN
func (s *StreamingServer) monitorSegments(streamID, quality, outputDir string) {
	for {
		stream := s.getActiveStream(streamID)
		if stream == nil {
			break
		}

		// Check for new segments
		segments, err := s.segmentStore.GetNewSegments(outputDir)
		if err == nil && len(segments) > 0 {
			// Upload to CDN if enabled
			if s.config.CDNUploadEnabled && s.cdnManager != nil {
				for _, segment := range segments {
					s.cdnManager.UploadSegment(streamID, quality, segment)
				}
			}

			// Update stream statistics
			stream.mu.Lock()
			stream.SegmentsGenerated += int64(len(segments))
			stream.LastSegmentTime = time.Now()
			stream.mu.Unlock()
		}

		time.Sleep(1 * time.Second)
	}
}

// setupHTTPRoutes configures HTTP endpoints
func (s *StreamingServer) setupHTTPRoutes(router *mux.Router) {
	// HLS endpoints
	router.HandleFunc("/hls/{streamID}/master.m3u8", s.handleMasterPlaylist).Methods("GET")
	router.HandleFunc("/hls/{streamID}/{quality}/playlist.m3u8", s.handlePlaylist).Methods("GET")
	router.HandleFunc("/hls/{streamID}/{quality}/{segment}", s.handleSegment).Methods("GET")

	// DASH endpoints
	router.HandleFunc("/dash/{streamID}/manifest.mpd", s.handleDASHManifest).Methods("GET")
	router.HandleFunc("/dash/{streamID}/{quality}/{segment}", s.handleDASHSegment).Methods("GET")

	// Stream info endpoint
	router.HandleFunc("/streams/{streamID}/info", s.handleStreamInfo).Methods("GET")

	// Health check
	router.HandleFunc("/health", s.handleHealth).Methods("GET")
}

// handleMasterPlaylist serves the HLS master playlist
func (s *StreamingServer) handleMasterPlaylist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	streamID := vars["streamID"]

	stream := s.getActiveStream(streamID)
	if stream == nil {
		http.Error(w, "Stream not found", http.StatusNotFound)
		return
	}

	// Generate master playlist
	masterPlaylist := s.manifestGenerator.GenerateMasterPlaylist(streamID, s.qualityProfiles)

	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(masterPlaylist))
}

// handlePlaylist serves quality-specific playlists
func (s *StreamingServer) handlePlaylist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	streamID := vars["streamID"]
	quality := vars["quality"]

	playlistPath := filepath.Join(s.config.StoragePath, streamID, quality, "playlist.m3u8")
	
	// Check if playlist exists
	if _, err := os.Stat(playlistPath); os.IsNotExist(err) {
		http.Error(w, "Playlist not found", http.StatusNotFound)
		return
	}

	// Serve playlist file
	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	w.Header().Set("Cache-Control", "no-cache")
	http.ServeFile(w, r, playlistPath)
}

// handleSegment serves video segments
func (s *StreamingServer) handleSegment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	streamID := vars["streamID"]
	quality := vars["quality"]
	segment := vars["segment"]

	// Validate segment name
	if !strings.HasSuffix(segment, ".ts") {
		http.Error(w, "Invalid segment", http.StatusBadRequest)
		return
	}

	segmentPath := filepath.Join(s.config.StoragePath, streamID, quality, segment)

	// Check if segment exists locally
	if _, err := os.Stat(segmentPath); os.IsNotExist(err) {
		// Try to fetch from CDN
		if s.cdnManager != nil {
			cdnURL := s.cdnManager.GetSegmentURL(streamID, quality, segment)
			http.Redirect(w, r, cdnURL, http.StatusFound)
			return
		}
		http.Error(w, "Segment not found", http.StatusNotFound)
		return
	}

	// Serve segment file
	w.Header().Set("Content-Type", "video/mp2t")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	http.ServeFile(w, r, segmentPath)
}

// handleDASHManifest serves DASH manifest
func (s *StreamingServer) handleDASHManifest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	streamID := vars["streamID"]

	stream := s.getActiveStream(streamID)
	if stream == nil {
		http.Error(w, "Stream not found", http.StatusNotFound)
		return
	}

	// Generate DASH manifest
	manifest := s.manifestGenerator.GenerateDASHManifest(streamID, s.qualityProfiles)

	w.Header().Set("Content-Type", "application/dash+xml")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(manifest))
}

// handleDASHSegment serves DASH segments
func (s *StreamingServer) handleDASHSegment(w http.ResponseWriter, r *http.Request) {
	// Similar to HLS segment handling
	s.handleSegment(w, r)
}

// handleStreamInfo returns stream information
func (s *StreamingServer) handleStreamInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	streamID := vars["streamID"]

	stream := s.getActiveStream(streamID)
	if stream == nil {
		http.Error(w, "Stream not found", http.StatusNotFound)
		return
	}

	stream.mu.RLock()
	info := map[string]interface{}{
		"streamId":          stream.StreamID,
		"startTime":         stream.StartTime,
		"currentViewers":    stream.CurrentViewers,
		"peakViewers":       stream.PeakViewers,
		"segmentsGenerated": stream.SegmentsGenerated,
		"lastSegmentTime":   stream.LastSegmentTime,
		"isLive":            true,
	}
	stream.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	// JSON encoding omitted for brevity
}

// handleHealth returns server health status
func (s *StreamingServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// startRTMPServer starts RTMP ingest server
func (s *StreamingServer) startRTMPServer(ctx context.Context) {
	// RTMP server implementation
	// Uses github.com/yutopp/go-rtmp or similar library
}

// startSRTServer starts SRT ingest server
func (s *StreamingServer) startSRTServer(ctx context.Context) {
	// SRT server implementation
	// Uses github.com/Haivision/srtgo or similar library
}

// cleanupOldSegments removes old segments based on DVR window
func (s *StreamingServer) cleanupOldSegments(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.segmentStore.CleanupOldSegments(s.config.DVRWindowMinutes)
		}
	}
}

// getActiveStream safely retrieves an active stream
func (s *StreamingServer) getActiveStream(streamID string) *ActiveStream {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.activeStreams[streamID]
}

// getDefaultQualityProfiles returns standard quality profiles
func getDefaultQualityProfiles() []QualityProfile {
	return []QualityProfile{
		{
			Name:      "1080p",
			Width:     1920,
			Height:    1080,
			Bitrate:   8000,
			Framerate: 60,
			Codec:     "h264",
			Profile:   "high",
			Level:     "4.2",
		},
		{
			Name:      "720p",
			Width:     1280,
			Height:    720,
			Bitrate:   4000,
			Framerate: 60,
			Codec:     "h264",
			Profile:   "main",
			Level:     "4.0",
		},
		{
			Name:      "480p",
			Width:     854,
			Height:    480,
			Bitrate:   2000,
			Framerate: 30,
			Codec:     "h264",
			Profile:   "main",
			Level:     "3.1",
		},
		{
			Name:      "360p",
			Width:     640,
			Height:    360,
			Bitrate:   1000,
			Framerate: 30,
			Codec:     "h264",
			Profile:   "baseline",
			Level:     "3.0",
		},
	}
}