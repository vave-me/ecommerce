'use client';

import React from 'react';
import { ShoppingBag, Filter, RefreshCw } from '@/icons';
import styles from './FeedEmptyState.module.css';

const FeedEmptyState = ({ hasFilters = false, onClearFilters, onRefresh }) => {
  return (
    <div className={styles.emptyContainer}>
      <div className={styles.emptyContent}>
        <div className={styles.iconWrapper}>
          <ShoppingBag size={64} className={styles.emptyIcon} />
        </div>
        
        <h3 className={styles.emptyTitle}>
          {hasFilters ? "No Results Found" : "Your Feed is Empty"}
        </h3>
        
        <p className={styles.emptyMessage}>
          {hasFilters 
            ? "We couldn't find any items matching your current filters."
            : "There are no items to display at the moment."}
        </p>
        
        <div className={styles.actionButtons}>
          {hasFilters && onClearFilters && (
            <button 
              onClick={onClearFilters} 
              className={styles.clearButton}
              aria-label="Clear all filters"
            >
              <Filter size={16} />
              <span>Clear Filters</span>
            </button>
          )}
          
          {onRefresh && (
            <button 
              onClick={onRefresh} 
              className={styles.refreshButton}
              aria-label="Refresh feed"
            >
              <RefreshCw size={16} />
              <span>Refresh</span>
            </button>
          )}
        </div>
        
        <div className={styles.suggestions}>
          <p>You might want to:</p>
          <ul>
            {hasFilters && <li>Try adjusting your filters</li>}
            <li>Check back later for new items</li>
            <li>Browse popular categories</li>
            <li>Search for something specific</li>
          </ul>
        </div>
      </div>
    </div>
  );
};

export default FeedEmptyState;