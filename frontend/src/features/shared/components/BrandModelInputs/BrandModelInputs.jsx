import React from 'react';
import PropTypes from 'prop-types';
/**
 * Shared BrandModelInputs component for consistent brand and model input fields across all modals
 */
export function BrandModelInputs({
    brandId = "brand",
    modelId = "model",
    brandLabel = "Brand",
    modelLabel = "Model",
    brandValue,
    modelValue,
    onBrandChange,
    onModelChange,
    brandPlaceholder = "Brand name",
    modelPlaceholder = "Model number/name",
    brandError,
    modelError,
    className = "",
    styles,
    required = false,
    disabled = false,
    ...props
}) {
    const handleBrandChange = (e) => {
        if (onBrandChange) {
            onBrandChange(e.target.value);
        }
    };
    const handleModelChange = (e) => {
        if (onModelChange) {
            onModelChange(e.target.value);
        }
    };
    return (
        <div className={`${styles?.formRow || 'form-row'} ${className}`} {...props}>
            <div className={styles?.formGroup || 'form-group'}>
                <label 
                    className={styles?.formLabel || 'form-label'} 
                    htmlFor={brandId}
                >
                    {brandLabel}
                    {required && <span className={styles?.requiredMark || 'required-mark'}>*</span>}
                </label>
                <input
                    id={brandId}
                    className={`${styles?.formInput || 'form-input'} ${brandError ? styles?.inputError || 'input-error' : ''}`}
                    type="text"
                    value={brandValue || ''}
                    onChange={handleBrandChange}
                    placeholder={brandPlaceholder}
                    disabled={disabled}
                    aria-invalid={!!brandError}
                />
                {brandError && (
                    <div className={styles?.fieldError || 'field-error'}>{brandError}</div>
                )}
            </div>
            <div className={styles?.formGroup || 'form-group'}>
                <label 
                    className={styles?.formLabel || 'form-label'} 
                    htmlFor={modelId}
                >
                    {modelLabel}
                    {required && <span className={styles?.requiredMark || 'required-mark'}>*</span>}
                </label>
                <input
                    id={modelId}
                    className={`${styles?.formInput || 'form-input'} ${modelError ? styles?.inputError || 'input-error' : ''}`}
                    type="text"
                    value={modelValue || ''}
                    onChange={handleModelChange}
                    placeholder={modelPlaceholder}
                    disabled={disabled}
                    aria-invalid={!!modelError}
                />
                {modelError && (
                    <div className={styles?.fieldError || 'field-error'}>{modelError}</div>
                )}
            </div>
        </div>
    );
}
BrandModelInputs.propTypes = {
    brandId: PropTypes.string,
    modelId: PropTypes.string,
    brandLabel: PropTypes.string,
    modelLabel: PropTypes.string,
    brandValue: PropTypes.string,
    modelValue: PropTypes.string,
    onBrandChange: PropTypes.func.isRequired,
    onModelChange: PropTypes.func.isRequired,
    brandPlaceholder: PropTypes.string,
    modelPlaceholder: PropTypes.string,
    brandError: PropTypes.string,
    modelError: PropTypes.string,
    className: PropTypes.string,
    styles: PropTypes.object,
    required: PropTypes.bool,
    disabled: PropTypes.bool,
};
export default BrandModelInputs; 