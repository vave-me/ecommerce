// TextEditor/components/common/ErrorMessage.jsx
import React from 'react';
import PropTypes from 'prop-types';
import styles from '../../TextEditor.module.css';
export function ErrorMessage({ message }) {
    if (!message) return null;
    return (
        <div className={styles.errorMessage} role="alert">
            {message}
        </div>
    );
}
ErrorMessage.propTypes = {
    message: PropTypes.string
};