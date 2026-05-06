"use client";
import React, { useEffect, useState, useCallback } from "react";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { toast } from "react-toastify";
import { 
    PlusCircle, 
    Heart, 
    Search, 
    Filter, 
    Grid3X3, 
    List,
    MoreVertical,
    ShoppingBag,
    Trash2,
    Edit3,
    Star,
    Eye
} from "@/icons";
import useWishlist from "../../../hooks/useWishlist";
import { useAuth } from "../../../context/AuthContext";
import WishlistItemCard from "../../../components/wishlist/WishlistItemCard";
import WishlistSelector from "../../../components/wishlist/WishlistSelector";
import WishlistSelectorModal from "../../../components/wishlist/WishlistSelectorModal";
import WishlistErrorBoundary from "../../../components/wishlist/WishlistErrorBoundary";
import styles from "./Wishlist.module.css";
const LAYOUT_MODES = {
    GRID: 'grid',
    LIST: 'list'
};
const SORT_OPTIONS = {
    NEWEST: 'newest',
    OLDEST: 'oldest',
    NAME: 'name'
};
function WishlistsPageContent() {
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
        showWishlistSelector,
        pendingItemToAdd,
        loadDefaultWishlist,
        createNewWishlist,
        deleteWishlist,
        selectWishlist,
        confirmAddToPendingWishlist,
        cancelAddToPendingWishlist
    } = useWishlist();
    // Local UI state
    const [newWishlistName, setNewWishlistName] = useState('');
    const [showCreateForm, setShowCreateForm] = useState(false);
    const [searchQuery, setSearchQuery] = useState('');
    const [layoutMode, setLayoutMode] = useState(LAYOUT_MODES.GRID);
    const [sortBy, setSortBy] = useState(SORT_OPTIONS.NEWEST);
    const [showFilters, setShowFilters] = useState(false);
    // Load wishlists on component mount
    useEffect(() => {
        if (userId) {
            loadDefaultWishlist();
        }
    }, [userId, loadDefaultWishlist]);
    // Handle creating a new wishlist
    const handleCreateWishlist = useCallback(async (e) => {
        e.preventDefault();
        if (!newWishlistName.trim()) {
            toast.warn(t('nameRequired'));
            return;
        }
        const wishlistId = await createNewWishlist(newWishlistName);
        if (wishlistId) {
            setNewWishlistName('');
            setShowCreateForm(false);
        }
    }, [newWishlistName, createNewWishlist, t]);
    // Handle deleting a wishlist with confirmation
    const handleDeleteWishlist = useCallback(async (wishlistId, name) => {
        if (window.confirm(t('confirmDelete', { name }))) {
            await deleteWishlist(wishlistId);
        }
    }, [deleteWishlist, t]);
    // Handle wishlist selection
    const handleSelectWishlist = useCallback(async (wishlistId) => {
        await selectWishlist(wishlistId);
    }, [selectWishlist]);
    // Filter and sort items
    const filteredAndSortedItems = React.useMemo(() => {
        let filtered = items;
        // Apply search filter
        if (searchQuery.trim()) {
            filtered = filtered.filter(item => 
                item.notes?.toLowerCase().includes(searchQuery.toLowerCase()) ||
                item.entityType?.toLowerCase().includes(searchQuery.toLowerCase())
            );
        }
        // Apply sorting
        const sorted = [...filtered].sort((a, b) => {
            switch (sortBy) {
                case SORT_OPTIONS.NEWEST:
                    return new Date(b.createdAt || 0) - new Date(a.createdAt || 0);
                case SORT_OPTIONS.OLDEST:
                    return new Date(a.createdAt || 0) - new Date(b.createdAt || 0);
                case SORT_OPTIONS.NAME:
                    return (a.notes || '').localeCompare(b.notes || '');
                default:
                    return 0;
            }
        });
        return sorted;
    }, [items, searchQuery, sortBy]);
    // Render loading state
    if (loading && !wishlists.length && !items.length) {
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
    if (error && !wishlists.length) {
        return (
            <div className={styles.container}>
                <div className={styles.errorState}>
                    <Heart size={48} className={styles.errorIcon} />
                    <h2>{t('errorTitle')}</h2>
                    <p>{error}</p>
                    <button 
                        onClick={() => loadDefaultWishlist()} 
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
            {/* Header Section */}
            <div className={styles.header}>
                <div className={styles.headerTop}>
                    <div className={styles.titleSection}>
                        <Heart size={32} className={styles.titleIcon} />
                        <h1 className={styles.pageTitle}>{t('myWishlists')}</h1>
                    </div>
                    <button 
                        className={styles.createButton}
                        onClick={() => setShowCreateForm(true)}
                    >
                        <PlusCircle size={18} />
                        <span>{t('createNew')}</span>
                    </button>
                </div>
                {/* Wishlist Selector */}
                {wishlists.length > 0 && (
                    <div className={styles.wishlistSelector}>
                        <WishlistSelector
                            wishlists={wishlists}
                            currentWishlist={currentWishlist}
                            onSelectWishlist={handleSelectWishlist}
                            onDeleteWishlist={handleDeleteWishlist}
                        />
                    </div>
                )}
            </div>
            {/* Create Wishlist Form */}
            {showCreateForm && (
                <div className={styles.createFormOverlay}>
                    <form onSubmit={handleCreateWishlist} className={styles.createForm}>
                        <h3>{t('createNewWishlist')}</h3>
                        <input
                            type="text"
                            value={newWishlistName}
                            onChange={(e) => setNewWishlistName(e.target.value)}
                            placeholder={t('wishlistNamePlaceholder')}
                            className={styles.nameInput}
                            autoFocus
                        />
                        <div className={styles.formActions}>
                            <button type="submit" className={styles.submitButton}>
                                {t('create')}
                            </button>
                            <button 
                                type="button" 
                                className={styles.cancelButton}
                                onClick={() => {
                                    setShowCreateForm(false);
                                    setNewWishlistName('');
                                }}
                            >
                                {t('cancel')}
                            </button>
                        </div>
                    </form>
                </div>
            )}
            {/* Main Content */}
            {currentWishlist ? (
                <div className={styles.mainContent}>
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
                            <select 
                                value={sortBy} 
                                onChange={(e) => setSortBy(e.target.value)}
                                className={styles.sortSelect}
                            >
                                <option value={SORT_OPTIONS.NEWEST}>{t('sortNewest')}</option>
                                <option value={SORT_OPTIONS.OLDEST}>{t('sortOldest')}</option>
                                <option value={SORT_OPTIONS.NAME}>{t('sortName')}</option>
                            </select>
                        </div>
                    </div>
                    {/* Items List */}
                    {filteredAndSortedItems.length > 0 ? (
                        <div className={`${styles.itemsContainer} ${styles[layoutMode]}`}>
                            {filteredAndSortedItems.map(item => (
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
            ) : (
                /* No Current Wishlist Selected */
                <div className={styles.noWishlistState}>
                    {wishlists.length === 0 ? (
                        <>
                            <Heart size={64} className={styles.emptyIcon} />
                            <h2>{t('noWishlists')}</h2>
                            <p>{t('noWishlistsMessage')}</p>
                        </>
                    ) : (
                        <>
                            <Heart size={64} className={styles.emptyIcon} />
                            <h2>{t('selectWishlist')}</h2>
                            <p>{t('selectWishlistMessage')}</p>
                        </>
                    )}
                </div>
            )}
            {/* Wishlist Selector Modal */}
            {showWishlistSelector && (
                <WishlistSelectorModal
                    isOpen={showWishlistSelector}
                    wishlists={wishlists}
                    onSelect={confirmAddToPendingWishlist}
                    onClose={cancelAddToPendingWishlist}
                    pendingItem={pendingItemToAdd}
                />
            )}
        </div>
    );
}
export default function WishlistsPage() {
    return (
        <WishlistErrorBoundary>
            <WishlistsPageContent />
        </WishlistErrorBoundary>
    );
}
