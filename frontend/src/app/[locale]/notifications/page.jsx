"use client";
import React, { useState, useMemo, memo } from 'react';
import { 
    Bell, 
    Search, 
    Filter, 
    RefreshCw, 
    Check, 
    CheckCheck, 
    X, 
    AlertCircle, 
    Clock,
    MessageCircle,
    CreditCard,
    Package,
    Store,
    Settings,
    Tag,
    Shield
} from '@/icons';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import { useAuth } from '../../../context/AuthContext';
import { useNotifications } from '../../../hooks/useNotifications';
import { getNotificationTypes, markAllAlertsAsRead } from '../../../api/client/notificationsApi';
import styles from './page.module.css';

dayjs.extend(relativeTime);

// Icon mapping for notification types
const iconMap = {
    MessageCircle,
    CreditCard,
    Package,
    Store,
    Settings,
    Tag,
    Shield,
    Clock
};

/**
 * NotificationItem Component
 * Individual notification item with actions
 */
const NotificationItem = memo(({ notification, onMarkAsRead, onRemove }) => {
    const { id, type, message, payload, isRead, createdAt } = notification;
    const typeConfig = getNotificationTypes().find(t => t.value === type) || 
                      { label: type, icon: 'Bell' };
    const timeAgo = createdAt ? dayjs(createdAt).fromNow() : '';
    
    // Get the icon component
    const IconComponent = iconMap[typeConfig.icon] || Bell;
    
    const handleMarkAsRead = () => {
        if (!isRead) {
            onMarkAsRead(id);
        }
    };
    
    const handleRemove = () => {
        onRemove(id);
    };
    
    return (
        <div className={`${styles.notificationItem} ${isRead ? styles.read : styles.unread}`}>
            <div 
                className={styles.statusBorder}
                style={{ backgroundColor: isRead ? '#94a3b8' : '#3b82f6' }}
            />
            
            <div className={styles.notificationIcon}>
                <IconComponent className={styles.typeIcon} />
            </div>
            
            <div className={styles.notificationContent}>
                <div className={styles.notificationHeader}>
                    <span className={styles.notificationType}>{typeConfig.label}</span>
                    {timeAgo && (
                        <span className={styles.notificationTime}>
                            <Clock className={styles.timeIcon} />
                            {timeAgo}
                        </span>
                    )}
                </div>
                <p className={styles.notificationMessage}>{message}</p>
                {payload && Object.keys(payload).length > 0 && (
                    <div className={styles.notificationPayload}>
                        {Object.entries(payload).map(([key, value]) => (
                            <span key={key} className={styles.payloadItem}>
                                <strong>{key}:</strong> {value}
                            </span>
                        ))}
                    </div>
                )}
            </div>
            
            <div className={styles.notificationActions}>
                {!isRead && (
                    <button
                        onClick={handleMarkAsRead}
                        className={`${styles.actionButton} ${styles.markReadButton}`}
                        title="Mark as read"
                        aria-label="Mark as read"
                    >
                        <Check className={styles.actionIcon} />
                    </button>
                )}
                <button
                    onClick={handleRemove}
                    className={`${styles.actionButton} ${styles.removeButton}`}
                    title="Remove notification"
                    aria-label="Remove notification"
                >
                    <X className={styles.actionIcon} />
                </button>
            </div>
            
            {!isRead && <div className={styles.unreadIndicator} />}
        </div>
    );
});

NotificationItem.displayName = 'NotificationItem';

/**
 * FilterBar Component
 * Filtering and search controls
 */
const FilterBar = memo(({ 
    filters, 
    onUpdateFilters, 
    onClearFilters, 
    searchTerm, 
    onSearchChange,
    unreadCount 
}) => {
    const notificationTypes = getNotificationTypes();
    
    return (
        <div className={styles.filterBar}>
            <div className={styles.searchContainer}>
                <div className={styles.searchInputWrapper}>
                    <Search className={styles.searchIcon} />
                    <input
                        type="text"
                        placeholder="Search notifications..."
                        value={searchTerm}
                        onChange={(e) => onSearchChange(e.target.value)}
                        className={styles.searchInput}
                    />
                </div>
            </div>
            
            <div className={styles.filtersContainer}>
                <div className={styles.filterGroup}>
                    <label className={styles.filterLabel}>Type:</label>
                    <select
                        value={filters.type || ''}
                        onChange={(e) => onUpdateFilters({ type: e.target.value || undefined })}
                        className={styles.filterSelect}
                    >
                        <option value="">All Types</option>
                        {notificationTypes.map(type => (
                            <option key={type.value} value={type.value}>
                                {type.label}
                            </option>
                        ))}
                    </select>
                </div>
                
                <div className={styles.filterGroup}>
                    <label className={styles.filterLabel}>Status:</label>
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
                        <option value="">All</option>
                        <option value="false">Unread ({unreadCount})</option>
                        <option value="true">Read</option>
                    </select>
                </div>
                
                <button
                    onClick={onClearFilters}
                    className={styles.clearButton}
                    title="Clear filters"
                >
                    <Filter className={styles.iconSmall} />
                    Clear
                </button>
            </div>
        </div>
    );
});

FilterBar.displayName = 'FilterBar';

/**
 * EmptyState Component
 */
const EmptyState = ({ hasFilters, onClearFilters }) => (
    <div className={styles.centerContainer}>
        <div className={styles.emptyContainer}>
            <Bell className={styles.emptyIcon} />
            <h3 className={styles.emptyTitle}>
                {hasFilters ? 'No notifications match your filters' : 'No notifications'}
            </h3>
            <p className={styles.emptyMessage}>
                {hasFilters 
                    ? "Try adjusting your search criteria to see more results." 
                    : "You're all caught up! Check back later for new updates."}
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
 * ErrorState Component
 */
const ErrorState = ({ error, onRetry }) => (
    <div className={styles.centerContainer}>
        <div className={styles.errorContainer}>
            <AlertCircle className={styles.errorIcon} />
            <h3 className={styles.errorTitle}>Unable to load notifications</h3>
            <p className={styles.errorMessage}>{error}</p>
            <button onClick={onRetry} className={styles.retryButton}>
                <RefreshCw className={styles.actionIcon} />
                Try Again
            </button>
        </div>
    </div>
);

/**
 * Login Required Component
 */
const LoginRequired = () => (
    <div className={styles.centerContainer}>
        <div className={styles.loginContainer}>
            <Bell className={styles.loginIcon} />
            <h3 className={styles.loginTitle}>Login Required</h3>
            <p className={styles.loginMessage}>
                Please log in to view your notifications.
            </p>
        </div>
    </div>
);

/**
 * Main Notifications Page Component
 */
export default function NotificationsPage() {
    const { user, isLoading: authLoading, authChecked } = useAuth();
    const userId = user?.userId || user?.id;
    
    // Local state for search
    const [searchTerm, setSearchTerm] = useState('');
    
    // Use notifications hook without auto-refresh to avoid unnecessary loading
    const {
        notifications,
        unreadCount,
        loading,
        error,
        filters,
        hasNotifications,
        isEmpty,
        updateFilters,
        clearFilters,
        markAsRead,
        markAllAsRead,
        removeNotification,
        refresh
    } = useNotifications(userId, {
        autoRefresh: false // Disable auto-refresh to prevent constant loading
    });
    
    // Filter notifications based on search term
    const filteredNotifications = useMemo(() => {
        if (!searchTerm) return notifications;
        
        const lowerSearchTerm = searchTerm.toLowerCase();
        return notifications.filter(notification => 
            notification.message?.toLowerCase().includes(lowerSearchTerm) ||
            notification.type?.toLowerCase().includes(lowerSearchTerm) ||
            Object.values(notification.payload || {}).some(value => 
                String(value).toLowerCase().includes(lowerSearchTerm)
            )
        );
    }, [notifications, searchTerm]);
    
    // Check if any filters are active
    const hasActiveFilters = Object.keys(filters).length > 0 || searchTerm;
    
    // Handle mark all as read
    const handleMarkAllAsRead = async () => {
        try {
            await markAllAlertsAsRead();
            markAllAsRead(); // Update local state
            refresh(); // Refresh the list
        } catch (error) {
        // Handle error silently for better UX
        if (process.env.NODE_ENV === 'development') {
            console.error('Event handler error:', error);
        }
    }
    };
    
    // Show loading only on initial load
    if (authLoading || !authChecked) {
        return null; // Return nothing while auth is loading
    }
    
    // Show login required state
    if (!user) {
        return <LoginRequired />;
    }
    
    // Show error state only if there's an error and no data
    if (error && !hasNotifications) {
        return <ErrorState error={error.message || 'Failed to load notifications'} onRetry={refresh} />;
    }
    
    return (
        <div className={styles.container}>
            <div className={styles.header}>
                <h1 className={styles.title}>
                    <Bell className={styles.titleIcon} />
                    Notifications
                    {unreadCount > 0 && (
                        <span className={styles.unreadBadge}>{unreadCount}</span>
                    )}
                </h1>
                
                <div className={styles.headerActions}>
                    <button 
                        onClick={refresh} 
                        className={styles.refreshButton}
                        disabled={loading}
                        title="Refresh notifications"
                    >
                        <RefreshCw className={`${styles.iconSmall} ${loading ? styles.spinning : ''}`} />
                    </button>
                    
                    {unreadCount > 0 && (
                        <button 
                            onClick={handleMarkAllAsRead}
                            className={styles.markAllButton}
                            title="Mark all as read"
                        >
                            <CheckCheck className={styles.iconSmall} />
                            Mark all read
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
            
            {/* Show error banner if there's an error but still have data */}
            {error && hasNotifications && (
                <div className={styles.errorBanner}>
                    <AlertCircle className={styles.errorIcon} />
                    Unable to refresh notifications
                </div>
            )}
            
            <div className={styles.notificationsList}>
                {isEmpty || filteredNotifications.length === 0 ? (
                    <EmptyState 
                        hasFilters={hasActiveFilters} 
                        onClearFilters={() => {
                            clearFilters();
                            setSearchTerm('');
                        }} 
                    />
                ) : (
                    filteredNotifications.map(notification => (
                        <NotificationItem
                            key={notification.id}
                            notification={notification}
                            onMarkAsRead={markAsRead}
                            onRemove={removeNotification}
                        />
                    ))
                )}
            </div>
        </div>
    );
}