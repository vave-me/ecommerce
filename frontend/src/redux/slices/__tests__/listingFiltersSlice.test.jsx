import listingFiltersReducer, {
  setListingType,
  setFilters,
  updateFilter,
  resetFilters
} from '../listingFiltersSlice';

describe('listingFiltersSlice', () => {
  const initialState = {
    listingType: 'products',
    category: '',
    displayMode: 'mobile',
    searchText: '',
    sortBy: '',
    sortOrder: 'asc',
    location: '',
    minPrice: '',
    maxPrice: '',
    tags: [],
    page: 1,
    pageSize: 20,
    offset: 0,
    limit: 20,
    condition: '',
    brand: '',
    model: '',
    rating: '',
    negotiable: false,
    userType: '',
    manageStock: false,
    minStock: '',
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
    sku: '',
    title: '',
    content: '',
    status: '',
  };

  test('should return the initial state', () => {
    expect(listingFiltersReducer(undefined, { type: undefined })).toEqual(initialState);
  });

  test('should handle setListingType', () => {
    expect(
      listingFiltersReducer(initialState, setListingType('posts'))
    ).toEqual({
      ...initialState,
      listingType: 'posts'
    });
  });

  test('should handle setFilters', () => {
    const newFilters = {
      category: 'electronics',
      minPrice: '100',
      maxPrice: '500'
    };
    
    expect(
      listingFiltersReducer(initialState, setFilters(newFilters))
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
      listingFiltersReducer(initialState, updateFilter(payload))
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
      listingFiltersReducer(modifiedState, resetFilters())
    ).toEqual(initialState);
  });

  test('should handle multiple filter updates', () => {
    let state = initialState;
    
    state = listingFiltersReducer(state, updateFilter({ key: 'category', value: 'electronics' }));
    state = listingFiltersReducer(state, updateFilter({ key: 'minPrice', value: '200' }));
    state = listingFiltersReducer(state, updateFilter({ key: 'maxPrice', value: '800' }));
    
    expect(state).toEqual({
      ...initialState,
      category: 'electronics',
      minPrice: '200',
      maxPrice: '800'
    });
  });

  test('should handle changing listing type and then setting filters', () => {
    let state = initialState;
    
    state = listingFiltersReducer(state, setListingType('posts'));
    state = listingFiltersReducer(state, setFilters({ 
      title: 'Test Post',
      content: 'Post content'
    }));
    
    expect(state).toEqual({
      ...initialState,
      listingType: 'posts',
      title: 'Test Post',
      content: 'Post content'
    });
  });
}); 