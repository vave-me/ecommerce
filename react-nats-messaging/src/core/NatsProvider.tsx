import React, { createContext, useContext, useEffect, useState, useCallback, useRef } from 'react';
import { Subscription } from 'nats.ws';
import { NatsConnectionManager } from './NatsConnection';
import { MessageEncoder, EncoderConfig } from './MessageEncoder';
import { 
  NatsConfig, 
  NatsContextValue, 
  ConnectionStatus, 
  MessageHandler,
  PublishOptions,
  RequestOptions
} from './types';

const NatsContext = createContext<NatsContextValue | null>(null);

export interface NatsProviderProps {
  children: React.ReactNode;
  config: NatsConfig;
  encoderConfig?: EncoderConfig;
  autoConnect?: boolean;
  onConnectionChange?: (status: ConnectionStatus) => void;
  onError?: (error: Error) => void;
}

export function NatsProvider({ 
  children, 
  config, 
  encoderConfig,
  autoConnect = false,
  onConnectionChange,
  onError 
}: NatsProviderProps) {
  const [connectionStatus, setConnectionStatus] = useState<ConnectionStatus>('disconnected');
  const [isConnecting, setIsConnecting] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  
  const connectionManagerRef = useRef<NatsConnectionManager>();
  const encoderRef = useRef<MessageEncoder>();
  const subscriptionsRef = useRef<Set<Subscription>>(new Set());

  // Initialize connection manager and encoder
  useEffect(() => {
    let mounted = true;
    
    connectionManagerRef.current = new NatsConnectionManager(config);
    encoderRef.current = new MessageEncoder(encoderConfig);

    // Set up status change listener
    const unsubscribe = connectionManagerRef.current.onStatusChange((status) => {
      if (mounted) {
        setConnectionStatus(status);
        onConnectionChange?.(status);
      }
    });

    // Auto-connect if enabled
    if (autoConnect) {
      connect().catch(err => {
        if (mounted) {
          setError(err);
          onError?.(err);
        }
      });
    }

    return () => {
      mounted = false;
      unsubscribe();
      disconnect();
    };
  }, []); // Only run on mount

  const connect = useCallback(async () => {
    if (!connectionManagerRef.current || connectionManagerRef.current.isConnected) {
      return;
    }

    setIsConnecting(true);
    setError(null);

    try {
      await connectionManagerRef.current.connect();
    } catch (err) {
      const error = err as Error;
      setError(error);
      onError?.(error);
      throw error;
    } finally {
      setIsConnecting(false);
    }
  }, [onError]);

  const disconnect = useCallback(async () => {
    // Clean up all subscriptions
    subscriptionsRef.current.forEach(sub => {
      sub.unsubscribe();
    });
    subscriptionsRef.current.clear();

    if (connectionManagerRef.current) {
      await connectionManagerRef.current.disconnect();
    }
  }, []);

  const publish = useCallback(async (
    subject: string, 
    data: Uint8Array, 
    _options?: PublishOptions // eslint-disable-line @typescript-eslint/no-unused-vars
  ) => {
    await connectionManagerRef.current?.ensureConnected();
    const connection = connectionManagerRef.current?.getConnection();
    
    if (!connection) {
      throw new Error('Not connected to NATS');
    }

    if (connectionManagerRef.current?.getJetStream()) {
      // Use JetStream for publishing
      const js = connectionManagerRef.current.getJetStream();
      await js!.publish(subject, data);
    } else {
      // Use regular NATS publish
      connection.publish(subject, data);
    }
  }, []);

  const subscribe = useCallback(async (
    subject: string, 
    callback: MessageHandler
  ): Promise<Subscription> => {
    await connectionManagerRef.current?.ensureConnected();
    const connection = connectionManagerRef.current?.getConnection();
    
    if (!connection) {
      throw new Error('Not connected to NATS');
    }

    const subscription = connection.subscribe(subject);
    subscriptionsRef.current.add(subscription);

    // Process messages
    (async () => {
      for await (const msg of subscription) {
        try {
          await callback(msg.data, msg.subject);
        } catch (err) {
          onError?.(err as Error);
        }
      }
    })();

    // Clean up on unsubscribe
    subscription.closed.then(() => {
      subscriptionsRef.current.delete(subscription);
    }).catch(() => {
      // Ignore errors
    });

    return subscription;
  }, [onError]);

  const request = useCallback(async (
    subject: string,
    data: Uint8Array,
    options?: RequestOptions
  ): Promise<Uint8Array> => {
    await connectionManagerRef.current?.ensureConnected();
    const connection = connectionManagerRef.current?.getConnection();
    
    if (!connection) {
      throw new Error('Not connected to NATS');
    }

    const response = await connection.request(
      subject, 
      data, 
      { 
        timeout: options?.timeout ?? 5000
      }
    );

    return response.data;
  }, []);

  const value: NatsContextValue = {
    connection: connectionManagerRef.current?.getConnection() ?? null,
    jetstream: connectionManagerRef.current?.getJetStream() ?? null,
    isConnected: connectionStatus === 'connected',
    isConnecting,
    error,
    connectionStatus,
    connect,
    disconnect,
    publish,
    subscribe,
    request
  };

  return (
    <NatsContext.Provider value={value}>
      {children}
    </NatsContext.Provider>
  );
}

export function useNatsContext(): NatsContextValue {
  const context = useContext(NatsContext);
  if (!context) {
    throw new Error('useNatsContext must be used within a NatsProvider');
  }
  return context;
}