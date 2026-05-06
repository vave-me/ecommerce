// src/components/ToggleSwitch.jsx
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import styles from './ToggleSwitch.module.css';
const ToggleSwitch = memo(({ isOn, handleToggle, label, id, ariaLabel }) => {
    // Generate a unique ID if none is provided
    const switchId = id || `switch-${Math.random().toString(36).substr(2, 9)}`;
    return (
        <div className={styles.container}>
            <label className={styles.switch} htmlFor={switchId}>
                <input
                    id={switchId}
                    type="checkbox"
                    checked={isOn}
                    onChange={handleToggle}
                    className={styles.checkbox}
                    aria-label={ariaLabel || label}
                />
                <span className={styles.slider}></span>
            </label>
            {label && <span className={styles.label}>{label}</span>}
        </div>
    );
});
ToggleSwitch.displayName = 'ToggleSwitch';
ToggleSwitch.propTypes = {
    isOn: PropTypes.bool.isRequired,
    handleToggle: PropTypes.func.isRequired,
    label: PropTypes.string,
    id: PropTypes.string,
    ariaLabel: PropTypes.string
};
export default ToggleSwitch;