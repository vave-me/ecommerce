import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useItemMedia, useItemThumbnail, useBulkItemMedia } from '@/hooks/queries/useMediaQueries.jsx';
import { getMediaByItem, getAllMediaImages } from '@/api/mediaApi.jsx';

// Mock the API functions
jest.mock('@/api/mediaApi.jsx', () => ({
  getMediaByItem: jest.fn(),
  getAllMediaImages: jest.fn()
}));

// Create a fresh QueryClient for each test
const createTestQueryClient = () => new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
      cacheTime: 1000 * 60 // 1 minute cache
    },
  },
  logger: {
    log: console.log,
    warn: console.warn,
    error: () => {},
  }
});

// Create a wrapper with the test query client
const createWrapper = () => {
  const testQueryClient = createTestQueryClient();
  return ({ children }) => (
    <QueryClientProvider client={testQueryClient}>{children}</QueryClientProvider>
  );
};

describe('useItemMedia hook', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should return null when itemId is not provided', async () => {
    const wrapper = createWrapper();
    
    // Render the hook with no itemId
    const { result } = renderHook(() => useItemMedia(null), { wrapper });

    // Should not be loading when disabled
    expect(result.current.isLoading).toBe(false);

    // API should not be called
    expect(getMediaByItem).not.toHaveBeenCalled();
    expect(getAllMediaImages).not.toHaveBeenCalled();
    
    // Should return undefined for data when query is not executed
    expect(result.current.data).toBe(undefined);
  });

  it('should fetch media for a valid itemId', async () => {
    const wrapper = createWrapper();
    
    // Mock API responses
    getMediaByItem.mockResolvedValue({ media: { id: 'media-1' } });
    getAllMediaImages.mockResolvedValue({
      images: [{ url: 'https://example.com/image1.jpg' }]
    });

    // Render the hook with an itemId
    const { result } = renderHook(() => useItemMedia('item-1'), { wrapper });

    // Wait for the query to complete
    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    // API should be called with the correct parameters
    expect(getMediaByItem).toHaveBeenCalledWith('item-1');
    expect(getAllMediaImages).toHaveBeenCalledWith('media-1');
    
    // Should return the images data
    expect(result.current.data).toEqual([{ url: 'https://example.com/image1.jpg' }]);
  });

  it('should handle empty responses from the API', async () => {
    const wrapper = createWrapper();
    
    // Mock API responses with empty data
    getMediaByItem.mockResolvedValue({ media: { id: 'media-2' } });
    getAllMediaImages.mockResolvedValue({ images: [] });

    // Render the hook
    const { result } = renderHook(() => useItemMedia('item-2'), { wrapper });

    // Wait for the query to complete
    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    // Should return an empty array
    expect(result.current.data).toEqual([]);
  });

  it('should handle API errors', async () => {
    const wrapper = createWrapper();
    
    // Mock API error
    getMediaByItem.mockRejectedValue(new Error('API Error'));

    // Render the hook
    const { result } = renderHook(() => useItemMedia('item-3'), { wrapper });

    // Wait for the query to fail
    await waitFor(() => {
      expect(result.current.isError).toBe(true);
    });

    // Should have error state
    expect(result.current.error).toBeTruthy();
    expect(result.current.error.message).toBe('API Error');
  });
});

describe('useItemThumbnail hook', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should return a thumbnail URL when media is available', async () => {
    const wrapper = createWrapper();
    
    // Mock API responses
    getMediaByItem.mockResolvedValue({ media: { id: 'media-1' } });
    getAllMediaImages.mockResolvedValue({
      images: [{ url: 'https://example.com/thumbnail.jpg' }]
    });

    // Render the hook
    const { result } = renderHook(() => useItemThumbnail('item-1'), { wrapper });

    // Wait for the query to complete
    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    // Should return the first image URL as thumbnail
    expect(result.current.thumbnail).toBe('https://example.com/thumbnail.jpg');
  });

  it('should return fallback URL when no media is available', async () => {
    const wrapper = createWrapper();
    
    // Mock empty API response
    getMediaByItem.mockResolvedValue({ media: { id: 'media-2' } });
    getAllMediaImages.mockResolvedValue({ images: [] });

    // Render the hook with a custom fallback
    const { result } = renderHook(
      () => useItemThumbnail('item-2', '/custom-fallback.jpg'),
      { wrapper }
    );

    // Wait for the query to complete
    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    // Should return the fallback URL
    expect(result.current.thumbnail).toBe('/custom-fallback.jpg');
  });
});

describe('useBulkItemMedia hook', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should fetch media for multiple items in parallel', async () => {
    const wrapper = createWrapper();
    
    // Mock API responses for different items
    getMediaByItem.mockImplementation((id) => {
      if (id === 'item-1') {
        return Promise.resolve({ media: { id: 'media-1' } });
      } else if (id === 'item-2') {
        return Promise.resolve({ media: { id: 'media-2' } });
      }
      return Promise.reject(new Error('Unknown item'));
    });

    getAllMediaImages.mockImplementation((id) => {
      if (id === 'media-1') {
        return Promise.resolve({ images: [{ url: 'https://example.com/img1.jpg' }] });
      } else if (id === 'media-2') {
        return Promise.resolve({ images: [{ url: 'https://example.com/img2.jpg' }] });
      }
      return Promise.reject(new Error('Unknown media'));
    });

    // Render the hook with multiple items
    const { result } = renderHook(
      () => useBulkItemMedia(['item-1', 'item-2']),
      { wrapper }
    );

    // Wait for all queries to complete
    await waitFor(() => {
      expect(result.current[0].isSuccess).toBe(true);
      expect(result.current[1].isSuccess).toBe(true);
    });
    
    // Verify correct data is returned for each item
    expect(result.current[0].data).toEqual([{ url: 'https://example.com/img1.jpg' }]);
    expect(result.current[1].data).toEqual([{ url: 'https://example.com/img2.jpg' }]);
    
    // Verify API was called with correct IDs
    expect(getMediaByItem).toHaveBeenCalledWith('item-1');
    expect(getMediaByItem).toHaveBeenCalledWith('item-2');
    expect(getAllMediaImages).toHaveBeenCalledWith('media-1');
    expect(getAllMediaImages).toHaveBeenCalledWith('media-2');
  });

  it('should handle API failures for individual items in bulk query', async () => {
    const wrapper = createWrapper();
    
    // Mock mixed success/failure responses
    getMediaByItem.mockImplementation((id) => {
      if (id === 'item-success') {
        return Promise.resolve({ media: { id: 'media-success' } });
      } else if (id === 'item-fail') {
        return Promise.reject(new Error('Failed to fetch media'));
      }
      return Promise.reject(new Error('Unknown item'));
    });

    getAllMediaImages.mockImplementation((id) => {
      if (id === 'media-success') {
        return Promise.resolve({ images: [{ url: 'https://example.com/success.jpg' }] });
      }
      return Promise.reject(new Error('Unknown media'));
    });

    // Render the hook with mixed items
    const { result } = renderHook(
      () => useBulkItemMedia(['item-success', 'item-fail']),
      { wrapper }
    );

    // Wait for all queries to settle
    await waitFor(() => {
      expect(result.current[0].isSuccess || result.current[0].isError).toBe(true);
      expect(result.current[1].isSuccess || result.current[1].isError).toBe(true);
    });
    
    // First item should succeed
    expect(result.current[0].isSuccess).toBe(true);
    expect(result.current[0].data).toEqual([{ url: 'https://example.com/success.jpg' }]);
    
    // Second item should have error state
    expect(result.current[1].isError).toBe(true);
    expect(result.current[1].error).toBeTruthy();
  });
}); 