import { combineReducers } from 'redux';
import rootReducer from '../rootReducer';
import listingFiltersReducer from '../../slices/listingFiltersSlice';
import modalsReducer from '../../slices/modalsSlice';

// Mock dependencies
jest.mock('redux', () => ({
  combineReducers: jest.fn(reducers => ({ mockCombinedReducer: true, reducers }))
}));

jest.mock('../../slices/listingFiltersSlice', () => 'mockListingFiltersReducer');
jest.mock('../../slices/modalsSlice', () => 'mockModalsReducer');

describe('rootReducer', () => {
  test('should combine all reducers correctly', () => {
    expect(combineReducers).toHaveBeenCalledWith({
      listingFilters: listingFiltersReducer,
      modals: modalsReducer
    });
  });

  test('should export the combined reducer', () => {
    expect(rootReducer).toEqual({
      mockCombinedReducer: true,
      reducers: {
        listingFilters: listingFiltersReducer,
        modals: modalsReducer
      }
    });
  });
}); 