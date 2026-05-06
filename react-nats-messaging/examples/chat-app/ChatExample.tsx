import React, { useState, useEffect } from 'react';
import { 
  NatsProvider, 
  useSubscription, 
  usePublish, 
  ConnectionIndicator,
  createMessageType,
  MessageEncoder
} from 'react-nats-messaging';

// Define message types
interface ChatMessage {
  id: string;
  conversationId: string;
  userId: string;
  messageText: string;
  createdAt: string;
}

// Create simple JSON-based encoder for the example
const encoder = new MessageEncoder({
  useStreamWrapper: false,
  useWebsocketWrapper: false
});

// Create message type with JSON encoding
const chatMessageType = createMessageType<ChatMessage>(
  'messenger.SendMessage',
  {
    encode: (msg: ChatMessage) => new TextEncoder().encode(JSON.stringify(msg)),
    decode: (data: Uint8Array) => JSON.parse(new TextDecoder().decode(data))
  }
);

function ChatWindow({ conversationId }: { conversationId: string }) {
  const [inputText, setInputText] = useState('');
  const [messages, setMessages] = useState<ChatMessage[]>([]);

  // Subscribe to messages for this conversation
  const { messages: incomingMessages, error: subError } = useSubscription(
    {
      ...chatMessageType,
      subject: `${chatMessageType.subject}.${conversationId}`
    },
    {
      encoder,
      deduplicate: true,
      onError: (err) => console.error('Subscription error:', err)
    }
  );

  // Publish hook
  const { publish, isPublishing, error: pubError } = usePublish(
    chatMessageType,
    { encoder }
  );

  // Update local messages when new messages arrive
  useEffect(() => {
    if (incomingMessages.length > 0) {
      setMessages(prev => [...prev, ...incomingMessages]);
    }
  }, [incomingMessages]);

  const handleSend = async () => {
    if (!inputText.trim()) return;

    const message: ChatMessage = {
      id: `msg-${Date.now()}`,
      conversationId,
      userId: 'current-user',
      messageText: inputText,
      createdAt: new Date().toISOString()
    };

    try {
      await publish(message);
      setInputText('');
    } catch (err) {
      console.error('Failed to send message:', err);
    }
  };

  return (
    <div style={{ display: 'flex', flexDirection: 'column', height: '400px' }}>
      <div style={{ 
        flex: 1, 
        overflowY: 'auto', 
        padding: '16px',
        backgroundColor: '#f3f4f6'
      }}>
        {messages.map(msg => (
          <div key={msg.id} style={{ marginBottom: '8px' }}>
            <strong>{msg.userId}:</strong> {msg.messageText}
            <span style={{ fontSize: '12px', color: '#666', marginLeft: '8px' }}>
              {new Date(msg.createdAt).toLocaleTimeString()}
            </span>
          </div>
        ))}
        {(subError || pubError) && (
          <div style={{ color: 'red', marginTop: '8px' }}>
            Error: {(subError || pubError)?.message}
          </div>
        )}
      </div>
      
      <div style={{ 
        display: 'flex', 
        padding: '16px',
        borderTop: '1px solid #e5e7eb'
      }}>
        <input
          type="text"
          value={inputText}
          onChange={(e) => setInputText(e.target.value)}
          onKeyPress={(e) => e.key === 'Enter' && handleSend()}
          placeholder="Type a message..."
          style={{ 
            flex: 1, 
            padding: '8px',
            borderRadius: '4px',
            border: '1px solid #d1d5db'
          }}
        />
        <button
          onClick={handleSend}
          disabled={isPublishing || !inputText.trim()}
          style={{
            marginLeft: '8px',
            padding: '8px 16px',
            borderRadius: '4px',
            backgroundColor: '#3b82f6',
            color: 'white',
            border: 'none',
            cursor: isPublishing ? 'not-allowed' : 'pointer',
            opacity: isPublishing ? 0.6 : 1
          }}
        >
          {isPublishing ? 'Sending...' : 'Send'}
        </button>
      </div>
    </div>
  );
}

export function ChatExample() {
  const [conversationId] = useState('conv-123');

  return (
    <NatsProvider
      config={{
        servers: 'wss://nats-ws.example.com',
        options: {
          reconnect: true,
          maxReconnectAttempts: 10
        },
        jetstream: {
          enabled: true
        }
      }}
      autoConnect={true}
    >
      <div style={{ maxWidth: '600px', margin: '0 auto', padding: '20px' }}>
        <div style={{ 
          display: 'flex', 
          justifyContent: 'space-between', 
          alignItems: 'center',
          marginBottom: '20px'
        }}>
          <h1>Chat Example</h1>
          <ConnectionIndicator />
        </div>
        
        <ChatWindow conversationId={conversationId} />
      </div>
    </NatsProvider>
  );
}