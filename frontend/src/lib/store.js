import { configureStore } from '@reduxjs/toolkit';
import debugFilters from '../utils/debugFilters';

// Import all slices
import listingFiltersSlice from '../redux/slices/listingFiltersSlice';
import modalsSlice from '../redux/slices/modalsSlice';
// REMOVED: filterSlice - using listingFiltersSlice instead to avoid confusion
import postFilterSlice from '../redux/slices/postFilterSlice';
import listingSlice from '../redux/slices/listingSlice';
import adsSlice from '../redux/slices/adsSlice';
import appModeSlice from '../redux/slices/appModeSlice';
import aiUISlice from '../redux/slices/aiUISlice';
import uiPreferencesSlice from '../redux/slices/uiPreferencesSlice';

/**
 * Performance monitoring middleware
 * Measures reducer execution time in development
 */
const measureReducers = store => next => action => {
  if (process.env.NODE_ENV === 'development') {
    const start = performance.now();
    const result = next(action);
    const end = performance.now();
    const diff = end - start;
    if (diff > 10) {
      // Slow reducer execution logged for debugging
    }
    return result;
  }
  return next(action);
};

/**
 * Filter debugging middleware
 * Logs all filter-related actions and state changes
 */
const filterDebugMiddleware = store => next => action => {
  // Check if this is a filter-related action
  if (action.type.includes('Filters/') || action.type.includes('filter/')) {
    const stateBefore = store.getState();
    debugFilters.logFilterFlow('REDUX_DISPATCH', {
      action: action.type,
      payload: action.payload,
      filtersBefore: stateBefore.listingFilters,
      timestamp: new Date().toISOString()
    });
    
    const result = next(action);
    
    const stateAfter = store.getState();
    debugFilters.compareFilters(stateBefore.listingFilters, stateAfter.listingFilters);
    
    return result;
  }
  
  return next(action);
};

// Unified root reducer with all slices
const rootReducer = {
  listingFilters: listingFiltersSlice,
  modals: modalsSlice,
  // REMOVED: filter - using listingFilters instead to avoid confusion
  postFilter: postFilterSlice,
  listing: listingSlice,
  ads: adsSlice,
  appMode: appModeSlice,
  aiUI: aiUISlice,
  uiPreferences: uiPreferencesSlice,
};

/**
 * Create Redux store with unified configuration
 * Combines the best features from both store configurations
 */
const createStore = () => {
  return configureStore({
    reducer: rootReducer,
    middleware: (getDefaultMiddleware) => {
      const middleware = getDefaultMiddleware({
        // Serialization checks only in development
        serializableCheck: process.env.NODE_ENV === 'development' ? {
          // Ignore these action types
          ignoredActions: ['persist/PERSIST', 'persist/REHYDRATE'],
        } : false,
        // Immutability checking only in development
        immutableCheck: process.env.NODE_ENV === 'development',
        // Enhanced thunk configuration
        thunk: {
          extraArgument: {
            // Services can be injected here if needed
            api: {},
          }
        }
      });
      
      // Add performance monitoring in development
      if (process.env.NODE_ENV === 'development') {
        middleware.push(measureReducers);
        middleware.push(filterDebugMiddleware);
      }
      
      return middleware;
    },
    // Enable Redux DevTools only in development
    devTools: process.env.NODE_ENV !== 'production',
  });
};

// Create store instance - SINGLETON pattern
let store;

const getStore = () => {
  if (!store) {
    store = createStore();
    
    // Enable hot module replacement for reducers in development
    if (process.env.NODE_ENV === 'development' && module.hot) {
      // Accept changes to individual slice files
      module.hot.accept([
        '../redux/slices/listingFiltersSlice',
        '../redux/slices/modalsSlice',
        '../redux/slices/postFilterSlice',
        '../redux/slices/listingSlice',
        '../redux/slices/adsSlice',
        '../redux/slices/appModeSlice',
        '../redux/slices/aiUISlice',
        '../redux/slices/uiPreferencesSlice'
      ], () => {
        // Re-import the slices and replace the reducer
        const newRootReducer = {
          listingFilters: require('../redux/slices/listingFiltersSlice').default,
          modals: require('../redux/slices/modalsSlice').default,
          postFilter: require('../redux/slices/postFilterSlice').default,
          listing: require('../redux/slices/listingSlice').default,
          ads: require('../redux/slices/adsSlice').default,
          appMode: require('../redux/slices/appModeSlice').default,
          aiUI: require('../redux/slices/aiUISlice').default,
          uiPreferences: require('../redux/slices/uiPreferencesSlice').default,
        };
        store.replaceReducer(newRootReducer);
      });
    }
  }
  return store;
};

// Export the store instance
store = getStore();

export { store, rootReducer };
export default store; 