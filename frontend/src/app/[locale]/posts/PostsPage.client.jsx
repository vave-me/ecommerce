// File: src/app/posts/PostsPage.client.jsx
"use client";
import React, {useState, useEffect, memo} from 'react';
import {useQuery} from '@tanstack/react-query';
import {useSelector} from 'react-redux';
import styles from './PostsPage.module.css'; // Assuming CSS module exists
// Import Components
// Leftside removed - using HorizontalFilters mobile modal instead
// Import API function & Hooks
import {searchPostsWithFilters} from '../../../api/searchApi'; // Changed from postsApi to searchApi
import { useIsMobile } from '../../../hooks/useMobileDetection';
import ImprovedPostCard from "../design/posts/page"; // Adjust path if needed
import ClientOnly from '../../../components/ClientOnly';
// --- Loading & Error Placeholders ---
// Simple placeholders, can be replaced with more sophisticated Skeleton Loaders
const LoadingPlaceholder = memo(function LoadingPlaceholder() {
    return <div className={styles.loadingState}>Loading posts...</div>;
});
const ErrorPlaceholder = memo(function ErrorPlaceholder({error, onRetry, labels}) {
    return (
        <div className={styles.errorState}>
            <h2>{labels?.errorTitle || "Something went wrong"}</h2>
            <p>{labels?.errorMessage || "We couldn't load posts right now. Please try again later."}</p>
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
});
const EmptyPlaceholder = memo(function EmptyPlaceholder({labels, currentCategory}) {
    const emptyMessage = currentCategory 
        ? (labels?.emptyMsg || `No posts found in the ${currentCategory} category.`)
        : (labels?.emptyMsg || "No posts found matching your criteria.");
    return (
        <div className={styles.emptyState}>
            <h3>{labels?.emptyTitle || "No Posts Found"}</h3>
            <p>{emptyMessage}</p>
        </div>
    );
});
/**
 * Client component for the Posts Listing Page.
 * Uses React Query (`useQuery`) for client-side data fetching,
 * similar to the working homepage pattern.
 * Relies on Redux state (`listingFilters`) for filtering criteria.
 */
const PostsPageClient = memo(function PostsPageClient(props) {
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
        queryKey: ['posts', listingFilters, currentCategory],
        // The function that performs the fetch
        queryFn: () => currentCategory 
            ? searchPostsWithFilters({...listingFilters, category: currentCategory})
            : searchPostsWithFilters(listingFilters),
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
        return <div className={styles.loadingState}>Loading posts...</div>;
    }
    // --- Render Logic ---
    // Now we're guaranteed to be on the client, we can use the full rendering
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
});
// --- Layout Components ---
const ListLayout = memo(function ListLayout({isMobile, posts}) {
    return (
        <div className={styles.layoutListGrid}>
            {/* Leftside removed - using HorizontalFilters mobile modal instead */}
            <section className={styles.contentListArea}>
                {posts.map((post) => (
                    <PostListItem key={post.id} post={post}/>
                ))}
            </section>
            {/* Optionally include Rightside if there's space */}
        </div>
    );
});
const GridLayout = memo(function GridLayout({isMobile, posts}) {
    return (
        <div className={styles.layoutGrid}>
            {/* Leftside removed - using HorizontalFilters mobile modal instead */}
            <section className={styles.contentArea}>
                <ul className={styles.postList}>
                    {posts.map((post) => (
                        <li key={post.id}>
                            <ImprovedPostCard post={post}/>
                        </li>
                    ))}
                </ul>
            </section>
        </div>
    );
});
const PostListItem = memo(function PostListItem({post}) {
    return (
        <div className={styles.postListItem}>
            <ImprovedPostCard post={post} displayMode="list" />
        </div>
    );
});
export default PostsPageClient;