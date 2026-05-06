/**
 * Centralized mapping between content types and entity types
 * Use this utility across the application to ensure consistency
 */
// Map from content type (plural/UI) to entity type (singular/API)
export const CONTENT_TYPE_TO_ENTITY_MAP = {
  'products': 'product',
  'posts': 'post',
  'videos': 'video',
  'tweets': 'tweet',
  'jobs': 'job',
  'services': 'service',
  'deals': 'deal',
  'shorts': 'short',
  'vehicles': 'vehicle',
  'properties': 'property',
  'news': 'post' // News uses post search
};
/**
 * Convert a content type to its corresponding entity type
 * @param {string} contentType - The content type (e.g., 'products', 'posts')
 * @returns {string} - The corresponding entity type (e.g., 'product', 'post')
 */
export function getEntityTypeFromContentType(contentType) {
  if (!contentType) return null;
  // Check if we have a direct mapping
  const entityType = CONTENT_TYPE_TO_ENTITY_MAP[contentType];
  if (entityType) return entityType;
  // Fallback to removing trailing 's' for unknown types
  return contentType.endsWith('s') 
    ? contentType.slice(0, -1) 
    : contentType;
}
/**
 * Get the API function name for a given entity type
 * @param {string} entityType - The entity type
 * @returns {string} - The corresponding API function name
 */
export function getApiMethodName(entityType) {
  if (!entityType) return null;
  return `search${entityType.charAt(0).toUpperCase() + entityType.slice(1)}sWithFilters`;
}
/**
 * Debug helper to log entity type conversion
 * @param {string} contentType - The content type
 * @returns {string} - The entity type
 */
export function logEntityConversion(contentType) {
  const entityType = getEntityTypeFromContentType(contentType);
  return entityType;
} 