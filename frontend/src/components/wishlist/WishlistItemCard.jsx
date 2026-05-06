"use client";
import React, {useState, useEffect, Suspense, memo} from 'react';
import {Trash2, AlertCircle, Heart} from '@/icons';
import {useTranslations} from 'next-intl';
import useWishlist from '../../hooks/useWishlist';
import styles from './WishlistItemCard.module.css';
// Import available card components
import { ClassifiedCard } from "../classified";
import ImprovedPostCard from "@/app/[locale]/design/posts/page";
import ImprovedVideoCard from "@/app/[locale]/design/videos/page";
/**
 * Wishlist Item Card - Uses existing card components like unified feed
 * Fetches entity data and renders the appropriate existing card component
 *
 * @param {Object} item - The wishlist item with itemId and entityType
 * @param {string} layoutMode - 'grid' or 'list' layout mode (optional)
 */
const WishlistItemCard = memo(function WishlistItemCard({item, layoutMode = 'grid'}) {
    const t = useTranslations('Wishlist');
    const {removeItem, loading: wishlistLoading} = useWishlist();
    const [entityData, setEntityData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(false);
    const [removing, setRemoving] = useState(false);
    // Fetch the actual entity data based on item_id and entity_type
    useEffect(() => {
        let isMounted = true;
        const fetchEntityData = async () => {
            if (!item?.item_id || !item?.entity_type) {
                setError(true);
                setLoading(false);
                return;
            }
            try {
                let response;
                // Use the same API pattern as unified feed
                switch (item.entity_type.toLowerCase()) {
                    case 'product':
                        const {getProduct} = await import('../../api/searchApi');
                        response = await getProduct(item.item_id);
                        break;
                    case 'post':
                        const {getPost} = await import('../../api/searchApi');
                        response = await getPost(item.item_id);
                        break;
                    case 'video':
                        const {getVideo} = await import('../../api/searchApi');
                        response = await getVideo(item.item_id);
                        break;
                    default:
                        throw new Error(`Unsupported entity type: ${item.entity_type}`);
                }
                // Transform response to match feed item format
                const transformedData = {
                    id: response.id || item.item_idresponse,
                    entityType: item.entity_type,
                    createdAt: response.createdAt || new Date().toISOString(),
                    // Add wishlist-specific metadata
                    wishlistItemId: item.id,
                    wishlistNotes: item.notes,
                    addedToWishlistAt: item.addedAt
                };
                if (isMounted) {
                    setEntityData(transformedData);
                    setError(false);
                }
            } catch (err) {
                if (process.env.NODE_ENV === 'development') {
                }
                if (isMounted) {
                    setError(true);
                }
            } finally {
                if (isMounted) {
                    setLoading(false);
                }
            }
        };
        fetchEntityData();
        return () => {
            isMounted = false;
        };
    }, [item?.itemId, item?.entityType]);
    // Handle removing item from wishlist
    const handleRemoveFromWishlist = async (e) => {
        e.stopPropagation();
        e.preventDefault();
        if (removing || wishlistLoading) return;
        setRemoving(true);
        try {
            await removeItem(item.id);
        } catch (err) {
            if (process.env.NODE_ENV === 'development') {
            }
            // Error: t('failedToRemove'...);
        } finally {
            setRemoving(false);
        }
    };
    // Render the appropriate card component based on entityType
    const renderEntityCard = () => {
        if (!entityData) return null;
        try {
            // Only support video, post, and product cards
            switch (item.entity_type.toLowerCase()) {
                case 'product':
                    // Use ClassifiedCard for products
                    return <ClassifiedCard product={entityData}/>;
                case 'post':
                    return <ImprovedPostCard post={entityData}/>;
                case 'video':
                    return <ImprovedVideoCard video={entityData}/>;
                default:
                    return (
                        <div className={styles.errorState}>
                            <p>{t('unsupportedType', {type: item.entity_type})}</p>
                        </div>
                    );
            }
        } catch (err) {
            return (
                <div className={styles.errorState}>
                    <p>{t('renderError', {type: item.entity_type, message: err.message})}</p>
                </div>
            );
        }
    };
    // Loading state
    if (loading) {
        return (
            <div className={`${styles.wishlistCard} ${styles.loading}`}>
                <div className={styles.loadingContainer}>
                    <div className={styles.loadingSpinner}></div>
                    <p>{t('loadingItem')}</p>
                </div>
            </div>
        );
    }
    // Error state
    if (error || !entityData) {
        return (
            <div className={`${styles.wishlistCard} ${styles.errorCard}`}>
                <div className={styles.errorContainer}>
                    <AlertCircle className={styles.errorIcon}/>
                    <div className={styles.errorContent}>
                        <h3>{t('itemNotFound')}</h3>
                        <p>{t('itemMayBeRemoved')}</p>
                        <small>
                            {item?.entity_type} ID: {item?.item_id}
                        </small>
                    </div>
                </div>
                <button
                    className={`${styles.removeButton} ${removing ? styles.removing : ''}`}
                    onClick={handleRemoveFromWishlist}
                    disabled={removing || wishlistLoading}
                    title={t('removeFromWishlist')}
                    aria-label={t('removeFromWishlist')}
                >
                    <Trash2 size={18}/>
                </button>
            </div>
        );
    }
    // Main render - existing card component wrapped with wishlist functionality
    return (
        <div className={styles.wishlistCard}>
            {/* Wishlist metadata header */}
            {item.notes && (
                <div className={styles.wishlistHeader}>
                    <div className={styles.wishlistNotes}>
                        <Heart size={14} className={styles.noteIcon}/>
                        <span>{item.notes}</span>
                    </div>
                </div>
            )}
            {/* Existing card component */}
            <div className={styles.cardWrapper}>
                <Suspense fallback={
                    <div className={styles.cardLoading}>
                        <div className={styles.loadingSpinner}></div>
                    </div>
                }>
                    {renderEntityCard()}
                </Suspense>
            </div>
            {/* Wishlist remove button overlay */}
            <button
                className={`${styles.removeButton} ${removing ? styles.removing : ''}`}
                onClick={handleRemoveFromWishlist}
                disabled={removing || wishlistLoading}
                title={t('removeFromWishlist')}
                aria-label={t('removeFromWishlist')}
            >
                {removing ? (
                    <div className={styles.removeSpinner}></div>
                ) : (
                    <Trash2 size={18}/>
                )}
            </button>
            {/* Wishlist metadata footer */}
            {item.addedAt && (
                <div className={styles.wishlistFooter}>
                    <span className={styles.addedDate}>
                        {t('addedOn', {date: new Date(item.addedAt).toLocaleDateString()})}
                    </span>
                </div>
            )}
        </div>
    );
});
export default WishlistItemCard; 