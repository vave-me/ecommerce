'use client';

import React from 'react';
import { RefreshCw, AlertCircle, WifiOff } from '@/icons';
import styles from './FeedErrorState.module.css';

const FeedErrorState = ({ error, onRetry }) => {
  // Determine error type based on error message or code
  const isNetworkError = error?.message?.toLowerCase().includes('network') || 
                        error?.message?.toLowerCase().includes('fetch') ||
                        error?.code === 'NETWORK_ERROR';
  
  const isServerError = error?.response?.status >= 500 || 
                       error?.message?.toLowerCase().includes('server');

  return (
    <div className={styles.errorContainer}>
      <div className={styles.errorContent}>
        <div className={styles.iconWrapper}>
          {isNetworkError ? (
            <WifiOff size={48} className={styles.errorIcon} />
          ) : (
            <AlertCircle size={48} className={styles.errorIcon} />
          )}
        </div>
        
        <h3 className={styles.errorTitle}>
          {isNetworkError 
            ? "No Internet Connection" 
            : isServerError 
              ? "Server Error" 
              : "Unable to Load Feed"}
        </h3>
        
        <p className={styles.errorMessage}>
          {isNetworkError 
            ? "Please check your internet connection and try again."
            : isServerError
              ? "Our servers are experiencing issues. Please try again in a few moments."
              : "We couldn't load the feed. This might be temporary."}
        </p>
        
        {onRetry && (
          <button 
            onClick={onRetry} 
            className={styles.retryButton}
            aria-label="Retry loading feed"
          >
            <RefreshCw size={16} />
            <span>Try Again</span>
          </button>
        )}
        
        <div className={styles.helpText}>
          <p>If the problem persists, you can:</p>
          <ul>
            <li>Refresh the page</li>
            <li>Clear your browser cache</li>
            <li>Try again later</li>
          </ul>
        </div>
      </div>
    </div>
  );
};

export default FeedErrorState;