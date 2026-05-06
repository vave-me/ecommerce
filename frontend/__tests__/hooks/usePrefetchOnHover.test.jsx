import { renderHook, act } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { usePrefetchOnHover, useImagePrefetch, useQueryPrefetch } from '@/hooks/usePrefetchOnHover.jsx';
import React from 'react';

// Mock timers - remove global setup
// jest.useFakeTimers();

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
  });
  
  return ({ children }) => (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  );
};

describe('usePrefetchOnHover hooks', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    jest.useFakeTimers();
  });

  afterEach(() => {
    jest.clearAllTimers();
    jest.useRealTimers();
  });

  describe('usePrefetchOnHover', () => {
    test('should return correct event handlers and initial state', () => {
      const mockPrefetchFn = jest.fn();
      
      const { result } = renderHook(() => usePrefetchOnHover(mockPrefetchFn));

      expect(typeof result.current.onMouseEnter).toBe('function');
      expect(typeof result.current.onMouseLeave).toBe('function');
      expect(result.current.hasPrefetched).toBe(false);
    });

    test('should trigger prefetch after default delay on mouse enter', () => {
      const mockPrefetchFn = jest.fn();
      
      const { result } = renderHook(() => usePrefetchOnHover(mockPrefetchFn));

      // Trigger mouse enter
      act(() => {
        result.current.onMouseEnter();
      });

      // Should not have been called yet
      expect(mockPrefetchFn).not.toHaveBeenCalled();
      expect(result.current.hasPrefetched).toBe(false);

      // Fast-forward time by 300ms (default delay)
      act(() => {
        jest.advanceTimersByTime(300);
      });

      expect(mockPrefetchFn).toHaveBeenCalledTimes(1);
      expect(result.current.hasPrefetched).toBe(true);
    });

    test('should trigger prefetch after custom delay', () => {
      const mockPrefetchFn = jest.fn();
      const customDelay = 500;
      
      const { result } = renderHook(() => 
        usePrefetchOnHover(mockPrefetchFn, customDelay)
      );

      act(() => {
        result.current.onMouseEnter();
      });

      // Should not trigger after default delay
      act(() => {
        jest.advanceTimersByTime(300);
      });
      expect(mockPrefetchFn).not.toHaveBeenCalled();

      // Should trigger after custom delay
      act(() => {
        jest.advanceTimersByTime(200); // Total: 500ms
      });
      expect(mockPrefetchFn).toHaveBeenCalledTimes(1);
    });

    test('should cancel prefetch on mouse leave', () => {
      const mockPrefetchFn = jest.fn();
      
      const { result } = renderHook(() => usePrefetchOnHover(mockPrefetchFn));

      // Enter and immediately leave
      act(() => {
        result.current.onMouseEnter();
        result.current.onMouseLeave();
      });

      // Fast-forward time
      act(() => {
        jest.advanceTimersByTime(300);
      });

      expect(mockPrefetchFn).not.toHaveBeenCalled();
      expect(result.current.hasPrefetched).toBe(false);
    });

    test('should not trigger prefetch if already prefetched', () => {
      const mockPrefetchFn = jest.fn();
      
      const { result } = renderHook(() => usePrefetchOnHover(mockPrefetchFn));

      // First hover
      act(() => {
        result.current.onMouseEnter();
      });

      act(() => {
        jest.advanceTimersByTime(300);
      });

      expect(mockPrefetchFn).toHaveBeenCalledTimes(1);
      expect(result.current.hasPrefetched).toBe(true);

      // Second hover should not trigger prefetch
      act(() => {
        result.current.onMouseEnter();
      });

      act(() => {
        jest.advanceTimersByTime(300);
      });

      expect(mockPrefetchFn).toHaveBeenCalledTimes(1); // Still only once
    });

    test('should handle multiple mouse enter/leave events correctly', () => {
      const mockPrefetchFn = jest.fn();
      
      const { result } = renderHook(() => usePrefetchOnHover(mockPrefetchFn));

      // Enter, leave, enter again
      act(() => {
        result.current.onMouseEnter();
      });

      act(() => {
        jest.advanceTimersByTime(100);
      });

      act(() => {
        result.current.onMouseLeave();
      });

      act(() => {
        result.current.onMouseEnter();
      });

      // Should trigger after full delay from second enter
      act(() => {
        jest.advanceTimersByTime(300);
      });

      expect(mockPrefetchFn).toHaveBeenCalledTimes(1);
    });

    test('should handle prefetch function that throws error', () => {
      const mockPrefetchFn = jest.fn().mockImplementation(() => {
        throw new Error('Prefetch failed');
      });
      
      const { result } = renderHook(() => usePrefetchOnHover(mockPrefetchFn));

      act(() => {
        result.current.onMouseEnter();
      });

      // Should not throw error when executing
      expect(() => {
        act(() => {
          jest.advanceTimersByTime(300);
        });
      }).not.toThrow();

      expect(mockPrefetchFn).toHaveBeenCalledTimes(1);
      expect(result.current.hasPrefetched).toBe(true);
    });
  });

  describe('useImagePrefetch', () => {
    let originalImage;

    beforeEach(() => {
      originalImage = global.Image;
      global.Image = jest.fn().mockImplementation(() => ({
        src: '',
        onload: null,
        onerror: null,
      }));
    });

    afterEach(() => {
      global.Image = originalImage;
    });

    test('should return hover handlers and initial state', () => {
      const imageUrls = ['image1.jpg', 'image2.jpg'];
      
      const { result } = renderHook(() => useImagePrefetch(imageUrls));

      expect(typeof result.current.onMouseEnter).toBe('function');
      expect(typeof result.current.onMouseLeave).toBe('function');
      expect(result.current.hasPrefetched).toBe(false);
    });

    test('should prefetch first 3 images immediately on hover', () => {
      const imageUrls = ['image1.jpg', 'image2.jpg', 'image3.jpg', 'image4.jpg'];
      
      const { result } = renderHook(() => useImagePrefetch(imageUrls));

      act(() => {
        result.current.onMouseEnter();
      });

      act(() => {
        jest.advanceTimersByTime(300);
      });

      // Should create 3 Image objects immediately
      expect(global.Image).toHaveBeenCalledTimes(3);
    });

    test('should prefetch remaining images with delay', () => {
      const imageUrls = ['image1.jpg', 'image2.jpg', 'image3.jpg', 'image4.jpg', 'image5.jpg'];
      
      const { result } = renderHook(() => useImagePrefetch(imageUrls));

      act(() => {
        result.current.onMouseEnter();
      });

      act(() => {
        jest.advanceTimersByTime(300);
      });

      // First 3 images should be prefetched
      expect(global.Image).toHaveBeenCalledTimes(3);

      // Advance by additional 500ms for remaining images
      act(() => {
        jest.advanceTimersByTime(500);
      });

      // Should have prefetched all 5 images
      expect(global.Image).toHaveBeenCalledTimes(5);
    });

    test('should handle empty image array', () => {
      const { result } = renderHook(() => useImagePrefetch([]));

      act(() => {
        result.current.onMouseEnter();
      });

      act(() => {
        jest.advanceTimersByTime(300);
      });

      expect(global.Image).not.toHaveBeenCalled();
    });

    test('should handle arrays with null/undefined URLs', () => {
      const imageUrls = ['image1.jpg', null, 'image2.jpg', undefined, 'image3.jpg'];
      
      const { result } = renderHook(() => useImagePrefetch(imageUrls));

      act(() => {
        result.current.onMouseEnter();
      });

      act(() => {
        jest.advanceTimersByTime(300);
      });

      // Should only create Image objects for valid URLs (3 valid ones)
      expect(global.Image).toHaveBeenCalledTimes(3);
    });

    test('should handle custom hover delay', () => {
      const imageUrls = ['image1.jpg', 'image2.jpg'];
      const customDelay = 100;
      
      const { result } = renderHook(() => useImagePrefetch(imageUrls, customDelay));

      act(() => {
        result.current.onMouseEnter();
      });

      // Should not prefetch before custom delay
      act(() => {
        jest.advanceTimersByTime(50);
      });
      expect(global.Image).not.toHaveBeenCalled();

      // Should prefetch after custom delay
      act(() => {
        jest.advanceTimersByTime(50);
      });
      expect(global.Image).toHaveBeenCalledTimes(2);
    });
  });

  describe('useQueryPrefetch', () => {
    let wrapper;
    let queryClient;

    beforeEach(() => {
      wrapper = createWrapper();
      queryClient = new QueryClient({
        defaultOptions: {
          queries: { retry: false, cacheTime: 0, staleTime: 0 },
        },
      });
      
      // Mock the prefetchQuery method
      queryClient.prefetchQuery = jest.fn().mockResolvedValue(undefined);
    });

    test('should return hover handlers and initial state', () => {
      const queryKey = ['data', '1'];
      const queryFn = jest.fn().mockResolvedValue({ data: 'test' });
      
      const { result } = renderHook(
        () => useQueryPrefetch(queryKey, queryFn),
        { wrapper }
      );

      expect(typeof result.current.onMouseEnter).toBe('function');
      expect(typeof result.current.onMouseLeave).toBe('function');
      expect(result.current.hasPrefetched).toBe(false);
    });

    test('should trigger prefetch query on hover after delay', () => {
      const queryKey = ['data', '1'];
      const queryFn = jest.fn().mockResolvedValue({ data: 'test' });
      
      // Mock useQueryClient to return our mock client
      const useQueryClientSpy = jest.spyOn(require('@tanstack/react-query'), 'useQueryClient');
      useQueryClientSpy.mockReturnValue(queryClient);
      
      const { result } = renderHook(
        () => useQueryPrefetch(queryKey, queryFn),
        { wrapper }
      );

      act(() => {
        result.current.onMouseEnter();
      });

      act(() => {
        jest.advanceTimersByTime(300);
      });

      expect(queryClient.prefetchQuery).toHaveBeenCalledWith({
        queryKey,
        queryFn,
        staleTime: 5 * 60 * 1000, // 5 minutes
      });
      expect(result.current.hasPrefetched).toBe(true);

      useQueryClientSpy.mockRestore();
    });

    test('should cancel prefetch if mouse leaves before delay', () => {
      const queryKey = ['data', '1'];
      const queryFn = jest.fn().mockResolvedValue({ data: 'test' });
      
      const useQueryClientSpy = jest.spyOn(require('@tanstack/react-query'), 'useQueryClient');
      useQueryClientSpy.mockReturnValue(queryClient);
      
      const { result } = renderHook(
        () => useQueryPrefetch(queryKey, queryFn),
        { wrapper }
      );

      act(() => {
        result.current.onMouseEnter();
        result.current.onMouseLeave();
      });

      act(() => {
        jest.advanceTimersByTime(300);
      });

      expect(queryClient.prefetchQuery).not.toHaveBeenCalled();
      expect(result.current.hasPrefetched).toBe(false);

      useQueryClientSpy.mockRestore();
    });

    test('should handle custom hover delay', () => {
      const queryKey = ['data', '1'];
      const queryFn = jest.fn().mockResolvedValue({ data: 'test' });
      const customDelay = 100;
      
      const useQueryClientSpy = jest.spyOn(require('@tanstack/react-query'), 'useQueryClient');
      useQueryClientSpy.mockReturnValue(queryClient);
      
      const { result } = renderHook(
        () => useQueryPrefetch(queryKey, queryFn, customDelay),
        { wrapper }
      );

      act(() => {
        result.current.onMouseEnter();
      });

      // Should not prefetch before custom delay
      act(() => {
        jest.advanceTimersByTime(50);
      });
      expect(queryClient.prefetchQuery).not.toHaveBeenCalled();

      // Should prefetch after custom delay
      act(() => {
        jest.advanceTimersByTime(50);
      });
      expect(queryClient.prefetchQuery).toHaveBeenCalledWith({
        queryKey,
        queryFn,
        staleTime: 5 * 60 * 1000,
      });

      useQueryClientSpy.mockRestore();
    });

    test('should not prefetch twice for same query', () => {
      const queryKey = ['data', '1'];
      const queryFn = jest.fn().mockResolvedValue({ data: 'test' });
      
      const useQueryClientSpy = jest.spyOn(require('@tanstack/react-query'), 'useQueryClient');
      useQueryClientSpy.mockReturnValue(queryClient);
      
      const { result } = renderHook(
        () => useQueryPrefetch(queryKey, queryFn),
        { wrapper }
      );

      // First hover
      act(() => {
        result.current.onMouseEnter();
      });

      act(() => {
        jest.advanceTimersByTime(300);
      });

      expect(queryClient.prefetchQuery).toHaveBeenCalledTimes(1);

      // Second hover should not trigger prefetch again
      act(() => {
        result.current.onMouseEnter();
      });

      act(() => {
        jest.advanceTimersByTime(300);
      });

      expect(queryClient.prefetchQuery).toHaveBeenCalledTimes(1); // Still only once

      useQueryClientSpy.mockRestore();
    });

    test('should handle prefetch errors gracefully', () => {
      const queryKey = ['data', '1'];
      const queryFn = jest.fn().mockRejectedValue(new Error('Prefetch failed'));
      
      queryClient.prefetchQuery = jest.fn().mockRejectedValue(new Error('Prefetch failed'));
      
      const useQueryClientSpy = jest.spyOn(require('@tanstack/react-query'), 'useQueryClient');
      useQueryClientSpy.mockReturnValue(queryClient);
      
      const { result } = renderHook(
        () => useQueryPrefetch(queryKey, queryFn),
        { wrapper }
      );

      expect(() => {
        act(() => {
          result.current.onMouseEnter();
        });

        act(() => {
          jest.advanceTimersByTime(300);
        });
      }).not.toThrow();

      expect(result.current.hasPrefetched).toBe(true);

      useQueryClientSpy.mockRestore();
    });
  });
}); 