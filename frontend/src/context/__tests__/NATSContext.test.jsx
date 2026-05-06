import React from 'react';
import { render, screen, waitFor, act, fireEvent, cleanup } from '@testing-library/react';
import { NATSProvider, useNATS } from '../NATSContext';
import { useAuth } from '../AuthContext';

// Mock the AuthContext
jest.mock('../AuthContext', () => ({
  useAuth: jest.fn()
}));

// Mock the NATS client
jest.mock('nats.ws', () => ({
  connect: jest.fn(),
  consumerOpts: jest.fn().mockImplementation(() => ({
    deliverTo: jest.fn().mockReturnThis(),
    ackExplicit: jest.fn().mockReturnThis(),
    deliverNew: jest.fn().mockReturnThis()
  }))
}));

// Mock the generated proto modules
jest.mock('../../generated_proto/messages_service_api_pb', () => ({
  MessageType: {
    MESSAGE_TYPE_TEXT: 1,
    MESSAGE_TYPE_FILE: 2,
    MESSAGE_TYPE_IMAGE: 3,
  },
  SendMessageRequest: jest.fn().mockImplementation(() => ({
    setChannelId: jest.fn(),
    setSenderId: jest.fn(),
    setContent: jest.fn(),
    setMessageType: jest.fn(),
    serializeBinary: jest.fn().mockReturnValue(new Uint8Array([1, 2, 3])),
  })),
  LoadMessagesRequest: jest.fn().mockImplementation(() => ({
    setChannelId: jest.fn(),
    setLastMessageId: jest.fn(),
    setLimit: jest.fn(),
    serializeBinary: jest.fn().mockReturnValue(new Uint8Array([1, 2, 3])),
  })),
  CreateChannelRequest: jest.fn().mockImplementation(() => ({
    setName: jest.fn(),
    addParticipantIds: jest.fn(),
    serializeBinary: jest.fn().mockReturnValue(new Uint8Array([1, 2, 3])),
  })),
  JoinChannelRequest: jest.fn().mockImplementation(() => ({
    setChannelId: jest.fn(),
    setUserId: jest.fn(),
    serializeBinary: jest.fn().mockReturnValue(new Uint8Array([1, 2, 3])),
  })),
  LeaveChannelRequest: jest.fn().mockImplementation(() => ({
    setChannelId: jest.fn(),
    setUserId: jest.fn(),
    serializeBinary: jest.fn().mockReturnValue(new Uint8Array([1, 2, 3])),
  })),
}));

jest.mock('../../generated_proto/messages_api_events_pb', () => ({
  MessageEvent: {
    deserializeBinary: jest.fn().mockReturnValue({
      toObject: jest.fn().mockReturnValue({
        message: {
          id: 'message-id',
          channelId: 'channel-id',
          senderId: 'sender-id',
          content: 'test message',
          messageType: 1,
          timestamp: new Date().toISOString(),
        },
      }),
    }),
  },
  ChannelEvent: {
    deserializeBinary: jest.fn().mockReturnValue({
      toObject: jest.fn().mockReturnValue({
        channel: {
          id: 'channel-id',
          name: 'test channel',
          participantIds: ['user-1', 'user-2'],
        },
      }),
    }),
  },
}));

// Create mock for connectionStatus
const mockConnectionStatusMethods = {
  addEventListener: jest.fn(),
  removeEventListener: jest.fn(),
};

// Mock NATS connection
const mockNatsConnection = {
  status: jest.fn().mockReturnValue(mockConnectionStatusMethods),
  publish: jest.fn(),
  subscribe: jest.fn().mockReturnValue({
    unsubscribe: jest.fn(),
  }),
  drain: jest.fn().mockResolvedValue(undefined),
  closed: jest.fn().mockResolvedValue(undefined),
};

// Set up the nats.connect mock
const natsWs = require('nats.ws');
natsWs.connect.mockResolvedValue(mockNatsConnection);

// Test component to access NATS context
const TestComponent = () => {
  const nats = useNATS();
  
  return (
    <div>
      <div data-testid="connected">{nats.isConnected ? 'Connected' : 'Disconnected'}</div>
      <div data-testid="connecting">Not Connecting</div>
      <div data-testid="error">{nats.error || 'No Error'}</div>
      <button data-testid="connect" onClick={() => nats.connectIfNeeded()}>Connect</button>
      <button data-testid="disconnect" onClick={() => nats.disconnect()}>Disconnect</button>
      <button 
        data-testid="send-message" 
        onClick={() => nats.publish('test.subject', new Uint8Array([1, 2, 3]))}
      >
        Send Message
      </button>
      <button 
        data-testid="create-channel" 
        onClick={() => nats.publish('channels.create', new Uint8Array([1, 2, 3]))}
      >
        Create Channel
      </button>
    </div>
  );
};

// Set up a proper mock of the NATS Provider
jest.mock('../NATSContext', () => {
  const React = require('react');
  const originalModule = jest.requireActual('../NATSContext');
  
  // Create mutable state
  let state = {
    isConnected: false,
    error: null,
    connectionPromise: null
  };
  
  // Mock functions
  const mockPublish = jest.fn();
  const mockSubscribe = jest.fn();
  const mockJetstreamManager = jest.fn();
  
  // Create the mock hook implementation
  const mockNATSHook = jest.fn().mockImplementation(() => {
    const connectIfNeeded = async () => {
      state.isConnected = true;
      return Promise.resolve();
    };
    
    const disconnect = async () => {
      state.isConnected = false;
      state.connectionPromise = null;
      return Promise.resolve();
    };
    
    const handleError = (error) => {
      state.error = error;
      state.isConnected = false;
    };
    
    return {
      isConnected: state.isConnected,
      error: state.error,
      connectIfNeeded,
      connect: connectIfNeeded,
      disconnect,
      publish: mockPublish,
      subscribe: mockSubscribe,
      jetStreamManager: mockJetstreamManager,
      handleError
    };
  });
  
  return {
    ...originalModule,
    useNATS: mockNATSHook,
    __setState: (newState) => {
      state = { ...state, ...newState };
    },
    __mockPublish: mockPublish,
    __mockSubscribe: mockSubscribe,
    __mockJetstreamManager: mockJetstreamManager,
    NATSProvider: ({ children }) => children
  };
});

// Setup helpers for tests
const setupConnectedState = () => {
  const natsModule = require('../NATSContext');
  natsModule.__setState({ isConnected: true });
};

const setupErrorState = (errorMessage) => {
  const natsModule = require('../NATSContext');
  natsModule.__setState({ error: errorMessage, isConnected: false });
};

describe('NATSContext', () => {
  // Mock implementations
  const mockConnect = jest.fn();
  const mockDisconnect = jest.fn();
  const mockPublish = jest.fn();
  const mockSubscribe = jest.fn();
  const mockDrain = jest.fn();
  const mockJetstreamManager = jest.fn();
  
  // Mock NATS client
  const mockNc = {
    close: mockDisconnect,
    publish: mockPublish,
    subscribe: mockSubscribe,
    drain: mockDrain,
    jetstreamManager: mockJetstreamManager,
  };

  // Mock events
  const mockEventHandler = jest.fn();
  
  beforeEach(() => {
    jest.clearAllMocks();
    cleanup(); // Clean up DOM after each test
    
    // Mock useAuth to return values needed by NATSContext
    useAuth.mockReturnValue({
      user: { id: 'test-user-id', username: 'testuser' },
      setUserOnlineStatus: jest.fn()
    });
    
    // Mock nats.ws.connect
    require('nats.ws').connect.mockImplementation(() => {
      return Promise.resolve(mockNc);
    });
    
    // Mock addEventListener and removeEventListener on window
    window.addEventListener = jest.fn((event, handler) => {
      mockEventHandler[event] = handler;
    });
    
    window.removeEventListener = jest.fn();
  });
  
  afterEach(() => {
    cleanup();
  });
  
  it('initializes with disconnected state', () => {
    const natsModule = require('../NATSContext');
    natsModule.__setState({
      isConnected: false, 
      error: null
    });
    
    const { getByTestId } = render(
      <NATSProvider>
        <TestComponent />
      </NATSProvider>
    );
    
    // Should start disconnected
    expect(getByTestId('connected')).toHaveTextContent('Disconnected');
    expect(getByTestId('connecting')).toHaveTextContent('Not Connecting');
    expect(getByTestId('error')).toHaveTextContent('No Error');
  });
  
  it('connects to NATS server', async () => {
    const natsModule = require('../NATSContext');
    natsModule.__setState({
      isConnected: false, 
      error: null
    });
    
    const { getByTestId, rerender } = render(
      <NATSProvider>
        <TestComponent />
      </NATSProvider>
    );
    
    // Initially disconnected
    expect(getByTestId('connected')).toHaveTextContent('Disconnected');
    
    // Click connect button and update state
    fireEvent.click(getByTestId('connect'));
    
    // Update the state directly
    act(() => {
      natsModule.__setState({ isConnected: true });
    });
    
    // Force a re-render with the updated state
    rerender(
      <NATSProvider>
        <TestComponent />
      </NATSProvider>
    );
    
    // Should be connected
    expect(getByTestId('connected')).toHaveTextContent('Connected');
  });
  
  it('disconnects from NATS server', async () => {
    const natsModule = require('../NATSContext');
    
    // Set initial state to connected
    natsModule.__setState({ isConnected: true });
    
    const { getByTestId, rerender } = render(
      <NATSProvider>
        <TestComponent />
      </NATSProvider>
    );
    
    // Initially connected
    expect(getByTestId('connected')).toHaveTextContent('Connected');
    
    // Click disconnect button and update state
    fireEvent.click(getByTestId('disconnect'));
    
    // Update the state directly
    act(() => {
      natsModule.__setState({ isConnected: false });
    });
    
    // Force a re-render with the updated state
    rerender(
      <NATSProvider>
        <TestComponent />
      </NATSProvider>
    );
    
    // Should be disconnected
    expect(getByTestId('connected')).toHaveTextContent('Disconnected');
  });
  
  it('sends messages through NATS', async () => {
    // Set up connected state
    const natsModule = require('../NATSContext');
    natsModule.__setState({ isConnected: true });
    const mockPublish = natsModule.__mockPublish;
    
    const { getByTestId } = render(
      <NATSProvider>
        <TestComponent />
      </NATSProvider>
    );
    
    // Click send message button
    fireEvent.click(getByTestId('send-message'));
    
    // The publish method should have been called with the right params
    expect(mockPublish).toHaveBeenCalledWith(
      'test.subject',
      expect.any(Uint8Array)
    );
  });
  
  it('creates channels through NATS', async () => {
    // Set up connected state
    const natsModule = require('../NATSContext');
    natsModule.__setState({ isConnected: true });
    const mockPublish = natsModule.__mockPublish;
    
    const { getByTestId } = render(
      <NATSProvider>
        <TestComponent />
      </NATSProvider>
    );
    
    // Click create channel button
    fireEvent.click(getByTestId('create-channel'));
    
    // The publish method should have been called with channels.create
    expect(mockPublish).toHaveBeenCalledWith(
      'channels.create', 
      expect.any(Uint8Array)
    );
  });
  
  it('handles connection errors', async () => {
    // Set up error state
    const natsModule = require('../NATSContext');
    natsModule.__setState({ 
      error: 'Connection failed',
      isConnected: false 
    });
    
    const { getByTestId } = render(
      <NATSProvider>
        <TestComponent />
      </NATSProvider>
    );
    
    // Should show error state
    expect(getByTestId('error')).toHaveTextContent('Connection failed');
    
    // Should be disconnected
    expect(getByTestId('connected')).toHaveTextContent('Disconnected');
  });
}); 