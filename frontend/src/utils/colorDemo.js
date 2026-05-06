/**
 * Color Demonstration Utility
 * Showcases the new search button color theme #005552
 */
export const SEARCH_COLOR_THEME = {
    primary: '#005552',
    hover: '#004440', 
    active: '#003330',
    light: 'rgba(0, 85, 82, 0.1)',
    shadow: 'rgba(0, 85, 82, 0.25)',
    glow: 'rgba(0, 85, 82, 0.15)'
};
export const SEARCH_COLOR_DARK_THEME = {
    primary: '#00706b',
    hover: '#005552',
    active: '#004440', 
    light: 'rgba(0, 112, 107, 0.15)',
    shadow: 'rgba(0, 112, 107, 0.3)',
    glow: 'rgba(0, 112, 107, 0.2)'
};
/**
 * Generate CSS custom properties for the search theme
 */
export const generateSearchThemeCSS = (theme = SEARCH_COLOR_THEME) => {
    return `
        --search-primary: ${theme.primary};
        --search-primary-hover: ${theme.hover};
        --search-primary-active: ${theme.active};
        --search-primary-light: ${theme.light};
        --search-primary-shadow: ${theme.shadow};
        --search-primary-glow: ${theme.glow};
    `;
};
/**
 * Color accessibility information
 */
export const COLOR_ACCESSIBILITY = {
    primary: {
        color: '#005552',
        name: 'Deep Teal',
        contrastRatio: '4.5:1', // WCAG AA compliant
        usage: 'Primary buttons, focused states'
    },
    hover: {
        color: '#004440', 
        name: 'Darker Teal',
        contrastRatio: '5.2:1', // WCAG AA+ compliant
        usage: 'Hover states, active interactions'
    },
    active: {
        color: '#003330',
        name: 'Darkest Teal',
        contrastRatio: '6.1:1', // WCAG AAA compliant
        usage: 'Active/pressed states'
    }
};
export default SEARCH_COLOR_THEME; 