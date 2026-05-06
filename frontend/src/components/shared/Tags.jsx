"use client";
import React, { memo } from 'react';
import styles from './Tags.module.css';
/**
 * Tags - Shared tags component for all card types
 * Memoized for performance optimization
 * 
 * Displays a list of tags with optional truncation and overflow indicator.
 * Used consistently across all card components (deals, vehicles, products, etc.)
 * 
 * @param {Array} tags - Array of tag strings
 * @param {number} maxTags - Maximum number of tags to display before showing overflow (default: 5)
 * @param {string} ariaLabel - Aria label for accessibility (default: "Tags")
 * @param {string} className - Additional CSS classes
 * @param {function} onTagClick - Optional callback when a tag is clicked (receives tag and index)
 */
const Tags = memo(({
    tags = [],
    maxTags = 5,
    ariaLabel = "Tags",
    className = '',
    onTagClick
}) => {
    // Early return if no tags
    if (!tags || tags.length === 0) {
        return null;
    }
    const displayTags = tags.slice(0, maxTags);
    const remainingCount = tags.length - maxTags;
    const hasOverflow = remainingCount > 0;
    const handleTagClick = (tag, index) => {
        if (onTagClick && typeof onTagClick === 'function') {
            onTagClick(tag, index);
        }
    };
    return (
        <div className={`${styles.tagsContainer} ${className}`} aria-label={ariaLabel}>
            {displayTags.map((tag, index) => (
                <span 
                    key={index} 
                    className={styles.tagBadge}
                    onClick={() => handleTagClick(tag, index)}
                    role={onTagClick ? "button" : undefined}
                    tabIndex={onTagClick ? 0 : undefined}
                    onKeyDown={onTagClick ? (e) => {
                        if (e.key === 'Enter' || e.key === ' ') {
                            e.preventDefault();
                            handleTagClick(tag, index);
                        }
                    } : undefined}
                >
                    {tag}
                </span>
            ))}
            {hasOverflow && (
                <span 
                    className={styles.moreTagsBadge} 
                    aria-label={`${remainingCount} more tags`}
                    title={`${remainingCount} more tags: ${tags.slice(maxTags).join(', ')}`}
                >
                    +{remainingCount}
                </span>
            )}
        </div>
    );
});
export default Tags; 