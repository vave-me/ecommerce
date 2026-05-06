import React from 'react';
import PropTypes from 'prop-types';
import { Tag } from "@/icons";
/**
 * Shared TagsInput component for consistent tags input field across all modals
 */
export function TagsInput({
    id = "tags",
    label = "Tags (comma separated)",
    value,
    onChange,
    placeholder = "summer, discount, limited",
    error,
    className = "",
    styles,
    required = false,
    disabled = false,
    showTagIcon = true,
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
                {showTagIcon && <Tag size={16} className={styles?.labelIcon || 'label-icon'} />}
                {label}
                {required && <span className={styles?.requiredMark || 'required-mark'}>*</span>}
            </label>
            <input
                id={id}
                className={`${styles?.formInput || 'form-input'} ${error ? styles?.inputError || 'input-error' : ''}`}
                type="text"
                value={value || ''}
                onChange={handleChange}
                placeholder={placeholder}
                disabled={disabled}
                aria-invalid={!!error}
            />
            {error && (
                <div className={styles?.fieldError || 'field-error'}>{error}</div>
            )}
        </div>
    );
}
TagsInput.propTypes = {
    id: PropTypes.string,
    label: PropTypes.string,
    value: PropTypes.string,
    onChange: PropTypes.func.isRequired,
    placeholder: PropTypes.string,
    error: PropTypes.string,
    className: PropTypes.string,
    styles: PropTypes.object,
    required: PropTypes.bool,
    disabled: PropTypes.bool,
    showTagIcon: PropTypes.bool,
};
export default TagsInput; 