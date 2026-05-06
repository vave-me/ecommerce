"use client"
/**
 * Admin User API
 * This file re-exports all authentication functions from the main userApi
 * to ensure admins and regular users use the same authentication flow
 */

// Re-export all authentication and token management functions from main userApi
export {
    // Token management
    setAccessToken,
    getAccessToken,
    isTokenValid,
    getUserIdFromToken,
    getRefreshToken,
    setRefreshToken,
    clearTokens,
    refreshAccessToken,
    initializeAuth,
    initFromLocalStorage,
    
    // Authentication methods
    registerUser,
    loginUser,
    loginWithGoogle,
    loginWithGoogleMobile,
    forgotPassword,
    resetPassword,
    logoutUser,
    
    // User management
    getCurrentUser,
    updateUser,
    uploadUserAvatar,
    getUserNotifications,
    updateNotificationPreferences,
    deleteUserAccount,
    
    // User queries
    getBaseUserById,
    listUsers,
    createUser,
    getUserById,
    updateUserById,
    renameUser,
    enableUser,
    disableUser,
    listParticipatingUsers,
    listEnabledUsers,
    refreshAuthToken,
    clearUserTokens,
    getSocialLinks,
    getReviewsByUserId
} from '../userApi';

// Import axiosInstance for admin-specific endpoints
import axiosInstance from '../../axiosInstance';

// ===== ADMIN-SPECIFIC ENDPOINTS =====

/**
 * Add a new admin user (admin only)
 * Based on Swagger spec: POST /api/users/add-admin
 * @param {Object} adminData - Admin user data
 * @returns {Promise<Object>} Create admin response
 */
export const addAdmin = async (adminData) => {
    try {
        const response = await axiosInstance.post('/users/admin/add', adminData);
        return response.data;
    } catch (error) {
        // Error: 'Error adding admin:', error...
        throw error;
    }
};

/**
 * Get all users with advanced filtering (admin only)
 * @param {Object} filters - Advanced filter options
 * @returns {Promise<Object>} Filtered users list
 */
export const getFilteredUsers = async (filters = {}) => {
    try {
        const response = await axiosInstance.get('/users', {
            params: filters
        });
        return response.data;
    } catch (error) {
        // Error: 'Error fetching filtered users:', error...
        throw error;
    }
};

/**
 * Bulk update users (admin only)
 * @param {Array<string>} userIds - Array of user IDs to update
 * @param {Object} updateData - Data to update for all users
 * @returns {Promise<Object>} Bulk update response
 */
export const bulkUpdateUsers = async (userIds, updateData) => {
    try {
        const response = await axiosInstance.patch('/users/bulk-update', {
            userIds,
            ...updateData
        });
        return response.data;
    } catch (error) {
        // Error: 'Error bulk updating users:', error...
        throw error;
    }
};

/**
 * Get user activity logs (admin only)
 * @param {string} userId - User ID
 * @param {Object} options - Query options (limit, offset, etc.)
 * @returns {Promise<Object>} User activity logs
 */
export const getUserActivityLogs = async (userId, options = {}) => {
    try {
        const response = await axiosInstance.get(`/users/${userId}/activity-logs`, {
            params: options
        });
        return response.data;
    } catch (error) {
        // Error: 'Error fetching user activity logs:', error...
        throw error;
    }
};

/**
 * Get system-wide user statistics (admin only)
 * @returns {Promise<Object>} User statistics
 */
export const getUserStatistics = async () => {
    try {
        const response = await axiosInstance.get('/users/statistics');
        return response.data;
    } catch (error) {
        // Error: 'Error fetching user statistics:', error...
        throw error;
    }
};

/**
 * Impersonate a user (admin only - use with extreme caution)
 * @param {string} userId - User ID to impersonate
 * @returns {Promise<Object>} Impersonation tokens
 */
export const impersonateUser = async (userId) => {
    try {
        const response = await axiosInstance.post(`/users/${userId}/impersonate`);
        // Store the impersonation tokens
        if (response.data.accessToken) {
            setAccessToken(response.data.accessToken);
        }
        if (response.data.refreshToken) {
            setRefreshToken(response.data.refreshToken);
        }
        return response.data;
    } catch (error) {
        // Error: 'Error impersonating user:', error...
        throw error;
    }
};

/**
 * Stop impersonation and return to admin account
 * @returns {Promise<Object>} Admin tokens
 */
export const stopImpersonation = async () => {
    try {
        const response = await axiosInstance.post('/users/stop-impersonation');
        // Restore admin tokens
        if (response.data.accessToken) {
            setAccessToken(response.data.accessToken);
        }
        if (response.data.refreshToken) {
            setRefreshToken(response.data.refreshToken);
        }
        return response.data;
    } catch (error) {
        // Error: 'Error stopping impersonation:', error...
        throw error;
    }
};