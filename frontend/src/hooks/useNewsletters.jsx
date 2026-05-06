import { useState, useEffect, useRef, useCallback } from 'react';
import { 
    listSubscriptions,
    subscribeNewsletter,
    unsubscribeNewsletter,
    updateSubscription,
    getSubscription,
    fetchNewsletters,
    getNewsletter,
    createEdition,
    sendEdition
} from '../api/client/newsletterApi';
import { useAuth } from '../context/AuthContext';
const useNewsletters = () => {
    const { user } = useAuth();
    const [subscriptions, setSubscriptions] = useState([]);
    const [availableNewsletters, setAvailableNewsletters] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [searchTerm, setSearchTerm] = useState('');
    const [statusFilter, setStatusFilter] = useState('all');
    const [categoryFilter, setCategoryFilter] = useState('all');
    const [actionLoading, setActionLoading] = useState({});
    const [autoRefresh, setAutoRefresh] = useState(true);
    const [stats, setStats] = useState({
        total: 0,
        active: 0,
        paused: 0,
        cancelled: 0,
        pending: 0
    });
    const isMountedRef = useRef(true);
    const refreshIntervalRef = useRef(null);
    /**
     * Calculate statistics from subscriptions
     */
    const calculateStats = useCallback((subscriptionsList) => {
        const stats = {
            total: subscriptionsList.length,
            active: subscriptionsList.filter(s => s.status === 'active').length,
            paused: subscriptionsList.filter(s => s.status === 'paused').length,
            cancelled: subscriptionsList.filter(s => s.status === 'cancelled' || s.status === 'unsubscribed').length,
            pending: subscriptionsList.filter(s => s.status === 'pending').length
        };
        setStats(stats);
        return stats;
    }, []);
    /**
     * Fetch subscriptions data
     */
    const fetchSubscriptions = useCallback(async (showLoading = true) => {
        if (showLoading) {
            setLoading(true);
        }
        setError(null);
        try {
            const params = {};
            if (user?.id || user?.userId) {
                params.userId = user.id || user.userId;
            }
            const response = await listSubscriptions(params);
            if (!isMountedRef.current) return;
            if (isNewsletterResponseSuccess(response)) {
                const subscriptionsList = response.subscriptions || [];
                setSubscriptions(subscriptionsList);
                calculateStats(subscriptionsList);
            } else {
                const errorMessage = getNewsletterErrorMessage(response);
                setError(errorMessage);
            }
        } catch (err) {
            if (!isMountedRef.current) return;
            const errorMessage = 'Failed to fetch subscriptions';
            setError(errorMessage);
        } finally {
            if (isMountedRef.current && showLoading) {
                setLoading(false);
            }
        }
    }, [user, calculateStats]);
    /**
     * Load available newsletters
     */
    const loadAvailableNewsletters = useCallback(async () => {
        try {
            const response = await fetchNewsletters({ activeOnly: true });
            const newsletters = response.newsletters || [];
            setAvailableNewsletters(newsletters);
        } catch (err) {
            // Error: 'Failed to load newsletters:', err...
            setAvailableNewsletters([]);
        }
    }, []);
    /**
     * Subscribe to a newsletter
     */
    const handleSubscribe = useCallback(async (newsletterId, preferences = 'weekly_digest') => {
        if (!user?.id && !user?.userId) {
            setError('Please log in to subscribe to newsletters');
            return;
        }
        const currentUserId = user.id || user.userId;
        setActionLoading(prev => ({ ...prev, [newsletterId]: 'subscribing' }));
        try {
            const response = await subscribeNewsletter(newsletterId, {
                frequencyOverride: preferences,
                topics: [],
                format: 'html'
            });
            // Refresh the subscriptions list
            await fetchSubscriptions(false);
        } catch (err) {
            setError('Failed to subscribe to newsletter');
        } finally {
            setActionLoading(prev => {
                const newState = { ...prev };
                delete newState[newsletterId];
                return newState;
            });
        }
    }, [user, fetchSubscriptions]);
    /**
     * Unsubscribe from a newsletter with optimistic updates
     */
    const handleUnsubscribe = useCallback(async (subscriptionId) => {
        if (!subscriptionId) return;
        setActionLoading(prev => ({ ...prev, [subscriptionId]: 'unsubscribing' }));
        try {
            // Optimistic update - remove from list
            const originalSubscriptions = subscriptions;
            setSubscriptions(prev => prev.filter(sub => sub.subscriptionId !== subscriptionId));
            await unsubscribeNewsletter(subscriptionId);
        } catch (err) {
            // Revert optimistic update on error
            setSubscriptions(subscriptions);
            setError('Failed to unsubscribe from newsletter');
        } finally {
            setActionLoading(prev => {
                const newState = { ...prev };
                delete newState[subscriptionId];
                return newState;
            });
        }
    }, [subscriptions]);
    /**
     * Update subscription preferences with optimistic updates
     */
    const handleUpdateSubscription = useCallback(async (subscriptionId, updateData) => {
        if (!subscriptionId) return;
        setActionLoading(prev => ({ ...prev, [subscriptionId]: 'updating' }));
        try {
            // Optimistic update
            setSubscriptions(prev => prev.map(subscription => 
                subscription.subscriptionId === subscriptionId 
                    ? { ...subscription, ...updateData }
                    : subscription
            ));
            const response = await updateSubscription(subscriptionId, updateData);
            // Update with server response if available
            if (response.subscription) {
                setSubscriptions(prev => prev.map(subscription => 
                    subscription.id === subscriptionId 
                        ? { ...subscription, ...response.subscription }
                        : subscription
                ));
            }
        } catch (err) {
            // Revert optimistic update on error - refresh data
            await fetchSubscriptions(false);
            setError('Failed to update subscription');
        } finally {
            setActionLoading(prev => {
                const newState = { ...prev };
                delete newState[subscriptionId];
                return newState;
            });
        }
    }, [fetchSubscriptions]);
    /**
     * Pause/Resume subscription
     */
    const handleToggleSubscription = useCallback(async (subscriptionId, currentStatus) => {
        const newStatus = currentStatus === 'active' ? 'paused' : 'active';
        await handleUpdateSubscription(subscriptionId, { subscriptionStatus: newStatus });
    }, [handleUpdateSubscription]);
    /**
     * Send newsletter edition (admin function)
     */
    const handleSendNewsletter = useCallback(async (editionId, testMode = false) => {
        setActionLoading(prev => ({ ...prev, [`send_${editionId}`]: 'sending' }));
        try {
            const response = await sendEdition(editionId, testMode);
            return response;
        } catch (err) {
            setError(err.message || 'Failed to send newsletter');
            return { success: false };
        } finally {
            setActionLoading(prev => {
                const newState = { ...prev };
                delete newState[`send_${editionId}`];
                return newState;
            });
        }
    }, []);
    /**
     * Filter subscriptions based on current filters
     */
    const filteredSubscriptions = useCallback(() => {
        let filtered = [...subscriptions];
        // Search filter
        if (searchTerm.trim()) {
            const term = searchTerm.toLowerCase();
            filtered = filtered.filter(subscription =>
                subscription.newsletterName?.toLowerCase().includes(term) ||
                subscription.newsletterDescription?.toLowerCase().includes(term) ||
                subscription.subscriptionPreferences?.toLowerCase().includes(term)
            );
        }
        // Status filter
        if (statusFilter !== 'all') {
            filtered = filtered.filter(subscription => subscription.subscriptionStatus === statusFilter);
        }
        // Category filter (based on available newsletters)
        if (categoryFilter !== 'all') {
            const categoryNewsletters = availableNewsletters
                .filter(nl => nl.category === categoryFilter)
                .map(nl => nl.newsletterId);
            filtered = filtered.filter(subscription => 
                categoryNewsletters.includes(subscription.newsletterId)
            );
        }
        return filtered.sort((a, b) => new Date(b.updatedAt || b.createdAt) - new Date(a.updatedAt || a.createdAt));
    }, [subscriptions, searchTerm, statusFilter, categoryFilter, availableNewsletters]);
    /**
     * Check if user is subscribed to a newsletter
     */
    const isSubscribed = useCallback((newsletterId) => {
        return subscriptions.some(s => 
            s.newsletterId === newsletterId && 
            (s.subscriptionStatus === 'active' || s.subscriptionStatus === 'pending')
        );
    }, [subscriptions]);
    /**
     * Get subscription by newsletter ID
     */
    const getSubscriptionByNewsletter = useCallback((newsletterId) => {
        return subscriptions.find(s => s.newsletterId === newsletterId);
    }, [subscriptions]);
    /**
     * Get subscription count for a specific status
     */
    const getSubscriptionCount = useCallback((status = 'active') => {
        if (status === 'all') return subscriptions.length;
        return subscriptions.filter(s => s.subscriptionStatus === status).length;
    }, [subscriptions]);
    /**
     * Refresh subscriptions data
     */
    const refresh = useCallback(() => {
        fetchSubscriptions(false);
    }, [fetchSubscriptions]);
    /**
     * Clear error
     */
    const clearError = useCallback(() => {
        setError(null);
    }, []);
    /**
     * Setup auto-refresh
     */
    useEffect(() => {
        if (autoRefresh) {
            refreshIntervalRef.current = setInterval(() => {
                if (isMountedRef.current) {
                    fetchSubscriptions(false);
                }
            }, 30000); // Refresh every 30 seconds
            return () => {
                if (refreshIntervalRef.current) {
                    clearInterval(refreshIntervalRef.current);
                    refreshIntervalRef.current = null;
                }
            };
        }
    }, [autoRefresh, fetchSubscriptions]);
    /**
     * Initial data fetch
     */
    useEffect(() => {
        fetchSubscriptions();
        loadAvailableNewsletters();
    }, [fetchSubscriptions, loadAvailableNewsletters]);
    /**
     * Cleanup on unmount
     */
    useEffect(() => {
        return () => {
            isMountedRef.current = false;
            if (refreshIntervalRef.current) {
                clearInterval(refreshIntervalRef.current);
            }
        };
    }, []);
    /**
     * Recalculate stats when subscriptions change
     */
    useEffect(() => {
        calculateStats(subscriptions);
    }, [subscriptions, calculateStats]);
    return {
        // Data
        subscriptions: filteredSubscriptions(),
        allSubscriptions: subscriptions,
        availableNewsletters,
        loading,
        error,
        stats,
        actionLoading,
        // Filters
        searchTerm,
        setSearchTerm,
        statusFilter,
        setStatusFilter,
        categoryFilter,
        setCategoryFilter,
        // Actions
        handleSubscribe,
        handleUnsubscribe,
        handleUpdateSubscription,
        handleToggleSubscription,
        handleSendNewsletter,
        refresh,
        clearError,
        // Utility functions
        isSubscribed,
        getSubscriptionByNewsletter,
        getSubscriptionCount,
        // Fetch methods
        fetchSubscriptions,
        // Settings
        autoRefresh,
        setAutoRefresh
    };
};
export default useNewsletters; 