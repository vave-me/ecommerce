"use client";
import React, { useMemo, useState, useCallback, memo } from "react";
import PropTypes from "prop-types";
import {useDispatch, useSelector} from "react-redux";
import {useTranslations} from "next-intl"; //  Import hook
import {updateFilter} from "../../redux/slices/listingFiltersSlice";
import {useCategory} from "../../context/CategoryContext";
import styles from "./CategoryBar.module.css";
const MAX_CATS = 10;
/**
 * CategoryBar: horizontal row of category chips with optional "Show More" toggle.
 * - Reads categories from useCategory()
 * - Dispatches updateFilter to Redux on selection
 * - Uses next-intl for translations
 */
function CategoryBar({useTopPadding = false}) {
    const t = useTranslations('CategoryBar'); //  Instantiate hook with namespace
    const {categories: mainCategories, loading, error} = useCategory() || {};
    const dispatch = useDispatch();
    const {category} = useSelector((state) => state.listingFilters);
    const [expanded, setExpanded] = useState(false);
    const displayedCategories = useMemo(() => {
        if (!mainCategories) return [];
        return expanded ? mainCategories : mainCategories.slice(0, MAX_CATS);
    }, [mainCategories, expanded]);
    const handleSelectCategory = useCallback((cat) => {
        dispatch(updateFilter({key: "category", value: cat.id}));
    }, [dispatch]);
    const toggleShowAll = useCallback(() => setExpanded((prev) => !prev), []);
    const shouldShowToggle = useMemo(() => 
        mainCategories && mainCategories.length > MAX_CATS, 
        [mainCategories]
    );
    const barClasses = useMemo(() => 
        `${styles.barContainer} ${useTopPadding ? styles.topPadding : ""}`, 
        [useTopPadding]
    );
    // Render states with translations
    if (loading) {
        return (
            <div className={barClasses}>
                {t('loading')} {/*  Use translation */}
            </div>
        );
    }
    if (error) {
        // Log the actual error for debugging, show generic translated message
        return (
            <div className={barClasses}>
                {t('error')} {/*  Use translation */}
            </div>
        );
    }
    if (!mainCategories || mainCategories.length === 0) {
        return (
            <div className={barClasses}>
                {t('empty')} {/*  Use translation */}
            </div>
        );
    }
    // Main render with translations
    return (
        <div className={barClasses}>
            <div className={styles.categoryRow}>
                {displayedCategories.map((cat) => {
                    const isSelected = cat.id === category;
                    const chipClass = [
                        styles.categoryChip,
                        isSelected && styles.selectedChip,
                    ]
                        .filter(Boolean)
                        .join(" ");
                    return (
                        <div
                            key={cat.id}
                            onClick={() => handleSelectCategory(cat)}
                            className={chipClass}
                            role="button"
                            tabIndex={0}
                            onKeyPress={(e) => {
                                if (e.key === "Enter" || e.key === " ") {
                                    handleSelectCategory(cat);
                                }
                            }}
                            //   Use translation with interpolation
                            aria-label={t('ariaCategory', {categoryName: cat.name || 'Unnamed'})}
                        >
                            {cat.name || 'Unnamed'} {/* Display category name, provide fallback */}
                        </div>
                    );
                })}
                {shouldShowToggle && (
                    <button onClick={toggleShowAll} className={styles.toggleButton}>
                        {expanded
                            ? t('showLess') //  Use translation
                            //   Use translation with interpolation
                            : t('showMore', {count: mainCategories.length - MAX_CATS})}
                    </button>
                )}
            </div>
        </div>
    );
}
CategoryBar.propTypes = {
    useTopPadding: PropTypes.bool,
};
export default memo(CategoryBar);