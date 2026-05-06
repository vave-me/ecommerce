"use client";
import React, { useState, useCallback, useEffect, useMemo, memo, useRef } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useTranslations, useLocale } from 'next-intl';
import { useRouter } from 'next/navigation';
import { 
    Filter, X, ChevronDown, Grid, Zap, Check
} from '@/icons';

// Import hooks and utilities
import { resetFilters, setFilters } from "../../redux/slices/listingFiltersSlice";
import { useIsMobile } from '../../hooks/useMobileDetection';
// Only import debugFilters in development
const debugFilters = process.env.NODE_ENV === 'development' ? require('../../utils/debugFilters').default : {
    logFilterState: () => {},
    logFilterFlow: () => {},
    compareFilters: () => {}
};

// Import categories hook
import { useMainCategories, useCategories, useSubCategories } from '../../hooks/useCategories';

// Import products API for fetching brands
import { fetchProductsByFilters } from '../../api/productsApi';

// Import UnifiedComposer for AI section
import Composer from '../Feed/Composer';

import styles from './HorizontalFilters.module.css';

/**
 * HorizontalFilters - Professional filter component that integrates with Redux
 */
const HorizontalFilters = memo(({ categoryType = "marketplace", onFiltersChange }) => {
    const t = useTranslations('Filters');
    const dispatch = useDispatch();
    const router = useRouter();
    const locale = useLocale();
    const isMobile = useIsMobile();
    
    // Redux state for filters - use specific selectors to prevent unnecessary re-renders
    const appliedCategories = useSelector(state => state.listingFilters.categories);
    const appliedMinPrice = useSelector(state => state.listingFilters.minPrice);
    const appliedMaxPrice = useSelector(state => state.listingFilters.maxPrice);
    const appliedCategory = useSelector(state => state.listingFilters.category);
    const appliedCategoryID = useSelector(state => state.listingFilters.categoryID);
    const appliedCategorySlug = useSelector(state => state.listingFilters.categorySlug);
    const sortBy = useSelector(state => state.listingFilters.sortBy);
    const sortOrder = useSelector(state => state.listingFilters.sortOrder);
    
    // State management
    const [showAiAssistant, setShowAiAssistant] = useState(true);
    const [activePopovers, setActivePopovers] = useState({});
    const [expandedGroup, setExpandedGroup] = useState(null);
    const [groupOrder, setGroupOrder] = useState(['all', 'marketplace', 'services', 'posts']);
    const [selectedEntity, setSelectedEntity] = useState('all'); // Track which entity filter is active
    const [showMobileFilters, setShowMobileFilters] = useState(false);
    const [localFilters, setLocalFilters] = useState({
        categories: [],
        categoryData: [], // Store full category info (id, slug, name)
        price: [0, 5000],
        brands: [],
        condition: '',
        rating: 0
    });
    const [availableBrands, setAvailableBrands] = useState([]);
    const [brandsLoading, setBrandsLoading] = useState(false);
    
    // Refs
    const containerRef = useRef(null);
    const brandFetchTimeout = useRef(null);
    
    // Fetch categories for all groups
    const { 
        data: marketplaceData, 
        isLoading: marketplaceLoading 
    } = useMainCategories('marketplace');
    
    const { 
        data: servicesData, 
        isLoading: servicesLoading 
    } = useMainCategories('service');
    
    const { 
        data: postsData, 
        isLoading: postsLoading 
    } = useMainCategories('posts');
    
    const mainCategoriesLoading = marketplaceLoading || servicesLoading || postsLoading;
    const mainCategoriesData = expandedGroup === 'services' ? servicesData : expandedGroup === 'posts' ? postsData : marketplaceData;
    
    // Fetch all categories for the tree based on expanded group
    const {
        data: allCategoriesData,
        isLoading: allCategoriesLoading
    } = useCategories(expandedGroup || 'marketplace');
    
    // Check if we have a selected category to show its subcategories
    // Use local state first to avoid re-renders from Redux changes
    const selectedCategoryId = localFilters.categoryData?.[0]?.id || appliedCategoryID;
    const selectedCategorySlug = localFilters.categories[0] || appliedCategorySlug || appliedCategory;
    
    // Fetch subcategories if a category is selected
    const {
        data: subCategoriesData,
        isLoading: subCategoriesLoading
    } = useSubCategories(selectedCategoryId, {
        enabled: !!selectedCategoryId
    });
    
    // Use categories from the expanded group for the popover
    const allCategories = expandedGroup === 'services' ? servicesData?.categories : 
                          expandedGroup === 'posts' ? postsData?.categories : 
                          marketplaceData?.categories || [];

    // Initialize local filters from Redux on mount only
    useEffect(() => {
        // Only run once on mount to set initial values
        if (appliedCategories?.length > 0 || appliedMinPrice || appliedMaxPrice) {
            setLocalFilters(prev => ({
                ...prev,
                categories: appliedCategories || [],
                price: [appliedMinPrice || 0, appliedMaxPrice || 5000]
            }));
        }
    }, []); // Empty dependency array - only run on mount

    // Prepare categories for display
    const displayCategories = useMemo(() => {
        if (!mainCategoriesData?.categories) return [];
        return mainCategoriesData.categories.slice(0, 8).map(cat => ({
            value: cat.value || cat.slug || cat.id,
            label: cat.label || cat.name,
            id: cat.id
        }));
    }, [mainCategoriesData]);
    
    // Fetch available brands only when Apply is clicked
    const fetchBrands = useCallback(async (category) => {
        if (!category) {
            setAvailableBrands([
                'Siemens', 'Schneider Electric', 'ABB', 'Zaptec', 
                'Philips', 'Osram', 'Bosch', 'Samsung'
            ]);
            return;
        }

        setBrandsLoading(true);
        try {
            const result = await fetchProductsByFilters({
                category: category,
                pageSize: 50,
                skipCache: true
            });
            
            const brands = new Set();
            if (result?.products) {
                result.products.forEach(product => {
                    if (product.brand?.trim()) {
                        brands.add(product.brand.trim());
                    }
                });
            }
            
            const brandArray = Array.from(brands).sort();
            setAvailableBrands(brandArray.length > 0 ? brandArray : [
                'Siemens', 'Schneider Electric', 'ABB', 'Zaptec'
            ]);
        } catch (error) {
            if (error.name !== 'CanceledError') {
                // Error: 'Brand fetch error:', error...
            }
            setAvailableBrands(['Siemens', 'Schneider Electric', 'ABB', 'Zaptec']);
        } finally {
            setBrandsLoading(false);
        }
    }, []);
    
    // Only fetch brands when filters are applied (from Redux)
    useEffect(() => {
        if (appliedCategory) {
            fetchBrands(appliedCategory);
        }
    }, [appliedCategory, fetchBrands]);

    // Popover handlers
    const togglePopover = useCallback((popoverId) => {
        setActivePopovers(prev => {
            const newState = {};
            Object.keys(prev).forEach(key => {
                newState[key] = key === popoverId ? !prev[key] : false;
            });
            if (!prev[popoverId]) {
                newState[popoverId] = true;
            }
            return newState;
        });
    }, []);
    
    const closeAllPopovers = useCallback(() => {
        setActivePopovers({});
    }, []);
    
    // Click outside handler
    useEffect(() => {
        const handleClickOutside = (event) => {
            if (containerRef.current && !containerRef.current.contains(event.target)) {
                closeAllPopovers();
            }
        };
        
        document.addEventListener('click', handleClickOutside);
        return () => document.removeEventListener('click', handleClickOutside);
    }, [closeAllPopovers]);
    
    // Category click handler (in filter dropdown)
    const handleCategoryClick = useCallback((categoryValue, categoryData = null) => {
        const isSelected = localFilters.categories.includes(categoryValue);
        const newCategories = isSelected
            ? localFilters.categories.filter(c => c !== categoryValue)
            : [...localFilters.categories, categoryValue];
        
        // Store category metadata for proper API calls
        const newCategoryData = isSelected
            ? localFilters.categoryData?.filter(cd => cd.value !== categoryValue) || []
            : [...(localFilters.categoryData || []), categoryData || { value: categoryValue }];
        
        setLocalFilters(prev => ({ 
            ...prev, 
            categories: newCategories,
            categoryData: newCategoryData
        }));
        
        // Don't update Redux immediately - wait for Apply button
    }, [localFilters.categories, localFilters.categoryData]);
    
    // Price handlers
    const handlePriceChange = useCallback((min, max) => {
        setLocalFilters(prev => ({ ...prev, price: [min, max] }));
        // Don't update Redux immediately - wait for Apply button
    }, []);
    
    // Brand handler
    const handleBrandChange = useCallback((brand) => {
        const isSelected = localFilters.brands.includes(brand);
        const newBrands = isSelected
            ? localFilters.brands.filter(b => b !== brand)
            : [...localFilters.brands, brand];
        
        setLocalFilters(prev => ({ ...prev, brands: newBrands }));
        // Don't update Redux immediately - wait for Apply button
    }, [localFilters.brands]);
    
    // Condition handler
    const handleConditionChange = useCallback((condition) => {
        setLocalFilters(prev => ({ ...prev, condition }));
        // Don't update Redux immediately - wait for Apply button
    }, []);
    
    // Rating handler
    const handleRatingChange = useCallback((rating) => {
        setLocalFilters(prev => ({ ...prev, rating }));
        // Don't update Redux immediately - wait for Apply button
    }, []);
    
    // Apply all filters
    const applyFilters = useCallback(() => {
        // Find the first selected category data
        const selectedCategoryData = localFilters.categoryData?.[0];
        
        const reduxFilters = {
            category: localFilters.categories[0] || '',
            categories: localFilters.categories,
            categoryID: selectedCategoryData?.id || '',
            categorySlug: selectedCategoryData?.slug || '',
            minPrice: localFilters.price[0] > 0 ? localFilters.price[0] : '',
            maxPrice: localFilters.price[1] < 5000 ? localFilters.price[1] : '',
            brands: localFilters.brands,
            brand: localFilters.brands[0] || '',
            condition: localFilters.condition,
            rating: localFilters.rating
        };
        
        debugFilters.logFilterState('HorizontalFilters.applyFilters', reduxFilters);
        
        dispatch(setFilters(reduxFilters));
        closeAllPopovers();
        
        // Call callback if provided
        if (onFiltersChange) {
            onFiltersChange(reduxFilters);
        }
    }, [localFilters, dispatch, closeAllPopovers, onFiltersChange]);    
    // Clear all filters
    const handleClearAll = useCallback(() => {
        const cleanFilters = {
            categories: [],
            categoryData: [],
            price: [0, 5000],
            brands: [],
            condition: '',
            rating: 0,
            // negotiable removed
            freeShipping: false,
            onlyPickup: false,
            verifiedSeller: false,
            isOnline: false,
            hasPortfolio: false
        };
        setLocalFilters(cleanFilters);
        
        // Clear Redux filters
        dispatch(resetFilters());
    }, [dispatch]);
    
    // Remove specific filter
    const handleRemoveFilter = useCallback((filterType, value) => {
        const updatedFilters = { ...localFilters };
        
        switch (filterType) {
            case 'category':
                updatedFilters.categories = localFilters.categories.filter(c => c !== value);
                updatedFilters.categoryData = localFilters.categoryData?.filter(cd => cd.value !== value) || [];
                break;
            case 'brand':
                updatedFilters.brands = localFilters.brands.filter(b => b !== value);
                break;
            case 'price':
                updatedFilters.price = [0, 5000];
                break;
            case 'condition':
                updatedFilters.condition = '';
                break;
            case 'rating':
                updatedFilters.rating = 0;
                break;
        }
        
        setLocalFilters(updatedFilters);
        // Don't update Redux immediately - wait for Apply button
    }, [localFilters]);
    
    // Sort handler - IMMEDIATE Redux dispatch (sorting doesn't need Apply button)
    const handleSortChange = useCallback((e) => {
        const value = e.target.value;
        let sortUpdate = {};
        
        switch (value) {
            case 'price_low':
                sortUpdate.sortBy = 'price';
                sortUpdate.sortOrder = 'asc';
                break;
            case 'price_high':
                sortUpdate.sortBy = 'price';
                sortUpdate.sortOrder = 'desc';
                break;
            case 'date_new':
                sortUpdate.sortBy = 'created_at';
                sortUpdate.sortOrder = 'desc';
                break;
            case 'rating_high':
                sortUpdate.sortBy = 'rating';
                sortUpdate.sortOrder = 'desc';
                break;
            case '':
                sortUpdate.sortBy = '';
                sortUpdate.sortOrder = '';
                break;
            default:
                sortUpdate.sortBy = value;
                sortUpdate.sortOrder = 'asc';
        }
        
        debugFilters.logFilterState('HorizontalFilters.handleSortChange', sortUpdate);
        dispatch(setFilters(sortUpdate));
    }, [dispatch]);

    // Check if more filters are active
    const hasMoreFilters = useMemo(() => {
        return false; // negotiable removed
    }, [localFilters]);
    
    // Check if any filters are active
    const hasActiveFilters = useMemo(() => {
        return localFilters.categories.length > 0 || 
               localFilters.price[0] !== 0 || 
               localFilters.price[1] !== 5000;
    }, [localFilters.categories.length, localFilters.price[0], localFilters.price[1]]);
    
    // Check if local filters differ from applied filters
    const hasUnappliedChanges = useMemo(() => {
        const priceChanged = localFilters.price[0] !== (appliedMinPrice || 0) || 
                           localFilters.price[1] !== (appliedMaxPrice || 5000);
        const categoriesChanged = JSON.stringify(localFilters.categories) !== JSON.stringify(appliedCategories || []);
        return priceChanged || categoriesChanged;
    }, [localFilters.price, localFilters.categories, appliedMinPrice, appliedMaxPrice, appliedCategories]);
    
    // Filter button states
    const filterButtonStates = useMemo(() => ({
        categories: localFilters.categories.length > 0,
        price: localFilters.price[0] !== 0 || localFilters.price[1] !== 5000
        // negotiable removed
        // userType removed
    }), [localFilters]);

    // Build category tree
    const buildCategoryTree = useCallback((categories) => {
        if (!categories) return [];
        return categories.map(cat => ({
            name: cat.name || cat.label || cat.description,     // Use name/label for display
            value: cat.slug || cat.value || cat.id,           // Use slug or value
            id: cat.id,                // Preserve the ID
            slug: cat.slug,            // Preserve the slug
            children: []
        }));
    }, []);

    // Render category tree item
    const renderCategoryTreeItem = (category) => {
        const isChecked = localFilters.categories.includes(category.value);
        
        return (
            <li key={category.value} className={styles.categoryTreeItem}>
                <div className={styles.categoryTreeItemContent}>
                    <input
                        type="checkbox"
                        id={`cat-${category.value}`}
                        className={styles.categoryCheckbox}
                        checked={isChecked}
                        onChange={() => handleCategoryClick(category.value, category)}
                    />
                    <label htmlFor={`cat-${category.value}`} className={styles.categoryLabel}>
                        {category.name}
                    </label>
                </div>
            </li>
        );
    };
    
    return (
        <>
            {/* Mobile filter modal */}
            {isMobile && showMobileFilters && (
                <div 
                    className={styles.mobileOverlay} 
                    onClick={() => setShowMobileFilters(false)}
                    role="dialog"
                    aria-modal="true"
                    aria-labelledby="mobile-filters-title"
                >
                    <div className={styles.mobileFilterModal} onClick={(e) => e.stopPropagation()}>
                        <div className={styles.mobileFilterHeader}>
                            <h3 id="mobile-filters-title">{t('filters', 'Filters')}</h3>
                            <button 
                                onClick={() => setShowMobileFilters(false)}
                                aria-label={t('close', 'Close filters')}
                                type="button"
                            >
                                ✕
                            </button>
                        </div>
                        <div className={styles.mobileFilterContent}>
                            {/* Categories Filter */}
                            <div className={styles.mobileFilterGroup}>
                                <h4>{t('categories', 'Categories')}</h4>
                                <div className={styles.mobileFilterOptions}>
                                    {displayCategories.map(cat => (
                                        <label key={cat.value} className={styles.mobileFilterItem}>
                                            <input
                                                type="checkbox"
                                                checked={localFilters.categories.includes(cat.value)}
                                                onChange={() => handleCategoryClick(cat.value, cat)}
                                            />
                                            <span>{cat.label}</span>
                                        </label>
                                    ))}
                                </div>
                            </div>
                            
                            {/* Price Filter */}
                            <div className={styles.mobileFilterGroup}>
                                <h4>{t('price', 'Price Range')}</h4>
                                <div className={styles.priceInputs}>
                                    <input
                                        type="number"
                                        placeholder="Min"
                                        value={localFilters.price[0]}
                                        onChange={(e) => handlePriceChange(parseInt(e.target.value) || 0, localFilters.price[1])}
                                        className={styles.priceInput}
                                    />
                                    <span>-</span>
                                    <input
                                        type="number"
                                        placeholder="Max"
                                        value={localFilters.price[1]}
                                        onChange={(e) => handlePriceChange(localFilters.price[0], parseInt(e.target.value) || 5000)}
                                        className={styles.priceInput}
                                    />
                                </div>
                            </div>

                            {/* Brands Filter */}
                            {availableBrands.length > 0 && (
                                <div className={styles.mobileFilterGroup}>
                                    <h4>{t('brands', 'Brands')}</h4>
                                    <div className={styles.mobileFilterOptions}>
                                        {availableBrands.slice(0, 6).map(brand => (
                                            <label key={brand} className={styles.mobileFilterItem}>
                                                <input
                                                    type="checkbox"
                                                    checked={localFilters.brands.includes(brand)}
                                                    onChange={() => handleBrandChange(brand)}
                                                />
                                                <span>{brand}</span>
                                            </label>
                                        ))}
                                    </div>
                                </div>
                            )}

                            {/* Condition Filter */}
                            <div className={styles.mobileFilterGroup}>
                                <h4>{t('condition', 'Condition')}</h4>
                                <div className={styles.mobileFilterOptions}>
                                    {['new', 'used', 'refurbished'].map(condition => (
                                        <label key={condition} className={styles.mobileFilterItem}>
                                            <input
                                                type="radio"
                                                name="condition"
                                                checked={localFilters.condition === condition}
                                                onChange={() => handleConditionChange(condition)}
                                            />
                                            <span>{t(condition, condition)}</span>
                                        </label>
                                    ))}
                                </div>
                            </div>
                        </div>
                        <div className={styles.mobileFilterActions}>
                            <button onClick={handleClearAll} className={styles.clearButton}>
                                {t('clearAll', 'Clear All')}
                            </button>
                            <button 
                                onClick={() => { 
                                    applyFilters(); 
                                    setShowMobileFilters(false); 
                                }} 
                                className={styles.applyButton}
                            >
                                {t('apply', 'Apply Filters')}
                            </button>
                        </div>
                    </div>
                </div>
            )}

            <div 
                className={styles.container} 
                ref={containerRef}
                onClick={isMobile ? (e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    setShowMobileFilters(true);
                } : undefined}
                role={isMobile ? "button" : undefined}
                tabIndex={isMobile ? 0 : undefined}
                aria-label={isMobile ? t('filters', 'Open filters') : undefined}
                aria-expanded={isMobile ? showMobileFilters : undefined}
                aria-haspopup={isMobile ? "dialog" : undefined}
                onKeyDown={isMobile ? (e) => {
                    if (e.key === 'Enter' || e.key === ' ') {
                        e.preventDefault();
                        setShowMobileFilters(true);
                    }
                } : undefined}
            >
            <div className={styles.containerInner}>
                {/* Mobile Filter Icon - Only render on mobile */}
                {isMobile && <Filter size={20} className={styles.mobileFilterIcon} />}
                
                {/* Main Category Bar with Group Labels */}
                <div className={styles.categoryBar}>
                    <div className={styles.categoryBarInner}>
                        {/* ALL Label - Always first */}
                        <button
                            className={`${styles.entityLabel} ${selectedEntity === 'all' ? styles.activeEntity : ''}`}
                            onClick={() => {
                                setSelectedEntity('all');
                                // Dispatch filter to show all entity types
                                dispatch(setFilters({ 
                                    entityTypes: ['product', 'post', 'service'],
                                    contentType: 'all' 
                                }));
                            }}
                        >
                            {t('all', 'ALL')}
                        </button>
                        
                        {/* Group Labels */}
                        {groupOrder.filter(g => g !== 'all').map((groupName) => (
                            <React.Fragment key={groupName}>
                                <button
                                    className={`${styles.groupLabel} ${expandedGroup === groupName ? styles.activeGroup : ''}`}
                                    onClick={() => {
                                        // Toggle if clicking the same group, otherwise expand new group
                                        if (expandedGroup === groupName) {
                                            setExpandedGroup(null);
                                        } else {
                                            setExpandedGroup(groupName);
                                            // Move clicked group to first position after 'all'
                                            const currentOrder = groupOrder.filter(g => g !== 'all');
                                            if (currentOrder[0] !== groupName) {
                                                const newOrder = ['all', groupName, ...currentOrder.filter(g => g !== groupName)];
                                                setGroupOrder(newOrder);
                                            }
                                        }
                                        // Set selected entity to the group when clicked
                                        setSelectedEntity(groupName);
                                        
                                        // Set entity types based on selected group
                                        const entityTypeMap = {
                                            'marketplace': ['product'],
                                            'services': ['service'],
                                            'posts': ['post']
                                        };
                                        
                                        dispatch(setFilters({ 
                                            entityTypes: entityTypeMap[groupName] || ['product'],
                                            contentType: groupName
                                        }));
                                    }}
                                >
                                    {t(groupName, groupName.charAt(0).toUpperCase() + groupName.slice(1))}
                                </button>
                                
                                {expandedGroup === groupName && (
                                    <>
                                        {displayCategories.map(cat => {
                                            // Determine the route based on the group
                                            const routePrefix = groupName === 'marketplace' ? 'products' : groupName === 'services' ? 'services' : 'posts';
                                            
                                            return (
                                                <button
                                                    key={cat.value}
                                                    className={`${styles.categoryLink} ${appliedCategory === cat.value ? styles.active : ''}`}
                                                    onClick={() => {
                                                        router.push(`/${locale}/${routePrefix}/${cat.value}`);
                                                        // Keep the group selected when a category is clicked
                                                        setSelectedEntity(groupName);
                                                    }}
                                                >
                                                    {cat.label}
                                                </button>
                                            );
                                        })}
                                    </>
                                )}
                            </React.Fragment>
                        ))}
                    </div>
                </div>
                
                {/* Filter Section */}
                <div className={styles.filterSection}>
                    <div className={styles.filterControls}>
                        {/* Filter buttons */}
                        <div className={styles.filterButtons}>
                            {/* Categories Filter */}
                            <div className={styles.popoverContainer}>
                                <button
                                    className={`${styles.filterButton} ${filterButtonStates.categories ? styles.active : ''}`}
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        togglePopover('categories');
                                    }}
                                >
                                    {selectedCategoryId && subCategoriesData?.length > 0 ? 
                                        t('subcategories', 'Subcategories') : 
                                        expandedGroup ? t(expandedGroup, expandedGroup.charAt(0).toUpperCase() + expandedGroup.slice(1)) + ' ' + t('categories', 'Categories') : t('allCategories')}
                                    <ChevronDown className={styles.chevronIcon} />
                                </button>
                                {activePopovers.categories && (
                                    <div className={`${styles.popover} ${styles.visible}`}>
                                        <div className={styles.popoverContent}>
                                            <div className={styles.categoryTree}>
                                                <ul>
                                                    {/* Show subcategories if a category is selected, otherwise show main categories */}
                                                    {selectedCategoryId && subCategoriesData?.length > 0 ? (
                                                        <>
                                                            {/* Back button to main categories */}
                                                            <li className={styles.categoryBackButton}>
                                                                <button
                                                                    onClick={() => {
                                                                        // Clear category selection to go back to main categories
                                                                        setLocalFilters(prev => ({ 
                                                                            ...prev, 
                                                                            categories: [],
                                                                            categoryData: []
                                                                        }));
                                                                        // Apply the change immediately
                                                                        setTimeout(() => {
                                                                            dispatch(setFilters({
                                                                                categories: [],
                                                                                category: '',
                                                                                categoryID: '',
                                                                                categorySlug: ''
                                                                            }));
                                                                        }, 0);
                                                                    }}
                                                                    className={styles.backButton}
                                                                >
                                                                    ← {t('backToCategories', 'Back to Categories')}
                                                                </button>
                                                            </li>
                                                            {/* Selected category name */}
                                                            <li className={styles.selectedCategoryHeader}>
                                                                {allCategories?.find(c => c.id === selectedCategoryId)?.name || selectedCategorySlug}
                                                            </li>
                                                            {/* Subcategories */}
                                                            {buildCategoryTree(subCategoriesData || []).map(cat => renderCategoryTreeItem(cat))}
                                                        </>
                                                    ) : (
                                                        /* Main categories */
                                                        buildCategoryTree(allCategories || []).map(cat => renderCategoryTreeItem(cat))
                                                    )}
                                                </ul>
                                            </div>
                                            <button
                                                type="button"
                                                onClick={applyFilters}
                                                className={styles.applyFilterButton}
                                            >
                                                {t('apply', 'Apply')}
                                            </button>
                                        </div>
                                    </div>
                                )}
                            </div>
                            
                            {/* Price Filter */}
                            <div className={styles.popoverContainer}>
                                <button
                                    className={`${styles.filterButton} ${filterButtonStates.price ? styles.active : ''}`}
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        togglePopover('price');
                                    }}
                                >
                                    {t('price')}
                                    <ChevronDown className={styles.chevronIcon} />
                                </button>
                                {activePopovers.price && (
                                    <div className={`${styles.popover} ${styles.visible}`}>
                                        <div className={styles.popoverContent}>
                                            <div className={styles.priceLabels}>
                                                <span>€{localFilters.price[0]}</span>
                                                <span>€{localFilters.price[1]}</span>
                                            </div>
                                            <div className={styles.priceSliderContainer}>
                                                <div className={styles.sliderTrack} />
                                                <div 
                                                    className={styles.sliderRange}
                                                    style={{
                                                        left: `${(localFilters.price[0] / 5000) * 100}%`,
                                                        width: `${((localFilters.price[1] - localFilters.price[0]) / 5000) * 100}%`
                                                    }}
                                                />
                                                <input
                                                    type="range"
                                                    className={styles.priceSlider}
                                                    min="0"
                                                    max="5000"
                                                    value={localFilters.price[0]}
                                                    onChange={(e) => {
                                                        const val = parseInt(e.target.value);
                                                        if (val < localFilters.price[1]) {
                                                            handlePriceChange(val, localFilters.price[1]);
                                                        }
                                                    }}
                                                />
                                                <input
                                                    type="range"
                                                    className={styles.priceSlider}
                                                    min="0"
                                                    max="5000"
                                                    value={localFilters.price[1]}
                                                    onChange={(e) => {
                                                        const val = parseInt(e.target.value);
                                                        if (val > localFilters.price[0]) {
                                                            handlePriceChange(localFilters.price[0], val);
                                                        }
                                                    }}
                                                />
                                            </div>
                                            <button
                                                type="button"
                                                onClick={applyFilters}
                                                className={styles.applyFilterButton}
                                            >
                                                {t('apply', 'Apply')}
                                            </button>
                                        </div>
                                    </div>
                                )}
                            </div>
                            
                            {/* User Type Filter removed - not used */}
                            
                            {/* Negotiable Toggle removed - not used */}
                            {/* Apply All Filters Button - visible when there are unapplied changes */}
                            {hasUnappliedChanges && (
                                <button
                                    type="button"
                                    onClick={applyFilters}
                                    className={styles.applyAllButton}
                                >
                                    <Check className={styles.applyIcon} />
                                    {t('applyFilters', 'Apply Filters')}
                                </button>
                            )}
                        </div>
                        
                        {/* Actions */}
                        <div className={styles.actions}>
                            <select
                                className={styles.sortSelect}
                                onChange={handleSortChange}
                                value={
                                    sortBy === 'price' && sortOrder === 'asc' ? 'price_low' :
                                    sortBy === 'price' && sortOrder === 'desc' ? 'price_high' :
                                    sortBy === 'created_at' && sortOrder === 'desc' ? 'date_new' :
                                    sortBy === 'rating' && sortOrder === 'desc' ? 'rating_high' :
                                    sortBy || ""
                                }
                            >
                                <option value="">{t('sortByRelevance') || 'Relevance'}</option>
                                <option value="price_low">{t('sortByPriceLow') || 'Price: Low to High'}</option>
                                <option value="price_high">{t('sortByPriceHigh') || 'Price: High to Low'}</option>
                                <option value="date_new">{t('sortByNewest') || 'Newest First'}</option>
                                <option value="rating_high">{t('sortByRatingHigh') || 'Highest Rated'}</option>
                            </select>
                            
                            {!showAiAssistant && !isMobile && (
                                <button
                                    className={styles.aiButton}
                                    onClick={() => setShowAiAssistant(true)}
                                    aria-label="Open AI Assistant"
                                    title="Open AI Assistant"
                                >
                                    <Zap className={styles.aiIcon} />
                                </button>
                            )}
                        </div>
                    </div>
                    
                    {/* Filter Pills - Show applied filters from Redux */}
                    {(appliedCategories?.length > 0 || appliedMinPrice || appliedMaxPrice) && (
                        <div className={styles.filterPills}>
                            <span className={styles.filterPillsLabel}>{t('active')}:</span>
                            {appliedCategories?.map(cat => (
                                <span key={cat} className={styles.filterPill}>
                                    {allCategories?.find(c => c.slug === cat || c.id === cat)?.name || cat}
                                    <button
                                        className={styles.filterPillRemove}
                                        onClick={() => {
                                            handleRemoveFilter('category', cat);
                                            // Apply immediately when removing from active filters
                                            setTimeout(() => applyFilters(), 0);
                                        }}
                                    >
                                        ×
                                    </button>
                                </span>
                            ))}
                            {(appliedMinPrice || appliedMaxPrice) && (
                                <span className={styles.filterPill}>
                                    €{appliedMinPrice || 0} - €{appliedMaxPrice || 5000}
                                    <button
                                        className={styles.filterPillRemove}
                                        onClick={() => {
                                            handleRemoveFilter('price');
                                            // Apply immediately when removing from active filters
                                            setTimeout(() => applyFilters(), 0);
                                        }}
                                    >
                                        ×
                                    </button>
                                </span>
                            )}
                            {/* Negotiable pill removed - not used */}
                            {/* User Type pill removed - not used */}
                            <button className={styles.clearAllButton} onClick={handleClearAll}>
                                {t('clearAll')}
                            </button>
                        </div>
                    )}
                </div>
                
                {/* AI Assistant Section - Desktop only */}
                {showAiAssistant && !isMobile && (
                    <Composer 
                        onFilterUpdate={(filters) => {
                            setLocalFilters(prev => ({ ...prev, ...filters }));
                            applyFilters();
                        }}
                        onClose={() => setShowAiAssistant(false)}
                    />
                )}
            </div>
        </div>
        </>
    );
});

HorizontalFilters.displayName = 'HorizontalFilters';

export default HorizontalFilters;