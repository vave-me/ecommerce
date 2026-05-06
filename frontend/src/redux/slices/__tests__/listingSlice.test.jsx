import listingReducer, { setListingType } from '../listingSlice';

describe('listingSlice', () => {
  const initialState = {
    listingType: 'products',
  };

  test('should return the initial state', () => {
    expect(listingReducer(undefined, { type: undefined })).toEqual(initialState);
  });

  test('should handle setListingType with products', () => {
    expect(
      listingReducer(initialState, setListingType('products'))
    ).toEqual({
      listingType: 'products',
    });
  });

  test('should handle setListingType with posts', () => {
    expect(
      listingReducer(initialState, setListingType('posts'))
    ).toEqual({
      listingType: 'posts',
    });
  });

  test('should handle setListingType with custom value', () => {
    expect(
      listingReducer(initialState, setListingType('cars'))
    ).toEqual({
      listingType: 'cars',
    });
  });

  test('should handle changing from one listing type to another', () => {
    // Start with posts
    const postsState = listingReducer(initialState, setListingType('posts'));
    expect(postsState).toEqual({
      listingType: 'posts',
    });
    
    // Then change to cars
    const carsState = listingReducer(postsState, setListingType('cars'));
    expect(carsState).toEqual({
      listingType: 'cars',
    });
  });
}); 