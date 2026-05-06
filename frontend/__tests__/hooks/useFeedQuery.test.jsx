import { renderHook, act, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useFeedQuery } from '@/hooks/queries/useFeedQuery.jsx';
import React from 'react';

// Mock the media API functions
jest.mock('@/api/mediaApi.jsx', () => ({
  getMediaByItem: jest.fn().mockResolvedValue({ media: { id: 'media123' } }),
  getAllMediaImages: jest.fn().mockResolvedValue({ images: ['image1.jpg', 'image2.jpg'] })
}));

// Create a custom wrapper for React Query
const createWrapper = () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        cacheTime: 0,
        staleTime: 0,
      },
    },
    logger: {
      log: () => {},
      warn: () => {},
      error: () => {},
    },
  });
  
  return ({ children }) => (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  );
};

// Mock global fetch
const mockFetch = jest.fn();
global.fetch = mockFetch;

describe('useFeedQuery hook', () => {
  let wrapper;
  
  beforeEach(() => {
    wrapper = createWrapper();
    jest.clearAllMocks();
    mockFetch.mockClear();
  });

  test('should initialize with correct loading state', () => {
    // Mock fetch to return immediately but never resolve
    mockFetch.mockImplementation(() => new Promise(() => {}));
    
    const { result } = renderHook(
      () => useFeedQuery({ limit: 10 }),
      { wrapper }
    );

    expect(result.current.isLoading).toBe(true);
    expect(result.current.data).toBeUndefined();
    expect(result.current.error).toBe(null);
    expect(result.current.hasNextPage).toBe(false);
    expect(result.current.isFetchingNextPage).toBe(false);
  });

  test('should fetch initial page successfully', async () => {
    const mockData = {
      items: [
        { id: '1', title: 'Post 1', content: 'Content 1' },
        { id: '2', title: 'Post 2', content: 'Content 2' },
      ],
      hasMore: true,
      page: 1
    };
    
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockData)
    });
    
    const { result } = renderHook(
      () => useFeedQuery({ limit: 10 }),
      { wrapper }
    );

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(result.current.data.pages).toHaveLength(1);
    expect(result.current.data.pages[0].items).toEqual(mockData.items);
    expect(result.current.hasNextPage).toBe(true);
    expect(mockFetch).toHaveBeenCalledWith('/api/feed?page=1&limit=10');
  });

  test('should handle fetch errors', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 500
    });
    
    const { result } = renderHook(
      () => useFeedQuery({ limit: 10 }),
      { wrapper }
    );

    await waitFor(() => {
      expect(result.current.isError).toBe(true);
    });

    expect(result.current.error).toBeTruthy();
    expect(result.current.data).toBeUndefined();
    expect(result.current.isLoading).toBe(false);
  });

  test('should fetch next page when hasNextPage is true', async () => {
    const page1Data = {
      items: [
        { id: '1', title: 'Post 1', content: 'Content 1' },
        { id: '2', title: 'Post 2', content: 'Content 2' },
      ],
      hasMore: true,
      page: 1
    };
    
    const page2Data = {
      items: [
        { id: '3', title: 'Post 3', content: 'Content 3' },
        { id: '4', title: 'Post 4', content: 'Content 4' },
      ],
      hasMore: false,
      page: 2
    };
    
    mockFetch
      .mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(page1Data)
      })
      .mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(page2Data)
      });
    
    const { result } = renderHook(
      () => useFeedQuery({ limit: 10 }),
      { wrapper }
    );

    // Wait for initial page
    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(result.current.hasNextPage).toBe(true);
    expect(result.current.data.pages).toHaveLength(1);

    // Fetch next page
    await act(async () => {
      await result.current.fetchNextPage();
    });

    await waitFor(() => {
      expect(result.current.data.pages).toHaveLength(2);
    });

    expect(result.current.data.pages[0].items).toEqual(page1Data.items);
    expect(result.current.data.pages[1].items).toEqual(page2Data.items);
    expect(result.current.hasNextPage).toBe(false);
    
    // Verify API calls
    expect(mockFetch).toHaveBeenCalledTimes(2);
    expect(mockFetch).toHaveBeenNthCalledWith(1, '/api/feed?page=1&limit=10');
    expect(mockFetch).toHaveBeenNthCalledWith(2, '/api/feed?page=2&limit=10');
  });

  test('should not fetch next page when hasNextPage is false', async () => {
    const mockData = {
      items: [{ id: '1', title: 'Post 1' }],
      hasMore: false,
      page: 1
    };
    
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockData)
    });
    
    const { result } = renderHook(
      () => useFeedQuery({ limit: 10 }),
      { wrapper }
    );

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(result.current.hasNextPage).toBe(false);

    // Attempt to fetch next page
    await act(async () => {
      await result.current.fetchNextPage();
    });

    // Should still only have one page
    expect(result.current.data.pages).toHaveLength(1);
    expect(mockFetch).toHaveBeenCalledTimes(1);
  });

  test('should handle custom options', async () => {
    const mockData = {
      items: [{ id: '1', title: 'Post 1' }],
      hasMore: false,
      page: 1
    };
    
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockData)
    });
    
    const { result } = renderHook(
      () => useFeedQuery({ 
        limit: 5,
        category: 'electronics',
        type: 'product'
      }),
      { wrapper }
    );

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(result.current.data.pages[0].items).toEqual(mockData.items);
    expect(mockFetch).toHaveBeenCalledWith('/api/feed?page=1&limit=5&category=electronics&type=product');
  });

  test('should provide refetch functionality', async () => {
    const initialData = {
      items: [{ id: '1', title: 'Post 1' }],
      hasMore: false,
      page: 1
    };
    
    const updatedData = {
      items: [{ id: '1', title: 'Updated Post 1' }, { id: '2', title: 'New Post' }],
      hasMore: false,
      page: 1
    };
    
    mockFetch
      .mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(initialData)
      })
      .mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(updatedData)
      });
    
    const { result } = renderHook(
      () => useFeedQuery({ limit: 10 }),
      { wrapper }
    );

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(result.current.data.pages[0].items).toEqual(initialData.items);

    // Refetch
    await act(async () => {
      await result.current.refetch();
    });

    await waitFor(() => {
      expect(result.current.data.pages[0].items).toEqual(updatedData.items);
    });

    expect(mockFetch).toHaveBeenCalledTimes(2);
  });

  test('should handle different cursor types', async () => {
    const page1Data = {
      items: [{ id: '1', title: 'Post 1' }],
      hasMore: true,
      page: 1
    };
    
    const page2Data = {
      items: [{ id: '2', title: 'Post 2' }],
      hasMore: false,
      page: 2
    };
    
    mockFetch
      .mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(page1Data)
      })
      .mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(page2Data)
      });
    
    const { result } = renderHook(
      () => useFeedQuery({ limit: 10 }),
      { wrapper }
    );

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(result.current.hasNextPage).toBe(true);

    await act(async () => {
      await result.current.fetchNextPage();
    });

    // Wait specifically for the second page to be added
    await waitFor(() => {
      expect(result.current.data.pages).toHaveLength(2);
    }, { timeout: 6000 });

    expect(result.current.data.pages[0].items).toEqual(page1Data.items);
    expect(result.current.data.pages[1].items).toEqual(page2Data.items);
  });

  test('should handle empty feed data', async () => {
    const mockData = {
      items: [],
      hasMore: false,
      page: 1
    };
    
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockData)
    });
    
    const { result } = renderHook(
      () => useFeedQuery({ limit: 10 }),
      { wrapper }
    );

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(result.current.data.pages[0].items).toEqual([]);
    expect(result.current.hasNextPage).toBe(false);
  });

  test('should handle partial fetch failures on subsequent pages', async () => {
    const page1Data = {
      items: [{ id: '1', title: 'Post 1' }],
      hasMore: true,
      page: 1
    };
    
    mockFetch
      .mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(page1Data)
      })
      .mockResolvedValueOnce({
        ok: false,
        status: 500
      });
    
    const { result } = renderHook(
      () => useFeedQuery({ limit: 10 }),
      { wrapper }
    );

    // Wait for initial page
    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(result.current.data.pages).toHaveLength(1);

    // Try to fetch next page (should fail)
    await act(async () => {
      try {
        await result.current.fetchNextPage();
      } catch (error) {
        // Expected to fail
      }
    });

    // Should still have first page
    expect(result.current.data.pages).toHaveLength(1);
    expect(result.current.data.pages[0].items).toEqual(page1Data.items);
  });

  test('should provide flat data transformation', async () => {
    const page1Data = {
      items: [{ id: '1', title: 'Post 1' }],
      hasMore: true,
      page: 1
    };
    
    const page2Data = {
      items: [{ id: '2', title: 'Post 2' }],
      hasMore: false,
      page: 2
    };
    
    mockFetch
      .mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(page1Data)
      })
      .mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(page2Data)
      });
    
    const { result } = renderHook(
      () => useFeedQuery({ limit: 10 }),
      { wrapper }
    );

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    await act(async () => {
      await result.current.fetchNextPage();
    });

    // Wait specifically for the second page to be added
    await waitFor(() => {
      expect(result.current.data.pages).toHaveLength(2);
    }, { timeout: 6000 });

    // Test flat data transformation
    const flatData = result.current.data.pages.flatMap(page => page.items);
    expect(flatData).toHaveLength(2);
    expect(flatData[0]).toEqual({ id: '1', title: 'Post 1' });
    expect(flatData[1]).toEqual({ id: '2', title: 'Post 2' });
  });

  test('should handle concurrent fetchNextPage calls', async () => {
    const page1Data = {
      items: [{ id: '1', title: 'Post 1' }],
      hasMore: true,
      page: 1
    };
    
    const page2Data = {
      items: [{ id: '2', title: 'Post 2' }],
      hasMore: false,
      page: 2
    };
    
    mockFetch
      .mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(page1Data)
      })
      .mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(page2Data)
      });
    
    const { result } = renderHook(
      () => useFeedQuery({ limit: 10 }),
      { wrapper }
    );

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    // Make concurrent calls - React Query should handle deduplication
    await act(async () => {
      const promises = [
        result.current.fetchNextPage(),
        result.current.fetchNextPage(),
        result.current.fetchNextPage(), // Third call to test deduplication
      ];
      await Promise.allSettled(promises); // Use allSettled to handle any rejections gracefully
    });

    // Wait for the pages to be loaded
    await waitFor(() => {
      expect(result.current.data.pages).toHaveLength(2);
    }, { timeout: 3000 });

    // Should only make one additional request due to deduplication
    expect(mockFetch).toHaveBeenCalledTimes(2);
  }, 15000); // Increase timeout for this test

  test('should handle query disabled state', () => {
    const { result } = renderHook(
      () => useFeedQuery({ limit: 10, enabled: false }),
      { wrapper }
    );

    // When query is disabled, result should be in initial state
    expect(result.current.isLoading).toBe(false);
    expect(result.current.data).toBeUndefined();
    expect(result.current.status).toBe('idle');
    expect(mockFetch).not.toHaveBeenCalled();
  });
}); 