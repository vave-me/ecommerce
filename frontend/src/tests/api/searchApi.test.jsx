import axios from 'axios';
import MockAdapter from 'axios-mock-adapter';
import * as searchApi from '../../api/searchApi';

// Access the internal validateEntityPayload function for testing
// We need to use this approach because it's not exported directly
const validateEntityPayload = searchApi.default 
  ? searchApi.default.validateEntityPayload 
  : Object.values(searchApi).find(fn => 
      fn.toString().includes('validateEntityPayload')
    );

// Create a mock for axios
const mockAxios = new MockAdapter(axios);

describe('searchApi', () => {
  beforeEach(() => {
    // Reset mocks before each test
    mockAxios.reset();
  });

  afterAll(() => {
    // Clean up
    mockAxios.restore();
  });

  describe('validateEntityPayload', () => {
    test('filters out undocumented fields for products', () => {
      const input = {
        name: 'Test Product',
        category: 'electronics',
        minPrice: '100',
        maxPrice: '500',
        brand: 'TestBrand',
        // Undocumented fields
        minShipping: '10',
        maxShipping: '50',
        dimensionFilter: true,
        // Standard fields
        page: 1,
        sortBy: 'price',
        sortOrder: 'asc'
      };

      const result = validateEntityPayload(input, 'product');

      // Check that documented fields are kept
      expect(result).toHaveProperty('name', 'Test Product');
      expect(result).toHaveProperty('category', 'electronics');
      expect(result).toHaveProperty('minPrice', 100); // Should be converted to number
      expect(result).toHaveProperty('maxPrice', 500); // Should be converted to number
      expect(result).toHaveProperty('brand', 'TestBrand');
      expect(result).toHaveProperty('page', 1);
      expect(result).toHaveProperty('sortBy', 'price');
      expect(result).toHaveProperty('sortOrder', 'asc');

      // Check that undocumented fields are removed
      expect(result).not.toHaveProperty('minShipping');
      expect(result).not.toHaveProperty('maxShipping');
      expect(result).not.toHaveProperty('dimensionFilter');
    });

    test('removes listingType and contentType for deals', () => {
      const input = {
        name: 'Test Deal',
        category: 'electronics',
        minPrice: '100',
        maxPrice: '500',
        listingType: 'deals',
        contentType: 'deals',
        page: 1
      };

      const result = validateEntityPayload(input, 'deal');

      // Check that documented fields are kept
      expect(result).toHaveProperty('name', 'Test Deal');
      expect(result).toHaveProperty('category', 'electronics');
      expect(result).toHaveProperty('minPrice', 100);
      expect(result).toHaveProperty('maxPrice', 500);
      expect(result).toHaveProperty('page', 1);

      // Check that listingType and contentType are removed
      expect(result).not.toHaveProperty('listingType');
      expect(result).not.toHaveProperty('contentType');
    });

    test('converts string boolean values to actual booleans', () => {
      const input = {
        negotiable: 'true',
        hasVariants: 'false',
        manageStock: true,
        page: 1
      };

      const result = validateEntityPayload(input, 'product');

      // Check boolean conversions
      expect(result).toHaveProperty('negotiable', true);
      expect(result).toHaveProperty('hasVariants', false);
      expect(result).toHaveProperty('manageStock', true);
      expect(result).toHaveProperty('page', 1);
    });

    test('removes empty values', () => {
      const input = {
        name: 'Test',
        category: '',
        minPrice: null,
        maxPrice: undefined,
        tags: []
      };

      const result = validateEntityPayload(input, 'product');

      // Check empty value removal
      expect(result).toHaveProperty('name', 'Test');
      expect(result).not.toHaveProperty('category');
      expect(result).not.toHaveProperty('minPrice');
      expect(result).not.toHaveProperty('maxPrice');
      // Empty arrays should be kept
      expect(result).toHaveProperty('tags');
      expect(result.tags).toEqual([]);
    });
    
    test('handles numeric fields correctly', () => {
      const input = {
        minStock: '10',    // valid string number -> should convert to integer
        maxStock: '',      // empty string -> should be removed
        minPrice: 25.5,    // float number -> should remain as is
        maxPrice: '50.75'  // string float -> should convert to number
      };
      
      const result = validateEntityPayload(input, 'product');
      
      expect(result).toHaveProperty('minStock', 10);
      expect(result).not.toHaveProperty('maxStock');
      expect(result).toHaveProperty('minPrice', 25.5);
      expect(result).toHaveProperty('maxPrice', 50.75);
    });
    
    test('handles tags properly', () => {
      const input = {
        // Test different tag formats
        tags: '',                      // Empty string -> should become empty array
        name: 'Test Product'
      };
      
      const result = validateEntityPayload(input, 'product');
      
      // Tags should be preserved as an array
      expect(result).toHaveProperty('tags');
      expect(Array.isArray(result.tags)).toBe(true);
      expect(result.tags).toHaveLength(0);
    });
    
    test('converts string tag to array', () => {
      const input = {
        tags: 'electronics,gaming,computers',  // Comma-separated string -> should become array
        name: 'Test Product'
      };
      
      const result = validateEntityPayload(input, 'product');
      
      expect(result).toHaveProperty('tags');
      expect(Array.isArray(result.tags)).toBe(true);
      expect(result.tags).toEqual(['electronics', 'gaming', 'computers']);
    });
    
    test('removes minStock and maxStock when zero or empty', () => {
      const input = {
        name: 'Test Product',
        minStock: 0,              // Zero -> should be removed
        maxStock: '0',            // String zero -> should be removed
        minPrice: 10,             // Valid -> should be kept
        hasStock: true            // Valid -> should be kept
      };
      
      const result = validateEntityPayload(input, 'product');
      
      // Stock fields should be removed
      expect(result).not.toHaveProperty('minStock');
      expect(result).not.toHaveProperty('maxStock');
      // Other fields should be kept
      expect(result).toHaveProperty('name', 'Test Product');
      expect(result).toHaveProperty('minPrice', 10);
      expect(result).toHaveProperty('hasStock', true);
    });
    
    test('only includes user-selected fields without adding defaults', () => {
      // A minimal input with just one field
      const input = {
        hasVariants: false
      };
      
      const result = validateEntityPayload(input, 'product');
      
      // Should only contain the one field that was provided
      expect(Object.keys(result).length).toBe(1);
      expect(result).toHaveProperty('hasVariants', false);
      
      // Should NOT add default values for these fields
      expect(result).not.toHaveProperty('page');
      expect(result).not.toHaveProperty('pageSize');
      expect(result).not.toHaveProperty('limit');
      expect(result).not.toHaveProperty('offset');
      expect(result).not.toHaveProperty('sortOrder');
    });
    
    test('handles empty input correctly', () => {
      // Empty input
      const input = {};
      
      const result = validateEntityPayload(input, 'product');
      
      // Result should be an empty object
      expect(Object.keys(result).length).toBe(0);
    });
    
    test('only allows sortOrder to be included in payload', () => {
      // Input with multiple fields
      const input = {
        hasVariants: false,
        limit: 20,
        manageStock: false,
        negotiable: false,
        offset: 0,
        page: 1,
        pageSize: 20,
        sortOrder: "asc"
      };
      
      const result = validateEntityPayload(input, 'product');
      
      // Should only contain sortOrder and nothing else
      expect(Object.keys(result).length).toBe(1);
      expect(result).toHaveProperty('sortOrder', 'asc');
      
      // Verify none of the other fields were included
      expect(result).not.toHaveProperty('hasVariants');
      expect(result).not.toHaveProperty('limit');
      expect(result).not.toHaveProperty('manageStock');
      expect(result).not.toHaveProperty('negotiable');
      expect(result).not.toHaveProperty('offset');
      expect(result).not.toHaveProperty('page');
      expect(result).not.toHaveProperty('pageSize');
    });
  });

  // Test for actual API methods if needed
  // These would use mockAxios to mock the API responses
}); 