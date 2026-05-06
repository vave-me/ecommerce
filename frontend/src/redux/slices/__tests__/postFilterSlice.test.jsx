import postFiltersReducer, {
  setPostFilters,
  updatePostFilter,
  resetPostFilters,
  clearPostFilters
} from '../postFilterSlice';

describe('postFilterSlice', () => {
  const initialState = {
    displayMode: 'mobile',
    title: '',
    content: '',
    status: '',
    tags: [],
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
    expect(postFiltersReducer(undefined, { type: undefined })).toEqual(initialState);
  });

  test('should handle setPostFilters', () => {
    const newFilters = {
      title: 'Test Post',
      content: 'Test content',
      status: 'published'
    };
    
    expect(
      postFiltersReducer(initialState, setPostFilters(newFilters))
    ).toEqual({
      ...initialState,
      ...newFilters
    });
  });

  test('should handle updatePostFilter', () => {
    const payload = {
      key: 'title',
      value: 'Updated Title'
    };
    
    expect(
      postFiltersReducer(initialState, updatePostFilter(payload))
    ).toEqual({
      ...initialState,
      title: 'Updated Title'
    });
  });

  test('should handle resetPostFilters', () => {
    const modifiedState = {
      ...initialState,
      title: 'Test Post',
      content: 'Test content',
      status: 'published'
    };
    
    expect(
      postFiltersReducer(modifiedState, resetPostFilters())
    ).toEqual(initialState);
  });

  test('should handle clearPostFilters', () => {
    const modifiedState = {
      ...initialState,
      title: 'Test Post',
      content: 'Test content',
      status: 'published'
    };
    
    expect(
      postFiltersReducer(modifiedState, clearPostFilters())
    ).toEqual({});
  });

  test('should handle multiple filter updates', () => {
    let state = initialState;
    
    state = postFiltersReducer(state, updatePostFilter({ key: 'title', value: 'Test Post' }));
    state = postFiltersReducer(state, updatePostFilter({ key: 'content', value: 'Test Content' }));
    state = postFiltersReducer(state, updatePostFilter({ key: 'status', value: 'draft' }));
    
    expect(state).toEqual({
      ...initialState,
      title: 'Test Post',
      content: 'Test Content',
      status: 'draft'
    });
  });

  test('should handle pagination updates', () => {
    let state = initialState;
    
    state = postFiltersReducer(state, updatePostFilter({ key: 'page', value: 2 }));
    state = postFiltersReducer(state, updatePostFilter({ key: 'pageSize', value: 50 }));
    
    expect(state).toEqual({
      ...initialState,
      page: 2,
      pageSize: 50
    });
  });

  test('should handle location-based filters', () => {
    const locationFilters = {
      lat: 40.7128,
      lng: -74.0060,
      radius: 10
    };
    
    expect(
      postFiltersReducer(initialState, setPostFilters(locationFilters))
    ).toEqual({
      ...initialState,
      ...locationFilters
    });
  });
}); 