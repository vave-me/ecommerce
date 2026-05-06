import { createSelector } from '@reduxjs/toolkit';
// Base selectors - extract slices from state
const selectModalsState = (state) => state.modals;
const selectListingFiltersState = (state) => state.listingFilters;
// ===== Modal Selectors =====
// Select all modal states at once (useful for GlobalModals component)
export const selectAllModalStates = createSelector(
  [selectModalsState],
  (modalsState) => modalsState
);
// Individual modal selectors
export const selectCommentsModalState = createSelector(
  [selectModalsState],
  (modals) => ({
    isOpen: modals.commentsFullModalOpen,
    itemId: modals.commentsFullModalItemId,
    itemType: modals.commentsFullItemType,
    categoryId: modals.commentsFullCategoryId
  })
);
export const selectMessageModalState = createSelector(
  [selectModalsState],
  (modals) => ({
    isOpen: modals.messageModalOpen,
    itemId: modals.messageModalItemId,
    recipientId: modals.messageRecipientId
  })
);
export const selectProductModalState = createSelector(
  [selectModalsState],
  (modals) => ({
    isOpen: modals.isProductModalOpen,
    selectedProduct: modals.selectedProduct
  })
);
export const selectAddModalStates = createSelector(
  [selectModalsState],
  (modals) => ({
    product: modals.addProductModalOpen,
    post: modals.addPostModalOpen,
    vehicle: modals.addVehicleModalOpen,
    deal: modals.addDealModalOpen,
    property: modals.addPropertyModalOpen,
    job: modals.addJobModalOpen,
    service: modals.addServiceModalOpen,
    video: modals.isVideoModalOpen
  })
);
export const selectOpenModalsCount = createSelector(
  [selectModalsState],
  (modals) => modals.openModalsCount
);
// ===== Listing Filters Selectors =====
// Select all filters (complete state)
export const selectAllListingFilters = createSelector(
  [selectListingFiltersState],
  (filters) => filters
);
// Main listing type
export const selectListingType = createSelector(
  [selectListingFiltersState],
  (filters) => filters.listingType
);
// Filter parameters grouped by type
export const selectPaginationFilters = createSelector(
  [selectListingFiltersState],
  (filters) => ({
    page: filters.page,
    pageSize: filters.pageSize,
    offset: filters.offset,
    limit: filters.limit,
  })
);
export const selectDisplayFilters = createSelector(
  [selectListingFiltersState],
  (filters) => ({
    displayMode: filters.displayMode,
    sortBy: filters.sortBy,
    sortOrder: filters.sortOrder,
    excludeNoMedia: filters.excludeNoMedia
  })
);
export const selectSearchFilters = createSelector(
  [selectListingFiltersState],
  (filters) => ({
    searchText: filters.searchText,
    category: filters.category,
    tags: filters.tags,
    location: filters.location
  })
);
export const selectPriceFilters = createSelector(
  [selectListingFiltersState],
  (filters) => ({
    minPrice: filters.minPrice,
    maxPrice: filters.maxPrice,
    negotiable: filters.negotiable
  })
);
export const selectProductSpecificFilters = createSelector(
  [selectListingFiltersState],
  (filters) => ({
    condition: filters.condition,
    brand: filters.brand,
    model: filters.model,
    hasVariants: filters.hasVariants,
    minStock: filters.minStock,
    manageStock: filters.manageStock,
    rating: filters.rating,
    userType: filters.userType,
    sku: filters.sku
  })
); 