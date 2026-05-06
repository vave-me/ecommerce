"use client";
import React, {useState, useMemo, memo} from 'react';
import {
    ShoppingBag,
    Search,
    Filter,
    RefreshCw,
    Check,
    X,
    AlertCircle,
    Clock,
    FileText,
    Eye,
    CheckCircle,
    XCircle,
    User,
    Package,
    Euro,
    Calendar,
    Play,
    Ban
} from '@/icons';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import {useAuth} from '../../../context/AuthContext';
import {useOffers} from '../../../hooks/queries/useOffersQuery';
import {getOfferStatusTypes, formatPrice} from '../../../api/offersApi';
import styles from './page.module.css';

dayjs.extend(relativeTime);
// Icon mapping for offer statuses
const statusIconMap = {
    FileText,
    Eye,
    CheckCircle,
    XCircle
};
/**
 * OfferItem Component
 * Individual offer item with actions
 */
const OfferItem = memo(({offer, onAccept, onActivate, onClose, currentUserId}) => {
    const {
        id,
        userSellerId,
        userCustomerId,
        productId,
        productName,
        productDescription,
        sellerName,
        price,
        offerStatus,
        createdAt
    } = offer;
    const statusConfig = getOfferStatusTypes().find(s => s.value === offerStatus) ||
        {label: offerStatus, icon: 'FileText', color: '#6b7280'};
    const timeAgo = createdAt ? dayjs(createdAt).fromNow() : '';
    const formattedPrice = formatPrice(price);
    // Get the icon component
    const StatusIcon = statusIconMap[statusConfig.icon] || FileText;
    const canAccept = offerStatus === 'active' && !userCustomerId && currentUserId !== userSellerId;
    const canActivate = offerStatus === 'draft' && currentUserId === userSellerId;
    const canClose = (offerStatus === 'active' || offerStatus === 'draft') && currentUserId === userSellerId;
    const handleAccept = () => {
        onAccept(id, currentUserId);
    };
    const handleActivate = () => {
        onActivate(id);
    };
    const handleClose = () => {
        onClose(id, 'Closed by seller');
    };
    return (
        <div className={`${styles.offerItem} ${styles[offerStatus]}`}>
            <div className={styles.offerIcon}>
                <StatusIcon size={16} style={{color: statusConfig.color}}/>
            </div>
            <div className={styles.offerContent}>
                <div className={styles.offerHeader}>
                    <div className={styles.offerTitle}>
                        <h3 className={styles.productName}>{productName}</h3>
                        <span className={styles.offerStatus} style={{color: statusConfig.color}}>
                            {statusConfig.label}
                        </span>
                    </div>
                    <div className={styles.offerPrice}>{formattedPrice}</div>
                </div>
                <p className={styles.productDescription}>{productDescription}</p>
                <div className={styles.offerDetails}>
                    <div className={styles.offerMeta}>
                        <span className={styles.metaItem}>
                            <User size={12}/>
                            {sellerName}
                        </span>
                        <span className={styles.metaItem}>
                            <Package size={12}/>
                            ID: {productId}
                        </span>
                        {timeAgo && (
                            <span className={styles.metaItem}>
                                <Clock size={12}/>
                                {timeAgo}
                            </span>
                        )}
                    </div>
                </div>
            </div>
            <div className={styles.offerActions}>
                {canAccept && (
                    <button
                        onClick={handleAccept}
                        className={`${styles.actionButton} ${styles.acceptButton}`}
                        title="Accept offer"
                        aria-label="Accept offer"
                    >
                        <Check size={16}/>
                    </button>
                )}
                {canActivate && (
                    <button
                        onClick={handleActivate}
                        className={`${styles.actionButton} ${styles.activateButton}`}
                        title="Activate offer"
                        aria-label="Activate offer"
                    >
                        <Play size={16}/>
                    </button>
                )}
                {canClose && (
                    <button
                        onClick={handleClose}
                        className={`${styles.actionButton} ${styles.closeButton}`}
                        title="Close offer"
                        aria-label="Close offer"
                    >
                        <Ban size={16}/>
                    </button>
                )}
            </div>
        </div>
    );
});
OfferItem.displayName = 'OfferItem';
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
                            offerCounts
                        }) => {
    const statusTypes = getOfferStatusTypes();
    return (
        <div className={styles.filterBar}>
            <div className={styles.searchContainer}>
                <Search size={20} className={styles.searchIcon}/>
                <input
                    type="text"
                    placeholder="Search offers..."
                    value={searchTerm}
                    onChange={(e) => onSearchChange(e.target.value)}
                    className={styles.searchInput}
                />
            </div>
            <div className={styles.filterControls}>
                <select
                    value={filters.offerStatus || ''}
                    onChange={(e) => onUpdateFilters({offerStatus: e.target.value || undefined})}
                    className={styles.filterSelect}
                >
                    <option value="">All Status</option>
                    {statusTypes.map(status => (
                        <option key={status.value} value={status.value}>
                            {status.label} ({offerCounts[status.value] || 0})
                        </option>
                    ))}
                </select>
                <select
                    value={filters.userSellerId || ''}
                    onChange={(e) => onUpdateFilters({userSellerId: e.target.value || undefined})}
                    className={styles.filterSelect}
                >
                    <option value="">All Sellers</option>
                    <option value="seller_123">My Offers</option>
                </select>
                <button
                    onClick={onClearFilters}
                    className={styles.clearFiltersButton}
                    title="Clear filters"
                >
                    <Filter size={16}/>
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
const EmptyState = ({hasFilters, onClearFilters}) => (
    <div className={styles.emptyState}>
        <ShoppingBag size={48} className={styles.emptyIcon}/>
        <h3>No offers found</h3>
        <p>
            {hasFilters
                ? "No offers match your current filters."
                : "No offers available. Create your first offer to get started!"}
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
 */
const ErrorState = ({error, onRetry}) => (
    <div className={styles.errorState}>
        <AlertCircle size={48} className={styles.errorIcon}/>
        <h3>Failed to load offers</h3>
        <p>{error || 'Something went wrong while loading offers.'}</p>
        <button onClick={onRetry} className={styles.retryButton}>
            <RefreshCw size={16}/>
            Try again
        </button>
    </div>
);
/**
 * LoadingState Component
 */
const LoadingState = () => (
    <div className={styles.loadingState}>
        <div className={styles.loadingSpinner}>
            <RefreshCw size={32} className={styles.spinIcon}/>
        </div>
        <h3>Loading offers...</h3>
        <p>Please wait while we fetch your offers.</p>
    </div>
);
/**
 * LoginRequired Component
 */
const LoginRequired = () => (
    <div className={styles.loginRequired}>
        <User size={48} className={styles.emptyIcon}/>
        <h2>Login Required</h2>
        <p>Please log in to view and manage offers.</p>
    </div>
);
/**
 * Main Offers Page Component
 */
export default function OffersPage() {
    const {user, isLoading, authChecked} = useAuth();
    const {
        offers,
        loading,
        error,
        filters,
        searchTerm,
        setSearchTerm,
        updateFilters,
        clearFilters,
        refreshOffers,
        acceptOffer,
        activateOffer,
        closeOffer,
        offerCounts,
        hasFilters,
        isEmpty
    } = useOffers({
        autoRefresh: true,
        refreshInterval: 30000
    });
    // Handle offer actions
    const handleAcceptOffer = async (offerId, userCustomerId) => {
        const success = await acceptOffer(offerId, userCustomerId);
        if (success) {
            // Could show success notification here
        }
    };
    const handleActivateOffer = async (offerId) => {
        const success = await activateOffer(offerId);
        if (success) {
            // Could show success notification here
        }
    };
    const handleCloseOffer = async (offerId, reason) => {
        const success = await closeOffer(offerId, reason);
        if (success) {
            // Could show success notification here
        }
    };
    const handleRefresh = () => {
        refreshOffers();
    };
    // Show loading state while auth is being checked
    if (isLoading || !authChecked) {
        return (
            <div className={styles.container}>
                <LoadingState/>
            </div>
        );
    }
    // Show login required state if user is not authenticated
    if (!user) {
        return (
            <div className={styles.container}>
                <LoginRequired/>
            </div>
        );
    }
    return (
        <div className={styles.container}>
            <div className={styles.header}>
                <div className={styles.titleSection}>
                    <h1 className={styles.title}>
                        <ShoppingBag size={28}/>
                        Offers
                        {offerCounts.active > 0 && (
                            <span className={styles.activeBadge}>{offerCounts.active}</span>
                        )}
                    </h1>
                    <p className={styles.subtitle}>
                        Manage your offers and negotiations • {offerCounts.total} total offers
                    </p>
                </div>
                <div className={styles.headerActions}>
                    <button
                        onClick={handleRefresh}
                        disabled={loading}
                        className={styles.refreshButton}
                        title="Refresh offers"
                    >
                        <RefreshCw size={16} className={loading ? styles.spinning : ''}/>
                        Refresh
                    </button>
                </div>
            </div>
            <FilterBar
                filters={filters}
                onUpdateFilters={updateFilters}
                onClearFilters={clearFilters}
                searchTerm={searchTerm}
                onSearchChange={setSearchTerm}
                offerCounts={offerCounts}
            />
            <div className={styles.content}>
                {loading ? (
                    <LoadingState/>
                ) : error ? (
                    <ErrorState error={error} onRetry={handleRefresh}/>
                ) : isEmpty ? (
                    <EmptyState hasFilters={hasFilters} onClearFilters={clearFilters}/>
                ) : (
                    <div className={styles.offersList}>
                        {offers.map(offer => (
                            <OfferItem
                                key={offer.id}
                                offer={offer}
                                onAccept={handleAcceptOffer}
                                onActivate={handleActivateOffer}
                                onClose={handleCloseOffer}
                                currentUserId={user?.id || user?.userId}
                            />
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
}
