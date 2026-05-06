// File: src/redux/slices/filterSlice.js
import {createSlice} from '@reduxjs/toolkit';
/**
 * Consolidated filter fields for all listing types
 */
const createInitialFiltersState = (listingType = 'marketplace') => ({
    // Identifier for which type of listing we're filtering
    listingType,
    // Content type for unified catalog filtering (following unified feed pattern)
    contentType: 'all',
    // Entity types for unified feed (which types to include)
    entityTypes: ['product', 'post', 'service'],
    // Common display fields
    displayMode: 'grid',
    // Common search fields
    searchText: '',
    title: '', // Alias for searchText for post filters (for backward compatibility)
    content: '', // Content search for posts
    // Category and location
    category: '',
    categoryID: '',
    categorySlug: '',
    location: '',
    lat: null,
    lng: null,
    radius: null,
    // Common filter fields
    status: '',
    tags: [],
    // Pagination
    page: 1,
    pageSize: 20,
    offset: 0,
    limit: 20,
    // Sorting
    sortBy: '',
    sortOrder: 'asc',
    // Price range
    minPrice: '',
    maxPrice: '',
    // Boolean filters
    excludeNoMedia: false,
    // Product-specific fields
    ...(listingType === 'marketplace' || listingType === 'products' ? {
        condition: '',
        brand: '',
        model: '',
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
        dimensionFilter: false,
        minWeight: '',
        maxWeight: '',
        minHeight: '',
        maxHeight: '',
        minWidth: '',
        maxWidth: '',
        minDepth: '',
        maxDepth: '',
        // Add common filter fields that were missing
        freeShipping: false,
        verifiedSeller: false,
    } : {}),
    // Video-specific fields
    ...(listingType === 'videos' ? {
        duration: '',
        quality: '',
        resolution: '',
        hasSound: false,
        hasCaptions: false,
    } : {}),
    // Tweet-specific fields
    ...(listingType === 'tweets' ? {
        hasMedia: false,
        hasLinks: false,
        hasMentions: false,
        isVerified: false,
    } : {}),
    // Job-specific fields
    ...(listingType === 'jobs' ? {
        jobType: '',
        salaryRange: '',
        experienceLevel: '',
        remote: false,
        employmentType: '',
    } : {}),
    // Service-specific fields
    ...(listingType === 'services' ? {
        serviceType: '',
        availability: '',
        isOnline: false,
        hasPortfolio: false,
    } : {}),
    // Deal-specific fields
    ...(listingType === 'deals' ? {
        dealType: '',
        minDiscount: '',
        maxDiscount: '',
        freeShipping: false,
        verifiedSeller: false,
        expiry: '',
    } : {}),
    // Short-specific fields
    ...(listingType === 'shorts' ? {
        duration: '',
        quality: '',
        hasSound: false,
        hasCaptions: false,
    } : {}),
    // Vehicle-specific fields
    ...(listingType === 'vehicles' ? {
        vehicleType: '',
        make: '',
        model: '',
        yearMin: '',
        yearMax: '',
        mileageMin: '',
        mileageMax: '',
        fuelType: '',
        transmission: '',
    } : {}),
    // Property-specific fields
    ...(listingType === 'properties' ? {
        propertyType: '',
        bedroomsMin: '',
        bedroomsMax: '',
        bathroomsMin: '',
        bathroomsMax: '',
        areaMin: '',
        areaMax: '',
        hasParking: false,
        furnished: false,
    } : {}),
    // News-specific fields
    ...(listingType === 'news' ? {
        newsCategory: '',
        timePeriod: '',
        sourceType: '',
        hasImages: false,
        hasVideos: false,
    } : {}),
});
/**
 * Create and configure a filter slice for a specific type
 * 
 * @param {string} name - Name of the slice
 * @param {string} listingType - Type of listing (products, posts, etc)
 * @returns {Object} Configured Redux slice
 */
export const createFilterSlice = (name, listingType = 'products') => {
    const initialState = createInitialFiltersState(listingType);
    return createSlice({
        name,
        initialState,
        reducers: {
            setFilters(state, action) {
                return {...state, ...action.payload};
            },
            updateFilter(state, action) {
                const {key, value} = action.payload;
                state[key] = value;
            },
            resetFilters(state) {
                // Reset but maintain the listing type and entity types
                const type = state.listingType;
                const entityTypes = state.entityTypes;
                const newState = createInitialFiltersState(type);
                // Preserve entity types if they were customized
                if (entityTypes && entityTypes.length > 0) {
                    newState.entityTypes = entityTypes;
                }
                return newState;
            },
            clearFilters() {
                return {};
            },
            // Special action for switching listing types
            setListingType(state, action) {
                state.listingType = action.payload;
            },
            // Unified action for setting content type (entity type)
            setContentType(state, action) {
                state.contentType = action.payload;
            }
        },
    });
};
// Product filters slice
const productFiltersSlice = createFilterSlice('productFilters', 'products');
export const {
    setFilters: setProductFilters,
    updateFilter: updateProductFilter,
    resetFilters: resetProductFilters,
    clearFilters: clearProductFilters,
} = productFiltersSlice.actions;
// Create a unified listing filters slice that can handle all types
const listingFiltersSlice = createFilterSlice('listingFilters', 'products');
export const {
    setFilters: setListingFilters,
    updateFilter: updateListingFilter,
    resetFilters: resetListingFilters,
    clearFilters: clearListingFilters,
    setListingType,
    setContentType: setListingContentType,
} = listingFiltersSlice.actions;
// Post filters slice - maintained for backward compatibility
const postFiltersSlice = createFilterSlice('postFilters', 'posts');
export const {
    setFilters: setPostFilters,
    updateFilter: updatePostFilter,
    resetFilters: resetPostFilters,
    clearFilters: clearPostFilters,
} = postFiltersSlice.actions;
// Default export - DEPRECATED: Use named exports instead
export default productFiltersSlice.reducer;
// Named exports for the slices
export const productFiltersReducer = productFiltersSlice.reducer;
export const listingFiltersReducer = listingFiltersSlice.reducer;
export const postFiltersReducer = postFiltersSlice.reducer;
