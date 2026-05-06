// common/components/ErrorAlert.jsx
import React from 'react';
import PropTypes from 'prop-types';
import {AlertCircle} from "@/icons";
import styles from './ErrorAlert.module.css'; // <-- IMPORT THE CSS MODULE
export function ErrorAlert({message}) {
    // Return null if there's no message to display
    if (!message) {
        return null;
    }
    return (
        // Apply the class name from the imported styles object
        <div className={styles.errorAlert} role="alert"> {/* Added role="alert" for accessibility */}
            <AlertCircle size={18} aria-hidden="true"/> {/* Hide decorative icon from screen readers */}
            <span>{message}</span>
        </div>
    );
}
ErrorAlert.propTypes = {
    // Ensure message is required and is a string
    message: PropTypes.string.isRequired
};