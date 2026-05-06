import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import Feed from '../../../src/components/Feed/Feed.client';
import { FeedContext } from '../../../src/components/Feed/FeedProvider.client';
import mockFeedData from '../../../__mocks__/api/mockFeedData';

// Mock components that aren't necessary for our tests
jest.mock('../../../src/app/[locale]/design/PostCard1', () => {
  return function MockPostCard({ post }) {
    return <div data-testid={`post-card-${post.id}`}>{post.title}</div>;
  };
});

jest.mock('../../../src/components/services/ServiceCard', () => {
  return function MockServiceCard({ service }) {
    return <div data-testid={`service-card-${service.id}`}>{service.title}</div>;
  };
});

jest.mock('../../../src/app/[locale]/design/JobCard', () => {
  return function MockJobCard({ job }) {
    return <div data-testid={`job-card-${job.id}`}>{job.title}</div>;
  };
});

jest.mock('../../../src/app/[locale]/design/TweetCard', () => {
  return function MockTweetCard({ tweet }) {
    return <div data-testid={`tweet-card-${tweet.id}`}>{tweet.title}</div>;
  };
});

jest.mock('../../../src/app/[locale]/design/vehicleCard1', () => {
  return function MockVehicleCard({ vehicle }) {
    return <div data-testid={`vehicle-card-${vehicle.id}`}>{vehicle.title}</div>;
  };
});

jest.mock('../../../src/app/[locale]/design/videos/page', () => {
  return function MockVideoCard({ video }) {
    return <div data-testid={`video-card-${video.id}`}>{video.title}</div>;
  };
});

jest.mock('../../../src/components/property/PropertyCard', () => {
  return function MockPropertyCard({ property }) {
    return <div data-testid={`property-card-${property.id}`}>{property.title}</div>;
  };
});

jest.mock('../../../src/components/deals/DealCard', () => {
  return function MockDealCard({ deal }) {
    return <div data-testid={`deal-card-${deal.id}`}>{deal.title}</div>;
  };
});

jest.mock('../../../src/app/[locale]/design/shortsSingle/page', () => {
  return function MockShortVideo({ video }) {
    return <div data-testid={`short-video-${video.id}`}>{video.title}</div>;
  };
});

// Create a wrapper component that provides mock context values
const renderWithFeedContext = (ui, contextValue) => {
  return render(
    <FeedContext.Provider value={contextValue}>
      {ui}
    </FeedContext.Provider>
  );
};

describe('Feed Component', () => {
  it('should render loading state when isLoading is true and no items', () => {
    renderWithFeedContext(<Feed />, {
      feedItems: [],
      isLoading: true,
      hasMore: false,
      loadMore: jest.fn(),
    });

    expect(screen.getByText('Loading feed...')).toBeInTheDocument();
  });

  it('should render empty state when no items and not loading', () => {
    renderWithFeedContext(<Feed />, {
      feedItems: [],
      isLoading: false,
      hasMore: false,
      loadMore: jest.fn(),
    });

    expect(screen.getByText('No items found')).toBeInTheDocument();
    expect(screen.getByText('Try adjusting your filters or search criteria')).toBeInTheDocument();
  });

  it('should render feed items when they are available', () => {
    // Use our mock feed data
    const feedItems = mockFeedData.allFeedItems;
    
    renderWithFeedContext(<Feed />, {
      feedItems,
      isLoading: false,
      hasMore: true,
      loadMore: jest.fn(),
    });

    // We should have an element for each feed item
    expect(screen.getAllByTestId(/item-/)).toHaveLength(feedItems.length);
  });

  it('should show loading indicator when loading more items', () => {
    const feedItems = mockFeedData.allFeedItems.slice(0, 5);
    
    renderWithFeedContext(<Feed />, {
      feedItems,
      isLoading: true,
      hasMore: true,
      loadMore: jest.fn(),
    });

    // Should show items and loading indicator
    expect(screen.getAllByTestId(/item-/)).toHaveLength(feedItems.length);
    expect(screen.getByText('Loading more...')).toBeInTheDocument();
  });

  it('should not call loadMore when hasMore is false', async () => {
    const loadMore = jest.fn();
    const feedItems = mockFeedData.allFeedItems.slice(0, 5);
    
    renderWithFeedContext(<Feed />, {
      feedItems,
      isLoading: false,
      hasMore: false,
      loadMore,
    });

    // Wait for any async effects to complete
    await waitFor(() => {
      expect(loadMore).not.toHaveBeenCalled();
    });
  });

  it('should not show loading more indicator when hasMore is false', () => {
    const feedItems = mockFeedData.allFeedItems.slice(0, 5);
    
    renderWithFeedContext(<Feed />, {
      feedItems,
      isLoading: false,
      hasMore: false,
      loadMore: jest.fn(),
    });

    // Should not show loading more indicator
    expect(screen.queryByText('Loading more...')).not.toBeInTheDocument();
  });

  it('should render items with the correct entity type', () => {
    // Use specific entity type items for clearer testing
    const productItems = mockFeedData.allFeedItems.filter(item => item.entityType === 'product');
    
    renderWithFeedContext(<Feed />, {
      feedItems: productItems,
      isLoading: false,
      hasMore: false,
      loadMore: jest.fn(),
    });

    // Should render all product items
    expect(screen.getAllByTestId(/item-product/)).toHaveLength(productItems.length);
  });
}); 