import { render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import React from 'react';
import { useItemThumbnail } from '@/hooks/queries/useMediaQueries.jsx';
import { getMediaByItem, getAllMediaImages } from '@/api/mediaApi.jsx';

// Import our components
import CarListItemClient from '../../../components/Items/CarListItem.client';
import ServiceListItemClient from '../../../components/Items/ServiceListItem.client';
import DealListItemClient from '../../../components/Items/DealListItem.client';
import JobListItemClient from '../../../components/Items/JobListItem.client';
import PropertyListItemClient from '../../../components/Items/PropertyListItem.client';
import ProductListItemClient from '../../../components/Items/ProductListItem.client';
import PostListItem from '../../../components/Items/PostListItem.client';

// Mock the API functions
jest.mock('@/api/mediaApi.jsx', () => ({
  getMediaByItem: jest.fn(),
  getAllMediaImages: jest.fn()
}));

// Mock the useItemThumbnail hook to isolate component tests from the hook
jest.mock('@/hooks/queries/useMediaQueries.jsx', () => ({
  useItemThumbnail: jest.fn()
}));

// Mock other dependencies
jest.mock('@/hooks/useWishlist.jsx', () => ({
  __esModule: true,
  default: () => ({
    cars: [],
    services: [],
    jobs: [],
    properties: [],
    products: [],
    deals: [],
    toggleItem: jest.fn().mockResolvedValue(undefined)
  })
}));

jest.mock('@/hooks/useBasketActions.jsx', () => ({
  __esModule: true,
  default: () => ({
    handleAddToCart: jest.fn()
  })
}));

jest.mock('@/hooks/useActivityApi.jsx', () => ({
  __esModule: true,
  default: () => ({
    handleLike: jest.fn().mockResolvedValue(undefined),
    handleDislike: jest.fn().mockResolvedValue(undefined)
  })
}));

jest.mock('react-redux', () => ({
  useDispatch: () => jest.fn()
}));

jest.mock('../../context/AuthContext.jsx', () => ({
  useAuth: () => ({
    user: { userId: 'test-user-id' }
  })
}));

// Create a test wrapper with QueryClientProvider
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

describe('List Item Components with React Query', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    
    // Setup default mock for useItemThumbnail
    useItemThumbnail.mockReturnValue({
      thumbnail: '/test-thumbnail.jpg',
      isLoading: false,
      error: null
    });
  });

  test('CarListItem uses useItemThumbnail hook', async () => {
    const carData = {
      id: 'car-123',
      name: 'Test Car',
      basePrice: '10000',
      categoryId: 'cat-123',
      canUseMiddleman: true
    };

    render(<CarListItemClient car={carData} />, { wrapper: createWrapper() });
    
    // Verify the hook was called with correct parameters
    expect(useItemThumbnail).toHaveBeenCalledWith(
      'car-123',
      expect.any(String)
    );
    
    // Verify the image uses the thumbnail from the hook
    await waitFor(() => {
      const image = screen.getByAltText('Test Car');
      expect(image.src).toContain('url=%2Ftest-thumbnail.jpg');
    });
  });

  test('PropertyListItem uses useItemThumbnail hook', async () => {
    const propertyData = {
      id: 'property-123',
      name: 'Test Property',
      price: '500000',
      location: 'Berlin',
      categoryId: 'cat-123',
      canUseMiddleman: true
    };

    render(<PropertyListItemClient property={propertyData} />, { wrapper: createWrapper() });
    
    expect(useItemThumbnail).toHaveBeenCalledWith(
      'property-123',
      expect.any(String)
    );
    
    await waitFor(() => {
      const image = screen.getByAltText('Test Property');
      expect(image.src).toContain('url=%2Ftest-thumbnail.jpg');
    });
  });

  test('ServiceListItem uses useItemThumbnail hook', async () => {
    const serviceData = {
      id: 'service-123',
      name: 'Test Service',
      price: '100',
      categoryId: 'cat-123',
      canUseMiddleman: true
    };

    render(<ServiceListItemClient service={serviceData} />, { wrapper: createWrapper() });
    
    expect(useItemThumbnail).toHaveBeenCalledWith(
      'service-123',
      expect.any(String)
    );
  });

  test('JobListItem uses useItemThumbnail hook', async () => {
    const jobData = {
      id: 'job-123',
      name: 'Test Job',
      salary: '50000',
      categoryId: 'cat-123'
    };

    render(<JobListItemClient job={jobData} />, { wrapper: createWrapper() });
    
    expect(useItemThumbnail).toHaveBeenCalledWith(
      'job-123',
      expect.any(String)
    );
  });

  test('DealListItem uses useItemThumbnail hook', async () => {
    const dealData = {
      id: 'deal-123',
      name: 'Test Deal',
      price: '99',
      categoryId: 'cat-123',
      canUseMiddleman: true
    };

    render(<DealListItemClient deal={dealData} />, { wrapper: createWrapper() });
    
    expect(useItemThumbnail).toHaveBeenCalledWith(
      'deal-123',
      expect.any(String)
    );
  });

  test('ProductListItem uses useItemThumbnail hook', async () => {
    const productData = {
      id: 'product-123',
      name: 'Test Product',
      price: '29.99',
      categoryId: 'cat-123',
      canUseMiddleman: true
    };

    render(<ProductListItemClient product={productData} />, { wrapper: createWrapper() });
    
    expect(useItemThumbnail).toHaveBeenCalledWith(
      'product-123',
      expect.any(String)
    );
  });

  test('PostListItem uses useItemThumbnail hook', async () => {
    const postData = {
      id: 'post-123',
      name: 'Test Post',
      categoryId: 'cat-123'
    };

    render(<PostListItem product={postData} />, { wrapper: createWrapper() });
    
    expect(useItemThumbnail).toHaveBeenCalledWith(
      'post-123',
      expect.any(String)
    );
  });

  test('Components handle loading state correctly', async () => {
    // Update mock to simulate loading state
    useItemThumbnail.mockReturnValue({
      thumbnail: '/placeholder.jpg',
      isLoading: true,
      error: null
    });
    
    const carData = {
      id: 'car-123',
      name: 'Test Car',
      basePrice: '10000',
      categoryId: 'cat-123',
      canUseMiddleman: true
    };

    render(<CarListItemClient car={carData} />, { wrapper: createWrapper() });
    
    // Even in loading state, the component should render with the placeholder
    await waitFor(() => {
      const image = screen.getByAltText('Test Car');
      expect(image.src).toContain('url=%2Fplaceholder.jpg');
    });
  });

  test('Components handle error state gracefully', async () => {
    // Update mock to simulate error state
    useItemThumbnail.mockReturnValue({
      thumbnail: '/fallback-image.jpg',
      isLoading: false,
      error: new Error('Failed to load image')
    });
    
    const carData = {
      id: 'car-123',
      name: 'Test Car',
      basePrice: '10000',
      categoryId: 'cat-123'
    };

    render(<CarListItemClient car={carData} />, { wrapper: createWrapper() });
    
    // Even with an error, the component should render with the fallback
    await waitFor(() => {
      const image = screen.getByAltText('Test Car');
      expect(image.src).toContain('url=%2Ffallback-image.jpg');
    });
  });
}); 