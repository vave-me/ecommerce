"use client";
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import {useTranslations} from 'next-intl';
import CategoryListItem from './CategoryListItem';
import styles from './CategoryList.module.css';
/**
 * CategoryList Component
 * Renders list of categories with loading, error, and empty states
 * Extracted from SelectTopic to improve component modularity
 */
const CategoryList = memo(({
                          displayCategories,
                          isDisplayLoading,
                          error,
                          selectedParentId,
                          selectedCategoryValue,
                          onSelectCategory,
                          onShowSubcategories,
                          onBackToMain
                      }) => {
    const t = useTranslations('Topics');
    const categoriesTitle = t('categories');
    // Loading state
    if (isDisplayLoading) {
        return (
            <div className={styles.loadingIndicator}>
                {t('loadingCategories')}...
            </div>
        );
    }
    // Error state
    if (error) {
        return (
            <div className={styles.errorIndicator}>
                {t('errorLoadingCategories')}
            </div>
        );
    }
    // Empty state
    if (!displayCategories || displayCategories.length === 0) {
        return (
            <div className={styles.noDataIndicator}>
                {t('noCategoriesFound')}
            </div>
        );
    }
    // Invalid data state
    if (!Array.isArray(displayCategories)) {
        return (
            <div className={styles.errorIndicator}>
                {t('errorLoadingCategories')}
            </div>
        );
    }
    return (
        <ul className={styles.categoryList} role="menu" aria-label={categoriesTitle}>
            {/* "All" option only shows in main category view */}
            {!selectedParentId && (
                <li
                    key="all"
                    className={`${styles.categoryRow} ${selectedCategoryValue === "all" ? styles.categorySelected : ''}`}
                    onClick={() => onSelectCategory("all")}
                    onKeyDown={(e) => (e.key === 'Enter' || e.key === ' ') && onSelectCategory("all")}
                    role="menuitem"
                    tabIndex={0}
                    aria-current={selectedCategoryValue === "all" ? 'true' : undefined}
                >
                    <span>{t('category_all_label')}</span>
                </li>
            )}
            {/* Back button when showing subcategories */}
            {selectedParentId && (
                <li
                    key="back"
                    className={`${styles.categoryRow} ${styles.backButton}`}
                    onClick={onBackToMain}
                    onKeyDown={(e) => (e.key === 'Enter' || e.key === ' ') && onBackToMain()}
                    role="menuitem"
                    tabIndex={0}
                >
                    <span>← {t('backToCategories')}</span>
                </li>
            )}
            {/* Map over categories */}
            {displayCategories.map((category, index) => {
                if (!category || typeof category !== 'object') {
                    return null;
                }
                const isSelected = selectedCategoryValue === category.value;
                // Always assume categories have subcategories unless we're already showing subcategories
                const hasSubcategories = selectedParentId ? false : true;
                const handleClick = () => {
                    if (hasSubcategories) {
                        onShowSubcategories(category.id, category.value);
                    } else {
                        onSelectCategory(category.value);
                    }
                };
                const handleKeyDown = (e) => {
                    if (e.key === 'Enter' || e.key === ' ') {
                        if (hasSubcategories) {
                            onShowSubcategories(category.id, category.value);
                        } else {
                            onSelectCategory(category.value);
                        }
                    }
                };
                return (
                    <CategoryListItem
                        key={category.id || category.value || index}
                        category={category}
                        isSelected={isSelected}
                        hasSubcategories={hasSubcategories}
                        onClick={handleClick}
                        onKeyDown={handleKeyDown}
                    />
                );
            })}
        </ul>
    );
});
CategoryList.displayName = 'CategoryList';
CategoryList.propTypes = {
    displayCategories: PropTypes.array,
    isDisplayLoading: PropTypes.bool.isRequired,
    error: PropTypes.object,
    selectedParentId: PropTypes.string,
    selectedCategoryValue: PropTypes.string.isRequired,
    onSelectCategory: PropTypes.func.isRequired,
    onShowSubcategories: PropTypes.func.isRequired,
    onBackToMain: PropTypes.func.isRequired,
};
export default CategoryList; 