"use client";
import React, { memo, useCallback } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useRouter, usePathname } from 'next/navigation';
import { useLocale } from 'next-intl';
import styles from './ModeSwitcher.module.css';

// Simple icons
const ClassicIcon = () => (
    <svg width="20" height="20" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth="2">
        <path strokeLinecap="round" strokeLinejoin="round" d="M4 6h16M4 12h16M4 18h7" />
    </svg>
);

const AIIcon = () => (
    <svg width="20" height="20" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth="2">
        <path strokeLinecap="round" strokeLinejoin="round" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
    </svg>
);

const ModeSwitcher = memo(() => {
    const dispatch = useDispatch();
    const router = useRouter();
    const pathname = usePathname();
    const locale = useLocale();
    
    // Get current mode from Redux
    const currentMode = useSelector(state => state.appMode?.currentMode || 'classic');
    const isAiMode = currentMode === 'ai';
    
    // Remove console.log to reduce noise
    // 
    
    const handleModeToggle = useCallback(() => {
        if (isAiMode) {
            // Switch to Classic mode
            // Update Redux state
            dispatch({ type: 'appMode/switchToClassicMode' });
            
            // Navigate away from AI page if we're on it
            if (pathname === `/${locale}/ai`) {
                router.push(`/${locale}`);
            }
        } else {
            // Switch to AI mode
            // Update Redux state
            dispatch({ type: 'appMode/switchToAiMode' });
            
            // Navigate to AI page
            router.push(`/${locale}/ai`);
        }
    }, [isAiMode, dispatch, pathname, locale, router]);
    
    return (
        <button
            onClick={handleModeToggle}
            className={`${styles.compactButton} ${isAiMode ? styles.classicMode : styles.aiMode}`}
            type="button"
            aria-label={isAiMode ? 'Switch to Classic Mode' : 'Switch to AI Mode'}
        >
            {isAiMode ? <ClassicIcon /> : <AIIcon />}
            <span className={styles.buttonText}>
                {isAiMode ? 'Classic' : 'AI Mode'}
            </span>
        </button>
    );
});

ModeSwitcher.displayName = 'ModeSwitcher';

export default ModeSwitcher;