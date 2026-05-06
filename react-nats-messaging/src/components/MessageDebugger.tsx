import React, { useState, useEffect } from 'react';
import { useNats } from '../hooks/useNats';
import { Subscription } from 'nats.ws';

export interface MessageDebuggerProps {
  subjects?: string[];
  maxMessages?: number;
  style?: React.CSSProperties;
  className?: string;
}

interface DebugMessage {
  id: string;
  timestamp: Date;
  subject: string;
  data: string;
  size: number;
}

export function MessageDebugger({ 
  subjects = ['*'], 
  maxMessages = 50,
  style,
  className 
}: MessageDebuggerProps) {
  const [messages, setMessages] = useState<DebugMessage[]>([]);
  const [isMonitoring, setIsMonitoring] = useState(false);
  const { subscribe, isConnected } = useNats();

  useEffect(() => {
    if (!isConnected || !isMonitoring) return;

    const subscriptions: Subscription[] = [];
    let cancelled = false;

    const startMonitoring = async () => {
      for (const subject of subjects) {
        try {
          const sub = await subscribe(subject, (data, subject) => {
            if (cancelled) return;

            const debugMsg: DebugMessage = {
              id: `${Date.now()}-${Math.random()}`,
              timestamp: new Date(),
              subject,
              data: new TextDecoder().decode(data),
              size: data.byteLength
            };

            setMessages(prev => {
              const updated = [debugMsg, ...prev];
              return updated.slice(0, maxMessages);
            });
          });

          subscriptions.push(sub);
        } catch (err) {
          console.error(`Failed to subscribe to ${subject}:`, err);
        }
      }
    };

    startMonitoring();

    return () => {
      cancelled = true;
      subscriptions.forEach(sub => sub.unsubscribe());
    };
  }, [isConnected, isMonitoring, subjects, subscribe, maxMessages]);

  const defaultStyle: React.CSSProperties = {
    backgroundColor: '#f3f4f6',
    border: '1px solid #e5e7eb',
    borderRadius: '8px',
    padding: '16px',
    maxHeight: '400px',
    overflow: 'auto',
    fontFamily: 'monospace',
    fontSize: '12px',
    ...style
  };

  const handleToggle = () => {
    if (isMonitoring) {
      setMessages([]);
    }
    setIsMonitoring(!isMonitoring);
  };

  const handleClear = () => {
    setMessages([]);
  };

  return (
    <div className={className} style={defaultStyle}>
      <div style={{ marginBottom: '12px', display: 'flex', gap: '8px' }}>
        <button
          onClick={handleToggle}
          style={{
            padding: '4px 12px',
            borderRadius: '4px',
            border: '1px solid #d1d5db',
            backgroundColor: isMonitoring ? '#ef4444' : '#10b981',
            color: 'white',
            cursor: 'pointer'
          }}
        >
          {isMonitoring ? 'Stop' : 'Start'} Monitoring
        </button>
        <button
          onClick={handleClear}
          style={{
            padding: '4px 12px',
            borderRadius: '4px',
            border: '1px solid #d1d5db',
            backgroundColor: '#6b7280',
            color: 'white',
            cursor: 'pointer'
          }}
        >
          Clear
        </button>
        <span style={{ marginLeft: 'auto', color: '#6b7280' }}>
          {messages.length} messages
        </span>
      </div>
      
      <div>
        {messages.length === 0 && (
          <div style={{ color: '#9ca3af', textAlign: 'center', padding: '20px' }}>
            {isMonitoring ? 'Waiting for messages...' : 'Click Start to begin monitoring'}
          </div>
        )}
        
        {messages.map(msg => (
          <div
            key={msg.id}
            style={{
              borderBottom: '1px solid #e5e7eb',
              paddingBottom: '8px',
              marginBottom: '8px'
            }}
          >
            <div style={{ color: '#6b7280', marginBottom: '4px' }}>
              [{msg.timestamp.toLocaleTimeString()}] {msg.subject} ({msg.size} bytes)
            </div>
            <div style={{ 
              backgroundColor: 'white', 
              padding: '4px 8px', 
              borderRadius: '4px',
              wordBreak: 'break-all'
            }}>
              {msg.data}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}