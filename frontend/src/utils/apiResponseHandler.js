/**
 * Enhanced API Response Handler with retry and timeout management
 */

/**
 * Default configuration for API requests
 */
export const API_CONFIG = {
    DEFAULT_TIMEOUT: 10000, // 10 seconds
    MAX_RETRIES: 2,
    RETRY_DELAY: 1000, // 1 second
    TIMEOUT_ERRORS: [504, 502, 503],
    NETWORK_ERRORS: ['ECONNABORTED', 'ETIMEDOUT', 'ERR_NETWORK']
};

/**
 * Check if an error is a timeout error
 */
export function isTimeoutError(error) {
    if (!error) return false;
    
    // Check for specific timeout status codes
    if (error.response && API_CONFIG.TIMEOUT_ERRORS.includes(error.response.status)) {
        return true;
    }
    
    // Check for network timeout errors
    if (error.code && API_CONFIG.NETWORK_ERRORS.includes(error.code)) {
        return true;
    }
    
    // Check for timeout in error message
    if (error.message && error.message.toLowerCase().includes('timeout')) {
        return true;
    }
    
    return false;
}

/**
 * Handle API response with enhanced error handling
 */
export async function handleApiResponse(apiCall, options = {}) {
    const {
        retries = 0,
        retryDelay = API_CONFIG.RETRY_DELAY,
        showError = true,
        errorMessage = 'An error occurred',
        onTimeout = null
    } = options;
    
    let lastError = null;
    let attempt = 0;
    
    while (attempt <= retries) {
        try {
            const response = await apiCall();
            return {
                success: true,
                data: response.data,
                error: null
            };
        } catch (error) {
            lastError = error;
            
            // Handle timeout specifically
            if (isTimeoutError(error)) {

                // Call timeout handler if provided
                if (onTimeout) {
                    onTimeout(error);
                }
                
                // Don't retry on first attempt if it's a timeout
                if (attempt === 0 && !retries) {
                    return {
                        success: false,
                        data: null,
                        error: 'Request timed out. Please try again.',
                        isTimeout: true,
                        status: error.response?.status
                    };
                }
            }
            
            // If we have retries left and it's a retryable error
            if (attempt < retries && (isTimeoutError(error) || error.response?.status >= 500)) {
                attempt++;
                ...`);
                await new Promise(resolve => setTimeout(resolve, retryDelay));
                continue;
            }
            
            // Log non-timeout errors
            if (!isTimeoutError(error) && showError) {
                // Error: 'API Error:', error...
            }
            
            break;
        }
    }
    
    // Return error response
    return {
        success: false,
        data: null,
        error: lastError.response?.data?.message || lastError.message || errorMessage,
        status: lastError.response?.status,
        isTimeout: isTimeoutError(lastError)
    };
}

/**
 * Create a timeout-aware axios instance
 */
export function createTimeoutAxios(axiosInstance, defaultTimeout = API_CONFIG.DEFAULT_TIMEOUT) {
    // Add request interceptor to set timeout
    axiosInstance.interceptors.request.use(
        (config) => {
            // Set timeout if not already set
            if (!config.timeout) {
                config.timeout = defaultTimeout;
            }
            return config;
        },
        (error) => Promise.reject(error)
    );
    
    // Add response interceptor for timeout handling
    axiosInstance.interceptors.response.use(
        (response) => response,
        (error) => {
            if (isTimeoutError(error)) {
                // Enhance error message for timeouts
                error.message = `Request timed out after ${error.config?.timeout || defaultTimeout}ms`;
            }
            return Promise.reject(error);
        }
    );
    
    return axiosInstance;
}

/**
 * Debounce API calls to prevent rapid successive calls
 */
export function debounceApi(apiFunction, delay = 300) {
    let timeoutId = null;
    let pendingPromise = null;
    
    return function(...args) {
        // If there's a pending promise, return it
        if (pendingPromise) {
            return pendingPromise;
        }
        
        // Clear existing timeout
        if (timeoutId) {
            clearTimeout(timeoutId);
        }
        
        // Create new promise
        pendingPromise = new Promise((resolve, reject) => {
            timeoutId = setTimeout(async () => {
                try {
                    const result = await apiFunction(...args);
                    resolve(result);
                } catch (error) {
                    reject(error);
                } finally {
                    pendingPromise = null;
                    timeoutId = null;
                }
            }, delay);
        });
        
        return pendingPromise;
    };
}