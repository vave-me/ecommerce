/**
 * Unified Basket Hook with React Query and Offline Support
 * 
 * This hook provides unified basket functionality with React Query 
 * and offline support for cart management.
 */
import { useEffect, useCallback, useRef } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useAuth } from '../context/AuthContext';
import * as basketApi from '../api/client/basketApi';
import { offlineBasketStorage, isClientSide } from '../utils/offlineStorage';

// Query keys
const BASKET_QUERY_KEY = 'basket';

// Create a global cache for basket creation promises to prevent race conditions
const basketCreationPromises = new Map();

/**
 * Main basket hook that handles both online and offline scenarios
 */
export function useBasket() {
  const { user } = useAuth();
  const queryClient = useQueryClient();
  const userId = user?.userId;
  const isLoggedIn = !!userId;
  
  // Fetch basket query
  const basketQuery = useQuery({
    queryKey: [BASKET_QUERY_KEY, userId],
    queryFn: async () => {
      if (!isLoggedIn) {
        // Return offline basket for non-logged users
        const offlineItems = offlineBasketStorage.getItems();
        return {
          id: 'offline-basket',
          items: offlineItems,
          basket_status: 'open',
          itemCount: offlineItems.reduce((sum, item) => sum + item.quantity, 0),
          totalAmount: offlineItems.reduce((sum, item) => 
            sum + ((item.product_price || item.price || 0) * item.quantity), 0
          )
        };
      }
      
      // Check if we already have a valid basket in cache
      const cachedBasket = queryClient.getQueryData([BASKET_QUERY_KEY, userId]);
      if (cachedBasket && cachedBasket.id && cachedBasket.id !== 'offline-basket') {
        // Validate the cached basket still exists
        try {
          const response = await basketApi.getBasket(cachedBasket.id);
          if (response?.basket) {
            const basket = response.basket;
            const itemCount = basket.items?.reduce((sum, item) => sum + (item.quantity || 0), 0) || 0;
            const totalAmount = basket.items?.reduce((sum, item) => 
              sum + ((item.product_price || 0) * (item.quantity || 0)), 0
            ) || 0;
            
            return {
              ...basket,
              itemCount,
              totalAmount
            };
          }
        } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
      }
      
      // Try to get current basket WITHOUT creating a new one
      try {
        const currentBasket = await basketApi.getCurrentBasket(userId);
        
        if (currentBasket && currentBasket.basket_id) {
          // We have an existing basket, fetch its details
          const response = await basketApi.getBasket(currentBasket.basket_id);
          const basket = response?.basket || { 
            id: currentBasket.basket_id,
            items: [], 
            basket_status: currentBasket.basket_status || 'open'
          };
          
          // Calculate totals
          const itemCount = basket.items?.reduce((sum, item) => sum + (item.quantity || 0), 0) || 0;
          const totalAmount = basket.items?.reduce((sum, item) => 
            sum + ((item.product_price || 0) * (item.quantity || 0)), 0
          ) || 0;
          
          return {
            ...basket,
            itemCount,
            totalAmount
          };
        }
      } catch (error) {
        if (error.response?.status !== 404) {
          // Only log non-404 errors
          // Error: 'Error fetching current basket:', error...
        }
      }
      
      // No basket exists - return empty basket state
      // Basket will be created only when user adds first item
      return {
        id: null,
        items: [],
        basket_status: 'empty',
        itemCount: 0,
        totalAmount: 0
      };
    },
    enabled: true, // Always enabled, handles both online/offline
    staleTime: 5 * 60 * 1000, // 5 minutes - increased from 2 minutes
    gcTime: 10 * 60 * 1000, // 10 minutes - using gcTime instead of cacheTime
    refetchOnWindowFocus: false, // Disabled to prevent excessive basket creation
    refetchOnMount: 'always',
    refetchOnReconnect: 'always',
  });
  
  // Add to basket mutation
  const addToBasketMutation = useMutation({
    mutationFn: async ({ productId, quantity = 1, price, metadata = {} }) => {
      if (!isLoggedIn) {
        // Handle offline addition
        const success = offlineBasketStorage.addItem(productId, quantity, price, metadata);
        // Item added to basket (will sync when you log in)
        return { success, offline: true };
      }
      
      // Get current basket ID
      const currentBasket = queryClient.getQueryData([BASKET_QUERY_KEY, userId]);
      let basketId = currentBasket?.id;
      
      // Only create basket if we don't have one and are adding an item
      if (!basketId || basketId === 'offline-basket') {
        // Check if we're already creating a basket for this user
        if (basketCreationPromises.has(userId)) {
          const result = await basketCreationPromises.get(userId);
          basketId = result.basketId;
        } else {
          // Create basket only when adding first item
          const creationPromise = basketApi.startBasket(userId);
          basketCreationPromises.set(userId, creationPromise);
          
          try {
            const result = await creationPromise;
            basketId = result.id;
            
            // Update the query cache with the new basket
            queryClient.setQueryData([BASKET_QUERY_KEY, userId], {
              id: basketId,
              items: [],
              basket_status: 'open',
              itemCount: 0,
              totalAmount: 0
            });
          } finally {
            setTimeout(() => {
              basketCreationPromises.delete(userId);
            }, 1000);
          }
        }
      }
      
      const response = await basketApi.addItemToBasket(basketId, productId, quantity);
      return response;
    },
    onMutate: async ({ productId, quantity, price }) => {
      await queryClient.cancelQueries([BASKET_QUERY_KEY]);
      
      const previousData = queryClient.getQueryData([BASKET_QUERY_KEY, userId]);
      
      queryClient.setQueryData([BASKET_QUERY_KEY, userId], (old) => {
        if (!old) return old;
        
        const existingItemIndex = old.items?.findIndex(item => 
          (item.product_id === productId || item.productId === productId)
        );
        
        let newItems = [...(old.items || [])];
        
        if (existingItemIndex >= 0) {
          // Update existing item
          newItems[existingItemIndex] = {
            ...newItems[existingItemIndex],
            quantity: newItems[existingItemIndex].quantity + quantity
          };
        } else {
          // Add new item
          newItems.push({
            product_id: productId,
            quantity,
            product_price: price || 0,
            addedAt: new Date().toISOString()
          });
        }
        
        const itemCount = newItems.reduce((sum, item) => sum + item.quantity, 0);
        const totalAmount = newItems.reduce((sum, item) => 
          sum + ((item.product_price || 0) * item.quantity), 0
        );
        
        return {
          ...old,
          items: newItems,
          itemCount,
          totalAmount
        };
      });
      
      return { previousData };
    },
    onError: (err, variables, context) => {
      if (context?.previousData) {
        queryClient.setQueryData([BASKET_QUERY_KEY, userId], context.previousData);
      }
      // Error: 'Failed to add item to basket:', err...
    },
    onSettled: (data) => {
      if (!data?.offline && isLoggedIn) {
        queryClient.invalidateQueries([BASKET_QUERY_KEY]);
      }
    }
  });
  
  // Update quantity mutation
  const updateQuantityMutation = useMutation({
    mutationFn: async ({ itemId, quantity }) => {
      if (!isLoggedIn) {
        const success = offlineBasketStorage.updateQuantity(itemId, quantity);
        // Quantity updated
        return { success, offline: true };
      }
      
      const currentBasket = queryClient.getQueryData([BASKET_QUERY_KEY, userId]);
      if (!currentBasket?.id || currentBasket.id === 'offline-basket') {
        throw new Error('No basket found');
      }
      
      const response = await basketApi.updateItemQuantity(currentBasket.id, itemId, quantity);
      return response;
    },
    onMutate: async ({ itemId, quantity }) => {
      await queryClient.cancelQueries([BASKET_QUERY_KEY]);
      
      const previousData = queryClient.getQueryData([BASKET_QUERY_KEY, userId]);
      
      queryClient.setQueryData([BASKET_QUERY_KEY, userId], (old) => {
        if (!old) return old;
        
        const newItems = old.items?.map(item => {
          if (item.item_id === itemId || item.id === itemId || item.product_id === itemId) {
            return { ...item, quantity };
          }
          return item;
        }) || [];
        
        const itemCount = newItems.reduce((sum, item) => sum + item.quantity, 0);
        const totalAmount = newItems.reduce((sum, item) => 
          sum + ((item.product_price || 0) * item.quantity), 0
        );
        
        return {
          ...old,
          items: newItems,
          itemCount,
          totalAmount
        };
      });
      
      return { previousData };
    },
    onError: (err, variables, context) => {
      if (context?.previousData) {
        queryClient.setQueryData([BASKET_QUERY_KEY, userId], context.previousData);
      }
      // Error: 'Failed to update quantity:', err...
    },
    onSettled: (data) => {
      if (!data?.offline && isLoggedIn) {
        queryClient.invalidateQueries([BASKET_QUERY_KEY]);
      }
    }
  });
  
  // Remove item mutation
  const removeItemMutation = useMutation({
    mutationFn: async ({ itemId }) => {
      if (!isLoggedIn) {
        const success = offlineBasketStorage.removeItem(itemId);
        // Item removed from basket
        return { success, offline: true };
      }
      
      const currentBasket = queryClient.getQueryData([BASKET_QUERY_KEY, userId]);
      if (!currentBasket?.id || currentBasket.id === 'offline-basket') {
        throw new Error('No basket found');
      }
      
      const response = await basketApi.removeItemFromBasket(currentBasket.id, itemId);
      return response;
    },
    onMutate: async ({ itemId }) => {
      await queryClient.cancelQueries([BASKET_QUERY_KEY]);
      
      const previousData = queryClient.getQueryData([BASKET_QUERY_KEY, userId]);
      
      queryClient.setQueryData([BASKET_QUERY_KEY, userId], (old) => {
        if (!old) return old;
        
        const newItems = old.items?.filter(item => 
          item.item_id !== itemId && item.id !== itemId && item.product_id !== itemId
        ) || [];
        
        const itemCount = newItems.reduce((sum, item) => sum + item.quantity, 0);
        const totalAmount = newItems.reduce((sum, item) => 
          sum + ((item.product_price || 0) * item.quantity), 0
        );
        
        return {
          ...old,
          items: newItems,
          itemCount,
          totalAmount
        };
      });
      
      return { previousData };
    },
    onError: (err, variables, context) => {
      if (context?.previousData) {
        queryClient.setQueryData([BASKET_QUERY_KEY, userId], context.previousData);
      }
      // Error: 'Failed to remove item:', err...
    },
    onSettled: (data) => {
      if (!data?.offline && isLoggedIn) {
        queryClient.invalidateQueries([BASKET_QUERY_KEY]);
      }
    }
  });
  
  // Sync offline basket mutation
  const syncOfflineBasketMutation = useMutation({
    mutationFn: async () => {
      if (!userId) return null;
      
      const offlineItems = offlineBasketStorage.getItems();
      if (offlineItems.length === 0) return null;
      
      // Create basket for syncing offline items
      let basketId;
      
      // First check if user already has a basket
      try {
        const currentBasket = await basketApi.getCurrentBasket(userId);
        if (currentBasket && currentBasket.basket_id) {
          basketId = currentBasket.basket_id;
        }
      } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    }
      
      // Only create if no basket exists
      if (!basketId) {
        if (basketCreationPromises.has(userId)) {
          const result = await basketCreationPromises.get(userId);
          basketId = result.id;
        } else {
          const creationPromise = basketApi.startBasket(userId);
          basketCreationPromises.set(userId, creationPromise);
          
          try {
            const result = await creationPromise;
            basketId = result.id;
          } finally {
            setTimeout(() => {
              basketCreationPromises.delete(userId);
            }, 1000);
          }
        }
      }
      
      let successCount = 0;
      let failCount = 0;
      
      for (const item of offlineItems) {
        try {
          await basketApi.addItemToBasket(basketId, item.productId || item.product_id, item.quantity);
          successCount++;
        } catch (error) {
          // Error: 'Failed to sync basket item:', item, error...
          failCount++;
        }
      }
      
      // Clear offline storage after sync
      offlineBasketStorage.clear();
      
      return { successCount, failCount };
    },
    onSuccess: (result) => {
      if (result) {
        const { successCount, failCount } = result;
        if (successCount > 0 && failCount === 0) {
          
        } else if (successCount > 0 && failCount > 0) {
          
        } else if (failCount > 0) {
          // Error: 'Failed to sync offline basket items'...
        }
      }
      queryClient.invalidateQueries([BASKET_QUERY_KEY]);
    },
    onError: () => {
      // Error: 'Failed to sync offline basket'...
    }
  });
  
  // Auto-sync offline basket when user logs in
  useEffect(() => {
    if (isLoggedIn && isClientSide() && syncOfflineBasketMutation.isIdle) {
      const offlineItems = offlineBasketStorage.getItems();
      if (offlineItems.length > 0) {
        syncOfflineBasketMutation.mutate();
      }
    }
  }, [isLoggedIn, syncOfflineBasketMutation]);
  
  // Helper functions
  const isInBasket = useCallback((productId) => {
    const basket = basketQuery.data;
    return basket?.items?.some(item => 
      item.product_id === productId || item.productId === productId
    ) || false;
  }, [basketQuery.data]);
  
  const getItemQuantity = useCallback((productId) => {
    const basket = basketQuery.data;
    const item = basket?.items?.find(item => 
      item.product_id === productId || item.productId === productId
    );
    return item?.quantity || 0;
  }, [basketQuery.data]);
  
  const initiateCheckout = useCallback(() => {
    const basket = basketQuery.data;
    if (!basket?.id || basket.id === 'offline-basket') {
      
      return null;
    }
    
    if (!basket.items || basket.items.length === 0) {
      
      return null;
    }
    
    return {
      basketId: basket.id,
      userCustomerId: userId,
      amount: basket.totalAmount || 0, // Already in cents
      items: basket.items
    };
  }, [basketQuery.data, userId]);
  
  const checkoutBasket = useCallback(async (paymentIntentId) => {
    const basket = basketQuery.data;
    if (!basket?.id || basket.id === 'offline-basket') {
      throw new Error('No basket to checkout');
    }
    
    if (!paymentIntentId) {
      throw new Error('Payment intent is required for checkout');
    }
    
    const response = await basketApi.checkoutBasket(basket.id, userId, paymentIntentId);
    
    // Clear the basket from cache after checkout
    queryClient.setQueryData([BASKET_QUERY_KEY, userId], null);
    
    return response;
  }, [basketQuery.data, userId, queryClient]);
  
  return {
    // Data
    basket: basketQuery.data,
    items: basketQuery.data?.items || [],
    itemCount: basketQuery.data?.itemCount || 0,
    totalAmount: basketQuery.data?.totalAmount || 0,
    hasItems: (basketQuery.data?.items?.length || 0) > 0,
    
    // Loading states
    isLoading: basketQuery.isLoading,
    error: basketQuery.error,
    isAddingToBasket: addToBasketMutation.isLoading,
    isUpdatingQuantity: updateQuantityMutation.isLoading,
    isRemovingItem: removeItemMutation.isLoading,
    isSyncing: syncOfflineBasketMutation.isLoading,
    
    // Actions
    addToBasket: (productId, quantity = 1, price = 0, metadata = {}) => 
      addToBasketMutation.mutate({ productId, quantity, price, metadata }),
    updateQuantity: (itemId, quantity) => 
      updateQuantityMutation.mutate({ itemId, quantity }),
    removeItem: (itemId) => 
      removeItemMutation.mutate({ itemId }),
    refetchBasket: () => basketQuery.refetch(),
    syncOfflineBasket: () => syncOfflineBasketMutation.mutate(),
    
    // Checkout
    initiateCheckout,
    checkoutBasket,
    
    // Utilities
    isInBasket,
    getItemQuantity,
    isLoggedIn,
  };
}

// Export as default for backward compatibility
export default useBasket;