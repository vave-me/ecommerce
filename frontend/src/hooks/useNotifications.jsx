import { useState, useEffect, useCallback, useRef, useMemo } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import axiosInstance from '../api/axiosInstance';
import {
    listAlerts,
    getAlertsByType,
    markAlertAsRead,
    markAllAlertsAsRead,
    deleteAlert,
    isNotificationsResponseSuccess,
    getNotificationsErrorMessage
} from '../api/client/notificationsApi';
// Query Keys
const NOTIFICATION_KEYS = {
  all: ['notifications'],
  lists: () => [...NOTIFICATION_KEYS.all, 'list'],
  list: (filters) => [...NOTIFICATION_KEYS.all, 'list', filters],
  unreadCount: () => [...NOTIFICATION_KEYS.all, 'unreadCount'],
  detail: (id) => [...NOTIFICATION_KEYS.all, 'detail', id],
};

// Modern API functions using axiosInstance
const notificationsApi = {
  getNotifications: async (params = {}) => {
    const { data } = await axiosInstance.get('/notifications/alerts', { params });
    return data;
  },
  
  markAsRead: async (notificationId) => {
    return await markAlertAsRead(notificationId);
  },
  
  markAllAsRead: async (type = null) => {
    return await markAllAlertsAsRead(type);
  },
  
  deleteNotification: async (notificationId) => {
    return await deleteAlert(notificationId);
  },
};

/**
 * Unified notifications hook with React Query support and backward compatibility
 * Supports both modern React Query pattern and legacy API
 * @param {string} userId - User ID to fetch notifications for
 * @param {Object} options - Configuration options
 * @returns {Object} Notifications state and management functions
 */
export const useNotifications = (userId, options = {}) => {
    const {
        autoRefresh = false,
        refreshInterval = 30000, // 30 seconds
        initialFilters = {},
        useReactQuery = true // New option to use React Query
    } = options;
    
    const queryClient = useQueryClient();
    
    // React Query hooks
    const { 
        data: queryData, 
        isLoading: queryLoading, 
        error: queryError, 
        refetch: queryRefetch 
    } = useQuery({
        queryKey: NOTIFICATION_KEYS.list(initialFilters),
        queryFn: async () => {
            // Try modern API first, fallback to legacy
            try {
                return await notificationsApi.getNotifications(initialFilters);
            } catch (error) {
                // Fallback to legacy API
                const response = await listAlerts(initialFilters);
                if (!isNotificationsResponseSuccess(response)) {
                    throw new Error(getNotificationsErrorMessage(response));
                }
                return response;
            }
        },
        enabled: useReactQuery && !!userId,
        refetchInterval: autoRefresh ? refreshInterval : false,
    });
    
    // Mutations
    const markAsReadMutation = useMutation({
        mutationFn: notificationsApi.markAsRead,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: NOTIFICATION_KEYS.lists() });
            queryClient.invalidateQueries({ queryKey: NOTIFICATION_KEYS.unreadCount() });
        },
    });
    
    const markAllAsReadMutation = useMutation({
        mutationFn: notificationsApi.markAllAsRead,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: NOTIFICATION_KEYS.all });
        },
    });
    
    const deleteNotificationMutation = useMutation({
        mutationFn: notificationsApi.deleteNotification,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: NOTIFICATION_KEYS.lists() });
            queryClient.invalidateQueries({ queryKey: NOTIFICATION_KEYS.unreadCount() });
        },
    });
    
    // Legacy state management (for backward compatibility)
    const [notifications, setNotifications] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [filters, setFilters] = useState(initialFilters);
    const [unreadCount, setUnreadCount] = useState(0);
    // Refs for cleanup
    const isMountedRef = useRef(true);
    const refreshTimeoutRef = useRef(null);
    // Sync React Query data with legacy state
    useEffect(() => {
        if (useReactQuery && queryData) {
            const notificationsArray = Array.isArray(queryData) 
                ? queryData 
                : (queryData.alerts || queryData.notifications || []);
            setNotifications(notificationsArray);
            setLoading(queryLoading);
            setError(queryError);
        }
    }, [useReactQuery, queryData, queryLoading, queryError]);
    
    // Calculate unread count whenever notifications change
    useEffect(() => {
        const count = notifications.filter(notification => !notification.isRead).length;
        setUnreadCount(count);
    }, [notifications]);
    // Cleanup on unmount
    useEffect(() => {
        return () => {
            isMountedRef.current = false;
            if (refreshTimeoutRef.current) {
                clearTimeout(refreshTimeoutRef.current);
            }
        };
    }, []);
    /**
     * Fetch notifications with current filters
     */
    const fetchNotifications = useCallback(async (silent = false) => {
        if (!userId) {
            setNotifications([]);
            setError(null);
            return;
        }
        if (!silent) {
            setLoading(true);
        }
        setError(null);
        try {
            const response = await listAlerts(filters);
            if (!isMountedRef.current) return;
            if (!isNotificationsResponseSuccess(response)) {
                const errorMessage = getNotificationsErrorMessage(response);
                setError(new Error(errorMessage));
                // Only log non-warning errors
                if (response?.severity !== 'warning') {
                }
                setNotifications([]);
                return;
            }
            const alertsArray = response.alerts || [];
            setNotifications(alertsArray);
        } catch (err) {
            if (isMountedRef.current) {
                setError(err);
                setNotifications([]);
            }
        } finally {
            if (isMountedRef.current && !silent) {
                setLoading(false);
            }
        }
    }, [userId, filters]);
    /**
     * Fetch notifications by specific type
     */
    const fetchNotificationsByType = useCallback(async (type, silent = false) => {
        if (!userId || !type) {
            return [];
        }
        if (!silent) {
            setLoading(true);
        }
        try {
            const response = await getAlertsByType({ type });
            if (!isNotificationsResponseSuccess(response)) {
                const errorMessage = getNotificationsErrorMessage(response);
                if (response?.severity !== 'warning') {
                }
                return [];
            }
            return response.alerts || [];
        } catch (err) {
            return [];
        } finally {
            if (!silent) {
                setLoading(false);
            }
        }
    }, [userId]);
    /**
     * Update filters and refetch notifications
     */
    const updateFilters = useCallback((newFilters) => {
        setFilters(prev => ({ ...prev, ...newFilters }));
    }, []);
    /**
     * Clear all filters
     */
    const clearFilters = useCallback(() => {
        setFilters({});
    }, []);
    /**
     * Mark notification as read (optimistic update)
     */
    const markAsRead = useCallback((notificationId) => {
        if (useReactQuery) {
            markAsReadMutation.mutate(notificationId);
        } else {
            setNotifications(prev => 
                prev.map(notification => 
                    notification.id === notificationId 
                        ? { ...notification, isRead: true }
                        : notification
                )
            );
        }
    }, [useReactQuery, markAsReadMutation]);
    
    /**
     * Mark all notifications as read (optimistic update)
     */
    const markAllAsRead = useCallback(() => {
        if (useReactQuery) {
            markAllAsReadMutation.mutate();
        } else {
            setNotifications(prev => 
                prev.map(notification => ({ ...notification, isRead: true }))
            );
        }
    }, [useReactQuery, markAllAsReadMutation]);
    
    /**
     * Remove notification from list (optimistic update)
     */
    const removeNotification = useCallback((notificationId) => {
        if (useReactQuery) {
            deleteNotificationMutation.mutate(notificationId);
        } else {
            setNotifications(prev => 
                prev.filter(notification => notification.id !== notificationId)
            );
        }
    }, [useReactQuery, deleteNotificationMutation]);
    /**
     * Refresh notifications
     */
    const refresh = useCallback(() => {
        if (useReactQuery) {
            queryRefetch();
        } else {
            fetchNotifications(false);
        }
    }, [useReactQuery, queryRefetch, fetchNotifications]);
    /**
     * Setup auto-refresh if enabled
     */
    useEffect(() => {
        if (autoRefresh && userId) {
            const scheduleRefresh = () => {
                refreshTimeoutRef.current = setTimeout(() => {
                    if (isMountedRef.current) {
                        fetchNotifications(true); // Silent refresh
                        scheduleRefresh(); // Schedule next refresh
                    }
                }, refreshInterval);
            };
            scheduleRefresh();
            return () => {
                if (refreshTimeoutRef.current) {
                    clearTimeout(refreshTimeoutRef.current);
                }
            };
        }
    }, [autoRefresh, refreshInterval, userId, fetchNotifications]);
    // Initial fetch when userId or filters change (only for legacy mode)
    useEffect(() => {
        if (!useReactQuery && userId) {
            fetchNotifications();
        }
    }, [useReactQuery, userId, filters, fetchNotifications]);
    // Derived state
    const hasNotifications = notifications.length > 0;
    const hasUnread = unreadCount > 0;
    const isEmpty = !hasNotifications && !loading && !error;
    return {
        // Data
        notifications,
        unreadCount,
        // State
        loading,
        error,
        filters,
        // Derived state
        hasNotifications,
        hasUnread,
        isEmpty,
        // Actions
        fetchNotifications,
        fetchNotificationsByType,
        updateFilters,
        clearFilters,
        markAsRead,
        markAllAsRead,
        removeNotification,
        refresh,
        // React Query specific exports
        isUsingReactQuery: useReactQuery,
        queryClient: useReactQuery ? queryClient : null
    };
};

// Export React Query hooks for direct use
export {
    NOTIFICATION_KEYS,
    notificationsApi
};

// React Query specific hooks
export function useUnreadNotificationsCount(options = {}) {
    return useQuery({
        queryKey: NOTIFICATION_KEYS.unreadCount(),
        queryFn: async () => {
            const { data } = await axiosInstance.get('/notifications/unread-count');
            return data;
        },
        refetchInterval: 30000, // Refetch every 30 seconds
        ...options,
    });
}

export function useNotificationStats(options = {}) {
    return useQuery({
        queryKey: ['notifications', 'stats'],
        queryFn: async () => {
            const { data } = await axiosInstance.get('/notifications/stats');
            return data;
        },
        refetchInterval: 60000, // Refetch every minute
        ...options,
    });
}