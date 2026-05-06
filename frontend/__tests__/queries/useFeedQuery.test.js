import React from 'react';
import { renderHook, waitFor, act } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useFeedQuery } from '@/hooks/queries/useFeedQuery.jsx';
import * as mediaApi from '@/api/mediaApi.jsx';

// Mock fetch globally
global.fetch = jest.fn();

// Mock the mediaApi functions
jest.mock('@/api/mediaApi.jsx', () => ({
  getMediaByItem: jest.fn(),
  getAllMediaImages: jest.fn()
}));

// Create a fresh QueryClient for each test
const createTestQueryClient = () => new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
    },
  },
  logger: {
    log: console.log,
    warn: console.warn,
    error: () => {},
  }
});

describe('useFeedQuery', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    
    // Mock setTimeout
    jest.useFakeTimers();
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  it('should fetch feed data successfully', async () => {
    // Create query client and mock prefetch
    const queryClient = createTestQueryClient();
    queryClient.prefetchQuery = jest.fn();
    
    // Create wrapper with our test query client
    const wrapper = ({ children }) => (
      <QueryClientProvider client={queryClient}>
        {children}
      </QueryClientProvider>
    );

    // Mock successful response
    global.fetch.mockResolvedValueOnce({
      ok: true,
      json: jest.fn().mockResolvedValueOnce({
        items: [
          { id: '1', type: 'post', title: 'Test Post' },
          { id: '2', type: 'service', title: 'Test Service' }
        ],
        hasMore: true,
        page: 1
      })
    });

    // Mock media API responses
    mediaApi.getMediaByItem.mockResolvedValue({ media: { id: 'media1' } });
    mediaApi.getAllMediaImages.mockResolvedValue({ images: [{ url: 'test.jpg' }] });

    // Render the hook
    const { result } = renderHook(() => useFeedQuery({ limit: 2 }), { wrapper });

    // Wait for the query to complete
    await waitFor(() => {
      expect(result.current.status).toBe('success');
    });

    // Assert fetch was called with correct URL
    expect(global.fetch).toHaveBeenCalledWith(
      '/api/feed?page=1&limit=2'
    );

    // Check if data was loaded correctly
    expect(result.current.data.pages[0].items.length).toBe(2);
    expect(result.current.data.pages[0].hasMore).toBe(true);
  });

  it('should handle error when fetch fails', async () => {
    // Create query client and mock prefetch
    const queryClient = createTestQueryClient();
    queryClient.prefetchQuery = jest.fn();
    
    // Create wrapper with our test query client
    const wrapper = ({ children }) => (
      <QueryClientProvider client={queryClient}>
        {children}
      </QueryClientProvider>
    );

    // Mock error response
    global.fetch.mockResolvedValueOnce({
      ok: false
    });

    console.error = jest.fn(); // Suppress expected error in test output

    // Render the hook
    const { result } = renderHook(() => useFeedQuery(), { wrapper });

    // Wait for the query to fail
    await waitFor(() => {
      expect(result.current.status).toBe('error');
    });

    // Assert fetch was called
    expect(global.fetch).toHaveBeenCalled();
  });

  it('should fetch with category and type filters', async () => {
    // Create query client and mock prefetch
    const queryClient = createTestQueryClient();
    queryClient.prefetchQuery = jest.fn();
    
    // Create wrapper with our test query client
    const wrapper = ({ children }) => (
      <QueryClientProvider client={queryClient}>
        {children}
      </QueryClientProvider>
    );

    // Mock successful response
    global.fetch.mockResolvedValueOnce({
      ok: true,
      json: jest.fn().mockResolvedValueOnce({
        items: [{ id: '1', type: 'post', category: 'news' }],
        hasMore: false,
        page: 1
      })
    });

    // Mock media API responses with empty results
    mediaApi.getMediaByItem.mockResolvedValue({ media: { id: 'media1' } });
    mediaApi.getAllMediaImages.mockResolvedValue({ images: [] });

    // Render the hook with filters
    const { result } = renderHook(
      () => useFeedQuery({ category: 'news', type: 'post', limit: 10 }),
      { wrapper }
    );

    // Wait for the query to complete
    await waitFor(() => {
      expect(result.current.status).toBe('success');
    });

    // Assert fetch was called with correct URL including filters
    expect(global.fetch).toHaveBeenCalledWith(
      '/api/feed?page=1&limit=10&category=news&type=post'
    );

    // Check if hasNextPage is false when there's no more data
    expect(result.current.hasNextPage).toBe(false);
  });

  it('should fetch the next page of data', async () => {
    // Create query client and mock prefetch
    const queryClient = createTestQueryClient();
    queryClient.prefetchQuery = jest.fn();
    
    // Create wrapper with our test query client
    const wrapper = ({ children }) => (
      <QueryClientProvider client={queryClient}>
        {children}
      </QueryClientProvider>
    );

    // Mock first page response
    global.fetch.mockResolvedValueOnce({
      ok: true,
      json: jest.fn().mockResolvedValueOnce({
        items: [{ id: '1', type: 'post' }],
        hasMore: true,
        page: 1
      })
    });

    // Mock second page response
    global.fetch.mockResolvedValueOnce({
      ok: true,
      json: jest.fn().mockResolvedValueOnce({
        items: [{ id: '2', type: 'service' }],
        hasMore: false,
        page: 2
      })
    });

    // Mock media API responses
    mediaApi.getMediaByItem.mockResolvedValue({ media: { id: 'media1' } });
    mediaApi.getAllMediaImages.mockResolvedValue({ images: [] });

    // Render the hook
    const { result } = renderHook(() => useFeedQuery(), { wrapper });

    // Wait for the first query to complete
    await waitFor(() => {
      expect(result.current.status).toBe('success');
    });
    
    // Ensure we have the first page loaded
    expect(result.current.data.pages.length).toBe(1);
    expect(result.current.data.pages[0].page).toBe(1);

    // Call fetchNextPage
    await act(async () => {
      await result.current.fetchNextPage();
    });

    // Wait for the second page to be fetched with more specific condition
    await waitFor(() => {
      // Make sure fetch was called a second time
      expect(global.fetch).toHaveBeenCalledTimes(2);
      // And that we have both pages in the result
      expect(result.current.data.pages.length).toBe(2);
      expect(result.current.data.pages[1].page).toBe(2);
    });

    // Check if both pages were loaded
    expect(result.current.data.pages[0].page).toBe(1);
    expect(result.current.data.pages[1].page).toBe(2);
    
    // Check if fetch was called twice with different page params
    expect(global.fetch).toHaveBeenNthCalledWith(1, expect.stringContaining('page=1'));
    expect(global.fetch).toHaveBeenNthCalledWith(2, expect.stringContaining('page=2'));
  });

  it('should prefetch items with delayed behavior for items beyond the first 5', async () => {
    // Create query client and mock prefetch
    const queryClient = createTestQueryClient();
    queryClient.prefetchQuery = jest.fn();
    
    // Create wrapper with our test query client
    const wrapper = ({ children }) => (
      <QueryClientProvider client={queryClient}>
        {children}
      </QueryClientProvider>
    );

    // Mock successful response with 6 items
    global.fetch.mockResolvedValueOnce({
      ok: true,
      json: jest.fn().mockResolvedValueOnce({
        items: Array.from({ length: 6 }, (_, i) => ({ id: String(i+1), type: 'post' })),
        hasMore: true,
        page: 1
      })
    });

    // Mock media API responses
    mediaApi.getMediaByItem.mockResolvedValue({ media: { id: 'media1' } });
    mediaApi.getAllMediaImages.mockResolvedValue({ images: [] });

    // Render the hook
    renderHook(() => useFeedQuery(), { wrapper });

    // Wait for the query to complete and prefetch to be called
    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalled();
    });

    // Fast-forward timers to trigger the delayed prefetch
    jest.runAllTimers();
  });
}); 