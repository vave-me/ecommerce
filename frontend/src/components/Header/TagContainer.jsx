"use client";
import React, { memo } from "react";
import {useSelector} from "react-redux";
import {useTranslations} from "next-intl";
import {Bot, Grid3X3, Sparkles, TrendingUp} from "@/icons";
import {selectCurrentMode, selectIsAiMode, APP_MODES} from "../../redux/slices/appModeSlice";
import styles from "./TagContainer.module.css";
const TagContainer = memo(function TagContainer({showStats = true, compact = false}) {
    const t = useTranslations('TagContainer');
    const currentMode = useSelector(selectCurrentMode);
    const isAiMode = useSelector(selectIsAiMode);
    // Dynamic content based on current mode
    const modeConfig = {
        [APP_MODES.CLASSIC]: {
            icon: Grid3X3,
            primaryText: t('classic.primary'),
            secondaryText: t('classic.secondary'),
            accent: 'classic'
        },
        [APP_MODES.AI]: {
            icon: Bot,
            primaryText: t('ai.primary'),
            secondaryText: t('ai.secondary'),
            accent: 'ai'
        }
    };
    const config = modeConfig[currentMode] || modeConfig[APP_MODES.CLASSIC];
    // Sample dynamic stats - replace with real data
    const stats = {
        activeUsers: "12k+",
        listings: "50k+",
        categories: "8"
    };
    const containerClasses = [
        styles.container,
        compact ? styles.compact : '',
        styles[config.accent]
    ].filter(Boolean).join(' ');
    return (
        <div className={containerClasses}>
            {/* Mode indicator with primary message */}
            <div className={styles.modeIndicator}>
                <span className={styles.primaryText}>
                    {config.secondaryText}
                </span>
            </div>
            {/* Contextual info */}
            <div className={styles.contextInfo}>
                <span className={styles.secondaryText}>
                    {config.primaryText}
                </span>
            </div>
            {/* Dynamic stats (optional) */}
            {showStats && !compact && (
                <div className={styles.statsRow}>
                    <div className={styles.stat}>
                        <TrendingUp size={12} aria-hidden="true"/>
                        <span>{stats.activeUsers} {t('stats.users')}</span>
                    </div>
                    <div className={styles.stat}>
                        <Sparkles size={12} aria-hidden="true"/>
                        <span>{stats.listings} {t('stats.listings')}</span>
                    </div>
                </div>
            )}
        </div>
    );
});
export default TagContainer;