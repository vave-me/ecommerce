// TextEditor/components/common/ProgressIndicator.jsx
import React from 'react';
import PropTypes from 'prop-types';
import styles from '../../TextEditor.module.css';
export function ProgressIndicator({ progress }) {
    return (
        <div
            className={styles.progressContainer}
            role="progressbar"
            aria-valuenow={progress}
            aria-valuemin="0"
            aria-valuemax="100"
        >
            <div className={styles.progressBarContainer}>
                <div
                    className={styles.progressBarFill}
                    style={{ width: `${progress}%` }}
                ></div>
            </div>
            <span className={styles.progressText}>{progress}%</span>
        </div>
    );
}
ProgressIndicator.propTypes = {
    progress: PropTypes.number.isRequired
};