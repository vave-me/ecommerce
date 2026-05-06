'use client';
import React from 'react';
import { useSelector } from 'react-redux';
import CompactVaveAd from './CompactVaveAd';
import styles from './ExpandableVaveAds.module.css';

/**
 * ExpandableVaveAds - Container for vaveme ads with Redux state management
 * Shows a compact ad below the UnifiedComposer if not dismissed
 */
const ExpandableVaveAds = () => {
    const isVaveAdVisible = useSelector(state => state.ads?.isVaveAdVisible);
    
    // Only show if explicitly true (not undefined or false)
    if (isVaveAdVisible !== true) {
        return null;
    }

    return (
        <div className={styles.container}>
            <div className={styles.compactSection}>
                <CompactVaveAd />
            </div>
        </div>
    );
};

export default ExpandableVaveAds;