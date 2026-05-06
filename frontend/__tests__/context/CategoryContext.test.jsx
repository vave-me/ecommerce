import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { CategoriesProvider, useCategories } from '@/hooks/useCategories.jsx';
import * as apiModule from '@/api/categories.jsx';

// Mock the API
jest.mock('@/api/categories.jsx', () => ({
  getCategories: jest.fn(),
  getCategoryBySlug: jest.fn(),
  getCategoryById: jest.fn(),
  fetchMainCategories: jest.fn()
}));

// Mock hierarchical category data
const mockCategories = [
  {
    id: '1',
    name: 'Electronics',
    slug: 'electronics',
    children: [
      {
        id: '2',
        name: 'Phones',
        slug: 'phones',
        children: [
          { id: '3', name: 'Smartphones', slug: 'smartphones', children: [] }
        ]
      },
      { id: '4', name: 'Computers', slug: 'computers', children: [] }
    ]
  },
  {
    id: '5',
    name: 'Clothing',
    slug: 'clothing',
    children: []
  }
];

// Create a new QueryClient for each test
const createTestQueryClient = () => new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
    },
  },
});

// Test component that uses the category context
const TestComponent = () => {
  const { 
    categoryData,
    prefetchCategories,
    loading
  } = useCategories();

  // Check for marketplace categories
  const isLoading = categoryData?.marketplace?.isLoading || Object.keys(loading || {}).length > 0;
  const error = categoryData?.marketplace?.error;
  const categories = categoryData?.marketplace?.categories || [];

  return (
    <div>
      <div data-testid="loading">{isLoading ? 'Loading' : 'Not Loading'}</div>
      <div data-testid="error">{error ? 'Error occurred' : 'No Error'}</div>
      <div data-testid="categories">{categories ? categories.length : 0}</div>
    </div>
  );
};

describe('CategoryContext', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    // Default mock implementation
    apiModule.fetchMainCategories.mockResolvedValue({ categories: mockCategories });
  });

  it('initializes with loading state and fetches categories', async () => {
    const queryClient = createTestQueryClient();
    
    render(
      <QueryClientProvider client={queryClient}>
        <CategoriesProvider categoryTypes={['marketplace']}>
          <TestComponent />
        </CategoriesProvider>
      </QueryClientProvider>
    );

    // Initially should be loading
    expect(screen.getByTestId('loading')).toHaveTextContent('Loading');

    // Wait for categories to load
    await waitFor(() => {
      expect(screen.getByTestId('loading')).toHaveTextContent('Not Loading');
    });

    // Should have fetched categories
    expect(apiModule.fetchMainCategories).toHaveBeenCalledTimes(1);
    expect(screen.getByTestId('categories')).toHaveTextContent('2');
  });

  it('handles API errors', async () => {
    // Mock API error
    apiModule.fetchMainCategories.mockRejectedValueOnce(new Error('Failed to fetch categories'));
    
    const queryClient = createTestQueryClient();
    
    render(
      <QueryClientProvider client={queryClient}>
        <CategoriesProvider categoryTypes={['marketplace']}>
          <TestComponent />
        </CategoriesProvider>
      </QueryClientProvider>
    );

    // Wait for error to be displayed
    await waitFor(() => {
      expect(screen.getByTestId('error')).toHaveTextContent('Error occurred');
    }, { timeout: 3000 });

    // Categories should be empty
    expect(screen.getByTestId('categories')).toHaveTextContent('0');
  });
}); 