'use client';

import React from 'react';
import styles from './FeedLoadingState.module.css';

const FeedLoadingState = ({ message = "Loading feed..." }) => {
  return (
    <div className={styles.loadingContainer}>
      <div className={styles.loadingContent}>
        <div className={styles.spinnerWrapper}>
          <div className={styles.spinner}>
            <div className={styles.spinnerDot}></div>
            <div className={styles.spinnerDot}></div>
            <div className={styles.spinnerDot}></div>
          </div>
        </div>
        <p className={styles.loadingMessage}>{message}</p>
      </div>
    </div>
  );
};

export default FeedLoadingState;