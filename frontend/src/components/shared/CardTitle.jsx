"use client";
import React, { memo } from 'react';
import Link from 'next/link';
import styles from './CardTitle.module.css';
/**
 * CardTitle - Reusable component for card titles/names with consistent styling
 * Memoized for performance optimization
 * 
 * @param {string} title - The title text to display
 * @param {string} href - URL for the title link (optional)
 * @param {string} className - Additional CSS classes for container
 * @param {string} titleClassName - Additional CSS classes for title element
 * @param {string} linkClassName - Additional CSS classes for link element
 * @param {string} headingLevel - HTML heading level (default: 'h2')
 * @param {function} onClick - Optional click handler
 * @param {boolean} truncate - Whether to apply 2-line truncation (default: true)
 */
const CardTitle = memo(({ 
    title,
    href,
    className = "",
    titleClassName = "",
    linkClassName = "",
    headingLevel = 'h2',
    onClick,
    truncate = true
}) => {
    if (!title || title.trim() === '') {
        return null;
    }
    const HeadingTag = headingLevel;
    const titleElement = (
        <HeadingTag
            className={`${styles.title} ${titleClassName} ${truncate ? styles.truncated : ''}`}
            onClick={onClick}
            title={title} // Full title on hover for truncated content
        >
            {title}
        </HeadingTag>
    );
    return (
        <div className={`${styles.titleContainer} ${className}`}>
            {href ? (
                <Link href={href} className={`${styles.titleLink} ${linkClassName}`}>
                    {titleElement}
                </Link>
            ) : (
                titleElement
            )}
        </div>
    );
});
export default CardTitle; 