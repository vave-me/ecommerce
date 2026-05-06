import React, { useCallback } from 'react';
import PropTypes from 'prop-types';
import CategorySelectionArea from "../../../../components/Category/CategorySelectionArea";

/**
 * Shared CategorySelector Component
 * Wraps CategorySelectionArea with consistent props and error handling
 * Used across all creation modals for category selection
 */
export function CategorySelector({
    categories = [],
    selectedCategoryId = '',
    selectedCategorySlug = '',
    selectedCategoryName = '',
    onCategorySelect,
    onCategoryClear,
    isLoading = false,
    error = null,
    styles,
    label = "Category",
    placeholder = "Select a category",
    required = true,
    categoryType = 'marketplace' // marketplace, jobs, services, etc.
}) {
    const handleCategorySelect = useCallback((category) => {
        if (category && onCategorySelect) {
            const categoryData = {
                id: category.id,
                slug: category.slug || '',
                name: category.name || ''
            };
            onCategorySelect(categoryData);
        }
    }, [onCategorySelect]);

    const handleCategoryClear = useCallback(() => {
        if (onCategoryClear) {
            onCategoryClear();
        }
    }, [onCategoryClear]);

    return (
        <div className={styles.formGroup}>
            <label className={styles.formLabel}>
                {label} {required && <span className={styles.requiredMark}>*</span>}
            </label>
            
            <CategorySelectionArea
                categories={categories}
                selectedCategoryId={selectedCategoryId}
                onSelectCategory={handleCategorySelect}
                onClear={handleCategoryClear}
                isLoading={isLoading}
                placeholder={placeholder}
                categoryType={categoryType}
                className={error ? styles.inputError : ''}
            />
            
            {error && (
                <div className={styles.fieldError} role="alert">
                    {error}
                </div>
            )}

            {/* Debugging section - only visible when categories aren't working */}
            {(!categories || categories.length === 0) && (
                <div style={{
                    marginTop: '10px',
                    padding: '8px',
                    background: '#fff4e5',
                    border: '1px solid #ffcc80',
                    borderRadius: '4px'
                }}>
                    <p style={{margin: '0 0 8px 0', fontWeight: 'bold'}}>Debug Info - No Categories Found</p>
                    <p style={{margin: '0 0 4px 0', fontSize: '14px'}}>Raw categories
                        data: {JSON.stringify(categories)}</p>
                    {categories && categories.length === 0 && (
                        <button
                            type="button"
                            onClick={() => {/* Debug click handler */}}
                            style={{
                                background: '#ff9800',
                                border: 'none',
                                padding: '4px 8px',
                                borderRadius: '4px',
                                color: 'white',
                                cursor: 'pointer',
                                fontSize: '12px'
                            }}
                        >
                            Log Categories to Console
                        </button>
                    )}
                </div>
            )}
            
            {selectedCategoryName && (
                <div className={styles.selectedCategoryDisplay}>
                    <span className={styles.selectedCategoryLabel}>Selected:</span>
                    <span className={styles.selectedCategoryName}>{selectedCategoryName}</span>
                </div>
            )}
        </div>
    );
}

CategorySelector.propTypes = {
    categories: PropTypes.array,
    selectedCategoryId: PropTypes.string,
    selectedCategorySlug: PropTypes.string,
    selectedCategoryName: PropTypes.string,
    onCategorySelect: PropTypes.func.isRequired,
    onCategoryClear: PropTypes.func.isRequired,
    isLoading: PropTypes.bool,
    error: PropTypes.string,
    styles: PropTypes.object.isRequired,
    label: PropTypes.string,
    placeholder: PropTypes.string,
    required: PropTypes.bool,
    categoryType: PropTypes.string
}; 