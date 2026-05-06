import { connect, NatsConnection, ConnectionOptions, JetStreamClient, JetStreamManager } from 'nats.ws';
import { NatsConfig, ConnectionStatus } from './types';

export class NatsConnectionManager {
  private connection: NatsConnection | null = null;
  private jetstream: JetStreamClient | null = null;
  private jetstreamManager: JetStreamManager | null = null;
  private config: NatsConfig;
  private reconnectTimer?: NodeJS.Timeout;
  private statusCallbacks: Set<(status: ConnectionStatus) => void> = new Set();
  private _status: ConnectionStatus = 'disconnected';

  constructor(config: NatsConfig) {
    this.config = config;
  }

  get status(): ConnectionStatus {
    return this._status;
  }

  get isConnected(): boolean {
    return this._status === 'connected';
  }

  private setStatus(status: ConnectionStatus) {
    if (this._status !== status) {
      this._status = status;
      this.statusCallbacks.forEach(cb => cb(status));
    }
  }

  onStatusChange(callback: (status: ConnectionStatus) => void): () => void {
    this.statusCallbacks.add(callback);
    return () => this.statusCallbacks.delete(callback);
  }

  async connect(): Promise<void> {
    if (this.connection && this._status === 'connected') {
      return;
    }

    this.setStatus('connecting');

    try {
      const options: ConnectionOptions = {
        servers: this.config.servers,
        reconnect: this.config.options?.reconnect ?? true,
        maxReconnectAttempts: this.config.options?.maxReconnectAttempts ?? 10,
        reconnectTimeWait: this.config.options?.reconnectTimeWait ?? 2000,
        pingInterval: this.config.options?.pingInterval ?? 120000,
        maxPingOut: this.config.options?.maxPingOut ?? 2,
        debug: this.config.options?.debug ?? false,
        name: this.config.options?.name,
      };

      this.connection = await connect(options);

      // Set up connection event handlers
      this.setupConnectionHandlers();

      // Initialize JetStream if enabled
      if (this.config.jetstream?.enabled) {
        this.jetstream = this.connection.jetstream(this.config.jetstream.options);
        this.jetstreamManager = await this.connection.jetstreamManager();
      }

      this.setStatus('connected');
    } catch (error) {
      this.setStatus('error');
      throw error;
    }
  }

  private statusMonitorAbort?: AbortController;

  private setupConnectionHandlers() {
    if (!this.connection) return;

    // Handle disconnection
    this.connection.closed().then(() => {
      this.setStatus('disconnected');
      this.connection = null;
      this.jetstream = null;
      this.jetstreamManager = null;
    }).catch((error) => {
      console.error('Connection closed with error:', error);
    });

    // Monitor connection status
    this.statusMonitorAbort = new AbortController();
    
    (async () => {
      if (!this.connection) return;
      
      try {
        for await (const status of this.connection.status()) {
          if (this.statusMonitorAbort?.signal.aborted) break;
          
          switch (status.type) {
            case 'disconnect':
            case 'error':
              this.setStatus('disconnected');
              break;
            case 'reconnecting':
              this.setStatus('reconnecting');
              break;
            case 'reconnect':
              this.setStatus('connected');
              break;
          }
        }
      } catch (error) {
        if (!this.statusMonitorAbort?.signal.aborted) {
          console.error('Status monitoring error:', error);
        }
      }
    })();
  }

  async disconnect(): Promise<void> {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = undefined;
    }

    // Abort status monitoring
    if (this.statusMonitorAbort) {
      this.statusMonitorAbort.abort();
      this.statusMonitorAbort = undefined;
    }

    if (this.connection) {
      await this.connection.close();
      this.connection = null;
      this.jetstream = null;
      this.jetstreamManager = null;
    }

    this.setStatus('disconnected');
  }

  getConnection(): NatsConnection | null {
    return this.connection;
  }

  getJetStream(): JetStreamClient | null {
    return this.jetstream;
  }

  getJetStreamManager(): JetStreamManager | null {
    return this.jetstreamManager;
  }

  async ensureConnected(): Promise<void> {
    if (!this.isConnected) {
      await this.connect();
    }
  }
}