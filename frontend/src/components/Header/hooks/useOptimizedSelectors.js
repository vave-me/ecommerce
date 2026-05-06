import { useSelector } from 'react-redux';
import { shallowEqual } from 'react-redux';
/**
 * Optimized Redux Selectors for Header Components
 * Reduces unnecessary re-renders and improves performance
 */
// Memoized selector for listing filters with shallow comparison
export const useListingFilters = () => {
    return useSelector(
        (state) => ({
            searchText: state.listingFilters.searchText,
            radius: state.listingFilters.radius,
            category: state.listingFilters.category,
            listingType: state.listingFilters.listingType,
            lat: state.listingFilters.lat,
            lng: state.listingFilters.lng,
        }),
        shallowEqual
    );
};
// Memoized selector for app mode with shallow comparison
export const useAppMode = () => {
    return useSelector(
        (state) => ({
            currentMode: state.appMode.currentMode,
            isAiMode: state.appMode.currentMode === 'ai',
            isTransitioning: state.appMode.isTransitioning,
        }),
        shallowEqual
    );
};
// Memoized selector for user data with shallow comparison
export const useUserData = () => {
    return useSelector(
        (state) => ({
            user: state.auth.user,
            isAuthenticated: !!state.auth.user,
            userId: state.auth.user?.userId,
        }),
        shallowEqual
    );
};
// Memoized selector for badge counts with shallow comparison
export const useBadgeCounts = () => {
    return useSelector(
        (state) => ({
            messages: state.notifications?.messages || 0,
            notifications: state.notifications?.notifications || 0,
            wishlist: state.wishlist?.count || 0,
            cart: state.cart?.count || 0,
        }),
        shallowEqual
    );
};
// Computed selector for active filters
export const useActiveFilters = () => {
    const filters = useListingFilters();
    return {
        hasActiveFilters: !!(filters.searchText || filters.lat || filters.lng || filters.radius !== 10 || filters.category),
        activeFilterCount: [
            filters.searchText,
            filters.lat,
            filters.lng,
            filters.radius !== 10,
            filters.category
        ].filter(Boolean).length,
        filters
    };
}; 