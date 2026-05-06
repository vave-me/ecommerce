import React from 'react';
import { render, act } from '@testing-library/react';
import { renderHook, waitFor } from '@testing-library/react';
import { 
  useCategories, 
  useTopicNavigation, 
  CategoriesProvider,
  TOPIC_CONFIG 
} from '../useCategories';
import * as categoryApi from '../../api/categories';

// Mock the API module
jest.mock('../../api/categories', () => ({
  getMainCategories: jest.fn(),
  prefetchCategoryTypes: jest.fn().mockResolvedValue([]),
  clearCategoryCache: jest.fn()
}));

// Mock next-intl hooks
jest.mock('next-intl', () => ({
  useLocale: jest.fn().mockReturnValue('en')
}));

// Test categories data
const mockCategories = [
  { id: 'cat1', name: 'Electronics' },
  { id: 'cat2', name: 'Clothing' }
];

// Wrapper component to provide context for hooks
const TestWrapper = ({ children, prefetchTopics = [] }) => (
  <CategoriesProvider prefetchTopics={prefetchTopics}>
    {children}
  </CategoriesProvider>
);

describe('Categories Context and Hooks', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    categoryApi.getMainCategories.mockResolvedValue({
      categories: mockCategories
    });
    categoryApi.prefetchCategoryTypes.mockResolvedValue([
      { status: 'fulfilled', value: { categories: mockCategories } }
    ]);
  });

  describe('useCategories hook', () => {
    test('should throw error if used outside provider', () => {
      // Suppress expected error console
      const consoleError = jest.spyOn(console, 'error').mockImplementation(() => {});
      
      // Trying to use hook outside provider should throw
      expect(() => {
        renderHook(() => useCategories());
      }).toThrow('useCategories must be used within a CategoriesProvider');
      
      consoleError.mockRestore();
    });

    test('should provide access to categories data', async () => {
      const { result } = renderHook(() => useCategories(), {
        wrapper: TestWrapper
      });
      
      // Check initial state
      expect(result.current.categories).toEqual({});
      expect(result.current.loading).toEqual({});
      expect(result.current.errors).toEqual({});
      
      // Fetch categories for a topic
      await act(async () => {
        await result.current.getCategoriesForTopic('marketplace');
      });
      
      // Should call the API with correct params
      expect(categoryApi.getMainCategories).toHaveBeenCalledWith({
        categoryType: 'marketplace',
        lang: 'en',
        page: 0,
        pageSize: 50
      });
      
      // Should update categories in context
      expect(result.current.categories).toEqual({
        marketplace: mockCategories
      });
    });

    test('should handle API errors', async () => {
      // Mock API to reject
      const error = new Error('API failure');
      categoryApi.getMainCategories.mockRejectedValueOnce(error);
      
      const { result } = renderHook(() => useCategories(), {
        wrapper: TestWrapper
      });
      
      await act(async () => {
        const categories = await result.current.getCategoriesForTopic('marketplace');
        // Function should return empty array on error
        expect(categories).toEqual([]);
      });
      
      // Should record error
      expect(result.current.errors).toEqual({
        marketplace: 'API failure'
      });
    });

    test('should prefetch topics specified in provider', async () => {
      renderHook(() => useCategories(), {
        wrapper: ({ children }) => (
          <TestWrapper prefetchTopics={['dashboard', 'marketplace']}>
            {children}
          </TestWrapper>
        )
      });
      
      // Should call prefetch API
      await waitFor(() => {
        expect(categoryApi.prefetchCategoryTypes).toHaveBeenCalledWith(
          ['dashboard', 'marketplace'], 
          'en'
        );
      });
    });

    test('should clear cache', async () => {
      const { result } = renderHook(() => useCategories(), {
        wrapper: TestWrapper
      });
      
      // First load some categories
      await act(async () => {
        await result.current.getCategoriesForTopic('marketplace');
      });
      
      // Verify they're in cache
      expect(result.current.categories.marketplace).toEqual(mockCategories);
      
      // Clear cache
      act(() => {
        result.current.clearCache();
      });
      
      // Cache should be empty
      expect(result.current.categories).toEqual({});
      expect(categoryApi.clearCategoryCache).toHaveBeenCalled();
    });
  });

  describe('useTopicNavigation hook', () => {
    test('should initialize with default topic', () => {
      const onChange = jest.fn();
      
      const { result } = renderHook(
        () => useTopicNavigation({ defaultTopic: 'marketplace', onChange }),
        { wrapper: TestWrapper }
      );
      
      expect(result.current.activeTopic).toBe('marketplace');
      expect(result.current.topics).toEqual(TOPIC_CONFIG);
      expect(result.current.getTopicByValue('marketplace')).toEqual(
        TOPIC_CONFIG.find(t => t.value === 'marketplace')
      );
    });

    test('should change topic', async () => {
      const onChange = jest.fn();
      
      const { result } = renderHook(
        () => useTopicNavigation({ defaultTopic: 'dashboard', onChange }),
        { wrapper: TestWrapper }
      );
      
      // Change topic
      await act(async () => {
        result.current.setActiveTopic('marketplace');
      });
      
      // Topic should be updated
      expect(result.current.activeTopic).toBe('marketplace');
      expect(onChange).toHaveBeenCalledWith('marketplace');
    });

    test('should select category', () => {
      const onChange = jest.fn();
      
      const { result } = renderHook(
        () => useTopicNavigation({ defaultTopic: 'marketplace', onChange }),
        { wrapper: TestWrapper }
      );
      
      // Select a category
      act(() => {
        result.current.selectCategory('cat1');
      });
      
      // Category should be updated
      expect(result.current.selectedCategory).toBe('cat1');
      expect(onChange).toHaveBeenCalledWith('marketplace', 'cat1');
    });
  });
}); 