package monitoring

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// StreamingMetrics contains all streaming-related metrics
type StreamingMetrics struct {
	// Stream metrics
	ActiveStreams      prometheus.Gauge
	TotalStreams       prometheus.Counter
	StreamDuration     prometheus.Histogram
	StreamErrors       prometheus.CounterVec
	
	// Viewer metrics
	ConcurrentViewers  prometheus.GaugeVec
	TotalViewers       prometheus.CounterVec
	ViewerDuration     prometheus.HistogramVec
	ViewerGeography    prometheus.CounterVec
	
	// Quality metrics
	StreamBitrate      prometheus.GaugeVec
	DroppedFrames      prometheus.CounterVec
	BufferingEvents    prometheus.CounterVec
	StreamLatency      prometheus.HistogramVec
	StartupTime        prometheus.HistogramVec
	
	// CDN metrics
	CDNBandwidth       prometheus.GaugeVec
	CDNCacheHitRate    prometheus.GaugeVec
	CDNErrors          prometheus.CounterVec
	SegmentUploadTime  prometheus.HistogramVec
	
	// Infrastructure metrics
	TranscodingQueue   prometheus.Gauge
	TranscodingErrors  prometheus.Counter
	CPUUsage           prometheus.GaugeVec
	GPUUsage           prometheus.GaugeVec
	MemoryUsage        prometheus.GaugeVec
	DiskIOPS           prometheus.GaugeVec
	NetworkThroughput  prometheus.GaugeVec
	
	// DRM metrics
	LicenseRequests    prometheus.CounterVec
	LicenseErrors      prometheus.CounterVec
	LicenseLatency     prometheus.HistogramVec
	
	// WebRTC metrics
	WebRTCConnections  prometheus.Gauge
	WebRTCPacketLoss   prometheus.GaugeVec
	WebRTCJitter       prometheus.GaugeVec
	WebRTCRTT          prometheus.GaugeVec

	// OpenTelemetry
	tracer       trace.Tracer
	meter        metric.Meter
	logger       *zap.Logger
}

// NewStreamingMetrics creates a new metrics instance
func NewStreamingMetrics(namespace string, logger *zap.Logger) *StreamingMetrics {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}

	return &StreamingMetrics{
		// Stream metrics
		ActiveStreams: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "streams",
			Name:      "active_total",
			Help:      "Number of currently active streams",
		}),
		
		TotalStreams: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "streams",
			Name:      "total",
			Help:      "Total number of streams created",
		}),
		
		StreamDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "streams",
			Name:      "duration_seconds",
			Help:      "Stream duration in seconds",
			Buckets:   prometheus.ExponentialBuckets(60, 2, 10), // 1min to ~17hours
		}),
		
		StreamErrors: *promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "streams",
			Name:      "errors_total",
			Help:      "Total number of stream errors",
		}, []string{"stream_id", "error_type"}),
		
		// Viewer metrics
		ConcurrentViewers: *promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "viewers",
			Name:      "concurrent_total",
			Help:      "Number of concurrent viewers per stream",
		}, []string{"stream_id"}),
		
		TotalViewers: *promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "viewers",
			Name:      "total",
			Help:      "Total number of viewers",
		}, []string{"stream_id"}),
		
		ViewerDuration: *promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "viewers",
			Name:      "session_duration_seconds",
			Help:      "Viewer session duration in seconds",
			Buckets:   prometheus.ExponentialBuckets(10, 2, 10), // 10s to ~3hours
		}, []string{"stream_id"}),
		
		ViewerGeography: *promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "viewers",
			Name:      "by_country_total",
			Help:      "Total viewers by country",
		}, []string{"stream_id", "country"}),
		
		// Quality metrics
		StreamBitrate: *promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "quality",
			Name:      "bitrate_kbps",
			Help:      "Current stream bitrate in kilobits per second",
		}, []string{"stream_id", "quality"}),
		
		DroppedFrames: *promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "quality",
			Name:      "dropped_frames_total",
			Help:      "Total number of dropped frames",
		}, []string{"stream_id", "session_id"}),
		
		BufferingEvents: *promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "quality",
			Name:      "buffering_events_total",
			Help:      "Total number of buffering events",
		}, []string{"stream_id", "session_id"}),
		
		StreamLatency: *promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "quality",
			Name:      "latency_milliseconds",
			Help:      "Stream latency in milliseconds",
			Buckets:   prometheus.ExponentialBuckets(100, 2, 10), // 100ms to ~100s
		}, []string{"stream_id", "protocol"}),
		
		StartupTime: *promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "quality",
			Name:      "startup_time_milliseconds",
			Help:      "Time to first frame in milliseconds",
			Buckets:   prometheus.ExponentialBuckets(100, 1.5, 10), // 100ms to ~5s
		}, []string{"stream_id", "protocol"}),
		
		// CDN metrics
		CDNBandwidth: *promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "cdn",
			Name:      "bandwidth_mbps",
			Help:      "CDN bandwidth usage in megabits per second",
		}, []string{"provider", "region"}),
		
		CDNCacheHitRate: *promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "cdn",
			Name:      "cache_hit_rate",
			Help:      "CDN cache hit rate (0-1)",
		}, []string{"provider", "content_type"}),
		
		CDNErrors: *promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "cdn",
			Name:      "errors_total",
			Help:      "Total CDN errors",
		}, []string{"provider", "error_type"}),
		
		SegmentUploadTime: *promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "cdn",
			Name:      "segment_upload_milliseconds",
			Help:      "Time to upload segment to CDN",
			Buckets:   prometheus.ExponentialBuckets(10, 2, 10), // 10ms to ~10s
		}, []string{"provider", "region"}),
		
		// Infrastructure metrics
		TranscodingQueue: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "transcoding",
			Name:      "queue_size",
			Help:      "Number of segments in transcoding queue",
		}),
		
		TranscodingErrors: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "transcoding",
			Name:      "errors_total",
			Help:      "Total transcoding errors",
		}),
		
		CPUUsage: *promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "infrastructure",
			Name:      "cpu_usage_percent",
			Help:      "CPU usage percentage",
		}, []string{"core"}),
		
		GPUUsage: *promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "infrastructure",
			Name:      "gpu_usage_percent",
			Help:      "GPU usage percentage",
		}, []string{"device"}),
		
		MemoryUsage: *promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "infrastructure",
			Name:      "memory_usage_bytes",
			Help:      "Memory usage in bytes",
		}, []string{"type"}),
		
		DiskIOPS: *promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "infrastructure",
			Name:      "disk_iops",
			Help:      "Disk I/O operations per second",
		}, []string{"device", "operation"}),
		
		NetworkThroughput: *promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "infrastructure",
			Name:      "network_throughput_mbps",
			Help:      "Network throughput in megabits per second",
		}, []string{"interface", "direction"}),
		
		// DRM metrics
		LicenseRequests: *promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "drm",
			Name:      "license_requests_total",
			Help:      "Total DRM license requests",
		}, []string{"provider", "stream_id"}),
		
		LicenseErrors: *promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "drm",
			Name:      "license_errors_total",
			Help:      "Total DRM license errors",
		}, []string{"provider", "error_type"}),
		
		LicenseLatency: *promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "drm",
			Name:      "license_latency_milliseconds",
			Help:      "DRM license generation latency",
			Buckets:   prometheus.ExponentialBuckets(10, 2, 8), // 10ms to ~2.5s
		}, []string{"provider"}),
		
		// WebRTC metrics
		WebRTCConnections: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "webrtc",
			Name:      "connections_total",
			Help:      "Total WebRTC connections",
		}),
		
		WebRTCPacketLoss: *promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "webrtc",
			Name:      "packet_loss_percent",
			Help:      "WebRTC packet loss percentage",
		}, []string{"stream_id", "session_id"}),
		
		WebRTCJitter: *promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "webrtc",
			Name:      "jitter_milliseconds",
			Help:      "WebRTC jitter in milliseconds",
		}, []string{"stream_id", "session_id"}),
		
		WebRTCRTT: *promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "webrtc",
			Name:      "rtt_milliseconds",
			Help:      "WebRTC round trip time in milliseconds",
		}, []string{"stream_id", "session_id"}),

		// OpenTelemetry
		tracer: otel.Tracer(namespace),
		meter:  otel.Meter(namespace),
		logger: logger,
	}
}

// RecordStreamStart records a stream start event
func (m *StreamingMetrics) RecordStreamStart(streamID string) {
	m.ActiveStreams.Inc()
	m.TotalStreams.Inc()
	
	m.logger.Info("Stream started",
		zap.String("stream_id", streamID),
		zap.Time("timestamp", time.Now()))
}

// RecordStreamEnd records a stream end event
func (m *StreamingMetrics) RecordStreamEnd(streamID string, duration time.Duration) {
	m.ActiveStreams.Dec()
	m.StreamDuration.Observe(duration.Seconds())
	
	// Clear stream-specific gauges
	m.ConcurrentViewers.DeleteLabelValues(streamID)
	m.StreamBitrate.DeleteLabelValues(streamID)
	
	m.logger.Info("Stream ended",
		zap.String("stream_id", streamID),
		zap.Duration("duration", duration))
}

// RecordViewerJoin records a viewer joining a stream
func (m *StreamingMetrics) RecordViewerJoin(streamID, sessionID, country string) {
	m.ConcurrentViewers.WithLabelValues(streamID).Inc()
	m.TotalViewers.WithLabelValues(streamID).Inc()
	m.ViewerGeography.WithLabelValues(streamID, country).Inc()
}

// RecordViewerLeave records a viewer leaving a stream
func (m *StreamingMetrics) RecordViewerLeave(streamID, sessionID string, duration time.Duration) {
	m.ConcurrentViewers.WithLabelValues(streamID).Dec()
	m.ViewerDuration.WithLabelValues(streamID).Observe(duration.Seconds())
}

// RecordQualityMetrics records stream quality metrics
func (m *StreamingMetrics) RecordQualityMetrics(streamID string, quality string, bitrate float64) {
	m.StreamBitrate.WithLabelValues(streamID, quality).Set(bitrate)
}

// RecordBufferingEvent records a buffering event
func (m *StreamingMetrics) RecordBufferingEvent(streamID, sessionID string) {
	m.BufferingEvents.WithLabelValues(streamID, sessionID).Inc()
}

// RecordDroppedFrames records dropped frames
func (m *StreamingMetrics) RecordDroppedFrames(streamID, sessionID string, count float64) {
	m.DroppedFrames.WithLabelValues(streamID, sessionID).Add(count)
}

// RecordStreamLatency records stream latency
func (m *StreamingMetrics) RecordStreamLatency(streamID, protocol string, latency time.Duration) {
	m.StreamLatency.WithLabelValues(streamID, protocol).Observe(float64(latency.Milliseconds()))
}

// RecordStartupTime records time to first frame
func (m *StreamingMetrics) RecordStartupTime(streamID, protocol string, duration time.Duration) {
	m.StartupTime.WithLabelValues(streamID, protocol).Observe(float64(duration.Milliseconds()))
}

// RecordCDNMetrics records CDN performance metrics
func (m *StreamingMetrics) RecordCDNMetrics(provider, region string, bandwidth, cacheHitRate float64) {
	m.CDNBandwidth.WithLabelValues(provider, region).Set(bandwidth)
	m.CDNCacheHitRate.WithLabelValues(provider, "video").Set(cacheHitRate)
}

// RecordCDNError records a CDN error
func (m *StreamingMetrics) RecordCDNError(provider, errorType string) {
	m.CDNErrors.WithLabelValues(provider, errorType).Inc()
}

// RecordSegmentUpload records segment upload time
func (m *StreamingMetrics) RecordSegmentUpload(provider, region string, duration time.Duration) {
	m.SegmentUploadTime.WithLabelValues(provider, region).Observe(float64(duration.Milliseconds()))
}

// RecordTranscodingMetrics records transcoding queue metrics
func (m *StreamingMetrics) RecordTranscodingMetrics(queueSize float64) {
	m.TranscodingQueue.Set(queueSize)
}

// RecordTranscodingError records a transcoding error
func (m *StreamingMetrics) RecordTranscodingError() {
	m.TranscodingErrors.Inc()
}

// RecordInfrastructureMetrics records infrastructure metrics
func (m *StreamingMetrics) RecordInfrastructureMetrics(metrics InfrastructureMetrics) {
	// CPU usage
	for i, usage := range metrics.CPUUsage {
		m.CPUUsage.WithLabelValues(string(i)).Set(usage)
	}
	
	// GPU usage
	for device, usage := range metrics.GPUUsage {
		m.GPUUsage.WithLabelValues(device).Set(usage)
	}
	
	// Memory usage
	m.MemoryUsage.WithLabelValues("used").Set(float64(metrics.MemoryUsed))
	m.MemoryUsage.WithLabelValues("available").Set(float64(metrics.MemoryAvailable))
	
	// Disk IOPS
	m.DiskIOPS.WithLabelValues(metrics.DiskDevice, "read").Set(metrics.DiskReadIOPS)
	m.DiskIOPS.WithLabelValues(metrics.DiskDevice, "write").Set(metrics.DiskWriteIOPS)
	
	// Network throughput
	m.NetworkThroughput.WithLabelValues(metrics.NetworkInterface, "rx").Set(metrics.NetworkRxMbps)
	m.NetworkThroughput.WithLabelValues(metrics.NetworkInterface, "tx").Set(metrics.NetworkTxMbps)
}

// RecordDRMLicenseRequest records a DRM license request
func (m *StreamingMetrics) RecordDRMLicenseRequest(provider, streamID string) {
	m.LicenseRequests.WithLabelValues(provider, streamID).Inc()
}

// RecordDRMLicenseError records a DRM license error
func (m *StreamingMetrics) RecordDRMLicenseError(provider, errorType string) {
	m.LicenseErrors.WithLabelValues(provider, errorType).Inc()
}

// RecordDRMLicenseLatency records DRM license generation latency
func (m *StreamingMetrics) RecordDRMLicenseLatency(provider string, latency time.Duration) {
	m.LicenseLatency.WithLabelValues(provider).Observe(float64(latency.Milliseconds()))
}

// RecordWebRTCConnection records WebRTC connection metrics
func (m *StreamingMetrics) RecordWebRTCConnection(connected bool) {
	if connected {
		m.WebRTCConnections.Inc()
	} else {
		m.WebRTCConnections.Dec()
	}
}

// RecordWebRTCMetrics records WebRTC quality metrics
func (m *StreamingMetrics) RecordWebRTCMetrics(streamID, sessionID string, packetLoss, jitter, rtt float64) {
	m.WebRTCPacketLoss.WithLabelValues(streamID, sessionID).Set(packetLoss)
	m.WebRTCJitter.WithLabelValues(streamID, sessionID).Set(jitter)
	m.WebRTCRTT.WithLabelValues(streamID, sessionID).Set(rtt)
}

// StartSpan starts a new tracing span
func (m *StreamingMetrics) StartSpan(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return m.tracer.Start(ctx, name, trace.WithAttributes(attrs...))
}

// InfrastructureMetrics contains infrastructure metrics
type InfrastructureMetrics struct {
	CPUUsage         []float64
	GPUUsage         map[string]float64
	MemoryUsed       int64
	MemoryAvailable  int64
	DiskDevice       string
	DiskReadIOPS     float64
	DiskWriteIOPS    float64
	NetworkInterface string
	NetworkRxMbps    float64
	NetworkTxMbps    float64
}