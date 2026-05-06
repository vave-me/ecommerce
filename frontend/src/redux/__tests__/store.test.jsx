import { configureStore } from '@reduxjs/toolkit';
import rootReducer from '../reducers/rootReducer';
import store from '../store';

// Mock @reduxjs/toolkit and the root reducer
jest.mock('@reduxjs/toolkit', () => ({
  configureStore: jest.fn().mockReturnValue({ mockStore: true })
}));

jest.mock('../reducers/rootReducer', () => ({ mockReducer: true }));

describe('Redux store', () => {
  test('should create store with root reducer', () => {
    expect(configureStore).toHaveBeenCalledWith({
      reducer: rootReducer
    });
  });

  test('should export the configured store', () => {
    expect(store).toEqual({ mockStore: true });
  });
}); 