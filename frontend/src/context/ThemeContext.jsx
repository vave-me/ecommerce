"use client";
import React, { createContext, useContext, useEffect, useState, useCallback } from 'react';

const ThemeContext = createContext();

export const useTheme = () => {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error('useTheme must be used within a ThemeProvider');
  }
  return context;
};

export const ThemeProvider = ({ children }) => {
  const [theme, setTheme] = useState('system'); // 'light', 'dark', 'system'
  const [resolvedTheme, setResolvedTheme] = useState('light'); // The actual theme being used
  const [isHydrated, setIsHydrated] = useState(false);

  // Apply theme to document
  const applyTheme = useCallback((themeToApply) => {
    if (typeof document !== 'undefined') {
      const root = document.documentElement;
      
      // Remove existing theme classes
      root.classList.remove('theme-light', 'theme-dark');
      
      // Add new theme class
      root.classList.add(`theme-${themeToApply}`);
      
      // Set data attribute for CSS targeting
      root.setAttribute('data-theme', themeToApply);
      
      // Update meta theme-color for mobile browsers
      const metaThemeColor = document.querySelector('meta[name="theme-color"]');
      if (metaThemeColor) {
        metaThemeColor.setAttribute('content', themeToApply === 'dark' ? '#111827' : '#ffffff');
      }
    }
  }, []);

  // Get system theme preference
  const getSystemTheme = useCallback(() => {
    if (typeof window === 'undefined') return 'light';
    
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
  }, []);

  // Resolve theme based on user preference and system
  const resolveTheme = useCallback((userTheme) => {
    if (userTheme === 'system') {
      return getSystemTheme();
    }
    return userTheme;
  }, [getSystemTheme]);

  // Update theme
  const updateTheme = useCallback((newTheme) => {
    setTheme(newTheme);
    
    // Save to localStorage
    try {
      localStorage.setItem('theme', newTheme);
    } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
    
    // Resolve and apply theme
    const resolved = resolveTheme(newTheme);
    setResolvedTheme(resolved);
    applyTheme(resolved);
  }, [resolveTheme, applyTheme]);

  // Initialize theme on mount
  useEffect(() => {
    setIsHydrated(true);
    
    if (typeof window === 'undefined') return;

    let initialTheme = 'system';
    
    // Try to get saved theme from localStorage
    try {
      const savedTheme = localStorage.getItem('theme');
      if (savedTheme && ['light', 'dark', 'system'].includes(savedTheme)) {
        initialTheme = savedTheme;
      }
    } catch (error) {
        // Initialization error - log but continue
        if (process.env.NODE_ENV === 'development') {
            console.error('Initialization error:', error);
        }
        // Continue with default behavior
    }
    
    setTheme(initialTheme);
    const resolved = resolveTheme(initialTheme);
    setResolvedTheme(resolved);
    applyTheme(resolved);
  }, [resolveTheme, applyTheme]);

  // Listen for system theme changes
  useEffect(() => {
    if (typeof window === 'undefined') return;
    
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    
    const handleSystemThemeChange = (e) => {
      if (theme === 'system') {
        const newResolvedTheme = e.matches ? 'dark' : 'light';
        setResolvedTheme(newResolvedTheme);
        applyTheme(newResolvedTheme);
      }
    };
    
    mediaQuery.addEventListener('change', handleSystemThemeChange);
    
    return () => {
      mediaQuery.removeEventListener('change', handleSystemThemeChange);
    };
  }, [theme, applyTheme]);

  // Theme switching utilities
  const toggleTheme = useCallback(() => {
    const newTheme = resolvedTheme === 'light' ? 'dark' : 'light';
    updateTheme(newTheme);
  }, [resolvedTheme, updateTheme]);

  const setLightTheme = useCallback(() => updateTheme('light'), [updateTheme]);
  const setDarkTheme = useCallback(() => updateTheme('dark'), [updateTheme]);
  const setSystemTheme = useCallback(() => updateTheme('system'), [updateTheme]);

  const value = {
    theme, // User preference: 'light', 'dark', 'system'
    resolvedTheme, // Actual theme: 'light' or 'dark'
    isHydrated,
    updateTheme,
    toggleTheme,
    setLightTheme,
    setDarkTheme,
    setSystemTheme,
    isDark: resolvedTheme === 'dark',
    isLight: resolvedTheme === 'light',
    isSystem: theme === 'system'
  };

  return (
    <ThemeContext.Provider value={value}>
      {children}
    </ThemeContext.Provider>
  );
};

export default ThemeProvider; 