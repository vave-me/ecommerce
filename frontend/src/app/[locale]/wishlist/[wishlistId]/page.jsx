"use client";
import React, { useEffect, useState, useCallback } from "react";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { ArrowLeft, RefreshCw, Heart, Search, Grid3X3, List, ShoppingBag } from "@/icons";
import useWishlist from "../../../../hooks/useWishlist";
import { useAuth } from "../../../../context/AuthContext";
import WishlistItemCard from "../../../../components/wishlist/WishlistItemCard";
import WishlistErrorBoundary from "../../../../components/wishlist/WishlistErrorBoundary";
import styles from "./WishlistDetail.module.css";
const LAYOUT_MODES = {
    GRID: 'grid',
    LIST: 'list'
};
function WishlistDetailPageContent({ params }) {
    const { wishlistId } = params;
    const t = useTranslations('Wishlist');
    const router = useRouter();
    const { user } = useAuth();
    const userId = user?.userId;
    // Get wishlist data and functions from our hook
    const { 
        wishlists,
        currentWishlist,
        items,
        loading,
        error,
        selectWishlist,
        loadDefaultWishlist,
        loadWishlistItems
    } = useWishlist();
    // Local UI state
    const [searchQuery, setSearchQuery] = useState('');
    const [layoutMode, setLayoutMode] = useState(LAYOUT_MODES.GRID);
    const [pageLoading, setPageLoading] = useState(true);
    // Load specific wishlist data
    useEffect(() => {
        const loadData = async () => {
            if (!userId || !wishlistId) return;
            setPageLoading(true);
            try {
                // Load all wishlists first if needed
                if (wishlists.length === 0) {
                    await loadDefaultWishlist();
                }
                // Select this specific wishlist
                await selectWishlist(wishlistId);
            } catch (err) {
                if (process.env.NODE_ENV === 'development') {
                }
                // Error: t('loadError'...);
            } finally {
                setPageLoading(false);
            }
        };
        loadData();
    }, [wishlistId, userId, wishlists.length, loadDefaultWishlist, selectWishlist, t]);
    // Handle navigating back to all wishlists
    const handleBackToWishlists = useCallback(() => {
        router.push('/wishlist');
    }, [router]);
    // Handle refreshing the wishlist
    const handleRefresh = useCallback(async () => {
        if (!wishlistId) return;
        try {
            await loadWishlistItems(wishlistId);
            // Wishlist items loaded successfully
        } catch (err) {
            if (process.env.NODE_ENV === 'development') {
            }
            // Error: t('refreshError'...);
        }
    }, [wishlistId, loadWishlistItems, t]);
    // Filter items based on search query
    const filteredItems = React.useMemo(() => {
        if (!searchQuery.trim()) return items;
        return items.filter(item => 
            item.notes?.toLowerCase().includes(searchQuery.toLowerCase()) ||
            item.entityType?.toLowerCase().includes(searchQuery.toLowerCase())
        );
    }, [items, searchQuery]);
    // Get current wishlist details
    const wishlistDetails = wishlists.find(w => w.id === wishlistId) || currentWishlist;
    // Render loading state
    if (pageLoading || (loading && !items.length)) {
        return (
            <div className={styles.container}>
                <div className={styles.loadingState}>
                    <div className={styles.spinner}></div>
                    <p>{t('loading')}</p>
                </div>
            </div>
        );
    }
    // Render error state
    if (error) {
        return (
            <div className={styles.container}>
                <div className={styles.errorState}>
                    <Heart size={48} className={styles.errorIcon} />
                    <h2>{t('errorTitle')}</h2>
                    <p>{error}</p>
                    <button 
                        onClick={handleRefresh}
                        className={styles.retryButton}
                    >
                        {t('retry')}
                    </button>
                </div>
            </div>
        );
    }
    // Render when user is not logged in
    if (!userId) {
        return (
            <div className={styles.container}>
                <div className={styles.notLoggedInState}>
                    <Heart size={64} className={styles.heartIcon} />
                    <h2>{t('loginRequired')}</h2>
                    <p>{t('loginMessage')}</p>
                    <button 
                        onClick={() => router.push('/login')} 
                        className={styles.loginButton}
                    >
                        {t('login')}
                    </button>
                </div>
            </div>
        );
    }
    return (
        <div className={styles.container}>
            {/* Header */}
            <div className={styles.header}>
                <div className={styles.headerTop}>
                    <button 
                        className={styles.backButton} 
                        onClick={handleBackToWishlists}
                        aria-label={t('backToWishlists')}
                    >
                        <ArrowLeft size={20} />
                        <span>{t('backToWishlists')}</span>
                    </button>
                    <div className={styles.titleSection}>
                        <Heart size={28} className={styles.titleIcon} />
                        <h1 className={styles.pageTitle}>
                            {wishlistDetails?.name || t('wishlist')}
                        </h1>
                    </div>
                    <button 
                        className={styles.refreshButton}
                        onClick={handleRefresh}
                        disabled={loading}
                        aria-label={t('refresh')}
                    >
                        <RefreshCw 
                            size={18} 
                            className={loading ? styles.spinning : ''} 
                        />
                        <span>{t('refresh')}</span>
                    </button>
                </div>
                {/* Description */}
                {wishlistDetails?.description && (
                    <p className={styles.description}>
                        {wishlistDetails.description}
                    </p>
                )}
            </div>
            {/* Toolbar */}
            <div className={styles.toolbar}>
                <div className={styles.searchSection}>
                    <div className={styles.searchInput}>
                        <Search size={18} className={styles.searchIcon} />
                        <input
                            type="text"
                            placeholder={t('searchItems')}
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                        />
                    </div>
                </div>
                <div className={styles.viewControls}>
                    <button
                        className={`${styles.viewButton} ${layoutMode === LAYOUT_MODES.GRID ? styles.active : ''}`}
                        onClick={() => setLayoutMode(LAYOUT_MODES.GRID)}
                        title={t('gridView')}
                    >
                        <Grid3X3 size={18} />
                    </button>
                    <button
                        className={`${styles.viewButton} ${layoutMode === LAYOUT_MODES.LIST ? styles.active : ''}`}
                        onClick={() => setLayoutMode(LAYOUT_MODES.LIST)}
                        title={t('listView')}
                    >
                        <List size={18} />
                    </button>
                </div>
            </div>
            {/* Items Content */}
            <div className={styles.mainContent}>
                {filteredItems.length > 0 ? (
                    <div className={`${styles.itemsContainer} ${styles[layoutMode]}`}>
                        {filteredItems.map(item => (
                            <WishlistItemCard
                                key={item.id}
                                item={item}
                                layoutMode={layoutMode}
                            />
                        ))}
                    </div>
                ) : (
                    <div className={styles.emptyState}>
                        {searchQuery ? (
                            <>
                                <Search size={48} className={styles.emptyIcon} />
                                <h3>{t('noSearchResults')}</h3>
                                <p>{t('noSearchResultsMessage')}</p>
                                <button 
                                    onClick={() => setSearchQuery('')} 
                                    className={styles.clearSearchButton}
                                >
                                    {t('clearSearch')}
                                </button>
                            </>
                        ) : (
                            <>
                                <ShoppingBag size={48} className={styles.emptyIcon} />
                                <h3>{t('emptyWishlist')}</h3>
                                <p>{t('emptyWishlistMessage')}</p>
                                <button 
                                    onClick={() => router.push('/explore')} 
                                    className={styles.exploreButton}
                                >
                                    {t('startShopping')}
                                </button>
                            </>
                        )}
                    </div>
                )}
            </div>
        </div>
    );
}
export default function WishlistDetailPage({ params }) {
    return (
        <WishlistErrorBoundary>
            <WishlistDetailPageContent params={params} />
        </WishlistErrorBoundary>
    );
}
