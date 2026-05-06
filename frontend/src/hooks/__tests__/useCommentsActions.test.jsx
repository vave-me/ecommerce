import { renderHook, act } from '@testing-library/react';
import { addCommentToCache } from '../useCommentsActions';

// Jest will automatically hoist these mocks to the top, so they're defined before imports
jest.mock('../../context/NATSContext', () => {
  return {
    useNATS: jest.fn()
  };
});

jest.mock('@tanstack/react-query', () => {
  return {
    useMutation: jest.fn(),
    useQueryClient: jest.fn()
  };
});

jest.mock('uuid', () => {
  return {
    v4: jest.fn()
  };
});

// Now import the mocked dependencies
import { useNATS } from '../../context/NATSContext';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { v4 as uuidv4 } from 'uuid';
import { useCommentsActions } from '../useCommentsActions';

describe('addCommentToCache function', () => {
  test('should add top-level comment correctly', () => {
    const allComments = [
      { id: 'comment1', content: 'Existing comment', replies: [] }
    ];
    
    const newComment = {
      id: 'comment2',
      content: 'New top-level comment',
      parentId: '',
      replies: []
    };
    
    const result = addCommentToCache(allComments, newComment);
    
    expect(result).toEqual([
      { id: 'comment1', content: 'Existing comment', replies: [] },
      { id: 'comment2', content: 'New top-level comment', parentId: '', replies: [] }
    ]);
  });
  
  test('should add reply to correct parent', () => {
    const allComments = [
      { id: 'comment1', content: 'Parent comment', replies: [] },
      { id: 'comment2', content: 'Another comment', replies: [] }
    ];
    
    const newComment = {
      id: 'reply1',
      content: 'Reply to comment1',
      parentId: 'comment1',
      replies: []
    };
    
    const result = addCommentToCache(allComments, newComment);
    
    expect(result).toEqual([
      { 
        id: 'comment1', 
        content: 'Parent comment', 
        replies: [
          { id: 'reply1', content: 'Reply to comment1', parentId: 'comment1', replies: [] }
        ] 
      },
      { id: 'comment2', content: 'Another comment', replies: [] }
    ]);
  });
  
  test('should add nested reply correctly', () => {
    const allComments = [
      { 
        id: 'comment1', 
        content: 'Top comment', 
        replies: [
          { id: 'reply1', content: 'First reply', replies: [] }
        ] 
      }
    ];
    
    const newComment = {
      id: 'reply2',
      content: 'Reply to first reply',
      parentId: 'reply1',
      replies: []
    };
    
    const result = addCommentToCache(allComments, newComment);
    
    expect(result).toEqual([
      { 
        id: 'comment1', 
        content: 'Top comment', 
        replies: [
          { 
            id: 'reply1', 
            content: 'First reply', 
            replies: [
              { id: 'reply2', content: 'Reply to first reply', parentId: 'reply1', replies: [] }
            ] 
          }
        ] 
      }
    ]);
  });
});

describe('useCommentsActions hook', () => {
  // Mocks for testing
  const mockPublish = jest.fn().mockResolvedValue(undefined);
  const mockMutate = jest.fn();
  const mockCancelQueries = jest.fn().mockResolvedValue();
  const mockGetQueryData = jest.fn();
  const mockSetQueryData = jest.fn();
  const mockInvalidateQueries = jest.fn();
  
  // Mock callback reference for useMutation
  let mutationCallback;
  let onMutateCallback;
  let onErrorCallback;
  let onSettledCallback;
  
  // Setup before each test
  beforeEach(() => {
    jest.clearAllMocks();
    
    // Setup NATS mock
    useNATS.mockReturnValue({
      isConnected: true,
      publish: mockPublish
    });
    
    // Setup UUID mock
    uuidv4.mockReturnValue('mock-uuid-1234');
    
    // Setup query client mock
    const mockQueryClient = {
      cancelQueries: mockCancelQueries,
      getQueryData: mockGetQueryData,
      setQueryData: mockSetQueryData,
      invalidateQueries: mockInvalidateQueries
    };
    useQueryClient.mockReturnValue(mockQueryClient);
    
    // Setup useMutation to capture the callback functions
    useMutation.mockImplementation((options) => {
      mutationCallback = options.mutationFn;
      onMutateCallback = options.onMutate;
      onErrorCallback = options.onError;
      onSettledCallback = options.onSettled;
      
      return {
        mutate: mockMutate,
        isLoading: false,
        error: null
      };
    });
    
    // Mock Date for consistent tests
    jest.spyOn(global, 'Date').mockImplementation(() => ({
      toISOString: () => '2020-09-13T12:26:40.000Z'
    }));
    global.Date.now = jest.fn().mockReturnValue(1600000000000);
  });
  
  afterEach(() => {
    jest.restoreAllMocks();
  });
  
  test('should initialize with correct values', () => {
    const { result } = renderHook(() => 
      useCommentsActions('item123', 'user123', 'category123')
    );
    
    expect(result.current.isSubmitting).toBe(false);
    expect(result.current.error).toBe(null);
    expect(typeof result.current.createComment).toBe('function');
  });
  
  test('createComment should call mutation.mutate with correct parameters', () => {
    const { result } = renderHook(() => 
      useCommentsActions('item123', 'user123', 'category123')
    );
    
    act(() => {
      result.current.createComment('Test comment', 'parent123');
    });
    
    expect(mockMutate).toHaveBeenCalledWith({
      content: 'Test comment',
      parentId: 'parent123'
    });
  });
  
  test('createComment should use empty string for top-level comments', () => {
    const { result } = renderHook(() => 
      useCommentsActions('item123', 'user123', 'category123')
    );
    
    act(() => {
      result.current.createComment('Test comment');
    });
    
    expect(mockMutate).toHaveBeenCalledWith({
      content: 'Test comment',
      parentId: ''
    });
  });
  
  test('mutation function should throw error when NATS is not connected', async () => {
    // Setup NATS as disconnected
    useNATS.mockReturnValue({
      isConnected: false,
      publish: mockPublish
    });
    
    renderHook(() => useCommentsActions('item123', 'user123', 'category123'));
    
    // Call the mutation function directly
    await expect(mutationCallback({ 
      content: 'Test comment',
      parentId: ''
    })).rejects.toThrow('Not connected to NATS');
    
    expect(mockPublish).not.toHaveBeenCalled();
  });
  
  test('mutation function should publish to NATS when connected', async () => {
    renderHook(() => useCommentsActions('item123', 'user123', 'category123'));
    
    // Call the mutation function directly
    await mutationCallback({ 
      content: 'Test comment',
      parentId: ''
    });
    
    expect(mockPublish).toHaveBeenCalled();
    // Check the first argument contains the correct subject
    expect(mockPublish.mock.calls[0][0]).toContain('comments.AddComment.item123');
  });
  
  test('onMutate function should perform optimistic update', async () => {
    // Setup mock data
    mockGetQueryData.mockReturnValue([{ id: 'existing1', content: 'Existing comment' }]);
    
    renderHook(() => useCommentsActions('item123', 'user123', 'category123'));
    
    // Call onMutate directly
    await onMutateCallback({ 
      content: 'New comment',
      parentId: ''
    });
    
    // Verify query data manipulation
    expect(mockCancelQueries).toHaveBeenCalledWith(['comments', 'item123']);
    expect(mockSetQueryData).toHaveBeenCalled();
    
    // Extract the update function from the setQueryData call
    const updateFn = mockSetQueryData.mock.calls[0][1];
    const oldData = [{ id: 'existing1', content: 'Existing comment' }];
    const updatedData = updateFn(oldData);
    
    // Verify the update function works correctly
    expect(updatedData.length).toBe(2);
    expect(updatedData[0]).toEqual({ id: 'existing1', content: 'Existing comment' });
    expect(updatedData[1].content).toBe('New comment');
    expect(updatedData[1].id).toBe('mock-uuid-1234');
  });
  
  test('onError function should restore previous data', async () => {
    renderHook(() => useCommentsActions('item123', 'user123', 'category123'));
    
    const previousData = [{ id: 'previous' }];
    const context = { previous: previousData };
    const error = new Error('Test error');
    
    // Call onError directly
    onErrorCallback(error, {}, context);
    
    expect(mockSetQueryData).toHaveBeenCalledWith(['comments', 'item123'], previousData);
  });
  
  test('onSettled function should invalidate queries', () => {
    renderHook(() => useCommentsActions('item123', 'user123', 'category123'));
    
    // Call onSettled directly
    onSettledCallback();
    
    expect(mockInvalidateQueries).toHaveBeenCalledWith(['comments', 'item123']);
  });
}); 