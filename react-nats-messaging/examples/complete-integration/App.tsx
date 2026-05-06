import React from 'react';
import { 
  NatsProvider, 
  useCommentsActions,
  useChatHistory,
  ConnectionIndicator,
  createMessagingApiClient,
  createCommentsApiClient
} from 'react-nats-messaging';
import axios from 'axios';

// Import your generated protobuf types
// import { commentspb } from './generated_proto/comments_api_pb';
// import { mes } from './generated_proto/messages_api_events_pb';
// import { message_type } from './generated_proto/message_types_pb';
// import { jetstream } from './generated_proto/message_api_pb';

// Create API clients
const axiosInstance = axios.create({
  baseURL: 'https://api.example.com',
  headers: {
    'Content-Type': 'application/json'
  }
});

const messagingApi = createMessagingApiClient(axiosInstance);
const commentsApi = createCommentsApiClient(axiosInstance);

// Example: Comments Component
function CommentsSection({ itemId, userId, categoryId }) {
  const { createComment } = useCommentsActions({
    itemId,
    userId,
    categoryId,
    // Uncomment when you have protobuf types
    // protobufTypes: {
    //   AddComment: commentspb.AddComment,
    //   WebsocketMessageData: message_type.WebsocketMessageData,
    //   StreamMessage: jetstream.StreamMessage
    // },
    onOptimisticUpdate: (comment) => {
      console.log('Optimistic update:', comment);
      // Update your local state/cache here
    },
    onError: (error) => {
      console.error('Comment error:', error);
    }
  });

  const handleSubmit = async (content: string, parentId?: string) => {
    try {
      await createComment(content, parentId);
      // Comment sent successfully
    } catch (error) {
      // Handle error
    }
  };

  return (
    <div>
      <h3>Comments</h3>
      <textarea 
        placeholder="Add a comment..."
        onKeyPress={(e) => {
          if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            handleSubmit(e.currentTarget.value);
            e.currentTarget.value = '';
          }
        }}
      />
    </div>
  );
}

// Example: Chat Component
function ChatWindow({ conversationId, userId, recipientId, itemId }) {
  const { messages, sendMessage, isLoading, error } = useChatHistory({
    conversationId,
    userId,
    recipientId,
    itemId,
    // Uncomment when you have protobuf types
    // protobufTypes: {
    //   SendMessage: mes.SendMessage,
    //   WebsocketMessageData: message_type.WebsocketMessageData,
    //   StreamMessage: jetstream.StreamMessage
    // },
    fetchHistory: async (convId) => {
      // Fetch message history from your API
      return await messagingApi.getMessagesByConversation(convId);
    },
    onError: (error) => {
      console.error('Chat error:', error);
    }
  });

  if (isLoading) return <div>Loading messages...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div>
      <div style={{ height: '400px', overflowY: 'auto' }}>
        {messages.map((msg) => (
          <div key={msg.id} style={{
            textAlign: msg.isUserMessage ? 'right' : 'left',
            margin: '8px',
            padding: '8px',
            backgroundColor: msg.isUserMessage ? '#007bff' : '#f0f0f0',
            color: msg.isUserMessage ? 'white' : 'black',
            borderRadius: '8px'
          }}>
            {msg.text}
          </div>
        ))}
      </div>
      
      <input
        type="text"
        placeholder="Type a message..."
        onKeyPress={async (e) => {
          if (e.key === 'Enter') {
            const text = e.currentTarget.value;
            if (text.trim()) {
              await sendMessage(text);
              e.currentTarget.value = '';
            }
          }
        }}
      />
    </div>
  );
}

// Main App
export function App() {
  const userId = 'user-123';
  const conversationId = 'conv-456';
  const recipientId = 'user-789';
  const itemId = 'item-101';
  const categoryId = 'cat-1';

  return (
    <NatsProvider
      config={{
        servers: process.env.NEXT_PUBLIC_NATS_URL || 'wss://nats-ws.example.com',
        options: {
          reconnect: true,
          maxReconnectAttempts: 10,
          reconnectTimeWait: 2000,
          pingInterval: 120000,
          maxPingOut: 2
        },
        jetstream: {
          enabled: true
        }
      }}
      autoConnect={true}
      onConnectionChange={(status) => {
        console.log('Connection status:', status);
      }}
      onError={(error) => {
        console.error('NATS error:', error);
      }}
    >
      <div style={{ padding: '20px' }}>
        <div style={{ marginBottom: '20px', display: 'flex', justifyContent: 'space-between' }}>
          <h1>NATS Messaging Example</h1>
          <ConnectionIndicator />
        </div>

        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px' }}>
          <div>
            <h2>Chat</h2>
            <ChatWindow
              conversationId={conversationId}
              userId={userId}
              recipientId={recipientId}
              itemId={itemId}
            />
          </div>

          <div>
            <h2>Comments</h2>
            <CommentsSection
              itemId={itemId}
              userId={userId}
              categoryId={categoryId}
            />
          </div>
        </div>
      </div>
    </NatsProvider>
  );
}