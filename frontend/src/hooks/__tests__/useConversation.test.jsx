import { renderHook, act } from '@testing-library/react';
import { useConversation } from '../useConversation';
import {
  getConversationByRecipientAndProduct,
  startConversation
} from '../../api/client/messagingApi';
import { useAuth } from '../../context/AuthContext';

// Mock dependencies
jest.mock('@tanstack/react-query', () => ({
  useQuery: jest.fn(),
  useMutation: jest.fn(),
  useQueryClient: jest.fn()
}));

jest.mock('../../api/client/messagingApi', () => ({
  getConversationByRecipientAndProduct: jest.fn(),
  startConversation: jest.fn()
}));

jest.mock('../../context/AuthContext', () => ({
  useAuth: jest.fn()
}));

describe('useConversation hook', () => {
  const mockUseQuery = require('@tanstack/react-query').useQuery;
  const mockUseMutation = require('@tanstack/react-query').useMutation;
  const mockUseQueryClient = require('@tanstack/react-query').useQueryClient;
  const mockQueryClient = {
    setQueryData: jest.fn(),
    getQueryData: jest.fn(),
    invalidateQueries: jest.fn()
  };
  
  beforeEach(() => {
    jest.clearAllMocks();
    
    // Default auth mock
    useAuth.mockReturnValue({ userId: 'user123' });
    
    // Default query client mock
    mockUseQueryClient.mockReturnValue(mockQueryClient);
    
    // Default useQuery response
    mockUseQuery.mockReturnValue({
      data: [],
      isLoading: false,
      isError: false,
      error: null
    });
    
    // Default useMutation implementation
    mockUseMutation.mockReturnValue({
      mutateAsync: jest.fn().mockResolvedValue({ id: 'new-convo-123' }),
      isLoading: false,
      isError: false,
      error: null
    });
  });
  
  test('should return initial state correctly', () => {
    const { result } = renderHook(() => useConversation('recipient123', 'item123'));
    
    expect(result.current.conversationList).toEqual([]);
    expect(result.current.isLoading).toBe(false);
    expect(result.current.isError).toBe(false);
    expect(typeof result.current.ensureConversationId).toBe('function');
  });
  
  test('should query conversations with correct parameters', () => {
    renderHook(() => useConversation('recipient123', 'item123'));
    
    expect(mockUseQuery).toHaveBeenCalledWith({
      queryKey: ['conversation', 'user123', 'recipient123', 'item123'],
      queryFn: expect.any(Function),
      enabled: true
    });
  });
  
  test('should disable query when required parameters are missing', () => {
    // Test with missing itemId
    renderHook(() => useConversation('recipient123', null));
    
    const lastCall = mockUseQuery.mock.calls[mockUseQuery.mock.calls.length - 1];
    expect(lastCall[0].enabled).toBe(false);
    
    // Test with missing userId
    useAuth.mockReturnValue({ userId: null });
    renderHook(() => useConversation('recipient123', 'item123'));
    
    const lastCall2 = mockUseQuery.mock.calls[mockUseQuery.mock.calls.length - 1];
    expect(lastCall2[0].enabled).toBe(false);
    
    // Test with missing recipientId
    useAuth.mockReturnValue({ userId: 'user123' });
    renderHook(() => useConversation(null, 'item123'));
    
    const lastCall3 = mockUseQuery.mock.calls[mockUseQuery.mock.calls.length - 1];
    expect(lastCall3[0].enabled).toBe(false);
  });
  
  test('ensureConversationId should return existing conversation ID from cache', async () => {
    mockQueryClient.getQueryData.mockReturnValue([{ id: 'existing-convo-123' }]);
    
    const { result } = renderHook(() => useConversation('recipient123', 'item123'));
    
    let conversationId;
    await act(async () => {
      conversationId = await result.current.ensureConversationId();
    });
    
    expect(conversationId).toBe('existing-convo-123');
    // Should not create new conversation
    expect(startConversation).not.toHaveBeenCalled();
  });
  
  test('ensureConversationId should create new conversation when none exists', async () => {
    mockQueryClient.getQueryData.mockReturnValue(null);
    startConversation.mockResolvedValue({ id: 'new-convo-456' });
    
    const mutateAsyncMock = jest.fn().mockResolvedValue({ id: 'new-convo-456' });
    mockUseMutation.mockReturnValue({
      mutateAsync: mutateAsyncMock,
      isLoading: false,
      isError: false,
      error: null
    });
    
    const { result } = renderHook(() => useConversation('recipient123', 'item123'));
    
    let conversationId;
    await act(async () => {
      conversationId = await result.current.ensureConversationId();
    });
    
    expect(mutateAsyncMock).toHaveBeenCalled();
    expect(conversationId).toBe('new-convo-456');
  });
  
  test('createConversation mutation should update cache on success', async () => {
    // Get the mutation options by calling the mock and extracting the first argument
    renderHook(() => useConversation('recipient123', 'item123'));
    
    const mutationOptions = mockUseMutation.mock.calls[0][0];
    
    // Call the onSuccess handler manually
    await mutationOptions.onSuccess({ id: 'new-convo-789' });
    
    // Verify the cache was updated
    expect(mockQueryClient.setQueryData).toHaveBeenCalledWith(
      ['conversation', 'user123', 'recipient123', 'item123'],
      [{ id: 'new-convo-789' }]
    );
  });
  
  test('ensureConversationId should throw error when creation fails', async () => {
    mockQueryClient.getQueryData.mockReturnValue(null);
    
    // Mock mutation to return a result without an ID
    const mutateAsyncMock = jest.fn().mockResolvedValue({ error: 'Creation failed' });
    mockUseMutation.mockReturnValue({
      mutateAsync: mutateAsyncMock,
      isLoading: false,
      isError: false,
      error: null
    });
    
    const { result } = renderHook(() => useConversation('recipient123', 'item123'));
    
    let error;
    await act(async () => {
      try {
        await result.current.ensureConversationId();
      } catch (e) {
        error = e;
      }
    });
    
    expect(error.message).toBe('Conversation creation failed (no id returned)');
  });
}); 