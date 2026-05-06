/**
 * HTML Sanitization Utility
 * Uses DOMPurify to prevent XSS attacks
 * Works on both client and server side
 */

import DOMPurify from 'isomorphic-dompurify';

/**
 * Default configuration for DOMPurify
 * Allows only safe HTML tags and attributes
 */
const DEFAULT_CONFIG = {
  ALLOWED_TAGS: [
    'p', 'br', 'span', 'div',
    'strong', 'b', 'em', 'i', 'u',
    'a', 
    'ul', 'ol', 'li',
    'h1', 'h2', 'h3', 'h4', 'h5', 'h6',
    'blockquote', 'code', 'pre',
    'table', 'thead', 'tbody', 'tr', 'th', 'td',
    'img', 'figure', 'figcaption'
  ],
  ALLOWED_ATTR: [
    'href', 'target', 'rel', 'title',
    'class', 'id', 'style',
    'src', 'alt', 'width', 'height',
    'colspan', 'rowspan'
  ],
  ALLOWED_URI_REGEXP: /^(?:(?:(?:f|ht)tps?|mailto|tel|callto|cid|xmpp):|[^a-z]|[a-z+.\-]+(?:[^a-z+.\-:]|$))/i,
  KEEP_CONTENT: true,
  FORCE_BODY: true
};

/**
 * Strict configuration for user-generated content
 * More restrictive than default
 */
const STRICT_CONFIG = {
  ALLOWED_TAGS: ['p', 'br', 'strong', 'em', 'a', 'ul', 'ol', 'li'],
  ALLOWED_ATTR: ['href', 'target', 'rel'],
  ALLOW_DATA_ATTR: false,
  KEEP_CONTENT: true,
  FORCE_BODY: true
};

/**
 * Configuration for rich content (like blog posts)
 */
const RICH_CONFIG = {
  ...DEFAULT_CONFIG,
  ALLOWED_TAGS: [
    ...DEFAULT_CONFIG.ALLOWED_TAGS,
    'video', 'source', 'iframe',
    'mark', 'del', 'ins', 'sub', 'sup',
    'hr', 'abbr', 'time'
  ],
  ALLOWED_ATTR: [
    ...DEFAULT_CONFIG.ALLOWED_ATTR,
    'controls', 'autoplay', 'loop', 'muted',
    'datetime', 'cite'
  ],
  ADD_TAGS: ['iframe'],
  ADD_ATTR: ['allowfullscreen', 'frameborder'],
  ALLOW_DATA_ATTR: true
};

/**
 * Sanitize HTML content to prevent XSS attacks
 * 
 * @param {string} html - The HTML string to sanitize
 * @param {Object} config - Optional DOMPurify configuration
 * @returns {string} - Sanitized HTML string
 */
export function sanitizeHtml(html, config = DEFAULT_CONFIG) {
  if (!html) return '';
  
  try {
    return DOMPurify.sanitize(html, config);
  } catch (error) {
    // Log error but return empty string to prevent breaking the app
    if (process.env.NODE_ENV === 'development') {
      console.error('HTML sanitization error:', error);
    }
    return '';
  }
}

/**
 * Sanitize HTML with strict rules for user-generated content
 * 
 * @param {string} html - The HTML string to sanitize
 * @returns {string} - Sanitized HTML string
 */
export function sanitizeUserHtml(html) {
  return sanitizeHtml(html, STRICT_CONFIG);
}

/**
 * Sanitize HTML for rich content like blog posts
 * 
 * @param {string} html - The HTML string to sanitize
 * @returns {string} - Sanitized HTML string
 */
export function sanitizeRichHtml(html) {
  return sanitizeHtml(html, RICH_CONFIG);
}

/**
 * Sanitize and strip all HTML tags, returning plain text
 * 
 * @param {string} html - The HTML string to sanitize
 * @returns {string} - Plain text without HTML tags
 */
export function sanitizeToText(html) {
  if (!html) return '';
  
  const config = {
    ALLOWED_TAGS: [],
    ALLOWED_ATTR: [],
    KEEP_CONTENT: true
  };
  
  return DOMPurify.sanitize(html, config);
}

/**
 * Check if HTML contains potentially dangerous content
 * 
 * @param {string} html - The HTML string to check
 * @returns {boolean} - True if content was sanitized (dangerous content found)
 */
export function containsDangerousHtml(html) {
  if (!html) return false;
  
  const clean = DOMPurify.sanitize(html, DEFAULT_CONFIG);
  return clean !== html;
}

/**
 * Sanitize URL to prevent javascript: and data: protocols
 * 
 * @param {string} url - The URL to sanitize
 * @returns {string} - Sanitized URL or empty string if dangerous
 */
export function sanitizeUrl(url) {
  if (!url) return '';
  
  // Check against allowed protocols
  const allowedProtocols = ['http:', 'https:', 'mailto:', 'tel:'];
  try {
    const urlObj = new URL(url);
    if (allowedProtocols.includes(urlObj.protocol)) {
      return url;
    }
  } catch {
    // If URL parsing fails, check if it's a relative URL
    if (url.startsWith('/') || url.startsWith('#')) {
      return url;
    }
  }
  
  return '';
}

/**
 * Create a sanitized HTML element with React
 * Use this instead of dangerouslySetInnerHTML
 * 
 * @param {string} html - The HTML string to render
 * @param {Object} config - Optional DOMPurify configuration
 * @returns {Object} - Props object for React element
 */
export function createSafeHtml(html, config = DEFAULT_CONFIG) {
  return {
    dangerouslySetInnerHTML: {
      __html: sanitizeHtml(html, config)
    }
  };
}

// Export configurations for custom use
export const SANITIZE_CONFIG = {
  DEFAULT: DEFAULT_CONFIG,
  STRICT: STRICT_CONFIG,
  RICH: RICH_CONFIG
};

// Default export
export default sanitizeHtml;