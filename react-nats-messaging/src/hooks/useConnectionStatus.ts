import { useNatsContext } from '../core/NatsProvider';

export function useConnectionStatus() {
  const { connectionStatus, isConnected, isConnecting, error } = useNatsContext();
  
  return {
    status: connectionStatus,
    isConnected,
    isConnecting,
    isDisconnected: connectionStatus === 'disconnected',
    isReconnecting: connectionStatus === 'reconnecting',
    isError: connectionStatus === 'error',
    error
  };
}