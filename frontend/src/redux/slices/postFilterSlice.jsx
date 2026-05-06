// File: src/redux/slices/postFilterSlice.js
// Re-export from the consolidated implementation for backward compatibility
import {
  postFiltersReducer,
  setPostFilters,
  updatePostFilter,
  resetPostFilters,
  clearPostFilters
} from './filterSlice';
export {
  setPostFilters,
  updatePostFilter,
  resetPostFilters,
  clearPostFilters
};
export default postFiltersReducer;
