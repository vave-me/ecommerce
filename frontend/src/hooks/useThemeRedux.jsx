import { useTheme } from '../context/ThemeContext';
// Mirror of the previous constant so existing imports don\'t break
export const THEME_MODES = {
    LIGHT: 'light',
    DARK: 'dark',
    SYSTEM: 'system',
};
/**
 * Compatibility hook – previously Redux-based, now forwards
 * to ThemeContext while keeping the same public API so that
 * existing components continue to work without code changes.
 */
export const useThemeRedux = () => {
    const {
        theme,
        resolvedTheme,
        isHydrated,
        updateTheme,
        toggleTheme,
        isDark,
        isLight,
        isSystem,
    } = useTheme();
    // Provide identical names/aliases as before
    const setTheme = updateTheme;
    return {
        // Theme state
        theme,
        resolvedTheme,
        systemTheme: resolvedTheme, // kept for API parity; not used
        isHydrated,
        // Booleans
        isDark,
        isLight,
        isSystem,
        // Actions
        updateTheme,
        toggleTheme,
        setTheme,
        // Consts
        THEME_MODES,
    };
};
export default useThemeRedux; // Preserve default export 