import React from 'react';
import PropTypes from 'prop-types';
/**
 * Shared PriceInput component for consistent price input across all modals
 * Includes currency symbol, proper number formatting, and error handling
 */
export function PriceInput({
    id,
    label,
    value,
    onChange,
    error,
    className = "",
    styles,
    required = false,
    disabled = false,
    placeholder = "0.00",
    min = "0",
    step = "0.01",
    currency = "€",
    ...props
}) {
    const handleChange = (e) => {
        if (onChange) {
            onChange(e.target.value);
        }
    };
    return (
        <div className={`${styles?.formGroup || 'form-group'} ${className}`} {...props}>
            <label 
                className={styles?.formLabel || 'form-label'} 
                htmlFor={id}
            >
                {label}
                {required && <span className={styles?.requiredMark || 'required-mark'}> *</span>}
            </label>
            <div
                className={`${styles?.priceInputWrapper || 'price-input-wrapper'} ${error ? styles?.inputError || 'input-error' : ''}`}
            >
                <span className={styles?.currencySymbol || 'currency-symbol'}>{currency}</span>
                <input
                    id={id}
                    className={styles?.priceInput || 'price-input'}
                    type="number"
                    min={min}
                    step={step}
                    placeholder={placeholder}
                    value={value || ''}
                    onChange={handleChange}
                    disabled={disabled}
                    aria-invalid={!!error}
                    aria-describedby={error ? `${id}-error` : undefined}
                />
            </div>
            {error && (
                <div 
                    className={styles?.fieldError || 'field-error'} 
                    id={`${id}-error`} 
                    role="alert"
                >
                    {error}
                </div>
            )}
        </div>
    );
}
PriceInput.propTypes = {
    id: PropTypes.string.isRequired,
    label: PropTypes.string.isRequired,
    value: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
    onChange: PropTypes.func.isRequired,
    error: PropTypes.string,
    className: PropTypes.string,
    styles: PropTypes.object,
    required: PropTypes.bool,
    disabled: PropTypes.bool,
    placeholder: PropTypes.string,
    min: PropTypes.string,
    step: PropTypes.string,
    currency: PropTypes.string,
};
export default PriceInput; 