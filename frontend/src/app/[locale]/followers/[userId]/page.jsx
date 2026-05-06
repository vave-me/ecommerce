"use client";
import React, { useState, useEffect, useCallback, useMemo, memo } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { getFollowers } from '@/api/followingApi.jsx';
import { useAuth } from '@/context/AuthContext';
import {
    Users, User, UserPlus, UserMinus, MessageCircle, Eye, EyeOff,
    Filter, Loader2, AlertCircle, RefreshCw, ChevronDown,
    Calendar, MapPin, Verified, Star, Bell, BellOff
} from 'lucide-react';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import EmptyPlaceholder from '@/components/Utils/EmptyPlaceholder';
import styles from './Followers.module.css';
import { 
    CheckCircle, XCircle, Clock, Search, ArrowLeft, AlertTriangle
} from '@/icons';
import useFollowing from '../../../../hooks/useFollowing';
import { getFollowStatusTypes, formatDate } from '../../../../api/followingApi';
// Icon mapping for dynamic icon rendering
const iconMap = {
    Users, CheckCircle, XCircle, Clock, Search, RefreshCw, User, ArrowLeft, 
    AlertTriangle, Loader2, UserPlus, UserMinus, MessageCircle, Eye, EyeOff,
    Filter, AlertCircle, ChevronDown, Calendar, MapPin, Verified, Star, Bell, BellOff
};
// Get icon component by name
const getIcon = (iconName, props = {}) => {
    const IconComponent = iconMap[iconName];
    return IconComponent ? <IconComponent {...props} /> : null;
};
/**
 * Individual Follower Item Component
 */
const FollowerItem = memo(({ follower, onApprove, onReject, onFollow, actionLoading, currentUserId, isOwner = false }) => {
    const statusTypes = getFollowStatusTypes();
    const statusInfo = statusTypes.find(s => s.value === follower.followStatus) || statusTypes[0];
    const isActionLoading = actionLoading[follower.id] || actionLoading[follower.userId];
    return (
        <div className={styles.followerItem}>
            {/* Status indicator */}
            <div 
                className={styles.statusBorder}
                style={{ backgroundColor: statusInfo.color }}
            />
            <div className={styles.followerContent}>
                <div className={styles.followerHeader}>
                    <div className={styles.followerInfo}>
                        <div className={styles.followerAvatar}>
                            {follower.followerAvatar ? (
                                <img 
                                    src={follower.followerAvatar} 
                                    alt={follower.followerName || 'User'} 
                                    className={styles.avatarImage}
                                />
                            ) : (
                                getIcon('User', { className: styles.avatarIcon })
                            )}
                        </div>
                        <div className={styles.followerMeta}>
                            <h3 className={styles.followerName}>
                                {follower.followerName || 'Anonymous User'}
                            </h3>
                            {follower.followerBio && (
                                <p className={styles.followerBio}>{follower.followerBio}</p>
                            )}
                            <div className={styles.followerDate}>
                                {follower.followStatus === 'approved' ? 'Following since' : 'Requested'} {formatDate(follower.createdAt)}
                            </div>
                        </div>
                    </div>
                    <div className={styles.statusBadge} style={{ color: statusInfo.color }}>
                        {getIcon(statusInfo.icon, { className: styles.iconSmall })}
                        {statusInfo.label}
                    </div>
                </div>
                {/* Actions */}
                <div className={styles.followerActions}>
                    {isOwner && follower.followStatus === 'pending' && (
                        <>
                            <button
                                onClick={() => onApprove(follower.id)}
                                disabled={isActionLoading}
                                className={`${styles.actionButton} ${styles.approveButton}`}
                                title="Approve follow request"
                            >
                                {isActionLoading === 'approving' ? 
                                    getIcon('Loader2', { className: `${styles.iconSmall} ${styles.spinning}` }) : 
                                    getIcon('CheckCircle', { className: styles.iconSmall })
                                }
                                Approve
                            </button>
                            <button
                                onClick={() => onReject(follower.id)}
                                disabled={isActionLoading}
                                className={`${styles.actionButton} ${styles.rejectButton}`}
                                title="Reject follow request"
                            >
                                {isActionLoading === 'rejecting' ? 
                                    getIcon('Loader2', { className: `${styles.iconSmall} ${styles.spinning}` }) : 
                                    getIcon('XCircle', { className: styles.iconSmall })
                                }
                                Reject
                            </button>
                        </>
                    )}
                    {!isOwner && currentUserId !== follower.userId && (
                        <button
                            onClick={() => onFollow(follower.userId)}
                            disabled={isActionLoading}
                            className={`${styles.actionButton} ${styles.followButton}`}
                            title="Follow back"
                        >
                            {isActionLoading === 'following' ? 
                                getIcon('Loader2', { className: `${styles.iconSmall} ${styles.spinning}` }) : 
                                getIcon('UserPlus', { className: styles.iconSmall })
                            }
                            Follow Back
                        </button>
                    )}
                </div>
            </div>
        </div>
    );
});
FollowerItem.displayName = 'FollowerItem';
/**
 * Filter Bar Component
 */
const FilterBar = memo(({ 
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
                        placeholder="Search followers..."
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
                    title="Refresh followers"
                >
                    {getIcon('RefreshCw', { className: `${styles.iconSmall} ${loading ? styles.spinning : ''}` })}
                </button>
            </div>
        </div>
    );
});
FilterBar.displayName = 'FilterBar';
/**
 * Stats Display Component
 */
const StatsDisplay = memo(({ stats, autoRefresh }) => (
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
));
StatsDisplay.displayName = 'StatsDisplay';
/**
 * Loading State Component
 */
const LoadingState = memo(() => (
    <div className={styles.centerContainer}>
        <div className={styles.loadingSpinner}>
            {getIcon('Loader2', { className: styles.spinningLarge })}
        </div>
        <p className={styles.loadingText}>Loading followers...</p>
    </div>
));
LoadingState.displayName = 'LoadingState';
/**
 * Error State Component
 */
const ErrorState = memo(({ error, onRetry, onClear }) => (
    <div className={styles.centerContainer}>
        <div className={styles.errorContainer}>
            {getIcon('AlertTriangle', { className: styles.errorIcon })}
            <h3 className={styles.errorTitle}>Unable to load followers</h3>
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
));
ErrorState.displayName = 'ErrorState';
/**
 * Empty State Component
 */
const EmptyState = memo(({ hasFilters, onClearFilters, isOwner }) => (
    <div className={styles.centerContainer}>
        <div className={styles.emptyContainer}>
            {getIcon('Users', { className: styles.emptyIcon })}
            <h3 className={styles.emptyTitle}>
                {hasFilters ? 'No followers match your filters' : isOwner ? 'No followers yet' : 'This user has no followers'}
            </h3>
            <p className={styles.emptyMessage}>
                {hasFilters 
                    ? 'Try adjusting your search criteria to see more results.'
                    : isOwner 
                        ? 'Start sharing great content to attract followers!'
                        : 'Be the first to follow this user.'
                }
            </p>
            {hasFilters && (
                <button onClick={onClearFilters} className={styles.clearFiltersButton}>
                    Clear Filters
                </button>
            )}
        </div>
    </div>
));
EmptyState.displayName = 'EmptyState';
/**
 * Login Required Component
 */
const LoginRequired = memo(() => (
    <div className={styles.centerContainer}>
        <div className={styles.loginContainer}>
            {getIcon('User', { className: styles.loginIcon })}
            <h3 className={styles.loginTitle}>Login Required</h3>
            <p className={styles.loginMessage}>
                Please log in to view followers.
            </p>
        </div>
    </div>
));
LoginRequired.displayName = 'LoginRequired';
/**
 * Main Followers Page Component
 */
const FollowersPage = () => {
    const params = useParams();
    const router = useRouter();
    const { user, isLoading: authLoading, authChecked } = useAuth();
    const userId = params.userId;
    const {
        followers,
        loading,
        error,
        stats,
        actionLoading,
        searchTerm,
        setSearchTerm,
        statusFilter,
        setStatusFilter,
        handleApproveFollow,
        handleRejectFollow,
        handleFollowUser,
        refresh,
        clearError,
        autoRefresh
    } = useFollowing(userId, 'followers');
    // Check if current user is the owner of this followers page
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
    if (followers.length === 0) {
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
                            {getIcon('Users', { className: styles.titleIcon })}
                            {isOwner ? 'Your Followers' : 'Followers'}
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
                        {getIcon('Users', { className: styles.titleIcon })}
                        {isOwner ? 'Your Followers' : 'Followers'}
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
            <div className={styles.followersList}>
                {followers.map((follower) => (
                    <FollowerItem
                        key={follower.id}
                        follower={follower}
                        onApprove={handleApproveFollow}
                        onReject={handleRejectFollow}
                        onFollow={handleFollowUser}
                        actionLoading={actionLoading}
                        currentUserId={currentUserId}
                        isOwner={isOwner}
                    />
                ))}
            </div>
        </div>
    );
};
export default memo(FollowersPage); 