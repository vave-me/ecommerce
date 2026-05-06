"use client";
import React, { memo } from "react";
import { Clock, X } from "@/icons";
import styles from "./SearchBar.module.css";
/**
 * RecentSearches Component
 * Placeholder for recent searches functionality
 * This component is created for future enhancement as it wasn't present in the original SearchBar
 */
const RecentSearches = memo(({ 
    recentSearches = [], 
    onSearchClick, 
    onClearRecent,
    isVisible = false 
}) => {
    // Return null if no recent searches or not visible
    if (!isVisible || !recentSearches.length) {
        return null;
    }
    return (
        <div className={styles.recentSearches}>
            <div className={styles.recentHeader}>
                <Clock size={16} className={styles.recentIcon} />
                <span className={styles.recentTitle}>Recent Searches</span>
                <button
                    type="button"
                    onClick={onClearRecent}
                    aria-label="Clear recent searches"
                    className={styles.clearRecentButton}
                >
                    <X size={14} />
                </button>
            </div>
            <ul className={styles.recentList}>
                {recentSearches.map((search, index) => (
                    <li key={index} className={styles.recentItem}>
                        <button
                            type="button"
                            onClick={() => onSearchClick(search)}
                            className={styles.recentSearchButton}
                            aria-label={`Search for ${search}`}
                        >
                            <Clock size={14} className={styles.recentItemIcon} />
                            <span>{search}</span>
                        </button>
                    </li>
                ))}
            </ul>
        </div>
    );
});
RecentSearches.displayName = 'RecentSearches';
export default RecentSearches; 