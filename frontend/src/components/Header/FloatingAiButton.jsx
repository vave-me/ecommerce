"use client";
import React, { useCallback, memo } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useTranslations } from 'next-intl';
import { 
    switchToAiMode, 
    setModeTransitioning,
    selectCurrentMode,
    selectIsAiMode,
    selectIsTransitioning,
    APP_MODES
} from '../../redux/slices/appModeSlice';
import styles from './FloatingAiButton.module.css';
// Exact same AIIcon as in ModeSwitcher
const AIIcon = memo(({ size = 24, className }) => (
    <svg 
        className={className} 
        width={size} 
        height={size} 
        fill="none" 
        viewBox="0 0 24 24" 
        stroke="currentColor" 
        strokeWidth="2"
    >
        <path strokeLinecap="round" strokeLinejoin="round" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
    </svg>
));
AIIcon.displayName = 'AIIcon';
/**
 * Enhanced FloatingAiButton Component
 * 
 * Floating action button with always visible "AI MODE" label
 * Uses blue color scheme without gradient
 * Fixed size, no hover expansion
 */
const FloatingAiButton = ({ 
    variant = 'pill' // 'pill', 'round', 'minimal'
}) => {
    const dispatch = useDispatch();
    const t = useTranslations('ModeSwitcher');
    const currentMode = useSelector(selectCurrentMode);
    const isAiMode = useSelector(selectIsAiMode);
    const isTransitioning = useSelector(selectIsTransitioning);
    /**
     * Handle switch to AI mode with enhanced feedback
     */
    const handleSwitchToAi = useCallback(async () => {
        if (currentMode === APP_MODES.AI || isTransitioning) return;
        try {
            // Start transition
            dispatch(setModeTransitioning(true));
            // Small delay for smooth UX
            await new Promise(resolve => setTimeout(resolve, 150));
            // Switch to AI mode - layout will handle content rendering
            dispatch(switchToAiMode());
            // End transition after animation
            setTimeout(() => {
                dispatch(setModeTransitioning(false));
            }, 350);
        } catch (error) {
            dispatch(setModeTransitioning(false));
        }
    }, [currentMode, isTransitioning, dispatch]);
    // Don't show in AI mode
    if (isAiMode) {
        return null;
    }
    // Dynamic button classes - always expanded for pill variant
    const buttonClasses = [
        styles.floatingButton,
        styles[variant],
        isTransitioning ? styles.disabled : '',
        variant === 'pill' ? styles.expanded : '' // Always expanded for pill variant
    ].filter(Boolean).join(' ');
    // Pill variant - with always visible label
    if (variant === 'pill') {
        return (
            <button
                className={buttonClasses}
                onClick={handleSwitchToAi}
                disabled={isTransitioning}
                aria-label={t('aiModeLabel')}
                title={t('aiModeTooltip')}
            >
                <div className={styles.buttonContent}>
                    <div className={styles.iconContainer}>
                        <AIIcon 
                            size={24} 
                            className={styles.icon}
                        />
                    </div>
                    <div className={styles.labelContainer}>
                        <span className={styles.labelText}>
                            AI MODE
                        </span>
                    </div>
                </div>
                <span className={styles.pulse} aria-hidden="true" />
            </button>
        );
    }
    // Minimal variant - icon only, smaller
    if (variant === 'minimal') {
        return (
            <button
                className={buttonClasses}
                onClick={handleSwitchToAi}
                disabled={isTransitioning}
                aria-label={t('aiModeLabel')}
                title={t('aiModeTooltip')}
            >
                <AIIcon 
                    size={18} 
                    className={styles.icon}
                />
            </button>
        );
    }
    // Default round variant (original behavior)
    return (
        <button
            className={buttonClasses}
            onClick={handleSwitchToAi}
            disabled={isTransitioning}
            aria-label={t('aiModeLabel')}
            title={t('aiModeTooltip')}
        >
            <AIIcon 
                size={24} 
                className={styles.icon}
            />
            <span className={styles.pulse} aria-hidden="true" />
        </button>
    );
};
export default memo(FloatingAiButton); 