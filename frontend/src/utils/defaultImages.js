/**
 * Default Images Configuration
 * Provides fallback images for the application to prevent 404 errors
 */
export const DEFAULT_IMAGES = {
  // Vehicle/Car related
  vehicle: '/images/default-vehicle.jpg',
  car: '/images/default-car.webp',
  // User/Profile related
  avatar: '/images/placeholder-avatar.jpg',
  user: '/images/user-user.webp',
  // Product/Item related
  product: '/images/default-product.webp',
  item: '/images/default-deal.png',
  thumbnail: '/images/placeholder-thumbnail.jpg',
  // Real Estate
  realEstate: '/images/default-real-estate.webp',
  property: '/images/default-real-estate.webp',
  // Jobs
  job: '/images/default-job.webp',
  // General placeholders
  placeholder: '/images/placeholder-500x300.png',
  search: '/images/default-search.webp',
  // Logos
  logo: '/images/logo-black.png',
  logoMobile: '/images/logo-mobile.png',
  logoSmall: '/images/logo-small.png',
  // Icons
  video: '/images/video-icon.webp',
  photo: '/images/photo.svg',
  // Backgrounds
  background: '/images/back.jpg',
  hero: '/images/psyche.jpg',
};
/**
 * Get default image for a specific type
 * @param {string} type - The type of image needed
 * @param {string} fallback - Optional fallback if type not found
 * @returns {string} - Image path
 */
export function getDefaultImage(type, fallback = DEFAULT_IMAGES.placeholder) {
  return DEFAULT_IMAGES[type] || fallback;
}
/**
 * Get image with fallback handling
 * @param {string} src - Original image source
 * @param {string} type - Type for fallback
 * @returns {string} - Image path with fallback
 */
export function getImageWithFallback(src, type = 'placeholder') {
  if (!src || src === '' || src === null || src === undefined) {
    return getDefaultImage(type);
  }
  return src;
}
/**
 * Image categories for easy reference
 */
export const IMAGE_CATEGORIES = {
  VEHICLES: ['vehicle', 'car'],
  USERS: ['avatar', 'user'],
  PRODUCTS: ['product', 'item', 'thumbnail'],
  REAL_ESTATE: ['realEstate', 'property'],
  JOBS: ['job'],
  GENERAL: ['placeholder', 'search'],
  LOGOS: ['logo', 'logoMobile', 'logoSmall'],
  ICONS: ['video', 'photo'],
  BACKGROUNDS: ['background', 'hero'],
};
/**
 * Preload critical images for better performance
 */
export function preloadCriticalImages() {
  const criticalImages = [
    DEFAULT_IMAGES.logo,
    DEFAULT_IMAGES.placeholder,
    DEFAULT_IMAGES.avatar,
    DEFAULT_IMAGES.vehicle,
    DEFAULT_IMAGES.product,
  ];
  if (typeof window !== 'undefined') {
    criticalImages.forEach(src => {
      const link = document.createElement('link');
      link.rel = 'preload';
      link.as = 'image';
      link.href = src;
      document.head.appendChild(link);
    });
  }
} 