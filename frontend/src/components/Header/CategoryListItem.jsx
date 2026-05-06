"use client";
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import styles from './CategoryListItem.module.css';
/**
 * CategoryListItem Component
 * Individual category item with selection, count, and subcategory indicator
 * Extracted from SelectTopic to improve component modularity
 */
const CategoryListItem = memo(({ 
    category, 
    isSelected, 
    hasSubcategories, 
    onClick, 
    onKeyDown 
}) => {
    if (!category || typeof category !== 'object') {
        return null;
    }
    return (
        <li
            key={category.id || category.value}
            className={`${styles.categoryRow} ${isSelected ? styles.categorySelected : ''} ${category.featured ? styles.categoryFeatured : ''}`}
            onClick={onClick}
            onKeyDown={onKeyDown}
            role="menuitem"
            tabIndex={0}
            aria-current={isSelected ? 'true' : undefined}
        >
            <span>{category.label}</span>
            {hasSubcategories && (
                <span className={styles.hasSubcategories}>❯</span>
            )}
            {typeof category.count === 'number' && (
                <span className={styles.categoryCount}>
                    {category.count.toLocaleString()}
                </span>
            )}
        </li>
    );
});
CategoryListItem.displayName = 'CategoryListItem';
CategoryListItem.propTypes = {
    category: PropTypes.shape({
        id: PropTypes.string,
        value: PropTypes.string.isRequired,
        label: PropTypes.string.isRequired,
        featured: PropTypes.bool,
        count: PropTypes.number,
    }).isRequired,
    isSelected: PropTypes.bool.isRequired,
    hasSubcategories: PropTypes.bool.isRequired,
    onClick: PropTypes.func.isRequired,
    onKeyDown: PropTypes.func.isRequired,
};
export default CategoryListItem; 