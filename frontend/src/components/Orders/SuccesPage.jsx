// SuccessPage.jsx
"use client";
import { FaCheckCircle } from '../../utils/iconImports';
import React, { memo } from 'react';
// Import the CSS module
import styles from './SuccessPage.module.css';
const SuccessPage = memo(() => {
    return (
        <div className={styles.container}>
            <div className={styles.iconWrapper}>
                <FaCheckCircle />
            </div>
            <h1 className={styles.heading}>Payment Successful!</h1>
            <p className={styles.message}>
                Thank you for your purchase. Your payment was processed successfully.
            </p>
            <a className={styles.button} href="/">
                Go to Homepage
            </a>
        </div>
    );
});
SuccessPage.displayName = 'SuccessPage';
export default SuccessPage;
