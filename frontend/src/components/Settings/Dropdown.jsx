// src/components/Dropdown.jsx
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import { ChevronDown } from '@/icons';
import styles from './Dropdown.module.css';
const Dropdown = memo(({ label, options, value, onChange, id }) => {
    // Generate a unique ID if none is provided
    const selectId = id || `dropdown-${Math.random().toString(36).substr(2, 9)}`;
    return (
        <div className={styles.container}>
            {label && (
                <label htmlFor={selectId} className={styles.label}>
                    {label}
                </label>
            )}
            <div className={styles.selectWrapper}>
                <select
                    id={selectId}
                    className={styles.select}
                    value={value}
                    onChange={onChange}
                >
                    {options.map((option, index) => (
                        <option key={index} value={option}>
                            {option}
                        </option>
                    ))}
                </select>
                <ChevronDown className={styles.icon} aria-hidden="true" />
            </div>
        </div>
    );
});
Dropdown.displayName = 'Dropdown';
Dropdown.propTypes = {
    label: PropTypes.string,
    options: PropTypes.arrayOf(PropTypes.string).isRequired,
    value: PropTypes.string.isRequired,
    onChange: PropTypes.func.isRequired,
    id: PropTypes.string
};
export default Dropdown;