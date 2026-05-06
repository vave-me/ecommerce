"use client";
import React, { useState } from 'react';
import { 
    MessageCircle, CheckCircle, XCircle, Clock, Flag, 
    Search, Filter, RefreshCw, User, Package, AlertTriangle,
    Loader2, Edit3, Trash2, Eye
} from '@/icons';
import useReviews from '../../../hooks/useReviews';
import { useAuth } from '../../../context/AuthContext';
import { getReviewStatusTypes, formatDate } from '../../../api/reviewsApi';
import styles from './page.module.css';
// Icon mapping for dynamic icon rendering
const iconMap = {
    MessageCircle, CheckCircle, XCircle, Clock, Flag, Search, Filter, 
    RefreshCw, User, Package, AlertTriangle, Loader2, Edit3, Trash2, Eye
};
// Get icon component by name
const getIcon = (iconName, props = {}) => {
    const IconComponent = iconMap[iconName];
    return IconComponent ? <IconComponent {...props} /> : null;
};
/**
 * Individual Review Item Component
 */
const ReviewItem = ({ review, onApprove, onReject, onEdit, actionLoading, isUserReview = false }) => {
    const [isEditing, setIsEditing] = useState(false);
    const [editContent, setEditContent] = useState(review.content || '');
    const statusTypes = getReviewStatusTypes();
    const statusInfo = statusTypes.find(s => s.value === review.reviewStatus) || statusTypes[0];
    const handleEdit = () => {
        if (isEditing) {
            if (editContent.trim() && editContent !== review.content) {
                onEdit(review.id, editContent.trim());
            }
            setIsEditing(false);
        } else {
            setEditContent(review.content || '');
            setIsEditing(true);
        }
    };
    const handleCancelEdit = () => {
        setEditContent(review.content || '');
        setIsEditing(false);
    };
    const isActionLoading = actionLoading[review.id];
    return (
        <div className={`${styles.reviewItem} ${review.flagged ? styles.flagged : ''}`}>
            {/* Status indicator */}
            <div 
                className={styles.statusBorder}
                style={{ backgroundColor: statusInfo.color }}
            />
            <div className={styles.reviewHeader}>
                <div className={styles.reviewMeta}>
                    <div className={styles.reviewUser}>
                        {getIcon('User', { className: styles.icon })}
                        <span className={styles.userName}>{review.senderName || 'Anonymous'}</span>
                        {review.flagged && (
                            <span className={styles.flaggedBadge}>
                                {getIcon('Flag', { className: styles.iconSmall })}
                                Flagged
                            </span>
                        )}
                    </div>
                    <div className={styles.reviewInfo}>
                        <span className={styles.itemInfo}>
                            {getIcon('Package', { className: styles.iconSmall })}
                            {review.itemName || review.itemId} ({review.itemType})
                        </span>
                        <span className={styles.reviewDate}>
                            {formatDate(review.createdAt)}
                        </span>
                    </div>
                </div>
                <div className={styles.statusBadge} style={{ color: statusInfo.color }}>
                    {getIcon(statusInfo.icon, { className: styles.iconSmall })}
                    {statusInfo.label}
                </div>
            </div>
            <div className={styles.reviewContent}>
                {isEditing ? (
                    <div className={styles.editForm}>
                        <textarea
                            value={editContent}
                            onChange={(e) => setEditContent(e.target.value)}
                            className={styles.editTextarea}
                            rows={3}
                            maxLength={1000}
                            placeholder="Edit your review..."
                        />
                        <div className={styles.editActions}>
                            <button
                                onClick={handleEdit}
                                disabled={!editContent.trim() || isActionLoading}
                                className={`${styles.actionButton} ${styles.saveButton}`}
                            >
                                {isActionLoading === 'editing' ? getIcon('Loader2', { className: `${styles.iconSmall} ${styles.spinning}` }) : getIcon('CheckCircle', { className: styles.iconSmall })}
                                Save
                            </button>
                            <button
                                onClick={handleCancelEdit}
                                disabled={isActionLoading}
                                className={`${styles.actionButton} ${styles.cancelButton}`}
                            >
                                Cancel
                            </button>
                        </div>
                        <div className={styles.charCount}>
                            {editContent.length}/1000
                        </div>
                    </div>
                ) : (
                    <p className={styles.reviewText}>{review.content}</p>
                )}
            </div>
            {/* Actions */}
            <div className={styles.reviewActions}>
                {isUserReview && review.reviewStatus !== 'rejected' && (
                    <button
                        onClick={handleEdit}
                        disabled={isActionLoading}
                        className={`${styles.actionButton} ${styles.editButton}`}
                        title="Edit review"
                    >
                        {getIcon('Edit3', { className: styles.iconSmall })}
                        Edit
                    </button>
                )}
                {!isUserReview && review.reviewStatus === 'pending' && (
                    <>
                        <button
                            onClick={() => onApprove(review.id)}
                            disabled={isActionLoading}
                            className={`${styles.actionButton} ${styles.approveButton}`}
                            title="Approve review"
                        >
                            {isActionLoading === 'approving' ? getIcon('Loader2', { className: `${styles.iconSmall} ${styles.spinning}` }) : getIcon('CheckCircle', { className: styles.iconSmall })}
                            Approve
                        </button>
                        <button
                            onClick={() => onReject(review.id)}
                            disabled={isActionLoading}
                            className={`${styles.actionButton} ${styles.rejectButton}`}
                            title="Reject review"
                        >
                            {isActionLoading === 'rejecting' ? getIcon('Loader2', { className: `${styles.iconSmall} ${styles.spinning}` }) : getIcon('XCircle', { className: styles.iconSmall })}
                            Reject
                        </button>
                    </>
                )}
            </div>
        </div>
    );
};
/**
 * Filter Bar Component
 */
const FilterBar = ({ 
    searchTerm, 
    setSearchTerm, 
    statusFilter, 
    setStatusFilter,
    flaggedFilter,
    setFlaggedFilter,
    onRefresh,
    loading 
}) => {
    const statusTypes = getReviewStatusTypes();
    return (
        <div className={styles.filterBar}>
            <div className={styles.searchContainer}>
                <div className={styles.searchInputWrapper}>
                    {getIcon('Search', { className: styles.searchIcon })}
                    <input
                        type="text"
                        placeholder="Search reviews..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        className={styles.searchInput}
                    />
                </div>
            </div>
            <div className={styles.filtersContainer}>
                <div className={styles.filterGroup}>
                    <label className={styles.filterLabel}>Status:</label>
                    <select
                        value={statusFilter}
                        onChange={(e) => setStatusFilter(e.target.value)}
                        className={styles.filterSelect}
                    >
                        <option value="all">All Statuses</option>
                        {statusTypes.map(status => (
                            <option key={status.value} value={status.value}>
                                {status.label}
                            </option>
                        ))}
                    </select>
                </div>
                <div className={styles.filterGroup}>
                    <label className={styles.filterLabel}>Flagged:</label>
                    <select
                        value={flaggedFilter}
                        onChange={(e) => setFlaggedFilter(e.target.value)}
                        className={styles.filterSelect}
                    >
                        <option value="all">All Reviews</option>
                        <option value="flagged">Flagged Only</option>
                        <option value="not_flagged">Not Flagged</option>
                    </select>
                </div>
                <button
                    onClick={onRefresh}
                    disabled={loading}
                    className={styles.refreshButton}
                    title="Refresh reviews"
                >
                    {getIcon('RefreshCw', { className: `${styles.iconSmall} ${loading ? styles.spinning : ''}` })}
                </button>
            </div>
        </div>
    );
};
/**
 * Stats Display Component
 */
const StatsDisplay = ({ stats, autoRefresh }) => (
    <div className={styles.statsContainer}>
        <div className={styles.statItem}>
            <span className={styles.statValue}>{stats.total}</span>
            <span className={styles.statLabel}>Total</span>
        </div>
        <div className={styles.statItem}>
            <span className={styles.statValue} style={{ color: '#f59e0b' }}>{stats.pending}</span>
            <span className={styles.statLabel}>Pending</span>
        </div>
        <div className={styles.statItem}>
            <span className={styles.statValue} style={{ color: '#059669' }}>{stats.approved}</span>
            <span className={styles.statLabel}>Approved</span>
        </div>
        <div className={styles.statItem}>
            <span className={styles.statValue} style={{ color: '#dc2626' }}>{stats.rejected}</span>
            <span className={styles.statLabel}>Rejected</span>
        </div>
        <div className={styles.statItem}>
            <span className={styles.statValue} style={{ color: '#7c3aed' }}>{stats.flagged}</span>
            <span className={styles.statLabel}>Flagged</span>
        </div>
        {autoRefresh && (
            <div className={styles.autoRefreshIndicator}>
                <div className={styles.pulsingDot} />
                Auto-refresh
            </div>
        )}
    </div>
);
/**
 * Loading State Component
 */
const LoadingState = () => (
    <div className={styles.centerContainer}>
        <div className={styles.loadingSpinner}>
            {getIcon('Loader2', { className: styles.spinningLarge })}
        </div>
        <p className={styles.loadingText}>Loading reviews...</p>
    </div>
);
/**
 * Error State Component
 */
const ErrorState = ({ error, onRetry, onClear }) => (
    <div className={styles.centerContainer}>
        <div className={styles.errorContainer}>
            {getIcon('AlertTriangle', { className: styles.errorIcon })}
            <h3 className={styles.errorTitle}>Unable to load reviews</h3>
            <p className={styles.errorMessage}>{error}</p>
            <div className={styles.errorActions}>
                <button onClick={onRetry} className={styles.retryButton}>
                    {getIcon('RefreshCw', { className: styles.iconSmall })}
                    Try Again
                </button>
                <button onClick={onClear} className={styles.dismissButton}>
                    Dismiss
                </button>
            </div>
        </div>
    </div>
);
/**
 * Empty State Component
 */
const EmptyState = ({ hasFilters, onClearFilters }) => (
    <div className={styles.centerContainer}>
        <div className={styles.emptyContainer}>
            {getIcon('MessageCircle', { className: styles.emptyIcon })}
            <h3 className={styles.emptyTitle}>
                {hasFilters ? 'No reviews match your filters' : 'No reviews found'}
            </h3>
            <p className={styles.emptyMessage}>
                {hasFilters 
                    ? 'Try adjusting your search criteria to see more results.'
                    : 'Reviews will appear here once they are submitted.'
                }
            </p>
            {hasFilters && (
                <button onClick={onClearFilters} className={styles.clearFiltersButton}>
                    Clear Filters
                </button>
            )}
        </div>
    </div>
);
/**
 * Login Required Component
 */
const LoginRequired = () => (
    <div className={styles.centerContainer}>
        <div className={styles.loginContainer}>
            {getIcon('User', { className: styles.loginIcon })}
            <h3 className={styles.loginTitle}>Login Required</h3>
            <p className={styles.loginMessage}>
                Please log in to view and manage reviews.
            </p>
        </div>
    </div>
);
/**
 * Main Reviews Page Component
 */
const ReviewsPage = () => {
    const { user, isLoading: authLoading, authChecked } = useAuth();
    const {
        reviews,
        loading,
        error,
        stats,
        actionLoading,
        searchTerm,
        setSearchTerm,
        statusFilter,
        setStatusFilter,
        flaggedFilter,
        setFlaggedFilter,
        handleApproveReview,
        handleRejectReview,
        handleEditReview,
        refresh,
        clearError,
        autoRefresh
    } = useReviews();
    // Check if user is authenticated
    if (authLoading || !authChecked) {
        return <LoadingState />;
    }
    if (!user) {
        return <LoginRequired />;
    }
    // Handle error state
    if (error) {
        return (
            <ErrorState 
                error={error} 
                onRetry={refresh} 
                onClear={clearError} 
            />
        );
    }
    // Handle loading state
    if (loading) {
        return <LoadingState />;
    }
    // Check if filters are applied
    const hasFilters = searchTerm.trim() || statusFilter !== 'all' || flaggedFilter !== 'all';
    // Clear all filters
    const clearAllFilters = () => {
        setSearchTerm('');
        setStatusFilter('all');
        setFlaggedFilter('all');
    };
    // Handle empty state
    if (reviews.length === 0) {
        return (
            <div className={styles.container}>
                <div className={styles.header}>
                    <h1 className={styles.title}>
                        {getIcon('MessageCircle', { className: styles.titleIcon })}
                        Reviews Management
                    </h1>
                    <StatsDisplay stats={stats} autoRefresh={autoRefresh} />
                </div>
                <FilterBar
                    searchTerm={searchTerm}
                    setSearchTerm={setSearchTerm}
                    statusFilter={statusFilter}
                    setStatusFilter={setStatusFilter}
                    flaggedFilter={flaggedFilter}
                    setFlaggedFilter={setFlaggedFilter}
                    onRefresh={refresh}
                    loading={loading}
                />
                <EmptyState hasFilters={hasFilters} onClearFilters={clearAllFilters} />
            </div>
        );
    }
    return (
        <div className={styles.container}>
            <div className={styles.header}>
                <h1 className={styles.title}>
                    {getIcon('MessageCircle', { className: styles.titleIcon })}
                    Reviews Management
                </h1>
                <StatsDisplay stats={stats} autoRefresh={autoRefresh} />
            </div>
            <FilterBar
                searchTerm={searchTerm}
                setSearchTerm={setSearchTerm}
                statusFilter={statusFilter}
                setStatusFilter={setStatusFilter}
                flaggedFilter={flaggedFilter}
                setFlaggedFilter={setFlaggedFilter}
                onRefresh={refresh}
                loading={loading}
            />
            <div className={styles.reviewsList}>
                {reviews.map((review) => (
                    <ReviewItem
                        key={review.id}
                        review={review}
                        onApprove={handleApproveReview}
                        onReject={handleRejectReview}
                        onEdit={handleEditReview}
                        actionLoading={actionLoading}
                        isUserReview={user && (user.id === review.senderId || user.userId === review.senderId)}
                    />
                ))}
            </div>
        </div>
    );
};
export default ReviewsPage;