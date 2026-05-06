import { useNatsContext } from '../core/NatsProvider';

export function useNats() {
  const context = useNatsContext();
  
  return {
    connection: context.connection,
    jetstream: context.jetstream,
    isConnected: context.isConnected,
    isConnecting: context.isConnecting,
    connectionStatus: context.connectionStatus,
    error: context.error,
    connect: context.connect,
    disconnect: context.disconnect,
    publish: context.publish,
    subscribe: context.subscribe,
    request: context.request
  };
}