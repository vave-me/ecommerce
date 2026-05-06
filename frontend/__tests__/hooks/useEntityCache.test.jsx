import { renderHook, act, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useEntity, useEntityUpdate } from '@/hooks/useEntityCache.jsx';
import React from 'react';

// Create a custom wrapper for React Query
const createWrapper = () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        cacheTime: 0,
        staleTime: 0,
      },
      mutations: {
        retry: false,
      },
    },
  });
  
  return ({ children }) => (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  );
};

describe('useEntityCache hooks', () => {
  let wrapper;
  
  beforeEach(() => {
    wrapper = createWrapper();
    jest.clearAllMocks();
  });

  describe('useEntity', () => {
    test('should initialize with correct loading state', () => {
      const mockFetchFn = jest.fn().mockResolvedValue({ id: '1', name: 'Test' });
      
      const { result } = renderHook(
        () => useEntity('product', '1', mockFetchFn),
        { wrapper }
      );

      expect(result.current.isLoading).toBe(true);
      expect(result.current.data).toBeUndefined();
      expect(result.current.error).toBe(null);
    });

    test('should fetch and return data successfully', async () => {
      const mockData = { id: '1', name: 'Test Product', price: 100 };
      const mockFetchFn = jest.fn().mockResolvedValue(mockData);
      
      const { result } = renderHook(
        () => useEntity('product', '1', mockFetchFn),
        { wrapper }
      );

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(result.current.data).toEqual(mockData);
      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBe(null);
      expect(mockFetchFn).toHaveBeenCalledWith('1');
      expect(mockFetchFn).toHaveBeenCalledTimes(1);
    });

    test('should handle fetch errors', async () => {
      const mockError = new Error('Fetch failed');
      const mockFetchFn = jest.fn().mockRejectedValue(mockError);
      
      const { result } = renderHook(
        () => useEntity('product', '1', mockFetchFn),
        { wrapper }
      );

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });

      expect(result.current.error).toBeTruthy();
      expect(result.current.data).toBeUndefined();
      expect(result.current.isLoading).toBe(false);
      expect(mockFetchFn).toHaveBeenCalledWith('1');
    });

    test('should not fetch when entityId is null or undefined and be in idle state', () => {
      const mockFetchFn = jest.fn().mockResolvedValue({ id: '1', name: 'Product' });
      
      // Test with null entityId
      const { result: nullResult } = renderHook(
        () => useEntity('product', null, mockFetchFn),
        { wrapper }
      );

      expect(nullResult.current.isLoading).toBe(false);
      // In React Query v5, disabled queries might start as pending
      expect(['idle', 'pending']).toContain(nullResult.current.status);
      expect(nullResult.current.data).toBeUndefined();
      expect(mockFetchFn).not.toHaveBeenCalled();

      // Test with undefined entityId
      const { result: undefinedResult } = renderHook(
        () => useEntity('product', undefined, mockFetchFn),
        { wrapper }
      );

      expect(undefinedResult.current.isLoading).toBe(false);
      // In React Query v5, disabled queries might start as pending
      expect(['idle', 'pending']).toContain(undefinedResult.current.status);
      expect(undefinedResult.current.data).toBeUndefined();
      expect(mockFetchFn).not.toHaveBeenCalled();
    });

    test('should accept custom options', async () => {
      const mockData = { id: '1', name: 'Test' };
      const mockFetchFn = jest.fn().mockResolvedValue(mockData);
      const onSuccessMock = jest.fn();
      
      const { result } = renderHook(
        () => useEntity('product', '1', mockFetchFn, { onSuccess: onSuccessMock }),
        { wrapper }
      );

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(result.current.data).toEqual(mockData);
      expect(mockFetchFn).toHaveBeenCalledWith('1');
    });

    test('should handle different entity types with same ID', async () => {
      const productData = { id: '1', name: 'Product 1', type: 'product' };
      const serviceData = { id: '1', name: 'Service 1', type: 'service' };
      
      const productFetchFn = jest.fn().mockResolvedValue(productData);
      const serviceFetchFn = jest.fn().mockResolvedValue(serviceData);
      
      const { result: productResult } = renderHook(
        () => useEntity('product', '1', productFetchFn),
        { wrapper }
      );
      
      const { result: serviceResult } = renderHook(
        () => useEntity('service', '1', serviceFetchFn),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(productResult.current.isSuccess).toBe(true);
        expect(serviceResult.current.isSuccess).toBe(true);
      });

      expect(productResult.current.data).toEqual(productData);
      expect(serviceResult.current.data).toEqual(serviceData);
      expect(productFetchFn).toHaveBeenCalledWith('1');
      expect(serviceFetchFn).toHaveBeenCalledWith('1');
    });

    test('should provide refetch functionality', async () => {
      const mockData1 = { id: '1', name: 'Original Data' };
      const mockData2 = { id: '1', name: 'Updated Data' };
      const mockFetchFn = jest.fn()
        .mockResolvedValueOnce(mockData1)
        .mockResolvedValueOnce(mockData2);
      
      const { result } = renderHook(
        () => useEntity('product', '1', mockFetchFn),
        { wrapper }
      );

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(result.current.data).toEqual(mockData1);

      // Refetch - force the cache to be invalidated
      await act(async () => {
        await result.current.refetch();
      });

      // Wait for refetch to complete
      await waitFor(() => {
        expect(mockFetchFn).toHaveBeenCalledTimes(2);
      });

      expect(result.current.data).toEqual(mockData2);
    });
  });

  describe('useEntityUpdate', () => {
    test('should initialize with correct idle state', () => {
      const mockUpdateFn = jest.fn().mockResolvedValue({ success: true });
      
      const { result } = renderHook(
        () => useEntityUpdate('product', mockUpdateFn),
        { wrapper }
      );

      expect(result.current.mutate).toBeInstanceOf(Function);
      expect(result.current.mutateAsync).toBeInstanceOf(Function);
      expect(result.current.isPending).toBe(false);
      expect(result.current.isSuccess).toBe(false);
      expect(result.current.isError).toBe(false);
      expect(result.current.data).toBeUndefined();
      expect(result.current.error).toBe(null);
    });

    test('should execute update successfully', async () => {
      const mockUpdateData = { id: '1', name: 'Updated Product' };
      const mockUpdateFn = jest.fn().mockResolvedValue(mockUpdateData);
      
      const { result } = renderHook(
        () => useEntityUpdate('product', mockUpdateFn),
        { wrapper }
      );

      await act(async () => {
        result.current.mutate({ id: '1', name: 'Updated Product' });
      });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(result.current.data).toEqual(mockUpdateData);
      expect(result.current.isPending).toBe(false);
      expect(result.current.error).toBe(null);
      expect(mockUpdateFn).toHaveBeenCalledWith({ id: '1', name: 'Updated Product' });
    });

    test('should handle update errors', async () => {
      const mockError = new Error('Update failed');
      const mockUpdateFn = jest.fn().mockRejectedValue(mockError);
      
      const { result } = renderHook(
        () => useEntityUpdate('product', mockUpdateFn),
        { wrapper }
      );

      await act(async () => {
        result.current.mutate({ id: '1', name: 'Updated Product' });
      });

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });

      expect(result.current.error).toBeTruthy();
      expect(result.current.data).toBeUndefined();
      expect(result.current.isPending).toBe(false);
    });

    test('should support mutateAsync for promise-based usage', async () => {
      const mockUpdateData = { id: '1', name: 'Updated Product' };
      const mockUpdateFn = jest.fn().mockResolvedValue(mockUpdateData);
      
      const { result } = renderHook(
        () => useEntityUpdate('product', mockUpdateFn),
        { wrapper }
      );

      let mutationResult;
      await act(async () => {
        mutationResult = await result.current.mutateAsync({ 
          id: '1', 
          name: 'Updated Product' 
        });
      });

      // Wait for the mutation to complete and update state
      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(mutationResult).toEqual(mockUpdateData);
      expect(mockUpdateFn).toHaveBeenCalledWith({ id: '1', name: 'Updated Product' });
    });

    test('should accept custom options with callbacks', async () => {
      const mockUpdateData = { id: '1', name: 'Updated Product' };
      const mockUpdateFn = jest.fn().mockResolvedValue(mockUpdateData);
      const onSuccessMock = jest.fn();
      const onErrorMock = jest.fn();
      
      const customOptions = {
        onSuccess: onSuccessMock,
        onError: onErrorMock,
      };
      
      const { result } = renderHook(
        () => useEntityUpdate('product', mockUpdateFn, customOptions),
        { wrapper }
      );

      await act(async () => {
        result.current.mutate({ id: '1', name: 'Updated Product' });
      });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      // React Query v5 callback signature - includes undefined context parameter
      expect(onSuccessMock).toHaveBeenCalledWith(
        mockUpdateData,
        { id: '1', name: 'Updated Product' },
        undefined
      );
      expect(onErrorMock).not.toHaveBeenCalled();
    });

    test('should call onError when mutation fails', async () => {
      const mockError = new Error('Update failed');
      const mockUpdateFn = jest.fn().mockRejectedValue(mockError);
      const onErrorMock = jest.fn();
      
      const customOptions = {
        onError: onErrorMock,
      };
      
      const { result } = renderHook(
        () => useEntityUpdate('product', mockUpdateFn, customOptions),
        { wrapper }
      );

      await act(async () => {
        result.current.mutate({ id: '1', name: 'Updated Product' });
      });

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });

      // React Query v5 callback signature - includes undefined context parameter
      expect(onErrorMock).toHaveBeenCalledWith(
        mockError,
        { id: '1', name: 'Updated Product' },
        undefined
      );
    });

    test('should reset mutation state', async () => {
      const mockUpdateData = { id: '1', name: 'Updated Product' };
      const mockUpdateFn = jest.fn().mockResolvedValue(mockUpdateData);
      
      const { result } = renderHook(
        () => useEntityUpdate('product', mockUpdateFn),
        { wrapper }
      );

      // Execute mutation
      await act(async () => {
        result.current.mutate({ id: '1', name: 'Updated Product' });
      });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      // Verify success state before reset
      expect(result.current.isSuccess).toBe(true);
      expect(result.current.data).toEqual(mockUpdateData);

      // Reset
      act(() => {
        result.current.reset();
      });

      // Wait for reset to complete
      await waitFor(() => {
        expect(result.current.isSuccess).toBe(false);
      });

      // After reset, the state should be back to initial
      expect(result.current.isError).toBe(false);
      expect(result.current.data).toBeUndefined();
      expect(result.current.error).toBe(null);
    });
  });

  describe('Integration tests', () => {
    test('should work together for fetch and update operations', async () => {
      const initialData = { id: '1', name: 'Original Product', version: 1 };
      const updatedData = { id: '1', name: 'Updated Product', version: 2 };
      
      const mockFetchFn = jest.fn().mockResolvedValue(initialData);
      const mockUpdateFn = jest.fn().mockResolvedValue(updatedData);
      
      const { result: fetchResult } = renderHook(
        () => useEntity('product', '1', mockFetchFn),
        { wrapper }
      );
      
      const { result: updateResult } = renderHook(
        () => useEntityUpdate('product', mockUpdateFn),
        { wrapper: createWrapper() }
      );

      // Wait for initial fetch
      await waitFor(() => {
        expect(fetchResult.current.isSuccess).toBe(true);
      });

      expect(fetchResult.current.data).toEqual(initialData);

      // Perform update
      await act(async () => {
        updateResult.current.mutate({ id: '1', name: 'Updated Product' });
      });

      await waitFor(() => {
        expect(updateResult.current.isSuccess).toBe(true);
      });

      expect(updateResult.current.data).toEqual(updatedData);
      expect(mockFetchFn).toHaveBeenCalledWith('1');
      expect(mockUpdateFn).toHaveBeenCalledWith({ id: '1', name: 'Updated Product' });
    });
  });
}); 