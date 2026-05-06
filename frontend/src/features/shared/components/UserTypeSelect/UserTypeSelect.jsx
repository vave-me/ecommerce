import React from 'react';
import PropTypes from 'prop-types';
// Centralized seller types - this should be the single source of truth
export const SELLER_TYPES = ["private", "business", "dealer", "manufacturer"];
/**
 * Shared UserTypeSelect component for consistent seller type selection across all modals
 */
export function UserTypeSelect({
    id = "seller-type",
    label = "Seller Type",
    value,
    onChange,
    error,
    className = "",
    styles,
    required = false,
    disabled = false,
    ...props
}) {
    const handleChange = (e) => {
        if (onChange) {
            onChange(e.target.value);
        }
    };
    const formatLabel = (type) => {
        return type.charAt(0).toUpperCase() + type.slice(1);
    };
    return (
        <div className={styles?.formGroup || 'form-group'}>
            <label 
                className={styles?.formLabel || 'form-label'} 
                htmlFor={id}
            >
                {label}
                {required && <span className={styles?.requiredMark || 'required-mark'}>*</span>}
            </label>
            <select
                id={id}
                className={`${styles?.formSelect || 'form-select'} ${error ? styles?.inputError || 'input-error' : ''} ${className}`}
                value={value}
                onChange={handleChange}
                disabled={disabled}
                aria-invalid={!!error}
                {...props}
            >
                {SELLER_TYPES.map((type) => (
                    <option key={type} value={type}>
                        {formatLabel(type)}
                    </option>
                ))}
            </select>
            {error && (
                <div className={styles?.fieldError || 'field-error'}>{error}</div>
            )}
        </div>
    );
}
UserTypeSelect.propTypes = {
    id: PropTypes.string,
    label: PropTypes.string,
    value: PropTypes.string.isRequired,
    onChange: PropTypes.func.isRequired,
    error: PropTypes.string,
    className: PropTypes.string,
    styles: PropTypes.object,
    required: PropTypes.bool,
    disabled: PropTypes.bool,
};
export default UserTypeSelect;