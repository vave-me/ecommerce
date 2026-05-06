/**
 * Utility functions for handling login redirects
 */
/**
 * Get the current page URL for use as a redirect parameter
 * @returns {string} Current URL (pathname + search params)
 */
export const getCurrentPageUrl = () => {
    if (typeof window === 'undefined') return '/';
    return window.location.pathname + window.location.search;
};
/**
 * Redirect to login page with current page as return URL
 * @param {string} currentUrl - Optional current URL, will auto-detect if not provided
 */
export const redirectToLogin = (currentUrl = null) => {
    const returnUrl = currentUrl || getCurrentPageUrl();
    const loginUrl = `/login?redirect=${encodeURIComponent(returnUrl)}`;
    if (typeof window !== 'undefined') {
        window.location.href = loginUrl;
    }
};
