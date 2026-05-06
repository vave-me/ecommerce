import { NatsConnection, Subscription, JetStreamClient } from 'nats.ws';

export interface NatsConfig {
  servers: string | string[];
  options?: {
    reconnect?: boolean;
    maxReconnectAttempts?: number;
    reconnectTimeWait?: number;
    pingInterval?: number;
    maxPingOut?: number;
    debug?: boolean;
    name?: string;
  };
  jetstream?: {
    enabled: boolean;
    options?: any;
  };
}

export interface NatsContextValue {
  connection: NatsConnection | null;
  jetstream: JetStreamClient | null;
  isConnected: boolean;
  isConnecting: boolean;
  error: Error | null;
  connectionStatus: ConnectionStatus;
  connect: () => Promise<void>;
  disconnect: () => Promise<void>;
  publish: (subject: string, data: Uint8Array, options?: PublishOptions) => Promise<void>;
  subscribe: (subject: string, callback: MessageHandler) => Promise<Subscription>;
  request: (subject: string, data: Uint8Array, options?: RequestOptions) => Promise<Uint8Array>;
}

export type ConnectionStatus = 'disconnected' | 'connecting' | 'connected' | 'reconnecting' | 'error';

export interface PublishOptions {
  headers?: Record<string, string>;
  timeout?: number;
}

export interface RequestOptions {
  timeout?: number;
  headers?: Record<string, string>;
}

export type MessageHandler = (data: Uint8Array, subject: string) => void | Promise<void>;

export interface ProtobufConfig<T = any> {
  encode: (message: T) => Uint8Array;
  decode: (data: Uint8Array) => T;
}

export interface MessageEnvelope<T = any> {
  id: string;
  subject: string;
  data: T;
  timestamp: Date;
  headers?: Record<string, string>;
}

export interface SubscriptionOptions {
  queue?: string;
  deduplicate?: boolean;
  deduplicationWindow?: number;
  onError?: (error: Error) => void;
}

export interface UseSubscriptionResult<T> {
  messages: T[];
  latestMessage: T | null;
  error: Error | null;
  isSubscribed: boolean;
  clear: () => void;
}

export interface UsePublishResult<T> {
  publish: (data: T) => Promise<void>;
  isPublishing: boolean;
  error: Error | null;
  lastPublished: T | null;
}

export interface MessageType<T = any> {
  subject: string;
  encode: (message: T) => Uint8Array;
  decode: (data: Uint8Array) => T;
}

export interface StreamMessage {
  id: string;
  name: string;
  data: Uint8Array;
  metadata?: Record<string, any>;
  sentAt: Date;
}

export interface WebsocketMessageData {
  payload: Uint8Array;
  occurredAt: Date;
}