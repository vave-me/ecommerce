import React from 'react';
import PropTypes from 'prop-types';
import { Tag } from "@/icons";
/**
 * Shared TagsSkuInputs component for consistent tags and SKU input fields across all modals
 */
export function TagsSkuInputs({
    tagsId = "tags",
    skuId = "sku",
    tagsLabel = "Tags (comma separated)",
    skuLabel = "SKU",
    tagsValue,
    skuValue,
    onTagsChange,
    onSkuChange,
    tagsPlaceholder = "summer, discount, limited",
    skuPlaceholder = "Stock keeping unit",
    tagsError,
    skuError,
    className = "",
    styles,
    required = false,
    disabled = false,
    showTagIcon = true,
    ...props
}) {
    const handleTagsChange = (e) => {
        if (onTagsChange) {
            onTagsChange(e.target.value);
        }
    };
    const handleSkuChange = (e) => {
        if (onSkuChange) {
            onSkuChange(e.target.value);
        }
    };
    return (
        <div className={`${styles?.formRow || 'form-row'} ${className}`} {...props}>
            <div className={styles?.formGroup || 'form-group'}>
                <label 
                    className={styles?.formLabel || 'form-label'} 
                    htmlFor={tagsId}
                >
                    {showTagIcon && <Tag size={16} className={styles?.labelIcon || 'label-icon'} />}
                    {tagsLabel}
                    {required && <span className={styles?.requiredMark || 'required-mark'}>*</span>}
                </label>
                <input
                    id={tagsId}
                    className={`${styles?.formInput || 'form-input'} ${tagsError ? styles?.inputError || 'input-error' : ''}`}
                    type="text"
                    value={tagsValue || ''}
                    onChange={handleTagsChange}
                    placeholder={tagsPlaceholder}
                    disabled={disabled}
                    aria-invalid={!!tagsError}
                />
                {tagsError && (
                    <div className={styles?.fieldError || 'field-error'}>{tagsError}</div>
                )}
            </div>
            <div className={styles?.formGroup || 'form-group'}>
                <label 
                    className={styles?.formLabel || 'form-label'} 
                    htmlFor={skuId}
                >
                    {skuLabel}
                    {required && <span className={styles?.requiredMark || 'required-mark'}>*</span>}
                </label>
                <input
                    id={skuId}
                    className={`${styles?.formInput || 'form-input'} ${skuError ? styles?.inputError || 'input-error' : ''}`}
                    type="text"
                    value={skuValue || ''}
                    onChange={handleSkuChange}
                    placeholder={skuPlaceholder}
                    disabled={disabled}
                    aria-invalid={!!skuError}
                />
                {skuError && (
                    <div className={styles?.fieldError || 'field-error'}>{skuError}</div>
                )}
            </div>
        </div>
    );
}
TagsSkuInputs.propTypes = {
    tagsId: PropTypes.string,
    skuId: PropTypes.string,
    tagsLabel: PropTypes.string,
    skuLabel: PropTypes.string,
    tagsValue: PropTypes.string,
    skuValue: PropTypes.string,
    onTagsChange: PropTypes.func.isRequired,
    onSkuChange: PropTypes.func.isRequired,
    tagsPlaceholder: PropTypes.string,
    skuPlaceholder: PropTypes.string,
    tagsError: PropTypes.string,
    skuError: PropTypes.string,
    className: PropTypes.string,
    styles: PropTypes.object,
    required: PropTypes.bool,
    disabled: PropTypes.bool,
    showTagIcon: PropTypes.bool,
};
export default TagsSkuInputs; 