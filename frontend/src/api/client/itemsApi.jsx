import axiosInstance from '../axiosInstance';

/**
 * Fetch a generic item by its ID and type
 * GET /{entityType}/{itemId}
 * 
 * @param {string} itemId - ID of the item to fetch
 * @param {string} entityType - Type of entity (product, service, deal, etc.)
 * @returns {Promise<Object>} - Item details
 */
export const getGenericItem = async (itemId, entityType = 'product') => {
    if (!itemId) {
        throw new Error('itemId is required for getGenericItem');
    }
    
    // Convert singular entity type to plural form for API endpoint
    const endpoint = `${entityType}${entityType.endsWith('s') ? '' : 's'}/${itemId}`;
    
    try {
        const response = await axiosInstance.get(`/${endpoint}`);
        return response.data;
    } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    }:`,...
        throw error;
    }
};

/**
 * Fetch multiple items of the same type
 * GET /{entityType}?ids=id1,id2,id3
 * 
 * @param {string[]} itemIds - Array of item IDs to fetch
 * @param {string} entityType - Type of entity (product, service, deal, etc.)
 * @returns {Promise<Object[]>} - Array of item details
 */
export const getMultipleItems = async (itemIds, entityType = 'product') => {
    if (!itemIds || !Array.isArray(itemIds) || itemIds.length === 0) {
        throw new Error('itemIds array is required for getMultipleItems');
    }
    
    // Convert singular entity type to plural form for API endpoint
    const pluralType = entityType.endsWith('s') ? entityType : `${entityType}s`;
    const endpoint = `/${pluralType}`;
    
    try {
        const response = await axiosInstance.get(endpoint, {
            params: {
                ids: itemIds.join(',')
            }
        });
        return response.data;
    } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    } items:`, er...
        throw error;
    }
};

/**
 * Helper function to check if an item exists
 * HEAD /{entityType}/{itemId}
 * 
 * @param {string} itemId - ID of the item to check
 * @param {string} entityType - Type of entity
 * @returns {Promise<boolean>} - True if item exists
 */
export const checkItemExists = async (itemId, entityType = 'product') => {
    if (!itemId) {
        throw new Error('itemId is required for checkItemExists');
    }
    
    // Convert singular entity type to plural form for API endpoint
    const pluralType = entityType.endsWith('s') ? entityType : `${entityType}s`;
    const endpoint = `/${pluralType}/${itemId}`;
    
    try {
        await axiosInstance.head(endpoint);
        return true; // If no error is thrown, the item exists
    } catch (error) {
        if (error.response?.status === 404) {
            return false; // Item doesn't exist
        }
        // Error: `Error checking if ${entityType} exists:`, error.r...
        throw error;
    }
}; 