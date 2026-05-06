// src/api/userApi.jsx - PUBLIC USER API (No Authentication Required)
import axios from "axios";
const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL;
// Create a simple axios instance for public endpoints
const publicAxios = axios.create({
    baseURL: apiBaseUrl,
    timeout: 10000,
    headers: {
        'Content-Type': 'application/json',
    },
});
/**
 * Get base user information by ID (public endpoint)
 * @param {string} id - User ID
 * @returns {Promise<Object>} Base user data
 */
export const getBaseUserById = async (id) => {
    try {
        const response = await publicAxios.get(`/users/${id}/base`);
        return response.data;
    } catch (error) {
        throw error;
    }
};
/**
 * Get user's public social links (no auth required)
 * @param {string} id - User ID
 * @returns {Promise<Object>} User social links
 */
export const getUserSocialLinks = async (id) => {
    try {
        const response = await publicAxios.get(`/users/${id}/social-links`);
        return response.data;
    } catch (error) {
        throw error;
    }
};
/**
 * Get public reviews for a user (no auth required)
 * @param {string} userId - User ID
 * @returns {Promise<Object>} User reviews
 */
export const getReviewsByUserId = async (userId) => {
    try {
        const response = await publicAxios.get(`/users/${userId}/reviews`);
        return response.data;
    } catch (error) {
        throw error;
    }
};
/**
 * Get public user profile (no auth required)
 * @param {string} userId - User ID
 * @returns {Promise<Object>} Public user profile
 */
export const getPublicUserProfile = async (userId) => {
    try {
        const response = await publicAxios.get(`/users/${userId}/public`);
        return response.data;
    } catch (error) {
        throw error;
    }
};
/**
 * Search users publicly (no auth required)
 * @param {string} query - Search query
 * @param {Object} filters - Search filters
 * @returns {Promise<Object>} Search results
 */
export const searchUsers = async (query, filters = {}) => {
    try {
        const response = await publicAxios.get('/users/search', {
            params: { query, ...filters }
        });
        return response.data;
    } catch (error) {
        throw error;
    }
};
/**
 * Get user's public statistics (no auth required)
 * @param {string} userId - User ID
 * @returns {Promise<Object>} User statistics
 */
export const getUserPublicStats = async (userId) => {
    try {
        const response = await publicAxios.get(`/users/${userId}/stats`);
        return response.data;
    } catch (error) {
        throw error;
    }
};
// ===== PUBLIC AUTHENTICATION ENDPOINTS =====
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
        throw error;
    }
};
/**
 * Login user with email/password (public endpoint)
 * @param {Object} credentials - Login credentials
 * @returns {Promise<Object>} Login response
 */
export const loginUser = async (credentials) => {
    try {
        const response = await publicAxios.post('/users/login', credentials);
        return response.data;
    } catch (error) {
        throw error;
    }
};
/**
 * Login with Google (public endpoint)
 * @param {string} idToken - Google ID token
 * @returns {Promise<Object>} Login response
 */
export const loginWithGoogle = async (idToken) => {
    try {
        const response = await publicAxios.post('/users/google-login', { idToken });
        return response.data;
    } catch (error) {
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
        throw error;
    }
};