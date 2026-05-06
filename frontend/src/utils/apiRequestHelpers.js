/**
 * SHARED API REQUEST UTILITIES
 * Eliminates duplicate code across API files
 * - Common request body cleaning
 * - Parameter parsing utilities  
 * - Request configuration helpers
 */
/**
 * Safely parse numeric values with fallback
 */
export const safeParseFloat = (value, fallback = 0) => {
  if (value === '' || value === null || value === undefined) return fallback;
  const parsed = parseFloat(value);
  return isNaN(parsed) ? fallback : parsed;
};
export const safeParseInt = (value, fallback = 0) => {
  if (value === '' || value === null || value === undefined) return fallback;
  const parsed = parseInt(value, 10);
  return isNaN(parsed) ? fallback : parsed;
};
/**
 * Clean request body by removing empty/null/undefined values
 * This eliminates the duplicate logic found in productsApi, vehiclesApi, propertiesApi
 */
export const cleanRequestBody = (requestBody) => {
  const cleaned = {};
  for (const [key, value] of Object.entries(requestBody)) {
    // Skip empty values
    if (value === '' || value === null || value === undefined) {
      continue;
    }
    // Skip empty arrays
    if (Array.isArray(value) && value.length === 0) {
      continue;
    }
    // Keep valid values
    cleaned[key] = value;
  }
  return cleaned;
};
/**
 * Create standardized search request body
 * Used by products, vehicles, properties, services, etc.
 */
export const createSearchRequestBody = (filters = {}) => {
  const {
    name = '',
    category = '',
    minPrice = '',
    maxPrice = '',
    condition = '',
    brand = '',
    model = '',
    status = '',
    tags = [],
    manageStock = false,
    minStock = 0,
    maxStock = 9999,
    SKU = '',
    negotiable = false,
    userType = '',
    middlemanService = '',
    hasVariants = false,
    shippingCost = '',
    // Dimensions & weight
    minWeight = '',
    maxWeight = '',
    minHeight = '',
    maxHeight = '',
    minWidth = '',
    maxWidth = '',
    minDepth = '',
    maxDepth = '',
    // Location
    lat = 0,
    lng = 0,
    radius = 0,
    // Pagination
    page = 1,
    pageSize = 20,
    offset = 0,
    limit = 20,
    sortBy = '',
    sortOrder = 'asc',
    // Entity-specific fields (optional)
    ...entitySpecificFields
  } = filters;
  const requestBody = {
    name,
    category,
    minPrice: safeParseFloat(minPrice),
    maxPrice: safeParseFloat(maxPrice),
    brand,
    condition,
    model,
    status,
    tags: Array.isArray(tags) ? tags : [],
    manageStock,
    minStock: safeParseInt(minStock),
    maxStock: safeParseInt(maxStock),
    sku: SKU,
    negotiable,
    userType,
    middlemanService,
    hasVariants,
    shippingCost,
    minWeight: safeParseFloat(minWeight),
    maxWeight: safeParseFloat(maxWeight),
    minHeight: safeParseFloat(minHeight),
    maxHeight: safeParseFloat(maxHeight),
    minWidth: safeParseFloat(minWidth),
    maxWidth: safeParseFloat(maxWidth),
    minDepth: safeParseFloat(minDepth),
    maxDepth: safeParseFloat(maxDepth),
    lat: safeParseFloat(lat),
    lng: safeParseFloat(lng),
    radius: safeParseInt(radius),
    page: safeParseInt(page, 1),
    pageSize: safeParseInt(pageSize, 20),
    offset: safeParseInt(offset),
    limit: safeParseInt(limit, 20),
    sortBy,
    sortOrder,
    ...entitySpecificFields
  };
  return cleanRequestBody(requestBody);
};
/**
 * Generate cache key for API requests
 * Eliminates duplicate cache key generation logic
 */
export const generateCacheKey = (endpoint, params = {}) => {
  // Sort keys for consistent cache keys
  const sortedParams = Object.keys(params)
    .sort()
    .reduce((acc, key) => {
      acc[key] = params[key];
      return acc;
    }, {});
  return `${endpoint}:${JSON.stringify(sortedParams)}`;
};
/**
 * File processing utilities for deduplication
 * Eliminates duplicate file handling logic
 */
export const generateFileHash = async (file) => {
  const arrayBuffer = await file.arrayBuffer();
  const hashBuffer = await crypto.subtle.digest('SHA-256', arrayBuffer);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
};
export const processFileWithDeduplication = async (file, existingHashes = new Set()) => {
  const preview = file.type.startsWith('image/') ? URL.createObjectURL(file) : null;
  const hash = await generateFileHash(file);
  return {
    file,
    name: file.name,
    type: file.type,
    size: file.size,
    preview,
    hash,
    isDuplicate: existingHashes.has(hash),
    uploadProgress: 0
  };
};
/**
 * Common validation patterns
 */
export const VALIDATION_PATTERNS = {
  email: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
  phone: /^[\+]?[\s\-\(\)]?[\d\s\-\(\)]{10,}$/,
  url: /^https?:\/\/.+\..+/,
  zipcode: /^\d{5}(-\d{4})?$/,
  creditcard: /^\d{4}\s?\d{4}\s?\d{4}\s?\d{4}$/
};
/**
 * Import unified utilities to eliminate duplications
 */
import { UnifiedUtils } from './duplicateEliminator';
/**
 * Debounce utility - from unified utilities
 */
export const createDebounceFunction = UnifiedUtils.debounce;
/**
 * Local storage utilities - from unified utilities
 */
export const DraftStorageUtils = UnifiedUtils.drafts;
/**
 * Common form utilities for entity creation
 */
export const EntityFormUtils = {
  buildBasicPayload: (basicInfo, user) => ({
    name: basicInfo.name || '',
    description: basicInfo.description || '',
    category_id: basicInfo.categoryId || '',
    category_slug: basicInfo.categorySlug || '',
    category_name: basicInfo.categoryName || '',
    price: safeParseFloat(basicInfo.price),
    negotiable: !!basicInfo.negotiable,
    user_type: basicInfo.userType || 'private',
    tags: Array.isArray(basicInfo.tags) 
      ? basicInfo.tags 
      : (typeof basicInfo.tags === 'string' && basicInfo.tags)
        ? basicInfo.tags.split(',').map(tag => tag.trim()).filter(Boolean)
        : [],
    user_id: user?.userId || user?.id
  }),
  extractLocation: (locationData) => ({
    lat: safeParseFloat(locationData?.lat),
    lng: safeParseFloat(locationData?.lng),
    address: locationData?.address || '',
    city: locationData?.city || '',
    country: locationData?.country || '',
    postalCode: locationData?.postalCode || ''
  }),
  extractMedia: (mediaData) => ({
    thumbnail: mediaData?.thumbnail || '',
    video_url: mediaData?.videoUrl || '',
    images: Array.isArray(mediaData?.images) ? mediaData.images : []
  })
};
export default {
  safeParseFloat,
  safeParseInt,
  cleanRequestBody,
  createSearchRequestBody,
  generateCacheKey,
  generateFileHash,
  processFileWithDeduplication,
  VALIDATION_PATTERNS,
  createDebounceFunction,
  DraftStorageUtils,
  EntityFormUtils
}; 