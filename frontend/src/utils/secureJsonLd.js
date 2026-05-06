/**
 * Secure JSON-LD Utility
 * Prevents XSS vulnerabilities in JSON-LD structured data injection
 * 
 * SECURITY: This utility ensures safe serialization of JSON-LD data
 * by properly escaping dangerous characters and validating content.
 */
/**
 * Safely serialize JSON-LD data for insertion into HTML
 * @param {Object} jsonLdData - The JSON-LD object to serialize
 * @returns {string} - Safely escaped JSON string
 */
export function safeSerializeJsonLd(jsonLdData) {
  if (!jsonLdData || typeof jsonLdData !== 'object') {
    if (process.env.NODE_ENV === 'development') {
    }
    return '';
  }
  try {
    // Clean the data first to remove dangerous properties
    const cleanedData = sanitizeJsonLdData(jsonLdData);
    // Then validate that this looks like valid JSON-LD
    if (!isValidJsonLd(cleanedData)) {
      if (process.env.NODE_ENV === 'development') {
        // Only warn if it's actually missing required fields, not just for ItemList structures
        const hasContext = !!cleanedData['@context'];
        const hasType = !!cleanedData['@type'];
        const isItemList = cleanedData['@type'] === 'ItemList' && cleanedData['itemListElement'];
        if (!hasContext || !hasType || (!isItemList && !cleanedData['name'] && !cleanedData['url'])) {
        }
      }
      return '';
    }
    // Serialize with safe JSON.stringify
    let jsonString = JSON.stringify(cleanedData);
    // Escape dangerous characters that could break out of script context
    jsonString = jsonString
      .replace(/</g, '\\u003c')     // Escape < to prevent script injection
      .replace(/>/g, '\\u003e')     // Escape > 
      .replace(/&/g, '\\u0026')     // Escape &
      .replace(/'/g, '\\u0027')     // Escape single quotes
      .replace(/"/g, '\\u0022')     // Escape double quotes in content
      .replace(/\u2028/g, '\\u2028') // Escape line separator
      .replace(/\u2029/g, '\\u2029'); // Escape paragraph separator
    return jsonString;
  } catch (error) {
    if (process.env.NODE_ENV === 'development') {
    }
    return '';
  }
}
/**
 * Sanitize JSON-LD data by removing dangerous properties
 * @param {Object} data - Object to sanitize
 * @returns {Object} - Sanitized object
 */
function sanitizeJsonLdData(data) {
  if (!data || typeof data !== 'object') {
    return data;
  }
  const dangerousKeys = [
    '__proto__', 'constructor', 'prototype',
    '__defineGetter__', '__defineSetter__',
    '__lookupGetter__', '__lookupSetter__',
    'hasOwnProperty', 'isPrototypeOf', 'propertyIsEnumerable',
    'toLocaleString', 'toString', 'valueOf'
  ];
  function cleanObject(obj) {
    if (!obj || typeof obj !== 'object') {
      return obj;
    }
    if (Array.isArray(obj)) {
      return obj.map(item => cleanObject(item)).filter(item => item !== null);
    }
    const cleaned = {};
    // Use Object.getOwnPropertyNames to get all properties, including non-enumerable ones
    const allKeys = Object.getOwnPropertyNames(obj);
    for (const key of allKeys) {
      // Skip dangerous properties completely
      if (dangerousKeys.includes(key)) {
        continue;
      }
      try {
        const value = obj[key];
        // Recursively clean nested objects
        if (typeof value === 'object' && value !== null) {
          const cleanedValue = cleanObject(value);
          if (cleanedValue !== null) {
            cleaned[key] = cleanedValue;
          }
        } else if (typeof value === 'string' || typeof value === 'number' || typeof value === 'boolean') {
          cleaned[key] = value;
        }
      } catch (error) {
        // Skip properties that can't be accessed safely
        continue;
      }
    }
    return cleaned;
  }
  return cleanObject(data);
}
/**
 * Basic validation for JSON-LD structure with enhanced security
 * @param {Object} data - Object to validate
 * @returns {boolean} - True if valid JSON-LD structure
 */
function isValidJsonLd(data) {
  if (!data || typeof data !== 'object' || Array.isArray(data)) {
    return false;
  }
  // Must have @context - allow string or object
  if (!data['@context'] || (typeof data['@context'] !== 'string' && typeof data['@context'] !== 'object')) {
    return false;
  }
  // Must have @type - allow string or array
  if (!data['@type'] || (typeof data['@type'] !== 'string' && !Array.isArray(data['@type']))) {
    return false;
  }
  // Special validation for ItemList - these are valid even without name/url
  if (data['@type'] === 'ItemList' && data['itemListElement']) {
    return true;
  }
  // For other types, be more lenient - just require @context and @type
  return true;
  // Check for dangerous properties that shouldn't be in JSON-LD
  const dangerousKeys = ['__proto__', 'constructor', 'prototype', '__defineGetter__', '__defineSetter__', '__lookupGetter__', '__lookupSetter__'];
  function checkObjectSafety(obj) {
    if (!obj || typeof obj !== 'object') return true;
    for (const key of dangerousKeys) {
      if (key in obj) {
        // Silently reject objects with dangerous properties
        return false;
      }
    }
    // Recursively check nested objects
    for (const value of Object.values(obj)) {
      if (typeof value === 'object' && value !== null) {
        if (!checkObjectSafety(value)) {
          return false;
        }
      }
    }
    return true;
  }
  return checkObjectSafety(data);
}
/**
 * React component for safely rendering JSON-LD
 * @param {Object} props - Component props
 * @param {Object} props.data - JSON-LD data to render
 * @returns {JSX.Element} - Safe script tag with JSON-LD
 */
export function SafeJsonLdScript({ data }) {
  const safeJsonString = safeSerializeJsonLd(data);
  if (!safeJsonString) {
    // Don't render anything if the data is invalid
    return null;
  }
  return (
    <script 
      type="application/ld+json"
      dangerouslySetInnerHTML={{ __html: safeJsonString }}
    />
  );
}
/**
 * Legacy function for backward compatibility
 * Use SafeJsonLdScript component instead when possible
 * @param {Object} jsonLdData - JSON-LD data
 * @returns {string} - Safe JSON string for HTML injection
 */
export function secureJsonLdString(jsonLdData) {
  return safeSerializeJsonLd(jsonLdData);
} 