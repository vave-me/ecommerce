"use client";
export const dynamic = 'force-dynamic';
import React, { useState, useMemo, memo } from 'react';
import { 
    Activity, 
    Search, 
    Filter, 
    RefreshCw, 
    Check, 
    CheckCheck, 
    X, 
    AlertCircle, 
    Clock,
    Heart,
    ThumbsUp,
    ThumbsDown,
    Eye,
    ShoppingCart,
    Bookmark,
    MessageSquare,
    Share,
    UserPlus,
    Star,
    Package,
    Store,
    User,
    Folder,
    Tag
} from '@/icons';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import { useAuth } from '../../../context/AuthContext';
import { useActivity } from '../../../hooks/useActivity';
import { 
    getActivityTypes, 
    getActionTypes, 
    getActionTypeConfig, 
    getItemTypeConfig,
    formatActivityMessage 
} from '../../../api/activityApi';
import styles from './page.module.css';
import { useTranslations } from "next-intl";
import ActivityDebug from '../../../components/Activity/ActivityDebug';
dayjs.extend(relativeTime);
// Icon mapping for activity types
const iconMap = {
    Heart,
    ThumbsUp,
    ThumbsDown,
    Eye,
    ShoppingCart,
    Bookmark,
    MessageSquare,
    Share,
    UserPlus,
    Star,
    Package,
    Store,
    User,
    Folder,
    Tag,
    Activity,
    Clock
};
/**
 * ActivityItem Component
 * Individual activity item with actions (matches NotificationItem structure)
 */
const ActivityItem = memo(({ activity, onMarkAsRead, onRemove }) => {
    const { id, actionType, itemType, itemId, target, message, createdAt, isRead } = activity;
    const actionConfig = getActionTypeConfig(actionType);
    const itemConfig = getItemTypeConfig(itemType);
    const formattedMessage = formatActivityMessage(activity);
    const timeAgo = createdAt ? dayjs(createdAt).fromNow() : '';
    // Get the icon component
    const IconComponent = iconMap[actionConfig.icon] || Activity;
    const handleMarkAsRead = () => {
        if (!isRead) {
            onMarkAsRead(id);
        }
    };
    const handleRemove = () => {
        onRemove(id);
    };
    return (
        <div className={`${styles.activityItem} ${isRead ? styles.read : styles.unread}`}>
            <div className={styles.activityIcon}>
                <IconComponent size={16} />
            </div>
            <div className={styles.activityContent}>
                <div className={styles.activityHeader}>
                    <span className={styles.activityType}>{actionConfig.label}</span>
                    {timeAgo && (
                        <span className={styles.activityTime}>
                            <Clock size={12} />
                            {timeAgo}
                        </span>
                    )}
                </div>
                <p className={styles.activityMessage}>{formattedMessage}</p>
                {target && (
                    <div className={styles.activityTarget}>
                        <span className={styles.targetItem}>
                            <strong>{itemConfig.label}:</strong> {target.itemName || itemId}
                        </span>
                    </div>
                )}
            </div>
            <div className={styles.activityActions}>
                {!isRead && (
                    <button
                        onClick={handleMarkAsRead}
                        className={styles.actionButton}
                        title="Mark as read"
                        aria-label="Mark as read"
                    >
                        <Check size={16} />
                    </button>
                )}
                <button
                    onClick={handleRemove}
                    className={`${styles.actionButton} ${styles.removeButton}`}
                    title="Remove activity"
                    aria-label="Remove activity"
                >
                    <X size={16} />
                </button>
            </div>
            {!isRead && <div className={styles.unreadIndicator} />}
        </div>
    );
});
ActivityItem.displayName = 'ActivityItem';
/**
 * FilterBar Component
 * Filtering and search controls (matches notifications pattern)
 */
const FilterBar = memo(({ 
    filters, 
    onUpdateFilters, 
    onClearFilters, 
    searchTerm, 
    onSearchChange,
    unreadCount 
}) => {
    const actionTypes = getActionTypes();
    const activityTypes = getActivityTypes();
    return (
        <div className={styles.filterBar}>
            <div className={styles.searchContainer}>
                <Search size={20} className={styles.searchIcon} />
                <input
                    type="text"
                    placeholder="Search activities..."
                    value={searchTerm}
                    onChange={(e) => onSearchChange(e.target.value)}
                    className={styles.searchInput}
                />
            </div>
            <div className={styles.filterControls}>
                <select
                    value={filters.actionType || ''}
                    onChange={(e) => onUpdateFilters({ actionType: e.target.value || undefined })}
                    className={styles.filterSelect}
                >
                    <option value="">All Actions</option>
                    {actionTypes.map(type => (
                        <option key={type.value} value={type.value}>
                            {type.label}
                        </option>
                    ))}
                </select>
                <select
                    value={filters.type || ''}
                    onChange={(e) => onUpdateFilters({ type: e.target.value || undefined })}
                    className={styles.filterSelect}
                >
                    <option value="">All Types</option>
                    {activityTypes.map(type => (
                        <option key={type.value} value={type.value}>
                            {type.label}
                        </option>
                    ))}
                </select>
                <select
                    value={filters.isRead?.toString() || ''}
                    onChange={(e) => {
                        const value = e.target.value;
                        onUpdateFilters({ 
                            isRead: value === '' ? undefined : value === 'true' 
                        });
                    }}
                    className={styles.filterSelect}
                >
                    <option value="">All Status</option>
                    <option value="false">Unread ({unreadCount})</option>
                    <option value="true">Read</option>
                </select>
                <button
                    onClick={onClearFilters}
                    className={styles.clearFiltersButton}
                    title="Clear filters"
                >
                    <Filter size={16} />
                    Clear
                </button>
            </div>
        </div>
    );
});
FilterBar.displayName = 'FilterBar';
/**
 * EmptyState Component
 * Shown when no activities are found
 */
const EmptyState = ({ hasFilters, onClearFilters }) => (
    <div className={styles.emptyState}>
        <Activity size={48} className={styles.emptyIcon} />
        <h3>No activities found</h3>
        <p>
            {hasFilters 
                ? "No activities match your current filters." 
                : "You haven't performed any activities yet. Start exploring!"}
        </p>
        {hasFilters && (
            <button onClick={onClearFilters} className={styles.clearFiltersButton}>
                Clear filters
            </button>
        )}
    </div>
);
/**
 * ErrorState Component
 * Shown when there's an error loading activities
 */
const ErrorState = ({ error, onRetry }) => (
    <div className={styles.errorState}>
        <AlertCircle size={48} className={styles.errorIcon} />
        <h3>Failed to load activities</h3>
        <p>{error?.message || 'Something went wrong while loading your activities.'}</p>
        <button onClick={onRetry} className={styles.retryButton}>
            <RefreshCw size={16} />
            Try again
        </button>
    </div>
);
/**
 * LoadingState Component
 * Shown while activities are loading
 */
const LoadingState = () => (
    <div className={styles.loadingState}>
        <div className={styles.loadingSpinner}>
            <RefreshCw size={24} className={styles.spinIcon} />
        </div>
        <p>Loading activities...</p>
    </div>
);
/**
 * Main Activity Page Component
 */
export default function ActivityPage() {
    const t = useTranslations('ActivityPage');
    const { user } = useAuth();
    const userId = user?.userId || user?.id;
    // Local state for search
    const [searchTerm, setSearchTerm] = useState('');
    // Use activity hook with auto-refresh
    const {
        activities,
        unreadCount,
        loading,
        error,
        filters,
        hasActivities,
        isEmpty,
        updateFilters,
        clearFilters,
        markAsRead,
        markAllAsRead,
        archiveActivityItem,
        refresh
    } = useActivity(userId, {
        autoRefresh: true,
        refreshInterval: 30000 // 30 seconds
    });
    // Filter activities based on search term
    const filteredActivities = useMemo(() => {
        if (!searchTerm) return activities;
        const lowerSearchTerm = searchTerm.toLowerCase();
        return activities.filter(activity => 
            formatActivityMessage(activity).toLowerCase().includes(lowerSearchTerm) ||
            activity.actionType?.toLowerCase().includes(lowerSearchTerm) ||
            activity.itemType?.toLowerCase().includes(lowerSearchTerm) ||
            activity.target?.itemName?.toLowerCase().includes(lowerSearchTerm)
        );
    }, [activities, searchTerm]);
    // Check if any filters are active
    const hasActiveFilters = Object.keys(filters).length > 0 || searchTerm;
    // Handle mark all as read
    const handleMarkAllAsRead = () => {
        markAllAsRead();
    };
    // Handle remove activity
    const handleRemoveActivity = async (activityId) => {
        try {
            await archiveActivityItem(activityId, 'User removed activity');
        } catch (err) {
        // Handle error silently for better UX
        if (process.env.NODE_ENV === 'development') {
            console.error('Event handler error:', err);
        }
    }
    };
    // Show login required state
    if (!user) {
        return (
            <div className={styles.container}>
                <div className={styles.loginRequired}>
                    <Activity size={48} />
                    <h2>Login Required</h2>
                    <p>Please log in to view your activity history.</p>
                </div>
            </div>
        );
    }
    return (
        <div className={styles.container}>
            {/* Debug component (development only) */}
            <ActivityDebug />
            <div className={styles.header}>
                <div className={styles.titleSection}>
                    <h1 className={styles.title}>
                        <Activity size={28} />
                        {t("title")}
                        {unreadCount > 0 && (
                            <span className={styles.unreadBadge}>{unreadCount}</span>
                        )}
                    </h1>
                    <p className={styles.subtitle}>
                        Track your interactions and engagement history
                    </p>
                </div>
                <div className={styles.headerActions}>
                    <button 
                        onClick={refresh} 
                        className={styles.refreshButton}
                        disabled={loading}
                        title="Refresh activities"
                    >
                        <RefreshCw size={16} className={loading ? styles.spinning : ''} />
                        Refresh
                    </button>
                    {unreadCount > 0 && (
                        <button 
                            onClick={handleMarkAllAsRead}
                            className={styles.markAllButton}
                            title="Mark all as read"
                        >
                            <CheckCheck size={16} />
                            {t("markAllAsRead")}
                        </button>
                    )}
                </div>
            </div>
            <FilterBar
                filters={filters}
                onUpdateFilters={updateFilters}
                onClearFilters={clearFilters}
                searchTerm={searchTerm}
                onSearchChange={setSearchTerm}
                unreadCount={unreadCount}
            />
            <div className={styles.content}>
                {loading && <LoadingState />}
                {error && (
                    <ErrorState error={error} onRetry={refresh} />
                )}
                {!loading && !error && isEmpty && (
                    <EmptyState 
                        hasFilters={hasActiveFilters} 
                        onClearFilters={() => {
                            clearFilters();
                            setSearchTerm('');
                        }} 
                    />
                )}
                {!loading && !error && hasActivities && (
                    <div className={styles.activitiesList}>
                        {filteredActivities.length === 0 ? (
                            <EmptyState 
                                hasFilters={hasActiveFilters} 
                                onClearFilters={() => {
                                    clearFilters();
                                    setSearchTerm('');
                                }} 
                            />
                        ) : (
                            filteredActivities.map(activity => (
                                <ActivityItem
                                    key={activity.id}
                                    activity={activity}
                                    onMarkAsRead={markAsRead}
                                    onRemove={handleRemoveActivity}
                                />
                            ))
                        )}
                    </div>
                )}
            </div>
        </div>
    );
}
