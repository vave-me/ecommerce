import { useState, useCallback, useEffect, useMemo } from 'react';
import { 
    getActivity,
    getInteractions,
    archiveActivity, 
    restoreActivity, 
    createActivity
} from '../api/client/activityApi';
/**
 * useActivity Hook
 * Manages activity state, loading, filtering, and CRUD operations
 * Follows the same pattern as useNotifications for consistency
 */
export const useActivity = (userId, options = {}) => {
    const {
        autoRefresh = false,
        refreshInterval = 30000,
        initialFilters = {}
    } = options;
    // Core state
    const [activities, setActivities] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [filters, setFilters] = useState(initialFilters);
    // Computed values
    const unreadCount = useMemo(() => 
        activities.filter(activity => !activity.isRead).length, 
        [activities]
    );
    const hasActivities = activities.length > 0;
    const isEmpty = !loading && !error && activities.length === 0;
    // Filter activities based on current filters
    const filteredActivities = useMemo(() => {
        return activities.filter(activity => {
            // Filter by type
            if (filters.type && activity.type !== filters.type) {
                return false;
            }
            // Filter by read status
            if (filters.isRead !== undefined && activity.isRead !== filters.isRead) {
                return false;
            }
            // Filter by action type
            if (filters.actionType && activity.actionType !== filters.actionType) {
                return false;
            }
            return true;
        });
    }, [activities, filters]);
    /**
     * Load activities from the API
     */
    const loadActivities = useCallback(async (silent = false) => {
        if (!userId) {
            setActivities([]);
            setLoading(false);
            return;
        }
        if (!silent) {
            setLoading(true);
        }
        setError(null);
        try {
            // First get the user's activity
            const activityResponse = await getActivity(userId);
            const activityId = activityResponse?.activityId;
            
            if (activityId) {
                // Then get interactions for that activity
                const response = await getInteractions(activityId);
                // Transform API response to match frontend expectations
                const interactionsData = response.interactions || [];
                const transformedActivities = interactionsData.map(interaction => ({
                    ...interaction,
                    // Map interaction fields to activity-like structure
                    id: interaction.id,
                    userId: userId,
                    isRead: interaction.isRead ?? false,
                    createdAt: interaction.createdAt || interaction.timestamp || new Date().toISOString(),
                    message: interaction.message || `${interaction.actionType || 'Action'} on ${interaction.itemType || 'item'}`,
                    actionType: interaction.actionType || 'unknown',
                    itemType: interaction.itemType || 'unknown',
                    itemId: interaction.itemId
                }));
                setActivities(transformedActivities);
            } else {
                // No activity found for user, set empty array
                setActivities([]);
            }
        } catch (err) {
            setError(err);
        } finally {
            setLoading(false);
        }
    }, [userId]);
    /**
     * Refresh activities (public method)
     */
    const refresh = useCallback(() => {
        loadActivities(false);
    }, [loadActivities]);
    /**
     * Mark activity as read/unread
     */
    const markAsRead = useCallback((activityId, isRead = true) => {
        setActivities(prev => 
            prev.map(activity => 
                activity.id === activityId 
                    ? { ...activity, isRead } 
                    : activity
            )
        );
        // TODO: Call API endpoint to mark as read on server if available
    }, []);
    /**
     * Mark all activities as read
     */
    const markAllAsRead = useCallback(() => {
        setActivities(prev => 
            prev.map(activity => ({ ...activity, isRead: true }))
        );
        // TODO: Call API endpoint to mark all as read on server if available
    }, []);
    /**
     * Archive (soft delete) an activity
     */
    const archiveActivityItem = useCallback(async (activityId, reason = 'User archived activity') => {
        try {
            await archiveActivity(activityId, reason);
            setActivities(prev => 
                prev.filter(activity => activity.id !== activityId)
            );
        } catch (err) {
            throw err;
        }
    }, []);
    /**
     * Restore an archived activity
     */
    const restoreActivityItem = useCallback(async (activityId, reason = 'User restored activity') => {
        try {
            await restoreActivity(activityId, reason);
            // Refresh to get the restored activity
            await loadActivities(true);
        } catch (err) {
            throw err;
        }
    }, [loadActivities]);
    /**
     * Create a new activity
     */
    const createNewActivity = useCallback(async () => {
        if (!userId) return null;
        try {
            const response = await createActivity(userId);
            // Refresh activities to include the new one
            await loadActivities(true);
            return response;
        } catch (err) {
            throw err;
        }
    }, [userId, loadActivities]);
    /**
     * Update filters
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
     * Toggle activity read status
     */
    const toggleRead = useCallback((activityId) => {
        setActivities(prev => 
            prev.map(activity => 
                activity.id === activityId 
                    ? { ...activity, isRead: !activity.isRead } 
                    : activity
            )
        );
    }, []);
    // Initial load effect
    useEffect(() => {
        loadActivities();
    }, [loadActivities]);
    // Auto-refresh effect
    useEffect(() => {
        if (!autoRefresh || !userId) return;
        const interval = setInterval(() => {
            loadActivities(true); // Silent refresh
        }, refreshInterval);
        return () => clearInterval(interval);
    }, [autoRefresh, userId, refreshInterval, loadActivities]);
    // Return the hook interface
    return {
        // State
        activities: filteredActivities,
        allActivities: activities,
        loading,
        error,
        filters,
        // Computed values
        unreadCount,
        hasActivities,
        isEmpty,
        // Actions
        refresh,
        loadActivities,
        markAsRead,
        markAllAsRead,
        archiveActivityItem,
        restoreActivityItem,
        createNewActivity,
        toggleRead,
        // Filters
        updateFilters,
        clearFilters
    };
}; 