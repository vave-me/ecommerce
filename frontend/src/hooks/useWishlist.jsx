import { useState, useCallback, useEffect } from 'react';
import { toast } from 'react-toastify';
import { useAuth } from '../context/AuthContext';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import * as wishlistApi from '../api/client/wishlistApi';

// Query Keys
const WISHLIST_KEYS = {
  all: ['wishlists'],
  lists: () => [...WISHLIST_KEYS.all, 'list'],
  list: (id) => [...WISHLIST_KEYS.all, 'list', id],
  details: (id) => [...WISHLIST_KEYS.all, 'details', id],
  check: (productId) => [...WISHLIST_KEYS.all, 'check', productId],
};

/**
 * Internal React Query hooks
 */

// Get all wishlists for the current user
function useWishlists(options = {}) {
  const { user } = useAuth();
  const isAuthenticated = !!user?.userId;
  
  return useQuery({
    queryKey: WISHLIST_KEYS.lists(),
    queryFn: wishlistApi.getWishlists,
    enabled: isAuthenticated && (options.enabled !== false),
    ...options,
  });
}

// Get wishlist items by ID
function useWishlistItems(wishlistId, options = {}) {
  const { user } = useAuth();
  const isAuthenticated = !!user?.userId;
  
  return useQuery({
    queryKey: WISHLIST_KEYS.details(wishlistId),
    queryFn: () => wishlistApi.getWishlistItems(wishlistId),
    enabled: !!wishlistId && isAuthenticated && (options.enabled !== false),
    ...options,
  });
}

// Create a new wishlist
function useCreateWishlist() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: wishlistApi.createWishlist,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: WISHLIST_KEYS.lists() });
    },
  });
}

// Add item to wishlist
function useAddToWishlist() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ wishlistId, itemId, entityType = 'product', notes = '' }) => 
      wishlistApi.addWishlistItem(wishlistId, itemId, entityType, notes),
    onSuccess: (_, { wishlistId, itemId }) => {
      queryClient.invalidateQueries({ queryKey: WISHLIST_KEYS.details(wishlistId) });
      queryClient.invalidateQueries({ queryKey: WISHLIST_KEYS.check(itemId) });
      queryClient.invalidateQueries({ queryKey: WISHLIST_KEYS.lists() });
    },
  });
}

// Remove item from wishlist
function useRemoveFromWishlist() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ wishlistId, itemId }) => 
      wishlistApi.removeWishlistItem(wishlistId, itemId),
    onSuccess: (_, { wishlistId, itemId }) => {
      queryClient.invalidateQueries({ queryKey: WISHLIST_KEYS.details(wishlistId) });
      queryClient.invalidateQueries({ queryKey: WISHLIST_KEYS.check(itemId) });
      queryClient.invalidateQueries({ queryKey: WISHLIST_KEYS.lists() });
    },
  });
}

// Check if item is in any wishlist
function useIsInWishlist(userId, itemId, entityType = 'product', options = {}) {
  return useQuery({
    queryKey: WISHLIST_KEYS.check(itemId),
    queryFn: () => wishlistApi.checkItemInWishlists(userId, itemId, entityType),
    enabled: !!itemId && !!userId,
    ...options,
  });
}

// Delete a wishlist
function useDeleteWishlist() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: wishlistApi.removeWishlist,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: WISHLIST_KEYS.all });
    },
  });
}

/**
 * Main Wishlist Hook
 * Provides complete wishlist management with caching and optimistic updates
 */
export default function useWishlist() {
  const { user } = useAuth();
  const userId = user?.userId;
  const isLoggedIn = !!userId;
  
  // React Query hooks
  const { data: wishlistData, isLoading: wishlistsLoading, error: wishlistsError } = useWishlists();
  const createWishlistMutation = useCreateWishlist();
  const addToWishlistMutation = useAddToWishlist();
  const removeFromWishlistMutation = useRemoveFromWishlist();
  const deleteWishlistMutation = useDeleteWishlist();
  
  // Local state
  const [currentWishlist, setCurrentWishlist] = useState(null);
  const [showWishlistSelector, setShowWishlistSelector] = useState(false);
  const [pendingItemToAdd, setPendingItemToAdd] = useState(null);
  
  // Extract wishlists array from response
  const wishlists = Array.isArray(wishlistData) ? wishlistData : 
                   (wishlistData?.wishlists || wishlistData?.items || []);
  
  // Get default wishlist or first wishlist
  const defaultWishlist = Array.isArray(wishlists) ? 
    (wishlists.find(w => w.isDefault || w.name === 'Default') || wishlists[0] || null) : null;
  
  // Get items for current wishlist
  const { data: itemsData, isLoading: itemsLoading } = useWishlistItems(currentWishlist?.id);
  const items = itemsData?.items || [];
  
  // Combined loading state
  const loading = wishlistsLoading || itemsLoading || 
                 createWishlistMutation.isLoading ||
                 addToWishlistMutation.isLoading ||
                 removeFromWishlistMutation.isLoading ||
                 deleteWishlistMutation.isLoading;
  
  // Combined error state
  const error = wishlistsError || createWishlistMutation.error || 
               addToWishlistMutation.error || removeFromWishlistMutation.error;
  
  // ===== UTILITY FUNCTIONS =====
  
  /**
   * Clear all state - used when user logs out
   */
  const clearState = useCallback(() => {
    setCurrentWishlist(null);
    setShowWishlistSelector(false);
    setPendingItemToAdd(null);
  }, []);
  
  /**
   * Handle errors consistently
   */
  const handleError = useCallback((error, message, showToast = true) => {
    if (process.env.NODE_ENV === 'development') {
      // Error: 'Wishlist Error:', error...
    }
    if (showToast) {
      toast.error(message);
    }
  }, []);
  
  // ===== WISHLIST MANAGEMENT =====
  
  /**
   * Load default wishlist
   */
  const loadDefaultWishlist = useCallback(async (silent = false) => {
    if (!isLoggedIn) {
      clearState();
      return;
    }
    
    if (defaultWishlist) {
      setCurrentWishlist(defaultWishlist);
    } else if (!wishlistsLoading && wishlists.length === 0) {
      // Create default wishlist if none exists
      try {
        const response = await createWishlistMutation.mutateAsync('Default');
        // API returns { id: string }
        const newWishlist = { 
          id: response.id, 
          name: 'Default',
          userId: userId 
        };
        setCurrentWishlist(newWishlist);
      } catch (err) {
        handleError(err, 'Failed to create default wishlist', !silent);
      }
    }
  }, [isLoggedIn, defaultWishlist, wishlists, wishlistsLoading, createWishlistMutation, clearState, handleError]);
  
  /**
   * Load a specific wishlist by name
   */
  const loadWishlistByName = useCallback(async (name = 'Default', silent = false) => {
    if (!isLoggedIn) {
      return null;
    }
    
    const wishlist = wishlists.find(w => w.name === name);
    if (wishlist) {
      setCurrentWishlist(wishlist);
      return wishlist.id;
    } else {
      if (!silent) {
        handleError(null, `Wishlist "${name}" not found`, true);
      }
      return null;
    }
  }, [isLoggedIn, wishlists, handleError]);
  
  /**
   * Create a new wishlist
   */
  const createNewWishlist = useCallback(async (name = 'New Wishlist') => {
    if (!isLoggedIn) {
      toast.warn('Please log in to create a wishlist');
      return null;
    }
    
    try {
      // API only accepts name parameter
      const response = await createWishlistMutation.mutateAsync(name);
      toast.success(`Wishlist "${name}" created`);
      return response.id;
    } catch (err) {
      handleError(err, 'Failed to create wishlist');
      return null;
    }
  }, [isLoggedIn, createWishlistMutation, handleError]);
  
  /**
   * Select a wishlist to work with
   */
  const selectWishlist = useCallback(async (wishlistId) => {
    const selected = wishlists.find(w => w.id === wishlistId);
    if (selected) {
      setCurrentWishlist(selected);
    }
  }, [wishlists]);
  
  /**
   * Delete a wishlist
   */
  const deleteWishlist = useCallback(async (wishlistId) => {
    if (!isLoggedIn) {
      toast.warn('Please log in to manage your wishlists');
      return;
    }
    
    try {
      await deleteWishlistMutation.mutateAsync(wishlistId);
      
      // If we deleted the current wishlist, reset it
      if (currentWishlist?.id === wishlistId) {
        setCurrentWishlist(null);
      }
      
      toast.success('Wishlist deleted');
    } catch (err) {
      handleError(err, 'Failed to delete wishlist');
    }
  }, [isLoggedIn, currentWishlist, deleteWishlistMutation, handleError]);
  
  // ===== ITEM MANAGEMENT =====
  
  /**
   * Add item to a specific wishlist
   */
  const addItemToWishlist = useCallback(async (wishlistId, itemId, entityType = 'product', notes = '') => {
    if (!isLoggedIn || !wishlistId) {
      return false;
    }
    
    try {
      await addToWishlistMutation.mutateAsync({
        wishlistId,
        itemId,
        entityType,
        notes
      });
      
      // Find wishlist name for better user feedback
      const wishlist = wishlists.find(w => w.id === wishlistId);
      const wishlistName = wishlist?.name || 'wishlist';
      toast.success(`Item added to ${wishlistName}`);
      return true;
    } catch (err) {
      handleError(err, 'Failed to add item to wishlist');
      return false;
    }
  }, [isLoggedIn, wishlists, addToWishlistMutation, handleError]);
  
  /**
   * Add item to default wishlist
   */
  const addToDefaultWishlist = useCallback(async (itemId, entityType = 'product', notes = '') => {
    if (!isLoggedIn) {
      toast.warn('Please log in to add items to wishlist');
      return false;
    }
    
    let wishlistId = defaultWishlist?.id;
    
    // Create default wishlist if none exists
    if (!wishlistId) {
      const newWishlistId = await createNewWishlist('Default', 'My default wishlist');
      if (!newWishlistId) return false;
      wishlistId = newWishlistId;
    }
    
    return addItemToWishlist(wishlistId, itemId, entityType, notes);
  }, [isLoggedIn, defaultWishlist, createNewWishlist, addItemToWishlist]);
  
  /**
   * Smart add item - handles wishlist selection and default creation
   */
  const addItem = useCallback(async (itemId, entityType = 'product', notes = '') => {
    if (!isLoggedIn) {
      toast.warn('Please log in to add items to wishlist');
      return;
    }
    
    // Validate entity type
    const allowedTypes = ['product', 'post', 'video', 'service', 'property', 'vehicle'];
    if (!allowedTypes.includes(entityType.toLowerCase())) {
      toast.error(`Only ${allowedTypes.join(', ')} can be added to wishlists`);
      return;
    }
    
    // Check if item already exists in any wishlist
    if (isInAnyWishlist(itemId, entityType)) {
      const wishlistWithItem = wishlists.find(wishlist =>
        wishlist.items?.some(item => 
          item.itemId === itemId && item.entityType === entityType
        )
      );
      if (wishlistWithItem) {
        toast.info(`Item is already in "${wishlistWithItem.name}"`);
        return;
      }
    }
    
    // If no wishlists or only one, add to default/first
    if (wishlists.length <= 1) {
      await addToDefaultWishlist(itemId, entityType, notes);
      return;
    }
    
    // Multiple wishlists - show selector
    setPendingItemToAdd({ itemId, entityType, notes });
    setShowWishlistSelector(true);
  }, [isLoggedIn, wishlists, addToDefaultWishlist]);
  
  /**
   * Confirm adding pending item to selected wishlist
   */
  const confirmAddToPendingWishlist = useCallback(async (wishlistId) => {
    if (!pendingItemToAdd) return;
    
    const success = await addItemToWishlist(
      wishlistId, 
      pendingItemToAdd.itemId, 
      pendingItemToAdd.entityType,
      pendingItemToAdd.notes || ''
    );
    
    if (success) {
      setShowWishlistSelector(false);
      setPendingItemToAdd(null);
    }
  }, [pendingItemToAdd, addItemToWishlist]);
  
  /**
   * Cancel adding pending item
   */
  const cancelAddToPendingWishlist = useCallback(() => {
    setShowWishlistSelector(false);
    setPendingItemToAdd(null);
  }, []);
  
  /**
   * Remove item from wishlist
   */
  const removeItem = useCallback(async (productId) => {
    if (!isLoggedIn || !currentWishlist?.id) {
      toast.warn('Please select a wishlist first');
      return;
    }
    
    try {
      // Find the wishlist item by product ID
      const wishlistItem = items.find(item => item.itemId === productId);
      if (!wishlistItem) {
        toast.error('Item not found in wishlist');
        return;
      }
      
      // Pass the wishlist item ID, not the product ID
      await removeFromWishlistMutation.mutateAsync({
        wishlistId: currentWishlist.id,
        itemId: wishlistItem.id  // This is the wishlist item ID
      });
      toast.success('Item removed from wishlist');
    } catch (err) {
      handleError(err, 'Failed to remove item from wishlist');
    }
  }, [isLoggedIn, currentWishlist, items, removeFromWishlistMutation, handleError]);
  
  /**
   * Remove item from all wishlists
   */
  const removeFromAllWishlists = useCallback(async (itemId, entityType = 'product') => {
    if (!isLoggedIn) {
      return;
    }
    
    // Find which wishlists contain this item
    const wishlistsWithItem = wishlists.filter(wishlist =>
      wishlist.items?.some(item => 
        item.itemId === itemId && item.entityType === entityType
      )
    );
    
    // Remove from all wishlists
    const promises = wishlistsWithItem.map(wishlist => {
      const wishlistItem = wishlist.items.find(item => 
        item.itemId === itemId && item.entityType === entityType
      );
      return wishlistItem ? removeFromWishlistMutation.mutateAsync({
        wishlistId: wishlist.id,
        itemId: wishlistItem.id  // Use the wishlist item ID
      }) : Promise.resolve();
    });
    
    await Promise.allSettled(promises);
    toast.success('Item removed from all wishlists');
  }, [isLoggedIn, wishlists, removeFromWishlistMutation]);
  
  /**
   * Toggle item in current wishlist
   */
  const toggleItem = useCallback(async (itemId, entityType = 'product') => {
    const isInList = items.some(item => item.itemId === itemId);
    
    if (isInList) {
      await removeItem(itemId);
    } else {
      await addItem(itemId, entityType);
    }
  }, [items, addItem, removeItem]);
  
  /**
   * Check if item is in current wishlist
   */
  const isInWishlist = useCallback((itemId) => {
    return items.some(item => item.itemId === itemId);
  }, [items]);
  
  /**
   * Check if item is in any wishlist
   */
  const isInAnyWishlist = useCallback((itemId, entityType = 'product') => {
    return wishlists.some(wishlist =>
      wishlist.items?.some(item => 
        item.itemId === itemId && 
        item.entityType === entityType
      )
    );
  }, [wishlists]);
  
  /**
   * Toggle item in wishlist (enhanced version)
   */
  const toggleWishlist = useCallback(async (itemId, entityType = 'product', notes = '') => {
    if (isInAnyWishlist(itemId, entityType)) {
      await removeFromAllWishlists(itemId, entityType);
    } else {
      await addToDefaultWishlist(itemId, entityType, notes);
    }
  }, [isInAnyWishlist, removeFromAllWishlists, addToDefaultWishlist]);
  
  // Get total wishlist count
  const totalItemsCount = wishlists.reduce((sum, wishlist) => 
    sum + (wishlist.itemCount || wishlist.items?.length || 0), 0
  );
  
  // ===== EFFECTS =====
  
  // Set default wishlist on load
  useEffect(() => {
    if (!currentWishlist && defaultWishlist && !wishlistsLoading) {
      setCurrentWishlist(defaultWishlist);
    }
  }, [currentWishlist, defaultWishlist, wishlistsLoading]);
  
  // Clear state on logout
  useEffect(() => {
    if (!isLoggedIn) {
      clearState();
    }
  }, [isLoggedIn, clearState]);
  
  // ===== RETURN HOOK INTERFACE =====
  return {
    // State
    wishlists,
    currentWishlist,
    items,
    loading,
    error,
    showWishlistSelector,
    pendingItemToAdd,
    defaultWishlist,
    totalItemsCount,
    
    // Wishlist operations
    loadDefaultWishlist,
    loadWishlistByName,
    createNewWishlist,
    selectWishlist,
    deleteWishlist,
    createWishlist: createNewWishlist, // Alias for compatibility
    
    // Item operations
    loadWishlistItems: () => {}, // No-op, React Query handles this
    addItem,
    addItemToWishlist,
    addToDefaultWishlist,
    removeItem,
    removeFromAllWishlists,
    toggleItem,
    toggleWishlist,
    
    // Wishlist selector operations
    confirmAddToPendingWishlist,
    cancelAddToPendingWishlist,
    
    // Utility functions
    isInWishlist,
    isInAnyWishlist,
    clearState,
    
    // Loading states (for compatibility)
    isLoading: loading,
    isAddingToWishlist: addToWishlistMutation.isLoading,
    isRemovingFromWishlist: removeFromWishlistMutation.isLoading,
    isCreatingWishlist: createWishlistMutation.isLoading,
  };
}

// Export named function for imports that use it
export { default as useWishlist } from './useWishlist';