import productFiltersReducer, {
  setFilters,
  updateFilter,
  resetFilters,
  clearFilters
} from '../filterSlice';

describe('filterSlice', () => {
  const initialState = {
    displayMode: 'mobile',
    name: '',
    category: '',
    minPrice: '',
    maxPrice: '',
    condition: '',
    brand: '',
    model: '',
    status: '',
    rating: '',
    negotiable: false,
    userType: '',
    middlemanService: '',
    manageStock: false,
    minStock: '',
    maxStock: '',
    sku: '',
    hasVariants: false,
    minShipping: '',
    maxShipping: '',
    onlyPickup: false,
    excludeNoMedia: false,
    dimensionFilter: false,
    minWeight: '',
    maxWeight: '',
    minHeight: '',
    maxHeight: '',
    minWidth: '',
    maxWidth: '',
    minDepth: '',
    maxDepth: '',
    tags: [],
    location: '',
    speaksLanguage: '',
    lat: null,
    lng: null,
    radius: null,
    page: 1,
    pageSize: 20,
    offset: 0,
    limit: 20,
    sortBy: '',
    sortOrder: 'asc',
  };

  test('should return the initial state', () => {
    expect(productFiltersReducer(undefined, { type: undefined })).toEqual(initialState);
  });

  test('should handle setFilters', () => {
    const newFilters = {
      category: 'electronics',
      minPrice: '100',
      maxPrice: '500'
    };
    
    expect(
      productFiltersReducer(initialState, setFilters(newFilters))
    ).toEqual({
      ...initialState,
      ...newFilters
    });
  });

  test('should handle updateFilter', () => {
    const payload = {
      key: 'category',
      value: 'clothing'
    };
    
    expect(
      productFiltersReducer(initialState, updateFilter(payload))
    ).toEqual({
      ...initialState,
      category: 'clothing'
    });
  });

  test('should handle resetFilters', () => {
    const modifiedState = {
      ...initialState,
      category: 'electronics',
      minPrice: '100',
      maxPrice: '500'
    };
    
    expect(
      productFiltersReducer(modifiedState, resetFilters())
    ).toEqual(initialState);
  });

  test('should handle clearFilters', () => {
    const modifiedState = {
      ...initialState,
      category: 'electronics',
      minPrice: '100',
      maxPrice: '500'
    };
    
    expect(
      productFiltersReducer(modifiedState, clearFilters())
    ).toEqual({});
  });
}); 