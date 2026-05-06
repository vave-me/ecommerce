"use client";
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { CACHE_TIMES, QUERY_KEYS, invalidateEntityQueries, defaultQueryOptions } from '../lib/reactQuery';
import { getEntity, updateEntity, createEntity, deleteEntity, getEntities } from '../api/client/entityApi';

/**
 * Unified hook system for entity management (products, services, deals, etc.)
 * Combines functionality from both useEntity and useEntityCache
 */

/**
 * Generic hook for fetching and caching entity data
 * @param {string} entityType - Type of entity ('product', 'service', 'deal')
 * @param {string} entityId - ID of the entity to fetch
 * @param {Function} fetchFn - Optional custom fetch function (defaults to getEntity API)
 * @param {Object} options - Additional options for the query
 * @returns {Object} Query result with entity data
 */
export function useEntity(entityType, entityId, fetchFn, options = {}) {
  // If fetchFn is actually options (backward compatibility)
  if (typeof fetchFn === 'object' && !options.queryFn) {
    options = fetchFn;
    fetchFn = null;
  }
  
  const defaultOptions = defaultQueryOptions?.entity || {};
  const { ...otherOptions } = options;
  
  return useQuery({
    queryKey: QUERY_KEYS.entity(entityType, entityId),
    queryFn: async () => {
      if (!entityId) {
        throw new Error(`${entityType} ID is required to fetch data`);
      }
      // Use custom fetch function or default to API
      if (fetchFn) {
        return await fetchFn(entityId);
      }
      return await getEntity(entityType, entityId);
    },
    enabled: !!entityId,
    ...defaultOptions,
    ...otherOptions,
    onError: (error) => {
      if (options.onError) {
        options.onError(error);
      }
    },
  });
}

/**
 * Generic hook for updating entity data with automatic cache invalidation
 * @param {string} entityType - Type of entity ('product', 'service', 'deal')
 * @param {Function} updateFn - Optional custom update function (defaults to updateEntity API)
 * @param {Object} options - Additional options for the mutation
 * @returns {Object} Mutation result with mutate function
 */
export function useEntityUpdate(entityType, updateFn, options = {}) {
  const queryClient = useQueryClient();
  
  // If updateFn is actually options (backward compatibility)
  if (typeof updateFn === 'object' && !options.mutationFn) {
    options = updateFn;
    updateFn = null;
  }
  
  return useMutation({
    mutationFn: async (data) => {
      const entityId = data.id || data.entityId;
      if (!entityId) {
        throw new Error(`${entityType} ID is required to update ${entityType}`);
      }
      
      // Use custom update function or default to API
      if (updateFn) {
        return await updateFn(data);
      }
      
      const updateData = data.updateData || data;
      return await updateEntity(entityType, entityId, updateData);
    },
    onSuccess: (data, variables) => {
      // Get the ID from the variables
      const entityId = variables.id || variables.entityId;
      
      // Use the centralized invalidation helper
      invalidateEntityQueries(entityType, entityId);
      
      // Also invalidate list queries
      queryClient.invalidateQueries(QUERY_KEYS.listings({ entityType }));
      
      // Call any additional success handlers from options
      if (options.onSuccess) {
        options.onSuccess(data, variables);
      }
    },
    onError: (error, variables) => {
      // Call any additional error handlers from options
      if (options.onError) {
        options.onError(error, variables);
      }
    },
  });
}

/**
 * Create a new entity
 * @param {string} entityType - Type of entity ('product', 'service', 'deal', etc.)
 * @param {Object} options - Additional mutation options
 * @returns {Object} Mutation result with mutate function
 */
export function useEntityCreate(entityType, options = {}) {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (data) => createEntity(entityType, data),
    onSuccess: (data, variables) => {
      // Invalidate entity list queries
      queryClient.invalidateQueries(QUERY_KEYS.listings({ entityType }));
      
      // Call any additional success handlers from options
      if (options.onSuccess) {
        options.onSuccess(data, variables);
      }
    },
    onError: (error, variables) => {
      // Call any additional error handlers from options
      if (options.onError) {
        options.onError(error, variables);
      }
    },
  });
}

/**
 * Delete an entity
 * @param {string} entityType - Type of entity ('product', 'service', 'deal', etc.)
 * @param {Object} options - Additional mutation options
 * @returns {Object} Mutation result with mutate function
 */
export function useEntityDelete(entityType, options = {}) {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (id) => deleteEntity(entityType, id),
    onSuccess: (_, entityId) => {
      // Invalidate the specific entity
      queryClient.invalidateQueries(QUERY_KEYS.entity(entityType, entityId));
      // Also invalidate list queries
      queryClient.invalidateQueries(QUERY_KEYS.listings({ entityType }));
      
      // Call any additional success handlers from options
      if (options.onSuccess) {
        options.onSuccess(_, entityId);
      }
    },
    onError: (error, entityId) => {
      // Call any additional error handlers from options
      if (options.onError) {
        options.onError(error, entityId);
      }
    },
  });
}

/**
 * Get a list of entities with optional filtering
 * @param {string} entityType - Type of entity ('product', 'service', 'deal', etc.)
 * @param {Object} params - Query parameters for filtering
 * @param {Object} options - Additional query options
 * @returns {Object} Query result with entity list
 */
export function useEntityList(entityType, params = {}, options = {}) {
  return useQuery({
    queryKey: QUERY_KEYS.listings({ entityType, ...params }),
    queryFn: () => getEntities(entityType, params),
    ...options,
  });
}

// Type-specific hooks for backwards compatibility
// Product hooks
export function useProduct(id, options = {}) {
  return useEntity('product', id, options);
}

export function useProductUpdate(options = {}) {
  return useEntityUpdate('product', options);
}

export function useProductCreate(options = {}) {
  return useEntityCreate('product', options);
}

export function useProductDelete(options = {}) {
  return useEntityDelete('product', options);
}

export function useProductList(params = {}, options = {}) {
  return useEntityList('product', params, options);
}

// Service hooks
export function useService(id, options = {}) {
  return useEntity('service', id, options);
}

export function useServiceUpdate(options = {}) {
  return useEntityUpdate('service', options);
}

export function useServiceCreate(options = {}) {
  return useEntityCreate('service', options);
}

export function useServiceDelete(options = {}) {
  return useEntityDelete('service', options);
}

export function useServiceList(params = {}, options = {}) {
  return useEntityList('service', params, options);
}

// Deal hooks
export function useDeal(id, options = {}) {
  return useEntity('deal', id, options);
}

export function useDealUpdate(options = {}) {
  return useEntityUpdate('deal', options);
}

export function useDealCreate(options = {}) {
  return useEntityCreate('deal', options);
}

export function useDealDelete(options = {}) {
  return useEntityDelete('deal', options);
}

export function useDealList(params = {}, options = {}) {
  return useEntityList('deal', params, options);
}

// Property hooks
export function useProperty(id, options = {}) {
  return useEntity('property', id, options);
}

export function usePropertyUpdate(options = {}) {
  return useEntityUpdate('property', options);
}

export function usePropertyCreate(options = {}) {
  return useEntityCreate('property', options);
}

export function usePropertyDelete(options = {}) {
  return useEntityDelete('property', options);
}

export function usePropertyList(params = {}, options = {}) {
  return useEntityList('property', params, options);
}

// Job hooks
export function useJob(id, options = {}) {
  return useEntity('job', id, options);
}

export function useJobUpdate(options = {}) {
  return useEntityUpdate('job', options);
}

export function useJobCreate(options = {}) {
  return useEntityCreate('job', options);
}

export function useJobDelete(options = {}) {
  return useEntityDelete('job', options);
}

export function useJobList(params = {}, options = {}) {
  return useEntityList('job', params, options);
}

// Vehicle hooks
export function useVehicle(id, options = {}) {
  return useEntity('vehicle', id, options);
}

export function useVehicleUpdate(options = {}) {
  return useEntityUpdate('vehicle', options);
}

export function useVehicleCreate(options = {}) {
  return useEntityCreate('vehicle', options);
}

export function useVehicleDelete(options = {}) {
  return useEntityDelete('vehicle', options);
}

export function useVehicleList(params = {}, options = {}) {
  return useEntityList('vehicle', params, options);
}

// Export the generic cache functions for backward compatibility
export { useEntity as useEntityCache, useEntityUpdate as useEntityCacheUpdate };