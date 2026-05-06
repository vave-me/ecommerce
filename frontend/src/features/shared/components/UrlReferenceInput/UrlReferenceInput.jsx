import React from 'react';
import PropTypes from 'prop-types';
export function UrlReferenceInput({
    id,
    value,
    onChange,
    label = "URL Reference",
    placeholder = "Enter reference URL",
    error,
    styles,
    required = false,
    disabled = false,
    ...props
}) {
    return (
        <div className={styles.formGroup}>
            <label className={styles.formLabel} htmlFor={id}>
                {label}
                {required && <span className={styles.requiredMark}>*</span>}
            </label>
            <input
                id={id}
                type="url"
                className={`${styles.formInput} ${error ? styles.inputError : ''}`}
                value={value}
                onChange={(e) => onChange(e.target.value)}
                placeholder={placeholder}
                disabled={disabled}
                {...props}
            />
            {error && (
                <div className={styles.fieldError}>{error}</div>
            )}
        </div>
    );
}
UrlReferenceInput.propTypes = {
    id: PropTypes.string.isRequired,
    value: PropTypes.string.isRequired,
    onChange: PropTypes.func.isRequired,
    label: PropTypes.string,
    placeholder: PropTypes.string,
    error: PropTypes.string,
    styles: PropTypes.object.isRequired,
    required: PropTypes.bool,
    disabled: PropTypes.bool,
}; 