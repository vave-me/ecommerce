import React from 'react';
import { render, screen, fireEvent, waitFor, act } from '@testing-library/react';
import '@testing-library/jest-dom';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import { listingFiltersReducer } from '../../../src/redux/slices/filterSlice';
import Leftside from '../../../src/components/Leftside/Leftside';
import Feed from '../../../src/components/Feed/Feed.client';
import FeedProvider from '../../../src/components/Feed/FeedProvider.client';
import { QueryClientProvider } from '@tanstack/react-query';
import { queryClient } from '../../../src/lib/reactQuery';

// Mock the searchApi module
jest.mock('../../../src/api/searchApi', () => {
  return require('../../../__mocks__/api/searchApi');
});

// Mock components that aren't necessary for our tests
jest.mock('../../../src/components/Feed/FeedItem.client', () => {
  return function MockFeedItem({ item }) {
    return (
      <div data-testid={`feed-item-${item.entityType}`}>
        <div data-testid="item-entity-type">{item.entityType}</div>
        <div data-testid="item-id">{item[item.entityType]?.id}</div>
        <div data-testid="item-location">{item[item.entityType]?.location}</div>
        <div data-testid="item-tags">{JSON.stringify(item[item.entityType]?.tags)}</div>
      </div>
    );
  };
});

// Mock the next-intl provider since we don't need actual translations for these tests
jest.mock('next-intl', () => ({
  useTranslations: () => (key) => key,
}));

// Performance monitoring is not essential for tests
jest.mock('../../../src/utils/performance', () => ({
  useRenderCount: jest.fn(),
}));

// Mock TypeFilters components, as they're not essential for our tests
jest.mock('../../../src/components/Leftside/TypeFilters/VideoFilters', () => () => <div>VideoFilters</div>);
jest.mock('../../../src/components/Leftside/TypeFilters/TweetFilters', () => () => <div>TweetFilters</div>);
jest.mock('../../../src/components/Leftside/TypeFilters/JobFilters', () => () => <div>JobFilters</div>);
jest.mock('../../../src/components/Leftside/TypeFilters/ServiceFilters', () => () => <div>ServiceFilters</div>);
jest.mock('../../../src/components/Leftside/TypeFilters/DealFilters', () => () => <div>DealFilters</div>);
jest.mock('../../../src/components/Leftside/TypeFilters/ShortFilters', () => () => <div>ShortFilters</div>);
jest.mock('../../../src/components/Leftside/TypeFilters/VehicleFilters', () => () => <div>VehicleFilters</div>);
jest.mock('../../../src/components/Leftside/TypeFilters/PropertyFilters', () => () => <div>PropertyFilters</div>);
jest.mock('../../../src/components/Leftside/TypeFilters/NewsFilters', () => () => <div>NewsFilters</div>);
jest.mock('../../../src/components/Leftside/AdvancedFilters', () => () => <div>AdvancedFilters</div>);

// Setup a test Redux store
const createTestStore = () => {
  return configureStore({
    reducer: {
      listingFilters: listingFiltersReducer,
    },
  });
};

describe('Leftside and Feed Integration', () => {
  let store;

  beforeEach(() => {
    jest.clearAllMocks();
    store = createTestStore();
  });

  const renderTestComponents = () => {
    return render(
      <Provider store={store}>
        <QueryClientProvider client={queryClient}>
          <FeedProvider initialParams={{
            feed_type: "latest",
            entity_types: ["product", "post", "deal", "vehicle", "property", "service", "job"],
            page: 1,
            page_size: 30
          }}>
            <div style={{ display: 'flex' }}>
              <div style={{ width: '300px' }}>
                <Leftside />
              </div>
              <div style={{ flex: 1 }}>
                <Feed />
              </div>
            </div>
          </FeedProvider>
        </QueryClientProvider>
      </Provider>
    );
  };

  it('should render both components without errors', async () => {
    await act(async () => {
      renderTestComponents();
    });

    // Check if both components are rendered
    expect(screen.getByText('contentTypeLabel')).toBeInTheDocument(); // From Leftside
    
    // Wait for feed items to load
    await waitFor(() => {
      expect(screen.getAllByTestId(/feed-item-/)).toHaveLength(15); // All mock items
    });
  });

  it('should filter feed items when location filter is applied', async () => {
    await act(async () => {
      renderTestComponents();
    });

    // Wait for initial feed items to load
    await waitFor(() => {
      expect(screen.getAllByTestId(/feed-item-/)).toHaveLength(15);
    });

    // Get initial item count for Berlin location
    const initialBerlinItems = screen.getAllByTestId('item-location')
      .filter(el => el.textContent === 'Berlin').length;

    // Enter "Berlin" in the location input
    const locationInput = screen.getByPlaceholderText('locationPlaceholder');
    await act(async () => {
      fireEvent.change(locationInput, { target: { value: 'Berlin' } });
    });

    // Click the apply button
    const applyButton = screen.getByText('applyButton');
    await act(async () => {
      fireEvent.click(applyButton);
    });

    // Wait for filtered feed items to load
    await waitFor(() => {
      // All items should now be from Berlin
      const berlinItems = screen.getAllByTestId('item-location')
        .filter(el => el.textContent === 'Berlin').length;
      
      // All displayed items should be Berlin items
      expect(berlinItems).toEqual(screen.getAllByTestId(/feed-item-/).length);
      
      // And there should be fewer items than before
      expect(screen.getAllByTestId(/feed-item-/).length).toBeLessThan(15);
    });
  });

  it('should filter feed items when content type filter is applied', async () => {
    await act(async () => {
      renderTestComponents();
    });

    // Wait for initial feed items to load
    await waitFor(() => {
      expect(screen.getAllByTestId(/feed-item-/)).toHaveLength(15);
    });

    // Select "videos" in the content type dropdown
    const contentTypeSelect = screen.getByLabelText('contentTypeLabel');
    await act(async () => {
      fireEvent.change(contentTypeSelect, { target: { value: 'videos' } });
    });

    // Click the apply button
    const applyButton = screen.getByText('applyButton');
    await act(async () => {
      fireEvent.click(applyButton);
    });

    // Wait for filtered feed items to load
    await waitFor(() => {
      // Only video items should remain
      const entityTypes = screen.getAllByTestId('item-entity-type')
        .map(el => el.textContent);
      
      // All displayed items should be of type 'video'
      expect(entityTypes.every(type => type === 'video')).toBeTruthy();
    });
  });

  it('should filter feed items when tag filter is applied', async () => {
    await act(async () => {
      renderTestComponents();
    });

    // Wait for initial feed items to load
    await waitFor(() => {
      expect(screen.getAllByTestId(/feed-item-/)).toHaveLength(15);
    });

    // Enter "tech" in the tags input
    const tagsInput = screen.getByPlaceholderText('tagsPlaceholder');
    await act(async () => {
      fireEvent.change(tagsInput, { target: { value: 'tech' } });
    });

    // Click the apply button
    const applyButton = screen.getByText('applyButton');
    await act(async () => {
      fireEvent.click(applyButton);
    });

    // Wait for filtered feed items to load
    await waitFor(() => {
      // Only items with 'tech' tag should remain
      const tagElements = screen.getAllByTestId('item-tags');
      
      // Check that each remaining item has the 'tech' tag
      tagElements.forEach(tagEl => {
        const tags = JSON.parse(tagEl.textContent);
        expect(tags.some(tag => tag.toLowerCase().includes('tech'))).toBeTruthy();
      });
      
      // And there should be fewer items than before
      expect(screen.getAllByTestId(/feed-item-/).length).toBeLessThan(15);
    });
  });

  it('should clear filters when reset button is clicked', async () => {
    await act(async () => {
      renderTestComponents();
    });

    // Wait for initial feed items to load
    await waitFor(() => {
      expect(screen.getAllByTestId(/feed-item-/)).toHaveLength(15);
    });

    // Apply a filter first
    const locationInput = screen.getByPlaceholderText('locationPlaceholder');
    await act(async () => {
      fireEvent.change(locationInput, { target: { value: 'Berlin' } });
    });

    // Click the apply button
    const applyButton = screen.getByText('applyButton');
    await act(async () => {
      fireEvent.click(applyButton);
    });

    // Wait for filtered feed items to load
    await waitFor(() => {
      const berlinItems = screen.getAllByTestId('item-location')
        .filter(el => el.textContent === 'Berlin').length;
      expect(berlinItems).toEqual(screen.getAllByTestId(/feed-item-/).length);
    });

    // Get the filtered item count
    const filteredItemCount = screen.getAllByTestId(/feed-item-/).length;

    // Click the clear button
    const clearButton = screen.getByText('clearButton');
    await act(async () => {
      fireEvent.click(clearButton);
    });

    // Wait for reset feed items to load
    await waitFor(() => {
      // There should be more items after reset
      expect(screen.getAllByTestId(/feed-item-/).length).toBeGreaterThan(filteredItemCount);
      
      // Location input should be reset
      expect(screen.getByPlaceholderText('locationPlaceholder').value).toBe('');
    });
  });
}); 