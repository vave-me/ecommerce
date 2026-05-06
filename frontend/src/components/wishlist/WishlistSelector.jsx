"use client";
import React, { useState, useRef, useEffect, memo } from "react";
import { useTranslations } from "next-intl";
import { ChevronDown, Heart, Trash2, MoreVertical, Plus, X, Check } from "@/icons";
import styles from "./WishlistSelector.module.css";
/**
 * Modern Wishlist Selector - Dropdown component for selecting wishlists
 * Matches header/topbar design patterns
 */
const WishlistSelector = memo(function WishlistSelector({ 
    wishlists = [], 
    currentWishlist, 
    onSelectWishlist, 
    onDeleteWishlist,
    onCreateNew,
    className = ''
}) {
    const t = useTranslations('Wishlist');
    const [isOpen, setIsOpen] = useState(false);
    const [showDeleteConfirm, setShowDeleteConfirm] = useState(null);
    const dropdownRef = useRef(null);
    const triggerRef = useRef(null);
    const [showCreateForm, setShowCreateForm] = useState(false);
    const [newWishlistName, setNewWishlistName] = useState('');
    // Close dropdown when clicking outside
    useEffect(() => {
        const handleClickOutside = (event) => {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
                setIsOpen(false);
                setShowDeleteConfirm(null);
            }
        };
        if (isOpen) {
            document.addEventListener('mousedown', handleClickOutside);
            return () => document.removeEventListener('mousedown', handleClickOutside);
        }
    }, [isOpen]);
    // Handle keyboard navigation
    useEffect(() => {
        const handleKeyDown = (event) => {
            if (!isOpen) return;
            switch (event.key) {
                case 'Escape':
                    setIsOpen(false);
                    setShowDeleteConfirm(null);
                    triggerRef.current?.focus();
                    break;
                case 'ArrowDown':
                    event.preventDefault();
                    // Focus first item
                    const firstItem = dropdownRef.current?.querySelector('[role="menuitem"]');
                    firstItem?.focus();
                    break;
            }
        };
        if (isOpen) {
            document.addEventListener('keydown', handleKeyDown);
            return () => document.removeEventListener('keydown', handleKeyDown);
        }
    }, [isOpen]);
    const handleSelectWishlist = (wishlistId) => {
        onSelectWishlist(wishlistId);
        setIsOpen(false);
        setShowDeleteConfirm(null);
    };
    const handleDeleteClick = (event, wishlistId, wishlistName) => {
        event.stopPropagation();
        setShowDeleteConfirm({ id: wishlistId, name: wishlistName });
    };
    const confirmDelete = () => {
        if (showDeleteConfirm) {
            onDeleteWishlist(showDeleteConfirm.id, showDeleteConfirm.name);
            setShowDeleteConfirm(null);
        }
    };
    const cancelDelete = () => {
        setShowDeleteConfirm(null);
    };
    const handleCreateWishlist = () => {
        if (newWishlistName.trim()) {
            onCreateNew(newWishlistName.trim());
            setNewWishlistName('');
            setShowCreateForm(false);
        }
    };
    if (!wishlists.length) {
        return (
            <div className={styles.noWishlists}>
                <Heart size={20} className={styles.emptyIcon} />
                <span>{t('noWishlistsAvailable')}</span>
            </div>
        );
    }
    return (
        <div className={`${styles.container} ${className}`} ref={dropdownRef}>
            {/* Dropdown Trigger */}
            <button
                ref={triggerRef}
                className={`${styles.trigger} ${isOpen ? styles.open : ''}`}
                onClick={() => setIsOpen(!isOpen)}
                aria-haspopup="listbox"
                aria-expanded={isOpen}
                aria-label={t('selectWishlistAriaLabel')}
            >
                <div className={styles.triggerContent}>
                    <div className={styles.wishlistInfo}>
                        <Heart size={18} className={styles.heartIcon} />
                        <div className={styles.wishlistDetails}>
                            <span className={styles.wishlistName}>
                                {currentWishlist?.name || t('selectWishlist')}
                            </span>
                            {currentWishlist?.description && (
                                <span className={styles.wishlistDescription}>
                                    {currentWishlist.description}
                                </span>
                            )}
                        </div>
                    </div>
                    <ChevronDown 
                        size={18} 
                        className={`${styles.chevron} ${isOpen ? styles.rotated : ''}`} 
                    />
                </div>
            </button>
            {/* Dropdown Menu */}
            {isOpen && (
                <div className={styles.dropdown} role="listbox">
                    <div className={styles.dropdownHeader}>
                        <span className={styles.dropdownTitle}>{t('selectWishlist')}</span>
                        <span className={styles.wishlistCount}>
                            {wishlists.length} {wishlists.length === 1 ? t('wishlist') : t('wishlists')}
                        </span>
                    </div>
                    <div className={styles.dropdownContent}>
                        {wishlists.map((wishlist) => (
                            <div
                                key={wishlist.id}
                                className={`${styles.wishlistItem} ${
                                    currentWishlist?.id === wishlist.id ? styles.active : ''
                                }`}
                                role="menuitem"
                                tabIndex={0}
                                onClick={() => handleSelectWishlist(wishlist.id)}
                                onKeyDown={(e) => {
                                    if (e.key === 'Enter' || e.key === ' ') {
                                        e.preventDefault();
                                        handleSelectWishlist(wishlist.id);
                                    }
                                }}
                            >
                                <div className={styles.itemContent}>
                                    <div className={styles.itemInfo}>
                                        <Heart 
                                            size={16} 
                                            className={`${styles.itemIcon} ${
                                                currentWishlist?.id === wishlist.id ? styles.activeIcon : ''
                                            }`} 
                                        />
                                        <div className={styles.itemDetails}>
                                            <span className={styles.itemName}>
                                                {wishlist.name || t('untitled')}
                                            </span>
                                            {wishlist.description && (
                                                <span className={styles.itemDescription}>
                                                    {wishlist.description}
                                                </span>
                                            )}
                                        </div>
                                    </div>
                                    {/* Delete Button - Only show if not the default or if there are multiple wishlists */}
                                    {(wishlist.name !== 'Default' || wishlists.length > 1) && (
                                        <button
                                            className={styles.deleteButton}
                                            onClick={(e) => handleDeleteClick(e, wishlist.id, wishlist.name)}
                                            aria-label={t('deleteWishlist')}
                                            title={t('deleteWishlist')}
                                        >
                                            <Trash2 size={14} />
                                        </button>
                                    )}
                                </div>
                                {/* Active indicator */}
                                {currentWishlist?.id === wishlist.id && (
                                    <div className={styles.activeIndicator}></div>
                                )}
                            </div>
                        ))}
                    </div>
                    {/* Add New Wishlist Option in Dropdown */}
                    <div className={styles.dropdownFooter}>
                        {!showCreateForm ? (
                            <button
                                className={styles.createNewOption}
                                onClick={(e) => {
                                    e.stopPropagation();
                                    setShowCreateForm(true);
                                }}
                                aria-label={t('createNewWishlist')}
                            >
                                <Plus size={16} />
                                <span>{t('createNewWishlist')}</span>
                            </button>
                        ) : (
                            <div className={styles.createFormInline}>
                                <input
                                    type="text"
                                    value={newWishlistName}
                                    onChange={(e) => setNewWishlistName(e.target.value)}
                                    onKeyDown={(e) => {
                                        if (e.key === 'Enter' && newWishlistName.trim()) {
                                            handleCreateWishlist();
                                        } else if (e.key === 'Escape') {
                                            setShowCreateForm(false);
                                            setNewWishlistName('');
                                        }
                                    }}
                                    placeholder={t('wishlistNamePlaceholder')}
                                    className={styles.nameInputInline}
                                    autoFocus
                                    aria-label={t('wishlistNamePlaceholder')}
                                />
                                <div className={styles.createActionsInline}>
                                    <button
                                        className={styles.cancelButtonIcon}
                                        onClick={() => {
                                            setShowCreateForm(false);
                                            setNewWishlistName('');
                                        }}
                                        aria-label={t('cancel')}
                                        type="button"
                                    >
                                        <X size={16} />
                                    </button>
                                    <button
                                        className={styles.createButtonIcon}
                                        onClick={handleCreateWishlist}
                                        disabled={!newWishlistName.trim()}
                                        aria-label={t('create')}
                                        type="button"
                                    >
                                        <Check size={16} />
                                    </button>
                                </div>
                            </div>
                        )}
                    </div>
                </div>
            )}
            {/* Delete Confirmation Modal */}
            {showDeleteConfirm && (
                <div className={styles.deleteModal}>
                    <div className={styles.deleteModalContent}>
                        <div className={styles.deleteModalHeader}>
                            <Trash2 size={24} className={styles.deleteModalIcon} />
                            <h3>{t('confirmDeleteTitle')}</h3>
                        </div>
                        <p className={styles.deleteModalMessage}>
                            {t('confirmDeleteMessage', { name: showDeleteConfirm.name })}
                        </p>
                        <div className={styles.deleteModalActions}>
                            <button 
                                className={styles.deleteConfirmButton}
                                onClick={confirmDelete}
                            >
                                {t('delete')}
                            </button>
                            <button 
                                className={styles.deleteCancelButton}
                                onClick={cancelDelete}
                            >
                                {t('cancel')}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
});
export default WishlistSelector;
