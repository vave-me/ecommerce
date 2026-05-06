import React from 'react';
import { render, screen, fireEvent, waitFor, act } from '@testing-library/react';
import '@testing-library/jest-dom';
import FeedProvider, { useFeed } from '../../../src/components/Feed/FeedProvider.client';

// Mock the searchApi module
jest.mock('../../../src/api/searchApi', () => {
  return require('../../../__mocks__/api/searchApi');
});

// Mock the next-intl provider since we don't need actual translations for these tests
jest.mock('next-intl', () => ({
  useTranslations: () => (key) => key,
}));

// Create a simple test component that consumes the FeedContext
const TestConsumer = () => {
  const { 
    feedItems, 
    isLoading, 
    hasMore, 
    filterParams, 
    updateFilters, 
    loadMore, 
    error 
  } = useFeed();

  return (
    <div>
      <div data-testid="loading-state">{isLoading ? 'Loading' : 'Not Loading'}</div>
      <div data-testid="has-more">{hasMore ? 'Has More' : 'No More'}</div>
      <div data-testid="item-count">{feedItems.length}</div>
      <div data-testid="feed-params">{JSON.stringify(filterParams)}</div>
      {error && <div data-testid="error-message">{error.message}</div>}
      
      <ul>
        {feedItems.map(item => (
          <li key={item.id} data-testid={`item-${item.entityType}`}>
            {item.entityType}: {item[item.entityType]?.id}
          </li>
        ))}
      </ul>
      
      <button 
        data-testid="filter-by-location" 
        onClick={() => updateFilters({ location: 'Berlin' })}
      >
        Filter by Berlin
      </button>
      
      <button 
        data-testid="filter-by-products" 
        onClick={() => updateFilters({ contentType: 'products' })}
      >
        Show Only Products
      </button>
      
      <button 
        data-testid="filter-by-tag" 
        onClick={() => updateFilters({ tags: ['tech'] })}
      >
        Filter by Tech Tag
      </button>
      
      <button 
        data-testid="load-more-button" 
        onClick={loadMore}
      >
        Load More
      </button>
      
      <button 
        data-testid="reset-filters" 
        onClick={() => updateFilters({})}
      >
        Reset Filters
      </button>
    </div>
  );
};

describe('FeedProvider Component', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should initialize and load feed items', async () => {
    // Use act to handle the useEffect in FeedProvider
    act(() => {
      render(
        <FeedProvider>
          <TestConsumer />
        </FeedProvider>
      );
    });

    // Initially shows loading state
    expect(screen.getByTestId('loading-state')).toHaveTextContent('Loading');
    
    // Wait for feed to load
    await waitFor(() => {
      expect(screen.getByTestId('loading-state')).toHaveTextContent('Not Loading');
    });
    
    // Verify items are loaded
    expect(screen.getByTestId('item-count')).not.toHaveTextContent('0');
    
    // Check that we have items of different entity types
    expect(screen.getAllByTestId(/item-product/)).toHaveLength(3);
    expect(screen.getAllByTestId(/item-post/)).toHaveLength(2);
    expect(screen.getAllByTestId(/item-deal/)).toHaveLength(2);
    expect(screen.getAllByTestId(/item-vehicle/)).toHaveLength(2);
    expect(screen.getAllByTestId(/item-property/)).toHaveLength(2);
    expect(screen.getAllByTestId(/item-service/)).toHaveLength(2);
    expect(screen.getAllByTestId(/item-job/)).toHaveLength(2);
  });

  it('should filter feed items by location', async () => {
    act(() => {
      render(
        <FeedProvider>
          <TestConsumer />
        </FeedProvider>
      );
    });

    // Wait for initial feed to load
    await waitFor(() => {
      expect(screen.getByTestId('loading-state')).toHaveTextContent('Not Loading');
    });
    
    // Get the initial item count
    const initialItemCount = screen.getAllByTestId(/item-/).length;
    
    // Filter by location
    act(() => {
      fireEvent.click(screen.getByTestId('filter-by-location'));
    });
    
    // Wait for filtered feed to load
    await waitFor(() => {
      expect(screen.getByTestId('loading-state')).toHaveTextContent('Not Loading');
    });
    
    // Verify that filter was applied
    expect(screen.getByTestId('feed-params')).toHaveTextContent('Berlin');
    
    // Verify that the item count changed (fewer items)
    const filteredItemCount = screen.getAllByTestId(/item-/).length;
    expect(filteredItemCount).toBeLessThan(initialItemCount);
  });

  it('should filter feed items by entity type', async () => {
    act(() => {
      render(
        <FeedProvider>
          <TestConsumer />
        </FeedProvider>
      );
    });

    // Wait for initial feed to load
    await waitFor(() => {
      expect(screen.getByTestId('loading-state')).toHaveTextContent('Not Loading');
    });
    
    // Filter by products only
    act(() => {
      fireEvent.click(screen.getByTestId('filter-by-products'));
    });
    
    // Wait for filtered feed to load
    await waitFor(() => {
      expect(screen.getByTestId('loading-state')).toHaveTextContent('Not Loading');
    });
    
    // Verify that only product items are shown
    expect(screen.getAllByTestId(/item-product/)).toHaveLength(3);
    expect(screen.queryByTestId(/item-post/)).not.toBeInTheDocument();
    expect(screen.queryByTestId(/item-deal/)).not.toBeInTheDocument();
    expect(screen.queryByTestId(/item-vehicle/)).not.toBeInTheDocument();
  });

  it('should filter feed items by tag', async () => {
    act(() => {
      render(
        <FeedProvider>
          <TestConsumer />
        </FeedProvider>
      );
    });

    // Wait for initial feed to load
    await waitFor(() => {
      expect(screen.getByTestId('loading-state')).toHaveTextContent('Not Loading');
    });
    
    // Get the initial item count
    const initialItemCount = screen.getAllByTestId(/item-/).length;
    
    // Filter by tech tag
    act(() => {
      fireEvent.click(screen.getByTestId('filter-by-tag'));
    });
    
    // Wait for filtered feed to load
    await waitFor(() => {
      expect(screen.getByTestId('loading-state')).toHaveTextContent('Not Loading');
    });
    
    // Verify that the item count changed (fewer items)
    const filteredItemCount = screen.getAllByTestId(/item-/).length;
    expect(filteredItemCount).toBeLessThan(initialItemCount);
  });

  it('should load more items when requested', async () => {
    act(() => {
      render(
        <FeedProvider>
          <TestConsumer />
        </FeedProvider>
      );
    });

    // Wait for initial feed to load
    await waitFor(() => {
      expect(screen.getByTestId('loading-state')).toHaveTextContent('Not Loading');
    });
    
    // Get the initial item count
    const initialItemCount = screen.getAllByTestId(/item-/).length;
    
    // Load more items
    act(() => {
      fireEvent.click(screen.getByTestId('load-more-button'));
    });
    
    // Wait for more items to load
    await waitFor(() => {
      expect(screen.getByTestId('loading-state')).toHaveTextContent('Not Loading');
    });
    
    // After loading more, hasMore should be false because our mock returns page 2 with hasMore: false
    expect(screen.getByTestId('has-more')).toHaveTextContent('No More');
  });

  it('should reset filters', async () => {
    act(() => {
      render(
        <FeedProvider>
          <TestConsumer />
        </FeedProvider>
      );
    });

    // Wait for initial feed to load
    await waitFor(() => {
      expect(screen.getByTestId('loading-state')).toHaveTextContent('Not Loading');
    });
    
    // Filter by location
    act(() => {
      fireEvent.click(screen.getByTestId('filter-by-location'));
    });
    
    // Wait for filtered feed to load
    await waitFor(() => {
      expect(screen.getByTestId('loading-state')).toHaveTextContent('Not Loading');
    });
    
    // Verify that filter was applied
    expect(screen.getByTestId('feed-params')).toHaveTextContent('Berlin');
    
    // Get the filtered item count
    const filteredItemCount = screen.getAllByTestId(/item-/).length;
    
    // Reset filters
    act(() => {
      fireEvent.click(screen.getByTestId('reset-filters'));
    });
    
    // Wait for reset feed to load
    await waitFor(() => {
      expect(screen.getByTestId('loading-state')).toHaveTextContent('Not Loading');
    });
    
    // Verify that filter was reset
    expect(screen.getByTestId('feed-params')).not.toHaveTextContent('Berlin');
    
    // Verify that the item count changed back to original count
    const resetItemCount = screen.getAllByTestId(/item-/).length;
    expect(resetItemCount).toBeGreaterThan(filteredItemCount);
  });
}); 