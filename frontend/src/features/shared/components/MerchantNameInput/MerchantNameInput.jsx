import React from 'react';
import PropTypes from 'prop-types';
export function MerchantNameInput({
    id,
    value,
    onChange,
    label = "Merchant Name",
    placeholder = "Enter merchant or business name",
    error,
    styles,
    required = false,
    disabled = false,
    maxLength = 100,
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
                type="text"
                className={`${styles.formInput} ${error ? styles.inputError : ''}`}
                value={value}
                onChange={(e) => onChange(e.target.value)}
                placeholder={placeholder}
                disabled={disabled}
                maxLength={maxLength}
                {...props}
            />
            {error && (
                <div className={styles.fieldError}>{error}</div>
            )}
        </div>
    );
}
MerchantNameInput.propTypes = {
    id: PropTypes.string.isRequired,
    value: PropTypes.string.isRequired,
    onChange: PropTypes.func.isRequired,
    label: PropTypes.string,
    placeholder: PropTypes.string,
    error: PropTypes.string,
    styles: PropTypes.object.isRequired,
    required: PropTypes.bool,
    disabled: PropTypes.bool,
    maxLength: PropTypes.number,
}; 