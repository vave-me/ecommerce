"use client"
export const dynamic = 'force-dynamic';
import React, { useState } from 'react';
import { 
    Mail, CheckCircle, XCircle, Clock, Pause, Play,
    Search, Settings, RefreshCw, User, AlertTriangle,
    Loader2, Edit3, Plus, Send, Eye, EyeOff
} from '@/icons';
import useNewsletters from '../../../hooks/useNewsletters';
import { useAuth } from '../../../context/AuthContext';
import { 
    getSubscriptionStatusTypes, 
    getNewsletterPreferences, 
    formatDate 
} from '../../../api/newslettersApi';
import styles from './page.module.css';
// Icon mapping for dynamic icon rendering
const iconMap = {
    Mail, CheckCircle, XCircle, Clock, Pause, Play, Search, Settings, 
    RefreshCw, User, AlertTriangle, Loader2, Edit3, Plus, Send, Eye, EyeOff
};
// Get icon component by name
const getIcon = (iconName, props = {}) => {
    const IconComponent = iconMap[iconName];
    return IconComponent ? <IconComponent {...props} /> : null;
};
/**
 * Individual Subscription Item Component
 */
const SubscriptionItem = ({ 
    subscription, 
    onUnsubscribe, 
    onToggleStatus, 
    onUpdatePreferences, 
    actionLoading,
    availableNewsletters 
}) => {
    const [isEditing, setIsEditing] = useState(false);
    const [editPreferences, setEditPreferences] = useState(subscription.subscriptionPreferences || '');
    const statusTypes = getSubscriptionStatusTypes();
    const statusInfo = statusTypes.find(s => s.value === subscription.subscriptionStatus) || statusTypes[0];
    const newsletterInfo = availableNewsletters.find(nl => nl.newsletterId === subscription.newsletterId);
    const isActionLoading = actionLoading[subscription.subscriptionId];
    const handleEditPreferences = () => {
        if (isEditing) {
            if (editPreferences.trim() && editPreferences !== subscription.subscriptionPreferences) {
                onUpdatePreferences(subscription.subscriptionId, {
                    subscriptionPreferences: editPreferences.trim()
                });
            }
            setIsEditing(false);
        } else {
            setEditPreferences(subscription.subscriptionPreferences || '');
            setIsEditing(true);
        }
    };
    const handleCancelEdit = () => {
        setEditPreferences(subscription.subscriptionPreferences || '');
        setIsEditing(false);
    };
    const preferences = getNewsletterPreferences();
    return (
        <div className={styles.subscriptionItem}>
            {/* Status indicator */}
            <div 
                className={styles.statusBorder}
                style={{ backgroundColor: statusInfo.color }}
            />
            <div className={styles.subscriptionHeader}>
                <div className={styles.subscriptionMeta}>
                    <div className={styles.subscriptionInfo}>
                        {getIcon('Mail', { className: styles.icon })}
                        <div className={styles.newsletterDetails}>
                            <h3 className={styles.newsletterName}>
                                {subscription.newsletterName || newsletterInfo?.name || 'Newsletter'}
                            </h3>
                            <p className={styles.newsletterDescription}>
                                {subscription.newsletterDescription || newsletterInfo?.description || 'No description available'}
                            </p>
                            {newsletterInfo && (
                                <div className={styles.newsletterMeta}>
                                    <span className={styles.category}>{newsletterInfo.category}</span>
                                    <span className={styles.frequency}>{newsletterInfo.frequency}</span>
                                </div>
                            )}
                        </div>
                    </div>
                    <div className={styles.subscriptionDate}>
                        Subscribed {formatDate(subscription.createdAt)}
                        {subscription.updatedAt && subscription.updatedAt !== subscription.createdAt && (
                            <span className={styles.updatedDate}>
                                Updated {formatDate(subscription.updatedAt)}
                            </span>
                        )}
                    </div>
                </div>
                <div className={styles.statusBadge} style={{ color: statusInfo.color }}>
                    {getIcon(statusInfo.icon, { className: styles.iconSmall })}
                    {statusInfo.label}
                </div>
            </div>
            {/* Preferences */}
            <div className={styles.preferencesSection}>
                <div className={styles.preferencesHeader}>
                    <span className={styles.preferencesLabel}>Preferences:</span>
                    <button
                        onClick={handleEditPreferences}
                        disabled={isActionLoading}
                        className={styles.editPreferencesButton}
                        title="Edit preferences"
                    >
                        {getIcon('Edit3', { className: styles.iconSmall })}
                    </button>
                </div>
                {isEditing ? (
                    <div className={styles.editForm}>
                        <select
                            value={editPreferences}
                            onChange={(e) => setEditPreferences(e.target.value)}
                            className={styles.preferencesSelect}
                        >
                            <option value="">Select preferences</option>
                            {preferences.map(pref => (
                                <option key={pref.value} value={pref.value}>
                                    {pref.label} - {pref.description}
                                </option>
                            ))}
                        </select>
                        <div className={styles.editActions}>
                            <button
                                onClick={handleEditPreferences}
                                disabled={!editPreferences.trim() || isActionLoading}
                                className={`${styles.actionButton} ${styles.saveButton}`}
                            >
                                {isActionLoading === 'updating' ? 
                                    getIcon('Loader2', { className: `${styles.iconSmall} ${styles.spinning}` }) : 
                                    getIcon('CheckCircle', { className: styles.iconSmall })
                                }
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
                    </div>
                ) : (
                    <div className={styles.currentPreferences}>
                        {subscription.subscriptionPreferences ? (
                            <span className={styles.preferencesValue}>
                                {preferences.find(p => p.value === subscription.subscriptionPreferences)?.label || 
                                 subscription.subscriptionPreferences}
                            </span>
                        ) : (
                            <span className={styles.noPreferences}>No preferences set</span>
                        )}
                    </div>
                )}
            </div>
            {/* Actions */}
            <div className={styles.subscriptionActions}>
                {subscription.subscriptionStatus === 'active' && (
                    <button
                        onClick={() => onToggleStatus(subscription.subscriptionId, subscription.subscriptionStatus)}
                        disabled={isActionLoading}
                        className={`${styles.actionButton} ${styles.pauseButton}`}
                        title="Pause subscription"
                    >
                        {isActionLoading === 'updating' ? 
                            getIcon('Loader2', { className: `${styles.iconSmall} ${styles.spinning}` }) : 
                            getIcon('Pause', { className: styles.iconSmall })
                        }
                        Pause
                    </button>
                )}
                {subscription.subscriptionStatus === 'paused' && (
                    <button
                        onClick={() => onToggleStatus(subscription.subscriptionId, subscription.subscriptionStatus)}
                        disabled={isActionLoading}
                        className={`${styles.actionButton} ${styles.resumeButton}`}
                        title="Resume subscription"
                    >
                        {isActionLoading === 'updating' ? 
                            getIcon('Loader2', { className: `${styles.iconSmall} ${styles.spinning}` }) : 
                            getIcon('Play', { className: styles.iconSmall })
                        }
                        Resume
                    </button>
                )}
                <button
                    onClick={() => onUnsubscribe(subscription.subscriptionId)}
                    disabled={isActionLoading}
                    className={`${styles.actionButton} ${styles.unsubscribeButton}`}
                    title="Unsubscribe"
                >
                    {isActionLoading === 'unsubscribing' ? 
                        getIcon('Loader2', { className: `${styles.iconSmall} ${styles.spinning}` }) : 
                        getIcon('XCircle', { className: styles.iconSmall })
                    }
                    Unsubscribe
                </button>
            </div>
        </div>
    );
};
/**
 * Available Newsletter Item Component
 */
const AvailableNewsletterItem = ({ newsletter, onSubscribe, actionLoading, isSubscribed }) => {
    const isActionLoading = actionLoading[newsletter.newsletterId];
    return (
        <div className={styles.availableNewsletterItem}>
            <div className={styles.newsletterIcon}>
                {getIcon('Mail', { className: styles.icon })}
            </div>
            <div className={styles.newsletterInfo}>
                <h3 className={styles.newsletterName}>{newsletter.name}</h3>
                <p className={styles.newsletterDescription}>{newsletter.description}</p>
                <div className={styles.newsletterMeta}>
                    <span className={styles.category}>{newsletter.category}</span>
                    <span className={styles.frequency}>{newsletter.frequency}</span>
                </div>
            </div>
            <div className={styles.subscribeAction}>
                {isSubscribed ? (
                    <span className={styles.subscribedIndicator}>
                        {getIcon('CheckCircle', { className: styles.iconSmall })}
                        Subscribed
                    </span>
                ) : (
                    <button
                        onClick={() => onSubscribe(newsletter.newsletterId)}
                        disabled={isActionLoading}
                        className={`${styles.actionButton} ${styles.subscribeButton}`}
                        title="Subscribe to newsletter"
                    >
                        {isActionLoading === 'subscribing' ? 
                            getIcon('Loader2', { className: `${styles.iconSmall} ${styles.spinning}` }) : 
                            getIcon('Plus', { className: styles.iconSmall })
                        }
                        Subscribe
                    </button>
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
    categoryFilter,
    setCategoryFilter,
    onRefresh,
    loading,
    availableNewsletters 
}) => {
    const statusTypes = getSubscriptionStatusTypes();
    const categories = [...new Set(availableNewsletters.map(nl => nl.category))];
    return (
        <div className={styles.filterBar}>
            <div className={styles.searchContainer}>
                <div className={styles.searchInputWrapper}>
                    {getIcon('Search', { className: styles.searchIcon })}
                    <input
                        type="text"
                        placeholder="Search newsletters..."
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
                    <label className={styles.filterLabel}>Category:</label>
                    <select
                        value={categoryFilter}
                        onChange={(e) => setCategoryFilter(e.target.value)}
                        className={styles.filterSelect}
                    >
                        <option value="all">All Categories</option>
                        {categories.map(category => (
                            <option key={category} value={category}>
                                {category}
                            </option>
                        ))}
                    </select>
                </div>
                <button
                    onClick={onRefresh}
                    disabled={loading}
                    className={styles.refreshButton}
                    title="Refresh subscriptions"
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
            <span className={styles.statValue} style={{ color: '#059669' }}>{stats.active}</span>
            <span className={styles.statLabel}>Active</span>
        </div>
        <div className={styles.statItem}>
            <span className={styles.statValue} style={{ color: '#f59e0b' }}>{stats.paused}</span>
            <span className={styles.statLabel}>Paused</span>
        </div>
        <div className={styles.statItem}>
            <span className={styles.statValue} style={{ color: '#dc2626' }}>{stats.cancelled}</span>
            <span className={styles.statLabel}>Cancelled</span>
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
        <p className={styles.loadingText}>Loading newsletters...</p>
    </div>
);
/**
 * Error State Component
 */
const ErrorState = ({ error, onRetry, onClear }) => (
    <div className={styles.centerContainer}>
        <div className={styles.errorContainer}>
            {getIcon('AlertTriangle', { className: styles.errorIcon })}
            <h3 className={styles.errorTitle}>Unable to load newsletters</h3>
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
const EmptyState = ({ hasFilters, onClearFilters, showAvailable }) => (
    <div className={styles.centerContainer}>
        <div className={styles.emptyContainer}>
            {getIcon('Mail', { className: styles.emptyIcon })}
            <h3 className={styles.emptyTitle}>
                {hasFilters ? 'No subscriptions match your filters' : 'No newsletter subscriptions'}
            </h3>
            <p className={styles.emptyMessage}>
                {hasFilters 
                    ? 'Try adjusting your search criteria to see more results.'
                    : 'Subscribe to newsletters to stay updated with the latest content.'
                }
            </p>
            {hasFilters ? (
                <button onClick={onClearFilters} className={styles.clearFiltersButton}>
                    Clear Filters
                </button>
            ) : (
                <button onClick={showAvailable} className={styles.browseButton}>
                    Browse Available Newsletters
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
                Please log in to manage your newsletter subscriptions.
            </p>
        </div>
    </div>
);
/**
 * Main Newsletters Page Component
 */
const NewslettersPage = () => {
    const { user, isLoading: authLoading, authChecked } = useAuth();
    const [showAvailable, setShowAvailable] = useState(false);
    const {
        subscriptions,
        availableNewsletters,
        loading,
        error,
        stats,
        actionLoading,
        searchTerm,
        setSearchTerm,
        statusFilter,
        setStatusFilter,
        categoryFilter,
        setCategoryFilter,
        handleSubscribe,
        handleUnsubscribe,
        handleUpdateSubscription,
        handleToggleSubscription,
        refresh,
        clearError,
        autoRefresh,
        isSubscribed
    } = useNewsletters();
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
    const hasFilters = searchTerm.trim() || statusFilter !== 'all' || categoryFilter !== 'all';
    // Clear all filters
    const clearAllFilters = () => {
        setSearchTerm('');
        setStatusFilter('all');
        setCategoryFilter('all');
    };
    return (
        <div className={styles.container}>
            <div className={styles.header}>
                <div className={styles.headerTop}>
                    <h1 className={styles.title}>
                        {getIcon('Mail', { className: styles.titleIcon })}
                        Newsletter Subscriptions
                    </h1>
                    <div className={styles.headerActions}>
                        <button
                            onClick={() => setShowAvailable(!showAvailable)}
                            className={styles.toggleViewButton}
                            title={showAvailable ? 'Show my subscriptions' : 'Browse available newsletters'}
                        >
                            {showAvailable ? getIcon('Eye', { className: styles.iconSmall }) : getIcon('Plus', { className: styles.iconSmall })}
                            {showAvailable ? 'My Subscriptions' : 'Browse Newsletters'}
                        </button>
                    </div>
                </div>
                <StatsDisplay stats={stats} autoRefresh={autoRefresh} />
            </div>
            <FilterBar
                searchTerm={searchTerm}
                setSearchTerm={setSearchTerm}
                statusFilter={statusFilter}
                setStatusFilter={setStatusFilter}
                categoryFilter={categoryFilter}
                setCategoryFilter={setCategoryFilter}
                onRefresh={refresh}
                loading={loading}
                availableNewsletters={availableNewsletters}
            />
            {showAvailable ? (
                // Available Newsletters View
                <div className={styles.availableNewslettersSection}>
                    <h2 className={styles.sectionTitle}>Available Newsletters</h2>
                    <div className={styles.availableNewslettersList}>
                        {availableNewsletters.map((newsletter) => (
                            <AvailableNewsletterItem
                                key={newsletter.newsletterId}
                                newsletter={newsletter}
                                onSubscribe={handleSubscribe}
                                actionLoading={actionLoading}
                                isSubscribed={isSubscribed(newsletter.newsletterId)}
                            />
                        ))}
                    </div>
                </div>
            ) : (
                // User Subscriptions View
                subscriptions.length === 0 ? (
                    <EmptyState 
                        hasFilters={hasFilters} 
                        onClearFilters={clearAllFilters}
                        showAvailable={() => setShowAvailable(true)}
                    />
                ) : (
                    <div className={styles.subscriptionsList}>
                        {subscriptions.map((subscription) => (
                            <SubscriptionItem
                                key={subscription.subscriptionId}
                                subscription={subscription}
                                onUnsubscribe={handleUnsubscribe}
                                onToggleStatus={handleToggleSubscription}
                                onUpdatePreferences={handleUpdateSubscription}
                                actionLoading={actionLoading}
                                availableNewsletters={availableNewsletters}
                            />
                        ))}
                    </div>
                )
            )}
        </div>
    );
};
export default NewslettersPage;