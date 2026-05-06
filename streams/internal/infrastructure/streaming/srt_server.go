package streaming

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/haivision/srtgo"
	"github.com/stackus/errors"
	"go.uber.org/zap"
)

// SRTServer handles SRT (Secure Reliable Transport) stream ingestion
type SRTServer struct {
	config        *SRTConfig
	socket        *srtgo.SrtSocket
	activeStreams map[string]*SRTStream
	streamHandler StreamHandler
	authenticator StreamAuthenticator
	logger        *zap.Logger
	mu            sync.RWMutex
}

// SRTConfig contains SRT server configuration
type SRTConfig struct {
	ListenAddr           string
	MaxBandwidth         int64  // bits per second
	Latency              int    // milliseconds
	PeerLatency          int    // milliseconds
	RecvBuffer           int    // bytes
	SendBuffer           int    // bytes
	PayloadSize          int    // bytes
	PassPhrase           string // encryption
	KeyLength            int    // 16, 24, or 32 bytes
	StreamIDRequired     bool
	MaxConcurrentStreams int
	StatsInterval        time.Duration
}

// SRTStream represents an active SRT stream
type SRTStream struct {
	StreamID        string
	StreamKey       string
	Socket          *srtgo.SrtSocket
	PublisherAddr   string
	StartTime       time.Time
	BytesReceived   int64
	PacketsReceived int64
	PacketsLost     int64
	RTT             int64 // microseconds
	Bandwidth       int64 // bits per second
	LastPacketTime  time.Time
	ctx             context.Context
	cancel          context.CancelFunc
}

// NewSRTServer creates a new SRT server
func NewSRTServer(config *SRTConfig, handler StreamHandler, auth StreamAuthenticator, logger *zap.Logger) *SRTServer {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}

	// Set default values
	if config.Latency == 0 {
		config.Latency = 120 // 120ms default latency
	}
	if config.PayloadSize == 0 {
		config.PayloadSize = 1316 // Default MPEG-TS packet size
	}

	return &SRTServer{
		config:        config,
		activeStreams: make(map[string]*SRTStream),
		streamHandler: handler,
		authenticator: auth,
		logger:        logger,
	}
}

// Start starts the SRT server
func (ss *SRTServer) Start(ctx context.Context) error {
	// Initialize SRT library
	srtgo.InitSRT()
	defer srtgo.CleanupSRT()

	// Create SRT socket
	socket := srtgo.NewSrtSocket("", "", srtgo.DefaultOptions())
	if socket == nil {
		return errors.Wrap(errors.ErrInternalServerError, "failed to create SRT socket")
	}
	ss.socket = socket

	// Configure socket options
	ss.configureSocket(socket)

	// Parse listen address
	host, port, err := net.SplitHostPort(ss.config.ListenAddr)
	if err != nil {
		return err
	}

	// Bind socket
	bindErr := socket.Bind(host, port)
	if bindErr != nil {
		return errors.Wrap(bindErr, "failed to bind SRT socket")
	}

	// Listen for connections
	listenErr := socket.Listen(ss.config.MaxConcurrentStreams)
	if listenErr != nil {
		return errors.Wrap(listenErr, "failed to listen on SRT socket")
	}

	ss.logger.Info("SRT server started",
		zap.String("address", ss.config.ListenAddr),
		zap.Int("latency", ss.config.Latency),
		zap.Bool("encrypted", ss.config.PassPhrase != ""))

	// Accept connections
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// Accept with timeout
				socket.SetListenCallback(ss.handleListenCallback)
				
				clientSocket, addr, err := socket.Accept()
				if err != nil {
					if ctx.Err() != nil {
						return
					}
					ss.logger.Error("Failed to accept SRT connection", zap.Error(err))
					continue
				}

				go ss.handleConnection(ctx, clientSocket, addr)
			}
		}
	}()

	<-ctx.Done()
	socket.Close()
	return nil
}

// configureSocket configures SRT socket options
func (ss *SRTServer) configureSocket(socket *srtgo.SrtSocket) {
	socketOptions := []srtgo.SocketOption{
		{
			Name:  "latency",
			Val:   ss.config.Latency,
			Level: srtgo.Binding,
		},
		{
			Name:  "peerlatency",
			Val:   ss.config.PeerLatency,
			Level: srtgo.Binding,
		},
		{
			Name:  "rcvbuf",
			Val:   ss.config.RecvBuffer,
			Level: srtgo.Binding,
		},
		{
			Name:  "sndbuf",
			Val:   ss.config.SendBuffer,
			Level: srtgo.Binding,
		},
		{
			Name:  "payloadsize",
			Val:   ss.config.PayloadSize,
			Level: srtgo.Binding,
		},
	}

	// Configure encryption if passphrase is set
	if ss.config.PassPhrase != "" {
		socketOptions = append(socketOptions, 
			srtgo.SocketOption{
				Name:  "passphrase",
				Val:   ss.config.PassPhrase,
				Level: srtgo.Binding,
			},
			srtgo.SocketOption{
				Name:  "pbkeylen",
				Val:   ss.config.KeyLength,
				Level: srtgo.Binding,
			},
		)
	}

	for _, opt := range socketOptions {
		socket.SetSocketOption(opt)
	}
}

// handleListenCallback handles pre-connection validation
func (ss *SRTServer) handleListenCallback(socket *srtgo.SrtSocket, version int, addr *net.UDPAddr, streamID string) bool {
	// Validate stream ID format
	if ss.config.StreamIDRequired && streamID == "" {
		ss.logger.Warn("Connection rejected: missing stream ID",
			zap.String("addr", addr.String()))
		return false
	}

	// Check concurrent stream limit
	ss.mu.RLock()
	activeCount := len(ss.activeStreams)
	ss.mu.RUnlock()

	if activeCount >= ss.config.MaxConcurrentStreams {
		ss.logger.Warn("Connection rejected: max streams reached",
			zap.Int("current", activeCount),
			zap.Int("max", ss.config.MaxConcurrentStreams))
		return false
	}

	// Parse stream ID to extract stream key
	streamInfo, streamKey := ss.parseStreamID(streamID)

	// Validate authentication
	if err := ss.authenticator.ValidateStreamKey(streamInfo, streamKey); err != nil {
		ss.logger.Warn("Authentication failed",
			zap.String("streamId", streamInfo),
			zap.String("addr", addr.String()),
			zap.Error(err))
		return false
	}

	return true
}

// handleConnection handles an established SRT connection
func (ss *SRTServer) handleConnection(ctx context.Context, socket *srtgo.SrtSocket, addr *net.UDPAddr) {
	defer socket.Close()

	// Get stream ID
	streamID, err := socket.GetSockOptString(srtgo.StreamID)
	if err != nil {
		ss.logger.Error("Failed to get stream ID", zap.Error(err))
		return
	}

	streamInfo, streamKey := ss.parseStreamID(streamID)

	ss.logger.Info("SRT publisher connected",
		zap.String("streamId", streamInfo),
		zap.String("addr", addr.String()))

	// Create stream context
	streamCtx, cancel := context.WithCancel(ctx)

	// Create SRT stream
	stream := &SRTStream{
		StreamID:      streamInfo,
		StreamKey:     streamKey,
		Socket:        socket,
		PublisherAddr: addr.String(),
		StartTime:     time.Now(),
		ctx:           streamCtx,
		cancel:        cancel,
	}

	// Register stream
	ss.mu.Lock()
	ss.activeStreams[streamInfo] = stream
	ss.mu.Unlock()

	// Start statistics collection
	go ss.collectStats(stream)

	// Handle stream
	defer func() {
		ss.stopStream(streamInfo)
	}()

	// Notify stream start
	metadata := StreamMetadata{
		StreamID: streamInfo,
		// SRT doesn't provide codec info upfront
	}

	if err := ss.streamHandler.OnStreamStart(streamInfo, metadata); err != nil {
		ss.logger.Error("Failed to start stream handling",
			zap.String("streamId", streamInfo),
			zap.Error(err))
		return
	}

	// Create SRT demuxer adapter
	demuxer := &SRTDemuxer{
		socket: socket,
		logger: ss.logger,
	}

	// Handle stream data
	if err := ss.streamHandler.HandleStream(streamInfo, demuxer); err != nil {
		ss.logger.Error("Stream handling error",
			zap.String("streamId", streamInfo),
			zap.Error(err))
	}
}

// collectStats collects statistics for a stream
func (ss *SRTServer) collectStats(stream *SRTStream) {
	ticker := time.NewTicker(ss.config.StatsInterval)
	defer ticker.Stop()

	for {
		select {
		case <-stream.ctx.Done():
			return
		case <-ticker.C:
			stats, err := stream.Socket.Stats()
			if err != nil {
				continue
			}

			// Update stream statistics
			stream.BytesReceived = stats.ByteRecv
			stream.PacketsReceived = stats.PktRecv
			stream.PacketsLost = stats.PktRcvLoss
			stream.RTT = stats.MsRTT * 1000 // Convert to microseconds
			stream.Bandwidth = stats.MbpsRecvRate * 1000000 // Convert to bps

			// Log statistics
			ss.logger.Debug("SRT stream statistics",
				zap.String("streamId", stream.StreamID),
				zap.Int64("bytesReceived", stream.BytesReceived),
				zap.Int64("packetsLost", stream.PacketsLost),
				zap.Int64("rtt", stream.RTT),
				zap.Int64("bandwidth", stream.Bandwidth))
		}
	}
}

// stopStream stops a stream and cleans up resources
func (ss *SRTServer) stopStream(streamID string) {
	ss.mu.Lock()
	stream, exists := ss.activeStreams[streamID]
	if !exists {
		ss.mu.Unlock()
		return
	}
	delete(ss.activeStreams, streamID)
	ss.mu.Unlock()

	// Cancel stream context
	if stream.cancel != nil {
		stream.cancel()
	}

	// Close socket
	if stream.Socket != nil {
		stream.Socket.Close()
	}

	// Calculate final statistics
	duration := time.Since(stream.StartTime)
	stats := StreamStats{
		Duration:        duration,
		BytesReceived:   stream.BytesReceived,
		PacketsReceived: stream.PacketsReceived,
		AverageBitrate:  int64(float64(stream.BytesReceived) * 8 / duration.Seconds()),
		DroppedFrames:   stream.PacketsLost,
	}

	// Notify stream end
	if err := ss.streamHandler.OnStreamEnd(streamID, stats); err != nil {
		ss.logger.Error("Failed to handle stream end",
			zap.String("streamId", streamID),
			zap.Error(err))
	}

	ss.logger.Info("SRT stream ended",
		zap.String("streamId", streamID),
		zap.Duration("duration", duration),
		zap.Int64("bytesReceived", stream.BytesReceived),
		zap.Int64("packetsLost", stream.PacketsLost))
}

// parseStreamID extracts stream info and key from SRT stream ID
// Format: "streamId:streamKey" or just "streamId"
func (ss *SRTServer) parseStreamID(streamID string) (info, key string) {
	parts := strings.SplitN(streamID, ":", 2)
	info = parts[0]
	if len(parts) > 1 {
		key = parts[1]
	}
	return
}

// GetActiveStreams returns currently active SRT streams
func (ss *SRTServer) GetActiveStreams() map[string]StreamMetadata {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	streams := make(map[string]StreamMetadata)
	for id, stream := range ss.activeStreams {
		streams[id] = StreamMetadata{
			StreamID: stream.StreamID,
			Bitrate:  stream.Bandwidth / 1000, // Convert to kbps
		}
	}
	return streams
}

// GetStreamStats returns statistics for a specific stream
func (ss *SRTServer) GetStreamStats(streamID string) (*StreamStats, error) {
	ss.mu.RLock()
	stream, exists := ss.activeStreams[streamID]
	ss.mu.RUnlock()

	if !exists {
		return nil, errors.Wrap(errors.ErrNotFound, "stream not found")
	}

	duration := time.Since(stream.StartTime)
	return &StreamStats{
		Duration:        duration,
		BytesReceived:   stream.BytesReceived,
		PacketsReceived: stream.PacketsReceived,
		AverageBitrate:  stream.Bandwidth / 1000, // kbps
		DroppedFrames:   stream.PacketsLost,
	}, nil
}

// SRTDemuxer adapts SRT socket to av.Demuxer interface
type SRTDemuxer struct {
	socket *srtgo.SrtSocket
	logger *zap.Logger
	buffer []byte
}

// Streams returns codec data (not implemented for raw SRT)
func (d *SRTDemuxer) Streams() ([]av.CodecData, error) {
	// SRT doesn't provide codec information
	// This would need to be detected from the payload
	return nil, nil
}

// ReadPacket reads a packet from the SRT stream
func (d *SRTDemuxer) ReadPacket() (av.Packet, error) {
	if d.buffer == nil {
		d.buffer = make([]byte, 1316) // Default MPEG-TS packet size
	}

	n, err := d.socket.Read(d.buffer)
	if err != nil {
		return av.Packet{}, err
	}

	// Create packet
	// Note: This would need proper MPEG-TS demuxing in production
	pkt := av.Packet{
		Data: d.buffer[:n],
		Time: time.Now().UnixNano() / int64(time.Millisecond),
	}

	return pkt, nil
}

// Close closes the demuxer
func (d *SRTDemuxer) Close() error {
	return nil
}