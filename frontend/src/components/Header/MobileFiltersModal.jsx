"use client";
import React, { useState, useEffect, useRef, useCallback, useMemo } from "react";
import { useDispatch, useSelector } from "react-redux";
import * as Dialog from "@radix-ui/react-dialog";
import { X, ChevronDown } from "@/icons";
import { useTranslations } from "next-intl";
import { setFilters, clearFilters } from "../../redux/slices/listingFiltersSlice";
import { useCategories } from "../../hooks/useCategories";
import styles from "./SearchBar.module.css";
/**
 * MobileFiltersModal Component
 * Mobile-optimized filter modal connected to Redux
 * Features:
 * - Category selection with API integration
 * - Connected to Redux for state management
 * - Responsive design optimized for mobile
 */
const MobileFiltersModal = function MobileFiltersModal({
    showMobileFilters,
    onOpenChange,
    onSubmit
}) {
    const dispatch = useDispatch();
    const t = useTranslations('Filters');
    const tWishlist = useTranslations('Wishlist');
    const filters = useSelector(state => state.listingFilters);
    
    // Fetch categories from API
    const { data: categories, isLoading: categoriesLoading } = useCategories('marketplace');
    
    // Local state for immediate UI feedback
    const [localFilters, setLocalFilters] = useState(filters);
    // Sync local state with Redux state
    useEffect(() => {
        setLocalFilters(filters);
    }, [filters]);
    // Handle category selection
    const handleCategorySelect = useCallback((categorySlug) => {
        if (categorySlug === 'all') {
            const { category: _, ...restFilters } = localFilters;
            setLocalFilters(restFilters);
        } else {
            setLocalFilters(prev => ({ ...prev, category: categorySlug }));
        }
    }, [localFilters]);
    // Handle clear all filters
    const handleClearAll = useCallback(() => {
        setLocalFilters({});
    }, []);
    // Handle apply filters
    const handleApply = useCallback(() => {
        dispatch(setFilters(localFilters));
        onSubmit?.();
        onOpenChange(false);
    }, [dispatch, localFilters, onSubmit, onOpenChange]);
    
    // Handle cancel
    const handleCancel = useCallback(() => {
        setLocalFilters(filters); // Reset to Redux state
        onOpenChange(false);
    }, [filters, onOpenChange]);
    
    // Get selected category object for display
    const selectedCategory = useMemo(() => {
        if (!localFilters.category || !categories) return null;
        return categories.find(cat => cat.slug === localFilters.category || cat.id === localFilters.category);
    }, [localFilters.category, categories]);
    return (
        <Dialog.Root
            open={showMobileFilters}
            onOpenChange={onOpenChange}
        >
            <Dialog.Portal>
                <Dialog.Overlay className={styles.modalOverlay}/>
                <Dialog.Content className={styles.modalContent}>
                    {/* Visual handle for dragging */}
                    <div className={styles.modalHandle}>
                        <div className={styles.handleBar}></div>
                    </div>
                    <div className={styles.modalHeader}>
                        <Dialog.Title className={styles.modalTitle}>
                            {t('sidebarTitle')}
                        </Dialog.Title>
                        <Dialog.Description className="sr-only">
                            {t('mobileDescription', 'Configure search filters')}
                        </Dialog.Description>
                        <Dialog.Close 
                            aria-label={t('closeButtonAriaLabel')}
                            className={styles.modalCloseButton}
                        >
                            <X size={18}/>
                        </Dialog.Close>
                    </div>
                    <div className={styles.modalBody}>
                        {/* Category selection */}
                        <div className={styles.filterGroup}>
                            <label className={styles.filterLabel}>
                                {t('categoryLabel', 'Category')}
                            </label>
                            <div className={styles.categoryGrid}>
                                {/* All Categories button */}
                                <button
                                    type="button"
                                    className={`${styles.categoryButton} ${!localFilters.category ? styles.categoryActive : ''}`}
                                    onClick={() => handleCategorySelect('all')}
                                    aria-pressed={!localFilters.category}
                                >
                                    {t('allCategories')}
                                </button>
                                
                                {/* Category buttons */}
                                {categoriesLoading ? (
                                    <div className={styles.loadingText}>
                                        {t('loadingCategories')}
                                    </div>
                                ) : (
                                    categories?.map(category => (
                                        <button
                                            key={category.id}
                                            type="button"
                                            className={`${styles.categoryButton} ${
                                                (localFilters.category === category.slug || localFilters.category === category.id)
                                                    ? styles.categoryActive : ''
                                            }`}
                                            onClick={() => handleCategorySelect(category.slug || category.id)}
                                            aria-pressed={localFilters.category === category.slug || localFilters.category === category.id}
                                        >
                                            {category.name || category.description}
                                        </button>
                                    ))
                                )}
                            </div>
                        </div>
                        {/* Filter summary */}
                        <div className={styles.filterSummary}>
                            <div className={styles.summaryTitle}>
                                {t('currentFilters', 'Current filters:')}
                            </div>
                            <div className={styles.filterChips}>
                                {selectedCategory && (
                                    <div className={styles.filterChip}>
                                        <span>{selectedCategory.name || selectedCategory.description}</span>
                                        <X
                                            size={14}
                                            className={styles.chipIcon}
                                            onClick={() => handleCategorySelect('all')}
                                        />
                                    </div>
                                )}
                                {Object.keys(localFilters).length === 0 && (
                                    <div className={styles.emptyFilters}>
                                        {t('noFiltersSelected', 'No filters selected')}
                                    </div>
                                )}
                            </div>
                        </div>
                    </div>
                    <div className={styles.modalFooter}>
                        <button
                            type="button"
                            onClick={handleClearAll}
                            className={styles.clearAllButton}
                            disabled={Object.keys(localFilters).length === 0}
                        >
                            {t('clearButton')}
                        </button>
                        <div className={styles.modalActions}>
                            <button
                                type="button"
                                onClick={handleCancel}
                                className={styles.cancelButton}
                            >
                                {tWishlist('cancel')}
                            </button>
                            <button
                                type="button"
                                onClick={handleApply}
                                className={styles.applyButton}
                            >
                                {t('applyButton')}
                            </button>
                        </div>
                    </div>
                </Dialog.Content>
            </Dialog.Portal>
        </Dialog.Root>
    );
};
export default MobileFiltersModal; 