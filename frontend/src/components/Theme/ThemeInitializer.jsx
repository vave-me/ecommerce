"use client";
import { useEffect } from 'react';
import { useTheme } from '../../context/ThemeContext';
/**
 * Redux Theme Initializer Component
 * Ensures the Redux theme system is properly initialized and synchronized
 */
const ThemeInitializer = () => {
    const { theme, resolvedTheme, isHydrated } = useTheme();
    // Debug logging in development
    useEffect(() => {
        if (process.env.NODE_ENV === 'development' && isHydrated) {
        }
    }, [theme, resolvedTheme, isHydrated]);
    // This component doesn't render anything - it's just for initialization
    return null;
};
export default ThemeInitializer; 