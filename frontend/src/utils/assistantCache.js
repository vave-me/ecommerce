/**
 * Assistant Cache Management Utilities
 * 
 * Provides functions to manage and clear assistant-related caches
 * to prevent stale data issues in the AI components
 */

/**
 * Clear all assistant-related data from localStorage and sessionStorage
 */
export const clearAssistantCache = () => {
    if (typeof window === 'undefined') return;
    
    try {
        // Clear specific assistant cache keys
        const assistantKeys = [
            'assistant_cache',
            'conversation_cache',
            'ai_mode_state',
            'selected_assistant',
            'active_conversation'
        ];
        
        assistantKeys.forEach(key => {
            localStorage.removeItem(key);
            sessionStorage.removeItem(key);
        });
        
        // Clear any keys that start with 'assistant_' or 'ai_'
        const allKeys = Object.keys(localStorage);
        allKeys.forEach(key => {
            if (key.startsWith('assistant_') || key.startsWith('ai_')) {
                localStorage.removeItem(key);
            }
        });

    } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
};

/**
 * Clear React Query cache for assistant queries
 */
export const clearAssistantQueryCache = (queryClient) => {
    if (!queryClient) return;
    
    try {
        // Invalidate all assistant-related queries
        queryClient.invalidateQueries(['assistants']);
        queryClient.invalidateQueries(['conversations']);
        queryClient.invalidateQueries(['ai']);
        
        // Remove the queries from cache entirely
        queryClient.removeQueries(['assistants']);
        queryClient.removeQueries(['conversations']);
        queryClient.removeQueries(['ai']);

    } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
};

/**
 * Force refresh all assistant data
 */
export const forceRefreshAssistantData = async (hooks) => {
    const { loadAssistants, clearCache } = hooks;
    
    try {
        // Clear all caches
        clearAssistantCache();
        if (clearCache) clearCache();
        
        // Force reload assistants
        if (loadAssistants) {
            await loadAssistants({ skipCache: true });
        }

    } catch (error) {
        // Error: 'Error refreshing assistant data:', error...
        throw error;
    }
};

export default {
    clearAssistantCache,
    clearAssistantQueryCache,
    forceRefreshAssistantData
};