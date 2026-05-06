import React from 'react';
import { render, screen, act } from '@testing-library/react';
import { renderHook, waitFor } from '@testing-library/react';
import { 
  useCategories, 
  CategoriesProvider
} from '@/hooks/useCategories.jsx';
import { QueryWrapper, createTestQueryClient } from '../test-utils.jsx';

// Mock the API functions
jest.mock('@/api/categories.jsx', () => ({
  fetchMainCategories: jest.fn()
}));

// Mock next-intl hooks
jest.mock('next-intl', () => ({
  useLocale: jest.fn().mockReturnValue('en')
}));

import { fetchMainCategories } from '@/api/categories.jsx';

describe('Categories Context and Hooks', () => {
  let queryClient;
  
  beforeEach(() => {
    jest.clearAllMocks();
    queryClient = createTestQueryClient();
    
    // Default mock implementations
    fetchMainCategories.mockResolvedValue({
      categories: [
        { id: 1, name: 'Electronics', slug: 'electronics' },
        { id: 2, name: 'Clothing', slug: 'clothing' },
      ]
    });
  });

  describe('CategoriesProvider', () => {
    test('should render children within provider', () => {
      render(
        <QueryWrapper queryClient={queryClient}>
          <CategoriesProvider>
            <div data-testid="test-child">Test Child</div>
          </CategoriesProvider>
        </QueryWrapper>
      );
      
      expect(screen.getByTestId('test-child')).toBeInTheDocument();
    });

    test('should handle prefetchTopics prop', () => {
      render(
        <QueryWrapper queryClient={queryClient}>
          <CategoriesProvider prefetchTopics={['electronics']}>
            <div data-testid="test-child">Test Child</div>
          </CategoriesProvider>
        </QueryWrapper>
      );
      
      expect(screen.getByTestId('test-child')).toBeInTheDocument();
    });
  });

  describe('useCategories hook', () => {
    const wrapper = ({ children }) => (
      <QueryWrapper queryClient={queryClient}>
        <CategoriesProvider>
          {children}
        </CategoriesProvider>
      </QueryWrapper>
    );

    test('should provide access to categories data', async () => {
      const { result } = renderHook(() => useCategories('all'), { wrapper });
      
      // Initial loading state
      expect(result.current.isLoading).toBe(true);
      expect(result.current.data).toEqual([]);
      
      // Wait for data to load
      await act(async () => {
        await new Promise(resolve => setTimeout(resolve, 0));
      });
      
      // Should have categories data
      expect(fetchMainCategories).toHaveBeenCalledWith({ 
        categoryType: 'all', 
        lang: 'en'
      });
    });

    test('should handle different category types', async () => {
      const { result } = renderHook(() => useCategories('electronics'), { wrapper });
      
      await act(async () => {
        await new Promise(resolve => setTimeout(resolve, 0));
      });
      
      expect(fetchMainCategories).toHaveBeenCalledWith({ 
        categoryType: 'electronics', 
        lang: 'en'
      });
    });

    test('should handle API errors', async () => {
      fetchMainCategories.mockRejectedValue(new Error('API Error'));
      
      const { result } = renderHook(() => useCategories('all'), { wrapper });
      
      await act(async () => {
        await new Promise(resolve => setTimeout(resolve, 0));
      });
      
      expect(result.current.isError).toBe(true);
    });

    test('should prefetch topics specified in provider', async () => {
      const wrapperWithPrefetch = ({ children }) => (
        <QueryWrapper queryClient={queryClient}>
          <CategoriesProvider prefetchTopics={['electronics', 'clothing']}>
            {children}
          </CategoriesProvider>
        </QueryWrapper>
      );

      renderHook(() => useCategories('all'), { wrapper: wrapperWithPrefetch });
      
      await act(async () => {
        await new Promise(resolve => setTimeout(resolve, 100));
      });
      
      // Should prefetch the specified topics
      expect(fetchMainCategories).toHaveBeenCalledWith({ 
        categoryType: 'electronics', 
        lang: 'en'
      });
      expect(fetchMainCategories).toHaveBeenCalledWith({ 
        categoryType: 'clothing', 
        lang: 'en'
      });
    });
  });
}); 