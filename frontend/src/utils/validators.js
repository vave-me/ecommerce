/**
 * Shared validation utilities
 */

/**
 * Check if value is empty (null, undefined, empty string, or empty array)
 * @param {*} value - Value to check
 * @returns {boolean} True if empty
 */
export function isEmpty(value) {
    return (
        value == null ||
        value === '' ||
        (Array.isArray(value) && value.length === 0) ||
        (typeof value === 'object' && Object.keys(value).length === 0)
    );
}

/**
 * Validate email format
 * @param {string} email - Email to validate
 * @returns {boolean} True if valid email
 */
export function isValidEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
}

/**
 * Validate phone number
 * @param {string} phone - Phone number to validate
 * @returns {boolean} True if valid phone
 */
export function isValidPhone(phone) {
    // Remove spaces and dashes
    const cleaned = phone.replace(/[\s-]/g, '');
    // Basic validation for international format
    const phoneRegex = /^[\+]?[(]?[0-9]{1,4}[)]?[-\s\.]?[(]?[0-9]{1,4}[)]?[-\s\.]?[0-9]{1,9}$/;
    return phoneRegex.test(cleaned);
}

/**
 * Validate URL format
 * @param {string} url - URL to validate
 * @returns {boolean} True if valid URL
 */
export function isValidUrl(url) {
    try {
        new URL(url);
        return true;
    } catch {
        return false;
    }
}

/**
 * Validate price (positive number with up to 2 decimals)
 * @param {number|string} price - Price to validate
 * @returns {boolean} True if valid price
 */
export function isValidPrice(price) {
    const num = parseFloat(price);
    return !isNaN(num) && num >= 0 && /^\d+(\.\d{1,2})?$/.test(price.toString());
}

/**
 * Validate required fields in an object
 * @param {Object} obj - Object to validate
 * @param {Array<string>} requiredFields - List of required field names
 * @returns {Object} { isValid: boolean, missingFields: Array<string> }
 */
export function validateRequiredFields(obj, requiredFields) {
    const missingFields = requiredFields.filter(field => isEmpty(obj[field]));
    return {
        isValid: missingFields.length === 0,
        missingFields
    };
}

/**
 * Validate string length
 * @param {string} str - String to validate
 * @param {number} min - Minimum length
 * @param {number} max - Maximum length
 * @returns {boolean} True if within range
 */
export function isValidLength(str, min = 0, max = Infinity) {
    if (!str) return min === 0;
    return str.length >= min && str.length <= max;
}

/**
 * Sanitize input to prevent XSS
 * @param {string} input - Input to sanitize
 * @returns {string} Sanitized input
 */
export function sanitizeInput(input) {
    if (!input) return '';
    return input
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;');
}

/**
 * Validate file type
 * @param {File} file - File to validate
 * @param {Array<string>} allowedTypes - Allowed MIME types
 * @returns {boolean} True if valid file type
 */
export function isValidFileType(file, allowedTypes) {
    if (!file || !file.type) return false;
    return allowedTypes.some(type => {
        if (type.endsWith('/*')) {
            // Handle wildcard types like 'image/*'
            const baseType = type.slice(0, -2);
            return file.type.startsWith(baseType);
        }
        return file.type === type;
    });
}

/**
 * Validate file size
 * @param {File} file - File to validate
 * @param {number} maxSizeInMB - Maximum size in megabytes
 * @returns {boolean} True if within size limit
 */
export function isValidFileSize(file, maxSizeInMB) {
    if (!file) return false;
    const maxSizeInBytes = maxSizeInMB * 1024 * 1024;
    return file.size <= maxSizeInBytes;
}