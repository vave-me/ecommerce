import { useState, useEffect, useRef, useCallback } from 'react';
import { 
    getFollowers, 
    getFollowing, 
    getApprovedFollowing,
    getMostFollowed,
    addFollow,
    approveFollow, 
    rejectFollow,
    isFollowingResponseSuccess,
    getFollowingErrorMessage
} from '../api/followingApi';
import { useAuth } from '../context/AuthContext';
const useFollowing = (userId = null, type = 'followers') => {
    const { user } = useAuth();
    const [followers, setFollowers] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [searchTerm, setSearchTerm] = useState('');
    const [statusFilter, setStatusFilter] = useState('all');
    const [actionLoading, setActionLoading] = useState({});
    const [autoRefresh, setAutoRefresh] = useState(true);
    const [stats, setStats] = useState({
        total: 0,
        pending: 0,
        approved: 0,
        rejected: 0
    });
    const isMountedRef = useRef(true);
    const refreshIntervalRef = useRef(null);
    const targetUserId = userId || user?.id || user?.userId;
    /**
     * Calculate statistics from followers
     */
    const calculateStats = useCallback((followersList) => {
        const stats = {
            total: followersList.length,
            pending: followersList.filter(f => f.followStatus === 'pending').length,
            approved: followersList.filter(f => f.followStatus === 'approved').length,
            rejected: followersList.filter(f => f.followStatus === 'rejected').length
        };
        setStats(stats);
        return stats;
    }, []);
    /**
     * Fetch followers data
     */
    const fetchFollowers = useCallback(async (showLoading = true) => {
        if (!targetUserId) {
            setFollowers([]);
            setLoading(false);
            return;
        }
        if (showLoading) {
            setLoading(true);
        }
        setError(null);
        try {
            let response;
            if (type === 'followers') {
                // Get users following this user
                response = await getFollowers(targetUserId);
            } else if (type === 'following') {
                // Get users this user is following
                response = await getFollowing(targetUserId);
            } else {
                // Get approved following for general use
                response = await getApprovedFollowing();
            }
            if (!isMountedRef.current) return;
            if (isFollowingResponseSuccess(response)) {
                const followersList = response.following || [];
                setFollowers(followersList);
                calculateStats(followersList);
            } else {
                const errorMessage = getFollowingErrorMessage(response);
                setError(errorMessage);
            }
        } catch (err) {
            if (!isMountedRef.current) return;
            const errorMessage = 'Failed to fetch followers';
            setError(errorMessage);
        } finally {
            if (isMountedRef.current && showLoading) {
                setLoading(false);
            }
        }
    }, [targetUserId, type, calculateStats]);
    /**
     * Follow a user
     */
    const handleFollowUser = useCallback(async (followedUserId) => {
        if (!user?.id && !user?.userId) {
            setError('Please log in to follow users');
            return;
        }
        const currentUserId = user.id || user.userId;
        if (currentUserId === followedUserId) {
            setError('You cannot follow yourself');
            return;
        }
        setActionLoading(prev => ({ ...prev, [followedUserId]: 'following' }));
        try {
            const followData = {
                userId: currentUserId,
                followedUserId: followedUserId,
                followedUserType: 'regular_user',
                content: null,
                categoryId: 'cat_general',
                parentId: null
            };
            const response = await addFollow(followData);
            if (isFollowingResponseSuccess(response)) {
                // Refresh the followers list
                await fetchFollowers(false);
            } else {
                const errorMessage = getFollowingErrorMessage(response);
                setError(`Failed to follow user: ${errorMessage}`);
            }
        } catch (err) {
            setError('Failed to follow user');
        } finally {
            setActionLoading(prev => {
                const newState = { ...prev };
                delete newState[followedUserId];
                return newState;
            });
        }
    }, [user, fetchFollowers]);
    /**
     * Approve a follow request with optimistic updates
     */
    const handleApproveFollow = useCallback(async (followId) => {
        if (!followId) return;
        setActionLoading(prev => ({ ...prev, [followId]: 'approving' }));
        try {
            // Optimistic update
            setFollowers(prev => prev.map(follower => 
                follower.id === followId 
                    ? { ...follower, followStatus: 'approved' }
                    : follower
            ));
            const response = await approveFollow(followId);
            if (!isFollowingResponseSuccess(response)) {
                // Revert optimistic update on failure
                setFollowers(prev => prev.map(follower => 
                    follower.id === followId 
                        ? { ...follower, followStatus: 'pending' }
                        : follower
                ));
                const errorMessage = getFollowingErrorMessage(response);
                setError(`Failed to approve follow: ${errorMessage}`);
            } else {
                // Update with server response if available
                if (response.id) {
                    setFollowers(prev => prev.map(follower => 
                        follower.id === followId 
                            ? { ...follower, ...response }
                            : follower
                    ));
                }
            }
        } catch (err) {
            // Revert optimistic update on error
            setFollowers(prev => prev.map(follower => 
                follower.id === followId 
                    ? { ...follower, followStatus: 'pending' }
                    : follower
            ));
            setError('Failed to approve follow');
        } finally {
            setActionLoading(prev => {
                const newState = { ...prev };
                delete newState[followId];
                return newState;
            });
        }
    }, []);
    /**
     * Reject a follow request with optimistic updates
     */
    const handleRejectFollow = useCallback(async (followId) => {
        if (!followId) return;
        setActionLoading(prev => ({ ...prev, [followId]: 'rejecting' }));
        try {
            // Optimistic update - remove from list
            const originalFollowers = followers;
            setFollowers(prev => prev.filter(follower => follower.id !== followId));
            const response = await rejectFollow(followId);
            if (!isFollowingResponseSuccess(response)) {
                // Revert optimistic update on failure
                setFollowers(originalFollowers);
                const errorMessage = getFollowingErrorMessage(response);
                setError(`Failed to reject follow: ${errorMessage}`);
            }
        } catch (err) {
            // Revert optimistic update on error
            setFollowers(followers);
            setError('Failed to reject follow');
        } finally {
            setActionLoading(prev => {
                const newState = { ...prev };
                delete newState[followId];
                return newState;
            });
        }
    }, [followers]);
    /**
     * Filter followers based on current filters
     */
    const filteredFollowers = useCallback(() => {
        let filtered = [...followers];
        // Search filter
        if (searchTerm.trim()) {
            const term = searchTerm.toLowerCase();
            filtered = filtered.filter(follower =>
                follower.followerName?.toLowerCase().includes(term) ||
                follower.followerBio?.toLowerCase().includes(term) ||
                follower.userId?.toLowerCase().includes(term)
            );
        }
        // Status filter
        if (statusFilter !== 'all') {
            filtered = filtered.filter(follower => follower.followStatus === statusFilter);
        }
        return filtered.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt));
    }, [followers, searchTerm, statusFilter]);
    /**
     * Check if current user is following someone
     */
    const isFollowing = useCallback((userId) => {
        return followers.some(f => 
            f.followedUserId === userId && 
            f.followStatus === 'approved' &&
            (f.userId === user?.id || f.userId === user?.userId)
        );
    }, [followers, user]);
    /**
     * Get follower count for a specific status
     */
    const getFollowerCount = useCallback((status = 'approved') => {
        if (status === 'all') return followers.length;
        return followers.filter(f => f.followStatus === status).length;
    }, [followers]);
    /**
     * Refresh followers data
     */
    const refresh = useCallback(() => {
        fetchFollowers(false);
    }, [fetchFollowers]);
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
                    fetchFollowers(false);
                }
            }, 30000); // Refresh every 30 seconds
            return () => {
                if (refreshIntervalRef.current) {
                    clearInterval(refreshIntervalRef.current);
                    refreshIntervalRef.current = null;
                }
            };
        }
    }, [autoRefresh, fetchFollowers]);
    /**
     * Initial data fetch
     */
    useEffect(() => {
        fetchFollowers();
    }, [fetchFollowers]);
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
     * Recalculate stats when followers change
     */
    useEffect(() => {
        calculateStats(followers);
    }, [followers, calculateStats]);
    return {
        // Data
        followers: filteredFollowers(),
        allFollowers: followers,
        loading,
        error,
        stats,
        actionLoading,
        targetUserId,
        // Filters
        searchTerm,
        setSearchTerm,
        statusFilter,
        setStatusFilter,
        // Actions
        handleFollowUser,
        handleApproveFollow,
        handleRejectFollow,
        refresh,
        clearError,
        // Utility functions
        isFollowing,
        getFollowerCount,
        // Fetch methods
        fetchFollowers,
        // Settings
        autoRefresh,
        setAutoRefresh
    };
};
export default useFollowing; 