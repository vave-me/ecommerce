package streaming

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
	"github.com/stackus/errors"
	"go.uber.org/zap"
)

// WebRTCServer handles WebRTC streaming
type WebRTCServer struct {
	config          *WebRTCConfig
	peerConnections map[string]*WebRTCPeer
	publishers      map[string]*WebRTCPublisher
	mu              sync.RWMutex
	upgrader        websocket.Upgrader
	mediaEngine     *webrtc.MediaEngine
	api             *webrtc.API
	logger          *zap.Logger
}

// WebRTCConfig contains WebRTC server configuration
type WebRTCConfig struct {
	ICEServers        []webrtc.ICEServer
	EnableTURN        bool
	TURNServer        string
	TURNUsername      string
	TURNPassword      string
	MaxBitrate        uint64
	EnableSimulcast   bool
	EnableABR         bool // Adaptive Bitrate
	StatsInterval     time.Duration
}

// WebRTCPeer represents a WebRTC peer connection
type WebRTCPeer struct {
	ID               string
	StreamID         string
	PeerConnection   *webrtc.PeerConnection
	DataChannel      *webrtc.DataChannel
	IsPublisher      bool
	UserID           string
	JoinedAt         time.Time
	LastPingTime     time.Time
	ConnectionState  webrtc.PeerConnectionState
	BitrateKbps      uint64
	PacketsLost      uint64
	Jitter           float64
	RTT              time.Duration
	VideoTrack       *webrtc.TrackLocalStaticRTP
	AudioTrack       *webrtc.TrackLocalStaticRTP
	mu               sync.RWMutex
}

// WebRTCPublisher represents a stream publisher
type WebRTCPublisher struct {
	StreamID        string
	PublisherID     string
	VideoTrack      *webrtc.TrackRemote
	AudioTrack      *webrtc.TrackRemote
	Subscribers     map[string]*WebRTCPeer
	MaxSubscribers  int
	CreatedAt       time.Time
	mu              sync.RWMutex
}

// SignalingMessage represents WebRTC signaling messages
type SignalingMessage struct {
	Type      string                     `json:"type"`
	StreamID  string                     `json:"streamId,omitempty"`
	PeerID    string                     `json:"peerId,omitempty"`
	SDP       *webrtc.SessionDescription `json:"sdp,omitempty"`
	Candidate *webrtc.ICECandidate       `json:"candidate,omitempty"`
	Error     string                     `json:"error,omitempty"`
}

// NewWebRTCServer creates a new WebRTC server
func NewWebRTCServer(config *WebRTCConfig, logger *zap.Logger) (*WebRTCServer, error) {
	// Create media engine with codecs
	mediaEngine := &webrtc.MediaEngine{}
	
	// Register video codecs
	if err := mediaEngine.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:     webrtc.MimeTypeH264,
			ClockRate:    90000,
			Channels:     0,
			SDPFmtpLine:  "level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=42001f",
		},
		PayloadType: 102,
	}, webrtc.RTPCodecTypeVideo); err != nil {
		return nil, err
	}

	// VP8 as fallback
	if err := mediaEngine.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:  webrtc.MimeTypeVP8,
			ClockRate: 90000,
		},
		PayloadType: 96,
	}, webrtc.RTPCodecTypeVideo); err != nil {
		return nil, err
	}

	// Register audio codec
	if err := mediaEngine.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:  webrtc.MimeTypeOpus,
			ClockRate: 48000,
			Channels:  2,
		},
		PayloadType: 111,
	}, webrtc.RTPCodecTypeAudio); err != nil {
		return nil, err
	}

	// Create API with media engine
	api := webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine))

	return &WebRTCServer{
		config:          config,
		peerConnections: make(map[string]*WebRTCPeer),
		publishers:      make(map[string]*WebRTCPublisher),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return checkOrigin(r)
			},
		},
		mediaEngine: mediaEngine,
		api:         api,
		logger:      logger,
	}, nil
}

// HandleWebSocket handles WebSocket connections for signaling
func (ws *WebRTCServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Add panic recovery
	defer func() {
		if r := recover(); r != nil {
			ws.logger.Error("Panic in WebRTC handler",
				zap.Any("panic", r),
				zap.Stack("stack"))
		}
	}()
	
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer conn.Close()

	peerID := generatePeerID()
	userID := r.URL.Query().Get("userId")
	
	// Ensure cleanup on exit
	defer ws.removePeer(peerID)

	// Handle signaling messages
	for {
		var msg SignalingMessage
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}

		switch msg.Type {
		case "publish":
			ws.handlePublish(conn, peerID, userID, msg)
		case "subscribe":
			ws.handleSubscribe(conn, peerID, userID, msg)
		case "offer":
			ws.handleOffer(conn, peerID, msg)
		case "answer":
			ws.handleAnswer(conn, peerID, msg)
		case "candidate":
			ws.handleCandidate(conn, peerID, msg)
		case "stop":
			ws.handleStop(peerID)
		}
	}

	// Clean up on disconnect
	ws.removePeer(peerID)
}

// handlePublish handles publisher connections
func (ws *WebRTCServer) handlePublish(conn *websocket.Conn, peerID, userID string, msg SignalingMessage) {
	// Check if stream already has a publisher
	ws.mu.RLock()
	if _, exists := ws.publishers[msg.StreamID]; exists {
		ws.mu.RUnlock()
		conn.WriteJSON(SignalingMessage{
			Type:  "error",
			Error: "Stream already has a publisher",
		})
		return
	}
	ws.mu.RUnlock()

	// Create peer connection
	pc, err := ws.createPeerConnection()
	if err != nil {
		conn.WriteJSON(SignalingMessage{
			Type:  "error",
			Error: err.Error(),
		})
		return
	}

	peer := &WebRTCPeer{
		ID:             peerID,
		StreamID:       msg.StreamID,
		PeerConnection: pc,
		IsPublisher:    true,
		UserID:         userID,
		JoinedAt:       time.Now(),
	}

	// Handle incoming tracks
	pc.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		ws.handleIncomingTrack(peer, track, receiver)
	})

	// Handle ICE connection state changes
	pc.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		ws.handleICEStateChange(peer, state)
	})

	// Handle data channel for stats and control
	dc, err := pc.CreateDataChannel("control", nil)
	if err == nil {
		peer.DataChannel = dc
		ws.setupDataChannel(peer)
	}

	ws.mu.Lock()
	ws.peerConnections[peerID] = peer
	ws.mu.Unlock()

	// Send success response
	conn.WriteJSON(SignalingMessage{
		Type:   "publish_ready",
		PeerID: peerID,
	})
}

// handleSubscribe handles subscriber connections
func (ws *WebRTCServer) handleSubscribe(conn *websocket.Conn, peerID, userID string, msg SignalingMessage) {
	// Check if publisher exists
	ws.mu.RLock()
	publisher, exists := ws.publishers[msg.StreamID]
	ws.mu.RUnlock()

	if !exists {
		conn.WriteJSON(SignalingMessage{
			Type:  "error",
			Error: "Stream not found",
		})
		return
	}

	// Check subscriber limit
	publisher.mu.RLock()
	if len(publisher.Subscribers) >= publisher.MaxSubscribers {
		publisher.mu.RUnlock()
		conn.WriteJSON(SignalingMessage{
			Type:  "error",
			Error: "Maximum subscribers reached",
		})
		return
	}
	publisher.mu.RUnlock()

	// Create peer connection
	pc, err := ws.createPeerConnection()
	if err != nil {
		conn.WriteJSON(SignalingMessage{
			Type:  "error",
			Error: err.Error(),
		})
		return
	}

	peer := &WebRTCPeer{
		ID:             peerID,
		StreamID:       msg.StreamID,
		PeerConnection: pc,
		IsPublisher:    false,
		UserID:         userID,
		JoinedAt:       time.Now(),
	}

	// Add tracks to subscriber
	if err := ws.addTracksToSubscriber(peer, publisher); err != nil {
		conn.WriteJSON(SignalingMessage{
			Type:  "error",
			Error: err.Error(),
		})
		return
	}

	// Handle ICE connection state changes
	pc.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		ws.handleICEStateChange(peer, state)
	})

	// Add to connections and publisher's subscribers
	ws.mu.Lock()
	ws.peerConnections[peerID] = peer
	ws.mu.Unlock()

	publisher.mu.Lock()
	publisher.Subscribers[peerID] = peer
	publisher.mu.Unlock()

	// Send success response
	conn.WriteJSON(SignalingMessage{
		Type:   "subscribe_ready",
		PeerID: peerID,
	})
}

// handleIncomingTrack processes incoming media tracks from publishers
func (ws *WebRTCServer) handleIncomingTrack(peer *WebRTCPeer, track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
	// Create publisher if not exists
	ws.mu.Lock()
	publisher, exists := ws.publishers[peer.StreamID]
	if !exists {
		publisher = &WebRTCPublisher{
			StreamID:       peer.StreamID,
			PublisherID:    peer.ID,
			Subscribers:    make(map[string]*WebRTCPeer),
			MaxSubscribers: 1000,
			CreatedAt:      time.Now(),
		}
		ws.publishers[peer.StreamID] = publisher
	}
	ws.mu.Unlock()

	// Store track reference
	if track.Kind() == webrtc.RTPCodecTypeVideo {
		publisher.VideoTrack = track
	} else if track.Kind() == webrtc.RTPCodecTypeAudio {
		publisher.AudioTrack = track
	}

	// Start forwarding to subscribers
	go ws.forwardTrackToSubscribers(publisher, track)

	// Start RTCP processing
	go ws.processRTCP(receiver)
}

// forwardTrackToSubscribers forwards media to all subscribers
func (ws *WebRTCServer) forwardTrackToSubscribers(publisher *WebRTCPublisher, track *webrtc.TrackRemote) {
	// Create a local track
	localTrack, err := webrtc.NewTrackLocalStaticRTP(
		track.Codec().RTPCodecCapability,
		track.ID(),
		track.StreamID(),
	)
	if err != nil {
		return
	}

	// Read RTP packets and forward
	rtpBuf := make([]byte, 1500)
	for {
		n, _, err := track.Read(rtpBuf)
		if err != nil {
			break
		}

		// Forward to all subscribers
		publisher.mu.RLock()
		subscribers := make([]*WebRTCPeer, 0, len(publisher.Subscribers))
		for _, sub := range publisher.Subscribers {
			subscribers = append(subscribers, sub)
		}
		publisher.mu.RUnlock()

		for _, subscriber := range subscribers {
			if track.Kind() == webrtc.RTPCodecTypeVideo && subscriber.VideoTrack != nil {
				subscriber.VideoTrack.Write(rtpBuf[:n])
			} else if track.Kind() == webrtc.RTPCodecTypeAudio && subscriber.AudioTrack != nil {
				subscriber.AudioTrack.Write(rtpBuf[:n])
			}
		}
	}
}

// processRTCP handles RTCP packets for statistics
func (ws *WebRTCServer) processRTCP(receiver *webrtc.RTPReceiver) {
	rtcpBuf := make([]byte, 1500)
	for {
		n, _, err := receiver.Read(rtcpBuf)
		if err != nil {
			break
		}

		// Parse RTCP packets
		pkts, err := rtcp.Unmarshal(rtcpBuf[:n])
		if err != nil {
			continue
		}

		for _, pkt := range pkts {
			switch p := pkt.(type) {
			case *rtcp.ReceiverReport:
				// Process receiver reports for quality metrics
				ws.processReceiverReport(p)
			case *rtcp.SenderReport:
				// Process sender reports
				ws.processSenderReport(p)
			}
		}
	}
}

// createPeerConnection creates a new WebRTC peer connection
func (ws *WebRTCServer) createPeerConnection() (*webrtc.PeerConnection, error) {
	config := webrtc.Configuration{
		ICEServers: ws.config.ICEServers,
	}

	// Add TURN server if enabled
	if ws.config.EnableTURN && ws.config.TURNServer != "" {
		config.ICEServers = append(config.ICEServers, webrtc.ICEServer{
			URLs:       []string{ws.config.TURNServer},
			Username:   ws.config.TURNUsername,
			Credential: ws.config.TURNPassword,
		})
	}

	return ws.api.NewPeerConnection(config)
}

// addTracksToSubscriber adds media tracks to a subscriber peer connection
func (ws *WebRTCServer) addTracksToSubscriber(peer *WebRTCPeer, publisher *WebRTCPublisher) error {
	// Create local tracks for forwarding
	if publisher.VideoTrack != nil {
		videoTrack, err := webrtc.NewTrackLocalStaticRTP(
			publisher.VideoTrack.Codec().RTPCodecCapability,
			"video",
			peer.StreamID,
		)
		if err != nil {
			return err
		}
		
		if _, err := peer.PeerConnection.AddTrack(videoTrack); err != nil {
			return err
		}
		peer.VideoTrack = videoTrack
	}

	if publisher.AudioTrack != nil {
		audioTrack, err := webrtc.NewTrackLocalStaticRTP(
			publisher.AudioTrack.Codec().RTPCodecCapability,
			"audio",
			peer.StreamID,
		)
		if err != nil {
			return err
		}
		
		if _, err := peer.PeerConnection.AddTrack(audioTrack); err != nil {
			return err
		}
		peer.AudioTrack = audioTrack
	}

	return nil
}

// setupDataChannel configures the data channel for control messages
func (ws *WebRTCServer) setupDataChannel(peer *WebRTCPeer) {
	dc := peer.DataChannel

	dc.OnOpen(func() {
		// Start sending periodic stats
		go ws.sendPeerStats(peer)
	})

	dc.OnMessage(func(msg webrtc.DataChannelMessage) {
		// Handle control messages from peer
		var ctrl map[string]interface{}
		if err := json.Unmarshal(msg.Data, &ctrl); err != nil {
			return
		}

		switch ctrl["type"] {
		case "request_keyframe":
			ws.requestKeyframe(peer)
		case "set_bitrate":
			if bitrate, ok := ctrl["bitrate"].(float64); ok {
				ws.setBitrate(peer, uint64(bitrate))
			}
		case "ping":
			peer.LastPingTime = time.Now()
			dc.SendText(`{"type":"pong"}`)
		}
	})
}

// sendPeerStats sends periodic statistics to peer
func (ws *WebRTCServer) sendPeerStats(peer *WebRTCPeer) {
	ticker := time.NewTicker(ws.config.StatsInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stats := peer.PeerConnection.GetStats()
			
			// Process and send relevant stats
			statsData := map[string]interface{}{
				"type":        "stats",
				"timestamp":   time.Now().Unix(),
				"bitrate":     peer.BitrateKbps,
				"packetsLost": peer.PacketsLost,
				"jitter":      peer.Jitter,
				"rtt":         peer.RTT.Milliseconds(),
			}

			data, _ := json.Marshal(statsData)
			if peer.DataChannel != nil && peer.DataChannel.ReadyState() == webrtc.DataChannelStateOpen {
				peer.DataChannel.Send(data)
			}
		}
	}
}

// handleICEStateChange handles ICE connection state changes
func (ws *WebRTCServer) handleICEStateChange(peer *WebRTCPeer, state webrtc.ICEConnectionState) {
	peer.mu.Lock()
	peer.ConnectionState = webrtc.PeerConnectionState(state)
	peer.mu.Unlock()

	switch state {
	case webrtc.ICEConnectionStateConnected:
		fmt.Printf("Peer %s connected\n", peer.ID)
	case webrtc.ICEConnectionStateDisconnected:
		fmt.Printf("Peer %s disconnected\n", peer.ID)
		// Start reconnection timer
		go ws.handleDisconnection(peer)
	case webrtc.ICEConnectionStateFailed:
		fmt.Printf("Peer %s connection failed\n", peer.ID)
		ws.removePeer(peer.ID)
	}
}

// handleDisconnection handles temporary disconnections
func (ws *WebRTCServer) handleDisconnection(peer *WebRTCPeer) {
	// Wait for reconnection
	time.Sleep(10 * time.Second)

	peer.mu.RLock()
	state := peer.ConnectionState
	peer.mu.RUnlock()

	if state == webrtc.PeerConnectionStateDisconnected {
		// Still disconnected, remove peer
		ws.removePeer(peer.ID)
	}
}

// removePeer removes a peer and cleans up resources
func (ws *WebRTCServer) removePeer(peerID string) {
	ws.mu.Lock()
	peer, exists := ws.peerConnections[peerID]
	if !exists {
		ws.mu.Unlock()
		return
	}
	delete(ws.peerConnections, peerID)
	ws.mu.Unlock()

	// Close peer connection
	if peer.PeerConnection != nil {
		peer.PeerConnection.Close()
	}

	// Remove from publisher if subscriber
	if !peer.IsPublisher {
		ws.mu.RLock()
		publisher, exists := ws.publishers[peer.StreamID]
		ws.mu.RUnlock()

		if exists {
			publisher.mu.Lock()
			delete(publisher.Subscribers, peerID)
			publisher.mu.Unlock()
		}
	} else {
		// Remove publisher
		ws.mu.Lock()
		delete(ws.publishers, peer.StreamID)
		ws.mu.Unlock()
	}
}

// Helper functions

func (ws *WebRTCServer) handleOffer(conn *websocket.Conn, peerID string, msg SignalingMessage) {
	ws.mu.RLock()
	peer, exists := ws.peerConnections[peerID]
	ws.mu.RUnlock()

	if !exists {
		return
	}

	if err := peer.PeerConnection.SetRemoteDescription(*msg.SDP); err != nil {
		return
	}

	answer, err := peer.PeerConnection.CreateAnswer(nil)
	if err != nil {
		return
	}

	if err := peer.PeerConnection.SetLocalDescription(answer); err != nil {
		return
	}

	conn.WriteJSON(SignalingMessage{
		Type: "answer",
		SDP:  &answer,
	})
}

func (ws *WebRTCServer) handleAnswer(conn *websocket.Conn, peerID string, msg SignalingMessage) {
	ws.mu.RLock()
	peer, exists := ws.peerConnections[peerID]
	ws.mu.RUnlock()

	if !exists {
		return
	}

	peer.PeerConnection.SetRemoteDescription(*msg.SDP)
}

func (ws *WebRTCServer) handleCandidate(conn *websocket.Conn, peerID string, msg SignalingMessage) {
	ws.mu.RLock()
	peer, exists := ws.peerConnections[peerID]
	ws.mu.RUnlock()

	if !exists {
		return
	}

	peer.PeerConnection.AddICECandidate(msg.Candidate.ToJSON())
}

func (ws *WebRTCServer) handleStop(peerID string) {
	ws.removePeer(peerID)
}

func (ws *WebRTCServer) processReceiverReport(report *rtcp.ReceiverReport) {
	// Extract quality metrics from receiver report
}

func (ws *WebRTCServer) processSenderReport(report *rtcp.SenderReport) {
	// Extract timing information from sender report
}

func (ws *WebRTCServer) requestKeyframe(peer *WebRTCPeer) {
	// Send PLI (Picture Loss Indication) to request keyframe
}

func (ws *WebRTCServer) setBitrate(peer *WebRTCPeer, bitrate uint64) {
	peer.mu.Lock()
	peer.BitrateKbps = bitrate / 1000
	peer.mu.Unlock()
}

func generatePeerID() string {
	return fmt.Sprintf("peer_%d", time.Now().UnixNano())
}

// checkOrigin validates WebSocket origin for security
func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return false
	}

	// Parse origin URL
	originURL, err := url.Parse(origin)
	if err != nil {
		return false
	}

	// Check against whitelist of allowed origins
	allowedOrigins := []string{
		"https://app.middleman.com",
		"https://middleman.com",
		"https://staging.middleman.com",
		"https://streams.middleman.com",
	}
	
	// In development, allow localhost
	if os.Getenv("ENVIRONMENT") == "development" {
		allowedOrigins = append(allowedOrigins,
			"http://localhost:3000",
			"http://localhost:3001",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:3001",
		)
	}

	// Check if origin is in allowed list
	for _, allowed := range allowedOrigins {
		allowedURL, err := url.Parse(allowed)
		if err != nil {
			continue
		}

		// Match protocol, host, and port
		if originURL.Scheme == allowedURL.Scheme &&
			originURL.Host == allowedURL.Host {
			return true
		}
	}

	// Check if it's same origin
	if originURL.Host == r.Host {
		return true
	}

	return false
}