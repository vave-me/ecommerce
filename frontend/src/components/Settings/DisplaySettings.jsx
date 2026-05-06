// src/components/DisplaySettings.jsx
"use client"
import { FaCheckCircle, FaColumns, FaMoon, FaPaintBrush, FaSun, FaTextHeight, FaPalette, FaEye } from '../../utils/iconImports';
import React, {useCallback, useMemo, useState, memo, useEffect} from "react";
import {useTranslations} from 'next-intl'; //  Import hook
import useThemeRedux from '../../hooks/useThemeRedux';
import Dropdown from "./Dropdown"; // Assume Dropdown receives translated label/options
import styles from "./DisplaySettings.module.css";
// Define keys for options
const FONT_SIZE_KEYS = ['small', 'medium', 'large'];
const LAYOUT_KEYS = ['default', 'compact', 'spacious'];
const THEME_KEYS = ['light', 'dark', 'system'];
const TOGGLE_KEYS = ['reduceMotion', 'highContrast'];
/**
 * DisplaySettings Component with Translations
 * Allows users to configure display settings.
 */
const DisplaySettings = memo(() => {
    const t = useTranslations('DisplaySettings'); //  Instantiate hook
    const { theme, updateTheme, resolvedTheme } = useThemeRedux();
    // State for display settings - use keys where applicable
    const [display, setDisplay] = useState({
        theme: theme || "system", // Use actual theme from context
        fontSize: "medium", // Use key: small, medium, large
        layout: "default", // Use key: default, compact, spacious
        reduceMotion: false,
        highContrast: false
    });
    // State for feedback message (stores key and values)
    const [feedback, setFeedback] = useState({
        show: false,
        key: "",
        values: {}
    });
    // Sync local display state with theme context
    useEffect(() => {
        setDisplay(prev => ({...prev, theme: theme}));
    }, [theme]);
    // Show feedback message using translation key
    const showFeedback = useCallback((key, values = {}, duration = 3000) => {
        setFeedback({show: true, key, values});
        const timer = setTimeout(() => {
            setFeedback({show: false, key: "", values: {}});
        }, duration);
        // Optional: Store timer to clear on unmount if needed
        // return () => clearTimeout(timer);
    }, []); // t is stable, no need to include if not directly used here
    // Translate options for dropdowns
    const translatedFontSizeOptions = useMemo(() => FONT_SIZE_KEYS.map(key => t(`fontSize_${key}`)), [t]);
    const translatedLayoutOptions = useMemo(() => LAYOUT_KEYS.map(key => t(`layout_${key}`)), [t]);
    // Map keys back to translated values for feedback messages if needed, or use keys directly
    // This map helps get the *translated* value for feedback interpolation
    const keyToTranslatedValueMap = useMemo(() => ({
        ...FONT_SIZE_KEYS.reduce((acc, key) => ({...acc, [key]: t(`fontSize_${key}`)}), {}),
        ...LAYOUT_KEYS.reduce((acc, key) => ({...acc, [key]: t(`layout_${key}`)}), {}),
        ...THEME_KEYS.reduce((acc, key) => ({...acc, [key]: t(`theme_${key}`)}), {}),
        ...TOGGLE_KEYS.reduce((acc, key) => ({...acc, [key]: t(`setting_${key}`)}), {}),
    }), [t]);
    // Handle theme change using keys
    const handleThemeChange = useCallback((themeKey) => {
        setDisplay(prev => ({...prev, theme: themeKey}));
        // Update the actual theme context
        updateTheme(themeKey);
        // Use specific keys for theme feedback
        showFeedback(`feedbackTheme${themeKey.charAt(0).toUpperCase() + themeKey.slice(1)}`);
    }, [showFeedback, updateTheme]); // Added updateTheme dependency
    // Handle font size change using keys
    const handleFontSizeChange = useCallback((translatedValue) => {
        // Find the key corresponding to the translated value
        const key = FONT_SIZE_KEYS.find(k => t(`fontSize_${k}`) === translatedValue) || display.fontSize;
        setDisplay(prev => ({...prev, fontSize: key}));
        showFeedback('feedbackFontSizeChange', {value: translatedValue});
    }, [t, display.fontSize, showFeedback]);
    // Handle layout change using keys
    const handleLayoutChange = useCallback((translatedValue) => {
        // Find the key corresponding to the translated value
        const key = LAYOUT_KEYS.find(k => t(`layout_${k}`) === translatedValue) || display.layout;
        setDisplay(prev => ({...prev, layout: key}));
        showFeedback('feedbackLayoutChange', {value: translatedValue});
    }, [t, display.layout, showFeedback]);
    // Handle toggle changes using keys
    const handleToggleChange = useCallback((settingKey) => { // settingKey is 'reduceMotion' or 'highContrast'
        const newValue = !display[settingKey];
        setDisplay(prev => ({...prev, [settingKey]: newValue}));
        const feedbackKey = newValue ? 'feedbackToggleEnabled' : 'feedbackToggleDisabled';
        const settingName = keyToTranslatedValueMap[settingKey] || settingKey; // Get translated setting name
        showFeedback(feedbackKey, {settingName: settingName});
    }, [display, showFeedback, keyToTranslatedValueMap]);
    // Handle save
    const handleSave = useCallback(() => {
        // Save display settings logic...
        // showFeedback("feedbackSaveSuccess"); //  Use translation key
    }, [showFeedback]); // Removed 'display' dependency unless save logic needs it immediately
    return (
        <div className={styles.container}>
            {/*   Use translation */}
            <h2 className={styles.pageTitle}>{t('pageTitle')}</h2>
            {/*   Use translation */}
            <p className={styles.pageDescription}>{t('pageDescription')}</p>
            {/* Theme Selection Section */}
            <section className={styles.section} aria-labelledby="theme-section-title">
                <h3 id="theme-section-title" className={styles.sectionTitle}>
                    <FaPaintBrush className={styles.sectionIcon} aria-hidden="true"/>
                    {/*   Use translation */}
                    {t('themeTitle')}
                </h3>
                <div className={styles.themeOptions} role="radiogroup" aria-labelledby="theme-section-title">
                    {/* Light Theme Card */}
                    <div
                        className={`${styles.themeCard} ${display.theme === 'light' ? styles.selectedTheme : ''}`}
                        onClick={() => handleThemeChange('light')}
                        onKeyDown={(e) => (e.key === 'Enter' || e.key === ' ') && handleThemeChange('light')}
                        role="radio"
                        aria-checked={display.theme === 'light'}
                        tabIndex={0} // Make focusable
                        aria-label={t('theme_light')} // Add ARIA label
                    >
                        <div className={styles.themePreview} aria-hidden="true">
                            <div className={styles.lightPreview}>
                                <div className={styles.previewHeader}></div>
                                <div className={styles.previewContent}>
                                    <div className={styles.previewLine}></div>
                                    <div className={styles.previewLine}></div>
                                    <div className={styles.previewLine}></div>
                                </div>
                            </div>
                        </div>
                        <div className={styles.themeInfo}>
                            <FaSun className={styles.themeIcon} aria-hidden="true"/>
                            {/*   Use translation */}
                            <span className={styles.themeName}>{t('theme_light')}</span>
                        </div>
                    </div>
                    {/* Dark Theme Card */}
                    <div
                        className={`${styles.themeCard} ${display.theme === 'dark' ? styles.selectedTheme : ''}`}
                        onClick={() => handleThemeChange('dark')}
                        onKeyDown={(e) => (e.key === 'Enter' || e.key === ' ') && handleThemeChange('dark')}
                        role="radio"
                        aria-checked={display.theme === 'dark'}
                        tabIndex={0}
                        aria-label={t('theme_dark')}
                    >
                        <div className={styles.themePreview} aria-hidden="true">
                            <div className={styles.darkPreview}>
                                <div className={styles.previewHeader}></div>
                                <div className={styles.previewContent}>
                                    <div className={styles.previewLine}></div>
                                    <div className={styles.previewLine}></div>
                                    <div className={styles.previewLine}></div>
                                </div>
                            </div>
                        </div>
                        <div className={styles.themeInfo}>
                            <FaMoon className={styles.themeIcon} aria-hidden="true"/>
                            <span className={styles.themeName}>{t('theme_dark')}</span>
                        </div>
                    </div>
                    {/* System Theme Card */}
                    <div
                        className={`${styles.themeCard} ${display.theme === 'system' ? styles.selectedTheme : ''}`}
                        onClick={() => handleThemeChange('system')}
                        onKeyDown={(e) => (e.key === 'Enter' || e.key === ' ') && handleThemeChange('system')}
                        role="radio"
                        aria-checked={display.theme === 'system'}
                        tabIndex={0}
                        aria-label={t('theme_system')}
                    >
                        <div className={styles.themePreview} aria-hidden="true">
                            <div className={styles.systemPreview}>
                                <div className={styles.previewSplit}>
                                    <div className={styles.previewLight}></div>
                                    <div className={styles.previewDark}></div>
                                </div>
                            </div>
                        </div>
                        <div className={styles.themeInfo}>
                            {/*   Use translation */}
                            <span className={styles.themeName}>{t('theme_system')}</span>
                        </div>
                    </div>
                </div>
            </section>
            {/* Font Size Section */}
            <section className={styles.section} aria-labelledby="font-section-title">
                <h3 id="font-section-title" className={styles.sectionTitle}>
                    <FaTextHeight className={styles.sectionIcon} aria-hidden="true"/>
                    {/*   Use translation */}
                    {t('fontTitle')}
                </h3>
                <div className={styles.fontSizeOptions}>
                    <div className={styles.optionCard}>
                        <div className={styles.dropdownContainer}>
                            <Dropdown
                                //   Use translation for label and options
                                label={t('fontLabel')}
                                options={translatedFontSizeOptions} // Pass translated strings
                                // Pass currently selected *translated* value for Dropdown display
                                value={keyToTranslatedValueMap[display.fontSize]}
                                onChange={handleFontSizeChange} // Handles mapping back to key
                            />
                        </div>
                        <div className={styles.sizePreview} aria-live="polite"> {/* Announce changes */}
                            {/*   Use translation */}
                            <p className={styles.sizeExampleSmall} hidden={display.fontSize !== 'small'}>
                                {t('fontPreviewSmall')}
                            </p>
                            <p className={styles.sizeExampleMedium} hidden={display.fontSize !== 'medium'}>
                                {t('fontPreviewMedium')}
                            </p>
                            <p className={styles.sizeExampleLarge} hidden={display.fontSize !== 'large'}>
                                {t('fontPreviewLarge')}
                            </p>
                        </div>
                    </div>
                </div>
            </section>
            {/* Layout Preferences Section */}
            <section className={styles.section} aria-labelledby="layout-section-title">
                <h3 id="layout-section-title" className={styles.sectionTitle}>
                    <FaColumns className={styles.sectionIcon} aria-hidden="true"/>
                    {/*   Use translation */}
                    {t('layoutTitle')}
                </h3>
                <div className={styles.layoutOptions}>
                    <div className={styles.optionCard}>
                        <div className={styles.dropdownContainer}>
                            <Dropdown
                                //   Use translation for label and options
                                label={t('layoutLabel')}
                                options={translatedLayoutOptions} // Pass translated strings
                                // Pass currently selected *translated* value for Dropdown display
                                value={keyToTranslatedValueMap[display.layout]}
                                onChange={handleLayoutChange} // Handles mapping back to key
                            />
                        </div>
                        <div className={styles.layoutInfo} aria-live="polite"> {/* Announce changes */}
                            {/*   Use translation based on selected layout key */}
                            {display.layout === 'default' &&
                                <p className={styles.layoutDescription}>{t('layoutDescDefault')}</p>}
                            {display.layout === 'compact' &&
                                <p className={styles.layoutDescription}>{t('layoutDescCompact')}</p>}
                            {display.layout === 'spacious' &&
                                <p className={styles.layoutDescription}>{t('layoutDescSpacious')}</p>}
                        </div>
                    </div>
                </div>
            </section>
            {/* Accessibility Options */}
            <section className={styles.section} aria-labelledby="accessibility-section-title">
                <h3 id="accessibility-section-title" className={styles.sectionTitle}>
                    {/* Note: Re-using FaPaintBrush icon, consider a dedicated accessibility icon */}
                    <FaPaintBrush className={styles.sectionIcon} aria-hidden="true"/>
                    {/*   Use translation */}
                    {t('accessibilityTitle')}
                </h3>
                <div className={styles.accessibilityOptions}>
                    {/* Reduce Motion */}
                    <div className={styles.optionItem}>
                        <div className={styles.optionInfo}>
                            {/*   Use translation */}
                            <h4 className={styles.optionTitle}>{t('setting_reduceMotion')}</h4>
                            {/*   Use translation */}
                            <p className={styles.optionDescription}>{t('settingDescReduceMotion')}</p>
                        </div>
                        <label className={styles.switch} htmlFor="toggle-reduceMotion">
                            <input
                                type="checkbox"
                                id="toggle-reduceMotion"
                                checked={display.reduceMotion}
                                onChange={() => handleToggleChange('reduceMotion')}
                                className={styles.switchInput}
                                //   Use translation for aria-label
                                aria-label={t('setting_reduceMotion')}
                            />
                            <span className={styles.switchSlider} aria-hidden="true"></span>
                        </label>
                    </div>
                    {/* High Contrast */}
                    <div className={styles.optionItem}>
                        <div className={styles.optionInfo}>
                            {/*   Use translation */}
                            <h4 className={styles.optionTitle}>{t('setting_highContrast')}</h4>
                            {/*   Use translation */}
                            <p className={styles.optionDescription}>{t('settingDescHighContrast')}</p>
                        </div>
                        <label className={styles.switch} htmlFor="toggle-highContrast">
                            <input
                                type="checkbox"
                                id="toggle-highContrast"
                                checked={display.highContrast}
                                onChange={() => handleToggleChange('highContrast')}
                                className={styles.switchInput}
                                //   Use translation for aria-label
                                aria-label={t('setting_highContrast')}
                            />
                            <span className={styles.switchSlider} aria-hidden="true"></span>
                        </label>
                    </div>
                </div>
            </section>
            <div className={styles.actions}>
                {/*   Use translation */}
                <button className={styles.saveButton} onClick={handleSave}>
                    {t('saveButton')}
                </button>
            </div>
            {/* Feedback message */}
            {feedback.show && feedback.key && (
                // Use translation with key and values
                <div className={styles.feedback} role="status"> {/* Use status role */}
                    <FaCheckCircle className={styles.feedbackIcon} aria-hidden="true"/>
                    <span>{t(feedback.key, feedback.values)}</span>
                </div>
            )}
        </div>
    );
});
DisplaySettings.displayName = 'DisplaySettings';
// Add PropTypes if needed for Dropdown or other passed props
// DisplaySettings.propTypes = { ... };
export default DisplaySettings;