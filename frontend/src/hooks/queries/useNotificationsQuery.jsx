import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import axiosInstance from '../../api/axiosInstance';

// Query Keys
const NOTIFICATION_KEYS = {
  all: ['notifications'],
  lists: () => [...NOTIFICATION_KEYS.all, 'list'],
  list: (filters) => [...NOTIFICATION_KEYS.all, 'list', filters],
  unreadCount: () => [...NOTIFICATION_KEYS.all, 'unreadCount'],
  detail: (id) => [...NOTIFICATION_KEYS.all, 'detail', id],
};

// API functions
const notificationsApi = {
  getNotifications: async (params = {}) => {
    const { data } = await axiosInstance.get('/notifications', { params });
    return data;
  },
  
  getUnreadCount: async () => {
    const { data } = await axiosInstance.get('/notifications/unread-count');
    return data;
  },
  
  markAsRead: async (notificationId) => {
    const { data } = await axiosInstance.put(`/notifications/${notificationId}/read`);
    return data;
  },
  
  markAllAsRead: async () => {
    const { data } = await axiosInstance.put('/notifications/read-all');
    return data;
  },
  
  deleteNotification: async (notificationId) => {
    const { data } = await axiosInstance.delete(`/notifications/${notificationId}`);
    return data;
  },
  
  updatePreferences: async (preferences) => {
    const { data } = await axiosInstance.put('/notifications/preferences', preferences);
    return data;
  },
};

/**
 * Get notifications list
 */
export function useNotifications(filters = {}, options = {}) {
  return useQuery({
    queryKey: NOTIFICATION_KEYS.list(filters),
    queryFn: () => notificationsApi.getNotifications(filters),
    ...options,
  });
}

/**
 * Get unread notifications count
 */
export function useUnreadNotificationsCount(options = {}) {
  return useQuery({
    queryKey: NOTIFICATION_KEYS.unreadCount(),
    queryFn: notificationsApi.getUnreadCount,
    refetchInterval: 30000, // Refetch every 30 seconds
    ...options,
  });
}

/**
 * Mark notification as read
 */
export function useMarkNotificationAsRead() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: notificationsApi.markAsRead,
    onSuccess: (_, notificationId) => {
      queryClient.invalidateQueries({ queryKey: NOTIFICATION_KEYS.lists() });
      queryClient.invalidateQueries({ queryKey: NOTIFICATION_KEYS.unreadCount() });
      queryClient.invalidateQueries({ queryKey: NOTIFICATION_KEYS.detail(notificationId) });
    },
  });
}

/**
 * Mark all notifications as read
 */
export function useMarkAllNotificationsAsRead() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: notificationsApi.markAllAsRead,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: NOTIFICATION_KEYS.all });
    },
  });
}

/**
 * Delete notification
 */
export function useDeleteNotification() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: notificationsApi.deleteNotification,
    onSuccess: (_, notificationId) => {
      queryClient.invalidateQueries({ queryKey: NOTIFICATION_KEYS.lists() });
      queryClient.invalidateQueries({ queryKey: NOTIFICATION_KEYS.unreadCount() });
      queryClient.removeQueries({ queryKey: NOTIFICATION_KEYS.detail(notificationId) });
    },
  });
}

/**
 * Update notification preferences
 */
export function useUpdateNotificationPreferences() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: notificationsApi.updatePreferences,
    onSuccess: () => {
      // Optionally invalidate related queries
      queryClient.invalidateQueries({ queryKey: ['user', 'preferences'] });
    },
  });
}

/**
 * Get notification statistics
 */
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

// Export the API for direct use if needed
export { notificationsApi };