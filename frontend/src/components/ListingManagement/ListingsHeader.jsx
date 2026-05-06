// src/components/ListingManagement/ListingsHeader.jsx
import React, { memo } from 'react';
import { ClipboardList } from '@/icons';
import styles from './ListingsHeader.module.css';
/**
 * ListingsHeader - Atomic Design Component
 * Header component for listing management page using CSS modules
 * 
 * @param {Object} props - Component props
 * @param {string} props.title - Header title text
 * @param {string} props.subtitle - Optional subtitle text
 * @returns {JSX.Element} Rendered header component
 */
const ListingsHeader = memo(({ 
    title = 'Listing Management', 
    subtitle = 'Manage your listings and inventory'
}) => {
    return (
        <header className={styles.headerContainer}>
            <div className={styles.headerContent}>
                <div className={styles.iconWrapper}>
                    <ClipboardList className={styles.icon} size={32} />
                </div>
                <div className={styles.textContent}>
                    <h1 className={styles.title}>{title}</h1>
                    {subtitle && (
                        <p className={styles.subtitle}>{subtitle}</p>
                    )}
                </div>
            </div>
        </header>
    );
});
ListingsHeader.displayName = 'ListingsHeader';
export default ListingsHeader;
