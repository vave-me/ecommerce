/**
 * Production-Ready Video Player SDK for DAZN-like Live Football Streaming
 * Supports HLS, DASH, WebRTC with advanced features
 */

import Hls from 'hls.js';
import dashjs, { MediaPlayerClass } from 'dashjs';
import { EventEmitter } from 'eventemitter3';
import adapter from 'webrtc-adapter';

// Ensure cross-browser WebRTC compatibility
console.log('WebRTC adapter browserDetails:', adapter.browserDetails);

export interface PlayerOptions {
  container?: HTMLElement | string;
  autoplay?: boolean;
  muted?: boolean;
  controls?: boolean;
  preferredProtocol?: 'auto' | 'hls' | 'dash' | 'webrtc';
  lowLatencyMode?: boolean;
  adaptiveBitrate?: boolean;
  maxBufferLength?: number;
  maxMaxBufferLength?: number;
  bufferForPlayback?: number;
  bufferForPlaybackAfterRebuffer?: number;
  drmConfig?: DRMConfig;
  webrtcConfig?: WebRTCConfig;
  statsInterval?: number;
  debug?: boolean;
  retryConfig?: RetryConfig;
  networkConfig?: NetworkConfig;
}

export interface DRMConfig {
  widevine?: {
    licenseUrl: string;
    headers?: Record<string, string>;
  };
  fairplay?: {
    licenseUrl: string;
    certificateUrl: string;
    headers?: Record<string, string>;
  };
  playready?: {
    licenseUrl: string;
    headers?: Record<string, string>;
  };
}

export interface WebRTCConfig {
  signalingUrl?: string;
  iceServers?: RTCIceServer[];
  maxBitrate?: number;
  simulcast?: boolean;
}

export interface RetryConfig {
  maxRetries?: number;
  retryDelay?: number;
  backoffFactor?: number;
  maxRetryDelay?: number;
}

export interface NetworkConfig {
  timeout?: number;
  withCredentials?: boolean;
  headers?: Record<string, string>;
}

export interface StreamStats {
  bitrate: number;
  bufferedDuration: number;
  droppedFrames: number;
  latency: number;
  downloadSpeed: number;
  qualityLevel: number;
  qualityLevels: QualityLevel[];
  protocol: string;
  connectionState: string;
}

export interface QualityLevel {
  bitrate: number;
  width: number;
  height: number;
  framerate: number;
  codec: string;
}

export type PlayerEvent = 
  | 'ready' | 'play' | 'pause' | 'ended' | 'error' 
  | 'loadstart' | 'loadedmetadata' | 'canplay' | 'waiting'
  | 'stats' | 'qualitychanged' | 'bufferwarning' | 'latencywarning';

export class VideoPlayerSDK extends EventEmitter<PlayerEvent> {
  private options: Required<PlayerOptions>;
  private container: HTMLElement | null = null;
  private videoElement: HTMLVideoElement | null = null;
  private hlsInstance: Hls | null = null;
  private dashInstance: MediaPlayerClass | null = null;
  private webrtcPeer: RTCPeerConnection | null = null;
  private webSocket: WebSocket | null = null;
  private protocol: string = '';
  private stats: StreamStats;
  private statsTimer: number | null = null;
  private retryCount: number = 0;
  private isDestroyed: boolean = false;
  private qualityController: QualityController;
  private networkMonitor: NetworkMonitor;
  private errorRecovery: ErrorRecovery;

  constructor(options: PlayerOptions = {}) {
    super();
    
    this.options = {
      container: null,
      autoplay: true,
      muted: false,
      controls: true,
      preferredProtocol: 'auto',
      lowLatencyMode: false,
      adaptiveBitrate: true,
      maxBufferLength: 30,
      maxMaxBufferLength: 600,
      bufferForPlayback: 1.5,
      bufferForPlaybackAfterRebuffer: 3,
      drmConfig: {},
      webrtcConfig: {
        iceServers: [{ urls: 'stun:stun.l.google.com:19302' }],
        maxBitrate: 8000000,
        simulcast: true,
      },
      statsInterval: 1000,
      debug: false,
      retryConfig: {
        maxRetries: 5,
        retryDelay: 1000,
        backoffFactor: 2,
        maxRetryDelay: 30000,
      },
      networkConfig: {
        timeout: 10000,
        withCredentials: false,
        headers: {},
      },
      ...options,
    } as Required<PlayerOptions>;

    this.stats = this.initializeStats();
    this.qualityController = new QualityController(this);
    this.networkMonitor = new NetworkMonitor(this);
    this.errorRecovery = new ErrorRecovery(this);
  }

  private initializeStats(): StreamStats {
    return {
      bitrate: 0,
      bufferedDuration: 0,
      droppedFrames: 0,
      latency: 0,
      downloadSpeed: 0,
      qualityLevel: -1,
      qualityLevels: [],
      protocol: '',
      connectionState: 'disconnected',
    };
  }

  public async init(container: HTMLElement | string): Promise<void> {
    if (typeof container === 'string') {
      this.container = document.querySelector(container);
    } else {
      this.container = container;
    }

    if (!this.container) {
      throw new Error('Container element not found');
    }

    this.createVideoElement();
    this.setupEventListeners();
    this.startStatsCollection();
    this.networkMonitor.start();

    this.emit('ready');
  }

  public async load(streamUrl: string, options: Partial<PlayerOptions> = {}): Promise<void> {
    if (this.isDestroyed) {
      throw new Error('Player has been destroyed');
    }

    const config = { ...this.options, ...options };
    this.protocol = this.detectProtocol(streamUrl, config.preferredProtocol);
    this.stats.protocol = this.protocol;

    try {
      this.emit('loadstart');

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
          throw new Error(`Unsupported protocol: ${this.protocol}`);
      }

      if (config.autoplay) {
        await this.play();
      }
    } catch (error) {
      this.handleError(error as Error);
      throw error;
    }
  }

  private async loadHLS(url: string, config: Required<PlayerOptions>): Promise<void> {
    if (!this.videoElement) throw new Error('Video element not initialized');

    // Check native HLS support (Safari)
    if (this.videoElement.canPlayType('application/vnd.apple.mpegurl')) {
      this.videoElement.src = url;
      await this.setupNativeHLS(config);
      return;
    }

    if (!Hls.isSupported()) {
      throw new Error('HLS is not supported in this browser');
    }

    this.destroyExistingPlayer();

    const hlsConfig: Partial<Hls.Config> = {
      debug: config.debug,
      enableWorker: true,
      lowLatencyMode: config.lowLatencyMode,
      backBufferLength: 90,
      maxBufferLength: config.maxBufferLength,
      maxMaxBufferLength: config.maxMaxBufferLength,
      maxBufferSize: 60 * 1000 * 1000,
      maxBufferHole: 0.5,
      highBufferWatchdogPeriod: 2,
      nudgeOffset: 0.1,
      nudgeMaxRetry: 10,
      maxFragLookUpTolerance: 0.25,
      enableSoftwareAES: false,
      startLevel: -1,
      testBandwidth: true,
      progressive: true,
      lowLatencyMode: config.lowLatencyMode,
    };

    // Configure low latency
    if (config.lowLatencyMode) {
      Object.assign(hlsConfig, {
        liveSyncDurationCount: 2,
        liveMaxLatencyDurationCount: 4,
        liveDurationInfinity: true,
        highBufferWatchdogPeriod: 1,
      });
    }

    this.hlsInstance = new Hls(hlsConfig);
    this.setupHLSEventHandlers();
    
    // Configure DRM
    if (config.drmConfig?.widevine) {
      this.setupHLSDRM(config.drmConfig);
    }

    this.hlsInstance.loadSource(url);
    this.hlsInstance.attachMedia(this.videoElement);
  }

  private async loadDASH(url: string, config: Required<PlayerOptions>): Promise<void> {
    if (!this.videoElement) throw new Error('Video element not initialized');

    this.destroyExistingPlayer();

    this.dashInstance = dashjs.MediaPlayer().create();

    const dashConfig = {
      debug: {
        logLevel: config.debug ? dashjs.Debug.LOG_LEVEL_DEBUG : dashjs.Debug.LOG_LEVEL_NONE,
      },
      streaming: {
        lowLatencyEnabled: config.lowLatencyMode,
        liveDelay: config.lowLatencyMode ? 2 : 4,
        liveCatchUpMinDrift: 0.05,
        liveCatchUpPlaybackRate: 0.5,
        stableBufferTime: config.bufferForPlayback,
        bufferTimeAtTopQuality: config.maxBufferLength,
        bufferTimeAtTopQualityLongForm: config.maxMaxBufferLength,
        abr: {
          autoSwitchBitrate: {
            video: config.adaptiveBitrate,
            audio: config.adaptiveBitrate,
          },
          limitBitrateByPortal: true,
        },
        retryIntervals: {
          MPD: 500,
          XLinkExpansion: 500,
          InitializationSegment: 1000,
          BitstreamSwitchingSegment: 1000,
          IndexSegment: 1000,
          MediaSegment: 1000,
        },
        retryAttempts: {
          MPD: 3,
          XLinkExpansion: 1,
          InitializationSegment: 3,
          BitstreamSwitchingSegment: 3,
          IndexSegment: 3,
          MediaSegment: 3,
        },
      },
    };

    // Configure DRM
    if (config.drmConfig) {
      this.setupDASHDRM(config.drmConfig);
    }

    this.dashInstance.updateSettings(dashConfig);
    this.dashInstance.initialize(this.videoElement, url, config.autoplay);
    this.setupDASHEventHandlers();
  }

  private async loadWebRTC(streamId: string, config: Required<PlayerOptions>): Promise<void> {
    if (!this.videoElement) throw new Error('Video element not initialized');

    this.destroyExistingPlayer();

    const signalingUrl = config.webrtcConfig?.signalingUrl || 
      `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/ws`;

    this.webSocket = new WebSocket(
      `${signalingUrl}?streamId=${streamId}&userId=${this.generateUserId()}`
    );

    this.webSocket.onopen = () => {
      this.stats.connectionState = 'connecting';
      this.emit('loadstart');
      this.webSocket?.send(JSON.stringify({ type: 'subscribe', streamId }));
    };

    this.webSocket.onmessage = async (event) => {
      const msg = JSON.parse(event.data);
      await this.handleWebRTCSignaling(msg, config);
    };

    this.webSocket.onerror = (error) => {
      this.handleError(new Error('WebSocket error'));
    };

    this.webSocket.onclose = () => {
      this.stats.connectionState = 'disconnected';
      if (!this.isDestroyed) {
        this.handleDisconnection();
      }
    };
  }

  private async handleWebRTCSignaling(msg: any, config: Required<PlayerOptions>): Promise<void> {
    switch (msg.type) {
      case 'subscribe_ready':
        await this.createWebRTCPeerConnection(config);
        break;
      case 'offer':
        await this.handleWebRTCOffer(msg.sdp);
        break;
      case 'candidate':
        await this.handleWebRTCCandidate(msg.candidate);
        break;
      case 'error':
        this.handleError(new Error(msg.error));
        break;
    }
  }

  private async createWebRTCPeerConnection(config: Required<PlayerOptions>): Promise<void> {
    const pcConfig: RTCConfiguration = {
      iceServers: config.webrtcConfig?.iceServers || [],
      bundlePolicy: 'max-bundle',
      rtcpMuxPolicy: 'require',
    };

    this.webrtcPeer = new RTCPeerConnection(pcConfig);

    // Add transceiver for receiving video/audio
    this.webrtcPeer.addTransceiver('video', { direction: 'recvonly' });
    this.webrtcPeer.addTransceiver('audio', { direction: 'recvonly' });

    this.webrtcPeer.ontrack = (event) => {
      if (this.videoElement && event.streams[0]) {
        this.videoElement.srcObject = event.streams[0];
        this.stats.connectionState = 'connected';
        this.emit('canplay');
      }
    };

    this.webrtcPeer.onicecandidate = (event) => {
      if (event.candidate) {
        this.webSocket?.send(JSON.stringify({
          type: 'candidate',
          candidate: event.candidate,
        }));
      }
    };

    this.webrtcPeer.onconnectionstatechange = () => {
      this.stats.connectionState = this.webrtcPeer?.connectionState || 'disconnected';
      if (this.webrtcPeer?.connectionState === 'failed') {
        this.handleError(new Error('WebRTC connection failed'));
      }
    };

    // Create and send offer
    const offer = await this.webrtcPeer.createOffer();
    await this.webrtcPeer.setLocalDescription(offer);
    
    this.webSocket?.send(JSON.stringify({
      type: 'offer',
      sdp: offer,
    }));
  }

  private async handleWebRTCOffer(sdp: RTCSessionDescriptionInit): Promise<void> {
    if (!this.webrtcPeer) return;

    await this.webrtcPeer.setRemoteDescription(sdp);
    const answer = await this.webrtcPeer.createAnswer();
    await this.webrtcPeer.setLocalDescription(answer);

    this.webSocket?.send(JSON.stringify({
      type: 'answer',
      sdp: answer,
    }));
  }

  private async handleWebRTCCandidate(candidate: RTCIceCandidateInit): Promise<void> {
    if (!this.webrtcPeer) return;
    await this.webrtcPeer.addIceCandidate(candidate);
  }

  private setupHLSEventHandlers(): void {
    if (!this.hlsInstance) return;

    this.hlsInstance.on(Hls.Events.MANIFEST_PARSED, (event, data) => {
      this.stats.qualityLevels = data.levels.map(level => ({
        bitrate: level.bitrate,
        width: level.width || 0,
        height: level.height || 0,
        framerate: parseFloat(level.attrs?.FRAMERATE || '30'),
        codec: level.videoCodec || '',
      }));
    });

    this.hlsInstance.on(Hls.Events.LEVEL_SWITCHED, (event, data) => {
      this.stats.qualityLevel = data.level;
      this.emit('qualitychanged', data);
    });

    this.hlsInstance.on(Hls.Events.ERROR, (event, data) => {
      this.handleHLSError(data);
    });

    this.hlsInstance.on(Hls.Events.FRAG_LOADED, (event, data) => {
      this.updateNetworkStats(data);
    });
  }

  private setupDASHEventHandlers(): void {
    if (!this.dashInstance) return;

    this.dashInstance.on(dashjs.MediaPlayer.events.QUALITY_CHANGE_RENDERED, (e) => {
      this.stats.qualityLevel = e.newQuality;
      this.emit('qualitychanged', e);
    });

    this.dashInstance.on(dashjs.MediaPlayer.events.ERROR, (e) => {
      this.handleError(new Error(e.error.message));
    });

    this.dashInstance.on(dashjs.MediaPlayer.events.METRIC_CHANGED, (e) => {
      this.updateDASHMetrics(e);
    });
  }

  private setupVideoEventListeners(): void {
    if (!this.videoElement) return;

    const events = [
      'play', 'pause', 'ended', 'waiting', 'canplay',
      'loadedmetadata', 'timeupdate', 'volumechange', 'error'
    ];

    events.forEach(event => {
      this.videoElement?.addEventListener(event, (e) => {
        this.emit(event as PlayerEvent, e);
      });
    });

    // Monitor dropped frames
    if ('getVideoPlaybackQuality' in this.videoElement) {
      setInterval(() => {
        const quality = this.videoElement?.getVideoPlaybackQuality();
        if (quality) {
          this.stats.droppedFrames = quality.droppedVideoFrames;
        }
      }, 1000);
    }
  }

  private updateStats(): void {
    if (!this.videoElement) return;

    // Buffer calculation
    const buffered = this.videoElement.buffered;
    if (buffered.length > 0) {
      const currentTime = this.videoElement.currentTime;
      const bufferedEnd = buffered.end(buffered.length - 1);
      this.stats.bufferedDuration = Math.max(0, bufferedEnd - currentTime);
    }

    // Protocol-specific stats
    if (this.protocol === 'hls' && this.hlsInstance) {
      this.stats.bitrate = this.hlsInstance.bandwidthEstimate || 0;
      this.stats.latency = this.hlsInstance.latency || 0;
    } else if (this.protocol === 'dash' && this.dashInstance) {
      this.updateDASHStats();
    } else if (this.protocol === 'webrtc' && this.webrtcPeer) {
      this.updateWebRTCStats();
    }

    // Check for warnings
    if (this.stats.bufferedDuration < 2) {
      this.emit('bufferwarning', { buffered: this.stats.bufferedDuration });
    }

    if (this.stats.latency > 5000) {
      this.emit('latencywarning', { latency: this.stats.latency });
    }

    this.emit('stats', { ...this.stats });
  }

  private async updateWebRTCStats(): Promise<void> {
    if (!this.webrtcPeer) return;

    try {
      const stats = await this.webrtcPeer.getStats();
      stats.forEach(report => {
        if (report.type === 'inbound-rtp' && report.mediaType === 'video') {
          this.stats.bitrate = (report.bytesReceived * 8) / 1000;
          this.stats.droppedFrames = report.framesDropped || 0;
        }
      });
    } catch (error) {
      console.error('Failed to get WebRTC stats:', error);
    }
  }

  private updateDASHStats(): void {
    if (!this.dashInstance) return;

    const dashMetrics = this.dashInstance.getDashMetrics();
    const streamInfo = this.dashInstance.getActiveStream();
    
    if (streamInfo) {
      const periodIdx = streamInfo.getStreamInfo().index;
      const repSwitch = dashMetrics.getCurrentRepresentationSwitch('video', true);
      
      if (repSwitch) {
        this.stats.bitrate = repSwitch.to.bandwidth / 1000;
      }
    }
  }

  public async play(): Promise<void> {
    if (!this.videoElement) throw new Error('Video element not initialized');

    try {
      await this.videoElement.play();
      this.emit('play');
    } catch (error) {
      this.handleError(error as Error);
      throw error;
    }
  }

  public pause(): void {
    if (!this.videoElement) return;
    this.videoElement.pause();
    this.emit('pause');
  }

  public seek(time: number): void {
    if (!this.videoElement) return;
    if (this.protocol === 'webrtc') {
      console.warn('Seeking not supported in WebRTC mode');
      return;
    }
    this.videoElement.currentTime = time;
  }

  public setQuality(level: number): void {
    if (this.protocol === 'hls' && this.hlsInstance) {
      this.hlsInstance.currentLevel = level;
    } else if (this.protocol === 'dash' && this.dashInstance) {
      this.dashInstance.setQualityFor('video', level);
    }
  }

  public getQuality(): number {
    return this.stats.qualityLevel;
  }

  public getQualityLevels(): QualityLevel[] {
    return this.stats.qualityLevels;
  }

  public setVolume(volume: number): void {
    if (!this.videoElement) return;
    this.videoElement.volume = Math.max(0, Math.min(1, volume));
  }

  public getVolume(): number {
    return this.videoElement?.volume || 0;
  }

  public setMuted(muted: boolean): void {
    if (!this.videoElement) return;
    this.videoElement.muted = muted;
  }

  public isMuted(): boolean {
    return this.videoElement?.muted || false;
  }

  public getStats(): StreamStats {
    return { ...this.stats };
  }

  public destroy(): void {
    this.isDestroyed = true;
    this.destroyExistingPlayer();
    
    if (this.statsTimer) {
      clearInterval(this.statsTimer);
    }
    
    this.networkMonitor.stop();
    
    if (this.videoElement) {
      this.videoElement.remove();
    }
    
    this.removeAllListeners();
  }

  private destroyExistingPlayer(): void {
    if (this.hlsInstance) {
      this.hlsInstance.destroy();
      this.hlsInstance = null;
    }
    
    if (this.dashInstance) {
      this.dashInstance.reset();
      this.dashInstance = null;
    }
    
    if (this.webrtcPeer) {
      this.webrtcPeer.close();
      this.webrtcPeer = null;
    }
    
    if (this.webSocket) {
      this.webSocket.close();
      this.webSocket = null;
    }
  }

  private createVideoElement(): void {
    this.videoElement = document.createElement('video');
    this.videoElement.controls = this.options.controls;
    this.videoElement.muted = this.options.muted;
    this.videoElement.autoplay = this.options.autoplay;
    this.videoElement.playsInline = true;
    this.videoElement.style.width = '100%';
    this.videoElement.style.height = '100%';
    this.videoElement.style.backgroundColor = '#000';
    
    if (this.container) {
      this.container.appendChild(this.videoElement);
    }
    
    this.setupVideoEventListeners();
  }

  private startStatsCollection(): void {
    this.statsTimer = window.setInterval(() => {
      this.updateStats();
    }, this.options.statsInterval);
  }

  private detectProtocol(url: string, preferred: string): string {
    if (preferred !== 'auto') return preferred;
    
    if (url.includes('.m3u8') || url.includes('/hls/')) return 'hls';
    if (url.includes('.mpd') || url.includes('/dash/')) return 'dash';
    if (url.startsWith('webrtc:') || url.includes('/webrtc/')) return 'webrtc';
    
    return 'hls'; // Default
  }

  private handleError(error: Error): void {
    console.error('Player error:', error);
    this.emit('error', error);
    this.errorRecovery.handleError(error);
  }

  private handleHLSError(data: any): void {
    console.error('HLS error:', data);
    
    if (data.fatal) {
      switch (data.type) {
        case Hls.ErrorTypes.NETWORK_ERROR:
          this.errorRecovery.handleNetworkError(data);
          break;
        case Hls.ErrorTypes.MEDIA_ERROR:
          this.errorRecovery.handleMediaError(data);
          break;
        default:
          this.handleError(new Error(data.details));
      }
    }
  }

  private handleDisconnection(): void {
    if (this.retryCount < (this.options.retryConfig?.maxRetries || 5)) {
      const delay = Math.min(
        (this.options.retryConfig?.retryDelay || 1000) * 
        Math.pow(this.options.retryConfig?.backoffFactor || 2, this.retryCount),
        this.options.retryConfig?.maxRetryDelay || 30000
      );
      
      this.retryCount++;
      setTimeout(() => {
        if (!this.isDestroyed) {
          // Attempt reconnection
        }
      }, delay);
    }
  }

  private generateUserId(): string {
    return `user_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }

  private setupHLSDRM(drmConfig: DRMConfig): void {
    if (!this.hlsInstance) return;
    
    // Configure Widevine for HLS
    if (drmConfig.widevine) {
      this.hlsInstance.config.emeEnabled = true;
      this.hlsInstance.config.widevineLicenseUrl = drmConfig.widevine.licenseUrl;
      this.hlsInstance.config.licenseXhrSetup = (xhr, url) => {
        if (drmConfig.widevine?.headers) {
          Object.entries(drmConfig.widevine.headers).forEach(([key, value]) => {
            xhr.setRequestHeader(key, value);
          });
        }
      };
    }
  }

  private setupDASHDRM(drmConfig: DRMConfig): void {
    if (!this.dashInstance) return;
    
    const protectionData: any = {};
    
    if (drmConfig.widevine) {
      protectionData['com.widevine.alpha'] = {
        serverURL: drmConfig.widevine.licenseUrl,
        httpRequestHeaders: drmConfig.widevine.headers || {},
      };
    }
    
    if (drmConfig.playready) {
      protectionData['com.microsoft.playready'] = {
        serverURL: drmConfig.playready.licenseUrl,
        httpRequestHeaders: drmConfig.playready.headers || {},
      };
    }
    
    this.dashInstance.setProtectionData(protectionData);
  }

  private async setupNativeHLS(config: Required<PlayerOptions>): Promise<void> {
    // Safari native HLS configuration
    if ('WebKitMediaKeys' in window && config.drmConfig?.fairplay) {
      // FairPlay DRM setup for Safari
      await this.setupFairPlayDRM(config.drmConfig.fairplay);
    }
  }

  private async setupFairPlayDRM(fairplayConfig: DRMConfig['fairplay']): Promise<void> {
    if (!this.videoElement || !fairplayConfig) return;
    
    // FairPlay implementation would go here
    // This requires Safari-specific APIs
  }

  private updateNetworkStats(data: any): void {
    if (data.stats) {
      this.stats.downloadSpeed = (data.stats.loaded * 8) / (data.stats.duration * 1000);
    }
  }

  private updateDASHMetrics(event: any): void {
    if (event.metric === 'BufferLevel') {
      this.stats.bufferedDuration = event.value.level;
    }
  }
}

// Quality Controller for adaptive bitrate management
class QualityController {
  constructor(private player: VideoPlayerSDK) {}

  public selectOptimalQuality(bandwidth: number, bufferLevel: number): number {
    const levels = this.player.getQualityLevels();
    if (levels.length === 0) return -1;

    // Advanced ABR algorithm
    const safetyFactor = bufferLevel < 10 ? 0.7 : 0.9;
    const targetBitrate = bandwidth * safetyFactor;

    let selectedLevel = 0;
    for (let i = levels.length - 1; i >= 0; i--) {
      if (levels[i].bitrate <= targetBitrate) {
        selectedLevel = i;
        break;
      }
    }

    return selectedLevel;
  }
}

// Network Monitor for connection quality
class NetworkMonitor {
  private connectionType: string = 'unknown';
  private effectiveBandwidth: number = 0;
  
  constructor(private player: VideoPlayerSDK) {}

  public start(): void {
    if ('connection' in navigator) {
      const connection = (navigator as any).connection;
      this.connectionType = connection.effectiveType;
      this.effectiveBandwidth = connection.downlink * 1000; // Convert to kbps
      
      connection.addEventListener('change', () => {
        this.handleConnectionChange();
      });
    }
  }

  public stop(): void {
    // Clean up listeners
  }

  private handleConnectionChange(): void {
    const connection = (navigator as any).connection;
    this.connectionType = connection.effectiveType;
    this.effectiveBandwidth = connection.downlink * 1000;
    
    // Adjust quality based on connection
    if (this.connectionType === '2g' || this.connectionType === 'slow-2g') {
      this.player.setQuality(0); // Lowest quality
    }
  }
}

// Error Recovery for resilient playback
class ErrorRecovery {
  private errorCount: number = 0;
  private lastErrorTime: number = 0;
  
  constructor(private player: VideoPlayerSDK) {}

  public handleError(error: Error): void {
    const now = Date.now();
    if (now - this.lastErrorTime < 1000) {
      this.errorCount++;
    } else {
      this.errorCount = 1;
    }
    this.lastErrorTime = now;

    if (this.errorCount > 5) {
      console.error('Too many errors, stopping recovery attempts');
      return;
    }

    // Implement recovery strategies
  }

  public handleNetworkError(data: any): void {
    console.log('Attempting network error recovery:', data.details);
    // Implement network-specific recovery
  }

  public handleMediaError(data: any): void {
    console.log('Attempting media error recovery:', data.details);
    // Implement media-specific recovery
  }
}

export default VideoPlayerSDK;