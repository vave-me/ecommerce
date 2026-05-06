"use client";
import React, { memo } from 'react';
import styles from './ProductsPage.module.css';
// Leftside removed - using HorizontalFilters mobile modal instead
import { useIsMobile } from '../../../hooks/useMobileDetection';
// Removed context providers - using direct hooks
import FeedDisplay from '../../../components/Feed/FeedDisplay';
import { Flame } from '@/icons';

/**
 * Enhanced Products Page with AI integration
 * Replaces the standard UnifiedComposer with EnhancedUnifiedComposer
 */
const EnhancedProductsPage = memo(function EnhancedProductsPage({
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
  currentPage = 1,
  initialHasMore = true
}) {
  const isMobile = useIsMobile();

  // Configure initial params for products
  const initialParams = {
    feedType: "latest",
    entityTypes: ["product"],
    contentType: "products",
    page: currentPage,
    pageSize: 20,
    category: currentCategory,
    ...serverFilters
  };

  // Error state handling
  if (fetchError) {
    return (
      <div className={styles.errorState}>
        <h2>Error Loading Products</h2>
        <p>{fetchError}</p>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <main className={styles.mainContent}>
        <div className={styles.layoutGrid}>
          {/* Leftside removed - using HorizontalFilters mobile modal instead */}
          <section className={styles.contentArea}>
            {/* Trending section for desktop */}
            {!isMobile && (
              <div className={styles.trendingSection}>
                <div className={styles.trendingCard}>
                  <div className={styles.trendingHeader}>
                    <Flame size={14} strokeWidth={1.5} />
                    <h3>Trending Now</h3>
                  </div>
                  <p className={styles.trendingDescription}>
                    Today's popular products and searches.
                  </p>
                </div>
              </div>
            )}
            
            <FeedDisplay 
              showComposer={true}
              composerPlaceholder="Search products or ask AI for recommendations..."
              feedParams={{ entityType: 'product' }}
            />
          </section>
        </div>
      </main>
    </div>
  );
});

export default EnhancedProductsPage;