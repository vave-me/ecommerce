/**
 * CENTRALIZED IMAGE CONSTANTS
 * 
 * Production-ready image path management to prevent 404 errors
 * and improve maintainability across the application.
 * 
 * USAGE:
 * import { IMAGES } from '@/utils/imageConstants';
 * <Image src={IMAGES.LOGOS.SFX} alt="Logo" />
 */
// Base paths for different image categories
const BASE_PATHS = {
    IMAGES: '/images',
    ICONS: '/images',
    LOGOS: '/images',
    PLACEHOLDERS: '/images'
};
// Logo configurations
export const LOGOS = {
    SFX: `${BASE_PATHS.LOGOS}/sfx.png`,
    SFX_SVG: `${BASE_PATHS.LOGOS}/sfx.svg`,
    MOBILE: `${BASE_PATHS.LOGOS}/logo-mobile.webp`,
    SMALL: `${BASE_PATHS.LOGOS}/logo-small.webp`,
    BLACK: `${BASE_PATHS.LOGOS}/logo-black.webp`,
    // Fallback configurations
    FALLBACKS: {
        DESKTOP: `${BASE_PATHS.LOGOS}/logo-black.webp`,
        MOBILE: `${BASE_PATHS.LOGOS}/logo-black.webp`
    }
};
// Default/placeholder images
export const PLACEHOLDERS = {
    AVATAR: `${BASE_PATHS.PLACEHOLDERS}/placeholder-avatar.jpg`,
    THUMBNAIL: `${BASE_PATHS.PLACEHOLDERS}/placeholder-thumbnail.jpg`,
    PRODUCT: `${BASE_PATHS.PLACEHOLDERS}/default-product.webp`,
    VEHICLE: `${BASE_PATHS.PLACEHOLDERS}/default-vehicle.webp`,
    PROPERTY: `${BASE_PATHS.PLACEHOLDERS}/default-real-estate.webp`,
    JOB: `${BASE_PATHS.PLACEHOLDERS}/default-job.webp`,
    DEAL: `${BASE_PATHS.PLACEHOLDERS}/default-deal.webp`,
    SEARCH: `${BASE_PATHS.PLACEHOLDERS}/default-search.webp`,
    GENERIC: `${BASE_PATHS.PLACEHOLDERS}/placeholder-500x300.png`
};
// Icon paths
export const ICONS = {
    VIDEO: `${BASE_PATHS.ICONS}/video-icon.webp`,
    PHOTO: `${BASE_PATHS.ICONS}/photo.svg`,
    USER: `${BASE_PATHS.ICONS}/user.svg`,
    SEARCH: `${BASE_PATHS.ICONS}/search-icon.svg`,
    MESSAGE: `${BASE_PATHS.ICONS}/message-icon.svg`,
    HEART_FILLED: `${BASE_PATHS.ICONS}/heart-filled.svg`,
    HEART_OUTLINE: `${BASE_PATHS.ICONS}/heart-outline.svg`
};
// Background images
export const BACKGROUNDS = {
    HERO: `${BASE_PATHS.IMAGES}/psyche.webp`,
    DEFAULT: `${BASE_PATHS.IMAGES}/back.webp`,
    LOGIN: `${BASE_PATHS.IMAGES}/login-hero.svg`
};
// Consolidated export for easy access
export const IMAGES = {
    LOGOS,
    PLACEHOLDERS,
    ICONS,
    BACKGROUNDS,
    BASE_PATHS
};
// Utility function to get image with fallback
export const getImageWithFallback = (primaryPath, fallbackPath) => {
    return {
        src: primaryPath,
        fallback: fallbackPath,
        onError: (e) => {
            if (e.currentTarget.src !== fallbackPath) {
                e.currentTarget.src = fallbackPath;
            }
        }
    };
};
// Utility function to check if image exists (client-side)
export const preloadImage = (src) => {
    return new Promise((resolve, reject) => {
        const img = new Image();
        img.onload = () => resolve(src);
        img.onerror = () => reject(new Error(`Failed to load image: ${src}`));
        img.src = src;
    });
};
// Logo configuration for components
export const LOGO_CONFIG = {
    SFX: {
        src: LOGOS.SFX,
        alt: "SFX Logo - Go to home page",
        fallback: LOGOS.SFX, // Use SFX as fallback too
        sizes: {
            default: { width: 120, height: 50 }, // Adjusted for actual SFX dimensions
            aiMode: { width: 200, height: 85 },
            small: { width: 60, height: 25 },
            large: { width: 140, height: 60 }
        }
    },
    MOBILE: {
        src: LOGOS.SFX,
        alt: "SFX Logo - Go to home page", 
        fallback: LOGOS.SFX, // Use SFX as fallback too
        sizes: {
            default: { width: 80, height: 35 }, // Adjusted for actual SFX dimensions
            aiMode: { width: 100, height: 42 },
            small: { width: 55, height: 23 }
        }
    }
};
export default IMAGES; 