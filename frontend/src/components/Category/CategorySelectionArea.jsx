"use client";
import React, { useState, useEffect, memo } from 'react';
import PropTypes from 'prop-types';
import { Check, X } from '@/icons';
import { useTranslations } from 'next-intl';
import { useCategories } from '../../hooks/useCategories';
import CategoryTree from './CategoryTree';
import styles from './CategorySelectionArea.module.css';
/**
 * Enhanced category selection area
 * Designed to fit seamlessly in product creation flows
 * Uses next-intl for translations.
 */
const CategorySelectionArea = memo(function CategorySelectionArea({
    selectedCategoryId,
    onSelectCategory,
    onClear,
    categoryType,
    categories: externalCategories
}) {
    // Fetch categories with our hook if not provided externally
    const { 
        data: hookedCategories = [], 
        isLoading: isLoadingCategories 
    } = useCategories(categoryType);
    // Use provided categories or those from the hook
    const categories = externalCategories && externalCategories.length > 0 
        ? externalCategories 
        : hookedCategories;
    // Create a fallback t function in case translations aren't available
    const defaultT = (key, options = {}) => {
        const fallbacks = {
            'instructionSelect': 'Please select a category',
            'instructionCurrent': `Current category: ${options.categoryName || 'Unknown'}`,
            'unnamedCategoryFallback': 'Unnamed Category',
            'clearSelectionAria': 'Clear category selection',
            'loading': 'Loading categories...'
        };
        return fallbacks[key] || key;
    };
    // Try to use next-intl translations, fall back to our defaults if it fails
    let t;
    try {
        t = useTranslations('CategorySelectionArea');
    } catch (error) {
        t = defaultT;
    }
    const [selectedInfo, setSelectedInfo] = useState(null);
    // Handle category selection
    const handleCategorySelect = (category) => {
        setSelectedInfo(category);
        if (onSelectCategory) {
            onSelectCategory(category);
        }
    };
    // Handle clearing selection
    const handleClear = () => {
        setSelectedInfo(null);
        if (onClear) {
            onClear();
        }
    };
    // Find and update selected category info when selectedCategoryId changes
    useEffect(() => {
        // If we have a selectedCategoryId but no selectedInfo, try to find it
        if (selectedCategoryId && !selectedInfo && categories?.length > 0) {
            // Simple flat search for the category (could be improved with recursion for subcategories)
            const findCategory = (cats, id) => {
                for (const cat of cats) {
                    if (cat.id === id) return cat;
                    if (cat.subcategories?.length > 0) {
                        const found = findCategory(cat.subcategories, id);
                        if (found) return found;
                    }
                }
                return null;
            };
            const found = findCategory(categories, selectedCategoryId);
            if (found) setSelectedInfo(found);
        }
        // If selectedCategoryId is cleared, clear our local state too
        else if (!selectedCategoryId && selectedInfo) {
            setSelectedInfo(null);
        }
    }, [selectedCategoryId, selectedInfo, categories]);
    return (
        <div className={styles.container}>
            {/* Selected category display */}
            {selectedInfo && (
                <div className={styles.selectionDisplay} style={{
                    backgroundColor: '#f0f9ff',
                    color: '#1e40af',
                    borderColor: '#e0e7ff'
                }}>
                    <div className={styles.selectionInfo}>
                        <Check size={18} className={styles.checkIcon} style={{color: '#2980b9'}} />
                        <span className={styles.selectedName} style={{color: '#1e40af'}}>
                            {selectedInfo.name || t('unnamedCategoryFallback')}
                        </span>
                    </div>
                    <button
                        type="button"
                        className={styles.clearButton}
                        onClick={handleClear}
                        aria-label={t('clearSelectionAria')}
                        style={{color: '#2980b9'}}
                    >
                        <X size={18} />
                    </button>
                </div>
            )}
            {/* Category browser */}
            <div className={styles.categoryBrowser}>
                {isLoadingCategories ? (
                    <p className={styles.loadingText}>{t('loading')}</p>
                ) : !selectedInfo ? (
                    <p className={styles.instructionText}>{t('instructionSelect')}</p>
                ) : null}
                {/* Pass selection handlers to CategoryTree */}
                <CategoryTree
                    selectedCategoryId={selectedInfo ? selectedInfo.id : selectedCategoryId}
                    onSelectCategory={handleCategorySelect}
                    maxHeight={240}
                    modalContext={true}
                    searchable={true}
                    categories={categories}
                    categoryType={categoryType}
                />
            </div>
        </div>
    );
});
// Proper PropTypes
CategorySelectionArea.propTypes = {
    selectedCategoryId: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
    onSelectCategory: PropTypes.func.isRequired,
    onClear: PropTypes.func.isRequired,
    categoryType: PropTypes.string,
    categories: PropTypes.array
};
export default CategorySelectionArea;