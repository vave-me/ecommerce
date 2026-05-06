"use client";
import React, { memo } from 'react';
import styles from './HeaderIcon.module.css';

/**
 * Unified HeaderIcon Component
 * Provides consistent styling for all header icons
 */
const HeaderIcon = memo(({ 
    icon: Icon, 
    onClick, 
    badge = 0, 
    ariaLabel, 
    title,
    className = '',
    isMobile = false,
    ariaExpanded,
    ariaHaspopup 
}) => {
    const handleClick = (e) => {
        if (onClick) {
            onClick(e);
        }
    };
    
    const formatBadge = (count) => {
        if (count === 0) return null;
        if (count > 99) return '99+';
        return count.toString();
    };
    
    const displayBadge = formatBadge(badge);
    
    return (
        <button
            type="button"
            onClick={handleClick}
            className={`${styles.iconButton} ${className}`}
            aria-label={ariaLabel}
            title={title}
            aria-expanded={ariaExpanded}
            aria-haspopup={ariaHaspopup}
        >
            <Icon 
                size={24} 
                className={styles.icon}
                aria-hidden="true"
            />
            {displayBadge && (
                <span className={styles.badge} aria-hidden="true">
                    {displayBadge}
                </span>
            )}
        </button>
    );
});

HeaderIcon.displayName = 'HeaderIcon';

export default HeaderIcon;