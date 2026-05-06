package streaming

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/nareix/joy5/av"
	"github.com/nareix/joy5/format/flv"
	"github.com/nareix/joy5/format/rtmp"
	"github.com/stackus/errors"
	"go.uber.org/zap"
)

// RTMPServer handles RTMP stream ingestion
type RTMPServer struct {
	config          *RTMPConfig
	server          *rtmp.Server
	activeStreams   map[string]*RTMPStream
	streamHandler   StreamHandler
	authenticator   StreamAuthenticator
	logger          *zap.Logger
	mu              sync.RWMutex
}

// RTMPConfig contains RTMP server configuration
type RTMPConfig struct {
	ListenAddr           string
	PublishTimeout       time.Duration
	ReadTimeout          time.Duration
	WriteTimeout         time.Duration
	MaxMessageSize       int
	EnableAuth           bool
	EnableSSL            bool
	CertFile             string
	KeyFile              string
	MaxConcurrentStreams int
}

// RTMPStream represents an active RTMP stream
type RTMPStream struct {
	StreamID        string
	StreamKey       string
	PublisherAddr   string
	StartTime       time.Time
	BytesReceived   int64
	PacketsReceived int64
	VideoCodec      string
	AudioCodec      string
	Width           int
	Height          int
	Framerate       float64
	Bitrate         int64
	LastPacketTime  time.Time
	conn            *rtmp.Conn
	writeChan       chan av.Packet
	ctx             context.Context
	cancel          context.CancelFunc
}

// StreamHandler processes incoming stream data
type StreamHandler interface {
	HandleStream(streamID string, reader av.Demuxer) error
	OnStreamStart(streamID string, metadata StreamMetadata) error
	OnStreamEnd(streamID string, stats StreamStats) error
}

// StreamAuthenticator validates stream keys
type StreamAuthenticator interface {
	ValidateStreamKey(streamID, streamKey string) error
	GetStreamConfig(streamID string) (*StreamConfig, error)
}

// StreamMetadata contains stream information
type StreamMetadata struct {
	StreamID   string
	VideoCodec string
	AudioCodec string
	Width      int
	Height     int
	Framerate  float64
	Bitrate    int64
}

// StreamStats contains stream statistics
type StreamStats struct {
	Duration        time.Duration
	BytesReceived   int64
	PacketsReceived int64
	AverageBitrate  int64
	DroppedFrames   int64
}

// StreamConfig contains stream configuration
type StreamConfig struct {
	MaxBitrate     int64
	MaxResolution  string
	RequireAuth    bool
	AllowedIPs     []string
	RecordingEnabled bool
}

// NewRTMPServer creates a new RTMP server
func NewRTMPServer(config *RTMPConfig, handler StreamHandler, auth StreamAuthenticator, logger *zap.Logger) *RTMPServer {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}

	return &RTMPServer{
		config:        config,
		activeStreams: make(map[string]*RTMPStream),
		streamHandler: handler,
		authenticator: auth,
		logger:        logger,
	}
}

// Start starts the RTMP server
func (rs *RTMPServer) Start(ctx context.Context) error {
	rs.server = &rtmp.Server{
		Addr: rs.config.ListenAddr,
		HandlePublish: rs.handlePublish,
		HandlePlay:    rs.handlePlay,
		HandleConn:    rs.handleConnection,
	}

	// Configure timeouts
	rs.server.ReadTimeout = rs.config.ReadTimeout
	rs.server.WriteTimeout = rs.config.WriteTimeout

	// Start listener
	var listener net.Listener
	var err error

	if rs.config.EnableSSL {
		// TLS configuration for RTMPS
		listener, err = rs.createTLSListener()
		if err != nil {
			return err
		}
	} else {
		listener, err = net.Listen("tcp", rs.config.ListenAddr)
		if err != nil {
			return err
		}
	}

	rs.logger.Info("RTMP server started",
		zap.String("address", rs.config.ListenAddr),
		zap.Bool("ssl", rs.config.EnableSSL))

	// Accept connections
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
					rs.logger.Error("Failed to accept connection", zap.Error(err))
					continue
				}
			}

			go rs.handleClientConnection(ctx, conn)
		}
	}()

	<-ctx.Done()
	return listener.Close()
}

// handleClientConnection handles individual client connections
func (rs *RTMPServer) handleClientConnection(ctx context.Context, netConn net.Conn) {
	defer netConn.Close()

	// Apply connection limits
	rs.mu.RLock()
	activeCount := len(rs.activeStreams)
	rs.mu.RUnlock()

	if activeCount >= rs.config.MaxConcurrentStreams {
		rs.logger.Warn("Maximum concurrent streams reached", 
			zap.Int("current", activeCount),
			zap.Int("max", rs.config.MaxConcurrentStreams))
		return
	}

	conn := rtmp.NewConn(netConn)
	if err := conn.HandshakeServer(); err != nil {
		rs.logger.Error("RTMP handshake failed", zap.Error(err))
		return
	}

	// Process connection
	if err := rs.server.HandleConn(conn); err != nil {
		if err != io.EOF {
			rs.logger.Error("Connection handling error", zap.Error(err))
		}
	}
}

// handlePublish handles incoming stream publishing
func (rs *RTMPServer) handlePublish(conn *rtmp.Conn, streams []av.CodecData) error {
	// Extract stream information from URL
	// URL format: rtmp://server/app/streamID?key=streamKey
	url := conn.URL
	streamID, streamKey := rs.parseStreamInfo(url.Path, url.RawQuery)

	rs.logger.Info("Publisher connected",
		zap.String("streamId", streamID),
		zap.String("addr", conn.NetConn().RemoteAddr().String()))

	// Validate stream key
	if rs.config.EnableAuth {
		if err := rs.authenticator.ValidateStreamKey(streamID, streamKey); err != nil {
			rs.logger.Warn("Authentication failed",
				zap.String("streamId", streamID),
				zap.Error(err))
			return errors.Wrap(errors.ErrUnauthorized, "invalid stream key")
		}
	}

	// Get stream configuration
	streamConfig, err := rs.authenticator.GetStreamConfig(streamID)
	if err != nil {
		return err
	}

	// Validate source IP if configured
	if len(streamConfig.AllowedIPs) > 0 {
		clientIP := conn.NetConn().RemoteAddr().(*net.TCPAddr).IP.String()
		if !rs.isIPAllowed(clientIP, streamConfig.AllowedIPs) {
			rs.logger.Warn("Unauthorized IP",
				zap.String("streamId", streamID),
				zap.String("ip", clientIP))
			return errors.Wrap(errors.ErrForbidden, "IP not allowed")
		}
	}

	// Create stream context
	streamCtx, cancel := context.WithCancel(context.Background())

	// Create RTMP stream
	stream := &RTMPStream{
		StreamID:      streamID,
		StreamKey:     streamKey,
		PublisherAddr: conn.NetConn().RemoteAddr().String(),
		StartTime:     time.Now(),
		conn:          conn,
		writeChan:     make(chan av.Packet, 100),
		ctx:           streamCtx,
		cancel:        cancel,
	}

	// Extract codec information
	for _, codec := range streams {
		switch codec.Type() {
		case av.H264:
			stream.VideoCodec = "H264"
			// Extract resolution and framerate from SPS
			if h264, ok := codec.(av.VideoCodecData); ok {
				stream.Width = h264.Width()
				stream.Height = h264.Height()
			}
		case av.AAC:
			stream.AudioCodec = "AAC"
		}
	}

	// Register stream
	rs.mu.Lock()
	rs.activeStreams[streamID] = stream
	rs.mu.Unlock()

	// Notify stream start
	metadata := StreamMetadata{
		StreamID:   streamID,
		VideoCodec: stream.VideoCodec,
		AudioCodec: stream.AudioCodec,
		Width:      stream.Width,
		Height:     stream.Height,
		Framerate:  stream.Framerate,
		Bitrate:    stream.Bitrate,
	}

	if err := rs.streamHandler.OnStreamStart(streamID, metadata); err != nil {
		rs.logger.Error("Failed to start stream handling", 
			zap.String("streamId", streamID),
			zap.Error(err))
		return err
	}

	// Start packet reader
	go rs.readPackets(stream)

	// Handle stream
	defer func() {
		rs.stopStream(streamID)
	}()

	// Create demuxer for the connection
	return rs.streamHandler.HandleStream(streamID, conn)
}

// readPackets reads packets from the RTMP connection
func (rs *RTMPServer) readPackets(stream *RTMPStream) {
	defer close(stream.writeChan)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var bytesLastSecond int64
	var startTime = time.Now()

	for {
		select {
		case <-stream.ctx.Done():
			return
		case <-ticker.C:
			// Calculate bitrate
			duration := time.Since(startTime).Seconds()
			if duration > 0 {
				stream.Bitrate = int64(float64(stream.BytesReceived) * 8 / duration / 1000) // kbps
			}
			bytesLastSecond = 0
		default:
			// Read packet with timeout
			stream.conn.NetConn().SetReadDeadline(time.Now().Add(rs.config.ReadTimeout))
			
			pkt, err := stream.conn.ReadPacket()
			if err != nil {
				if err == io.EOF {
					return
				}
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				rs.logger.Error("Error reading packet",
					zap.String("streamId", stream.StreamID),
					zap.Error(err))
				return
			}

			// Update statistics
			stream.PacketsReceived++
			stream.BytesReceived += int64(len(pkt.Data))
			stream.LastPacketTime = time.Now()
			bytesLastSecond += int64(len(pkt.Data))

			// Forward packet
			select {
			case stream.writeChan <- pkt:
			default:
				// Channel full, drop packet
				rs.logger.Warn("Dropping packet due to full buffer",
					zap.String("streamId", stream.StreamID))
			}
		}
	}
}

// handlePlay handles play requests (not implemented for ingest-only server)
func (rs *RTMPServer) handlePlay(conn *rtmp.Conn) error {
	return errors.Wrap(errors.ErrNotImplemented, "play not supported on ingest server")
}

// handleConnection handles new connections
func (rs *RTMPServer) handleConnection(conn *rtmp.Conn) {
	// Connection-level handling
	rs.logger.Debug("New RTMP connection",
		zap.String("addr", conn.NetConn().RemoteAddr().String()))
}

// stopStream stops a stream and cleans up resources
func (rs *RTMPServer) stopStream(streamID string) {
	rs.mu.Lock()
	stream, exists := rs.activeStreams[streamID]
	if !exists {
		rs.mu.Unlock()
		return
	}
	delete(rs.activeStreams, streamID)
	rs.mu.Unlock()

	// Cancel stream context
	if stream.cancel != nil {
		stream.cancel()
	}

	// Calculate final statistics
	duration := time.Since(stream.StartTime)
	stats := StreamStats{
		Duration:        duration,
		BytesReceived:   stream.BytesReceived,
		PacketsReceived: stream.PacketsReceived,
		AverageBitrate:  int64(float64(stream.BytesReceived) * 8 / duration.Seconds() / 1000),
		DroppedFrames:   0, // TODO: Track dropped frames
	}

	// Notify stream end
	if err := rs.streamHandler.OnStreamEnd(streamID, stats); err != nil {
		rs.logger.Error("Failed to handle stream end",
			zap.String("streamId", streamID),
			zap.Error(err))
	}

	rs.logger.Info("Stream ended",
		zap.String("streamId", streamID),
		zap.Duration("duration", duration),
		zap.Int64("bytesReceived", stream.BytesReceived))
}

// GetActiveStreams returns currently active streams
func (rs *RTMPServer) GetActiveStreams() map[string]StreamMetadata {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	streams := make(map[string]StreamMetadata)
	for id, stream := range rs.activeStreams {
		streams[id] = StreamMetadata{
			StreamID:   stream.StreamID,
			VideoCodec: stream.VideoCodec,
			AudioCodec: stream.AudioCodec,
			Width:      stream.Width,
			Height:     stream.Height,
			Framerate:  stream.Framerate,
			Bitrate:    stream.Bitrate,
		}
	}
	return streams
}

// GetStreamStats returns statistics for a specific stream
func (rs *RTMPServer) GetStreamStats(streamID string) (*StreamStats, error) {
	rs.mu.RLock()
	stream, exists := rs.activeStreams[streamID]
	rs.mu.RUnlock()

	if !exists {
		return nil, errors.Wrap(errors.ErrNotFound, "stream not found")
	}

	duration := time.Since(stream.StartTime)
	return &StreamStats{
		Duration:        duration,
		BytesReceived:   stream.BytesReceived,
		PacketsReceived: stream.PacketsReceived,
		AverageBitrate:  int64(float64(stream.BytesReceived) * 8 / duration.Seconds() / 1000),
	}, nil
}

// parseStreamInfo extracts stream ID and key from URL
func (rs *RTMPServer) parseStreamInfo(path, query string) (streamID, streamKey string) {
	// Extract streamID from path (e.g., /live/streamID)
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 2 {
		streamID = parts[1]
	}

	// Extract stream key from query parameters
	params, _ := url.ParseQuery(query)
	streamKey = params.Get("key")

	return
}

// isIPAllowed checks if an IP is in the allowed list
func (rs *RTMPServer) isIPAllowed(ip string, allowedIPs []string) bool {
	for _, allowed := range allowedIPs {
		if ip == allowed {
			return true
		}
		// Check CIDR ranges
		if strings.Contains(allowed, "/") {
			_, cidr, err := net.ParseCIDR(allowed)
			if err == nil && cidr.Contains(net.ParseIP(ip)) {
				return true
			}
		}
	}
	return false
}

// createTLSListener creates a TLS listener for RTMPS
func (rs *RTMPServer) createTLSListener() (net.Listener, error) {
	cert, err := tls.LoadX509KeyPair(rs.config.CertFile, rs.config.KeyFile)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
	}

	return tls.Listen("tcp", rs.config.ListenAddr, tlsConfig)
}