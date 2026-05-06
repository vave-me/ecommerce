import React from 'react';
import { renderHook } from '@testing-library/react';
import { NatsProvider } from '../../core/NatsProvider';
import { useNats } from '../useNats';

describe('useNats', () => {
  const wrapper = ({ children }: { children: React.ReactNode }) => (
    <NatsProvider
      config={{
        servers: 'wss://test.nats.io',
        options: { reconnect: true }
      }}
    >
      {children}
    </NatsProvider>
  );

  it('should provide NATS context values', () => {
    const { result } = renderHook(() => useNats(), { wrapper });

    expect(result.current).toHaveProperty('isConnected');
    expect(result.current).toHaveProperty('isConnecting');
    expect(result.current).toHaveProperty('connectionStatus');
    expect(result.current).toHaveProperty('connect');
    expect(result.current).toHaveProperty('disconnect');
    expect(result.current).toHaveProperty('publish');
    expect(result.current).toHaveProperty('subscribe');
  });

  it('should throw error when used outside provider', () => {
    // This should throw, so we wrap in try-catch
    let error: Error | null = null;
    try {
      renderHook(() => useNats());
    } catch (e) {
      error = e as Error;
    }
    
    expect(error).not.toBeNull();
    expect(error?.message).toBe('useNatsContext must be used within a NatsProvider');
  });
});