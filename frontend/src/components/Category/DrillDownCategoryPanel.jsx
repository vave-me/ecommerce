"use client";
import { ChevronLeft, ChevronRight } from '../../../icons';
import React, { useEffect, useState, memo } from "react";
import PropTypes from "prop-types";
import { useTranslations } from "next-intl"; //  Import hook
import styles from "./DrillDownCategoryPanel.module.css";
import { fetchSubCategories } from "../../api/categories"; // Assuming API returns appropriate language if possible
/**
 * A single-level drill-down panel that shows:
 * 1) A "Back" button + category name
 * 2) Subcategories in a list
 * 3) Loading/Error/Empty states
 * Uses next-intl for translations.
 * 
 * OPTIMIZED: Memoized for better category navigation performance
 */
const DrillDownCategoryPanel = memo(function DrillDownCategoryPanel({
                                                   category,
                                                   onSelectCategory,
                                                   onBack,
                                                   // level = 0, // level prop is passed but not used
                                                   parent = null,
                                               }) {
    const t = useTranslations('DrillDownCategoryPanel'); //  Instantiate hook with namespace
    const [subcategories, setSubcategories] = useState([]);
    const [hasFetched, setHasFetched] = useState(false);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null); // Store actual error object for logging
    // Assume category.name is already translated if needed
    const categoryName = category ? category.name : t('fallbackRootName'); //  Use translation fallback
    // On mount or category change => fetch subcategories
    useEffect(() => {
        // Reset state when category changes
        setSubcategories([]);
        setHasFetched(false);
        setError(null);
        if (!category || !category.id) {
            // If no category or no ID, don't attempt to fetch
            setHasFetched(true); // Mark as "fetched" to potentially show "no subcategories"
            return;
        }
        let isMounted = true; // Prevent state updates on unmounted component
        const fetchSubs = async () => {
            setLoading(true);
            try {
                const data = await fetchSubCategories(category.id);
                if (isMounted) {
                    setSubcategories(data || []);
                    setHasFetched(true);
                }
            } catch (err) {
                if (isMounted) {
                    setError(err); // Store the actual error object
                }
            } finally {
                if (isMounted) {
                    setLoading(false);
                }
            }
        };
        fetchSubs();
        return () => {
            isMounted = false; // Cleanup function for async calls
        };
    }, [category, setSubcategories]); // Complete dependency array
    return (
        <div className={styles.panelContainer}>
            {/* Header row: "Back" button + current category name */}
            <div className={styles.headerRow}>
                {parent ? (
                    <button
                        className={styles.backButton}
                        onClick={onBack}
                        aria-label={t('backButtonAria', {parentName: parent.name || ''})} //  Added ARIA label
                    >
                        <ChevronLeft aria-hidden="true"/> {/* Icon is decorative */}
                        {/*   Use translation */}
                        {t('backButton')}
                    </button>
                ) : (
                    // preserve spacing if no back button
                    <span style={{minWidth: '60px', display: 'inline-block'}}/> // Adjusted for typical button width
                )}
                {/* Display current category name (assumed translated) */}
                <div className={styles.currentCategoryTitle}>{categoryName}</div>
            </div>
            {/* Body => subcategory list or loading/error */}
            <div className={styles.contentArea}>
                {loading && (
                    //   Use translation
                    <div className={styles.loadingMsg}>{t('loading')}</div>
                )}
                {error && (
                    //   Use translation (show generic message, log specific error)
                    <div className={styles.errorMsg}>{t('errorLoading')}</div>
                )}
                {!loading && !error && subcategories.length > 0 && (
                    <div className={styles.subcatList}>
                        {subcategories.map((sub) => (
                            <div
                                key={sub.id}
                                className={styles.subcatItem}
                                onClick={() => onSelectCategory(sub)}
                                role="button"
                                tabIndex={0}
                                onKeyPress={(e) => {
                                    if (e.key === "Enter" || e.key === " ") {
                                        onSelectCategory(sub);
                                    }
                                }}
                                // Add ARIA label for better accessibility
                                aria-label={t('selectSubcategoryAria', {subcategoryName: sub.name || ''})}
                            >
                                {/* Assume sub.name is already translated */}
                                <span>{sub.name}</span>
                                <ChevronRight aria-hidden="true"/> {/* Icon is decorative */}
                            </div>
                        ))}
                    </div>
                )}
                {/* Only show "No further subcategories" if fetch completed successfully and result is empty */}
                {!loading && !error && hasFetched && subcategories.length === 0 && (
                    //   Use translation
                    <div className={styles.noSubcatMsg}>{t('noSubcategories')}</div>
                )}
            </div>
        </div>
    );
}, (prevProps, nextProps) => {
    // Only re-render if category or parent changed
    return (
        prevProps.category?.id === nextProps.category?.id &&
        prevProps.category?.name === nextProps.category?.name &&
        prevProps.parent?.id === nextProps.parent?.id &&
        prevProps.parent?.name === nextProps.parent?.name
        // Skip function comparisons as they're callbacks
    );
});
DrillDownCategoryPanel.propTypes = {
    category: PropTypes.shape({
        id: PropTypes.oneOfType([PropTypes.number, PropTypes.string]).isRequired, // ID is required to fetch
        name: PropTypes.string,
    }),
    onSelectCategory: PropTypes.func.isRequired,
    onBack: PropTypes.func,
    // level: PropTypes.number, // Prop not used
    parent: PropTypes.shape({ // Define parent shape for ARIA label
        name: PropTypes.string,
    }),
};
DrillDownCategoryPanel.defaultProps = {
    category: null, // Parent component should ideally always provide a category
    onBack: null,   // Explicitly null if no back action is possible
    parent: null,
};
export default DrillDownCategoryPanel;