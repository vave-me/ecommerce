// File: src/api/client/searchApi.jsx (Recommended new name)

// Assuming this is your configured Axios client instance
import axiosInstance from '../axiosInstance';

// --- Helper Functions (Consider moving to a shared utils file) ---

/**
 * Helper to safely encode URI components.
 * @param {string | number | undefined | null} component - The path component to encode.
 * @returns {string} The encoded component.
 */
const safeEncode = (component) => {
    if (component === null || typeof component === 'undefined') {
        
        return ''; // Or throw? Depends on whether it's truly optional upstream.
    }
    return encodeURIComponent(String(component));
};

/**
 * Helper to create a clean query params object, removing null/undefined/empty strings.
 * @param {object} filters - Input filters object.
 * @returns {object} Cleaned filters object for query params.
 */
function cleanFilters(filters = {}) {
    const cleaned = {};
    for (const key in filters) {
        const value = filters[key];
        // Keep 0, false, but remove null, undefined, empty strings
        if (value !== null && typeof value !== 'undefined' && value !== '') {
            cleaned[key] = value;
        }
    }
    return cleaned;
}

// --- Endpoint Constants ---
const SEARCH_ENDPOINT = '/search';

/**
 * UNIFIED CATALOG
 * GET /api/search/{user_id}/catalog
 * Returns unified catalog with all entities for a user
 */
export const getUnifiedCatalog = async (userId, catalogParams = {}) => {
    if (!userId) throw new Error('userId is required for getUnifiedCatalog.');

    const safeUserId = encodeURIComponent(userId);
    const endpoint = `${SEARCH_ENDPOINT}/${safeUserId}/catalog`;
    const queryParams = cleanFilters(catalogParams);

    try {
        if (process.env.NODE_ENV === 'development') {
            }
        
        const response = await axiosInstance.get(endpoint, { params: queryParams });
        
        if (process.env.NODE_ENV === 'development') {
            // Response logged for debugging
        }
        
        // Extra check: if the response contains userId information, verify it matches
        if (process.env.NODE_ENV === 'development' && response.data?.items && Array.isArray(response.data.items)) {
            const userIds = response.data.items.map(item => item.userId || item.ownerId || item.createdBy).filter(Boolean);
            const uniqueUserIds = [...new Set(userIds)];
            // User IDs checked for debugging
        }
        
        return response.data;
    } catch (error) {
        if (process.env.NODE_ENV === 'development') {
            // API error logged for debugging
        }
        throw error;
    }
};

/**
 * DEBUG FUNCTION: Test the unified catalog endpoint directly
 * Call this function in the browser console to test server response
 * Usage: window.debugUnifiedCatalog('your-user-id')
 */
export const debugUnifiedCatalog = async (userId) => {
    if (process.env.NODE_ENV !== 'development') return;
    try {
        const response = await getUnifiedCatalog(userId, {});
        // Analyze the response
        const analysis = {
            responseType: typeof response,
            isArray: Array.isArray(response),
            hasItems: !!response?.items,
            itemCount: response?.items?.length || (Array.isArray(response) ? response.length : 0),
            structure: response ? Object.keys(response) : null
        };
        
        if (response?.items && Array.isArray(response.items)) {
            // Check user ownership
            const userAnalysis = response.items.map(item => ({
                id: item.id,
                title: item.title || item.name,
                userId: item.userId || item.ownerId || item.createdBy,
                entityType: item.entityType || item.type
            }));
            
            const uniqueUsers = [...new Set(userAnalysis.map(item => item.userId).filter(Boolean))];
            
            analysis.userOwnership = {
                requestedUserId: userId,
                foundUsers: uniqueUsers,
                itemsFromRequestedUser: userAnalysis.filter(item => item.userId === userId).length,
                itemsFromOtherUsers: userAnalysis.filter(item => item.userId && item.userId !== userId).length,
                itemsWithoutUserId: userAnalysis.filter(item => !item.userId).length,
                sampleItems: userAnalysis.slice(0, 5)
            };
        }
        
        return { response, analysis };
    } catch (error) {
        // Error: '❌ [DEBUG TEST] Error:', error...
        return { error };
    }
};

// Make debug function available globally for testing
if (typeof window !== 'undefined') {
    window.debugUnifiedCatalog = debugUnifiedCatalog;
}

