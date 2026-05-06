/**
 * Mock Feed Data
 * Contains dummy data for all entity types in the unified feed
 */

// Helper to create random IDs
const generateId = () => Math.random().toString(36).substring(2, 10);

// Mock product items
const productItems = [
  {
    id: generateId(),
    entityType: 'product',
    product: {
      id: 'prod-001',
      name: 'High-end Laptop',
      description: 'Powerful laptop with latest specs',
      price: 1299.99,
      currency: 'EUR',
      category: 'electronics',
      tags: ['tech', 'laptop', 'gaming'],
      location: 'Berlin',
      condition: 'new',
      createdAt: '2023-10-12T10:30:00Z',
      seller: {
        id: 'user-123',
        name: 'TechStore',
        verified: true,
        rating: 4.8
      },
      thumbnail: 'https://example.com/laptop.jpg'
    },
    createdAt: '2023-10-12T10:30:00Z'
  },
  {
    id: generateId(),
    entityType: 'product',
    product: {
      id: 'prod-002',
      name: 'Vintage Guitar',
      description: 'Beautiful vintage acoustic guitar',
      price: 450,
      currency: 'EUR',
      category: 'music',
      tags: ['music', 'guitar', 'vintage'],
      location: 'Munich',
      condition: 'used',
      createdAt: '2023-10-10T14:20:00Z',
      seller: {
        id: 'user-456',
        name: 'MusicLover',
        verified: false,
        rating: 4.6
      },
      thumbnail: 'https://example.com/guitar.jpg'
    },
    createdAt: '2023-10-10T14:20:00Z'
  },
  {
    id: generateId(),
    entityType: 'product',
    product: {
      id: 'prod-003',
      name: 'Smart Watch',
      description: 'Latest model with health tracking',
      price: 299.99,
      currency: 'EUR',
      category: 'electronics',
      tags: ['tech', 'wearable', 'smartwatch'],
      location: 'Berlin',
      condition: 'new',
      createdAt: '2023-10-15T09:15:00Z',
      seller: {
        id: 'user-123',
        name: 'TechStore',
        verified: true,
        rating: 4.8
      },
      thumbnail: 'https://example.com/smartwatch.jpg'
    },
    createdAt: '2023-10-15T09:15:00Z'
  }
];

// Mock post items
const postItems = [
  {
    id: generateId(),
    entityType: 'post',
    post: {
      id: 'post-001',
      title: 'My Travel Experience',
      content: 'Just got back from an amazing trip to Portugal...',
      likes: 42,
      comments: 7,
      tags: ['travel', 'portugal', 'vacation'],
      location: 'Lisbon',
      createdAt: '2023-10-18T16:45:00Z',
      author: {
        id: 'user-789',
        name: 'TravelExplorer',
        verified: true,
        avatar: 'https://example.com/avatar1.jpg'
      },
      thumbnail: 'https://example.com/portugal.jpg'
    },
    createdAt: '2023-10-18T16:45:00Z'
  },
  {
    id: generateId(),
    entityType: 'post',
    post: {
      id: 'post-002',
      title: 'Tech Review: Latest Smartphone',
      content: 'Here\'s my honest review of the newest flagship phone...',
      likes: 124,
      comments: 35,
      tags: ['tech', 'review', 'smartphone'],
      location: 'Berlin',
      createdAt: '2023-10-17T11:30:00Z',
      author: {
        id: 'user-101',
        name: 'TechReviewer',
        verified: true,
        avatar: 'https://example.com/avatar2.jpg'
      },
      thumbnail: 'https://example.com/smartphone.jpg'
    },
    createdAt: '2023-10-17T11:30:00Z'
  }
];

// Mock deal items
const dealItems = [
  {
    id: generateId(),
    entityType: 'deal',
    deal: {
      id: 'deal-001',
      title: 'Flash Sale: 50% Off Electronics',
      description: 'Limited time offer on all electronics',
      discount: 50,
      originalPrice: 500,
      salePrice: 250,
      currency: 'EUR',
      category: 'electronics',
      tags: ['sale', 'electronics', 'discount'],
      location: 'Frankfurt',
      expiresAt: '2023-11-01T23:59:59Z',
      createdAt: '2023-10-15T08:00:00Z',
      merchant: {
        id: 'merchant-123',
        name: 'ElectroSale',
        verified: true
      },
      thumbnail: 'https://example.com/electronicsale.jpg'
    },
    createdAt: '2023-10-15T08:00:00Z'
  },
  {
    id: generateId(),
    entityType: 'deal',
    deal: {
      id: 'deal-002',
      title: 'Buy One Get One Free: Books',
      description: 'Special promotion on selected books',
      discount: 50,
      originalPrice: 40,
      salePrice: 20,
      currency: 'EUR',
      category: 'books',
      tags: ['books', 'sale', 'bogo'],
      location: 'Hamburg',
      expiresAt: '2023-10-31T23:59:59Z',
      createdAt: '2023-10-14T10:15:00Z',
      merchant: {
        id: 'merchant-456',
        name: 'BookHaven',
        verified: false
      },
      thumbnail: 'https://example.com/booksale.jpg'
    },
    createdAt: '2023-10-14T10:15:00Z'
  }
];

// Mock vehicle items
const vehicleItems = [
  {
    id: generateId(),
    entityType: 'vehicle',
    vehicle: {
      id: 'vehicle-001',
      title: '2020 Tesla Model 3',
      description: 'Fully electric vehicle in excellent condition',
      price: 42000,
      currency: 'EUR',
      category: 'cars',
      tags: ['electric', 'tesla', 'sedan'],
      location: 'Munich',
      make: 'Tesla',
      model: 'Model 3',
      year: 2020,
      mileage: 25000,
      fuelType: 'electric',
      transmission: 'automatic',
      condition: 'excellent',
      createdAt: '2023-10-10T15:30:00Z',
      seller: {
        id: 'user-222',
        name: 'CarDealer',
        verified: true,
        rating: 4.9
      },
      thumbnail: 'https://example.com/tesla.jpg'
    },
    createdAt: '2023-10-10T15:30:00Z'
  },
  {
    id: generateId(),
    entityType: 'vehicle',
    vehicle: {
      id: 'vehicle-002',
      title: '2018 BMW 3 Series',
      description: 'Well-maintained luxury sedan',
      price: 28500,
      currency: 'EUR',
      category: 'cars',
      tags: ['luxury', 'bmw', 'sedan'],
      location: 'Berlin',
      make: 'BMW',
      model: '3 Series',
      year: 2018,
      mileage: 45000,
      fuelType: 'gasoline',
      transmission: 'automatic',
      condition: 'good',
      createdAt: '2023-10-08T09:45:00Z',
      seller: {
        id: 'user-333',
        name: 'LuxuryAutos',
        verified: true,
        rating: 4.7
      },
      thumbnail: 'https://example.com/bmw.jpg'
    },
    createdAt: '2023-10-08T09:45:00Z'
  }
];

// Mock property items
const propertyItems = [
  {
    id: generateId(),
    entityType: 'property',
    property: {
      id: 'property-001',
      title: 'Modern Apartment in City Center',
      description: 'Spacious 2-bedroom apartment with great view',
      price: 1200,
      currency: 'EUR',
      category: 'apartments',
      tags: ['apartment', 'city-center', 'modern'],
      location: 'Berlin',
      type: 'apartment',
      bedrooms: 2,
      bathrooms: 1,
      area: 85,
      areaUnit: 'm²',
      furnished: true,
      parking: true,
      createdAt: '2023-10-14T12:00:00Z',
      agent: {
        id: 'agent-123',
        name: 'CityProperties',
        verified: true,
        rating: 4.8
      },
      thumbnail: 'https://example.com/apartment.jpg'
    },
    createdAt: '2023-10-14T12:00:00Z'
  },
  {
    id: generateId(),
    entityType: 'property',
    property: {
      id: 'property-002',
      title: 'Suburban Family House',
      description: 'Beautiful family house with garden',
      price: 450000,
      currency: 'EUR',
      category: 'houses',
      tags: ['house', 'family', 'garden'],
      location: 'Munich',
      type: 'house',
      bedrooms: 4,
      bathrooms: 2,
      area: 180,
      areaUnit: 'm²',
      furnished: false,
      parking: true,
      createdAt: '2023-10-13T10:30:00Z',
      agent: {
        id: 'agent-456',
        name: 'FamilyHomes',
        verified: true,
        rating: 4.6
      },
      thumbnail: 'https://example.com/house.jpg'
    },
    createdAt: '2023-10-13T10:30:00Z'
  }
];

// Mock service items
const serviceItems = [
  {
    id: generateId(),
    entityType: 'service',
    service: {
      id: 'service-001',
      title: 'Professional Photography Services',
      description: 'High-quality photography for events and portraits',
      price: 150,
      currency: 'EUR',
      category: 'photography',
      tags: ['photography', 'events', 'portraits'],
      location: 'Berlin',
      availability: 'weekends',
      onlineService: false,
      createdAt: '2023-10-16T11:00:00Z',
      provider: {
        id: 'provider-123',
        name: 'CaptureMoments',
        verified: true,
        rating: 4.9
      },
      thumbnail: 'https://example.com/photography.jpg'
    },
    createdAt: '2023-10-16T11:00:00Z'
  },
  {
    id: generateId(),
    entityType: 'service',
    service: {
      id: 'service-002',
      title: 'Web Development & Design',
      description: 'Custom website development and UI/UX design',
      price: 80,
      priceType: 'hourly',
      currency: 'EUR',
      category: 'web-development',
      tags: ['web-dev', 'design', 'programming'],
      location: 'Remote',
      availability: 'weekdays',
      onlineService: true,
      createdAt: '2023-10-15T09:30:00Z',
      provider: {
        id: 'provider-456',
        name: 'WebWizards',
        verified: true,
        rating: 4.8
      },
      thumbnail: 'https://example.com/webdev.jpg'
    },
    createdAt: '2023-10-15T09:30:00Z'
  }
];

// Mock job items
const jobItems = [
  {
    id: generateId(),
    entityType: 'job',
    job: {
      id: 'job-001',
      title: 'Senior Frontend Developer',
      description: 'Looking for an experienced React developer to join our team',
      salary: '70,000 - 90,000',
      currency: 'EUR',
      category: 'software-development',
      tags: ['react', 'javascript', 'frontend'],
      location: 'Berlin',
      type: 'full-time',
      remote: true,
      experienceLevel: 'senior',
      createdAt: '2023-10-16T08:45:00Z',
      company: {
        id: 'company-123',
        name: 'TechInnovators',
        verified: true,
        logo: 'https://example.com/company1.jpg'
      },
      thumbnail: 'https://example.com/frontend-job.jpg'
    },
    createdAt: '2023-10-16T08:45:00Z'
  },
  {
    id: generateId(),
    entityType: 'job',
    job: {
      id: 'job-002',
      title: 'Marketing Manager',
      description: 'Seeking a creative marketing professional',
      salary: '55,000 - 65,000',
      currency: 'EUR',
      category: 'marketing',
      tags: ['marketing', 'digital', 'social-media'],
      location: 'Munich',
      type: 'full-time',
      remote: false,
      experienceLevel: 'mid',
      createdAt: '2023-10-15T14:30:00Z',
      company: {
        id: 'company-456',
        name: 'BrandBuilders',
        verified: false,
        logo: 'https://example.com/company2.jpg'
      },
      thumbnail: 'https://example.com/marketing-job.jpg'
    },
    createdAt: '2023-10-15T14:30:00Z'
  }
];

// Combine all feed items
const allFeedItems = [
  ...productItems,
  ...postItems,
  ...dealItems,
  ...vehicleItems,
  ...propertyItems,
  ...serviceItems,
  ...jobItems
];

// Mock unified feed success response
export const mockUnifiedFeedResponse = {
  items: allFeedItems,
  hasMore: true,
  total: allFeedItems.length,
  page: 1,
  page_size: 30
};

// Mock unified feed with pagination (page 2)
export const mockUnifiedFeedPage2Response = {
  items: allFeedItems.slice(5, 10),
  hasMore: false,
  total: allFeedItems.length,
  page: 2,
  page_size: 5
};

// Mock responses for specific entity types
export const getMockFeedResponseByEntityType = (entityType) => {
  const items = allFeedItems.filter(item => item.entityType === entityType);
  return {
    items,
    hasMore: false,
    total: items.length,
    page: 1,
    page_size: 30
  };
};

// Mock response for location filter
export const getMockFeedResponseByLocation = (location) => {
  const items = allFeedItems.filter(item => {
    const locationField = item[item.entityType]?.location;
    return locationField && locationField.toLowerCase().includes(location.toLowerCase());
  });
  
  return {
    items,
    hasMore: false,
    total: items.length,
    page: 1,
    page_size: 30
  };
};

// Mock response for tag filter
export const getMockFeedResponseByTags = (tags) => {
  const tagArray = typeof tags === 'string' ? tags.split(',') : tags;
  
  const items = allFeedItems.filter(item => {
    const itemTags = item[item.entityType]?.tags || [];
    return tagArray.some(tag => 
      itemTags.some(itemTag => itemTag.toLowerCase().includes(tag.toLowerCase()))
    );
  });
  
  return {
    items,
    hasMore: false,
    total: items.length,
    page: 1,
    page_size: 30
  };
};

// Mock error response
export const mockErrorResponse = {
  error: true,
  message: 'An error occurred while fetching feed data',
  status: 500
};

export default {
  allFeedItems,
  mockUnifiedFeedResponse,
  mockUnifiedFeedPage2Response,
  getMockFeedResponseByEntityType,
  getMockFeedResponseByLocation,
  getMockFeedResponseByTags,
  mockErrorResponse
}; 