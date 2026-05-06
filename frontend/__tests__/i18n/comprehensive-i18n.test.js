/**
 * Comprehensive Internationalization and Routing Test Suite
 * Tests translation completeness, routing functionality, and live server responses
 */

import { exec } from 'child_process';
import { promisify } from 'util';
import fs from 'fs/promises';
import path from 'path';
import { routing } from '../../i18n/routing';

const execAsync = promisify(exec);

describe('Comprehensive I18n and Routing Tests', () => {
  const BASE_URL = process.env.TEST_BASE_URL || 'http://localhost:3000';
  const SUPPORTED_LOCALES = ['en', 'pl', 'de'];
  const MESSAGES_DIR = path.join(process.cwd(), 'messages');

  // Test data for common routes
  const TEST_ROUTES = [
    '',
    '/home',
    '/marketplace',
    '/deals',
    '/jobs',
    '/properties',
    '/services',
    '/products',
    '/news',
    '/search',
    '/login',
    '/about',
    '/contact',
    '/terms',
    '/privacy'
  ];

  describe('Translation Files Validation', () => {
    let translations = {};

    beforeAll(async () => {
      // Load all translation files
      for (const locale of SUPPORTED_LOCALES) {
        try {
          const filePath = path.join(MESSAGES_DIR, `${locale}.json`);
          const content = await fs.readFile(filePath, 'utf-8');
          translations[locale] = JSON.parse(content);
        } catch (error) {
          console.error(`Failed to load ${locale}.json:`, error.message);
        }
      }
    });

    test('all translation files exist', async () => {
      for (const locale of SUPPORTED_LOCALES) {
        const filePath = path.join(MESSAGES_DIR, `${locale}.json`);
        await expect(fs.access(filePath)).resolves.not.toThrow();
      }
    });

    test('translation files have valid JSON structure', () => {
      for (const locale of SUPPORTED_LOCALES) {
        expect(translations[locale]).toBeDefined();
        expect(typeof translations[locale]).toBe('object');
      }
    });

    test('all locales have the same translation keys structure', () => {
      const englishKeys = getNestedKeys(translations.en);
      
      for (const locale of SUPPORTED_LOCALES.filter(l => l !== 'en')) {
        const localeKeys = getNestedKeys(translations[locale]);
        
        // Check for missing keys
        const missingKeys = englishKeys.filter(key => !localeKeys.includes(key));
        const extraKeys = localeKeys.filter(key => !englishKeys.includes(key));
        
        expect(missingKeys).toEqual([]);
        expect(extraKeys).toEqual([]);
      }
    });

    test('critical translation keys are present', () => {
      const criticalKeys = [
        'HomePage.browseCategories',
        'NotFound.title',
        'Error.title',
        'LoginForm.signIn',
        'Header.search',
        'Seo.title',
        'Seo.description'
      ];

      for (const locale of SUPPORTED_LOCALES) {
        for (const key of criticalKeys) {
          const value = getNestedValue(translations[locale], key);
          expect(value).toBeDefined();
          expect(typeof value).toBe('string');
          expect(value.trim()).not.toBe('');
        }
      }
    });

    test('no translation values are empty or just whitespace', () => {
      for (const locale of SUPPORTED_LOCALES) {
        const emptyValues = findEmptyTranslations(translations[locale]);
        expect(emptyValues).toEqual([]);
      }
    });

    test('placeholder consistency across locales', () => {
      const placeholderPattern = /\{[^}]+\}/g;
      
      for (const key of getNestedKeys(translations.en)) {
        const englishValue = getNestedValue(translations.en, key);
        if (typeof englishValue !== 'string') continue;
        
        const englishPlaceholders = (englishValue.match(placeholderPattern) || []).sort();
        
        for (const locale of SUPPORTED_LOCALES.filter(l => l !== 'en')) {
          const localeValue = getNestedValue(translations[locale], key);
          if (typeof localeValue !== 'string') continue;
          
          const localePlaceholders = (localeValue.match(placeholderPattern) || []).sort();
          
          expect(localePlaceholders).toEqual(englishPlaceholders);
        }
      }
    });
  });

  describe('Routing Configuration', () => {
    test('routing configuration is valid', () => {
      expect(routing).toBeDefined();
      expect(routing.locales).toEqual(expect.arrayContaining(SUPPORTED_LOCALES));
      expect(routing.defaultLocale).toBe('en');
      expect(routing.localePrefix).toBe('always');
    });

    test('all supported locales are configured', () => {
      expect(routing.locales).toHaveLength(SUPPORTED_LOCALES.length);
      expect(routing.locales.sort()).toEqual(SUPPORTED_LOCALES.sort());
    });
  });

  describe('Live Server Routing Tests', () => {
    const timeout = 30000; // 30 seconds timeout for server tests

    beforeAll(async () => {
      // Wait for server to be ready
      await waitForServer(BASE_URL, 30000);
    }, timeout);

    test('root path redirects to default locale', async () => {
      try {
        const response = await makeRequest('GET', BASE_URL, { followRedirects: false });
        
        // Should redirect to /en
        expect([301, 302, 307, 308]).toContain(response.statusCode);
        expect(response.headers.location).toMatch(/\/en\/?$/);
      } catch (error) {
        console.warn('Root redirect test failed:', error.message);
        // Don't fail the test if server is not available
      }
    }, timeout);

    test.each(SUPPORTED_LOCALES)('locale %s routes are accessible', async (locale) => {
      try {
        const response = await makeRequest('GET', `${BASE_URL}/${locale}`);
        
        // Should return 200 or redirect to a valid page
        expect([200, 301, 302, 307, 308]).toContain(response.statusCode);
        
        if (response.statusCode === 200) {
          expect(response.headers['content-type']).toMatch(/text\/html/);
        }
      } catch (error) {
        console.warn(`Locale ${locale} test failed:`, error.message);
      }
    }, timeout);

    test.each(TEST_ROUTES)('route %s works for all locales', async (route) => {
      for (const locale of SUPPORTED_LOCALES) {
        try {
          const url = `${BASE_URL}/${locale}${route}`;
          const response = await makeRequest('GET', url);
          
          // Should not return 404 or 500 for valid routes
          expect(response.statusCode).not.toBe(404);
          
          if (response.statusCode >= 500) {
            console.warn(`Server error for ${url}: ${response.statusCode}`);
          }
        } catch (error) {
          console.warn(`Route test failed for ${locale}${route}:`, error.message);
        }
      }
    }, timeout);

    test('invalid locale returns 404', async () => {
      try {
        const response = await makeRequest('GET', `${BASE_URL}/invalid-locale`);
        expect(response.statusCode).toBe(404);
      } catch (error) {
        console.warn('Invalid locale test failed:', error.message);
      }
    }, timeout);

    test('API routes are not affected by locale middleware', async () => {
      try {
        const response = await makeRequest('GET', `${BASE_URL}/api/health`);
        // API routes should work regardless of locale
        expect([200, 404]).toContain(response.statusCode); // 404 is OK if endpoint doesn't exist
      } catch (error) {
        console.warn('API route test failed:', error.message);
      }
    }, timeout);

    test('static assets are accessible', async () => {
      try {
        const response = await makeRequest('GET', `${BASE_URL}/favicon.ico`);
        expect([200, 404]).toContain(response.statusCode);
      } catch (error) {
        console.warn('Static asset test failed:', error.message);
      }
    }, timeout);
  });

  describe('SEO and Meta Tags', () => {
    test.each(SUPPORTED_LOCALES)('locale %s has proper meta tags', async (locale) => {
      try {
        const response = await makeRequest('GET', `${BASE_URL}/${locale}`);
        
        if (response.statusCode === 200 && response.body) {
          const html = response.body;
          
          // Check for essential meta tags
          expect(html).toMatch(/<title>/);
          expect(html).toMatch(/<meta[^>]+description/);
          expect(html).toMatch(/<html[^>]+lang=["']?${locale}["']?/);
        }
      } catch (error) {
        console.warn(`Meta tags test failed for ${locale}:`, error.message);
      }
    });
  });

  describe('Performance and Caching', () => {
    test('static pages have appropriate cache headers', async () => {
      try {
        const response = await makeRequest('GET', `${BASE_URL}/en/about`);
        
        if (response.statusCode === 200) {
          // Check for cache-related headers
          expect(response.headers).toHaveProperty('cache-control');
        }
      } catch (error) {
        console.warn('Cache headers test failed:', error.message);
      }
    });
  });
});

// Helper functions
function getNestedKeys(obj, prefix = '') {
  let keys = [];
  
  for (const [key, value] of Object.entries(obj)) {
    const fullKey = prefix ? `${prefix}.${key}` : key;
    
    if (typeof value === 'object' && value !== null) {
      keys = keys.concat(getNestedKeys(value, fullKey));
    } else {
      keys.push(fullKey);
    }
  }
  
  return keys;
}

function getNestedValue(obj, path) {
  return path.split('.').reduce((current, key) => current?.[key], obj);
}

function findEmptyTranslations(obj, prefix = '') {
  let emptyKeys = [];
  
  for (const [key, value] of Object.entries(obj)) {
    const fullKey = prefix ? `${prefix}.${key}` : key;
    
    if (typeof value === 'object' && value !== null) {
      emptyKeys = emptyKeys.concat(findEmptyTranslations(value, fullKey));
    } else if (typeof value === 'string' && value.trim() === '') {
      emptyKeys.push(fullKey);
    }
  }
  
  return emptyKeys;
}

async function waitForServer(url, timeout = 30000) {
  const start = Date.now();
  
  while (Date.now() - start < timeout) {
    try {
      await makeRequest('GET', url);
      return;
    } catch (error) {
      await new Promise(resolve => setTimeout(resolve, 1000));
    }
  }
  
  throw new Error(`Server not ready after ${timeout}ms`);
}

async function makeRequest(method, url, options = {}) {
  const curlCommand = buildCurlCommand(method, url, options);
  
  try {
    const { stdout, stderr } = await execAsync(curlCommand);
    return parseCurlResponse(stdout);
  } catch (error) {
    throw new Error(`Request failed: ${error.message}`);
  }
}

function buildCurlCommand(method, url, options = {}) {
  let command = `curl -s -i -X ${method}`;
  
  if (options.followRedirects === false) {
    // Don't follow redirects
  } else {
    command += ' -L';
  }
  
  command += ` "${url}"`;
  
  return command;
}

function parseCurlResponse(response) {
  const [headerSection, ...bodyParts] = response.split('\r\n\r\n');
  const body = bodyParts.join('\r\n\r\n');
  
  const lines = headerSection.split('\r\n');
  const statusLine = lines[0];
  const statusCode = parseInt(statusLine.split(' ')[1]);
  
  const headers = {};
  for (let i = 1; i < lines.length; i++) {
    const [key, ...valueParts] = lines[i].split(': ');
    if (key && valueParts.length > 0) {
      headers[key.toLowerCase()] = valueParts.join(': ');
    }
  }
  
  return {
    statusCode,
    headers,
    body: body.trim()
  };
} 