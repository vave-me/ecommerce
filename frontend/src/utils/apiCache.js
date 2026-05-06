/**
 * Simple in-memory cache for API requests with deduplication
 */
class APICache {
    constructor() {
        this.cache = new Map();
        this.pendingRequests = new Map();
        this.cacheTimeout = 5 * 60 * 1000; // 5 minutes default
    }

    /**
     * Get cache key from URL and params
     */
    getCacheKey(url, params = {}) {
        const sortedParams = Object.keys(params)
            .sort()
            .reduce((acc, key) => {
                acc[key] = params[key];
                return acc;
            }, {});
        return `${url}:${JSON.stringify(sortedParams)}`;
    }

    /**
     * Check if cache entry is still valid
     */
    isValid(entry) {
        return entry && Date.now() - entry.timestamp < this.cacheTimeout;
    }

    /**
     * Get from cache
     */
    get(url, params) {
        const key = this.getCacheKey(url, params);
        const entry = this.cache.get(key);
        
        if (this.isValid(entry)) {
            
            return entry.data;
        }
        
        // Remove expired entry
        if (entry) {
            this.cache.delete(key);
        }
        
        return null;
    }

    /**
     * Set cache entry
     */
    set(url, params, data) {
        const key = this.getCacheKey(url, params);
        this.cache.set(key, {
            data,
            timestamp: Date.now()
        });
    }

    /**
     * Get or fetch with deduplication
     * Prevents multiple identical requests from being made simultaneously
     */
    async getOrFetch(url, params, fetchFn) {
        const key = this.getCacheKey(url, params);
        
        // Check cache first
        const cached = this.get(url, params);
        if (cached) {
            return cached;
        }
        
        // Check if request is already pending
        if (this.pendingRequests.has(key)) {
            
            return this.pendingRequests.get(key);
        }
        
        // Make the request and store the promise
        const requestPromise = fetchFn().then(
            (data) => {
                this.set(url, params, data);
                this.pendingRequests.delete(key);
                return data;
            },
            (error) => {
                this.pendingRequests.delete(key);
                throw error;
            }
        );
        
        this.pendingRequests.set(key, requestPromise);
        return requestPromise;
    }

    /**
     * Clear cache for specific URL pattern
     */
    clear(urlPattern) {
        if (!urlPattern) {
            this.cache.clear();
            return;
        }
        
        for (const [key] of this.cache) {
            if (key.includes(urlPattern)) {
                this.cache.delete(key);
            }
        }
    }

    /**
     * Set cache timeout
     */
    setTimeout(timeout) {
        this.cacheTimeout = timeout;
    }
}

// Create singleton instance
const apiCache = new APICache();

// Specific caches for different API endpoints
export const categoryCache = new APICache();
categoryCache.setTimeout(10 * 60 * 1000); // 10 minutes for categories

export const assistantCache = new APICache();
assistantCache.setTimeout(5 * 60 * 1000); // 5 minutes for assistants

export const wishlistCache = new APICache();
wishlistCache.setTimeout(2 * 60 * 1000); // 2 minutes for wishlists

export default apiCache;