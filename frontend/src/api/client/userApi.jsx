"use client"
import axiosInstance, { injectTokenFunctions } from '../axiosInstance';
import {jwtDecode} from "jwt-decode";
import axios from 'axios';
import { secureTokenStorage } from '../../utils/secureTokenStorage';

// Create a public axios instance for non-authenticated endpoints
const publicAxios = axios.create({
    baseURL: process.env.NEXT_PUBLIC_API_BASE_URL || 'http://192.168.178.84:8080/api',
    timeout: 10000,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Token storage is now handled by secureTokenStorage utility

// ===== TOKEN MANAGEMENT =====

/**
 * Store access token in localStorage
 * @param {string} token - JWT access token
 */
export const setAccessToken = (token) => {
    if (token) {
        secureTokenStorage.setAccessToken(token);
    } else {
        secureTokenStorage.clearTokens();
    }
};

/**
 * Get access token from localStorage
 * @returns {string|null} Current access token or null
 */
export const getAccessToken = () => {
    return secureTokenStorage.getAccessToken();
};

/**
 * Validate JWT token and check expiration
 * @param {string} token - JWT access token
 * @returns {boolean} Whether token is valid
 */
export const isTokenValid = (token) => {
    if (!token || typeof token !== 'string') {
        if (process.env.NODE_ENV === 'development') {
            
        }
        return false;
    }

    try {
        const decoded = jwtDecode(token);
        
        if (!decoded.exp) {
            if (process.env.NODE_ENV === 'development') {
                
            }
            return false;
        }
        
        const isValid = decoded.exp * 1000 > Date.now() + 10000;
        
        if (process.env.NODE_ENV === 'development') {
            const expiresAt = new Date(decoded.exp * 1000);
            // Token validation logged for debugging
        }
        
        return isValid;
    } catch (error) {
        if (process.env.NODE_ENV === 'development') {
            // Error: '🔐 Token validation error:', error.message
        }
        return false;
    }
};

/**
 * Extract user ID from JWT token
 * @param {string} token - JWT token
 * @returns {string|null} User ID or null if token invalid
 */
export const getUserIdFromToken = (token) => {
    if (!token) return null;

    try {
        const decoded = jwtDecode(token);
        return decoded.userId || decoded.sub || null;
    } catch (error) {
        return null;
    }
};

/**
 * Get refresh token from localStorage
 * @returns {string|null} Refresh token or null
 */
export const getRefreshToken = () => {
    return secureTokenStorage.getRefreshToken();
};

/**
 * Set refresh token in localStorage
 * @param {string} token - Refresh token
 */
export const setRefreshToken = (token) => {
    if (token) {
        secureTokenStorage.setRefreshToken(token);
    } else {
        secureTokenStorage.clearTokens();
    }
};

/**
 * Clear all tokens from localStorage
 */
export const clearTokens = async () => {
    if (typeof window !== 'undefined') {
        localStorage.removeItem(ACCESS_TOKEN_KEY);
        localStorage.removeItem(REFRESH_TOKEN_KEY);
    }

    try {
        const currentToken = getAccessToken();
        if (currentToken) {
            // Try to clear server-side tokens
            await axiosInstance.post('/users/clear-tokens', {});
        }
    } catch (error) {
        // Silently ignore 401 errors during token clearing as it's expected when tokens are invalid
        if (error.response?.status !== 401) {
            
        }
    }
};

/**
 * Clear user tokens with specific parameters (security endpoint)
 * @param {string} userId - User ID
 * @param {string} tokenId - Token ID (optional)
 * @param {string} refreshToken - Refresh token (optional)  
 * @param {string} reason - Reason for clearing tokens
 * @returns {Promise<Object>} Clear tokens response
 */
export const clearUserTokens = async (userId, tokenId = '', refreshToken = '', reason = 'logout') => {
    try {
        const currentToken = getAccessToken();
        
        if (!currentToken) {
            throw new Error('No access token available for authentication');
        }
        
        if (!isTokenValid(currentToken)) {
            throw new Error('Current access token is invalid or expired');
        }
        
        if (!userId) {
            throw new Error('User ID is required to clear tokens');
        }

        // Include current access token as tokenId if not provided (common pattern)
        const actualTokenId = tokenId || currentToken;
        
        const requestBody = {
            userId,
            tokenId: actualTokenId,
            refreshToken: refreshToken || getRefreshToken() || '',
            reason
        };

        // The interceptor should now handle the Authorization header properly
        const response = await axiosInstance.post('/users/clear-tokens', requestBody);
        return response.data;
    } catch (error) {
        if (process.env.NODE_ENV === 'development') {
            // Error logged for debugging
        }
        throw error;
    }
};

/**
 * Refresh access token using refresh token
 * @returns {Promise<string>} New access token
 */
export const refreshAccessToken = async () => {
    const refreshToken = getRefreshToken();
    const currentAccessToken = getAccessToken();
    
    if (!refreshToken) {
        // Error: '🔐 refreshAccessToken: No refresh token available
        throw new Error('No refresh token available');
    }

    try {
        
        // Add Authorization header to refresh request as requested by user
        const headers = {};
        if (currentAccessToken) {
            headers.Authorization = `Bearer ${currentAccessToken}`;
            
        }
        
        const response = await publicAxios.post('/users/refresh-token', {
            refreshToken
        }, { headers });

        const { accessToken, refreshToken: newRefreshToken } = response.data;

        setAccessToken(accessToken);
        if (newRefreshToken) {
            setRefreshToken(newRefreshToken);
        }

        return accessToken;
    } catch (error) {
        // Error logged for debugging
        // Clear tokens if refresh fails
        await clearTokens();
        throw error;
    }
};

// Note: Token injection is no longer needed - axiosInstance uses direct token access
// This avoids circular dependencies and timing issues

// ===== AUTHENTICATION METHODS =====

/**
 * Register a new user (public endpoint)
 * @param {Object} userData - User registration data
 * @returns {Promise<Object>} Registration response
 */
export const registerUser = async (userData) => {
    try {
        const response = await publicAxios.post('/users', userData);
        return response.data;
    } catch (error) {
        // Error: 'Registration error:', error
        throw error;
    }
};

/**
 * Login user with email/password (public endpoint)
 * @param {Object} credentials - Login credentials
 * @returns {Promise<Object>} Login response with tokens
 */
export const loginUser = async (credentials) => {
    try {
        
        const response = await publicAxios.post('/users/login', credentials);

        // Store tokens if login successful
        if (response.data.accessToken || response.data.token) {
            const token = response.data.accessToken || response.data.token;
            // Token stored successfully
            setAccessToken(token);
            
            if (response.data.refreshToken) {
                
                setRefreshToken(response.data.refreshToken);
            }
        }
        
        return response.data;
    } catch (error) {
        // Error: 'Login error:', error
        throw error;
    }
};

/**
 * Login with Google (public endpoint)
 * @param {string} idToken - Google ID token
 * @returns {Promise<Object>} Login response with tokens
 */
export const loginWithGoogle = async (idToken) => {
    try {
        const response = await publicAxios.post('/users/google-login', { idToken });
        
        // Store tokens if login successful
        if (response.data.accessToken || response.data.token) {
            const token = response.data.accessToken || response.data.token;
            setAccessToken(token);
            
            if (response.data.refreshToken) {
                setRefreshToken(response.data.refreshToken);
            }
        }
        
        return response.data;
    } catch (error) {
        // Error: 'Google login error:', error
        throw error;
    }
};

/**
 * Mobile login with Google (public endpoint)
 * @param {string} idToken - Google ID token
 * @returns {Promise<Object>} Login response with tokens
 */
export const loginWithGoogleMobile = async (idToken) => {
    try {
        const response = await publicAxios.post('/users/google-login/mobile', { idToken });
        
        // Store tokens if login successful
        if (response.data.accessToken || response.data.token) {
            const token = response.data.accessToken || response.data.token;
            setAccessToken(token);
            
            if (response.data.refreshToken) {
                setRefreshToken(response.data.refreshToken);
            }
        }
        
        return response.data;
    } catch (error) {
        // Error: 'Mobile Google login error:', error
        throw error;
    }
};

/**
 * Forgot password (public endpoint)
 * @param {string} email - User email
 * @returns {Promise<Object>} Forgot password response
 */
export const forgotPassword = async (email) => {
    try {
        const response = await publicAxios.post('/users/forgot-password', { email });
        return response.data;
    } catch (error) {
        // Error: 'Forgot password error:', error
        throw error;
    }
};

/**
 * Reset password (public endpoint)
 * @param {string} token - Reset token
 * @param {string} newPassword - New password
 * @returns {Promise<Object>} Reset password response
 */
export const resetPassword = async (token, newPassword) => {
    try {
        const response = await publicAxios.post('/users/reset-password', { 
            token, 
            newPassword 
        });
        return response.data;
    } catch (error) {
        // Error: 'Reset password error:', error
        throw error;
    }
};

/**
 * Logout user and clear tokens (requires auth)
 * @returns {Promise<void>}
 */
export const logoutUser = async () => {
    try {
        const userId = getUserIdFromToken(getAccessToken());
        const authToken = getAccessToken();
        const refreshToken = getRefreshToken();
        
        if (userId) {
            await axiosInstance.post('/users/logout', { 
                id: userId,
                authToken,
                refreshToken
            });
        }
    } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    } finally {
        await clearTokens();
    }
};

// ===== USER MANAGEMENT METHODS =====

/**
 * Get current user profile (requires auth)
 * @returns {Promise<Object>} User profile
 */
export const getCurrentUser = async () => {
    try {
        const response = await axiosInstance.get('/users/me');
        return response.data;
    } catch (error) {
        // Error: 'Get current user error:', error
        throw error;
    }
};

/**
 * Update user profile (requires auth)
 * @param {Object} updateData - Profile update data
 * @returns {Promise<Object>} Updated user data
 */
export const updateUser = async (updateData) => {
    try {
        const response = await axiosInstance.put('/users/me', updateData);
        return response.data;
    } catch (error) {
        // Error: 'Update user error:', error
        throw error;
    }
};

/**
 * Update user by ID (admin/privileged operation)
 * @param {string} id - User ID
 * @param {Object} updateData - Update data matching Swagger spec
 * @returns {Promise<Object>} Updated user data
 */
export const updateUserById = async (id, updateData) => {
    try {
        const response = await axiosInstance.patch(`/users/${id}/update`, updateData);
        return response.data;
    } catch (error) {
        // Error: 'Update user by ID error:', error
        throw error;
    }
};

/**
 * Upload user avatar (requires auth)
 * @param {File} file - Avatar image file
 * @returns {Promise<Object>} Upload response
 */
export const uploadUserAvatar = async (file) => {
    try {
        const formData = new FormData();
        formData.append('avatar', file);
        
        const response = await axiosInstance.post('/users/avatar', formData, {
            headers: {
                'Content-Type': 'multipart/form-data',
            }
        });
        return response.data;
    } catch (error) {
        // Error: 'Avatar upload error:', error
        throw error;
    }
};

/**
 * Get user notifications (requires auth)
 * @returns {Promise<Object>} User notifications
 */
export const getUserNotifications = async () => {
    try {
        const response = await axiosInstance.get('/users/notifications');
        return response.data;
    } catch (error) {
        // Error: 'Get notifications error:', error
        throw error;
    }
};

/**
 * Update notification preferences (requires auth)
 * @param {Object} preferences - Notification preferences
 * @returns {Promise<Object>} Updated preferences
 */
export const updateNotificationPreferences = async (preferences) => {
    try {
        const response = await axiosInstance.put('/users/notification-preferences', preferences);
        return response.data;
    } catch (error) {
        // Error: 'Update notification preferences error:', error
        throw error;
    }
};

/**
 * Delete user account (requires auth)
 * @param {string} reason - Deletion reason
 * @returns {Promise<Object>} Deletion response
 */
export const deleteUserAccount = async (reason = '') => {
    try {
        const response = await axiosInstance.delete('/users/me', {
            data: { reason }
        });
        
        // Clear tokens after successful deletion
        await clearTokens();
        
        return response.data;
    } catch (error) {
        // Error: 'Delete account error:', error
        throw error;
    }
};

// ===== INITIALIZATION =====

// Note: Token injection removed - axiosInstance now uses direct token access

/**
 * Initialize authentication from localStorage
 * @returns {Object|null} Current user data or null
 */
export const initializeAuth = () => {
    const token = getAccessToken();
    
    if (token && isTokenValid(token)) {
        return getUserIdFromToken(token);
    }
    
    // Clear invalid tokens
    clearTokens();
    return null;
};

/**
 * Initialize authentication from legacy localStorage format
 * This function should only be used during initial migration from old format
 */
export const initFromLocalStorage = () => {
    if (typeof window !== 'undefined') {
        // Check for token in old format
        const legacyToken = localStorage.getItem('jwtToken');
        if (legacyToken) {
            // Only set if token is still valid
            if (isTokenValid(legacyToken)) {
                setAccessToken(legacyToken);
            }

            // Schedule removal from localStorage regardless
            setTimeout(() => {
                localStorage.removeItem('jwtToken');
            }, 1000);
        }
    }
};

/**
 * Get base user data (public endpoint, limited user information)
 * Based on Swagger spec: GET /api/users/{id}/base
 * Returns: { user: BaseUser }
 * @param {string} id - User ID
 * @returns {Promise<Object>} Base user data response
 */
export const getBaseUserById = async (id) => {
    try {
        const response = await publicAxios.get(`/users/${id}/base`);
        // According to Swagger spec, response should be: { user: BaseUser }
        return response.data;
    } catch (error) {
        // Error: 'Error fetching base user data:', error
        throw error;
    }
};

/**
 * List all users, optionally filtered by userIds (public endpoint)
 * Based on Swagger spec: GET /api/users
 * @param {string[]} [userIds] - Optional array of user IDs to filter by
 * @returns {Promise<Object>} List users response
 */
export const listUsers = async (userIds = []) => {
    try {

        const params = {};
        if (userIds && userIds.length > 0) {
            params.userIds = userIds;
        }

        // Use axiosInstance for authenticated requests instead of publicAxios
        const response = await axiosInstance.get('/users', {params});

        return response.data;
    } catch (error) {
        // Error listing users logged for debugging
        throw error;
    }
};

/**
 * Create a new user
 * Based on Swagger spec: POST /api/users
 * @param {Object} userData - User data for creation
 * @returns {Promise<Object>} Create user response
 */
export const createUser = async (userData) => {
    try {
        const response = await axiosInstance.post('/users', userData);
        return response.data;
    } catch (error) {
        // Error: 'Error creating user:', error
        throw error;
    }
};

/**
 * Get user by ID (full user information, requires auth)
 * Based on Swagger spec: GET /api/users/{id}
 * @param {string} id - User ID
 * @returns {Promise<Object>} Get user response
 */
export const getUserById = async (id) => {
    try {
        const response = await axiosInstance.get(`/users/${id}`);
        return response.data;
    } catch (error) {
        // Error: 'Error fetching user data:', error
        throw error;
    }
};

/**
 * Update user's username
 * Based on Swagger spec: PATCH /api/users/{id}/name
 * @param {string} id - User ID
 * @param {string} userName - New username
 * @returns {Promise<Object>} Rename user response
 */
export const renameUser = async (id, userName) => {
    try {
        const response = await axiosInstance.patch(`/users/${id}/name`, {userName});
        return response.data;
    } catch (error) {
        // Error: 'Error renaming user:', error
        throw error;
    }
};

/**
 * Enable a user account
 * Based on Swagger spec: PATCH /api/users/{id}/enable
 * @param {string} id - User ID
 * @param {string} [verificationToken] - Optional verification token
 * @returns {Promise<Object>} Enable user response
 */
export const enableUser = async (id, verificationToken = '') => {
    try {
        const response = await axiosInstance.patch(`/users/${id}/enable`, {verificationToken});
        return response.data;
    } catch (error) {
        // Error: 'Error enabling user:', error
        throw error;
    }
};

/**
 * Disable a user account
 * Based on Swagger spec: PATCH /api/users/{id}/disable
 * @param {string} id - User ID
 * @returns {Promise<Object>} Disable user response
 */
export const disableUser = async (id) => {
    try {
        const response = await axiosInstance.patch(`/users/${id}/disable`, {});
        return response.data;
    } catch (error) {
        // Error: 'Error disabling user:', error
        throw error;
    }
};

/**
 * Get list of participating users
 * Based on Swagger spec: GET /api/users/participating
 * @returns {Promise<Object>} List participating users response
 */
export const listParticipatingUsers = async () => {
    try {
        const response = await axiosInstance.get('/users/enabled');
        return response.data;
    } catch (error) {
        // Error: 'Error listing participating users:', error
        throw error;
    }
};

// Alias for the renamed function
export const listEnabledUsers = listParticipatingUsers;

/**
 * Refresh auth token
 * Based on Swagger spec: POST /api/users/refresh-token
 * @param {string} refreshToken - Refresh token
 * @param {string} userId - User ID
 * @returns {Promise<Object>} Refresh token response
 */
export const refreshAuthToken = async (refreshToken, userId) => {
    try {
        const response = await axiosInstance.post('/users/refresh-token', {
            refreshToken,
            userId
        });
        return response.data;
    } catch (error) {
        // Error: 'Error refreshing auth token:', error
        throw error;
    }
};

// Social Links API
export const getSocialLinks = async (id) => {
    try {
        const response = await fetch(`/api/users/${id}/social-links`);
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        
        return {
            links: data.links || {}
        };
    } catch (error) {
        // Error: '[getSocialLinks] API Error:', error
        // Return empty response instead of dummy data
        return {
            links: {}
        };
    }
};

// Review APIs
export const getReviewsByUserId = async (userId) => {
    try {
        const response = await fetch(`/api/users/${userId}/reviews`);
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        
        return {
            reviews: data.reviews || []
        };
    } catch (error) {
        // Error: '[getReviewsByUserId] API Error:', error
        // Return empty response instead of dummy data
        return {
            reviews: []
        };
    }
};

 