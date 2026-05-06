// src/api/client/wishlistsApi.jsx
import axiosInstance from '../axiosInstance';
import { secureTokenStorage } from '../../utils/secureTokenStorage';
/**
 * Enhanced Wishlist API - Updated with new endpoint structure
 * Handles all wishlist and wishlist item operations with proper error handling
 */
// ===== WISHLIST OPERATIONS =====
/**
 * Create a new wishlist
 * POST /api/wishlists/create
 * @param {string} name - Name of the wishlist
 * @returns {Promise<{id: string}>} New wishlist ID
 */
export const createWishlist = async (name = 'Default') => {
    try {
        const response = await axiosInstance.post('/wishlists/create', {
            name
        });
        return response.data;
    } catch (error) {
        if (process.env.NODE_ENV === 'development') {
        }
        throw error;
    }
};
/**
 * Get all wishlists for the authenticated user
 * GET /api/wishlists
 * @returns {Promise<{wishlists: Array}>} List of user's wishlists
 */
export const getWishlists = async () => {
    try {
        // Check if user is authenticated before making the request
        const token = secureTokenStorage.getAccessToken();
        
        if (!token) {
            // Return empty wishlists if not authenticated
            return { wishlists: [] };
        }
        
        const response = await axiosInstance.get('/wishlists');
        return response.data;
    } catch (error) {
        if (process.env.NODE_ENV === 'development') {
            // Error: 'Error fetching wishlists:', error...
        }
        // Return empty wishlists on error to prevent app crash
        if (error.response?.status === 401 || error.response?.status === 500) {
            return { wishlists: [] };
        }
        throw error;
    }
};
/**
 * Get a specific wishlist by name
 * GET /api/wishlists/name/{name}
 * @param {string} userId - User ID (not used in new API)
 * @param {string} name - Wishlist name
 * @returns {Promise<{wishlistId: string}>} Wishlist details
 */
export const getWishlist = async (userId, name = 'Default') => {
    try {
        const response = await axiosInstance.get(`/wishlists/name/${encodeURIComponent(name)}`);
        return response.data;
    } catch (error) {
        if (process.env.NODE_ENV === 'development') {
        }
        throw error;
    }
};
/**
 * Remove a wishlist
 * DELETE /api/wishlists/{id}
 * @param {string} id - Wishlist ID to remove
 * @returns {Promise<Object>} Empty response
 */
export const removeWishlist = async (id) => {
    try {
        const response = await axiosInstance.delete(`/wishlists/${id}`);
        return response.data;
    } catch (error) {
        if (process.env.NODE_ENV === 'development') {
        }
        throw error;
    }
};
// ===== WISHLIST ITEM OPERATIONS =====
/**
 * Get all items in a wishlist
 * GET /api/wishlists/{wishlistId}/items
 * @param {string} wishlistId - Wishlist ID
 * @returns {Promise<{items: Array}>} List of wishlist items
 */
export const getWishlistItems = async (wishlistId) => {
    try {
        // Check if user is authenticated
        const token = secureTokenStorage.getAccessToken();
        
        if (!token || !wishlistId) {
            return { items: [] };
        }
        
        const response = await axiosInstance.get(`/wishlists/${wishlistId}/items`);
        return response.data;
    } catch (error) {
        if (process.env.NODE_ENV === 'development') {
            // Error: 'Error fetching wishlist items:', error...
        }
        // Return empty items on error to prevent app crash
        if (error.response?.status === 401 || error.response?.status === 500) {
            return { items: [] };
        }
        throw error;
    }
};
/**
 * Get a specific item in a wishlist
 * GET /api/wishlists/{wishlistId}/items/{id}
 * @param {string} wishlistId - Wishlist ID
 * @param {string} id - Item ID in wishlist
 * @param {string} [itemId] - Optional item ID query parameter
 * @returns {Promise<{item: Object}>} Item details
 */
export const getWishlistItem = async (wishlistId, id, itemId = null) => {
    try {
        const url = `/wishlists/${wishlistId}/items/${id}`;
        const params = itemId ? { itemId } : {};
        const response = await axiosInstance.get(url, { params });
        return response.data;
    } catch (error) {
        if (process.env.NODE_ENV === 'development') {
        }
        throw error;
    }
};
/**
 * Add an item to a wishlist
 * POST /api/wishlists/{wishlistId}/items
 * @param {string} wishlistId - Wishlist ID
 * @param {string} itemId - Item ID to add
 * @param {string} entityType - Type of entity (product, post, video only)
 * @param {string} [notes] - Optional notes about the item
 * @returns {Promise<{id: string}>} Added item ID
 */
export const addWishlistItem = async (wishlistId, itemId, entityType = 'product', notes = '') => {
    // Validate entity type
    const allowedTypes = ['product', 'post', 'video'];
    if (!allowedTypes.includes(entityType.toLowerCase())) {
        throw new Error(`Invalid entity type: ${entityType}. Only ${allowedTypes.join(', ')} are supported.`);
    }
    
    try {
        const response = await axiosInstance.post(`/wishlists/${wishlistId}/items`, {
            itemId: itemId,
            entityType: entityType,
            ...(notes && { notes })
        });
        return response.data;
    } catch (error) {
        if (process.env.NODE_ENV === 'development') {
        }
        throw error;
    }
};
/**
 * Remove an item from a wishlist
 * DELETE /api/wishlists/{wishlistId}/items/{id}
 * @param {string} wishlistId - Wishlist ID
 * @param {string} id - Item ID to remove (wishlist item ID, not the product ID)
 * @returns {Promise<Object>} Empty response
 */
export const removeWishlistItem = async (wishlistId, id) => {
    try {
        const response = await axiosInstance.delete(`/wishlists/${wishlistId}/items/${id}`);
        return response.data;
    } catch (error) {
        if (process.env.NODE_ENV === 'development') {
        }
        throw error;
    }
};
// ===== UTILITY FUNCTIONS =====
// Note: Default wishlist is automatically created by backend when user is created
/**
 * Check if an item exists in any of the user's wishlists
 * @param {string} userId - User ID
 * @param {string} itemId - Item ID to check
 * @param {string} entityType - Entity type to check (product, post, video only)
 * @returns {Promise<{exists: boolean, wishlistId?: string, wishlistItemId?: string}>}
 */
export const checkItemInWishlists = async (userId, itemId, entityType = 'product') => {
    // Check if user is authenticated
    const token = secureTokenStorage.getAccessToken();
    
    if (!token || !userId) {
        return { exists: false };
    }
    
    // Validate entity type
    const allowedTypes = ['product', 'post', 'video'];
    if (!allowedTypes.includes(entityType.toLowerCase())) {
        return { exists: false }; // Unsupported types can't be in wishlists
    }
    
    try {
        const userWishlists = await getWishlists();
        for (const wishlist of userWishlists.wishlists || []) {
            try {
                const items = await getWishlistItems(wishlist.id);
                const foundItem = items.items?.find(
                    item => item.itemId === itemId && item.entityType === entityType
                );
                if (foundItem) {
                    return {
                        exists: true,
                        wishlistId: wishlist.id,
                        wishlistItemId: foundItem.id,
                        wishlistName: wishlist.name
                    };
                }
            } catch (itemsError) {
                // Skip this wishlist if we can't load its items
                if (process.env.NODE_ENV === 'development') {
                }
                continue;
            }
        }
        return { exists: false };
    } catch (error) {
        if (process.env.NODE_ENV === 'development') {
        }
        throw error;
    }
};
