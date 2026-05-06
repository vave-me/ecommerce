import React from 'react';
import PropTypes from 'prop-types';
/**
 * Lightweight inline PriceInput component for use within complex layouts
 * This doesn't include form group wrapper or labels - just the input with currency symbol
 */
export function PriceInputInline({
    id,
    value,
    onChange,
    className = "",
    styles,
    disabled = false,
    placeholder = "0.00",
    min = "0",
    step = "0.01",
    currency = "€",
    style,
    ...props
}) {
    const handleChange = (e) => {
        if (onChange) {
            onChange(e.target.value);
        }
    };
    return (
        <div
            className={`${styles?.priceInputWrapper || 'price-input-wrapper'} ${className}`}
            style={style}
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
                {...props}
            />
        </div>
    );
}
PriceInputInline.propTypes = {
    id: PropTypes.string,
    value: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
    onChange: PropTypes.func.isRequired,
    className: PropTypes.string,
    styles: PropTypes.object,
    disabled: PropTypes.bool,
    placeholder: PropTypes.string,
    min: PropTypes.string,
    step: PropTypes.string,
    currency: PropTypes.string,
    style: PropTypes.object,
};
export default PriceInputInline; 