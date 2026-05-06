import {
  flattenTree,
  filterCategories,
  processCategoryData,
  buildCategoryTree,
  findCategoryById,
  getCategoryPath
} from '@/utils/categoryUtils.jsx';

describe('Category Utilities', () => {
  
  // Sample data for tests
  const sampleCategories = [
    {
      id: '1',
      name: 'Electronics',
      subcategories: [
        {
          id: '11',
          name: 'Computers',
          subcategories: [
            { id: '111', name: 'Laptops', subcategories: [] },
            { id: '112', name: 'Desktops', subcategories: [] }
          ]
        },
        {
          id: '12',
          name: 'Phones',
          subcategories: [
            { id: '121', name: 'Smartphones', subcategories: [] },
            { id: '122', name: 'Accessories', subcategories: [] }
          ]
        }
      ]
    },
    {
      id: '2',
      name: 'Clothing',
      subcategories: [
        { id: '21', name: 'Men', subcategories: [] },
        { id: '22', name: 'Women', subcategories: [] }
      ]
    }
  ];

  // Sample flat categories with parent references
  const sampleFlatCategories = [
    { id: '1', name: 'Electronics' },
    { id: '11', name: 'Computers', parentId: '1' },
    { id: '111', name: 'Laptops', parentId: '11' },
    { id: '112', name: 'Desktops', parentId: '11' },
    { id: '12', name: 'Phones', parentId: '1' },
    { id: '121', name: 'Smartphones', parentId: '12' },
    { id: '122', name: 'Accessories', parentId: '12' },
    { id: '2', name: 'Clothing' },
    { id: '21', name: 'Men', parentId: '2' },
    { id: '22', name: 'Women', parentId: '2' }
  ];

  describe('flattenTree', () => {
    test('should flatten a nested tree to a single level array', () => {
      const result = flattenTree(sampleCategories);
      
      expect(result).toHaveLength(2); // Only root categories when no expanded map
      expect(result[0].id).toBe('1'); // Electronics
      expect(result[1].id).toBe('2'); // Clothing
    });

    test('should include subcategories when they are expanded', () => {
      const expandedMap = { '1': true };
      const result = flattenTree(sampleCategories, expandedMap);
      
      expect(result).toHaveLength(4); // Root categories + direct subcategories of expanded node
      expect(result[0].id).toBe('1'); // Electronics
      expect(result[1].id).toBe('11'); // Computers
      expect(result[2].id).toBe('12'); // Phones
      expect(result[3].id).toBe('2'); // Clothing
    });

    test('should include nested subcategories when multiple levels are expanded', () => {
      const expandedMap = { '1': true, '11': true };
      const result = flattenTree(sampleCategories, expandedMap);
      
      expect(result).toHaveLength(6); // Root + subcats of Electronics + subcats of Computers
      expect(result[0].id).toBe('1'); // Electronics
      expect(result[1].id).toBe('11'); // Computers
      expect(result[2].id).toBe('111'); // Laptops
      expect(result[3].id).toBe('112'); // Desktops
      expect(result[4].id).toBe('12'); // Phones
      expect(result[5].id).toBe('2'); // Clothing
    });

    test('should handle empty or null input', () => {
      expect(flattenTree(null)).toEqual([]);
      expect(flattenTree([])).toEqual([]);
      expect(flattenTree(undefined)).toEqual([]);
    });
  });

  describe('filterCategories', () => {
    test('should filter categories based on name', () => {
      const result = filterCategories(sampleCategories, 'phone');
      
      expect(result).toHaveLength(1); // Only Electronics should match
      expect(result[0].id).toBe('1'); // Electronics
      expect(result[0].subcategories).toHaveLength(1); // Only Phones subcategory
      expect(result[0].subcategories[0].id).toBe('12'); // Phones
    });

    test('should be case-insensitive', () => {
      const result = filterCategories(sampleCategories, 'LAPTOP');
      
      expect(result).toHaveLength(1); // Only Electronics should match
      expect(result[0].subcategories[0].subcategories).toHaveLength(1); // Only Laptops
      expect(result[0].subcategories[0].subcategories[0].id).toBe('111'); // Laptops
    });

    test('should return all categories if query is empty', () => {
      expect(filterCategories(sampleCategories, '')).toEqual(sampleCategories);
      expect(filterCategories(sampleCategories, '  ')).toEqual(sampleCategories);
    });

    test('should handle empty or null input', () => {
      expect(filterCategories(null, 'test')).toEqual([]);
      expect(filterCategories([], 'test')).toEqual([]);
      expect(filterCategories(undefined, 'test')).toEqual([]);
    });
  });

  describe('processCategoryData', () => {
    test('should normalize category data', () => {
      const inputCategories = [
        { id: '1', name: 'Electronics' },
        { id: '2', description: 'Fashion items' }, // No name
        { id: '3' } // No name or description
      ];
      
      const result = processCategoryData(inputCategories);
      
      expect(result).toHaveLength(3);
      expect(result[0].name).toBe('Electronics');
      expect(result[1].name).toBe('Fashion items'); // Should use description as name
      expect(result[2].name).toBe('Unnamed Category'); // Should use default name
      
      // Should initialize empty subcategories arrays
      expect(result[0].subcategories).toEqual([]);
      expect(result[1].subcategories).toEqual([]);
      expect(result[2].subcategories).toEqual([]);
    });

    test('should handle empty or null input', () => {
      expect(processCategoryData(null)).toEqual([]);
      expect(processCategoryData([])).toEqual([]);
      expect(processCategoryData(undefined)).toEqual([]);
    });

    test('should filter out null/undefined categories', () => {
      const input = [{ id: '1' }, null, undefined, { id: '2' }];
      const result = processCategoryData(input);
      
      expect(result).toHaveLength(2);
      expect(result[0].id).toBe('1');
      expect(result[1].id).toBe('2');
    });
  });

  describe('buildCategoryTree', () => {
    test('should build a nested tree from flat categories', () => {
      const result = buildCategoryTree(sampleFlatCategories);
      
      // Should have 2 root categories
      expect(result).toHaveLength(2);
      
      // Check first root category
      expect(result[0].id).toBe('1');
      expect(result[0].name).toBe('Electronics');
      expect(result[0].subcategories).toHaveLength(2);
      
      // Check subcategories of first root
      expect(result[0].subcategories[0].id).toBe('11');
      expect(result[0].subcategories[0].name).toBe('Computers');
      expect(result[0].subcategories[0].subcategories).toHaveLength(2);
      
      // Check second level subcategories
      expect(result[0].subcategories[0].subcategories[0].id).toBe('111');
      expect(result[0].subcategories[0].subcategories[0].name).toBe('Laptops');
    });

    test('should handle empty or null input', () => {
      expect(buildCategoryTree(null)).toEqual([]);
      expect(buildCategoryTree([])).toEqual([]);
      expect(buildCategoryTree(undefined)).toEqual([]);
    });

    test('should handle invalid parent references', () => {
      const invalidParentCategories = [
        { id: '1', name: 'Electronics' },
        { id: '2', name: 'Computers', parentId: 'non-existent' }
      ];
      
      const result = buildCategoryTree(invalidParentCategories);
      
      // Both should be treated as root categories
      expect(result).toHaveLength(2);
      expect(result[0].id).toBe('1');
      expect(result[1].id).toBe('2');
    });
  });

  describe('findCategoryById', () => {
    test('should find a category by ID in a nested tree', () => {
      const result = findCategoryById(sampleCategories, '121');
      
      expect(result).not.toBeNull();
      expect(result.id).toBe('121');
      expect(result.name).toBe('Smartphones');
    });

    test('should return null if category not found', () => {
      const result = findCategoryById(sampleCategories, 'non-existent');
      expect(result).toBeNull();
    });

    test('should handle empty or null input', () => {
      expect(findCategoryById(null, '1')).toBeNull();
      expect(findCategoryById([], '1')).toBeNull();
      expect(findCategoryById(undefined, '1')).toBeNull();
      expect(findCategoryById(sampleCategories, null)).toBeNull();
      expect(findCategoryById(sampleCategories, undefined)).toBeNull();
    });
  });

  describe('getCategoryPath', () => {
    test('should return path from root to target category', () => {
      const result = getCategoryPath(sampleCategories, '111');
      
      expect(result).toHaveLength(3);
      expect(result[0].id).toBe('1'); // Electronics
      expect(result[1].id).toBe('11'); // Computers
      expect(result[2].id).toBe('111'); // Laptops
    });

    test('should return empty array if category not found', () => {
      const result = getCategoryPath(sampleCategories, 'non-existent');
      expect(result).toEqual([]);
    });

    test('should handle empty or null input', () => {
      expect(getCategoryPath(null, '1')).toEqual([]);
      expect(getCategoryPath([], '1')).toEqual([]);
      expect(getCategoryPath(undefined, '1')).toEqual([]);
      expect(getCategoryPath(sampleCategories, null)).toEqual([]);
      expect(getCategoryPath(sampleCategories, undefined)).toEqual([]);
    });
  });
}); 