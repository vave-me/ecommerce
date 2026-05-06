import { render, screen, waitFor, act } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import React from 'react';
import Feed from '@/components/Feed/Feed.client.jsx';
import { useFeedQuery } from '@/hooks/queries/useFeedQuery.jsx';
import { queryClient as appQueryClient } from '../../lib/reactQuery.jsx';
import { useInView } from 'react-intersection-observer';

// Mock the hooks and modules
jest.mock('@/hooks/queries/useFeedQuery.jsx');
jest.mock('react-intersection-observer');
jest.mock('../../lib/reactQuery.jsx', () => ({
  queryClient: {
    prefetchQuery: jest.fn()
  }
}));

// Mock the card components
jest.mock('../../app/[locale]/design/PostCard1.jsx', () => ({
  __esModule: true,
  default: ({ post }) => <div data-testid="post-card">{post.title}</div>
}));

jest.mock('../../src/components/services/ServiceCard.jsx', () => ({
  __esModule: true,
  default: ({ service }) => <div data-testid="service-card">{service.title}</div>
}));

jest.mock('../../app/[locale]/design/JobCard.jsx', () => ({
  __esModule: true,
  default: ({ job }) => <div data-testid="job-card">{job.title}</div>
}));

jest.mock('../../app/[locale]/design/TweetCard.jsx', () => ({
  __esModule: true,
  default: ({ tweet }) => <div data-testid="tweet-card">{tweet.title}</div>
}));

jest.mock('../../app/[locale]/design/vehicleCard1.jsx', () => ({
  __esModule: true,
  default: ({ vehicle }) => <div data-testid="vehicle-card">{vehicle.title}</div>
}));

jest.mock('../../app/[locale]/design/property/page.jsx', () => ({
  __esModule: true,
  default: ({ property }) => <div data-testid="property-card">{property.title}</div>
}));

jest.mock('../../app/[locale]/design/deal/page.jsx', () => ({
  __esModule: true,
  default: ({ deal }) => <div data-testid="deal-card">{deal.title}</div>
}));

jest.mock('../../app/[locale]/design/shortsSingle/page.jsx', () => ({
  __esModule: true,
  default: ({ video }) => <div data-testid="video-card">{video.title}</div>
}));

// Create a wrapper with QueryClientProvider
const createWrapper = () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
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

describe('Feed Component with React Query Integration', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    
    // Mock useInView to trigger pagination easily
    useInView.mockReturnValue({
      ref: (el) => el,
      inView: false
    });
    
    // Setup default mock for useFeedQuery
    useFeedQuery.mockReturnValue({
      data: {
        pages: [
          {
            items: [
              { id: '1', type: 'post', title: 'First Post' },
              { id: '2', type: 'service', title: 'Service Offering' }
            ],
            hasMore: true,
            page: 1
          }
        ]
      },
      fetchNextPage: jest.fn(),
      hasNextPage: true,
      isFetchingNextPage: false,
      status: 'success'
    });
  });

  test('Feed renders items based on their type', async () => {
    render(<Feed />, { wrapper: createWrapper() });
    
    await waitFor(() => {
      expect(screen.getByTestId('post-card')).toHaveTextContent('First Post');
      expect(screen.getByTestId('service-card')).toHaveTextContent('Service Offering');
    });
  });

  test('Feed handles all item types correctly', async () => {
    // Update the mock to include all types of items
    useFeedQuery.mockReturnValue({
      data: {
        pages: [
          {
            items: [
              { id: '1', type: 'post', title: 'Blog Post' },
              { id: '2', type: 'service', title: 'Service Offering' },
              { id: '3', type: 'job', title: 'Job Listing' },
              { id: '4', type: 'tweet', title: 'Tweet Post' },
              { id: '5', type: 'vehicle', title: 'Car Listing' },
              { id: '6', type: 'video', title: 'Video Content' },
              { id: '7', type: 'property', title: 'Property Listing' },
              { id: '8', type: 'deal', title: 'Deal Offer' }
            ],
            hasMore: false,
            page: 1
          }
        ]
      },
      fetchNextPage: jest.fn(),
      hasNextPage: false,
      isFetchingNextPage: false,
      status: 'success'
    });
    
    render(<Feed />, { wrapper: createWrapper() });
    
    // Check that each type renders correctly
    await waitFor(() => {
      expect(screen.getByTestId('post-card')).toBeInTheDocument();
      expect(screen.getByTestId('service-card')).toBeInTheDocument();
      expect(screen.getByTestId('job-card')).toBeInTheDocument();
      expect(screen.getByTestId('tweet-card')).toBeInTheDocument();
      expect(screen.getByTestId('vehicle-card')).toBeInTheDocument();
      expect(screen.getByTestId('video-card')).toBeInTheDocument();
      expect(screen.getByTestId('property-card')).toBeInTheDocument();
      expect(screen.getByTestId('deal-card')).toBeInTheDocument();
    });
  });
  
  test('Feed triggers fetchNextPage when intersection observer detects visibility', async () => {
    const mockFetchNextPage = jest.fn();
    
    useFeedQuery.mockReturnValue({
      data: {
        pages: [{ items: [{ id: '1', type: 'post', title: 'Test Post' }], hasMore: true, page: 1 }]
      },
      fetchNextPage: mockFetchNextPage,
      hasNextPage: true,
      isFetchingNextPage: false,
      status: 'success'
    });
    
    render(<Feed />, { wrapper: createWrapper() });
    
    // Initially fetch should not be called
    expect(mockFetchNextPage).not.toHaveBeenCalled();
    
    // Now simulate the intersection observer triggering
    useInView.mockReturnValue({
      ref: (el) => el,
      inView: true
    });
    
    // Re-render with updated inView value
    render(<Feed />, { wrapper: createWrapper() });
    
    // The effect should trigger fetchNextPage
    await waitFor(() => {
      expect(mockFetchNextPage).toHaveBeenCalled();
    });
  });
  
  test('Feed shows loading state when initially loading', async () => {
    useFeedQuery.mockReturnValue({
      status: 'loading',
      data: undefined
    });
    
    render(<Feed />, { wrapper: createWrapper() });
    
    expect(screen.getByText('Loading feed...')).toBeInTheDocument();
  });
  
  test('Feed shows error state when query fails', async () => {
    useFeedQuery.mockReturnValue({
      status: 'error',
      error: new Error('Failed to load feed'),
      data: undefined
    });
    
    render(<Feed />, { wrapper: createWrapper() });
    
    expect(screen.getByText('Error loading feed. Please try again.')).toBeInTheDocument();
  });
  
  test('Feed applies filters from props', async () => {
    render(<Feed category="cars" type="vehicle" limit={20} />, { wrapper: createWrapper() });
    
    // Check that the useFeedQuery was called with the right filters
    expect(useFeedQuery).toHaveBeenCalledWith({ 
      category: "cars", 
      type: "vehicle", 
      limit: 20 
    });
  });
  
  test('Feed shows loading indicator when fetching next page', async () => {
    useFeedQuery.mockReturnValue({
      data: {
        pages: [{ items: [{ id: '1', type: 'post', title: 'Test Post' }], hasMore: true, page: 1 }]
      },
      fetchNextPage: jest.fn(),
      hasNextPage: true,
      isFetchingNextPage: true, // Set to true to show loading indicator
      status: 'success'
    });
    
    render(<Feed />, { wrapper: createWrapper() });
    
    expect(screen.getByText('Loading more...')).toBeInTheDocument();
  });
  
  test('Feed correctly combines data from multiple pages', async () => {
    useFeedQuery.mockReturnValue({
      data: {
        pages: [
          { 
            items: [{ id: '1', type: 'post', title: 'Page 1 Post' }], 
            hasMore: true, 
            page: 1 
          },
          { 
            items: [{ id: '2', type: 'post', title: 'Page 2 Post' }], 
            hasMore: false, 
            page: 2 
          }
        ]
      },
      fetchNextPage: jest.fn(),
      hasNextPage: false,
      isFetchingNextPage: false,
      status: 'success'
    });
    
    render(<Feed />, { wrapper: createWrapper() });
    
    // Both posts from different pages should be rendered
    await waitFor(() => {
      const posts = screen.getAllByTestId('post-card');
      expect(posts.length).toBe(2);
      expect(posts[0]).toHaveTextContent('Page 1 Post');
      expect(posts[1]).toHaveTextContent('Page 2 Post');
    });
  });
  
  test('Feed does not attempt to load more pages when hasNextPage is false', async () => {
    const mockFetchNextPage = jest.fn();
    
    useFeedQuery.mockReturnValue({
      data: {
        pages: [{ items: [{ id: '1', type: 'post', title: 'Test Post' }], hasMore: false, page: 1 }]
      },
      fetchNextPage: mockFetchNextPage,
      hasNextPage: false, // No more pages to fetch
      isFetchingNextPage: false,
      status: 'success'
    });
    
    // Simulate intersection observer trigger
    useInView.mockReturnValue({
      ref: (el) => el,
      inView: true
    });
    
    render(<Feed />, { wrapper: createWrapper() });
    
    // Function should not be called since hasNextPage is false
    expect(mockFetchNextPage).not.toHaveBeenCalled();
  });
  
  test('Feed handles race conditions when loading changes during render', async () => {
    // First return loading state
    useFeedQuery.mockReturnValueOnce({
      status: 'loading',
      data: undefined
    });
    
    const { rerender } = render(<Feed />, { wrapper: createWrapper() });
    
    // Should show loading state
    expect(screen.getByText('Loading feed...')).toBeInTheDocument();
    
    // Then return success state
    useFeedQuery.mockReturnValue({
      data: {
        pages: [{ items: [{ id: '1', type: 'post', title: 'Race Condition Test' }], hasMore: false, page: 1 }]
      },
      fetchNextPage: jest.fn(),
      hasNextPage: false,
      isFetchingNextPage: false,
      status: 'success'
    });
    
    // Rerender to simulate React Query updating its state
    rerender(<Feed />);
    
    // Should now show data
    await waitFor(() => {
      expect(screen.getByTestId('post-card')).toHaveTextContent('Race Condition Test');
    });
  });
  
  test('Feed gracefully handles empty data pages', async () => {
    useFeedQuery.mockReturnValue({
      data: {
        pages: [
          { items: [], hasMore: false, page: 1 }
        ]
      },
      fetchNextPage: jest.fn(),
      hasNextPage: false,
      isFetchingNextPage: false,
      status: 'success'
    });
    
    const { container } = render(<Feed />, { wrapper: createWrapper() });
    
    // There should be no card elements rendered
    expect(screen.queryByTestId('post-card')).not.toBeInTheDocument();
    expect(screen.queryByTestId('service-card')).not.toBeInTheDocument();
    // But the feed container should still be rendered
    expect(container.querySelector('.feed')).toBeInTheDocument();
  });
}); 