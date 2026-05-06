import React from 'react';
import PropTypes from 'prop-types';
// Default condition options - can be overridden via props
const DEFAULT_CONDITIONS = ["new", "like-new", "excellent", "good", "fair"];
// Alternative condition sets for different item types
export const CONDITION_SETS = {
    GENERAL: ["new", "like-new", "excellent", "good", "fair"],
    PRODUCT: ["new", "used", "broken", "refurbished"],
    VEHICLE: ["new", "like-new", "excellent", "good", "fair"],
    ELECTRONICS: ["new", "used", "broken", "refurbished"]
};
/**
 * Shared ConditionSelect component for consistent condition selection across all modals
 */
export function ConditionSelect({
    id = "condition",
    label = "Condition",
    value,
    onChange,
    conditions = DEFAULT_CONDITIONS,
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
    return (
        <div className={`${styles?.formGroup || 'form-group'} ${className}`} {...props}>
            <label 
                className={styles?.formLabel || 'form-label'} 
                htmlFor={id}
            >
                {label}
                {required && <span className={styles?.requiredMark || 'required-mark'}> *</span>}
            </label>
            <select
                id={id}
                className={`${styles?.formSelect || 'form-select'} ${error ? styles?.inputError || 'input-error' : ''}`}
                value={value || ''}
                onChange={handleChange}
                disabled={disabled}
                aria-invalid={!!error}
                aria-describedby={error ? `${id}-error` : undefined}
            >
                {conditions.map((condition) => (
                    <option key={condition} value={condition}>
                        {condition.charAt(0).toUpperCase() + condition.slice(1)}
                    </option>
                ))}
            </select>
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
ConditionSelect.propTypes = {
    id: PropTypes.string,
    label: PropTypes.string,
    value: PropTypes.string,
    onChange: PropTypes.func.isRequired,
    conditions: PropTypes.arrayOf(PropTypes.string),
    error: PropTypes.string,
    className: PropTypes.string,
    styles: PropTypes.object,
    required: PropTypes.bool,
    disabled: PropTypes.bool,
};
export default ConditionSelect; 