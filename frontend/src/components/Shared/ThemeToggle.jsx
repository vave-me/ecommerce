"use client";
import React from 'react';
import useThemeRedux from '../../hooks/useThemeRedux';
import styles from './ThemeToggle.module.css';
const ThemeToggle = ({ 
  showLabel = true, 
  size = 'medium',
  variant = 'button' // 'button', 'switch', 'dropdown'
}) => {
  const { theme, resolvedTheme, updateTheme, isDark, isLight, isSystem } = useThemeRedux();
  if (variant === 'dropdown') {
    return (
      <div className={`${styles.themeDropdown} ${styles[size]}`}>
        <select 
          value={theme} 
          onChange={(e) => updateTheme(e.target.value)}
          className={styles.select}
          aria-label="Choose theme"
        >
          <option value="light">☀️ Light</option>
          <option value="dark">🌙 Dark</option>
          <option value="system">💻 System</option>
        </select>
      </div>
    );
  }
  if (variant === 'switch') {
    return (
      <div className={`${styles.themeSwitch} ${styles[size]}`}>
        {showLabel && (
          <span className={styles.label}>
            {isDark ? '🌙' : '☀️'} {resolvedTheme === 'dark' ? 'Dark' : 'Light'}
          </span>
        )}
        <button
          onClick={() => updateTheme(isDark ? 'light' : 'dark')}
          className={`${styles.switchButton} ${isDark ? styles.dark : styles.light}`}
          aria-label={`Switch to ${isDark ? 'light' : 'dark'} mode`}
          title={`Currently ${resolvedTheme} mode. Click to switch.`}
        >
          <span className={styles.switchIcon}>
            {isDark ? '☀️' : '🌙'}
          </span>
        </button>
      </div>
    );
  }
  // Default button variant
  return (
    <div className={`${styles.themeToggle} ${styles[size]}`}>
      <div className={styles.buttonGroup}>
        <button
          onClick={() => updateTheme('light')}
          className={`${styles.themeButton} ${isLight && !isSystem ? styles.active : ''}`}
          aria-label="Light theme"
          title="Switch to light theme"
        >
          ☀️
          {showLabel && <span>Light</span>}
        </button>
        <button
          onClick={() => updateTheme('dark')}
          className={`${styles.themeButton} ${isDark && !isSystem ? styles.active : ''}`}
          aria-label="Dark theme"
          title="Switch to dark theme"
        >
          🌙
          {showLabel && <span>Dark</span>}
        </button>
        <button
          onClick={() => updateTheme('system')}
          className={`${styles.themeButton} ${isSystem ? styles.active : ''}`}
          aria-label="System theme"
          title="Use system theme preference"
        >
          💻
          {showLabel && <span>Auto</span>}
        </button>
      </div>
      {isSystem && (
        <div className={styles.systemIndicator}>
          <small>Following system ({resolvedTheme})</small>
        </div>
      )}
    </div>
  );
};
export default ThemeToggle; 