"use client";
import React, { memo } from 'react';
import { useTranslations } from 'next-intl';
import WishlistSelector from './WishlistSelector';
import useWishlist from '../../hooks/useWishlist';
/**
 * Container component that connects WishlistSelector to the useWishlist hook
 * This will automatically appear when showWishlistSelector is true in the useWishlist hook
 */
const WishlistSelectorModal = memo(() => {
    const t = useTranslations('Wishlist');
    const {
        wishlists,
        isLoading: loading,
        showWishlistSelector,
        pendingItemToAdd,
        createWishlist: createNewWishlist,
        confirmAddToPendingWishlist,
        cancelAddToPendingWishlist
    } = useWishlist();
    // Don't render anything if the selector shouldn't be shown
    if (!showWishlistSelector || !pendingItemToAdd) return null;
    return (
        <div className="wishlist-selector-modal-overlay" onClick={cancelAddToPendingWishlist}>
            <div className="wishlist-selector-modal-content" onClick={(e) => e.stopPropagation()}>
                <h3>{t('selectWishlist')}</h3>
                <p>{t('selectWishlistForItem')}</p>
                <WishlistSelector
                    wishlists={wishlists}
                    currentWishlist={null}
                    onSelectWishlist={(wishlistId) => confirmAddToPendingWishlist(wishlistId)}
                    onDeleteWishlist={() => {}}
                    onCreateNew={(name) => {
                        // Create new wishlist and add the item to it
                        createNewWishlist(name).then(wishlistId => {
                            if (wishlistId) {
                                confirmAddToPendingWishlist(wishlistId);
                            }
                        });
                    }}
                />
                <button onClick={cancelAddToPendingWishlist} className="cancel-button">
                    {t('cancel')}
                </button>
            </div>
        </div>
    );
});
WishlistSelectorModal.displayName = 'WishlistSelectorModal';
export default WishlistSelectorModal; 