"use client";
import React, { createContext, useContext, useMemo, useCallback, useState, memo } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useLocale } from 'next-intl';
import { fetchMainCategories, fetchSubCategories } from '../api/categories';
import { QUERY_KEYS } from '../lib/reactQuery';
// Define category types used throughout the application
export const CATEGORY_TYPES = [
  'marketplace',
  'jobs', 
  'automotive', 
  'property', 
  'posts',
  'deals',
  'services'
];
// Create a context for categories
export const CategoriesContext = createContext(null);
/**
 * Categories Provider component that fetches and caches categories data
 */
export const CategoriesProvider = memo(function CategoriesProvider({ children, prefetchTopics = [] }) {
    const locale = useLocale();
    const queryClient = useQueryClient();
    // Keep track of which category types have been requested
    const [requestedTypes, setRequestedTypes] = useState(new Set(prefetchTopics));
    // Method to request a specific category type
    const requestCategoryType = useCallback((type) => {
        if (!requestedTypes.has(type)) {
            setRequestedTypes(prev => new Set([...prev, type]));
        }
    }, [requestedTypes]);
    // Helper to get query key for a category type
    const getCategoryQueryKey = useCallback((type) => {
        return QUERY_KEYS.categories.all(`${type}-${locale}`);
    }, [locale]);
    // Prefetch categories for given types
    const prefetchCategories = useCallback(async (types = []) => {
        const promises = types.map(type => 
            queryClient.prefetchQuery({
                queryKey: getCategoryQueryKey(type),
                queryFn: () => fetchMainCategories({ categoryType: type, lang: locale })
            })
        );
        await Promise.all(promises);
        // Mark these types as requested
        setRequestedTypes(prev => new Set([...prev, ...types]));
    }, [queryClient, getCategoryQueryKey, locale]);
    // Prefetch specified topics on mount or locale change
    React.useEffect(() => {
        if (prefetchTopics.length > 0) {
            prefetchCategories(prefetchTopics);
        }
    }, [locale, prefetchTopics, prefetchCategories]); // No categoryCache dependency
    // Create the context value
    const contextValue = useMemo(() => ({
        requestCategoryType,
        prefetchCategories,
        getCategoryQueryKey,
        requestedTypes: Array.from(requestedTypes)
    }), [requestCategoryType, prefetchCategories, getCategoryQueryKey, requestedTypes]);
    return (
        <CategoriesContext.Provider value={contextValue}>
            {children}
        </CategoriesContext.Provider>
    );
});
/**
 * Hook for accessing categories data with proper caching
 * @param {string} categoryType - Type of categories to fetch
 * @returns {Object} Category data and related state
 */
export function useCategories(categoryType = 'marketplace') {
    const locale = useLocale();
    const context = useContext(CategoriesContext);
    if (!context) {
        throw new Error('useCategories must be used within a CategoriesProvider');
    }
    const { requestCategoryType, getCategoryQueryKey } = context;
    // Request this category type if not already requested
    React.useEffect(() => {
        requestCategoryType(categoryType);
    }, [categoryType, requestCategoryType]);
    // Fetch categories using React Query
    const queryKey = getCategoryQueryKey(categoryType);
    const queryResult = useQuery({
        queryKey,
        queryFn: () => fetchMainCategories({ 
            categoryType, 
            lang: locale
        }),
        // Only fetch when this category type has been requested
        enabled: context.requestedTypes.includes(categoryType),
        // Transform the response
        select: (data) => {
            // Map categories to match the expected format
            const categories = (data?.categories || []).map(cat => ({
                slug: cat.slug || '',
                id: cat.id,
                name: cat.name || cat.description || 'Unknown Category',
                description: cat.description,
                featured: cat.featured,
                count: cat.count,
                subcategories: cat.subcategories || []
            }));
            return categories;
        },
        // Avoid unnecessary fetches
        staleTime: 5 * 60 * 1000, // 5 minutes
        // Don't refetch on window focus to prevent unnecessary requests
        refetchOnWindowFocus: false
    });
    return {
        ...queryResult,
        data: queryResult.data || []
    };
}
// Export for backward compatibility
export function usePrefetchCategories() {
    const context = useContext(CategoriesContext);
    if (!context) {
        throw new Error('usePrefetchCategories must be used within a CategoriesProvider');
    }
    return context.prefetchCategories;
}

/**
 * Hook to fetch and cache main categories with React Query
 * @param {string} categoryType - Type of categories to fetch
 * @param {Object} options - Additional query options
 */
export function useMainCategories(categoryType = "all", options = {}) {
  const locale = useLocale();
  const queryClient = useQueryClient();

  const queryKey = QUERY_KEYS.categories.main(categoryType, locale);

  return useQuery({
    queryKey,
    queryFn: async () => {
      const filters = {
        categoryType,
        lang: locale,
        page: 0,
        pageSize: 50,
      };
      const result = await fetchMainCategories(filters);
      
      // Transform the data for easier consumption
      const categories = (result?.categories || []).map(cat => ({
        value: cat.slug || "",
        label: cat.name || cat.description || "Unknown Category",
        id: cat.id,
        featured: cat.featured,
        count: cat.count,
        hasSubcategories: true,
        subcategories: cat.subcategories || [],
      }));

      return {
        categories,
        totalCount: result.totalCount,
        totalPages: result.totalPages,
      };
    },
    staleTime: 30 * 60 * 1000, // 30 minutes
    gcTime: 60 * 60 * 1000, // 1 hour
    refetchOnWindowFocus: false,
    refetchOnMount: false,
    ...options,
  });
}

/**
 * Hook to fetch subcategories for a parent category
 * @param {string} parentId - Parent category ID
 * @param {Object} options - Additional query options
 */
export function useSubCategories(parentId, options = {}) {
  const locale = useLocale();

  const queryKey = QUERY_KEYS.categories.sub(parentId, locale);

  return useQuery({
    queryKey,
    queryFn: async () => {
      const filters = { lang: locale };
      const subcats = await fetchSubCategories(parentId, filters);
      
      // Transform the data
      return subcats.map(cat => ({
        value: cat.slug || "",
        label: cat.name || cat.description || "Unknown Category",
        id: cat.id,
        featured: cat.featured,
        count: cat.count,
        hasSubcategories: false,
        parentId: parentId,
        subcategories: cat.subcategories || [],
      }));
    },
    enabled: !!parentId,
    staleTime: 30 * 60 * 1000,
    gcTime: 60 * 60 * 1000,
    refetchOnWindowFocus: false,
    ...options,
  });
} 