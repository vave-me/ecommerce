"use client";
import React, { memo } from 'react';
import { QueryClientProvider } from '@tanstack/react-query';
import FilteredFeed from './FilteredFeed';
import { queryClient } from '../../lib/reactQuery';
// FeedProvider removed - using direct hooks

/**
 * Client-side wrapper for FilteredFeed component
 * Handles providers for context, Redux and React Query
 * OPTIMIZED: React.memo for performance
 */
const FilteredFeedWrapper = memo(function FilteredFeedWrapper() {
  return (
    <QueryClientProvider client={queryClient}>
      <FilteredFeed />
    </QueryClientProvider>
  );
});

export default FilteredFeedWrapper; 