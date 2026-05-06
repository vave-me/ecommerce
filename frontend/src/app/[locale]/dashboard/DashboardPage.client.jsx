// File: src/app/dashboard/DashboardPage.client.jsx
"use client";
import React, { useState } from 'react';
import {useSelector} from 'react-redux';
import styles from './DashboardPage.module.css';
// Import Components
// Leftside removed - using HorizontalFilters mobile modal instead
// Import Hooks
import { useIsMobile } from '../../../hooks/useMobileDetection';
import PostCard from "./[category]/[slug]/PostCard";
// --- Loading & Error Placeholders ---
const LoadingPlaceholder = () => <div className={styles.loadingState}>Loading dashboard posts...</div>;
const ErrorPlaceholder = ({error, fetchError, onRetry}) => {
    const [isRetrying, setIsRetrying] = useState(false);
    const handleRetry = async () => {
        if (onRetry) {
            setIsRetrying(true);
            try {
                await onRetry();
            } finally {
                setIsRetrying(false);
            }
        }
    };
    return (
        <div className={styles.errorState}>
            <h2>Something went wrong</h2>
            <p>{fetchError || "We couldn't load dashboard posts right now. Please try again later."}</p>
            {onRetry && (
                <button 
                    onClick={handleRetry} 
                    disabled={isRetrying}
                    className={styles.retryButton}
                >
                    {isRetrying ? 'Retrying...' : 'Try Again'}
                </button>
            )}
            {/* Show technical error details only in development */}
            {process.env.NODE_ENV === 'development' && error?.message && (
                <pre className={styles.devError}>Error: {error.message}</pre>
            )}
        </div>
    );
};
const EmptyPlaceholder = ({labels}) => (
    <div className={styles.emptyState}>
        <h3>No Dashboard Posts Found</h3>
        <p>{labels?.emptyMsg || "No dashboard posts found matching your criteria."}</p>
    </div>
);
/**
 * Client component for the Dashboard Listing Page.
 * Uses server-side data passed as props instead of client-side fetching
 * to avoid dual data fetching and improve performance.
 */
export default function DashboardPageClient({
                                            serverPosts = [],
                                            serverFilters = {},
                                            labels = {},
                                            fetchError = null,
                                            totalPages = 0,
                                            totalCount = 0,
                                            currentPage = 1,
                                            currentCategory = null,
                                            onRetry = null
                                        }) {
    // --- Hooks ---
    const isMobile = useIsMobile();
    // Get current filters from Redux store for display mode
    const listingFilters = useSelector((state) => state.listingFilters);
    // Extract display mode for layout switching
    const {displayMode = 'list'} = listingFilters || {};
    // --- Render Logic ---
    // 1. Handle Error State
    if (fetchError) {
        return <ErrorPlaceholder fetchError={fetchError} onRetry={onRetry} />;
    }
    // 2. Handle Empty or Invalid Data
    if (!Array.isArray(serverPosts) || serverPosts.length === 0) {
        return <EmptyPlaceholder labels={labels}/>;
    }
    // Filter out any potentially invalid post entries
    const validPosts = serverPosts.filter(post => post && post.id);
    if (validPosts.length === 0) {
        return <EmptyPlaceholder labels={labels}/>;
    }
    // 3. Render Post List/Grid
    return (
        <div className={styles.container}>
            <main className={styles.mainContent}>
                {/* Optional title section */}
                {labels?.pageTitle && (
                    <div className={styles.pageHeader}>
                        <h1 className={styles.pageTitle}>{labels.pageTitle}</h1>
                        {totalCount > 0 && (
                            <p className={styles.resultCount}>
                                Showing {validPosts.length} of {totalCount} dashboard posts
                                {currentCategory && ` in ${currentCategory}`}
                            </p>
                        )}
                    </div>
                )}
                <GridLayout
                    isMobile={isMobile}
                    posts={validPosts}
                    displayMode={displayMode}
                />
                {/* Pagination Info */}
                {totalPages > 1 && (
                    <div className={styles.paginationInfo}>
                        <p>Page {currentPage} of {totalPages}</p>
                    </div>
                )}
            </main>
        </div>
    );
}
function GridLayout({isMobile, posts, displayMode}) {
    return (
        <div className={displayMode === 'list' ? styles.layoutListGrid : styles.layoutGrid}>
            {/* Leftside removed - using HorizontalFilters mobile modal instead */}
            <section className={displayMode === 'list' ? styles.contentListArea : styles.contentArea}>
                <ul className={styles.postList}>
                    {posts.map((post) => (
                        <li key={post.id}>
                            <PostCard post={post}/>
                        </li>
                    ))}
                </ul>
            </section>
        </div>
    );
}