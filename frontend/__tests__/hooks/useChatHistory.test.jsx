import { renderHook, act } from '@testing-library/react';
import useChatHistory from '@/hooks/useChatHistory.jsx';
import { getMessagesByConversation } from '@/api/client/messagingApi.jsx';
import { useNATS } from '../../context/NATSContext.jsx';
import { useAuth } from '../../context/AuthContext.jsx';

// Mock dependencies
jest.mock('@/api/client/messagingApi.jsx', () => ({
  getMessagesByConversation: jest.fn()
}));

jest.mock('../../context/NATSContext.jsx', () => ({
  useNATS: jest.fn()
}));

jest.mock('../../context/AuthContext.jsx', () => ({
  useAuth: jest.fn()
}));

// Mock the protobuf imports
jest.mock('../../generated_proto/messages_api_events_pb.js', () => ({
  mes: {
    SendMessage: {
      create: jest.fn(obj => obj),
      encode: jest.fn(() => ({
        finish: jest.fn(() => new Uint8Array([1, 2, 3]))
      })),
      decode: jest.fn()
    }
  }
}));

jest.mock('../../generated_proto/message_types_pb.js', () => ({
  message_type: {
    WebsocketMessageData: {
      create: jest.fn(obj => obj),
      encode: jest.fn(() => ({
        finish: jest.fn(() => new Uint8Array([4, 5, 6]))
      })),
      decode: jest.fn()
    }
  }
}));

jest.mock('../../generated_proto/message_api_pb.js', () => ({
  jetstream: {
    StreamMessage: {
      create: jest.fn(obj => obj),
      decode: jest.fn()
    }
  }
}));

// Mock UUID generator
jest.mock('uuid', () => ({
  v4: jest.fn(() => 'mocked-uuid-123')
}));

describe('useChatHistory hook', () => {
  const mockConversationId = 'conv-123';
  const mockUserId = 'user-123';
  const mockRecipientId = 'user-456';
  const mockItemId = 'item-789';
  
  const mockHistoricalMessages = [
    {
      id: 'msg1',
      body: 'Hello there',
      senderId: mockUserId,
      recipientId: mockRecipientId,
      conversationId: mockConversationId,
      createdAt: '2023-01-01T10:00:00Z'
    },
    {
      id: 'msg2',
      body: 'Hi, how are you?',
      senderId: mockRecipientId,
      recipientId: mockUserId,
      conversationId: mockConversationId,
      createdAt: '2023-01-01T10:01:00Z'
    }
  ];

  // Default mock implementations
  beforeEach(() => {
    jest.clearAllMocks();
    
    // Mock auth context
    useAuth.mockReturnValue({
      user: { userId: mockUserId }
    });
    
    // Mock getMessagesByConversation API
    getMessagesByConversation.mockResolvedValue(mockHistoricalMessages);
    
    // Mock NATS context
    const mockUnsubscribe = jest.fn();
    const mockSubscribe = jest.fn().mockImplementation((subject, callback) => {
      // Store callback for later use in tests
      mockSubscribe.callback = callback;
      return Promise.resolve(mockUnsubscribe);
    });
    
    const mockPublish = jest.fn().mockResolvedValue({ success: true });
    
    useNATS.mockReturnValue({
      isConnected: true,
      publish: mockPublish,
      subscribe: mockSubscribe
    });
    
    // Set environment variable
    process.env.NEXT_NATS_SM_NAME = 'messenger.SendMessage';
    
    // Configure mocks for protobuf decoding to simulate an incoming message
    require('../../generated_proto/message_api_pb.js').jetstream.StreamMessage.decode.mockImplementation(() => ({
      data: new Uint8Array([10, 11, 12])
    }));
    
    require('../../generated_proto/message_types_pb.js').message_type.WebsocketMessageData.decode.mockImplementation(() => ({
      payload: new Uint8Array([7, 8, 9])
    }));
    
    require('../../generated_proto/messages_api_events_pb.js').mes.SendMessage.decode.mockImplementation(() => ({
      id: 'msg-id-123',
      body: 'Hello world',
      senderId: 'sender-123',
      recipientId: 'recipient-456',
      conversationId: 'conv-789',
      itemId: 'item-001'
    }));
  });

  test('should load historical messages on mount', async () => {
    const { result } = renderHook(() => 
      useChatHistory(mockConversationId, { 
        recipientId: mockRecipientId, 
        itemId: mockItemId 
      })
    );
    
    // Initially should be loading
    expect(result.current.isLoading).toBe(true);
    
    // Wait for data to load
    await act(async () => {
      await Promise.resolve();
    });
    
    // Should have called the API
    expect(getMessagesByConversation).toHaveBeenCalledWith(mockConversationId);
    
    // Should have messages in state
    expect(result.current.messages.length).toBe(2);
    expect(result.current.messages[0].text).toBe('Hello there');
    expect(result.current.messages[0].isUserMessage).toBe(true); // From current user
    expect(result.current.messages[1].isUserMessage).toBe(false); // From other user
    
    // Should no longer be loading
    expect(result.current.isLoading).toBe(false);
    
    // Should not have error
    expect(result.current.error).toBeNull();
  });

  test('should subscribe to real-time updates', async () => {
    const { result } = renderHook(() => 
      useChatHistory(mockConversationId, { 
        recipientId: mockRecipientId
      })
    );
    
    // Wait for hook to initialize
    await act(async () => {
      await Promise.resolve();
    });
    
    // Should have subscribed to the correct subject
    const { subscribe } = useNATS();
    expect(subscribe).toHaveBeenCalledWith(
      `messenger.SendMessage.${mockConversationId}`,
      expect.any(Function)
    );
  });

  test('should not allow empty messages', async () => {
    const { result } = renderHook(() => 
      useChatHistory(mockConversationId, { 
        recipientId: mockRecipientId
      })
    );
    
    // Wait for hook to initialize
    await act(async () => {
      await Promise.resolve();
    });
    
    const initialMessageCount = result.current.messages.length;
    
    // Try to send an empty message
    await act(async () => {
      await result.current.sendMessage('  ');
    });
    
    // Publish should not have been called
    const { publish } = useNATS();
    expect(publish).not.toHaveBeenCalled();
    
    // Message count should not have changed
    expect(result.current.messages.length).toBe(initialMessageCount);
  });

  test('should handle incoming messages', async () => {
    const { result } = renderHook(() => 
      useChatHistory(mockConversationId, { 
        recipientId: mockRecipientId
      })
    );
    
    // Wait for hook to initialize
    await act(async () => {
      await Promise.resolve();
    });
    
    // Simulate an incoming message via the subscription callback
    const { subscribe } = useNATS();
    
    await act(async () => {
      // Call the subscription callback with mock binary data
      subscribe.callback(new Uint8Array([1, 2, 3]));
      await Promise.resolve();
    });
    
    // The last message should be the new incoming one
    const messages = result.current.messages;
    const lastMessage = messages[messages.length - 1];
    
    expect(lastMessage.id).toBe('msg-id-123');
    expect(lastMessage.text).toBe('Hello world');
    expect(lastMessage.senderId).toBe('sender-123');
    expect(lastMessage.isUserMessage).toBe(false); // Not from current user
  });

  test('should not load messages if not logged in', async () => {
    // Mock user as not logged in
    useAuth.mockReturnValue({ user: null });
    
    const { result } = renderHook(() => 
      useChatHistory(mockConversationId, { 
        recipientId: mockRecipientId
      })
    );
    
    // Wait a bit
    await act(async () => {
      await Promise.resolve();
    });
    
    // Should not have called the API
    expect(getMessagesByConversation).not.toHaveBeenCalled();
    
    // Should have empty messages
    expect(result.current.messages.length).toBe(0);
  });
}); 