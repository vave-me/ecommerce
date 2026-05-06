import mockFeedData from './mockFeedData';

// Mock implementation of the searchApi's unifiedFeed function
export const unifiedFeed = jest.fn().mockImplementation((feedParams = {}) => {
  // Process feedParams to determine what mock data to return
  
  // Extract and normalize parameters
  const entityTypes = feedParams.entity_types || [];
  const page = feedParams.page || 1;
  const location = feedParams.location || '';
  const tags = feedParams.tags || '';
  const category = feedParams.category || '';
  
  // Check if we should return the second page
  if (page === 2) {
    return Promise.resolve(mockFeedData.mockUnifiedFeedPage2Response);
  }
  
  // Check if we should filter by entity type
  if (entityTypes.length === 1) {
    const entityType = entityTypes[0];
    if (entityType && entityType !== 'all') {
      return Promise.resolve(mockFeedData.getMockFeedResponseByEntityType(entityType));
    }
  }
  
  // Check if we should filter by location
  if (location) {
    return Promise.resolve(mockFeedData.getMockFeedResponseByLocation(location));
  }
  
  // Check if we should filter by tags
  if (tags) {
    return Promise.resolve(mockFeedData.getMockFeedResponseByTags(tags));
  }
  
  // Default response - return all items
  return Promise.resolve(mockFeedData.mockUnifiedFeedResponse);
});

// Error case for testing error handling
export const unifiedFeedWithError = jest.fn().mockImplementation(() => {
  return Promise.reject(mockFeedData.mockErrorResponse);
});

// Other searchApi functions we might want to mock
export const searchProductsWithFilters = jest.fn().mockImplementation(() => {
  return Promise.resolve({ products: mockFeedData.productItems });
});

export const searchPostsWithFilters = jest.fn().mockImplementation(() => {
  return Promise.resolve({ posts: mockFeedData.postItems });
});

export const searchVehiclesWithFilters = jest.fn().mockImplementation(() => {
  return Promise.resolve({ vehicles: mockFeedData.vehicleItems });
});

export const searchPropertiesWithFilters = jest.fn().mockImplementation(() => {
  return Promise.resolve({ properties: mockFeedData.propertyItems });
});

export const searchServicesWithFilters = jest.fn().mockImplementation(() => {
  return Promise.resolve({ services: mockFeedData.serviceItems });
});

export const searchJobsWithFilters = jest.fn().mockImplementation(() => {
  return Promise.resolve({ jobs: mockFeedData.jobItems });
});

export const searchDealsWithFilters = jest.fn().mockImplementation(() => {
  return Promise.resolve({ deals: mockFeedData.dealItems });
});

// Default export
export default {
  unifiedFeed,
  unifiedFeedWithError,
  searchProductsWithFilters,
  searchPostsWithFilters,
  searchVehiclesWithFilters,
  searchPropertiesWithFilters,
  searchServicesWithFilters,
  searchJobsWithFilters,
  searchDealsWithFilters
}; 