// src/api/client/metricsApi.jsx - Metrics API endpoints
import axiosInstance from '../axiosInstance';

/**
 * Get user metrics
 * @param {string} userId - User ID
 * @returns {Promise<Object>} User metrics
 */
export const getUserMetrics = async (userId) => {
  try {
    const response = await axiosInstance.get(`/metrics/user/${userId}`);
    return response.data;
  } catch (error) {
    // Error: 'Failed to fetch user metrics:', error...
    throw error;
  }
};

/**
 * Get item metrics
 * @param {string} itemId - Item ID
 * @returns {Promise<Object>} Item metrics
 */
export const getItemMetrics = async (itemId) => {
  try {
    const response = await axiosInstance.get(`/metrics/item/${itemId}`);
    return response.data;
  } catch (error) {
    // Error: 'Failed to fetch item metrics:', error...
    throw error;
  }
};

/**
 * Get multiple items metrics
 * @param {Array<string>} itemIds - Array of item IDs
 * @returns {Promise<Object>} Items metrics
 */
export const getItemsMetrics = async (itemIds) => {
  try {
    const response = await axiosInstance.post('/metrics/items', { itemIds });
    return response.data;
  } catch (error) {
    // Error: 'Failed to fetch items metrics:', error...
    throw error;
  }
};

/**
 * Update user metric
 * @param {string} userId - User ID
 * @param {Object} metric - Metric data
 * @returns {Promise<Object>} Updated metric
 */
export const updateUserMetric = async (userId, metric) => {
  try {
    const response = await axiosInstance.post(`/metrics/user/${userId}`, metric);
    return response.data;
  } catch (error) {
    // Error: 'Failed to update user metric:', error...
    throw error;
  }
};