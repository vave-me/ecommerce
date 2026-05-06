// File: src/app/news/PostsPage.client.jsx
"use client";
import React, {useState, useEffect} from 'react';
import {useQuery} from '@tanstack/react-query';
import {useSelector} from 'react-redux';
import styles from './PostsPage.module.css'; // Assuming CSS module exists
// Import Components
// Leftside removed - using HorizontalFilters mobile modal instead
// Import API function & Hooks
import {searchPostsWithFilters} from '../../../api/searchApi'; // Changed from postsApi to searchApi
import { useIsMobile } from '../../../hooks/useMobileDetection';
import ImprovedPostCard from "../design/posts/page"; // Fixed path to posts design component
import ClientOnly from '../../../components/ClientOnly';
// --- Loading & Error Placeholders ---
// Simple placeholders, can be replaced with more sophisticated Skeleton Loaders
const LoadingPlaceholder = () => <div className={styles.loadingState}>Loading news...</div>;
const ErrorPlaceholder = ({error, onRetry, labels}) => (
    <div className={styles.errorState}>
        <h2>{labels?.errorTitle || "Something went wrong"}</h2>
        <p>{labels?.errorMessage || "We couldn't load news right now. Please try again later."}</p>
        {onRetry && (
            <button 
                onClick={onRetry} 
                className={styles.retryButton}
                type="button"
            >
                {labels?.retryButton || "Try Again"}
            </button>
        )}
        {/* Show technical error details only in development */}
        {process.env.NODE_ENV === 'development' && error?.message && (
            <pre className={styles.devError}>Error: {error.message}</pre>
        )}
    </div>
);
const EmptyPlaceholder = ({labels, currentCategory}) => {
    const emptyMessage = currentCategory 
        ? (labels?.emptyMsg || `No news found in the ${currentCategory} category.`)
        : (labels?.emptyMsg || "No news found matching your criteria.");
    return (
        <div className={styles.emptyState}>
            <h3>{labels?.emptyTitle || "No News Found"}</h3>
            <p>{emptyMessage}</p>
        </div>
    );
};
/**
 * Client component for the News Listing Page.
 * Uses React Query (`useQuery`) for client-side data fetching,
 * similar to the working homepage pattern.
 * Relies on Redux state (`listingFilters`) for filtering criteria.
 */
export default function PostsPageClient(props) {
    const {
        serverPosts, 
        serverFilters, 
        fetchError, 
        labels = {}, 
        totalPages = 0, 
        totalCount = 0, 
        currentPage = 1,
        currentCategory
    } = props;
    // State to track if component is mounted (client-side)
    const [mounted, setMounted] = useState(false);
    // Set mounted state on client-side only
    useEffect(() => {
        setMounted(true);
    }, []);
    // --- Hooks ---
    const isMobile = useIsMobile();
    // Get current filters from Redux store
    const listingFilters = useSelector((state) => state.listingFilters);
    // Extract display mode for layout switching
    const {displayMode = 'list'} = listingFilters || {}; // Default to 'list'
    // --- Data Fetching (Client-Side) ---
    const {
        data: featuredData, // Data returned from the API { posts: [...] }
        isLoading,
        isError,
        error,          // Error object if the fetch fails
        isFetching,     // Useful for showing loading indicators on refetch
        refetch,        // Function to manually refetch data
    } = useQuery({
        // Query key includes filters: TanStack Query automatically refetches
        // when the `listingFilters` object changes identity.
        queryKey: ['news', listingFilters, currentCategory],
        // The function that performs the fetch
        queryFn: () => {
            const filters = currentCategory 
                ? {...listingFilters, category: currentCategory, type: 'news'}
                : {...listingFilters, type: 'news'};
            return searchPostsWithFilters(filters);
        },
        // Optional: Keep previous data visible while refetching for smoother UX
        // keepPreviousData: true,
        // Optional: Set staletime if data doesn't change frequently
        // staleTime: 5 * 60 * 1000, // 5 minutes
        // Initialize with server data if available
        initialData: serverPosts ? {posts: serverPosts, totalPages, totalCount, currentPage} : undefined,
    });
    // Server-side error handling
    if (fetchError) {
        return (
            <div className={styles.container}>
                <main className={styles.mainContent}>
                    <ErrorPlaceholder 
                        error={{message: fetchError}} 
                        onRetry={() => window.location.reload()} 
                        labels={labels}
                    />
                </main>
            </div>
        );
    }
    // On server or during first client render, show a minimal loading state
    if (!mounted) {
        return <div className={styles.loadingState}>Loading news...</div>;
    }
    // --- Render Logic ---
    // 1. Handle Loading State
    if (isLoading) {
        return (
            <div className={styles.container}>
                <main className={styles.mainContent}>
                    <LoadingPlaceholder/>
                </main>
            </div>
        );
    }
    // 2. Handle Error State
    if (isError) {
        return (
            <div className={styles.container}>
                <main className={styles.mainContent}>
                    <ErrorPlaceholder 
                        error={error} 
                        onRetry={refetch} 
                        labels={labels}
                    />
                </main>
            </div>
        );
    }
    // 3. Handle Empty or Invalid Data
    const posts = featuredData?.posts;
    if (!Array.isArray(posts) || posts.length === 0) {
        return (
            <div className={styles.container}>
                <main className={styles.mainContent}>
                    <EmptyPlaceholder labels={labels} currentCategory={currentCategory} />
                </main>
            </div>
        );
    }
    // Filter out any potentially invalid post entries just in case
    const validPosts = posts.filter(post => post && post.id);
    if (validPosts.length === 0) {
        return (
            <div className={styles.container}>
                <main className={styles.mainContent}>
                    <EmptyPlaceholder labels={labels} currentCategory={currentCategory} />
                </main>
            </div>
        );
    }
    // 4. Render Post List/Grid
    return (
        <div className={styles.container}>
            {/* Optional: Indicate background refetching */}
            {isFetching && !isLoading && <div className={styles.refetchingIndicator}>Updating...</div>}
            <main className={styles.mainContent}>
                {displayMode === 'list' ? (
                    <ListLayout
                        isMobile={isMobile}
                        posts={validPosts}
                    />
                ) : (
                    <GridLayout
                        isMobile={isMobile}
                        posts={validPosts}
                    />
                )}
            </main>
        </div>
    );
}
// --- Layout Components --- (Can be moved to separate files if preferred)
function ListLayout({isMobile, posts}) {
    return (
        <div className={styles.layoutListGrid}>
            {/* Leftside removed - using HorizontalFilters mobile modal instead */}
            <section className={styles.contentListArea}>
                {posts.map((post) => (
                    <div key={post.id} className={styles.postListItem}>
                        <ImprovedPostCard post={post}/>
                    </div>
                ))}
            </section>
            {/* Optionally include Rightside for list view if design requires */}
        </div>
    );
}
function GridLayout({isMobile, posts}) {
    return (
        <div className={styles.layoutGrid}>
            {/* Leftside removed - using HorizontalFilters mobile modal instead */}
            <section className={styles.contentArea}>
                <ul className={styles.postList}>
                    {posts.map((post) => (
                        // Ensure ImprovedPostCard handles potential missing fields gracefully
                        <li key={post.id}>
                            <ImprovedPostCard post={post}/>
                        </li>
                    ))}
                </ul>
            </section>
        </div>
    );
}