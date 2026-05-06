import React from 'react';
import { useConnectionStatus } from '../hooks/useConnectionStatus';

export interface ConnectionIndicatorProps {
  className?: string;
  showText?: boolean;
  customLabels?: {
    connected?: string;
    connecting?: string;
    disconnected?: string;
    reconnecting?: string;
    error?: string;
  };
  style?: React.CSSProperties;
}

export function ConnectionIndicator({ 
  className, 
  showText = true,
  customLabels = {},
  style 
}: ConnectionIndicatorProps) {
  const { status, error } = useConnectionStatus();

  const labels = {
    connected: 'Connected',
    connecting: 'Connecting...',
    disconnected: 'Disconnected',
    reconnecting: 'Reconnecting...',
    error: 'Connection Error',
    ...customLabels
  };

  const getStatusColor = () => {
    switch (status) {
      case 'connected':
        return '#10b981'; // green
      case 'connecting':
      case 'reconnecting':
        return '#f59e0b'; // amber
      case 'disconnected':
        return '#6b7280'; // gray
      case 'error':
        return '#ef4444'; // red
      default:
        return '#6b7280';
    }
  };

  const defaultStyle: React.CSSProperties = {
    display: 'inline-flex',
    alignItems: 'center',
    gap: '8px',
    padding: '4px 12px',
    borderRadius: '16px',
    backgroundColor: 'rgba(0, 0, 0, 0.05)',
    fontSize: '14px',
    ...style
  };

  return (
    <div className={className} style={defaultStyle}>
      <span
        style={{
          display: 'inline-block',
          width: '8px',
          height: '8px',
          borderRadius: '50%',
          backgroundColor: getStatusColor(),
        }}
      />
      {showText && (
        <span>
          {labels[status]}
          {status === 'error' && error && `: ${error.message}`}
        </span>
      )}
    </div>
  );
}