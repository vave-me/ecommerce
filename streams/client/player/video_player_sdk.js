/**
 * DAZN-like Video Player SDK for Live Football Streaming
 * Supports HLS, DASH, and WebRTC with DRM
 */

import Hls from 'hls.js';
import dashjs from 'dashjs';
import { EventEmitter } from 'events';

class VideoPlayerSDK extends EventEmitter {
  constructor(options = {}) {
    super();
    
    this.options = {
      container: null,
      autoplay: true,
      muted: false,
      controls: true,
      preferredProtocol: 'auto', // 'hls', 'dash', 'webrtc', 'auto'
      lowLatencyMode: false,
      adaptiveBitrate: true,
      maxBufferLength: 30,
      maxMaxBufferLength: 600,
      bufferForPlayback: 1.5,
      bufferForPlaybackAfterRebuffer: 3,
      drmConfig: {},
      webrtcConfig: {
        iceServers: [
          { urls: 'stun:stun.l.google.com:19302' }
        ]
      },
      statsInterval: 1000,
      debug: false,
      ...options
    };

    this.videoElement = null;
    this.player = null;
    this.protocol = null;
    this.stats = {
      bitrate: 0,
      bufferedDuration: 0,
      droppedFrames: 0,
      latency: 0,
      downloadSpeed: 0,
      qualityLevel: -1,
      qualityLevels: []
    };
    this.isLive = true;
    this.webrtcPeer = null;
    this.statsTimer = null;
  }

  /**
   * Initialize the player with a container element
   */
  init(container) {
    if (typeof container === 'string') {
      this.options.container = document.querySelector(container);
    } else {
      this.options.container = container;
    }

    if (!this.options.container) {
      throw new Error('Container element not found');
    }

    this.createVideoElement();
    this.setupEventListeners();
    this.startStatsCollection();

    return this;
  }

  /**
   * Load and play a stream
   */
  async load(streamUrl, options = {}) {
    const config = { ...this.options, ...options };
    
    // Detect protocol from URL or use specified
    this.protocol = this.detectProtocol(streamUrl, config.preferredProtocol);

    try {
      switch (this.protocol) {
        case 'hls':
          await this.loadHLS(streamUrl, config);
          break;
        case 'dash':
          await this.loadDASH(streamUrl, config);
          break;
        case 'webrtc':
          await this.loadWebRTC(streamUrl, config);
          break;
        default:
          throw new Error('Unsupported protocol');
      }

      this.emit('loadstart', { protocol: this.protocol, url: streamUrl });

      if (config.autoplay) {
        await this.play();
      }
    } catch (error) {
      this.emit('error', error);
      throw error;
    }
  }

  /**
   * Load HLS stream
   */
  async loadHLS(url, config) {
    if (!Hls.isSupported()) {
      // Fallback to native HLS if supported
      if (this.videoElement.canPlayType('application/vnd.apple.mpegurl')) {
        this.videoElement.src = url;
        return;
      }
      throw new Error('HLS is not supported');
    }

    // Clean up existing player
    if (this.player) {
      this.player.destroy();
    }

    const hlsConfig = {
      debug: config.debug,
      enableWorker: true,
      lowLatencyMode: config.lowLatencyMode,
      backBufferLength: 30,
      maxBufferLength: config.maxBufferLength,
      maxMaxBufferLength: config.maxMaxBufferLength,
      maxBufferSize: 60 * 1000 * 1000, // 60 MB
      maxBufferHole: 0.5,
      highBufferWatchdogPeriod: 2,
      nudgeOffset: 0.1,
      nudgeMaxRetry: 10,
      maxFragLookUpTolerance: 0.25,
      enableSoftwareAES: false,
      startLevel: -1, // Auto
      fragLoadingTimeOut: 20000,
      fragLoadingMaxRetry: 6,
      fragLoadingRetryDelay: 1000,
      fragLoadingMaxRetryTimeout: 64000,
      ...config.hlsConfig
    };

    // Configure DRM if needed
    if (config.drmConfig && config.drmConfig.widevine) {
      hlsConfig.emeEnabled = true;
      hlsConfig.widevineLicenseUrl = config.drmConfig.widevine.licenseUrl;
    }

    this.player = new Hls(hlsConfig);
    
    // Set up HLS event handlers
    this.setupHLSEvents();
    
    this.player.loadSource(url);
    this.player.attachMedia(this.videoElement);
  }

  /**
   * Load DASH stream
   */
  async loadDASH(url, config) {
    // Clean up existing player
    if (this.player) {
      this.player.reset();
    }

    this.player = dashjs.MediaPlayer().create();

    const dashConfig = {
      debug: {
        logLevel: config.debug ? dashjs.Debug.LOG_LEVEL_DEBUG : dashjs.Debug.LOG_LEVEL_NONE
      },
      streaming: {
        lowLatencyEnabled: config.lowLatencyMode,
        liveDelay: config.lowLatencyMode ? 3 : 12,
        liveCatchUpMinDrift: 0.05,
        liveCatchUpPlaybackRate: 0.5,
        abr: {
          autoSwitchBitrate: {
            video: config.adaptiveBitrate,
            audio: config.adaptiveBitrate
          }
        },
        buffer: {
          bufferTimeAtTopQuality: config.maxBufferLength,
          bufferTimeAtTopQualityLongForm: config.maxMaxBufferLength,
          initialBufferLevel: config.bufferForPlayback,
          bufferForPlaybackAfterSeek: config.bufferForPlaybackAfterRebuffer
        }
      }
    };

    // Configure DRM
    if (config.drmConfig) {
      const protectionData = {};
      
      if (config.drmConfig.widevine) {
        protectionData['com.widevine.alpha'] = {
          serverURL: config.drmConfig.widevine.licenseUrl,
          httpRequestHeaders: config.drmConfig.widevine.headers || {}
        };
      }
      
      if (config.drmConfig.playready) {
        protectionData['com.microsoft.playready'] = {
          serverURL: config.drmConfig.playready.licenseUrl,
          httpRequestHeaders: config.drmConfig.playready.headers || {}
        };
      }

      this.player.setProtectionData(protectionData);
    }

    this.player.updateSettings(dashConfig);
    this.player.initialize(this.videoElement, url, config.autoplay);
    
    // Set up DASH event handlers
    this.setupDASHEvents();
  }

  /**
   * Load WebRTC stream for ultra-low latency
   */
  async loadWebRTC(streamId, config) {
    // Clean up existing connection
    if (this.webrtcPeer) {
      this.webrtcPeer.close();
    }

    // Create WebSocket connection for signaling
    const wsUrl = config.webrtcSignalingUrl || `wss://${window.location.host}/ws`;
    const ws = new WebSocket(`${wsUrl}?streamId=${streamId}&userId=${config.userId || 'anonymous'}`);

    ws.onopen = () => {
      this.emit('webrtc:connected');
      
      // Send subscribe message
      ws.send(JSON.stringify({
        type: 'subscribe',
        streamId: streamId
      }));
    };

    ws.onmessage = async (event) => {
      const msg = JSON.parse(event.data);
      
      switch (msg.type) {
        case 'subscribe_ready':
          await this.createWebRTCPeer(ws, config);
          break;
        case 'offer':
          await this.handleWebRTCOffer(ws, msg.sdp);
          break;
        case 'candidate':
          await this.handleWebRTCCandidate(msg.candidate);
          break;
        case 'error':
          this.emit('error', new Error(msg.error));
          break;
      }
    };

    ws.onerror = (error) => {
      this.emit('error', error);
    };

    ws.onclose = () => {
      this.emit('webrtc:disconnected');
    };
  }

  /**
   * Create WebRTC peer connection
   */
  async createWebRTCPeer(ws, config) {
    const pcConfig = {
      iceServers: config.webrtcConfig.iceServers,
      iceCandidatePoolSize: 10
    };

    this.webrtcPeer = new RTCPeerConnection(pcConfig);

    // Handle incoming stream
    this.webrtcPeer.ontrack = (event) => {
      if (event.track.kind === 'video') {
        this.videoElement.srcObject = event.streams[0];
        this.emit('webrtc:trackadded', event.track);
      }
    };

    // Handle ICE candidates
    this.webrtcPeer.onicecandidate = (event) => {
      if (event.candidate) {
        ws.send(JSON.stringify({
          type: 'candidate',
          candidate: event.candidate
        }));
      }
    };

    // Handle connection state changes
    this.webrtcPeer.onconnectionstatechange = () => {
      this.emit('webrtc:connectionstate', this.webrtcPeer.connectionState);
    };

    // Create offer
    const offer = await this.webrtcPeer.createOffer({
      offerToReceiveVideo: true,
      offerToReceiveAudio: true
    });

    await this.webrtcPeer.setLocalDescription(offer);

    ws.send(JSON.stringify({
      type: 'offer',
      sdp: offer
    }));
  }

  /**
   * Setup HLS event handlers
   */
  setupHLSEvents() {
    this.player.on(Hls.Events.MANIFEST_PARSED, (event, data) => {
      this.emit('levelsloaded', data.levels);
      this.stats.qualityLevels = data.levels.map(level => ({
        bitrate: level.bitrate,
        width: level.width,
        height: level.height,
        framerate: level.attrs.FRAMERATE
      }));
    });

    this.player.on(Hls.Events.LEVEL_SWITCHED, (event, data) => {
      this.stats.qualityLevel = data.level;
      this.emit('qualitychanged', {
        level: data.level,
        bitrate: this.player.levels[data.level].bitrate
      });
    });

    this.player.on(Hls.Events.FRAG_BUFFERED, (event, data) => {
      this.emit('fragmentloaded', data);
    });

    this.player.on(Hls.Events.ERROR, (event, data) => {
      if (data.fatal) {
        switch (data.type) {
          case Hls.ErrorTypes.NETWORK_ERROR:
            this.handleNetworkError(data);
            break;
          case Hls.ErrorTypes.MEDIA_ERROR:
            this.handleMediaError(data);
            break;
          default:
            this.emit('error', new Error(data.details));
            break;
        }
      }
    });
  }

  /**
   * Setup DASH event handlers
   */
  setupDASHEvents() {
    this.player.on(dashjs.MediaPlayer.events.STREAM_INITIALIZED, () => {
      this.emit('streaminitialized');
    });

    this.player.on(dashjs.MediaPlayer.events.QUALITY_CHANGE_RENDERED, (e) => {
      this.stats.qualityLevel = e.newQuality;
      this.emit('qualitychanged', {
        level: e.newQuality,
        mediaType: e.mediaType
      });
    });

    this.player.on(dashjs.MediaPlayer.events.BUFFER_LOADED, () => {
      this.emit('bufferloaded');
    });

    this.player.on(dashjs.MediaPlayer.events.ERROR, (e) => {
      this.emit('error', e.error);
    });

    this.player.on(dashjs.MediaPlayer.events.METRIC_CHANGED, (e) => {
      if (e.metric === 'BufferLevel') {
        this.stats.bufferedDuration = e.value.level;
      }
    });
  }

  /**
   * Handle network errors with retry logic
   */
  handleNetworkError(data) {
    console.error('Network error:', data);
    
    if (data.details === Hls.ErrorDetails.MANIFEST_LOAD_ERROR ||
        data.details === Hls.ErrorDetails.MANIFEST_LOAD_TIMEOUT) {
      // Retry loading manifest
      setTimeout(() => {
        this.player.startLoad();
      }, 1000);
    } else if (data.details === Hls.ErrorDetails.FRAG_LOAD_ERROR ||
               data.details === Hls.ErrorDetails.FRAG_LOAD_TIMEOUT) {
      // Skip problematic fragment
      this.player.startLoad();
    }
  }

  /**
   * Handle media errors
   */
  handleMediaError(data) {
    console.error('Media error:', data);
    
    // Try to recover
    this.player.recoverMediaError();
  }

  /**
   * Play the video
   */
  async play() {
    try {
      await this.videoElement.play();
      this.emit('play');
    } catch (error) {
      this.emit('error', error);
      throw error;
    }
  }

  /**
   * Pause the video
   */
  pause() {
    this.videoElement.pause();
    this.emit('pause');
  }

  /**
   * Seek to a specific time (for DVR)
   */
  seek(time) {
    if (this.protocol === 'webrtc') {
      console.warn('Seeking not supported in WebRTC mode');
      return;
    }
    
    this.videoElement.currentTime = time;
    this.emit('seek', time);
  }

  /**
   * Set quality level
   */
  setQuality(level) {
    if (this.protocol === 'hls' && this.player) {
      if (level === -1) {
        this.player.currentLevel = -1; // Auto
      } else {
        this.player.currentLevel = level;
      }
    } else if (this.protocol === 'dash' && this.player) {
      const mediaInfo = this.player.getActiveStream().getMediaInfo('video');
      this.player.setQualityFor('video', level, mediaInfo);
    }
    
    this.emit('qualitychangerequest', level);
  }

  /**
   * Get current quality level
   */
  getQuality() {
    if (this.protocol === 'hls' && this.player) {
      return this.player.currentLevel;
    } else if (this.protocol === 'dash' && this.player) {
      return this.player.getQualityFor('video');
    }
    return -1;
  }

  /**
   * Get available quality levels
   */
  getQualityLevels() {
    return this.stats.qualityLevels;
  }

  /**
   * Set volume
   */
  setVolume(volume) {
    this.videoElement.volume = Math.max(0, Math.min(1, volume));
    this.emit('volumechange', this.videoElement.volume);
  }

  /**
   * Get current stats
   */
  getStats() {
    return { ...this.stats };
  }

  /**
   * Start collecting statistics
   */
  startStatsCollection() {
    this.statsTimer = setInterval(() => {
      this.updateStats();
    }, this.options.statsInterval);
  }

  /**
   * Update statistics
   */
  updateStats() {
    // Get buffer information
    const buffered = this.videoElement.buffered;
    if (buffered.length > 0) {
      this.stats.bufferedDuration = buffered.end(buffered.length - 1) - this.videoElement.currentTime;
    }

    // Get quality metrics
    if (this.videoElement.getVideoPlaybackQuality) {
      const quality = this.videoElement.getVideoPlaybackQuality();
      this.stats.droppedFrames = quality.droppedVideoFrames;
    }

    // Protocol-specific stats
    if (this.protocol === 'hls' && this.player) {
      this.stats.bitrate = this.player.bandwidthEstimate || 0;
      this.stats.latency = this.player.latency || 0;
    } else if (this.protocol === 'dash' && this.player) {
      const dashMetrics = this.player.getDashMetrics();
      const streamInfo = this.player.getActiveStream();
      if (streamInfo) {
        const periodIdx = streamInfo.getStreamInfo().index;
        const repSwitch = dashMetrics.getCurrentRepresentationSwitch('video', true);
        if (repSwitch) {
          this.stats.bitrate = repSwitch.to.bandwidth;
        }
      }
    } else if (this.protocol === 'webrtc' && this.webrtcPeer) {
      // WebRTC stats would be collected via getStats API
      this.collectWebRTCStats();
    }

    this.emit('stats', this.stats);
  }

  /**
   * Collect WebRTC statistics
   */
  async collectWebRTCStats() {
    if (!this.webrtcPeer) return;

    try {
      const stats = await this.webrtcPeer.getStats();
      stats.forEach(report => {
        if (report.type === 'inbound-rtp' && report.mediaType === 'video') {
          this.stats.bitrate = report.bytesReceived * 8 / 1000; // kbps
          this.stats.droppedFrames = report.framesDropped || 0;
          this.stats.latency = report.jitter || 0;
        }
      });
    } catch (error) {
      console.error('Error collecting WebRTC stats:', error);
    }
  }

  /**
   * Create video element
   */
  createVideoElement() {
    this.videoElement = document.createElement('video');
    this.videoElement.controls = this.options.controls;
    this.videoElement.muted = this.options.muted;
    this.videoElement.autoplay = this.options.autoplay;
    this.videoElement.playsInline = true;
    this.videoElement.style.width = '100%';
    this.videoElement.style.height = '100%';
    
    this.options.container.appendChild(this.videoElement);
  }

  /**
   * Setup video element event listeners
   */
  setupEventListeners() {
    const events = [
      'loadstart', 'loadedmetadata', 'loadeddata', 'canplay', 'canplaythrough',
      'play', 'pause', 'playing', 'waiting', 'seeking', 'seeked', 'ended',
      'durationchange', 'timeupdate', 'progress', 'volumechange', 'ratechange',
      'error', 'stalled', 'suspend'
    ];

    events.forEach(event => {
      this.videoElement.addEventListener(event, (e) => {
        this.emit(event, e);
      });
    });
  }

  /**
   * Detect streaming protocol from URL
   */
  detectProtocol(url, preferred) {
    if (preferred !== 'auto') {
      return preferred;
    }

    if (url.includes('.m3u8') || url.includes('/hls/')) {
      return 'hls';
    } else if (url.includes('.mpd') || url.includes('/dash/')) {
      return 'dash';
    } else if (url.startsWith('webrtc:') || url.includes('/webrtc/')) {
      return 'webrtc';
    }

    // Default to HLS
    return 'hls';
  }

  /**
   * Destroy the player and clean up
   */
  destroy() {
    if (this.statsTimer) {
      clearInterval(this.statsTimer);
    }

    if (this.player) {
      if (this.protocol === 'hls') {
        this.player.destroy();
      } else if (this.protocol === 'dash') {
        this.player.reset();
      }
    }

    if (this.webrtcPeer) {
      this.webrtcPeer.close();
    }

    if (this.videoElement) {
      this.videoElement.remove();
    }

    this.removeAllListeners();
  }
}

// Export for different module systems
if (typeof module !== 'undefined' && module.exports) {
  module.exports = VideoPlayerSDK;
} else if (typeof define === 'function' && define.amd) {
  define([], () => VideoPlayerSDK);
} else {
  window.VideoPlayerSDK = VideoPlayerSDK;
}