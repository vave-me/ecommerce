/**
 * Consolidated React Query configuration and utilities
 * Combines functionality from both reactQuery.js and queryClient.js
 */
import { QueryClient } from '@tanstack/react-query';
import { createSyncStoragePersister } from '@tanstack/query-sync-storage-persister';
import { persistQueryClient } from '@tanstack/react-query-persist-client';
/**
 * Cache time configurations for different data types
 * Optimized based on data characteristics
 */
export const CACHE_TIMES = {
  STATIC: 24 * 60 * 60 * 1000,     // 24 hours - rarely changing data
  SEMI_STATIC: 60 * 60 * 1000,     // 1 hour - infrequently changing data
  STANDARD: 5 * 60 * 1000,         // 5 minutes - normal data
  DYNAMIC: 30 * 1000,              // 30 seconds - frequently changing data
  REALTIME: 0,                     // Always fresh - critical data
  // Legacy aliases for backward compatibility
  SHORT: 1 * 60 * 1000,            // 1 minute
  EXTENDED: 30 * 60 * 1000,        // 30 minutes
  LONG: 60 * 60 * 1000,            // 1 hour
  PERSIST: 24 * 60 * 60 * 1000,    // 24 hours
};
/**
 * Standardized query key factory for consistent cache management
 * Structured to enable granular invalidation
 */
export const QUERY_KEYS = {
  listings: (filters = {}) => ['listings', filters],
  product: (id) => ['product', id],
  products: (params = {}) => ['products', params],
  categories: {
    all: (type = 'all') => ['categories', type],
    main: (type, locale) => ['categories', 'main', type, locale],
    sub: (parentId, locale) => ['categories', 'sub', parentId, locale],
  },
  search: (params = {}) => ['search', params],
  config: (name = 'app') => ['config', name],
  user: {
    current: () => ['user', 'current'],        // current logged-in user
    byId: (id) => ['user', id],               // any user by id
    settings: () => ['user', 'settings'],
    orders: (status) => ['user', 'orders', status],
    messages: (id) => ['user', 'messages', id],
  },
  comments: (itemId, itemType) => ['comments', itemType, itemId],
  entity: (type, id) => [type, id],
  media: (itemId, type) => ['media', type, itemId],
};
/**
 * Data-specific cache configurations based on data characteristics
 */
export const QUERY_CONFIGS = {
  // Static reference data
  categories: {
    staleTime: CACHE_TIMES.STATIC,
    gcTime: CACHE_TIMES.STATIC * 2,
    retry: 1,
    refetchOnMount: false,
    refetchOnWindowFocus: false,
    refetchOnReconnect: false,
  },
  // Configuration data
  config: {
    staleTime: CACHE_TIMES.SEMI_STATIC,
    gcTime: CACHE_TIMES.SEMI_STATIC * 2,
    retry: 2,
    refetchOnWindowFocus: false,
  },
  // Dynamic entity data
  listings: {
    staleTime: CACHE_TIMES.STANDARD,
    gcTime: CACHE_TIMES.STANDARD * 2, 
    keepPreviousData: true,
    refetchOnWindowFocus: true,
  },
  // User data
  user: {
    staleTime: CACHE_TIMES.DYNAMIC, 
    gcTime: CACHE_TIMES.STANDARD,
    retry: 1,
    refetchOnMount: true,
    refetchOnWindowFocus: true,
  },
  // Real-time data
  messages: {
    staleTime: CACHE_TIMES.REALTIME,
    gcTime: CACHE_TIMES.DYNAMIC,
    refetchInterval: 15000, // Poll every 15 seconds
  }
};
/**
 * Default options for common query types
 * For backward compatibility
 */
export const defaultQueryOptions = {
  entity: {
    staleTime: CACHE_TIMES.STANDARD,
    cacheTime: CACHE_TIMES.STANDARD * 2,
  },
  categories: {
    staleTime: CACHE_TIMES.STANDARD,
    cacheTime: CACHE_TIMES.EXTENDED,
    refetchOnWindowFocus: false,
    select: (data) => data?.categories || [],
    retry: 2,
  },
  prefetch: {
    staleTime: CACHE_TIMES.EXTENDED,
    cacheTime: CACHE_TIMES.LONG,
  },
};
/**
 * Create optimized query client with performance settings
 * @returns {QueryClient} - Configured query client
 */
export function createQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: CACHE_TIMES.STANDARD,
        gcTime: CACHE_TIMES.STANDARD * 2, // Using gcTime instead of cacheTime
        refetchOnWindowFocus: false,
        retry: (failureCount, error) => {
          // Don't retry on 404s
          if (error?.response?.status === 404) return false;
          // Don't retry on 401s - needs authentication
          if (error?.response?.status === 401) return false;
          // Retry up to 2 times for other errors
          return failureCount < 2;
        },
        // Improve UI experience during loading states
        keepPreviousData: true,
        // Prevent excessive network requests
        refetchOnMount: "always",
        refetchOnReconnect: "always",
      },
      mutations: {
        // Prevent redundant mutation attempts
        retry: 1,
        // Reduce network request during navigation
        throwOnError: true,
      },
    },
  });
}
/**
 * Query filter function for persister
 * Determines which queries should be persisted to storage
 * 
 * @param {Object} query - Query to filter
 * @returns {boolean} - Whether to persist this query
 */
export function shouldDehydrateQuery(query) {
  // Only persist queries with stable data
  const [queryType] = query.queryKey;
  return ['categories', 'config'].includes(queryType);
}
// For SSR/test environments, use a different approach
const isServer = typeof window === 'undefined';
const isTest = process.env.NODE_ENV === 'test';
// Create the singleton instance
let queryClient = createQueryClient();
// Setup persistence for the singleton in browser environments
if (!isServer && !isTest) {
  try {
    const persister = createSyncStoragePersister({
      storage: window.localStorage,
      key: 'CLASSIFIED_REACT_QUERY_CACHE',
      throttleTime: 1000,
    });
    persistQueryClient({
      queryClient,
      persister,
      maxAge: CACHE_TIMES.PERSIST,
      dehydrateOptions: {
        shouldDehydrateQuery,
      },
    });
  } catch (e) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', e);
        }
    }
}
// Export the singleton instance of QueryClient
export { queryClient };
/**
 * Helper function to invalidate entity queries when mutations occur
 * @param {string} entityType - Type of entity ('product', 'service', 'deal', etc.)
 * @param {string} entityId - ID of the entity that was mutated
 */
export const invalidateEntityQueries = (entityType, entityId) => {
  // Invalidate the specific entity
  queryClient.invalidateQueries(QUERY_KEYS.entity(entityType, entityId));
  // Also invalidate any list queries for this entity type
  queryClient.invalidateQueries(QUERY_KEYS.listings({ entityType }));
  // Optional plural invalidation is covered by partial matching; no need.
};