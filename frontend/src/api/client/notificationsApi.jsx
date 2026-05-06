// src/api/notificationsApi.js

import axiosInstance from "../axiosInstance";

// Base path from swagger => all endpoints start with: /api/notifications
const NOTIFICATIONS_BASE_URL = '/notifications';

/**
 * List All alerts for a user
 * GET /api/notifications/alerts
 * Uses auth token for user identification
 */
export const listAlerts = async (queryParams = {}) => {
    const response = await axiosInstance.get(`${NOTIFICATIONS_BASE_URL}/alerts`, {
        params: queryParams, // type, isRead, limit, offset
    });
    return response.data; // shape => notificationspbListAlertsResponse
};

/**
 * Retrieve Alerts by Type for a User
 * GET /api/notifications/alerts/type
 * Uses auth token for user identification
 */
export const getAlertsByType = async (queryParams = {}) => {
    const response = await axiosInstance.get(`${NOTIFICATIONS_BASE_URL}/alerts/type`, {
        params: queryParams, // type (required), isRead, limit, offset
    });
    return response.data; // shape => notificationspbGetAlertsByTypeResponse
};

/**
 * Mark a specific alert as read
 * PUT /api/notifications/alerts/{alertId}/read
 */
export const markAlertAsRead = async (alertId) => {
    const response = await axiosInstance.put(`${NOTIFICATIONS_BASE_URL}/alerts/${alertId}/read`);
    return response.data;
};

/**
 * Mark all alerts as read
 * PUT /api/notifications/alerts/read-all
 * Optional: type parameter to mark only specific type as read
 */
export const markAllAlertsAsRead = async (type = null) => {
    const body = type ? { type } : {};
    const response = await axiosInstance.put(`${NOTIFICATIONS_BASE_URL}/alerts/read-all`, body);
    return response.data;
};

/**
 * Delete an alert (remove from user's view)
 * DELETE /api/notifications/alerts/{alertId}
 */
export const deleteAlert = async (alertId) => {
    const response = await axiosInstance.delete(`${NOTIFICATIONS_BASE_URL}/alerts/${alertId}`);
    return response.data;
};

/**
 * Get unread alerts count
 * GET /api/notifications/alerts/unread-count
 */
export const getUnreadCount = async () => {
    const response = await axiosInstance.get(`${NOTIFICATIONS_BASE_URL}/alerts/unread-count`);
    return response.data;
};

/**
 * Get notification statistics
 * GET /api/notifications/alerts/stats
 */
export const getNotificationStats = async () => {
    const response = await axiosInstance.get(`${NOTIFICATIONS_BASE_URL}/alerts/stats`);
    return response.data;
};

// Legacy API functions (for backward compatibility)
export const createNotification = async (createNotificationData) => {
    const response = await axiosInstance.post(`${NOTIFICATIONS_BASE_URL}`, createNotificationData);
    return response.data;
};

export const getNotification = async (id, userId) => {
    const response = await axiosInstance.get(`${NOTIFICATIONS_BASE_URL}/${id}`, {
        params: {userId},
    });
    return response.data;
};

export const deleteNotification = async (id) => {
    const response = await axiosInstance.delete(`${NOTIFICATIONS_BASE_URL}/${id}`);
    return response.data;
};

export const updateNotification = async (id, updateBody) => {
    const response = await axiosInstance.patch(
        `${NOTIFICATIONS_BASE_URL}/${id}/update`,
        updateBody
    );
    return response.data;
};

export const getNotificationByUser = async (userId) => {
    const response = await axiosInstance.get(`${NOTIFICATIONS_BASE_URL}/user/${userId}`);
    return response.data;
};

export const updatePreferences = async (userId, preferencesBody) => {
    const response = await axiosInstance.patch(
        `${NOTIFICATIONS_BASE_URL}/preferences/${userId}`,
        preferencesBody
    );
    return response.data;
};

// Notification types enum
export const NotificationTypes = {
    MESSAGE: 'message',
    COMMENT: 'comment',
    OFFER: 'offer',
    ORDER: 'order',
    PAYMENT: 'payment',
    REVIEW: 'review',
    FOLLOWING: 'following',
    PRODUCT: 'product',
    WISHLIST: 'wishlist',
    SUPPORT: 'support',
    BASKET: 'basket',
    INTERACTION: 'interaction',
    SYSTEM: 'system',
};

// Get notification type configuration
export const getNotificationTypes = () => [
    { value: NotificationTypes.MESSAGE, label: 'Messages', icon: 'MessageCircle' },
    { value: NotificationTypes.COMMENT, label: 'Comments', icon: 'MessageCircle' },
    { value: NotificationTypes.OFFER, label: 'Offers', icon: 'Tag' },
    { value: NotificationTypes.ORDER, label: 'Orders', icon: 'Package' },
    { value: NotificationTypes.PAYMENT, label: 'Payments', icon: 'CreditCard' },
    { value: NotificationTypes.REVIEW, label: 'Reviews', icon: 'Star' },
    { value: NotificationTypes.FOLLOWING, label: 'Following', icon: 'Users' },
    { value: NotificationTypes.PRODUCT, label: 'Products', icon: 'Package' },
    { value: NotificationTypes.WISHLIST, label: 'Wishlist', icon: 'Heart' },
    { value: NotificationTypes.SUPPORT, label: 'Support', icon: 'HelpCircle' },
    { value: NotificationTypes.BASKET, label: 'Basket', icon: 'ShoppingCart' },
    { value: NotificationTypes.INTERACTION, label: 'Interactions', icon: 'Activity' },
    { value: NotificationTypes.SYSTEM, label: 'System', icon: 'Settings' },
];

// Helper to check if notifications response is successful
export const isNotificationsResponseSuccess = (response) => {
    return response && !response.error && response.alerts !== undefined;
};

// Helper to get error message from notifications response
export const getNotificationsErrorMessage = (response) => {
    if (!response) return 'No response from server';
    if (response.error) return response.error;
    if (response.message) return response.message;
    return 'Unknown error occurred';
};
