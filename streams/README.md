# Live Streaming Service

A production-ready live streaming service designed for DAZN-like football streaming with support for HLS, DASH, and WebRTC protocols.

## Features

### Streaming Protocols
- **HLS (HTTP Live Streaming)**: Apple's adaptive bitrate streaming protocol
- **DASH (Dynamic Adaptive Streaming)**: MPEG standard for adaptive streaming
- **WebRTC**: Ultra-low latency streaming (<1 second delay)

### Video Quality
- Multiple quality profiles (360p, 480p, 720p, 1080p, 4K)
- Adaptive bitrate streaming
- Hardware-accelerated transcoding (NVIDIA, Intel QSV, VAAPI)
- H.264/H.265 codec support

### CDN Integration
- Multi-CDN support (CloudFlare, Akamai, Fastly, AWS CloudFront)
- Automatic failover between CDN providers
- Edge server distribution
- Cache management

### DRM (Digital Rights Management)
- Widevine (Google)
- FairPlay (Apple)
- PlayReady (Microsoft)
- ClearKey
- Token-based authentication
- Geo-blocking

### Low Latency Features
- WebRTC for sub-second latency
- Low-latency HLS (LL-HLS)
- Low-latency DASH (LL-DASH)
- Configurable buffer sizes

### Analytics & Monitoring
- Real-time viewer statistics
- Quality metrics (bitrate, buffering, dropped frames)
- CDN performance monitoring
- Stream health monitoring

## Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   Broadcaster   │────▶│ Ingest Server   │────▶│   Transcoder    │
│  (RTMP/SRT)    │     │                 │     │   (FFmpeg)      │
└─────────────────┘     └─────────────────┘     └─────────────────┘
                                                          │
                                                          ▼
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   Video Player  │◀────│      CDN        │◀────│ Segment Store   │
│   (HLS/DASH)    │     │  (Multi-CDN)    │     │                 │
└─────────────────┘     └─────────────────┘     └─────────────────┘
```

## Quick Start

### 1. Start the Streaming Server

```go
// Create streaming infrastructure
streamingConfig := &streaming.StreamingConfig{
    IngestPort:       1935,  // RTMP
    StreamingPort:    8080,  // HTTP
    RTMPEnabled:      true,
    SRTEnabled:       true,
    WebRTCEnabled:    true,
    SegmentDuration:  4,
    PlaylistSize:     10,
    StoragePath:      "/var/streams",
    CDNUploadEnabled: true,
    LowLatencyMode:   true,
    DVRWindowMinutes: 120,
}

// Initialize CDN manager
cdnConfigs := []streaming.CDNConfig{
    {
        Provider:     streaming.CDNProviderCloudflare,
        Endpoint:     "https://cdn.cloudflare.com",
        CustomDomain: "stream.example.com",
    },
}
cdnManager := streaming.NewCDNManager(cdnConfigs, streaming.CDNProviderCloudflare, true)

// Initialize DRM manager
drmManager := streaming.NewDRMManager()
drmManager.RegisterProvider(streaming.DRMProviderWidevine, 
    streaming.NewWidevineDRM(licenseURL, signingKey, encKey, keyStore))

// Start streaming server
streamingServer := streaming.NewStreamingServer(streamingConfig, cdnManager)
go streamingServer.Start(ctx)

// Initialize WebRTC server
webrtcConfig := &streaming.WebRTCConfig{
    ICEServers: []webrtc.ICEServer{
        {URLs: []string{"stun:stun.l.google.com:19302"}},
    },
    EnableTURN:      true,
    MaxBitrate:      8000000, // 8 Mbps
    EnableSimulcast: true,
}
webrtcServer, _ := streaming.NewWebRTCServer(webrtcConfig)
```

### 2. Create a Live Stream

```go
// Create stream command
cmd := commands.CreateLiveStream{
    ID:                 "premier-league-match-1",
    Title:              "Manchester United vs Liverpool",
    Description:        "Premier League Match Day 25",
    EventType:          "football_match",
    HomeTeam:           "Manchester United",
    AwayTeam:           "Liverpool",
    Competition:        "Premier League",
    Season:             "2023/24",
    MatchDay:           25,
    Stadium:            "Old Trafford",
    ScheduledStartTime: time.Now().Add(2 * time.Hour),
    ScheduledEndTime:   time.Now().Add(4 * time.Hour),
}

err := app.CreateLiveStream(ctx, cmd)
```

### 3. Configure Streaming

```go
// Configure streaming protocols and quality
configCmd := commands.ConfigureStreaming{
    StreamID: "premier-league-match-1",
    Protocols: []domain.StreamingProtocol{
        domain.ProtocolHLS,
        domain.ProtocolDASH,
        domain.ProtocolWebRTC,
    },
    QualityProfiles: []domain.StreamingQualityProfile{
        {Name: "1080p", Resolution: "1920x1080", Bitrate: 8000, Framerate: 60},
        {Name: "720p", Resolution: "1280x720", Bitrate: 4000, Framerate: 60},
        {Name: "480p", Resolution: "854x480", Bitrate: 2000, Framerate: 30},
    },
    AdaptiveBitrate:  true,
    LowLatencyMode:   true,
    DVREnabled:       true,
    DVRWindowMinutes: 120,
    CDNEndpoints: []domain.CDNEndpoint{
        {Provider: "Cloudflare", Region: "global", Active: true},
    },
    PrimaryCDN:     "Cloudflare",
    EnableFailover: true,
    EnableDRM:      true,
    DRMProviders: map[string]domain.DRMConfig{
        "widevine": {Provider: "Widevine", Enabled: true},
    },
}

err := app.ConfigureStreaming(ctx, configCmd)
```

### 4. Start Broadcasting

```go
// Start the live stream
startCmd := commands.StartLiveStream{
    StreamID:        "premier-league-match-1",
    IngestProtocol:  domain.ProtocolRTMP,
    IngestURL:       "rtmp://ingest.example.com/live",
    BackupIngestURL: "rtmp://backup-ingest.example.com/live",
    StreamKey:       "secret-stream-key-123",
}

err := app.StartLiveStream(ctx, startCmd)
```

### 5. Client Integration

```javascript
// Initialize video player
const player = new VideoPlayerSDK({
    autoplay: true,
    controls: true,
    lowLatencyMode: true,
    adaptiveBitrate: true,
    drmConfig: {
        widevine: {
            licenseUrl: 'https://license.example.com/widevine',
            headers: {
                'Authorization': 'Bearer token'
            }
        }
    }
});

// Initialize player with container
player.init('#video-container');

// Load stream
player.load('/hls/premier-league-match-1/master.m3u8', {
    preferredProtocol: 'hls' // or 'dash', 'webrtc'
});

// Listen to events
player.on('stats', (stats) => {
    console.log(`Bitrate: ${stats.bitrate} kbps`);
    console.log(`Buffer: ${stats.bufferedDuration}s`);
    console.log(`Latency: ${stats.latency}ms`);
});

player.on('error', (error) => {
    console.error('Playback error:', error);
});
```

## Broadcasting Tools

### OBS Studio Configuration

1. **Stream Settings**:
   - Service: Custom
   - Server: `rtmp://your-server.com/live`
   - Stream Key: Your unique stream key

2. **Output Settings**:
   - Video Bitrate: 6000-8000 Kbps
   - Audio Bitrate: 192 Kbps
   - Keyframe Interval: 2s

3. **Video Settings**:
   - Base Resolution: 1920x1080
   - Output Resolution: 1920x1080
   - FPS: 60

### FFmpeg Broadcasting

```bash
# RTMP broadcasting
ffmpeg -re -i input.mp4 \
  -c:v libx264 -preset medium -b:v 8000k \
  -c:a aac -b:a 192k \
  -f flv rtmp://your-server.com/live/stream-key

# SRT broadcasting (lower latency)
ffmpeg -re -i input.mp4 \
  -c:v libx264 -preset medium -b:v 8000k \
  -c:a aac -b:a 192k \
  -f mpegts "srt://your-server.com:9999?streamid=stream-key"
```

## API Endpoints

### Stream Management

```http
POST /api/streams
GET /api/streams/{streamId}
PUT /api/streams/{streamId}/start
PUT /api/streams/{streamId}/stop
DELETE /api/streams/{streamId}
```

### Playback URLs

```http
# HLS
GET /hls/{streamId}/master.m3u8
GET /hls/{streamId}/{quality}/playlist.m3u8
GET /hls/{streamId}/{quality}/{segment}.ts

# DASH
GET /dash/{streamId}/manifest.mpd
GET /dash/{streamId}/{quality}/{segment}.m4s

# WebRTC Signaling
WS /ws?streamId={streamId}&userId={userId}
```

### DRM License

```http
POST /api/drm/{provider}/license
{
  "contentId": "stream-id",
  "challenge": "base64-encoded-challenge"
}
```

## Performance Optimization

### Transcoding

1. **Hardware Acceleration**:
   ```bash
   # NVIDIA GPU
   ffmpeg -hwaccel cuda -hwaccel_output_format cuda ...
   
   # Intel QuickSync
   ffmpeg -hwaccel qsv ...
   ```

2. **Preset Tuning**:
   - `ultrafast`: Lowest latency, higher bandwidth
   - `medium`: Balanced latency and quality
   - `slow`: Best quality, higher latency

### CDN Configuration

1. **Edge Locations**: Deploy close to viewers
2. **Cache Headers**: Optimize TTL for segments
3. **HTTP/2 Push**: Push segments proactively
4. **Connection Pooling**: Reuse connections

### Client Optimization

1. **Buffer Management**:
   ```javascript
   player.options.maxBufferLength = 30; // seconds
   player.options.bufferForPlayback = 1.5; // seconds
   ```

2. **Quality Selection**:
   - Start with lower quality
   - Ramp up based on bandwidth
   - Consider device capabilities

## Monitoring

### Metrics to Track

- **Stream Health**:
  - Ingest stability
  - Transcoding performance
  - Segment generation rate

- **Viewer Experience**:
  - Startup time
  - Rebuffering ratio
  - Average bitrate
  - Stream latency

- **Infrastructure**:
  - CPU/GPU utilization
  - Network bandwidth
  - Storage usage
  - CDN hit ratio

### Grafana Dashboard

```sql
-- Concurrent viewers
SELECT COUNT(DISTINCT session_id) 
FROM viewer_sessions 
WHERE stream_id = $streamId 
  AND last_heartbeat > NOW() - INTERVAL '30 seconds';

-- Average bitrate
SELECT AVG(bitrate_kbps) 
FROM viewer_metrics 
WHERE stream_id = $streamId 
  AND timestamp > NOW() - INTERVAL '5 minutes';

-- Buffering events
SELECT COUNT(*) 
FROM buffering_events 
WHERE stream_id = $streamId 
  AND timestamp > NOW() - INTERVAL '1 hour';
```

## Production Checklist

- [ ] SSL/TLS certificates configured
- [ ] DRM licenses obtained
- [ ] CDN contracts in place
- [ ] Monitoring alerts configured
- [ ] Backup ingest servers ready
- [ ] Load testing completed
- [ ] Geo-blocking rules defined
- [ ] GDPR compliance verified
- [ ] Stream recording enabled
- [ ] Disaster recovery plan

## Troubleshooting

### High Latency

1. Enable low-latency mode
2. Reduce segment duration
3. Use WebRTC for critical streams
4. Optimize CDN routing

### Buffering Issues

1. Check origin bandwidth
2. Verify CDN performance
3. Analyze client logs
4. Adjust quality ladder

### DRM Errors

1. Verify license server
2. Check certificate validity
3. Confirm device compatibility
4. Review geo-restrictions

## License

This streaming service is proprietary software. All rights reserved.