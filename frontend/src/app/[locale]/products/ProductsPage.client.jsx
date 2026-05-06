"use client";
import React, { memo, useEffect, useState } from 'react';
import styles from './ProductsPage.module.css';
// import Leftside from '../../../components/Leftside/Leftside'; // Removed per user request
import { useIsMobile } from '../../../hooks/useMobileDetection';
import { useSelector, useDispatch } from 'react-redux';
// Import the actual component from classified directory
import ClassifiedCard from "../../../components/classified/ClassifiedCard";
import ExpandableVaveAds from '../../../components/ads/ExpandableVaveAds';
import { Flame } from '@/icons';
import { selectShowComposerOnProducts } from '../../../redux/slices/uiPreferencesSlice';
import { useFilteredSearch } from '../../../hooks/useFilteredSearch';
import { setFilters } from '../../../redux/slices/listingFiltersSlice';
/**
 * ProductsPageClient - Client component for the products listing page
 *
 * @param {Object} props - Component props
 * @param {Array} props.serverProducts - Products data fetched from the server
 * @param {Object} props.serverFilters - Filters data fetched from the server
 * @param {Array} props.availableCategories - Available categories for navigation
 * @param {string} props.fetchError - Error message if any
 * @param {string} props.currentCategory - Current category slug
 * @param {string} props.categoryName - Current category display name
 * @param {string} props.locale - Current locale
 * @param {Object} props.labels - Translated labels
 * @returns {JSX.Element} Rendered component
 */
const ProductsPageClient = memo(function ProductsPageClient({
    serverProducts = [], 
    serverFilters = {},
    availableCategories = [],
    fetchError = null,
    currentCategory = null,
    categoryName = null,
    locale = 'en',
    labels = {},
    totalPages = 0,
    totalCount = 0,
    currentPage = 1
}) {
    // Hooks
    const isMobile = useIsMobile();
    const dispatch = useDispatch();
    const listingFilters = useSelector((state) => state.listingFilters);
    const {displayMode} = listingFilters;
    
    // State for products
    const [products, setProducts] = useState(serverProducts);
    
    // Track if category has been synced to Redux
    const [categoryReady, setCategoryReady] = useState(false);
    
    // Debug logging in development
    useEffect(() => {
        if (process.env.NODE_ENV === 'development') {
            
        }
    }, []);
    
    // Sync category from URL with Redux state
    useEffect(() => {
        if (currentCategory) {
            // Always set the category when we have one
            dispatch(setFilters({ 
                ...listingFilters,
                category: currentCategory 
            }));
            setCategoryReady(true);
            
            if (process.env.NODE_ENV === 'development') {
                
            }
        } else {
            // No category means we're on the main products page
            setCategoryReady(true);
        }
        
        // Clear category filter when component unmounts
        return () => {
            if (currentCategory) {
                dispatch(setFilters({ 
                    ...listingFilters,
                    category: '' 
                }));
            }
        };
    }, [currentCategory, dispatch]); // Fixed dependencies
    
    // Use filtered search hook to get products based on Redux filters
    const { 
        data: searchResults, 
        isLoading, 
        error: searchError 
    } = useFilteredSearch({
        entityType: 'product',
        enabled: categoryReady, // Only enable after category is set
        onSuccess: (data) => {
            // Update products when search results change
            if (data?.products) {
                setProducts(data.products);
            }
        }
    });
    
    // Debug logging for search results
    useEffect(() => {
        if (process.env.NODE_ENV === 'development') {
            
        }
    }, [searchResults, isLoading, searchError, categoryReady, currentCategory, listingFilters]);
    
    // Update products when search results change
    useEffect(() => {
        // If we have a category filter and search is complete, always use search results
        if (currentCategory && categoryReady && !isLoading && searchResults) {
            // Use search results when filtering by category
            setProducts(searchResults.products || []);
            
            if (process.env.NODE_ENV === 'development') {
                
            }
        } else if (!currentCategory && searchResults?.products) {
            // For the main products page (no category), update with search results
            setProducts(searchResults.products);
        }
    }, [searchResults, isLoading, currentCategory, categoryReady]);
    
    // Error state handling
    if (fetchError || searchError) {
        return (
            <div className={styles.errorState}>
                <h2>Error Loading Products</h2>
                <p>{fetchError || searchError?.message || 'An error occurred'}</p>
            </div>
        );
    }
    
    // Loading state
    if (isLoading && !products.length) {
        return (
            <div className={styles.loadingState}>
                <div className={styles.spinner}></div>
                <p>Loading products...</p>
            </div>
        );
    }
    
    // Empty state handling
    if (!products.length && !isLoading) {
        return (
            <div className={styles.emptyState}>
                {labels.empty || "No products found matching your filters."}
            </div>
        );
    }
    
    return (
        <div className={styles.container}>
            <main className={styles.mainContent}>
                <GridLayout
                    isMobile={isMobile}
                    products={products}
                    locale={locale}
                    availableCategories={availableCategories}
                    isLoading={isLoading}
                />
            </main>
        </div>
    );
});
/**
 * GridLayout - Component for grid view of products
 */
const GridLayout = memo(function GridLayout({isMobile, products, locale, availableCategories, isLoading}) {
    const showComposerOnProducts = useSelector(selectShowComposerOnProducts);
    
    return (
        <div className={styles.layoutGrid}>
            <section className={styles.contentArea}>
                {/* Top Section: Expandable Ads (includes UnifiedComposer) */}
                {showComposerOnProducts && (
                    <div className={styles.composerSection}>
                        <ExpandableVaveAds />
                    </div>
                )}
                
                {/* Trending Section (Desktop only) */}
                {!isMobile && (
                    <div className={styles.trendingSection}>
                        <div className={styles.trendingCard}>
                            <div className={styles.trendingHeader}>
                                <Flame size={14} strokeWidth={1.5} />
                                <h3>Trending Now</h3>
                            </div>
                            <p className={styles.trendingDescription}>
                                Today's popular videos, deals, and searches.
                            </p>
                        </div>
                    </div>
                )}
                
                {/* Loading overlay */}
                {isLoading && (
                    <div className={styles.loadingOverlay}>
                        <div className={styles.spinnerSmall}></div>
                    </div>
                )}
                
                {/* Products List below */}
                <ul className={styles.productList}>
                    {products.map((prod) => (
                        <li key={prod.id}>
                            <ClassifiedCard 
                                product={prod} 
                                locale={locale}
                                availableCategories={availableCategories}
                            />
                        </li>
                    ))}
                </ul>
            </section>
        </div>
    );
});
export default ProductsPageClient; 