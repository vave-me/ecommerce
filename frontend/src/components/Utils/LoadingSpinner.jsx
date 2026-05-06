// src/components/LoadingSpinner.jsx
import React, {memo} from 'react';
import styles from './LoadingSpinner.module.css';
const LoadingSpinner = memo(() => {
    return <div className={styles.spinner}/>;
});
LoadingSpinner.displayName = 'LoadingSpinner';
export default LoadingSpinner;
