"use client";
import React, { useState, useEffect, memo } from "react";
import { useTranslations } from "next-intl";
import { Heart, Plus, Check, Loader2 } from "@/icons";
import useWishlist from "../../hooks/useWishlist";
import styles from "./WishlistButton.module.css";
/**
 * Modern Wishlist Button Component
 * A reusable button for adding/removing items from wishlists
 * Matches header/topbar design patterns
 * OPTIMIZED: React.memo with custom comparison for performance
 */
const WishlistButton = memo(function WishlistButton({ 
    itemId, 
    entityType = 'product', 
    size = 'medium', 
    variant = 'default',
    className = '',
    showText = true,
    notes = '',
    onWishlistChange = null
}) {
    const t = useTranslations('Wishlist');
    const {
        isLoading: loading,
        isInWishlist,
        addToDefaultWishlist: addItem,
        wishlists
    } = useWishlist();
    const [isInList, setIsInList] = useState(false);
    const [isChecking, setIsChecking] = useState(false);
    // Check if item is in current wishlist (optimized - no async checks)
    useEffect(() => {
        if (!itemId) return;
        // Only check current wishlist (synchronous)
        const inCurrentWishlist = isInWishlist(itemId);
        setIsInList(inCurrentWishlist);
    }, [itemId, isInWishlist, wishlists]); // Re-check when wishlists change
    const handleClick = async () => {
        if (loading || isChecking || !itemId) return;
        try {
            await addItem(itemId, entityType, notes);
            // Update local state
            setIsInList(true);
            // Notify parent component
            if (onWishlistChange) {
                onWishlistChange(true);
            }
        } catch (error) {
            if (process.env.NODE_ENV === 'development') {
            }
        }
    };
    const getSizeClasses = () => {
        switch (size) {
            case 'small':
                return styles.small;
            case 'large':
                return styles.large;
            case 'medium':
            default:
                return styles.medium;
        }
    };
    const getVariantClasses = () => {
        switch (variant) {
            case 'outline':
                return styles.outline;
            case 'ghost':
                return styles.ghost;
            case 'minimal':
                return styles.minimal;
            case 'default':
            default:
                return styles.default;
        }
    };
    const getIconSize = () => {
        switch (size) {
            case 'small': return 16;
            case 'large': return 24;
            case 'medium':
            default: return 20;
        }
    };
    const buttonClasses = [
        styles.button,
        getSizeClasses(),
        getVariantClasses(),
        isInList ? styles.active : '',
        (loading || isChecking) ? styles.loading : '',
        className
    ].filter(Boolean).join(' ');
    const renderIcon = () => {
        if (loading || isChecking) {
            return <Loader2 size={getIconSize()} className={styles.spinner} />;
        }
        if (isInList) {
            return <Heart size={getIconSize()} className={styles.heartFilled} />;
        }
        return variant === 'minimal' ? 
            <Plus size={getIconSize()} className={styles.plusIcon} /> :
            <Heart size={getIconSize()} className={styles.heartOutline} />;
    };
    const getButtonText = () => {
        if (loading || isChecking) {
            return t('checking');
        }
        if (isInList) {
            return variant === 'minimal' ? t('added') : t('inWishlist');
        }
        return variant === 'minimal' ? t('addToWishlist') : t('addToWishlist');
    };
    const getAriaLabel = () => {
        if (isInList) {
            return t('itemInWishlistAriaLabel');
        }
        return t('addToWishlistAriaLabel');
    };
    return (
        <>
            <button
                className={buttonClasses}
                onClick={handleClick}
                disabled={loading || isChecking}
                aria-label={getAriaLabel()}
                title={getButtonText()}
            >
                <span className={styles.iconWrapper}>
                    {renderIcon()}
                </span>
                {showText && (
                    <span className={styles.text}>
                        {getButtonText()}
                    </span>
                )}
                {/* Ripple effect container */}
                <span className={styles.ripple}></span>
            </button>
        </>
    );
}, (prevProps, nextProps) => {
    // Custom comparison for optimal performance
    // Skip callback comparison as it should be stable
    return prevProps.itemId === nextProps.itemId &&
           prevProps.entityType === nextProps.entityType &&
           prevProps.size === nextProps.size &&
           prevProps.variant === nextProps.variant &&
           prevProps.className === nextProps.className &&
           prevProps.showText === nextProps.showText &&
           prevProps.notes === nextProps.notes;
});
export default WishlistButton; 