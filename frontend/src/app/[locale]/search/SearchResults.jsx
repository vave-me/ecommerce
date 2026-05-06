"use client";
import React, { useEffect, useState, useCallback } from "react";
import { useSearchParams } from "next/navigation";
import { useDispatch } from "react-redux";
import { setFilters } from "../../../redux/slices/listingFiltersSlice";
import { useFeed } from "../../../hooks/useFeed";
import Feed from "../../../components/Feed/Feed";
import HorizontalFilters from "../../../components/Filters/HorizontalFilters";
import styles from "./SearchResults.module.css";

const SearchResultsContent = () => {
    const searchParams = useSearchParams();
    const dispatch = useDispatch();
    const { items: feedItems, isLoading, hasMore, loadMore, error } = useFeed();
    
    // Get search query from URL params
    const searchQuery = searchParams.get("query") || searchParams.get("q") || "";
    const category = searchParams.get("category") || "";
    const lat = searchParams.get("lat");
    const lng = searchParams.get("lng") || searchParams.get("lon");
    const radius = searchParams.get("radius");
    
    // Update Redux filters when search params change
    useEffect(() => {
        const filters = {
            searchText: searchQuery,
            ...(category && { category }),
            ...(lat && { lat: parseFloat(lat) }),
            ...(lng && { lng: parseFloat(lng) }),
            ...(radius && { radius: parseInt(radius, 10) }),
        };
        
        // Update Redux filters only
        dispatch(setFilters(filters));
    }, [searchQuery, category, lat, lng, radius, dispatch]);
    
    return (
        <div className={styles.container}>
            {/* Add HorizontalFilters for additional filtering */}
            <HorizontalFilters />
            
            <div className={styles.header}>
                <h1 className={styles.title}>
                    {searchQuery ? `Search Results for "${searchQuery}"` : "Search Results"}
                </h1>
                <p className={styles.resultCount}>
                    {feedItems.length > 0 && `Found ${feedItems.length} items`}
                </p>
            </div>
            
            {error && (
                <div className={styles.errorMessage}>
                    Failed to load search results. Please try again.
                </div>
            )}
            
            {!isLoading && feedItems.length === 0 && !error && (
                <div className={styles.noResultsMessage}>
                    <p>No results found for your search.</p>
                    <p>Try adjusting your filters or search terms.</p>
                </div>
            )}
            
            {/* Use Feed component to display results */}
            <div className={styles.resultsGrid}>
                <Feed />
            </div>
        </div>
    );
};

const SearchResults = () => {
    return <SearchResultsContent />;
};

SearchResults.displayName = 'SearchResults';
export default SearchResults;