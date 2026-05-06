import axiosInstance from './axiosInstance';
const NOTIFICATIONS_ENDPOINT = '/notifications';
/**
 * Safe URI component encoding
 * @param {string} component - Component to encode
 * @returns {string} Safely encoded component
 */
const safeEncode = (component) => {
    if (component === null || typeof component === 'undefined') {
        return '';
    }
    return encodeURIComponent(String(component));
};
/**
 * Handle notifications API errors with detailed logging and user-friendly messages
 * @param {Error} error - The error object from axios
 * @param {string} endpoint - The API endpoint that failed
 * @param {string} operation - The operation being performed
 * @returns {Object} Standardized error response
 */
const handleNotificationsError = (error, endpoint, operation) => {
    const errorDetails = {
        success: false,
        operation,
        endpoint,
        timestamp: new Date().toISOString()
    };
    if (error.response) {
        // Server responded with error status
        const { status, data } = error.response;
        errorDetails.status = status;
        errorDetails.message = data?.message || `Server error (${status})`;
        errorDetails.data = data;
        // Handle specific HTTP errors
        switch (status) {
            case 404:
                errorDetails.userMessage = 'Notifications not found';
                errorDetails.severity = 'warning';
                break;
            case 500:
                errorDetails.userMessage = 'Server error - notifications temporarily unavailable';
                errorDetails.severity = 'error';
                break;
            case 403:
                errorDetails.userMessage = 'Access denied to notifications';
                errorDetails.severity = 'error';
                break;
            case 401:
                errorDetails.userMessage = 'Authentication required to view notifications';
                errorDetails.severity = 'error';
                break;
            default:
                errorDetails.userMessage = 'Unable to load notifications';
                errorDetails.severity = 'error';
        }
        // Log only in development with safe data serialization
        if (process.env.NODE_ENV === 'development') {
            const safeData = data ? (typeof data === 'object' ? JSON.stringify(data) : String(data)) : 'No response data';
        }
    } else if (error.request) {
        // Network error
        errorDetails.userMessage = 'Network error - please check your connection';
        errorDetails.severity = 'error';
        errorDetails.network = true;
        if (process.env.NODE_ENV === 'development') {
        }
    } else {
        // Other error
        errorDetails.userMessage = 'Unexpected error occurred';
        errorDetails.severity = 'error';
        errorDetails.message = error.message;
        if (process.env.NODE_ENV === 'development') {
        }
    }
    return errorDetails;
};
/**
 * LIST ALL ALERTS FOR A USER
 * GET /api/notifications/alerts/{userId}
 * @param {string} userId - User ID to get alerts for
 * @param {Object} filters - Optional filters { type, isRead }
 * @returns {Object} Success/error response with alerts array
 */
export const listAlerts = async (userId, filters = {}) => {
    const endpoint = `${NOTIFICATIONS_ENDPOINT}/ListAlerts`;
    try {
        const response = await axiosInstance.post(endpoint, {
            type: filters.type || '',
            is_read: filters.isRead
        });
        return { success: true, ...response.data };
    } catch (error) {
        const errorResponse = handleNotificationsError(error, endpoint, 'listAlerts');
        // For 404 or network errors, return empty array
        if (error.response?.status === 404 || !error.response) {
            return { success: true, alerts: [] };
        }
        return errorResponse;
    }
};
/**
 * GET ALERTS BY TYPE FOR A USER
 * GET /api/notifications/{userId}/alerts/type/{type}
 * @param {string} userId - User ID to get alerts for
 * @param {string} type - Alert type to filter by
 * @returns {Object} Success/error response with alerts array
 */
export const getAlertsByType = async (userId, type) => {
    if (!type) {
        throw new Error('type is required for getAlertsByType.');
    }
    const endpoint = `${NOTIFICATIONS_ENDPOINT}/GetAlertsByType`;
    try {
        const response = await axiosInstance.post(endpoint, {
            type: type,
            is_read: false
        });
        return { success: true, ...response.data };
    } catch (error) {
        const errorResponse = handleNotificationsError(error, endpoint, 'getAlertsByType');
        // For 404 or network errors, return empty array
        if (error.response?.status === 404 || !error.response) {
            return { success: true, alerts: [] };
        }
        return errorResponse;
    }
};
/**
 * Utility function to check if a notifications API response was successful
 * @param {Object} response - Response from any notifications API function
 * @returns {boolean} True if successful, false otherwise
 */
export const isNotificationsResponseSuccess = (response) => {
    return response && response.success === true;
};
/**
 * Get user-friendly error message from notifications API response
 * @param {Object} response - Error response from notifications API
 * @returns {string|null} User-friendly error message
 */
export const getNotificationsErrorMessage = (response) => {
    if (isNotificationsResponseSuccess(response)) {
        return null;
    }
    return response?.userMessage || 'Failed to load notifications';
};
/**
 * Helper function to get notification types with user-friendly labels
 * @returns {Array} Array of notification types
 */
export const getNotificationTypes = () => [
    { value: 'message', label: 'Messages', icon: 'MessageCircle' },
    { value: 'transaction', label: 'Transactions', icon: 'CreditCard' },
    { value: 'product', label: 'Products', icon: 'Package' },
    { value: 'store', label: 'Store Updates', icon: 'Store' },
    { value: 'system', label: 'System Alerts', icon: 'Settings' },
    { value: 'promotion', label: 'Promotions', icon: 'Tag' },
    { value: 'security', label: 'Security', icon: 'Shield' },
    { value: 'reminder', label: 'Reminders', icon: 'Clock' }
]; 