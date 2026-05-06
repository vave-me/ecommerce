// common/components/FormField.jsx
import React from 'react';
import PropTypes from 'prop-types';
import styles from './FormField.module.css';
export function FormField({
                              id,
                              label,
                              type = 'text',
                              value,
                              onChange,
                              error,
                              placeholder,
                              required = false,
                              maxLength,
                              hint,
                              icon,
                              showCharCount = false
                          }) {
    return (
        <div className={styles.formGroup}>
            <label className={styles.formLabel} htmlFor={id}>
                {label} {required && <span className={styles.requiredMark}>*</span>}
            </label>
            <div className={icon ? styles.inputWithIcon : ''}>
                {icon && <span className={styles.inputIcon}>{icon}</span>}
                <input
                    id={id}
                    className={`${styles.formInput} ${error ? styles.inputError : ''}`}
                    type={type}
                    value={value}
                    onChange={onChange}
                    placeholder={placeholder}
                    required={required}
                    maxLength={maxLength}
                    aria-invalid={error ? 'true' : 'false'}
                    aria-describedby={error ? `${id}-error` : hint ? `${id}-hint` : undefined}
                />
            </div>
            {error ? (
                <div className={styles.fieldError} id={`${id}-error`} role="alert">
                    {error}
                </div>
            ) : hint ? (
                <div className={styles.inputHint} id={`${id}-hint`}>
                    {hint}
                </div>
            ) : null}
            {showCharCount && maxLength && (
                <div className={styles.charCount} aria-live="polite">
                    {value.length}/{maxLength}
                </div>
            )}
        </div>
    );
}
FormField.propTypes = {
    id: PropTypes.string.isRequired,
    label: PropTypes.string.isRequired,
    type: PropTypes.string,
    value: PropTypes.string.isRequired,
    onChange: PropTypes.func.isRequired,
    error: PropTypes.string,
    placeholder: PropTypes.string,
    required: PropTypes.bool,
    maxLength: PropTypes.number,
    hint: PropTypes.string,
    icon: PropTypes.node,
    showCharCount: PropTypes.bool
};