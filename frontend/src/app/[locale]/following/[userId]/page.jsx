"use client";
import React, { useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { 
    UserCheck, CheckCircle, XCircle, Clock, 
    Search, RefreshCw, User, ArrowLeft, AlertTriangle,
    Loader2, UserMinus, MessageCircle
} from '@/icons';
import useFollowing from '../../../../hooks/useFollowing';
import { useAuth } from '../../../../context/AuthContext';
import { getFollowStatusTypes, formatDate } from '../../../../api/followingApi';
import styles from './page.module.css';
// Icon mapping for dynamic icon rendering
const iconMap = {
    UserCheck, CheckCircle, XCircle, Clock, Search, RefreshCw, User, ArrowLeft, 
    AlertTriangle, Loader2, UserMinus, MessageCircle
};
// Get icon component by name
const getIcon = (iconName, props = {}) => {
    const IconComponent = iconMap[iconName];
    return IconComponent ? <IconComponent {...props} /> : null;
};
/**
 * Individual Following Item Component
 */
const FollowingItem = ({ following, actionLoading, currentUserId }) => {
    const statusTypes = getFollowStatusTypes();
    const statusInfo = statusTypes.find(s => s.value === following.followStatus) || statusTypes[0];
    const isActionLoading = actionLoading[following.id] || actionLoading[following.followedUserId];
    return (
        <div className={styles.followingItem}>
            {/* Status indicator */}
            <div 
                className={styles.statusBorder}
                style={{ backgroundColor: statusInfo.color }}
            />
            <div className={styles.followingContent}>
                <div className={styles.followingHeader}>
                    <div className={styles.followingInfo}>
                        <div className={styles.followingAvatar}>
                            {following.followedUserAvatar ? (
                                <img 
                                    src={following.followedUserAvatar} 
                                    alt={following.followedUserName || 'User'} 
                                    className={styles.avatarImage}
                                />
                            ) : (
                                getIcon('User', { className: styles.avatarIcon })
                            )}
                        </div>
                        <div className={styles.followingMeta}>
                            <h3 className={styles.followingName}>
                                {following.followedUserName || following.followedUserId || 'Anonymous User'}
                            </h3>
                            {following.followedUserBio && (
                                <p className={styles.followingBio}>{following.followedUserBio}</p>
                            )}
                            <div className={styles.followingDate}>
                                {following.followStatus === 'approved' ? 'Following since' : 'Requested'} {formatDate(following.createdAt)}
                            </div>
                        </div>
                    </div>
                    <div className={styles.statusBadge} style={{ color: statusInfo.color }}>
                        {getIcon(statusInfo.icon, { className: styles.iconSmall })}
                        {statusInfo.label}
                    </div>
                </div>
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
    onRefresh,
    loading 
}) => {
    const statusTypes = getFollowStatusTypes();
    return (
        <div className={styles.filterBar}>
            <div className={styles.searchContainer}>
                <div className={styles.searchInputWrapper}>
                    {getIcon('Search', { className: styles.searchIcon })}
                    <input
                        type="text"
                        placeholder="Search following..."
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
                <button
                    onClick={onRefresh}
                    disabled={loading}
                    className={styles.refreshButton}
                    title="Refresh following"
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
            <span className={styles.statLabel}>Following</span>
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
        <p className={styles.loadingText}>Loading following...</p>
    </div>
);
/**
 * Error State Component
 */
const ErrorState = ({ error, onRetry, onClear }) => (
    <div className={styles.centerContainer}>
        <div className={styles.errorContainer}>
            {getIcon('AlertTriangle', { className: styles.errorIcon })}
            <h3 className={styles.errorTitle}>Unable to load following</h3>
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
const EmptyState = ({ hasFilters, onClearFilters, isOwner }) => (
    <div className={styles.centerContainer}>
        <div className={styles.emptyContainer}>
            {getIcon('UserCheck', { className: styles.emptyIcon })}
            <h3 className={styles.emptyTitle}>
                {hasFilters ? 'No following match your filters' : isOwner ? 'Not following anyone yet' : 'This user is not following anyone'}
            </h3>
            <p className={styles.emptyMessage}>
                {hasFilters 
                    ? 'Try adjusting your search criteria to see more results.'
                    : isOwner 
                        ? 'Start following interesting users to see them here!'
                        : 'This user has not followed anyone yet.'
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
                Please log in to view following.
            </p>
        </div>
    </div>
);
/**
 * Main Following Page Component
 */
const FollowingPage = () => {
    const params = useParams();
    const router = useRouter();
    const { user, isLoading: authLoading, authChecked } = useAuth();
    const userId = params.userId;
    const {
        followers: following,
        loading,
        error,
        stats,
        actionLoading,
        searchTerm,
        setSearchTerm,
        statusFilter,
        setStatusFilter,
        refresh,
        clearError,
        autoRefresh
    } = useFollowing(userId, 'following');
    // Check if current user is the owner of this following page
    const isOwner = user && (user.id === userId || user.userId === userId);
    const currentUserId = user?.id || user?.userId;
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
    const hasFilters = searchTerm.trim() || statusFilter !== 'all';
    // Clear all filters
    const clearAllFilters = () => {
        setSearchTerm('');
        setStatusFilter('all');
    };
    // Handle empty state
    if (following.length === 0) {
        return (
            <div className={styles.container}>
                <div className={styles.header}>
                    <div className={styles.headerTop}>
                        <button
                            onClick={() => router.back()}
                            className={styles.backButton}
                            title="Go back"
                        >
                            {getIcon('ArrowLeft', { className: styles.iconSmall })}
                            Back
                        </button>
                        <h1 className={styles.title}>
                            {getIcon('UserCheck', { className: styles.titleIcon })}
                            {isOwner ? 'Who You Follow' : 'Following'}
                        </h1>
                    </div>
                    <StatsDisplay stats={stats} autoRefresh={autoRefresh} />
                </div>
                <FilterBar
                    searchTerm={searchTerm}
                    setSearchTerm={setSearchTerm}
                    statusFilter={statusFilter}
                    setStatusFilter={setStatusFilter}
                    onRefresh={refresh}
                    loading={loading}
                />
                <EmptyState hasFilters={hasFilters} onClearFilters={clearAllFilters} isOwner={isOwner} />
            </div>
        );
    }
    return (
        <div className={styles.container}>
            <div className={styles.header}>
                <div className={styles.headerTop}>
                    <button
                        onClick={() => router.back()}
                        className={styles.backButton}
                        title="Go back"
                    >
                        {getIcon('ArrowLeft', { className: styles.iconSmall })}
                        Back
                    </button>
                    <h1 className={styles.title}>
                        {getIcon('UserCheck', { className: styles.titleIcon })}
                        {isOwner ? 'Who You Follow' : 'Following'}
                    </h1>
                </div>
                <StatsDisplay stats={stats} autoRefresh={autoRefresh} />
            </div>
            <FilterBar
                searchTerm={searchTerm}
                setSearchTerm={setSearchTerm}
                statusFilter={statusFilter}
                setStatusFilter={setStatusFilter}
                onRefresh={refresh}
                loading={loading}
            />
            <div className={styles.followingList}>
                {following.map((followingItem) => (
                    <FollowingItem
                        key={followingItem.id}
                        following={followingItem}
                        actionLoading={actionLoading}
                        currentUserId={currentUserId}
                    />
                ))}
            </div>
        </div>
    );
};
export default FollowingPage; 