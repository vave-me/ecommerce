# React NATS Messaging

A powerful React library for building real-time messaging applications with NATS and Protocol Buffers support.

## Features

- 🚀 **Easy Integration** - Simple React hooks and components for NATS messaging
- 📦 **Protocol Buffers** - Full protobuf support with flexible encoding/decoding
- 🔄 **Real-time Updates** - Subscribe to messages with automatic UI updates
- 🛡️ **Type Safety** - Full TypeScript support
- 🔌 **Auto-reconnect** - Resilient connection management with exponential backoff
- 💾 **Message Deduplication** - Built-in duplicate detection
- 🎯 **JetStream Support** - Optional JetStream integration for persistence
- 🔧 **Developer Tools** - Debug components for monitoring messages

## Installation

```bash
npm install react-nats-messaging nats.ws protobufjs
# or
yarn add react-nats-messaging nats.ws protobufjs
```

## Quick Start

### 1. Set up the Provider

```tsx
import { NatsProvider } from 'react-nats-messaging';

function App() {
  return (
    <NatsProvider 
      config={{
        servers: 'wss://your-nats-server.com',
        options: {
          reconnect: true,
          maxReconnectAttempts: 10
        }
      }}
      autoConnect={true}
    >
      <YourApp />
    </NatsProvider>
  );
}
```

### 2. Define Message Types

```tsx
import { createMessageType } from 'react-nats-messaging';
import { ChatMessage } from './proto/generated'; // Your protobuf types

const chatMessageType = createMessageType<ChatMessage>(
  'chat.messages', // NATS subject
  ChatMessage      // Protobuf type
);
```

### 3. Subscribe to Messages

```tsx
import { useSubscription } from 'react-nats-messaging';

function ChatRoom({ roomId }) {
  const { messages, error, isSubscribed } = useSubscription(
    {
      ...chatMessageType,
      subject: `chat.messages.${roomId}`
    },
    {
      deduplicate: true,
      onError: (err) => console.error('Subscription error:', err)
    }
  );

  return (
    <div>
      {messages.map(msg => (
        <div key={msg.id}>{msg.text}</div>
      ))}
    </div>
  );
}
```

### 4. Publish Messages

```tsx
import { usePublish } from 'react-nats-messaging';

function MessageInput({ roomId }) {
  const [text, setText] = useState('');
  const { publish, isPublishing } = usePublish(chatMessageType);

  const handleSend = async () => {
    await publish({
      id: Date.now().toString(),
      text,
      userId: 'current-user',
      timestamp: new Date()
    });
    setText('');
  };

  return (
    <input 
      value={text}
      onChange={(e) => setText(e.target.value)}
      onKeyPress={(e) => e.key === 'Enter' && handleSend()}
    />
  );
}
```

## API Reference

### Components

#### `<NatsProvider>`

The root provider component that manages NATS connection.

```tsx
interface NatsProviderProps {
  config: NatsConfig;
  encoderConfig?: EncoderConfig;
  autoConnect?: boolean;
  onConnectionChange?: (status: ConnectionStatus) => void;
  onError?: (error: Error) => void;
}
```

#### `<ConnectionIndicator>`

Visual indicator for connection status.

```tsx
<ConnectionIndicator 
  showText={true}
  customLabels={{
    connected: 'Online',
    disconnected: 'Offline'
  }}
/>
```

#### `<MessageDebugger>`

Development tool for monitoring messages.

```tsx
<MessageDebugger 
  subjects={['chat.*', 'notifications.*']}
  maxMessages={100}
/>
```

### Hooks

#### `useNats()`

Access the NATS context directly.

```tsx
const { 
  isConnected, 
  connectionStatus, 
  connect, 
  disconnect 
} = useNats();
```

#### `useSubscription(messageType, options)`

Subscribe to messages with automatic decoding.

```tsx
const { 
  messages,      // Array of received messages
  latestMessage, // Most recent message
  error,         // Subscription error
  isSubscribed,  // Subscription status
  clear          // Clear message history
} = useSubscription(messageType, {
  deduplicate: true,
  deduplicationWindow: 100,
  onError: (err) => console.error(err)
});
```

#### `usePublish(messageType, options)`

Publish messages with automatic encoding.

```tsx
const { 
  publish,       // Async function to publish
  isPublishing,  // Publishing status
  error,         // Last publish error
  lastPublished  // Last published message
} = usePublish(messageType);
```

#### `useConnectionStatus()`

Monitor connection status.

```tsx
const { 
  status,
  isConnected,
  isConnecting,
  isDisconnected,
  isReconnecting,
  error
} = useConnectionStatus();
```

### Advanced Usage

#### Custom Message Encoding

```tsx
import { MessageEncoder } from 'react-nats-messaging';

const encoder = new MessageEncoder({
  streamMessageType: StreamMessage,
  websocketMessageType: WebsocketMessageData,
  useStreamWrapper: true,
  useWebsocketWrapper: true
});

// Use with hooks
const { messages } = useSubscription(messageType, { encoder });
```

#### Message Deduplication

```tsx
import { MessageDeduplicator } from 'react-nats-messaging';

const deduplicator = new MessageDeduplicator({
  windowMs: 5000,
  maxSize: 1000,
  getKey: (msg) => msg.id
});

// Check for duplicates
if (!deduplicator.isDuplicate(message)) {
  processMessage(message);
}
```

#### Retry Logic

```tsx
import { withRetry, ExponentialBackoff } from 'react-nats-messaging';

// Retry failed operations
const result = await withRetry(
  () => publish(message),
  {
    maxAttempts: 3,
    initialDelayMs: 1000,
    backoffMultiplier: 2
  }
);

// Custom backoff
const backoff = new ExponentialBackoff({
  initialDelayMs: 1000,
  maxDelayMs: 30000,
  jitterMs: 500
});
```

## Complete Integration Example

### With Protocol Buffers (Recommended)

```tsx
import { 
  NatsProvider, 
  useCommentsActions,
  useChatHistory,
  createMessagingApiClient
} from 'react-nats-messaging';

// Import your generated protobuf types
import { commentspb } from './generated_proto/comments_api_pb';
import { mes } from './generated_proto/messages_api_events_pb';
import { message_type } from './generated_proto/message_types_pb';
import { jetstream } from './generated_proto/message_api_pb';

// 1. Setup Provider
function App() {
  return (
    <NatsProvider
      config={{
        servers: 'wss://nats-ws.example.com',
        jetstream: { enabled: true }
      }}
      autoConnect={true}
    >
      <YourApp />
    </NatsProvider>
  );
}

// 2. Comments Integration
function Comments({ itemId, userId }) {
  const { createComment } = useCommentsActions({
    itemId,
    userId,
    protobufTypes: {
      AddComment: commentspb.AddComment,
      WebsocketMessageData: message_type.WebsocketMessageData,
      StreamMessage: jetstream.StreamMessage
    },
    onOptimisticUpdate: (comment) => {
      // Update UI optimistically
    }
  });

  return (
    <button onClick={() => createComment('Great post!')}>
      Add Comment
    </button>
  );
}

// 3. Chat Integration
function Chat({ conversationId, userId, recipientId }) {
  const { messages, sendMessage } = useChatHistory({
    conversationId,
    userId,
    recipientId,
    protobufTypes: {
      SendMessage: mes.SendMessage,
      WebsocketMessageData: message_type.WebsocketMessageData,
      StreamMessage: jetstream.StreamMessage
    },
    fetchHistory: messagingApi.getMessagesByConversation
  });

  return (
    <div>
      {messages.map(msg => (
        <div key={msg.id}>{msg.text}</div>
      ))}
      <input onKeyPress={(e) => {
        if (e.key === 'Enter') sendMessage(e.target.value);
      }} />
    </div>
  );
}
```

### Without Protocol Buffers (Simple JSON)

```tsx
// Works out of the box with JSON encoding
const { createComment } = useCommentsActions({
  itemId: 'item-123',
  userId: 'user-456',
  // No protobufTypes needed - uses JSON
});

const { messages, sendMessage } = useChatHistory({
  conversationId: 'conv-789',
  userId: 'user-456',
  recipientId: 'user-123',
  // No protobufTypes needed - uses JSON
});
```

## API Integration

The library includes API client factories for easy REST integration:

```tsx
import axios from 'axios';
import { createMessagingApiClient, createCommentsApiClient } from 'react-nats-messaging';

// Create your axios instance
const api = axios.create({
  baseURL: 'https://api.example.com',
  headers: { Authorization: `Bearer ${token}` }
});

// Create API clients
const messagingApi = createMessagingApiClient(api);
const commentsApi = createCommentsApiClient(api);

// Use in your components
const conversations = await messagingApi.getConversations(userId);
const comments = await commentsApi.getCommentsByItem(itemId);
```

## Best Practices

1. **Connection Management**
   - Use `autoConnect` for automatic connection on mount
   - Handle connection errors gracefully
   - Show connection status to users

2. **Message Deduplication**
   - Always enable deduplication for real-time subscriptions
   - Use meaningful message IDs
   - Configure appropriate deduplication windows

3. **Error Handling**
   - Provide `onError` callbacks to hooks
   - Show user-friendly error messages
   - Implement retry logic for critical operations

4. **Performance**
   - Limit message history size
   - Use message pagination for large datasets
   - Implement virtual scrolling for message lists

5. **Security**
   - Use WSS (WebSocket Secure) connections
   - Validate messages on both client and server
   - Implement proper authentication

## TypeScript Support

The library is written in TypeScript and provides full type definitions.

```tsx
import { 
  NatsConfig, 
  ConnectionStatus, 
  MessageType,
  UseSubscriptionResult 
} from 'react-nats-messaging';

// Type-safe message definitions
interface MyMessage {
  id: string;
  content: string;
  timestamp: Date;
}

const messageType: MessageType<MyMessage> = {
  subject: 'my.messages',
  encode: (msg) => encode(msg),
  decode: (data) => decode(data)
};
```

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details.

## License

MIT © [Your Name]