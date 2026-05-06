import React from 'react';
import PropTypes from 'prop-types';
const OFFER_TYPES = [
    { value: 'sell', label: 'Sell' },
    { value: 'rent', label: 'Rent' },
    { value: 'lease', label: 'Lease' },
    { value: 'auction', label: 'Auction' }
];
export function OfferTypeSelect({
    id,
    value,
    onChange,
    label = "Type of Offer",
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
            <select
                id={id}
                className={`${styles.formSelect} ${error ? styles.inputError : ''}`}
                value={value}
                onChange={(e) => onChange(e.target.value)}
                disabled={disabled}
                {...props}
            >
                <option value="">Select offer type</option>
                {OFFER_TYPES.map((type) => (
                    <option key={type.value} value={type.value}>
                        {type.label}
                    </option>
                ))}
            </select>
            {error && (
                <div className={styles.fieldError}>{error}</div>
            )}
        </div>
    );
}
OfferTypeSelect.propTypes = {
    id: PropTypes.string.isRequired,
    value: PropTypes.string.isRequired,
    onChange: PropTypes.func.isRequired,
    label: PropTypes.string,
    error: PropTypes.string,
    styles: PropTypes.object.isRequired,
    required: PropTypes.bool,
    disabled: PropTypes.bool,
}; 