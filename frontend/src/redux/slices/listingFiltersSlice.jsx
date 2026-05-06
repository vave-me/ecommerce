// File: src/redux/slices/listingFiltersSlice.js
// Re-export from the consolidated implementation
import {
  listingFiltersReducer,
  setListingType,
  setListingFilters,
  updateListingFilter,
  resetListingFilters,
  clearListingFilters,
  setListingContentType
} from './filterSlice';

// Re-export with cleaner names for components
export {
  setListingType,
  setListingFilters as setFilters,
  updateListingFilter as updateFilter,
  resetListingFilters as resetFilters,
  clearListingFilters as clearFilters,
  setListingContentType
};

export default listingFiltersReducer;
