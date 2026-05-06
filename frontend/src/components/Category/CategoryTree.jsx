"use client";
import React, { useCallback, useEffect, useRef, useState, memo } from "react";
import {ChevronDown, ChevronRight, Loader, Search, X} from "@/icons";
import PropTypes from "prop-types";
import styles from "./CategoryTree.module.css";
import {filterCategories} from "../../utils/categoryUtils";
import {fetchSubCategories} from "../../api/categories";
// Optional import to handle cases where next-intl might not be fully configured
let useTranslations;
try {
    useTranslations = require("next-intl").useTranslations;
} catch (e) {
    useTranslations = () => ({
        // Provide fallbacks for all used translation keys
        loading: "Loading...",
        errorLoading: "Error loading categories",
        retry: "Retry",
        empty: "No categories available",
        searchPlaceholder: "Search categories",
        clearSearchAria: "Clear search",
        noResults: "No results for '{query}'",
        unnamedCategory: "Unnamed Category",
    });
}
/**
 * Enhanced CategoryTree Component
 * Optimized for use in modals with improved UX/UI and mobile responsiveness
 * Uses next-intl for translations with fallbacks for missing translations.
 */
const CategoryTree = memo(function CategoryTree({
    selectedCategoryId = null,
    onSelectCategory,
    maxHeight = 240,
    modalContext = true,
    searchable = true,
    categories = [],
    categoryType = 'deal'
}) {
    // Create a fallback t function in case translations aren't available
    const defaultT = (key, options = {}) => {
        const fallbacks = {
            'unnamedCategory': 'Unnamed Category',
            'loading': 'Loading categories...',
            'errorLoading': 'Failed to load categories',
            'empty': 'No categories available',
            'retry': 'Retry',
            'noResults': `No results for "${options.query || ''}"`,
            'searchPlaceholder': 'Search categories...',
            'clearSearchAria': 'Clear search'
        };
        return fallbacks[key] || key;
    };
    // Try to use next-intl translations, fall back to our defaults if it fails
    let t;
    try {
        t = useTranslations('CategoryTree');
    } catch (error) {
        t = defaultT;
    }
    // Define the unnamed category fallback early so it can be used in normalizeCategory
    const unnamedCategoryText = t('unnamedCategory');
    const [expandedMap, setExpandedMap] = useState({});
    const [searchQuery, setSearchQuery] = useState("");
    const [loading, setLoading] = useState(false);
    const [localCategories, setLocalCategories] = useState([]);
    const [error, setError] = useState(null);
    const [renderedCategories, setRenderedCategories] = useState([]);
    const searchInputRef = useRef(null);
    const treeContainerRef = useRef(null);
    // Use useCallback to prevent recreating this function on each render
    const normalizeCategory = useCallback((cat, parentId = null) => {
        // First ensure we have a valid category object
        if (!cat || typeof cat !== 'object') {
            return null;
        }
        return {
            ...cat,
            // Support both name and description fields, with fallback
            name: cat.name || cat.description || unnamedCategoryText,
            // Ensure subcategories array exists
            subcategories: Array.isArray(cat.subcategories) ? cat.subcategories : [],
            parentId: parentId,
        };
    }, [unnamedCategoryText]);
    // Initialize local categories when prop categories change
    useEffect(() => {
        if (categories && categories.length > 0) {
            const normalized = categories.map(cat => normalizeCategory(cat));
            setLocalCategories(normalized);
        } else {
            // Clear local categories when props are empty
            setLocalCategories([]);
        }
    }, [categories, normalizeCategory]);
    // Build rendered tree considering search and expansion
    useEffect(() => {
        const buildTree = async () => {
            // Don't try to process if we have no categories
            if (!localCategories || localCategories.length === 0) {
                setRenderedCategories([]);
                return;
            }
            // Filter categories based on search
            let filtered = localCategories;
            if (searchQuery) {
                filtered = filterCategories(localCategories, searchQuery);
            }
            const renderList = [];
            const processCategory = async (cat, level) => {
                if (!cat) return; // Protect against null/undefined
                renderList.push({...cat, level});
                if (expandedMap[cat.id]) {
                    try {
                        // Fetch subcategories if they haven't been fetched yet for this category object
                        let subcatsToProcess = cat.subcategories;
                        if (!cat.subcategoriesFetched && cat.subcategories.length === 0) {
                            // Set local loading state
                            setLoading(true);
                            try {
                                const fetchedSubcats = await fetchSubCategories(cat.id);
                                // Update the category in the main state
                                cat.subcategories = fetchedSubcats.map(sub => normalizeCategory(sub, cat.id));
                                cat.subcategoriesFetched = true; // Mark as fetched
                                subcatsToProcess = cat.subcategories;
                            } catch (err) {
                                cat.subcategoriesFetched = true; // Still mark as fetched to prevent repeated errors
                                cat.error = err.message;
                            } finally {
                                setLoading(false);
                            }
                        }
                        // Process subcategories if we have any
                        if (subcatsToProcess && subcatsToProcess.length > 0) {
                            for (const subcat of subcatsToProcess) {
                                await processCategory(subcat, level + 1);
                            }
                        }
                    } catch (err) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', err);
        }
    }
                }
            };
            // Process all top-level filtered categories
            for (const cat of filtered) {
                await processCategory(cat, 0);
            }
            setRenderedCategories(renderList);
        };
        buildTree();
    }, [localCategories, expandedMap, searchQuery, normalizeCategory]);
    // Toggle expand/collapse of a category
    const toggleExpand = useCallback((categoryId) => {
        setExpandedMap(prev => ({
            ...prev,
            [categoryId]: !prev[categoryId],
        }));
    }, []);
    // Handle selecting a category
    const handleSelect = useCallback((category) => {
        if (onSelectCategory) {
            onSelectCategory(category);
        }
    }, [onSelectCategory]);
    // Handle search input
    const handleSearchChange = useCallback((e) => {
        setSearchQuery(e.target.value);
    }, []);
    // Clear search
    const clearSearch = useCallback(() => {
        setSearchQuery("");
        if (searchInputRef.current) {
            searchInputRef.current.focus();
        }
    }, []);
    // Show error UI
    if (error) {
        return (
            <div className={styles.errorContainer}>
                <div className={styles.errorMessage}>{t('errorLoading')}</div>
                <button
                    className={styles.retryButton}
                    onClick={() => setError(null)}
                >
                    {t('retry')}
                </button>
            </div>
        );
    }
    // Determine if there are actually categories to show
    const hasCategories = localCategories && localCategories.length > 0;
    const hasFilteredCategories = renderedCategories && renderedCategories.length > 0;
    const isEmptySearch = searchQuery && !hasFilteredCategories;
    return (
        <div
            className={`${styles.container} ${modalContext ? styles.modalContext : ""}`}
            style={{
                maxHeight: `${maxHeight}px`, 
                backgroundColor: '#ffffff'
            }}
        >
            {/* Search bar */}
            {searchable && (
                <div className={styles.searchContainer}>
                    <Search size={16} className={styles.searchIcon}/>
                    <input
                        ref={searchInputRef}
                        className={styles.searchInput}
                        value={searchQuery}
                        onChange={handleSearchChange}
                        placeholder={t('searchPlaceholder')}
                        aria-label={t('searchPlaceholder')}
                        style={{backgroundColor: '#ffffff', color: '#1f2937'}}
                    />
                    {searchQuery && (
                        <button
                            className={styles.clearSearchButton}
                            onClick={clearSearch}
                            aria-label={t('clearSearchAria')}
                        >
                            <X size={16}/>
                        </button>
                    )}
                </div>
            )}
            {/* Main tree content */}
            <div
                ref={treeContainerRef}
                className={styles.treeContainer}
                style={{
                    maxHeight: searchable ? `${maxHeight - 40}px` : `${maxHeight}px`,
                    backgroundColor: '#ffffff',
                    color: '#1f2937'
                }}
            >
                {loading && (
                    <div className={styles.loadingOverlay}>
                        <Loader className={styles.spinner}/>
                        <span>{t('loading')}</span>
                    </div>
                )}
                {/* Empty state when no categories are available */}
                {!hasCategories && !loading && (
                    <div className={styles.emptyState}>
                        {t('empty')}
                    </div>
                )}
                {/* Empty search results */}
                {isEmptySearch && (
                    <div className={styles.noResults}>
                        {t('noResults', {query: searchQuery})}
                    </div>
                )}
                {/* Render the actual tree */}
                {hasFilteredCategories && (
                    <ul className={styles.categoriesList} style={{backgroundColor: '#ffffff'}}>
                        {renderedCategories.map((category) => (
                            <li
                                key={category.id}
                                className={`${styles.categoryItem} ${
                                    selectedCategoryId === category.id ? styles.selectedItem : ""
                                }`}
                                style={{
                                    paddingLeft: `${category.level * 16}px`,
                                    backgroundColor: '#ffffff'
                                }}
                            >
                                {/* Only show expand toggle if category has or might have children */}
                                {(category.subcategories?.length > 0 || 
                                   !category.subcategoriesFetched) && (
                                    <button
                                        type="button"
                                        className={styles.expandButton}
                                        onClick={() => toggleExpand(category.id)}
                                        aria-expanded={expandedMap[category.id]}
                                        style={{color: '#6b7280'}}
                                    >
                                        {expandedMap[category.id] ? (
                                            <ChevronDown size={16}/>
                                        ) : (
                                            <ChevronRight size={16}/>
                                        )}
                                    </button>
                                )}
                                {/* Leaf node without toggle */}
                                {(!category.subcategories || category.subcategories.length === 0) && 
                                  category.subcategoriesFetched && (
                                    <span className={styles.leafNode}></span>
                                )}
                                {/* Category name and selection */}
                                <button
                                    type="button"
                                    className={styles.categoryButton}
                                    onClick={() => handleSelect(category)}
                                    aria-selected={selectedCategoryId === category.id}
                                    style={{
                                        backgroundColor: selectedCategoryId === category.id ? '#f0f9ff' : '#ffffff',
                                        color: selectedCategoryId === category.id ? '#2980b9' : '#1f2937'
                                    }}
                                >
                                    {category.name || t('unnamedCategory')}
                                </button>
                            </li>
                        ))}
                    </ul>
                )}
            </div>
        </div>
    );
});
CategoryTree.displayName = 'CategoryTree';
CategoryTree.propTypes = {
    selectedCategoryId: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
    onSelectCategory: PropTypes.func.isRequired,
    maxHeight: PropTypes.number,
    modalContext: PropTypes.bool,
    searchable: PropTypes.bool,
    categories: PropTypes.array,
    categoryType: PropTypes.string
};
export default CategoryTree;