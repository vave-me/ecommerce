"use client";
import React, {useEffect, memo} from 'react';
import {usePrefetchCategories} from '../hooks/useCategories';
/**
 * Component that handles app initialization tasks like data prefetching
 * This should be mounted once near the root of the app
 *
 * OPTIMIZED: Memoized to prevent unnecessary re-initialization
 */
const AppInitializer = memo(function AppInitializer() {
    // Initialize category prefetching
    usePrefetchCategories();
    // Additional initialization can be added here
    // This component doesn't render anything visible
    return null;
});
export default AppInitializer; 