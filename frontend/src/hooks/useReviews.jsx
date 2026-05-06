import { useState, useEffect, useRef, useCallback } from 'react';
import { 
    getApprovedReviews, 
    getReviewsBySender, 
    getReviewsForItem,
    getMostReviewed,
    approveReview, 
    rejectReview,
    editReview,
    isReviewsResponseSuccess,
    getReviewsErrorMessage
} from '../api/reviewsApi';
import { useAuth } from '../context/AuthContext';
const useReviews = () => {
    const { user } = useAuth();
    const [reviews, setReviews] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [searchTerm, setSearchTerm] = useState('');
    const [statusFilter, setStatusFilter] = useState('all');
    const [flaggedFilter, setFlaggedFilter] = useState('all');
    const [itemFilter, setItemFilter] = useState('all');
    const [actionLoading, setActionLoading] = useState({});
    const [autoRefresh, setAutoRefresh] = useState(true);
    const [stats, setStats] = useState({
        total: 0,
        pending: 0,
        approved: 0,
        rejected: 0,
        flagged: 0
    });
    const isMountedRef = useRef(true);
    const refreshIntervalRef = useRef(null);
    /**
     * Calculate statistics from reviews
     */
    const calculateStats = useCallback((reviewsList) => {
        const stats = {
            total: reviewsList.length,
            pending: reviewsList.filter(r => r.reviewStatus === 'pending').length,
            approved: reviewsList.filter(r => r.reviewStatus === 'approved').length,
            rejected: reviewsList.filter(r => r.reviewStatus === 'rejected').length,
            flagged: reviewsList.filter(r => r.flagged === true).length
        };
        setStats(stats);
        return stats;
    }, []);
    /**
     * Fetch reviews data
     */
    const fetchReviews = useCallback(async (showLoading = true) => {
        if (showLoading) {
            setLoading(true);
        }
        setError(null);
        try {
            // Default to approved reviews if user is not authenticated
            // If user is admin/moderator, they can see all reviews
            const response = await getApprovedReviews();
            if (!isMountedRef.current) return;
            if (isReviewsResponseSuccess(response)) {
                const reviewsList = response.reviews || [];
                setReviews(reviewsList);
                calculateStats(reviewsList);
            } else {
                const errorMessage = getReviewsErrorMessage(response);
                setError(errorMessage);
            }
        } catch (err) {
            if (!isMountedRef.current) return;
            const errorMessage = 'Failed to fetch reviews';
            setError(errorMessage);
        } finally {
            if (isMountedRef.current && showLoading) {
                setLoading(false);
            }
        }
    }, [calculateStats]);
    /**
     * Fetch reviews by current user
     */
    const fetchMyReviews = useCallback(async () => {
        if (!user?.id && !user?.userId) {
            setReviews([]);
            setLoading(false);
            return;
        }
        setLoading(true);
        setError(null);
        try {
            const userId = user.id || user.userId;
            const response = await getReviewsBySender(userId);
            if (!isMountedRef.current) return;
            if (isReviewsResponseSuccess(response)) {
                const reviewsList = response.reviews || [];
                setReviews(reviewsList);
                calculateStats(reviewsList);
            } else {
                const errorMessage = getReviewsErrorMessage(response);
                setError(errorMessage);
            }
        } catch (err) {
            if (!isMountedRef.current) return;
            const errorMessage = 'Failed to fetch your reviews';
            setError(errorMessage);
        } finally {
            if (isMountedRef.current) {
                setLoading(false);
            }
        }
    }, [user, calculateStats]);
    /**
     * Fetch reviews for specific item
     */
    const fetchItemReviews = useCallback(async (itemId) => {
        if (!itemId) {
            setReviews([]);
            setLoading(false);
            return;
        }
        setLoading(true);
        setError(null);
        try {
            const response = await getReviewsForItem(itemId);
            if (!isMountedRef.current) return;
            if (isReviewsResponseSuccess(response)) {
                const reviewsList = response.reviews || [];
                setReviews(reviewsList);
                calculateStats(reviewsList);
            } else {
                const errorMessage = getReviewsErrorMessage(response);
                setError(errorMessage);
            }
        } catch (err) {
            if (!isMountedRef.current) return;
            const errorMessage = 'Failed to fetch item reviews';
            setError(errorMessage);
        } finally {
            if (isMountedRef.current) {
                setLoading(false);
            }
        }
    }, [calculateStats]);
    /**
     * Approve a review with optimistic updates
     */
    const handleApproveReview = useCallback(async (reviewId) => {
        if (!reviewId) return;
        setActionLoading(prev => ({ ...prev, [reviewId]: 'approving' }));
        try {
            // Optimistic update
            setReviews(prev => prev.map(review => 
                review.id === reviewId 
                    ? { ...review, reviewStatus: 'approved' }
                    : review
            ));
            const response = await approveReview(reviewId);
            if (!isReviewsResponseSuccess(response)) {
                // Revert optimistic update on failure
                setReviews(prev => prev.map(review => 
                    review.id === reviewId 
                        ? { ...review, reviewStatus: 'pending' }
                        : review
                ));
                const errorMessage = getReviewsErrorMessage(response);
                setError(`Failed to approve review: ${errorMessage}`);
            } else {
                // Update with server response if available
                if (response.id) {
                    setReviews(prev => prev.map(review => 
                        review.id === reviewId 
                            ? { ...review, ...response }
                            : review
                    ));
                }
            }
        } catch (err) {
            // Revert optimistic update on error
            setReviews(prev => prev.map(review => 
                review.id === reviewId 
                    ? { ...review, reviewStatus: 'pending' }
                    : review
            ));
            setError('Failed to approve review');
        } finally {
            setActionLoading(prev => {
                const newState = { ...prev };
                delete newState[reviewId];
                return newState;
            });
        }
    }, []);
    /**
     * Reject a review with optimistic updates
     */
    const handleRejectReview = useCallback(async (reviewId) => {
        if (!reviewId) return;
        setActionLoading(prev => ({ ...prev, [reviewId]: 'rejecting' }));
        try {
            // Optimistic update
            setReviews(prev => prev.map(review => 
                review.id === reviewId 
                    ? { ...review, reviewStatus: 'rejected' }
                    : review
            ));
            const response = await rejectReview(reviewId);
            if (!isReviewsResponseSuccess(response)) {
                // Revert optimistic update on failure
                setReviews(prev => prev.map(review => 
                    review.id === reviewId 
                        ? { ...review, reviewStatus: 'pending' }
                        : review
                ));
                const errorMessage = getReviewsErrorMessage(response);
                setError(`Failed to reject review: ${errorMessage}`);
            } else {
                // Update with server response if available
                if (response.id) {
                    setReviews(prev => prev.map(review => 
                        review.id === reviewId 
                            ? { ...review, ...response }
                            : review
                    ));
                }
            }
        } catch (err) {
            // Revert optimistic update on error
            setReviews(prev => prev.map(review => 
                review.id === reviewId 
                    ? { ...review, reviewStatus: 'pending' }
                    : review
            ));
            setError('Failed to reject review');
        } finally {
            setActionLoading(prev => {
                const newState = { ...prev };
                delete newState[reviewId];
                return newState;
            });
        }
    }, []);
    /**
     * Edit a review with optimistic updates
     */
    const handleEditReview = useCallback(async (reviewId, newContent) => {
        if (!reviewId || !newContent) return;
        setActionLoading(prev => ({ ...prev, [reviewId]: 'editing' }));
        const originalContent = reviews.find(r => r.id === reviewId)?.content;
        try {
            // Optimistic update
            setReviews(prev => prev.map(review => 
                review.id === reviewId 
                    ? { ...review, content: newContent }
                    : review
            ));
            const response = await editReview(reviewId, newContent);
            if (!isReviewsResponseSuccess(response)) {
                // Revert optimistic update on failure
                setReviews(prev => prev.map(review => 
                    review.id === reviewId 
                        ? { ...review, content: originalContent }
                        : review
                ));
                const errorMessage = getReviewsErrorMessage(response);
                setError(`Failed to edit review: ${errorMessage}`);
            } else {
                // Update with server response if available
                if (response.id) {
                    setReviews(prev => prev.map(review => 
                        review.id === reviewId 
                            ? { ...review, ...response }
                            : review
                    ));
                }
            }
        } catch (err) {
            // Revert optimistic update on error
            setReviews(prev => prev.map(review => 
                review.id === reviewId 
                    ? { ...review, content: originalContent }
                    : review
            ));
            setError('Failed to edit review');
        } finally {
            setActionLoading(prev => {
                const newState = { ...prev };
                delete newState[reviewId];
                return newState;
            });
        }
    }, [reviews]);
    /**
     * Filter reviews based on current filters
     */
    const filteredReviews = useCallback(() => {
        let filtered = [...reviews];
        // Search filter
        if (searchTerm.trim()) {
            const term = searchTerm.toLowerCase();
            filtered = filtered.filter(review =>
                review.content?.toLowerCase().includes(term) ||
                review.senderName?.toLowerCase().includes(term) ||
                review.itemName?.toLowerCase().includes(term) ||
                review.itemType?.toLowerCase().includes(term)
            );
        }
        // Status filter
        if (statusFilter !== 'all') {
            filtered = filtered.filter(review => review.reviewStatus === statusFilter);
        }
        // Flagged filter
        if (flaggedFilter !== 'all') {
            const isFlagged = flaggedFilter === 'flagged';
            filtered = filtered.filter(review => review.flagged === isFlagged);
        }
        // Item filter (could be item type or specific item)
        if (itemFilter !== 'all') {
            filtered = filtered.filter(review => 
                review.itemType === itemFilter || 
                review.itemId === itemFilter
            );
        }
        return filtered.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt));
    }, [reviews, searchTerm, statusFilter, flaggedFilter, itemFilter]);
    /**
     * Refresh reviews data
     */
    const refresh = useCallback(() => {
        fetchReviews(false);
    }, [fetchReviews]);
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
                    fetchReviews(false);
                }
            }, 30000); // Refresh every 30 seconds
            return () => {
                if (refreshIntervalRef.current) {
                    clearInterval(refreshIntervalRef.current);
                    refreshIntervalRef.current = null;
                }
            };
        }
    }, [autoRefresh, fetchReviews]);
    /**
     * Initial data fetch
     */
    useEffect(() => {
        fetchReviews();
    }, [fetchReviews]);
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
     * Recalculate stats when reviews change
     */
    useEffect(() => {
        calculateStats(reviews);
    }, [reviews, calculateStats]);
    return {
        // Data
        reviews: filteredReviews(),
        allReviews: reviews,
        loading,
        error,
        stats,
        actionLoading,
        // Filters
        searchTerm,
        setSearchTerm,
        statusFilter,
        setStatusFilter,
        flaggedFilter,
        setFlaggedFilter,
        itemFilter,
        setItemFilter,
        // Actions
        handleApproveReview,
        handleRejectReview,
        handleEditReview,
        refresh,
        clearError,
        // Fetch methods
        fetchReviews,
        fetchMyReviews,
        fetchItemReviews,
        // Settings
        autoRefresh,
        setAutoRefresh
    };
};
export default useReviews; 