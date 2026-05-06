
// src/api/axiosInstance.jsx
import axios from 'axios';
import { secureTokenStorage } from '../utils/secureTokenStorage';

// Direct token access using secure storage
const getAuthToken = () => {
    return secureTokenStorage.getAccessToken();
};

const setAuthToken = (token, expiresIn) => {
    if (token) {
        secureTokenStorage.setAccessToken(token, expiresIn);
    } else {
        secureTokenStorage.clearTokens();
    }
};

const getRefreshToken = () => {
    return secureTokenStorage.getRefreshToken();
};

const setRefreshToken = (token) => {
    if (token) {
        secureTokenStorage.setRefreshToken(token);
    } else {
        secureTokenStorage.clearTokens();
    }
};

const isTokenExpired = (token) => {
    if (!token) return true;
    
    try {
        // Simple JWT decode - just get the payload
        const parts = token.split('.');
        if (parts.length !== 3) return true;
        
        const payload = JSON.parse(atob(parts[1]));
        if (!payload.exp) return true;
        
        // Check if expired (with 10 second buffer)
        const isExpired = payload.exp * 1000 <= Date.now() + 10000;
        return isExpired;
    } catch (error) {
        // Error: '🔐 Error checking token expiration:', error...
        return true;
    }
};

const clearTokens = async () => {
    secureTokenStorage.clearTokens();
};

// Simple refresh token function without circular dependency
const refreshAccessToken = async () => {
    const refreshToken = getRefreshToken();
    const currentAccessToken = getAuthToken();
    
    if (!refreshToken) {
        // Error: '🔐 refreshAccessToken: No refresh token available...
        throw new Error('No refresh token available');
    }

    try {

        // Create a simple axios instance for refresh
        const refreshAxios = axios.create({
            baseURL: process.env.NEXT_PUBLIC_API_BASE_URL,
            timeout: 10000,
            headers: {
                'Content-Type': 'application/json',
            },
        });
        
        // Add Authorization header as requested by user
        const headers = {};
        if (currentAccessToken) {
            headers.Authorization = `Bearer ${currentAccessToken}`;
            
        }
        
        const response = await refreshAxios.post('/users/refresh-token', {
            refreshToken
        }, { headers });

        const { accessToken, refreshToken: newRefreshToken } = response.data;

        setAuthToken(accessToken);
        if (newRefreshToken) {
            setRefreshToken(newRefreshToken);
        }

        return accessToken;
    } catch (error) {
        // Token refresh failed - logged for debugging
        // Clear tokens if refresh fails
        await clearTokens();
        throw error;
    }
};

// Environment configuration
const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL;
const isDev = process.env.NODE_ENV === 'development';
const isClient = typeof window !== 'undefined';

// Validate environment
if (!apiBaseUrl) {
    // Error: 'NEXT_PUBLIC_API_BASE_URL is not set!'...
}

/**
 * Main Axios Instance
 * Base URL should be: http://domain:port/api (already includes /api in env)
 * So all endpoint calls should NOT add /api prefix
 */
const axiosInstance = axios.create({
    baseURL: apiBaseUrl, // Already includes /api from env
    withCredentials: false,
    // NO TIMEOUT - let the request complete
    headers: {
        'Content-Type': 'application/json',
    },
});

// Create a separate SSR-optimized instance for server-side rendering
export const ssrAxiosInstance = axios.create({
    baseURL: apiBaseUrl,
    withCredentials: false, // Not needed for SSR
    timeout: 5000, // Aggressive timeout for SSR
    headers: {
        'Content-Type': 'application/json',
        'User-Agent': 'NextJS-SSR-Bot/1.0',
    },
});

// Track refresh state to avoid multiple refresh calls
let isRefreshing = false;
let failedQueue = [];
const MAX_RETRIES = 3;

// Calculate exponential backoff delay
const getBackoffDelay = (retryAttempt) => {
    return Math.min(1000 * Math.pow(2, retryAttempt), 10000);
};

// Process queued requests
const processQueue = (error, token = null) => {
    failedQueue.forEach(prom => {
        if (error) {
            prom.reject(error);
        } else {
            prom.resolve(token);
        }
    });
    failedQueue = [];
};

// Request interceptor - add auth token (only on client side)
axiosInstance.interceptors.request.use(
    async (config) => {
        // Only handle auth on client side
        if (isClient) {
            // Use direct token access - no circular dependency!
            const token = getAuthToken();
            
            // Only log important requests in development
            if (process.env.NODE_ENV === 'development' && !config.url?.includes('/media/')) {
                // Don't log requests that will likely fail due to no auth
                if (token || !config.url?.includes('/assistants')) {
                    // Request logged for debugging
                }
            }
            // Token logging removed for cleaner output
            
            if (token && !isTokenExpired(token)) {
                config.headers.Authorization = `Bearer ${token}`;
            } else if (token && isTokenExpired(token)) {
                // Token refresh handled silently
                // Token expired, try refresh before request
                try {
                    const newToken = await refreshAccessToken();
                    if (newToken) {
                        config.headers.Authorization = `Bearer ${newToken}`;
                        // Refreshed token added silently
                    }
                } catch (error) {
                    // Token refresh failure handled silently
                    // Continue without auth header
                }
            }
            
            // Request sent
        }
        return config;
    },
    (error) => {
        // Error: '❌ Request interceptor error:', error...
        return Promise.reject(error);
    }
);

// Response interceptor - handle auth errors and retries (only on client side)
axiosInstance.interceptors.response.use(
    (response) => {
        return response;
    },
    async (error) => {
        const originalRequest = error.config;
        
        // Handle network errors with retry logic
        if (!error.response) {
            originalRequest._retryCount = originalRequest._retryCount || 0;
            const shouldRetry = originalRequest._retryCount < 2 && 
                               (error.code === 'ECONNABORTED' || 
                                error.code === 'ENOTFOUND' || 
                                error.code === 'ECONNREFUSED' ||
                                error.message?.includes('timeout'));
            
            if (shouldRetry) {
                originalRequest._retryCount++;
                const delay = Math.min(1000 * Math.pow(2, originalRequest._retryCount), 3000);
                // Retrying after delay
                await new Promise(resolve => setTimeout(resolve, delay));
                return axiosInstance(originalRequest);
            }
            
            // Error details logged for debugging
        }
        return Promise.reject(error);
    }
);

// Export token functions for compatibility (but we don't use injection anymore)
export const tokenFunctions = {
    getAccessToken: getAuthToken,
    setAccessToken: setAuthToken,
    getRefreshToken,
    setRefreshToken,
    clearTokens,
    refreshAccessToken,
    isTokenValid: (token) => token && !isTokenExpired(token)
};

// Legacy function for compatibility - does nothing now
export const injectTokenFunctions = (functions) => {
    
};

export default axiosInstance;