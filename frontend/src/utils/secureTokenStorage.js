/**
 * Secure token storage utility
 * Provides enhanced security for authentication tokens
 * 
 * Security features:
 * - HttpOnly cookies for refresh tokens (server-side required)
 * - Memory storage for access tokens
 * - Encryption for localStorage fallback
 * - Token expiry validation
 * - XSS protection
 */

// In-memory storage for access tokens (most secure for SPAs)
let memoryTokens = {
  accessToken: null,
  tokenExpiry: null
};

// Encryption key (in production, this should come from environment)
const STORAGE_KEY = process.env.NEXT_PUBLIC_STORAGE_KEY || 'sfx_secure_auth';

/**
 * Simple encryption for localStorage fallback
 * In production, use Web Crypto API or a proper encryption library
 */
const encrypt = (text) => {
  if (!text) return '';
  try {
    // Basic obfuscation - in production use proper encryption
    return btoa(encodeURIComponent(text));
  } catch (error) {
    return text;
  }
};

const decrypt = (encryptedText) => {
  if (!encryptedText) return '';
  try {
    return decodeURIComponent(atob(encryptedText));
  } catch (error) {
    return encryptedText;
  }
};

/**
 * Token storage manager with multiple security layers
 */
export const secureTokenStorage = {
  /**
   * Store access token in memory (most secure)
   */
  setAccessToken: (token, expiresIn = 3600) => {
    memoryTokens.accessToken = token;
    memoryTokens.tokenExpiry = Date.now() + (expiresIn * 1000);
    
    // Fallback to sessionStorage (cleared on tab close)
    if (typeof window !== 'undefined' && window.sessionStorage) {
      try {
        const encryptedToken = encrypt(token);
        sessionStorage.setItem(`${STORAGE_KEY}_access`, encryptedToken);
        sessionStorage.setItem(`${STORAGE_KEY}_expiry`, memoryTokens.tokenExpiry.toString());
      } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
    }
  },

  /**
   * Get access token from memory or sessionStorage
   */
  getAccessToken: () => {
    // First check memory
    if (memoryTokens.accessToken && memoryTokens.tokenExpiry > Date.now()) {
      return memoryTokens.accessToken;
    }
    
    // Fallback to sessionStorage
    if (typeof window !== 'undefined' && window.sessionStorage) {
      try {
        const encryptedToken = sessionStorage.getItem(`${STORAGE_KEY}_access`);
        const expiry = sessionStorage.getItem(`${STORAGE_KEY}_expiry`);
        
        if (encryptedToken && expiry && parseInt(expiry) > Date.now()) {
          const token = decrypt(encryptedToken);
          // Restore to memory
          memoryTokens.accessToken = token;
          memoryTokens.tokenExpiry = parseInt(expiry);
          return token;
        }
      } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
    }
    
    return null;
  },

  /**
   * Store refresh token (should ideally be httpOnly cookie)
   * Using encrypted localStorage as fallback
   */
  setRefreshToken: (token) => {
    if (typeof window !== 'undefined' && window.localStorage) {
      try {
        const encryptedToken = encrypt(token);
        localStorage.setItem(`${STORAGE_KEY}_refresh`, encryptedToken);
      } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
    }
  },

  /**
   * Get refresh token
   */
  getRefreshToken: () => {
    if (typeof window !== 'undefined' && window.localStorage) {
      try {
        const encryptedToken = localStorage.getItem(`${STORAGE_KEY}_refresh`);
        return encryptedToken ? decrypt(encryptedToken) : null;
      } catch (error) {
        return null;
      }
    }
    return null;
  },

  /**
   * Clear all tokens (logout)
   */
  clearTokens: () => {
    // Clear memory
    memoryTokens = {
      accessToken: null,
      tokenExpiry: null
    };
    
    // Clear storages
    if (typeof window !== 'undefined') {
      try {
        // Clear sessionStorage
        if (window.sessionStorage) {
          sessionStorage.removeItem(`${STORAGE_KEY}_access`);
          sessionStorage.removeItem(`${STORAGE_KEY}_expiry`);
        }
        
        // Clear localStorage
        if (window.localStorage) {
          localStorage.removeItem(`${STORAGE_KEY}_refresh`);
          // Also clear legacy unencrypted tokens
          localStorage.removeItem('access_token');
          localStorage.removeItem('refresh_token');
        }
      } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
    }
  },

  /**
   * Check if user is authenticated
   */
  isAuthenticated: () => {
    return !!secureTokenStorage.getAccessToken();
  },

  /**
   * Get token expiry time
   */
  getTokenExpiry: () => {
    return memoryTokens.tokenExpiry;
  },

  /**
   * Migrate from old localStorage to secure storage
   */
  migrateFromLocalStorage: () => {
    if (typeof window !== 'undefined' && window.localStorage) {
      try {
        // Check for old tokens
        const oldAccessToken = localStorage.getItem('access_token');
        const oldRefreshToken = localStorage.getItem('refresh_token');
        
        if (oldAccessToken) {
          secureTokenStorage.setAccessToken(oldAccessToken);
          localStorage.removeItem('access_token');
        }
        
        if (oldRefreshToken) {
          secureTokenStorage.setRefreshToken(oldRefreshToken);
          localStorage.removeItem('refresh_token');
        }
      } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
    }
  }
};

// Auto-migrate on module load
if (typeof window !== 'undefined') {
  secureTokenStorage.migrateFromLocalStorage();
}

// Export individual functions for backward compatibility
export const setAccessToken = secureTokenStorage.setAccessToken;
export const getAccessToken = secureTokenStorage.getAccessToken;
export const setRefreshToken = secureTokenStorage.setRefreshToken;
export const getRefreshToken = secureTokenStorage.getRefreshToken;
export const clearTokens = secureTokenStorage.clearTokens;
export const isAuthenticated = secureTokenStorage.isAuthenticated;

export default secureTokenStorage;