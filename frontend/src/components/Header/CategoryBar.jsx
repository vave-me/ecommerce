"use client";
import React, { useEffect, useState, useRef, memo, useCallback } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useRouter, usePathname } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { setFilters } from '../../redux/slices/listingFiltersSlice';
import { useMainCategories } from '../../hooks/useCategories';
import styles from './CategoryBar.module.css';

/**
 * Enhanced CategoryBar Component - Multiple category groups with labels
 * Supports Marketplace, Services, and Posts categories
 */
const CategoryBar = memo(({ categoryType = "marketplace" }) => {
    const t = useTranslations('Categories');
    const dispatch = useDispatch();
    const router = useRouter();
    const pathname = usePathname();
    const scrollContainerRef = useRef(null);
    
    // Redux state
    const filters = useSelector(state => state.listingFilters);
    const selectedCategory = filters?.category || '';
    
    // State for scroll buttons visibility
    const [showLeftScroll, setShowLeftScroll] = useState(false);
    const [showRightScroll, setShowRightScroll] = useState(false);
    
    // State for expanded category group
    const [expandedGroup, setExpandedGroup] = useState('marketplace');
    
    // Fetch categories for all groups
    const { 
        data: marketplaceData, 
        isLoading: marketplaceLoading,
        error: marketplaceError 
    } = useMainCategories('marketplace');
    
    const { 
        data: servicesData, 
        isLoading: servicesLoading,
        error: servicesError 
    } = useMainCategories('services');
    
    const { 
        data: postsData, 
        isLoading: postsLoading,
        error: postsError 
    } = useMainCategories('posts');
    
    // Check scroll position to show/hide scroll buttons
    const checkScroll = useCallback(() => {
        if (!scrollContainerRef.current) return;
        
        const { scrollLeft, scrollWidth, clientWidth } = scrollContainerRef.current;
        setShowLeftScroll(scrollLeft > 0);
        setShowRightScroll(scrollLeft < scrollWidth - clientWidth - 5);
    }, []);
    
    // Scroll handlers
    const scrollLeft = () => {
        if (!scrollContainerRef.current) return;
        scrollContainerRef.current.scrollBy({ 
            left: -200, 
            behavior: 'smooth' 
        });
    };
    
    const scrollRight = () => {
        if (!scrollContainerRef.current) return;
        scrollContainerRef.current.scrollBy({ 
            left: 200, 
            behavior: 'smooth' 
        });
    };
    
    // Handle category click
    const handleCategoryClick = (categoryValue, groupType) => {
        // Toggle category selection
        const newCategory = selectedCategory === categoryValue ? '' : categoryValue;
        
        // Update Redux state with both category and group type
        dispatch(setFilters({ 
            category: newCategory,
            categoryType: groupType 
        }));
    };
    
    // Handle group label click
    const handleGroupClick = (groupName) => {
        if (expandedGroup === groupName) {
            // Don't collapse if clicking the same group
            return;
        }
        setExpandedGroup(groupName);
        // Reset scroll position when switching groups
        if (scrollContainerRef.current) {
            scrollContainerRef.current.scrollLeft = 0;
        }
        // Check scroll after a small delay to ensure DOM is updated
        setTimeout(checkScroll, 100);
    };
    
    // Check scroll on mount and resize
    useEffect(() => {
        checkScroll();
        window.addEventListener('resize', checkScroll);
        
        return () => {
            window.removeEventListener('resize', checkScroll);
        };
    }, [checkScroll]);
    
    // Check scroll when expanded group changes
    useEffect(() => {
        const timeoutId = setTimeout(checkScroll, 100);
    
    return () => clearTimeout(timeoutId);
  }, [expandedGroup, checkScroll]);
    
    // Determine which categories to show based on expanded group
    const getExpandedCategories = () => {
        switch (expandedGroup) {
            case 'marketplace':
                return marketplaceData?.categories || [];
            case 'services':
                return servicesData?.categories || [];
            case 'posts':
                return postsData?.categories || [];
            default:
                return [];
        }
    };
    
    const isLoading = marketplaceLoading || servicesLoading || postsLoading;
    const categories = getExpandedCategories();
    
    if (isLoading) {
        return (
            <div className={styles.container}>
                <div className={styles.loadingState}>
                    {[...Array(8)].map((_, i) => (
                        <div key={i} className={styles.skeletonItem} />
                    ))}
                </div>
            </div>
        );
    }
    
    return (
        <div className={styles.container}>
            <div className={styles.wrapper}>
                {/* Left scroll button */}
                {showLeftScroll && (
                    <button 
                        className={`${styles.scrollButton} ${styles.scrollButtonLeft}`}
                        onClick={scrollLeft}
                        aria-label={t('scrollLeft')}
                    >
                        <svg className={styles.scrollIcon} viewBox="0 0 20 20" fill="currentColor">
                            <path fillRule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clipRule="evenodd" />
                        </svg>
                    </button>
                )}
                
                {/* Categories scroll container */}
                <div 
                    ref={scrollContainerRef}
                    className={styles.scrollContainer}
                    onScroll={checkScroll}
                >
                    <div className={styles.categoryList}>
                        {/* Marketplace Group */}
                        <button
                            className={`${styles.groupLabel} ${expandedGroup === 'marketplace' ? styles.activeGroup : ''}`}
                            onClick={() => handleGroupClick('marketplace')}
                            aria-expanded={expandedGroup === 'marketplace'}
                        >
                            {t('marketplace', 'Marketplace')}
                        </button>
                        
                        {/* Marketplace Categories */}
                        {expandedGroup === 'marketplace' && (
                            <>
                                {/* All categories button */}
                                <button
                                    className={`${styles.categoryLink} ${!selectedCategory ? styles.active : ''}`}
                                    onClick={() => handleCategoryClick('', 'marketplace')}
                                    aria-current={!selectedCategory ? 'page' : undefined}
                                >
                                    {t('allCategories', 'All')}
                                </button>
                                
                                {/* Category buttons */}
                                {categories.map((category) => {
                                    const categoryValue = category.slug || category.id;
                                    const categoryLabel = category.name || category.label;
                                    const isActive = selectedCategory === categoryValue;
                                    
                                    return (
                                        <button
                                            key={categoryValue}
                                            className={`${styles.categoryLink} ${isActive ? styles.active : ''}`}
                                            onClick={() => handleCategoryClick(categoryValue, 'marketplace')}
                                            aria-current={isActive ? 'page' : undefined}
                                            title={categoryLabel}
                                        >
                                            {categoryLabel}
                                        </button>
                                    );
                                })}
                                
                                {/* End label for Marketplace */}
                                <span className={styles.groupEndLabel}>
                                    {t('marketplace', 'Marketplace')}
                                </span>
                            </>
                        )}
                        
                        {/* Services Group */}
                        <button
                            className={`${styles.groupLabel} ${expandedGroup === 'services' ? styles.activeGroup : ''}`}
                            onClick={() => handleGroupClick('services')}
                            aria-expanded={expandedGroup === 'services'}
                        >
                            {t('services', 'Services')}
                        </button>
                        
                        {/* Services Categories */}
                        {expandedGroup === 'services' && (
                            <>
                                {/* All categories button */}
                                <button
                                    className={`${styles.categoryLink} ${!selectedCategory ? styles.active : ''}`}
                                    onClick={() => handleCategoryClick('', 'services')}
                                    aria-current={!selectedCategory ? 'page' : undefined}
                                >
                                    {t('allCategories', 'All')}
                                </button>
                                
                                {/* Category buttons */}
                                {categories.map((category) => {
                                    const categoryValue = category.slug || category.id;
                                    const categoryLabel = category.name || category.label;
                                    const isActive = selectedCategory === categoryValue;
                                    
                                    return (
                                        <button
                                            key={categoryValue}
                                            className={`${styles.categoryLink} ${isActive ? styles.active : ''}`}
                                            onClick={() => handleCategoryClick(categoryValue, 'services')}
                                            aria-current={isActive ? 'page' : undefined}
                                            title={categoryLabel}
                                        >
                                            {categoryLabel}
                                        </button>
                                    );
                                })}
                                
                                {/* End label for Services */}
                                <span className={styles.groupEndLabel}>
                                    {t('services', 'Services')}
                                </span>
                            </>
                        )}
                        
                        {/* Posts Group */}
                        <button
                            className={`${styles.groupLabel} ${expandedGroup === 'posts' ? styles.activeGroup : ''}`}
                            onClick={() => handleGroupClick('posts')}
                            aria-expanded={expandedGroup === 'posts'}
                        >
                            {t('posts', 'Posts')}
                        </button>
                        
                        {/* Posts Categories */}
                        {expandedGroup === 'posts' && (
                            <>
                                {/* All categories button */}
                                <button
                                    className={`${styles.categoryLink} ${!selectedCategory ? styles.active : ''}`}
                                    onClick={() => handleCategoryClick('', 'posts')}
                                    aria-current={!selectedCategory ? 'page' : undefined}
                                >
                                    {t('allCategories', 'All')}
                                </button>
                                
                                {/* Category buttons */}
                                {categories.map((category) => {
                                    const categoryValue = category.slug || category.id;
                                    const categoryLabel = category.name || category.label;
                                    const isActive = selectedCategory === categoryValue;
                                    
                                    return (
                                        <button
                                            key={categoryValue}
                                            className={`${styles.categoryLink} ${isActive ? styles.active : ''}`}
                                            onClick={() => handleCategoryClick(categoryValue, 'posts')}
                                            aria-current={isActive ? 'page' : undefined}
                                            title={categoryLabel}
                                        >
                                            {categoryLabel}
                                        </button>
                                    );
                                })}
                                
                                {/* End label for Posts */}
                                <span className={styles.groupEndLabel}>
                                    {t('posts', 'Posts')}
                                </span>
                            </>
                        )}
                    </div>
                </div>
                
                {/* Right scroll button */}
                {showRightScroll && (
                    <button 
                        className={`${styles.scrollButton} ${styles.scrollButtonRight}`}
                        onClick={scrollRight}
                        aria-label={t('scrollRight')}
                    >
                        <svg className={styles.scrollIcon} viewBox="0 0 20 20" fill="currentColor">
                            <path fillRule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clipRule="evenodd" />
                        </svg>
                    </button>
                )}
            </div>
        </div>
    );
});

CategoryBar.displayName = 'CategoryBar';

export default CategoryBar;