"use client";
import React from 'react';
import styles from './MarketplacePage.module.css';
// Leftside removed - using HorizontalFilters mobile modal instead
import { useIsMobile } from '../../../hooks/useMobileDetection';
// Removed context providers - using direct hooks
import FeedDisplay from '../../../components/Feed/FeedDisplay';

/**
 * Enhanced Marketplace Page with AI integration
 * This version integrates the UnifiedComposer and AI responses
 */
export default function EnhancedMarketplacePage({ 
  serverProducts = [], 
  serverFilters = {},
  initialHasMore = true 
}) {
  const isMobile = useIsMobile();

  // Configure initial params for marketplace (products only)
  const initialParams = {
    feedType: "latest",
    entityTypes: ["product"],
    contentType: "products",
    page: 1,
    pageSize: 20,
    ...serverFilters
  };

  return (
    <div className={styles.container}>
      <main className={styles.mainContent}>
        <div className={styles.layoutGrid}>
          {/* Leftside removed - using HorizontalFilters mobile modal instead */}
          <section className={styles.contentArea}>
            <FeedDisplay 
              showComposer={true}
              composerPlaceholder="Search products, ask about deals, or list an item..."
              initialItems={serverProducts}
              feedParams={{ entityType: 'product' }}
            />
          </section>
        </div>
      </main>
    </div>
  );
}