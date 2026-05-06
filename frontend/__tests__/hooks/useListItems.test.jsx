import { renderHook, act } from '@testing-library/react';
import useListItems from '@/hooks/useListItems.jsx';
import { fetchProductsByFilters } from "@/api/productsApi.jsx";
import { fetchPostsByFilters } from "@/api/postsApi.jsx";

// Mock dependencies
jest.mock('react-redux', () => ({
  useSelector: jest.fn()
}));

jest.mock('@tanstack/react-query', () => ({
  useInfiniteQuery: jest.fn()
}));

jest.mock('@/api/productsApi.jsx', () => ({
  fetchProductsByFilters: jest.fn()
}));

jest.mock('@/api/postsApi.jsx', () => ({
  fetchPostsByFilters: jest.fn()
}));

describe('useListItems hook', () => {
  const mockUseSelector = require('react-redux').useSelector;
  const mockUseInfiniteQuery = require('@tanstack/react-query').useInfiniteQuery;
  
  beforeEach(() => {
    jest.clearAllMocks();
    
    // Default mock state
    mockUseSelector.mockReturnValue({
      listingType: 'products',
      category: 'electronics',
      sortBy: 'newest'
    });
    
    // Default useInfiniteQuery response
    mockUseInfiniteQuery.mockReturnValue({
      data: {
        pages: [
          { products: [{ id: '1', name: 'Product 1' }, { id: '2', name: 'Product 2' }] }
        ]
      },
      isLoading: false,
      isError: false,
      error: null,
      fetchNextPage: jest.fn(),
      hasNextPage: true,
      isFetching: false
    });
  });
  
  test('should return products when listingType is products', () => {
    mockUseSelector.mockReturnValue({
      listingType: 'products',
      category: 'electronics'
    });
    
    const { result } = renderHook(() => useListItems());
    
    expect(result.current.items).toEqual([
      { id: '1', name: 'Product 1' }, 
      { id: '2', name: 'Product 2' }
    ]);
    expect(result.current.isLoading).toBe(false);
    expect(result.current.isError).toBe(false);
  });
  
  test('should return posts when listingType is not products', () => {
    mockUseSelector.mockReturnValue({
      listingType: 'posts',
      category: 'news'
    });
    
    mockUseInfiniteQuery.mockReturnValue({
      data: {
        pages: [
          { posts: [{ id: '1', title: 'Post 1' }, { id: '2', title: 'Post 2' }] }
        ]
      },
      isLoading: false,
      isError: false,
      error: null,
      fetchNextPage: jest.fn(),
      hasNextPage: true,
      isFetching: false
    });
    
    const { result } = renderHook(() => useListItems());
    
    expect(result.current.items).toEqual([
      { id: '1', title: 'Post 1' }, 
      { id: '2', title: 'Post 2' }
    ]);
  });
  
  test('should handle loading state correctly', () => {
    mockUseInfiniteQuery.mockReturnValue({
      data: null,
      isLoading: true,
      isError: false,
      error: null,
      fetchNextPage: jest.fn(),
      hasNextPage: false,
      isFetching: true
    });
    
    const { result } = renderHook(() => useListItems());
    
    expect(result.current.items).toEqual([]);
    expect(result.current.isLoading).toBe(true);
    expect(result.current.isFetching).toBe(true);
  });
  
  test('should handle error state correctly', () => {
    const testError = new Error('Failed to fetch');
    
    mockUseInfiniteQuery.mockReturnValue({
      data: null,
      isLoading: false,
      isError: true,
      error: testError,
      fetchNextPage: jest.fn(),
      hasNextPage: false,
      isFetching: false
    });
    
    const { result } = renderHook(() => useListItems());
    
    expect(result.current.items).toEqual([]);
    expect(result.current.isError).toBe(true);
    expect(result.current.error).toBe(testError);
  });
  
  test('should handle empty data', () => {
    mockUseInfiniteQuery.mockReturnValue({
      data: { pages: [] },
      isLoading: false,
      isError: false,
      error: null,
      fetchNextPage: jest.fn(),
      hasNextPage: false,
      isFetching: false
    });
    
    const { result } = renderHook(() => useListItems());
    
    expect(result.current.items).toEqual([]);
  });
}); 